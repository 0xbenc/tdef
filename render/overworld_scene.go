package render

import (
	"github.com/0xbenc/tdef/game"
	"math"
)

// OverworldLayout keeps the battlefield's scale steps, but fits a viewport
// rather than shrinking the entire lair. Camera offsets are in terminal
// columns, so scrolling at 2x+ can move by less than a whole world cell.
func OverworldLayout(w, h int, st OWState) Layout {
	l := GameLayout(45, OWH, w, h)
	span := OWW * l.Scale
	available := w - 2
	if span <= available {
		l.Ox = 1 + (available-span)/2
		return l
	}
	focus := float64(st.Cursor.X)
	if st.CameraSet {
		focus = st.CameraX
	}
	offset := int(math.Round((focus+0.5)*float64(l.Scale))) - available/2
	// Even after a browser jump or a resize, never leave Grak outside the
	// viewport while the camera is catching up.
	player := st.Cursor.X*l.Scale + l.Scale/2
	margin := 3 * l.Scale
	if player-offset < margin {
		offset = player - margin
	}
	if player-offset > available-margin-1 {
		offset = player - (available - margin - 1)
	}
	if offset < 0 {
		offset = 0
	}
	if offset > span-available {
		offset = span - available
	}
	l.Ox = 1 - offset
	return l
}

func owSweepRadius() float64 {
	n := owNodeByID("rotunda")
	dx := math.Max(float64(n.X), float64(OWW-1-n.X))
	dy := math.Max(float64(n.Y), float64(OWH-1-n.Y))
	return math.Hypot(dx, dy) + 6
}

// The background is built from large silhouettes in terminal space. Surface
// grain stays secondary to the cave lip, pillars, bridge and underground river.
func drawOWCavern(f *Frame, l Layout, frame int) {
	for y := 0; y < OWH*l.Scale; y++ {
		for x := 0; x < OWW*l.Scale; x++ {
			bg := 233
			if y > 9*l.Scale {
				bg = 232
			}
			f.Put(l.Ox+x, l.Oy+y, ' ', 236, bg)
		}
	}
	// Broad, irregular rock masses. Sloping edges describe stone instead
	// of repeating the same texture glyph over every expanded map cell.
	for _, rock := range []struct {
		x, width, height int
		ceiling          bool
	}{
		{0, 6, 3, true}, {15, 9, 2, true}, {30, 10, 3, true},
		{46, 11, 2, true}, {71, 7, 4, true},
		{0, 5, 3, false}, {12, 8, 2, false}, {25, 7, 2, false}, {72, 6, 3, false},
	} {
		width, height := rock.width*l.Scale, rock.height*l.Scale
		for dx := 0; dx < width; dx++ {
			depth := height - abs(dx-width/2)*height/(width/2+1)
			for dy := 0; dy < depth; dy++ {
				y := l.Oy + dy
				if !rock.ceiling {
					y = l.Y(OWH) - 1 - dy
				}
				c := Cell{R: ' ', FG: 237, BG: 235}
				if dy == depth-1 {
					c.R = '╲'
					if (dx < width/2) != rock.ceiling {
						c.R = '╱'
					}
				} else if dy > 0 && (dx+dy*3)%11 == 0 {
					c.R, c.FG = '─', 236
				}
				f.Set(l.X(rock.x)+dx, y, c)
			}
		}
	}
	// Far-wall ruins: large dark arches, partly lost behind the crossing.
	for _, x := range []int{33, 40, 47} {
		x0, y0 := l.X(x), l.Y(1)
		width, height := 5*l.Scale, 4*l.Scale
		f.Put(x0, y0+1, '╭', 237, 233)
		f.Put(x0+width-1, y0+1, '╮', 237, 233)
		for dx := 1; dx < width-1; dx++ {
			f.Put(x0+dx, y0, '─', 237, 233)
		}
		for dy := 2; dy < height; dy++ {
			f.Put(x0, y0+dy, '│', 237, 233)
			f.Put(x0+width-1, y0+dy, '│', 237, 233)
		}
	}
	// Black water winds below the paths; bright ripples are sparse enough
	// that it reads as a surface rather than another navigable route.
	for x := l.X(32); x < l.X(71); x++ {
		bank := l.Y(10) + (x/l.Scale/7)%2
		for y := bank; y < l.Y(OWH); y++ {
			c := Cell{R: ' ', BG: 232}
			if y == bank {
				c.R, c.FG = '≈', 24
			} else if (x+y*7+frame/12)%23 < 3 {
				c.R, c.FG = '~', 24
			}
			f.Set(x, y, c)
		}
	}
	// Bridge piers hold up the central span. Their feet disappear in water.
	for _, x := range []int{31, 37, 43} {
		px, py := l.X(x), l.Y(7)
		for y := py; y < l.Y(12); y++ {
			f.Put(px, y, '┃', 239, 233)
		}
		for bx := px + 1; bx < l.X(x+6); bx++ {
			f.Put(bx, py, '━', 240, 233)
		}
		f.Put(px+1, py+1, '╭', 240, 233)
		f.Put(l.X(x+6)-1, py+1, '╮', 240, 233)
	}
	// A waterfall behind the bridge, broken into drops by the draft.
	for y := l.Y(2); y < l.Y(12); y++ {
		for dx := 0; dx < max(1, l.Scale); dx++ {
			if (y-frame/3+dx)%5 == 0 {
				continue
			}
			f.Put(l.X(40)+dx, y, '┊', 24+dx%2, 233)
		}
	}
	// A fallen expedition's bones on the bank make the empty crossing a place.
	stampOW(f, l, 34, 11, []string{"(x)─┼─"}, 240, 232)
	// Torch sconces and the fork's signpost are scenery, off the walkable cells.
	for _, v := range [][2]int{{12, 7}, {32, 4}, {43, 5}, {54, 5}} {
		px, py := l.center(v[0], v[1])
		r := '*'
		if (frame/6+v[0])%3 == 0 {
			r = '+'
		}
		f.Put(px, py, r, 214, 233)
		f.Put(px, py+1, '│', 94, 233)
	}
	stampOW(f, l, 47, 7, []string{"↑╥↓", " ║ "}, 180, 233)
}

