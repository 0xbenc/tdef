package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderDragonMockup is a static composition for reviewing the enlarged
// Rotunda dragon. It has no cinematic timeline or simulation state.
func RenderDragonMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.dragon.enlarge_to_view_malgrath"), 180, 233, false)
	} else {
		// Preserve the room's proportions, but use nearly the whole screen.
		ph := min(h-7, (w-6)*7/15)
		pw := ph * 15 / 7
		x, y := (w-pw)/2, 4+(h-7-ph)/2
		drawDragonTableau(f, x, y, pw, ph, max(1, ph/7), 233,
			owPadView{status: OWOpen}, 90, true)
		putString(f, (w-8)/2, 1, copytext.Text("characters.dragon.malgrath"), 180, 233, true)
		putString(f, (w-16)/2, 2, copytext.Text("characters.dragon.the_dying_dragon"), 240, 233, false)
	}
	hint := copytext.Format("characters.dragon.esc_enter_return_q_quit", "enter", "enter", "escape", "esc", "quit", "q")
	if w >= 42 {
		hint = copytext.Format("characters.dragon.tab_grak", "tab", "tab") + hint
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}
