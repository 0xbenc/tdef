package tui

import (
	"testing"

	"github.com/0xbenc/tdef/game"
	"github.com/0xbenc/tdef/hiscore"
	"github.com/0xbenc/tdef/render"
)

// owTestApp builds an App standing on the lair map, with an isolated lair file.
func owTestApp(t *testing.T) *App {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return &App{
		screen: ScreenOverworld,
		diff:   game.Normal,
		lair:   hiscore.LoadLair(),
		ow:     render.NewOWState(),
	}
}

// The unseal chain must follow the lair's memory at the current renown:
// rotunda open at first; rift after rotunda; halls after rift; garden after the
// halls; the heart after all four; the depths after the Heart is also held.
func TestOWUnlockChain(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	if !a.ow.Unlocked["rotunda"] {
		t.Fatal("rotunda must be open at first")
	}
	if a.ow.Unlocked["rift"] || a.ow.Unlocked["halls"] || a.ow.Unlocked["garden"] || a.ow.Unlocked["depths"] {
		t.Fatal("the far floors must be sealed at first")
	}
	a.lair.Record("rotunda", 1, 20, true)
	a.owRefresh()
	if !a.ow.Unlocked["rift"] || a.ow.Unlocked["halls"] || a.ow.Unlocked["depths"] {
		t.Fatal("holding the rotunda should open only the rift")
	}
	a.lair.Record("rift", 1, 20, true)
	a.owRefresh()
	if !a.ow.Unlocked["halls"] {
		t.Fatal("the halls must open once the rift is held")
	}
	if a.ow.Unlocked["garden"] {
		t.Fatal("the garden must stay sealed until the halls are held")
	}
	if a.ow.Unlocked["depths"] || a.ow.BossReady {
		t.Fatal("two held floors must leave the depths and heart sealed")
	}
	a.lair.Record("halls", 1, 20, true)
	a.owRefresh()
	if !a.ow.Unlocked["garden"] {
		t.Fatal("the garden must open once the halls are held")
	}
	if a.ow.Unlocked["depths"] {
		t.Fatal("the depths must stay sealed with three floors held")
	}
	if a.ow.BossReady {
		t.Fatal("the heart must stay sealed with three floors held")
	}
	a.lair.Record("garden", 1, 20, true)
	a.owRefresh()
	if !a.ow.BossReady {
		t.Fatal("the heart must unseal once all four floors are held")
	}
	if a.ow.Hearts != 4 {
		t.Fatalf("hearts = %d, want 4", a.ow.Hearts)
	}
	if a.ow.Unlocked["depths"] {
		t.Fatal("the depths must stay sealed until the Heart is beaten")
	}
	a.lair.Record(hiscore.HeartFloor, 1, 20, true)
	a.owUnsealCheck(1)
	if a.ow.Unsealing["depths"] == 0 {
		t.Fatal("the Heart victory did not start the depths unseal")
	}
	a.owRefresh()
	if !a.ow.Unlocked["depths"] {
		t.Fatal("all fixed defenses held must open the depths")
	}
	a.ow.Diff = 0
	a.owRefresh()
	if a.ow.Unlocked["depths"] {
		t.Fatal("normal victories unlocked depths on easy")
	}
}

func TestOWHubWinUnsealsRift(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	a.lair.Record("rotunda", 1, 8, false)
	a.owUnsealCheck(1)
	if a.ow.Unlocked["rift"] || a.ow.Unsealing["rift"] > 0 {
		t.Fatal("losing the hub must not open the rift")
	}
	a.lair.Record("rotunda", 1, 20, true)
	a.owUnsealCheck(1)
	if a.ow.Unsealing["rift"] != render.OWUnsealFrames || a.ow.Unlocked["rift"] {
		t.Fatal("holding the hub should start the rift's unseal transition")
	}
	for i := 0; i < render.OWUnsealFrames-1; i++ {
		a.owTick()
	}
	if a.ow.Unlocked["rift"] {
		t.Fatal("rift opened before its transition finished")
	}
	a.owTick()
	if !a.ow.Unlocked["rift"] || a.ow.Unlocked["halls"] {
		t.Fatal("the rift alone should open after the hub win")
	}
	a.owVisitFloor("rift")
	a.owEnter()
	if a.ow.Descending != "rift" {
		t.Fatal("newly opened rift cannot be entered")
	}
}

