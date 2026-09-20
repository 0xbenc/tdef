package render

import (
	"fmt"
	"math"
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

const titleTagline = "— terminal tower defense —"

// The title screen is one long, fully scripted loop — a pure function of the
// frame counter (30fps), titleCycle frames ≈ 56s:
//
//	0-1500    a three-wave battle in the demo box. Wave 1 (minions) plays at
//	         true 1× gameplay speed so the opening feels like the real game;
//	         waves 2-3 are compressed to keep the loop short. The towers hold
//	         wave 1, mostly hold wave 2 (one runner leaks), then wave 3's
//	         boss — the slowest thing on the board — breaks through the exit.
//	         Tower fire-rates match the real 1× RoT.
//	1500-1513 the exit overloads and glows
//	1513-1585 a static whiteout, then a shockwave from the exit eats the
//	         whole frame; the burn trail cools into a glowing grid
//	1585-1605 a beat of the cooling grid
//	1605-1615 a black beat with a single ignition spark
//	1615-1685 reboot: the border draws itself, the logo drops in row by row,
//	         the tagline types on, the demo returns to standby — frame
//	         titleCycle-1 is exactly frame 0, so the loop is seamless.
const (
	titleOverloadEnd = 1513
	titleBlastEnd    = 1585
	titleGridEnd     = 1605
	titleVoidEnd     = 1615
	titleCycle       = 1685 // must not be a multiple of 15 (footer blink seam)
)

// Demo box band rows: border, upper tower line, path, lower tower line,
// border.
const (
	demoTop   = 9
	demoUpper = 10
	demoPath  = 11
	demoLower = 12
	demoRows  = 5
)

type demoTower struct {
	kind   int     // index into game.TowerSpecs
	u      float64 // position along the path, 0 (spawn) .. 1 (exit)
	above  bool
	cd     int // frames between shots
	beam   int // frames the beam stays visible (instant towers)
	shell  int // frames of shell flight (splash towers)
	rangeU float64
}

// Five different towers: three below the path, two above. The cooldowns are
// the real level-1 RoT converted to frames (30/RoT), so the shooting reads
// as 1× gameplay.
var demoTowers = []demoTower{
	{kind: 0, u: 0.14, cd: 23, beam: 3, rangeU: 0.13},                // Gunner: fast, short
	{kind: 1, u: 0.36, cd: 55, shell: 9, rangeU: 0.11},               // Cannon: slow splash
	{kind: 2, u: 0.56, cd: 33, beam: 4, rangeU: 0.11},                // Frost: steady
	{kind: 3, u: 0.30, above: true, cd: 86, beam: 2, rangeU: 0.18},   // Sniper: long, rare
	{kind: 5, u: 0.74, above: true, cd: 67, shell: 14, rangeU: 0.16}, // Mortar: big splash
}

type demoEnemy struct {
	kind  int // index into game.EnemySpecs
	spawn int
	cross int // frames from spawn to exit
	die   int // death frame; <0 = leaks through the exit
}

// The battle script. Wave 1 minions cross in 445 frames — the true 1×
// scale-2 speed (6.4 cells/s over the 95-cell path) — and are held
// completely. Wave 2 runners are compressed (300 frames) and one leaks.
// Wave 3's boss is the slowest thing on the board (600 frames, slower than
// a minion) and walks through the exit at frame 1500 while its minions are
// picked off. Held enemies die at u≈0.75-0.76, inside the mortar's range —
// the last line of defence.
var demoWaves = []demoEnemy{
	{0, 30, 445, 370}, {0, 75, 445, 415}, {0, 120, 445, 460},
	{0, 165, 445, 505}, {0, 210, 445, 550}, {0, 255, 445, 595},
	{1, 480, 300, 705}, {1, 505, 300, 730}, {1, 530, 300, 755},
	{1, 555, 300, -1},
	{5, 900, 600, -1}, {0, 950, 445, 1290}, {0, 1010, 445, 1350}, {0, 1070, 445, 1410},
}

// demoX maps a 0..1 path position to a column: 2 (spawn) .. w-3 (exit).
func demoX(w int, u float64) int { return 2 + int(u*float64(w-5)) }

// RenderTitle draws the animated title screen. It is a pure function of
// (w, h, frame, scores): the same frame index always yields the same frame,
// so the animation is deterministic and unit-testable. `frame` is the 30fps
// tick counter.
func RenderTitle(w, h, frame int, scores map[string]int, pal Colors) *Frame {
	fr := frame % titleCycle
	switch {
	case fr < titleOverloadEnd:
		lit := (fr/15)%2 == 0
		f := screenBox(w, h, "TDEF", []fseg{
			{key: "[enter]", text: " start", blink: true},
			{key: "q", text: " quit"},
		}, lit, pal)
		drawTitleBattle(f, w, h, fr, scores, pal)
		if fr >= titleOverloadEnd-13 {
			drawTitleOverload(f, w, h, fr-titleOverloadEnd+13, pal)
		}
		return f
	case fr < titleBlastEnd:
		f := blankFrame(w, h)
		t := fr - titleOverloadEnd
		drawTitleStatic(f, w, h, t)
		drawTitleShockwave(f, w, h, t)
		return f
	case fr < titleGridEnd:
		f := blankFrame(w, h)
		drawTitleGrid(f, w, h, fr-titleBlastEnd)
		return f
	case fr < titleVoidEnd:
		f := blankFrame(w, h)
		drawTitleSpark(f, w, h, fr-titleGridEnd)
		return f
	default:
		f := blankFrame(w, h)
		drawTitleReboot(f, w, h, fr-titleVoidEnd, scores, pal)
		return f
	}
}

// blankFrame returns an empty w×h frame (the base for the effect segments).
func blankFrame(w, h int) *Frame {
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' '}
	}
	return f
}

