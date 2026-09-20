package render

import (
	"fmt"
	"sort"

	"tdef/game"
)

// Rect is a frame-space rectangle used for mouse hit-testing. The renderers
// and the Rects helpers share one layout function, so clicks can never
// drift from the pixels.
type Rect struct{ X, Y, W, H int }

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// fseg is one segment of a screen's footer group: a bold hotkey and its
// label, e.g. ("↑↓", " move").
type fseg struct {
	key, text string
	blink     bool
}

// screenBox returns a w×h frame wrapped in a full-window rounded border
// (240), the screen title embedded in the top border (btop's ┐title┌
// grammar) and a footer segment group centered in the bottom border
// (┘group└). Menu screens are terminal-sized, so a resize is picked up on
// the next frame.
func screenBox(w, h int, title string, footer []fseg, lit bool, pal Colors) *Frame {
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' '}
	}
	drawRoundedBox(f, 0, 0, w, h, pal.Path)
	embedSegment(f, 0, 2, title, '┐', '┌', pal.Path, pal.Bright, true)
	drawFooter(f, footer, lit, pal)
	return f
}

// drawFooter centers the footer segment group in the bottom border:
// ┘key text ─ key text└, keys in bold hotkey red, text dim, segments
// joined by " ─ ". A blinking segment dims as a unit when !lit.
func drawFooter(f *Frame, segs []fseg, lit bool, pal Colors) {
	type run struct {
		s  string
		fg int
		b  bool
	}
	var runs []run
	for i, s := range segs {
		if i > 0 {
			runs = append(runs, run{" ─ ", pal.Path, false})
		}
		kfg, kbold := 167, true
		tfg := pal.Dim
		if s.blink && !lit {
			kfg, kbold, tfg = pal.Dim, false, pal.Dim
		}
		runs = append(runs, run{s.key, kfg, kbold}, run{s.text, tfg, false})
	}
	total := 2 // brackets
	for _, r := range runs {
		total += len([]rune(r.s))
	}
	x := (f.W - total) / 2
	if x < 1 {
		x = 1
	}
	f.Set(x, f.H-1, Cell{R: '┘', FG: pal.Path})
	x++
	for _, r := range runs {
		for _, ch := range r.s {
			f.Set(x, f.H-1, Cell{R: ch, FG: r.fg, Bold: r.b})
			x++
		}
	}
	f.Set(x, f.H-1, Cell{R: '└', FG: pal.Path})
}

// screenOff is the top offset of the 19-row content band inside a taller
// frame (0 at the minimum 62x19).
func screenOff(h int) int {
	off := (h - 19) / 2
	if off < 0 {
		off = 0
	}
	return off
}

func centerPut(f *Frame, y int, s string, fg int, bold bool) {
	putString(f, (f.W-len([]rune(s)))/2, y, s, fg, 0, bold)
}

// drawSubBox draws a content-sized rounded sub-box with a btop-style label
// (245 bold) embedded in its top border.
func drawSubBox(f *Frame, x, y, bw, bh int, title string, pal Colors) {
	drawRoundedBox(f, x, y, bw, bh, pal.Path)
	embedSegment(f, y, x+1, title, '┐', '┌', pal.Path, pal.Dim, true)
}

// drawMapPreview renders map terrain plus the spawn/exit markers at layout
// l. Shared by the playfield and the level-select preview.
func drawMapPreview(f *Frame, m *game.Map, pal Colors, l Layout) {
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			switch m.At(game.Vec{X: x, Y: y}) {
			case game.CellWall:
				l.block(f, x, y, Cell{R: ' ', FG: 0, BG: pal.Wall})
			case game.CellPath:
				l.block(f, x, y, Cell{R: '·', FG: pal.Path, BG: 0})
			case game.CellGrass:
				l.block(f, x, y, Cell{R: ' ', FG: 0, BG: pal.Grass})
			}
		}
	}
	sx, sy := l.center(m.Spawn.X, m.Spawn.Y)
	f.Set(sx, sy, Cell{R: '▶', FG: 46, Bold: true})
	ex, ey := l.center(m.Exit.X, m.Exit.Y)
	f.Set(ex, ey, Cell{R: 'E', FG: 196, Bold: true})
}

