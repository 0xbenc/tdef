package tui

import (
	"reflect"
	"testing"
)

func TestLightningPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenRangerMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenLightningMockup {
			t.Fatal("Ranger did not lead to Lightning Mage")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenTrebuchetMockup {
			t.Fatal("Lightning Mage did not lead to Trebuchet")
		}
		a.screen = ScreenLightningMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Lightning portrait changed the lair")
		}
	}
}

func TestLightningPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenLightningMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Lightning portrait ignored quit")
		}
	}
}
