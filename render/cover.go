package render

// RenderTerminalLogo returns the actual title-screen lettering on the terminal's
// default background, allowing a transparent terminal capture.
func RenderTerminalLogo() *Frame {
	f := &Frame{W: 44, H: 6, C: make([]Cell, 44*6)}
	drawTitleLogo(f, f.W, -2, -1)
	return f
}

// RenderPromo adapts the cover's game art to wide banner and social layouts.
func RenderPromo(w, h int) *Frame {
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 245, BG: 233}
	}
	drawRoundedBox(f, 0, 0, w, h, 238)
	embedSegment(f, 0, 2, "./tdef", '┐', '┌', 238, 180, false)
	leftW := w/2 - 4
	left := &Frame{W: leftW, H: h, C: make([]Cell, leftW*h)}
	drawTitleLogo(left, leftW, 3, -1)
	drawTitleTagline(left, leftW, 3, len([]rune(titleTagline)))
	centerPut(left, 14, "THE LAST MONSTER", 180, true)
	for y := 0; y < h; y++ {
		for x := 0; x < leftW; x++ {
			c := left.C[y*leftW+x]
			if c.R != 0 {
				c.BG = 233
				f.Set(3+x, y, c)
			}
		}
	}
	x := w / 2
	drawDragonTableau(f, x, 2, w-x-3, h-5, 2, 233,
		owPadView{status: OWOpen}, 90, true)
	return f
}

// RenderCover composes the game's title lettering and Malgrath tableau for
// the itch.io cover. Every mark is a normal terminal cell, using the same
// drawing routines and palette as the title screen and dragon portrait.
func RenderCover() *Frame {
	const w, h = 66, 25
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 245, BG: 233}
	}
	drawRoundedBox(f, 0, 0, w, h, 238)
	embedSegment(f, 0, 2, "./tdef", '┐', '┌', 238, 180, false)
	drawTitleLogo(f, w, 0, -1)
	drawTitleTagline(f, w, 0, len([]rune(titleTagline)))
	drawDragonTableau(f, 4, 10, 58, 13, 2, 233,
		owPadView{status: OWOpen}, 90, true)
	centerPut(f, 23, "THE LAST MONSTER", 180, true)
	// Helpers use the terminal's default background; make it explicit so the
	// capture looks identical regardless of the viewer's terminal theme.
	for i := range f.C {
		if f.C[i].BG == 0 {
			f.C[i].BG = 233
		}
	}
	return f
}
