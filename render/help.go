package render

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"strings"
)

type HelpState struct{ Page, Scroll int }
type helpSection struct {
	title, text string
	controls    [][2]string
}
type helpPage struct {
	title    string
	sections []helpSection
}

func helpPages() []helpPage {
	return []helpPage{
		{copytext.Text("ui.help_guide.start"), []helpSection{
			{copytext.Text("ui.help_guide.goal_title"), copytext.Text("ui.help_guide.goal"), nil},
			{copytext.Text("ui.help_guide.opening_title"), copytext.Text("ui.help_guide.opening"), nil},
			{copytext.Text("ui.help_guide.modes_title"), copytext.Text("ui.help_guide.modes"), nil},
		}},
		{copytext.Text("ui.help_guide.controls"), []helpSection{
			{copytext.Text("ui.help_guide.placement_title"), copytext.Text("ui.help_guide.placement_tip"), [][2]string{
				{copytext.Text("ui.help_bindings.move"), "Arrows|WASD"},
				{copytext.Text("ui.help_bindings.choose"), "1-7"},
				{copytext.Text("ui.help_bindings.roster"), "[|]"},
				{copytext.Text("ui.help_bindings.place"), "Enter|Click"},
			}},
			{copytext.Text("ui.help_guide.manage_title"), copytext.Text("ui.help_guide.manage_tip"), [][2]string{
				{copytext.Text("ui.help_bindings.select"), "Tab"},
				{copytext.Text("ui.help_bindings.upgrade"), "U"},
				{copytext.Text("ui.help_bindings.sell"), "X"},
				{copytext.Text("ui.help_bindings.target"), "T"},
				{copytext.Text("ui.help_bindings.rotate"), "R"},
			}},
			{copytext.Text("ui.help_guide.pace_title"), copytext.Text("ui.help_guide.pace_tip"), [][2]string{
				{copytext.Text("ui.help_bindings.pause"), "P"},
				{copytext.Text("ui.help_bindings.speed"), "F"},
				{copytext.Text("ui.help_bindings.wave"), "N"},
			}},
			{copytext.Text("ui.help_guide.other_title"), copytext.Text("ui.help_guide.other_tip"), [][2]string{
				{copytext.Text("ui.help_bindings.help"), "H"},
				{copytext.Text("ui.help_bindings.journal"), "J"},
				{copytext.Text("ui.help_bindings.cancel"), "Esc"},
				{copytext.Text("ui.help_bindings.quit"), "Q"},
			}},
			{copytext.Text("ui.help_guide.lair_controls_title"), copytext.Text("ui.help_guide.lair_controls_tip"), [][2]string{
				{copytext.Text("ui.help_bindings.descend"), "Enter"},
				{copytext.Text("ui.help_bindings.renown"), "Tab"},
				{copytext.Text("ui.help_bindings.relic"), "T"},
				{copytext.Text("ui.help_bindings.films"), "I|E"},
				{copytext.Text("ui.help_bindings.seed"), "0-9"},
			}},
		}},
		{copytext.Text("ui.help_guide.counters"), []helpSection{
			{copytext.Text("ui.help_guide.rogue_title"), copytext.Text("ui.help_guide.rogue"), nil},
			{copytext.Text("ui.help_guide.paladin_title"), copytext.Text("ui.help_guide.paladin"), nil},
			{copytext.Text("ui.help_guide.centurion_title"), copytext.Text("ui.help_guide.centurion"), nil},
		}},
		{copytext.Text("ui.help_guide.specialists"), []helpSection{
			{copytext.Text("ui.help_guide.forge_title"), copytext.Text("ui.help_guide.forge"), nil},
			{copytext.Text("ui.help_guide.ogre_title"), copytext.Text("ui.help_guide.ogre"), nil},
			{copytext.Text("ui.help_guide.sapper_title"), copytext.Text("ui.help_guide.sapper"), nil},
			{copytext.Text("ui.help_guide.witch_title"), copytext.Text("ui.help_guide.witch"), nil},
		}},
	}
}

func HelpPageCount() int { return len(helpPages()) }

type helpLine struct {
	text string
	fg   int
	bold bool
	keys string
}