// Entering an open floor starts the descent; the transition plays out over
// OWDescendFrames and then launches the matching level from the lair.
func TestOWDescendLaunchesLevel(t *testing.T) {
	a := owTestApp(t)
	a.ow.Cursor = game.Vec{X: 6, Y: 10} // the Rift
	a.ow.Unlocked["rift"] = true
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
	a.ow.Cursor = game.Vec{X: 63, Y: 3} // the (sealed) Long Halls
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
	a.ow.Unlocked["rift"] = true
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

// While the lair is waking (BootTTL > 0) Grak cannot walk: the arrival
// cinematic owns the screen. Nav keys (esc, q) still pass.
func TestOWInputGateDuringBoot(t *testing.T) {
	a := owTestApp(t)
	a.ow.BootTTL = render.OWBootFrames
	a.handleOverworld(Event{Rune: 'w'})
	a.handleOverworld(Event{Key: KeyRight})
	if a.ow.Cursor != (game.Vec{X: 22, Y: 6}) {
		t.Fatalf("cursor moved during boot: %v", a.ow.Cursor)
	}
	if len(a.ow.Trail) != 0 {
		t.Fatalf("trail left during boot: %v", a.ow.Trail)
	}
	a.handleOverworld(Event{Key: KeyEscape})
	if a.screen != ScreenTitle {
		t.Fatal("esc must pass during boot")
	}
}

// The arrival cinematic is armed exactly once per session: the first entry
// into the lair starts it, and later entries (returns from a defense) do not
// replay it.
func TestOWBootArmedOnce(t *testing.T) {
	a := owTestApp(t)
	a.lair.IntroSeen = true // Returning players retain the short lair arrival.
	hiscore.SaveLair(a.lair)
	a.toScreen(ScreenOverworld)
	if a.ow.BootTTL != render.OWBootFrames {
		t.Fatalf("first entry BootTTL = %d, want %d", a.ow.BootTTL, render.OWBootFrames)
	}
	for i := 0; i < 50; i++ {
		a.owTick()
	}
	if got := a.ow.BootTTL; got != render.OWBootFrames-50 {
		t.Fatalf("BootTTL after 50 ticks = %d, want %d", got, render.OWBootFrames-50)
	}
	a.toScreen(ScreenTitle)
	a.toScreen(ScreenOverworld)
	if a.ow.BootTTL != render.OWBootFrames-50 {
		t.Fatalf("re-entry re-armed the boot: BootTTL = %d, want %d", a.ow.BootTTL, render.OWBootFrames-50)
	}
}

// The heart-unseal shockwave arms exactly once, when a defense unseals the
// heart at the current renown; an already-unsealed heart never re-fires it.
func TestOWHeartBlastCheck(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh() // no floors held: the heart is sealed, marked unseen
	a.lair.Record("rift", 1, 20, true)
	a.lair.Record("halls", 1, 20, true)
	a.lair.Record("rotunda", 1, 20, true)
	a.owHeartBlastCheck(1) // three floors: still sealed
	if a.ow.BlastTTL != 0 {
		t.Fatalf("blast armed with three floors held: %d", a.ow.BlastTTL)
	}
	a.lair.Record("garden", 1, 20, true) // the fourth: the heart unseals
	a.owHeartBlastCheck(1)
	if a.ow.BlastTTL != render.OWBlastFrames {
		t.Fatalf("blast not armed on the unseal: %d, want %d", a.ow.BlastTTL, render.OWBlastFrames)
	}
	a.owHeartBlastCheck(1) // a later check must not re-arm it
	if a.ow.BlastTTL != render.OWBlastFrames {
		t.Fatalf("blast re-armed: %d", a.ow.BlastTTL)
	}
	// A renown whose heart was already unsealed at entry never fires it.
	a2 := owTestApp(t)
	a2.ow.Diff = 2
	a2.lair.Record("rift", 2, 20, true)
	a2.lair.Record("halls", 2, 20, true)
	a2.lair.Record("rotunda", 2, 20, true)
	a2.lair.Record("garden", 2, 20, true)
	a2.owRefresh() // marks renown 2 as seen with the heart unsealed
	a2.owHeartBlastCheck(2)
	if a2.ow.BlastTTL != 0 {
		t.Fatalf("blast fired for an already-unsealed heart: %d", a2.ow.BlastTTL)
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
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
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

func TestOWCameraFollowsWithoutMovingPlayer(t *testing.T) {
	a := owTestApp(t)
	a.owTick() // establish the camera at the Rift
	a.ow.Cursor = game.Vec{X: 45, Y: 6}
	before := a.ow.CameraX
	a.owTick()
	if a.ow.CameraX <= before || a.ow.CameraX >= 45 {
		t.Fatalf("camera did not ease toward Grak: %f -> %f", before, a.ow.CameraX)
	}
	for i := 0; i < 40; i++ {
		a.owTick()
	}
	if a.ow.CameraX != 45 || a.ow.Cursor != (game.Vec{X: 45, Y: 6}) {
		t.Fatalf("camera did not settle or moved player: camera=%f cursor=%v", a.ow.CameraX, a.ow.Cursor)
	}
	// Browser shortcuts reveal their destination immediately.
	a.owWheel(false)
	if a.ow.CameraX != float64(a.ow.Cursor.X) {
		t.Fatal("wheel jump left the camera at the previous floor")
	}
}

func TestOWMouseAfterScrolling(t *testing.T) {
	a := owTestApp(t)
	a.ow.Unlocked["garden"] = true
	a.ow.Cursor = game.Vec{X: 63, Y: 3}
	a.ow.CameraSet, a.ow.CameraX = true, 60.5
	w, h := a.termSize()
	l := render.OverworldLayout(w, h, a.ow)
	fl, _ := render.OWFloorOf("garden")
	x := l.X(fl.Center.X) + l.Scale/2
	y := l.Y(fl.Center.Y) + l.Scale/2
	a.handleOWMouse(Event{Mouse: true, Press: true, X: x, Y: y})
	if a.ow.Cursor != fl.Center {
		t.Fatalf("scrolled click went to %v, want %v", a.ow.Cursor, fl.Center)
	}
}

func TestOWWalkingFollowsTheFork(t *testing.T) {
	a := owTestApp(t)
	a.ow.Cursor = game.Vec{X: 45, Y: 6}
	a.handleOverworld(Event{Key: KeyRight})
	if a.ow.Cursor != (game.Vec{X: 45, Y: 6}) {
		t.Fatal("walked off the bridge along the old straight route")
	}
	a.handleOverworld(Event{Key: KeyUp})
	a.handleOverworld(Event{Key: KeyUp})
	a.handleOverworld(Event{Key: KeyRight})
	if a.ow.Cursor != (game.Vec{X: 46, Y: 4}) {
		t.Fatalf("upper approach ended at %v", a.ow.Cursor)
	}
	a.ow.Cursor = game.Vec{X: 45, Y: 6}
	a.handleOverworld(Event{Key: KeyDown})
	a.handleOverworld(Event{Key: KeyDown})
	a.handleOverworld(Event{Key: KeyRight})
	if a.ow.Cursor != (game.Vec{X: 46, Y: 8}) {
		t.Fatalf("lower approach ended at %v", a.ow.Cursor)
	}
}

func TestOWWalkingStopsAtSealedDoorway(t *testing.T) {
	for _, id := range []string{"rift", "halls", "garden", "depths"} {
		a := owTestApp(t)
		fl, _ := render.OWFloorOf(id)
		approach, _ := render.OWFloorApproach(id)
		a.ow.Cursor = approach
		dx, dy := 0, 0
		if fl.Center.X > approach.X {
			dx = 1
		}
		if fl.Center.X < approach.X {
			dx = -1
		}
		if fl.Center.Y > approach.Y {
			dy = 1
		}
		if fl.Center.Y < approach.Y {
			dy = -1
		}
		a.owWalk(dx, dy)
		if a.ow.Cursor != approach || len(a.ow.Trail) != 0 || a.ow.Msg == "" {
			t.Fatalf("%s: crossed a sealed doorway or failed to explain it", id)
		}
		a.owEnter()
		if a.ow.Descending != "" {
			t.Fatalf("%s: descended from a sealed doorstep", id)
		}
		a.ow.Unlocked[id], a.ow.Unsealing[id] = true, 1
		a.owWalk(dx, dy)
		if a.ow.Cursor != approach {
			t.Fatalf("%s: entered mid-unseal", id)
		}
		a.owTick()
		a.owWalk(dx, dy)
		if a.owCursorFloor() != id {
			t.Fatalf("%s: could not enter after unsealing", id)
		}
	}
}

func TestOWShortcutsRespectSeals(t *testing.T) {
	a := owTestApp(t)
	for _, id := range []string{"halls", "garden", "depths", "rift", "rotunda"} {
		a.owWheel(false)
		if a.owBrowseFloor() != id {
			t.Fatalf("wheel stopped at %q, want %q", a.owBrowseFloor(), id)
		}
		if !a.ow.Unlocked[id] {
			approach, _ := render.OWFloorApproach(id)
			if a.ow.Cursor != approach || a.owCursorFloor() != "" {
				t.Fatalf("wheel entered sealed %s", id)
			}
		}
	}
	for _, id := range []string{"rift", "halls", "garden", "depths"} {
		fl, _ := render.OWFloorOf(id)
		a.ow.CameraSet, a.ow.CameraX = true, float64(fl.Center.X)
		w, h := a.termSize()
		l := render.OverworldLayout(w, h, a.ow)
		x, y := l.X(fl.Center.X)+l.Scale/2, l.Y(fl.Center.Y)+l.Scale/2
		a.handleOWMouse(Event{Mouse: true, Press: true, X: x, Y: y})
		approach, _ := render.OWFloorApproach(id)
		if a.ow.Cursor != approach {
			t.Fatalf("mouse entered sealed %s", id)
		}
	}
}

func TestOWRenownChangeReturnsPlayerToDoorway(t *testing.T) {
	a := owTestApp(t)
	a.lair.Record("rift", 1, 20, true)
	a.owRefresh()
	fl, _ := render.OWFloorOf("halls")
	a.ow.Cursor = fl.Center
	a.ow.Unsealing["garden"] = 1
	a.owCycleDiff(1) // hard has no cleared floors
	approach, _ := render.OWFloorApproach("halls")
	if a.ow.Cursor != approach {
		t.Fatal("renown change left Grak inside a sealed room")
	}
	a.owTick()
	if a.ow.Unlocked["garden"] {
		t.Fatal("old renown's unseal opened a hard floor")
	}
}

func TestDepthsCannotBypassCampaign(t *testing.T) {
	a := owTestApp(t)
	for _, floor := range hiscore.LairFloors {
		a.lair.Record(floor, 1, game.MaxWaves, true)
	}
	a.owRefresh()
	a.owUnsealCheck(1)
	if a.ow.Unlocked["depths"] || a.ow.Unsealing["depths"] > 0 {
		t.Fatal("four floors without Heart opened depths")
	}
	// A pending/direct launch cannot spend relics or create a procedural run.
	a.ow.BonusGold = 123
	a.ow.Descending, a.ow.DescendTTL = "depths", 1
	a.owLaunch("depths")
	if a.screen != ScreenOverworld || a.g != nil || a.ow.BonusGold != 123 || a.ow.Descending != "" || a.ow.Msg != depthsLockedMessage {
		t.Fatal("locked depths launch bypassed progression or spent relics")
	}
	b := lsApp()
	b.lair = a.lair
	b.ls.Diff, b.ls.Cursor = 1, len(b.ls.Levels)
	if b.lsView().Err != depthsLockedMessage {
		t.Fatal("Quick Play preview omitted lock requirement")
	}
	b.handle(Event{Key: KeyEnter})
	if b.screen != ScreenLevelSelect || b.g != nil || b.ls.Err != depthsLockedMessage {
		t.Fatal("Quick Play bypassed depths lock")
	}
	// The command-line play entry must refuse the maze before opening a TTY.
	if err := Run(nil, "maze1234", game.Normal); err == nil || err.Error() != depthsLockedMessage {
		t.Fatalf("direct play bypassed depths lock: %v", err)
	}
}
