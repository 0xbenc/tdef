package tui

import (
	"github.com/0xbenc/tdef/game"
	"github.com/0xbenc/tdef/hiscore"
	"github.com/0xbenc/tdef/render"
	"strconv"
)

// RunJournal opens the tower collection directly, returning to the lair.
func RunJournal() error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, game.Normal)
	a.ow = render.NewOWState()
	a.owRefresh()
	a.owBootArmed = true
	a.screen = ScreenOverworld
	a.openJournal()
	return a.run()
}

func (a *App) ensureJournal() {
	if a.journal == nil {
		a.journal = hiscore.LoadJournal()
	}
	if a.journalMigrated || a.lair == nil {
		return
	}
	a.journalMigrated = true
	changed := false
	knownWins := 0
	// Old records prove visits and completed waves, but not exact repeat counts.
	for _, floor := range render.JournalFloorIDs {
		for d := 0; d < 3; d++ {
			rec := a.lair.Floor(floor, d)
			if floor == hiscore.HeartFloor && a.lair.BossHeld(d) {
				rec.Cleared = true
			}
			if rec.BestWave > 0 || rec.Cleared {
				if a.journal.DiscoverPlace(floor) {
					changed = true
				}
			}
			completed := max(0, rec.BestWave-1)
			if rec.Cleared {
				completed = game.MaxWaves
				if a.journal.DiscoverStory(floor + ":" + strconv.Itoa(d)) {
					changed = true
				}
				if floor == hiscore.HeartFloor || floor == hiscore.DepthsFloor {
					knownWins++
				}
			}
			for wave := 1; wave <= min(completed, game.MaxWaves); wave++ {
				for _, spawn := range game.BuildWave(wave) {
					if a.journal.DiscoverEnemy(render.EnemyJournalID(spawn.Kind)) {
						changed = true
					}
				}
			}
		}
	}
	if a.journal.EndgameWins < knownWins {
		a.journal.EndgameWins = knownWins
		changed = true
	}
	if a.journal.UnlockAfterward() {
		changed = true
	}
	if changed {
		hiscore.SaveJournal(a.journal)
	}
}

func (a *App) recordEnemyDiscoveries() {
	a.ensureJournal()
	changed := false
	for k, seen := range a.g.SeenEnemies {
		if seen && a.journal.DiscoverEnemy(render.EnemyJournalID(game.EnemyKind(k))) {
			changed = true
		}
	}
	for _, e := range a.g.Enemies {
		if !e.Dead && !e.Leaked && a.journal.DiscoverEnemy(render.EnemyJournalID(e.Kind)) {
			changed = true
		}
	}
	if changed {
		hiscore.SaveJournal(a.journal)
	}
}

func (a *App) discoverPlace(id string) {
	for _, known := range render.JournalFloorIDs {
		if id == known {
			a.ensureJournal()
			if a.journal.DiscoverPlace(id) {
				hiscore.SaveJournal(a.journal)
			}
			return
		}
	}
}

func (a *App) discoverCurrentPlace() {
	if floor, ok := render.OWFloorAt(a.ow.Cursor.X, a.ow.Cursor.Y); ok && render.OWFloorOpen(floor.ID, a.ow) {
		a.discoverPlace(floor.ID)
	}
}

func (a *App) journalPages(section int) map[string]bool {
	switch section {
	case render.JournalEnemies:
		return a.journal.Enemies
	case render.JournalPlaces:
		return a.journal.Places
	case render.JournalStories:
		return a.journal.Stories
	default:
		return a.journal.Towers
	}
}

func (a *App) changeJournalChapter(delta int) {
	section := (a.journalUI.Section + delta + render.JournalSectionCount) % render.JournalSectionCount
	a.journalUI = render.JournalState{Section: section, Unlocked: a.journalPages(section)}
	a.prev = nil
}

func (a *App) discoverTower(k game.TowerKind) {
	a.ensureJournal()
	if a.journal.DiscoverTower(render.TowerJournalID(k)) {
		hiscore.SaveJournal(a.journal)
	}
}

func (a *App) openJournal() {
	a.ensureJournal()
	a.journalReturn = a.screen
	a.journalUI = render.JournalState{Unlocked: a.journal.Towers}
	a.toScreen(ScreenJournal)
}

