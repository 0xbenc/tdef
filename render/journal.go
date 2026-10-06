package render

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"strings"

	"github.com/0xbenc/termtd/game"
)

type TowerJournalEntry struct {
	ID                         string
	Kind                       game.TowerKind
	Name, Subtitle, Lore, Note string
	Accent                     int
	paint                      func(*Frame, int, int, int, int)
}

// One definitive image per tower, independent of its gameplay upgrade level.
var TowerJournalEntries = []TowerJournalEntry{
	{"gunner", game.TowerGunner, copytext.Text("defenders.gunner.name"), copytext.Text("defenders.gunner.subtitle"),
		copytext.Text("defenders.gunner.lore"),
		copytext.Text("defenders.gunner.note"), 114, drawGunnerPortrait},
	{"cannonier", game.TowerCannon, copytext.Text("defenders.cannonier.name"), copytext.Text("defenders.cannonier.subtitle"),
		copytext.Text("defenders.cannonier.lore"),
		copytext.Text("defenders.cannonier.note"), 173, drawCannonierPortrait},
	{"frost", game.TowerFrost, copytext.Text("defenders.frost.name"), copytext.Text("defenders.frost.subtitle"),
		copytext.Text("defenders.frost.lore"),
		copytext.Text("defenders.frost.note"), 117, drawFrostPortrait},
	{"ranger", game.TowerSniper, copytext.Text("defenders.ranger.name"), copytext.Text("defenders.ranger.subtitle"),
		copytext.Text("defenders.ranger.lore"),
		copytext.Text("defenders.ranger.note"), 144, drawRangerPortrait},
	{"lightning", game.TowerTesla, copytext.Text("defenders.lightning.name"), copytext.Text("defenders.lightning.subtitle"),
		copytext.Text("defenders.lightning.lore"),
		copytext.Text("defenders.lightning.note"), 153, drawLightningPortrait},
	{"trebuchet", game.TowerMortar, copytext.Text("defenders.trebuchet.name"), copytext.Text("defenders.trebuchet.subtitle"),
		copytext.Text("defenders.trebuchet.lore"),
		copytext.Text("defenders.trebuchet.note"), 180, drawTrebuchetPortrait},
	{"gnolls", game.TowerFlak, copytext.Text("defenders.gnolls.name"), copytext.Text("defenders.gnolls.subtitle"),
		copytext.Text("defenders.gnolls.lore"),
		copytext.Text("defenders.gnolls.note"), 180, drawGnollPortrait},
}

func TowerJournalID(kind game.TowerKind) string {
	for _, e := range TowerJournalEntries {
		if e.Kind == kind {
			return e.ID
		}
	}
	return ""
}

const (
	JournalTowers = iota
	JournalEnemies
	JournalPlaces
	JournalStories
	JournalSectionCount
)

func JournalEntries(section int) []TowerJournalEntry {
	switch section {
	case JournalEnemies:
		return EnemyJournalEntries
	case JournalPlaces:
		return PlaceJournalEntries
	case JournalStories:
		return StoryJournalEntries
	default:
		return TowerJournalEntries
	}
}

func JournalChapter(section int) string {
	return []string{copytext.Text("ui.journal_chapter.i_the_towers"), copytext.Text("ui.journal_chapter.ii_the_guild"), copytext.Text("ui.journal_chapter.iii_the_lair"), copytext.Text("ui.journal_chapter.iv_story_moments")}[max(0, min(section, 3))]
}

type JournalState struct {
	Section  int
	Cursor   int
	Reading  bool
	Scroll   int
	Unlocked map[string]bool
}

// JournalCards is shared with mouse navigation, including the compact list.
func JournalCollectionCapacity(w, h, section int) int {
	n := len(JournalEntries(section))
	if w < 76 || h < 24 {
		return min(n, max(1, h-8))
	}
	return min(n, 8)
}
func JournalCollectionStart(w, h, section, cursor int) int {
	cap := JournalCollectionCapacity(w, h, section)
	return cursor / cap * cap
}
func JournalCards(w, h int, sections ...int) []Rect {
	section := JournalTowers
	if len(sections) > 0 {
		section = sections[0]
	}
	n := JournalCollectionCapacity(w, h, section)
	rects := make([]Rect, n)
	if w < 76 || h < 24 {
		for i := range rects {
			rects[i] = Rect{3, 5 + i, max(0, min(w-6, 29)), 1}
		}
		return rects
	}
	cols := 4
	rows := (n + cols - 1) / cols
	cw := (w - 8) / cols
	ch := (h - 8) / rows
	for i := range rects {
		rects[i] = Rect{4 + (i%cols)*cw, 5 + (i/cols)*ch, cw - 1, ch - 1}
	}
	return rects
}