// ---------------------------------------------------------------- title

// titleLetters is the 5x9 slab grid for "TDEF".
var titleLetters = [4][5]string{
	{"XXXXXXXXX", "XXXXXXXXX", "XXXXXXXXX", "...XXX...", "...XXX..."},
	{"XXXXXXX..", "XXX...XXX", "XXX...XXX", "XXX...XXX", "XXXXXXX.."},
	{"XXXXXXXXX", "XXX......", "XXXXXXXX.", "XXX......", "XXXXXXXXX"},
	{"XXXXXXXXX", "XXX......", "XXXXXXXX.", "XXX......", "XXX......"},
}

var titleColors = [4]int{46, 220, 203, 171}

// titleBest returns the highest score and its key (0, "" when empty).
func titleBest(scores map[string]int) (int, string) {
	best, name := 0, ""
	for k, v := range scores {
		if v > best || (v == best && name != "" && k < name) {
			best, name = v, k
		}
	}
	return best, name
}

// RenderTitle draws the animated title screen. It is a pure function of
// (w, h, frame, scores): the same frame index always yields the same frame,
// so the animation is deterministic and unit-testable. `frame` is the 30fps
// tick counter.
func RenderTitle(w, h, frame int, scores map[string]int, pal Colors) *Frame {
	lit := (frame/15)%2 == 0
	f := screenBox(w, h, "TDEF", []fseg{
		{key: "[enter]", text: " start", blink: true},
		{key: "q", text: " quit"},
	}, lit, pal)
	off := screenOff(h)
	// Band rows (19-row layout): 2-6 logo, 7 shadow, 8 tagline, 9-12
	// demo sub-box, 13-14 roster, 16 best.
	row := func(y int) int { return off + y }
	put := func(y int, s string, fg int, bold bool) {
		if yy := row(y); yy > 0 && yy < h-1 {
			centerPut(f, yy, s, fg, bold)
		}
	}
	const logoW = 42 // 4 letters * 9 + 3 gaps * 2
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	// Drop shadow under the logo's bottom row (drawn first, so the letters
	// win any overlap).
	for li, letters := range titleLetters {
		lx := x0 + li*11
		for ci, ch := range letters[4] {
			if ch == 'X' {
				f.Set(lx+ci+1, row(6)+1, Cell{R: '░', FG: 238})
			}
		}
	}
	bevel := [5]rune{'█', '▓', '▓', '▓', '▒'}
	for li, letters := range titleLetters {
		lx := x0 + li*11
		for ri, lrow := range letters {
			for ci, ch := range lrow {
				if ch != 'X' {
					continue
				}
				f.Set(lx+ci, row(2+ri), Cell{R: bevel[ri], FG: titleColors[li], Bold: true})
			}
		}
	}
	put(8, "— terminal tower defense —", 245, false)
	drawTitleDemo(f, w, h, frame, row, pal)
	drawTitleRoster(f, w, h, pal, row)
	best, name := titleBest(scores)
	if best > 0 {
		put(16, fmt.Sprintf("★ best %d — %s", best, name), 220, false)
	} else {
		put(16, " no scores yet ", 238, false)
	}
	return f
}

