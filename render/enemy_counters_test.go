package render

import (
	"testing"

	"github.com/0xbenc/termtd/game"
)

func TestCounterTraitsRemainVisibleAndFrostRemovesEvasionTell(t *testing.T) {
	l := GameLayout(10, 6, 80, 24)
	g := &game.State{Time: 10}
	for _, tc := range []struct {
		kind game.EnemyKind
		slow bool
		bg   int
	}{
		{game.EnemyTank, false, 58},
		{game.EnemyShield, false, 24},
		{game.EnemyRunner, false, 53},
		{game.EnemyRunner, true, 0},
	} {
		f := &Frame{W: 80, H: 24, C: make([]Cell, 80*24)}
		e := &game.Enemy{Kind: tc.kind, HP: 100, MaxHP: 100, Pos: game.Pos{X: 3.5, Y: 2.5}}
		if tc.slow {
			e.SlowUntil = 11
		}
		drawEnemy(f, g, Palette(), l, e)
		c := f.C[l.FY(e.Pos.Y)*f.W+l.FX(e.Pos.X)]
		if c.BG != tc.bg || c.R != game.EnemySpecs[tc.kind].Short {
			t.Fatalf("kind %d slowed=%v: got %+v", tc.kind, tc.slow, c)
		}
		if tc.slow && c.FG != Palette().Frost {
			t.Fatal("frosted Rogue lost its cyan slow indication")
		}
	}
}
