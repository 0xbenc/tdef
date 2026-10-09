package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/render"
	"testing"
)

func TestRosterSelectionRotationAndMouse(t *testing.T) {
	a := journalGameApp(t)
	a.ui.Selected = 42
	a.ui.PlacingOn = true
	a.changeRoster(1)
	if a.ui.RosterPage != 1 || a.ui.PlacingOn || a.ui.Selected != 42 {
		t.Fatal("page change lost selection or kept placement")
	}
	a.handle(Event{Rune: '1'})
	if !a.ui.PlacingOn || a.ui.Placing != game.TowerRuneforge {
		t.Fatal("specialist key used original roster")
	}
	a.rotateForge()
	if a.ui.Facing != game.FacingSouth {
		t.Fatal("preview did not rotate")
	}
	a.g.Gold = 10000
	for y := 1; y < a.g.Map.H; y++ {
		for x := 1; x < a.g.Map.W; x++ {
			v := game.Vec{X: x, Y: y}
			if a.g.CanBuild(v, game.TowerRuneforge) {
				a.ui.Cursor = v
				break
			}
		}
		if a.g.CanBuild(a.ui.Cursor, game.TowerRuneforge) {
			break
		}
	}
	a.place()
	tower := a.g.TowerAt(a.ui.Cursor)
	if tower == nil || tower.Facing != game.FacingSouth {
		t.Fatal("placement discarded preview facing")
	}
	if a.ui.RosterPage != 0 {
		t.Fatal("specialist placement did not return to regular defenders")
	}
	a.ui.PlacingOn = false
	a.ui.Selected = tower.ID
	a.rotateForge()
	if tower.Facing != game.FacingWest {
		t.Fatal("selected forge did not rotate")
	}
	a.layout = render.GameLayout(a.g.Map.W, a.g.Map.H, 80, 24)
	pager := render.RosterPager(80, 24)
	a.changeRoster(1)
	a.handleMenuClick(Event{X: pager.X + pager.W - 1, Y: pager.Y})
	if a.ui.RosterPage != 0 || a.ui.Selected != tower.ID {
		t.Fatal("mouse pager lost selected tower")
	}
}

func TestSpecialistPlacementReturnsToRegularRoster(t *testing.T) {
	for kind := game.TowerRuneforge; kind < game.TowerCount; kind++ {
		t.Run(game.TowerSpecs[kind].Name, func(t *testing.T) {
			a := journalGameApp(t)
			sites := keyboardSites(t, a)
			a.ui.RosterPage = 1
			a.ui.Cursor, a.ui.Placing, a.ui.PlacingOn = sites[0], kind, true
			a.g.Gold = 0
			a.handle(Event{Key: KeyEnter})
			if a.ui.RosterPage != 1 || !a.ui.PlacingOn {
				t.Fatal("failed placement changed the roster or canceled placement")
			}
			a.g.Gold = 10000
			a.handle(Event{Key: KeyEnter})
			if a.g.TowerAt(sites[0]) == nil || a.ui.RosterPage != 0 || a.ui.PlacingOn {
				t.Fatal("successful specialist placement did not return to regular roster")
			}
			a.handle(Event{Rune: '1'})
			if a.ui.Placing != game.TowerGunner || !a.ui.PlacingOn {
				t.Fatal("number key did not select a regular defender after placement")
			}
		})
	}
}