// stampOW uses terminal characters for objects, so their outlines stay crisp
// at every scale rather than becoming repeated blocks of the same glyph.
func stampOW(f *Frame, l Layout, mx, my int, art []string, fg, bg int) {
	cx, cy := l.center(mx, my)
	width := 0
	for _, line := range art {
		width = max(width, len([]rune(line)))
	}
	for row, line := range art {
		for i, r := range []rune(line) {
			if r != ' ' {
				f.Put(cx-width/2+i, cy+row, r, fg, bg)
			}
		}
	}
}

func owWalkwayRune(x, y int) rune {
	// Neighbour connectivity comes from the actual corridor set, including
	// shared forks and stair turns; decoration never defines a new route.
	n := owCorridor[game.Vec{X: x, Y: y - 1}]
	e := owCorridor[game.Vec{X: x + 1, Y: y}]
	s := owCorridor[game.Vec{X: x, Y: y + 1}]
	w := owCorridor[game.Vec{X: x - 1, Y: y}]
	switch {
	case n && e && s && w:
		return '╬'
	case n && e && s:
		return '╠'
	case n && s && w:
		return '╣'
	case e && s && w:
		return '╦'
	case n && e && w:
		return '╩'
	case n && e:
		return '╚'
	case n && w:
		return '╝'
	case s && e:
		return '╔'
	case s && w:
		return '╗'
	case n || s:
		return '║'
	default:
		return '═'
	}
}

