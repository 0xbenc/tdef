package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestGnollPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderGnollMockup(w, h)
		if !reflect.DeepEqual(f, RenderGnollMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "GNOLL SLINGERS") || !strings.Contains(f.Text(), "esc return") {
			t.Fatalf("%v: title/controls missing", size)
		}
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		x0, y0 := (w-pw)/2, 4+(h-7-ph)/2
		// Both characters must retain substantial fur masses; the loop stays
		// above them and ammunition stays below the supplier's outstretched arm.
		count := func(u0, v0, u1, v1 float64, fg int) int {
			n := 0
			for y := 0; y < ph; y++ {
				for x := 0; x < pw; x++ {
					u, v := float64(x)/float64(pw-1), float64(y)/float64(ph-1)
					c := f.C[(y0+y)*w+x0+x]
					if u >= u0 && u <= u1 && v >= v0 && v <= v1 && c.FG == fg && c.R != ' ' {
						n++
					}
				}
			}
			return n
		}
		for _, mass := range []struct {
			name           string
			u0, v0, u1, v1 float64
			fg             int
		}{
			{"thrower", .39, .30, .75, .70, 137},
			{"supplier", .10, .50, .38, .87, 137},
			{"sling", .30, .00, .90, .22, 180},
			{"pouch", .29, .81, .45, .92, 130},
		} {
			if count(mass.u0, mass.v0, mass.u1, mass.v1, mass.fg) == 0 {
				t.Fatalf("%v: lost %s", size, mass.name)
			}
		}
		for y := 0; y < h; y++ {
			if f.C[y*w].R != ' ' || f.C[y*w+w-1].R != ' ' {
				t.Fatalf("%v: art escaped framing", size)
			}
		}
	}
}

func TestGnollPortraitEyesAndThinCord(t *testing.T) {
	f := RenderGnollMockup(94, 47)
	eyes, halfCord := 0, 0
	for _, c := range f.C {
		if c.R == '━' && c.FG == 220 {
			eyes++
			if c.BG == 233 {
				t.Fatal("eye punched through face")
			}
		}
		if (c.R == '▀' || c.R == '▄') && c.FG == 180 {
			halfCord++
		}
	}
	if eyes == 0 || halfCord == 0 {
		t.Fatal("missing eye or half-tile sling cord")
	}
}

func TestGnollPortraitSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderGnollMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad dimensions %v", size)
		}
	}
}
