package tui

import (
	"reflect"
	"testing"
)

func TestNecromancerPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenTrebuchetMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenNecromancerMockup {
			t.Fatal("Trebuchet did not lead to Necromancer")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenPaladinMockup {
			t.Fatal("Necromancer did not lead to Paladin")
		}
		a.screen = ScreenNecromancerMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Necromancer portrait changed the lair")
		}
	}
}

func TestNecromancerPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenNecromancerMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Necromancer portrait ignored quit")
		}
	}
}
