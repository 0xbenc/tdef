package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestCenturionPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderCenturionMockup(w, h)
		if !reflect.DeepEqual(f, RenderCenturionMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "CENTURION") || !strings.Contains(f.Text(), "esc return") {
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
			{"transverse horsehair crest", .3, .02, .69, .16, []int{52, 131, 167}},
			{"human helmet opening", .41, .24, .59, .35, []int{137, 223, 94}},
			{"rectangular shield", .41, .36, .78, .9, []int{137, 180, 52, 95}},
			{"segmented iron armor", .27, .34, .44, .55, []int{239, 246, 109}},
			{"short sword", .18, .61, .31, .9, []int{239, 251, 109}},
			{"magical ward", .78, .35, .86, .86, []int{117, 24, 153}},
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

func TestCenturionSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderCenturionMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad small window %v", size)
		}
	}
}

func TestCenturionFaceAndWardGap(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}} {
		f := RenderCenturionMockup(size[0], size[1])
		eyes := 0
		for _, c := range f.C {
			if c.R == '━' && c.FG == 234 {
				eyes++
				if c.BG != 223 && c.BG != 137 {
					t.Fatal("eye escaped face")
				}
			}
		}
		want := 1
		if size[1] >= 40 {
			want = 2
		}
		if eyes != want {
			t.Fatalf("%v: lost face opening (%d)", size, eyes)
		}
	}
	f := RenderCenturionMockup(94, 47)
	ph := min(47-7, (94-6)*2/5)
	pw := ph * 5 / 2
	x0, y0 := (94-pw)/2, 4+(47-7-ph)/2
	gap := f.C[(y0+int(.67*float64(ph-1)+.5))*94+x0+int(.816*float64(pw-1)+.5)]
	if gap.R != ' ' || gap.BG != 233 {
		t.Fatal("ward arch lost its air gap")
	}
}
