package tui

import (
	"reflect"
	"testing"
)

func TestGrakPortraitNavigationPreservesLair(t *testing.T) {
	a := owTestApp(t)
	a.ow.Msg, a.ow.BonusGold = "waiting", 123
	before := a.ow
	a.handle(Event{Rune: 'g'})
	if a.screen != ScreenGrakMockup {
		t.Fatal("g did not open Grak in the Rotunda")
	}
	for _, screen := range []Screen{ScreenGnollMockup, ScreenGunnerMockup, ScreenFrostMockup, ScreenPlayerMockup, ScreenCannonierMockup, ScreenRangerMockup, ScreenLightningMockup, ScreenTrebuchetMockup, ScreenNecromancerMockup, ScreenPaladinMockup, ScreenRogueMockup, ScreenMercenaryMockup, ScreenWizardMockup, ScreenCenturionMockup, ScreenSquireMockup, ScreenDragonMockup, ScreenGrakMockup} {
		a.handle(Event{Key: KeyTab})
		if a.screen != screen {
			t.Fatalf("tab: got %v, want %v", a.screen, screen)
		}
	}
	a.handle(Event{Rune: 'w'})
	a.handle(Event{Key: KeyEnter})
	if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
		t.Fatal("portraits changed the lair")
	}
}

func TestGrakPortraitQuit(t *testing.T) {
	for _, e := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
		a := owTestApp(t)
		a.screen = ScreenGrakMockup
		a.handle(e)
		if !a.quitting {
			t.Fatal("Grak portrait ignored quit")
		}
	}
}
