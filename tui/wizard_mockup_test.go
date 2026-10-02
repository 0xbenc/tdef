package tui

import (
	"reflect"
	"testing"
)

func TestWizardPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenMercenaryMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenWizardMockup {
			t.Fatal("Mercenary did not lead to Wizard")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenCenturionMockup {
			t.Fatal("Wizard did not lead to Centurion")
		}
		a.screen = ScreenWizardMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Wizard portrait changed the lair")
		}
	}
}

func TestWizardPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenWizardMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Wizard portrait ignored quit")
		}
	}
}
