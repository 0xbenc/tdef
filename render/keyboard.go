package render

import (
	"fmt"
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
	"strings"
)

func focusedDefender(g *game.State, ui *UI) *game.Tower {
	if ui.PlacingOn {
		return nil
	}
	if t := g.Tower(ui.Selected); t != nil {
		return t
	}
	return g.TowerAt(ui.Cursor)
}

func drawPlacementSites(f *Frame, g *game.State, ui *UI, l Layout) {
	if !ui.PlacingOn || !ui.Placing.Valid() {
		return
	}
	mark := Cell{R: '·', FG: 46, BG: 22, Bold: true}
	if g.Gold < game.TowerSpecs[ui.Placing].Cost[0] {
		mark.FG, mark.BG = 214, 58
	}
	for y := 0; y < g.Map.H; y++ {
		for x := 0; x < g.Map.W; x++ {
			if !g.CanBuildSite(game.Vec{X: x, Y: y}, ui.Placing) {
				continue
			}
			cx, cy := l.center(x, y)
			f.Set(cx, cy, mark)
		}
	}
}

func drawKeyboardCursor(f *Frame, g *game.State, ui *UI, l Layout, frame int) {
	if g.Status != game.StatusRunning || !g.Map.InBounds(ui.Cursor) {
		return
	}
	cx, cy := l.center(ui.Cursor.X, ui.Cursor.Y)
	bg := 24
	if frame/15%2 != 0 {
		bg = 30
	}
	if ui.PlacingOn && !g.CanBuild(ui.Cursor, ui.Placing) {
		bg = 88
		if g.CanBuildSite(ui.Cursor, ui.Placing) {
			bg = 58
		}
	}
	for y := l.Y(ui.Cursor.Y); y < l.Y(ui.Cursor.Y)+l.Scale; y++ {
		for x := l.X(ui.Cursor.X); x < l.X(ui.Cursor.X)+l.Scale; x++ {
			if x < 1 || x >= f.W-1 || y < ChromeTop || y >= f.H-ChromeBot {
				continue
			}
			c := f.C[y*f.W+x]
			c.BG = bg
			c.Bold = true
			f.Set(x, y, c)
		}
	}
	if cx >= 1 && cx < f.W-1 && cy >= ChromeTop && cy < f.H-ChromeBot {
		c := f.C[cy*f.W+cx]
		c.FG, c.BG, c.Bold = 231, bg, true
		if ui.PlacingOn {
			c.R = game.TowerSpecs[ui.Placing].Short
			if ui.Placing == game.TowerRuneforge {
				c.R = ui.Facing.Arrow()
			}
		} else if t := g.TowerAt(ui.Cursor); t != nil {
			c.R = t.Spec().Short
			if t.Kind == game.TowerRuneforge {
				c.R = t.Facing.Arrow()
			}
		} else if g.Map.At(ui.Cursor) == game.CellWall || isGroundRune(c.R) {
			c.R = '◻'
		}
		f.Set(cx, cy, c)
		f.Set(0, cy, Cell{R: '▶', FG: 159, Bold: true})
		f.Set(f.W-1, cy, Cell{R: '◀', FG: 159, Bold: true})
	}
}

// Target explanations are shared by the action feedback and the help panel.
func TargetExplanation(mode game.TargetMode) string {
	switch mode {
	case game.TargetStrongest:
		return copytext.Text("ui.keyboard.strongest")
	case game.TargetClosest:
		return copytext.Text("ui.keyboard.closest")
	default:
		return copytext.Text("ui.keyboard.first")
	}
}

func drawTargetGuide(f *Frame, g *game.State, ui *UI, l Layout) {
	t := focusedDefender(g, ui)
	if t == nil {
		return
	}
	e := g.TargetFor(t)
	if e == nil {
		return
	}
	ax, ay, bx, by := l.FX(t.Pos().X), l.FY(t.Pos().Y), l.FX(e.Pos.X), l.FY(e.Pos.Y)
	n := max(abs(bx-ax), abs(by-ay))
	for i := 1; i < n; i++ {
		x, y := ax+(bx-ax)*i/max(1, n), ay+(by-ay)*i/max(1, n)
		if x >= 1 && x < f.W-1 && y >= ChromeTop && y < f.H-ChromeBot && isGroundRune(f.C[y*f.W+x].R) {
			putOver(f, x, y, '·', 159)
		}
	}
}

