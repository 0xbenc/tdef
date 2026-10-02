package render

// RenderFrostMockup is the Frost Mage's definitive lore portrait: a still,
// ancient monster holding a shard beside a crooked, forked crystal staff.
func RenderFrostMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the mage", 153, 233, false)
	} else {
		ph := min(h-7, (w-6)/2)
		pw := ph * 2
		drawFrostPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-10)/2, 1, "FROST MAGE", 153, 233, true)
		putString(f, (w-18)/2, 2, "the cold remembers", 240, 233, false)
	}
	hint := "tab Player · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawFrostPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.oval
	// Two quiet halos frame the face and staff without a field of sparkles.
	oval(234, .49, .29, .235, .245)
	oval(235, .49, .29, .195, .215)
	oval(234, .20, .18, .13, .15)
	oval(235, .20, .18, .09, .11)
	oval(234, .75, .46, .105, .105)
	oval(235, .75, .46, .075, .08)
	poly(236, portraitPoint{.10, .93}, portraitPoint{.80, .93}, portraitPoint{.85, .95}, portraitPoint{.07, .95})
	// The long staff bends twice. Its pale coating follows one lit edge.
	poly(60, portraitPoint{.205, .23}, portraitPoint{.23, .23}, portraitPoint{.25, .44}, portraitPoint{.23, .63}, portraitPoint{.245, .92}, portraitPoint{.21, .92}, portraitPoint{.20, .62}, portraitPoint{.22, .43})
	poly(109, portraitPoint{.21, .24}, portraitPoint{.222, .24}, portraitPoint{.241, .44}, portraitPoint{.216, .63}, portraitPoint{.228, .91}, portraitPoint{.215, .91}, portraitPoint{.20, .62}, portraitPoint{.224, .43})
	// A fork, a broken tine, and an upright spear crystal form one crown.
	poly(153, portraitPoint{.20, .27}, portraitPoint{.16, .23}, portraitPoint{.12, .17}, portraitPoint{.12, .08}, portraitPoint{.15, .055}, portraitPoint{.15, .15}, portraitPoint{.18, .19}, portraitPoint{.21, .205}, portraitPoint{.245, .17}, portraitPoint{.275, .10}, portraitPoint{.31, .08}, portraitPoint{.32, .115}, portraitPoint{.29, .14}, portraitPoint{.26, .21}, portraitPoint{.23, .26})
	poly(110, portraitPoint{.12, .17}, portraitPoint{.15, .15}, portraitPoint{.18, .19}, portraitPoint{.21, .205}, portraitPoint{.21, .26}, portraitPoint{.16, .23})
	poly(189, portraitPoint{.145, .065}, portraitPoint{.15, .055}, portraitPoint{.15, .15}, portraitPoint{.18, .19}, portraitPoint{.165, .19}, portraitPoint{.14, .155})
	poly(110, portraitPoint{.19, .16}, portraitPoint{.225, .04}, portraitPoint{.25, .16}, portraitPoint{.225, .235})
	poly(153, portraitPoint{.225, .04}, portraitPoint{.25, .16}, portraitPoint{.225, .235}, portraitPoint{.225, .15})
	poly(255, portraitPoint{.225, .04}, portraitPoint{.225, .15}, portraitPoint{.21, .19}, portraitPoint{.20, .155})
	// A wide hem, narrow shoulders, and a sloping hood make a dark triangle.
	poly(236, portraitPoint{.41, .33}, portraitPoint{.60, .33}, portraitPoint{.65, .48}, portraitPoint{.75, .88}, portraitPoint{.70, .925}, portraitPoint{.28, .925}, portraitPoint{.17, .90}, portraitPoint{.29, .55})
	poly(60, portraitPoint{.39, .34}, portraitPoint{.55, .34}, portraitPoint{.61, .48}, portraitPoint{.69, .89}, portraitPoint{.55, .93}, portraitPoint{.26, .92}, portraitPoint{.19, .88}, portraitPoint{.29, .55})
	poly(236, portraitPoint{.38, .46}, portraitPoint{.43, .51}, portraitPoint{.37, .90}, portraitPoint{.27, .92}, portraitPoint{.26, .84})
	poly(67, portraitPoint{.48, .45}, portraitPoint{.55, .44}, portraitPoint{.60, .55}, portraitPoint{.69, .89}, portraitPoint{.59, .92}, portraitPoint{.48, .91}, portraitPoint{.51, .69})
	poly(60, portraitPoint{.51, .58}, portraitPoint{.54, .57}, portraitPoint{.59, .88}, portraitPoint{.55, .91}, portraitPoint{.52, .79})
	poly(110, portraitPoint{.59, .53}, portraitPoint{.61, .53}, portraitPoint{.71, .88}, portraitPoint{.70, .905}, portraitPoint{.675, .87})
	// Deep sleeves keep the long hands and the floating ice visually separate.
	poly(60, portraitPoint{.34, .39}, portraitPoint{.40, .45}, portraitPoint{.34, .57}, portraitPoint{.22, .62}, portraitPoint{.18, .57}, portraitPoint{.20, .49})
	poly(67, portraitPoint{.32, .45}, portraitPoint{.35, .46}, portraitPoint{.29, .55}, portraitPoint{.22, .58}, portraitPoint{.20, .54}, portraitPoint{.24, .50})
	poly(236, portraitPoint{.20, .54}, portraitPoint{.28, .55}, portraitPoint{.25, .61}, portraitPoint{.22, .62}, portraitPoint{.18, .57})
	poly(60, portraitPoint{.58, .38}, portraitPoint{.65, .42}, portraitPoint{.73, .52}, portraitPoint{.80, .55}, portraitPoint{.79, .63}, portraitPoint{.69, .65}, portraitPoint{.59, .59}, portraitPoint{.56, .48})
	poly(67, portraitPoint{.62, .45}, portraitPoint{.65, .46}, portraitPoint{.74, .54}, portraitPoint{.79, .55}, portraitPoint{.77, .60}, portraitPoint{.69, .59})
	poly(236, portraitPoint{.66, .57}, portraitPoint{.71, .59}, portraitPoint{.77, .60}, portraitPoint{.79, .63}, portraitPoint{.69, .65}, portraitPoint{.61, .60})
	// Sparse frost along the hem belongs to the cloth, not the background.
	poly(110, portraitPoint{.22, .89}, portraitPoint{.27, .905}, portraitPoint{.46, .90}, portraitPoint{.46, .925}, portraitPoint{.27, .93}, portraitPoint{.20, .91})
	poly(153, portraitPoint{.30, .915}, portraitPoint{.33, .915}, portraitPoint{.315, .945})
	poly(153, portraitPoint{.41, .91}, portraitPoint{.44, .91}, portraitPoint{.425, .94})
	// His staff hand is a narrow, bone-long grip emerging from the sleeve.
	poly(109, portraitPoint{.18, .53}, portraitPoint{.24, .515}, portraitPoint{.26, .54}, portraitPoint{.255, .575}, portraitPoint{.22, .595}, portraitPoint{.19, .58})
	poly(152, portraitPoint{.185, .53}, portraitPoint{.23, .525}, portraitPoint{.235, .56}, portraitPoint{.20, .57})
	// The other hand turns upward under the suspended shard, with two fingers.
	poly(109, portraitPoint{.68, .555}, portraitPoint{.72, .56}, portraitPoint{.77, .55}, portraitPoint{.80, .525}, portraitPoint{.82, .53}, portraitPoint{.80, .575}, portraitPoint{.75, .595}, portraitPoint{.69, .59})
	poly(152, portraitPoint{.695, .56}, portraitPoint{.735, .57}, portraitPoint{.77, .555}, portraitPoint{.80, .525}, portraitPoint{.805, .545}, portraitPoint{.78, .58}, portraitPoint{.735, .585}, portraitPoint{.69, .575})
	poly(152, portraitPoint{.735, .568}, portraitPoint{.755, .54}, portraitPoint{.768, .54}, portraitPoint{.759, .572})
	// A peaked, heavy hood opens around a long, angular monster face.
	poly(60, portraitPoint{.44, .095}, portraitPoint{.54, .135}, portraitPoint{.61, .235}, portraitPoint{.63, .345}, portraitPoint{.60, .415}, portraitPoint{.47, .445}, portraitPoint{.34, .37}, portraitPoint{.31, .29}, portraitPoint{.345, .175})
	poly(236, portraitPoint{.44, .095}, portraitPoint{.45, .20}, portraitPoint{.38, .285}, portraitPoint{.40, .38}, portraitPoint{.47, .445}, portraitPoint{.34, .37}, portraitPoint{.31, .29}, portraitPoint{.345, .175})
	poly(67, portraitPoint{.44, .095}, portraitPoint{.54, .135}, portraitPoint{.61, .235}, portraitPoint{.60, .345}, portraitPoint{.58, .35}, portraitPoint{.55, .22}, portraitPoint{.47, .18})
	poly(234, portraitPoint{.44, .195}, portraitPoint{.53, .205}, portraitPoint{.58, .255}, portraitPoint{.59, .32}, portraitPoint{.56, .375}, portraitPoint{.48, .415}, portraitPoint{.39, .35}, portraitPoint{.38, .275})
	// A folded blade ear and three icy beard points signal ancient monster kin.
	poly(109, portraitPoint{.56, .285}, portraitPoint{.635, .235}, portraitPoint{.625, .305}, portraitPoint{.57, .34})
	poly(66, portraitPoint{.585, .29}, portraitPoint{.62, .26}, portraitPoint{.612, .30}, portraitPoint{.58, .32})
	poly(109, portraitPoint{.44, .235}, portraitPoint{.51, .22}, portraitPoint{.55, .245}, portraitPoint{.57, .295}, portraitPoint{.545, .36}, portraitPoint{.505, .405}, portraitPoint{.46, .365}, portraitPoint{.415, .29})
	poly(152, portraitPoint{.445, .24}, portraitPoint{.50, .23}, portraitPoint{.54, .255}, portraitPoint{.54, .29}, portraitPoint{.50, .31}, portraitPoint{.45, .29}, portraitPoint{.43, .27})
	poly(66, portraitPoint{.425, .29}, portraitPoint{.46, .31}, portraitPoint{.49, .335}, portraitPoint{.54, .325}, portraitPoint{.545, .36}, portraitPoint{.505, .405}, portraitPoint{.46, .365})
	poly(189, portraitPoint{.51, .27}, portraitPoint{.555, .33}, portraitPoint{.515, .34}, portraitPoint{.50, .30})
	poly(109, portraitPoint{.515, .34}, portraitPoint{.555, .33}, portraitPoint{.54, .35}, portraitPoint{.52, .355})
	poly(152, portraitPoint{.46, .36}, portraitPoint{.48, .365}, portraitPoint{.49, .445}, portraitPoint{.465, .415})
	poly(189, portraitPoint{.49, .375}, portraitPoint{.515, .37}, portraitPoint{.515, .47}, portraitPoint{.495, .44})
	poly(109, portraitPoint{.52, .365}, portraitPoint{.545, .35}, portraitPoint{.55, .41}, portraitPoint{.53, .455})
	// Cold light rests in a sharp, two-faced diamond above the open hand.
	poly(110, portraitPoint{.745, .35}, portraitPoint{.80, .425}, portraitPoint{.775, .515}, portraitPoint{.705, .47}, portraitPoint{.695, .425})
	poly(153, portraitPoint{.745, .35}, portraitPoint{.80, .425}, portraitPoint{.775, .515}, portraitPoint{.745, .445})
	poly(189, portraitPoint{.745, .35}, portraitPoint{.745, .445}, portraitPoint{.705, .47}, portraitPoint{.695, .425})
	poly(255, portraitPoint{.745, .35}, portraitPoint{.752, .39}, portraitPoint{.713, .438}, portraitPoint{.71, .416})
	// Three suspended chips repeat the crystal's geometry at a quieter scale.
	poly(110, portraitPoint{.73, .24}, portraitPoint{.75, .27}, portraitPoint{.74, .30}, portraitPoint{.72, .28})
	poly(153, portraitPoint{.87, .44}, portraitPoint{.89, .46}, portraitPoint{.88, .49}, portraitPoint{.86, .47})
	poly(67, portraitPoint{.12, .35}, portraitPoint{.14, .37}, portraitPoint{.13, .40}, portraitPoint{.11, .38})
	// The eyes sit in cold recesses, so their light is distinct from skin.
	poly(24, portraitPoint{.44, .278}, portraitPoint{.475, .278}, portraitPoint{.483, .302}, portraitPoint{.449, .301})
	poly(24, portraitPoint{.53, .285}, portraitPoint{.554, .28}, portraitPoint{.556, .312}, portraitPoint{.534, .31})
	p.detail(.46, .285, '━', 255)
	p.detail(.54, .30, '━', 255)
	p.detail(.50, .355, '─', 66)
	p.detail(.205, .55, '│', 66)
	p.detail(.23, .55, '│', 66)
	if h >= 24 {
		p.detail(.465, .325, '╲', 152)
		p.detail(.55, .76, '╲', 110)
	}
}
