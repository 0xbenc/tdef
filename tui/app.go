package tui

import (
	"fmt"
	"strconv"
	"time"

	"tdef/game"
	"tdef/hiscore"
	"tdef/render"
)

const tickRate = 20.0
const frameRate = 30.0

// Screen is the top-level state of the TUI.
type Screen int

const (
	ScreenTitle Screen = iota
	ScreenMenu
	ScreenHelp
	ScreenHiscores
	ScreenLevelSelect
	ScreenGame
)

type App struct {
	term     *Terminal
	g        *game.State
	ui       render.UI
	pal      render.Colors
	layout   render.Layout
	diff     game.Difficulty
	level    string
	scored   bool
	quitting bool

	// Pre-game screen state.
	screen       Screen
	frameNo      int // 30fps tick counter; animates the title
	menuSel      int // main-menu selection
	hsTop        int // high-score scroll offset
	ls           render.LSState
	scores       map[string]int
	lsPreview    *game.Map
	lsPreviewKey string

	events chan Event
	acc    float64
	msgTTL float64
	// prev holds the last blitted frame, indexed y*W+x. Nil (or a size
	// mismatch) forces a full redraw.
	prev []render.Cell
}

// Run starts a specific level directly, skipping the title/menu flow.
func Run(m *game.Map, name string, diff game.Difficulty) error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, diff)
	a.enterGame(m, name, diff)
	return a.run()
}

// RunMenu starts on the title screen; diff preselects the difficulty on the
// level-select screen.
func RunMenu(diff game.Difficulty) error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, diff)
	a.ls.Levels = game.LevelNames()
	a.ls.Diff = diffIndex(diff)
	return a.run()
}

func diffIndex(d game.Difficulty) int {
	for i, v := range render.Difficulties {
		if v == d {
			return i
		}
	}
	return 1
}

func newApp(term *Terminal, diff game.Difficulty) *App {
	return &App{
		term:   term,
		pal:    render.Palette(),
		diff:   diff,
		events: make(chan Event, 256),
		scores: hiscore.Load(),
	}
}

// termSize returns the current terminal size, or a default 80x24 for tests
// that build an App without a real terminal.
func (a *App) termSize() (int, int) {
	if a.term == nil {
		return 80, 24
	}
	return a.term.Size()
}

// enterGame starts a game on map m: it computes the render scale from the
// current terminal size and builds fresh state and UI.
func (a *App) enterGame(m *game.Map, name string, diff game.Difficulty) {
	scale := 1
	if tw, th := a.termSize(); tw > 0 && th > 0 {
		scale = render.ComputeScale(m.W, m.H, tw, th)
	}
	a.g = game.NewStateDiff(m, diff)
	a.diff = diff
	a.level = name
	a.ui = freshUI(m, scale)
	a.ui.Paused = false
	a.layout = render.ComputeLayout(m.W, m.H, a.ui.Scale)
	a.scored = false
	a.acc = 0
	a.msgTTL = 0
	a.prev = nil
	a.screen = ScreenGame
}

// run starts the input reader, switches the terminal into raw mode and runs
// the frame loop until the app quits.
func (a *App) run() error {
	defer a.term.Close()
	startReader(a.term.in, a.events)
	a.term.AltScreen(true)
	a.term.Cursor(false)
	a.term.Mouse(true)
	defer func() {
		a.term.Mouse(false)
		a.term.AltScreen(false)
		a.term.Cursor(true)
	}()
	return a.loop()
}

// freshUI builds the UI for a new game. Selected must be pinned to
// NoSelection (see render.NoSelection) — the zero value is tower ID 0.
func freshUI(m *game.Map, scale int) render.UI {
	return render.UI{
		Cursor:   game.Vec{X: m.W / 2, Y: m.H / 2},
		Placing:  game.TowerGunner,
		Selected: render.NoSelection,
		Speed:    1,
		Paused:   true,
		Scale:    scale,
	}
}