func (a *App) handleJournal(e Event) {
	st := &a.journalUI
	entries := render.JournalEntries(st.Section)
	if !e.Mouse && (e.Rune == '[' || e.Rune == ']') {
		delta := 1
		if e.Rune == '[' {
			delta = -1
		}
		a.changeJournalChapter(delta)
		return
	}
	if e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
		a.quit()
		return
	}
	if e.Mouse {
		if !e.Press {
			return
		}
		if st.Reading {
			if e.Btn == 64 {
				a.scrollJournal(-3)
			} else if e.Btn == 65 {
				a.scrollJournal(3)
			}
		} else if e.Btn == 0 {
			w, h := a.termSize()
			for i, r := range render.JournalCards(w, h, st.Section) {
				if r.Contains(e.X, e.Y) {
					index := render.JournalCollectionStart(w, h, st.Section, st.Cursor) + i
					if index >= len(entries) {
						return
					}
					st.Cursor = index
					a.readJournalPage()
					break
				}
			}
		}
		return
	}
	if e.Key == KeyEscape || e.Rune == 'j' || e.Rune == 'J' {
		if st.Reading {
			st.Reading = false
			st.Scroll = 0
		} else {
			// No lair refresh, boot replay, auto-resume or simulation step on return.
			a.screen = a.journalReturn
			a.prev = nil
			a.acc = 0
		}
		return
	}
	if st.Reading {
		switch e.Key {
		case KeyLeft:
			a.turnJournalPage(-1)
		case KeyRight, KeyTab:
			a.turnJournalPage(1)
		case KeyUp:
			a.scrollJournal(-1)
		case KeyDown:
			a.scrollJournal(1)
		}
		switch e.Rune {
		case 'a', 'A':
			a.turnJournalPage(-1)
		case 'd', 'D':
			a.turnJournalPage(1)
		case 'w', 'W':
			a.scrollJournal(-1)
		case 's', 'S':
			a.scrollJournal(1)
		}
	} else {
		delta := 0
		vertical := 0
		switch e.Key {
		case KeyLeft:
			delta = -1
		case KeyUp:
			vertical = -1
		case KeyRight, KeyTab:
			delta = 1
		case KeyDown:
			vertical = 1
		case KeyEnter:
			a.readJournalPage()
		}
		switch e.Rune {
		case 'a', 'A':
			delta = -1
		case 'w', 'W':
			vertical = -1
		case 'd', 'D':
			delta = 1
		case 's', 'S':
			vertical = 1
		case '1', '2', '3', '4', '5', '6', '7', '8':
			w, h := a.termSize()
			index := render.JournalCollectionStart(w, h, st.Section, st.Cursor) + int(e.Rune-'1')
			if index < len(entries) {
				st.Cursor = index
				a.readJournalPage()
			}
		}
		if vertical != 0 {
			w, h := a.termSize()
			if w >= 76 && h >= 24 {
				row := st.Cursor / 4
				if row+vertical >= 0 && row+vertical < (len(entries)+3)/4 {
					st.Cursor = min(len(entries)-1, st.Cursor+vertical*4)
				}
			} else {
				delta = vertical
			}
		}
		if delta != 0 {
			st.Cursor = (st.Cursor + delta + len(entries)) % len(entries)
		}
	}
}

func (a *App) readJournalPage() {
	st := &a.journalUI
	if st.Unlocked[render.JournalEntries(st.Section)[st.Cursor].ID] {
		st.Reading = true
		st.Scroll = 0
	}
}

func (a *App) turnJournalPage(delta int) {
	st := &a.journalUI
	entries := render.JournalEntries(st.Section)
	for i := 0; i < len(entries); i++ {
		st.Cursor = (st.Cursor + delta + len(entries)) % len(entries)
		if st.Unlocked[render.JournalEntries(st.Section)[st.Cursor].ID] {
			st.Scroll = 0
			return
		}
	}
}

func (a *App) scrollJournal(delta int) {
	w, h := a.termSize()
	a.journalUI.Scroll = max(0, min(a.journalUI.Scroll+delta, render.JournalMaxScroll(w, h, a.journalUI.Cursor, a.journalUI.Section)))
}
