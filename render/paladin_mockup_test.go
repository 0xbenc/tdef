package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestPaladinPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderPaladinMockup(w, h)
		if !reflect.DeepEqual(f, RenderPaladinMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "PALADIN") || !strings.Contains(f.Text(), "esc return") {
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
			{"closed helm", .38, .09, .58, .35, []int{109, 239, 246, 251}},
			{"dark visor", .40, .18, .56, .27, []int{234}},
			{"pointed shield", .15, .37, .57, .92, []int{109, 144, 223, 230}},
			{"guild seal", .28, .48, .44, .83, []int{234, 236}},
			{"flanged mace", .69, .04, .90, .27, []int{109, 239, 246, 251}},
			{"planted armor", .25, .84, .78, .95, []int{239, 246, 230}},
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

func TestPaladinSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderPaladinMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestPaladinFingerSeamsKeepTheirPlate(t *testing.T) {
	f := RenderPaladinMockup(94, 47)
	shieldGrip := 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 109 {
			shieldGrip++
			if c.BG != 251 && c.BG != 230 {
				t.Fatal("shield finger seam erased its plate")
			}
		}
	}
	if shieldGrip < 1 {
		t.Fatal("lost bracing shield grip")
	}
}