func (a *App) loop() error {
	frame := time.NewTicker(time.Second / time.Duration(frameRate))
	defer frame.Stop()
	last := time.Now()
	for {
		// Consume resizes before rendering; RefreshSize runs on this
		// goroutine, so the cached size never races with the winch notifier.
		select {
		case <-a.term.Winch():
			a.term.RefreshSize()
			a.prev = nil
		default:
		}
		a.drainInput()
		if a.quitting {
			return nil
		}
		now := time.Now()
		real := now.Sub(last).Seconds()
		last = now
		switch a.screen {
		case ScreenGame:
			a.stepGame(real)
		case ScreenTitle:
			a.frameNo++
		}
		a.drawScreen()
		<-frame.C
	}
}

// stepGame advances the simulation by real seconds, expires transient
// messages and records the final score exactly once.
func (a *App) stepGame(real float64) {
	if a.msgTTL > 0 {
		a.msgTTL -= real
		if a.msgTTL <= 0 {
			a.ui.Message = ""
		}
	}
	if a.g.Status == game.StatusRunning && !a.ui.Paused {
		prevLives := a.g.Lives
		prevWaveActive := a.g.WaveActive
		prevWave := a.g.Wave
		a.acc += real * float64(a.ui.Speed)
		const dt = 1.0 / tickRate
		steps := 0
		for a.acc >= dt && steps < 10 {
			a.g.Step(dt)
			a.acc -= dt
			steps++
		}
		if steps == 10 {
			a.acc = 0
		}
		if a.g.Lives < prevLives {
			n := prevLives - a.g.Lives
			plural := "life"
			if n > 1 {
				plural = "lives"
			}
			a.term.Write([]byte("\a"))
			a.msg("leak! -" + strconv.Itoa(n) + " " + plural)
		}
		if prevWaveActive && !a.g.WaveActive && a.g.Wave < game.MaxWaves {
			next := a.g.Wave + 1
			if tg := game.WaveTelegraph(next); tg != "" {
				// Hold the telegraph for the whole break so it can be read.
				a.ui.Message = tg
				a.msgTTL = game.AutoWaveDelayFor(a.g.Wave)
			} else {
				a.msg("wave " + strconv.Itoa(prevWave) + " cleared +" + strconv.Itoa(game.WaveBonus(prevWave)) + "g")
			}
		}
	} else {
		a.acc = 0
	}
	if a.g.Status != game.StatusRunning && !a.scored {
		a.scored = true
		best, isNew := hiscore.Update(a.level, a.g.Score)
		a.ui.BestScore = best
		a.ui.NewBest = isNew
	}
}

// drawScreen renders the current screen. The playfield frame is
// Layout-sized and guarded; the menu screens are terminal-sized. A switch
// changes the frame size, so the blit always full-redraws.
func (a *App) drawScreen() {
	w, h := a.termSize()
	switch a.screen {
	case ScreenTitle:
		a.blit(render.RenderTitle(w, h, a.frameNo, a.scores, a.pal))
	case ScreenMenu:
		a.blit(render.RenderMenu(w, h, a.menuSel, a.pal))
	case ScreenHelp:
		a.blit(render.RenderHelp(w, h, a.pal))
	case ScreenHiscores:
		a.blit(render.RenderHighScores(w, h, a.hsTop, a.scores, a.pal))
	case ScreenLevelSelect:
		a.blit(render.RenderLevelSelect(a.lsView(), w, h, a.pal))
	default:
		a.renderGuarded(func() *render.Frame { return render.Render(a.g, &a.ui, a.pal) })
	}
}

func (a *App) msg(s string) {
	a.ui.Message = s
	a.msgTTL = 2.5
}

func (a *App) drainInput() {
	for {
		select {
		case e := <-a.events:
			a.handle(e)
		default:
			return
		}
	}
}

func (a *App) handle(e Event) {
	switch a.screen {
	case ScreenTitle:
		a.handleTitle(e)
	case ScreenMenu:
		a.handleMenu(e)
	case ScreenHelp:
		a.handleHelp(e)
	case ScreenHiscores:
		a.handleHiscores(e)
	case ScreenLevelSelect:
		a.handleLevelSelect(e)
	default:
		a.handleGame(e)
	}
}

