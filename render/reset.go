package render

var resetProgressItems = []string{"Cancel", "Reset all progress"}

func ResetProgressRects(w, h int) []Rect {
	y := screenOff(h) + 10
	const width = 21
	x := (w - width) / 2
	return []Rect{{X: x, Y: y, W: width, H: 1}, {X: x, Y: y + 2, W: width, H: 1}}
}

func RenderResetProgress(w, h, sel int, done bool, errMsg string, pal Colors) *Frame {
	footer := []fseg{{key: "↑↓", text: " move"}, {key: "enter", text: " select"}, {key: "esc", text: " cancel"}}
	if done {
		footer = []fseg{{key: "enter", text: " main menu"}}
	}
	f := screenBox(w, h, "RESET PROGRESS", footer, true, pal)
	off := screenOff(h)
	if done {
		centerPut(f, off+7, "All progress has been reset.", pal.Bright, true)
		centerPut(f, off+9, "Your next Start begins a new journey.", pal.Dim, false)
		return f
	}
	centerPut(f, off+4, "Reset all progress?", pal.Bright, true)
	centerPut(f, off+6, "Clears floor progress, relics, and high scores.", pal.Dim, false)
	centerPut(f, off+7, "Also clears journal discoveries and story progress.", pal.Dim, false)
	centerPut(f, off+8, "This cannot be undone.", pal.Bright, true)
	for i, rect := range ResetProgressRects(w, h) {
		prefix, fg := "   ", pal.Dim
		if i == sel {
			prefix, fg = " ▸ ", pal.Bright
		}
		centerPut(f, rect.Y, prefix+resetProgressItems[i], fg, i == sel)
	}
	if errMsg != "" {
		centerPut(f, off+14, "Could not reset all progress. Try again.", pal.Bright, true)
		centerPut(f, off+15, fitMsg(errMsg, w-4), pal.Dim, false)
	}
	return f
}
