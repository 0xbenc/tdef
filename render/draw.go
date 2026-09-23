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

	Level string // level id for the header, e.g. "winding" or "maze1234"

	BestScore int
	NewBest   bool

	// ToLair marks a run that started from the overworld, so the game-over
	// box offers "esc lair" (back to the map) as well as restart/quit.
	ToLair bool

	// EndAtFrame is the ambient frame the run ended (victory/defeat) at, set
	// by the caller the moment the status flips; 0 while the run is going (or
	// unknown — e.g. headless capture — in which case the end box shows at
	// once). The end cinematics are pure functions of (frame - EndAtFrame),
	// so they play to the end even though g.Time has frozen.
	EndAtFrame int
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
	Wall, WallHi, WallLo   int
	Path, RoadBG, RoadLine int
	Grass, GrassTuft       int
	Spawn, Exit, Frost     int
	Gold, Dim, Bright      int
	Tower                  [game.TowerCount]int
	Enemy                  [game.EnemyCount]int
	Beam                   [game.TowerCount]int
}

func Palette() Colors {
	return Colors{
		Wall: 235, WallHi: 237, WallLo: 233,
		Path: 240, RoadBG: 238, RoadLine: 246,
		Grass: 23, GrassTuft: 34,
		Spawn: 201, Exit: 196, Frost: 51,
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
	var kind game.TowerKind
	var haveKind bool
	if ui.PlacingOn {
		c := ui.Cursor.Center()
		rng = game.TowerSpecs[ui.Placing].Range[0]
		center = &c
		kind, haveKind = ui.Placing, true
	} else if ui.Selected >= 0 {
		if t := g.Tower(ui.Selected); t != nil {
			c := t.Pos()
			rng = t.Range()
			center = &c
			kind, haveKind = t.Kind, true
		}
	}
	if center == nil {
		return
	}
	ringCol := pal.Dim
	if haveKind {
		ringCol = pal.Beam[kind]
	}
	for y := 0; y < g.Map.H; y++ {
		for x := 0; x < g.Map.W; x++ {
			if g.Map.At(game.Vec{X: x, Y: y}) == game.CellWall {
				continue
			}
			p := game.Pos{X: float64(x) + 0.5, Y: float64(y) + 0.5}
			d := p.Dist(*center)
			if d > rng {
				continue
			}
			cx, cy := l.center(x, y)
			if d > rng-0.55 {
				putOver(f, cx, cy, '·', ringCol) // the boundary ring
			} else if cellHash(x, y) < 0.22 {
				putOver(f, cx, cy, '·', pal.Path) // a faint interior so the disc reads as an area
			}
		}
	}
}

// putOverBold is putOver with a bold flag, for overlay marks that should read
// heavy (projectile shells, beam energy).
func putOverBold(f *Frame, x, y int, r rune, fg int, bold bool) {
	if x < 0 || y < 0 || x >= f.W || y >= f.H {
		return
	}
	cc := f.C[y*f.W+x]
	cc.R, cc.FG, cc.Bold = r, fg, bold
	f.C[y*f.W+x] = cc
}

// projGlyph returns a projectile's glyph and whether it is bold: shells
// (cannon, mortar) are heavy filled circles — mortar bolder still — while
// gunner and flak fire light tracers.
func projGlyph(k game.TowerKind) (rune, bool) {
	switch k {
	case game.TowerCannon:
		return '●', false
	case game.TowerMortar:
		return '●', true
	default:
		return '·', false
	}
}

// drawProjectile renders one projectile: a kind-specific glyph in its beam
// colour with a short dim trail behind it so the motion reads.
func drawProjectile(f *Frame, l Layout, pal Colors, p *game.Projectile) {
	x, y := l.FX(p.Pos.X), l.FY(p.Pos.Y)
	glyph, bold := projGlyph(p.Kind)
	putOverBold(f, x, y, glyph, pal.Beam[p.Kind], bold)
	dx, dy := p.LastPos.X-p.Pos.X, p.LastPos.Y-p.Pos.Y
	if d := math.Hypot(dx, dy); d > 0.05 {
		ux, uy := dx/d, dy/d
		putOver(f, l.FX(p.Pos.X-ux*0.45), l.FY(p.Pos.Y-uy*0.45), '·', 240)
	}
}

// drawBeamSeg draws one beam segment as a bright, continuous energy line. It
// overwrites the terrain glyphs but preserves the background, so a bolt reads
// cleanly over road and rock.
func drawBeamSeg(f *Frame, l Layout, a, b game.Pos, color int) {
	ax, ay := l.FX(a.X), l.FY(a.Y)
	bx, by := l.FX(b.X), l.FY(b.Y)
	dx, dy := bx-ax, by-ay
	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy)))) * 2
	if steps < 1 {
		steps = 1
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		putOverBold(f, ax+int(float64(dx)*t), ay+int(float64(dy)*t), '·', color, true)
	}
}

