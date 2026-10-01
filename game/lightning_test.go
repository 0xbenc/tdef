package game

import (
	"math"
	"testing"
)

// Exercise normal target acquisition and firing, with only the first hero
// inside the mage's range. Later heroes must be reached from the last hit.
func TestLightningBouncesBeyondMageRange(t *testing.T) {
	for level := 1; level <= 3; level++ {
		s := NewState(loadTestMap(t))
		tower := &Tower{ID: 100, Kind: TowerTesla, Level: level, Cell: Vec{X: 1, Y: 1}}
		s.Towers = []*Tower{tower}
		origin := tower.Pos()
		hops := tower.Spec().Chain[level-1]
		// Include one extra hero to verify that range does not remove the
		// per-shot bounce limit. All heroes survive their hit.
		for i := 0; i < hops+2; i++ {
			e := &Enemy{
				ID: i + 1, Kind: EnemyGrunt, HP: 1000, MaxHP: 1000,
				Pos:  Pos{X: origin.X + 3 + 2*float64(i), Y: origin.Y},
				Prog: float64(i + 1),
			}
			if i > 0 && e.Pos.Dist(origin) <= tower.Range() {
				t.Fatal("test requires bounce targets outside the mage's range")
			}
			s.Enemies = append(s.Enemies, e)
		}
		if target := s.acquireTarget(tower); target != s.Enemies[0] {
			t.Fatalf("level %d: acquired wrong first target", level)
		}
		s.fireTowers(0)
		if len(s.Beams) != 1 || len(s.Beams[0].From) != hops+2 {
			t.Fatalf("level %d: beam did not record the full chain", level)
		}
		for i, e := range s.Enemies {
			wantHP := 1000.0
			if i <= hops {
				wantHP -= tower.Dmg() * math.Pow(ChainFalloff, float64(i))
				if s.Beams[0].From[i+1] != e.Pos {
					t.Fatalf("level %d: beam missed hero %d", level, e.ID)
				}
			}
			if math.Abs(e.HP-wantHP) > 1e-9 {
				t.Fatalf("level %d hero %d: HP %f, want %f", level, e.ID, e.HP, wantHP)
			}
		}
	}
}

func TestLightningCannotBridgeGapBeyondBounceRange(t *testing.T) {
	s := NewState(loadTestMap(t))
	tower := &Tower{ID: 100, Kind: TowerTesla, Level: 3, Cell: Vec{X: 1, Y: 1}}
	s.Towers = []*Tower{tower}
	origin := tower.Pos()
	first := &Enemy{ID: 1, Kind: EnemyGrunt, HP: 1000, MaxHP: 1000, Pos: Pos{X: origin.X + 3, Y: origin.Y}}
	beyond := &Enemy{ID: 2, Kind: EnemyGrunt, HP: 1000, MaxHP: 1000, Pos: Pos{X: first.Pos.X + ChainRange + 0.1, Y: first.Pos.Y}}
	s.Enemies = []*Enemy{first, beyond}
	s.fireTowers(0)
	if first.HP == first.MaxHP || beyond.HP != beyond.MaxHP {
		t.Fatal("lightning must hit the first hero and stop at the oversized gap")
	}
	if len(s.Beams) != 1 || len(s.Beams[0].From) != 2 {
		t.Fatal("beam crossed a gap beyond the bounce range")
	}
}
