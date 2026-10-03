package game

import (
	"reflect"
	"strings"
	"testing"
)

func authoredMap(t *testing.T, name string) *Map {
	t.Helper()
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
	return m
}

func TestAuthoredEncountersPreserveIntroductionsAndBosses(t *testing.T) {
	for _, name := range []string{"hub", "canyon", "winding", "garden", "heart"} {
		t.Run(name, func(t *testing.T) {
			m := authoredMap(t, name)
			if TacticalBrief(m) == "" {
				t.Fatal("missing tactical brief")
			}
			for wave := 1; wave <= MaxWaves; wave++ {
				base := BuildWave(wave)
				counts, original := [EnemyCount]int{}, [EnemyCount]int{}
				entries := BuildWaveFor(m, wave)
				for i, e := range entries {
					counts[e.Kind]++
					if i > 0 && e.At <= entries[i-1].At {
						t.Fatalf("wave %d non-increasing spawn times", wave)
					}
				}
				for _, e := range base {
					original[e.Kind]++
				}
				for k := EnemyKind(0); k < EnemyCount; k++ {
					if (counts[k] > 0) != (original[k] > 0) {
						t.Fatalf("wave %d changed enemy introductions: %v", wave, k)
					}
				}
				if counts[EnemyBoss] != original[EnemyBoss] {
					t.Fatalf("wave %d changed boss milestone", wave)
				}
				if !reflect.DeepEqual(base, BuildWave(wave)) {
					t.Fatal("mutated shared waves")
				}
			}
			for _, wave := range []int{0, MaxWaves + 1} {
				if BuildWaveFor(m, wave) != nil {
					t.Fatal("invalid wave spawned enemies")
				}
			}
		})
	}
}

func TestEncounterQueueAndPreviewAgree(t *testing.T) {
	for _, name := range []string{"canyon", "winding", "garden", "heart"} {
		m := authoredMap(t, name)
		s := NewState(m)
		s.Wave = 6
		s.StartWave()
		if !reflect.DeepEqual(s.SpawnQueue, BuildWaveFor(m, 7)) {
			t.Fatalf("%s state ignored encounter", name)
		}
		counts := [EnemyCount]int{}
		for _, e := range s.SpawnQueue {
			counts[e.Kind]++
		}
		if counts[EnemyRunner] == 0 || WavePreviewFor(m, 7) == "" || WaveTelegraphFor(m, 7) == "" {
			t.Fatalf("%s missing raid preview", name)
		}
	}
	garden, halls := authoredMap(t, "garden"), authoredMap(t, "winding")
	g, h := BuildWaveFor(garden, 7), BuildWaveFor(halls, 7)
	if len(g) <= len(h) || g[1].At-g[0].At >= h[1].At-h[0].At {
		t.Fatal("garden should have denser, larger raids than the halls")
	}
	if !strings.Contains(WaveTelegraphFor(garden, 7), "packed") {
		t.Fatal("garden raid warning is misleading")
	}
}

func TestDefaultWavesRemainUnchanged(t *testing.T) {
	for _, m := range []*Map{nil, {}, {Encounter: EncounterRotunda}} {
		for w := 1; w <= MaxWaves; w++ {
			if !reflect.DeepEqual(BuildWaveFor(m, w), BuildWave(w)) || WavePreviewFor(m, w) != WavePreview(w) {
				t.Fatalf("default wave %d changed", w)
			}
		}
	}
}

func TestLayoutsHaveDifferentExposure(t *testing.T) {
	halls, garden := authoredMap(t, "winding"), authoredMap(t, "garden")
	longest := func(m *Map) int {
		best, run := 1, 1
		for i := 2; i < len(m.Path); i++ {
			if m.Path[i].Sub(m.Path[i-1]) == m.Path[i-1].Sub(m.Path[i-2]) {
				run++
			} else {
				run = 1
			}
			if run > best {
				best = run
			}
		}
		return best
	}
	if longest(halls) < 35 || longest(garden) >= longest(halls)/2 {
		t.Fatal("halls must offer substantially longer straight firing lanes")
	}
	for _, name := range []string{"hub", "canyon", "winding", "garden", "heart"} {
		m := authoredMap(t, name)
		for y := 0; y < m.H; y++ {
			for x := 0; x < m.W; x++ {
				v := Vec{x, y}
				if m.At(v) != CellGrass {
					continue
				}
				reachable := false
				for _, sp := range m.Samples {
					if v.Center().Dist(sp) <= TowerSpecs[TowerSniper].Range[0] {
						reachable = true
						break
					}
				}
				if !reachable {
					t.Errorf("%s has useless build tile %v", name, v)
				}
			}
		}
	}
}
