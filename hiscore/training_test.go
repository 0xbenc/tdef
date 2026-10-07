package hiscore

import (
	"github.com/0xbenc/termtd/game"
	"testing"
)

func TestTrainingMigrationAndSaveRoundTrip(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	fresh := emptyLair()
	p := fresh.EnsureTraining()
	if p.Stage != 1 || p.Unlocked != 1<<game.TowerGunner {
		t.Fatal("fresh save bypassed onboarding")
	}
	p.Unlocked |= 1 << game.TowerFrost
	if err := SaveLair(fresh); err != nil {
		t.Fatal(err)
	}
	loaded := LoadLair().EnsureTraining()
	if *loaded != *p {
		t.Fatal("earned recruits did not survive save/load")
	}
	legacy := emptyLair()
	legacy.Record("rotunda", 1, 3, false)
	migrated := legacy.EnsureTraining()
	if migrated.Stage != 3 || migrated.Unlocked != game.AllTowersMask {
		t.Fatal("legacy progress was restricted")
	}
}
