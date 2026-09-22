package tui

import (
	"testing"

	"tdef/game"
	"tdef/hiscore"
	"tdef/render"
)

// owTestApp builds an App standing on the lair map, with an isolated lair file.
func owTestApp(t *testing.T) *App {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	return &App{
		screen: ScreenOverworld,
		diff:   game.Normal,
		lair:   hiscore.LoadLair(),
		ow:     render.NewOWState(),
	}
}

// The unseal chain must follow the lair's memory at the current renown:
// rift+rotunda open at first; the halls after the rift; the garden after the
// halls; the depths after two built-in floors; the heart after all four.
func TestOWUnlockChain(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	if !a.ow.Unlocked["rift"] || !a.ow.Unlocked["rotunda"] {
		t.Fatal("rift and rotunda must be open at first")
	}
	if a.ow.Unlocked["halls"] || a.ow.Unlocked["garden"] || a.ow.Unlocked["depths"] {
		t.Fatal("the far floors must be sealed at first")
	}
	a.lair.Record("rift", 1, 20, true)
	a.owRefresh()
	if !a.ow.Unlocked["halls"] {
		t.Fatal("the halls must open once the rift is held")
	}
	if a.ow.Unlocked["garden"] {
		t.Fatal("the garden must stay sealed until the halls are held")
	}
	a.lair.Record("halls", 1, 20, true)
	a.owRefresh()
	if !a.ow.Unlocked["garden"] {
		t.Fatal("the garden must open once the halls are held")
	}
	if !a.ow.Unlocked["depths"] {
		t.Fatal("the depths must open once two built-in floors are held")
	}
	if a.ow.BossReady {
		t.Fatal("the heart must stay sealed with two floors held")
	}
	a.lair.Record("rotunda", 1, 20, true)
	a.lair.Record("garden", 1, 20, true)
	a.owRefresh()
	if !a.ow.BossReady {
		t.Fatal("the heart must unseal once all four floors are held")
	}
	if a.ow.Hearts != 4 {
		t.Fatalf("hearts = %d, want 4", a.ow.Hearts)
	}
}

// Entering an open floor starts the descent; the transition plays out over
// OWDescendFrames and then launches the matching level from the lair.
func TestOWDescendLaunchesLevel(t *testing.T) {
	a := owTestApp(t)
	a.ow.Cursor = game.Vec{X: 6, Y: 10} // the Rift
	a.owEnter()
	if a.ow.Descending != "rift" {
		t.Fatalf("descending = %q, want rift", a.ow.Descending)
	}
	for i := 0; i < render.OWDescendFrames && a.screen == ScreenOverworld; i++ {
		a.owTick()
	}
	if a.screen != ScreenGame {
		t.Fatalf("after the descent = %v, want game", a.screen)
	}
	if a.level != "canyon" || !a.fromOW || a.owFloorID != "rift" {
		t.Fatalf("launched %q fromOW=%v floor=%q, want canyon/true/rift", a.level, a.fromOW, a.owFloorID)
	}
}

// A sealed floor refuses the descent with a message and does not launch.
func TestOWSealedRefusesDescent(t *testing.T) {
	a := owTestApp(t)
	a.ow.Cursor = game.Vec{X: 35, Y: 3} // the (sealed) Long Halls
	a.owEnter()
	if a.ow.Descending != "" {
		t.Fatalf("sealed floor descended: %q", a.ow.Descending)
	}
	if a.ow.Msg == "" {
		t.Fatal("sealed descent left no message")
	}
	if a.screen != ScreenOverworld {
		t.Fatalf("sealed descent left the map: %v", a.screen)
	}
}

// The wheel hops Grak between floor pads in the lair's depth order.
func TestOWWheelHopsFloors(t *testing.T) {
	a := owTestApp(t)
	a.ow.Cursor = game.Vec{X: 6, Y: 10} // the Rift (first in depth order)
	a.owWheel(false)                    // next
	if a.owCursorFloor() != "rotunda" {
		t.Fatalf("wheel down from rift = %q, want rotunda", a.owCursorFloor())
	}
	a.owWheel(true) // back up
	if a.owCursorFloor() != "rift" {
		t.Fatalf("wheel up from rotunda = %q, want rift", a.owCursorFloor())
	}
}

