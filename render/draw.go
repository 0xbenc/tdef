package render

import (
	"fmt"
	"math"
	"strings"

	"tdef/game"
)

// NoSelection is "no tower selected". UI.Selected must be initialized to it
// in every constructor: the Go zero value (0) is a real tower ID, so a
// fresh UI would hide the cursor and highlight the first tower built.
const NoSelection = -1

type UI struct {
	Cursor    game.Vec
	Placing   game.TowerKind
	PlacingOn bool
	Selected  int
	Speed     int
	Paused    bool
	Help      bool
	Message   string

	Level string // level name for the header, e.g. "winding" or "maze1234"

	BestScore int
	NewBest   bool
}

const (
	FrameW = 62 // minimum frame width at scale 1

	// Chrome budget of the terminal-sized in-game frame: rows 0-1 on top
	// (header segments embedded in the top border, message row) and rows
	// th-4..th-1 on the bottom (two slot rows, hint row, bottom border).
	// The playfield region is rows ChromeTop..th-ChromeBot-1, cols 1..tw-2.
	ChromeTop = 2
	ChromeBot = 4
)

type Colors struct {
	Wall, Path, Grass int
	Gold, Dim, Bright int
	Tower             [game.TowerCount]int
	Enemy             [game.EnemyCount]int
	Beam              [game.TowerCount]int
}

func Palette() Colors {
	return Colors{
		Wall: 235, Path: 240, Grass: 234,
		Gold: 220, Dim: 245, Bright: 255,
		Tower: [game.TowerCount]int{46, 203, 51, 171, 220, 130, 226},
		Enemy: [game.EnemyCount]int{213, 214, 180, 204, 171, 199, 147, 75},
		Beam:  [game.TowerCount]int{255, 203, 51, 171, 226, 130, 226},
	}
}

// Layout maps map-space coordinates to frame-space at a given integer scale.
type Layout struct {
	Ox, Oy int
	Scale  int
	W, H   int
}

// ComputeScale picks the largest integer scale (1-4) whose map fits the
// playfield region of a tw×th terminal (rows ChromeTop..th-ChromeBot-1,
// cols 1..tw-2). Recomputed every frame, so resizes reflow live.
func ComputeScale(mW, mH, tw, th int) int {
	if tw <= 0 || th <= 0 {
		return 1
	}
	s := (tw - 2) / mW
	if v := (th - ChromeTop - ChromeBot) / mH; v < s {
		s = v
	}
	if s < 1 {
		s = 1
	}
	if s > 4 {
		s = 4
	}
	return s
}

// GameLayout is the frame layout for a tw×th terminal at the largest
// fitting stepped scale: the map centered in the playfield region (rows
// ChromeTop..th-ChromeBot-1, cols 1..tw-2). W,H are the frame (= terminal)
// size; the renderer and the mouse mapping both use this, so clicks can
// never drift from the pixels.
func GameLayout(mW, mH, tw, th int) Layout {
	s := ComputeScale(mW, mH, tw, th)
	ox := 1 + (tw-2-mW*s)/2
	if ox < 1 {
		ox = 1
	}
	oy := ChromeTop + (th-ChromeTop-ChromeBot-mH*s)/2
	if oy < ChromeTop {
		oy = ChromeTop
	}
	return Layout{Ox: ox, Oy: oy, Scale: s, W: tw, H: th}
}

// MinFrame returns the smallest terminal (scale 1) needed for a map of the
// given size. A terminal smaller than this cannot show the playfield.
func MinFrame(mW, mH int) (w, h int) {
	w = FrameW
	if mW+2 > w {
		w = mW + 2
	}
	h = ChromeTop + mH + ChromeBot
	return
}

