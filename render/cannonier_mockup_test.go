package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestCannonierPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderCannonierMockup(w, h)
		if !reflect.DeepEqual(f, RenderCannonierMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "CANNONIER") || !strings.Contains(f.Text(), "esc return") {
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
						if c.R != ' ' && c.FG == fg {
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
			{"loader", .05, .23, .32, .65, []int{107, 65, 150}},
			{"spotter", .59, .14, .84, .33, []int{143, 101, 186}},
			{"open bore", .28, .41, .48, .66, []int{232}},
			{"front wheel", .50, .61, .77, .945, []int{130, 137}},
			{"back wheel", .28, .67, .47, .92, []int{137, 95}},
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

func TestCannonierGripsRetainSkin(t *testing.T) {
	f := RenderCannonierMockup(94, 47)
	loader, spotter, curves := 0, 0, 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 65 {
			loader++
			if c.BG != 107 && c.BG != 150 {
				t.Fatal("loader grip erased skin")
			}
		}
		if c.R == '│' && c.FG == 101 {
			spotter++
			if c.BG != 143 && c.BG != 186 {
				t.Fatal("spotter grip erased skin")
			}
		}
		if c.R == '▀' {
			curves++
		}
	}
	if loader < 1 || spotter < 1 || curves < 10 {
		t.Fatalf("lost grips or round contours: %d %d %d", loader, spotter, curves)
	}
}

func TestCannonierSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderCannonierMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

// A partial curved tile must retain the underlying wood/metal, not replace
// the uncovered half with the scene background. Repeat for both edge halves.
func TestSmoothOvalCompositesUnderlyingHalves(t *testing.T) {
	for _, under := range []Cell{
		{R: '█', FG: 130, BG: 233}, {R: '▀', FG: 180, BG: 94}, {R: '▄', FG: 94, BG: 180},
	} {
		f := &Frame{W: 24, H: 14, C: make([]Cell, 24*14)}
		for i := range f.C {
			f.C[i] = under
		}
		p := portraitPainter{f, 2, 2, 20, 10}
		p.smoothOval(242, .5, .5, .35, .31)
		partial := 0
		for y := 0; y < f.H; y++ {
			for x := 0; x < f.W; x++ {
				c := f.C[y*f.W+x]
				if x < 2 || x >= 22 || y < 2 || y >= 12 {
					if c != under {
						t.Fatal("oval escaped its local rectangle")
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
				if c.FG != 242 && c.FG != top || c.BG != 242 && c.BG != bottom {
					t.Fatal("curve erased underlying half")
				}
			}
		}
		if partial == 0 {
			t.Fatal("test did not exercise curved half-cells")
		}
	}
}
