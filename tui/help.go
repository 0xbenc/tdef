package tui

import "github.com/0xbenc/termtd/render"

func (a *App) openHelp() {
	a.helpReturn = a.screen
	a.helpWasPaused = a.ui.Paused
	if a.screen == ScreenGame {
		a.ui.Paused = true
		a.acc = 0
	}
	a.help = render.HelpState{}
	a.toScreen(ScreenHelp)
}

func (a *App) closeHelp() {
	destination := a.helpReturn
	if destination != ScreenGame && destination != ScreenOverworld {
		destination = ScreenMenu
	}
	if destination == ScreenGame {
		a.ui.Paused = a.helpWasPaused
		a.acc = 0
	}
	a.screen = destination
	a.prev = nil
}

func (a *App) handleHelp(e Event) {
	if e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
		a.quit()
		return
	}
	if e.Key == KeyEscape || (!e.Mouse && (e.Rune == 'h' || e.Rune == 'H')) {
		a.closeHelp()
		return
	}
	page, scroll := a.help.Page, a.help.Scroll
	if e.Mouse {
		if !e.Press {
			return
		}
		switch e.Btn {
		case 64:
			scroll -= 3
		case 65:
			scroll += 3
		case 0:
			w, h := a.termSize()
			if e.Y == h-1 {
				a.closeHelp()
				return
			}
			if selected := render.HelpTopicAt(w, h, e.X, e.Y); selected >= 0 {
				page = selected
			}
		}
	} else {
		switch e.Key {
		case KeyLeft:
			page--
		case KeyRight, KeyTab:
			page++
		case KeyUp:
			scroll--
		case KeyDown:
			scroll++
		}
		switch e.Rune {
		case '[', 'a', 'A':
			page--
		case ']', 'd', 'D':
			page++
		case 'w', 'W':
			scroll--
		case 's', 'S':
			scroll++
		case '1', '2', '3', '4', '5':
			page = int(e.Rune - '1')
		}
	}
	count := render.HelpPageCount()
	page = (page + count) % count
	if page != a.help.Page {
		scroll = 0
	}
	w, h := a.termSize()
	a.help = render.HelpState{Page: page, Scroll: max(0, min(scroll, render.HelpMaxScroll(w, h, page)))}
}
