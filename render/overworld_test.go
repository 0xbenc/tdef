package render

import (
	"testing"

	"tdef/game"
)

// The overworld must use the exact playfield scale logic the game does: the
// same grid size (45x13) drives ComputeScale/GameLayout, so a landmark lands
// where GameLayout says it will at every scale step.
func TestOverworldReusesPlayfieldScale(t *testing.T) {
	pal := Palette()
	st := NewOWState()
	sizes := []struct {
		w, h, wantScale int
	}{{62, 19, 1}, {92, 32, 2}, {137, 45, 3}, {182, 58, 4}}
	for _, s := range sizes {
		if got := ComputeScale(OWW, OWH, s.w, s.h); got != s.wantScale {
			t.Fatalf("ComputeScale(%d,%d,%d,%d) = %d, want %d", OWW, OWH, s.w, s.h, got, s.wantScale)
		}
		l := GameLayout(OWW, OWH, s.w, s.h)
		if l.Scale != s.wantScale {
			t.Fatalf("GameLayout scale = %d, want %d at %dx%d", l.Scale, s.wantScale, s.w, s.h)
		}
		f := RenderOverworld(s.w, s.h, st, 30, pal)
		// The Rotunda is open by default, so its pad centre carries the heart.
		cx, cy := l.center(22, 6)
		if f.C[cy*f.W+cx].R != '♥' {
			t.Fatalf("rotunda landmark missing at %dx%d scale %d: got %q", s.w, s.h, s.wantScale, f.C[cy*f.W+cx].R)
		}
	}
}

// A 45x13 overworld must fit the minimum 62x19 frame, exactly like a level.
func TestOverworldFitsMinFrame(t *testing.T) {
	mw, mh := MinFrame(OWW, OWH)
	if mw != FrameW || mh != ChromeTop+OWH+ChromeBot {
		t.Fatalf("MinFrame(%d,%d) = %dx%d, want %dx%d", OWW, OWH, mw, mh, FrameW, ChromeTop+OWH+ChromeBot)
	}
	f := RenderOverworld(62, 19, NewOWState(), 0, Palette())
	if f.W != 62 || f.H != 19 {
		t.Fatalf("frame is %dx%d, want 62x19", f.W, f.H)
	}
	// The full-window border is intact at the corners.
	for _, c := range []struct{ x, y, r rune }{{0, 0, '╭'}, {61, 0, '╮'}, {0, 18, '╰'}, {61, 18, '╯'}} {
		if got := f.C[c.y*62+c.x].R; got != c.r {
			t.Fatalf("corner (%d,%d) = %q, want %q", c.x, c.y, got, c.r)
		}
	}
}

// The renderer is a pure function of its inputs: identical inputs give an
// identical frame, so the animation is deterministic and testable.
func TestOverworldDeterministic(t *testing.T) {
	pal := Palette()
	st := NewOWState()
	a := RenderOverworld(92, 32, st, 123, pal).Text()
	b := RenderOverworld(92, 32, st, 123, pal).Text()
	if a != b {
		t.Fatal("RenderOverworld is not deterministic for identical inputs")
	}
}

// Every floor must be reachable from the Rift by walking corridors/pads — the
// "connected points" guarantee. A BFS over walkable cells from the Rift must
// land on every node's pad.
func TestOverworldAllNodesReachable(t *testing.T) {
	start := game.Vec{X: 6, Y: 10} // the Rift
	if !OWWalkable(start.X, start.Y) {
		t.Fatal("the Rift start cell is not walkable")
	}
	seen := map[game.Vec]bool{start: true}
	queue := []game.Vec{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range [4]game.Vec{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}} {
			n := game.Vec{X: cur.X + d.X, Y: cur.Y + d.Y}
			if seen[n] || !OWWalkable(n.X, n.Y) {
				continue
			}
			seen[n] = true
			queue = append(queue, n)
		}
	}
	for i := range owNodes {
		n := &owNodes[i]
		if !seen[game.Vec{X: n.X, Y: n.Y}] {
			t.Fatalf("floor %q at (%d,%d) is not reachable from the Rift", n.Name, n.X, n.Y)
		}
	}
}

// OWRects must return one rect per node, sized to the pad at the current
// scale, so mouse clicks cannot drift from the rendered pads.
func TestOWRectsMatchNodes(t *testing.T) {
	w, h := 92, 32
	l := GameLayout(OWW, OWH, w, h)
	rects := OWRects(w, h)
	if len(rects) != len(owNodes) {
		t.Fatalf("got %d rects, want %d", len(rects), len(owNodes))
	}
	for i, r := range rects {
		n := &owNodes[i]
		want := Rect{X: l.X(n.X - n.PW/2), Y: l.Y(n.Y - n.PH/2), W: n.PW * l.Scale, H: n.PH * l.Scale}
		if r != want {
			t.Fatalf("rect[%d] = %+v, want %+v", i, r, want)
		}
	}
}

// NewOWState must start Grak on the Rift with the Rift and Rotunda open.
func TestNewOWStateStart(t *testing.T) {
	st := NewOWState()
	floor, ok := OWFloorAt(st.Cursor.X, st.Cursor.Y)
	if !ok || floor.ID != "rift" {
		t.Fatalf("cursor starts at (%d,%d), not on the Rift", st.Cursor.X, st.Cursor.Y)
	}
	if !st.Unlocked["rift"] || !st.Unlocked["rotunda"] {
		t.Fatal("Rift and Rotunda must be unlocked at start")
	}
	if st.Unlocked["halls"] || st.Unlocked["garden"] || st.Unlocked["depths"] {
		t.Fatal("the far floors must be sealed at start")
	}
}
