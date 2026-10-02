package tui

import (
	"reflect"
	"testing"
)

func TestRoguePortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenPaladinMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenRogueMockup {
			t.Fatal("Paladin did not lead to Rogue")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenMercenaryMockup {
			t.Fatal("Rogue did not lead to Mercenary")
		}
		a.screen = ScreenRogueMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Rogue portrait changed the lair")
		}
	}
}

func TestRoguePortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenRogueMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Rogue portrait ignored quit")
		}
	}
}
