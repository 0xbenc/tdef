package tui

import (
	"math"
	"strings"
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
)

func TestGrakIntroductionPersistsUntilSuccessfulInput(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	if !a.ow.PlayerIntro {
		t.Fatal("new player has no orientation marker")
	}
	a.owStep(100, 100)
	if !a.ow.PlayerIntro || a.lair.GrakLocated {
		t.Fatal("blocked movement dismissed orientation")
	}
	a.handleOverworld(Event{Key: KeyRight})
	if a.ow.PlayerIntro || !hiscore.LoadLair().GrakLocated {
		t.Fatal("successful movement did not remember orientation")
	}
	a.owRefresh()
	if a.ow.PlayerIntro {
		t.Fatal("screen refresh replayed orientation")
	}
	restarted := newApp(nil, game.Normal)
	restarted.ow = render.NewOWState()
	restarted.owRefresh()
	if restarted.ow.PlayerIntro {
		t.Fatal("restart replayed orientation")
	}
}

func TestUnlockRevealWaitsForVictoryAndDoesNotMoveGrak(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	cursor := a.ow.Cursor
	a.lair.Record("rotunda", 1, game.MaxWaves, true)
	a.owUnsealCheck(1)
	a.g = &game.State{Wave: game.MaxWaves}
	a.owFloorID = "rotunda"
	a.owSetBanner(true, 0, false)
	a.owRefresh() // Returning from the battlefield must preserve the pending reveal.
	for i := 0; i < render.OWVictoryBeatFrames-1; i++ {
		a.owTick()
		if a.ow.RevealFloor != "" || a.ow.Unsealing["rift"] != render.OWUnsealFrames {
			t.Fatal("new area stole the victory beat")
		}
	}
	a.owTick()
	if a.ow.RevealFloor != "rift" {
		t.Fatal("rift did not get a dedicated reveal")
	}
	for _, event := range []Event{{Key: KeyRight}, {Key: KeyTab}, {Key: KeyEnter}, {Mouse: true, Press: true, Btn: 65}} {
		a.handleOverworld(event)
	}
	if a.ow.Cursor != cursor || a.ow.Diff != 1 || a.ow.Descending != "" {
		t.Fatal("input interrupted reveal")
	}
	for i := 0; i < render.OWRevealPanFrames; i++ {
		a.owTick()
	}
	if math.Abs(a.ow.CameraX-6) > 0.1 {
		t.Fatalf("camera never focused the Rift: %f", a.ow.CameraX)
	}
	if a.ow.Cursor != cursor {
		t.Fatal("camera reveal teleported Grak")
	}
	for i := render.OWRevealPanFrames; i < render.OWRevealFrames; i++ {
		a.owTick()
	}
	if a.ow.RevealBusy() || !render.OWFloorOpen("rift", a.ow) || a.ow.Cursor != cursor {
		t.Fatal("reveal did not finish cleanly")
	}
	if !strings.Contains(a.ow.ReturnMsg, "Rift") {
		t.Fatal("new destination not announced")
	}
	for i := 0; i < 45; i++ {
		a.owTick()
	}
	if math.Abs(a.ow.CameraX-float64(cursor.X)) > 0.02 {
		t.Fatal("camera did not return to Grak")
	}
}

func TestUnlocksRevealSequentiallyAndRenownClearsThem(t *testing.T) {
	a := owTestApp(t)
	a.owQueueReveal("rift")
	a.owQueueReveal("halls")
	a.owQueueReveal("rift")
	if len(a.ow.RevealQueue) != 2 {
		t.Fatal("duplicate reveal queued")
	}
	a.owTick()
	for i := 0; i < render.OWRevealFrames; i++ {
		if a.ow.Unsealing["halls"] != render.OWUnsealFrames {
			t.Fatal("offscreen room animated before its turn")
		}
		a.owTick()
	}
	a.owTick()
	if a.ow.RevealFloor != "halls" {
		t.Fatal("second room did not get its own reveal")
	}
	a.owCycleDiff(1)
	if a.ow.RevealBusy() || len(a.ow.Unsealing) != 0 || a.ow.BlastTTL != 0 {
		t.Fatal("renown inherited old reveals")
	}
	for i := 0; i < render.OWRevealFrames+1; i++ {
		a.owTick()
	}
	if a.ow.Unlocked["halls"] {
		t.Fatal("old reveal opened a room at another difficulty")
	}
}

func TestUnlockRevealPausesWhenTerminalTooSmall(t *testing.T) {
	a := owTestApp(t)
	a.owQueueReveal("garden")
	a.owTick()
	ttl := a.ow.RevealTTL
	a.term = &Terminal{size: [2]int{40, 15}}
	for i := 0; i < 200; i++ {
		a.owTick()
	}
	if a.ow.RevealTTL != ttl {
		t.Fatal("unseen reveal expired while terminal was too small")
	}
	a.term.size = [2]int{80, 24}
	a.owTick()
	if a.ow.RevealTTL != ttl-1 {
		t.Fatal("reveal did not resume after resize")
	}
}

func TestOldUnlockedAreasDoNotReplayReveals(t *testing.T) {
	a := owTestApp(t)
	for _, id := range hiscore.LairFloors {
		a.lair.Record(id, 1, game.MaxWaves, true)
	}
	a.lair.Record(hiscore.HeartFloor, 1, game.MaxWaves, true)
	a.owRefresh()
	a.owUnsealCheck(1)
	a.owHeartBlastCheck(1)
	if a.ow.RevealBusy() {
		t.Fatal("existing campaign replayed its unlocks")
	}
}
