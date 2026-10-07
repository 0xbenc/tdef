package hiscore

import "github.com/0xbenc/termtd/game"

type TrainingProgress struct {
	Stage    int    `json:"stage"`
	Unlocked uint32 `json:"unlocked"`
}

// A missing training record on an established save must never relock units.
func (l *Lair) EnsureTraining() *TrainingProgress {
	if l.Training == nil {
		l.Training = &TrainingProgress{Stage: 1, Unlocked: 1 << game.TowerGunner}
		if l.AnyRecord() || len(l.Boss) > 0 {
			l.Training.Stage = 3
			l.Training.Unlocked = game.AllTowersMask
		}
	}
	return l.Training
}
