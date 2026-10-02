package tui

import (
	"reflect"
	"testing"
)

func TestPaladinPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenNecromancerMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenPaladinMockup {
			t.Fatal("Necromancer did not lead to Paladin")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenRogueMockup {
			t.Fatal("Paladin did not lead to Rogue")
		}
		a.screen = ScreenPaladinMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Paladin portrait changed the lair")
		}
	}
}

func TestPaladinPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenPaladinMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Paladin portrait ignored quit")
		}
	}
}
