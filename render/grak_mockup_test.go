package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestGrakMockupComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}} {
		w, h := size[0], size[1]
		f := RenderGrakMockup(w, h)
		if !reflect.DeepEqual(f, RenderGrakMockup(w, h)) {
			t.Fatalf("%v: portrait is not static", size)
		}
		if !strings.Contains(f.Text(), "GRAK") || !strings.Contains(f.Text(), "the last monster") || !strings.Contains(f.Text(), "esc return") {
			t.Fatalf("%v: missing title or controls", size)
		}
		counts := map[int]int{}
		for _, c := range f.C {
			if c.R == '█' {
				counts[c.FG]++
			}
		}
		for _, fg := range []int{107, 65, 95, 94, 238} {
			if counts[fg] == 0 {
				t.Fatalf("%v: missing skin/mantle/boots/mallet mass %d", size, fg)
			}
		}
		for y := 0; y < h; y++ {
			if f.C[y*w].R != ' ' || f.C[y*w+w-1].R != ' ' {
				t.Fatalf("%v: portrait escaped frame", size)
			}
		}
	}
}

func TestGrakMockupSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderGrakMockup(size[0], size[1])
		if len(f.C) != size[0]*size[1] {
			t.Fatalf("bad frame size %v", size)
		}
	}
}

func TestGrakDetailsRetainSkin(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {94, 47}} {
		f := RenderGrakMockup(size[0], size[1])
		eyes := 0
		for _, c := range f.C {
			if c.R == '━' && c.FG == 180 {
				eyes++
				if c.BG == 233 {
					t.Fatal("eye erased its supporting face")
				}
			}
		}
		if eyes == 0 {
			t.Fatalf("%v: lost both eyes", size)
		}
	}
}
