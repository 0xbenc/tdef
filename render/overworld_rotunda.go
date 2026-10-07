package render

import "math"

// The hub reads as three masses: a stone vault, a sleeping horned dragon,
// and the hoard inside its curled tail. All geometry is local to the room's
// terminal footprint, so its silhouette survives the integer scale steps.
func drawOWRotunda(f *Frame, l Layout, n *owNode, bg int, v owPadView, frame int) {
	x0, y0 := l.X(n.X-n.PW/2), l.Y(n.Y-n.PH/2)
	w, h := n.PW*l.Scale, n.PH*l.Scale
	drawDragonTableau(f, x0, y0, w, h, l.Scale, bg, v, frame, false)
}

// drawDragonTableau shares the pose between the room and its portrait. The
// larger study uses filled horns, a resting foreleg and broad wing folds.
func drawDragonTableau(f *Frame, x0, y0, w, h, scale, bg int, v owPadView, frame int, portrait bool) {
	drawRotundaTableau(f, x0, y0, w, h, scale, bg, v, frame, portrait, true)
}

func drawRotundaTableau(f *Frame, x0, y0, w, h, scale, bg int, v owPadView, frame int, portrait, dragon bool) {
	if w < 2 || h < 2 {
		return
	}
	open := v.status != OWSealed || v.unseal > 0
	color := func(c int) int {
		if !open {
			c = dim(c)
		}
		if v.unseal > 0 {
			c = lerp(dim(c), c, v.unseal)
		}
		if v.flash > 0 {
			c = lerp(c, v.flashC, v.flash)
		}
		return c
	}
	put := func(x, y int, r rune, fg int) {
		if x >= 0 && x < w && y >= 0 && y < h {
			f.Set(x0+x, y0+y, Cell{R: r, FG: color(fg), BG: bg})
		}
	}
	// A broad, rounded vault with thick piers and one quiet keystone. Its
	// stone is subdued; the dragon and gold carry the figure-ground contrast.
	roof := max(1, int(float64(h)*.28))
	for x := 0; x < w; x++ {
		u := float64(x)/float64(w-1)*2 - 1
		y := int(math.Round(float64(roof) * (1 - math.Sqrt(math.Max(0, 1-u*u)))))
		r := '━'
		if x < w/5 {
			r = '╱'
		} else if x > w*4/5 {
			r = '╲'
		}
		put(x, y, r, 240)
		if scale >= 2 {
			put(x, y+1, '▀', 236)
		}
	}
	put(w/2, 0, '◆', 244)
	for y := roof + 1; y < h-1; y++ {
		for dx := 0; dx < max(1, scale); dx++ {
			if portrait && dx > 0 {
				put(dx, y, '█', 236)
				put(w-1-dx, y, '█', 236)
			} else {
				put(dx, y, '▐', 238)
				put(w-1-dx, y, '▌', 238)
			}
		}
	}
	for x := 0; x < w; x++ {
		put(x, h-1, '━', 238)
	}
	// Gold is a single low hill, rather than scattered dots. A small number
	// of facets and moving glints give the pile life without dissolving it.
	facetPeriod, glintPeriod := 7, 29
	for x := scale + 1; x < w-scale-1; x++ {
		u := (float64(x)/float64(w-1) - .54) / .30
		height := int(math.Round(float64(scale) * (.6 + 1.7*math.Max(0, 1-u*u))))
		for y := h - 2 - height; y < h-1; y++ {
			r, fg := '█', 130
			if y == h-2-height {
				r = '▄'
				fg = 178
			}
			facet := (x+2*y)%facetPeriod == 0
			glint := (x+2*y+frame/10)%glintPeriod == 0
			if portrait {
				// Irregular, sparse facets avoid the diagonal wallpaper pattern.
				facet = owHash(x, y, 17)%47 == 0
				glint = owHash(x, y, 31)%113 == 0
			}
			if y > h-2-height && facet {
				r = '◆'
				fg = 178
			}
			if glint {
				r = '◆'
				fg = 220
			}
			if portrait && r == '◆' {
				f.Set(x0+x, y0+y, Cell{R: r, FG: color(fg), BG: color(130)})
			} else {
				put(x, y, r, fg)
			}
		}
	}
	if !dragon {
		return
	}
	// The body and projecting muzzle form one solid silhouette. The folded
	// triangular wing, pale horns and closed eye identify it as a dragon.
	for y := 1; y < h-1; y++ {
		for x := scale; x < w-scale; x++ {
			u, vv := float64(x)/float64(w-1), float64(y)/float64(h-1)
			body := math.Pow((u-.64)/.24, 2)+math.Pow((vv-.46)/.25, 2) <= 1
			neck := rotundaTriangle(u, vv, [6]float64{.27, .48, .50, .34, .57, .62})
			head := u >= .19 && u <= .38 && vv >= .40 && vv <= .57
			muzzle := u >= .10 && u <= .25 && vv >= .48 && vv <= .57
			if portrait && u < .15 {
				muzzle = muzzle && u >= .10+math.Abs(vv-.52)*.7
			}
			if body || neck || head || muzzle {
				fg := 131
				if portrait {
					// Curved underside shade gives the torso volume without texture.
					if (body && vv > .46 && math.Hypot((u-.60)/.24, (vv-.37)/.25) > 1) || (!body && vv > .55) {
						fg = 95
					}
					if body && u > .66 && vv < .42 && math.Hypot((u-.61)/.24, (vv-.49)/.25) > 1 {
						fg = 173
					}
				}
				put(x, y, '█', fg)
			}
			wing := rotundaTriangle(u, vv, [6]float64{.40, .50, .60, .15, .79, .55})
			if wing {
				put(x, y, '█', 95)
				edge := .045
				if portrait {
					edge = .022
				}
				if portrait && rotundaTriangle(u, vv, [6]float64{.60, .19, .57, .51, .70, .53}) {
					put(x, y, '█', 131)
				}
				if portrait && rotundaTriangle(u, vv, [6]float64{.60, .20, .60, .51, .64, .52}) {
					put(x, y, '█', 95)
				}
				if math.Abs(vv-(.50-(u-.40)*1.75)) < edge {
					r := '╱'
					if portrait {
						r = '█'
					}
					put(x, y, r, 173)
				}
				if math.Abs(vv-(.15+(u-.60)*2.10)) < edge {
					r := '╲'
					if portrait {
						r = '█'
					}
					put(x, y, r, 173)
				}
			}
			if portrait {
				// A bent foreleg joins the shoulder to the paw under the head.
				arm := rotundaTriangle(u, vv, [6]float64{.48, .52, .54, .69, .38, .70})
				paw := u >= .28 && u <= .45 && vv >= .64 && vv <= .71
				if arm || paw {
					put(x, y, '█', 131)
				}
				// Two swept horns are tapered ivory wedges, rather than outlines.
				if rotundaTriangle(u, vv, [6]float64{.27, .41, .40, .20, .34, .41}) ||
					rotundaTriangle(u, vv, [6]float64{.36, .40, .49, .25, .42, .43}) {
					put(x, y, '█', 180)
				}
			}
			// The tail curls around the gold, leaving its center exposed.
			tail := math.Hypot((u-.56)/.32, (vv-.67)/.20)
			inner, end := .78, .67
			if portrait && u < .42 {
				inner, end = .87, .73 // taper the tip below the resting paw
			}
			if tail >= inner && tail <= 1.05 && (vv > end || u > .76) {
				fg := 131
				if portrait && tail > .97 && vv > .75 {
					fg = 95
				}
				put(x, y, '█', fg)
			}
		}
	}
	// Horns lean back from the brow. At 1x the same gestures reduce to a
	// few clear diagonal strokes instead of a miniature detailed drawing.
	hx, hy := int(math.Round(float64(w-1)*.29)), int(math.Ceil(float64(h-1)*.40))
	for d := 0; d < max(1, scale) && !portrait; d++ {
		put(hx+d, hy-1-d, '╱', 180)
		put(hx+2*scale+d, hy-1-d, '╱', 180)
	}
	ex, ey := int(math.Round(float64(w-1)*.23)), int(math.Round(float64(h-1)*.47))
	eye := '─'
	if open && frame%150 >= 144 {
		eye = '·'
	}
	if portrait {
		// Details inherit the skin beneath them. A dark cell background would
		// punch rectangular holes through the eye and paw in a real terminal.
		surface := func(x, y int, r rune, fg int) {
			if x < 0 || x >= w || y < 0 || y >= h {
				return
			}
			px, py := x0+x, y0+y
			if px < 0 || px >= f.W || py < 0 || py >= f.H {
				return
			}
			c := f.C[py*f.W+px]
			skin := c.FG
			if c.R != '█' {
				if c.BG == bg {
					return
				}
				skin = c.BG
			}
			f.Set(px, py, Cell{R: r, FG: color(fg), BG: skin})
		}
		for dx := 0; dx < max(2, w/30); dx++ {
			surface(ex+dx, ey, '━', 180)
		}
		surface(int(float64(w-1)*.12), int(float64(h-1)*.52), '▪', 52)
		// The mouth stays on the muzzle; the tooth is anchored in the jaw.
		jy := int(float64(h-1) * .55)
		for dx := 0; dx < max(2, w/12); dx++ {
			surface(int(float64(w-1)*.13)+dx, jy, '━', 95)
		}
		surface(int(float64(w-1)*.15), jy, '▼', 180)
		for claw := 0; claw < 3; claw++ {
			surface(int(float64(w)*.29)+claw*max(1, w/35), int(float64(h-1)*.69), '╲', 180)
		}
	} else {
		put(ex, ey, eye, 180)
		put(int(math.Round(float64(w-1)*.12)), int(math.Round(float64(h-1)*.56)), '╲', 173)
	}
	// A single wing rib and a foreclaw supply the last identifying details.
	if scale >= 2 && !portrait {
		for d := 0; d < scale+1; d++ {
			put(int(float64(w)*.59)+d, int(float64(h)*.31)+d, '╲', 131)
		}
		put(int(float64(w)*.36), int(float64(h)*.64), '╰', 173)
		put(int(float64(w)*.36)+1, int(float64(h)*.64), '─', 173)
	}
}

func rotundaTriangle(x, y float64, p [6]float64) bool {
	cross := func(ax, ay, bx, by float64) float64 { return (x-bx)*(ay-by) - (ax-bx)*(y-by) }
	a, b, c := cross(p[0], p[1], p[2], p[3]), cross(p[2], p[3], p[4], p[5]), cross(p[4], p[5], p[0], p[1])
	return !((a < 0 || b < 0 || c < 0) && (a > 0 || b > 0 || c > 0))
}
