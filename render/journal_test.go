package render

import (
	"reflect"
	"strings"
	"testing"

	"github.com/0xbenc/tdef/game"
)

func TestJournalPagesResponsiveAndComplete(t *testing.T) {
	unlocked := map[string]bool{}
	kinds := map[game.TowerKind]bool{}
	for _, e := range TowerJournalEntries {
		if e.ID == "" || kinds[e.Kind] || TowerJournalID(e.Kind) != e.ID {
			t.Fatal("missing or duplicated tower identity")
		}
		kinds[e.Kind] = true
		unlocked[e.ID] = true
	}
	if len(kinds) != int(game.TowerCount) {
		t.Fatal("missing tower")
	}
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		for i, e := range TowerJournalEntries {
			st := JournalState{Cursor: i, Reading: true, Unlocked: unlocked}
			f := RenderJournal(w, h, st)
			if len(f.C) != w*h || !reflect.DeepEqual(f, RenderJournal(w, h, st)) {
				t.Fatal("page not static or responsive")
			}
			if !strings.Contains(f.Text(), e.Name) || !strings.Contains(f.Text(), "esc") {
				t.Fatalf("%v %s: missing title or controls", size, e.Name)
			}
			_, r := journalPageLayout(w, h)
			lines := journalLines(e, r.W)
			for _, line := range lines {
				if len([]rune(line)) > r.W {
					t.Fatal("prose escaped viewport")
				}
			}
			// Every lore line and field note can be reached by scrolling.
			seen := map[string]bool{}
			for top := 0; top <= JournalMaxScroll(w, h, i); top++ {
				st.Scroll = top
				text := RenderJournal(w, h, st).Text()
				for _, line := range lines {
					if strings.Contains(text, line) {
						seen[line] = true
					}
				}
			}
			for _, line := range lines {
				if !seen[line] {
					t.Fatalf("%v %s: inaccessible prose %q", size, e.Name, line)
				}
			}
		}
	}
}

func TestJournalCollectionLocksHideNamesAndArt(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {120, 40}} {
		w, h := size[0], size[1]
		f := RenderJournal(w, h, JournalState{})
		for _, e := range TowerJournalEntries {
			if strings.Contains(f.Text(), e.Name) {
				t.Fatal("locked entry spoiled identity")
			}
		}
		if !strings.Contains(f.Text(), "Unwritten") {
			t.Fatal("missing blank pages")
		}
		for _, c := range f.C {
			if c.R == '█' || c.R == '▀' || c.R == '▄' {
				t.Fatal("locked entry revealed portrait")
			}
		}
		f = RenderJournal(w, h, JournalState{Unlocked: map[string]bool{"frost": true}})
		if !strings.Contains(f.Text(), "Frost Mage") {
			t.Fatal("discovery missing from collection")
		}
	}
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		RenderJournal(size[0], size[1], JournalState{})
	}
}

func TestExpandedJournalCatalogAndResponsivePages(t *testing.T) {
	for section := JournalEnemies; section <= JournalStories; section++ {
		entries := JournalEntries(section)
		unlocked := map[string]bool{}
		for _, e := range entries {
			if unlocked[e.ID] {
				t.Fatal("duplicate page identity")
			}
			unlocked[e.ID] = true
		}
		for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {120, 46}} {
			for i, e := range entries {
				st := JournalState{Section: section, Cursor: i, Reading: true, Unlocked: unlocked}
				f := RenderJournal(size[0], size[1], st)
				if len(f.C) != size[0]*size[1] || !reflect.DeepEqual(f, RenderJournal(size[0], size[1], st)) {
					t.Fatal("expanded page is not static and responsive")
				}
				_, text := journalPageLayout(size[0], size[1])
				seen := map[string]bool{}
				lines := journalLines(e, text.W, section)
				for top := 0; top <= JournalMaxScroll(size[0], size[1], i, section); top++ {
					st.Scroll = top
					frame := RenderJournal(size[0], size[1], st).Text()
					for _, line := range lines {
						if strings.Contains(frame, line) {
							seen[line] = true
						}
					}
				}
				for _, line := range lines {
					if !seen[line] {
						t.Fatalf("%v %s: inaccessible lore %q", size, e.ID, line)
					}
				}
				if section == JournalStories && !unlocked[e.ID] {
					t.Fatal("missing story")
				}
			}
		}
	}
	if len(EnemyJournalEntries) != int(game.EnemyCount) || len(PlaceJournalEntries) != 6 || len(StoryJournalEntries) != 22 {
		t.Fatal("incomplete journal roster")
	}
	for k := game.EnemyKind(0); k < game.EnemyCount; k++ {
		if EnemyJournalID(k) == "" {
			t.Fatal("missing enemy mapping")
		}
	}
}

func TestLockedStoryRequirementsDoNotSpoilTitles(t *testing.T) {
	for _, index := range []int{0, 7, 17, 18, 21} {
		f := RenderJournal(120, 40, JournalState{Section: JournalStories, Cursor: index})
		if !strings.Contains(f.Text(), JournalRequirement(JournalStories, index)) {
			t.Fatal("missing exact unlock condition")
		}
		for _, e := range StoryJournalEntries {
			if strings.Contains(f.Text(), e.Name) {
				t.Fatal("locked scene revealed its title")
			}
		}
		for _, c := range f.C {
			if c.R == '█' || c.R == '▀' || c.R == '▄' {
				t.Fatal("locked scene revealed its illustration")
			}
		}
	}
}

func TestEveryDifficultyMomentHasDistinctIllustration(t *testing.T) {
	// Compare only the illustration, so different titles cannot mask reused art.
	art := func(e TowerJournalEntry) []Cell {
		f := &Frame{W: 100, H: 40, C: make([]Cell, 4000)}
		for i := range f.C {
			f.C[i] = Cell{R: ' ', BG: 233}
		}
		e.paint(f, 0, 0, 100, 40)
		return f.C
	}
	for floor := 0; floor < 6; floor++ {
		for a := 0; a < 3; a++ {
			for b := a + 1; b < 3; b++ {
				if reflect.DeepEqual(art(StoryJournalEntries[floor*3+a]), art(StoryJournalEntries[floor*3+b])) {
					t.Fatalf("floor %d: difficulty illustrations reused", floor)
				}
			}
		}
	}
	for a := 18; a < 22; a++ {
		for b := a + 1; b < 22; b++ {
			if reflect.DeepEqual(art(StoryJournalEntries[a]), art(StoryJournalEntries[b])) {
				t.Fatal("bonus illustrations reused")
			}
		}
	}
}