func RenderJournal(w, h int, st JournalState) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("ui.render_journal.enlarge_to_read_the_journal"), 180, 233, false)
		return f
	}
	entries := JournalEntries(st.Section)
	st.Cursor = max(0, min(st.Cursor, len(entries)-1))
	journalCenter(f, 1, copytext.Text("ui.render_journal.grak_s_journal"), 180, true)
	chapter := JournalChapter(st.Section)
	if st.Section == JournalStories && st.Reading && strings.HasPrefix(entries[st.Cursor].ID, "afterward:") {
		chapter = copytext.Text("ui.render_journal.iv_afterward")
	}
	journalCenter(f, 2, chapter, 240, false)
	if !st.Reading || !st.Unlocked[entries[st.Cursor].ID] {
		drawJournalCollection(f, st)
		journalCenter(f, 4, copytext.Text("ui.render_journal.discover_entries_in_the_campaign"), 240, false)
	} else {
		drawJournalPage(f, st)
	}
	return f
}

func journalCenter(f *Frame, y int, s string, fg int, bold bool) {
	putString(f, max(0, (f.W-len([]rune(s)))/2), y, s, fg, 233, bold)
}

func journalArt(f *Frame, e TowerJournalEntry, r Rect) {
	if r.W < 5 || r.H < 3 {
		return
	}
	// Cannonier's established silhouette needs a wider canvas than the others.
	ph := min(r.H, r.W*2/5)
	pw := ph * 5 / 2
	if e.ID == "cannonier" || e.ID == "rogue" {
		ph = min(r.H, r.W/3)
		pw = ph * 3
	}
	if ph < 3 {
		return
	}
	e.paint(f, r.X+(r.W-pw)/2, r.Y+(r.H-ph)/2, pw, ph)
}

func drawJournalCollection(f *Frame, st JournalState) {
	entries := JournalEntries(st.Section)
	base := JournalCollectionStart(f.W, f.H, st.Section, st.Cursor)
	cards := JournalCards(f.W, f.H, st.Section)
	compact := f.W < 76 || f.H < 24
	for local, r := range cards {
		i := base + local
		if i >= len(entries) {
			break
		}
		e := entries[i]
		active := i == st.Cursor
		fg := 240
		if active {
			fg = 180
		}
		name := copytext.Text("ui.draw_journal_collection.unwritten")
		if st.Unlocked[e.ID] {
			name = e.Name
		}
		if compact {
			mark := "  "
			if active {
				mark = "› "
			}
			putString(f, r.X, r.Y, fitMsg(fmt.Sprintf("%s%d  %s", mark, local+1, name), r.W), fg, 233, active)
			continue
		}
		// Frame only the selected card, leaving the collection light and spacious.
		if active {
			for x := r.X; x < r.X+r.W; x++ {
				f.Put(x, r.Y, '─', 94, 233)
				f.Put(x, r.Y+r.H-1, '─', 94, 233)
			}
			for y := r.Y; y < r.Y+r.H; y++ {
				f.Put(r.X, y, '│', 94, 233)
				f.Put(r.X+r.W-1, y, '│', 94, 233)
			}
			f.Put(r.X, r.Y, '┌', 180, 233)
			f.Put(r.X+r.W-1, r.Y, '┐', 180, 233)
			f.Put(r.X, r.Y+r.H-1, '└', 180, 233)
			f.Put(r.X+r.W-1, r.Y+r.H-1, '┘', 180, 233)
		}
		if st.Unlocked[e.ID] {
			journalArt(f, e, Rect{r.X + 2, r.Y + 1, r.W - 4, r.H - 3})
		}
		label := fitMsg(name, r.W-4)
		putString(f, r.X+max(2, (r.W-len([]rune(label)))/2), r.Y+r.H-2, label, fg, 233, active)
	}
	e := JournalEntries(st.Section)[st.Cursor]
	if compact && f.W >= 52 && st.Unlocked[e.ID] {
		journalArt(f, e, Rect{33, 5, f.W - 36, f.H - 9})
	}
	if len(entries) > len(cards) {
		journalCenter(f, 3, fmt.Sprintf("%d–%d / %d", base+1, min(base+len(cards), len(entries)), len(entries)), 240, false)
	}
	hint := copytext.Format("ui.draw_journal_collection.chapters_arrows_browse_enter_read_esc_return", "chapters", "[ ]", "enter", "enter", "arrows", "arrows", "escape", "esc")
	if f.W < 76 {
		hint = copytext.Format("ui.draw_journal_collection.chapter_arrows_browse_enter_read_esc", "chapters", "[ ]", "enter", "enter", "arrows", "arrows", "escape", "esc")
	}
	if !st.Unlocked[e.ID] {
		journalCenter(f, f.H-3, fitMsg(JournalRequirement(st.Section, st.Cursor), f.W-4), 240, false)
	}
	if f.W < 52 {
		hint = copytext.Format("ui.draw_journal_collection.chapter_enter_read_esc", "chapters", "[ ]", "enter", "enter", "escape", "esc")
	}
	journalCenter(f, f.H-2, hint, 240, false)
}

