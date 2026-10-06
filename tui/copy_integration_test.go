package tui

import (
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/render"
)

func TestCampaignMessagesCarryPresentation(t *testing.T) {
	a := owTestApp(t)
	a.g = &game.State{Wave: 7}
	a.owFloorID = "rotunda"
	for _, won := range []bool{true, false} {
		a.owSetBanner(won, 300, true)
		want := render.OWMessageDefeat
		if won {
			want = render.OWMessageSuccess
		}
		if a.ow.ReturnKind != want {
			t.Fatalf("won %v: return kind %d, want %d", won, a.ow.ReturnKind, want)
		}
	}
	for _, floor := range []string{"rift", "depths"} {
		a.owSealedMessage(floor)
		if a.ow.MsgKind != render.OWMessageLocked {
			t.Fatalf("%s should carry locked presentation", floor)
		}
	}
	a.ow.Unsealing["rift"] = 30
	a.owSealedMessage("rift")
	if a.ow.MsgKind != render.OWMessageNeutral {
		t.Fatal("opening a door must clear locked presentation")
	}
}
