package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestSquirePortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderSquireMockup(w, h)
		if !reflect.DeepEqual(f, RenderSquireMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "SQUIRE") || !strings.Contains(f.Text(), "esc return") {
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
			{"tilted kettle helm", .48, .19, .78, .35, []int{239, 246, 109, 251}},
			{"youthful human face", .53, .33, .69, .43, []int{137, 223, 180}},
			{"oversized diagonal blade", .15, .07, .51, .55, []int{239, 251, 109, 246}},
			{"two-handed grip", .45, .54, .58, .65, []int{137, 223}},
			{"short guild surcoat", .47, .42, .68, .78, []int{95, 131, 173, 52}},
			{"back-strapped shield", .75, .35, .9, .66, []int{94, 137, 180}},
			{"braced boots", .23, .87, .8, .97, []int{94, 137, 180}},
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

func TestSquireSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderSquireMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestSquireFaceAndTwoHandedGrip(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}} {
		f := RenderSquireMockup(size[0], size[1])
		eyes := 0
		for _, c := range f.C {
			if c.R == '━' && c.FG == 234 {
				eyes++
				if c.BG != 223 && c.BG != 137 {
					t.Fatal("eye escaped face")
				}
			}
		}
		if eyes != 1 {
			t.Fatalf("%v: lost youthful face", size)
		}
	}
	f := RenderSquireMockup(94, 47)
	hands := 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 137 {
			hands++
			if c.BG != 223 && c.BG != 137 {
				t.Fatal("finger seam erased a hand")
			}
		}
	}
	if hands != 2 {
		t.Fatalf("lost two-handed grip: %d", hands)
	}
}
