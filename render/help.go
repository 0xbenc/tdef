package render

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
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
	c := helpControls()
	return []helpPage{
		{copytext.Text("ui.help_guide.start"), []helpSection{
			{copytext.Text("ui.help_guide.goal_title"), copytext.Text("ui.help_guide.goal"), nil},
			{copytext.Text("ui.help_guide.opening_title"), copytext.Text("ui.help_guide.opening"), nil},
			{copytext.Text("ui.help_guide.modes_title"), copytext.Text("ui.help_guide.modes"), nil},
			{copytext.Text("ui.help_guide.saves_title"), copytext.Text("ui.help_guide.saves"), nil},
		}},
		{copytext.Text("ui.help_guide.controls"), []helpSection{
			{copytext.Text("ui.help_guide.placement_title"), copytext.Text("ui.help_guide.placement_tip"), c[:2]},
			{copytext.Text("ui.help_guide.manage_title"), copytext.Text("ui.help_guide.manage_tip"), c[2:5]},
			{copytext.Text("ui.help_guide.pace_title"), copytext.Text("ui.help_guide.pace_tip"), [][2]string{c[5], c[6], c[9]}},
			{copytext.Text("ui.help_guide.other_title"), copytext.Text("ui.help_guide.other_tip"), [][2]string{c[7], c[8], c[10], c[11], c[12], c[13], c[14]}},
		}},
		{copytext.Text("ui.help_guide.counters"), []helpSection{
			{copytext.Text("ui.help_guide.rogue_title"), copytext.Text("ui.help_guide.rogue"), nil},
			{copytext.Text("ui.help_guide.paladin_title"), copytext.Text("ui.help_guide.paladin"), nil},
			{copytext.Text("ui.help_guide.centurion_title"), copytext.Text("ui.help_guide.centurion"), nil},
			{copytext.Text("ui.help_guide.others_title"), copytext.Text("ui.help_guide.others"), nil},
		}},
		{copytext.Text("ui.help_guide.specialists"), []helpSection{
			{copytext.Text("ui.help_guide.forge_title"), copytext.Text("ui.help_guide.forge"), nil},
			{copytext.Text("ui.help_guide.ogre_title"), copytext.Text("ui.help_guide.ogre"), nil},
			{copytext.Text("ui.help_guide.sapper_title"), copytext.Text("ui.help_guide.sapper"), nil},
			{copytext.Text("ui.help_guide.witch_title"), copytext.Text("ui.help_guide.witch"), nil},
		}},
		{copytext.Text("ui.help_guide.lair"), []helpSection{
			{copytext.Text("ui.help_guide.route_title"), copytext.Text("ui.help_guide.route"), nil},
			{copytext.Text("ui.help_guide.relic_title"), copytext.Text("ui.help_guide.relic"), nil},
			{copytext.Text("ui.help_guide.seed_title"), copytext.Text("ui.help_guide.seed"), nil},
			{copytext.Text("ui.help_guide.replay_title"), copytext.Text("ui.help_guide.replay"), nil},
		}},
	}
}

func HelpPageCount() int { return len(helpPages()) }

type helpLine struct {
	text string
	fg   int
	bold bool
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
		lines := []helpLine{{"", 252, false}, {section.title, 220, true}, {"", 252, false}}
		for _, control := range section.controls {
			lines = append(lines, helpLine{control[0], 245, false})
			for _, line := range recruitmentLines(control[1], cardW-4) {
				lines = append(lines, helpLine{line, 167, true})
			}
		}
		if len(section.controls) > 0 {
			lines = append(lines, helpLine{"", 252, false})
		}
		for _, line := range recruitmentLines(section.text, cardW-4) {
			lines = append(lines, helpLine{line, 252, false})
		}
		lines = append(lines, helpLine{"", 252, false})
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