// drawTitleBattle draws the title contents (logo, tagline, demo, roster,
// best line) on top of a screenBox frame at battle frame fr.
func drawTitleBattle(f *Frame, w, h, fr int, scores map[string]int, pal Colors) {
	off := screenOff(h)
	drawTitleLogo(f, w, off, -1)
	drawTitleTagline(f, w, off, len(titleTagline))
	drawTitleDemo(f, w, h, off, fr, pal)
	drawTitleRoster(f, w, h, off, pal)
	if yy := off + 17; yy > 0 && yy < h-1 {
		best, name := titleBest(scores)
		if best > 0 {
			centerPut(f, yy, fmt.Sprintf("★ best %d — %s", best, name), 220, false)
		} else {
			centerPut(f, yy, " no scores yet ", 238, false)
		}
	}
}

// drawTitleLogo draws the beveled TDEF slab (band rows 2-6) and its drop
// shadow (row 7). flashRow (0-4, or -1 for none) renders one logo row white
// — the reboot uses it for the letter drop.
func drawTitleLogo(f *Frame, w, off, flashRow int) {
	drawTitleLogoShadow(f, w, off)
	for i := 0; i < 5; i++ {
		drawTitleLogoRow(f, w, off, i, i == flashRow)
	}
}

func drawTitleLogoShadow(f *Frame, w, off int) {
	const logoW = 42 // 4 letters * 9 + 3 gaps * 2
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	for li, letters := range titleLetters {
		lx := x0 + li*11
		for ci, ch := range letters[4] {
			if ch == 'X' {
				f.Set(lx+ci+1, off+7, Cell{R: '░', FG: 238})
			}
		}
	}
}

func drawTitleLogoRow(f *Frame, w, off, row int, flash bool) {
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	bevel := [5]rune{'█', '▓', '▓', '▓', '▒'}
	for li, letters := range titleLetters {
		lx := x0 + li*11
		fg := titleColors[li]
		if flash {
			fg = 255
		}
		for ci, ch := range letters[row] {
			if ch == 'X' {
				f.Set(lx+ci, off+2+row, Cell{R: bevel[row], FG: fg, Bold: true})
			}
		}
	}
}

// drawTitleTagline types the tagline on character by character (n = number
// of characters shown).
func drawTitleTagline(f *Frame, w, off, n int) {
	if n <= 0 {
		return
	}
	r := []rune(titleTagline)
	if n > len(r) {
		n = len(r)
	}
	centerPut(f, off+8, string(r[:n]), 245, false)
}

