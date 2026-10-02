package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestDragonMockupResponsiveComposition(t *testing.T) {
	for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {120, 40}, {182, 58}} {
		w, h := size[0], size[1]
		f := RenderDragonMockup(w, h)
		if !reflect.DeepEqual(f, RenderDragonMockup(w, h)) {
			t.Fatalf("%v: study is not static", size)
		}
		if !strings.Contains(f.Text(), "MALGRATH") || !strings.Contains(f.Text(), "esc / enter return") {
			t.Fatalf("%v: missing title or exit hint", size)
		}
		counts := map[int]int{}
		for y := 4; y < h-3; y++ {
			for x := 1; x < w-1; x++ {
				c := f.C[y*w+x]
				if c.R == '█' {
					counts[c.FG]++
				}
			}
		}
		for _, fg := range []int{131, 95, 180, 130} {
			if counts[fg] == 0 {
				t.Fatalf("%v: missing dragon/wing/horn/hoard mass %d", size, fg)
			}
		}
		for y := 0; y < h; y++ {
			if f.C[y*w].R != ' ' || f.C[y*w+w-1].R != ' ' {
				t.Fatalf("%v: composition escaped framing on row %d", size, y)
			}
		}
	}
}

func TestDragonMockupSmallWindows(t *testing.T) {
	for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
		f := RenderDragonMockup(size[0], size[1])
		if f.W != size[0] || f.H != size[1] || len(f.C) != size[0]*size[1] {
			t.Fatalf("bad dimensions for %v", size)
		}
	}
}

func TestRotundaAdvertisesPortrait(t *testing.T) {
	st := NewOWState()
	st.Tokens, st.BonusGold, st.BonusLives, st.BonusTower = 3, 200, 1, true
	f := RenderOverworld(62, 19, st, 30, Palette())
	if !strings.Contains(f.Text(), "v view Malgrath") {
		t.Fatal("missing Rotunda portrait hint")
	}
	if f.C[(f.H-2)*f.W+f.W-1].R != '│' {
		t.Fatal("context hint overwrote frame")
	}
}

// Sparse line details must retain their supporting mass in the cell background.
func TestDragonPortraitDetailsStayOnTheirSurface(t *testing.T) {
	f := RenderDragonMockup(94, 47)
	eye, claws, gold := 0, 0, 0
	for _, c := range f.C {
		if c.R == '━' && c.FG == 180 {
			eye++
			if c.BG != 131 {
				t.Fatal("eye punched a hole through the head")
			}
		}
		if c.R == '╲' && c.FG == 180 {
			claws++
			if c.BG != 131 {
				t.Fatal("claw punched a hole through the paw")
			}
		}
		if c.R == '◆' && (c.FG == 178 || c.FG == 220) {
			gold++
			if c.BG != 130 {
				t.Fatal("facet punched a hole through the hoard")
			}
		}
	}
	if eye < 2 || claws != 3 || gold == 0 {
		t.Fatalf("missing portrait details: eye=%d claws=%d gold=%d", eye, claws, gold)
	}
}