// RenderTooSmall builds a frame that fills tw×th with a centered "enlarge
// your terminal" notice, used when the window is smaller than MinFrame.
func RenderTooSmall(tw, th, needW, needH int) *Frame {
	if tw <= 0 || th <= 0 {
		tw, th = 80, 24
	}
	f := &Frame{W: tw, H: th, C: make([]Cell, tw*th)}
	m1 := " tdef: terminal too small"
	m2 := fmt.Sprintf(" needs at least %dx%d — enlarge the window", needW, needH)
	// Kept short: the notice renders in a terminal NARROWER than the
	// frame, so long lines get clipped at the edges.
	m3 := "(paused — resize to resume)"
	for i := range f.C {
		f.C[i] = Cell{R: ' '}
	}
	putString(f, (tw-len(m1))/2, th/2-1, m1, 220, 0, true)
	putString(f, (tw-len(m2))/2, th/2+1, m2, 245, 0, false)
	putString(f, (tw-len(m3))/2, th/2+3, m3, 240, 0, false)
	return f
}

func (l Layout) X(mx int) int { return l.Ox + mx*l.Scale }
func (l Layout) Y(my int) int { return l.Oy + my*l.Scale }

// FX/FY map float (map) positions to frame space, with sub-block precision.
func (l Layout) FX(px float64) int { return l.Ox + int(px*float64(l.Scale)) }
func (l Layout) FY(py float64) int { return l.Oy + int(py*float64(l.Scale)) }

func (l Layout) MenuTop() int { return l.H - 4 }

func (l Layout) center(mx, my int) (int, int) {
	return l.X(mx) + l.Scale/2, l.Y(my) + l.Scale/2
}

// block fills a map cell with a Scale×Scale frame block.
func (l Layout) block(f *Frame, mx, my int, c Cell) {
	for dy := 0; dy < l.Scale; dy++ {
		for dx := 0; dx < l.Scale; dx++ {
			f.Set(l.X(mx)+dx, l.Y(my)+dy, c)
		}
	}
}

func drawRange(f *Frame, g *game.State, ui *UI, pal Colors, l Layout) {
	var center *game.Pos
	var rng float64
	if ui.PlacingOn {
		c := ui.Cursor.Center()
		rng = game.TowerSpecs[ui.Placing].Range[0]
		center = &c
	} else if ui.Selected >= 0 {
		if t := g.Tower(ui.Selected); t != nil {
			c := t.Pos()
			rng = t.Range()
			center = &c
		}
	}
	if center == nil {
		return
	}
	for y := 0; y < g.Map.H; y++ {
		for x := 0; x < g.Map.W; x++ {
			if g.Map.At(game.Vec{X: x, Y: y}) == game.CellWall {
				continue
			}
			p := game.Pos{X: float64(x) + 0.5, Y: float64(y) + 0.5}
			d := p.Dist(*center)
			if d <= rng && d > rng-0.6 {
				cx, cy := l.center(x, y)
				f.Put(cx, cy, '·', pal.Dim, 0)
			}
		}
	}
}

func drawLine(f *Frame, ax, ay, bx, by, color int) {
	dx, dy := float64(bx-ax), float64(by-ay)
	steps := int(math.Sqrt(dx*dx+dy*dy)*2) + 1
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := ax + int(dx*t)
		y := ay + int(dy*t)
		if x < 0 || y < 0 || x >= f.W || y >= f.H {
			continue
		}
		if f.C[y*f.W+x].R != ' ' && f.C[y*f.W+x].R != 0 {
			continue
		}
		f.Put(x, y, '·', color, 0)
	}
}

// drawHPBar renders a 3-segment health bar centered above (x,y).
func drawHPBar(f *Frame, x, y int, hp float64) {
	if y < 0 {
		return
	}
	filled := int(math.Round(hp * 3))
	if filled < 0 {
		filled = 0
	}
	if filled > 3 {
		filled = 3
	}
	col := 46
	if hp < 0.34 {
		col = 196
	} else if hp < 0.67 {
		col = 214
	}
	for i := -1; i <= 1; i++ {
		if i+1 < filled {
			f.Put(x+i, y, '■', col, 0)
		} else {
			f.Put(x+i, y, '■', 238, 0)
		}
	}
}

