package render

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilmCopyRejectsBrokenEdits(t *testing.T) {
	for _, data := range []string{
		`[{`,
		`[]`,
		`[{"title":"Scene","speaker":"TYPO","text":"Words"}]`,
		`[{"title":"Scene","speaker":"","text":""}]`,
		`[{"title":"Scene","speaker":"","text":"Words","typo":true}]`,
		`[{"title":"Scene","speaker":"","text":"Words"}] []`,
	} {
		t.Run(data, func(t *testing.T) {
			defer func() {
				failure := recover()
				if failure == nil || !strings.Contains(fmt.Sprint(failure), "ending cinematic copy") {
					t.Fatalf("broken edit did not identify its source: %v", failure)
				}
			}()
			withFilmCopy("ending", []byte(data), []FilmShot{{art: "ending-fallen"}})
		})
	}
}

func TestEndingPanoramasAtCameraHoldsAndTravel(t *testing.T) {
	for i, shot := range FilmShots(FilmEnding) {
		if !filmIsPanorama(shot.art) {
			continue
		}
		for _, size := range [][2]int{{62, 19}, {120, 36}} {
			for _, frame := range []int{filmPanStart, (filmPanStart + filmPanEnd) / 2, filmPanEnd + filmPanHold} {
				f := RenderCutscene(size[0], size[1], CutsceneState{Film: FilmEnding, Shot: i, Frame: frame, Revealed: true})
				text := strings.Join(strings.Fields(f.Text()), " ")
				if !strings.Contains(text, shot.Dialogue) || !strings.Contains(text, "esc skip") {
					t.Fatalf("shot %d frame %d at %v obscured its copy or skip control", i, frame, size)
				}
				for _, c := range f.C {
					if c.FG < 0 || c.FG > 255 || c.BG < 0 || c.BG > 255 {
						t.Fatalf("shot %d leaked an internal layer color: %+v", i, c)
					}
				}
			}
		}
	}
}
