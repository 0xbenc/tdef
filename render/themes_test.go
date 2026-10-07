package render

import (
	"testing"

	"github.com/0xbenc/termtd/game"
)

// firstWall returns the first wall cell of a map (row-major). Every level has
// at least one wall (the border), so this always finds one.
func firstWall(m *game.Map) game.Vec {
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if m.At(game.Vec{X: x, Y: y}) == game.CellWall {
				return game.Vec{X: x, Y: y}
			}
		}
	}
	return game.Vec{}
}

// Each floor's stinger must actually render — its signature glyph appears
// somewhere on the map (checked over a few frames, since the light animates).
func TestStingerSignatureGlyphs(t *testing.T) {
	cases := []struct {
		level string
		load  func() (*game.Map, error)
		glyph rune
	}{
		{"winding", func() (*game.Map, error) { return game.LoadLevel("winding") }, '✦'},
		{"hub", func() (*game.Map, error) { return game.LoadLevel("hub") }, '●'},
		{"garden", func() (*game.Map, error) { return game.LoadLevel("garden") }, '≈'},
		{"canyon", func() (*game.Map, error) { return game.LoadLevel("canyon") }, '║'},
		{"maze1", func() (*game.Map, error) { return game.MazeFromSeed(7) }, '░'},
		{"heart", func() (*game.Map, error) { return game.LoadBoss() }, '○'},
	}
	pal := Palette()
	for _, c := range cases {
		m, err := c.load()
		if err != nil {
			t.Fatal(err)
		}
		g := game.NewState(m)
		seen := false
		for _, frame := range []int{0, 15, 30, 45, 60} {
			f := Render(g, &UI{Level: c.level}, pal, 62, 19, frame)
			for i := range f.C {
				if f.C[i].R == c.glyph {
					seen = true
					break
				}
			}
			if seen {
				break
			}
		}
		if !seen {
			t.Errorf("%s: stinger glyph %q never rendered", c.level, c.glyph)
		}
	}
}

// No two rooms of the lair may wear the same stone: each named level (plus the
// seeded maze and the heart) must render its walls in its own theme's palette,
// and no two levels may share a wall colour.
func TestLevelThemesAreDistinct(t *testing.T) {
	maps := map[string]*game.Map{}
	for _, name := range []string{"winding", "hub", "garden", "canyon"} {
		m, err := game.LoadLevel(name)
		if err != nil {
			t.Fatal(err)
		}
		maps[name] = m
	}
	mz, err := game.MazeFromSeed(1)
	if err != nil {
		t.Fatal(err)
	}
	maps["maze1"] = mz
	heart, err := game.LoadBoss()
	if err != nil {
		t.Fatal(err)
	}
	maps["heart"] = heart

	pal := Palette()
	seen := map[int]string{} // wall base colour -> level that first used it
	for _, level := range []string{"winding", "hub", "garden", "canyon", "maze1", "heart"} {
		m := maps[level]
		th := themeForLevel(level)
		g := game.NewState(m)
		l := GameLayout(m.W, m.H, 62, 19)
		f := Render(g, &UI{Level: level, Selected: NoSelection, Cursor: game.Vec{X: -1, Y: -1}}, pal, 62, 19, 0)
		w := firstWall(m)
		x, y := l.X(w.X), l.Y(w.Y)
		bg := f.C[y*f.W+x].BG
		if bg != th.Wall && bg != th.WallHi && bg != th.WallLo {
			t.Errorf("%s: wall %v renders BG %d, want one of {Wall %d, WallHi %d, WallLo %d}",
				level, w, bg, th.Wall, th.WallHi, th.WallLo)
		}
		if other, dup := seen[th.Wall]; dup {
			t.Errorf("wall colour %d is shared by %s and %s", th.Wall, other, level)
		}
		seen[th.Wall] = level
	}
}
