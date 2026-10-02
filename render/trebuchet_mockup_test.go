package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestTrebuchetPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderTrebuchetMockup(w, h)
		if !reflect.DeepEqual(f, RenderTrebuchetMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "TREBUCHET") || !strings.Contains(f.Text(), "esc return") {
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
			{"throwing arm", .10, .04, .49, .35, []int{94, 95, 137, 180}},
			{"loaded sling", .06, .34, .24, .47, []int{237, 239, 243, 247}},
			{"ballast", .63, .54, .83, .76, []int{94, 95, 137, 239}},
			{"timber frame", .25, .39, .78, .92, []int{94, 130, 137, 180}},
			{"winding crew", .10, .59, .27, .82, []int{65, 107, 150}},
			{"spotting crew", .77, .57, .91, .78, []int{101, 143, 186}},
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

func TestTrebuchetSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderTrebuchetMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestTrebuchetGripsRetainSkin(t *testing.T) {
	f := RenderTrebuchetMockup(94, 47)
	grips := 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 65 {
			grips++
			if c.BG != 107 && c.BG != 150 {
				t.Fatal("winding grip erased skin")
			}
		}
	}
	if grips < 1 {
		t.Fatal("lost winding crew grip")
	}
}
