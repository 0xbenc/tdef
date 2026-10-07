package render

import (
	"fmt"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/internal/copytext"
)

func specialistInfo(t *game.Tower, cost, refund int) string {
	switch t.Kind {
	case game.TowerRuneforge:
		return copytext.Format("ui.specialists.forge_info", "tower", t.Spec().Name, "level", fmt.Sprint(t.Level), "damage", fmt.Sprint(t.Dmg()), "range", fmt.Sprint(t.Range()))
	case game.TowerHookmaster:
		return copytext.Format("ui.specialists.hook_info", "tower", t.Spec().Name, "level", fmt.Sprint(t.Level), "pull", fmt.Sprint(t.Spec().Pull[t.Level-1]))
	case game.TowerSappers:
		return copytext.Format("ui.specialists.sapper_info", "tower", t.Spec().Name, "level", fmt.Sprint(t.Level), "damage", fmt.Sprint(t.Dmg()), "cap", fmt.Sprint(t.Spec().MineCap[t.Level-1]))
	default:
		return copytext.Format("ui.specialists.witch_info", "tower", t.Spec().Name, "level", fmt.Sprint(t.Level), "bonus", fmt.Sprintf("%.0f", t.Spec().Mark[t.Level-1]*100))
	}
}

func drawForgeRay(f *Frame, l Layout, start, end game.Pos, color int, preview bool) {
	ax, ay, bx, by := l.FX(start.X), l.FY(start.Y), l.FX(end.X), l.FY(end.Y)
	steps := max(abs(bx-ax), abs(by-ay))
	for i := 1; i <= steps; i++ {
		x, y := ax+(bx-ax)*i/max(1, steps), ay+(by-ay)*i/max(1, steps)
		r := '═'
		if ax == bx {
			r = '║'
		}
		if preview {
			r = '·'
		}
		putOverBold(f, x, y, r, color, !preview)
	}
}

func drawSpecialistField(f *Frame, g *game.State, ui *UI, pal Colors, l Layout, frame int) {
	for _, t := range g.Towers {
		if t.Kind == game.TowerRuneforge && (t.Active || (ui.Selected == t.ID || (!ui.PlacingOn && ui.Cursor == t.Cell))) {
			color := 30
			if t.Active {
				color = 117
			}
			drawForgeRay(f, l, t.Pos(), t.BeamEnd, color, !t.Active)
			// Only the mouth changes; the long line remains steady and readable.
			if t.Active {
				d := t.Facing.Vector()
				p := game.Pos{X: t.Pos().X + float64(d.X)*.65, Y: t.Pos().Y + float64(d.Y)*.65}
				putOverBold(f, l.FX(p.X), l.FY(p.Y), '✦', []int{81, 117, 159}[frame/4%3], true)
			}
		}
		if t.Kind == game.TowerSappers && (ui.Selected == t.ID || (!ui.PlacingOn && ui.Cursor == t.Cell)) {
			for _, site := range g.MineSites(t) {
				putOver(f, l.FX(site.Pos.X), l.FY(site.Pos.Y), '·', 214)
			}
		}
	}
	if ui.PlacingOn && ui.Placing == game.TowerSappers {
		preview := &game.Tower{Kind: game.TowerSappers, Level: 1, Cell: ui.Cursor}
		for _, site := range g.MineSites(preview) {
			putOver(f, l.FX(site.Pos.X), l.FY(site.Pos.Y), '.', 214)
		}
	}
	for _, mine := range g.Mines {
		r, color := '✳', 214
		if mine.ArmAt > g.Time {
			r, color = '○', 137
		}
		putOverBold(f, l.FX(mine.Pos.X), l.FY(mine.Pos.Y), r, color, true)
	}
	if ui.PlacingOn && ui.Placing == game.TowerRuneforge {
		start := ui.Cursor.Center()
		end := game.ForgeEnd(g.Map, ui.Cursor, ui.Facing, game.TowerSpecs[game.TowerRuneforge].Range[0])
		drawForgeRay(f, l, start, end, 81, true)
		putOverBold(f, l.FX(end.X), l.FY(end.Y), ui.Facing.Arrow(), 117, true)
	}
}

func drawSpecialistBeam(f *Frame, l Layout, beam *game.Beam, pal Colors) {
	for i := 1; i < len(beam.From); i++ {
		if beam.Kind != game.TowerHookmaster && beam.Kind != game.TowerWitch {
			drawBeamSeg(f, l, beam.From[i-1], beam.From[i], pal.Beam[beam.Kind])
			continue
		}
		a, b := beam.From[i-1], beam.From[i]
		ax, ay, bx, by := l.FX(a.X), l.FY(a.Y), l.FX(b.X), l.FY(b.Y)
		n := max(abs(bx-ax), abs(by-ay))
		for j := 1; j <= n; j++ {
			r := '╌'
			if beam.Kind == game.TowerWitch {
				r = '·'
			} else if abs(by-ay) > abs(bx-ax) {
				r = '╎'
			}
			x, y := ax+(bx-ax)*j/max(1, n), ay+(by-ay)*j/max(1, n)
			putOverBold(f, x, y, r, pal.Beam[beam.Kind], true)
		}
		if beam.Kind == game.TowerHookmaster {
			putOverBold(f, bx, by, 'ʒ', 223, true)
		}
	}
}
