package tui

import (
	"testing"

	"github.com/0xbenc/tdef/game"
	"github.com/0xbenc/tdef/hiscore"
	"github.com/0xbenc/tdef/render"
)

// lsApp builds an App on the level-select screen with the built-in level
// list.
func lsApp() *App {
	return &App{
		screen: ScreenLevelSelect,
		ls:     render.LSState{Levels: game.LevelNames()},
		scores: map[string]int{},
	}
}

func TestMenuNavWraps(t *testing.T) {
	a := &App{screen: ScreenMenu}
	a.handle(Event{Key: KeyUp}) // 0 wraps to last
	if a.menuSel != len(render.MenuItems)-1 {
		t.Fatalf("up from 0 = %d, want %d", a.menuSel, len(render.MenuItems)-1)
	}
	a.handle(Event{Key: KeyDown}) // wraps back to 0
	if a.menuSel != 0 {
		t.Fatalf("down from last = %d, want 0", a.menuSel)
	}
	a.handle(Event{Rune: 's'})
	if a.menuSel != 1 {
		t.Errorf("s = %d, want 1", a.menuSel)
	}
	a.handle(Event{Rune: 'w'})
	if a.menuSel != 0 {
		t.Errorf("w = %d, want 0", a.menuSel)
	}
}

func TestMenuTransitions(t *testing.T) {
	t.Setenv("HOME", t.TempDir()) // the lair's memory is read on Start
	a := &App{screen: ScreenMenu, scores: map[string]int{"canyon": 1}}
	// Each case re-enters the menu, since activating an item leaves it.
	for sel, want := range map[int]Screen{
		0: ScreenOverworld,
		1: ScreenLevelSelect,
		2: ScreenHelp,
		3: ScreenHiscores,
	} {
		a.screen = ScreenMenu
		a.menuSel = sel
		a.handle(Event{Key: KeyEnter})
		if a.screen != want {
			t.Errorf("enter at item %d = %v, want %v", sel, a.screen, want)
		}
	}
	a.screen = ScreenMenu
	a.menuSel = 4
	a.handle(Event{Key: KeyEnter})
	if !a.quitting {
		t.Error("enter at Quit did not quit")
	}
	// esc goes back to the title screen.
	b := &App{screen: ScreenMenu}
	b.handle(Event{Key: KeyEscape})
	if b.screen != ScreenTitle {
		t.Errorf("esc in menu = %v, want title", b.screen)
	}
}

// A click on a menu item must activate it, using the same geometry the
// renderer draws with.
func TestMenuClickActivates(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	a := &App{screen: ScreenMenu}
	w, h := a.termSize()
	rects := render.MenuRects(w, h)
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: rects[2].X + 1, Y: rects[2].Y})
	if a.screen != ScreenHelp {
		t.Errorf("click on Help row = %v, want help", a.screen)
	}
	// Wheel scrolls the selection.
	b := &App{screen: ScreenMenu}
	b.handle(Event{Mouse: true, Btn: 64, Press: true})
	if b.menuSel != len(render.MenuItems)-1 {
		t.Errorf("wheel up from 0 = %d, want last", b.menuSel)
	}
}

func TestTitleAnyKeyToMenu(t *testing.T) {
	for _, e := range []Event{{Rune: 'z'}, {Key: KeyEnter}, {Key: KeyEscape}} {
		a := &App{screen: ScreenTitle}
		a.handle(e)
		if a.screen != ScreenMenu {
			t.Errorf("title + %+v = %v, want menu", e, a.screen)
		}
	}
	a := &App{screen: ScreenTitle}
	a.handle(Event{Mouse: true, Btn: 0, Press: false})
	if a.screen != ScreenTitle {
		t.Error("mouse release must not advance the title")
	}
	a.handle(Event{Mouse: true, Btn: 0, Press: true})
	if a.screen != ScreenMenu {
		t.Errorf("mouse press = %v, want menu", a.screen)
	}
	q := &App{screen: ScreenTitle}
	q.handle(Event{Rune: 'q'})
	if !q.quitting {
		t.Error("q on the title did not quit")
	}
}

func TestHelpAndHiscoresReturnToMenu(t *testing.T) {
	h := &App{screen: ScreenHelp}
	h.handle(Event{Key: KeyEnter})
	if h.screen != ScreenMenu {
		t.Errorf("help + enter = %v, want menu", h.screen)
	}
	s := &App{screen: ScreenHiscores, scores: map[string]int{"a": 1, "b": 2, "c": 3}}
	s.handle(Event{Key: KeyDown})
	s.handle(Event{Key: KeyDown})
	if s.hsTop != 2 {
		t.Errorf("two downs = hsTop %d, want 2", s.hsTop)
	}
	s.handle(Event{Key: KeyDown})
	if s.hsTop != 2 {
		t.Errorf("hsTop scrolled past the last entry (%d)", s.hsTop)
	}
	s.handle(Event{Key: KeyUp})
	if s.hsTop != 1 {
		t.Errorf("up = hsTop %d, want 1", s.hsTop)
	}
	s.handle(Event{Key: KeyEscape})
	if s.screen != ScreenMenu {
		t.Errorf("hiscores + esc = %v, want menu", s.screen)
	}
}

