package tui

import (
	"github.com/0xbenc/termtd/render"
	"testing"
)

func TestHelpPausesDefenseAndRestoresState(t *testing.T) {
	for _, paused := range []bool{false, true} {
		a := journalGameApp(t)
		a.screen = ScreenGame
		a.ui.Paused = paused
		gold, wave := a.g.Gold, a.g.Wave
		a.handle(Event{Rune: 'h'})
		if a.screen != ScreenHelp || !a.ui.Paused {
			t.Fatal("help did not pause defense")
		}
		a.handle(Event{Rune: '3'})
		a.handle(Event{Key: KeyDown})
		a.handle(Event{Rune: ']'})
		if a.help.Page != 3 || a.help.Scroll != 0 {
			t.Fatal("topic navigation did not reset scroll")
		}
		a.handle(Event{Key: KeyEscape})
		if a.screen != ScreenGame || a.ui.Paused != paused || a.g.Gold != gold || a.g.Wave != wave {
			t.Fatal("help changed defense state")
		}
	}
}

func TestHelpScrollAndMouseStayInGuide(t *testing.T) {
	a := &App{screen: ScreenMenu}
	a.openHelp()
	a.handle(Event{Key: KeyDown})
	if a.help.Scroll != 1 {
		t.Fatal("help did not scroll")
	}
	a.handle(Event{Mouse: true, Press: true, Btn: 65})
	if a.screen != ScreenHelp || a.help.Scroll != 4 {
		t.Fatal("wheel dismissed help")
	}
	w, _ := a.termSize()
	for x := 0; x < w; x++ {
		if render.HelpTopicAt(w, 24, x, 4) == 2 {
			a.handle(Event{Mouse: true, Press: true, Btn: 0, X: x, Y: 4})
			break
		}
	}
	if a.help.Page != 2 || a.help.Scroll != 0 {
		t.Fatal("mouse did not select topic")
	}
	a.handle(Event{Key: KeyEscape})
	if a.screen != ScreenMenu {
		t.Fatal("help did not return to menu")
	}
}
