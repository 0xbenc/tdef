package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/render"
)

// RunCutscene is a review/replay entry point. It never grants progression or
// marks the campaign introduction seen; leaving it returns to the lair.
func RunCutscene(film render.Film) error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, game.Normal)
	a.ow = render.NewOWState()
	a.owRefresh()
	a.owBootArmed = true
	a.startCutscene(film, ScreenOverworld, false)
	return a.run()
}

func (a *App) startCutscene(film render.Film, destination Screen, remember bool) {
	a.film = render.CutsceneState{Film: film}
	a.filmReturn = destination
	a.filmRemember = remember
	a.screen = ScreenCutscene
	a.prev = nil
	a.acc = 0
}

func (a *App) finishCutscene() {
	if a.filmRemember {
		a.lair.IntroSeen = true
		a.saveCampaign()
	}
	a.screen = a.filmReturn
	a.prev = nil
	a.acc = 0
	// Restore the completed defense's result, past the short combat outro.
	// Rewards were committed before the film, so skipping and replaying cannot
	// score twice, restart the simulation, or spend a relic.
	if a.screen == ScreenGame && a.g != nil && a.g.Status != game.StatusRunning {
		a.ui.EndAtFrame = a.frameNo - 90
	}
}

// Keep the film clock at 30 frames per second even when artwork takes longer
// to render. Caption speed must not depend on the shot's rendering cost.
func (a *App) tickCutsceneElapsed(seconds float64) {
	w, h := a.termSize()
	if w < 62 || h < 19 {
		return
	}
	a.acc += seconds * 30
	frames := int(a.acc)
	a.acc -= float64(frames)
	a.film.Frame += frames
}

func (a *App) tickCutscene() {
	a.tickCutsceneElapsed(1.0 / 30)
}

func (a *App) handleCutscene(e Event) {
	if e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
		a.quit()
		return
	}
	if !e.Mouse && e.Key == KeyEscape {
		a.finishCutscene()
		return
	}
	w, h := a.termSize()
	if w < 62 || h < 19 {
		return
	}
	if !e.Mouse && e.Key == KeyLeft {
		if a.film.Shot > 0 {
			a.film.Shot--
			a.film.Frame = 0
			a.film.Revealed = false
			a.prev = nil
		}
		return
	}
	advance := (!e.Mouse && (e.Key == KeyEnter || e.Key == KeyRight || e.Rune == ' ')) || (e.Mouse && e.Press && e.Btn == 0)
	if !advance {
		return
	}
	shots := render.FilmShots(a.film.Film)
	if render.CutsceneVisible(a.film) < len([]rune(shots[a.film.Shot].Dialogue)) {
		a.film.Revealed = true
		return
	}
	if !render.CutsceneCanAdvance(a.film) {
		return
	}
	if a.film.Shot+1 >= len(shots) {
		a.finishCutscene()
		return
	}
	a.film.Shot++
	a.film.Frame = 0
	a.film.Revealed = false
	a.prev = nil
}

// Returning campaigns keep their familiar arrival. A genuinely fresh lair
// gets the opening once; old saves with progress are never forced through it.
func (a *App) maybeOpening() bool {
	if a.g == nil && a.lair != nil && !a.lair.IntroSeen && !a.lair.AnyRecord() && !a.owBootArmed {
		a.owBootArmed = true
		a.ow.BootTTL = 0
		a.startCutscene(render.FilmOpening, ScreenOverworld, true)
		return true
	}
	return false
}