// handleTitle: any key (or click) moves to the menu; q quits.
func (a *App) handleTitle(e Event) {
	if e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
		a.quit()
		return
	}
	if e.Mouse {
		if !e.Press {
			return
		}
	} else if e.Rune == 0 && e.Key == KeyNone {
		return
	}
	a.toScreen(ScreenMenu)
}

func (a *App) handleMenu(e Event) {
	if e.Mouse {
		a.handleMenuMouse(e)
		return
	}
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	switch e.Rune {
	case 'q', 'Q':
		a.quit()
	case 'w', 'W':
		a.moveMenu(-1)
	case 's', 'S':
		a.moveMenu(1)
	case '1', '2', '3', '4':
		a.activateMenu(int(e.Rune - '1'))
	}
	switch e.Key {
	case KeyUp:
		a.moveMenu(-1)
	case KeyDown:
		a.moveMenu(1)
	case KeyEnter:
		a.activateMenu(a.menuSel)
	case KeyEscape:
		a.toScreen(ScreenTitle)
	}
}

func (a *App) handleMenuMouse(e Event) {
	if !e.Press {
		return
	}
	w, h := a.termSize()
	switch e.Btn {
	case 64:
		a.moveMenu(-1)
		return
	case 65:
		a.moveMenu(1)
		return
	}
	for i, r := range render.MenuRects(w, h) {
		if r.Contains(e.X, e.Y) {
			a.activateMenu(i)
			return
		}
	}
}

// handleHelp: any key (or click) goes back to the menu; q quits.
func (a *App) handleHelp(e Event) {
	if e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
		a.quit()
		return
	}
	if e.Mouse {
		if !e.Press {
			return
		}
	} else if e.Rune == 0 && e.Key == KeyNone {
		return
	}
	a.toScreen(ScreenMenu)
}

func (a *App) handleHiscores(e Event) {
	if e.Mouse {
		if !e.Press {
			return
		}
		switch e.Btn {
		case 64:
			a.scrollScores(-1)
		case 65:
			a.scrollScores(1)
		}
		return
	}
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	if e.Rune == 'q' || e.Rune == 'Q' {
		a.quit()
		return
	}
	switch e.Rune {
	case 'k', 'w', 'W':
		a.scrollScores(-1)
	case 'j', 's', 'S':
		a.scrollScores(1)
	}
	switch e.Key {
	case KeyUp:
		a.scrollScores(-1)
	case KeyDown:
		a.scrollScores(1)
	case KeyEnter, KeyEscape:
		a.toScreen(ScreenMenu)
	}
}

func (a *App) scrollScores(delta int) {
	if len(a.scores) == 0 {
		return
	}
	a.hsTop += delta
	if a.hsTop < 0 {
		a.hsTop = 0
	}
	if max := len(a.scores) - 1; a.hsTop > max {
		a.hsTop = max
	}
}

func (a *App) handleLevelSelect(e Event) {
	if e.Mouse {
		a.handleLSMouse(e)
		return
	}
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	if e.Rune == 'q' || e.Rune == 'Q' {
		a.quit()
		return
	}
	// On the maze row, digits edit the seed, 'r' clears it, backspace
	// deletes.
	if a.ls.Cursor == len(a.ls.Levels) {
		switch {
		case e.Rune >= '0' && e.Rune <= '9':
			if len(a.ls.Seed) < 18 {
				a.ls.Seed += string(e.Rune)
			}
			return
		case e.Rune == 'r' || e.Rune == 'R':
			a.ls.Seed = ""
			return
		case e.Key == KeyBackspace:
			if a.ls.Seed != "" {
				a.ls.Seed = a.ls.Seed[:len(a.ls.Seed)-1]
			}
			return
		}
	}
	switch e.Key {
	case KeyUp:
		a.moveLevel(-1)
	case KeyDown:
		a.moveLevel(1)
	case KeyLeft:
		a.cycleDiff(-1)
	case KeyRight:
		a.cycleDiff(1)
	case KeyEnter:
		a.startGame()
	case KeyEscape:
		a.toScreen(ScreenMenu)
	}
}

