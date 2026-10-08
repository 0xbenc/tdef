package tui

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"strconv"
	"time"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
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
	ScreenOverworld
	ScreenDragonMockup
	ScreenGrakMockup
	ScreenGnollMockup
	ScreenGunnerMockup
	ScreenFrostMockup
	ScreenPlayerMockup
	ScreenCannonierMockup
	ScreenRangerMockup
	ScreenLightningMockup
	ScreenTrebuchetMockup
	ScreenNecromancerMockup
	ScreenPaladinMockup
	ScreenRogueMockup
	ScreenMercenaryMockup
	ScreenWizardMockup
	ScreenCenturionMockup
	ScreenSquireMockup
	ScreenJournal
	ScreenCutscene
	ScreenResetProgress
	ScreenCredits
)

type App struct {
	help               render.HelpState
	helpReturn         Screen
	helpWasPaused      bool
	loadDetailTop      int
	loadDetails        bool
	loadErrors         [3]error
	saveErrors         [3]error
	saveRetryAt        time.Time
	saveQuitArmed      bool
	trainingReplay     int
	trainingReplayMask uint32
	recruitPlanning    bool
	term               *Terminal
	g                  *game.State
	ui                 render.UI
	pal                render.Colors
	layout             render.Layout
	diff               game.Difficulty
	level              string
	scored             bool
	quitting           bool

	// Pre-game screen state.
	screen       Screen
	frameNo      int // 30fps tick counter; animates the title
	titleBootAt  int // frameNo when this title visit started (boot clock)
	menuSel      int // main-menu selection
	resetSel     int // 0 = cancel; 1 = confirm reset
	resetErr     string
	resetDone    bool
	hsTop        int // high-score scroll offset
	ls           render.LSState
	scores       map[string]int
	lsPreview    *game.Map
	lsPreviewKey string
	ow           render.OWState // overworld (the lair map) state

	// Lair (overworld) plumbing: the persistent lair memory, and the run's
	// provenance, so a defense's result comes back to the map it left.
	film               render.CutsceneState
	filmReturn         Screen
	filmRemember       bool
	heartEndingPending bool
	journalMigrated    bool
	journal            *hiscore.Journal
	journalUI          render.JournalState
	journalReturn      Screen
	lair               *hiscore.Lair
	fromOW             bool // the run started from the overworld
	owFloorID          string
	defenseStartLives  int          // includes any HP bonuses granted before the first wave
	owBootArmed        bool         // the arrival cinematic has been armed this session
	owBossSeen         map[int]bool // renown -> the heart's unseal has been seen

	// Relic bonuses carried into the next defense.
	bonusGold  int
	bonusTower bool
	bonusLives int

	events chan Event
	acc    float64
	msgTTL float64
	// prev holds the last blitted frame, indexed y*W+x. Nil (or a size
	// mismatch) forces a full redraw.
	prev []render.Cell
}

// Run starts a scores-only game directly, skipping the title/menu flow.
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
	a.ow = render.NewOWState()
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
	a := &App{
		term:   term,
		pal:    render.Palette(),
		diff:   diff,
		events: make(chan Event, 256),
	}
	a.reloadScores()
	a.reloadCampaign()
	a.journal, a.loadErrors[saveJournal] = hiscore.LoadJournalWithError()
	a.ensureJournal()
	return a
}

// termSize returns the current terminal size, or a default 80x24 for tests
// that build an App without a real terminal.
func (a *App) termSize() (int, int) {
	if a.term == nil {
		return 80, 24
	}
	return a.term.Size()
}

