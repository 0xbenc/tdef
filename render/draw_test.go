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
	g := game.NewState(m, 1, false)
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
