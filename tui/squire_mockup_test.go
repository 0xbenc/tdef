package tui

import (
	"reflect"
	"testing"
)

func TestSquirePortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenCenturionMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenSquireMockup {
			t.Fatal("Centurion did not lead to Squire")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenDragonMockup {
			t.Fatal("Squire did not lead to Malgrath")
		}
		a.screen = ScreenSquireMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Squire portrait changed the lair")
		}
	}
}

func TestSquirePortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenSquireMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Squire portrait ignored quit")
		}
	}
}
