package tui

import (
	"testing"

	"tdef/render"
)

func TestCycleSpeedWraps(t *testing.T) {
	a := &App{ui: render.UI{Speed: 1}}
	// three ups wrap 1->2->4->1, three downs wrap 1->4->2->1
	steps := []struct {
		dir  int
		want int
	}{
		{1, 2}, {1, 4}, {1, 1},
		{-1, 4}, {-1, 2}, {-1, 1},
	}
	for i, s := range steps {
		a.cycleSpeed(s.dir)
		if a.ui.Speed != s.want {
			t.Fatalf("step %d (dir %d): speed = %d, want %d", i, s.dir, a.ui.Speed, s.want)
		}
	}
}