// drawTitleDemo animates the scripted battle inside a full-width
// [ BATTLE ] sub-box on band rows 9-13: five towers (three below the path,
// two above) fire on scripted cooldowns at scripted enemy waves. Everything
// is a pure function of fr.
func drawTitleDemo(f *Frame, w, h, off, fr int, pal Colors) {
	y0 := off + demoTop
	if y0 <= 0 || y0+demoRows >= h-1 {
		return
	}
	drawSubBox(f, 1, y0, w-2, demoRows, "BATTLE", pal)
	waveText, waveFG := "WAVE 1", pal.Bright
	switch {
	case fr >= titleOverloadEnd-13: // the boss reaches the exit
		waveText, waveFG = "BREACH", 167
	case fr >= 900:
		waveText, waveFG = "WAVE 3", 167
	case fr >= 480:
		waveText, waveFG = "WAVE 2", pal.Bright
	}
	// The BATTLE label occupies x 2..9 (┐ + 6 + ┌); the wave segment
	// follows at x 11.
	embedSegment(f, y0, 11, waveText, '┐', '┌', pal.Path, waveFG, true)

	pathY := off + demoPath
	L := w - 5 // path spans columns 2..w-3
	for x := 2; x <= w-3; x++ {
		f.Set(x, pathY, Cell{R: '·', FG: 240})
	}
	// An energy packet flows toward the exit so the standby screen is never
	// static. One cell per frame — ambient, not gameplay-speed.
	pk := fr % L
	for i := -1; i <= 1; i++ {
		if px := 2 + (pk+i+L)%L; px >= 2 && px <= w-3 {
			f.Set(px, pathY, Cell{R: '·', FG: 251})
		}
	}
	f.Set(2, pathY, Cell{R: '▶', FG: 46, Bold: true})

	// Enemy states at this frame: alive (0) or in its death burst (1).
	type est struct {
		e  demoEnemy
		u  float64
		bf int // frames since death
	}
	var alive []est
	leakFlash := false
	for _, e := range demoWaves {
		end := e.die
		if e.die < 0 {
			end = e.spawn + e.cross
		}
		switch {
		case fr < e.spawn || fr >= end+6:
			continue
		case fr < end:
			alive = append(alive, est{e: e, u: float64(fr-e.spawn) / float64(e.cross)})
		default:
			if e.die < 0 {
				if fr-end < 8 {
					leakFlash = true
				}
				continue
			}
			alive = append(alive, est{
				e:  e,
				u:  float64(e.die-e.spawn) / float64(e.cross),
				bf: fr - end,
			})
		}
	}
	for _, s := range alive {
		x := demoX(w, s.u)
		if s.bf > 0 { // death burst
			switch {
			case s.bf == 1:
				f.Set(x, pathY, Cell{R: '*', FG: 255, Bold: true})
			case s.bf <= 3:
				f.Set(x, pathY, Cell{R: '*', FG: pal.Enemy[s.e.kind], Bold: true})
				f.Set(x-1, pathY, Cell{R: '·', FG: 244})
				f.Set(x+1, pathY, Cell{R: '·', FG: 244})
			default:
				f.Set(x, pathY, Cell{R: '·', FG: 244})
			}
			continue
		}
		f.Set(x, pathY, Cell{R: game.EnemySpecs[s.e.kind].Short, FG: pal.Enemy[s.e.kind], Bold: s.e.kind == 5})
	}

	// Towers: glyph, then the visual for the most recent shot (beam, shell,
	// or splash ring).
	for _, tw := range demoTowers {
		tx := demoX(w, tw.u)
		ty := off + demoLower
		if tw.above {
			ty = off + demoUpper
		}
		shot := (fr / tw.cd) * tw.cd
		age := fr - shot
		tgt, has := demoTarget(shot, tw)
		g := game.TowerSpecs[tw.kind].Short
		if has && age < 2 { // recoil flash
			f.Set(tx, ty, Cell{R: g, FG: 255, Bold: true})
		} else {
			f.Set(tx, ty, Cell{R: g, FG: pal.Tower[tw.kind], Bold: true})
		}
		if !has {
			continue
		}
		// The target keeps walking (or sits at its death point), so the
		// shot tracks its current position.
		end := tgt.die
		if tgt.die < 0 {
			end = tgt.spawn + tgt.cross
		}
		tfr := fr
		if tfr > end {
			tfr = end
		}
		exx := demoX(w, float64(tfr-tgt.spawn)/float64(tgt.cross))
		if tw.beam > 0 && age < tw.beam {
			drawDemoBeam(f, tx, ty, exx, pathY, pal.Beam[tw.kind])
		} else if tw.shell > 0 {
			if age < tw.shell {
				k := float64(age) / float64(tw.shell)
				sx := int(float64(tx) + (float64(exx)-float64(tx))*k)
				sy := ty
				if k >= 0.5 {
					sy = pathY
				}
				f.Set(sx, sy, Cell{R: '·', FG: pal.Beam[tw.kind], Bold: true})
			} else if age < tw.shell+10 {
				drawDemoRing(f, exx, pathY, 0.6+float64(age-tw.shell)*0.22, pal.Beam[tw.kind], age-tw.shell)
			}
		}
	}

	// The exit marker last, so anything reaching it sits under it.
	if leakFlash && fr%2 == 0 {
		f.Set(w-3, pathY, Cell{R: 'E', FG: 255, BG: 167, Bold: true})
	} else {
		f.Set(w-3, pathY, Cell{R: 'E', FG: 196, Bold: true})
	}
}

