package render

import (
	"strings"
	"testing"

	"tdef/game"
)

func TestComputeScale(t *testing.T) {
	cases := []struct {
		tw, th, want int
	}{
		{80, 24, 1},   // small terminal -> 1x
		{120, 40, 2},  // typical big -> 2x
		{150, 50, 3},  // large -> 3x
		{400, 100, 4}, // huge -> capped at 4x
		{30, 10, 1},   // tiny -> 1x (min)
		{0, 0, 1},     // unknown size -> 1x
	}
	for _, c := range cases {
		if got := ComputeScale(45, 13, c.tw, c.th); got != c.want {
			t.Errorf("ComputeScale(45,13,%d,%d) = %d, want %d", c.tw, c.th, got, c.want)
		}
	}
}

func TestComputeLayout(t *testing.T) {
	l := ComputeLayout(45, 13, 2)
	if l.Scale != 2 {
		t.Errorf("Scale = %d, want 2", l.Scale)
	}
	// map is 45*2=90 wide; frame must be at least 92 (mapW+2) and >= FrameW(62)
	if l.W < 92 {
		t.Errorf("W = %d, want >= 92", l.W)
	}
	// height = HUDRows(2) + mapH(13*2=26) + 4 = 32
	if l.H != 32 {
		t.Errorf("H = %d, want 32", l.H)
	}
	// map origin: centered horizontally, below HUD
	if l.Oy != 2 {
		t.Errorf("Oy = %d, want 2", l.Oy)
	}
	// X/Y mapping
	if got := l.X(10); got != l.Ox+20 {
		t.Errorf("X(10) = %d, want %d", got, l.Ox+20)
	}
	if got := l.Y(5); got != l.Oy+10 {
		t.Errorf("Y(5) = %d, want %d", got, l.Oy+10)
	}
	// menu top is the last 4 rows
	if got := l.MenuTop(); got != l.H-4 {
		t.Errorf("MenuTop() = %d, want %d", got, l.H-4)
	}
	// scale clamps to >=1
	if c := ComputeLayout(45, 13, 0).Scale; c != 1 {
		t.Errorf("clamp Scale = %d, want 1", c)
	}
}

func TestDrawRing(t *testing.T) {
	l := ComputeLayout(45, 13, 1)
	f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
	fx := &game.Fx{Pos: game.Pos{X: 20, Y: 6}, TTL: 0.25, Max: 0.25, Ring: 2.0, Color: 203}
	drawRing(f, l, fx)
	dots := 0
	for _, c := range f.C {
		if c.R == '·' && c.FG == 203 {
			dots++
		}
	}
	// a ring of radius ~2 should draw several dots (a circle, not a blob)
	if dots < 6 {
		t.Errorf("ring dots = %d, want >= 6", dots)
	}
}

func TestMinFrame(t *testing.T) {
	w, h := MinFrame(45, 13)
	if w != 62 || h != 19 {
		t.Errorf("MinFrame(45,13) = %dx%d, want 62x19", w, h)
	}
	if w2, _ := MinFrame(100, 13); w2 != 102 {
		t.Errorf("MinFrame(100,13) w = %d, want 102", w2)
	}
}

func TestRenderTooSmall(t *testing.T) {
	f := RenderTooSmall(30, 10, 62, 19)
	if f.W != 30 || f.H != 10 {
		t.Errorf("size = %dx%d, want 30x10", f.W, f.H)
	}
	text := f.Text()
	if !strings.Contains(text, "too small") {
		t.Errorf("missing 'too small' message:\n%s", text)
	}
	if !strings.Contains(text, "62x19") {
		t.Errorf("missing size hint:\n%s", text)
	}
	// 30 columns is narrower than the notice lines, so check for the
	// substring that survives clipping, not the full sentence.
	if !strings.Contains(text, "paused") {
		t.Errorf("missing pause note:\n%s", text)
	}
}

