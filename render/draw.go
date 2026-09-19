package render

import (
	"fmt"

	"tdef/game"
)

type UI struct {
	Cursor    game.Vec
	Placing   game.TowerKind
	PlacingOn bool
	Selected  int
	Speed     int
	Paused    bool
	Help      bool
	Message   string

	BestScore int
	NewBest   bool
}

const (
	HUDRows = 2
	MapH    = 13
	FrameW  = 62
	MenuTop = HUDRows + MapH
	FrameH  = MenuTop + 4
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
		Tower: [game.TowerCount]int{46, 203, 51, 171, 220},
		Enemy: [game.EnemyCount]int{213, 114, 180, 204},
		Beam:  [game.TowerCount]int{255, 203, 51, 171, 226},
	}
}

type bounds struct{ ox, oy int }

func (b bounds) x(mx int) int { return b.ox + mx }
func (b bounds) y(my int) int { return b.oy + my }

func Render(g *game.State, ui *UI, pal Colors) *Frame {
	f := &Frame{W: FrameW, H: FrameH, C: make([]Cell, FrameW*FrameH)}
	b := bounds{ox: (FrameW - g.Map.W) / 2, oy: HUDRows}
	for y := 0; y < g.Map.H; y++ {
		for x := 0; x < g.Map.W; x++ {
			v := game.Vec{X: x, Y: y}
			switch g.Map.At(v) {
			case game.CellWall:
				f.Put(b.x(x), b.y(y), ' ', 0, pal.Wall)
			case game.CellPath:
				f.Put(b.x(x), b.y(y), '·', pal.Path, 0)
			case game.CellGrass:
				f.Put(b.x(x), b.y(y), ' ', pal.Grass, 0)
			}
		}
	}
	f.C[b.y(g.Map.Spawn.Y)*f.W+b.x(g.Map.Spawn.X)] = Cell{R: '▶', FG: 46, Bold: true}
	f.C[b.y(g.Map.Exit.Y)*f.W+b.x(g.Map.Exit.X)] = Cell{R: 'E', FG: 196, Bold: true}
	drawRange(f, g, ui, pal, b)
	for _, t := range g.Towers {
		x, y := b.x(t.Cell.X), b.y(t.Cell.Y)
		c := pal.Tower[t.Kind]
		f.Put(x, y, t.Spec().Short, c, 0)
		if ui.Selected == t.ID {
			f.C[y*f.W+x] = Cell{R: t.Spec().Short, FG: pal.Bright, BG: c, Bold: true}
		}
		for i := 1; i < t.Level; i++ {
			f.Put(x-i, y, '▪', c, 0)
		}
	}
	for _, p := range g.Projectiles {
		f.Put(b.x(int(p.Pos.X)), b.y(int(p.Pos.Y)), '+', pal.Beam[p.Kind], 0)
	}
	for _, bm := range g.Beams {
		for i := 1; i < len(bm.From); i++ {
			drawLine(f, bm.From[i-1], bm.From[i], b, pal.Beam[bm.Kind])
		}
	}
	for _, e := range g.Enemies {
		x, y := b.x(int(e.Pos.X)), b.y(int(e.Pos.Y))
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
		f.C[y*f.W+x] = Cell{R: game.EnemySpecs[e.Kind].Short, FG: fg, Bold: e.Kind == game.EnemyBoss}
	}
	for _, fx := range g.Fx {
		x, y := b.x(int(fx.Pos.X)), b.y(int(fx.Pos.Y))
		if x >= 0 && y >= 0 && x < f.W && y < f.H {
			f.Put(x, y, fx.R, fx.Color, 0)
		}
	}
	if g.LeakFlash > 0 {
		ex, ey := b.x(g.Map.Exit.X), b.y(g.Map.Exit.Y)
		f.C[ey*f.W+ex] = Cell{R: 'E', FG: 231, BG: 196, Bold: true}
	}
	if ui.PlacingOn {
		x, y := b.x(ui.Cursor.X), b.y(ui.Cursor.Y)
		c := 196
		if g.CanBuild(ui.Cursor, ui.Placing) {
			c = 46
		}
		f.C[y*f.W+x] = Cell{R: game.TowerSpecs[ui.Placing].Short, FG: c, Bold: true}
	} else if ui.Selected < 0 {
		x, y := b.x(ui.Cursor.X), b.y(ui.Cursor.Y)
		if x >= 0 && y >= 0 && x < f.W && y < f.H {
			cc := f.C[y*f.W+x]
			if cc.R == 0 || cc.R == ' ' {
				cc.R = '◻'
				cc.FG = pal.Dim
			} else {
				cc.Bold = true
			}
			f.C[y*f.W+x] = cc
		}
	}
	drawHUD(f, g, ui, pal)
	drawMenu(f, g, ui, pal)
	if g.Status != game.StatusRunning {
		drawGameOver(f, g, ui, pal)
	}
	return f
}

