package render

import (
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
)

const (
	OWRevealPanFrames   = 24
	OWRevealHoldFrames  = 36
	OWRevealFrames      = OWRevealPanFrames + OWUnsealFrames + OWRevealHoldFrames
	OWVictoryBeatFrames = 60 // existing return ring completes before a new room wakes
)

func drawOWPlayerIntro(f *Frame, l Layout, st OWState) {
	if !st.PlayerIntro || st.BootTTL > 0 || st.RevealBusy() || st.Descending != "" || st.ReturnFX.Floor != "" {
		return
	}
	cx, cy := l.center(st.Cursor.X, st.Cursor.Y)
	label := copytext.Text("overworld.orientation.you")
	hint := copytext.Text("overworld.orientation.move")
	putString(f, cx-len([]rune(label))/2, cy-2, label, 231, 24, true)
	f.Set(cx, cy-1, Cell{R: '▼', FG: 231, BG: 24, Bold: true})
	putString(f, cx-len([]rune(hint))/2, cy+2, hint, 231, 24, false)
}

// The same timed reveal has a distinct material for each room: broken amber
// seals, lit hall columns, water rings, lifting violet fog, or the heart blast.
// Everything is drawn in world space before labels and viewport cropping.
func drawOWReveal(f *Frame, l Layout, st OWState) {
	if st.RevealFloor == "" {
		return
	}
	id := st.RevealFloor
	if id == HeartFloorID {
		id = "rotunda"
	}
	n := owNodeByID(id)
	if n == nil {
		return
	}
	age := OWRevealFrames - st.RevealTTL
	for y := ChromeTop; y < f.H-ChromeBot; y++ {
		for x := 1; x < f.W-1; x++ {
			wx, wy := (x-l.Ox)/l.Scale, (y-l.Oy)/l.Scale
			if n.contains(wx, wy) {
				continue
			}
			cell := f.C[y*f.W+x]
			cell.FG, cell.Bold = 239, false
			if cell.BG >= 232 {
				cell.BG = 232 + (cell.BG-232)/3
			} else if cell.BG > 0 {
				cell.BG = 233
			}
			f.Set(x, y, cell)
		}
	}
	// A live ember runs from the Rotunda along this room's actual corridor.
	if ri := owRouteIndexFor(id); ri >= 0 && age < OWRevealPanFrames+OWUnsealFrames {
		route := owRouteCells[ri]
		if len(route) > 0 {
			reverse := abs(route[0].X-n.X)+abs(route[0].Y-n.Y) < abs(route[len(route)-1].X-n.X)+abs(route[len(route)-1].Y-n.Y)
			front := min(len(route)-1, age*(len(route)-1)/(OWRevealPanFrames+OWUnsealFrames/2))
			for k := 0; k <= front; k++ {
				idx := k
				if reverse {
					idx = len(route) - 1 - k
				}
				point := route[idx]
				if n.contains(point.X, point.Y) {
					continue
				}
				x, y := l.center(point.X, point.Y)
				if x < 0 || x >= f.W || y < ChromeTop || y >= f.H-ChromeBot {
					continue
				}
				cell := f.C[y*f.W+x]
				cell.FG, cell.Bold = n.Chrome, true
				if front-k < 3 {
					cell.R, cell.FG = '·', 231
				}
				f.Set(x, y, cell)
			}
		}
	}
	if age < OWRevealPanFrames || st.RevealTTL <= OWRevealHoldFrames {
		return
	}
	progress := float64(age-OWRevealPanFrames) / OWUnsealFrames
	radius := progress * float64(n.PW/2+3)
	for dy := -n.PH / 2; dy <= n.PH/2; dy++ {
		for dx := -n.PW / 2; dx <= n.PW/2; dx++ {
			distance := math.Hypot(float64(dx), float64(dy)*1.8)
			wave := math.Abs(distance-radius) < 0.7
			if id == "halls" {
				wave = abs(dx) == int(progress*float64(n.PW/2+1))
			}
			if id == "depths" {
				wave = math.Abs(float64(dx)-progress*float64(n.PW)+float64(n.PW/2)) < 1
			}
			if !wave {
				continue
			}
			x, y := l.center(n.X+dx, n.Y+dy)
			if x < 0 || x >= f.W || y < ChromeTop || y >= f.H-ChromeBot {
				continue
			}
			cell := f.C[y*f.W+x]
			cell.FG, cell.Bold = 231, true
			switch id {
			case "rift":
				if (dx+dy+age/3)%3 == 0 {
					cell.R = '╱'
				} else {
					cell.R = '╲'
				}
				cell.FG = 214
			case "halls":
				cell.FG = 159
			case "garden":
				if dy >= 0 {
					cell.R = '≈'
					cell.FG = 159
				}
			case "depths":
				cell.R = '~'
				cell.FG = 189
			case "rotunda":
				cell.R = '·'
				cell.FG = 226
			}
			f.Set(x, y, cell)
		}
	}
}
