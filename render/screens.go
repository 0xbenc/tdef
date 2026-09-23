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
// drawMapPreview renders the battlefield terrain: mottled rock walls, grassy
// clearings, and the road — drawn as directional connectors so the route the
// horde marches reads at a glance. The spawn rift and the lair heart sit on the
// road's ends. Everything is scale-aware: at 1x each map cell is one terminal
// cell; at 2x+ a block carries its own texture.
func drawMapPreview(f *Frame, m *game.Map, pal Colors, l Layout, frame int) {
	isRoad := func(x, y int) bool { return m.At(game.Vec{X: x, Y: y}) == game.CellPath }
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			switch m.At(game.Vec{X: x, Y: y}) {
			case game.CellWall:
				drawWallBlock(f, pal, l, x, y)
			case game.CellPath:
				g := roadGlyph(isRoad(x, y-1), isRoad(x+1, y), isRoad(x, y+1), isRoad(x-1, y))
				drawRoadBlock(f, pal, l, x, y, g)
			case game.CellGrass:
				drawGrassBlock(f, pal, l, x, y)
			}
		}
	}
	drawSpawnRift(f, pal, l, m.Spawn, frame)
	drawLairHeart(f, pal, l, m.Exit, frame)
}

// roadGlyph picks the box-drawing connector for a road cell from which of its
// four neighbours are also road: a straight, a corner, a tee or a cross.
func roadGlyph(n, e, s, w bool) rune {
	switch {
	case n && e && s && w:
		return '┼'
	case n && e && w:
		return '┬'
	case e && s && w:
		return '┴'
	case n && s && w:
		return '├'
	case n && s && e:
		return '┤'
	case n && s:
		return '│'
	case e && w:
		return '─'
	case n && w:
		return '┘' // up + left
	case n && e:
		return '└' // up + right
	case s && w:
		return '┐' // down + left
	case s && e:
		return '┌' // down + right
	case n || s:
		return '│'
	case e || w:
		return '─'
	}
	return '·'
}

// cellHash is a stable per-cell value in [0,1) that drives ambient variation
// (rock mottling, grass tufts). Deterministic, so renders reproduce exactly.
func cellHash(x, y int) float64 {
	h := uint32(x+0x9e37)*73856093 ^ uint32(y+0x27f0)*19349663
	h ^= h >> 13
	h *= 0x5bd1e995
	h ^= h >> 15
	return float64(h%1000) / 1000.0
}

// drawWallBlock fills one wall cell. At 1x a single mottled rock cell; at 2x+ a
// block of stone with a few lighter/darker speckles so it reads as rock, not a
// flat slab.
func drawWallBlock(f *Frame, pal Colors, l Layout, x, y int) {
	base := pal.Wall
	if h := cellHash(x, y); h < 0.15 {
		base = pal.WallHi
	} else if h > 0.85 {
		base = pal.WallLo
	}
	if l.Scale == 1 {
		f.Set(l.X(x), l.Y(y), Cell{R: ' ', BG: base})
		return
	}
	for dy := 0; dy < l.Scale; dy++ {
		for dx := 0; dx < l.Scale; dx++ {
			sh := base
			if hs := cellHash(x*31+dx, y*17+dy); hs < 0.10 {
				sh = pal.WallHi
			} else if hs > 0.92 {
				sh = pal.WallLo
			}
			f.Set(l.X(x)+dx, l.Y(y)+dy, Cell{R: ' ', BG: sh})
		}
	}
}

// drawGrassBlock fills one clearing cell: a dark-green floor with sparse
// lighter tufts so it reads as grass set into the rock.
func drawGrassBlock(f *Frame, pal Colors, l Layout, x, y int) {
	if l.Scale == 1 {
		c := Cell{R: ' ', BG: pal.Grass}
		if cellHash(x, y) < 0.18 {
			c = Cell{R: '·', FG: pal.GrassTuft, BG: pal.Grass}
		}
		f.Set(l.X(x), l.Y(y), c)
		return
	}
	for dy := 0; dy < l.Scale; dy++ {
		for dx := 0; dx < l.Scale; dx++ {
			c := Cell{BG: pal.Grass}
			if cellHash(x*13+dx, y*29+dy) < 0.12 {
				c = Cell{R: '·', FG: pal.GrassTuft, BG: pal.Grass}
			}
			f.Set(l.X(x)+dx, l.Y(y)+dy, c)
		}
	}
}

// drawRoadBlock fills one road cell with the road surface and its directional
// centerline connector. At 2x+ the whole block is road, so adjacent cells merge
// into a continuous band with the connector marking the route.
func drawRoadBlock(f *Frame, pal Colors, l Layout, x, y int, g rune) {
	if l.Scale == 1 {
		f.Set(l.X(x), l.Y(y), Cell{R: g, FG: pal.RoadLine, BG: pal.RoadBG})
		return
	}
	for dy := 0; dy < l.Scale; dy++ {
		for dx := 0; dx < l.Scale; dx++ {
			f.Set(l.X(x)+dx, l.Y(y)+dy, Cell{BG: pal.RoadBG})
		}
	}
	cx, cy := l.center(x, y)
	f.Set(cx, cy, Cell{R: g, FG: pal.RoadLine, BG: pal.RoadBG})
}

// drawSpawnRift marks where the horde pours in: a bright rift glyph that
// breathes on the ambient clock.
func drawSpawnRift(f *Frame, pal Colors, l Layout, v game.Vec, frame int) {
	x, y := l.center(v.X, v.Y)
	fg := pal.Spawn
	if (frame/16)%2 == 0 {
		fg = pal.Bright
	}
	f.Set(x, y, Cell{R: '▶', FG: fg, Bold: true})
}

