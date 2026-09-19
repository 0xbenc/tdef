package render

import (
	"testing"

	"tdef/game"
)

// The keyboard cursor must be visible on every ground cell — grass, road,
// and wall — not just buildable grass.
func TestCursorVisibleOnAllGround(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m, 1, false)
	pal := Palette()
	l := ComputeLayout(m.W, m.H, 1)
	cases := []struct {
		name string
		v    game.Vec
	}{
		{"grass", game.Vec{X: 7, Y: 2}},
		{"path", game.Vec{X: 10, Y: 1}},
		{"wall", game.Vec{X: 5, Y: 0}},
	}
	for _, c := range cases {
		ui := &UI{Cursor: c.v, Placing: game.TowerGunner, Selected: -1, Scale: 1}
		f := Render(g, ui, pal)
		x, y := l.X(c.v.X), l.Y(c.v.Y)
		got := f.C[y*f.W+x]
		if got.R != '◻' {
			t.Errorf("%s cell %v: cursor cell = %q (fg %d), want '◻'", c.name, c.v, got.R, got.FG)
		}
	}
}

// The cursor must not clobber entity glyphs (e.g. the spawn arrow).
func TestCursorDoesNotClobberGlyphs(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m, 1, false)
	pal := Palette()
	l := ComputeLayout(m.W, m.H, 1)
	ui := &UI{Cursor: m.Spawn, Placing: game.TowerGunner, Selected: -1, Scale: 1}
	f := Render(g, ui, pal)
	x, y := l.X(m.Spawn.X), l.Y(m.Spawn.Y)
	if got := f.C[y*f.W+x].R; got != '▶' {
		t.Errorf("spawn cell = %q, want '▶'", got)
	}
}
