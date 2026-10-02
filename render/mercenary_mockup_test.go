package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestMercenaryPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderMercenaryMockup(w, h)
		if !reflect.DeepEqual(f, RenderMercenaryMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "MERCENARY") || !strings.Contains(f.Text(), "esc return") {
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
			{"dented open helmet", .38, .09, .62, .25, []int{239, 241, 246, 109}},
			{"scarred human face", .42, .24, .59, .38, []int{137, 223, 95, 173, 131}},
			{"salvaged shoulder plate", .22, .32, .44, .48, []int{239, 244, 246, 109}},
			{"quilted shoulder", .56, .34, .75, .5, []int{65, 101, 144, 58}},
			{"resting hands", .43, .46, .59, .56, []int{137, 223}},
			{"planted cleaver", .39, .65, .66, .95, []int{239, 244, 251, 109}},
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

func TestMercenarySmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderMercenaryMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestMercenaryExpressionAndHands(t *testing.T) {
	f := RenderMercenaryMockup(94, 47)
	eyes, hands := 0, 0
	for _, c := range f.C {
		if c.R == '━' && c.FG == 234 && c.BG == 223 {
			eyes++
		}
		if c.R == '│' && c.FG == 137 {
			hands++
			if c.BG != 137 && c.BG != 223 {
				t.Fatal("finger seam erased the hand")
			}
		}
	}
	if eyes != 2 || hands != 2 {
		t.Fatalf("expression or grip missing: eyes=%d hands=%d", eyes, hands)
	}
}
