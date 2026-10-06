package render

import (
	"testing"

	"github.com/0xbenc/termtd/game"
)

func TestEditableNamesMatchGameplayAndJournal(t *testing.T) {
	for i, entry := range TowerJournalEntries {
		if entry.Kind != game.TowerKind(i) || entry.Name != game.TowerSpecs[i].Name {
			t.Fatalf("defender %s no longer matches its game entity", entry.ID)
		}
	}
	for i, entry := range EnemyJournalEntries {
		if entry.Name != game.EnemySpecs[i].Name {
			t.Fatalf("enemy %s no longer matches its game entity", entry.ID)
		}
	}
	// Copy changes must not become changes to saved discovery IDs.
	for i, floor := range JournalFloorIDs {
		if PlaceJournalEntries[i].ID != floor {
			t.Fatalf("place %d changed its persistent ID", i)
		}
	}
	for i, wins := range []string{"1", "3", "5", "10"} {
		if got := StoryJournalEntries[18+i].ID; got != "afterward:"+wins {
			t.Fatalf("Afterward ID = %q", got)
		}
	}
}

func TestOverworldMessageColorsIgnoreWording(t *testing.T) {
	for _, tc := range []struct {
		kind OWMessageKind
		want int
	}{
		{OWMessageNeutral, 244},
		{OWMessageSuccess, 114},
		{OWMessageDefeat, 167},
		{OWMessageLocked, 174},
	} {
		for _, wording := range []string{"All quiet.", "held endures broke sealed breaks open"} {
			for _, banner := range []bool{false, true} {
				st := NewOWState()
				st.BootTTL = 0
				if banner {
					st.ReturnTTL, st.ReturnMsg, st.ReturnKind = 180, wording, tc.kind
				} else {
					st.Msg, st.MsgKind = wording, tc.kind
				}
				f := &Frame{W: 100, H: 24, C: make([]Cell, 100*24)}
				drawOWVoice(f, 100, st, 0)
				if got := f.C[102].FG; got != tc.want {
					t.Fatalf("kind %d, wording %q, banner %v: color %d, want %d", tc.kind, wording, banner, got, tc.want)
				}
			}
		}
	}
}