func TestFormatTime(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0m00s"},
		{59, "0m59s"},
		{60, "1m00s"},
		{155, "2m35s"},
	}
	for _, c := range cases {
		if got := formatTime(c.in); got != c.want {
			t.Errorf("formatTime(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDrawGameOver(t *testing.T) {
	l := ComputeLayout(45, 13, 1)
	f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
	g := &game.State{
		Status:     game.StatusVictory,
		Wave:       20,
		TotalKills: 100,
		TotalLeaks: 3,
		Score:      500,
		MaxCombo:   5,
		Time:       150,
		Towers:     []*game.Tower{{}, {}},
	}
	drawGameOver(f, g, &UI{BestScore: 900}, Palette())
	text := f.Text()
	for _, want := range []string{"VICTORY", "20/20", "100", "2m30s"} {
		if !strings.Contains(text, want) {
			t.Errorf("game-over missing %q:\n%s", want, text)
		}
	}
}

func TestUpgradePipsStyledLikeMenuSlot(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m)
	g.Gold = 1000
	v := game.Vec{X: 7, Y: 2}
	tw := g.Build(v, game.TowerGunner)
	g.Upgrade(tw)
	g.Upgrade(tw)
	pal := Palette()
	l := ComputeLayout(m.W, m.H, 1)
	f := Render(g, &UI{Placing: game.TowerGunner, Selected: -1, Scale: 1}, pal)
	x, y := l.center(v.X, v.Y)
	want := Cell{R: '▪', FG: pal.Bright, BG: pal.Tower[game.TowerGunner], Bold: true}
	for i := 1; i < 3; i++ {
		if got := f.C[y*f.W+x-i]; got != want {
			t.Errorf("pip %d = %+v, want %+v", i, got, want)
		}
	}
}

// The selected-tower info line renders on the 62-col menu row at 1x scale;
// it must never be wider than the frame or it clips.
func TestTowerInfoFitsFrame(t *testing.T) {
	for k := game.TowerKind(0); k < game.TowerCount; k++ {
		invested := 0
		for level := 1; level <= 3; level++ {
			invested += game.TowerSpecs[k].Cost[level-1]
			refund := int(float64(invested) * game.SellRefund)
			upCost := 0
			if level < 3 {
				upCost = game.TowerSpecs[k].Cost[level]
			}
			for mode := game.TargetMode(0); mode < game.TargetModeCount; mode++ {
				tw := &game.Tower{Kind: k, Level: level, TargetMode: mode}
				line := towerInfo(tw, upCost, refund)
				if n := len([]rune(line)); n > FrameW {
					t.Errorf("%s lv%d %s: %d runes, want <= %d: %q",
						game.TowerSpecs[k].Name, level, mode.Name(), n, FrameW, line)
				}
			}
		}
	}
}

// A tower's level pips are drawn to its LEFT, so an upgraded tower sitting
// right of a neighbor used to overwrite the neighbor's glyph (the pip pass
// and glyph pass were interleaved in build order). Pips must skip occupied
// cells so both towers stay visible.
func TestAdjacentUpgradedTowersDontEraseEachOther(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	var left, right game.Vec
	found := false
	for y := 0; y < m.H && !found; y++ {
		for x := 0; x+1 < m.W && !found; x++ {
			if m.At(game.Vec{X: x, Y: y}) == game.CellGrass && m.At(game.Vec{X: x + 1, Y: y}) == game.CellGrass {
				left, right = game.Vec{X: x, Y: y}, game.Vec{X: x + 1, Y: y}
				found = true
			}
		}
	}
	if !found {
		t.Fatal("no adjacent grass cells found")
	}
	g := game.NewState(m)
	g.Gold = 1000
	// Build the LEFT tower first: with the old interleaved pass the right
	// tower's pip (drawn later) clobbered its glyph.
	tl := g.Build(left, game.TowerGunner)
	tr := g.Build(right, game.TowerCannon)
	if tl == nil || tr == nil {
		t.Fatal("build failed")
	}
	for i := 0; i < 2; i++ { // left -> level 3 (2 pips)
		g.Upgrade(tl)
	}
	g.Upgrade(tr) // right -> level 2 (1 pip, aimed at the left tower's cell)
	l := ComputeLayout(m.W, m.H, 1)
	f := Render(g, &UI{Selected: NoSelection, Scale: 1}, Palette())
	lx, ly := l.center(left.X, left.Y)
	rx, ry := l.center(right.X, right.Y)
	if got := f.C[ly*f.W+lx].R; got != 'G' {
		t.Errorf("left tower glyph = %q, want 'G' (clobbered by right tower's pip?)", got)
	}
	if got := f.C[ry*f.W+rx].R; got != 'C' {
		t.Errorf("right tower glyph = %q, want 'C'", got)
	}
	// The left tower's own pips still render on the empty cells to its left.
	if got := f.C[ly*f.W+lx-1].R; got != '▪' {
		t.Errorf("left tower pip = %q, want '▪'", got)
	}
}

func TestANSIPositionsRowsWithCUP(t *testing.T) {
	f := &Frame{W: 4, H: 3, C: make([]Cell, 12)}
	f.Put(0, 0, 'a', 220, 0)
	f.Put(0, 1, 'b', 220, 0)
	f.Put(0, 2, 'c', 220, 0)
	s := f.ANSI()
	if strings.Contains(s, "\x1b[M") {
		t.Error("ANSI still uses bare cursor-up for row advance")
	}
	for _, want := range []string{"\x1b[2;1H", "\x1b[3;1H"} {
		if !strings.Contains(s, want) {
			t.Errorf("ANSI missing CUP sequence %q", want)
		}
	}
}

// The intro's help lines and prompt must each land on their own row inside
// the frame: at scale 1 the last line used to be drawn one row past the
// bottom and the prompt smashes into the third help line.
func TestRenderIntroFitsFrame(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	for _, scale := range []int{1, 2, 4} {
		f := RenderIntro(m, "winding", game.Normal, Palette(), scale)
		lines := strings.Split(f.Text(), "\n")
		markers := []string{
			"Enemies walk the path",
			"1-7 pick tower",
			"n start wave early",
			"Don't let them reach E",
			"press any key to start",
		}
		seen := map[int]bool{}
		for _, want := range markers {
			found := -1
			for i, ln := range lines {
				if strings.Contains(ln, want) {
					found = i
					break
				}
			}
			if found < 0 {
				t.Errorf("scale %d: intro missing %q:\n%s", scale, want, f.Text())
				continue
			}
			if seen[found] {
				t.Errorf("scale %d: two intro lines share row %d", scale, found)
			}
			seen[found] = true
		}
	}
}

func TestDrawHPBar(t *testing.T) {
	f := &Frame{W: 20, H: 5, C: make([]Cell, 20*5)}
	drawHPBar(f, 10, 2, 0.5)
	// 3 segments drawn at (9..11, 2)
	if f.C[2*20+9].R != '■' || f.C[2*20+10].R != '■' || f.C[2*20+11].R != '■' {
		t.Errorf("HP bar not drawn: %+v %+v %+v", f.C[2*20+9], f.C[2*20+10], f.C[2*20+11])
	}
	// hp=0.5 -> bar color is yellow (214); 2 of 3 segments filled, last empty (238)
	if f.C[2*20+9].FG != 214 || f.C[2*20+10].FG != 214 {
		t.Errorf("filled segments FG = %d,%d, want 214,214", f.C[2*20+9].FG, f.C[2*20+10].FG)
	}
	if f.C[2*20+11].FG != 238 {
		t.Errorf("empty segment FG = %d, want 238", f.C[2*20+11].FG)
	}
}
