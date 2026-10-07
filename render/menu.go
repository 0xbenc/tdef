package render

import (
	"strings"

	"github.com/0xbenc/termtd/internal/buildinfo"
	"github.com/0xbenc/termtd/internal/copytext"
)

// Leave a full-width target beside the Rotunda, even at 62 columns.
func menuColumns(w int) (left, x, width int) {
	width = 23
	left = max(12, min(80, w-32))
	base := max(2, (w-left-width-7)/2)
	return left, base + left + 4, width
}

// The main menu looks into the same vault as the lair, with its own quiet clock.
func RenderMenuAnimated(w, h, sel, frame int, revealed bool, pal Colors) *Frame {
	f := screenBox(w, h, copytext.Text("ui.render_menu.main_menu"), []fseg{
		{key: "↑↓", text: copytext.Text("ui.render_menu.move")},
		{key: "enter", text: copytext.Text("ui.render_menu.select")},
		{key: "q", text: copytext.Text("ui.render_menu.quit")},
	}, true, pal)
	w, h = f.W, f.H
	label := buildinfo.Version
	if label != "dev" {
		label = "v" + strings.TrimPrefix(label, "v")
	}
	if len([]rune(label)) <= w-22 {
		embedSegment(f, 0, w-len([]rune(label))-4, label, '┐', '┌', pal.Path, 244, false)
	}
	left, mx, mw := menuColumns(w)
	x := mx - left - 4
	for y := 1; y < h-1; y++ {
		for xx := 1; xx < w-1; xx++ {
			f.Set(xx, y, Cell{R: ' ', FG: 244, BG: 233})
		}
	}
	sel = max(0, min(sel, len(MenuItems)-1))
	top := max(2, (h-25)/2)
	artY := top + 6
	if left >= 44 && h >= 30 {
		logo := &Frame{W: left, H: 9, C: make([]Cell, left*9)}
		drawTitleLogo(logo, left, 0, -1)
		for yy := 0; yy < logo.H; yy++ {
			for xx := 0; xx < left; xx++ {
				c := logo.C[yy*left+xx]
				if c.R != 0 {
					c.BG = 233
					f.Set(x+xx, top-1+yy, c)
				}
			}
		}
		artY = top + 11
	} else {
		putString(f, x+1, top, "T E R M  T D", 180, 233, true)
	}
	artH := max(4, min(h-artY-2, left*7/15))
	artW := min(left, artH*15/7)
	drawRotundaTableau(f, x+(left-artW)/2, artY, artW, artH, max(1, artH/7), 233,
		owPadView{status: OWOpen}, frame, true, revealed)
	if left-artW >= 8 && artH >= 7 {
		for i, bx := range []int{x + 1, x + left - 4} {
			drawMenuBrazier(f, bx, artY+artH-2, frame+i*17, revealed)
		}
	}
	for y := max(1, menuLayout(h)[0]-1); y <= min(h-2, menuLayout(h)[7]+1); y++ {
		f.Set(mx-2, y, Cell{R: '│', FG: 237, BG: 233})
	}
	for i, r := range MenuRects(w, h) {
		bg, fg := 233, 245
		if i == sel {
			bg, fg = 52, 230
		}
		for xx := r.X; xx < r.X+mw; xx++ {
			f.Set(xx, r.Y, Cell{R: ' ', FG: fg, BG: bg})
		}
		marker := "   "
		if i == sel {
			marker = " ▸ "
		}
		putString(f, r.X, r.Y, marker+MenuItems[i], fg, bg, i == sel)
		if i == sel {
			f.Set(r.X+1, r.Y, Cell{R: '▸', FG: 167, Bold: true})
		}
		f.Set(r.X+mw-2, r.Y, Cell{R: rune('1' + i), FG: 167, BG: bg, Bold: true})
	}
	// Respect the terminal's background, including the selection and artwork.
	for i := range f.C {
		f.C[i].BG = 0
	}
	return f
}

func drawMenuBrazier(f *Frame, x, floor, frame int, lit bool) {
	if !lit {
		putString(f, x, floor-2, "╲━╱", 240, 233, false)
		f.Set(x+1, floor-1, Cell{R: '┃', FG: 238, BG: 233})
		putString(f, x, floor, "━┻━", 238, 233, false)
		return
	}
	// Fixed geometry and a low flame keep the resting dragon in focus.
	phase := (frame / 7) % 4
	for y := floor - 4; y <= floor; y++ {
		for xx := x; xx < x+3; xx++ {
			f.Set(xx, y, Cell{R: ' ', FG: 130, BG: 234})
		}
	}
	f.Set(x+1, floor-4, Cell{R: []rune{'╵', '╱', '╵', '╲'}[phase], FG: 208, BG: 234})
	putString(f, x, floor-3, "▟█▙", 166, 234, false)
	f.Set(x+1, floor-3, Cell{R: '▓', FG: 220, BG: 166, Bold: true})
	putString(f, x, floor-2, "╲━╱", 137, 234, false)
	f.Set(x+1, floor-1, Cell{R: '┃', FG: 238, BG: 234})
	putString(f, x, floor, "━┻━", 238, 234, false)
	if phase != 0 {
		f.Set(x+phase%3, floor-6, Cell{R: '·', FG: 130, BG: 233})
	}
}
