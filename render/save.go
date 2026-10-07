package render

import "github.com/0xbenc/termtd/internal/copytext"

// DrawSaveFailure stays visible across screens until the pending save succeeds.
func DrawSaveFailure(f *Frame, detail string, quitting bool) {
	message := copytext.Text("ui.save.failed")
	if quitting {
		message = copytext.Text("ui.save.quit_warning")
	}
	for i, text := range []string{message, detail} {
		y := f.H - 3 + i
		if y < 0 {
			continue
		}
		for x := 1; x < f.W-1; x++ {
			f.Set(x, y, Cell{R: ' ', FG: 255, BG: 52})
		}
		putString(f, 2, y, fitMsg(text, max(0, f.W-4)), 255, 52, i == 0)
	}
}
