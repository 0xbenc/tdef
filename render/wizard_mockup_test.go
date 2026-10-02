package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestWizardPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderWizardMockup(w, h)
		if !reflect.DeepEqual(f, RenderWizardMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "WIZARD") || !strings.Contains(f.Text(), "esc return") {
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
			{"crooked pointed hat", .19, .02, .64, .33, []int{60, 61, 97, 140}},
			{"human face and ivory beard", .39, .3, .57, .53, []int{137, 223, 250, 231}},
			{"heavy violet robe", .15, .41, .74, .96, []int{60, 61, 97, 140}},
			{"bowed wooden staff", .77, .17, .9, .96, []int{94, 137, 180}},
			{"diamond spell", .4, .45, .68, .71, []int{117, 153, 231, 81}},
			{"supporting hand", .4, .6, .55, .72, []int{137, 223}},
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

func TestWizardSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderWizardMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

// The spell needs three distinct masses: rim, negative space, and bright core.
func TestWizardSpellKeepsItsOpening(t *testing.T) {
	f := RenderWizardMockup(94, 47)
	ph := min(47-7, (94-6)*2/5)
	pw := ph * 5 / 2
	x0, y0 := (94-pw)/2, 4+(47-7-ph)/2
	colors := map[int]int{}
	for y := 0; y < ph; y++ {
		for x := 0; x < pw; x++ {
			u, v := float64(x)/float64(pw-1), float64(y)/float64(ph-1)
			if u < .41 || u > .67 || v < .45 || v > .64 {
				continue
			}
			c := f.C[(y0+y)*94+x0+x]
			if c.R != ' ' {
				colors[c.FG]++
				colors[c.BG]++
			}
		}
	}
	for _, color := range []int{117, 23, 231} {
		if colors[color] < 3 {
			t.Fatalf("spell lost rim, opening or core: %v", colors)
		}
	}
}

func TestWizardEyeStaysOnFace(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}} {
		f := RenderWizardMockup(size[0], size[1])
		eyes := 0
		for _, c := range f.C {
			if c.R == '━' && c.FG == 234 {
				eyes++
				if c.BG != 223 && c.BG != 137 && c.BG != 180 {
					t.Fatal("eye escaped face")
				}
			}
		}
		if eyes != 1 {
			t.Fatalf("%v: lost human eye", size)
		}
	}
}
