package tui

import (
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
	"time"
)

func (a *App) handleResetProgress(e Event) {
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	if a.resetDone {
		if e.Key == KeyEnter || e.Key == KeyEscape || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
			a.toScreen(ScreenMenu)
		}
		return
	}
	if e.Mouse {
		if !e.Press || e.Btn != 0 {
			return
		}
		w, h := a.termSize()
		for i, rect := range render.ResetProgressRects(w, h) {
			if rect.Contains(e.X, e.Y) {
				a.resetSel = i
				a.activateResetProgress()
				return
			}
		}
		return
	}
	switch {
	case e.Key == KeyEscape || e.Rune == 'n' || e.Rune == 'N' || e.Rune == 'q' || e.Rune == 'Q':
		a.toScreen(ScreenMenu)
	case e.Key == KeyUp || e.Key == KeyDown || e.Key == KeyTab || e.Rune == 'w' || e.Rune == 'W' || e.Rune == 's' || e.Rune == 'S':
		a.resetSel = 1 - a.resetSel
	case e.Key == KeyEnter:
		a.activateResetProgress()
	}
}

func (a *App) activateResetProgress() {
	if a.resetSel == 0 {
		a.toScreen(ScreenMenu)
		return
	}
	err := hiscore.ResetProgress()
	if err != nil {
		a.resetErr = err.Error()
		a.resetSel = 0
		return
	}
	a.scores = hiscore.Load()
	a.lair = hiscore.LoadLair()
	a.journal = hiscore.LoadJournal()
	a.saveErrors = [3]error{}
	a.saveRetryAt = time.Time{}
	a.saveQuitArmed = false
	a.journalUI = render.JournalState{}
	a.journalMigrated = false
	a.hsTop = 0
	a.ow = render.NewOWState()
	a.owBootArmed = false
	a.owBossSeen = nil
	a.film = render.CutsceneState{}
	a.filmRemember = false
	a.heartEndingPending = false
	a.bonusGold, a.bonusLives, a.bonusTower = 0, 0, false
	a.g, a.scored, a.fromOW = nil, false, false
	a.level, a.owFloorID = "", ""
	a.cleanWaves, a.waveStartLives = 0, 0
	a.resetErr = ""
	a.resetDone = true
}
