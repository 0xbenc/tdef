package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestLightningPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderLightningMockup(w, h)
		if !reflect.DeepEqual(f, RenderLightningMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "LIGHTNING MAGE") || !strings.Contains(f.Text(), "esc return") {
			t.Fatalf("%v: missing title or controls", size)
		}
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		x0, y0 := (w-pw)/2, 4+(h-7-ph)/2
		count := func(u0, v0, u1, v1 float64, colors ...int) int {
			n := 0
			for y := 0; y < ph; y++ {
				for x := 0; x < pw; x++ {
					u, v := float64(x)/float64(pw-1), float64(y)/float64(ph-1)
					if u < u0 || u > u1 || v < v0 || v > v1 {
						continue
					}
					c := f.C[(y0+y)*w+x0+x]
					for _, fg := range colors {
						if c.R != ' ' && (c.FG == fg || c.BG == fg) {
							n++
							break
						}
					}
				}
			}
			return n
		}
		for _, mass := range []struct {
			name           string
			u0, v0, u1, v1 float64
			colors         []int
		}{
			{"face", .43, .24, .63, .44, []int{103, 139, 182}},
			{"receiving claw", .20, .19, .36, .34, []int{139, 182, 195}},
			{"directing claw", .78, .41, .93, .56, []int{139, 182, 195}},
			{"descending fork", .15, 0, .59, .24, []int{147, 195}},
			{"outward fork", .84, .25, .99, .57, []int{147, 195}},
			{"split robe", .25, .60, .87, .88, []int{54, 60, 97, 139}},
		} {
			if count(mass.u0, mass.v0, mass.u1, mass.v1, mass.colors...) == 0 {
				t.Fatalf("%v: lost %s", size, mass.name)
			}
		}
		for y := 0; y < h; y++ {
			if f.C[y*w].R != ' ' || f.C[y*w+w-1].R != ' ' {
				t.Fatalf("%v: art escaped frame", size)
			}
		}
	}
}

func TestLightningSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderLightningMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

// Reflected facial and claw lines must preserve the material beneath them.
func TestLightningFeaturesKeepTheirSurfaces(t *testing.T) {
	f := RenderLightningMockup(94, 47)
	marks := 0
	for _, c := range f.C {
		if (c.R == '╱' || c.R == '╲') && c.FG == 103 {
			marks++
			if c.BG != 139 && c.BG != 182 {
				t.Fatal("claw mark erased skin")
			}
		}
	}
	if marks < 2 {
		t.Fatal("lost open claw features")
	}
}
