package render

import (
	"github.com/0xbenc/termtd/internal/copytext"
	"strings"
)

// DrawSaveFailure stays visible across screens until the pending save succeeds.
func DrawSaveFailure(f *Frame, detail string, quitting bool) {
	message := copytext.Text("ui.save.failed")
	if quitting {
		message = copytext.Text("ui.save.quit_warning")
	}
	drawPersistenceNotice(f, message, detail)
}

// DrawLoadFailure explains why saving is blocked without implying retries are safe.
func DrawLoadFailure(f *Frame, detail string, expanded, quitting bool, top int) {
	var names []string
	for _, failure := range strings.Split(detail, "\n") {
		name, _, _ := strings.Cut(failure, ":")
		names = append(names, name)
	}
	lines := []string{copytext.Text("ui.save.load_failed"),
		copytext.Format("ui.save.load_blocked", "names", strings.Join(names, ", ")),
		copytext.Text("ui.save.load_kept"),
		copytext.Text("ui.save.load_repair")}

	if expanded {
		lines = nil
		width := max(1, f.W-4)
		for _, paragraph := range strings.Split(detail, "\n") {
			runes := []rune(paragraph)
			for len(runes) > 0 {
				n := min(width, len(runes))
				lines = append(lines, string(runes[:n]))
				runes = runes[n:]
			}
		}
		visible := max(1, f.H-4)
		top = min(max(0, top), max(0, len(lines)-visible))
		lines = append([]string{copytext.Text("ui.save.load_details")}, lines[top:min(len(lines), top+visible)]...)
	}
	if quitting {
		lines = append(lines, copytext.Text("ui.save.load_quit"))
	}
	if len(lines) > max(0, f.H-2) {
		lines = lines[:max(0, f.H-2)]
	}
	for i, text := range lines {
		y := f.H - 1 - len(lines) + i
		for x := 1; x < f.W-1; x++ {
			f.Set(x, y, Cell{R: ' ', FG: 255, BG: 52})
		}
		putString(f, 2, y, fitMsg(text, max(0, f.W-4)), 255, 52, i == 0)
	}
}

func drawPersistenceNotice(f *Frame, message, detail string) {
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