// drawRing renders a shrinking impact ring for a splash Fx.
func drawRing(f *Frame, l Layout, fx *game.Fx) {
	frac := fx.TTL / fx.Max
	r := fx.Ring * frac
	if r < 0.4 {
		return
	}
	cx, cy := fx.Pos.X, fx.Pos.Y
	minx, maxx := int(cx-r)-1, int(cx+r)+1
	miny, maxy := int(cy-r)-1, int(cy+r)+1
	for my := miny; my <= maxy; my++ {
		for mx := minx; mx <= maxx; mx++ {
			p := game.Pos{X: float64(mx) + 0.5, Y: float64(my) + 0.5}
			if math.Abs(p.Dist(game.Pos{X: cx, Y: cy})-r) < 0.5 {
				sx, sy := l.X(mx)+l.Scale/2, l.Y(my)+l.Scale/2
				f.Put(sx, sy, '·', fx.Color, 0)
			}
		}
	}
}

type MenuSlot struct {
	Kind game.TowerKind
	X    int
	Y    int
	W    int // rendered label width ("k Name cost")
}

// towerInfo is the menu line for the selected tower. It must fit within
// FrameW (62) at 1x scale, so keep it short (TestTowerInfoFitsFrame).
func towerInfo(t *game.Tower, upCost, refund int) string {
	if t.Level >= 3 {
		return fmt.Sprintf(" ▸ %s Lv%d (max)  %s [t]  sell +%d", t.Spec().Name, t.Level, t.TargetMode.Name(), refund)
	}
	return fmt.Sprintf(" ▸ %s Lv%d  %.0fd %.1fr  %s [t]  up %d  sell +%d", t.Spec().Name, t.Level, t.Dmg(), t.Range(), t.TargetMode.Name(), upCost, refund)
}

func diffName(d game.Difficulty) string {
	switch d {
	case game.Easy:
		return "easy"
	case game.Hard:
		return "hard"
	}
	return "normal"
}

// drawGameOver renders the end-of-game box: a 46×12 rounded box centered in
// the frame, the state title embedded in its top border (btop grammar),
// two-column stats, the best-score line and the restart/quit hint.
func drawGameOver(f *Frame, g *game.State, ui *UI, pal Colors) {
	const bg, border = 236, 240
	title, tc := "VICTORY", 48
	if g.Status != game.StatusVictory {
		title, tc = "DEFEAT", 167
	}
	bw, bh := 46, 12
	bx, by := (f.W-bw)/2, (f.H-bh)/2
	for y := by; y < by+bh; y++ {
		for x := bx; x < bx+bw; x++ {
			f.Set(x, y, Cell{R: ' ', BG: bg})
		}
	}
	drawRoundedBox(f, bx, by, bw, bh, border)
	embedSegment(f, by, bx+(bw-len([]rune(title))-2)/2, title, '┐', '┌', border, tc, true)
	put := func(y, col int, k, v string) {
		putString(f, bx+col, y, fmt.Sprintf("%-7s", k), 254, bg, true)
		putString(f, bx+col+7, y, v, 251, bg, false)
	}
	put(by+2, 2, "wave", fmt.Sprintf("%d/%d", g.Wave, game.MaxWaves))
	put(by+2, 24, "score", fmt.Sprintf("%d", g.Score))
	put(by+3, 2, "kills", fmt.Sprintf("%d", g.TotalKills))
	put(by+3, 24, "combo", fmt.Sprintf("x%d", g.MaxCombo))
	put(by+4, 2, "leaks", fmt.Sprintf("%d", g.TotalLeaks))
	put(by+4, 24, "time", formatTime(g.Time))
	put(by+5, 2, "towers", fmt.Sprintf("%d", len(g.Towers)))
	put(by+5, 24, "best", fmt.Sprintf("%d", ui.BestScore))
	bestLine := fmt.Sprintf("best %d", ui.BestScore)
	bold := false
	if ui.NewBest {
		bestLine = fmt.Sprintf("★ NEW BEST %d ★", ui.BestScore)
		bold = true
	}
	putString(f, bx+(bw-len([]rune(bestLine)))/2, by+7, bestLine, 220, bg, bold)
	x := bx + (bw-len("r restart | q quit"))/2
	for _, part := range []struct {
		s  string
		fg int
	}{
		{"r", 167}, {" restart | ", 251}, {"q", 167}, {" quit", 251},
	} {
		for _, ch := range part.s {
			f.Set(x, by+9, Cell{R: ch, FG: part.fg, BG: bg, Bold: part.fg == 167})
			x++
		}
	}
}

