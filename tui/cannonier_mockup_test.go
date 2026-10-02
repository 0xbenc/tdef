package tui

import (
	"reflect"
	"testing"
)

func TestCannonierPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenPlayerMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenCannonierMockup {
			t.Fatal("Player did not lead to Cannonier")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenRangerMockup {
			t.Fatal("Cannonier did not lead to Ranger")
		}
		a.screen = ScreenCannonierMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Cannonier portrait changed the lair")
		}
	}
}

func TestCannonierPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenCannonierMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Cannonier portrait ignored quit")
		}
	}
}
