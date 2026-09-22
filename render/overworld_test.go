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

// owStateAllHeld is a fully-held lair for render tests: every built-in floor
// cleared, all corridors open, the heart unsealed.
func owStateAllHeld() OWState {
	st := NewOWState()
	for _, id := range []string{"rift", "halls", "garden", "rotunda"} {
		st.Records[id] = OWRec{Cleared: true, BestWave: 20, LastWave: 20, LastWon: true}
		st.Scores[id] = 9000
	}
	st.Unlocked["halls"] = true
	st.Unlocked["garden"] = true
	st.Unlocked["depths"] = true
	st.Hearts = 4
	st.BossReady = true
	st.FirstRun = false
	return st
}

// Pad edge glows must not leave the playfield: the bottom pads (Rift, Garden)
// touch grid row 12, so their "glow just outside the frame" used to land on
// the ledger row (h-4), and the Depths touches row 0, its glow the voice row.
func TestOWGlowClamped(t *testing.T) {
	pal := Palette()
	w, h := 182, 58
	l := GameLayout(OWW, OWH, w, h)
	f := RenderOverworld(w, h, owStateAllHeld(), 0, pal)
	// The Rift's bottom glow point: its centre column, the first row below
	// the map (the ledger row at this size).
	rx, ry := l.X(6)+l.Scale/2, l.Oy+OWH*l.Scale
	if c := f.C[ry*f.W+rx]; c.R != ' ' {
		t.Fatalf("rift bottom glow leaked to (%d,%d) = %q, want blank", rx, ry, c.R)
	}
	// The Depths' top glow point: its centre column, the first row above the
	// map (the voice row; blank there in this state).
	dx, dy := l.X(9)+l.Scale/2, l.Oy-1
	if c := f.C[dy*f.W+dx]; c.R != ' ' {
		t.Fatalf("depths top glow leaked to (%d,%d) = %q, want blank", dx, dy, c.R)
	}
}

// Sealed rooms must stay legible: dim() floors at 24, so the Sunken Garden's
// chrome reads as dark green against the void, not black.
func TestOWSealedDim(t *testing.T) {
	pal := Palette()
	w, h := 62, 19
	l := GameLayout(OWW, OWH, w, h)
	f := RenderOverworld(w, h, NewOWState(), 0, pal)
	n := owNodeByID("garden")
	c := f.C[l.Y(n.Y-n.PH/2)*f.W+l.X(n.X-n.PW/2)]
	if c.R != '╔' || c.FG != dim(45) {
		t.Fatalf("sealed garden corner = %q/%d, want ╔/%d", c.R, c.FG, dim(45))
	}
	if dim(45) < 24 {
		t.Fatalf("dim(45) = %d, want >= 24 (sealed rooms must stay legible)", dim(45))
	}
}

// The arrival cinematic: the heart's room is the ignition source (lit from
// the first frame), the unlit void is masked until the light front sweeps
// out across the map, and the last frame is the steady state (seam).
func TestOWBootPhases(t *testing.T) {
	pal := Palette()
	w, h := 137, 45
	l := GameLayout(OWW, OWH, w, h)
	st := NewOWState()

	// First frame of the waking: the far corner is unlit dark; the heart,
	// the ignition source, is already lit.
	st.BootTTL = OWBootFrames - 1
	f := RenderOverworld(w, h, st, 0, pal)
	cx, cy := l.center(44, 0)
	if c := f.C[cy*f.W+cx]; c.R != ' ' || c.BG != 233 {
		t.Fatalf("boot frame 1: far corner = %q bg %d, want unlit (space/233)", c.R, c.BG)
	}
	hx, hy := l.center(22, 6)
	if c := f.C[hy*f.W+hx]; c.R != '♥' {
		t.Fatalf("boot frame 1: heart = %q, want ♥ (the ignition source)", c.R)
	}

	// Mid-sweep (frame 65): the light has crossed the Rift's pad — a corner
	// of the pad is lit (its interior, not the unlit void's 233). Grak starts
	// on the Rift's centre landmark, so the landmark itself is occluded until
	// he walks off it (the C4 lantern addresses that).
	st.BootTTL = OWBootFrames - 65
	f2 := RenderOverworld(w, h, st, 0, pal)
	rx, ry := l.center(5, 8)
	if c := f2.C[ry*f2.W+rx]; c.BG == 233 {
		t.Fatalf("boot frame 65: rift pad corner still unlit (bg 233)")
	}
	// The front itself is a white dot at the cell 18 units from the heart.
	fx, fy := l.center(4, 6)
	if c := f2.C[fy*f2.W+fx]; c.R != '·' || c.FG != 255 {
		t.Fatalf("boot frame 65: front at (4,6) = %q/%d, want ·/255", c.R, c.FG)
	}
}

