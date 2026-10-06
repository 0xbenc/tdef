package render

import "math"

const (
	bootComboLaunch = 12
	bootComboImpact = 60
	bootComboLen    = 76
)

// The siege finale shares one clock: a visibly heavy launch, a cold front
// crossing the formation, a held frozen pose, then the returning stone.
func drawBootSiegeFinale(f *Frame, w, h, off, s int) {
	if s < 0 || s >= bootComboLen {
		return
	}
	x0 := max(0, (w-42)/2)
	cx, floor := w/2, h-3
	impactAge := s - bootComboImpact
	drawBootIceGround(f, cx, floor, s)
	drawBootFrostCast(f, x0, off, cx, floor, s)
	for i, dx := range []int{-9, 0, 9} {
		freezeAt := 42 - i*8
		// The center takes the stone; its neighbors break with the shockwave.
		hitAt := bootComboImpact + abs(dx)/3
		if s < hitAt {
			drawBootSiegeTarget(f, cx+dx, floor, s, freezeAt, i)
		}
	}
	drawBootTrebuchet(f, x0, off, cx, floor, s)
	if impactAge >= 0 {
		drawBootIceImpact(f, cx, floor, impactAge)
	}
	if impactAge >= 0 && impactAge < 4 {
		// A brief camera kick gives the stone weight without shaking the UI
		// once it has appeared. Every frame is derived from the same clock.
		dx := [4]int{1, -1, 1, 0}[impactAge]
		dy := [4]int{0, 1, 0, 0}[impactAge]
		cells := append([]Cell(nil), f.C...)
		for i := range f.C {
			f.C[i] = Cell{R: ' '}
		}
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				f.Set(x+dx, y+dy, cells[y*w+x])
			}
		}
	}
}

func clearBootLetter(f *Frame, x0, off, li int) {
	for row := 0; row < 6; row++ {
		for col := 0; col < 6; col++ {
			f.Set(x0+li*7+col, off+2+row, Cell{R: ' '})
		}
	}
}

// bootSiegeStone follows a high parabola, briefly passing above the terminal.
// It returns exactly to the middle champion when the frost has stopped them.
func bootSiegeStone(x0, off, cx, floor, s int) (int, int) {
	u := float64(s-bootComboLaunch) / float64(bootComboImpact-bootComboLaunch)
	mx, my := float64(x0+4*7-1), float64(off+1)
	arc := float64(floor+3)*0.48 + float64(off+2)
	x := mx + (float64(cx)-mx)*u
	y := my + (float64(floor-2)-my)*u - 4*arc*u*(1-u)
	return int(math.Round(x)), int(math.Round(y))
}

func drawBootTrebuchet(f *Frame, x0, off, cx, floor, s int) {
	const li = 4
	if s <= bootComboLaunch+6 {
		clearBootLetter(f, x0, off, li)
		dy, fg := 0, titleColors[li]
		if s < bootComboLaunch {
			if s >= 6 {
				fg = 220
			}
			if s >= 9 {
				fg = 229
			}
		} else {
			dy = max(0, 2-(s-bootComboLaunch)/3)
			fg = 255
		}
		drawBootLetter(f, f.W, off, li, 0, dy, fg, s >= 9)
	}
	if s < bootComboLaunch {
		// A taut sling pivots over the crossbar while its stone winds back.
		mx, my := x0+li*7+2, off+2
		a := -float64(s+2) / float64(bootComboLaunch+2) * math.Pi
		bx := mx + int(math.Round(3*math.Cos(a)))
		by := my + int(math.Round(2*math.Sin(a)))
		for k := 1; k <= 3; k++ {
			x, y := shotPos(mx, my, bx, by, float64(k)/3)
			f.Set(x, y, Cell{R: '·', FG: 180})
		}
		f.Set(bx, by, Cell{R: '◆', FG: 229, Bold: true})
		return
	}
	if s >= bootComboImpact {
		return
	}
	// A dim falling shadow tightens into the impact point during the freeze.
	if s >= 42 {
		radius := max(1, (bootComboImpact-s)/4)
		for dx := -radius; dx <= radius; dx++ {
			f.Set(cx+dx, floor+1, Cell{R: '━', FG: 60})
		}
	}
	for k := 5; k >= 1; k-- {
		if s-k < bootComboLaunch {
			continue
		}
		x, y := bootSiegeStone(x0, off, cx, floor, s-k)
		f.Set(x, y, Cell{R: '·', FG: 238 + (5-k)*2})
	}
	x, y := bootSiegeStone(x0, off, cx, floor, s)
	fg := 180
	if s > bootComboImpact-6 {
		fg = 229
	}
	putString(f, x-1, y-1, "▄█▄", fg, 0, false)
	putString(f, x-1, y, "▀█▀", fg, 0, true)
	f.Set(x-1, y-1, Cell{R: '▗', FG: 250})
}

