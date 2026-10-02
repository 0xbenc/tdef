package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestRoguePortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderRogueMockup(w, h)
		if !reflect.DeepEqual(f, RenderRogueMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "ROGUE") || !strings.Contains(f.Text(), "esc return") {
			t.Fatalf("%v: missing title or controls", size)
		}
		ph := min(h-7, (w-6)/3)
		pw := ph * 3
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
			{"pointed hood", .47, .10, .74, .42, []int{23, 66, 236}},
			{"masked human face", .58, .23, .76, .41, []int{137, 180, 223, 95, 131}},
			{"trailing scarf", .03, .09, .57, .43, []int{52, 95, 131}},
			{"reverse knife", .04, .70, .22, .94, []int{239, 251, 110}},
			{"forward knife", .82, .26, .99, .51, []int{239, 251, 110}},
			{"crouching legs", .18, .62, .73, .95, []int{60, 67, 94, 137}},
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

func TestRogueSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderRogueMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestRogueFingerSeamsRetainSkin(t *testing.T) {
	f := RenderRogueMockup(94, 47)
	grips := 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 94 {
			grips++
			if c.BG != 137 && c.BG != 223 {
				t.Fatal("knife grip erased the hand")
			}
		}
	}
	if grips < 2 {
		t.Fatal("lost one of the knife grips")
	}
}
