package render

import (
	"strings"
	"testing"
)

func TestResetProgressShowsConfirmationAndResult(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		w, h := size[0], size[1]
		f := RenderResetProgress(w, h, 0, false, "", Palette())
		text := f.Text()
		for _, want := range []string{"Reset all progress?", "high scores", "journal discoveries", "cannot be undone", "▸ Cancel"} {
			if !strings.Contains(text, want) {
				t.Fatalf("%dx%d: confirmation missing %q", w, h, want)
			}
		}
		rows := strings.Split(text, "\n")
		for i, r := range ResetProgressRects(w, h) {
			if !strings.Contains(rows[r.Y], resetProgressItems[i]) {
				t.Fatal("confirmation hit region does not match rendered choice")
			}
		}
		text = RenderResetProgress(w, h, 0, true, "", Palette()).Text()
		if !strings.Contains(text, "All progress has been reset.") || strings.Contains(text, "Reset all progress?") {
			t.Fatal("success must replace confirmation")
		}
		text = RenderResetProgress(w, h, 0, false, "permission denied", Palette()).Text()
		if !strings.Contains(text, "Could not reset all progress") || !strings.Contains(text, "permission denied") {
			t.Fatal("save removal error must be visible")
		}
	}
}
