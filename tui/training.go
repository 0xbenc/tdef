package tui

import (
	"fmt"
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/internal/copytext"
	"github.com/0xbenc/termtd/render"
	"math"
	"strconv"
)

// RunTutorial previews either guided defense without changing saved progress.
func RunTutorial(stage int) error {
	if stage != 1 && stage != 2 {
		return fmt.Errorf("tutorial stage must be 1 or 2")
	}
	term, err := Open()
	if err != nil {
		return err
	}
	a := &App{term: term, pal: render.Palette(), diff: game.Normal, events: make(chan Event, 256), scores: map[string]int{}, journal: &hiscore.Journal{Towers: map[string]bool{}}, trainingReplay: stage, trainingReplayMask: 1 << game.TowerGunner}
	name := "hub"
	if stage == 2 {
		name = "canyon"
	}
	m, err := game.LoadLevel(name)
	if err != nil {
		term.Close()
		return err
	}
	a.enterGame(m, name, game.Normal)
	return a.run()
}

func (a *App) configureTraining() {
	a.recruitPlanning = false
	if a.trainingReplay != 0 {
		a.g.ConfigureTraining(a.trainingReplay, a.trainingReplayMask)
		a.syncTraining()
		return
	}
	if !a.fromOW || a.lair == nil {
		return
	}
	progress := a.lair.EnsureTraining()
	if progress.Stage == 1 && a.level != "hub" {
		return
	}
	if progress.Stage < 1 || progress.Stage > 2 {
		return
	}
	a.g.ConfigureTraining(progress.Stage, progress.Unlocked)
	a.syncTraining()
}

func (a *App) syncTraining() {
	if a.trainingReplay != 0 {
		a.trainingReplayMask = a.g.UnlockedTowers
		if a.g.LessonPending {
			a.ui.Paused = true
			a.acc = 0
			a.journal.DiscoverTower(render.TowerJournalID(a.g.Lesson))
		}
		return
	}
	if a.g.TrainingStage == 0 || a.lair == nil {
		return
	}
	progress := a.lair.EnsureTraining()
	if progress.Unlocked != a.g.UnlockedTowers {
		progress.Unlocked = a.g.UnlockedTowers
		hiscore.SaveLair(a.lair)
	}
	if a.g.LessonPending {
		a.ui.Paused = true
		a.acc = 0
		a.discoverTower(a.g.Lesson)
	}
}

func (a *App) acknowledgeRecruit() {
	kind := a.g.Lesson
	a.g.LessonPending = false
	a.ui.Paused = true
	a.recruitPlanning = true
	a.ui.Placing, a.ui.PlacingOn = kind, true
	a.ui.Selected = -1
	a.suggestRecruitSite(kind)
	a.ui.RosterPage = 0
	if kind >= game.TowerRuneforge {
		a.ui.RosterPage = 1
	}
	a.msg(copytext.Text("ui.training.plan"))
}

// Move the preview to useful, unoccupied ground; the player still places it.
func (a *App) suggestRecruitSite(kind game.TowerKind) {
	s := a.g
	anchor := a.ui.Cursor.Center()
	best := -math.MaxFloat64
	for y := 1; y < s.Map.H-1; y++ {
		for x := 1; x < s.Map.W-1; x++ {
			cell := game.Vec{X: x, Y: y}
			if !s.CanBuild(cell, kind) {
				continue
			}
			for facing := game.FacingEast; facing <= game.FacingNorth; facing++ {
				if kind != game.TowerRuneforge && facing != game.FacingEast {
					break
				}
				end := game.ForgeEnd(s.Map, cell, facing, game.TowerSpecs[kind].Range[0])
				score := 0.
				for i, p := range s.Map.Samples {
					if kind == game.TowerGunner && s.Wave == 0 && i >= 36 {
						break
					}
					hits := cell.Center().Dist(p) <= game.TowerSpecs[kind].Range[0]
					if kind == game.TowerRuneforge {
						hits = game.ForgeContains(cell.Center(), end, p)
					}
					if !hits {
						continue
					}
					score++
					for _, tower := range s.Towers {
						if tower.Pos().Dist(p) <= tower.Range() {
							score += .5
						}
					}
				}
				score -= cell.Center().Dist(anchor) * .01
				if score > best {
					best = score
					a.ui.Cursor = cell
					a.ui.Facing = facing
				}
			}
		}
	}
}

func (a *App) completeTrainingDefense() {
	if a.g.TrainingStage == 0 || a.g.Status != game.StatusVictory || a.lair == nil {
		return
	}
	progress := a.lair.EnsureTraining()
	if a.g.TrainingStage == 1 {
		progress.Stage = 2
		progress.Unlocked |= game.MainTowersMask
	} else {
		progress.Stage = 3
		progress.Unlocked = game.AllTowersMask
	}
}

func (a *App) lockedRecruit(kind game.TowerKind) {
	if a.g.TrainingStage == 1 && kind >= game.TowerRuneforge {
		a.msg(copytext.Text("ui.training.next_defense"))
		return
	}
	for _, r := range game.RecruitmentOrder(a.g.TrainingStage) {
		if r.Kind == kind {
			a.msg(copytext.Format("ui.training.after_wave", "wave", strconv.Itoa(r.AfterWave)))
			return
		}
	}
	a.msg(copytext.Text("ui.training.locked"))
}