// enterGame starts a game on map m with fresh state and UI. The render
// layout is computed from the current terminal size; it is refreshed every
// frame (drawGame), so resizes reflow the playfield live.
func (a *App) enterGame(m *game.Map, name string, diff game.Difficulty) {
	a.g = game.NewStateDiff(m, diff)
	a.diff = diff
	a.level = name
	a.ui = freshUI(m)
	a.ui.Level = name
	a.ui.Paused = false
	a.ui.ToLair = a.fromOW
	// Relic bonuses earned on the lair, spent at the Rotunda: they ride into
	// this defense and are consumed here.
	if a.fromOW {
		a.g.Gold += a.bonusGold
		a.g.Lives += a.bonusLives
		if a.bonusTower {
			if cell, ok := firstGrass(m); ok {
				// Build normally so the defender keeps its usual sale value,
				// but the relic pays its cost instead of the player.
				cost := game.TowerSpecs[game.TowerGunner].Cost[0]
				a.g.Gold += cost
				if a.g.Build(cell, game.TowerGunner) != nil {
					a.discoverTower(game.TowerGunner)
				} else {
					a.g.Gold -= cost
				}
			}
		}
		a.bonusGold, a.bonusTower, a.bonusLives = 0, false, 0
		if a.lair != nil && (a.lair.BonusGold > 0 || a.lair.BonusTower || a.lair.BonusLives > 0) {
			a.lair.BonusGold, a.lair.BonusTower, a.lair.BonusLives = 0, false, 0
			a.saveCampaign()
		}
	}
	a.defenseStartLives = a.g.Lives
	tw, th := a.termSize()
	a.layout = render.GameLayout(m.W, m.H, tw, th)
	a.scored = false
	a.heartEndingPending = false
	a.acc = 0
	a.msgTTL = 0
	a.prev = nil
	a.screen = ScreenGame
	if a.fromOW {
		a.discoverPlace(floorForLevel(name))
	}
	a.configureTraining()
}

// firstGrass is the first grass cell in reading order: where a relic's free
// tower waits.
func firstGrass(m *game.Map) (game.Vec, bool) {
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if m.At(game.Vec{X: x, Y: y}) == game.CellGrass {
				return game.Vec{X: x, Y: y}, true
			}
		}
	}
	return game.Vec{}, false
}

// run starts the input reader, switches the terminal into raw mode and runs
// the frame loop until the app quits.
func (a *App) run() error {
	defer a.term.Close()
	a.term.startInput(a.events)
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
// The playfield scale is not stored: it follows the live terminal size.
func freshUI(m *game.Map) render.UI {
	return render.UI{
		Cursor:   game.Vec{X: m.W / 2, Y: m.H / 2},
		Placing:  game.TowerGunner,
		Selected: render.NoSelection,
		Speed:    1,
		Paused:   true,
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
		if !a.saveRetryAt.IsZero() && !now.Before(a.saveRetryAt) {
			a.retrySaves()
		}
		real := now.Sub(last).Seconds()
		last = now
		switch a.screen {
		case ScreenCutscene:
			a.frameNo++
			a.tickCutsceneElapsed(real)
		case ScreenGame:
			a.stepGame(real)
		case ScreenOverworld:
			a.frameNo++
			a.owTick()
		case ScreenTitle, ScreenMenu:
			a.frameNo++
		}
		a.drawScreen()
		<-frame.C
	}
}

// stepGame advances the gameplay animation clock even while paused or ended,
// steps the simulation, expires messages and records the final score once.
func (a *App) stepGame(real float64) {
	a.frameNo++ // victory/defeat cinematics must advance after simulation stops
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
			if a.g.LessonPending {
				a.syncTraining()
				break
			}
		}
		if steps == 10 {
			a.acc = 0
		}
		if a.g.Lives < prevLives {
			n := prevLives - a.g.Lives
			a.term.Write([]byte("\a"))
			if n == 6 {
				a.msg(copytext.Text("ui.step_game.the_player_breached_6"))
			} else {
				a.msg(copytext.Format("ui.status.breach", "lives", strconv.Itoa(n)))
			}
		}
		if prevWaveActive && !a.g.WaveActive && a.g.Wave < game.MaxWaves {
			next := a.g.Wave + 1
			if tg := game.WaveTelegraphFor(a.g.Map, next); tg != "" {
				// Hold the telegraph for the whole break so it can be read.
				a.ui.Message = tg
				a.msgTTL = game.AutoWaveDelayFor(a.g.Wave)
			} else {
				a.msg(copytext.Format("ui.status.wave_cleared", "wave", strconv.Itoa(prevWave), "gold", strconv.Itoa(a.g.Map.ScaleGold(game.WaveBonus(prevWave)))))
			}
		}
	} else {
		a.acc = 0
	}
	a.revealMenuAfterFirstWave()
	a.recordEnemyDiscoveries()
	if a.g.Status != game.StatusRunning && !a.scored {
		a.scored = true
		a.ui.EndAtFrame = a.frameNo // the end cinematics run off this clock
		won := a.g.Status == game.StatusVictory
		best, isNew := 0, false
		if a.trainingReplay == 0 {
			if a.scores == nil {
				a.reloadScores()
			}
			best, isNew = hiscore.RecordScore(hiscore.Table(a.scores), a.level, a.g.Score)
			if isNew {
				a.saveScores()
			}
		}
		a.ui.BestScore = best
		a.ui.NewBest = isNew
		if a.fromOW {
			// Only defenses entered through the lair change campaign progress.
			d := diffIndex(a.diff)
			floor := a.owFloorID
			if floor == "" {
				floor = floorForLevel(a.level)
			}
			a.heartEndingPending = won && floor == hiscore.HeartFloor && !a.lair.BossHeld(d)
			a.completeTrainingDefense()
			a.lair.RecordResult(floor, d, a.g.Wave, won)
			if won {
				a.journal.RecordVictory(floor, d)
				a.saveJournal()
			}
			a.lair.Tokens += a.earnedRelics()
			a.saveCampaign()
			a.owUnsealCheck(d)
			a.owSetBanner(won, best, isNew)
			a.owHeartBlastCheck(d)
		}
	}
	if a.heartEndingPending && a.frameNo-a.ui.EndAtFrame >= 90 {
		a.heartEndingPending = false
		a.startCutscene(render.FilmEnding, ScreenGame, false)
	}
}

