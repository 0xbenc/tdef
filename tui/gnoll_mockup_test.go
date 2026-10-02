package tui

import (
	"reflect"
	"testing"
)

func TestGnollPortraitReturnsWithoutChangingLair(t *testing.T) {
	for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
		a := owTestApp(t)
		a.ow.Msg, a.ow.BonusGold = "waiting", 123
		before := a.ow
		a.screen = ScreenGnollMockup
		a.handle(Event{Rune: 'w'})
		a.handle(exit)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("Gnoll portrait changed the lair")
		}
	}
}

func TestGnollPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenGnollMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Gnoll portrait ignored quit")
		}
	}
}
