package hiscore

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Table map[string]int

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tdef-hiscores.json"), nil
}

func Load() Table {
	t := Table{}
	p, err := Path()
	if err != nil {
		return t
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return t
	}
	json.Unmarshal(data, &t)
	return t
}

// Update records score under key and reports whether it is a new high score.
func Update(key string, score int) (best int, isNew bool) {
	t := Load()
	best = t[key]
	if score > best {
		t[key] = score
		isNew = true
		p, err := Path()
		if err == nil {
			data, _ := json.MarshalIndent(t, "", "  ")
			os.WriteFile(p, data, 0o644)
		}
		return score, true
	}
	return best, false
}
