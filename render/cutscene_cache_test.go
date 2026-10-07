package render

import (
	"reflect"
	"testing"
)

func TestEarthCachedPanMatchesFreshArtworkAfterResize(t *testing.T) {
	for _, size := range [][2]int{{60, 18}, {120, 35}, {60, 18}} {
		w, h := size[0], size[1]
		far, near := filmFrame(w*2, h), filmTransparentFrame(w*2, h)
		filmEarthBackdrop(far)
		filmEarthForeground(near)
		for _, frame := range []int{0, 90, 270} {
			got, want := filmFrame(w, h), filmFrame(w, h)
			stage := Rect{0, 0, w, h}
			paintEndingPanorama(got, stage, "ending-earth", frame)
			x := filmPanColumn("ending-earth", frame, w)
			filmCompositePan(want, stage, far, near, x/4, x)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("cached Earth differs at %v, frame %d", size, frame)
			}
		}
	}
}