func drawOWWalkway(f *Frame, l Layout, st OWState) {
	for v := range owCorridor {
		bg, fg := 237, 250
		// A cleared branch has warmer stone; shared branches keep the best
		// state regardless of the order in which the floors are rendered.
		for i, route := range owRouteCells {
			if !st.Records[owRouteFloors[i]].Cleared {
				continue
			}
			for _, cell := range route {
				if cell == v {
					fg = 180
					break
				}
			}
		}
		l.block(f, v.X, v.Y, Cell{R: ' ', FG: fg, BG: bg})
		cx, cy := l.center(v.X, v.Y)
		if l.Scale == 1 {
			f.Put(cx, cy, owWalkwayRune(v.X, v.Y), fg, bg)
			continue
		}
		// Paving seams join across block boundaries. The exposed edges form
		// parapets along horizontal spans and stair rails on vertical ones.
		for dy := 0; dy < l.Scale; dy++ {
			for dx := 0; dx < l.Scale; dx++ {
				x, y := l.X(v.X)+dx, l.Y(v.Y)+dy
				r, color := ' ', fg
				if dy == 0 && !owCorridor[game.Vec{X: v.X, Y: v.Y - 1}] {
					r, color = '▀', 242
				}
				if dy == l.Scale-1 && !owCorridor[game.Vec{X: v.X, Y: v.Y + 1}] {
					r, color = '▄', 242
				}
				if dx == 0 && !owCorridor[game.Vec{X: v.X - 1, Y: v.Y}] && (owCorridor[game.Vec{X: v.X, Y: v.Y - 1}] || owCorridor[game.Vec{X: v.X, Y: v.Y + 1}]) {
					r, color = '▌', 242
				}
				if dx == l.Scale-1 && !owCorridor[game.Vec{X: v.X + 1, Y: v.Y}] && (owCorridor[game.Vec{X: v.X, Y: v.Y - 1}] || owCorridor[game.Vec{X: v.X, Y: v.Y + 1}]) {
					r, color = '▐', 242
				}
				if dy == l.Scale/2 && (owCorridor[game.Vec{X: v.X - 1, Y: v.Y}] || owCorridor[game.Vec{X: v.X + 1, Y: v.Y}]) {
					r, color = '═', fg
				}
				if dx == l.Scale/2 && (owCorridor[game.Vec{X: v.X, Y: v.Y - 1}] || owCorridor[game.Vec{X: v.X, Y: v.Y + 1}]) {
					r, color = '║', fg
				}
				f.Put(x, y, r, color, bg)
			}
		}
		f.Put(cx, cy, owWalkwayRune(v.X, v.Y), fg, bg)
	}
}

