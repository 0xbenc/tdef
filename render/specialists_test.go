package render

import (
	"github.com/0xbenc/termtd/game"
	"testing"
)

func TestSpecialistRosterFitsSmallTerminal(t *testing.T) {
	for _, w := range []int{62, 80, 120} {
		slots := TowerSlots(w, 24, 1)
		if len(slots) != 4 {
			t.Fatal("missing specialists")
		}
		for i, slot := range slots {
			if slot.Key != i+1 || slot.X+slot.W > w-1 {
				t.Fatalf("invalid slot at %d: %+v", w, slot)
			}
			for j := 0; j < i; j++ {
				prev := slots[j]
				if prev.Y == slot.Y && prev.X+prev.W >= slot.X {
					t.Fatal("overlapping specialist labels")
				}
			}
		}
	}
}

func BenchmarkSpecialistBattlefield(b *testing.B) {
	m, _ := game.LoadLevel("winding")
	s := game.NewState(m)
	s.Gold = 100000
	for y := 1; y < m.H-1; y++ {
		for x := 1; x < m.W-1; x++ {
			if len(s.Towers) >= 24 {
				break
			}
			kind := game.TowerGunner
			if len(s.Towers) < 4 {
				kind = game.TowerKind(int(game.TowerRuneforge) + len(s.Towers))
			}
			s.Build(game.Vec{X: x, Y: y}, kind)
		}
	}
	for i := 0; i < 150; i++ {
		progress := float64(i % int(m.TotalLen))
		s.Enemies = append(s.Enemies, &game.Enemy{ID: i + 1000, Kind: game.EnemyGrunt, HP: 10000, MaxHP: 10000, Prog: progress, Pos: m.PointAt(progress), Speed: 1, SlowFactor: 1})
	}
	ui := UI{Selected: -1, Level: "winding"}
	pal := Palette()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Step(1. / 30)
		Render(s, &ui, pal, 120, 40, i)
	}
}
