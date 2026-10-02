package tui

import (
	"reflect"
	"testing"
)

func TestMercenaryPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenRogueMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenMercenaryMockup {
			t.Fatal("Rogue did not lead to Mercenary")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenWizardMockup {
			t.Fatal("Mercenary did not lead to Wizard")
		}
		a.screen = ScreenMercenaryMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Mercenary portrait changed the lair")
		}
	}
}

func TestMercenaryPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenMercenaryMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Mercenary portrait ignored quit")
		}
	}
}
