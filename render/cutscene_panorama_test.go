package render

import (
	"reflect"
	"testing"
)

func TestFilmPanRevealsWholeColumnsAndSettles(t *testing.T) {
	for _, art := range []string{"village", "village-search", "ruined-road", "capital-vault", "ending-pursuit", "ending-depths", "ending-earth"} {
		for _, travel := range []int{18, 60, 120} {
			first := filmPanColumn(art, 0, travel)
			last := filmPanColumn(art, filmPanEnd, travel)
			if abs(first-last) != travel {
				t.Fatalf("%s only traveled %d of %d columns", art, abs(first-last), travel)
			}
			prev := first
			for frame := 0; frame <= filmPanEnd+30; frame++ {
				x := filmPanColumn(art, frame, travel)
				if x < 0 || x > travel || (art == "village-search" && x > prev) || (art != "village-search" && x < prev) {
					t.Fatalf("%s camera reversed or left the scene at frame %d", art, frame)
				}
				if frame <= filmPanStart && x != first || frame >= filmPanEnd && x != last {
					t.Fatalf("%s camera did not hold at frame %d", art, frame)
				}
				prev = x
			}
		}
		if art != "capital-vault" {
			continue // Independently moving layers are checked below.
		}
		// A pan translates already rasterized art; overlapping columns remain
		// byte-for-byte identical rather than crawling through sample points.
		a, b := filmFrame(60, 18), filmFrame(60, 18)
		paintFilmPanorama(a, Rect{0, 0, 60, 18}, art, 60)
		paintFilmPanorama(b, Rect{0, 0, 60, 18}, art, 90)
		shift := filmPanColumn(art, 90, 60) - filmPanColumn(art, 60, 60)
		if reflect.DeepEqual(a, b) {
			t.Fatalf("%s pan revealed no new artwork", art)
		}
		for y := 0; y < a.H; y++ {
			for x := 0; x < a.W; x++ {
				bx := x - shift
				if bx >= 0 && bx < b.W && a.C[y*a.W+x] != b.C[y*b.W+bx] {
					t.Fatalf("%s reshaped a moving cell at %d,%d", art, x, y)
				}
			}
		}
	}
}

func TestFilmsHoldStillAfterStagedMotion(t *testing.T) {
	for film := FilmOpening; film <= FilmEnding; film++ {
		for i, shot := range FilmShots(film) {
			if filmIsOutdoor(shot.art) {
				continue // Smoke stays alive while the camera holds.
			}
			a := RenderCutscene(80, 24, CutsceneState{Film: film, Shot: i, Frame: filmPanEnd + filmPanHold, Revealed: true})
			b := RenderCutscene(80, 24, CutsceneState{Film: film, Shot: i, Frame: 360, Revealed: true})
			if !reflect.DeepEqual(a, b) {
				t.Errorf("film %d shot %d keeps shimmering after settling", film, i)
			}
		}
	}
}

func TestPanoramaHoldsBeforeAdvance(t *testing.T) {
	for i, shot := range FilmShots(FilmOpening) {
		st := CutsceneState{Film: FilmOpening, Shot: i, Frame: filmPanEnd + filmPanHold - 1}
		if CutsceneCanAdvance(st) == filmIsPanorama(shot.art) {
			t.Errorf("shot %d has incorrect advance timing", i)
		}
		st.Frame++
		if !CutsceneCanAdvance(st) {
			t.Errorf("shot %d cannot advance after the end hold", i)
		}
	}
}