// One keypress moves exactly one cell. The terminal driver emits no
// key-release events, so a tap cannot be told from a hold — the walk must not
// keep striding on later frames, or a single tap would move many cells.
func TestOWWalkOneTapOneCell(t *testing.T) {
	a := owTestApp(t)
	start := a.ow.Cursor
	a.handle(Event{Key: KeyRight})
	if a.ow.Cursor.Man(start) != 1 {
		t.Fatalf("one tap moved %d cell(s): %v -> %v, want exactly 1", a.ow.Cursor.Man(start), start, a.ow.Cursor)
	}
	// Many more frames must not move Grak any further.
	for i := 0; i < 200; i++ {
		a.owTick()
	}
	if a.ow.Cursor != (game.Vec{X: start.X + 1, Y: start.Y}) {
		t.Fatalf("walk drifted after the tap: %v -> %v", start, a.ow.Cursor)
	}
}

// Spending a relic trades a token for a bonus on the next defense.
func TestOWSpendRelic(t *testing.T) {
	a := owTestApp(t)
	a.ow.Tokens = 1
	a.owSpendRelic(0)
	if a.ow.Tokens != 0 || a.ow.BonusGold != 60 {
		t.Fatalf("after spend: tokens=%d bonusGold=%d, want 0/60", a.ow.Tokens, a.ow.BonusGold)
	}
	// The spend is persisted to the lair.
	if a.lair.Tokens != 0 {
		t.Fatalf("lair tokens = %d, want 0", a.lair.Tokens)
	}
	// A second spend with no tokens is refused (and leaves no bonus).
	a.owSpendRelic(0)
	if a.ow.BonusGold != 60 {
		t.Fatalf("overspend changed the bonus: %d", a.ow.BonusGold)
	}
}

// A run that started from the lair returns to the lair on game over (esc),
// not to the level select.
func TestOWGameOverReturnsToLair(t *testing.T) {
	a := owTestApp(t)
	a.fromOW = true
	a.owFloorID = "rift"
	m, err := game.LoadLevel("canyon")
	if err != nil {
		t.Fatal(err)
	}
	a.enterGame(m, "canyon", game.Normal)
	if !a.ui.ToLair {
		t.Fatal("a lair run must be flagged ToLair")
	}
	a.g.Status = game.StatusDefeat
	a.handle(Event{Key: KeyEscape})
	if a.screen != ScreenOverworld {
		t.Fatalf("esc on a lair game-over = %v, want the lair", a.screen)
	}
}

// A run that did not come from the lair (a.ow is the zero value) must still
// fold its result into the lair's memory without touching the overworld's
// transition maps — the unseal cascade only runs for lair runs.
func TestDirectPlayRecordsLairWithoutOverworld(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m, err := game.LoadLevel("canyon")
	if err != nil {
		t.Fatal(err)
	}
	a := &App{
		screen: ScreenGame,
		g:      game.NewStateDiff(m, game.Normal),
		ui:     freshUI(m),
		level:  "canyon",
		diff:   game.Normal,
		fromOW: false,
		// a.ow is deliberately the zero value (nil maps), as in a direct-play
		// session that skipped the lair.
		lair: hiscore.LoadLair(),
	}
	a.g.Status = game.StatusVictory
	a.g.Wave = game.MaxWaves
	// Must not panic: the floor clears (unsealing the halls) but the
	// overworld's transition maps are never written.
	a.stepGame(1.0 / 20.0)
	if !a.lair.Floor("rift", diffIndex(game.Normal)).Cleared {
		t.Fatal("a direct-play win did not record the floor in the lair")
	}
}

// A run that started from the level select returns to the level select on
// game over, never to the lair.
func TestLevelSelectGameOverReturnsToSelect(t *testing.T) {
	a := owTestApp(t)
	a.fromOW = false
	a.ls.Levels = game.LevelNames()
	m, err := game.LoadLevel("canyon")
	if err != nil {
		t.Fatal(err)
	}
	a.enterGame(m, "canyon", game.Normal)
	if a.ui.ToLair {
		t.Fatal("a level-select run must not be flagged ToLair")
	}
	a.g.Status = game.StatusDefeat
	a.handle(Event{Key: KeyEscape})
	if a.screen != ScreenLevelSelect {
		t.Fatalf("esc on a menu game-over = %v, want the level select", a.screen)
	}
}
