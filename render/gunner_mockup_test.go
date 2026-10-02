package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestGunnerPortraitComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderGunnerMockup(w, h)
		if !reflect.DeepEqual(f, RenderGunnerMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "ORC GUNNER") || !strings.Contains(f.Text(), "esc return") {
			t.Fatalf("%v: missing title or controls", size)
		}
		counts := map[int]int{}
		for _, c := range f.C {
			if c.R == '█' {
				counts[c.FG]++
			}
		}
		for _, fg := range []int{71, 65, 94, 180, 237, 242, 60} {
			if counts[fg] == 0 {
				t.Fatalf("%v: lost skin/wood/brass/iron/legs mass %d", size, fg)
			}
		}
		for y := 0; y < h; y++ {
			if f.C[y*w].R != ' ' || f.C[y*w+w-1].R != ' ' {
				t.Fatalf("%v: portrait escaped framing", size)
			}
		}
		// An open gap between the planted legs keeps the braced stance clear.
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		x0, y0 := (w-pw)/2, 4+(h-7-ph)/2
		x, y := x0+int(.40*float64(pw-1)), y0+int(.86*float64(ph-1))
		if f.C[y*w+x].R != ' ' {
			t.Fatalf("%v: planted legs merged", size)
		}
	}
}

func TestGunnerFaceAndHandsStaySolid(t *testing.T) {
	f := RenderGunnerMockup(94, 47)
	eye, fingers := 0, 0
	for _, c := range f.C {
		if c.R == '━' && c.FG == 230 {
			eye++
			if c.BG == 233 {
				t.Fatal("eye erased face")
			}
		}
		if c.R == '│' && c.FG == 65 {
			fingers++
			if c.BG != 71 && c.BG != 114 {
				t.Fatal("finger lines erased their grip")
			}
		}
	}
	if eye != 1 || fingers < 2 {
		t.Fatalf("lost squint or grips: eye=%d fingers=%d", eye, fingers)
	}
}

func TestGunnerSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderGunnerMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad dimensions %v", size)
		}
	}
}
