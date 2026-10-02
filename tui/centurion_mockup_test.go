package tui

import (
	"reflect"
	"testing"
)

func TestCenturionPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenWizardMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenCenturionMockup {
			t.Fatal("Wizard did not lead to Centurion")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenSquireMockup {
			t.Fatal("Centurion did not lead to Squire")
		}
		a.screen = ScreenCenturionMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Centurion portrait changed the lair")
		}
	}
}

func TestCenturionPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenCenturionMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Centurion portrait ignored quit")
		}
	}
}
