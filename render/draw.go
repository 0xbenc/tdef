package render

import (
	"fmt"
	"math"

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
	Scale     int // playfield scale (1-4), computed once at boot

	BestScore int
	NewBest   bool
}

const (
	HUDRows = 2
	FrameW  = 62
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

// ComputeScale picks the largest integer scale (1-4) whose frame fits the
// terminal. Called once at boot; no live resize support.
func ComputeScale(mW, mH, tw, th int) int {
	if tw <= 0 || th <= 0 {
		return 1
	}
	s := (tw - 2) / mW
	if v := (th - HUDRows - 4) / mH; v < s {
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

func ComputeLayout(mW, mH, scale int) Layout {
	if scale < 1 {
		scale = 1
	}
	mapW, mapH := mW*scale, mH*scale
	fw := FrameW
	if mapW+2 > fw {
		fw = mapW + 2
	}
	fh := HUDRows + mapH + 4
	return Layout{Ox: (fw - mapW) / 2, Oy: HUDRows, Scale: scale, W: fw, H: fh}
}

// MinFrame returns the smallest frame (scale 1) needed for a map of the given
// size. A terminal smaller than this cannot show the playfield.
func MinFrame(mW, mH int) (w, h int) {
	w = FrameW
	if mW+2 > w {
		w = mW + 2
	}
	h = HUDRows + mH + 4
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

func Render(g *game.State, ui *UI, pal Colors) *Frame {
	l := ComputeLayout(g.Map.W, g.Map.H, ui.Scale)
	f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
	for y := 0; y < g.Map.H; y++ {
		for x := 0; x < g.Map.W; x++ {
			v := game.Vec{X: x, Y: y}
			switch g.Map.At(v) {
			case game.CellWall:
				l.block(f, x, y, Cell{R: ' ', FG: 0, BG: pal.Wall})
			case game.CellPath:
				l.block(f, x, y, Cell{R: '·', FG: pal.Path, BG: 0})
			case game.CellGrass:
				l.block(f, x, y, Cell{R: ' ', FG: 0, BG: pal.Grass})
			}
		}
	}
	sx, sy := l.center(g.Map.Spawn.X, g.Map.Spawn.Y)
	f.Set(sx, sy, Cell{R: '▶', FG: 46, Bold: true})
	ex, ey := l.center(g.Map.Exit.X, g.Map.Exit.Y)
	f.Set(ex, ey, Cell{R: 'E', FG: 196, Bold: true})
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
		// Skip the bar when it would land on the HUD rows above the map.
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
	drawHUD(f, g, ui, pal)
	drawMenu(f, g, ui, pal)
	if g.Status != game.StatusRunning {
		drawGameOver(f, g, ui, pal)
	}
	return f
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

// drawHUD renders the two top rows. Both must stay within FrameW (62) at
// scale 1 — the break line used to exceed it from wave 17 on (long theme +
// preview + countdown). TestHUDLinesFitFrame pins the worst cases.
func drawHUD(f *Frame, g *game.State, ui *UI, pal Colors) {
	pause := ""
	if ui.Paused {
		pause = " ⏸ "
	}
	combo := ""
	if g.Combo >= 5 {
		combo = fmt.Sprintf("  ⚡%d", g.Combo)
	}
	inBreak := g.Status == game.StatusRunning && !g.WaveActive && g.Wave < game.MaxWaves
	// A telegraph during the break takes the whole top line (so it isn't
	// truncated); the wave indicator moves to the stat line.
	telegraph := inBreak && ui.Message != ""
	line0 := " tdef "
	line1 := ""
	switch {
	case g.Status == game.StatusRunning && g.WaveActive:
		line0 += fmt.Sprintf("Wave %d/%d%s", g.Wave, game.MaxWaves, combo)
	case inBreak && !telegraph:
		nw := g.Wave + 1
		in := int(g.NextWaveAt-g.Time) + 1
		if in < 0 {
			in = 0
		}
		bonus := game.EarlyBonus(g.Wave)
		line0 += fmt.Sprintf("next wave %d (%s) in %ds [n +%dg]", nw, game.WaveTheme(nw), in, bonus)
		if ui.Message == "" {
			// Composition preview on the stat row; a transient message
			// (e.g. "wave 18 cleared +93g") takes its place instead.
			line1 = fmt.Sprintf(" →%d: %s", nw, game.WavePreview(nw))
		}
	}
	if telegraph {
		line0 = "  " + ui.Message
		nw := g.Wave + 1
		line1 = fmt.Sprintf(" →%d %s: %s", nw, game.WaveTheme(nw), game.WavePreview(nw))
	}
	line0 += pause
	if ui.Message != "" && !inBreak {
		line0 += "  " + ui.Message
	} else if ui.Message != "" && inBreak && !telegraph {
		line1 = "  " + ui.Message
	}
	line1 += fmt.Sprintf("   ⛁ %d  ♥ %d  ★ %d  x%d", g.Gold, g.Lives, g.Score, ui.Speed)
	putString(f, 0, 0, line0, pal.Bright, 0, true)
	putString(f, 0, 1, line1, pal.Gold, 0, false)
}

type MenuSlot struct {
	Kind game.TowerKind
	X    int
	Y    int
}

var MenuSlots = []MenuSlot{
	{game.TowerGunner, 1, 0}, {game.TowerCannon, 14, 0}, {game.TowerFrost, 28, 0}, {game.TowerSniper, 40, 0},
	{game.TowerTesla, 1, 1}, {game.TowerMortar, 14, 1}, {game.TowerFlak, 28, 1},
}

func drawMenu(f *Frame, g *game.State, ui *UI, pal Colors) {
	menuTop := f.H - 4
	for _, slot := range MenuSlots {
		k := slot.Kind
		spec := game.TowerSpecs[k]
		afford := g.Gold >= spec.Cost[0]
		c := pal.Tower[k]
		if !afford {
			c = pal.Dim
		}
		y := menuTop + slot.Y
		label := fmt.Sprintf("%d %s %d", k+1, spec.Name, spec.Cost[0])
		if ui.PlacingOn && ui.Placing == k {
			for i := 0; i < len(label); i++ {
				f.Set(slot.X+i, y, Cell{R: rune(label[i]), FG: pal.Bright, BG: c, Bold: true})
			}
		} else {
			putString(f, slot.X, y, label, c, 0, false)
		}
	}
	y := menuTop + 2
	// Every line here must fit FrameW (62) at scale 1 — TestMenuLinesFitFrame.
	if ui.Help {
		putString(f, 0, y, " move: arrows/wasd  place: enter/click  select: click tower", pal.Dim, 0, false)
		y++
		putString(f, 0, y, " u up · x sell · t target · n wave (early=bonus) · p pause · f speed", pal.Dim, 0, false)
	} else {
		putString(f, 0, y, " enter place · u up · x sell · t target · n wave · p pause · q", pal.Dim, 0, false)
	}
	// With help on, both hint rows are taken, so the selected-tower info
	// line can't render (it would land off-frame at the bottom).
	if !ui.Help && ui.Selected >= 0 {
		if t := g.Tower(ui.Selected); t != nil {
			y++
			info := towerInfo(t, g.UpgradeCost(t), int(float64(t.Invested)*game.SellRefund))
			putString(f, 0, y, info, pal.Tower[t.Kind], 0, false)
		}
	}
}

func RenderIntro(m *game.Map, name string, diff game.Difficulty, pal Colors, scale int) *Frame {
	l := ComputeLayout(m.W, m.H, scale)
	f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
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
	putString(f, 0, 0, " tdef — terminal tower defense", pal.Bright, 0, true)
	putString(f, 0, 1, fmt.Sprintf(" map: %s   difficulty: %s", name, diffName(diff)), pal.Dim, 0, false)
	// Five stacked rows, all inside the frame: the four help lines occupy
	// H-5..H-2 (one row above the menu zone at scale 1), the prompt gets the
	// last row. Starting at MenuTop() instead pushed the final line off the
	// frame and let the prompt overwrite line 3.
	lines := []string{
		" Enemies walk the path. Build towers on grass to stop them.",
		" 1-7 pick tower · enter/click place · u upgrade · x sell",
		" n start wave early for bonus · p pause · f speed · q quit",
		" Don't let them reach E. Survive all 20 waves.",
	}
	y := l.H - 5
	for _, line := range lines {
		putString(f, 2, y, line, pal.Bright, 0, false)
		y++
	}
	prompt := " press any key to start"
	putString(f, (l.W-len(prompt))/2, l.H-1, prompt, 220, 0, true)
	return f
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

func drawGameOver(f *Frame, g *game.State, ui *UI, pal Colors) {
	won := g.Status == game.StatusVictory
	title, tc := " VICTORY ", 46
	if !won {
		title, tc = " DEFEAT ", 196
	}
	const bg = 236
	bw, bh := 40, 10
	bx, by := (f.W-bw)/2, (f.H-bh)/2
	for y := by; y < by+bh; y++ {
		for x := bx; x < bx+bw; x++ {
			f.Set(x, y, Cell{R: ' ', BG: bg})
		}
	}
	for x := bx; x < bx+bw; x++ {
		f.Set(x, by, Cell{R: '─', FG: pal.Bright, BG: bg})
		f.Set(x, by+bh-1, Cell{R: '─', FG: pal.Bright, BG: bg})
	}
	for y := by; y < by+bh; y++ {
		f.Set(bx, y, Cell{R: '│', FG: pal.Bright, BG: bg})
		f.Set(bx+bw-1, y, Cell{R: '│', FG: pal.Bright, BG: bg})
	}
	f.Set(bx, by, Cell{R: '╔', FG: pal.Bright, BG: bg})
	f.Set(bx+bw-1, by, Cell{R: '╗', FG: pal.Bright, BG: bg})
	f.Set(bx, by+bh-1, Cell{R: '╚', FG: pal.Bright, BG: bg})
	f.Set(bx+bw-1, by+bh-1, Cell{R: '╝', FG: pal.Bright, BG: bg})

	cy := by + 1
	putString(f, (f.W-len(title))/2, cy, title, tc, bg, true)
	cy++
	type kv struct{ k, v string }
	stats := []kv{
		{"wave", fmt.Sprintf("%d/%d", g.Wave, game.MaxWaves)},
		{"kills", fmt.Sprintf("%d", g.TotalKills)},
		{"leaks", fmt.Sprintf("%d", g.TotalLeaks)},
		{"towers", fmt.Sprintf("%d", len(g.Towers))},
		{"score", fmt.Sprintf("%d", g.Score)},
		{"combo", "x" + fmt.Sprintf("%d", g.MaxCombo)},
	}
	for i := 0; i < len(stats); i += 2 {
		line := fmt.Sprintf("%-7s %-9s %-7s %s", stats[i].k, stats[i].v, stats[i+1].k, stats[i+1].v)
		putString(f, (f.W-len(line))/2, cy, line, pal.Dim, bg, false)
		cy++
	}
	putString(f, (f.W-len("time "+formatTime(g.Time)))/2, cy, "time "+formatTime(g.Time), pal.Dim, bg, false)
	cy++
	bestLine := fmt.Sprintf("best %d", ui.BestScore)
	if ui.NewBest {
		bestLine = fmt.Sprintf("★ NEW BEST %d ★", ui.BestScore)
	}
	putString(f, (f.W-len(bestLine))/2, cy, bestLine, 220, bg, ui.NewBest)
	cy++
	putString(f, (f.W-22)/2, cy, "r restart   q quit", pal.Bright, bg, false)
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
