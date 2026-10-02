package render

import "math"

type portraitPoint struct{ x, y float64 }

// portraitPainter paints normalized, code-native tile shapes inside a local
// rectangle. Small facial marks keep the color of the mass beneath them.
type portraitPainter struct {
	f            *Frame
	x0, y0, w, h int
}

func (p portraitPainter) poly(fg int, points ...portraitPoint) {
	if p.w < 2 || p.h < 2 || len(points) < 3 {
		return
	}
	for y := 0; y < p.h; y++ {
		v := float64(y) / float64(p.h-1)
		for x := 0; x < p.w; x++ {
			u := float64(x) / float64(p.w-1)
			inside := false
			for i, a := range points {
				b := points[(i+1)%len(points)]
				if (a.y > v) != (b.y > v) && u < (b.x-a.x)*(v-a.y)/(b.y-a.y)+a.x {
					inside = !inside
				}
			}
			if inside {
				p.f.Set(p.x0+x, p.y0+y, Cell{R: '█', FG: fg, BG: 233})
			}
		}
	}
}

func (p portraitPainter) oval(fg int, cx, cy, rx, ry float64) {
	points := make([]portraitPoint, 32)
	for i := range points {
		a := float64(i) * math.Pi / 16
		points[i] = portraitPoint{cx + rx*math.Cos(a), cy + ry*math.Sin(a)}
	}
	p.poly(fg, points...)
}

// smoothOval samples two vertical half-cells. Partial edge cells retain the
// underlying half's color, so a round rim does not cut gaps into its carriage.
func (p portraitPainter) smoothOval(fg int, cx, cy, rx, ry float64) {
	if p.w < 2 || p.h < 2 || rx <= 0 || ry <= 0 {
		return
	}
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			px, py := p.x0+x, p.y0+y
			if px < 0 || px >= p.f.W || py < 0 || py >= p.f.H {
				continue
			}
			u := (float64(x)/float64(p.w-1) - cx) / rx
			hits := [2]bool{}
			for half := range hits {
				v := ((float64(y)+float64(half)*.5-.25)/float64(p.h-1) - cy) / ry
				hits[half] = u*u+v*v <= 1
			}
			if !hits[0] && !hits[1] {
				continue
			}
			old := p.f.C[py*p.f.W+px]
			top, bottom := old.BG, old.BG
			switch old.R {
			case '█':
				top, bottom = old.FG, old.FG
			case '▀':
				top = old.FG
			case '▄':
				bottom = old.FG
			}
			if hits[0] {
				top = fg
			}
			if hits[1] {
				bottom = fg
			}
			c := Cell{R: '▀', FG: top, BG: bottom}
			if top == bottom {
				c = Cell{R: '█', FG: top, BG: 233}
			}
			p.f.Set(px, py, c)
		}
	}
}

// stroke follows a path at a constant physical width. A vertical half-cell
// is approximately one horizontal cell, so curves and cords keep their weight
// in both directions. Partial tiles preserve the surface behind the stroke.
func (p portraitPainter) stroke(fg int, width float64, points ...portraitPoint) {
	if p.w < 2 || p.h < 2 || width <= 0 || len(points) < 2 {
		return
	}
	path := make([]portraitPoint, len(points))
	for i, a := range points {
		path[i] = portraitPoint{a.x * float64(p.w-1), a.y * float64(p.h-1) * 2}
	}
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			px, py := p.x0+x, p.y0+y
			if px < 0 || px >= p.f.W || py < 0 || py >= p.f.H {
				continue
			}
			hits := [2]bool{}
			for half := range hits {
				qx, qy := float64(x), float64(y)*2+float64(half)-.5
				for i, a := range path[:len(path)-1] {
					b := path[i+1]
					dx, dy := b.x-a.x, b.y-a.y
					t := 0.0
					if d := dx*dx + dy*dy; d > 0 {
						t = math.Max(0, math.Min(1, ((qx-a.x)*dx+(qy-a.y)*dy)/d))
					}
					ex, ey := qx-(a.x+t*dx), qy-(a.y+t*dy)
					if ex*ex+ey*ey <= width*width/4 {
						hits[half] = true
						break
					}
				}
			}
			if !hits[0] && !hits[1] {
				continue
			}
			old := p.f.C[py*p.f.W+px]
			top, bottom := old.BG, old.BG
			switch old.R {
			case '█':
				top, bottom = old.FG, old.FG
			case '▀':
				top = old.FG
			case '▄':
				bottom = old.FG
			}
			if hits[0] {
				top = fg
			}
			if hits[1] {
				bottom = fg
			}
			c := Cell{R: '▀', FG: top, BG: bottom}
			if top == bottom {
				c = Cell{R: '█', FG: top, BG: 233}
			}
			p.f.Set(px, py, c)
		}
	}
}

// portraitCurve samples a cubic path for a continuous recurved limb or cord.
func portraitCurve(a, b, c, d portraitPoint) []portraitPoint {
	points := make([]portraitPoint, 25)
	for i := range points {
		t := float64(i) / 24
		s := 1 - t
		points[i] = portraitPoint{s*s*s*a.x + 3*s*s*t*b.x + 3*s*t*t*c.x + t*t*t*d.x, s*s*s*a.y + 3*s*s*t*b.y + 3*s*t*t*c.y + t*t*t*d.y}
	}
	return points
}

// ring samples two vertical half-cells so a thin cord stays continuous without
// becoming a broad oval band on terminals whose cells are taller than wide.
func (p portraitPainter) ring(fg int, cx, cy, rx, ry float64) {
	if p.w < 2 || p.h < 2 {
		return
	}
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			u := float64(x) / float64(p.w-1)
			hits := [2]bool{}
			for half := range hits {
				v := (float64(y) + float64(half)*.5 - .25) / float64(p.h-1)
				nx, ny := (u-cx)/rx, (v-cy)/ry
				d := math.Hypot(nx, ny)
				band := .7 * math.Hypot(nx/(rx*float64(p.w-1)), ny/(ry*float64(p.h-1)*2))
				hits[half] = math.Abs(d-1) < band
			}
			r := '█'
			if !hits[0] && !hits[1] {
				continue
			}
			if !hits[0] {
				r = '▄'
			} else if !hits[1] {
				r = '▀'
			}
			p.f.Set(p.x0+x, p.y0+y, Cell{R: r, FG: fg, BG: 233})
		}
	}
}

func (p portraitPainter) detail(u, v float64, r rune, fg int) {
	x, y := p.x0+int(u*float64(p.w-1)), p.y0+int(v*float64(p.h-1))
	if x < 0 || x >= p.f.W || y < 0 || y >= p.f.H {
		return
	}
	c := p.f.C[y*p.f.W+x]
	if c.R == '█' {
		p.f.Set(x, y, Cell{R: r, FG: fg, BG: c.FG})
	}
}
