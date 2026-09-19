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

// screenFrame returns a w×h frame filled with spaces and a top/bottom
// border line. Menu screens are terminal-sized (unlike the playfield frame,
// which is Layout-sized), so a resize is picked up on the next frame.
func screenFrame(w, h int) *Frame {
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
	for x := 0; x < w; x++ {
		f.Set(x, 0, Cell{R: '═', FG: 235})
		f.Set(x, h-1, Cell{R: '═', FG: 235})
	}
	return f
}

// screenOff is the top offset of the 19-row content card inside a taller
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

// drawMapPreview renders map terrain plus the spawn/exit markers at layout
// l. Shared by the playfield, the intro, and the level-select preview.
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
	f := screenFrame(w, h)
	off := screenOff(h)
	// Card rows (19-row layout): 1-5 logo, 6 shadow, 7 tagline, 9-10 demo
	// strip, 12-13 roster, 15 best, 16 prompt.
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
				f.Set(lx+ci+1, row(1+4)+1, Cell{R: '░', FG: 238})
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
				f.Set(lx+ci, row(1+ri), Cell{R: bevel[ri], FG: titleColors[li], Bold: true})
			}
		}
	}
	put(7, "— terminal tower defense —", 245, false)
	drawTitleDemo(f, w, h, frame, row)
	drawTitleRoster(f, w, h, pal, row)
	best, name := titleBest(scores)
	if best > 0 {
		put(15, fmt.Sprintf("★ best %d — %s", best, name), 220, false)
	} else {
		put(15, " no scores yet ", 238, false)
	}
	// The prompt blinks once per second (30 ticks).
	s := "[enter] start    q quit"
	put(16, s, 245, false)
	if (frame/15)%2 == 0 {
		if yy := row(16); yy > 0 && yy < h-1 {
			putString(f, (w-len(s))/2, yy, "[enter] start", 255, 0, true)
		}
	}
	return f
}

// drawTitleDemo animates a miniature battle on card rows 9-10: an enemy
// walks the path from the spawn toward E, the tower G beams it while it is
// in range, and E flashes while the enemy "leaks". Two render ticks per
// step (~67ms) keeps the motion readable at 30fps.
func drawTitleDemo(f *Frame, w, h, frame int, row func(int) int) {
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
	y, ty := row(9), row(10)
	if y <= 0 || y >= h-1 {
		return
	}
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
	line(12, " towers  ",
		[]rune{'G', 'C', 'F', 'S', 'T', 'M', 'L'},
		[]int{pal.Tower[0], pal.Tower[1], pal.Tower[2], pal.Tower[3], pal.Tower[4], pal.Tower[5], pal.Tower[6]})
	line(13, " enemies  ",
		[]rune{'o', 'r', 'g', 't', 's', 'B', 'w', 'D'},
		[]int{pal.Enemy[0], pal.Enemy[1], pal.Enemy[2], pal.Enemy[3], pal.Enemy[4], pal.Enemy[5], pal.Enemy[6], pal.Enemy[7]})
}

// ---------------------------------------------------------------- menu

// MenuItems is the main menu, in order.
var MenuItems = []string{"Start", "Help", "High Scores", "Quit"}

// menuLayout returns the header row and the item rows of the main menu. The
// 4-item block (7 rows tall at 2-row stride) is centered between the header
// and the footer so the menu sits in the middle of the frame at any height.
func menuLayout(h int) (header int, items []int) {
	off := screenOff(h)
	header = off + 2
	top, bottom := header+1, h-3
	// Block height is 7 rows (4 items, 2-row stride); split the leftover
	// rows as evenly as possible above and below it.
	start := top + (bottom-top+1-7)/2
	if start < top {
		start = top
	}
	for i := 0; i < len(MenuItems); i++ {
		items = append(items, start+2*i)
	}
	return
}

// MenuRects returns the hit-test rectangle for each main-menu item.
func MenuRects(w, h int) []Rect {
	_, items := menuLayout(h)
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
	f := screenFrame(w, h)
	header, items := menuLayout(h)
	if header > 0 && header < h-1 {
		centerPut(f, header, " main menu ", 245, false)
	}
	for i, name := range MenuItems {
		if i >= len(items) {
			break
		}
		y := items[i]
		if y <= 0 || y >= h-1 {
			continue
		}
		if i == sel {
			centerPut(f, y, " ▸ "+name, 255, true)
		} else {
			centerPut(f, y, "   "+name, 245, false)
		}
	}
	centerPut(f, h-2, " ↑↓ move · enter select · q quit ", 240, false)
	return f
}

// ---------------------------------------------------------------- help

