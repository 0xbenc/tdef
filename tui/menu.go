package tui

import "github.com/0xbenc/termtd/game"

func (a *App) revealMenuAfterFirstWave() {
	if !a.fromOW || a.trainingReplay != 0 || a.level != "hub" || a.lair == nil || a.lair.HasMenuReveal() {
		return
	}
	// A cleared wave or a defeat counts as played; merely starting it does not.
	if a.g.Wave < 1 || (a.g.Wave == 1 && a.g.WaveActive && a.g.Status == game.StatusRunning) {
		return
	}
	a.lair.MenuRevealed = true
	a.saveCampaign()
}
