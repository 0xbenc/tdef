package tui

import (
	"testing"

	"github.com/0xbenc/termtd/hiscore"
)

func TestMenuRevealWaitsForFirstCampaignWaveAndPersists(t *testing.T) {
	a := journalGameApp(t)
	a.lair.Training.Stage = 1
	a.ui.Paused = true
	a.stepGame(0)
	if a.lair.HasMenuReveal() {
		t.Fatal("menu revealed before playing")
	}
	a.g.Wave, a.g.WaveActive = 1, true
	a.stepGame(0)
	if a.lair.HasMenuReveal() {
		t.Fatal("starting a wave revealed the menu")
	}
	a.g.WaveActive = false
	a.fromOW = false
	a.stepGame(0)
	if a.lair.HasMenuReveal() {
		t.Fatal("Quick Play changed the campaign reveal")
	}
	a.fromOW = true
	a.stepGame(0)
	if !a.lair.MenuRevealed || !hiscore.LoadLair().HasMenuReveal() {
		t.Fatal("played wave did not persist the reveal")
	}
	a.resetSel = 1
	a.activateResetProgress()
	if a.lair.HasMenuReveal() {
		t.Fatal("reset kept the reveal")
	}
}

func TestExistingCampaignRetainsMenuReveal(t *testing.T) {
	l := &hiscore.Lair{Floors: map[string]hiscore.FloorRec{"rotunda:1": {BestWave: 1}}}
	if !l.HasMenuReveal() {
		t.Fatal("existing player's dragon disappeared")
	}
}