func RenderHelp(w, h int, pal Colors) *Frame {
	f := screenFrame(w, h)
	off := screenOff(h)
	if off+2 > 0 && off+2 < h-1 {
		centerPut(f, off+2, " help ", 245, false)
	}
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
	for i, r := range rows {
		y := off + 4 + i
		if y <= 0 || y >= h-1 {
			continue
		}
		line := fmt.Sprintf("  %-12s  %s", r[0], r[1])
		x0 := (w - len([]rune(line))) / 2
		putString(f, x0, y, line, 245, 0, false)
		putString(f, x0+2+12+2, y, r[1], 255, 0, false)
	}
	centerPut(f, h-2, " esc back ", 240, false)
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
	f := screenFrame(w, h)
	off := screenOff(h)
	if off+2 > 0 && off+2 < h-1 {
		centerPut(f, off+2, " high scores ", 245, false)
	}
	bodyTop := off + 4
	bodyBottom := h - 4
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
	if len(entries) == 0 {
		if bodyTop > 0 && bodyTop < h-1 {
			centerPut(f, bodyTop, " no scores yet ", 238, false)
		}
	}
	for i := 0; i < visible && top+i < len(entries); i++ {
		e := entries[top+i]
		name := e.Name
		if r := []rune(name); len(r) > 20 {
			name = string(r[:20])
		}
		line := fmt.Sprintf("  %-20s %10d", name, e.Score)
		x0 := (w - len([]rune(line))) / 2
		putString(f, x0, bodyTop+i, line, 245, 0, false)
		// Redraw the score bright, starting after the name field and its
		// separator space.
		putString(f, x0+2+20+1, bodyTop+i, fmt.Sprintf("%10d", e.Score), 255, 0, false)
	}
	centerPut(f, h-2, " ↑↓ scroll · esc back ", 240, false)
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

// lsLayout returns the header, first item row, seed row, difficulty row and
// first preview row of the level-select screen. The list holds nLevels+1
// rows (the maze row last) and the chrome below it is packed tight, so a
// scale-1 map preview fits on a 32-row terminal.
func lsLayout(h, nLevels int) (header, first, seed, diff, prev int) {
	off := screenOff(h)
	header = off + 2
	first = off + 3
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
	_, first, seed, diff, _ := lsLayout(h, len(v.Levels))
	x := (w - (lsMaxRowWidth(v.Levels) + 3)) / 2
	for i := 0; i <= len(v.Levels); i++ {
		rows = append(rows, Rect{X: x, Y: first + i, W: lsMaxRowWidth(v.Levels) + 3, H: 1})
	}
	lineX, labels := diffLayout(w)
	for i, l := range labels {
		diffs[i] = Rect{X: l.X - 1, Y: diff, W: len(l.S) + 2, H: 1}
	}
	_ = lineX
	body := v.Seed
	if len(body) < len("(empty = random)") {
		body = "(empty = random)"
	}
	s := " seed: " + body + "▌"
	seedRect = Rect{X: (w - len([]rune(s))) / 2, Y: seed, W: len([]rune(s)), H: 1}
	return
}

// diffLayout returns where " difficulty: easy normal hard " is drawn: the
// line start x and each label's start x and text.
func diffLayout(w int) (lineX int, labels [3]struct {
	X int
	S string
}) {
	names := [3]string{"easy", "normal", "hard"}
	line := " difficulty:"
	for _, n := range names {
		line += " " + n + " "
	}
	lineX = (w - len(line)) / 2
	x := lineX + len(" difficulty:")
	for i, n := range names {
		labels[i] = struct {
			X int
			S string
		}{x + 1, n}
		x += len(n) + 2
	}
	return
}

func RenderLevelSelect(v LSState, w, h int, pal Colors) *Frame {
	f := screenFrame(w, h)
	header, first, seed, diff, prev := lsLayout(h, len(v.Levels))
	guard := func(y int) bool { return y > 0 && y < h-1 }
	if guard(header) {
		centerPut(f, header, " select level ", 245, false)
	}
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
			centerPut(f, y, " ▸ "+name, 255, true)
		} else {
			centerPut(f, y, "   "+name, 245, false)
		}
	}
	if guard(seed) {
		if v.Cursor == len(v.Levels) {
			// One centered string, caret appended, so the caret never
			// drifts from the placeholder on odd-width frames.
			body := v.Seed
			bodyFG := 255
			if body == "" {
				body, bodyFG = "(empty = random)", 240
			}
			full := " seed: " + body
			x0 := (w - len([]rune(full))) / 2
			putString(f, x0, seed, " seed: ", 240, 0, false)
			putString(f, x0+len(" seed: "), seed, body, bodyFG, 0, false)
			f.Set(x0+len([]rune(full)), seed, Cell{R: '▌', FG: 255})
		} else {
			centerPut(f, seed, " seed: (empty = random)", 240, false)
		}
	}
	if guard(diff) {
		lineX, labels := diffLayout(w)
		putString(f, lineX, diff, " difficulty:", 245, 0, false)
		for i, l := range labels {
			if i == v.Diff {
				putString(f, l.X-1, diff, " "+l.S+" ", 255, 0, true)
			} else {
				putString(f, l.X-1, diff, " "+l.S+" ", 240, 0, false)
			}
		}
	}
	if v.Err != "" && guard(diff+1) {
		centerPut(f, diff+1, v.Err, 196, false)
	}
	// Map preview, only if the frame has room (13 rows for scale 1, 26 for
	// scale 2) below the chrome and one row above the footer. An error
	// message claims the first preview row, so the preview shifts down.
	if v.Preview != nil {
		p := prev
		if v.Err != "" {
			p++
		}
		availH := (h - 3) - p
		scale := 0
		if availH >= v.Preview.H*2 && w >= v.Preview.W*2+4 {
			scale = 2
		} else if availH >= v.Preview.H && w >= v.Preview.W+4 {
			scale = 1
		}
		if scale > 0 {
			l := Layout{Ox: (w - v.Preview.W*scale) / 2, Oy: p, Scale: scale, W: w, H: h}
			drawMapPreview(f, v.Preview, pal, l)
		}
	}
	centerPut(f, h-2, " enter start · esc back ", 240, false)
	return f
}
