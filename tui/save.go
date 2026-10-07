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

func (a *App) saveCampaign() { a.saved(saveCampaign, hiscore.SaveLair(a.lair)) }
func (a *App) saveJournal()  { a.saved(saveJournal, hiscore.SaveJournal(a.journal)) }
func (a *App) saveScores()   { a.saved(saveScores, hiscore.Save(hiscore.Table(a.scores))) }

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