func drawRange(f *Frame, g *game.State, ui *UI, pal Colors, b bounds) {
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
				f.Put(b.x(x), b.y(y), '·', pal.Dim, 0)
			}
		}
	}
}

func drawLine(f *Frame, a, bp game.Pos, b bounds, color int) {
	steps := int(a.Dist(bp)*8) + 1
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := b.x(int(a.X + (bp.X-a.X)*t))
		y := b.y(int(a.Y + (bp.Y-a.Y)*t))
		if x < 0 || y < 0 || x >= f.W || y >= f.H {
			continue
		}
		if f.C[y*f.W+x].R != ' ' && f.C[y*f.W+x].R != 0 {
			continue
		}
		f.Put(x, y, '·', color, 0)
	}
}

func drawHUD(f *Frame, g *game.State, ui *UI, pal Colors) {
	pause := ""
	if ui.Paused {
		pause = " ⏸ "
	}
	msg := ""
	if ui.Message != "" {
		msg = "  " + ui.Message
	}
	combo := ""
	if g.Combo >= 5 {
		combo = fmt.Sprintf("  ⚡%d", g.Combo)
	}
	line0 := " tdef "
	line1 := ""
	switch {
	case g.Status == game.StatusRunning && g.WaveActive:
		line0 += fmt.Sprintf("Wave %d/%d%s", g.Wave, game.MaxWaves, combo)
	case g.Status == game.StatusRunning && g.Wave < game.MaxWaves:
		nw := g.Wave + 1
		in := int(g.NextWaveAt-g.Time) + 1
		if in < 0 {
			in = 0
		}
		bonus := game.EarlyBonusBase + g.Wave
		line0 += fmt.Sprintf("next wave %d: %s   in %ds  [n +%dg]", nw, game.WavePreview(nw), in, bonus)
		line1 += ""
	}
	line0 += pause + msg
	line1 += fmt.Sprintf(" ⛁ %d   ♥ %d   ★ %d   x%d", g.Gold, g.Lives, g.Score, ui.Speed)
	putString(f, 0, 0, line0, pal.Bright, 0, true)
	putString(f, 0, 1, line1, pal.Gold, 0, false)
}

type MenuSlot struct {
	Kind game.TowerKind
	X    int
	Y    int
}

var MenuSlots = []MenuSlot{
	{game.TowerGunner, 1, 0}, {game.TowerCannon, 17, 0}, {game.TowerFrost, 35, 0},
	{game.TowerSniper, 1, 1}, {game.TowerTesla, 19, 1},
}