// floorForLevel maps a level id to its corresponding lair floor.
func floorForLevel(level string) string {
	switch {
	case level == "canyon":
		return "rift"
	case level == "hub":
		return "rotunda"
	case level == "winding":
		return "halls"
	case level == "garden":
		return "garden"
	case len(level) > 4 && level[:4] == "maze":
		return hiscore.DepthsFloor
	}
	return level // "heart" passes through
}

// owUnsealCheck starts the unseal cascade for any floor that just opened on
// the current renown (the cascade plays when Grak returns to the map).
func (a *App) owUnsealCheck(d int) {
	st := &a.ow
	opened := func(id string, cond bool) {
		if cond && !st.Unlocked[id] && st.Unsealing[id] == 0 {
			st.Unsealing[id] = render.OWUnsealFrames
		}
	}
	opened("rift", a.lair.Floor("rotunda", d).Cleared)
	opened("halls", a.lair.Floor("rift", d).Cleared)
	opened("garden", a.lair.Floor("halls", d).Cleared)
	opened("depths", a.lair.DepthsReady(d))
}

// owHeartBlastCheck arms the heart-unseal shockwave when a defense just
// unsealed the heart at the current renown (BossReady false -> true). A
// renown whose heart is already unsealed never re-fires it: owRefresh marks
// each renown seen as the lair is entered.
func (a *App) owHeartBlastCheck(d int) {
	if a.lair.BossReady(d) && !a.owBossSeen[d] {
		a.ow.BlastTTL = render.OWBlastFrames
	}
	a.owBossSeen[d] = a.lair.BossReady(d)
}

// Relics reward a complete campaign defense with no HP lost. Checking leaks
// also prevents any future healing mechanic from concealing a breach.
func (a *App) earnedRelics() int {
	if !a.fromOW || a.trainingReplay != 0 || a.g == nil || a.g.Status != game.StatusVictory || a.g.Wave != game.MaxWaves {
		return 0
	}
	if a.g.TotalLeaks != 0 || a.g.Lives != a.defenseStartLives {
		return 0
	}
	return diffIndex(a.diff) + 1 // Easy, Normal, Hard: 1, 2, 3.
}

