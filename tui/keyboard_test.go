package tui

import (
	"testing"

	"github.com/0xbenc/termtd/game"
)

func keyboardSites(t *testing.T, a *App) []game.Vec {
	t.Helper()
	a.g.Gold = 10000
	var sites []game.Vec
	for y := 0; y < a.g.Map.H; y++ {
		for x := 0; x < a.g.Map.W; x++ {
			v := game.Vec{X: x, Y: y}
			if a.g.CanBuild(v, game.TowerGunner) {
				sites = append(sites, v)
			}
		}
	}
	if len(sites) < 2 {
		t.Fatal("need two build sites")
	}
	return sites
}

func TestKeyboardPlacementBuildsOnce(t *testing.T) {
	a := journalGameApp(t)
	sites := keyboardSites(t, a)
	a.ui.Cursor, a.ui.Placing, a.ui.PlacingOn = sites[0], game.TowerGunner, true
	a.handle(Event{Key: KeyEnter})
	if len(a.g.Towers) != 1 || a.ui.PlacingOn || a.ui.Selected >= 0 {
		t.Fatal("build must exit placement and selection")
	}
	a.ui.Cursor = sites[1]
	a.handle(Event{Key: KeyEnter})
	if len(a.g.Towers) != 1 {
		t.Fatal("Enter repeated the previous build")
	}
}

func TestKeyboardFocusActionsAndCycling(t *testing.T) {
	a := journalGameApp(t)
	sites := keyboardSites(t, a)
	first := a.g.Build(sites[0], game.TowerGunner)
	second := a.g.Build(sites[1], game.TowerGunner)
	a.ui.Cursor, a.ui.Selected = first.Cell, first.ID
	a.nextDefender()
	if a.ui.Selected != second.ID || a.ui.Cursor != second.Cell {
		t.Fatal("Tab did not move focus to next defender")
	}
	a.nextDefender()
	if a.ui.Selected != first.ID {
		t.Fatal("Tab did not wrap")
	}
	a.moveCursor(1, 0)
	if a.ui.Selected >= 0 {
		t.Fatal("moving retained stale selection")
	}
	a.ui.Cursor = second.Cell
	a.upgradeSelected()
	if second.Level != 2 || first.Level != 1 {
		t.Fatal("upgrade did not use hovered defender")
	}
	a.ui.Selected = -1
	before := second.TargetMode
	a.cycleTarget()
	if second.TargetMode != before.Next() || first.TargetMode != game.TargetFirst {
		t.Fatal("priority did not use hovered defender")
	}
	a.ui.Selected = -1
	a.sellSelected()
	if a.g.Tower(second.ID) != nil || a.g.Tower(first.ID) == nil {
		t.Fatal("sell did not use hovered defender")
	}
}
