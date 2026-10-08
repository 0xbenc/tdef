package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFailedSavesRetainLatestProgressAndRetry(t *testing.T) {
	a := journalGameApp(t)
	a.scores = map[string]int{} // Already loaded before the disk becomes unavailable.
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

func TestLoadFailuresBlockSessionWritesUntilRestart(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	for _, name := range []string{"lair", "journal", "hiscores"} {
		if err := os.WriteFile(filepath.Join(home, ".termtd-"+name+".json"), []byte(`{broken`), 0600); err != nil {
			t.Fatal(err)
		}
	}
	a := newApp(nil, game.Normal)
	if !a.hasLoadErrors() || !strings.Contains(a.loadErrorText(), "recovery copy:") {
		t.Fatal("load failures hidden")
	}
	a.lair.Tokens = 99
	a.journal.EndgameWins = 99
	a.scores["hub"] = 99
	a.saveCampaign()
	a.saveJournal()
	a.saveScores()
	a.retrySaves()
	for _, name := range []string{"lair", "journal", "hiscores"} {
		path := filepath.Join(home, ".termtd-"+name+".json")
		data, err := os.ReadFile(path)
		if err != nil || string(data) != `{broken` {
			t.Fatal("original overwritten")
		}
		if err := os.WriteFile(path, []byte(`{}`), 0600); err != nil {
			t.Fatal(err)
		}
	}
	// Even once the disk file is fixed, this session's empty-derived progress
	// must not replace recovered progress. Only a restart loads it safely.
	a.saveCampaign()
	a.saveJournal()
	a.saveScores()
	a.toScreen(ScreenHiscores)
	a.owRefresh()
	if a.lair.Tokens != 99 || a.journal.EndgameWins != 99 || a.scores["hub"] != 99 {
		t.Fatal("screen transition discarded session state")
	}
	for _, name := range []string{"lair", "journal", "hiscores"} {
		data, err := os.ReadFile(filepath.Join(home, ".termtd-"+name+".json"))
		if err != nil || string(data) != `{}` {
			t.Fatal("repaired save overwritten by blocked session")
		}
	}
	a.handle(Event{Key: KeyCtrlL})
	if !a.loadDetails {
		t.Fatal("recovery details unavailable")
	}
	a.handle(Event{Key: KeyDown})
	if a.loadDetailTop != 1 {
		t.Fatal("details do not scroll")
	}
	a.handle(Event{Key: KeyEscape})
	if a.loadDetails {
		t.Fatal("details did not close")
	}
	a.quit()
	if a.quitting || !a.saveQuitArmed {
		t.Fatal("quit hid unsaved session progress")
	}
	a.handle(Event{Rune: 'q'})
	if !a.quitting {
		t.Fatal("explicit quit failed")
	}
	restarted := newApp(nil, game.Normal)
	if restarted.hasLoadErrors() {
		t.Fatal("repair did not recover after restart")
	}
	restarted.lair.Tokens = 7
	restarted.saveCampaign()
	if hiscore.LoadLair().Tokens != 7 {
		t.Fatal("repaired save remained blocked")
	}
}

func TestResetClearsLoadBlockAndPreservesRecoveryCopies(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	path, _ := hiscore.LairPath()
	if err := os.WriteFile(path, []byte(`{broken`), 0600); err != nil {
		t.Fatal(err)
	}
	a := newApp(nil, game.Normal)
	copies, _ := filepath.Glob(path + ".corrupt-*")
	a.resetSel = 1
	a.activateResetProgress()
	if a.hasLoadErrors() || !a.resetDone {
		t.Fatal("confirmed reset did not unblock saves")
	}
	if len(copies) != 1 {
		t.Fatal("missing recovery copy")
	}
	if _, err := os.Stat(copies[0]); err != nil {
		t.Fatal("reset deleted recovery copy")
	}
	a.lair.Tokens = 8
	a.saveCampaign()
	if hiscore.LoadLair().Tokens != 8 {
		t.Fatal("reset failed to restore saving")
	}
}

func TestLoadWarningFitsAndDetailsScrollAtMinimumSize(t *testing.T) {
	f := &render.Frame{W: 62, H: 19, C: make([]render.Cell, 62*19)}
	render.DrawLoadFailure(f, "Campaign: "+strings.Repeat("detail ", 100)+"END", true, true, 1000)
	if !strings.Contains(f.Text(), "END") || !strings.Contains(f.Text(), "Q again") {
		t.Fatal("recovery details or quit warning unreachable")
	}
}

func TestLoadFailureOnlyBlocksAffectedSaveAndDetailsPause(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	path, _ := hiscore.JournalPath()
	if err := os.WriteFile(path, []byte(`{broken`), 0600); err != nil {
		t.Fatal(err)
	}
	a := newApp(nil, game.Normal)
	if a.loadErrors[saveJournal] == nil || a.loadErrors[saveCampaign] != nil || a.loadErrors[saveScores] != nil {
		t.Fatal("wrong save blocked")
	}
	a.lair.Tokens = 12
	a.scores["hub"] = 34
	a.saveCampaign()
	a.saveScores()
	if hiscore.LoadLair().Tokens != 12 || hiscore.Load()["hub"] != 34 {
		t.Fatal("healthy saves were blocked")
	}
	a.screen = ScreenGame
	a.ui.Paused = false
	a.acc = 1
	a.handle(Event{Key: KeyCtrlL})
	if !a.ui.Paused || a.acc != 0 {
		t.Fatal("game continued behind recovery details")
	}
}

func TestRefreshLoadFailureRetainsKnownProgress(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if err := hiscore.SaveLair(&hiscore.Lair{Tokens: 27}); err != nil {
		t.Fatal(err)
	}
	if err := hiscore.Save(hiscore.Table{"hub": 123}); err != nil {
		t.Fatal(err)
	}
	a := newApp(nil, game.Normal)
	for _, name := range []string{"lair", "hiscores"} {
		if err := os.WriteFile(filepath.Join(home, ".termtd-"+name+".json"), []byte(`{broken`), 0600); err != nil {
			t.Fatal(err)
		}
	}
	a.toScreen(ScreenHiscores)
	a.owRefresh()
	if !a.hasLoadErrors() || a.lair.Tokens != 27 || a.scores["hub"] != 123 {
		t.Fatal("failed refresh replaced known progress")
	}
}
