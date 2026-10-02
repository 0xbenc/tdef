package tui

import (
	"strings"
	"testing"

	"github.com/0xbenc/tdef/game"
	"github.com/0xbenc/tdef/render"
)

func TestWaveClearMessageShowsFloorReward(t *testing.T) {
	m, err := game.LoadLevel("garden")
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewState(m)
	g.Wave, g.WaveActive = 1, true
	// The last enemy has died and all spawns are complete.
	a := &App{g: g, ui: render.UI{Speed: 1}, waveStartLives: g.Lives}
	a.stepGame(1.0 / tickRate)
	if a.ui.Message != "wave 1 cleared +49g" || g.Gold != 357 {
		t.Errorf("wave clear: message=%q gold=%d", a.ui.Message, g.Gold)
	}
}

// A fresh game must start with no tower selected: the zero value of
// UI.Selected is tower ID 0, which would hide the cursor and, once the
// first tower exists, highlight it as selected.
func TestFreshUIHasNoSelection(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	ui := freshUI(m)
	if ui.Selected != render.NoSelection {
		t.Fatalf("fresh UI Selected = %d, want %d", ui.Selected, render.NoSelection)
	}
	if ui.Speed != 1 || !ui.Paused || ui.Placing != game.TowerGunner {
		t.Errorf("fresh UI defaults wrong: %+v", ui)
	}
	// The cursor glyph must be visible in the rendered frame.
	f := render.Render(game.NewState(m), &ui, render.Palette(), 80, 24, 0)
	if !strings.Contains(f.Text(), "◻") {
		t.Error("cursor glyph not drawn for a fresh UI")
	}
}

func TestRestartClearsSelection(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	a := &App{g: game.NewState(m), level: "winding", ui: render.UI{Selected: 3, Speed: 4, Help: true, Level: "winding"}}
	a.restart()
	if a.ui.Selected != render.NoSelection {
		t.Fatalf("Selected after restart = %d, want %d", a.ui.Selected, render.NoSelection)
	}
	if a.ui.Speed != 4 || !a.ui.Help || a.ui.Level != "winding" {
		t.Errorf("restart must keep preferences: speed=%d help=%v level=%q", a.ui.Speed, a.ui.Help, a.ui.Level)
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
		a := &App{g: game.NewState(m), screen: ScreenGame, layout: render.GameLayout(m.W, m.H, 62, 19), ui: freshUI(m)}
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

// The app's mouse mapping must agree with the renderer's layout at every
// size: same GameLayout, so a click at frame (x,y) lands on the map cell
// the pixels show.
func TestMouseMapBoundsAtScale2(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	a := &App{g: game.NewState(m), screen: ScreenGame, layout: render.GameLayout(m.W, m.H, 92, 32), ui: freshUI(m)}
	ox, oy, sc := a.mapBounds()
	if ox != 1 || oy != 2 || sc != 2 {
		t.Fatalf("mapBounds = (%d,%d,%d), want (1,2,2)", ox, oy, sc)
	}
	var cell game.Vec
	found := false
	for y := 0; y < m.H && !found; y++ {
		for x := 0; x < m.W && !found; x++ {
			if m.At(game.Vec{X: x, Y: y}) == game.CellGrass {
				cell, found = game.Vec{X: x, Y: y}, true
			}
		}
	}
	if !found {
		t.Fatal("no grass cell")
	}
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: ox + cell.X*sc + 1, Y: oy + cell.Y*sc + 1})
	if a.ui.Cursor != cell {
		t.Fatalf("cursor = %v, want %v", a.ui.Cursor, cell)
	}
	tw := a.g.Build(cell, game.TowerGunner)
	if tw == nil {
		t.Fatal("build failed")
	}
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: ox + cell.X*sc, Y: oy + cell.Y*sc})
	if a.ui.Selected != tw.ID {
		t.Errorf("selected = %d, want tower %d", a.ui.Selected, tw.ID)
	}
}

// Clicking a rendered tower slot label toggles placement; the geometry is
// TowerSlots, the same helper the renderer draws with.
func TestMenuSlotClickTogglesPlacing(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	a := &App{g: game.NewState(m), screen: ScreenGame, layout: render.GameLayout(m.W, m.H, 62, 19), ui: freshUI(m)}
	slots := render.TowerSlots(62, 19)
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: slots[1].X + 1, Y: slots[1].Y})
	if a.ui.Placing != game.TowerCannon || !a.ui.PlacingOn || a.ui.Selected != render.NoSelection {
		t.Fatalf("after slot click: %+v", a.ui)
	}
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: slots[1].X + 1, Y: slots[1].Y})
	if a.ui.PlacingOn {
		t.Fatal("second click must cancel placement")
	}
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: slots[6].X + 1, Y: slots[6].Y})
	if a.ui.Placing != game.TowerFlak || !a.ui.PlacingOn {
		t.Fatalf("flak slot click: %+v", a.ui)
	}
}