func formatTime(t float64) string {
	s := int(t)
	return fmt.Sprintf("%dm%02ds", s/60, s%60)
}

func putString(f *Frame, x, y int, s string, fg, bg int, bold bool) {
	for _, r := range s {
		if x >= f.W {
			break
		}
		f.Set(x, y, Cell{R: r, FG: fg, BG: bg, Bold: bold})
		x++
	}
}

// ---------------------------------------------------------------------------
// Terminal-sized in-game frame (btop-style chrome).
//
// Render draws a frame that is exactly tw×th — the full terminal —
// wrapped in one rounded box: a header of segments embedded in the top
// border (row 0), a message/telegraph row (row 1), the map centered in the
// playfield region, and the tower menu across the bottom (rows th-4..th-1,
// the selected-tower info embedded in the bottom border). GameLayout is the
// single source of truth for map placement, shared with the mouse mapping.

// TowerSlots distributes the seven tower menu slots across a tw×th frame:
// four on row th-4, three on row th-3. cell = (tw-2)/n and X = 1+i·cell, so
// the slots span the interior; the longest label ("4 Sniper 150", 12 cols)
// fits every cell at the minimum width 62. The renderer and the click
// handler share this function.
func TowerSlots(tw, th int) []MenuSlot {
	kinds := []game.TowerKind{
		game.TowerGunner, game.TowerCannon, game.TowerFrost, game.TowerSniper,
		game.TowerTesla, game.TowerMortar, game.TowerFlak,
	}
	var out []MenuSlot
	for r, row := range [][]game.TowerKind{kinds[:4], kinds[4:]} {
		y := th - ChromeBot + r
		cell := (tw - 2) / len(row)
		for i, k := range row {
			spec := game.TowerSpecs[k]
			w := len(fmt.Sprintf("%d %s %d", k+1, spec.Name, spec.Cost[0]))
			out = append(out, MenuSlot{Kind: k, X: 1 + i*cell, Y: y, W: w})
		}
	}
	return out
}

// fitMsg truncates s to at most max runes, preferring a word boundary: the
// longest prefix ending in a space that fits max-1, plus an ellipsis. A
// single word longer than max-1 is hard-cut.
func fitMsg(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	for i := max - 1; i >= 1; i-- {
		if r[i-1] == ' ' {
			return string(r[:i-1]) + "…"
		}
	}
	return string(r[:max-1]) + "…"
}

type headerRun struct {
	text string
	fg   int
	bold bool
}

type headerSeg struct {
	runs []headerRun
	drop int // drop priority: lower is elided first, -1 keeps the segment
}

func (s headerSeg) len() int {
	n := 0
	for _, r := range s.runs {
		n += len([]rune(r.text))
	}
	return n
}

