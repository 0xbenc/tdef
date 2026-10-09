package render

import (
	"strings"
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/internal/copytext"
)

func TestRuneforgePlacementHintFitsMinimumWidth(t *testing.T) {
	hint := copytext.Text("ui.keyboard.place") + copytext.Text("ui.keyboard.rotate")
	if len([]rune(hint)) > 62-4 {
		t.Fatalf("placement hint clips at minimum width: %q", hint)
	}
}

func TestKeyboardBuildSitesAndCursor(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m)
	g.Gold = 10000
	ui := UI{Selected: NoSelection, PlacingOn: true, Placing: game.TowerGunner}
	l := GameLayout(m.W, m.H, 120, 40)
	f := &Frame{W: 120, H: 40, C: make([]Cell, 120*40)}
	drawPlacementSites(f, g, &ui, l)
	var site game.Vec
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			v := game.Vec{X: x, Y: y}
			cx, cy := l.center(x, y)
			if cx < 1 || cx >= f.W-1 || cy < ChromeTop || cy >= f.H-ChromeBot {
				continue
			}
			if marked := f.C[cy*f.W+cx].R == '·'; marked != g.CanBuild(v, ui.Placing) {
				t.Fatalf("build marker disagrees with placement at %+v", v)
			}
			if g.CanBuild(v, ui.Placing) {
				site = v
			}
		}
	}
	tower := g.Build(site, game.TowerGunner)
	ui.Cursor, ui.Selected, ui.PlacingOn = site, tower.ID, false
	cx, cy := l.center(site.X, site.Y)
	for _, frame := range []int{0, 15} {
		drawKeyboardCursor(f, g, &ui, l, frame)
		c := f.C[cy*f.W+cx]
		if c.FG != 231 || c.BG == 0 || !c.Bold || c.R != tower.Spec().Short {
			t.Fatalf("selected defender cursor disappeared: %+v", c)
		}
	}
}

func TestPlacementSitesShowGoldShortage(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	site := game.Vec{X: 4, Y: 1}
	l := GameLayout(m.W, m.H, 120, 40)
	for _, tc := range []struct {
		name    string
		gold    int
		blocked string
		fg, bg  int
	}{
		{"short one coin", 49, "", 214, 58},
		{"exact price", 50, "", 46, 22},
		{"occupied", 0, "occupied", 0, 0},
		{"locked", 0, "locked", 0, 0},
		{"specialist already placed", 0, "specialist", 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewState(m)
			g.Gold = tc.gold
			ui := UI{Selected: NoSelection, PlacingOn: true, Placing: game.TowerGunner, Cursor: site}
			switch tc.blocked {
			case "occupied":
				g.Towers = []*game.Tower{{Cell: site, Kind: game.TowerGunner, Level: 1}}
			case "locked":
				g.TrainingStage = 1
			case "specialist":
				ui.Placing = game.TowerRuneforge
				g.Towers = []*game.Tower{{Cell: game.Vec{X: 5, Y: 1}, Kind: game.TowerRuneforge, Level: 1}}
			}
			f := &Frame{W: 120, H: 40, C: make([]Cell, 120*40)}
			drawPlacementSites(f, g, &ui, l)
			cx, cy := l.center(site.X, site.Y)
			mark := f.C[cy*f.W+cx]
			if tc.blocked == "" {
				if mark.R != '·' || mark.FG != tc.fg || mark.BG != tc.bg {
					t.Fatalf("unexpected build marker: %+v", mark)
				}
			} else if mark.R == '·' {
				t.Fatal("blocked placement showed a build marker")
			}
			for _, v := range []game.Vec{m.Path[0], {X: 0, Y: 0}} {
				x, y := l.center(v.X, v.Y)
				if f.C[y*f.W+x].R == '·' {
					t.Fatalf("road or wall marked buildable at %v", v)
				}
			}
			if tc.blocked == "" {
				drawKeyboardCursor(f, g, &ui, l, 0)
				bg := f.C[cy*f.W+cx].BG
				if tc.gold < 50 && bg != 58 {
					t.Fatalf("unaffordable valid cursor should be amber, got %d", bg)
				}
			}
		})
	}
}

func TestKeyboardActionsFitAndExplainPriority(t *testing.T) {
	m, _ := game.LoadLevel("winding")
	g := game.NewState(m)
	g.Gold = 10000
	var tower *game.Tower
	for y := 0; y < m.H && tower == nil; y++ {
		for x := 0; x < m.W && tower == nil; x++ {
			tower = g.Build(game.Vec{X: x, Y: y}, game.TowerGunner)
		}
	}
	ui := UI{Selected: NoSelection, Cursor: tower.Cell}
	for _, width := range []int{62, 80, 120} {
		f := Render(g, &ui, Palette(), width, 24, 0)
		text := f.Text()
		for _, label := range []string{"u upgrade", "t first", "x sell"} {
			if !strings.Contains(text, label) {
				t.Fatalf("width %d missing action %q", width, label)
			}
		}
		for _, label := range []string{"sell", "first"} {
			if count := strings.Count(text, label); count != 1 {
				t.Fatalf("width %d repeats %q %d times", width, label, count)
			}
		}
		if !strings.Contains(text, "Lv1") || !strings.Contains(text, "Tab next") {
			t.Fatalf("width %d lost level or defender cycling", width)
		}
		if f.C[(f.H-2)*f.W+f.W-1].R != '│' {
			t.Fatalf("width %d actions border: %+v", width, f.C[(f.H-2)*f.W+f.W-1])
		}
	}
	g.Gold = 0
	f := Render(g, &ui, Palette(), 80, 24, 0)
	if !strings.Contains(f.Text(), "u need") {
		t.Fatal("unaffordable upgrade not explained")
	}
}
