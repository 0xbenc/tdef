package main

import (
	"github.com/0xbenc/termtd/game"
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
	for _, strategy := range []string{"cannons", "siege", "frost-siege", "mixed-counters"} {
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
		{"hub", 0, []string{"frost-siege", "mixed-counters", "one-camp"}},
		{"canyon", 1, []string{"frost-siege", "ranged-siege"}},
		{"winding", 2, []string{"rangers", "ranged-siege"}},
		{"garden", 3, []string{"siege", "frost-siege"}},
		{"heart", 4, []string{"ranged-siege", "mixed-counters"}},
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

func TestSpecialistsHaveViableHardCampaignNiches(t *testing.T) {
	cases := []struct {
		floor, plan string
		hearts      int
		kind        game.TowerKind
	}{
		{"hub", "frost-forge", 0, game.TowerRuneforge},
		{"winding", "hook-siege", 2, game.TowerHookmaster},
		{"canyon", "minefield", 1, game.TowerSappers},
		{"winding", "hex-siege", 2, game.TowerWitch},
	}
	for _, c := range cases {
		t.Run(c.plan, func(t *testing.T) {
			r := probe(t, c.floor, c.plan, game.Hard, c.hearts, false)
			if !r.Won || r.Mix[c.kind] == 0 {
				t.Fatal("specialist composition cannot hold its intended hard floor")
			}
		})
	}
	r := probe(t, "heart", "specialists", game.Normal, 4, false)
	if !r.Won {
		t.Fatal("combined specialist defense cannot hold the Heart")
	}
	for k := game.TowerRuneforge; k < game.TowerCount; k++ {
		if r.Mix[k] == 0 {
			t.Fatal("combined defense omitted a specialist")
		}
	}
}

func TestGuidedDefensesRemainWinnable(t *testing.T) {
	for _, diff := range []game.Difficulty{game.Easy, game.Normal, game.Hard} {
		for stage := 1; stage <= 2; stage++ {
			name, hearts := "hub", 0
			wanted := []string{"mixed-counters", "frost-siege"}
			if stage == 2 {
				name, hearts = "canyon", 1
				wanted = []string{"mixed-counters", "minefield"}
			}
			m, err := game.LoadLevel(name)
			if err != nil {
				t.Fatal(err)
			}
			for _, strategy := range wanted {
				for _, p := range plans {
					if p.name != strategy {
						continue
					}
					r := runTraining(m, name, p, diff, hearts, stage)
					if !r.Won || r.TimedOut {
						t.Fatalf("stage %d %v %s cannot hold guided defense: %+v", stage, diff, strategy, r)
					}
				}
			}
		}
	}
}