// headerSegments builds the row-0 segments in canonical order:
// tdef | level·diff | wave | ⏸ | ⛁ g·♥ l·★ s·xN. The wave and stats
// segments carry the decisions, so they are never dropped; under pressure
// the level·diff, then wordmark, then pause mark are elided, in that order.
func headerSegments(g *game.State, ui *UI, pal Colors) []headerSeg {
	segs := []headerSeg{
		{runs: []headerRun{{"tdef", pal.Bright, true}}, drop: 1},
	}
	if ui.Level != "" {
		segs = append(segs, headerSeg{
			runs: []headerRun{{ui.Level + " · " + diffName(g.Diff), pal.Path, false}},
			drop: 0,
		})
	}
	var wave []headerRun
	switch {
	case g.Status == game.StatusRunning && g.WaveActive:
		wave = []headerRun{{fmt.Sprintf("wave %d/%d", g.Wave, game.MaxWaves), pal.Bright, true}}
		if g.Combo >= 5 {
			wave = append(wave, headerRun{fmt.Sprintf(" ⚡%d", g.Combo), pal.Gold, true})
		}
	case g.Status == game.StatusRunning: // inter-wave break
		nw := g.Wave + 1
		in := int(g.NextWaveAt-g.Time) + 1
		if in < 0 {
			in = 0
		}
		wave = []headerRun{
			{fmt.Sprintf("next %d in %ds (", nw, in), pal.Bright, true},
			{"n", 167, true},
			{fmt.Sprintf(" +%dg)", game.EarlyBonus(g.Wave)), pal.Bright, true},
		}
	default: // game over: show the final wave
		wave = []headerRun{{fmt.Sprintf("wave %d/%d", g.Wave, game.MaxWaves), pal.Bright, true}}
	}
	segs = append(segs, headerSeg{wave, -1})
	if ui.Paused {
		segs = append(segs, headerSeg{runs: []headerRun{{"⏸", 167, true}}, drop: 2})
	}
	stats := []headerRun{
		{fmt.Sprintf("⛁ %d", g.Gold), pal.Gold, false},
		{" · ", pal.Path, false},
		{fmt.Sprintf("♥ %d", g.Lives), 167, false},
		{" · ", pal.Path, false},
		{fmt.Sprintf("★ %d", g.Score), pal.Bright, false},
		{" · ", pal.Path, false},
	}
	if ui.Speed > 1 {
		stats = append(stats, headerRun{fmt.Sprintf("x%d", ui.Speed), 167, true})
	} else {
		stats = append(stats, headerRun{"x1", pal.Bright, false})
	}
	segs = append(segs, headerSeg{stats, -1})
	return segs
}