// owSetBanner writes the result line Grak reads back on the map.
func (a *App) owSetBanner(won bool, best int, isNew bool) {
	st := &a.ow
	name := render.OWFloorName(a.owFloorID)
	var msg string
	switch {
	case won && a.owFloorID == render.HeartFloorID:
		msg = copytext.Text("overworld.ow_set_banner.the_heart_is_held_malgrath_endures")
	case won:
		msg = copytext.Format("overworld.results.held", "floor", name)
	default:
		msg = copytext.Format("overworld.results.lost", "floor", name, "wave", strconv.Itoa(a.g.Wave))
	}
	if relics := a.earnedRelics(); relics > 0 {
		msg += copytext.Format("overworld.results.relic_reward", "count", strconv.Itoa(relics))
	}
	if isNew {
		msg += copytext.Format("overworld.results.new_best", "score", strconv.Itoa(best))
	}
	st.ReturnMsg = msg
	st.ReturnKind = render.OWMessageDefeat
	if won {
		st.ReturnKind = render.OWMessageSuccess
	}
	st.ReturnTTL = 180
	st.ReturnFX = render.OWReturnFX{Floor: a.owFloorID, Won: won}
}

// drawScreen renders the current screen. Every screen, game included, is
// terminal-sized. A screen switch resets prev, so the blit always
// full-redraws.
func (a *App) drawScreen() {
	w, h := a.termSize()
	switch a.screen {
	case ScreenCutscene:
		a.blit(render.RenderCutscene(w, h, a.film))
	case ScreenJournal:
		a.blit(render.RenderJournal(w, h, a.journalUI))
	case ScreenTitle:
		a.blit(render.RenderTitle(w, h, a.frameNo, a.frameNo-a.titleBootAt, a.scores, a.pal))
	case ScreenMenu:
		a.blit(render.RenderMenuAnimated(w, h, a.menuSel, a.frameNo, a.lair.HasMenuReveal(), a.pal))
	case ScreenResetProgress:
		a.blit(render.RenderResetProgress(w, h, a.resetSel, a.resetDone, a.resetErr, a.pal))
	case ScreenHelp:
		a.blit(render.RenderHelpPage(w, h, a.help, a.pal))
	case ScreenCredits:
		a.blit(render.RenderCredits(w, h, a.pal))
	case ScreenHiscores:
		a.blit(render.RenderHighScores(w, h, a.hsTop, a.scores, a.pal))
	case ScreenLevelSelect:
		a.blit(render.RenderLevelSelect(a.lsView(), w, h, a.pal))
	case ScreenFrostMockup:
		a.blit(render.RenderFrostMockup(w, h))
	case ScreenCenturionMockup:
		a.blit(render.RenderCenturionMockup(w, h))
	case ScreenSquireMockup:
		a.blit(render.RenderSquireMockup(w, h))
	case ScreenWizardMockup:
		a.blit(render.RenderWizardMockup(w, h))
	case ScreenMercenaryMockup:
		a.blit(render.RenderMercenaryMockup(w, h))
	case ScreenRogueMockup:
		a.blit(render.RenderRogueMockup(w, h))
	case ScreenPaladinMockup:
		a.blit(render.RenderPaladinMockup(w, h))
	case ScreenNecromancerMockup:
		a.blit(render.RenderNecromancerMockup(w, h))
	case ScreenTrebuchetMockup:
		a.blit(render.RenderTrebuchetMockup(w, h))
	case ScreenLightningMockup:
		a.blit(render.RenderLightningMockup(w, h))
	case ScreenRangerMockup:
		a.blit(render.RenderRangerMockup(w, h))
	case ScreenCannonierMockup:
		a.blit(render.RenderCannonierMockup(w, h))
	case ScreenPlayerMockup:
		a.blit(render.RenderPlayerMockup(w, h))
	case ScreenGunnerMockup:
		a.blit(render.RenderGunnerMockup(w, h))
	case ScreenGnollMockup:
		a.blit(render.RenderGnollMockup(w, h))
	case ScreenGrakMockup:
		a.blit(render.RenderGrakMockup(w, h))
	case ScreenDragonMockup:
		a.blit(render.RenderDragonMockup(w, h))
	case ScreenOverworld:
		a.blit(render.RenderOverworld(w, h, a.ow, a.frameNo, a.pal))
	default:
		a.drawGame()
	}
}

