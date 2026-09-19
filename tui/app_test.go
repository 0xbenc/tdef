package tui

import (
	"strings"
	"testing"

	"tdef/game"
	"tdef/render"
)

// A fresh game must start with no tower selected: the zero value of
// UI.Selected is tower ID 0, which would hide the cursor and, once the
// first tower exists, highlight it as selected.
func TestFreshUIHasNoSelection(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	ui := freshUI(m, 1)
	if ui.Selected != render.NoSelection {
		t.Fatalf("fresh UI Selected = %d, want %d", ui.Selected, render.NoSelection)
	}
	if ui.Speed != 1 || !ui.Paused || ui.Placing != game.TowerGunner {
		t.Errorf("fresh UI defaults wrong: %+v", ui)
	}
	// The cursor glyph must be visible in the rendered frame.
	f := render.Render(game.NewState(m), &ui, render.Palette())
	if !strings.Contains(f.Text(), "◻") {
		t.Error("cursor glyph not drawn for a fresh UI")
	}
}

func TestRestartClearsSelection(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	a := &App{g: game.NewState(m), level: "winding", ui: render.UI{Selected: 3, Speed: 4, Scale: 2}}
	a.restart()
	if a.ui.Selected != render.NoSelection {
		t.Fatalf("Selected after restart = %d, want %d", a.ui.Selected, render.NoSelection)
	}
	if a.ui.Speed != 4 || a.ui.Scale != 2 {
		t.Errorf("restart must keep preferences: speed=%d scale=%d", a.ui.Speed, a.ui.Scale)
	}
}

// A mouse click must not place a tower once the game is over: the overlay
// covers the field, and handle() used to route mouse events before the
// status check.
func TestMouseIgnoredAfterGameOver(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	var cell game.Vec
	found := false
	for y := 0; y < m.H && !found; y++ {
		for x := 0; x < m.W && !found; x++ {
			if m.At(game.Vec{X: x, Y: y}) == game.CellGrass {
				cell = game.Vec{X: x, Y: y}
				found = true
			}
		}
	}
	if !found {
		t.Fatal("no grass cell")
	}
	ox, oy, sc := 0, 0, 0
	newApp := func() *App {
		a := &App{g: game.NewState(m), screen: ScreenGame, layout: render.ComputeLayout(m.W, m.H, 1), ui: freshUI(m, 1)}
		ox, oy, sc = a.mapBounds()
		a.ui.Placing = game.TowerGunner
		a.ui.PlacingOn = true
		return a
	}
	dead := newApp()
	dead.g.Status = game.StatusDefeat
	dead.handle(Event{Mouse: true, Btn: 0, Press: true, X: ox + cell.X*sc, Y: oy + cell.Y*sc})
	if len(dead.g.Towers) != 0 {
		t.Fatalf("mouse built %d tower(s) after game over", len(dead.g.Towers))
	}
	// Sanity: the same click mid-game builds.
	live := newApp()
	live.g.Gold = 100
	live.handle(Event{Mouse: true, Btn: 0, Press: true, X: ox + cell.X*sc, Y: oy + cell.Y*sc})
	if len(live.g.Towers) != 1 {
		t.Fatal("mouse did not build mid-game")
	}
}

func TestCycleSpeedWraps(t *testing.T) {
	a := &App{ui: render.UI{Speed: 1}}
	// three ups wrap 1->2->4->1, three downs wrap 1->4->2->1
	steps := []struct {
		dir  int
		want int
	}{
		{1, 2}, {1, 4}, {1, 1},
		{-1, 4}, {-1, 2}, {-1, 1},
	}
	for i, s := range steps {
		a.cycleSpeed(s.dir)
		if a.ui.Speed != s.want {
			t.Fatalf("step %d (dir %d): speed = %d, want %d", i, s.dir, a.ui.Speed, s.want)
		}
	}
}