// drawHeader renders the btop-style top chrome: row 0 is the top border
// with embedded segments (╭─┐seg┌─┐seg┌──╮), row 1 is the message row.
func drawHeader(f *Frame, g *game.State, ui *UI, pal Colors) {
	const border = 240
	seg := func(x int, s headerSeg) int {
		f.Set(x, 0, Cell{R: '┐', FG: border})
		x++
		for _, r := range s.runs {
			for _, ch := range r.text {
				f.Set(x, 0, Cell{R: ch, FG: r.fg, Bold: r.bold})
				x++
			}
		}
		f.Set(x, 0, Cell{R: '┌', FG: border})
		return x + 1
	}
	segs := headerSegments(g, ui, pal)
	interior := f.W - 2
	total := func() int {
		n := 1 // the dash after the corner
		for _, s := range segs {
			n += s.len() + 2
		}
		return n + len(segs) - 1 + 1 // separators + one trailing dash
	}
	for total() > interior && len(segs) > 2 {
		idx, best := -1, 3
		for i, s := range segs {
			if s.drop >= 0 && s.drop < best {
				best, idx = s.drop, i
			}
		}
		if idx < 0 {
			break
		}
		segs = append(segs[:idx], segs[idx+1:]...)
	}
	if t := total() - interior; t > 0 {
		// Safety net for absurd stats: trim the last (never-dropped)
		// segment instead of painting over the corner.
		l := &segs[len(segs)-1]
		rr := &l.runs[len(l.runs)-1]
		r := []rune(rr.text)
		if t >= len(r) {
			rr.text = ""
		} else {
			rr.text = string(r[:len(r)-t])
		}
	}
	x := 2 // btop leaves one dash between the corner and the first title
	for i, s := range segs {
		if i > 0 {
			f.Set(x, 0, Cell{R: '─', FG: border})
			x++
		}
		x = seg(x, s)
	}

	// Row 1: message / telegraph / break preview.
	msg := ui.Message
	running := g.Status == game.StatusRunning
	active := running && g.WaveActive
	inBreak := running && !g.WaveActive && g.Wave < game.MaxWaves
	switch {
	case inBreak && msg != "":
		// A message held over the break (the wave telegraph, or any
		// transient) takes the whole row as a btop-style bar.
		for xx := 0; xx < f.W; xx++ {
			f.Set(xx, 1, Cell{R: ' ', BG: 167})
		}
		f.Set(0, 1, Cell{R: '├', FG: pal.Bright, BG: 167})
		f.Set(f.W-1, 1, Cell{R: '┤', FG: pal.Bright, BG: 167})
		text := fitMsg(msg, f.W-4)
		putString(f, (f.W-len([]rune(text)))/2, 1, text, pal.Bright, 167, true)
	case active && msg != "":
		fg, bold := pal.Bright, false
		switch {
		case strings.HasPrefix(msg, "leak!"):
			fg, bold = 167, true
		case strings.Contains(msg, "-> Lv"):
			fg, bold = 48, true
		case strings.Contains(msg, "cleared +"):
			fg, bold = pal.Gold, true
		}
		putString(f, 2, 1, fitMsg(msg, f.W-4), fg, 0, bold)
	case inBreak:
		nw := g.Wave + 1
		putString(f, 2, 1, fmt.Sprintf("→%d: %s", nw, game.WavePreview(nw)), pal.Dim, 0, false)
	}
}