// drawTitleDemo animates a miniature battle inside a full-width [ BATTLE ]
// sub-box on band rows 9-12: an enemy walks the path from the spawn toward
// E, the tower G beams it while it is in range, and E flashes while the
// enemy "leaks". Two render ticks per step (~67ms) keeps the motion
// readable at 30fps.
func drawTitleDemo(f *Frame, w, h, frame int, row func(int) int, pal Colors) {
	const towerX = 20
	const towerRange = 8
	walk := w - 7 // steps from x=3 to x=w-4
	if walk < 10 {
		walk = 10
	}
	const leakSteps = 6
	cycle := walk + leakSteps
	c := (frame / 2) % cycle
	ex, leak := 3, false
	if c < walk {
		ex = 3 + c
	} else {
		ex, leak = w-4, true
	}
	y, ty := row(10), row(11)
	if row(9) <= 0 || row(12) >= h-1 {
		return
	}
	drawSubBox(f, 1, row(9), w-2, 4, "BATTLE", pal)
	for x := 2; x < w-2; x++ {
		f.Set(x, y, Cell{R: '·', FG: 240})
	}
	f.Set(2, y, Cell{R: '▶', FG: 46, Bold: true})
	if ty > 0 && ty < h-1 {
		for x := towerX - 4; x <= towerX+4; x++ {
			f.Set(x, ty, Cell{R: ' ', FG: 0, BG: 234})
		}
		f.Set(towerX, ty, Cell{R: 'G', FG: 46, Bold: true, BG: 234})
	}
	inRange := !leak && ex >= towerX-towerRange && ex <= towerX+towerRange
	if inRange {
		lo, hi := ex, towerX
		if lo > hi {
			lo, hi = hi, lo
		}
		for x := lo + 1; x < hi; x++ {
			f.Set(x, y, Cell{R: '+', FG: 255})
		}
	}
	efg, ebold := 213, false
	if inRange {
		efg, ebold = 231, true
	}
	f.Set(ex, y, Cell{R: 'o', FG: efg, Bold: ebold})
	if leak {
		f.Set(w-3, y, Cell{R: 'E', FG: 231, BG: 196, Bold: true})
	} else {
		f.Set(w-3, y, Cell{R: 'E', FG: 196, Bold: true})
	}
}

// drawTitleRoster lists the tower and enemy glyphs in their colors, so the
// title doubles as a legend.
func drawTitleRoster(f *Frame, w, h int, pal Colors, row func(int) int) {
	// Both lines share one grid (label column 9 wide, glyphs 4 apart,
	// anchored to the 8-glyph enemy line) so the columns line up.
	const labelW, stride, maxGlyphs = 9, 4, 8
	widest := labelW + maxGlyphs*stride
	line := func(y int, label string, glyphs []rune, colors []int) {
		yy := row(y)
		if yy <= 0 || yy >= h-1 {
			return
		}
		x0 := (w - widest) / 2
		if x0 < 0 {
			x0 = 0
		}
		putString(f, x0, yy, label, 245, 0, false)
		for i, g := range glyphs {
			f.Set(x0+labelW+stride*i, yy, Cell{R: g, FG: colors[i], Bold: true})
		}
	}
	line(13, " towers  ",
		[]rune{'G', 'C', 'F', 'S', 'T', 'M', 'L'},
		[]int{pal.Tower[0], pal.Tower[1], pal.Tower[2], pal.Tower[3], pal.Tower[4], pal.Tower[5], pal.Tower[6]})
	line(14, " enemies  ",
		[]rune{'o', 'r', 'g', 't', 's', 'B', 'w', 'D'},
		[]int{pal.Enemy[0], pal.Enemy[1], pal.Enemy[2], pal.Enemy[3], pal.Enemy[4], pal.Enemy[5], pal.Enemy[6], pal.Enemy[7]})
}

// ---------------------------------------------------------------- menu

// MenuItems is the main menu, in order.
var MenuItems = []string{"Start", "Help", "High Scores", "Quit"}

// menuLayout returns the item rows of the main menu. The 4-item block (7
// rows tall at 2-row stride) is centered in the content band so the menu
// sits in the middle of the frame at any height.
func menuLayout(h int) []int {
	off := screenOff(h)
	top, bottom := off+1, off+17
	// Block height is 7 rows (4 items, 2-row stride); split the leftover
	// rows as evenly as possible above and below it.
	start := top + (bottom-top+1-7)/2
	if start < top {
		start = top
	}
	var items []int
	for i := 0; i < len(MenuItems); i++ {
		items = append(items, start+2*i)
	}
	return items
}

