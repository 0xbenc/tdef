package hiscore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MaxMazeEntries caps how many procedural-maze entries the table keeps.
// Maze entries are keyed by seed, so without a cap the file would grow one
// entry per maze ever played. Hand-crafted levels (no "maze" prefix) are
// always kept.
const MaxMazeEntries = 32

type Table map[string]int

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".termtd-hiscores.json"), nil
}

// Load is the compatibility API. Interactive callers use LoadWithError.
func Load() Table {
	value, _ := LoadWithError()
	return value
}

// LoadWithError returns initialized empty progress and an error when an
// existing save cannot be loaded. Such progress must not be saved over it.
func LoadWithError() (Table, error) {
	value := Table{}
	path, err := Path()
	if err != nil {
		return Table{}, err
	}
	if err := loadSave(path, &value); err != nil {
		return Table{}, err
	}
	if value == nil {
		value = Table{}
	}
	return value, nil
}

// Save writes the table atomically: temp file in the same directory, then
// rename over the real one. A crash mid-write leaves either the old file or
// a new one, never a truncated one.
func Save(t Table) error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := checkSave(p, &Table{}); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".termtd-hiscores*.tmp")
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

// pruneMaze keeps at most MaxMazeEntries maze entries (highest scores first,
// ties broken by key so the result is deterministic).
func pruneMaze(t Table) {
	var mazeKeys []string
	for k := range t {
		if strings.HasPrefix(k, "maze") {
			mazeKeys = append(mazeKeys, k)
		}
	}
	if len(mazeKeys) <= MaxMazeEntries {
		return
	}
	sort.Slice(mazeKeys, func(i, j int) bool {
		if t[mazeKeys[i]] != t[mazeKeys[j]] {
			return t[mazeKeys[i]] > t[mazeKeys[j]]
		}
		return mazeKeys[i] < mazeKeys[j]
	})
	for _, k := range mazeKeys[MaxMazeEntries:] {
		delete(t, k)
	}
}

// RecordScore updates an in-memory table, including the maze retention cap.
// The caller owns persistence and can retain this table after a failed save.
func RecordScore(t Table, key string, score int) (best int, isNew bool) {
	if score > t[key] {
		t[key] = score
		isNew = true
		pruneMaze(t)
	}
	return t[key], isNew
}

// Update is the legacy best-effort convenience API. The game uses
// RecordScore and Save separately to report failures and retry.
func Update(key string, score int) (best int, isNew bool) {
	t := Load()
	best, isNew = RecordScore(t, key, score)
	if isNew {
		Save(t)
	}
	return best, isNew
}
