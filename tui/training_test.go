package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"os"
	"testing"
)

func TestRecruitPlacementResumesAutomaticWaves(t *testing.T) {
	a := journalGameApp(t)
	a.g.ConfigureTraining(1, 1<<game.TowerGunner)
	a.acknowledgeRecruit()
	if !a.ui.Paused || !a.recruitPlanning {
		t.Fatal("recruit placement must allow safe planning")
	}
	a.place()
	if a.ui.Paused || a.recruitPlanning || a.ui.PlacingOn {
		t.Fatal("placing recruit did not resume play")
	}
	a.g.NextWaveAt = a.g.Time
	a.g.Step(.01)
	if !a.g.WaveActive {
		t.Fatal("wave required N after recruit placement")
	}
}

func TestTutorialReplayDoesNotWriteSaves(t *testing.T) {
	a := owTestApp(t)
	a.trainingReplay = 1
	a.trainingReplayMask = 1 << game.TowerGunner
	a.journal = &hiscore.Journal{Towers: map[string]bool{}}
	m, _ := game.LoadLevel("hub")
	a.enterGame(m, "hub", game.Normal)
	a.acknowledgeRecruit()
	a.g.Wave = 1
	a.g.WaveActive = true
	a.ui.Paused = false
	a.stepGame(.1)
	if !a.g.TowerAvailable(game.TowerFrost) || !a.journal.Towers["frost"] {
		t.Fatal("replay did not recruit")
	}
	a.g.Status = game.StatusVictory
	a.g.Wave = game.MaxWaves
	a.stepGame(.1)
	paths := []func() (string, error){hiscore.LairPath, hiscore.JournalPath, hiscore.Path}
	for _, path := range paths {
		p, err := path()
		if err != nil {
			t.Fatal(err)
		}
		if _, err = os.Stat(p); !os.IsNotExist(err) {
			t.Fatalf("replay wrote %s", p)
		}
	}
	a.leaveGame()
	if a.journal.Towers["frost"] {
		t.Fatal("replay discoveries leaked into campaign")
	}
	if a.trainingReplay != 0 {
		t.Fatal("replay mode escaped into menu")
	}
}

func TestFirstCampaignRecruitmentPersistsAndDoesNotLeakInput(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	a.owLaunch("rotunda")
	if a.g.TrainingStage != 1 || !a.g.LessonPending {
		t.Fatal("fresh campaign missing first lesson")
	}
	clock := a.g.Time
	a.stepGame(3)
	if a.g.Time != clock {
		t.Fatal("lesson did not freeze simulation")
	}
	a.handleGame(Event{Rune: 'n'})
	if a.g.WaveActive {
		t.Fatal("wave bypassed lesson")
	}
	a.handleGame(Event{Key: KeyEnter})
	if len(a.g.Towers) != 0 || a.g.LessonPending || !a.ui.Paused || !a.ui.PlacingOn {
		t.Fatal("acknowledgement placed a tower or resumed combat")
	}
	a.g.Wave = 1
	a.g.WaveActive = true
	a.ui.Paused = false
	a.stepGame(.1)
	if !a.g.LessonPending || a.g.Lesson != game.TowerFrost || !a.ui.Paused {
		t.Fatal("first clear did not safely recruit frost")
	}
	saved := hiscore.LoadLair()
	if saved.Training == nil || saved.Training.Unlocked&(1<<game.TowerFrost) == 0 || !hiscore.LoadJournal().Towers["frost"] {
		t.Fatal("recruit or journal did not persist immediately")
	}
	a.restart()
	if !a.g.TowerAvailable(game.TowerFrost) || a.g.TowerAvailable(game.TowerCannon) {
		t.Fatal("retry changed earned recruits")
	}
	a.g.Status = game.StatusVictory
	a.g.Wave = game.MaxWaves
	a.completeTrainingDefense()
	hiscore.SaveLair(a.lair)
	a.fromOW = true
	m, _ := game.LoadLevel("canyon")
	a.enterGame(m, "canyon", game.Normal)
	if a.g.TrainingStage != 2 || !a.g.LessonPending || a.g.Lesson != game.TowerRuneforge {
		t.Fatal("second defense did not introduce specialists")
	}
	if a.g.UnlockedTowers&game.MainTowersMask != game.MainTowersMask {
		t.Fatal("main defenders locked on second defense")
	}
	a.handleGame(Event{Key: KeyEnter})
	if a.ui.RosterPage != 1 || a.ui.Placing != game.TowerRuneforge {
		t.Fatal("specialist lesson did not open its roster")
	}
	a.g.Status = game.StatusVictory
	a.completeTrainingDefense()
	if a.lair.Training.Stage != 3 || a.lair.Training.Unlocked != game.AllTowersMask {
		t.Fatal("training did not finish")
	}
	a.enterQuickGame(m, "canyon", game.Hard)
	if a.g.TrainingStage != 0 || !a.g.TowerAvailable(game.TowerWitch) {
		t.Fatal("quick play inherited recruitment restrictions")
	}
}