func TestLevelSelectEdits(t *testing.T) {
	a := lsApp()
	// Digits are ignored off the maze row.
	a.handle(Event{Rune: '5'})
	if a.ls.Seed != "" {
		t.Errorf("digit off the maze row set seed %q", a.ls.Seed)
	}
	// Walk to the maze row (last).
	for i := 0; i < len(a.ls.Levels); i++ {
		a.handle(Event{Key: KeyDown})
	}
	if a.ls.Cursor != len(a.ls.Levels) {
		t.Fatalf("cursor = %d, want maze row %d", a.ls.Cursor, len(a.ls.Levels))
	}
	a.handle(Event{Rune: '1'})
	a.handle(Event{Rune: '2'})
	a.handle(Event{Rune: '3'})
	if a.ls.Seed != "123" {
		t.Errorf("seed = %q, want 123", a.ls.Seed)
	}
	a.handle(Event{Key: KeyBackspace})
	if a.ls.Seed != "12" {
		t.Errorf("seed after backspace = %q, want 12", a.ls.Seed)
	}
	a.handle(Event{Rune: 'r'})
	if a.ls.Seed != "" {
		t.Errorf("seed after r = %q, want empty", a.ls.Seed)
	}
	// Seed editing must not wrap the cursor.
	a.handle(Event{Key: KeyDown})
	if a.ls.Cursor != 0 {
		t.Errorf("down on the maze row wrapped to %d, want 0", a.ls.Cursor)
	}
	// Difficulty cycles with left/right and wraps.
	a.ls.Cursor = len(a.ls.Levels)
	a.handle(Event{Key: KeyRight})
	if a.ls.Diff != 1 {
		t.Errorf("right = diff %d, want 1", a.ls.Diff)
	}
	a.handle(Event{Key: KeyLeft})
	if a.ls.Diff != 0 {
		t.Errorf("left = diff %d, want 0", a.ls.Diff)
	}
	a.handle(Event{Key: KeyLeft})
	if a.ls.Diff != len(render.Difficulties)-1 {
		t.Errorf("left wrap = diff %d, want last", a.ls.Diff)
	}
	// Esc goes back to the menu.
	a.handle(Event{Key: KeyEscape})
	if a.screen != ScreenMenu {
		t.Errorf("esc = %v, want menu", a.screen)
	}
}

func TestLevelSelectMouse(t *testing.T) {
	a := lsApp()
	w, h := a.termSize()
	rows, diffs, seedRect := render.LSRects(a.lsView(), w, h)
	// Click a level row.
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: rows[2].X + 1, Y: rows[2].Y})
	if a.ls.Cursor != 2 {
		t.Errorf("row click = cursor %d, want 2", a.ls.Cursor)
	}
	// Click the seed row lands on the maze row.
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: seedRect.X + 1, Y: seedRect.Y})
	if a.ls.Cursor != len(a.ls.Levels) {
		t.Errorf("seed click = cursor %d, want maze row", a.ls.Cursor)
	}
	// Click a difficulty label.
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: diffs[2].X + 1, Y: diffs[2].Y})
	if a.ls.Diff != 2 {
		t.Errorf("difficulty click = %d, want 2", a.ls.Diff)
	}
}

func TestLevelSelectStartBuiltin(t *testing.T) {
	a := lsApp()
	a.ls.Diff = 2
	a.ls.Cursor = 2 // hub
	a.handle(Event{Key: KeyEnter})
	if a.screen != ScreenGame {
		t.Fatalf("enter = %v, want game", a.screen)
	}
	if a.level != "hub" || a.diff != game.Hard {
		t.Errorf("started %s on %v, want hub on hard", a.level, a.diff)
	}
	want, err := game.LoadLevel("hub")
	if err != nil {
		t.Fatal(err)
	}
	if a.g == nil || a.g.Map.W != want.W || a.g.Map.H != want.H {
		t.Errorf("game map %vx%v, want %dx%d", a.g.Map.W, a.g.Map.H, want.W, want.H)
	}
	if a.ui.Selected != render.NoSelection || a.ui.Paused {
		t.Errorf("game must start unpaused with no selection: %+v", a.ui)
	}
}

