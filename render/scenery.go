package render

import (
	"math"

	"github.com/0xbenc/termtd/game"
)

// sceneryCanvas restricts every decorative mark to blocked map cells. It is
// shared by the material pass and the larger illustrations, including previews.
// Frame-space drawing lets masonry, branches and shorelines join at 2x+.
type sceneryCanvas struct {
	f           *Frame
	m           *game.Map
	l           Layout
	th          Theme
	frame       int
	occupied    map[game.Vec]bool
	riverColumn float64
}

func (s sceneryCanvas) put(x, y int, r rune, fg, bg int, bold bool) {
	if x < 0 || y < 0 || x >= s.m.W*s.l.Scale || y >= s.m.H*s.l.Scale {
		return
	}
	if cell(s.m, x/s.l.Scale, y/s.l.Scale) != game.CellWall {
		return
	}
	fx, fy := s.l.Ox+x, s.l.Oy+y
	if fx < 1 || fx >= s.f.W-1 || fy < ChromeTop || fy >= s.f.H-ChromeBot {
		return
	}
	s.f.Set(fx, fy, Cell{R: r, FG: fg, BG: bg, Bold: bold})
}

func drawScenery(f *Frame, m *game.Map, th Theme, l Layout, frame int) {
	if th.Stinger == StingerNone {
		return
	}
	s := sceneryCanvas{f: f, m: m, l: l, th: th, frame: frame, occupied: map[game.Vec]bool{}}
	if th.Stinger == StingerRift {
		x, _ := longestWallRun(m)
		s.riverColumn = float64(x)
	}
	for y := 0; y < m.H*l.Scale; y++ {
		for x := 0; x < m.W*l.Scale; x++ {
			mx, my := x/l.Scale, y/l.Scale
			if cell(m, mx, my) != game.CellWall {
				continue
			}
			r, fg, bg := s.material(x, y)
			s.put(x, y, r, fg, bg, false)
		}
	}
	s.landmarks()
}

