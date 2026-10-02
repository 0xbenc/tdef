package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestRangerPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderRangerMockup(w, h)
		if !reflect.DeepEqual(f, RenderRangerMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "RANGER") || !strings.Contains(f.Text(), "esc return") {
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
			{"face", .35, .24, .59, .36, []int{108, 151}},
			{"cloak", .03, .46, .29, .83, []int{65, 66, 236}},
			{"bow", .69, .06, .85, .85, []int{94, 137, 180}},
			{"arrowhead", .90, .34, 1, .51, []int{244, 110}},
			{"planted boots", .14, .81, .65, .94, []int{94, 137}},
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

func TestRangerSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderRangerMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestRangerGripMarksRetainSkin(t *testing.T) {
	f := RenderRangerMockup(94, 47)
	grips := 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 65 {
			grips++
			if c.BG != 108 && c.BG != 151 {
				t.Fatal("grip erased the skin")
			}
		}
	}
	if grips < 2 {
		t.Fatal("lost the draw and bow hand grips")
	}
}

func TestPortraitStrokePreservesUnderlyingHalves(t *testing.T) {
	for _, under := range []Cell{{R: '█', FG: 108, BG: 233}, {R: '▀', FG: 151, BG: 108}, {R: '▄', FG: 108, BG: 151}} {
		f := &Frame{W: 24, H: 14, C: make([]Cell, 24*14)}
		for i := range f.C {
			f.C[i] = under
		}
		p := portraitPainter{f, 2, 2, 20, 10}
		p.stroke(230, 1, portraitPoint{-.1, .2}, portraitPoint{.5, .8}, portraitPoint{1.1, .2})
		partial := 0
		for y := 0; y < f.H; y++ {
			for x := 0; x < f.W; x++ {
				c := f.C[y*f.W+x]
				if x < 2 || x >= 22 || y < 2 || y >= 12 {
					if c != under {
						t.Fatal("stroke escaped its local rectangle")
					}
					continue
				}
				if c == under || c.R != '▀' {
					continue
				}
				partial++
				top, bottom := under.BG, under.BG
				switch under.R {
				case '█':
					top, bottom = under.FG, under.FG
				case '▀':
					top = under.FG
				case '▄':
					bottom = under.FG
				}
				if c.FG != 230 && c.FG != top || c.BG != 230 && c.BG != bottom {
					t.Fatal("stroke erased the uncovered half")
				}
			}
		}
		if partial == 0 {
			t.Fatal("no partial stroke cells exercised")
		}
	}
}