// demoTarget returns the enemy furthest along the path that is within tw's
// range at frame fr (the shot frame), or ok=false.
func demoTarget(fr int, tw demoTower) (demoEnemy, bool) {
	var best demoEnemy
	bu, ok := 0.0, false
	for _, e := range demoWaves {
		end := e.die
		if e.die < 0 {
			end = e.spawn + e.cross
		}
		if fr < e.spawn || fr >= end {
			continue
		}
		u := float64(fr-e.spawn) / float64(e.cross)
		if d := u - tw.u; d < -tw.rangeU || d > tw.rangeU {
			continue
		}
		if !ok || u > bu {
			best, bu, ok = e, u, true
		}
	}
	return best, ok
}

// drawDemoBeam draws a tower's shot as a line of '+' from (x0,y0) to
// (x1,y1); towers sit one row off the path, so the line spans two rows.
func drawDemoBeam(f *Frame, x0, y0, x1, y1, fg int) {
	steps := x1 - x0
	if steps < 0 {
		steps = -steps
	}
	for i := 1; i < steps; i++ {
		x := x0 + (x1-x0)*i/steps
		y := y0
		if i*2 >= steps {
			y = y1
		}
		f.Set(x, y, Cell{R: '+', FG: fg})
	}
}

// drawDemoRing draws a landed shell's splash ring: an expanding circle in
// the terminal's 1:2 aspect metric.
func drawDemoRing(f *Frame, cx, cy int, r float64, fg, age int) {
	col := fg
	if age > 4 {
		col = 51
	}
	dx := int(r/0.55) + 2
	for x := cx - dx; x <= cx+dx; x++ {
		if x < 0 || x >= f.W {
			continue
		}
		for y := cy - int(r) - 2; y <= cy+int(r)+2; y++ {
			if y < 0 || y >= f.H {
				continue
			}
			d := math.Hypot(float64(x-cx)*0.55, float64(y-cy))
			if math.Abs(d-r) < 0.6 {
				f.Set(x, y, Cell{R: '·', FG: col})
			}
		}
	}
}

