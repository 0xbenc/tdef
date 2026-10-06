package hiscore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResetProgressClearsCurrentAndLegacySaves(t *testing.T) {
	home := isolateHome(t)
	for _, prefix := range []string{".termtd-", ".tdef-"} {
		for _, name := range []string{"hiscores", "journal", "lair"} {
			if err := os.WriteFile(filepath.Join(home, prefix+name+".json"), []byte(`{"hub":42,"tokens":5,"endgame_wins":2}`), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	unrelated := filepath.Join(home, "other-game.json")
	if err := os.WriteFile(unrelated, []byte("keep me"), 0644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := ResetProgress(); err != nil {
			t.Fatal(err)
		}
		if len(Load()) != 0 || LoadJournal().EndgameWins != 0 || LoadLair().Tokens != 0 {
			t.Fatal("progress survived reset or returned from legacy saves")
		}
	}
	entries, err := os.ReadDir(home)
	if err != nil || len(entries) != 1 || entries[0].Name() != "other-game.json" {
		t.Fatalf("unexpected files after reset: %v, %v", entries, err)
	}
	if data, err := os.ReadFile(unrelated); err != nil || string(data) != "keep me" {
		t.Fatalf("unrelated file changed: %s, %v", data, err)
	}
}

func TestResetProgressRefusesDirectoriesBeforeRemovingSaves(t *testing.T) {
	home := isolateHome(t)
	if err := Save(Table{"hub": 42}); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(home, ".tdef-lair.json"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ResetProgress(); err == nil {
		t.Fatal("expected an error for directory at a save path")
	}
	if Load()["hub"] != 42 {
		t.Fatal("preflight failure must leave saves intact")
	}
}