func helpContent(w int, page helpPage) *Frame {
	width := min(112, max(8, w-6))
	columns := 1
	if width >= 94 {
		columns = 2
	}
	cardW := (width - (columns-1)*3) / columns
	var cards [][]helpLine
	for _, section := range page.sections {
		lines := []helpLine{{"", 252, false, ""}, {section.title, 220, true, ""}, {"", 252, false, ""}}
		keyWidth := 0
		for _, control := range section.controls {
			keyWidth = max(keyWidth, helpKeyWidth(control[1]))
		}
		for _, control := range section.controls {
			textWidth := cardW - 4 - keyWidth - 2
			if textWidth < 12 {
				lines = append(lines, helpLine{keys: control[1]})
				for _, line := range recruitmentLines(control[0], cardW-4) {
					lines = append(lines, helpLine{line, 252, false, ""})
				}
			} else {
				for i, line := range recruitmentLines(control[0], textWidth) {
					keys := ""
					if i == 0 {
						keys = control[1]
					}
					lines = append(lines, helpLine{strings.Repeat(" ", keyWidth+2) + line, 252, false, keys})
				}
			}
		}
		if len(section.controls) > 0 {
			lines = append(lines, helpLine{"", 252, false, ""})
		}
		for _, line := range recruitmentLines(section.text, cardW-4) {
			lines = append(lines, helpLine{line, 252, false, ""})
		}
		lines = append(lines, helpLine{"", 252, false, ""})
		cards = append(cards, lines)
	}
	height := 0
	for i := 0; i < len(cards); i += columns {
		rowH := len(cards[i])
		if columns == 2 && i+1 < len(cards) {
			rowH = max(rowH, len(cards[i+1]))
		}
		height += rowH + 1
	}
	f := &Frame{W: width, H: max(0, height-1), C: make([]Cell, width*max(0, height-1))}
	y := 0
	for i := 0; i < len(cards); i += columns {
		rowH := len(cards[i])
		if columns == 2 && i+1 < len(cards) {
			rowH = max(rowH, len(cards[i+1]))
		}
		for col := 0; col < columns && i+col < len(cards); col++ {
			x := col * (cardW + 3)
			for yy := 0; yy < rowH; yy++ {
				for xx := 0; xx < cardW; xx++ {
					f.Set(x+xx, y+yy, Cell{R: ' ', BG: 235, FG: 252})
				}
			}
			for line, text := range cards[i+col] {
				putString(f, x+2, y+line, text.text, text.fg, 235, text.bold)
				if page.sections[i+col].title == copytext.Text("ui.help_guide.manage_title") {
					for _, target := range []string{"First", "Strongest", "Closest"} {
						if start := strings.Index(text.text, target); start >= 0 {
							targetX := x + 2 + len([]rune(text.text[:start]))
							putString(f, targetX, y+line, target, text.fg, 235, true)
						}
					}
				}
				keyX := x + 2
				if text.keys != "" {
					for _, key := range strings.Split(text.keys, "|") {
						label := " " + key + " "
						putString(f, keyX, y+line, label, 231, 24, true)
						keyX += len([]rune(label)) + 1
					}
				}
			}
			for xx := 0; xx < cardW; xx++ {
				f.Set(x+xx, y, Cell{R: '─', FG: 240, BG: 235})
			}
		}
		y += rowH + 1
	}
	return f
}

func helpViewport(h int) (tabsY, top, visible int) {
	if h < 24 {
		return 2, 4, max(1, h-7)
	}
	return 4, 6, max(1, h-9)
}

func HelpMaxScroll(w, h, page int) int {
	pages := helpPages()
	page = max(0, min(page, len(pages)-1))
	_, _, visible := helpViewport(h)
	return max(0, helpContent(w, pages[page]).H-visible)
}

func RenderHelpPage(w, h int, st HelpState, pal Colors) *Frame {
	f := screenBox(w, h, copytext.Text("ui.render_help.grak_s_ledger"), []fseg{
		{key: "[ ] / left right", text: copytext.Text("ui.help_guide.topic")},
		{key: "up down", text: copytext.Text("ui.help_guide.scroll")},
		{key: "esc", text: copytext.Text("ui.render_help.back")},
	}, true, pal)
	pages := helpPages()
	page := max(0, min(st.Page, len(pages)-1))
	tabsY, top, visible := helpViewport(h)
	if h >= 24 {
		centerPut(f, 2, copytext.Text("ui.render_help.how_to_hold_the_lair_against_twenty"), 252, false)
	}
	if page == 3 && h >= 24 {
		centerPut(f, 2, copytext.Text("ui.help_guide.limit"), 252, false)
	}
	var tabs []owRun
	for i, p := range pages {
		if i > 0 {
			tabs = append(tabs, owRun{"  ", 240, false})
		}
		fg := 245
		if i == page {
			fg = 220
		}
		tabs = append(tabs, owRun{fmt.Sprint(i + 1), 167, true}, owRun{" " + p.title, fg, i == page})
	}
	putCenteredRuns(f, tabsY, tabs)
	content := helpContent(w, pages[page])

	offset := max(0, min(st.Scroll, max(0, content.H-visible)))
	x := (w - content.W) / 2
	for yy := 0; yy < visible && yy+offset < content.H; yy++ {
		for xx := 0; xx < content.W; xx++ {
			f.Set(x+xx, top+yy, content.C[(yy+offset)*content.W+xx])
		}
	}
	position := copytext.Format("ui.help_guide.position", "page", fmt.Sprint(page+1), "pages", fmt.Sprint(len(pages)), "first", fmt.Sprint(offset+1), "last", fmt.Sprint(min(content.H, offset+visible)), "lines", fmt.Sprint(content.H))
	centerPut(f, h-3, position, 245, false)
	return f
}

// HelpTopicAt shares the centered tab geometry with RenderHelpPage.
func HelpTopicAt(w, h, x, y int) int {
	tabsY, _, _ := helpViewport(h)
	if y != tabsY {
		return -1
	}
	pages := helpPages()
	total := 0
	for _, p := range pages {
		total += 2 + len([]rune(p.title))
	}
	total += (len(pages) - 1) * 2
	left := max(1, (w-total)/2)
	for i, p := range pages {
		width := 2 + len([]rune(p.title))
		if x >= left && x < left+width {
			return i
		}
		left += width + 2
	}
	return -1
}

func helpKeyWidth(keys string) int {
	width := 0
	for _, key := range strings.Split(keys, "|") {
		width += len([]rune(key)) + 3
	}
	return max(0, width-1)
}
