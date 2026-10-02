package tui

import (
	"reflect"
	"testing"

	"github.com/0xbenc/tdef/game"
)

func TestDragonMockupReturnsToUnchangedLair(t *testing.T) {
	for _, back := range []Event{{Key: KeyEscape}, {Key: KeyEnter}, {Rune: 'v'}} {
		a := owTestApp(t)
		a.ow.CameraSet, a.ow.CameraX = true, 17
		a.ow.Msg, a.ow.BonusGold = "waiting for Grak", 123
		before := a.ow
		a.handle(Event{Rune: 'v'})
		if a.screen != ScreenDragonMockup {
			t.Fatal("Rotunda did not open study")
		}
		a.handle(Event{Rune: 'w'})
		a.handle(back)
		if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
			t.Fatal("returning from study changed the lair")
		}
	}
}

func TestDragonMockupEntryGuards(t *testing.T) {
	for _, guard := range []string{"other floor", "arrival", "descent", "relic menu"} {
		a := owTestApp(t)
		switch guard {
		case "other floor":
			a.ow.Cursor = game.Vec{X: 65, Y: 9}
		case "arrival":
			a.ow.BootTTL = 1
		case "descent":
			a.ow.Descending = "rotunda"
		case "relic menu":
			a.ow.RelicMenu = true
		}
		a.handle(Event{Rune: 'v'})
		if a.screen != ScreenOverworld {
			t.Fatalf("opened study during %s", guard)
		}
	}
}

func TestDragonMockupQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenDragonMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("study did not honor quit")
		}
	}
}
