package hiscore

import (
	"os"
	"testing"
)

func TestLairLoadMissingYieldsEmpty(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	if l == nil || l.Floors == nil || l.Boss == nil {
		t.Fatalf("LoadLair on missing file = %+v, want an empty but initialised lair", l)
	}
	if l.AnyRecord() {
		t.Error("fresh lair reports AnyRecord")
	}
}

func TestLairLoadCorruptSelfHeals(t *testing.T) {
	home := isolateHome(t)
	p, err := LairPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	l := LoadLair()
	if l.AnyRecord() {
		t.Fatalf("LoadLair on corrupt file = %+v, want empty", l)
	}
	l.Record("rift", 1, 20, true)
	// the next load must see the repaired, persisted record
	if again := LoadLair(); !again.Floor("rift", 1).Cleared {
		t.Fatal("record lost after self-heal")
	}
	// no temp files left behind
	entries, err := os.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != ".termtd-lair.json" {
			t.Fatalf("unexpected file %q in %s", e.Name(), home)
		}
	}
}

func TestLairRecordTracksBestAndLast(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	l.Record("rift", 1, 7, false)
	r := l.Floor("rift", 1)
	if r.Cleared || r.LastWon || r.BestWave != 7 || r.LastWave != 7 {
		t.Fatalf("after a loss at wave 7: %+v", r)
	}
	l.Record("rift", 1, 12, false)
	r = l.Floor("rift", 1)
	if r.Cleared || r.BestWave != 12 || r.LastWave != 12 {
		t.Fatalf("after a loss at wave 12: %+v, want best 12 last 12 not cleared", r)
	}
	l.Record("rift", 1, 20, true)
	r = l.Floor("rift", 1)
	if !r.Cleared || !r.LastWon || r.BestWave != 20 {
		t.Fatalf("after a win: %+v, want cleared best 20", r)
	}
}

func TestLairRecordIsPerDifficulty(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	l.Record("rift", 0, 20, true) // easy
	if l.Floor("rift", 1).Cleared {
		t.Error("clearing rift on easy leaked into normal")
	}
	if l.ClearedCount(0) != 1 || l.ClearedCount(1) != 0 {
		t.Fatalf("cleared counts easy=%d normal=%d, want 1/0", l.ClearedCount(0), l.ClearedCount(1))
	}
}

func TestLairClearedCountAndBossReady(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	if l.BossReady(1) {
		t.Error("fresh lair has the heart unsealed")
	}
	for _, f := range []string{"rift", "rotunda", "halls"} {
		l.Record(f, 1, 20, true)
	}
	if l.ClearedCount(1) != 3 {
		t.Fatalf("cleared = %d, want 3", l.ClearedCount(1))
	}
	if l.BossReady(1) {
		t.Error("heart unsealed with 3 of 4 floors held")
	}
	l.Record("garden", 1, 20, true)
	if l.ClearedCount(1) != 4 || !l.BossReady(1) {
		t.Fatalf("cleared=%d bossReady=%v, want 4/true", l.ClearedCount(1), l.BossReady(1))
	}
}

// The Unmapped Depths is chaos, not progress: clearing it must not earn a
// heart or unseal the heart.
func TestLairDepthsEarnNoHeart(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	for _, f := range LairFloors {
		l.Record(f, 1, 20, true)
	}
	l.Record(DepthsFloor, 1, 20, true)
	if l.ClearedCount(1) != len(LairFloors) {
		t.Fatalf("depths counted toward hearts: %d", l.ClearedCount(1))
	}
}

func TestLairBossHeldSticks(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	l.Record(HeartFloor, 1, 9, false)
	if l.BossHeld(1) {
		t.Error("losing the heart marks it held")
	}
	l.Record(HeartFloor, 1, 20, true)
	if !l.BossHeld(1) {
		t.Error("holding the heart did not mark it held")
	}
	l.Record(HeartFloor, 1, 4, false)
	if !l.BossHeld(1) {
		t.Error("a later loss un-held the heart")
	}
}

func TestLairTokensPersist(t *testing.T) {
	isolateHome(t)
	l := LoadLair()
	l.Tokens = 3
	if err := SaveLair(l); err != nil {
		t.Fatal(err)
	}
	if again := LoadLair(); again.Tokens != 3 {
		t.Fatalf("tokens = %d, want 3", again.Tokens)
	}
}

func TestDepthsRequiresEveryFixedDefenseOnSameDifficulty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	l := LoadLair()
	// Even a Heart win from direct play cannot replace the other victories.
	l.Record(HeartFloor, 1, 20, true)
	if l.DepthsReady(1) {
		t.Fatal("Heart alone unlocked depths")
	}
	for i, floor := range LairFloors {
		l.Record(floor, 1, 20, true)
		if got, want := l.DepthsReady(1), i == len(LairFloors)-1; got != want {
			t.Fatalf("after %s: ready=%v want=%v", floor, got, want)
		}
	}
	if l.DepthsReady(0) || l.DepthsReady(2) {
		t.Fatal("victories unlocked other difficulties")
	}
	var missing *Lair
	if missing.DepthsReady(1) {
		t.Fatal("missing progression unlocked depths")
	}
}

func TestIntroMemorySurvivesProgressAndLegacySaves(t *testing.T) {
	home := isolateHome(t)
	if err := os.WriteFile(home+"/.termtd-lair.json", []byte(`{"floors":{"rotunda:1":{"best_wave":2}},"tokens":3}`), 0644); err != nil {
		t.Fatal(err)
	}
	l := LoadLair()
	if l.IntroSeen || !l.AnyRecord() {
		t.Fatal("legacy intro/progress migration changed the save")
	}
	l.IntroSeen = true
	if err := SaveLair(l); err != nil {
		t.Fatal(err)
	}
	l = LoadLair()
	l.Record("rift", 1, 20, true)
	if !LoadLair().IntroSeen || LoadLair().Tokens != 3 {
		t.Fatal("progress overwrote intro memory or tokens")
	}
}