func (s sceneryCanvas) material(x, y int) (rune, int, int) {
	scale := s.l.Scale
	mx, my := x/scale, y/scale
	edge := cell(s.m, mx-1, my) != game.CellWall || cell(s.m, mx+1, my) != game.CellWall ||
		cell(s.m, mx, my-1) != game.CellWall || cell(s.m, mx, my+1) != game.CellWall
	r, fg, bg := ' ', s.th.Wall, s.th.WallLo
	switch s.th.Stinger {
	case StingerHalls:
		// Bonded courses and recessed arches, with fluted pillar ends where
		// the masonry meets a chamber. Mortar stays dimmer than the road.
		fg = s.th.Wall
		if y%(2*scale) == 0 {
			r = '─'
			if (x+(y/(2*scale))*(3*scale))%(6*scale) == 0 {
				r = '┴'
			}
		} else if (x+(y/(2*scale))*(3*scale))%(6*scale) == 0 {
			r = '│'
		}
		if edge {
			bg = s.th.WallLo
			fg = s.th.Wall
			r = ' '
			if x%(2*scale) == scale/2 {
				r = '│'
			}
			if cell(s.m, mx, my+1) != game.CellWall && y%scale == scale-1 {
				r = '╨'
			}
		}
		if my == 0 && x%(8*scale) >= 2*scale && x%(8*scale) <= 5*scale {
			r = '╱'
			if x%(8*scale) > 3*scale {
				r = '╲'
			}
			fg = s.th.AccentDim
		}
	case StingerRotunda:
		// Concentric vault ribs surround the hoard; a carved cornice binds
		// the outside of the room into one structure.
		dx := float64(x)/float64(scale) - float64(s.m.W)/2
		dy := (float64(y)/float64(scale) - float64(s.m.H)/2) * 2.3
		rad := math.Hypot(dx, dy)
		if int(rad*1.3)%7 == 0 {
			r = '·'
			fg = s.th.WallHi
		}
		if edge {
			bg = s.th.WallLo
			fg = s.th.WallHi
			r = ' '
			if y%scale == scale-1 && cell(s.m, mx, my+1) != game.CellWall {
				r = '═'
			}
		}
		if my == 0 || my == s.m.H-1 {
			bg = s.th.WallLo
			r = '─'
			fg = s.th.WallHi
			if x%(6*scale) == 3*scale {
				r = '◆'
				fg = s.th.AccentDim
			}
		}
	case StingerGarden:
		// Unbuildable ground is a sunken basin: composed ripples and reeds
		// hug the bank, instead of random texture scattered over the floor.
		bg = 23
		fg = 24
		if (x+int(3*math.Sin(float64(y)/float64(scale))))%(9*scale) < 2*scale && y%(2*scale) == scale-1 {
			r = '≈'
			fg = 30
			if (s.frame/9+x/(3*scale)+y/scale)%7 == 0 {
				fg = 37
			}
		}
		if edge {
			bg = s.th.WallLo
			r = ' '
			fg = s.th.WallHi
			if x%(6*scale) == 0 {
				r = '╱'
			}
			if y%scale == scale-1 && x%(5*scale) == 0 {
				r = '╵'
				fg = 28
			}
		}
		if my == 0 || my == s.m.H-1 {
			r = ' '
			if (x+s.frame/20)%(12*scale) < 3*scale && y%scale == scale-1 {
				r = '~'
			}
			fg = 24
			bg = s.th.WallLo
		}
	case StingerRift:
		// A winding molten seam sits below basalt terraces. Its curves are
		// continuous in frame space and clipped by the playable ledges.
		river := s.riverColumn + .65*math.Sin(float64(y)/float64(scale)*.65)
		d := math.Abs(float64(x)/float64(scale) - river)
		bg = 232
		fg = 52
		if d < 1.2 {
			bg = 52
			fg = 130
			r = '≈'
			if (x+y+s.frame/6)%(5*scale) == 0 {
				r = '~'
				fg = 208
			}
		} else if d < 2.0 {
			r = '▒'
			fg = 88
		}
		if edge && d >= 1.2 {
			bg = s.th.WallLo
			r = ' '
			fg = 88
			if y%scale == scale-1 && (x/(2*scale)+y/scale)%3 == 0 {
				r = '╱'
				fg = 94
			}
		} else if y%(3*scale) == scale && (x+y/2)%(7*scale) < 3*scale {
			r = '─'
			fg = 237
		}
		if my == 0 || my == s.m.H-1 {
			bg = s.th.WallLo
			fg = 88
			r = '▴'
			if x%(4*scale) >= 2*scale {
				r = ' '
			}
		}
	case StingerHeart:
		// Rib arches and suspended chains frame a living, dark chamber.
		bg = s.th.WallLo
		fg = 89
		if (x/(3*scale)+y/(2*scale))%4 == 0 {
			r = '╱'
		}
		if (x/(3*scale)-y/(2*scale)+40)%4 == 0 {
			r = '╲'
		}
		if x%(10*scale) == 5*scale {
			r = '┆'
			fg = 95
		}
		if edge {
			r = ' '
			fg = 95
			bg = s.th.WallLo
			if y%scale == scale-1 {
				r = '═'
			}
		}
		if my == 0 || my == s.m.H-1 {
			bg = s.th.WallLo
			fg = 95
			r = '─'
			if x%(6*scale) == 0 {
				r = '╬'
				fg = 131
			}
		}
	case StingerDepths:
		// Cold, buried masonry, split by crystal seams. The relief follows
		// the seed's walls rather than relying on fixed landmark positions.
		bg = s.th.WallLo
		fg = 24
		if y%(2*scale) == 0 {
			r = '─'
		}
		if (x+y/(2*scale)*2)%(7*scale) == 0 {
			r = '│'
		}
		if edge {
			bg = s.th.WallLo
			fg = 31
			r = ' '
			if y%scale == scale-1 {
				r = '┄'
			}
		}
		if mx%11 == 5 && my%4 == 2 {
			r = '◇'
			fg = 37
			if (s.frame/12+mx+my)%6 == 0 {
				fg = 45
			}
		}
	}
	return r, fg, bg
}

// A stamp has a compact illustration and a more detailed terminal drawing.
// Both fit the same wall footprint. Larger terminal scales reveal detail
// without turning decorative structures into oversized gameplay glyphs.
type sceneryStamp struct {
	small, large []string
	fg           int
}

var dragonRelief = sceneryStamp{
	small: []string{" /\\/\\ ", "<__~_>"},
	large: []string{"   /\\__/\\  ", "__/ o  o \\_ ", "\\  __>   /  ", " \\_/\\__/   "}, fg: 178,
}
var treasureVault = sceneryStamp{
	small: []string{"╭───╮", "╰─◆─╯"},
	large: []string{" ╭────────╮ ", " │ ╭──╮ ◆ │ ", " │ ╰─◆╯ ● │ ", " ╰────────╯ "}, fg: 136,
}
var hallBanner = sceneryStamp{
	small: []string{"┬──┬", "│▼ │"},
	large: []string{"╥──────╥", "║ ╲  ╱ ║", "║  ◆   ║", "╨  ▼   ╨"}, fg: 130,
}
var gardenTree = sceneryStamp{
	small: []string{" ▄█▄ ", "╱ │ ╲"},
	large: []string{"  .╭──╮.  ", "╭──╯▓▓╰──╮", "╰──╮││╭──╯", "  ╱││││╲  "}, fg: 28,
}
var drownedArch = sceneryStamp{
	small: []string{"╭──╮", "│≈≈│"},
	large: []string{" ╭────╮ ", "╭╯ ╭╮ ╰╮", "│  ││  │", "┴≈≈┴┴≈≈┴"}, fg: 30,
}
var basaltPeak = sceneryStamp{
	small: []string{" /\\  ", "/╱╲\\_"},
	large: []string{"    /\\    ", " __/╱╲\\   ", "/ ╱│ │╲\\_ ", "_/─┴─┴─\\_"}, fg: 94,
}
var heartAltar = sceneryStamp{
	small: []string{"╲◆╱", "╘═╛"},
	large: []string{"╲    ╱", " ╲◆◆╱ ", "╭────╮", "╘════╛"}, fg: 131,
}
var boneMask = sceneryStamp{
	small: []string{"╭┬╮"},
	large: []string{"╭────╮", "╰┬┴┬╯"}, fg: 95,
}
var crystalRelief = sceneryStamp{
	small: []string{" /\\ ", "<◇◇>"},
	large: []string{"  /\\/\\ ", " /◇/◇/\\", "<◇/◇/◇/ ", " \\_\\_/ "}, fg: 37,
}

