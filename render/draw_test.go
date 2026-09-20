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

// rowWidth is the last non-space column of row y plus one (0 if blank).
func rowWidth(f *Frame, y int) int {
	w := 0
	for x := 0; x < f.W; x++ {
		if r := f.C[y*f.W+x].R; r != 0 && r != ' ' {
			w = x + 1
		}
	}
	return w
}

// The menu hint/help rows sit at the bottom of a frame that is exactly
// FrameW wide at scale 1; anything longer is silently clipped.
func TestMenuLinesFitFrame(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	l := ComputeLayout(m.W, m.H, 1)
	g := &game.State{Map: m, Gold: 1000}
	check := func(ui *UI, rows ...int) {
		f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
		drawMenu(f, g, ui, Palette())
		for _, y := range rows {
			if w := rowWidth(f, y); w > FrameW {
				t.Errorf("menu row %d is %d cols, want <= %d", y, w, FrameW)
			}
		}
	}
	hintRow := l.H - 4 + 2
	// compact hint
	check(&UI{Placing: game.TowerGunner, Scale: 1}, hintRow)
	// two help lines
	check(&UI{Placing: game.TowerGunner, Help: true, Scale: 1}, hintRow, hintRow+1)
}

// Every HUD line combination must fit the 62-col frame at scale 1. The old
// break line (theme + preview + countdown + bonus) exceeded it from wave 17
// on, hiding the early-start bonus hint.
func TestHUDLinesFitFrame(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	l := ComputeLayout(m.W, m.H, 1)
	check := func(g *game.State, ui *UI) {
		f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
		drawHUD(f, g, ui, Palette())
		for y := 0; y < 2; y++ {
			if w := rowWidth(f, y); w > FrameW {
				t.Errorf("wave %d msg=%q paused=%v: HUD row %d is %d cols, want <= %d",
					g.Wave, ui.Message, ui.Paused, y, w, FrameW)
			}
		}
	}
	msgs := []string{"", "wave 18 cleared +93g", "Faster ones are coming."}
	for wave := 1; wave < game.MaxWaves; wave++ {
		for _, paused := range []bool{false, true} {
			for _, msg := range msgs {
				// inter-wave break (telegraph or transient message)
				g := &game.State{
					Map: m, Status: game.StatusRunning, Wave: wave,
					NextWaveAt: 9.5, Time: 5.0, WaveActive: false,
					Gold: 9999, Lives: 15, Score: 999999,
				}
				check(g, &UI{Speed: 4, Paused: paused, Message: msg, Scale: 1})
			}
		}
	}
	// active wave, worst-case combo + leak message
	g := &game.State{
		Map: m, Status: game.StatusRunning, Wave: 20, WaveActive: true,
		Combo: 123, Gold: 9999, Lives: 15, Score: 999999,
	}
	check(g, &UI{Speed: 4, Paused: true, Message: "leak! -6 lives", Scale: 1})
}

func TestGameLayoutCentering(t *testing.T) {
	cases := []struct {
		tw, th, s, ox, oy int
	}{
		{62, 19, 1, 8, 2},
		{80, 24, 1, 17, 4},
		{92, 32, 2, 1, 2},
		{137, 45, 3, 1, 2},
		{200, 70, 4, 10, 8},
		{30, 10, 1, 1, 2}, // clamped origins
	}
	for _, c := range cases {
		l := GameLayout(45, 13, c.tw, c.th)
		if l.Scale != c.s || l.Ox != c.ox || l.Oy != c.oy {
			t.Errorf("GameLayout(45,13,%d,%d) = s%d at (%d,%d), want s%d at (%d,%d)",
				c.tw, c.th, l.Scale, l.Ox, l.Oy, c.s, c.ox, c.oy)
		}
		if l.W != c.tw || l.H != c.th {
			t.Errorf("frame size = %dx%d, want %dx%d", l.W, l.H, c.tw, c.th)
		}
		if got := l.MenuTop(); got != c.th-4 {
			t.Errorf("MenuTop() = %d, want %d", got, c.th-4)
		}
	}
}

func TestComputeScaleSteps(t *testing.T) {
	cases := []struct {
		tw, th, want int
	}{
		{62, 19, 1},   // minimum frame: exactly 1x
		{92, 32, 2},   // 2x threshold
		{137, 45, 3},  // 3x threshold
		{182, 58, 4},  // 4x threshold
		{91, 32, 1},   // one column short of 2x
		{92, 31, 1},   // one row short of 2x
		{136, 45, 2},  // one column short of 3x
		{400, 100, 4}, // capped at 4x
		{30, 10, 1},   // tiny -> 1x
		{0, 0, 1},     // unknown size -> 1x
	}
	for _, c := range cases {
		if got := ComputeScale(45, 13, c.tw, c.th); got != c.want {
			t.Errorf("ComputeScale(45,13,%d,%d) = %d, want %d", c.tw, c.th, got, c.want)
		}
	}
}