func drawDefenderActions(f *Frame, g *game.State, ui *UI) {
	if ui.Help || g.Status != game.StatusRunning {
		return
	}
	t := focusedDefender(g, ui)
	if t == nil && !ui.PlacingOn {
		return
	}
	y := f.H - 2
	for x := 1; x < f.W-1; x++ {
		f.Set(x, y, Cell{R: ' '})
	}
	if ui.PlacingOn {
		y = f.H - 1
		hint := copytext.Text("ui.keyboard.place")
		if ui.Placing == game.TowerRuneforge {
			hint += copytext.Text("ui.keyboard.rotate")
		}
		hint = fitMsg(hint, f.W-4)
		x := (f.W - len([]rune(hint))) / 2
		f.Set(x-1, y, Cell{R: '┘', FG: 240})
		for len(hint) > 0 {
			colored := false
			for _, token := range []struct {
				text  string
				color int
			}{{"Enter", 46}, {"Esc", 174}, {"green", 46}, {"amber", 214}, {"r aim", 159}} {
				if strings.HasPrefix(hint, token.text) {
					putString(f, x, y, token.text, token.color, 0, true)
					x += len([]rune(token.text))
					hint = strings.TrimPrefix(hint, token.text)
					colored = true
					break
				}
			}
			if !colored {
				r := []rune(hint)
				f.Set(x, y, Cell{R: r[0], FG: 252})
				x++
				hint = string(r[1:])
			}
		}
		f.Set(x, y, Cell{R: '└', FG: 240})
		return
	}
	x := 2
	action := func(key, text string, color int) {
		if x+len([]rune(key))+len([]rune(text)) > f.W-2 {
			return
		}
		putString(f, x, y, key, color, 0, true)
		x += len([]rune(key))
		putString(f, x, y, text, 252, 0, false)
		x += len([]rune(text))
	}
	cost := g.UpgradeCost(t)
	if cost == 0 {
		action("u", copytext.Text("ui.keyboard.max"), 252)
	} else {
		color := 46
		label := copytext.Format("ui.keyboard.upgrade", "gold", fmt.Sprint(cost))
		if g.Gold < cost {
			color = 174
			label = copytext.Format("ui.keyboard.need", "gold", fmt.Sprint(cost))
		}
		action("u", label, color)
	}
	switch t.Kind {
	case game.TowerRuneforge:
		action("r", copytext.Format("ui.keyboard.aim", "direction", string(t.Facing.Arrow())), 159)
	case game.TowerSappers:
		count := 0
		for _, mine := range g.Mines {
			if mine.Owner == t.ID {
				count++
			}
		}
		action("", copytext.Format("ui.keyboard.stock", "count", fmt.Sprintf("%d/%d", count, t.Spec().MineCap[t.Level-1])), 214)
	default:
		action("t", copytext.Format("ui.keyboard.target", "mode", t.TargetMode.Name()), 159)
	}
	refund := int(math.Floor(float64(t.Invested) * game.SellRefund))
	action("x", copytext.Format("ui.keyboard.sell", "gold", fmt.Sprint(refund)), 174)
	if x+10 < f.W {
		action("Tab", copytext.Text("ui.keyboard.next"), 159)
	}
	// Keep the colored controls together, centered above the framed stats.
	width := x - 2
	cells := append([]Cell(nil), f.C[y*f.W+2:y*f.W+x]...)
	for col := 1; col < f.W-1; col++ {
		f.Set(col, y, Cell{R: ' '})
	}
	start := (f.W - width) / 2
	for i, c := range cells {
		f.Set(start+i, y, c)
	}
}