func (s sceneryCanvas) landmarks() {
	switch s.th.Stinger {
	case StingerRotunda:
		s.stamp(dragonRelief, s.m.W/2, s.m.H/2)
		s.stamp(treasureVault, 3, 2)
		s.stamp(treasureVault, s.m.W-5, s.m.H-3)
	case StingerHalls:
		s.stamp(hallBanner, 10, 3)
		s.stamp(hallBanner, 24, 7)
		s.stamp(hallBanner, s.m.W-6, 3)
	case StingerGarden:
		s.stamp(gardenTree, 8, 3)
		s.stamp(gardenTree, 16, 9)
		s.stamp(drownedArch, 24, 2)
		s.stamp(drownedArch, s.m.W-7, s.m.H-3)
	case StingerRift:
		s.stamp(basaltPeak, 11, 6)
		s.stamp(basaltPeak, 29, 8)
		s.stamp(basaltPeak, s.m.W-4, s.m.H-3)
	case StingerHeart:
		s.stamp(heartAltar, 6, 2)
		s.stamp(boneMask, s.m.W-5, 5)
		s.stamp(boneMask, 22, 6)
		s.stamp(boneMask, 30, 8)
	case StingerDepths:
		s.stamp(crystalRelief, 12, 3)
		s.stamp(crystalRelief, s.m.W-10, s.m.H-3)
	}
}

// Fit each illustration into existing solid rock, preferring its authored
// anchor. If a maze has no suitable pocket, the material pass still supplies
// its architecture. No object is cut across a road or build pad.
func (s sceneryCanvas) stamp(st sceneryStamp, ax, ay int) {
	w, h := 0, len(st.small)
	for _, line := range st.small {
		if n := len([]rune(line)); n > w {
			w = n
		}
	}
	// Detailed versions must fit even at 2x.
	for _, line := range st.large {
		if n := (len([]rune(line)) + 1) / 2; n > w {
			w = n
		}
	}
	if n := (len(st.large) + 1) / 2; n > h {
		h = n
	}
	best, bx, by := int(^uint(0)>>1), -1, -1
	for y := 0; y+h <= s.m.H; y++ {
		for x := 0; x+w <= s.m.W; x++ {
			fits := true
			for dy := 0; dy < h && fits; dy++ {
				for dx := 0; dx < w; dx++ {
					if cell(s.m, x+dx, y+dy) != game.CellWall || s.occupied[game.Vec{X: x + dx, Y: y + dy}] {
						fits = false
						break
					}
				}
			}
			if !fits {
				continue
			}
			dx, dy := x+w/2-ax, (y+h/2-ay)*2
			score := dx*dx + dy*dy
			if score < best {
				best, bx, by = score, x, y
			}
		}
	}
	if bx < 0 {
		return
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s.occupied[game.Vec{X: bx + x, Y: by + y}] = true
			for dy := 0; dy < s.l.Scale; dy++ {
				for dx := 0; dx < s.l.Scale; dx++ {
					s.put((bx+x)*s.l.Scale+dx, (by+y)*s.l.Scale+dy, ' ', st.fg, s.th.WallLo, false)
				}
			}
		}
	}
	rows := st.small
	if s.l.Scale >= 2 {
		rows = st.large
	}
	rw := 0
	for _, line := range rows {
		if n := len([]rune(line)); n > rw {
			rw = n
		}
	}
	ox := bx*s.l.Scale + (w*s.l.Scale-rw)/2
	oy := by*s.l.Scale + (h*s.l.Scale-len(rows))/2
	for y, line := range rows {
		for x, r := range []rune(line) {
			if r == ' ' {
				continue
			}
			s.put(ox+x, oy+y, r, st.fg, s.th.WallLo, false)
		}
	}
}
