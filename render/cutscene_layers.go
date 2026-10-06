package render

func filmIsOutdoor(art string) bool {
	return art == "village" || art == "village-search" || art == "ruined-road"
}

// Negative colors are internal transparency, resolved before reaching the
// terminal. Half-cell edges preserve the hills beneath smoke and round shapes.
func filmTransparentFrame(w, h int) *Frame {
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: -1, BG: -1}
	}
	return f
}

func filmCellHalves(c Cell) (int, int) {
	switch c.R {
	case '█':
		return c.FG, c.FG
	case '▀':
		return c.FG, c.BG
	case '▄':
		return c.BG, c.FG
	}
	return c.BG, c.BG
}

func filmLayerOver(back, front Cell) Cell {
	if front.R == ' ' {
		return back
	}
	if front.R != '█' && front.R != '▀' && front.R != '▄' {
		if front.BG < 0 {
			_, front.BG = filmCellHalves(back)
		}
		return front
	}
	top, bottom := filmCellHalves(front)
	bt, bb := filmCellHalves(back)
	if top < 0 {
		top = bt
	}
	if bottom < 0 {
		bottom = bb
	}
	if top == bottom {
		return Cell{R: '█', FG: top, BG: 233}
	}
	return Cell{R: '▀', FG: top, BG: bottom}
}

func filmCompositePan(dst *Frame, stage Rect, far, near *Frame, farX, nearX int) {
	for y := 0; y < stage.H; y++ {
		for x := 0; x < stage.W; x++ {
			back := far.C[y*far.W+farX+x]
			front := near.C[y*near.W+nearX+x]
			dst.Set(stage.X+x, stage.Y+y, filmLayerOver(back, front))
		}
	}
}

func filmOutdoorBackdrop(f *Frame, art string) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.poly(235, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	// The sun and one broad ridge share a single distant layer. Strong peaks
	// make the slower travel visible without adding another busy silhouette.
	p.smoothOval(131, .39, .18, .027, .12)
	if art == "ruined-road" {
		p.poly(237, portraitPoint{0, .66}, portraitPoint{.12, .48}, portraitPoint{.23, .27}, portraitPoint{.37, .52}, portraitPoint{.5, .32}, portraitPoint{.65, .55}, portraitPoint{.83, .22}, portraitPoint{1, .57}, portraitPoint{1, 1}, portraitPoint{0, 1})
	} else {
		p.poly(237, portraitPoint{0, .58}, portraitPoint{.1, .43}, portraitPoint{.22, .24}, portraitPoint{.38, .49}, portraitPoint{.53, .29}, portraitPoint{.7, .51}, portraitPoint{.88, .21}, portraitPoint{1, .5}, portraitPoint{1, 1}, portraitPoint{0, 1})
	}
}

func filmSmoke(p portraitPainter, u, v float64, frame, seed int) {
	// Quiet, continuous plumes, staggered across ruins. Five updates per second
	// keep the secondary motion distinct from the deliberate camera travel.
	clock := max(0, frame) / 6 * 6
	// A narrow continuous stem keeps the rising curls reading as smoke,
	// rather than detached particles. Only a few roofs carry a plume.
	p.poly(238, portraitPoint{u - .004, v}, portraitPoint{u + .003, v - .16}, portraitPoint{u + .01, v - .21}, portraitPoint{u + .009, v - .08}, portraitPoint{u + .005, v})
	for i := 0; i < 3; i++ {
		phase := (clock + seed*31 + i*60) % 180
		t := float64(phase) / 180
		color := 238
		if t < .35 {
			color = 239
		}
		p.smoothOval(color, u+.018*t, v-.34*t, .007+.01*t, .055+.02*t)
	}
}
