package render

import "github.com/0xbenc/termtd/internal/copytext"

var resetProgressItems = []string{copytext.Text("ui.reset_progress_items.cancel"), copytext.Text("ui.reset_progress_items.reset_all_progress")}

func ResetProgressRects(w, h int) []Rect {
	y := screenOff(h) + 10
	width := 3
	for _, item := range resetProgressItems {
		width = max(width, 3+len([]rune(item)))
	}
	x := (w - width) / 2
	return []Rect{{X: x, Y: y, W: width, H: 1}, {X: x, Y: y + 2, W: width, H: 1}}
}

func RenderResetProgress(w, h, sel int, done bool, errMsg string, pal Colors) *Frame {
	footer := []fseg{{key: "↑↓", text: copytext.Text("ui.render_reset_progress.move")}, {key: "enter", text: copytext.Text("ui.render_reset_progress.select")}, {key: "esc", text: copytext.Text("ui.render_reset_progress.cancel")}}
	if done {
		footer = []fseg{{key: "enter", text: copytext.Text("ui.render_reset_progress.main_menu")}}
	}
	f := screenBox(w, h, copytext.Text("ui.render_reset_progress.reset_progress"), footer, true, pal)
	off := screenOff(h)
	if done {
		centerPut(f, off+7, copytext.Text("ui.render_reset_progress.all_progress_has_been_reset"), pal.Bright, true)
		centerPut(f, off+9, copytext.Text("ui.render_reset_progress.your_next_start_begins_a_new_journey"), pal.Dim, false)
		return f
	}
	centerPut(f, off+4, copytext.Text("ui.render_reset_progress.reset_all_progress"), pal.Bright, true)
	centerPut(f, off+6, copytext.Text("ui.render_reset_progress.clears_floor_progress_relics_and_high_scores"), pal.Dim, false)
	centerPut(f, off+7, copytext.Text("ui.render_reset_progress.also_clears_journal_discoveries_and_story_progress"), pal.Dim, false)
	centerPut(f, off+8, copytext.Text("ui.render_reset_progress.this_cannot_be_undone"), pal.Bright, true)
	for i, rect := range ResetProgressRects(w, h) {
		prefix, fg := "   ", pal.Dim
		if i == sel {
			prefix, fg = " ▸ ", pal.Bright
		}
		centerPut(f, rect.Y, prefix+resetProgressItems[i], fg, i == sel)
	}
	if errMsg != "" {
		centerPut(f, off+14, copytext.Text("ui.render_reset_progress.could_not_reset_all_progress_try_again"), pal.Bright, true)
		centerPut(f, off+15, fitMsg(errMsg, w-4), pal.Dim, false)
	}
	return f
}
