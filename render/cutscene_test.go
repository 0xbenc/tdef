package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestFilmCaptionsAndEveryShotAtPlayableSizes(t *testing.T) {
	for film := FilmOpening; film <= FilmEnding; film++ {
		shots := FilmShots(film)
		if len(shots) < 12 {
			t.Fatalf("film %d has too few directed shots", film)
		}
		for i, shot := range shots {
			for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 46}, {160, 54}} {
				st := CutsceneState{Film: film, Shot: i, Frame: 120, Revealed: true}
				f := RenderCutscene(size[0], size[1], st)
				text := strings.Join(strings.Fields(f.Text()), " ")
				if !strings.Contains(text, shot.Dialogue) {
					t.Errorf("film %d shot %d at %v lost dialogue: %s", film, i, size, shot.Dialogue)
				}
				if !strings.Contains(text, FilmTitle(film)) || !strings.Contains(text, "esc skip") {
					t.Errorf("film %d shot %d at %v lost navigation", film, i, size)
				}
				// Figures must occupy the middle of the frame, away from caption chrome.
				painted := 0
				for y := 3; y < f.H-8; y++ {
					for x := 0; x < f.W; x++ {
						c := f.C[y*f.W+x]
						if (c.R == '█' || c.R == '▀' || c.R == '▄') && c.FG != 233 {
							painted++
						}
					}
				}
				if painted < 25 {
					t.Errorf("film %d shot %d at %v lost artwork (%d cells)", film, i, size, painted)
				}
			}
			st := CutsceneState{Film: film, Shot: i, Frame: 120, Revealed: true}
			if !reflect.DeepEqual(RenderCutscene(80, 24, st), RenderCutscene(80, 24, st)) {
				t.Fatal("film is nondeterministic")
			}
		}
	}
}

func TestCutsceneRevealSettlesAndResizeNotice(t *testing.T) {
	st := CutsceneState{Film: FilmOpening, Shot: 1}
	if CutsceneVisible(st) != 0 {
		t.Fatal("caption appeared before the opening hold")
	}
	st.Frame = 40
	n := CutsceneVisible(st)
	if n <= 0 || n >= len([]rune(FilmShots(st.Film)[st.Shot].Dialogue)) {
		t.Fatal("typewriter did not reveal gradually")
	}
	st.Frame = 10000
	if CutsceneVisible(st) != len([]rune(FilmShots(st.Film)[st.Shot].Dialogue)) {
		t.Fatal("caption did not settle")
	}
	for _, size := range [][2]int{{0, 0}, {20, 10}, {61, 19}, {62, 18}} {
		f := RenderCutscene(size[0], size[1], st)
		if len(f.C) != size[0]*size[1] {
			t.Fatal("bad compact frame")
		}
		if size[0] >= 40 && !strings.Contains(f.Text(), "enlarge") {
			t.Fatal("missing resize instruction")
		}
	}
}

func TestFilmsHaveDistinctCompositionsAndQuietMotion(t *testing.T) {
	// The closing hoard and the opening builder use distinct compositions.
	a := filmFrame(240, 80)
	b := filmFrame(240, 80)
	paintFilmShot(a, CutsceneState{Film: FilmOpening, Frame: 120}, "resolve")
	paintFilmShot(b, CutsceneState{Film: FilmEnding, Frame: 120}, "ending-hoard")
	if reflect.DeepEqual(a, b) {
		t.Fatal("ending reused the opening pose")
	}
	for _, art := range []string{"eye", "mallet", "ending-rescue"} {
		a, b = filmFrame(240, 80), filmFrame(240, 80)
		paintFilmShot(a, CutsceneState{Film: FilmEnding, Shot: 7, Frame: 20}, art)
		paintFilmShot(b, CutsceneState{Film: FilmEnding, Shot: 7, Frame: 95}, art)
		if reflect.DeepEqual(a, b) {
			t.Errorf("%s has no staged motion", art)
		}
	}
	// Fade is monotonic darkness to the original warm palette, never an
	// unrelated palette-index sweep through green or cyan.
	for _, col := range []int{131, 173, 180, 223} {
		if filmFade(col, 0) != 233 || filmFade(col, 1) != col {
			t.Fatalf("bad fade endpoints for %d", col)
		}
		for _, v := range []float64{.2, .5, .8} {
			r, g, b := xtermRGB(filmFade(col, v))
			if g > r || b > r {
				t.Fatalf("warm fade flashed an unrelated hue for %d", col)
			}
		}
	}
}