func (a *App) handleLSMouse(e Event) {
	if !e.Press {
		return
	}
	w, h := a.termSize()
	switch e.Btn {
	case 64:
		a.moveLevel(-1)
		return
	case 65:
		a.moveLevel(1)
		return
	}
	rows, diffs, seedRect := render.LSRects(a.lsView(), w, h)
	if seedRect.Contains(e.X, e.Y) {
		a.ls.Cursor = len(a.ls.Levels)
		return
	}
	for i, d := range diffs {
		if d.Contains(e.X, e.Y) {
			a.ls.Diff = i
			return
		}
	}
	for i, r := range rows {
		if r.Contains(e.X, e.Y) {
			a.ls.Cursor = i
			return
		}
	}
}

func (a *App) moveMenu(dir int) {
	n := len(render.MenuItems)
	a.menuSel = (a.menuSel + dir + n) % n
}

func (a *App) activateMenu(i int) {
	switch i {
	case 0:
		a.toScreen(ScreenLevelSelect)
	case 1:
		a.toScreen(ScreenHelp)
	case 2:
		a.toScreen(ScreenHiscores)
	default:
		a.quit()
	}
}

func (a *App) moveLevel(dir int) {
	n := len(a.ls.Levels) + 1
	a.ls.Cursor = (a.ls.Cursor + dir + n) % n
}

func (a *App) cycleDiff(dir int) {
	n := len(render.Difficulties)
	a.ls.Diff = (a.ls.Diff + dir + n) % n
}

// startGame launches the level under the cursor: a built-in level, or a
// procedural maze (empty seed = random).
func (a *App) startGame() {
	diff := render.Difficulties[a.ls.Diff]
	if a.ls.Cursor == len(a.ls.Levels) {
		seed, ok := a.parseSeed()
		if !ok {
			a.ls.Err = "seed too large"
			return
		}
		m, err := game.MazeFromSeed(seed)
		if err != nil {
			a.ls.Err = err.Error()
			return
		}
		a.enterGame(m, fmt.Sprintf("maze%d", seed), diff)
		return
	}
	name := a.ls.Levels[a.ls.Cursor]
	m, err := game.LoadLevel(name)
	if err != nil {
		a.ls.Err = err.Error()
		return
	}
	a.enterGame(m, name, diff)
}

// parseSeed converts the typed seed to an int64. An empty string (or a seed
// of zero) means "random".
func (a *App) parseSeed() (int64, bool) {
	if a.ls.Seed == "" {
		return time.Now().UnixNano(), true
	}
	n, err := strconv.ParseInt(a.ls.Seed, 10, 64)
	if err != nil {
		return 0, false
	}
	if n == 0 {
		return time.Now().UnixNano(), true
	}
	return n, true
}

// lsView returns the level-select state with a live map preview attached.
// The preview is cached by "cursor:seed" so the map is not regenerated every
// frame; an empty or zero seed previews a fixed seed, which keeps the frame
// stable.
func (a *App) lsView() render.LSState {
	key := fmt.Sprintf("%d:%s", a.ls.Cursor, a.ls.Seed)
	if a.lsPreviewKey != key {
		a.lsPreviewKey = key
		a.lsPreview = nil
		if a.ls.Cursor == len(a.ls.Levels) {
			seed := a.ls.Seed
			if seed == "" || seed == "0" {
				seed = "1234"
			}
			if n, err := strconv.ParseInt(seed, 10, 64); err == nil {
				a.lsPreview, _ = game.MazeFromSeed(n)
			}
		} else {
			a.lsPreview, _ = game.LoadLevel(a.ls.Levels[a.ls.Cursor])
		}
	}
	v := a.ls
	v.Preview = a.lsPreview
	return v
}

