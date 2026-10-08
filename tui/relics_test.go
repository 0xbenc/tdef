package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
)

func TestRelicsRequirePerfectCompletedCampaignLevel(t *testing.T) {
	for _, difficulty := range []struct {
		name   string
		diff   game.Difficulty
		reward int
	}{
		{"easy", game.Easy, 1}, {"normal", game.Normal, 2}, {"hard", game.Hard, 3},
	} {
		for _, scenario := range []struct {
			name                            string
			status                          game.GameStatus
			lostHP, leaks                   int
			quickPlay, tutorial, incomplete bool
		}{
			{name: "perfect", status: game.StatusVictory},
			{name: "lost HP", status: game.StatusVictory, lostHP: 1, leaks: 1},
			{name: "lost bonus HP", status: game.StatusVictory, lostHP: 1, leaks: 1},
			{name: "healed breach", status: game.StatusVictory, leaks: 1},
			{name: "defeat", status: game.StatusDefeat, lostHP: 20, leaks: 20},
			{name: "quick play", status: game.StatusVictory, quickPlay: true},
			{name: "tutorial", status: game.StatusVictory, tutorial: true},
			{name: "incomplete", status: game.StatusVictory, incomplete: true},
		} {
			t.Run(difficulty.name+"/"+scenario.name, func(t *testing.T) {
				a := owTestApp(t)
				a.owBossSeen = map[int]bool{}
				a.journal = hiscore.LoadJournal()
				a.lair.Training = &hiscore.TrainingProgress{Stage: 3, Unlocked: game.AllTowersMask}
				a.fromOW = !scenario.quickPlay && !scenario.tutorial
				a.owFloorID = "rotunda"
				a.lair.Tokens = 7
				a.scores = map[string]int{}
				if scenario.name == "lost bonus HP" {
					a.bonusLives = 5
				}
				if scenario.tutorial {
					a.trainingReplay = 1
				}
				m, err := game.LoadLevel("hub")
				if err != nil {
					t.Fatal(err)
				}
				a.enterGame(m, "hub", difficulty.diff)
				a.g.Status = scenario.status
				a.g.Wave = game.MaxWaves
				if scenario.incomplete {
					a.g.Wave--
				}
				a.g.Lives = max(0, a.g.Lives-scenario.lostHP)
				a.g.TotalLeaks = scenario.leaks
				a.stepGame(0)
				expected := 0
				if scenario.name == "perfect" {
					expected = difficulty.reward
				}
				if a.lair.Tokens != 7+expected {
					t.Fatalf("tokens=%d, want %d", a.lair.Tokens, 7+expected)
				}
				a.stepGame(0)
				a.retrySaves()
				if a.lair.Tokens != 7+expected {
					t.Fatal("result or retry paid relics twice")
				}
				if a.fromOW && hiscore.LoadLair().Tokens != 7+expected {
					t.Fatal("relic reward not persisted")
				}
				if expected > 0 && !strings.Contains(a.ow.ReturnMsg, fmt.Sprintf("relics +%d", expected)) {
					t.Fatal("reward missing from result banner")
				}
			})
		}
	}
}

func TestClearingOneWaveDoesNotAwardRelics(t *testing.T) {
	a := journalGameApp(t)
	a.lair.Tokens = 7
	a.g.Wave = 1
	a.g.WaveActive = true
	a.g.SpawnQueue = nil
	a.g.Enemies = nil
	a.stepGame(1 / tickRate)
	if a.g.WaveActive || a.g.Status != game.StatusRunning {
		t.Fatal("wave did not finish")
	}
	if a.lair.Tokens != 7 {
		t.Fatal("individual wave awarded relics")
	}
}

func TestRestartResetsPerfectDefenseEligibility(t *testing.T) {
	a := journalGameApp(t)
	a.g.Lives--
	a.g.TotalLeaks = 1
	a.restart()
	a.g.Status = game.StatusVictory
	a.g.Wave = game.MaxWaves
	a.stepGame(0)
	if a.lair.Tokens != 2 {
		t.Fatalf("clean retry earned %d relics, want 2", a.lair.Tokens)
	}
}
