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

type App struct {
	term     *Terminal
	g        *game.State
	ui       render.UI
	pal      render.Colors
	layout   render.Layout
	diff     game.Difficulty
	level    string
	scored   bool
	started  bool
	quitting bool

	events chan Event
	acc    float64
	msgTTL float64
	overT  time.Time
	// prev holds the last blitted frame, indexed y*W+x. Nil (or a size
	// mismatch) forces a full redraw.
	prev []render.Cell
}

func Run(m *game.Map, name string, diff game.Difficulty) error {
	term, err := Open()
	if err != nil {
		return err
	}
	defer term.Close()
	a := &App{
		term:   term,
		g:      game.NewStateDiff(m, diff),
		pal:    render.Palette(),
		diff:   diff,
		level:  name,
		events: make(chan Event, 256),
	}
	scale := 1
	if tw, th := a.term.Size(); tw > 0 && th > 0 {
		scale = render.ComputeScale(m.W, m.H, tw, th)
	}
	a.ui = freshUI(m, scale)
	a.started = false
	a.layout = render.ComputeLayout(m.W, m.H, a.ui.Scale)
	startReader(a.term.in, a.events)
	term.AltScreen(true)
	term.Cursor(false)
	term.Mouse(true)
	defer func() {
		term.Mouse(false)
		term.AltScreen(false)
		term.Cursor(true)
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
		// Consume resizes before rendering in BOTH the intro gate and the
		// running path; RefreshSize runs on this goroutine, so the cached
		// size never races with the winch notifier.
		select {
		case <-a.term.Winch():
			a.term.RefreshSize()
			a.prev = nil
		default:
		}
		if !a.started {
			select {
			case e := <-a.events:
				if e.Key == KeyCtrlC || e.Rune == 'q' {
					a.quit()
					return nil
				}
				a.started = true
				a.ui.Paused = false
			default:
			}
			a.renderGuarded(func() *render.Frame {
				return render.RenderIntro(a.g.Map, a.level, a.diff, a.pal, a.ui.Scale)
			})
			<-frame.C
			continue
		}
		a.drainInput()
		if a.quitting {
			return nil
		}
		now := time.Now()
		real := now.Sub(last).Seconds()
		last = now
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
		a.draw()
		<-frame.C
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
	if e.Mouse {
		a.handleMouse(e)
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
	case '\r':
		a.activate()
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
	a.term.Mouse(false)
	a.term.AltScreen(false)
	a.term.Cursor(true)
	a.term.Close()
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

func (a *App) draw() {
	a.renderGuarded(func() *render.Frame { return render.Render(a.g, &a.ui, a.pal) })
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
