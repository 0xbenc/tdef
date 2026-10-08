package hiscore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

// FloorRec is the lair's memory of one floor on one difficulty.
type FloorRec struct {
	Cleared  bool `json:"cleared,omitempty"`
	BestWave int  `json:"best_wave,omitempty"` // furthest wave reached
	LastWave int  `json:"last_wave,omitempty"` // wave reached on the last attempt
	LastWon  bool `json:"last_won,omitempty"`
}

// Lair is the persistent overworld state: per-floor results, the dragon
// hearts they earn, and the relic tokens Grak can spend. It lives in its own
// file so the hiscore table stays a plain map[string]int (old builds keep
// reading it).
type Lair struct {
	MenuRevealed bool                `json:"menu_revealed,omitempty"`
	BonusGold    int                 `json:"bonus_gold,omitempty"`
	BonusTower   bool                `json:"bonus_tower,omitempty"`
	BonusLives   int                 `json:"bonus_lives,omitempty"`
	Training     *TrainingProgress   `json:"training,omitempty"`
	IntroSeen    bool                `json:"intro_seen,omitempty"`
	Floors       map[string]FloorRec `json:"floors,omitempty"` // key "floor:diffIdx"
	Boss         map[int]bool        `json:"boss,omitempty"`   // diffIdx -> the heart is held
	Tokens       int                 `json:"tokens,omitempty"`
}

// Existing campaign records predate the menu reveal flag.
func (l *Lair) HasMenuReveal() bool {
	return l != nil && (l.MenuRevealed || l.AnyRecord() || (l.Training != nil && l.Training.Stage > 1))
}

// LairFloors are the built-in floors, in the lair's depth order (mouth to
// core). The Unmapped Depths is chaos, not progress: it never unseals
// anything and earns no hearts.
var LairFloors = []string{"rift", "rotunda", "halls", "garden"}

// HeartFloor is the endgame defense's floor id (the heart chamber).
const HeartFloor = "heart"

// DepthsFloor is the Unmapped Depths' floor id (procedural mazes).
const DepthsFloor = "depths"

func LairPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".termtd-lair.json"), nil
}

func emptyLair() *Lair {
	return &Lair{Floors: map[string]FloorRec{}, Boss: map[int]bool{}}
}

// LoadLair is the compatibility API. Interactive callers use LoadLairWithError.
func LoadLair() *Lair {
	value, _ := LoadLairWithError()
	return value
}

// LoadLairWithError returns initialized empty progress and an error when an
// existing save cannot be loaded. Such progress must not be saved over it.
func LoadLairWithError() (*Lair, error) {
	value := emptyLair()
	path, err := LairPath()
	if err != nil {
		return emptyLair(), err
	}
	if err := loadSave(path, &value); err != nil {
		return emptyLair(), err
	}
	if value.Floors == nil {
		value.Floors = map[string]FloorRec{}
	}
	if value.Boss == nil {
		value.Boss = map[int]bool{}
	}
	return value, nil
}

// SaveLair writes the lair atomically (temp file + rename), like Save.
func SaveLair(l *Lair) error {
	p, err := LairPath()
	if err != nil {
		return err
	}
	if err := checkSave(p, &Lair{}); err != nil {
		return err
	}
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".termtd-lair*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name) // no-op once the rename succeeds
	if _, err := tmp.Write(append(data, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, p)
}

func floorKey(floor string, diff int) string {
	return floor + ":" + strconv.Itoa(diff)
}

// Floor returns the record for a floor on a difficulty (zero if unknown).
func (l *Lair) Floor(floor string, diff int) FloorRec {
	return l.Floors[floorKey(floor, diff)]
}

// Record is the legacy best-effort persistence API. Interactive callers use
// RecordResult and save separately so load failures cannot be bypassed.
func (l *Lair) Record(floor string, diff, wave int, won bool) {
	l.RecordResult(floor, diff, wave, won)
	SaveLair(l)
}

// RecordResult folds one finished defense into memory without writing to disk.
// wave is the last wave reached (MaxWaves on a win). A held heart sticks:
// losing it later does not un-hold it.
func (l *Lair) RecordResult(floor string, diff, wave int, won bool) {
	rec := l.Floors[floorKey(floor, diff)]
	rec.LastWave = wave
	rec.LastWon = won
	if won {
		rec.Cleared = true
	}
	if wave > rec.BestWave {
		rec.BestWave = wave
	}
	l.Floors[floorKey(floor, diff)] = rec
	if floor == HeartFloor && won {
		l.Boss[diff] = true
	}
}

// ClearedCount is the number of built-in floors cleared on a difficulty —
// the dragon hearts held at that renown (0..4).
func (l *Lair) ClearedCount(diff int) int {
	n := 0
	for _, f := range LairFloors {
		if l.Floors[floorKey(f, diff)].Cleared {
			n++
		}
	}
	return n
}

// BossReady reports whether the heart has unsealed on a difficulty: every
// built-in floor held.
func (l *Lair) BossReady(diff int) bool {
	return l.ClearedCount(diff) >= len(LairFloors)
}

// DepthsReady reports whether every fixed defense, including the Heart,
// has been won on this difficulty. Procedural defenses are post-campaign.
func (l *Lair) DepthsReady(diff int) bool {
	return l != nil && l.BossReady(diff) && l.BossHeld(diff)
}

// BossHeld reports whether the heart has been held on a difficulty (the
// game's end state).
func (l *Lair) BossHeld(diff int) bool {
	return l.Boss[diff]
}

// AnyRecord reports whether the lair remembers anything at all; drives the
// first-run hint.
func (l *Lair) AnyRecord() bool {
	return len(l.Floors) > 0 || l.Tokens > 0
}