// drawGame renders the in-game frame at the live terminal size, recomputing
// the layout each frame so resizes reflow (scale step, centering, menu
// slots) without a restart. Below MinFrame the sim is paused and an
// "enlarge" notice is shown: running blind behind it would silently lose
// lives. The player resumes with p once the window fits again.
func (a *App) drawGame() {
	tw, th := a.termSize()
	a.layout = render.GameLayout(a.g.Map.W, a.g.Map.H, tw, th)
	mw, mh := render.MinFrame(a.g.Map.W, a.g.Map.H)
	if a.tooSmall(tw, th, mw, mh) {
		a.prev = nil
		if a.g.Status == game.StatusRunning && !a.ui.Paused {
			a.ui.Paused = true
		}
		a.blit(render.RenderTooSmall(tw, th, mw, mh))
		return
	}
	a.blit(render.Render(a.g, &a.ui, a.pal, tw, th, a.frameNo))
}

// tooSmall reports whether a tw×th terminal cannot show the playfield at
// scale 1. A 0×0 (unknown) size is not "too small": it renders normally.
func (a *App) tooSmall(tw, th, mw, mh int) bool {
	return tw > 0 && th > 0 && (tw < mw || th < mh)
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
	if a.hasLoadErrors() && e.Key == KeyCtrlL {
		a.loadDetails = !a.loadDetails
		a.loadDetailTop = 0
		if a.loadDetails && a.screen == ScreenGame {
			a.ui.Paused = true
			a.acc = 0
		}
		return
	}
	if a.hasLoadErrors() && a.loadDetails && e.Key != KeyCtrlC && e.Rune != 'q' && e.Rune != 'Q' {
		if e.Key == KeyEscape {
			a.loadDetails = false
		}
		if e.Key == KeyUp || e.Rune == 'w' {
			a.loadDetailTop = max(0, a.loadDetailTop-1)
		}
		if e.Key == KeyDown || e.Rune == 's' {
			a.loadDetailTop++
		}
		return
	}
	if a.hasSaveErrors() && e.Key == KeyCtrlL {
		a.retrySaves()
		a.prev = nil
		return
	}
	if (a.hasSaveErrors() || a.hasLoadErrors()) && a.saveQuitArmed && (e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q'))) {
		a.quit()
		return
	}
	switch a.screen {
	case ScreenCutscene:
		a.handleCutscene(e)
	case ScreenJournal:
		a.handleJournal(e)
	case ScreenTitle:
		a.handleTitle(e)
	case ScreenMenu:
		a.handleMenu(e)
	case ScreenResetProgress:
		a.handleResetProgress(e)
	case ScreenHelp:
		a.handleHelp(e)
	case ScreenCredits:
		a.handleDismiss(e)
	case ScreenHiscores:
		a.handleHiscores(e)
	case ScreenLevelSelect:
		a.handleLevelSelect(e)
	case ScreenDragonMockup, ScreenGrakMockup, ScreenGnollMockup, ScreenGunnerMockup, ScreenFrostMockup, ScreenPlayerMockup, ScreenCannonierMockup, ScreenRangerMockup, ScreenLightningMockup, ScreenTrebuchetMockup, ScreenNecromancerMockup, ScreenPaladinMockup, ScreenRogueMockup, ScreenMercenaryMockup, ScreenWizardMockup, ScreenCenturionMockup, ScreenSquireMockup:
		a.handlePortraitMockup(e)
	case ScreenOverworld:
		a.handleOverworld(e)
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
	case '1', '2', '3', '4', '5', '6', '7', '8':
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

// Help and credits: any key (or click) goes back to the menu; q quits.
func (a *App) handleDismiss(e Event) {
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
	case render.MenuStart:
		a.toScreen(ScreenOverworld) // Start: the lair itself
	case render.MenuQuickPlay:
		a.toScreen(ScreenLevelSelect)
	case render.MenuHelp:
		a.openHelp()
	case render.MenuHighScores:
		a.toScreen(ScreenHiscores)
	case render.MenuJournal:
		a.openJournal()
	case render.MenuCredits:
		a.toScreen(ScreenCredits)
	case render.MenuResetProgress:
		a.toScreen(ScreenResetProgress)
	case render.MenuQuit:
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
			a.ls.Err = copytext.Text("ui.start_game.seed_too_large")
			return
		}
		m, err := game.MazeFromSeed(seed)
		if err != nil {
			a.ls.Err = err.Error()
			return
		}
		a.enterQuickGame(m, fmt.Sprintf("maze%d", seed), diff)
		return
	}
	name := a.ls.Levels[a.ls.Cursor]
	m, err := game.LoadLevel(name)
	if err != nil {
		a.ls.Err = err.Error()
		return
	}
	a.enterQuickGame(m, name, diff)
}

