package tui

import (
	"reflect"
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
)

func TestOpeningOnceAndLegacyCampaignContinues(t *testing.T) {
	a := owTestApp(t)
	a.toScreen(ScreenOverworld)
	if a.screen != ScreenCutscene || a.film.Film != render.FilmOpening || a.ow.BootTTL != 0 {
		t.Fatal("fresh campaign did not start the opening")
	}
	if hiscore.LoadLair().IntroSeen {
		t.Fatal("intro remembered before viewing/skipping")
	}
	a.handle(Event{Key: KeyEscape})
	if a.screen != ScreenOverworld || !hiscore.LoadLair().IntroSeen {
		t.Fatal("skip did not remember the intro and enter the lair")
	}
	a.toScreen(ScreenTitle)
	a.toScreen(ScreenOverworld)
	if a.screen != ScreenOverworld || a.ow.BootTTL != 0 {
		t.Fatal("opening or old arrival replayed")
	}
	// A new process with an intro-only save also avoids replaying the film.
	b := &App{lair: hiscore.LoadLair(), ow: render.NewOWState()}
	b.toScreen(ScreenOverworld)
	if b.screen != ScreenOverworld {
		t.Fatal("intro-only save replayed film")
	}
	c := owTestApp(t)
	c.lair.Record("rotunda", 1, 2, false)
	c.toScreen(ScreenOverworld)
	if c.screen != ScreenOverworld {
		t.Fatal("legacy campaign was forced through the opening")
	}
}

func TestCutsceneReviewInputAndExactReturn(t *testing.T) {
	a := journalGameApp(t)
	a.g.Gold = 827
	a.ui.Paused = true
	a.acc = 1.2
	g := *a.g
	ui := a.ui
	ow := a.ow
	a.startCutscene(render.FilmOpening, ScreenGame, false)
	a.handle(Event{Key: KeyEnter})
	if a.film.Shot != 0 || !a.film.Revealed {
		t.Fatal("first advance must reveal before cutting")
	}
	a.handle(Event{Mouse: true, Press: false, Btn: 0})
	a.handle(Event{Mouse: true, Press: true, Btn: 64})
	if a.film.Shot != 0 {
		t.Fatal("release/wheel advanced film")
	}
	a.handle(Event{Rune: ' '})
	if a.film.Shot != 0 {
		t.Fatal("advance cut the panorama before its end hold")
	}
	a.film.Frame = 10000
	a.handle(Event{Rune: ' '})
	if a.film.Shot != 1 || a.film.Revealed || a.film.Frame != 0 {
		t.Fatal("second advance did not cut cleanly")
	}
	a.handle(Event{Key: KeyLeft})
	if a.film.Shot != 0 {
		t.Fatal("previous shot inaccessible")
	}
	// Every directed shot can be read and advanced by either keyboard or click.
	for a.screen == ScreenCutscene {
		a.film.Frame = 10000
		a.handle(Event{Mouse: true, Press: true, Btn: 0})
	}
	if a.screen != ScreenGame || !reflect.DeepEqual(*a.g, g) || !reflect.DeepEqual(a.ui, ui) || !reflect.DeepEqual(a.ow, ow) || a.acc != 0 {
		t.Fatal("review mutated the caller")
	}
	if hiscore.LoadLair().IntroSeen {
		t.Fatal("preview changed campaign intro state")
	}
	a.startCutscene(render.FilmEnding, ScreenGame, false)
	a.handle(Event{Rune: 'q'})
	if !a.quitting {
		t.Fatal("q did not exit film")
	}
}