// MenuRects returns the hit-test rectangle for each main-menu item.
func MenuRects(w, h int) []Rect {
	items := menuLayout(h)
	maxW := 0
	for _, it := range MenuItems {
		if len(it) > maxW {
			maxW = len(it)
		}
	}
	x := (w - (maxW + 3)) / 2
	out := make([]Rect, len(items))
	for i, y := range items {
		out[i] = Rect{X: x, Y: y, W: maxW + 3, H: 1}
	}
	return out
}

func RenderMenu(w, h, sel int, pal Colors) *Frame {
	f := screenBox(w, h, "MAIN MENU", []fseg{
		{key: "↑↓", text: " move"},
		{key: "enter", text: " select"},
		{key: "q", text: " quit"},
	}, true, pal)
	items := menuLayout(h)
	for i, name := range MenuItems {
		if i >= len(items) {
			break
		}
		y := items[i]
		if y <= 0 || y >= h-1 {
			continue
		}
		if i == sel {
			centerPut(f, y, " ▸ "+name, pal.Bright, true)
		} else {
			centerPut(f, y, "   "+name, pal.Dim, false)
		}
	}
	return f
}

// ---------------------------------------------------------------- help

func RenderHelp(w, h int, pal Colors) *Frame {
	f := screenBox(w, h, "HELP", []fseg{{key: "esc", text: " back"}}, true, pal)
	off := screenOff(h)
	rows := [][2]string{
		{"move", "arrows / wasd"},
		{"place", "1-7 pick · enter or click"},
		{"upgrade", "u"},
		{"sell", "x"},
		{"target", "t"},
		{"wave", "n (early = bonus gold)"},
		{"pause", "p"},
		{"speed", "f or wheel"},
		{"cancel", "esc"},
		{"quit", "q"},
	}
	// Fixed-width two-column table: labels right-aligned in an 8-col field,
	// values in one column, so the rows line up instead of ragged-centering.
	tableW := 0
	for _, r := range rows {
		if n := len(fmt.Sprintf("%8s   %s", r[0], r[1])); n > tableW {
			tableW = n
		}
	}
	x0 := (w - tableW) / 2
	if x0 < 1 {
		x0 = 1
	}
	for i, r := range rows {
		y := off + 4 + i
		if y <= 0 || y >= h-1 {
			continue
		}
		putString(f, x0+8-len([]rune(r[0])), y, r[0], pal.Dim, 0, false)
		putString(f, x0+11, y, r[1], pal.Bright, 0, false)
	}
	return f
}

// ---------------------------------------------------------- high scores

// ScoreEntry is one row of the high-score table.
type ScoreEntry struct {
	Name  string
	Score int
}