// A prior lair run must not make a subsequent Quick Play run a campaign run.
func (a *App) enterQuickGame(m *game.Map, name string, diff game.Difficulty) {
	a.fromOW = false
	a.owFloorID = ""
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
	case ScreenResetProgress:
		a.resetSel, a.resetErr, a.resetDone = 0, "", false
	case ScreenTitle:
		a.titleBootAt = a.frameNo // replay the boot cinematic
	case ScreenHiscores:
		if a.saveErrors[saveScores] == nil && a.loadErrors[saveScores] == nil {
			a.reloadScores()
		}
		a.hsTop = 0
	case ScreenLevelSelect:
		a.ls.Err = ""
	case ScreenOverworld:
		a.owRefresh()
		if a.maybeOpening() {
			return
		}
		if !a.owBootArmed {
			// The arrival cinematic plays once per session, on first entry.
			a.owBootArmed = true
			a.ow.BootTTL = render.OWBootFrames
		}
	}
}

func (a *App) handleGame(e Event) {
	if a.g.LessonPending {
		if e.Key == KeyCtrlC || e.Rune == 'q' || e.Rune == 'Q' {
			a.quit()
			return
		}
		if e.Key == KeyEnter || e.Rune == ' ' || (e.Mouse && e.Press && e.Btn == 0) {
			a.acknowledgeRecruit()
		}
		return
	}

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
	if (e.Rune == 'j' || e.Rune == 'J') && (a.ui.Paused || a.g.Status != game.StatusRunning) {
		a.openJournal()
		return
	}
	if a.g.Status != game.StatusRunning {
		if a.heartEndingPending && (e.Key == KeyEscape || e.Rune == 'r' || e.Rune == 'R') {
			a.heartEndingPending = false
			a.startCutscene(render.FilmEnding, ScreenGame, false)
			return
		}
		switch {
		case (e.Rune == 'v' || e.Rune == 'V') && a.g.Status == game.StatusVictory && floorForLevel(a.level) == hiscore.HeartFloor:
			a.heartEndingPending = false
			a.startCutscene(render.FilmEnding, ScreenGame, false)
		case e.Rune == 'q' || e.Rune == 'Q':
			a.quit()
		case e.Rune == 'r' || e.Rune == 'R':
			a.restart()
		case e.Key == KeyEscape:
			a.leaveGame()
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
		a.openHelp()
	case 'n', 'N':
		a.startWave()
	case 'u', 'U':
		a.upgradeSelected()
	case 'x', 'X':
		a.sellSelected()
	case 't', 'T':
		a.cycleTarget()
	case '[', ']':
		delta := 1
		if e.Rune == '[' {
			delta = -1
		}
		a.changeRoster(delta)
	case 'r', 'R':
		a.rotateForge()
	case '1', '2', '3', '4', '5', '6', '7':
		k := game.TowerCount
		for _, slot := range render.TowerSlots(a.layout.W, a.layout.H, a.ui.RosterPage) {
			if slot.Key == int(e.Rune-'0') {
				k = slot.Kind
			}
		}
		if k < game.TowerCount {
			if !a.g.TowerAvailable(k) {
				a.lockedRecruit(k)
				return
			}
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
	case KeyTab:
		a.nextDefender()
	case KeyUp:
		a.moveCursor(0, -1)
	case KeyDown:
		a.moveCursor(0, 1)
	case KeyLeft:
		a.moveCursor(-1, 0)
	case KeyRight:
		a.moveCursor(1, 0)
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
	a.ui.Selected = -1
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
	// Same slot geometry the renderer draws (TowerSlots), so clicks can
	// never drift from the labels across a resize.
	pager := render.RosterPager(a.layout.W, a.layout.H)
	if e.Y == pager.Y && e.X >= pager.X && e.X < pager.X+pager.W {
		delta := 1
		if e.X < pager.X+pager.W/2 {
			delta = -1
		}
		a.changeRoster(delta)
		return
	}
	for _, slot := range render.TowerSlots(a.layout.W, a.layout.H, a.ui.RosterPage) {
		if e.Y == slot.Y && e.X >= slot.X && e.X < slot.X+slot.W {
			if !a.g.TowerAvailable(slot.Kind) {
				a.lockedRecruit(slot.Kind)
				return
			}
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
	if !a.g.TowerAvailable(a.ui.Placing) {
		a.lockedRecruit(a.ui.Placing)
		return
	}
	if a.g.SpecialistPlaced(a.ui.Placing) {
		a.msg(copytext.Text("ui.specialists.one_per_level"))
		a.ui.PlacingOn = false
		return
	}

	if tower := a.g.Build(a.ui.Cursor, a.ui.Placing); tower != nil {
		a.g.Aim(tower, a.ui.Facing)
		a.discoverTower(a.ui.Placing)
		a.ui.PlacingOn = false
		a.ui.Selected = -1
		if a.recruitPlanning {
			a.recruitPlanning = false
			a.ui.Paused = false
			a.acc = 0
		}
	} else {
		a.msg(copytext.Text("ui.place.can_t_build_there"))
	}

}

func (a *App) upgradeSelected() {
	t := a.focusedTower()
	if t == nil {
		a.msg(copytext.Text("ui.keyboard.select_defender"))
		return
	}
	a.ui.Selected = t.ID
	if t.Level >= 3 {
		a.msg(copytext.Text("ui.upgrade_selected.max_level"))
		return
	}
	if a.g.Upgrade(t) {
		a.msg(copytext.Format("ui.status.upgraded", "tower", t.Spec().Name, "level", strconv.Itoa(t.Level)))
	} else {
		a.msg(copytext.Format("ui.status.need_gold", "gold", strconv.Itoa(a.g.UpgradeCost(t))))
	}
}

func (a *App) sellSelected() {
	t := a.focusedTower()
	if t == nil {
		return
	}
	refund := a.g.Sell(t)
	a.ui.Selected = -1
	a.msg(copytext.Format("ui.status.sold", "gold", strconv.Itoa(refund)))
}

func (a *App) cycleTarget() {
	t := a.focusedTower()
	if t == nil {
		a.msg(copytext.Text("ui.cycle_target.no_tower_here"))
		return
	}
	if t.Kind == game.TowerRuneforge {
		a.msg(copytext.Text("ui.specialists.aim_prompt"))
		return
	}
	if t.Kind == game.TowerSappers {
		a.msg(copytext.Text("ui.keyboard.automatic_mines"))
		return
	}
	t.TargetMode = t.TargetMode.Next()
	a.ui.Selected = t.ID
	a.msg(copytext.Format("ui.status.target_changed", "tower", t.Spec().Name, "target", render.TargetExplanation(t.TargetMode)))
}

func (a *App) startWave() {
	if a.g.WaveActive {
		return
	}
	if a.g.Wave >= game.MaxWaves {
		return
	}
	// The siege waits for the player to commit a tower; the early-start key
	// can't bypass that.
	if a.g.Wave == 0 && len(a.g.Towers) == 0 {
		a.msg(copytext.Text("ui.start_wave.the_siege_waits_build_a_tower_to"))
		return
	}
	a.g.StartWave()
	if a.g.TrainingStage != 0 && a.g.WaveActive {
		a.ui.Paused = false
	}
	a.msg(copytext.Format("ui.status.wave_incoming", "wave", strconv.Itoa(a.g.Wave)))
}

func (a *App) quit() {
	a.retrySaves()
	if (a.hasSaveErrors() || a.hasLoadErrors()) && !a.saveQuitArmed {
		a.saveQuitArmed = true
		a.prev = nil
		return
	}
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
	a.heartEndingPending = false
	a.defenseStartLives = a.g.Lives
	// Rebuild the UI, but keep player preferences (speed, help). The
	// playfield scale follows the live terminal size, so there is nothing
	// to preserve there.
	a.ui = render.UI{
		Cursor:    game.Vec{X: m.W / 2, Y: m.H / 2},
		Placing:   game.TowerGunner,
		Selected:  render.NoSelection,
		Speed:     a.ui.Speed,
		Help:      a.ui.Help,
		Level:     a.level,
		BestScore: a.scores[a.level],
		ToLair:    a.fromOW,
	}
	a.acc = 0
	a.prev = nil
	a.configureTraining()
}

// leaveGame ends the run from the game-over box: back to the lair map (with
// the result banner) when the run started from the overworld, otherwise back
// to the level select (or the title, when there is no menu state).
func (a *App) leaveGame() {
	if a.trainingReplay != 0 {
		a.journal, a.loadErrors[saveJournal] = hiscore.LoadJournalWithError()
		a.journalMigrated = false
	}
	a.trainingReplay, a.trainingReplayMask = 0, 0
	if a.fromOW {
		a.toScreen(ScreenOverworld)
		return
	}
	// A finished Quick Play run must not suppress the first campaign opening.
	a.g = nil
	if len(a.ls.Levels) > 0 {
		a.toScreen(ScreenLevelSelect)
		return
	}
	a.toScreen(ScreenTitle)
}

type pen struct {
	x, y int
	fg   int
	bg   int
	bold bool
}

// blitWriter is the minimal sink the frame diff writes to: the terminal in
// production, a capture buffer in tests.
type blitWriter interface{ Write(p []byte) }

func (a *App) blit(f *render.Frame) {
	if a.hasLoadErrors() {
		render.DrawLoadFailure(f, a.loadErrorText(), a.loadDetails, a.saveQuitArmed, a.loadDetailTop)
	} else if a.hasSaveErrors() {
		render.DrawSaveFailure(f, a.saveErrorText(), a.saveQuitArmed)
	}
	a.blitTo(f, a.term)
}

// blitTo diffs the frame against the previous one and writes only the changed
// cells to w, as ANSI cursor-moves + 256-colour SGR + the rune.
func (a *App) blitTo(f *render.Frame, w blitWriter) {
	var buf []byte
	const chunk = 8192
	flush := func() {
		if len(buf) > 0 {
			w.Write(buf)
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
			// Symbols can occupy one or two columns depending on the terminal.
			// Anchor the next cell instead of letting a wide glyph shift HUD digits.
			if c.R > 127 && !(c.R >= 0x2500 && c.R <= 0x259f) && !(c.R >= 0x2800 && c.R <= 0x28ff) {
				p.x = -1
			}

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