// isGroundRune reports whether a rendered rune is terrain (buildable/roamable
// ground) rather than an entity: blank, a grass tuft, or one of the road
// connectors. The cursor overlays a marker on ground but only bolds entities
// (towers, enemies, the spawn rift and the lair heart).
func isGroundRune(r rune) bool {
	if r == 0 || r == ' ' || r == '·' {
		return true
	}
	for _, c := range "─│┌┐└┘├┤┬┴┼" {
		if r == c {
			return true
		}
	}
	return false
}

// putOver sets a rune and its colour on the cell at (x,y) while PRESERVING the
// cell's background — so overlay marks (range dots, muzzle tracers) sit on top
// of the textured terrain instead of punching a black hole in it.
func putOver(f *Frame, x, y int, r rune, fg int) {
	if x < 0 || y < 0 || x >= f.W || y >= f.H {
		return
	}
	cc := f.C[y*f.W+x]
	cc.R = r
	cc.FG = fg
	cc.Bold = false
	f.C[y*f.W+x] = cc
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

// drawTower renders one tower: its glyph on a dark pad (the pad fills the whole
// block at 2x+, so the tower reads as a structure set into the ground), corner
// brackets when it is the selected tower, and a muzzle flash + tracer while it
// is firing (t.Flash > 0).
func drawTower(f *Frame, g *game.State, ui *UI, pal Colors, l Layout, t *game.Tower) {
	sel := ui.Selected == t.ID
	glyph := t.Spec().Short
	// An upgraded tower's level lives on its own pedestal, never above it (the
	// cell above a tower is usually the road — the lane the horde walks). The
	// glyph goes white and the base charges in the tower's own colour: a dark
	// tint at level 2, the full colour once maxed. A selected tower just adds
	// its corner bracket — the colours are not what mark selection.
	fg, bg := pal.Tower[t.Kind], 235
	if t.Level >= 3 {
		fg, bg = pal.Bright, pal.Tower[t.Kind]
	} else if t.Level == 2 {
		fg, bg = pal.Bright, baseColor(pal.Tower[t.Kind], 55)
	}
	if sel {
		fg, bg = pal.Bright, pal.Tower[t.Kind]
	}
	l.block(f, t.Cell.X, t.Cell.Y, Cell{R: ' ', BG: bg})
	cx, cy := l.center(t.Cell.X, t.Cell.Y)
	f.Set(cx, cy, Cell{R: glyph, FG: fg, BG: bg, Bold: true})
	if sel {
		drawSelectionBracket(f, l, t.Cell)
	}
	if t.Flash > 0 {
		drawMuzzle(f, l, pal, t)
	}
}

// drawSelectionBracket draws bright corner brackets just outside the tower's
// block so the active tower is unmistakable at a glance.
func drawSelectionBracket(f *Frame, l Layout, v game.Vec) {
	const c = 255
	x0, y0 := l.X(v.X), l.Y(v.Y)
	x1, y1 := x0+l.Scale-1, y0+l.Scale-1
	f.Set(x0-1, y0-1, Cell{R: '╭', FG: c, Bold: true})
	f.Set(x1+1, y0-1, Cell{R: '╮', FG: c, Bold: true})
	f.Set(x0-1, y1+1, Cell{R: '╰', FG: c, Bold: true})
	f.Set(x1+1, y1+1, Cell{R: '╯', FG: c, Bold: true})
}

// drawMuzzle renders a firing tower's muzzle flash: the glyph flares white and
// a short tracer kicks out toward the target it is firing at (t.FlashTo). The
// tracer preserves the terrain background so it reads over road and rock.
func drawMuzzle(f *Frame, l Layout, pal Colors, t *game.Tower) {
	cx, cy := l.center(t.Cell.X, t.Cell.Y)
	f.Set(cx, cy, Cell{R: t.Spec().Short, FG: pal.Bright, Bold: true})
	sx, sy := float64(t.Cell.X)+0.5, float64(t.Cell.Y)+0.5
	dx, dy := t.FlashTo.X-sx, t.FlashTo.Y-sy
	d := math.Hypot(dx, dy)
	if d < 0.1 {
		return
	}
	ux, uy := dx/d, dy/d
	for i := 1; i <= 4; i++ {
		off := 0.5 * float64(i) // 0.5..2 cells out
		x, y := l.FX(sx+ux*off), l.FY(sy+uy*off)
		if x == cx && y == cy {
			continue // never overwrite the tower glyph itself
		}
		putOver(f, x, y, '·', pal.Beam[t.Kind])
	}
}

// tankKind reports whether an enemy kind is a "tank" — high enough HP that its
// health bar is worth showing even at full health, so the player can read how
// much a Paladin/Centurion/boss/Necromancer can take before it breaks.
func tankKind(k game.EnemyKind) bool {
	switch k {
	case game.EnemyTank, game.EnemyShield, game.EnemyBoss, game.EnemySplitter:
		return true
	}
	return false
}

// drawEnemy renders one enemy. Its glyph keeps its own kind colour (identity is
// the point — HP is read off the bar, not the glyph); it tints cyan while frost
// slowed, flashes white when hit, and the boss gets a caged presence. Tanks
// always show their HP bar; everyone else only when wounded.
func drawEnemy(f *Frame, g *game.State, pal Colors, l Layout, e *game.Enemy) {
	x, y := l.FX(e.Pos.X), l.FY(e.Pos.Y)
	if x < 0 || y < 0 || x >= f.W || y >= f.H {
		return
	}
	hp := e.HP / e.MaxHP
	kind := e.Kind
	fg, bold := pal.Enemy[kind], kind == game.EnemyBoss
	switch {
	case e.HitTTL > g.Time:
		fg, bold = 231, true // hit flash: brief white
	case e.Slowed(g.Time):
		fg, bold = pal.Frost, true // frost slowed: cyan
	}
	bg := 0
	if kind == game.EnemyBoss {
		bg = 53 // the boss sits on a dark pad for presence
	}
	f.Set(x, y, Cell{R: game.EnemySpecs[kind].Short, FG: fg, BG: bg, Bold: bold})
	if kind == game.EnemyBoss {
		// a cage of rails either side of the boss (never over the HP bar above)
		putOver(f, x-1, y, '│', 204)
		putOver(f, x+1, y, '│', 204)
	}
	showBar := hp < 1 || tankKind(kind)
	if showBar && y-1 >= l.Oy {
		drawHPBar(f, x, y-1, hp)
	}
}

type MenuSlot struct {
	Kind game.TowerKind
	X    int
	Y    int
	W    int // rendered label width ("k Name cost")
}

// towerInfo is the menu line for the selected tower. drawMenu elides it
// with fitMsg when it outgrows the bottom border (TestTowerInfoFitsFrame).
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

// levelDisplayName maps a level id to the name shown in the header and on
// the level-select screen; unknown ids pass through unchanged.
func levelDisplayName(id string) string {
	switch {
	case id == "hub":
		return "the Rotunda"
	case id == "winding":
		return "the Long Halls"
	case id == "garden":
		return "the Sunken Garden"
	case id == "canyon":
		return "the Rift"
	case strings.HasPrefix(id, "maze"):
		return "the Unmapped Depths"
	}
	return id
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
	lore := "The lair is held. Malgrath endures."
	if g.Status != game.StatusVictory {
		lore = "Malgrath has fallen. The lair is clean."
	}
	putString(f, bx+(bw-len([]rune(lore)))/2, by+8, lore, 244, bg, false)
	verdict := endVerdict(g)
	putString(f, bx+(bw-len([]rune(verdict)))/2, by+9, verdict, 245, bg, false)
	parts := []struct {
		s  string
		fg int
	}{
		{"r", 167}, {" restart | ", 251},
	}
	if ui.ToLair {
		parts = append(parts,
			struct {
				s  string
				fg int
			}{"esc", 220},
			struct {
				s  string
				fg int
			}{" lair | ", 251})
	}
	parts = append(parts,
		struct {
			s  string
			fg int
		}{"q", 167},
		struct {
			s  string
			fg int
		}{" quit", 251})
	total := 0
	for _, p := range parts {
		total += len(p.s)
	}
	x := bx + (bw-total)/2
	for _, part := range parts {
		for _, ch := range part.s {
			f.Set(x, by+10, Cell{R: ch, FG: part.fg, BG: bg, Bold: part.fg == 167 || part.fg == 220})
			x++
		}
	}
}

// endVerdict is the lair's assessment of the run, shown under the lore on the
// end box — it reads the result the way the lair would.
func endVerdict(g *game.State) string {
	if g.Status == game.StatusVictory {
		switch {
		case g.TotalLeaks == 0:
			return "A flawless hold — not one breach."
		case g.TotalLeaks <= 5:
			return "A steady hold. The lair endures."
		default:
			return "A hard-fought hold. The lair feels it."
		}
	}
	switch {
	case g.Wave >= 15:
		return fmt.Sprintf("The lair fell late, on wave %d.", g.Wave)
	case g.Wave >= 8:
		return fmt.Sprintf("The lair held to wave %d. So close.", g.Wave)
	default:
		return fmt.Sprintf("The lair fell early, on wave %d.", g.Wave)
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
// the slots span the interior; the longest label ("5 Lightning Mage 200",
// 18 cols) fits every cell at the minimum width 62. The renderer and the
// click handler share this function.
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
			runs: []headerRun{{levelDisplayName(ui.Level) + " · " + diffName(g.Diff), pal.Bright, true}},
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
	case g.Status == game.StatusRunning && g.Wave == 0 && len(g.Towers) == 0:
		// the siege is held until the player commits their first tower
		wave = []headerRun{{"build a tower", 220, true}}
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
		case strings.Contains(msg, "breach"):
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
			putString(f, x, y, p[0], 167, 0, true) // the slot-key red
			x += len([]rune(p[0]))
			putString(f, x, y, p[1], pal.Bright, 0, false)
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
// comment above). frame is the ambient 30fps tick (it keeps advancing while
// paused and after game-over, unlike g.Time, which freezes); the beats and
// cinematics are pure functions of it so they always play to the end.
func Render(g *game.State, ui *UI, pal Colors, tw, th, frame int) *Frame {
	l := GameLayout(g.Map.W, g.Map.H, tw, th)
	f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
	drawRoundedBox(f, 0, 0, l.W, l.H, pal.Path)
	drawHeader(f, g, ui, pal)
	drawMapPreview(f, g.Map, pal, themeForLevel(ui.Level), l, frame)
	drawRange(f, g, ui, pal, l)
	// Towers: a colored glyph on a dark pad (a small pedestal at 2x+); an
	// upgraded tower's pedestal charges in its colour and the glyph goes white,
	// corner brackets on the selected one, and a muzzle flash while it fires.
	for _, t := range g.Towers {
		drawTower(f, g, ui, pal, l, t)
	}
	for _, p := range g.Projectiles {
		drawProjectile(f, l, pal, p)
	}
	for _, bm := range g.Beams {
		for i := 1; i < len(bm.From); i++ {
			drawBeamSeg(f, l, bm.From[i-1], bm.From[i], pal.Beam[bm.Kind])
		}
	}
	for _, e := range g.Enemies {
		drawEnemy(f, g, pal, l, e)
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
		f.Set(ex, ey, Cell{R: '♥', FG: 231, BG: 196, Bold: true})
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
			// Terrain (wall/grass/road) gets the cursor marker; entity glyphs
			// (towers, enemies, the spawn rift, the lair heart) just bold.
			if isGroundRune(cc.R) {
				cc.R = '◻'
				cc.FG = pal.Dim
			} else {
				cc.Bold = true
			}
			f.Set(x, y, cc)
		}
	}
	drawBeats(f, g, pal, l, frame)
	drawMenu(f, g, ui, pal)
	drawEndSequence(f, g, ui, pal, l, frame)
	return f
}

// endBeatFrames is how long the end cinematic plays before the stats box
// settles in (3s at the 30fps ambient clock).
const endBeatFrames = 90

// drawBeats renders the in-run beats, all pure functions of the state: the
// leak edge pulse, the wave-start banner, and the boss entrance. They play
// while the run is going (the end cinematics are separate, in
// drawEndSequence).
func drawBeats(f *Frame, g *game.State, pal Colors, l Layout, frame int) {
	if g.Status != game.StatusRunning {
		return
	}
	midY := l.Oy + (g.Map.H*l.Scale)/2
	if g.LeakFlash > 0 {
		drawEdgePulse(f, 196) // the frame throbs red while the exit is struck
	}
	if g.WaveActive && g.Time-g.WaveStart < 1.8 {
		text := fmt.Sprintf("WAVE %d", g.Wave)
		if g.Wave == 1 {
			text = "THE SIEGE BEGINS"
		}
		fg := pal.Bright
		if g.Time-g.WaveStart > 1.1 {
			fg = pal.Dim // the banner fades as the wave gets under way
		}
		drawCenterBanner(f, midY, text, fg)
	}
	for _, e := range g.Enemies {
		if e.Kind == game.EnemyBoss && !e.Dead && !e.Leaked && e.Prog < 4 {
			drawCenterBanner(f, midY, "THE PLAYER", 204)
			drawEdgePulse(f, 204) // a regal purple pulse as the boss appears
			break
		}
	}
}

// drawCenterBanner centers a beat banner on row y with a dark backing so it
// reads over the terrain.
func drawCenterBanner(f *Frame, y int, text string, fg int) {
	w := len([]rune(text))
	x := (f.W - w) / 2
	if x < 1 {
		x = 1
	}
	for i := 0; i < w; i++ {
		f.Set(x+i, y, Cell{R: ' ', BG: 232})
	}
	putString(f, x, y, text, fg, 232, true)
}

// drawEdgePulse recolors the left/right frame edges (in the playfield region,
// leaving the header and menu alone) in color — the "the lair is being hit"
// and "the boss has arrived" pulses.
func drawEdgePulse(f *Frame, color int) {
	for y := ChromeTop; y < f.H-ChromeBot; y++ {
		f.Set(0, y, Cell{R: '│', FG: color, Bold: true})
		f.Set(f.W-1, y, Cell{R: '│', FG: color, Bold: true})
	}
}

// drawEndSequence renders the end of the run. When the caller recorded the end
// frame (ui.EndAtFrame > 0) a short cinematic plays — a sweep and the lair's
// verdict — before the stats box settles in; otherwise (a headless capture,
// where the end frame is unknown) the box shows at once.
func drawEndSequence(f *Frame, g *game.State, ui *UI, pal Colors, l Layout, frame int) {
	if g.Status == game.StatusRunning {
		return
	}
	won := g.Status == game.StatusVictory
	if ui.EndAtFrame <= 0 {
		drawGameOver(f, g, ui, pal)
		return
	}
	if beat := frame - ui.EndAtFrame; beat < endBeatFrames {
		drawEndBeat(f, g, pal, l, beat, won)
		return
	}
	drawGameOver(f, g, ui, pal)
}

// drawEndBeat is the end cinematic: a bright sweep races across the playfield
// while the lair's verdict decodes in — gold for a hold, red and final for a
// fall.
func drawEndBeat(f *Frame, g *game.State, pal Colors, l Layout, beat int, won bool) {
	mw := g.Map.W * l.Scale
	sx := l.Ox + beat*mw/endBeatFrames
	sweepCol := pal.Bright
	if !won {
		sweepCol = 167
	}
	for y := l.Oy; y < l.Oy+g.Map.H*l.Scale; y++ {
		putOver(f, sx, y, '│', sweepCol)
	}
	text, fg := "THE LAIR HOLDS", pal.Gold
	if !won {
		text, fg = "THE LAIR FALLS", 167
	}
	reveal := beat * len(text) / 45
	if reveal > len(text) {
		reveal = len(text)
	}
	midY := l.Oy + (g.Map.H*l.Scale)/2
	drawCenterBanner(f, midY, text[:reveal], fg)
}

// GameFrame renders the in-game frame for a tw×th terminal, or a centered
// "enlarge" notice when the terminal is smaller than MinFrame.
func GameFrame(g *game.State, ui *UI, pal Colors, tw, th, frame int) *Frame {
	minW, minH := MinFrame(g.Map.W, g.Map.H)
	if tw > 0 && th > 0 && (tw < minW || th < minH) {
		return RenderTooSmall(tw, th, minW, minH)
	}
	return Render(g, ui, pal, tw, th, frame)
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