// The boot's last frame must be the steady state: at BootTTL = 1 nothing is
// masked, no front is drawn, the chrome is up and the normal voice speaks.
func TestOWBootSeam(t *testing.T) {
	pal := Palette()
	for _, s := range []struct{ w, h int }{{62, 19}, {137, 45}} {
		st := NewOWState()
		st.BootTTL = 1
		a := RenderOverworld(s.w, s.h, st, 0, pal).Text()
		st.BootTTL = 0
		b := RenderOverworld(s.w, s.h, st, 0, pal).Text()
		if a != b {
			t.Fatalf("boot seam broken at %dx%d (BootTTL 1 != 0)", s.w, s.h)
		}
	}
}

// The walk trail fades: the newest step is brightest, the oldest dimmest,
// and a zero-age step is gone.
func TestOWTrailFade(t *testing.T) {
	pal := Palette()
	w, h := 92, 32
	l := GameLayout(OWW, OWH, w, h)
	st := NewOWState()
	st.Cursor = game.Vec{X: 14, Y: 6}
	st.Trail = []game.Vec{{X: 13, Y: 6}, {X: 12, Y: 6}, {X: 11, Y: 6}}
	st.TrailAge = []int{24, 12, 4}
	f := RenderOverworld(w, h, st, 0, pal)
	var fgs []int
	for _, v := range st.Trail {
		cx, cy := l.center(v.X, v.Y)
		c := f.C[cy*f.W+cx]
		if c.R != '·' {
			t.Fatalf("trail cell (%d,%d) = %q, want ·", v.X, v.Y, c.R)
		}
		fgs = append(fgs, c.FG)
	}
	for i := 1; i < len(fgs); i++ {
		if fgs[i-1] <= fgs[i] {
			t.Fatalf("trail must dim with age: %v", fgs)
		}
	}
	// A zero-age step is not drawn (the cell keeps its pad texture).
	st2 := NewOWState()
	st2.Cursor = game.Vec{X: 6, Y: 10}
	st2.Trail = []game.Vec{{X: 5, Y: 10}}
	st2.TrailAge = []int{0}
	f2 := RenderOverworld(w, h, st2, 0, pal)
	cx, cy := l.center(5, 10)
	if f2.C[cy*f2.W+cx].R == '·' {
		t.Fatal("zero-age trail cell still drawn")
	}
}

// The first-run hint must actually show on first entry: Grak starts on the
// Rift, so the old "only when the floor line is empty" trigger never fired.
func TestOWFirstRunHint(t *testing.T) {
	pal := Palette()
	f := RenderOverworld(62, 19, NewOWState(), 0, pal)
	row := strings.Split(f.Text(), "\n")[1]
	if !strings.Contains(row, "wasd walk the lair") {
		t.Fatalf("first-run hint missing from the voice row: %q", row)
	}
}

// The ledger's floor ids wear each room's accent colour (rift 208, halls 110,
// garden 45, rotunda 178, heart 220), so the ledger reads as the map.
func TestOwLedgerAccents(t *testing.T) {
	pal := Palette()
	w, h := 62, 19
	f := RenderOverworld(w, h, NewOWState(), 0, pal)
	seen := map[int]bool{}
	for x := 0; x < w; x++ {
		c := f.C[(h-4)*f.W+x]
		if c.R == 'r' || c.R == 'h' || c.R == 'g' {
			seen[c.FG] = true
		}
	}
	for _, fg := range []int{208, 110, 45, 178, 220} {
		if !seen[fg] {
			t.Fatalf("ledger row missing accent colour %d", fg)
		}
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
