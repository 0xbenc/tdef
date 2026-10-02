package hiscore

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestJournalPersistsOnlyFirstDiscovery(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	j := LoadJournal()
	if len(j.Towers) != 0 {
		t.Fatal("fresh journal was pre-unlocked")
	}
	if j.DiscoverTower("") || !j.DiscoverTower("gunner") || j.DiscoverTower("gunner") {
		t.Fatal("discovery is not idempotent")
	}
	if err := SaveJournal(j); err != nil {
		t.Fatal(err)
	}
	j = LoadJournal()
	if !j.Towers["gunner"] || len(j.Towers) != 1 {
		t.Fatal("discovery lost on reload")
	}
	j.DiscoverTower("frost")
	if err := SaveJournal(j); err != nil {
		t.Fatal(err)
	}
	if len(LoadJournal().Towers) != 2 {
		t.Fatal("new discovery erased old page")
	}
}

func TestJournalHandlesOldMissingOrCorruptSave(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	p, _ := JournalPath()
	for _, data := range []string{"{}", "null", "{broken", "{\"towers\":null}"} {
		if err := os.WriteFile(p, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
		j := LoadJournal()
		if !j.DiscoverTower("ranger") {
			t.Fatal("save could not accept discovery")
		}
		if err := SaveJournal(j); err != nil {
			t.Fatal(err)
		}
		if !LoadJournal().Towers["ranger"] {
			t.Fatal("save did not recover")
		}
	}
}

func TestJournalVictoryDifficultyAndBonusMilestones(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	j := LoadJournal()
	j.RecordVictory("rotunda", 0)
	j.RecordVictory("rotunda", 2)
	if !j.Stories["rotunda:0"] || !j.Stories["rotunda:2"] || j.Stories["rotunda:1"] || j.EndgameWins != 0 {
		t.Fatal("difficulty victories mixed or counted as endgame")
	}
	for n := 1; n <= 10; n++ {
		floor := HeartFloor
		if n%2 == 0 {
			floor = DepthsFloor
		}
		j.RecordVictory(floor, 1)
		if j.EndgameWins != n {
			t.Fatal("endgame win not counted")
		}
		for _, threshold := range AfterwardMilestones {
			if j.Stories["afterward:"+strconv.Itoa(threshold)] != (n >= threshold) {
				t.Fatalf("milestone %d at win %d", threshold, n)
			}
		}
	}
	if err := SaveJournal(j); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(j, LoadJournal()) {
		t.Fatal("expanded journal did not persist")
	}
	j.RecordVictory("unknown", 1)
	j.RecordVictory("heart", 99)
	if j.EndgameWins != 10 {
		t.Fatal("invalid victory affected journal")
	}
}
