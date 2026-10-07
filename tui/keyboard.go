package tui

import "github.com/0xbenc/termtd/game"

func (a *App) focusedTower() *game.Tower {
	if a.ui.PlacingOn {
		return nil
	}
	if t := a.g.Tower(a.ui.Selected); t != nil {
		return t
	}
	return a.g.TowerAt(a.ui.Cursor)
}

func (a *App) nextDefender() {
	if len(a.g.Towers) == 0 {
		return
	}
	index := -1
	if t := a.focusedTower(); t != nil {
		for i, tower := range a.g.Towers {
			if tower.ID == t.ID {
				index = i
				break
			}
		}
	}
	t := a.g.Towers[(index+1)%len(a.g.Towers)]
	a.ui.Cursor, a.ui.Selected = t.Cell, t.ID
	a.ui.PlacingOn = false
}