// journalPageLayout leaves enough room to read prose at narrow widths. The
// same text viewport drives scrolling, so resizing cannot strand the ending.
func journalPageLayout(w, h int) (art, text Rect) {
	// The illustration owns the page width. Keep the prose comfortably narrow
	// below it, with a small scrolling viewport instead of shrinking the art.
	tw := min(90, w-6)
	return Rect{3, 3, w - 6, max(4, h-12)}, Rect{(w - tw) / 2, h - 7, tw, 4}
}

func journalLines(e TowerJournalEntry, width int, sections ...int) []string {
	section := JournalTowers
	if len(sections) > 0 {
		section = sections[0]
	}
	var lines []string
	wrap := func(s string) {
		line := ""
		for _, word := range strings.Fields(s) {
			if line != "" && len([]rune(line+" "+word)) > width {
				lines = append(lines, line)
				line = ""
			}
			if line != "" {
				line += " "
			}
			line += word
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	wrap(e.Lore)
	heading := copytext.Text("ui.journal_lines.in_the_field")
	if section == JournalPlaces {
		heading = copytext.Text("ui.journal_lines.the_way_through")
	}
	if section == JournalStories {
		heading = copytext.Text("ui.journal_lines.remembered")
	}
	lines = append(lines, "", heading)
	wrap(e.Note)
	return lines
}

func JournalMaxScroll(w, h, cursor int, sections ...int) int {
	section := JournalTowers
	if len(sections) > 0 {
		section = sections[0]
	}
	entries := JournalEntries(section)
	_, text := journalPageLayout(w, h)
	if text.W < 1 || text.H < 1 {
		return 0
	}
	return max(0, len(journalLines(entries[max(0, min(cursor, len(entries)-1))], text.W, section))-text.H)
}

func drawJournalPage(f *Frame, st JournalState) {
	e := JournalEntries(st.Section)[st.Cursor]
	art, text := journalPageLayout(f.W, f.H)
	journalArt(f, e, art)
	heading := e.Name
	if f.W >= 76 {
		heading += " · " + e.Subtitle
	}
	x := max(0, (f.W-len([]rune(heading)))/2)
	putString(f, x, text.Y-1, e.Name, e.Accent, 233, true)
	if f.W >= 76 {
		putString(f, x+len([]rune(e.Name)), text.Y-1, " · "+e.Subtitle, 240, 233, false)
	}
	if text.H > 0 && text.W > 0 {
		lines := journalLines(e, text.W, st.Section)
		top := min(max(0, st.Scroll), JournalMaxScroll(f.W, f.H, st.Cursor, st.Section))
		for i := 0; i < text.H && top+i < len(lines); i++ {
			fg, bold := 250, false
			if lines[top+i] == "IN THE FIELD" || lines[top+i] == "THE WAY THROUGH" || lines[top+i] == "REMEMBERED" {
				fg, bold = e.Accent, true
			}
			putString(f, text.X, text.Y+i, lines[top+i], fg, 233, bold)
		}
		if top > 0 {
			f.Put(text.X+text.W, text.Y, '↑', 180, 233)
		}
		if top+text.H < len(lines) {
			f.Put(text.X+text.W, text.Y+text.H-1, '↓', 180, 233)
		}
	}
	hint := copytext.Format("ui.draw_journal_page.chapters_pages_text_esc_collection", "chapters", "[ ]", "scroll", "↑↓", "pages", "←→", "escape", "esc")
	if f.W < 76 {
		hint = copytext.Format("ui.draw_journal_page.chapter_text_esc", "chapters", "[ ]", "scroll", "↑↓", "escape", "esc")
	}
	journalCenter(f, f.H-2, hint, 240, false)
}
