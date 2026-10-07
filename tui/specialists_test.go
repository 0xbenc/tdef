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
	a.ui.PlacingOn = false
	a.ui.Selected = tower.ID
	a.rotateForge()
	if tower.Facing != game.FacingWest {
		t.Fatal("selected forge did not rotate")
	}
	a.layout = render.GameLayout(a.g.Map.W, a.g.Map.H, 80, 24)
	pager := render.RosterPager(80, 24)
	a.handleMenuClick(Event{X: pager.X + pager.W - 1, Y: pager.Y})
	if a.ui.RosterPage != 0 || a.ui.Selected != tower.ID {
		t.Fatal("mouse pager lost selected tower")
	}
}