// toScreen switches screens. prev is reset so the blit full-redraws, and
// per-screen state is refreshed (hiscores reloaded, stale errors cleared).
func (a *App) toScreen(s Screen) {
	a.screen = s
	a.prev = nil
	switch s {
	case ScreenHiscores:
		a.scores = hiscore.Load()
		a.hsTop = 0
	case ScreenLevelSelect:
		a.ls.Err = ""
	}
}

func (a *App) handleGame(e Event) {
	if e.Mouse {
		// Ignore all mouse input once the game is over: a click on the
		// grass behind the DEFEAT/VICTORY overlay would otherwise build a
		// tower (CanBuild/Build never check Status).
		if a.g.Status == game.StatusRunning {
			a.handleMouse(e)
		}
		return
	}
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	if a.g.Status != game.StatusRunning {
		switch e.Rune {
		case 'q', 'Q':
			a.quit()
		case 'r', 'R':
			a.restart()
		}
		return
	}
	switch e.Rune {
	case 'q', 'Q':
		a.quit()
	case 'p', 'P':
		a.ui.Paused = !a.ui.Paused
	case 'f', 'F':
		a.cycleSpeed(1)
	case 'h', 'H':
		a.ui.Help = !a.ui.Help
	case 'n', 'N':
		a.startWave()
	case 'u', 'U':
		a.upgradeSelected()
	case 'x', 'X':
		a.sellSelected()
	case 't', 'T':
		a.cycleTarget()
	case '1', '2', '3', '4', '5', '6', '7':
		k := game.TowerKind(e.Rune - '1')
		if k < game.TowerCount {
			if a.ui.PlacingOn && a.ui.Placing == k {
				a.ui.PlacingOn = false
			} else {
				a.ui.Placing = k
				a.ui.PlacingOn = true
				a.ui.Selected = -1
			}
		}
	}
	switch e.Key {
	case KeyEnter:
		a.activate()
	case KeyEscape:
		a.ui.PlacingOn = false
		a.ui.Selected = -1
	case KeyBackspace:
		a.ui.PlacingOn = false
	case KeyUp, KeyDown, KeyLeft, KeyRight:
		m := a.g.Map
		switch e.Key {
		case KeyUp:
			a.ui.Cursor.Y--
		case KeyDown:
			a.ui.Cursor.Y++
		case KeyLeft:
			a.ui.Cursor.X--
		case KeyRight:
			a.ui.Cursor.X++
		}
		if a.ui.Cursor.X < 0 {
			a.ui.Cursor.X = 0
		}
		if a.ui.Cursor.Y < 0 {
			a.ui.Cursor.Y = 0
		}
		if a.ui.Cursor.X >= m.W {
			a.ui.Cursor.X = m.W - 1
		}
		if a.ui.Cursor.Y >= m.H {
			a.ui.Cursor.Y = m.H - 1
		}
	case KeyCtrlL:
		a.prev = nil
	}
	if e.Rune >= 'a' && e.Rune <= 'z' {
		switch e.Rune {
		case 'w':
			a.moveCursor(0, -1)
		case 's':
			a.moveCursor(0, 1)
		case 'a':
			a.moveCursor(-1, 0)
		case 'd':
			a.moveCursor(1, 0)
		}
	}
}

func (a *App) moveCursor(dx, dy int) {
	a.ui.Cursor.X += dx
	a.ui.Cursor.Y += dy
	m := a.g.Map
	if a.ui.Cursor.X < 0 {
		a.ui.Cursor.X = 0
	}
	if a.ui.Cursor.Y < 0 {
		a.ui.Cursor.Y = 0
	}
	if a.ui.Cursor.X >= m.W {
		a.ui.Cursor.X = m.W - 1
	}
	if a.ui.Cursor.Y >= m.H {
		a.ui.Cursor.Y = m.H - 1
	}
}

func (a *App) mapBounds() (ox, oy, scale int) {
	return a.layout.Ox, a.layout.Oy, a.layout.Scale
}

