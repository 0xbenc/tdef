package tui

import (
	"fmt"
	"time"

	"tdef/game"
	"tdef/render"
)

// RunOverworld starts on the lair map (the overworld) instead of the title
// flow: a look-dev entry point for the Mario-style map of the lair's floors.
func RunOverworld() error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, game.Normal)
	a.ow = render.NewOWState()
	a.screen = ScreenOverworld
	return a.run()
}

// handleOverworld drives Grak around the lair map: arrows/wasd walk one cell
// onto walkable ground (corridors and floor pads), enter descends into the
// floor under the cursor, esc returns to the title, q quits. The mouse clicks
// a floor pad to step onto it.
func (a *App) handleOverworld(e Event) {
	if e.Mouse {
		a.handleOWMouse(e)
		return
	}
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	switch e.Rune {
	case 'q', 'Q':
		a.quit()
	case 'w', 'W':
		a.owMove(0, -1)
	case 's', 'S':
		a.owMove(0, 1)
	case 'a', 'A':
		a.owMove(-1, 0)
	case 'd', 'D':
		a.owMove(1, 0)
	case 'r', 'R':
		a.ow.RevealAll = !a.ow.RevealAll
		if a.ow.RevealAll {
			a.ow.Msg = "the whole lair, lit — every floor revealed"
		} else {
			a.ow.Msg = ""
		}
	}
	switch e.Key {
	case KeyUp:
		a.owMove(0, -1)
	case KeyDown:
		a.owMove(0, 1)
	case KeyLeft:
		a.owMove(-1, 0)
	case KeyRight:
		a.owMove(1, 0)
	case KeyEnter:
		a.owEnter()
	case KeyEscape:
		a.toScreen(ScreenTitle)
	}
}

// owMove steps the cursor one cell if the destination is walkable; a blocked
// step is a no-op, so Grak always stays on a corridor or a floor pad.
func (a *App) owMove(dx, dy int) {
	a.ow.Msg = ""
	n := game.Vec{X: a.ow.Cursor.X + dx, Y: a.ow.Cursor.Y + dy}
	if render.OWWalkable(n.X, n.Y) {
		a.ow.Cursor = n
	}
}

// owEnter descends into the floor under the cursor. Sealed floors are refused
// with a message; open ones launch the matching level (the Unmapped Depths
// roll a fresh maze seed, as they are uncharted).
func (a *App) owEnter() {
	fl, ok := render.OWFloorAt(a.ow.Cursor.X, a.ow.Cursor.Y)
	if !ok {
		return
	}
	if !a.ow.Unlocked[fl.ID] {
		a.ow.Msg = fl.Name + " is sealed — hold the lair to break it open"
		return
	}
	if fl.Level == "maze" {
		seed := time.Now().UnixNano()
		m, err := game.MazeFromSeed(seed)
		if err != nil {
			a.ow.Msg = err.Error()
			return
		}
		a.enterGame(m, fmt.Sprintf("maze%d", seed), a.diff)
		return
	}
	m, err := game.LoadLevel(fl.Level)
	if err != nil {
		a.ow.Msg = err.Error()
		return
	}
	a.enterGame(m, fl.Level, a.diff)
}

// handleOWMouse steps Grak onto the clicked floor pad (if any). Wheel is
// ignored on the map.
func (a *App) handleOWMouse(e Event) {
	if !e.Press {
		return
	}
	if e.Btn == 64 || e.Btn == 65 {
		return
	}
	w, h := a.termSize()
	if c, ok := render.OWFloorAtFrame(w, h, e.X, e.Y); ok {
		a.ow.Msg = ""
		a.ow.Cursor = c
	}
}