// drawLairHeart marks the lair the horde is marching on: a heart that beats — a
// soft double-thump — on the ambient clock.
func drawLairHeart(f *Frame, pal Colors, l Layout, v game.Vec, frame int) {
	x, y := l.center(v.X, v.Y)
	fg := pal.Exit
	switch c := frame % 48; {
	case c < 3, c >= 10 && c < 13:
		fg = 203
	}
	f.Set(x, y, Cell{R: '♥', FG: fg, Bold: true})
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

// The title screen is a fully scripted state machine, a pure function of
// (frame, boot): `frame` is the 30fps tick counter (ambient timing) and
// `boot` is frames since this visit to the title started.
//
//	boot 0-250   the cinematic (~8s): a black beat, then the slab's
//	             bounding box draws itself with a scan flicker (the light
//	             pen's curtain-up); a light pen traces the TDEF slab out of
//	             digital noise (each letter flashing white as it locks in);
//	             then each letter, left to right, fires a different weapon
//	             at a creature in a frame corner — Gunner tracer spray into
//	             a 2x2 braille blob, Cannon shell + AOE burst into a 3x5
//	             bug, Sniper charge + piercing beam through an (o_o) guy,
//	             Tesla chain lightning into a >_< guy — each dying with its
//	             own animation; the slab ignites white, the subtitle
//	             decodes, then the frame chrome and an empty battlefield
//	             fade in
//	boot 251+    the attract loop, every titleAttractCycle frames:
//
//	15s idle     standby: full UI, empty battlefield — no towers, no
//	             enemies, no THE SIEGE/WAVE text
//	then 120f    "┐THE SIEGE┌" and "┐WAVE 1┌" decode on, the five
//	             towers power up left to right
//	then 1685f   the battle script: wave 1 (minions) plays at true 1×
//	             gameplay speed; waves 2-3 are compressed. The towers hold
//	             wave 1, mostly hold wave 2 (one runner leaks), then wave 3's
//	             boss — the slowest thing on the board — breaks through the
//	             exit and the exit overloads (1500-1513); a static whiteout
//	             then a shockwave from the exit eats the whole frame
//	             (1513-1585); the burn trail cools into a glowing grid
//	             (1585-1605); a black beat with one ignition spark
//	             (1605-1615); the screen reboots (1615-1684) back to the
//	             standby state, where the 15s clock starts again.
const (
	titleBootIntroEnd   = 8   // 0-7:      the slab box self-draws with scan flicker; the pen tip lands
	titleBootPenEnd     = 73  // 8-72:     the light pen traces the TDEF slab
	titleBootWeaponsEnd = 203 // 73-202:   the letters fire: Gunner, Cannon, Sniper, Tesla
	titleBootFlashStart = 203 // 203-206:  the slab ignites white
	titleBootFlashEnd   = 207 //
	titleBootSubEnd     = 221 // 207-220:  the subtitle decodes
	titleBootLen        = 251 // 221-250:  the rest of the UI fades in
)

const (
	// The attract loop: 15s of inactivity, then the battle sequence.
	titleIdleWait     = 450 // 15s at 30fps
	titleBattleLead   = 120 // THE SIEGE/WAVE text + tower power-up before wave 1
	titleBattleLen    = titleBattleLead + titleCycle
	titleAttractCycle = titleIdleWait + titleBattleLen
)

const (
	// The battle + reset script, in internal frames 0-1684.
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

// ---------------------------------------------------------------- boot

// bootStep is one slab cell in the light pen's serpentine path: each
// letter's filled cells, row by row, direction alternating.
type bootStep struct {
	relX int // column relative to the slab's left edge
	row  int // slab row 0-4
	li   int // letter index (color)
}

var titlePenPath = buildBootPenPath()

func buildBootPenPath() []bootStep {
	var p []bootStep
	for li, letters := range titleLetters {
		for row := 0; row < 5; row++ {
			lo, hi, step := 0, 9, 1
			if row%2 == 1 {
				lo, hi, step = 8, -1, -1
			}
			for ci := lo; ci != hi; ci += step {
				if letters[row][ci] == 'X' {
					p = append(p, bootStep{relX: li*11 + ci, row: row, li: li})
				}
			}
		}
	}
	return p
}

// titleLetterDone[li] is the frame at which the pen finishes letter li.
var titleLetterDone = func() [4]int {
	var d [4]int
	n := 0
	for li, letters := range titleLetters {
		for row := range letters {
			for _, ch := range letters[row] {
				if ch == 'X' {
					d[li] = titleBootIntroEnd + (n+1)/2
					n++
				}
			}
		}
	}
	return d
}()

// bootPenIndex is the pen's path position at boot frame t: two cells per
// frame from titleBootIntroEnd.
func bootPenIndex(t int) int {
	pos := 2 * (t - titleBootIntroEnd)
	if pos > len(titlePenPath)-1 {
		pos = len(titlePenPath) - 1
	}
	return pos
}

// RenderTitle draws the animated title screen. It is a pure function of
// (w, h, frame, boot, scores): the same inputs always yield the same frame,
// so the animation is deterministic and unit-testable. `frame` is the 30fps
// tick counter and `boot` is the number of frames this title visit has been
// up (any keypress leaves the title, so idle time == boot time).
func RenderTitle(w, h, frame, boot int, scores map[string]int, pal Colors) *Frame {
	if boot < titleBootLen {
		f := blankFrame(w, h)
		drawTitleBoot(f, w, h, boot, frame, pal)
		return f
	}
	t := boot - titleBootLen
	local := t % titleAttractCycle
	if local < titleIdleWait {
		return drawTitleStandby(w, h, frame, scores, pal)
	}
	return drawTitleBattleSeq(w, h, frame, local-titleIdleWait, scores, pal)
}

// titleFooter is the title screen's footer group.
func titleFooter() []fseg {
	return []fseg{
		{key: "[enter]", text: " start", blink: true},
		{key: "q", text: " quit"},
	}
}

// drawTitleStandby is the idle title: full UI, empty battlefield — no
// towers, no enemies, no THE SIEGE/WAVE text, just the path and its
// ambient energy packet.
func drawTitleStandby(w, h, frame int, scores map[string]int, pal Colors) *Frame {
	lit := (frame/15)%2 == 0
	f := screenBox(w, h, "TDEF", titleFooter(), lit, pal)
	off := screenOff(h)
	drawTitleChrome(f, w, h, off, scores, pal)
	drawTitleEmptyBox(f, w, h, off, frame, pal)
	return f
}

// drawTitleBattleSeq runs the attract battle: the THE SIEGE/WAVE
// text decodes on, the towers power up, then the battle script plays
// (internal frame fr = local - titleBattleLead).
func drawTitleBattleSeq(w, h, frame, local int, scores map[string]int, pal Colors) *Frame {
	if local < titleBattleLead {
		f := drawTitleStandby(w, h, frame, scores, pal)
		off := screenOff(h)
		drawTitleBattleIntro(f, w, h, off, local, pal)
		return f
	}
	fr := local - titleBattleLead
	switch {
	case fr < titleOverloadEnd:
		lit := (frame/15)%2 == 0
		f := screenBox(w, h, "TDEF", titleFooter(), lit, pal)
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
		drawTitleReboot(f, w, h, frame, fr-titleVoidEnd, scores, pal)
		return f
	}
}

// drawTitleBattleIntro reveals "┐THE SIEGE┌" and "┐WAVE 1┌"
// character by character (each flashing white on arrival), then powers the
// five towers up left to right. At local == titleBattleLead the box matches
// battle frame 0 exactly.
func drawTitleBattleIntro(f *Frame, w, h, off, t int, pal Colors) {
	y0 := off + demoTop
	if y0 <= 0 || y0+demoRows >= h-1 {
		return
	}
	seg := func(x int, s string, fg, at int) {
		if t < at {
			return
		}
		for i, ch := range []rune(s) {
			ct := at + 2*i
			if t < ct {
				break
			}
			c := Cell{R: ch}
			switch {
			case ch == '┐' || ch == '┌':
				c = Cell{R: ch, FG: pal.Path}
			case t-ct < 2:
				c = Cell{R: ch, FG: 255, Bold: true}
			default:
				c = Cell{R: ch, FG: fg, Bold: true}
			}
			f.Set(x+i, y0, c)
		}
	}
	seg(2, "┐THE SIEGE┌", pal.Dim, 0)
	seg(14, "┐WAVE 1┌", pal.Bright, 18)

	const towerStart, towerDur = 32, 17
	// Left to right along the path.
	order := [5]int{0, 3, 1, 2, 4}
	for i, k := range order {
		tw := demoTowers[k]
		tt := t - towerStart - i*towerDur
		if tt < 0 {
			continue
		}
		tx := demoX(w, tw.u)
		ty := off + demoLower
		if tw.above {
			ty = off + demoUpper
		}
		g := game.TowerSpecs[tw.kind].Short
		switch {
		case tt < 8: // charge
			f.Set(tx, ty, Cell{R: '·', FG: 234 + tt})
		case tt < 10: // ignition
			f.Set(tx, ty, Cell{R: '█', FG: 255, Bold: true})
		default:
			f.Set(tx, ty, Cell{R: g, FG: pal.Tower[tw.kind], Bold: true})
		}
	}
}

// ---------------------------------------------------------------- boot (cont.)

// drawTitleBoot plays the one-shot cinematic: intro, pen, the weapon fan,
// flash, subtitle, UI fade-in.
func drawTitleBoot(f *Frame, w, h, t, frame int, pal Colors) {
	off := screenOff(h)
	switch {
	case t < titleBootIntroEnd:
		drawBootIntro(f, w, off, t)
	case t < titleBootPenEnd:
		drawBootPen(f, w, off, t)
	case t < titleBootWeaponsEnd:
		drawBootWeapons(f, w, h, off, t-titleBootPenEnd)
	case t < titleBootFlashEnd:
		drawBootFlash(f, w, h, off)
	case t < titleBootSubEnd:
		drawTitleLogo(f, w, off, -1)
		drawBootSubtitle(f, w, off, t-titleBootFlashEnd)
	default:
		drawTitleLogo(f, w, off, -1)
		drawTitleTagline(f, w, off, len(titleTagline))
		drawBootUI(f, w, h, off, t-titleBootSubEnd, frame, pal)
	}
}

// drawBootIntro is the light pen's curtain-up: a black beat (0), then the
// slab's bounding box draws itself (1-7), tip lit, over a scan flicker.
// The box ends one cell left of the pen's first slab cell, so the tip
// hands off cleanly when the trace starts.
func drawBootIntro(f *Frame, w, off, t int) {
	if t == 0 {
		return
	}
	// The box draws itself (14 cells/frame).
	pos := t * 98 / 7
	drawBootBox(f, w, off, pos, 1, 234)
	for y := off + 2; y <= off+6; y++ {
		for x := 0; x < w; x++ {
			n := titleHash(x, y, t)
			if n%13 == 0 {
				f.Set(x, y, Cell{R: '·', FG: 235})
			} else if n%17 == 0 {
				f.Set(x, y, Cell{R: '+', FG: 237})
			}
		}
	}
}

// boxCell is the i-th cell of a bw×bh box perimeter starting top-left and
// running clockwise.
func boxCell(i, bx, by int, bw, bh int) (int, int, rune) {
	switch {
	case i == 0:
		return bx, by, '╭'
	case i < bw:
		if i == bw-1 {
			return bx + bw - 1, by, '╮'
		}
		return bx + i, by, '─'
	case i < bw+(bh-2):
		j := i - bw
		return bx + bw - 1, by + 1 + j, '│'
	case i == 2*bw+(bh-2)-1:
		return bx + bw - 1, by + bh - 1, '╯'
	case i < 2*bw+(bh-2):
		j := i - bw - (bh - 2)
		return bx + bw - 1 - j, by + bh - 1, '─'
	default:
		j := i - 2*bw - (bh - 2)
		return bx, by + bh - 2 - j, '│'
	}
}

// drawBootBox draws the slab's bounding box (44×7, one cell around the
// slab). pos is how many of the 98 perimeter cells are drawn; tip lights
// the last cell while it is still being drawn.
func drawBootBox(f *Frame, w, off, pos, tip, fg int) {
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 1 {
		x0 = 1
	}
	bx, by := x0-1, off+1
	if pos > 98 {
		pos = 98
	}
	for i := 0; i < pos; i++ {
		x, y, r := boxCell(i, bx, by, 44, 7)
		c := fg
		if tip == 1 && i == pos-1 {
			c = 51
		}
		f.Set(x, y, Cell{R: r, FG: c})
	}
}

// drawBootLetter draws one slab letter, optionally shifted (dx, dy) and
// recolored (fg = 0 keeps the letter's own color).
func drawBootLetter(f *Frame, w, off, li, dx, dy, fg int, bold bool) {
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	if fg == 0 {
		fg = titleColors[li]
	}
	bevel := [5]rune{'█', '▓', '▓', '▓', '▒'}
	letters := titleLetters[li]
	for row := 0; row < 5; row++ {
		for ci := 0; ci < 9; ci++ {
			if letters[row][ci] != 'X' {
				continue
			}
			f.Set(x0+li*11+ci+dx, off+2+row+dy, Cell{R: bevel[row], FG: fg, Bold: bold})
		}
	}
}

// drawBootSlab is the complete, settled slab (all four letters in color).
func drawBootSlab(f *Frame, w, off int) {
	for li := 0; li < 4; li++ {
		drawBootLetter(f, w, off, li, 0, 0, 0, false)
	}
}

// drawBootPen traces the TDEF slab with a light pen: the tip is white, the
// recent trail cools white -> cyan -> the letter color, and untraced cells
// flicker as faint digital noise. The slab's box stays drawn around it.
func drawBootPen(f *Frame, w, off, t int) {
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	drawBootBox(f, w, off, 98, 0, 234)
	bevel := [5]rune{'█', '▓', '▓', '▓', '▒'}
	pos := bootPenIndex(t)
	for i, s := range titlePenPath {
		x, y := x0+s.relX, off+2+s.row
		if i <= pos {
			fg, bold := titleColors[s.li], false
			switch d := pos - i; {
			case d <= 2:
				fg, bold = 255, true
			case d <= 8:
				fg = 51
			}
			f.Set(x, y, Cell{R: bevel[s.row], FG: fg, Bold: bold})
		} else {
			n := titleHash(x, y, t/4)
			r := '·'
			if n%7 == 0 {
				r = '+'
			}
			if n%11 == 0 {
				r = '░'
			}
			f.Set(x, y, Cell{R: r, FG: 236 + n%3})
		}
	}
	// Each letter flashes white for three frames as the pen completes it.
	for li := 0; li < 4; li++ {
		if t >= titleLetterDone[li] && t < titleLetterDone[li]+3 {
			for _, s := range titlePenPath {
				if s.li == li {
					f.Set(x0+s.relX, off+2+s.row, Cell{R: bevel[s.row], FG: 255, Bold: true})
				}
			}
		}
	}
	// The pen tip.
	s := titlePenPath[pos]
	f.Set(x0+s.relX, off+2+s.row, Cell{R: '█', FG: 255, Bold: true})
}

// ---------------------------------------------------------------- weapon fan

// shotPos is the cell at fraction u (0..1) of the way from muzzle to target.
func shotPos(mx, my, tx, ty int, u float64) (int, int) {
	return mx + int(math.Round(float64(tx-mx)*u)), my + int(math.Round(float64(ty-my)*u))
}

// drawBootReticle draws the lock-on reticle at (x, y). shatter >= 0 means
// the reticle is breaking into six radial fragments (frames since the hit).
func drawBootReticle(f *Frame, x, y, shatter int) {
	if shatter >= 0 {
		c := 244
		if shatter == 1 {
			c = 236
		}
		if shatter == 2 {
			c = 234
		}
		for i := 0; i < 6; i++ {
			a := float64(i) * math.Pi / 3
			for k := 1; k <= shatter+1; k++ {
				fx := x + int(math.Round(math.Cos(a)*float64(k)))
				fy := y + int(math.Round(math.Sin(a)*float64(k)))
				f.Set(fx, fy, Cell{R: '·', FG: c})
			}
		}
		return
	}
	f.Set(x, y, Cell{R: '+', FG: 244})
	f.Set(x-1, y, Cell{R: '·', FG: 236})
	f.Set(x+1, y, Cell{R: '·', FG: 236})
	f.Set(x, y-1, Cell{R: '·', FG: 236})
	f.Set(x, y+1, Cell{R: '·', FG: 236})
}

// drawBootBracket draws a pulsing lock-on bracket one cell around the
// cols×rows target whose top-left corner is (x0, y0).
func drawBootBracket(f *Frame, x0, y0, cols, rows int, on bool) {
	if !on {
		return
	}
	for _, p := range [4][2]int{{x0 - 1, y0 - 1}, {x0 + cols, y0 - 1}, {x0 - 1, y0 + rows}, {x0 + cols, y0 + rows}} {
		f.Set(p[0], p[1], Cell{R: '·', FG: 234})
	}
}

// The fan's four targets, one per weapon: a 2x2 braille blob lower-left, a
// 3x5 bug upper-left, an (o_o) guy upper-right, a >_< guy lower-right. The
// corners stay clear of the slab box at any frame size.

// drawBootBrail is the Gunner's target: a 2x2 braille blob whose
// bottom-right cell sits at (tx,ty). It bobs between two dot patterns,
// flashes white while a tracer impact is fresh (hit = frames since, -1
// none), and on the killing blow (die = frames since, -1 alive) dissolves
// corner to corner into radial shards.
func drawBootBrail(f *Frame, tx, ty, s, hit, die int) {
	x0, y0 := tx-1, ty-1
	if die >= 0 {
		pos := [4][2]int{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
		for i, p := range pos {
			switch {
			case die == 0 || die < 2*i:
				f.Set(x0+p[0], y0+p[1], Cell{R: '⣿', FG: 255, Bold: true})
			case die < 2*i+2:
				f.Set(x0+p[0], y0+p[1], Cell{R: '⠿', FG: 240})
			case die < 2*i+3:
				f.Set(x0+p[0], y0+p[1], Cell{R: '·', FG: 236})
			}
		}
		if die < 8 {
			for k := 0; k < 8; k++ {
				a := float64(k) * math.Pi / 4
				r := die + 1
				c := 255
				if die >= 3 {
					c = 240
				}
				if die >= 6 {
					c = 236
				}
				f.Set(x0+1+int(math.Round(math.Cos(a)*float64(r))), y0+1+int(math.Round(math.Sin(a)*float64(r))), Cell{R: '·', FG: c})
			}
		}
		return
	}
	if hit >= 0 {
		for r := 0; r < 2; r++ {
			for c := 0; c < 2; c++ {
				f.Set(x0+c, y0+r, Cell{R: '⣿', FG: 255, Bold: true})
			}
		}
		return
	}
	ab := [2][2]rune{{'⣾', '⣾'}, {'⣿', '⣿'}}
	if s%6 < 3 {
		ab = [2][2]rune{{'⣿', '⣿'}, {'⣾', '⣾'}}
	}
	for r := 0; r < 2; r++ {
		for c := 0; c < 2; c++ {
			f.Set(x0+c, y0+r, Cell{R: ab[r][c], FG: 247})
		}
	}
}

// drawBootBug is the Cannon's target: a 3x5 bug centered at (tx,ty). It
// cycles its legs, and on the shell impact (squash = frames since, -1
// alive) collapses to a single row as the burst sprays it apart.
func drawBootBug(f *Frame, tx, ty, s, squash int) {
	if squash >= 0 {
		switch {
		case squash == 0:
			for c := -1; c <= 1; c++ {
				f.Set(tx+c, ty, Cell{R: '▓', FG: 255, Bold: true})
			}
		case squash == 1:
			for c := -1; c <= 1; c++ {
				f.Set(tx+c, ty, Cell{R: '░', FG: 240})
			}
			for i := 0; i < 6; i++ {
				a := float64(i) * math.Pi / 3
				f.Set(tx+int(math.Round(math.Cos(a))), ty+int(math.Round(math.Sin(a))), Cell{R: '▓', FG: 220})
			}
		case squash == 2:
			for i := 0; i < 6; i++ {
				a := float64(i) * math.Pi / 3
				f.Set(tx+int(math.Round(2*math.Cos(a))), ty+int(math.Round(2*math.Sin(a))), Cell{R: '▓', FG: 240})
			}
		}
		return
	}
	body := [4]string{"░▓░", "▓█▓", "█▓█", "▓█▓"}
	for r, row := range body {
		for c, ch := range []rune(row) {
			f.Set(tx-1+c, ty-2+r, Cell{R: ch, FG: 205})
		}
	}
	legs := "█ █"
	if s%8 < 4 {
		legs = "▓▓▓"
	}
	for c, ch := range []rune(legs) {
		f.Set(tx-1+c, ty+2, Cell{R: ch, FG: 205})
	}
}

// drawBootGuy is the Sniper's target: the ascii guy (o_o) centered at
// (tx,ty). It blinks, and on the beam hit (die = frames since, -1 alive)
// flashes white, then its five characters pop off one by one, spraying out
// and up.
func drawBootGuy(f *Frame, tx, ty, s, die int) {
	if die >= 0 {
		if die == 0 {
			for i, ch := range []rune("(o_o)") {
				f.Set(tx-2+i, ty, Cell{R: ch, FG: 255, Bold: true})
			}
			return
		}
		for i, ch := range []rune("(o_o)") {
			age := die - 1 - 2*i
			if age < 0 || age >= 9 {
				continue
			}
			c := 255
			if age >= 2 {
				c = 245
			}
			if age >= 5 {
				c = 236
			}
			x := tx - 2 + i + (i-2)*(age/2)
			y := ty - age/3
			f.Set(x, y, Cell{R: ch, FG: c, Bold: true})
		}
		return
	}
	guy := "(o_o)"
	if s%8 >= 6 {
		guy = "(-_-)"
	}
	for i, ch := range []rune(guy) {
		f.Set(tx-2+i, ty, Cell{R: ch, FG: 250})
	}
}

// drawBootZap is the Tesla's target: the ascii guy >_< centered at (tx,ty).
// It jitters and convulses under the arc, and on the lock (die = frames
// since, -1 alive) flashes white for two frames, then its characters pop
// off while a spark ring expands.
func drawBootZap(f *Frame, tx, ty, s, die int) {
	if die >= 0 {
		if die <= 1 {
			for i, ch := range []rune(">_<") {
				f.Set(tx-1+i, ty, Cell{R: ch, FG: 255, Bold: true})
			}
			return
		}
		for i, ch := range []rune(">_<") {
			age := die - 2 - i
			if age < 0 || age >= 5 {
				continue
			}
			c := 255
			if age >= 2 {
				c = 240
			}
			if age >= 4 {
				c = 236
			}
			x := tx - 1 + i + (i-1)*age
			y := ty - age/2
			f.Set(x, y, Cell{R: ch, FG: c, Bold: true})
		}
		if die < 8 {
			r := die - 1
			for k := 0; k < 8; k++ {
				a := float64(k) * math.Pi / 4
				c := 203
				if k%2 == 0 {
					c = 244
				}
				if die >= 5 {
					c = 236
				}
				f.Set(tx+int(math.Round(math.Cos(a)*float64(r))), ty+int(math.Round(math.Sin(a)*float64(r))), Cell{R: '·', FG: c})
			}
		}
		return
	}
	guy := ">_<"
	if s%7 == 5 {
		guy = "x_x"
	}
	jx := (s % 3) - 1
	for i, ch := range []rune(guy) {
		f.Set(tx-1+jx+i, ty, Cell{R: ch, FG: 244})
	}
}

// drawBootWeapons is the fan: each letter, left to right, fires a different
// game weapon at a creature in a frame corner — T (Gunner) at the braille
// blob lower-left, D (Cannon) at the 3x5 bug upper-left, E (Sniper) at the
// (o_o) guy upper-right, F (Tesla) at the >_< guy lower-right. t is 0..129.
func drawBootWeapons(f *Frame, w, h, off, t int) {
	drawBootBox(f, w, off, 98, 0, 234)
	drawBootSlab(f, w, off)
	drawBootGunner(f, w, h, off, t)
	drawBootCannon(f, w, h, off, t)
	drawBootSniper(f, w, h, off, t)
	drawBootTesla(f, w, h, off, t)
	// Settle beat: the box flickers once before the flash.
	if t == 126 || t == 127 {
		drawBootBox(f, w, off, 98, 0, 240)
	}
}

// drawBootGunner: T sprays five tracers at the braille blob in the
// lower-left corner. The letter flinches and flashes with every round; the
// blob dies on the fifth impact.
func drawBootGunner(f *Frame, w, h, off, t int) {
	const li = 0
	s := t
	if s < 0 || s > 59 {
		return
	}
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	mx, my := x0, off+4
	tx, ty := 3, h-4
	fl := int(math.Hypot(float64(tx-mx), float64(ty-my))/1.5) + 1
	if fl > 24 {
		fl = 24
	}
	last := 4*6 + fl
	die := s - (last + 1)
	if die < 0 {
		die = -1
	}
	hit := -1
	for r := 0; r < 4; r++ {
		if d := s - (6*r + fl + 1); d >= 0 && d < 3 {
			hit = d
		}
	}
	drawBootBrail(f, tx, ty, s, hit, die)
	if die < 0 {
		drawBootBracket(f, tx-1, ty-1, 2, 2, s%4 < 3)
	}
	for r := 0; r < 5; r++ {
		p := s - 6*r
		if p < 0 {
			continue
		}
		if p == 0 {
			drawBootLetter(f, w, off, li, 1, 0, 255, true)
			f.Set(mx, my, Cell{R: '+', FG: 255})
			continue
		}
		if p <= fl {
			x, y := shotPos(mx, my, tx, ty, float64(p-1)/float64(fl))
			f.Set(x, y, Cell{R: '█', FG: 255})
			if p >= 2 {
				if x1, y1 := shotPos(mx, my, tx, ty, float64(p-2)/float64(fl)); x1 >= 0 {
					f.Set(x1, y1, Cell{R: '·', FG: 46})
				}
			}
			if p >= 3 {
				if x2, y2 := shotPos(mx, my, tx, ty, float64(p-3)/float64(fl)); x2 >= 0 {
					f.Set(x2, y2, Cell{R: '·', FG: 236})
				}
			}
			continue
		}
		if p <= fl+2 {
			if p == fl+1 {
				f.Set(tx, ty, Cell{R: '*', FG: 255})
			} else {
				f.Set(tx, ty, Cell{R: '·', FG: 46})
			}
		}
	}
}

// drawBootCannon: D recoils, then fires a chunky shell at the 3x5 bug in
// the upper-left corner. The bug squashes flat under the impact as the AOE
// burst takes it.
func drawBootCannon(f *Frame, w, h, off, t int) {
	const li = 1
	s := t - 30
	if s < 0 || s > 60 {
		return
	}
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	mx, my := x0+15, off+2
	tx, ty := 3, 4
	fl := int(math.Hypot(float64(tx-mx), float64(ty-my))*2) + 1
	if fl > 24 {
		fl = 24
	}
	if s < 6+fl {
		drawBootBug(f, tx, ty, s, -1)
		drawBootBracket(f, tx-1, ty-2, 3, 5, s%4 < 3)
		if s < 6 {
			// Recoil: the letter is shoved down-right while the muzzle
			// gathers.
			drawBootLetter(f, w, off, li, 1, 1, 220, false)
			if s >= 3 {
				f.Set(mx, my, Cell{R: '•', FG: 220})
			}
			return
		}
		p := s - 6
		x, y := shotPos(mx, my, tx, ty, float64(p)/float64(fl))
		f.Set(x, y, Cell{R: '▓', FG: 220, Bold: true})
		if p >= 2 {
			if x1, y1 := shotPos(mx, my, tx, ty, float64(p-2)/float64(fl)); x1 >= 0 {
				f.Set(x1, y1, Cell{R: '░', FG: 240})
			}
		}
		if p >= 4 {
			if x2, y2 := shotPos(mx, my, tx, ty, float64(p-4)/float64(fl)); x2 >= 0 {
				f.Set(x2, y2, Cell{R: '·', FG: 236})
			}
		}
		return
	}
	if s > 6+fl+5 {
		return
	}
	q := s - 6 - fl
	drawBootBug(f, tx, ty, s, q)
	switch {
	case q == 0:
		f.Set(tx, ty, Cell{R: '█', FG: 255, Bold: true})
		f.Set(tx-2, ty, Cell{R: '·', FG: 220})
		f.Set(tx+2, ty, Cell{R: '·', FG: 220})
		f.Set(tx, ty-2, Cell{R: '·', FG: 220})
		f.Set(tx, ty+2, Cell{R: '·', FG: 220})
	case q == 1:
		f.Set(tx, ty, Cell{R: '*', FG: 220, Bold: true})
		for i := 0; i < 8; i++ {
			a := float64(i) * math.Pi / 4
			f.Set(tx+int(math.Round(math.Cos(a))), ty+int(math.Round(math.Sin(a))), Cell{R: '·', FG: 220})
		}
	case q == 2:
		for i := 0; i < 8; i++ {
			a := float64(i) * math.Pi / 4
			f.Set(tx+int(math.Round(2*math.Cos(a))), ty+int(math.Round(2*math.Sin(a))), Cell{R: '·', FG: 220})
		}
		f.Set(tx, ty, Cell{R: '·', FG: 240})
	case q == 3:
		for i := 0; i < 8; i++ {
			a := float64(i) * math.Pi / 4
			f.Set(tx+int(math.Round(3*math.Cos(a))), ty+int(math.Round(3*math.Sin(a))), Cell{R: '·', FG: 240})
		}
	case q <= 5:
		f.Set(tx, ty, Cell{R: '·', FG: 236})
	}
}

// drawBootSniper: E charges for fourteen frames (its top row filling with a
// ░▒▓█ ramp, the reticle locking on), then fires one piercing beam at the
// (o_o) guy in the upper-right corner, who pops off bead by bead.
func drawBootSniper(f *Frame, w, h, off, t int) {
	const li = 2
	s := t - 65
	if s < 0 || s > 24 {
		return
	}
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	mx, my := x0+26, off+2
	tx, ty := w-4, 3
	if s < 14 {
		fg := 203 + int(float64(s)/13*48)
		drawBootLetter(f, w, off, li, 0, 0, fg, false)
		for ci := 0; ci < 9; ci++ {
			at := ci * 14 / 9
			if s < at {
				continue
			}
			r := '░'
			switch age := s - at; {
			case age >= 9:
				r = '█'
			case age >= 6:
				r = '▓'
			case age >= 3:
				r = '▒'
			}
			f.Set(x0+li*11+ci, off+2, Cell{R: r, FG: 203})
		}
		drawBootGuy(f, tx, ty, s, -1)
		drawBootReticle(f, tx, ty, -1)
		return
	}
	die := s - 14
	if die < 2 {
		drawBootLetter(f, w, off, li, 0, 0, 255, true)
		f.Set(mx, my, Cell{R: '█', FG: 255, Bold: true})
	}
	if die < 4 {
		beamR, beamC := '█', 255
		if die == 1 {
			beamR, beamC = '▓', 203
		}
		if die >= 2 {
			beamR, beamC = '▓', 240
		}
		n := int(math.Hypot(float64(tx-mx), float64(ty-my))) + 1
		for k := 0; k <= n; k++ {
			x, y := shotPos(mx, my, tx, ty, float64(k)/float64(n))
			f.Set(x, y, Cell{R: beamR, FG: beamC})
		}
	}
	switch die {
	case 0:
		f.Set(tx-2, ty, Cell{R: '·', FG: 203})
		f.Set(tx+2, ty, Cell{R: '·', FG: 203})
	case 1:
		f.Set(tx, ty, Cell{R: '·', FG: 203})
	case 2, 3:
		f.Set(tx, ty, Cell{R: '·', FG: 236})
	case 4, 5:
		f.Set(tx, ty, Cell{R: '·', FG: 234})
	}
	if die <= 2 {
		drawBootReticle(f, tx, ty, die)
	}
	// The guy draws last so its white flash beats the impact markers.
	drawBootGuy(f, tx, ty, s, die)
}

// drawBootTesla: F charges, then chain lightning flickers toward the
// target at 2:15 — the main arc re-rolls its jitter every frame, with two
// branches forking off.
func drawBootTesla(f *Frame, w, h, off, t int) {
	const li = 3
	s := t - 100
	if s < 0 || s > 17 {
		return
	}
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	mx, my := x0+41, off+4
	tx, ty := w-4, h-4
	die := s - 12
	if die < 0 {
		die = -1
	}
	if s < 12 {
		if s%5 != 4 {
			drawBootReticle(f, tx, ty, -1)
		}
	} else if s < 15 {
		drawBootReticle(f, tx, ty, s-12)
	}
	if s < 4 {
		fg := 171 + int(float64(s)/3*80)
		drawBootLetter(f, w, off, li, 0, 0, fg, s == 3)
		if s >= 1 {
			for i := 0; i < 4; i++ {
				n := titleHash(mx, my, s*7+i)
				f.Set(mx+n%5-2, my+n%7-3, Cell{R: '·', FG: 203})
			}
		}
		drawBootZap(f, tx, ty, s, die)
		return
	}
	D := math.Hypot(float64(tx-mx), float64(ty-my))
	dx, dy := float64(tx-mx)/D, float64(ty-my)/D
	N := int(D / 2)
	if N < 4 {
		N = 4
	}
	for k := 0; k <= N; k++ {
		x, y := shotPos(mx, my, tx, ty, float64(k)/float64(N))
		j := titleHash(x, y, s*13+k)%3 - 1 // perpendicular jitter, re-rolled
		x, y = x-j, y+j
		n := titleHash(k, s*31, w) % 5
		r, c := '·', 171
		switch {
		case n < 2:
			r, c = '█', 255
		case n < 4:
			r, c = '▒', 203
		}
		f.Set(x, y, Cell{R: r, FG: c})
	}
	// Two branches fork off at 35% and 65% of the main arc.
	for bi := 0; bi < 2; bi++ {
		bx, by := shotPos(mx, my, tx, ty, [2]float64{0.35, 0.65}[bi])
		sgn := 1.0
		if bi == 1 {
			sgn = -1
		}
		vx, vy := sgn*-dy+0.5*dx, sgn*dx+0.5*dy
		norm := math.Hypot(vx, vy)
		vx, vy = vx/norm, vy/norm
		for k := 1; k <= 7; k++ {
			x := bx + int(math.Round(vx*1.3*float64(k)))
			y := by + int(math.Round(vy*1.3*float64(k)))
			c := 203
			if k == 7 {
				c = 244
			}
			f.Set(x, y, Cell{R: '·', FG: c})
		}
	}
	if s == 12 || s == 13 {
		f.Set(tx, ty, Cell{R: '*', FG: 255})
	} else if s <= 15 {
		f.Set(tx, ty, Cell{R: '·', FG: 203})
	}
	// The guy draws last so its white flash beats the impact markers.
	drawBootZap(f, tx, ty, s, die)
}

// drawBootFlash is the ignition: the whole frame glows as a dim grid while
// the slab burns white.
func drawBootFlash(f *Frame, w, h, off int) {
	for y := 0; y < h; y += 2 {
		for x := 0; x < w; x += 3 {
			f.Set(x, y, Cell{R: '·', FG: 234})
		}
	}
	const logoW = 42
	x0 := (w - logoW) / 2
	if x0 < 0 {
		x0 = 0
	}
	bevel := [5]rune{'█', '▓', '▓', '▓', '▒'}
	for _, s := range titlePenPath {
		f.Set(x0+s.relX, off+2+s.row, Cell{R: bevel[s.row], FG: 255, Bold: true})
	}
}

// drawBootSubtitle decodes the tagline left to right (two chars per frame)
// with a white head and a caret.
func drawBootSubtitle(f *Frame, w, off, t int) {
	r := []rune(titleTagline)
	n := 2 * t
	if n > len(r) {
		n = len(r)
	}
	x0 := (w - len(r)) / 2
	y := off + 8
	for i, ch := range r {
		switch {
		case i < n-1:
			f.Set(x0+i, y, Cell{R: ch, FG: 245})
		case i == n-1:
			f.Set(x0+i, y, Cell{R: ch, FG: 255, Bold: true})
		default:
			g := '·'
			if titleHash(x0+i, y, t)%9 == 0 {
				g = '+'
			}
			f.Set(x0+i, y, Cell{R: g, FG: 236})
		}
	}
	if n < len(r) {
		f.Set(x0+n, y, Cell{R: '█', FG: 251})
	}
}

// drawBootUI fades the rest of the chrome in over three staggered groups:
// frame border + footer, then the empty battlefield, then roster + best.
// Each group ghosts for eight frames, then lands on the exact standby form.
func drawBootUI(f *Frame, w, h, off, t, frame int, pal Colors) {
	lit := (frame/15)%2 == 0
	// Group A: the outer frame.
	if t >= 0 {
		if t < 8 {
			drawRoundedBox(f, 0, 0, w, h, 234)
			embedSegment(f, 0, 2, "TDEF", '┐', '┌', 234, 244, true)
			ghostRun(f, (w-44)/2, (w+44)/2-1, h-1)
		} else {
			drawRoundedBox(f, 0, 0, w, h, pal.Path)
			embedSegment(f, 0, 2, "TDEF", '┐', '┌', pal.Path, pal.Bright, true)
			drawFooter(f, titleFooter(), lit, pal)
			drawTitleSig(f)
		}
	}
	// Group B: the empty battlefield.
	if t >= 6 {
		y0 := off + demoTop
		if y0 > 0 && y0+demoRows < h-1 {
			if t < 14 {
				drawRoundedBox(f, 1, y0, w-2, demoRows, 234)
				pathY := off + demoPath
				for x := 2; x <= w-3; x++ {
					f.Set(x, pathY, Cell{R: '·', FG: 234})
				}
			} else {
				drawTitleEmptyBox(f, w, h, off, frame, pal)
			}
		}
	}
	// Group C: roster + best.
	if t >= 12 {
		if t < 20 {
			ghostRun(f, (w-41)/2, (w+41)/2-1, off+14)
			ghostRun(f, (w-41)/2, (w+41)/2-1, off+15)
			ghostRun(f, (w-30)/2, (w+30)/2-1, off+17)
		} else {
			drawTitleRoster(f, w, h, off, pal)
			drawTitleBest(f, w, off, nil)
		}
	}
}

// ghostRun lays a dim run of dots (a "pre-render" ghost of a text row).
// The footer call targets the bottom border row on purpose.
func ghostRun(f *Frame, x0, x1, y int) {
	for x := x0; x <= x1; x++ {
		if x < 0 || x >= f.W {
			continue
		}
		f.Set(x, y, Cell{R: '·', FG: 234})
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
	drawTitleChrome(f, w, h, off, scores, pal)
	drawTitleDemo(f, w, h, off, fr, pal)
}

// titleSig is the signature, embedded in the bottom border's right section
// — mirroring the TDEF embed in the top border, but muted.
const titleSig = "by 0xbenc"

// drawTitleSig draws the signature into the bottom border.
func drawTitleSig(f *Frame) {
	x := f.W - len(titleSig) - 2
	for i, ch := range []rune(titleSig) {
		f.Set(x+i, f.H-1, Cell{R: ch, FG: 238})
	}
}

// drawTitleChrome is the title content shared by the standby and the
// battle: logo, tagline, roster, best line, signature.
func drawTitleChrome(f *Frame, w, h, off int, scores map[string]int, pal Colors) {
	drawTitleLogo(f, w, off, -1)
	drawTitleTagline(f, w, off, len(titleTagline))
	drawTitleRoster(f, w, h, off, pal)
	drawTitleBest(f, w, off, scores)
	drawTitleSig(f)
}

func drawTitleBest(f *Frame, w, off int, scores map[string]int) {
	if yy := off + 17; yy > 0 && yy < f.H-1 {
		best, name := titleBest(scores)
		if best > 0 {
			centerPut(f, yy, fmt.Sprintf("★ best %d — %s", best, name), 220, false)
		} else {
			centerPut(f, yy, " no scores yet ", 238, false)
		}
	}
}

// drawTitleEmptyBox is the idle battlefield: the sub-box and its path with
// the ambient energy packet — no THE SIEGE/WAVE text, no towers, no
// enemies.
func drawTitleEmptyBox(f *Frame, w, h, off, frame int, pal Colors) {
	y0 := off + demoTop
	if y0 <= 0 || y0+demoRows >= h-1 {
		return
	}
	drawRoundedBox(f, 1, y0, w-2, demoRows, pal.Path)
	pathY := off + demoPath
	L := w - 5
	for x := 2; x <= w-3; x++ {
		f.Set(x, pathY, Cell{R: '·', FG: 240})
	}
	pk := frame % L
	for i := -1; i <= 1; i++ {
		if px := 2 + (pk+i+L)%L; px >= 2 && px <= w-3 {
			f.Set(px, pathY, Cell{R: '·', FG: 251})
		}
	}
	f.Set(2, pathY, Cell{R: '▶', FG: 46, Bold: true})
	f.Set(w-3, pathY, Cell{R: 'E', FG: 196, Bold: true})
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
// [ THE SIEGE ] sub-box on band rows 9-13: five towers (three below
// the path, two above) fire on scripted cooldowns at scripted enemy waves.
// Everything is a pure function of fr.
func drawTitleDemo(f *Frame, w, h, off, fr int, pal Colors) {
	y0 := off + demoTop
	if y0 <= 0 || y0+demoRows >= h-1 {
		return
	}
	drawSubBox(f, 1, y0, w-2, demoRows, "THE SIEGE", pal)
	waveText, waveFG := "WAVE 1", pal.Bright
	switch {
	case fr >= titleOverloadEnd-13: // the boss reaches the exit
		waveText, waveFG = "BREACH", 167
	case fr >= 900:
		waveText, waveFG = "WAVE 3", 167
	case fr >= 480:
		waveText, waveFG = "WAVE 2", pal.Bright
	}
	// The THE SIEGE label occupies x 2..12 (┐ + 9 + ┌); the wave
	// segment follows at x 14.
	embedSegment(f, y0, 14, waveText, '┐', '┌', pal.Path, waveFG, true)

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
// and the empty battlefield returns. At the last frame the screen is
// exactly the standby, so the attract loop is seamless.
func drawTitleReboot(f *Frame, w, h, frame, t int, scores map[string]int, pal Colors) {
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
		drawTitleEmptyBox(f, w, h, off, frame, pal)
	}
	if p >= 0.80 {
		drawTitleRoster(f, w, h, off, pal)
	}
	if p >= 0.88 {
		drawTitleBest(f, w, off, scores)
		drawTitleSig(f)
	}
	if p >= 0.95 {
		drawFooter(f, titleFooter(), (frame/15)%2 == 0, pal)
	}
}

// ---------------------------------------------------------------- menu

// MenuItems is the main menu, in order. Start is the lair itself: the
// overworld, where Grak walks the floors and descends into one.
var MenuItems = []string{"Start", "Quick Play", "Help", "High Scores", "Quit"}

// menuLayout returns the item rows of the main menu. The item block (2*len-1
// rows tall at 2-row stride) is centered in the content band so the menu
// sits in the middle of the frame at any height.
func menuLayout(h int) []int {
	off := screenOff(h)
	top, bottom := off+1, off+17
	bh := 2*len(MenuItems) - 1
	start := top + (bottom-top+1-bh)/2
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
	f := screenBox(w, h, "GRAK'S LEDGER", []fseg{{key: "esc", text: " back"}}, true, pal)
	off := screenOff(h)
	rows := [][2]string{
		{"move", "arrows / wasd"},
		{"place", "1-7 pick · enter or click"},
		{"upgrade", "u"},
		{"sell", "x"},
		{"target", "t (game) · t relics (lair)"},
		{"wave", "n (early = bonus gold)"},
		{"pause", "p"},
		{"speed", "f or wheel"},
		{"descend", "enter, on a floor (the lair)"},
		{"renown", "tab (the lair)"},
		{"seed", "0-9 on the Depths (the lair)"},
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
	if y := off + 3; y > 0 && y < h-1 {
		centerPut(f, y, "how to hold the lair against twenty expeditions", pal.Dim, false)
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
// len(Levels) is the maze row. Seed is raw typed text; the empty string
// means "random" at start time.
type LSState struct {
	Levels  []string
	Cursor  int
	Seed    string
	Diff    int
	Err     string
	Preview *game.Map
}

const lsMazeRow = "the Unmapped Depths"

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
		if w := len(levelDisplayName(n)); w > maxW {
			maxW = w
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
			name = levelDisplayName(v.Levels[i])
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
			drawMapPreview(f, v.Preview, pal, l, 0)
		} else if guard(p) {
			centerPut(f, p, " (preview needs more room) ", 238, false)
		}
	}
	return f
}