func TestTooSmallThreshold(t *testing.T) {
	a := &App{}
	cases := []struct {
		tw, th int
		want   bool
	}{
		{61, 19, true}, {62, 18, true}, {30, 10, true},
		{62, 19, false}, {80, 24, false}, {120, 40, false},
		{0, 0, false}, // unknown size renders normally
	}
	for _, c := range cases {
		if got := a.tooSmall(c.tw, c.th, 62, 19); got != c.want {
			t.Errorf("tooSmall(%d,%d) = %v, want %v", c.tw, c.th, got, c.want)
		}
	}
}

// In help mode the footer shows the expanded hint and the bottom border
// stays plain (no tower info), so nothing overflows at the minimum width.
func TestFooterHelpMode(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	g := &game.State{Map: m, Gold: 100, Status: game.StatusRunning}
	tw := g.Build(game.Vec{X: 7, Y: 2}, game.TowerGunner)
	if tw == nil {
		t.Fatal("build failed")
	}
	f := render.Render(g, &render.UI{Placing: game.TowerGunner, Selected: tw.ID, Help: true, Level: "winding"}, render.Palette(), 62, 19, 0)
	lines := strings.Split(f.Text(), "\n")
	if !strings.Contains(lines[17], "↑↓/wasd") {
		t.Errorf("help hint row missing: %q", lines[17])
	}
	if strings.Contains(lines[18], "▸") {
		t.Errorf("bottom border must stay plain in help mode: %q", lines[18])
	}
	f = render.Render(g, &render.UI{Placing: game.TowerGunner, Selected: tw.ID, Level: "winding"}, render.Palette(), 62, 19, 0)
	lines = strings.Split(f.Text(), "\n")
	if !strings.Contains(lines[18], "▸") {
		t.Errorf("selected-tower info missing from bottom border: %q", lines[18])
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

// Ending a run freezes simulation, not the animation that reveals its result.
func TestGameEndSequenceAdvancesToResult(t *testing.T) {
	for _, outcome := range []struct {
		status game.GameStatus
		title  string
	}{{game.StatusVictory, "VICTORY"}, {game.StatusDefeat, "DEFEAT"}} {
		t.Run(outcome.title, func(t *testing.T) {
			a := owTestApp(t)
			a.owRefresh()
			a.fromOW, a.owFloorID = true, "rotunda"
			m, err := game.LoadLevel("winding")
			if err != nil {
				t.Fatal(err)
			}
			a.enterGame(m, "winding", game.Normal)
			a.frameNo = 30 // as after arriving from the title and lair
			a.stepGame(1.0 / frameRate)
			if a.frameNo != 31 {
				t.Fatal("paused gameplay stopped its animation clock")
			}
			a.g.Wave = game.MaxWaves
			if outcome.status == game.StatusVictory {
				// Last wave cleared: no pending spawns or surviving enemies.
				a.g.WaveActive = true
				a.ui.Paused = false
			} else {
				a.g.Status = game.StatusDefeat
			}
			a.stepGame(1.0 / tickRate)
			if a.g.Status != outcome.status || !a.scored {
				t.Fatal("run did not finish and record its result")
			}
			end, simTime, tokens := a.ui.EndAtFrame, a.g.Time, a.lair.Tokens
			f := render.Render(a.g, &a.ui, render.Palette(), 94, 47, a.frameNo)
			if strings.Contains(f.Text(), outcome.title) {
				t.Fatal("result appeared before the end sequence")
			}
			for i := 0; i < 120; i++ {
				a.stepGame(1.0 / frameRate)
			}
			f = render.Render(a.g, &a.ui, render.Palette(), 94, 47, a.frameNo)
			if !strings.Contains(f.Text(), outcome.title) || !strings.Contains(f.Text(), "esc lair") {
				t.Fatalf("end sequence never revealed result and return controls:\n%s", f.Text())
			}
			if a.ui.EndAtFrame != end || a.g.Time != simTime || a.lair.Tokens != tokens {
				t.Fatal("end animation advanced simulation or recorded the run again")
			}
			a.handle(Event{Key: KeyEscape})
			if a.screen != ScreenOverworld {
				t.Fatal("result did not return to the lair")
			}
		})
	}
}
