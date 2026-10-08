package hiscore

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

// isolateHome points the hiscore file at a throwaway directory.
func isolateHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return home
}

func TestUpdateAndLoad(t *testing.T) {
	isolateHome(t)
	if _, isNew := Update("hub", 100); !isNew {
		t.Fatal("first score should be a new best")
	}
	if best, isNew := Update("hub", 50); isNew || best != 100 {
		t.Fatalf("best=%d isNew=%v, want 100 false", best, isNew)
	}
	if best, isNew := Update("hub", 200); !isNew || best != 200 {
		t.Fatalf("best=%d isNew=%v, want 200 true", best, isNew)
	}
	if got := Load()["hub"]; got != 200 {
		t.Fatalf("Load[hub] = %d, want 200", got)
	}
}

func TestSaveIsAtomicAndReadable(t *testing.T) {
	home := isolateHome(t)
	if err := Save(Table{"winding": 42, "hub": 7}); err != nil {
		t.Fatal(err)
	}
	p, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if got := Load(); got["winding"] != 42 || got["hub"] != 7 {
		t.Fatalf("Load after Save = %v", got)
	}
	// no temp files left behind
	entries, err := os.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != ".termtd-hiscores.json" {
			t.Fatalf("unexpected file %q in %s", e.Name(), home)
		}
	}
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && fi.Mode().Perm() != 0o644 {
		t.Fatalf("mode = %v, want 0644", fi.Mode().Perm())
	}
}

func TestCorruptFileCannotBeOverwrittenByUpdate(t *testing.T) {
	isolateHome(t)
	p, _ := Path()
	original := "{not json"
	if err := os.WriteFile(p, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	Update("hub", 10)
	data, err := os.ReadFile(p)
	if err != nil || string(data) != original {
		t.Fatalf("original changed: %q, %v", data, err)
	}
}

func TestPruneMazeKeepsTopScores(t *testing.T) {
	tbl := Table{"hub": 10, "canyon": 20}
	for i := 0; i < MaxMazeEntries+8; i++ {
		tbl[fmt.Sprintf("maze%d", i)] = i + 1
	}
	pruneMaze(tbl)
	if got := len(tbl) - 2; got != MaxMazeEntries {
		t.Fatalf("maze entries = %d, want %d", got, MaxMazeEntries)
	}
	if _, ok := tbl["hub"]; !ok {
		t.Fatal("hand-crafted level entry pruned")
	}
	// highest-scoring mazes survive, lowest-scoring ones are dropped
	if _, ok := tbl[fmt.Sprintf("maze%d", MaxMazeEntries+7)]; !ok {
		t.Fatal("highest-scoring maze pruned")
	}
	if _, ok := tbl["maze0"]; ok {
		t.Fatal("lowest-scoring maze not pruned")
	}
}

func TestUpdatePrunesBeyondCap(t *testing.T) {
	isolateHome(t)
	for i := 0; i < MaxMazeEntries+5; i++ {
		Update(fmt.Sprintf("maze%d", i), (i+1)*10)
	}
	tbl := Load()
	mazeCount := 0
	for k := range tbl {
		if strings.HasPrefix(k, "maze") {
			mazeCount++
		}
	}
	if mazeCount != MaxMazeEntries {
		t.Fatalf("maze entries = %d, want %d", mazeCount, MaxMazeEntries)
	}
}
