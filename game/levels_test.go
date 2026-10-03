package game

import "testing"

// A touching hairpin can silently become a shortcut in computePath's BFS.
// Every painted road cell must belong to one unbranched spawn-to-exit route.
func TestBuiltInRoutesAndBuildSpace(t *testing.T) {
	for _, name := range append(LevelNames(), "heart") {
		t.Run(name, func(t *testing.T) {
			var m *Map
			var err error
			if name == "heart" {
				m, err = LoadBoss()
			} else {
				m, err = LoadLevel(name)
			}
			if err != nil {
				t.Fatal(err)
			}
			roads, useful := 0, 0
			for y := 0; y < m.H; y++ {
				for x := 0; x < m.W; x++ {
					v := Vec{x, y}
					if m.At(v) == CellPath {
						roads++
						neighbors := 0
						for _, d := range [4]Vec{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
							if m.At(Vec{x + d.X, y + d.Y}) == CellPath {
								neighbors++
							}
						}
						want := 2
						if v == m.Spawn || v == m.Exit {
							want = 1
						}
						if neighbors != want {
							t.Errorf("road at %v has %d neighbors, want %d", v, neighbors, want)
						}
						continue
					}
					if m.At(v) == CellGrass {
						// Count meaningful placements even for the shortest-range
						// tower, rather than inflating the count with remote grass.
						for _, p := range m.Path {
							if v.Center().Dist(p.Center()) <= TowerSpecs[TowerFlak].Range[0] {
								useful++
								break
							}
						}
					}
				}
			}
			if roads != len(m.Path)+1 { // Path excludes the spawn cell.
				t.Errorf("%d painted roads but %d routed cells: shortcut or orphan road", roads, len(m.Path)+1)
			}
			minUseful := 100
			if name == "canyon" {
				minUseful = 80
			} // Narrow banks deliberately limit one camp.
			if useful < minUseful {
				t.Errorf("only %d useful short-range placements, want at least %d", useful, minUseful)
			}
			turns := 0
			for i := 2; i < len(m.Path); i++ {
				if m.Path[i].Sub(m.Path[i-1]) != m.Path[i-1].Sub(m.Path[i-2]) {
					turns++
				}
			}
			minTurns := 10
			if name == "winding" {
				minTurns = 4
			} // Long firing lanes are intentional here.
			if turns < minTurns {
				t.Errorf("only %d bends, want at least %d", turns, minTurns)
			}
		})
	}
}

func TestFloorEconomyPaysForSpawnedAndSplitEnemies(t *testing.T) {
	m := loadTestMap(t)
	m.GoldMul = 1.5
	s := NewState(m)
	if s.Gold != 330 {
		t.Fatalf("starting gold = %d, want 330", s.Gold)
	}
	s.Wave, s.WaveActive = 12, true
	s.SpawnQueue = []SpawnEntry{{Kind: EnemySplitter}}
	s.Step(0)
	if len(s.Enemies) != 1 {
		t.Fatalf("spawned %d enemies, want one necromancer", len(s.Enemies))
	}
	s.applyDamage(s.Enemies[0], 1e6, TowerSniper)
	if len(s.Enemies) != 3 {
		t.Fatalf("necromancer did not spawn two children")
	}
	for _, e := range s.Enemies {
		if !e.Dead {
			s.applyDamage(e, 1e6, TowerSniper)
		}
	}
	s.Step(0)
	// 330 starting + 18 necromancer + 2*3 children + 135 wave-clear.
	if s.Gold != 489 || s.TotalKills != 3 || s.WaveActive {
		t.Errorf("after clearing: gold=%d kills=%d active=%v", s.Gold, s.TotalKills, s.WaveActive)
	}
}

func TestFloorEconomyPreservesPricesAndScalesEarlyBonus(t *testing.T) {
	m := loadTestMap(t)
	m.GoldMul = 1.5
	s := NewStateDiff(m, Hard)
	if s.Gold != 300 || s.Lives != Hard.Lives() {
		t.Fatalf("hard start: gold=%d lives=%d", s.Gold, s.Lives)
	}
	tower := s.Build(Vec{4, 1}, TowerGunner)
	if tower == nil || s.Gold != 250 {
		t.Fatalf("gunner should still cost 50: gold=%d", s.Gold)
	}
	if !s.Upgrade(tower) || s.Gold != 205 {
		t.Fatalf("upgrade should still cost 45: gold=%d", s.Gold)
	}
	s.Wave, s.NextWaveAt = 1, 10
	if bonus := s.StartWave(); bonus != 12 || s.Gold != 217 {
		t.Errorf("early wave: bonus=%d gold=%d, want 12 and 217", bonus, s.Gold)
	}
	if refund := s.Sell(tower); refund != 66 {
		t.Errorf("refund = %d, want 66 (70%% of 95)", refund)
	}
}

func TestUnconfiguredFloorEconomy(t *testing.T) {
	for _, mul := range []float64{0, 1} {
		m := loadTestMap(t)
		m.GoldMul = mul
		for _, diff := range []Difficulty{Easy, Normal, Hard} {
			if got := NewStateDiff(m, diff).Gold; got != diff.Gold() {
				t.Errorf("multiplier %v difficulty %v: starting gold=%d", mul, diff, got)
			}
		}
		if got := m.ScaleGold(5); got != 5 {
			t.Errorf("multiplier %v: bounty=%d, want 5", mul, got)
		}
	}
}