// SortScores returns the table's entries sorted by score descending, then
// name ascending, so the on-screen order is deterministic.
func SortScores(scores map[string]int) []ScoreEntry {
	out := make([]ScoreEntry, 0, len(scores))
	for k, v := range scores {
		out = append(out, ScoreEntry{Name: k, Score: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// RenderHighScores draws the hiscore table. top is the scroll offset (first
// visible entry); it is clamped here, so callers can pass a stale value.
func RenderHighScores(w, h, top int, scores map[string]int, pal Colors) *Frame {
	f := screenBox(w, h, "HIGH SCORES", []fseg{
		{key: "↑↓", text: " scroll"},
		{key: "esc", text: " back"},
	}, true, pal)
	off := screenOff(h)
	const tableW = 37 // len("  %3d %-20s %10d")
	x0 := (w - tableW) / 2
	if x0 < 1 {
		x0 = 1
	}
	bodyTop := off + 4
	bodyBottom := h - 3
	visible := bodyBottom - bodyTop + 1
	if visible < 1 {
		visible = 1
	}
	entries := SortScores(scores)
	if top < 0 {
		top = 0
	}
	if maxTop := len(entries) - visible; top > maxTop && maxTop >= 0 {
		top = maxTop
	}
	if hdr := off + 2; hdr > 0 && hdr < h-1 {
		putString(f, x0, hdr, fmt.Sprintf("  %3s %-20s %10s", " #", "NAME", "SCORE"), pal.Path, 0, false)
		if top > 0 {
			f.Set(x0+37, hdr, Cell{R: '▲', FG: pal.Path})
		}
		if top+visible < len(entries) {
			f.Set(x0+39, hdr, Cell{R: '▼', FG: pal.Path})
		}
	}
	if len(entries) == 0 {
		if bodyTop > 0 && bodyTop < h-1 {
			centerPut(f, bodyTop, " no scores yet ", 238, false)
		}
		return f
	}
	for i := 0; i < visible && top+i < len(entries); i++ {
		e := entries[top+i]
		y := bodyTop + i
		if y > bodyBottom || y <= 0 || y >= h-1 {
			break
		}
		name := e.Name
		if r := []rune(name); len(r) > 20 {
			name = string(r[:20])
		}
		rankFG := pal.Path
		switch top + i {
		case 0:
			rankFG = 220
		case 1:
			rankFG = 251
		case 2:
			rankFG = 180
		}
		nameFG := pal.Dim
		if top+i < 3 {
			nameFG = pal.Bright
		}
		putString(f, x0, y, fmt.Sprintf("  %3d", top+i+1), rankFG, 0, false)
		putString(f, x0+6, y, name, nameFG, 0, false)
		putString(f, x0+27, y, fmt.Sprintf("%10d", e.Score), pal.Bright, 0, false)
	}
	return f
}

// ---------------------------------------------------------- level select

// Difficulties is the display order for the difficulty selector.
var Difficulties = []game.Difficulty{game.Easy, game.Normal, game.Hard}

// DiffName returns the display name for a Difficulties index.
func DiffName(i int) string {
	if i < 0 || i >= len(Difficulties) {
		i = 1
	}
	return diffName(Difficulties[i])
}

// LSState is the level-select view state. Cursor ranges 0..len(Levels);
// len(Levels) is the "maze (procedural)" row. Seed is raw typed text; the
// empty string means "random" at start time.
type LSState struct {
	Levels  []string
	Cursor  int
	Seed    string
	Diff    int
	Err     string
	Preview *game.Map
}

const lsMazeRow = "maze (procedural)"

// lsLayout returns the first item row, seed row, difficulty row and first
// preview row of the level-select screen. The list holds nLevels+1 rows
// (the maze row last) and the chrome below it is packed tight, so a
// scale-1 map preview fits on a 32-row terminal.
func lsLayout(h, nLevels int) (first, seed, diff, prev int) {
	off := screenOff(h)
	first = off + 2
	seed = first + nLevels + 1
	diff = seed + 1
	prev = diff + 1
	return
}

// lsMaxRowWidth is the widest selectable row (the maze row is the longest).
func lsMaxRowWidth(names []string) int {
	maxW := len(lsMazeRow)
	for _, n := range names {
		if len(n) > maxW {
			maxW = len(n)
		}
	}
	return maxW
}

// LSRects returns the hit-test rectangles for the level rows, the three
// difficulty labels and the seed row.
func LSRects(v LSState, w, h int) (rows []Rect, diffs [3]Rect, seedRect Rect) {
	first, seed, diff, _ := lsLayout(h, len(v.Levels))
	x := (w - (lsMaxRowWidth(v.Levels) + 3)) / 2
	for i := 0; i <= len(v.Levels); i++ {
		rows = append(rows, Rect{X: x, Y: first + i, W: lsMaxRowWidth(v.Levels) + 3, H: 1})
	}
	_, labels := diffLayout(w)
	for i, l := range labels {
		diffs[i] = Rect{X: l.X - 1, Y: diff, W: len(l.S) + 2, H: 1}
	}
	body := v.Seed
	if len(body) < len("(empty = random)") {
		body = "(empty = random)"
	}
	s := " seed: " + body + "▌"
	seedRect = Rect{X: (w - len([]rune(s))) / 2, Y: seed, W: len([]rune(s)), H: 1}
	return
}

// diffLayout returns where " difficulty: [easy] [normal] [hard]" is drawn:
// the line start x and each label's start x (first letter, after the
// bracket) and text. The fixed width keeps the rectangles
// selection-independent.
func diffLayout(w int) (lineX int, labels [3]struct {
	X int
	S string
}) {
	names := [3]string{"easy", "normal", "hard"}
	line := " difficulty:"
	for _, n := range names {
		line += " [" + n + "]"
	}
	lineX = (w - len(line)) / 2
	x := lineX + len(" difficulty:")
	for i, n := range names {
		labels[i] = struct {
			X int
			S string
		}{x + 2, n}
		x += len(n) + 3
	}
	return
}

func RenderLevelSelect(v LSState, w, h int, pal Colors) *Frame {
	f := screenBox(w, h, "SELECT LEVEL", []fseg{
		{key: "enter", text: " start"},
		{key: "esc", text: " back"},
	}, true, pal)
	first, seed, diff, prev := lsLayout(h, len(v.Levels))
	guard := func(y int) bool { return y > 0 && y < h-1 }
	for i := 0; i <= len(v.Levels); i++ {
		y := first + i
		if !guard(y) {
			continue
		}
		name := lsMazeRow
		if i < len(v.Levels) {
			name = v.Levels[i]
		}
		if i == v.Cursor {
			centerPut(f, y, " ▸ "+name, pal.Bright, true)
		} else {
			centerPut(f, y, "   "+name, pal.Dim, false)
		}
	}
	if guard(seed) {
		if v.Cursor == len(v.Levels) {
			// One centered string, caret appended, so the caret never
			// drifts from the placeholder on odd-width frames.
			body := v.Seed
			bodyFG := pal.Bright
			if body == "" {
				body, bodyFG = "(empty = random)", pal.Path
			}
			full := " seed: " + body
			x0 := (w - len([]rune(full))) / 2
			putString(f, x0, seed, " seed: ", pal.Path, 0, false)
			putString(f, x0+len(" seed: "), seed, body, bodyFG, 0, false)
			f.Set(x0+len([]rune(full)), seed, Cell{R: '▌', FG: pal.Bright})
		} else {
			centerPut(f, seed, " seed: (empty = random)", pal.Path, false)
		}
	}
	if guard(diff) {
		lineX, labels := diffLayout(w)
		putString(f, lineX, diff, " difficulty:", pal.Dim, 0, false)
		for i, l := range labels {
			f.Set(l.X-1, diff, Cell{R: '[', FG: pal.Path})
			fg, bold := pal.Dim, false
			if i == v.Diff {
				fg, bold = pal.Bright, true
			}
			putString(f, l.X, diff, l.S, fg, 0, bold)
			f.Set(l.X+len([]rune(l.S)), diff, Cell{R: ']', FG: pal.Path})
		}
	}
	if v.Err != "" && guard(diff+1) {
		centerPut(f, diff+1, v.Err, 196, false)
	}
	// Map preview in a content-sized [ PREVIEW ] sub-box, only if the
	// frame has room (H*s+2 rows, W*s+2 cols) below the chrome. An error
	// message claims the first preview row, so the preview shifts down.
	if v.Preview != nil {
		p := prev
		if v.Err != "" {
			p++
		}
		availH := (h - 2) - p + 1
		scale := 0
		if availH >= v.Preview.H*2+2 && w >= v.Preview.W*2+4 {
			scale = 2
		} else if availH >= v.Preview.H+2 && w >= v.Preview.W+4 {
			scale = 1
		}
		if scale > 0 {
			bw, bh := v.Preview.W*scale+2, v.Preview.H*scale+2
			bx := (w - bw) / 2
			drawSubBox(f, bx, p, bw, bh, "PREVIEW", pal)
			l := Layout{Ox: bx + 1, Oy: p + 1, Scale: scale, W: w, H: h}
			drawMapPreview(f, v.Preview, pal, l)
		} else if guard(p) {
			centerPut(f, p, " (preview needs more room) ", 238, false)
		}
	}
	return f
}
