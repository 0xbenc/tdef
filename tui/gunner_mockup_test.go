package tui

import (
	"reflect"
	"testing"
)

func TestGunnerPortraitNavigation(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenGnollMockup
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenGunnerMockup {
			t.Fatal("tab did not open gunner after gnolls")
		}
		a.handle(Event{Key: KeyTab})
		if a.screen != ScreenFrostMockup {
			t.Fatal("tab did not open the Frost Mage")
		}
		a.screen = ScreenGunnerMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Gunner portrait changed the lair")
		}
	}
}

func TestGunnerPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenGunnerMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Gunner portrait ignored quit")
		}
	}
}
