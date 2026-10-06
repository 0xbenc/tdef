package render

import "testing"

func TestFilmParallaxMovesLandmarksAtDifferentSpeeds(t *testing.T) {
	far := filmFrame(240, 24)
	near := filmTransparentFrame(240, 24)
	far.Put(60, 3, '█', 131, 235)
	near.Put(100, 18, '█', 214, 233)
	a, b := filmFrame(120, 24), filmFrame(120, 24)
	filmCompositePan(a, Rect{0, 0, 120, 24}, far, near, 0, 0)
	filmCompositePan(b, Rect{0, 0, 120, 24}, far, near, 20, 80)
	if a.C[3*a.W+60].FG != 131 || b.C[3*b.W+40].FG != 131 {
		t.Fatal("distant landmark did not move 20 columns")
	}
	if a.C[18*a.W+100].FG != 214 || b.C[18*b.W+20].FG != 214 {
		t.Fatal("foreground landmark did not move 80 columns")
	}
	for _, art := range []string{"village", "village-search", "ruined-road"} {
		start := filmPanColumn(art, 0, 120)
		end := filmPanColumn(art, filmPanEnd, 120)
		if abs(start/4-end/4) != 30 {
			t.Fatalf("%s background travel is not one quarter of foreground travel", art)
		}
	}
}

func TestFilmLayerEdgesPreserveBackdrop(t *testing.T) {
	back := Cell{R: '▀', FG: 131, BG: 237}
	for _, front := range []Cell{
		{R: '▀', FG: 239, BG: -1},
		{R: '▄', FG: 239, BG: -1},
		{R: ' ', FG: -1, BG: -1},
	} {
		got := filmLayerOver(back, front)
		top, bottom := filmCellHalves(got)
		if got.FG < 0 || got.BG < 0 {
			t.Fatal("internal transparency leaked to the terminal")
		}
		switch front.R {
		case '▀':
			if top != 239 || bottom != 237 {
				t.Fatal("upper smoke edge erased lower backdrop")
			}
		case '▄':
			if top != 131 || bottom != 239 {
				t.Fatal("lower smoke edge erased upper backdrop")
			}
		case ' ':
			if got != back {
				t.Fatal("empty layer obscured backdrop")
			}
		}
	}
}

func TestOutdoorSmokeMovesWhileSceneryHolds(t *testing.T) {
	for _, art := range []string{"village", "village-search", "ruined-road"} {
		for _, size := range [][2]int{{36, 9}, {60, 18}, {120, 36}} {
			a, b := filmFrame(size[0], size[1]), filmFrame(size[0], size[1])
			stage := Rect{0, 0, a.W, a.H}
			paintFilmPanorama(a, stage, art, filmPanEnd+filmPanHold)
			paintFilmPanorama(b, stage, art, filmPanEnd+filmPanHold+72)
			changed := 0
			for i, c := range a.C {
				d := b.C[i]
				if c.FG < 0 || c.BG < 0 || d.FG < 0 || d.BG < 0 {
					t.Fatalf("%s leaked transparency at %v", art, size)
				}
				if c != d {
					changed++
					if i/a.W >= a.H*3/4 {
						t.Fatalf("%s held foreground moved at %v", art, size)
					}
				}
			}
			if changed == 0 || changed > len(a.C)/6 {
				t.Fatalf("%s smoke changed %d of %d cells at %v", art, changed, len(a.C), size)
			}
		}
	}
}
