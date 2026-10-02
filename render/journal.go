package render

import (
	"fmt"
	"strings"

	"github.com/0xbenc/tdef/game"
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
	{"gunner", game.TowerGunner, "Orc Gunner", "the patient thunder",
		"He learned patience before he learned the gun. While the guild argues over who will strike first, he settles his elbows, lets the smoke clear, and chooses the next fool. I have never heard him boast. I have heard him reload.",
		"Give him a bend in the road. He will make good use of every second they spend there.", 114, drawGunnerPortrait},
	{"cannonier", game.TowerCannon, "Cannonier", "two hands, one argument",
		"One measures the powder. The other insists that more would be better. They have carried that argument through every chamber of the lair, along with a cannon neither could move alone. When the guild arrives shoulder to shoulder, they briefly agree.",
		"Let them bunch together. A cannonball has little respect for a formation.", 173, drawCannonierPortrait},
	{"frost", game.TowerFrost, "Frost Mage", "the winter we invited in",
		"The garden froze the night she arrived. By morning there was a narrow path through the ice, exactly wide enough for Malgrath's drinking bowl. She says she does not care for dragons. I leave the bowl where she can find it.",
		"Her cold buys time. Put someone beside her who knows how to spend it.", 117, drawFrostPortrait},
	{"ranger", game.TowerSniper, "Ranger", "a promise at a distance",
		"She can tell which guild banner is coming before I can see the boots beneath it. We used to argue about how far away a danger had to be before it became our business. Now she draws the bow, and I trust the answer.",
		"Leave her a long view. Every arrow should reach someone worth the wait.", 144, drawRangerPortrait},
	{"lightning", game.TowerTesla, "Lightning Mage", "the restless sky",
		"He claims the storm was following him long before he came underground. I believe him. Even the iron door rings hum when he passes. Malgrath sleeps through thunder now; it is the silence afterward that makes me look up.",
		"Lightning follows close company. Make the guild regret marching together.", 153, drawLightningPortrait},
	{"trebuchet", game.TowerMortar, "Trebuchet", "the mountain learns to throw",
		"We built it from bridge timbers the guild thought we would need to escape. One of the crew winds while the other watches the far road. There is no hurry in either of them. The stone will arrive, and everyone beneath it will notice.",
		"Give the crew room to see a crowd. Their answer is slow, and very large.", 180, drawTrebuchetPortrait},
	{"gnolls", game.TowerFlak, "Gnoll Slingers", "a handful of trouble",
		"They brought no banner, only pockets full of stones and an argument about whose throw was better. I offered them a place by the road. They offered to make it unpleasant. For once, a bargain meant exactly what it sounded like.",
		"Keep them close to the path. A little stone, thrown often, is still a problem.", 180, drawGnollPortrait},
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
	return []string{"I  ·  THE TOWERS", "II  ·  THE GUILD", "III  ·  THE LAIR", "IV  ·  STORY MOMENTS"}[max(0, min(section, 3))]
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
		putString(f, 1, h/2, "enlarge to read the journal", 180, 233, false)
		return f
	}
	entries := JournalEntries(st.Section)
	st.Cursor = max(0, min(st.Cursor, len(entries)-1))
	journalCenter(f, 1, "GRAK'S JOURNAL", 180, true)
	chapter := JournalChapter(st.Section)
	if st.Section == JournalStories && st.Reading && strings.HasPrefix(entries[st.Cursor].ID, "afterward:") {
		chapter = "IV  ·  AFTERWARD"
	}
	journalCenter(f, 2, chapter, 240, false)
	if !st.Reading || !st.Unlocked[entries[st.Cursor].ID] {
		drawJournalCollection(f, st)
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
		name := "Unwritten"
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
	hint := "[ ] chapters · arrows browse · enter read · esc return"
	if f.W < 76 {
		hint = "[ ] chapter · arrows browse · enter read · esc"
	}
	if !st.Unlocked[e.ID] {
		journalCenter(f, f.H-3, fitMsg(JournalRequirement(st.Section, st.Cursor), f.W-4), 240, false)
	}
	if f.W < 52 {
		hint = "[ ] chapter · enter read · esc"
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
	heading := "IN THE FIELD"
	if section == JournalPlaces {
		heading = "THE WAY THROUGH"
	}
	if section == JournalStories {
		heading = "REMEMBERED"
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
	hint := "[ ] chapters · ←→ pages · ↑↓ text · esc collection"
	if f.W < 76 {
		hint = "[ ] chapter · ↑↓ text · esc"
	}
	journalCenter(f, f.H-2, hint, 240, false)
}