// drawMenu renders the bottom chrome: the tower slots (rows th-4/th-3),
// the hint line (row th-2) and the bottom border (row th-1), which carries
// the selected-tower info embedded btop-style (╰──┘info└──╯).
func drawMenu(f *Frame, g *game.State, ui *UI, pal Colors) {
	const border = 240
	for _, slot := range TowerSlots(f.W, f.H) {
		k := slot.Kind
		spec := game.TowerSpecs[k]
		cost := fmt.Sprintf("%d", spec.Cost[0])
		// Label is "k Name cost": digit, space, name, space, cost.
		if ui.PlacingOn && ui.Placing == k {
			f.Set(slot.X, slot.Y, Cell{R: rune('0' + k + 1), FG: pal.Bright, BG: 95, Bold: true})
			f.Set(slot.X+1, slot.Y, Cell{R: ' ', BG: 95})
			putString(f, slot.X+2, slot.Y, spec.Name, pal.Bright, 95, true)
			f.Set(slot.X+2+len(spec.Name), slot.Y, Cell{R: ' ', BG: 95})
			putString(f, slot.X+2+len(spec.Name)+1, slot.Y, cost, pal.Gold, 95, true)
			continue
		}
		afford := g.Gold >= spec.Cost[0]
		nameFG, costFG, kFG, kBold := pal.Tower[k], pal.Gold, 167, true
		if !afford {
			nameFG, costFG, kFG, kBold = pal.Dim, pal.Dim, pal.Dim, false
		}
		f.Set(slot.X, slot.Y, Cell{R: rune('0' + k + 1), FG: kFG, Bold: kBold})
		putString(f, slot.X+2, slot.Y, spec.Name, nameFG, 0, false)
		putString(f, slot.X+2+len(spec.Name)+1, slot.Y, cost, costFG, 0, false)
	}

	hint := func(y int, pairs [][2]string) {
		x := 2
		for _, p := range pairs {
			putString(f, x, y, p[0], 167, 0, true)
			x += len([]rune(p[0]))
			putString(f, x, y, p[1], pal.Path, 0, false)
			x += len([]rune(p[1]))
		}
	}
	y := f.H - 2
	if ui.Help {
		hint(y, [][2]string{
			{"↑↓/wasd", " "}, {"⏎/1-7", " "}, {"u", " up "}, {"x", " sell "},
			{"t", " target "}, {"n", " wave "}, {"p", " pause "}, {"f", " speed"},
		})
	} else {
		hint(y, [][2]string{
			{"↑↓", "|move "}, {"⏎", "|place "}, {"u", "|up "}, {"x", "|sell "},
			{"t", "|target "}, {"n", "|wave "}, {"p", "|pause "}, {"q", "|quit"},
		})
	}

	if !ui.Help && ui.Selected >= 0 && g.Status == game.StatusRunning {
		if t := g.Tower(ui.Selected); t != nil {
			info := strings.TrimPrefix(towerInfo(t, g.UpgradeCost(t), int(float64(t.Invested)*game.SellRefund)), " ")
			if n := len([]rune(info)); n > f.W-6 {
				info = fitMsg(info, f.W-8)
			}
			r := []rune(info)
			extra := (f.W - 2) - len(r) - 2
			x := 1
			for i := 0; i < extra/2; i++ {
				f.Set(x, f.H-1, Cell{R: '─', FG: border})
				x++
			}
			f.Set(x, f.H-1, Cell{R: '┘', FG: border})
			x++
			for i, ch := range r {
				fg, bold := 251, false
				if i == 0 && ch == '▸' {
					fg, bold = 254, true
				}
				if i+3 <= len(r) && string(r[i:i+3]) == "[t]" {
					fg, bold = 167, true
				}
				f.Set(x, f.H-1, Cell{R: ch, FG: fg, Bold: bold})
				x++
			}
			f.Set(x, f.H-1, Cell{R: '└', FG: border})
			x++
			for i := 0; i < extra-extra/2; i++ {
				f.Set(x, f.H-1, Cell{R: '─', FG: border})
				x++
			}
		}
	}
}