func TestLevelSelectStartMaze(t *testing.T) {
	a := lsApp()
	a.ls.Cursor = len(a.ls.Levels)
	a.ls.Seed = "1234"
	a.handle(Event{Key: KeyEnter})
	if a.screen != ScreenGame {
		t.Fatalf("enter = %v, want game", a.screen)
	}
	if a.level != "maze1234" {
		t.Errorf("level = %q, want maze1234", a.level)
	}
	want, err := game.MazeFromSeed(1234)
	if err != nil {
		t.Fatal(err)
	}
	if a.g.Map.Spawn != want.Spawn || a.g.Map.Exit != want.Exit {
		t.Errorf("maze spawn/exit = %v/%v, want %v/%v", a.g.Map.Spawn, a.g.Map.Exit, want.Spawn, want.Exit)
	}
	// An unparseable (too long) seed shows an error instead of starting.
	b := lsApp()
	b.ls.Cursor = len(b.ls.Levels)
	b.ls.Seed = "99999999999999999999"
	b.handle(Event{Key: KeyEnter})
	if b.screen != ScreenLevelSelect {
		t.Errorf("bad seed = %v, want stay on level select", b.screen)
	}
	if b.ls.Err == "" {
		t.Error("bad seed did not set an error")
	}
}

// The full pre-game flow must be reachable: title -> menu -> the lair ->
// descend into a floor -> game, ending in a playable state.
func TestScreenSequenceReachable(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	a := &App{
		screen: ScreenTitle,
		ls:     render.LSState{Levels: game.LevelNames()},
		scores: map[string]int{},
		lair:   hiscore.LoadLair(),
		ow:     render.NewOWState(),
	}
	a.handle(Event{Rune: ' '})
	if a.screen != ScreenMenu {
		t.Fatalf("title -> %v, want menu", a.screen)
	}
	a.handle(Event{Key: KeyEnter}) // Start: the lair
	if a.screen != ScreenOverworld {
		t.Fatalf("menu -> %v, want the lair", a.screen)
	}
	// Grak starts on the Rotunda: wait for the lair to wake (the boot gates
	// input), then descend and let the transition play out.
	for a.ow.BootTTL > 0 {
		a.owTick()
	}
	a.handle(Event{Key: KeyEnter})
	if a.ow.Descending != "rotunda" {
		t.Fatalf("enter on the Rotunda = descending %q, want rotunda", a.ow.Descending)
	}
	for i := 0; i < render.OWDescendFrames && a.screen == ScreenOverworld; i++ {
		a.owTick()
	}
	if a.screen != ScreenGame {
		t.Fatalf("descent -> %v, want game", a.screen)
	}
	if a.g == nil || a.level != "hub" || !a.fromOW {
		t.Fatalf("game state after the sequence: level=%q fromOW=%v", a.level, a.fromOW)
	}
}

func TestParseSeed(t *testing.T) {
	cases := []struct {
		seed string
		ok   bool
		want int64
	}{
		{"", true, 0},        // random (nonzero)
		{"0", true, 0},       // random (nonzero)
		{"1234", true, 1234}, // exact
		{"99999999999999999999", false, 0},
	}
	for _, c := range cases {
		a := &App{}
		a.ls.Seed = c.seed
		got, ok := a.parseSeed()
		if ok != c.ok {
			t.Errorf("parseSeed(%q) ok = %v, want %v", c.seed, ok, c.ok)
			continue
		}
		if !ok {
			continue
		}
		if c.want != 0 && got != c.want {
			t.Errorf("parseSeed(%q) = %d, want %d", c.seed, got, c.want)
		}
		if c.want == 0 && got == 0 {
			t.Errorf("parseSeed(%q) = 0, want a random nonzero seed", c.seed)
		}
	}
}

// lsView must cache the preview: repeated calls return the same map without
// regenerating it, and a random-seed preview uses a fixed seed.
func TestLSViewCachesPreview(t *testing.T) {
	a := lsApp()
	a.ls.Cursor = len(a.ls.Levels) // maze, empty seed
	v1 := a.lsView()
	v2 := a.lsView()
	if v1.Preview == nil {
		t.Fatal("random-seed maze has no preview")
	}
	if v1.Preview != v2.Preview {
		t.Error("preview not cached")
	}
	want, err := game.MazeFromSeed(1234)
	if err != nil {
		t.Fatal(err)
	}
	if v1.Preview.Spawn != want.Spawn || v1.Preview.Exit != want.Exit {
		t.Error("random-seed preview does not use the fixed seed")
	}
	// Moving the cursor invalidates the cache.
	a.ls.Cursor = 0
	v3 := a.lsView()
	if v3.Preview == v1.Preview {
		t.Error("preview not invalidated on cursor change")
	}
}