func (a *App) handleMouse(e Event) {
	if !e.Press {
		return
	}
	// Wheel cycles through the speeds like 'f' (up = faster, wrapping),
	// instead of jumping straight to 4x.
	if e.Btn == 64 {
		a.cycleSpeed(1)
		return
	}
	if e.Btn == 65 {
		a.cycleSpeed(-1)
		return
	}
	ox, oy, sc := a.mapBounds()
	if e.Y < oy || e.Y >= oy+a.g.Map.H*sc || e.X < ox || e.X >= ox+a.g.Map.W*sc {
		a.handleMenuClick(e)
		return
	}
	cell := game.Vec{X: (e.X - ox) / sc, Y: (e.Y - oy) / sc}
	a.ui.Cursor = cell
	if a.ui.PlacingOn {
		a.place()
		return
	}
	if t := a.g.TowerAt(cell); t != nil {
		a.ui.Selected = t.ID
	} else {
		a.ui.Selected = -1
	}
}

func (a *App) handleMenuClick(e Event) {
	menuTop := a.layout.H - 4
	for _, slot := range render.MenuSlots {
		spec := game.TowerSpecs[slot.Kind]
		labelLen := len(fmt.Sprintf("%d %s %d", slot.Kind+1, spec.Name, spec.Cost[0]))
		if e.Y == menuTop+slot.Y && e.X >= slot.X && e.X < slot.X+labelLen {
			if a.ui.PlacingOn && a.ui.Placing == slot.Kind {
				a.ui.PlacingOn = false
			} else {
				a.ui.Placing = slot.Kind
				a.ui.PlacingOn = true
				a.ui.Selected = -1
			}
			return
		}
	}
}

// cycleSpeed moves to the next (dir > 0) or previous (dir < 0) speed in
// the 1x/2x/4x set, wrapping around.
func (a *App) cycleSpeed(dir int) {
	speeds := [3]int{1, 2, 4}
	for i, s := range speeds {
		if a.ui.Speed == s {
			a.ui.Speed = speeds[(i+dir+3)%3]
			return
		}
	}
	a.ui.Speed = 1
}

func (a *App) activate() {
	if a.ui.PlacingOn {
		a.place()
		return
	}
	if t := a.g.TowerAt(a.ui.Cursor); t != nil {
		a.ui.Selected = t.ID
	} else {
		a.ui.Selected = -1
	}
}

func (a *App) place() {
	if a.g.CanBuild(a.ui.Cursor, a.ui.Placing) {
		a.g.Build(a.ui.Cursor, a.ui.Placing)
		if a.g.Gold < game.TowerSpecs[a.ui.Placing].Cost[0] {
			a.ui.PlacingOn = false
		}
	} else {
		a.msg("can't build there")
	}
}

func (a *App) upgradeSelected() {
	if a.ui.Selected < 0 {
		return
	}
	t := a.g.Tower(a.ui.Selected)
	if t == nil {
		return
	}
	if t.Level >= 3 {
		a.msg("max level")
		return
	}
	if a.g.Upgrade(t) {
		a.msg(t.Spec().Name + " -> Lv" + strconv.Itoa(t.Level))
	} else {
		a.msg("need " + strconv.Itoa(a.g.UpgradeCost(t)) + " gold")
	}
}

func (a *App) sellSelected() {
	if a.ui.Selected < 0 {
		return
	}
	t := a.g.Tower(a.ui.Selected)
	if t == nil {
		return
	}
	refund := a.g.Sell(t)
	a.ui.Selected = -1
	a.msg("sold for " + strconv.Itoa(refund))
}

func (a *App) cycleTarget() {
	t, mode := a.g.CycleTarget(a.ui.Cursor)
	if t == nil {
		a.msg("no tower here")
		return
	}
	a.ui.Selected = t.ID
	a.msg(t.Spec().Name + " target: " + mode.Name())
}

func (a *App) startWave() {
	if a.g.WaveActive {
		return
	}
	if a.g.Wave >= game.MaxWaves {
		return
	}
	a.g.StartWave()
	a.msg("wave " + strconv.Itoa(a.g.Wave) + " incoming")
}

