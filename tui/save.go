package tui

import (
	"github.com/0xbenc/termtd/hiscore"
	"strings"
	"time"
)

const (
	saveCampaign = iota
	saveJournal
	saveScores
)

func (a *App) hasSaveErrors() bool {
	for _, err := range a.saveErrors {
		if err != nil {
			return true
		}
	}
	return false
}

func (a *App) saved(kind int, err error) {
	a.saveErrors[kind] = err
	if a.hasSaveErrors() {
		a.saveRetryAt = time.Now().Add(5 * time.Second)
	} else {
		a.saveRetryAt = time.Time{}
		a.saveQuitArmed = false
	}
}

func (a *App) saveCampaign() {
	if a.loadErrors[saveCampaign] == nil {
		a.saved(saveCampaign, hiscore.SaveLair(a.lair))
	}
}
func (a *App) saveJournal() {
	if a.loadErrors[saveJournal] == nil {
		a.saved(saveJournal, hiscore.SaveJournal(a.journal))
	}
}
func (a *App) saveScores() {
	if a.loadErrors[saveScores] == nil {
		a.saved(saveScores, hiscore.Save(hiscore.Table(a.scores)))
	}
}

// Retry the current in-memory snapshots, never replay victories or rewards.
func (a *App) retrySaves() {
	if a.saveErrors[saveCampaign] != nil {
		a.saveCampaign()
	}
	if a.saveErrors[saveJournal] != nil {
		a.saveJournal()
	}
	if a.saveErrors[saveScores] != nil {
		a.saveScores()
	}
}

func (a *App) saveErrorText() string {
	var failures []string
	for i, name := range []string{"Campaign", "Journal", "High scores"} {
		if err := a.saveErrors[i]; err != nil {
			failures = append(failures, name+": "+err.Error())
		}
	}
	return strings.Join(failures, "; ")
}

func (a *App) hasLoadErrors() bool {
	for _, err := range a.loadErrors {
		if err != nil {
			return true
		}
	}
	return false
}

func (a *App) loadErrorText() string {
	var failures []string
	for i, name := range []string{"Campaign", "Journal", "High scores"} {
		if err := a.loadErrors[i]; err != nil {
			failures = append(failures, name+": "+err.Error())
		}
	}
	return strings.Join(failures, "\n")
}

// Keep the last successfully loaded snapshot if a later refresh fails.
func (a *App) reloadCampaign() {
	value, err := hiscore.LoadLairWithError()
	a.loadErrors[saveCampaign] = err
	if err == nil || a.lair == nil {
		a.lair = value
	}
}

func (a *App) reloadScores() {
	value, err := hiscore.LoadWithError()
	a.loadErrors[saveScores] = err
	if err == nil || a.scores == nil {
		a.scores = value
	}
}