// Each room is an illustration with a distinct silhouette rather than the
// same framed pad. These structures occupy the existing walkable footprint;
// their centre remains the floor's interactive landmark.
func drawOWBuilding(f *Frame, l Layout, n *owNode, accent, bg int, v owPadView, frame int) {
	x0, y0 := l.X(n.X-n.PW/2), l.Y(n.Y-n.PH/2)
	x1, y1 := l.X(n.X+n.PW/2)+l.Scale-1, l.Y(n.Y+n.PH/2)+l.Scale-1
	cx, cy := l.center(n.X, n.Y)
	open := v.status != OWSealed || v.unseal > 0
	material := accent
	if open {
		switch n.ID {
		case "halls":
			material = 244
		case "garden":
			material = 143
		case "depths":
			material = 94
		case "rift":
			material = 240
		}
	}
	if v.unseal > 0 {
		material = lerp(dim(material), material, v.unseal)
	}
	if v.flash > 0 {
		material = lerp(material, v.flashC, v.flash)
	}
	put := func(x, y int, r rune) { f.Set(x, y, Cell{R: r, FG: material, BG: bg}) }
	hline := func(a, b, y int, r rune) {
		for x := a; x <= b; x++ {
			put(x, y, r)
		}
	}
	vline := func(x, a, b int, r rune) {
		for y := a; y <= b; y++ {
			put(x, y, r)
		}
	}
	switch n.ID {
	case "halls":
		// Crenellated towers, a vaulted gallery and suspended guild banners.
		hline(x0, x1, y1, '━')
		tower := max(3, 2*l.Scale)
		for _, x := range []int{x0, x1 - tower + 1} {
			for dx := 0; dx < tower; dx++ {
				r := '▄'
				if dx%2 == 0 {
					r = '▟'
				}
				put(x+dx, y0, r)
			}
			for y := y0 + 1; y < y1; y++ {
				for dx := 0; dx < tower; dx++ {
					r := ' '
					if dx == 0 || dx == tower-1 {
						r = '┃'
					}
					f.Put(x+dx, y, r, material, 236)
				}
			}
		}
		hline(x0+tower, x1-tower, y0+1, '─')
		step := max(4, 4*l.Scale)
		for x := x0 + tower; x+step-1 <= x1-tower; x += step {
			put(x, y0+2, '╭')
			hline(x+1, x+step-2, y0+2, '─')
			put(x+step-1, y0+2, '╮')
			vline(x, y0+3, y1-1, '│')
			vline(x+step-1, y0+3, y1-1, '│')
			if l.Scale >= 2 {
				c := accent
				if !open {
					c = dim(c)
				}
				f.Put(x+step/2, y0+3, '▐', c, bg)
				f.Put(x+step/2, y0+4, '▼', c, bg)
			}
		}
	case "garden":
		// A broad temple roof, weathered columns, trailing ivy and lily pads.
		roof := max(2, l.Scale+1)
		for y := 0; y <= roof; y++ {
			half := (x1 - x0 - 4) * y / (2 * roof)
			for x := cx - half; x <= cx+half; x++ {
				f.Put(x, y0+y, ' ', material, 235)
			}
			if y > 0 {
				put(cx-half, y0+y, '╱')
				put(cx+half, y0+y, '╲')
			}
		}
		put(cx, y0, '▲')
		hline(x0+2, x1-2, y0+roof, '═')
		for _, x := range []int{x0 + 3, x0 + 5, x1 - 5, x1 - 3} {
			vline(x, y0+roof+1, y1-1, '│')
			put(x, y1-1, '┴')
		}
		hline(x0+1, x1-1, y1, '≈')
		if open {
			for y := y0 + roof; y < y1; y++ {
				x := x0 + 1 + (y-y0)%2
				f.Put(x, y, 'Y', 65, bg)
				if y%2 == 0 {
					f.Put(x1-1, y, 'Y', 65, bg)
				}
			}
			for _, x := range []int{cx - 4, cx + 4} {
				f.Put(x, y1-1, 'o', 65, 233)
			}
		}
	case "depths":
		// A timbered mine mouth; rails converge into the dark.
		hline(x0+2, x1-2, y0, '━')
		put(x0+1, y0, '╭')
		put(x1-1, y0, '╮')
		vline(x0+1, y0+1, y1, '┃')
		vline(x1-1, y0+1, y1, '┃')
		for y := cy + 1; y <= y1; y++ {
			d := y - cy
			put(cx-d, y, '╱')
			put(cx+d, y, '╲')
			hline(cx-d+1, cx+d-1, y, '─')
		}
		if open && l.Scale >= 2 {
			stampOW(f, l, n.X-3, n.Y+1, []string{"┌──┐", "o──o"}, 137, bg)
			f.Put(x1-2, y0+2, '*', 214, bg)
			f.Put(x1-2, y0+3, '│', 94, bg)
		}
	case "rift":
		// Broken basalt teeth framing a breach, open at the top.
		for y := y0; y <= y1; y++ {
			d := (y - y0) % 3
			put(x0+d, y, '╲')
			put(x1-d, y, '╱')
			if y > y0 {
				put(x0+d+1, y, '▒')
				put(x1-d-1, y, '▒')
			}
		}
		hline(x0+2, x1-2, y1, '▄')
		if open {
			for x := x0 + 3; x < x1-2; x++ {
				fg := 202
				if (x+frame/6)%3 == 0 {
					fg = 214
				}
				f.Put(x, y1, '≈', fg, 52)
			}
			f.Put(cx, cy-1, '*', 214, bg)
		}

	case "rotunda":
		drawOWRotunda(f, l, n, bg, v, frame)

	}
	if n.ID == "rift" || n.ID == "halls" {
		x := x0
		if n.ID == "rift" {
			x = x1
		}
		for dy := -1; dy <= 1; dy++ {
			f.Put(x, cy+dy, ' ', material, 237)
		}
		f.Put(x, cy, '═', 250, 237)
	}

}
