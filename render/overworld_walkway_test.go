package render

import (
	"strings"
	"testing"

	"github.com/0xbenc/tdef/game"
)

// Every step in a route must be an orthogonal step that Grak can actually
// take, including the bridge's turns and the upper/lower approaches.
func TestOWRoutesFollowWalkableCells(t *testing.T) {
	for i, cells := range owRouteCells {
		for j, v := range cells {
			if !OWWalkable(v.X, v.Y) {
				t.Fatalf("%s route contains impassable cell %v", owRouteFloors[i], v)
			}
			if j > 0 && cells[j-1].Man(v) != 1 {
				t.Fatalf("%s route jumps from %v to %v", owRouteFloors[i], cells[j-1], v)
			}
		}
	}
	// The crossing forks into two separate approaches, with real turns.
	if !owCorridor[game.Vec{X: 45, Y: 5}] || !owCorridor[game.Vec{X: 45, Y: 7}] {
		t.Fatal("the bridge must lead to both the upper and lower approaches")
	}
	if owCorridor[game.Vec{X: 50, Y: 6}] {
		t.Fatal("the old straight route still bypasses the fork")
	}
}

// Scenery and labels must never disguise a walkable route as empty ground.
// Test the final composed scene at all four scales, not just the paving pass.
func TestOWWalkwayRemainsVisible(t *testing.T) {
	for scale := 1; scale <= 4; scale++ {
		w, h := OWW*scale+2, OWH*scale+ChromeTop+ChromeBot
		st := NewOWState()
		l := OverworldLayout(w, h, st)
		f := RenderOverworld(w, h, st, 100, Palette()) // procession is between passes
		for v := range owCorridor {
			if _, room := OWFloorAt(v.X, v.Y); room {
				continue // room interiors have their own entrance illustrations
			}
			x, y := l.center(v.X, v.Y)
			c := f.C[y*f.W+x]
			if !strings.ContainsRune("═║╔╗╚╝╠╣╦╩╬", c.R) || c.BG != 237 {
				t.Fatalf("scale %d: route %v hidden by %q bg=%d", scale, v, c.R, c.BG)
			}
		}
	}
}
