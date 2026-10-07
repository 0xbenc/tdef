package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/internal/copytext"
)

func (a *App) changeRoster(delta int) {
	a.ui.RosterPage = (a.ui.RosterPage + delta + 2) % 2
	a.ui.PlacingOn = false
}

func (a *App) rotateForge() {
	if a.ui.PlacingOn && a.ui.Placing == game.TowerRuneforge {
		a.ui.Facing = a.ui.Facing.Next()
		a.msg(copytext.Format("ui.specialists.aim_changed", "direction", string(a.ui.Facing.Arrow())))
		return
	}
	t := a.g.Tower(a.ui.Selected)
	if t == nil {
		t = a.g.TowerAt(a.ui.Cursor)
	}
	if t == nil || t.Kind != game.TowerRuneforge {
		a.msg(copytext.Text("ui.specialists.aim_prompt"))
		return
	}
	a.g.Aim(t, t.Facing.Next())
	a.ui.Selected = t.ID
	a.msg(copytext.Format("ui.specialists.aim_changed", "direction", string(t.Facing.Arrow())))
}
