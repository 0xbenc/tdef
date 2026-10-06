package hiscore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLegacyProgressSurvivesRename(t *testing.T) {
	home := isolateHome(t)
	files := map[string]string{
		"hiscores": `{"hub":42}`,
		"journal":  `{"towers":{"cannon":true},"endgame_wins":2}`,
		"lair":     `{"intro_seen":true,"tokens":3}`,
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(home, ".tdef-"+name+".json"), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if Load()["hub"] != 42 {
		t.Fatal("legacy score lost")
	}
	j := LoadJournal()
	if !j.Towers["cannon"] || j.EndgameWins != 2 {
		t.Fatal("legacy discoveries lost")
	}
	l := LoadLair()
	if !l.IntroSeen || l.Tokens != 3 {
		t.Fatal("legacy lair progress lost")
	}
	if err := Save(Table{"hub": 84}); err != nil {
		t.Fatal(err)
	}
	if err := SaveJournal(j); err != nil {
		t.Fatal(err)
	}
	if err := SaveLair(l); err != nil {
		t.Fatal(err)
	}
	for name := range files {
		path := filepath.Join(home, ".termtd-"+name+".json")
		if _, err := os.Stat(path); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(`{broken`), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if Load()["hub"] != 0 || LoadJournal().EndgameWins != 0 || LoadLair().Tokens != 0 {
		t.Fatal("existing new saves must take precedence over legacy files, even when corrupt")
	}
	for name, data := range files {
		got, err := os.ReadFile(filepath.Join(home, ".tdef-"+name+".json"))
		if err != nil || string(got) != data {
			t.Fatalf("legacy %s modified: %v", name, err)
		}
	}
}
