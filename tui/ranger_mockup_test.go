package tui

import (
	"reflect"
	"testing"
)

func TestRangerPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenCannonierMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenRangerMockup {
			t.Fatal("Cannonier did not lead to Ranger")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenLightningMockup {
			t.Fatal("Ranger did not lead to Lightning Mage")
		}
		a.screen = ScreenRangerMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Ranger portrait changed the lair")
		}
	}
}

func TestRangerPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenRangerMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Ranger portrait ignored quit")
		}
	}
}
