package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
	"os"
	"testing"
)

func TestFailedSavesRetainLatestProgressAndRetry(t *testing.T) {
	a := journalGameApp(t)
	lp, _ := hiscore.LairPath()
	jp, _ := hiscore.JournalPath()
	sp, _ := hiscore.Path()
	for _, p := range []string{lp, jp, sp} {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if err := os.Mkdir(p, 0700); err != nil {
			t.Fatal(err)
		}
	}
	a.g.Status = game.StatusVictory
	a.g.Wave = game.MaxWaves
	a.g.Score = 12345
	a.owFloorID = "rotunda"
	a.cleanWaves = 4
	a.stepGame(0)
	if !a.hasSaveErrors() || a.saveErrors[saveCampaign] == nil || a.saveErrors[saveJournal] == nil || a.saveErrors[saveScores] == nil {
		t.Fatal("failed writes were hidden")
	}
	tokens, wins := a.lair.Tokens, a.journal.EndgameWins
	a.stepGame(0)
	if a.lair.Tokens != tokens || a.journal.EndgameWins != wins {
		t.Fatal("failed save replayed rewards")
	}
	a.lair.Tokens++ // Later changes must be included in the retry.
	a.toScreen(ScreenHiscores)
	a.owRefresh()
	if a.lair.Tokens != tokens+1 || a.scores[a.level] != 12345 {
		t.Fatal("screen change reloaded stale disk data")
	}
	for _, p := range []string{lp, jp, sp} {
		if err := os.Remove(p); err != nil {
			t.Fatal(err)
		}
	}
	a.handle(Event{Key: KeyCtrlL})
	if a.hasSaveErrors() || hiscore.LoadLair().Tokens != tokens+1 || hiscore.Load()[a.level] != 12345 || !hiscore.LoadJournal().Stories["rotunda:1"] {
		t.Fatal("retry did not persist the current in-memory snapshots")
	}
}

func TestQuitWarnsBeforeDiscardingFailedSave(t *testing.T) {
	a := owTestApp(t)
	lp, _ := hiscore.LairPath()
	if err := os.Mkdir(lp, 0700); err != nil {
		t.Fatal(err)
	}
	a.saveCampaign()
	a.quit()
	if a.quitting || !a.saveQuitArmed {
		t.Fatal("quit silently discarded progress")
	}
	a.handle(Event{Rune: 'q'})
	if !a.quitting {
		t.Fatal("explicit second quit did not exit")
	}
}

func TestSuccessfulResetClearsPendingSaves(t *testing.T) {
	a, _ := resetTestApp(t)
	a.saveErrors[saveCampaign] = os.ErrPermission
	a.resetSel = 1
	a.activateResetProgress()
	a.retrySaves()
	if a.hasSaveErrors() || hiscore.LoadLair().Tokens != 0 {
		t.Fatal("retry restored reset progress")
	}
}

func TestSaveFailureNoticeFitsSmallFrame(t *testing.T) {
	f := &render.Frame{W: 62, H: 19, C: make([]render.Cell, 62*19)}
	render.DrawSaveFailure(f, "Campaign: permission denied", true)
	if f.C[16*62+2].BG != 52 {
		t.Fatal("save warning is not visible")
	}
}

func TestRelicReservationSurvivesReloadAndIsConsumedOnce(t *testing.T) {
	a := owTestApp(t)
	a.ow.Tokens = 3
	a.owSpendRelic(0)
	a.owSpendRelic(1)
	a.owSpendRelic(2)
	a.owRefresh()
	if a.ow.BonusGold != 60 || !a.ow.BonusTower || a.ow.BonusLives != 1 {
		t.Fatal("reserved rewards were lost on refresh")
	}
	a.bonusGold, a.bonusTower, a.bonusLives = a.ow.BonusGold, a.ow.BonusTower, a.ow.BonusLives
	a.fromOW = true
	m, _ := game.LoadLevel("hub")
	baseline := game.NewStateDiff(m, game.Normal)
	a.enterGame(m, "hub", game.Normal)
	if a.g.Gold != baseline.Gold+60 || a.g.Lives != baseline.Lives+1 || len(a.g.Towers) != 1 {
		t.Fatal("reserved rewards not granted")
	}
	saved := hiscore.LoadLair()
	if saved.BonusGold != 0 || saved.BonusTower || saved.BonusLives != 0 {
		t.Fatal("consumed rewards remained saved")
	}
}
