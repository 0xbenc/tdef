package main

import (
	"github.com/0xbenc/tdef/game"
	"testing"
)

func probe(t *testing.T, name, strategy string, diff game.Difficulty, hearts int, neutral bool) result {
	t.Helper()
	var m *game.Map
	var err error
	if name == "heart" {
		m, err = game.LoadBoss()
	} else {
		m, err = game.LoadLevel(name)
	}
	if err != nil {
		t.Fatal(err)
	}
	if neutral {
		m.Encounter = ""
		m.HPMul = .9
		m.GoldMul = 1.25
	}
	for _, p := range plans {
		if p.name == strategy {
			r := run(m, name, p, diff, hearts)
			if r.TimedOut {
				t.Fatalf("%s/%s stalled", name, strategy)
			}
			if strategy != "one-camp" && (r.FirstAttack <= 0 || r.FirstAttack > 8) {
				t.Fatalf("%s/%s opening has no timely action: %.1f", name, strategy, r.FirstAttack)
			}
			return r
		}
	}
	t.Fatal("unknown policy")
	return result{}
}

func TestShapesChangePlacementMetaWithIdenticalEncounters(t *testing.T) {
	rot := probe(t, "hub", "one-camp", game.Normal, 0, true)
	rift := probe(t, "canyon", "one-camp", game.Normal, 0, true)
	heart := probe(t, "heart", "one-camp", game.Normal, 0, true)
	if !rot.Won || rift.Won || heart.Won {
		t.Fatal("identical encounters should reward shared Rotunda coverage but require spreading on Rift/Heart")
	}
	halls := probe(t, "winding", "rangers", game.Normal, 0, true)
	garden := probe(t, "garden", "rangers", game.Normal, 0, true)
	if !halls.Won || halls.Leaks >= garden.Leaks {
		t.Fatal("Halls geometry must give long bows better exposure than Garden")
	}
}

func TestAuthoredRaidChangesViableCompositions(t *testing.T) {
	if !probe(t, "winding", "rangers", game.Normal, 0, false).Won {
		t.Fatal("Halls should support range")
	}
	if probe(t, "garden", "rangers", game.Normal, 0, false).Won {
		t.Fatal("Garden should challenge the same range-only policy")
	}
	for _, strategy := range []string{"cannons", "frost-chain", "frost-siege"} {
		if !probe(t, "garden", strategy, game.Normal, 0, false).Won {
			t.Fatalf("Garden must support %s as an alternative", strategy)
		}
	}
}

func TestHardCampaignSupportsAlternativeDefenses(t *testing.T) {
	cases := []struct {
		name       string
		hearts     int
		strategies []string
	}{
		{"hub", 0, []string{"rangers", "one-camp"}},
		{"canyon", 1, []string{"frost-siege", "ranged-siege"}},
		{"winding", 2, []string{"rangers", "ranged-siege"}},
		{"garden", 3, []string{"siege", "frost-siege"}},
		{"heart", 4, []string{"rangers", "ranged-siege", "frost-siege"}},
	}
	for _, c := range cases {
		for _, strategy := range c.strategies {
			t.Run(c.name+"/"+strategy, func(t *testing.T) {
				if !probe(t, c.name, strategy, game.Hard, c.hearts, false).Won {
					t.Fatal("intended defense cannot hold the hard campaign")
				}
			})
		}
	}
}