func (a *App) quit() {
	a.quitting = true
	if a.term != nil { // nil only in tests
		a.term.Mouse(false)
		a.term.AltScreen(false)
		a.term.Cursor(true)
		a.term.Close()
	}
}

func (a *App) restart() {
	m := a.g.Map
	a.g = game.NewStateDiff(m, a.diff)
	a.scored = false
	// Rebuild the UI, but keep player preferences: a zeroed UI would drop
	// the boot-computed Scale (ComputeLayout clamps it back to 1x), so the
	// playfield would shrink after every game-over restart.
	a.ui = render.UI{
		Cursor:    game.Vec{X: m.W / 2, Y: m.H / 2},
		Placing:   game.TowerGunner,
		Selected:  render.NoSelection,
		Speed:     a.ui.Speed,
		Help:      a.ui.Help,
		Scale:     a.ui.Scale,
		BestScore: hiscore.Load()[a.level],
	}
	a.acc = 0
	a.prev = nil
}

// renderGuarded blits the given frame, unless the terminal has shrunk below
// the (boot-fixed) frame size, in which case it shows an "enlarge" notice
// and pauses the game: running blind behind the notice would silently
// lose lives. The player resumes with p once the window fits again.
func (a *App) renderGuarded(makeFrame func() *render.Frame) {
	tw, th := a.term.Size()
	if tw > 0 && th > 0 && (tw < a.layout.W || th < a.layout.H) {
		a.prev = nil
		if a.g.Status == game.StatusRunning && !a.ui.Paused {
			a.ui.Paused = true
		}
		a.blit(render.RenderTooSmall(tw, th, a.layout.W, a.layout.H))
		return
	}
	a.blit(makeFrame())
}

type pen struct {
	x, y int
	fg   int
	bg   int
	bold bool
}

func (a *App) blit(f *render.Frame) {
	var buf []byte
	const chunk = 8192
	flush := func() {
		if len(buf) > 0 {
			a.term.Write(buf)
			buf = buf[:0]
		}
	}
	p := pen{x: -1, y: -1, fg: -1, bg: -1, bold: false}
	full := len(a.prev) != f.W*f.H
	if full {
		a.prev = make([]render.Cell, f.W*f.H)
		buf = append(buf, "\x1b[2J\x1b[H"...)
	}
	tw, th := a.term.Size()
	fw, fh := f.W, f.H
	if tw > 0 && tw < fw {
		fw = tw
	}
	if th > 0 && th < fh {
		fh = th
	}
	for y := 0; y < fh; y++ {
		for x := 0; x < fw; x++ {
			c := f.C[y*f.W+x]
			if !full && a.prev[y*f.W+x] == c {
				continue
			}
			if p.x != x || p.y != y {
				buf = appendf(buf, "\x1b[%d;%dH", y+1, x+1)
				p.x, p.y = x, y
			}
			fg := c.FG
			if fg == 0 {
				fg = 255
			}
			if fg != p.fg {
				buf = appendf(buf, "\x1b[38;5;%dm", fg)
				p.fg = fg
			}
			if c.BG != p.bg {
				if c.BG == 0 {
					buf = append(buf, "\x1b[49m"...)
				} else {
					buf = appendf(buf, "\x1b[48;5;%dm", c.BG)
				}
				p.bg = c.BG
			}
			if c.Bold != p.bold {
				if c.Bold {
					buf = append(buf, "\x1b[1m"...)
				} else {
					buf = append(buf, "\x1b[22m"...)
				}
				p.bold = c.Bold
			}
			if c.R == 0 {
				buf = append(buf, ' ')
			} else {
				buf = append(buf, []byte(string(c.R))...)
			}
			p.x, p.y = x+1, y
			if len(buf) > chunk {
				flush()
			}
		}
		p.x = 0
	}
	flush()
	copy(a.prev, f.C)
}

func appendf(b []byte, format string, args ...any) []byte {
	return fmt.Appendf(b, format, args...)
}
