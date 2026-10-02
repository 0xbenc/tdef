package tui

import (
	"reflect"
	"testing"
)

func TestTrebuchetPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenLightningMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenTrebuchetMockup {
			t.Fatal("Lightning Mage did not lead to Trebuchet")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenNecromancerMockup {
			t.Fatal("Trebuchet did not lead to Necromancer")
		}
		a.screen = ScreenTrebuchetMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Trebuchet portrait changed the lair")
		}
	}
}

func TestTrebuchetPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenTrebuchetMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Trebuchet portrait ignored quit")
		}
	}
}
