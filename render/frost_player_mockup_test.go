package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestFrostAndPlayerPortraitsResponsive(t *testing.T) {
	for _, portrait := range []struct {
		name, title string
		draw        func(int, int) *Frame
		colors      []int
	}{
		{"frost", "FROST MAGE", RenderFrostMockup, []int{60, 67, 109, 153}},
		{"player", "THE PLAYER", RenderPlayerMockup, []int{251, 110, 180, 95, 130}},
	} {
		t.Run(portrait.name, func(t *testing.T) {
			for _, size := range [][2]int{{32, 16}, {62, 19}, {80, 24}, {94, 47}, {120, 40}, {182, 58}} {
				w, h := size[0], size[1]
				f := portrait.draw(w, h)
				if !reflect.DeepEqual(f, portrait.draw(w, h)) {
					t.Fatalf("%v: portrait is not static", size)
				}
				if !strings.Contains(f.Text(), portrait.title) || !strings.Contains(f.Text(), "esc return") {
					t.Fatalf("%v: missing title or controls", size)
				}
				colors := map[int]int{}
				for _, c := range f.C {
					if c.R == '█' {
						colors[c.FG]++
					}
				}
				for _, fg := range portrait.colors {
					if colors[fg] == 0 {
						t.Fatalf("%v: lost defining mass %d", size, fg)
					}
				}
				for y := 0; y < h; y++ {
					if f.C[y*w].R != ' ' || f.C[y*w+w-1].R != ' ' {
						t.Fatalf("%v: illustration escaped its framing", size)
					}
				}
			}
			for _, size := range [][2]int{{0, 0}, {1, 1}, {20, 10}, {80, 15}} {
				f := portrait.draw(size[0], size[1])
				if len(f.C) != size[0]*size[1] {
					t.Fatalf("bad small frame: %v", size)
				}
			}
		})
	}
}

func TestFrostHandSurvivesGlowAndShard(t *testing.T) {
	f := RenderFrostMockup(94, 47)
	ph := min(47-7, (94-6)/2)
	pw := ph * 2
	x0, y0 := (94-pw)/2, 4+(47-7-ph)/2
	hand, crystal := 0, 0
	for y := 0; y < ph; y++ {
		for x := 0; x < pw; x++ {
			u, v := float64(x)/float64(pw-1), float64(y)/float64(ph-1)
			c := f.C[(y0+y)*94+x0+x]
			if u >= .68 && u <= .82 && v >= .53 && v <= .60 && c.R == '█' && (c.FG == 152 || c.FG == 109) {
				hand++
			}
			if u >= .69 && u <= .81 && v >= .35 && v <= .52 && c.R == '█' && (c.FG == 153 || c.FG == 189) {
				crystal++
			}
		}
	}
	if hand < 3 || crystal < 3 {
		t.Fatalf("lost open hand or floating ice: hand=%d crystal=%d", hand, crystal)
	}
	for _, c := range f.C {
		if c.R == '━' && c.FG == 255 && c.BG == 233 {
			t.Fatal("mage eye erased its supporting face")
		}
	}
}

func TestPlayerBoneSocketsRetainIvory(t *testing.T) {
	f := RenderPlayerMockup(94, 47)
	sockets := 0
	for _, c := range f.C {
		if c.R == '▀' && c.FG == 232 {
			sockets++
			if c.BG != 144 && c.BG != 230 && c.BG != 255 {
				t.Fatalf("bone socket lost its ivory surface: %+v", c)
			}
		}
	}
	if sockets != 2 {
		t.Fatalf("lost skull's empty sockets: %d", sockets)
	}
}