func TestHeartEndingAfterRewardsOncePerRenown(t *testing.T) {
	for _, diff := range render.Difficulties {
		a := owTestApp(t)
		a.owRefresh()
		a.owBootArmed = true
		a.fromOW = true
		a.owFloorID = "heart"
		m, err := game.LoadBoss()
		if err != nil {
			t.Fatal(err)
		}
		a.enterGame(m, "heart", diff)
		a.g.Status = game.StatusVictory
		a.g.Wave = game.MaxWaves
		a.g.Score = 12345
		a.stepGame(.05)
		if !a.scored || !a.heartEndingPending || a.screen != ScreenGame {
			t.Fatal("Heart did not arm its ending after scoring")
		}
		if !hiscore.LoadLair().BossHeld(diffIndex(diff)) || hiscore.LoadJournal().EndgameWins != 1 {
			t.Fatal("rewards missing before the film")
		}
		tokens, time, gold := a.lair.Tokens, a.g.Time, a.g.Gold
		ow := a.ow
		for a.screen == ScreenGame {
			a.stepGame(.05)
		}
		if a.screen != ScreenCutscene || a.film.Film != render.FilmEnding || a.heartEndingPending {
			t.Fatal("combat outro did not reach the ending")
		}
		for i := 0; i < 180; i++ {
			a.tickCutscene()
		}
		if a.g.Time != time || a.g.Gold != gold || !reflect.DeepEqual(a.ow, ow) {
			t.Fatal("film advanced gameplay or overworld")
		}
		a.handle(Event{Key: KeyEscape})
		if a.screen != ScreenGame || a.frameNo-a.ui.EndAtFrame < 90 {
			t.Fatal("skip did not reach the finished result")
		}
		for i := 0; i < 120; i++ {
			a.stepGame(.05)
		}
		if a.screen != ScreenGame || a.lair.Tokens != tokens || hiscore.LoadJournal().EndgameWins != 1 {
			t.Fatal("ending repeated or awarded twice")
		}
		a.handle(Event{Rune: 'v'})
		if a.screen != ScreenCutscene {
			t.Fatal("result ending replay missing")
		}
		a.handle(Event{Key: KeyEscape})
		a.handle(Event{Key: KeyEscape})
		if a.screen != ScreenOverworld {
			t.Fatal("ending did not restore the lair return route")
		}
		a.enterGame(m, "heart", diff)
		a.g.Status = game.StatusVictory
		a.g.Wave = 20
		for i := 0; i < 120; i++ {
			a.stepGame(.05)
		}
		if a.screen != ScreenGame || a.heartEndingPending || hiscore.LoadJournal().EndgameWins != 2 {
			t.Fatal("repeat Heart victory forced another film or counted wrong")
		}
	}
}

func TestEndingDoesNotPlayOnDefeatOrOtherFloor(t *testing.T) {
	for _, c := range []struct {
		level  string
		status game.GameStatus
	}{{"heart", game.StatusDefeat}, {"hub", game.StatusVictory}} {
		a := journalGameApp(t)
		m, err := game.LoadLevel("hub")
		if c.level == "heart" {
			m, err = game.LoadBoss()
		}
		if err != nil {
			t.Fatal(err)
		}
		a.enterGame(m, c.level, game.Normal)
		a.g.Status = c.status
		a.g.Wave = 20
		for i := 0; i < 120; i++ {
			a.stepGame(.05)
		}
		if a.screen != ScreenGame || a.heartEndingPending {
			t.Fatal("unearned Heart film played")
		}
		a.handle(Event{Rune: 'v'})
		if a.screen != ScreenGame {
			t.Fatal("unearned replay available")
		}
	}
}

func TestEarlyResultNavigationVisitsEndingAndRotundaReplay(t *testing.T) {
	a := journalGameApp(t)
	m, err := game.LoadBoss()
	if err != nil {
		t.Fatal(err)
	}
	a.enterGame(m, "heart", game.Normal)
	a.g.Status = game.StatusVictory
	a.g.Wave = 20
	a.stepGame(.05)
	a.handle(Event{Rune: 'r'})
	if a.screen != ScreenCutscene || a.g.Status != game.StatusVictory {
		t.Fatal("early restart bypassed the ending")
	}
	a.handle(Event{Key: KeyEscape})
	a.screen = ScreenOverworld
	a.ow.BootTTL = 0
	a.handle(Event{Rune: 'i'})
	if a.screen != ScreenCutscene || a.film.Film != render.FilmOpening {
		t.Fatal("Rotunda intro replay missing")
	}
	a.handle(Event{Key: KeyEscape})
	a.handle(Event{Rune: 'e'})
	if a.screen != ScreenCutscene || a.film.Film != render.FilmEnding {
		t.Fatal("held Heart ending replay missing")
	}
}

func TestCutsceneResizeHoldsCameraAndInput(t *testing.T) {
	a := owTestApp(t)
	a.term = &Terminal{size: [2]int{40, 12}}
	a.startCutscene(render.FilmOpening, ScreenOverworld, false)
	for i := 0; i < 100; i++ {
		a.tickCutscene()
	}
	a.handle(Event{Key: KeyEnter})
	if a.film.Frame != 0 || a.film.Revealed {
		t.Fatal("small terminal advanced film invisibly")
	}
	a.term.size = [2]int{80, 24}
	a.tickCutscene()
	if a.film.Frame != 1 {
		t.Fatal("film did not continue after resize")
	}
	a.term.size = [2]int{40, 12}
	a.handle(Event{Key: KeyEscape})
	if a.screen != ScreenOverworld {
		t.Fatal("resize notice trapped skip")
	}
}
