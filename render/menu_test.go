package render

import (
	"testing"
)

func TestMenuRevealAndTerminalBackground(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		for _, revealed := range []bool{false, true} {
			f := RenderMenuAnimated(size[0], size[1], 0, 90, revealed, Palette())
			dragon, flame, gold := false, false, false
			artY := max(2, (size[1]-25)/2) + 6
			if size[0] >= 76 && size[1] >= 30 {
				artY += 5
			}
			for i, c := range f.C {
				if c.BG != 0 {
					t.Fatal("menu overrides terminal background")
				}
				if i/f.W < artY {
					continue
				}
				dragon = dragon || (c.R == '█' && c.FG == 131)
				flame = flame || (c.R == '▓' && c.FG == 220)
				gold = gold || (c.R == '█' && c.FG == 130)
			}
			if dragon != revealed || flame != revealed || !gold {
				t.Fatalf("%v revealed=%v: dragon=%v flame=%v gold=%v", size, revealed, dragon, flame, gold)
			}

		}
	}
}
