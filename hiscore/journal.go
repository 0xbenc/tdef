package hiscore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

// Journal remembers discoveries across levels and difficulties. There is no
// unread state: discovery never asks the player to acknowledge anything.
type Journal struct {
	Towers      map[string]bool `json:"towers,omitempty"`
	Enemies     map[string]bool `json:"enemies,omitempty"`
	Places      map[string]bool `json:"places,omitempty"`
	Stories     map[string]bool `json:"stories,omitempty"`
	EndgameWins int             `json:"endgame_wins,omitempty"`
}

func JournalPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".termtd-journal.json"), nil
}

// LoadJournal is the compatibility API. Interactive callers use LoadJournalWithError.
func LoadJournal() *Journal {
	value, _ := LoadJournalWithError()
	return value
}

// LoadJournalWithError returns initialized empty progress and an error when an
// existing save cannot be loaded. Such progress must not be saved over it.
func LoadJournalWithError() (*Journal, error) {
	value := &Journal{Towers: map[string]bool{}}
	path, err := JournalPath()
	if err != nil {
		return &Journal{Towers: map[string]bool{}}, err
	}
	if err := loadSave(path, &value); err != nil {
		return &Journal{Towers: map[string]bool{}}, err
	}
	if value.Towers == nil {
		value.Towers = map[string]bool{}
	}
	return value, nil
}

// DiscoverTower reports only the first discovery. Callers persist that change
// without a notification, screen switch or change to the simulation clock.
func discoverPage(pages *map[string]bool, id string) bool {
	if id == "" || (*pages)[id] {
		return false
	}
	if *pages == nil {
		*pages = map[string]bool{}
	}
	(*pages)[id] = true
	return true
}
func (j *Journal) DiscoverTower(id string) bool { return discoverPage(&j.Towers, id) }
func (j *Journal) DiscoverEnemy(id string) bool { return discoverPage(&j.Enemies, id) }
func (j *Journal) DiscoverPlace(id string) bool { return discoverPage(&j.Places, id) }
func (j *Journal) DiscoverStory(id string) bool { return discoverPage(&j.Stories, id) }

var AfterwardMilestones = []int{1, 3, 5, 10}

// RecordVictory is called once per finished defense. Procedural seeds share
// one story per difficulty; both Heart and Depths victories count as endgame.
func (j *Journal) RecordVictory(floor string, diff int) {
	if diff < 0 || diff > 2 {
		return
	}
	valid := false
	for _, id := range []string{"rotunda", "rift", "halls", "garden", "heart", "depths"} {
		if id == floor {
			valid = true
		}
	}
	if !valid {
		return
	}
	j.DiscoverStory(floorKey(floor, diff))
	if floor == HeartFloor || floor == DepthsFloor {
		j.EndgameWins++
		j.UnlockAfterward()
	}
}

func (j *Journal) UnlockAfterward() bool {
	changed := false
	for _, n := range AfterwardMilestones {
		if j.EndgameWins >= n {
			if j.DiscoverStory("afterward:" + strconv.Itoa(n)) {
				changed = true
			}
		}
	}
	return changed
}

func SaveJournal(j *Journal) error {
	p, err := JournalPath()
	if err != nil {
		return err
	}
	if err := checkSave(p, &Journal{}); err != nil {
		return err
	}
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".termtd-journal*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err = tmp.Write(append(data, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, p)
}
