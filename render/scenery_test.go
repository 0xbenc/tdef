package render

import (
	"reflect"
	"testing"

	"github.com/0xbenc/termtd/game"
)

func sceneryMaps(t *testing.T) map[string]*game.Map {
	t.Helper()
	maps := map[string]*game.Map{}
	for _, name := range game.LevelNames() {
		m, err := game.LoadLevel(name)
		if err != nil {
			t.Fatal(err)
		}
		maps[name] = m
	}
	heart, err := game.LoadBoss()
	if err != nil {
		t.Fatal(err)
	}
	maps["heart"] = heart
	for _, seed := range []int64{7, 91} {
		m, err := game.MazeFromSeed(seed)
		if err != nil {
			t.Fatal(err)
		}
		name := "maze7"
		if seed == 91 {
			name = "maze91"
		}
		maps[name] = m
	}
	m, err := game.GenerateMap(19, 9, 7)
	if err != nil {
		t.Fatal(err)
	}
	maps["maze-small"] = m
	return maps
}

// Scenery must never paint a road, a build tile, or the HUD. Check every
// terminal cell at every scale, including arbitrary procedural wall shapes.
func TestSceneryOnlyChangesBlockedTerrain(t *testing.T) {
	for name, m := range sceneryMaps(t) {
		t.Run(name, func(t *testing.T) {
			originalCells := append([]game.CellKind(nil), m.Cell...)
			originalPath := append([]game.Vec(nil), m.Path...)
			for scale := 1; scale <= 4; scale++ {
				w, h := CaptureSize(m.W, m.H, scale)
				l := GameLayout(m.W, m.H, w, h)
				sentinel := Cell{R: '?', FG: 255, BG: 17, Bold: true}
				blank := func() *Frame {
					f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
					for i := range f.C {
						f.C[i] = sentinel
					}
					return f
				}
				a, b := blank(), blank()
				drawScenery(a, m, themeForLevel(name), l, 37)
				drawScenery(b, m, themeForLevel(name), l, 37)
				if !reflect.DeepEqual(a.C, b.C) {
					t.Fatalf("scale %d: scenery is nondeterministic", scale)
				}
				changes := 0
				for y := 0; y < h; y++ {
					for x := 0; x < w; x++ {
						inside := x >= l.Ox && x < l.Ox+m.W*scale && y >= l.Oy && y < l.Oy+m.H*scale
						wall := false
						if inside {
							wall = m.At(game.Vec{X: (x - l.Ox) / scale, Y: (y - l.Oy) / scale}) == game.CellWall
						}
						c := a.C[y*w+x]
						if !wall && c != sentinel {
							t.Fatalf("scale %d: scenery overwrote functional terrain or chrome at %d,%d", scale, x, y)
						}
						if c != sentinel {
							changes++
						}
					}
				}
				if changes == 0 {
					t.Fatalf("scale %d: no scenery drawn", scale)
				}
			}
			if !reflect.DeepEqual(m.Cell, originalCells) || !reflect.DeepEqual(m.Path, originalPath) {
				t.Fatal("rendering mutated map geometry")
			}
		})
	}
}

func TestSceneryClipsToPlayfield(t *testing.T) {
	m, err := game.LoadLevel("hub")
	if err != nil {
		t.Fatal(err)
	}
	f := &Frame{W: 30, H: 12, C: make([]Cell, 360)}
	// Deliberately cropped, like a compact map preview.
	l := Layout{Ox: -4, Oy: 1, Scale: 2, W: f.W, H: f.H}
	drawScenery(f, m, themeForLevel("hub"), l, 37)
	for y := 0; y < f.H; y++ {
		for x := 0; x < f.W; x++ {
			if x == 0 || x == f.W-1 || y < ChromeTop || y >= f.H-ChromeBot {
				if f.C[y*f.W+x] != (Cell{}) {
					t.Fatalf("scenery escaped playfield at %d,%d", x, y)
				}
			}
		}
	}
}

func TestCursorVisibleOnDecoratedWall(t *testing.T) {
	m, err := game.LoadLevel("hub")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m)
	for scale := 1; scale <= 4; scale++ {
		w, h := CaptureSize(m.W, m.H, scale)
		l := GameLayout(m.W, m.H, w, h)
		ui := &UI{Level: "hub", Selected: NoSelection, Cursor: game.Vec{X: 3, Y: 0}, Speed: 1}
		f := Render(g, ui, Palette(), w, h, 37)
		cx, cy := l.center(3, 0)
		if got := f.C[cy*w+cx].R; got != '◻' {
			t.Errorf("scale %d: wall cursor = %q", scale, got)
		}
	}
}
