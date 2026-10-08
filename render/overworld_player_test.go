package render

import "testing"

// Enlarging the map must never enlarge Grak's character footprint. Animation
// changes only its light, and must leave surrounding glyphs intact.
func TestOWPlayerSingleCharacterAtEveryScale(t *testing.T) {
	for scale := 1; scale <= 4; scale++ {
		st := NewOWState()
		st.TrailAge = []int{OWTrailMaxAge}
		l := GameLayout(OWW, OWH, OWW*scale+2, OWH*scale+ChromeTop+ChromeBot)
		seen := map[rune]bool{}
		for frame := 0; frame < 16; frame++ {
			f := &Frame{W: l.W, H: l.H, C: make([]Cell, l.W*l.H)}
			for i := range f.C {
				f.C[i] = Cell{R: 'x', FG: 123, BG: 234}
			}
			drawOWPlayer(f, l, st, frame)
			cx, cy := l.center(st.Cursor.X, st.Cursor.Y)
			for i, c := range f.C {
				if i == cy*f.W+cx {
					seen[c.R] = true
				} else if c.R != 'x' {
					t.Fatalf("scale %d frame %d: player overwrote neighbour at %d,%d", scale, frame, i%f.W, i/f.W)
				}
			}
		}
		if len(seen) != 1 || !seen['@'] {
			t.Fatalf("scale %d: walking glyph changed identity", scale)
		}
	}
}
