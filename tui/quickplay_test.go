package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
)

func TestQuickPlayOnlySavesScores(t *testing.T) {
	for _, level := range append(game.LevelNames(), "maze1234", "heart") {
		for _, status := range []game.GameStatus{game.StatusVictory, game.StatusDefeat} {
			name := "win"
			if status == game.StatusDefeat {
				name = "loss"
			}
			t.Run(level+"/"+name, func(t *testing.T) {
				home := t.TempDir()
				t.Setenv("HOME", home)
				t.Setenv("USERPROFILE", home)
				lair := hiscore.LoadLair()
				lair.IntroSeen, lair.Tokens = true, 11
				lair.Record("rotunda", 1, game.MaxWaves, true)
				if err := hiscore.SaveLair(lair); err != nil {
					t.Fatal(err)
				}
				if err := hiscore.SaveJournal(&hiscore.Journal{Towers: map[string]bool{"gunner": true}}); err != nil {
					t.Fatal(err)
				}
				a := newApp(nil, game.Normal)
				a.ls = render.LSState{Levels: game.LevelNames(), Diff: 2, Seed: "1234"}
				// Simulate switching to Quick Play after a campaign run.
				a.fromOW, a.owFloorID = true, "heart"
				a.bonusGold, a.bonusLives, a.bonusTower = 300, 5, true
				before := map[string][]byte{}
				for _, name := range []string{"lair", "journal"} {
					current := filepath.Join(home, ".termtd-"+name+".json")
					data, err := os.ReadFile(current)
					if err != nil {
						t.Fatal(err)
					}
					legacy := filepath.Join(home, ".tdef-"+name+".json")
					if err := os.WriteFile(legacy, data, 0644); err != nil {
						t.Fatal(err)
					}
					before[current], before[legacy] = data, data
				}
				lairBefore, journalBefore := hiscore.LoadLair(), hiscore.LoadJournal()
				if level == "heart" {
					m, err := game.LoadBoss()
					if err != nil {
						t.Fatal(err)
					}
					a.enterQuickGame(m, level, game.Hard) // Direct CLI play uses the same scores-only policy.
				} else {
					a.ls.Cursor = len(a.ls.Levels)
					for i, n := range a.ls.Levels {
						if n == level {
							a.ls.Cursor = i
						}
					}
					a.startGame()
				}
				if a.screen != ScreenGame || a.fromOW || a.owFloorID != "" || a.ui.ToLair {
					t.Fatal("Quick Play retained campaign provenance")
				}
				fresh := game.NewStateDiff(a.g.Map, game.Hard)
				if a.g.Gold != fresh.Gold || a.g.Lives != fresh.Lives || len(a.g.Towers) != 0 {
					t.Fatal("Quick Play received campaign relic bonuses")
				}
				if a.bonusGold != 300 || a.bonusLives != 5 || !a.bonusTower {
					t.Fatal("Quick Play consumed campaign relic bonuses")
				}
				cell, ok := firstGrass(a.g.Map)
				if !ok {
					t.Fatal("no build site")
				}
				a.g.Gold = 10000
				a.ui.Cursor, a.ui.Placing, a.ui.PlacingOn = cell, game.TowerFrost, true
				a.place()
				if len(a.g.Towers) != 1 {
					t.Fatal("Quick Play tower placement failed")
				}
				a.g.SeenEnemies[game.EnemyShield] = true
				a.ui.Paused = true
				a.stepGame(.05)
				a.openJournal()
				a.handle(Event{Key: KeyEscape})
				if a.screen != ScreenGame {
					t.Fatal("journal did not return to Quick Play")
				}
				a.g.Status, a.g.Wave, a.g.Score = status, game.MaxWaves, 500
				for i := 0; i < 120; i++ {
					a.stepGame(1 / frameRate)
				}
				if a.screen != ScreenGame || a.heartEndingPending {
					t.Fatal("Quick Play triggered a campaign ending")
				}
				if hiscore.Load()[level] != 500 || a.scores[level] != 500 || !a.ui.NewBest {
					t.Fatal("Quick Play did not persist/display its high score")
				}
				a.restart()
				if a.fromOW || a.ui.ToLair {
					t.Fatal("restart changed Quick Play into a campaign run")
				}
				a.g.Status, a.g.Score = game.StatusDefeat, 400
				a.stepGame(.05)
				if hiscore.Load()[level] != 500 {
					t.Fatal("lower score overwrote high score")
				}
				a.handle(Event{Key: KeyEscape})
				if a.screen != ScreenLevelSelect {
					t.Fatal("Quick Play must return to level selection")
				}
				for path, data := range before {
					got, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(got, data) {
						t.Fatalf("Quick Play changed campaign save %s: %v", path, err)
					}
				}
				if !reflect.DeepEqual(a.lair, lairBefore) || !reflect.DeepEqual(a.journal, journalBefore) {
					t.Fatal("Quick Play changed cached campaign progress")
				}
			})
		}
	}
}

func TestQuickPlayDoesNotCreateCampaignSaves(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	a := newApp(nil, game.Normal)
	a.ls = render.LSState{Levels: game.LevelNames(), Diff: 1}
	a.startGame()
	a.g.Status, a.g.Score = game.StatusVictory, 123
	a.stepGame(.05)
	for _, name := range []string{"lair", "journal"} {
		if _, err := os.Stat(filepath.Join(home, ".termtd-"+name+".json")); !os.IsNotExist(err) {
			t.Fatalf("Quick Play created %s save: %v", name, err)
		}
	}
	a.leaveGame()
	a.activateMenu(0)
	if a.screen != ScreenCutscene || !a.filmRemember {
		t.Fatal("Quick Play must preserve the first campaign introduction")
	}
}

func TestCampaignStillRecordsProgressAndScores(t *testing.T) {
	a := owTestApp(t)
	a.lair.Training = &hiscore.TrainingProgress{Stage: 3, Unlocked: game.AllTowersMask}
	hiscore.SaveLair(a.lair)
	a.owRefresh()
	a.owLaunch("rotunda")
	if !a.fromOW {
		t.Fatal("campaign launch lost its provenance")
	}
	a.g.Gold = 10000
	cell, ok := firstGrass(a.g.Map)
	if !ok {
		t.Fatal("no build site")
	}
	a.ui.Cursor, a.ui.Placing, a.ui.PlacingOn = cell, game.TowerFrost, true
	a.place()
	a.g.SeenEnemies[game.EnemyShield] = true
	a.g.Status, a.g.Wave, a.g.Score = game.StatusVictory, game.MaxWaves, 500
	a.stepGame(.05)
	lair, journal := hiscore.LoadLair(), hiscore.LoadJournal()
	if !lair.Floor("rotunda", 1).Cleared || lair.Tokens != 2 || !journal.Towers["frost"] || !journal.Enemies[render.EnemyJournalID(game.EnemyShield)] || !journal.Places["rotunda"] || !journal.Stories["rotunda:1"] {
		t.Fatal("campaign run failed to persist progress, relics, or journal entries")
	}
	if hiscore.Load()["hub"] != 500 {
		t.Fatal("campaign high score not saved")
	}
	a.leaveGame()
	if a.screen != ScreenOverworld || !a.ow.Unlocked["rift"] {
		t.Fatal("campaign win did not return to/unlock the lair")
	}
}
