package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestNecromancerPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderNecromancerMockup(w, h)
		if !reflect.DeepEqual(f, RenderNecromancerMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "NECROMANCER") || !strings.Contains(f.Text(), "esc return") {
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
			{"staff hook", .06, .02, .35, .25, []int{144, 180, 230}},
			{"hanging skull", .21, .19, .31, .32, []int{144, 230}},
			{"human face", .40, .22, .57, .40, []int{108, 144, 230}},
			{"summoning hand", .66, .45, .85, .61, []int{108, 144, 230}},
			{"kneeling squire", .22, .64, .44, .94, []int{109, 239, 246}},
			{"standing squire", .58, .62, .81, .94, []int{109, 239, 246}},
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

func TestNecromancerSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderNecromancerMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestNecromancerFeaturesRetainTheirSurfaces(t *testing.T) {
	f := RenderNecromancerMockup(94, 47)
	grips := 0
	for _, c := range f.C {
		if c.R == '│' && c.FG == 108 {
			grips++
			if c.BG != 108 && c.BG != 144 && c.BG != 230 {
				t.Fatal("finger mark erased the hand")
			}
		}
	}
	if grips < 1 {
		t.Fatal("lost the skeletal finger marks")
	}
	// The facial eye must survive half-cell sampling in the dark socket.
	ph := min(47-7, (94-6)*2/5)
	pw := ph * 5 / 2
	x0, y0 := (94-pw)/2, 4+(47-7-ph)/2
	eye := false
	for y := int(.27 * float64(ph-1)); y <= int(.32*float64(ph-1))+1; y++ {
		for x := int(.47 * float64(pw-1)); x <= int(.53*float64(pw-1)); x++ {
			c := f.C[(y0+y)*94+x0+x]
			if c.R != ' ' && (c.FG == 151 || c.BG == 151) {
				eye = true
			}
		}
	}
	if !eye {
		t.Fatal("the eye disappeared into its socket")
	}
}