func TestTowerSlotDistribution(t *testing.T) {
	slots := TowerSlots(62, 19)
	if len(slots) != 7 {
		t.Fatalf("slots = %d, want 7", len(slots))
	}
	wantRow0 := []int{1, 16, 31, 46}
	wantRow1 := []int{1, 21, 41}
	for i, x := range wantRow0 {
		if slots[i].Y != 15 || slots[i].X != x {
			t.Errorf("slot %d = (%d,%d), want (%d,15)", i, slots[i].X, slots[i].Y, x)
		}
	}
	for i, x := range wantRow1 {
		if slots[4+i].Y != 16 || slots[4+i].X != x {
			t.Errorf("slot %d = (%d,%d), want (%d,16)", 4+i, slots[4+i].X, slots[4+i].Y, x)
		}
	}
	// At every width from the minimum up, no label overruns its cell:
	// the next slot starts at least one column after the label ends.
	for _, tw := range []int{62, 63, 80, 120} {
		s := TowerSlots(tw, 20)
		rows := map[int][]MenuSlot{}
		for _, sl := range s {
			rows[sl.Y] = append(rows[sl.Y], sl)
		}
		for _, row := range rows {
			cell := (tw - 2) / len(row)
			for i, sl := range row {
				if want := 1 + i*cell; sl.X != want {
					t.Errorf("tw=%d slot %d X=%d, want %d", tw, i, sl.X, want)
				}
				if sl.W > cell {
					t.Errorf("tw=%d slot %d label %d cols > cell %d", tw, i, sl.W, cell)
				}
			}
		}
	}
}

func TestHeaderSegmentsAt62(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := &game.State{
		Map: m, Status: game.StatusRunning, Wave: 5, WaveActive: true,
		Combo: 7, Gold: 9999, Lives: 14, Score: 3120,
	}
	ui := &UI{Placing: game.TowerGunner, Selected: NoSelection, Speed: 1, Level: "winding"}
	check := func(tw, th int, present, absent []string) {
		f := RenderAt(g, ui, Palette(), tw, th)
		var b strings.Builder
		for x := 0; x < f.W; x++ {
			b.WriteRune(f.C[x].R)
		}
		for _, s := range present {
			if !strings.Contains(b.String(), s) {
				t.Errorf("tw=%d header missing %q: %q", tw, s, b.String())
			}
		}
		for _, s := range absent {
			if strings.Contains(b.String(), s) {
				t.Errorf("tw=%d header should drop %q: %q", tw, s, b.String())
			}
		}
	}
	// At 62 the level·diff segment is elided (first in the drop order).
	check(62, 19, []string{"tdef", "wave 5/20", "⛁"}, []string{"winding", "normal"})
	// At 80 everything fits.
	check(80, 24, []string{"tdef", "winding", "normal", "wave 5/20", "⛁"}, nil)
	if f := RenderAt(g, ui, Palette(), 62, 19); f.C[0].R != '╭' || f.C[f.W-1].R != '╮' {
		t.Errorf("top border corners missing: %q %q", f.C[0].R, f.C[f.W-1].R)
	}
}

// Every header combination must keep row 0 exactly full-width (the border
// is never broken) and row 1 inside the frame.
func TestHeaderNoOverflow(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	msgs := []string{"", "leak! -6 lives", "wave 18 cleared +93g", "Faster ones are coming. " + "and faster still, faster."}
	for _, tw := range []int{62, 63, 70, 80, 120} {
		for _, paused := range []bool{false, true} {
			for _, msg := range msgs {
				g := &game.State{
					Map: m, Status: game.StatusRunning, Wave: 19,
					NextWaveAt: 9.5, Time: 5.0, WaveActive: false,
					Gold: 9999, Lives: 15, Score: 999999,
				}
				ui := &UI{Speed: 4, Paused: paused, Message: msg, Level: "winding"}
				f := RenderAt(g, ui, Palette(), tw, 24)
				if w := rowWidth(f, 0); w != tw {
					t.Errorf("wave-active off paused=%v msg=%q: header row %d cols, want %d", paused, msg, w, tw)
				}
				if w := rowWidth(f, 1); w > tw {
					t.Errorf("paused=%v msg=%q: message row %d cols > %d", paused, msg, w, tw)
				}
				g.WaveActive = true
				g.Combo = 123
				f = RenderAt(g, ui, Palette(), tw, 24)
				if w := rowWidth(f, 0); w != tw {
					t.Errorf("wave-active paused=%v msg=%q: header row %d cols, want %d", paused, msg, w, tw)
				}
				if w := rowWidth(f, 1); w > tw {
					t.Errorf("wave-active paused=%v msg=%q: message row %d cols > %d", paused, msg, w, tw)
				}
			}
		}
	}
}

func TestRenderAtSmoke(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m)
	ui := &UI{Cursor: game.Vec{X: m.W / 2, Y: m.H / 2}, Placing: game.TowerGunner, Selected: NoSelection, Speed: 1, Level: "winding"}
	f := RenderAt(g, ui, Palette(), 80, 24)
	if f.W != 80 || f.H != 24 {
		t.Fatalf("frame = %dx%d, want 80x24", f.W, f.H)
	}
	text := f.Text()
	for _, want := range []string{"tdef", "1 Gunner 50", "⏎|place", "winding", "normal"} {
		if !strings.Contains(text, want) {
			t.Errorf("frame missing %q", want)
		}
	}
	if f.C[0].R != '╭' || f.C[f.W-1].R != '╮' || f.C[(f.H-1)*f.W].R != '╰' || f.C[f.H*f.W-1].R != '╯' {
		t.Errorf("rounded corners missing")
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
