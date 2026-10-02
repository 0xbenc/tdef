package render

import "testing"

// The illustration must stay inside the existing room footprint, leaving
// connected corridors and neighboring rooms untouched at every scale.
func TestOWRotundaIllustrationClippedToRoom(t *testing.T) {
	n := owNodeByID("rotunda")
	for scale := 1; scale <= 4; scale++ {
		w, h := OWW*scale+2, OWH*scale+ChromeTop+ChromeBot
		l := OverworldLayout(w, h, NewOWState())
		f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
		sentinel := Cell{R: '?', FG: 255, BG: 17}
		for i := range f.C {
			f.C[i] = sentinel
		}
		drawOWRotunda(f, l, n, 233, owPadView{status: OWOpen}, 90)
		x0, y0 := l.X(n.X-n.PW/2), l.Y(n.Y-n.PH/2)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				inside := x >= x0 && x < x0+n.PW*scale && y >= y0 && y < y0+n.PH*scale
				if !inside && f.C[y*w+x] != sentinel {
					t.Fatalf("scale %d: rotunda escaped its room at %d,%d", scale, x, y)
				}
			}
		}
	}
}
