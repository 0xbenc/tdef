package tui

import (
	"reflect"
	"testing"
)

func TestFrostAndPlayerPortraitNavigation(t *testing.T) {
	for _, screen := range []Screen{ScreenFrostMockup, ScreenPlayerMockup} {
		for _, exit := range []Event{{Key: KeyEscape}, {Key: KeyEnter}} {
			a := owTestApp(t)
			a.ow.Msg, a.ow.BonusGold = "waiting", 123
			before := a.ow
			a.screen = screen
			a.handle(Event{Rune: 'w'})
			a.handle(exit)
			if a.screen != ScreenOverworld || !reflect.DeepEqual(before, a.ow) {
				t.Fatal("lore portrait changed the lair")
			}
		}
		for _, quit := range []Event{{Rune: 'q'}, {Key: KeyCtrlC}} {
			a := owTestApp(t)
			a.screen = screen
			a.handle(quit)
			if !a.quitting {
				t.Fatal("lore portrait ignored quit")
			}
		}
	}
	a := owTestApp(t)
	a.screen = ScreenGunnerMockup
	for _, want := range []Screen{ScreenFrostMockup, ScreenPlayerMockup, ScreenCannonierMockup, ScreenRangerMockup, ScreenLightningMockup, ScreenTrebuchetMockup, ScreenNecromancerMockup, ScreenPaladinMockup, ScreenRogueMockup, ScreenMercenaryMockup, ScreenWizardMockup, ScreenCenturionMockup, ScreenSquireMockup, ScreenDragonMockup} {
		a.handle(Event{Key: KeyTab})
		if a.screen != want {
			t.Fatalf("tab: got %v, want %v", a.screen, want)
		}
	}
}