func drawBootFrostCast(f *Frame, x0, off, cx, floor, s int) {
	const li = 5
	if s < 14 || s >= 48 {
		return
	}
	charge := s - 14
	fg := 117
	if charge < 6 {
		fg = 153
	} else if charge < 10 {
		fg = 195
	} else if charge < 13 {
		fg = 255
	}
	drawBootLetter(f, f.W, off, li, 0, 0, fg, charge >= 6 && charge < 13)
	sx, sy := x0+li*7+3, off+5
	if charge < 10 {
		for i := 0; i < 5; i++ {
			x := sx + titleHash(i, charge/2, li)%7 - 3
			y := sy + titleHash(li, i, charge/2)%7 - 3
			f.Set(x, y, Cell{R: '·', FG: 117})
		}
		return
	}
	// A curved, widening frost stream hits the right flank first. The cold
	// then runs along the ground, climbing boots before encasing the figures.
	travel := min(1.0, float64(charge-9)/12)
	for k := 0; k <= 24; k++ {
		u := float64(k) / 24
		if u > travel {
			break
		}
		x := int(math.Round((1-u)*(1-u)*float64(sx) + 2*(1-u)*u*float64(cx+18) + u*u*float64(cx+12)))
		y := int(math.Round((1-u)*(1-u)*float64(sy) + 2*(1-u)*u*float64(floor-5) + u*u*float64(floor)))
		for spread := -1; spread <= 1; spread++ {
			if spread != 0 && (k+s)%3 != 0 {
				continue
			}
			r := [4]rune{'·', '+', '❄', '╱'}[(k+s/2)%4]
			color := 117
			if spread == 0 {
				color = 195
			}
			f.Set(x+spread, y, Cell{R: r, FG: color, Bold: spread == 0})
		}
	}
}

func drawBootIceGround(f *Frame, cx, floor, s int) {
	for x := cx - 22; x <= cx+22; x++ {
		f.Set(x, floor+1, Cell{R: '─', FG: 235})
	}
	if s < 25 {
		return
	}
	left := cx + 12 - min(38, (s-25)*2)
	fg := 60
	if s >= 34 {
		fg = 117
	}
	if s >= bootComboImpact+10 {
		fg = 60
	}
	for x := max(cx-22, left); x <= cx+22; x++ {
		f.Set(x, floor, Cell{R: '▁', FG: fg})
		f.Set(x, floor+1, Cell{R: '━', FG: fg})
		seed := titleHash(x, floor, 7)
		if seed%4 == 0 {
			rise := min(3, max(0, (s-28)/5))
			for y := 0; y < rise; y++ {
				f.Set(x, floor-y, Cell{R: [3]rune{'╱', '│', '╲'}[seed%3], FG: fg})
			}
		}
	}
}

func drawBootSiegeTarget(f *Frame, tx, floor, s, freezeAt, i int) {
	pose := min(s, freezeAt)
	approach := max(0, (freezeAt-pose)/6)
	if i == 0 {
		tx -= approach
	} else {
		tx += approach
	}
	rows := []string{" o ", "╱█╲", "╱ ╲"}
	if i == 1 {
		rows = []string{"[O]", "╞╬╡", "╱ ╲"}
	}
	if pose/4%2 == 0 {
		rows[2] = " ╳ "
	}
	freezeAge := s - freezeAt
	for row, text := range rows {
		for col, r := range []rune(text) {
			if r == ' ' {
				continue
			}
			fg, bold := 244, false
			if row == 1 {
				fg = 167 + i*6
			}
			if freezeAge >= 0 {
				fg = 117
				if row == 2 || freezeAge >= 3 {
					fg = 153
				}
				if freezeAge >= 6 {
					fg = 195
				}
				bold = true
				if row == 0 {
					r = '◇'
				}
				if row == 1 && col == 1 {
					r = '▓'
				}
			}
			f.Set(tx-1+col, floor-3+row, Cell{R: r, FG: fg, Bold: bold})
		}
	}
	if freezeAge >= 7 {
		// Stillness matters: these crystals and the trapped pose stop changing
		// until the stone descends. Cracks only appear just before impact.
		f.Set(tx-2, floor-2, Cell{R: '╱', FG: 117})
		f.Set(tx+2, floor-1, Cell{R: '╲', FG: 117})
		if s >= bootComboImpact-5 {
			f.Set(tx, floor-2, Cell{R: '╱', FG: 255, Bold: true})
		}
	}
}

func drawBootIceImpact(f *Frame, cx, floor, age int) {
	if age < 3 {
		// A compact white core, followed by the wide, low pressure ring.
		for y := -3; y <= 0; y++ {
			for x := -3; x <= 3; x++ {
				if abs(x)+abs(y) <= 4-age {
					f.Set(cx+x, floor+y, Cell{R: '█', FG: 255, Bold: true})
				}
			}
		}
	}
	if age < 12 {
		rx, ry := float64(2+age*2), float64(1+age)/3
		for k := 0; k < 48; k++ {
			a := math.Pi + math.Pi*float64(k)/47
			x := cx + int(math.Round(rx*math.Cos(a)))
			y := floor + int(math.Round(ry*math.Sin(a)))
			fg := 195
			if age < 2 {
				fg = 255
			} else if age > 7 {
				fg = 60
			}
			f.Set(x, y, Cell{R: '·', FG: fg, Bold: age < 4})
		}
	}
	for i, dx := range []int{-9, 0, 9} {
		a := age - abs(dx)/3
		if a < 0 {
			continue
		}
		for k := 0; k < 12; k++ {
			seed := titleHash(k, i, 19)
			vx := float64(seed%17-8) / 7
			vy := -float64(4+seed%9) / 9
			t := float64(a)
			x := cx + dx + int(math.Round(vx*t))
			y := floor - 2 + int(math.Round(vy*t+0.055*t*t))
			if y > floor+1 {
				continue
			}
			fg := [4]int{117, 153, 195, 180}[k%4]
			if a < 2 {
				fg = 255
			} else if a > 9 {
				fg = 60
			}
			f.Set(x, y, Cell{R: [5]rune{'▪', '╱', '╲', '▴', '·'}[seed%5], FG: fg, Bold: a < 5})
		}
	}
}