// Render draws the full terminal-sized in-game frame (see the section
// comment above).
func Render(g *game.State, ui *UI, pal Colors, tw, th int) *Frame {
	l := GameLayout(g.Map.W, g.Map.H, tw, th)
	f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
	drawRoundedBox(f, 0, 0, l.W, l.H, pal.Path)
	drawHeader(f, g, ui, pal)
	drawMapPreview(f, g.Map, pal, l)
	drawRange(f, g, ui, pal, l)
	for _, t := range g.Towers {
		x, y := l.center(t.Cell.X, t.Cell.Y)
		c := pal.Tower[t.Kind]
		f.Put(x, y, t.Spec().Short, c, 0)
		if ui.Selected == t.ID {
			f.Set(x, y, Cell{R: t.Spec().Short, FG: pal.Bright, BG: c, Bold: true})
		}
	}
	// Level pips in a second pass, left of the tower. Pips skip cells that
	// already hold an entity glyph (tower, enemy, spawn/exit), so adjacent
	// upgraded towers can't erase each other — previously the pip pass ran
	// interleaved with the glyph pass and clobbered the left neighbor.
	for _, t := range g.Towers {
		x, y := l.center(t.Cell.X, t.Cell.Y)
		c := pal.Tower[t.Kind]
		for i := 1; i < t.Level; i++ {
			if px := x - i; px >= 0 {
				if r := f.C[y*f.W+px].R; r == 0 || r == ' ' || r == '·' {
					f.Set(px, y, Cell{R: '▪', FG: pal.Bright, BG: c, Bold: true})
				}
			}
		}
	}
	for _, p := range g.Projectiles {
		f.Put(l.FX(p.Pos.X), l.FY(p.Pos.Y), '+', pal.Beam[p.Kind], 0)
	}
	for _, bm := range g.Beams {
		for i := 1; i < len(bm.From); i++ {
			drawLine(f, l.FX(bm.From[i-1].X), l.FY(bm.From[i-1].Y), l.FX(bm.From[i].X), l.FY(bm.From[i].Y), pal.Beam[bm.Kind])
		}
	}
	for _, e := range g.Enemies {
		x, y := l.FX(e.Pos.X), l.FY(e.Pos.Y)
		if x < 0 || y < 0 || x >= f.W || y >= f.H {
			continue
		}
		hp := e.HP / e.MaxHP
		fg := pal.Enemy[e.Kind]
		if hp < 0.34 {
			fg = 196
		} else if hp < 0.67 {
			fg = 214
		}
		bold := e.Kind == game.EnemyBoss
		if e.HitTTL > g.Time {
			fg = 231 // hit flash: brief white
			bold = true
		}
		f.Set(x, y, Cell{R: game.EnemySpecs[e.Kind].Short, FG: fg, Bold: bold})
		// Skip the bar when it would land on the chrome above the map.
		if hp < 1 && y-1 >= l.Oy {
			drawHPBar(f, x, y-1, hp)
		}
	}
	for _, fx := range g.Fx {
		if fx.Ring > 0 {
			drawRing(f, l, fx)
			continue
		}
		x, y := l.FX(fx.Pos.X), l.FY(fx.Pos.Y)
		if x >= 0 && y >= 0 && x < f.W && y < f.H {
			f.Put(x, y, fx.R, fx.Color, 0)
		}
	}
	if g.LeakFlash > 0 {
		ex, ey := l.center(g.Map.Exit.X, g.Map.Exit.Y)
		f.Set(ex, ey, Cell{R: 'E', FG: 231, BG: 196, Bold: true})
	}
	if ui.PlacingOn {
		c := 196
		if g.CanBuild(ui.Cursor, ui.Placing) {
			c = 46
		}
		cx, cy := l.center(ui.Cursor.X, ui.Cursor.Y)
		f.Set(cx, cy, Cell{R: game.TowerSpecs[ui.Placing].Short, FG: c, Bold: true})
	} else if ui.Selected < 0 {
		x, y := l.X(ui.Cursor.X), l.Y(ui.Cursor.Y)
		if x >= 0 && y >= 0 && x < f.W && y < f.H {
			cc := f.C[y*f.W+x]
			// Ground glyphs (grass/wall space, path dot) get the cursor
			// marker; entity glyphs (towers, enemies, spawn/exit) just bold.
			if cc.R == 0 || cc.R == ' ' || cc.R == '·' {
				cc.R = '◻'
				cc.FG = pal.Dim
			} else {
				cc.Bold = true
			}
			f.Set(x, y, cc)
		}
	}
	drawMenu(f, g, ui, pal)
	if g.Status != game.StatusRunning {
		drawGameOver(f, g, ui, pal)
	}
	return f
}

// GameFrame renders the in-game frame for a tw×th terminal, or a centered
// "enlarge" notice when the terminal is smaller than MinFrame.
func GameFrame(g *game.State, ui *UI, pal Colors, tw, th int) *Frame {
	minW, minH := MinFrame(g.Map.W, g.Map.H)
	if tw > 0 && th > 0 && (tw < minW || th < minH) {
		return RenderTooSmall(tw, th, minW, minH)
	}
	return Render(g, ui, pal, tw, th)
}

// CaptureSize is the minimum terminal size that yields the given playfield
// scale: capture pins its virtual terminal here so output is
// environment-deterministic.
func CaptureSize(mW, mH, scale int) (tw, th int) {
	if scale < 1 {
		scale = 1
	}
	if scale > 4 {
		scale = 4
	}
	tw = FrameW
	if mW*scale+2 > tw {
		tw = mW*scale + 2
	}
	th = ChromeTop + mH*scale + ChromeBot
	return
}
