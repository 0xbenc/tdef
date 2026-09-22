package render

import (
	"strings"
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

// owNodeStatus must read the lair's memory: sealed when locked, open when
// unlocked, current when Grak stands on it, and open mid-unseal.
func TestOWNodeStatus(t *testing.T) {
	st := NewOWState()
	if got := owNodeStatus(st, "halls"); got != OWSealed {
		t.Errorf("halls status = %v, want sealed", got)
	}
	if got := owNodeStatus(st, "rift"); got != OWCurrent {
		t.Errorf("rift status (Grak is there) = %v, want current", got)
	}
	if got := owNodeStatus(st, "rotunda"); got != OWOpen {
		t.Errorf("rotunda status = %v, want open", got)
	}
	st.Cursor = game.Vec{X: 35, Y: 3} // stand on the (sealed) halls
	if got := owNodeStatus(st, "halls"); got != OWSealed {
		t.Errorf("standing on a sealed floor = %v, want still sealed", got)
	}
	st.Unlocked["halls"] = true
	if got := owNodeStatus(st, "halls"); got != OWCurrent {
		t.Errorf("unlocked+standing = %v, want current", got)
	}
	st.Unsealing["garden"] = 20
	if got := owNodeStatus(st, "garden"); got != OWOpen {
		t.Errorf("mid-unseal garden = %v, want open", got)
	}
}

// A sealed floor wears a blinking lock on its top frame; a revealed lair has
// no locks at all.
func TestOWSealedLock(t *testing.T) {
	pal := Palette()
	w, h := 62, 19
	l := GameLayout(OWW, OWH, w, h)
	st := NewOWState()
	f := RenderOverworld(w, h, st, 0, pal) // frame 0: the lock is lit
	n := owNodeByID("halls")
	lx, ly := l.center(n.X, n.Y-n.PH/2)
	if got := f.C[ly*f.W+lx].R; got != '╳' {
		t.Fatalf("sealed halls lock at (%d,%d) = %q, want ╳", lx, ly, got)
	}
	// RevealAll unseals the presentation: no lock remains.
	st.RevealAll = true
	f2 := RenderOverworld(w, h, st, 0, pal)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if f2.C[y*f2.W+x].R == '╳' {
				t.Fatalf("revealed lair still shows a lock at (%d,%d)", x, y)
			}
		}
	}
}

// When the heart has unsealed, the Rotunda's landmark becomes the portal
// (◉) instead of the heart (♥).
func TestOWBossDoor(t *testing.T) {
	pal := Palette()
	w, h := 62, 19
	l := GameLayout(OWW, OWH, w, h)
	st := NewOWState()
	cx, cy := l.center(22, 6)
	if got := RenderOverworld(w, h, st, 0, pal).C[cy*w+cx].R; got != '♥' {
		t.Fatalf("unopened rotunda landmark = %q, want ♥", got)
	}
	st.BossReady = true
	if got := RenderOverworld(w, h, st, 0, pal).C[cy*w+cx].R; got != '◉' {
		t.Fatalf("unsealed rotunda landmark = %q, want ◉ (the portal)", got)
	}
}

// A held lair changes its title and calms the circulation.
func TestOWBossDoneTitle(t *testing.T) {
	pal := Palette()
	w, h := 62, 19
	st := NewOWState()
	if !strings.Contains(RenderOverworld(w, h, st, 0, pal).Text(), "THE LAIR") {
		t.Fatal("title missing")
	}
	st.BossDone = true
	if !strings.Contains(RenderOverworld(w, h, st, 0, pal).Text(), "THE LAIR — HELD") {
		t.Fatal("held title missing")
	}
}

// A memory-rich lair (records, hearts, an unsealed heart, scars) must still
// render deterministically.
func TestOverworldRichStateDeterministic(t *testing.T) {
	pal := Palette()
	st := NewOWState()
	st.Unlocked["halls"] = true
	st.Records["rift"] = OWRec{Cleared: true, BestWave: 20, LastWave: 20, LastWon: true}
	st.Records["halls"] = OWRec{Cleared: false, BestWave: 9, LastWave: 9, LastWon: false}
	st.Hearts = 2
	st.BossReady = true
	st.Tokens = 4
	a := RenderOverworld(92, 32, st, 55, pal).Text()
	b := RenderOverworld(92, 32, st, 55, pal).Text()
	if a != b {
		t.Fatal("rich-state RenderOverworld is not deterministic")
	}
}

// OWFloorOf / OWFloorName expose each floor's identity; the boss door reads
// "the Heart".
func TestOWFloorIdentity(t *testing.T) {
	fl, ok := OWFloorOf("garden")
	if !ok || fl.Level != "garden" || fl.Name != "the Sunken Garden" {
		t.Fatalf("OWFloorOf(garden) = %+v %v", fl, ok)
	}
	if _, ok := OWFloorOf("nope"); ok {
		t.Error("OWFloorOf(unknown) reported a floor")
	}
	if got := OWFloorName(HeartFloorID); got != "the Heart" {
		t.Errorf("OWFloorName(heart) = %q, want the Heart", got)
	}
	if got := OWFloorName("rift"); got != "the Rift" {
		t.Errorf("OWFloorName(rift) = %q, want the Rift", got)
	}
}