func drawMenu(f *Frame, g *game.State, ui *UI, pal Colors) {
	for _, slot := range MenuSlots {
		k := slot.Kind
		spec := game.TowerSpecs[k]
		afford := g.Gold >= spec.Cost[0]
		c := pal.Tower[k]
		if !afford {
			c = pal.Dim
		}
		y := MenuTop + slot.Y
		label := fmt.Sprintf("%d %s %d", k+1, spec.Name, spec.Cost[0])
		if ui.PlacingOn && ui.Placing == k {
			for i := 0; i < len(label); i++ {
				f.Set(slot.X+i, y, Cell{R: rune(label[i]), FG: pal.Bright, BG: c, Bold: true})
			}
		} else {
			putString(f, slot.X, y, label, c, 0, false)
		}
	}
	y := MenuTop + 2
	if ui.Help {
		putString(f, 0, y, " move: arrows/wasd  place: enter/click  select: click a tower  upgrade: u  sell: x", pal.Dim, 0, false)
		y++
		putString(f, 0, y, " start wave: n (early = bonus)  pause: p  speed: f/wheel  cancel: esc  quit: q", pal.Dim, 0, false)
	} else {
		putString(f, 0, y, " enter place · u up · x sell · t target · n wave · p pause · f speed · esc · q", pal.Dim, 0, false)
	}
	if ui.Selected >= 0 {
		if t := g.Tower(ui.Selected); t != nil {
			y++
			var info string
			if t.Level >= 3 {
				info = fmt.Sprintf(" ▸ %s Lv%d (max)  target %s [t]  sell +%d", t.Spec().Name, t.Level, t.TargetMode.Name(), int(float64(t.Invested)*game.SellRefund))
			} else {
				info = fmt.Sprintf(" ▸ %s Lv%d  dmg %.0f  rng %.1f  target %s [t]  upgrade %d  sell +%d", t.Spec().Name, t.Level, t.Dmg(), t.Range(), t.TargetMode.Name(), g.UpgradeCost(t), int(float64(t.Invested)*game.SellRefund))
			}
			putString(f, 0, y, info, pal.Tower[t.Kind], 0, false)
		}
	}
}

func RenderIntro(m *game.Map, name string, diff game.Difficulty, pal Colors) *Frame {
	f := &Frame{W: FrameW, H: FrameH, C: make([]Cell, FrameW*FrameH)}
	b := bounds{ox: (FrameW - m.W) / 2, oy: HUDRows}
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			switch m.At(game.Vec{X: x, Y: y}) {
			case game.CellWall:
				f.Put(b.x(x), b.y(y), ' ', 0, pal.Wall)
			case game.CellPath:
				f.Put(b.x(x), b.y(y), '·', pal.Path, 0)
			case game.CellGrass:
				f.Put(b.x(x), b.y(y), ' ', pal.Grass, 0)
			}
		}
	}
	f.C[b.y(m.Spawn.Y)*f.W+b.x(m.Spawn.X)] = Cell{R: '▶', FG: 46, Bold: true}
	f.C[b.y(m.Exit.Y)*f.W+b.x(m.Exit.X)] = Cell{R: 'E', FG: 196, Bold: true}
	putString(f, 0, 0, " tdef — terminal tower defense", pal.Bright, 0, true)
	putString(f, 0, 1, fmt.Sprintf(" map: %s   difficulty: %s", name, diffName(diff)), pal.Dim, 0, false)
	lines := []string{
		" Enemies walk the path. Build towers on grass to stop them.",
		" 1-5 pick tower · enter/click place · u upgrade · x sell",
		" n start wave early (bonus gold) · p pause · f speed · q quit",
		"",
		" Don't let them reach E. Survive all 20 waves.",
	}
	y := MenuTop
	for _, l := range lines {
		putString(f, 2, y, l, pal.Bright, 0, false)
		y++
	}
	putString(f, (FrameW-26)/2, FrameH-2, " press any key to start", 220, 0, true)
	return f
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
	title, c := " VICTORY ", 46
	if !won {
		title, c = " DEFEAT ", 196
	}
	cx := (f.W - len(title)) / 2
	putString(f, cx, f.H/2-2, title, pal.Bright, c, true)
	line := fmt.Sprintf(" wave %d/%d   kills %d   score %d", g.Wave, game.MaxWaves, g.TotalKills, g.Score)
	putString(f, (f.W-len(line))/2, f.H/2, line, pal.Dim, 0, false)
	bestLine := fmt.Sprintf(" best %d", ui.BestScore)
	if ui.NewBest {
		bestLine = fmt.Sprintf(" ★ NEW BEST %d ★", ui.BestScore)
	}
	putString(f, (f.W-len(bestLine))/2, f.H/2+2, bestLine, 220, 0, ui.NewBest)
	putString(f, (f.W-24)/2, f.H/2+4, " r restart    q quit", pal.Bright, 0, false)
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