// drawTitleRoster lists the tower and enemy glyphs in their colors, so the
// title doubles as a legend.
func drawTitleRoster(f *Frame, w, h, off int, pal Colors) {
	// Both lines share one grid (label column 9 wide, glyphs 4 apart,
	// anchored to the 8-glyph enemy line) so the columns line up.
	const labelW, stride, maxGlyphs = 9, 4, 8
	widest := labelW + maxGlyphs*stride
	line := func(y int, label string, glyphs []rune, colors []int) {
		yy := off + y
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
	line(14, " towers  ",
		[]rune{'G', 'C', 'F', 'S', 'T', 'M', 'L'},
		[]int{pal.Tower[0], pal.Tower[1], pal.Tower[2], pal.Tower[3], pal.Tower[4], pal.Tower[5], pal.Tower[6]})
	line(15, " enemies  ",
		[]rune{'o', 'r', 'g', 't', 's', 'B', 'w', 'D'},
		[]int{pal.Enemy[0], pal.Enemy[1], pal.Enemy[2], pal.Enemy[3], pal.Enemy[4], pal.Enemy[5], pal.Enemy[6], pal.Enemy[7]})
}

// ---------------------------------------------------------------- effects

// titleHash is a cheap deterministic per-cell/per-frame hash for effect
// flicker.
func titleHash(x, y, t int) int {
	n := x*73856093 ^ y*19349663 ^ t*83492791
	if n < 0 {
		n = -n
	}
	return n
}

// drawTitleOverload blooms the exit after the boss breaks through: a red
// heat wash that expands out from E with a bright front, swallowing the
// frame contents near the exit.
func drawTitleOverload(f *Frame, w, h, t int, pal Colors) {
	off := screenOff(h)
	ox, oy := w-3, off+demoPath
	rr := 1.0 + float64(t)*1.1
	dx := int(rr/0.55) + 2
	for x := ox - dx; x <= ox+dx; x++ {
		if x < 1 || x >= w-1 {
			continue
		}
		for y := oy - int(rr) - 2; y <= oy+int(rr)+2; y++ {
			if y < 1 || y >= h-1 {
				continue
			}
			d := math.Hypot(float64(x-ox)*0.55, float64(y-oy))
			if math.Abs(d-rr) < 0.8 {
				f.Set(x, y, Cell{R: '█', FG: 255})
			} else if d < rr {
				f.Set(x, y, Cell{R: ' ', FG: 0, BG: 124 + (12-t)*6})
			}
		}
	}
	f.Set(ox, oy, Cell{R: 'E', FG: 255, Bold: true})
}

// drawTitleStatic paints the whiteout: a field of flickering static that
// the shockwave carves through.
func drawTitleStatic(f *Frame, w, h, t int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n := titleHash(x, y, t)
			r, c := '░', 234
			switch {
			case n%37 == 0:
				r, c = '█', 255
			case n%13 == 0:
				r, c = '▒', 245
			case n%5 == 0:
				c = 238
			}
			f.Set(x, y, Cell{R: r, FG: c})
		}
	}
}

// drawTitleShockwave consumes the static with an expanding shockwave from
// the exit: a white-cyan front, a hot flickering trail with re-sparks, then
// a cool-down that leaves the glowing grid the reboot builds on.
func drawTitleShockwave(f *Frame, w, h, t int) {
	off := screenOff(h)
	ox, oy := w-3, off+demoPath
	r := math.Max(0, float64(t-2)) * 2.6
	if r <= 0 {
		return
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			d := math.Hypot(float64(x-ox)*0.55, float64(y-oy))
			if d > r {
				continue
			}
			age := r - d
			n := titleHash(x, y, t)
			switch {
			case d < 1.5 && age < 40: // the origin keeps pulsing
				f.Set(x, y, Cell{R: '*', FG: 255, Bold: true})
			case age < 1.2:
				f.Set(x, y, Cell{R: '█', FG: 255, Bold: true})
			case age < 3:
				f.Set(x, y, Cell{R: '█', FG: 51})
			case age < 6:
				f.Set(x, y, Cell{R: '▓', FG: 117})
			case age < 12:
				c := 33
				if n%3 == 0 {
					c = 39
				}
				f.Set(x, y, Cell{R: '▒', FG: c})
			default:
				if n%89 == 0 && age < 28 {
					f.Set(x, y, Cell{R: '·', FG: 231})
					continue
				}
				if x%3 == 0 && y%2 == 0 {
					f.Set(x, y, Cell{R: '·', FG: 24})
				} else {
					f.Set(x, y, Cell{R: ' '})
				}
			}
		}
	}
}

// drawTitleGrid holds a beat of the cooling grid left by the shockwave.
func drawTitleGrid(f *Frame, w, h, t int) {
	c := 24
	if t < 6 {
		c = 33
	}
	if t%3 == 0 {
		c += 4
	}
	for y := 0; y < h; y += 2 {
		for x := 0; x < w; x += 3 {
			f.Set(x, y, Cell{R: '·', FG: c})
		}
	}
}

// drawTitleSpark is the ignition spark at the center of the black beat.
func drawTitleSpark(f *Frame, w, h, t int) {
	x, y := w/2, h/2
	switch {
	case t >= 3 && t < 6:
		f.Set(x, y, Cell{R: '·', FG: 231})
	case t >= 6 && t < 9:
		f.Set(x, y, Cell{R: '*', FG: 255, Bold: true})
	}
}

// perimCell returns the (x, y, glyph) of the i-th cell of the frame
// perimeter, walked from the top-left clockwise: top row, right column,
// bottom row, left column.
func perimCell(i, w, h int) (int, int, rune) {
	if i < w {
		g := '─'
		if i == 0 {
			g = '╭'
		} else if i == w-1 {
			g = '╮'
		}
		return i, 0, g
	}
	i -= w
	if i < h-2 {
		return w - 1, i + 1, '│'
	}
	i -= h - 2
	if i < w {
		x := w - 1 - i
		g := '─'
		if x == w-1 {
			g = '╯'
		} else if x == 0 {
			g = '╰'
		}
		return x, h - 1, g
	}
	i -= w
	return 0, h - 2 - i, '│'
}

// drawTitleReboot draws the screen back into existence: the explosion's
// grid fades, the border draws itself from the top-left, the logo drops in
// row by row (each row flashing white on arrival), the tagline types on,
// and the demo returns to standby. At the last frame the screen is exactly
// the battle's frame 0, closing the loop seamlessly.
func drawTitleReboot(f *Frame, w, h, t int, scores map[string]int, pal Colors) {
	const dur = titleCycle - titleVoidEnd
	p := float64(t) / float64(dur)
	off := screenOff(h)
	if p < 0.2 { // the residual blast grid cools out first
		for y := 0; y < h; y += 2 {
			for x := 0; x < w; x += 3 {
				f.Set(x, y, Cell{R: '·', FG: 24})
			}
		}
	}
	// The border pen: top row, right column, bottom row, left column.
	P := 2*w + 2*h - 4
	pen := int(p / 0.55 * float64(P))
	if pen > P {
		pen = P
	}
	for i := 0; i < pen; i++ {
		x, y, r := perimCell(i, w, h)
		fg := 240
		if i == pen-1 && pen < P {
			fg = 51 // the pen tip glows while it draws
		}
		f.Set(x, y, Cell{R: r, FG: fg})
	}
	if p >= 0.58 {
		embedSegment(f, 0, 2, "TDEF", '┐', '┌', 240, pal.Bright, true)
	}
	pPrev := float64(t-1) / float64(dur)
	for i := 0; i < 5; i++ {
		th := 0.30 + 0.05*float64(i)
		if p < th {
			continue
		}
		drawTitleLogoRow(f, w, off, i, t > 0 && pPrev < th)
	}
	if p >= 0.55 {
		drawTitleLogoShadow(f, w, off)
	}
	drawTitleTagline(f, w, off, int((p-0.58)/0.14*float64(len(titleTagline))))
	if p >= 0.72 {
		drawTitleDemo(f, w, h, off, 0, pal)
	}
	if p >= 0.80 {
		drawTitleRoster(f, w, h, off, pal)
	}
	if p >= 0.88 {
		best, name := titleBest(scores)
		if best > 0 {
			centerPut(f, off+17, fmt.Sprintf("★ best %d — %s", best, name), 220, false)
		} else {
			centerPut(f, off+17, " no scores yet ", 238, false)
		}
	}
	if p >= 0.95 {
		drawFooter(f, []fseg{
			{key: "[enter]", text: " start"},
			{key: "q", text: " quit"},
		}, true, pal)
	}
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
