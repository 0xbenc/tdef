package render

func drawRuneforgePortrait(f *Frame, x, y, w, h int) {
	drawSpecialistPortrait(f, x, y, w, h, "runeforge")
}
func drawHookmasterPortrait(f *Frame, x, y, w, h int) {
	drawSpecialistPortrait(f, x, y, w, h, "hookmaster")
}
func drawSappersPortrait(f *Frame, x, y, w, h int) { drawSpecialistPortrait(f, x, y, w, h, "sappers") }
func drawWitchPortrait(f *Frame, x, y, w, h int)   { drawSpecialistPortrait(f, x, y, w, h, "witch") }

func drawSpecialistPortrait(f *Frame, x, y, w, h int, kind string) {
	art := filmCachedArtwork("portrait-"+kind, w, h, false, func(dst *Frame) {
		p := portraitPainter{dst, 0, 0, w, h}
		switch kind {
		case "runeforge":
			paintRuneforgePortrait(p)
		case "hookmaster":
			paintHookmasterPortrait(p)
		case "sappers":
			paintSappersPortrait(p)
		case "witch":
			paintWitchPortrait(p)
		}
	})
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			f.Set(x+xx, y+yy, art.C[yy*w+xx])
		}
	}
}

func paintRuneforgePortrait(p portraitPainter) {
	box := func(c int, x, y, xx, yy float64) { filmEarthRect(p, c, x, y, xx, yy) }
	// A squat brass furnace, tilted bellows and an unmistakable white-hot eye.
	p.poly(236, portraitPoint{.04, .85}, portraitPoint{.78, .82}, portraitPoint{.96, .95}, portraitPoint{.14, .99})
	p.poly(94, portraitPoint{.16, .42}, portraitPoint{.35, .27}, portraitPoint{.73, .3}, portraitPoint{.81, .45}, portraitPoint{.77, .87}, portraitPoint{.19, .88})
	box(137, .2, .42, .71, .85)
	p.poly(180, portraitPoint{.2, .42}, portraitPoint{.36, .3}, portraitPoint{.69, .32}, portraitPoint{.71, .42})
	p.poly(130, portraitPoint{.71, .42}, portraitPoint{.8, .45}, portraitPoint{.77, .87}, portraitPoint{.71, .85})
	box(239, .27, .06, .39, .33)
	box(246, .27, .06, .3, .33)
	box(94, .24, .04, .42, .09)
	box(95, .48, .1, .66, .3)
	box(180, .5, .12, .64, .15)
	box(180, .5, .25, .64, .28)
	for i := 0; i < 4; i++ {
		u := .51 + float64(i)*.033
		p.stroke(137, 1, portraitPoint{u, .16}, portraitPoint{u, .24})
	}
	p.ring(94, .45, .61, .22, .27)
	p.smoothOval(23, .45, .61, .18, .22)
	p.ring(30, .45, .61, .145, .18)
	p.ring(81, .45, .61, .11, .14)
	p.smoothOval(117, .45, .61, .075, .1)
	p.smoothOval(231, .45, .61, .035, .07)
	for i := 0; i < 6; i++ {
		u := .24 + float64(i)*.085
		p.smoothOval(223, u, .45, .009, .013)
		p.smoothOval(223, u, .83, .009, .013)
	}
	p.stroke(180, 3, portraitPoint{.72, .62}, portraitPoint{.86, .56}, portraitPoint{.9, .75})
	p.poly(95, portraitPoint{.8, .71}, portraitPoint{.91, .62}, portraitPoint{.98, .77}, portraitPoint{.86, .88})
	for i := 0; i < 4; i++ {
		v := .73 + float64(i)*.03
		p.stroke(137, 1, portraitPoint{.83, v}, portraitPoint{.93, v - .06})
	}
	box(239, .24, .87, .33, .95)
	box(239, .62, .87, .71, .95)
	p.stroke(81, 1, portraitPoint{.12, .53}, portraitPoint{.04, .52})
	p.stroke(117, 1, portraitPoint{.13, .65}, portraitPoint{.03, .7})
}

func paintHookmasterPortrait(p portraitPainter) {
	// Huge bare forearms, a working harness and a hook larger than his head.
	p.poly(94, portraitPoint{.24, .46}, portraitPoint{.7, .44}, portraitPoint{.78, .8}, portraitPoint{.2, .85})
	p.poly(65, portraitPoint{.18, .45}, portraitPoint{.32, .42}, portraitPoint{.35, .64}, portraitPoint{.2, .75}, portraitPoint{.09, .63})
	p.poly(107, portraitPoint{.7, .43}, portraitPoint{.85, .49}, portraitPoint{.87, .7}, portraitPoint{.69, .77}, portraitPoint{.62, .6})
	p.poly(107, portraitPoint{.34, .15}, portraitPoint{.56, .12}, portraitPoint{.65, .26}, portraitPoint{.59, .45}, portraitPoint{.38, .45}, portraitPoint{.29, .31})
	p.poly(150, portraitPoint{.34, .15}, portraitPoint{.55, .13}, portraitPoint{.59, .19}, portraitPoint{.33, .23})
	p.poly(65, portraitPoint{.14, .2}, portraitPoint{.34, .24}, portraitPoint{.32, .32}, portraitPoint{.2, .29})
	p.poly(107, portraitPoint{.6, .23}, portraitPoint{.75, .15}, portraitPoint{.72, .25}, portraitPoint{.62, .32})
	p.stroke(58, 4, portraitPoint{.34, .28}, portraitPoint{.44, .29})
	p.stroke(58, 4, portraitPoint{.51, .28}, portraitPoint{.61, .25})
	p.stroke(223, 1, portraitPoint{.36, .29}, portraitPoint{.42, .3})
	p.stroke(223, 1, portraitPoint{.53, .29}, portraitPoint{.59, .28})
	p.poly(223, portraitPoint{.36, .38}, portraitPoint{.4, .39}, portraitPoint{.36, .33})
	p.poly(223, portraitPoint{.55, .39}, portraitPoint{.6, .37}, portraitPoint{.61, .31})
	p.stroke(130, 5, portraitPoint{.3, .45}, portraitPoint{.67, .77})
	p.stroke(180, 1, portraitPoint{.31, .46}, portraitPoint{.67, .77})
	p.poly(58, portraitPoint{.24, .76}, portraitPoint{.45, .76}, portraitPoint{.43, .95}, portraitPoint{.22, .95})
	p.poly(58, portraitPoint{.52, .77}, portraitPoint{.7, .75}, portraitPoint{.76, .95}, portraitPoint{.54, .95})
	for i := 0; i < 7; i++ {
		u := .16 + float64(i)*.095
		p.ring(246, u, .67, .037, .026)
	}
	p.smoothOval(107, .24, .65, .045, .055)
	p.smoothOval(150, .69, .66, .047, .055)
	p.stroke(239, 7, portraitPoint{.84, .65}, portraitPoint{.9, .73}, portraitPoint{.89, .85}, portraitPoint{.8, .9}, portraitPoint{.75, .85}, portraitPoint{.78, .77})
	p.stroke(252, 2, portraitPoint{.86, .67}, portraitPoint{.92, .75}, portraitPoint{.89, .85}, portraitPoint{.81, .88})
}

func paintSappersPortrait(p portraitPainter) {
	// Two goblins, one open kettle-mine, and carefully separated powder stock.
	for i, u := range []float64{.12, .55} {
		p.poly(95, portraitPoint{u, .4}, portraitPoint{u + .25, .4}, portraitPoint{u + .3, .76}, portraitPoint{u - .025, .76})
		p.poly(107, portraitPoint{u + .04, .18}, portraitPoint{u + .19, .16}, portraitPoint{u + .24, .27}, portraitPoint{u + .17, .39}, portraitPoint{u + .03, .35})
		p.poly(150, portraitPoint{u - .08, .17}, portraitPoint{u + .06, .23}, portraitPoint{u + .05, .29})
		p.poly(107, portraitPoint{u + .18, .23}, portraitPoint{u + .34, .12}, portraitPoint{u + .26, .29}, portraitPoint{u + .21, .31})
		p.poly(130, portraitPoint{u + .02, .18}, portraitPoint{u + .08, .09}, portraitPoint{u + .18, .1}, portraitPoint{u + .22, .18})
		p.stroke(223, 1, portraitPoint{u + .05, .27}, portraitPoint{u + .09, .28})
		p.stroke(223, 1, portraitPoint{u + .15, .26}, portraitPoint{u + .19, .25})
		p.poly(65, portraitPoint{u + .2, .3}, portraitPoint{u + .3, .33}, portraitPoint{u + .2, .35})
		p.stroke(137, 3, portraitPoint{u + .03, .41}, portraitPoint{u + .22, .71})
		if i == 0 {
			p.stroke(246, 2, portraitPoint{u + .23, .45}, portraitPoint{.45, .65})
		}
	}
	p.smoothOval(94, .2, .83, .16, .15)
	p.ring(180, .2, .83, .16, .15)
	p.stroke(137, 2, portraitPoint{.08, .76}, portraitPoint{.31, .75})
	p.stroke(137, 2, portraitPoint{.06, .88}, portraitPoint{.34, .88})
	p.smoothOval(239, .55, .82, .19, .13)
	p.smoothOval(233, .55, .77, .15, .07)
	p.smoothOval(94, .55, .79, .1, .035)
	p.stroke(180, 2, portraitPoint{.55, .78}, portraitPoint{.6, .66}, portraitPoint{.69, .65})
	p.smoothOval(220, .69, .65, .014, .025)
	p.ring(246, .82, .87, .1, .07)
	p.smoothOval(239, .82, .87, .075, .045)
	p.stroke(94, 3, portraitPoint{.73, .96}, portraitPoint{.94, .76})
	p.poly(246, portraitPoint{.89, .74}, portraitPoint{.95, .72}, portraitPoint{.98, .77}, portraitPoint{.93, .82})
}

func paintWitchPortrait(p portraitPainter) {
	// A seated hag, pointed hood, crooked staff and an open book of names.
	p.poly(94, portraitPoint{.24, .6}, portraitPoint{.68, .6}, portraitPoint{.71, .95}, portraitPoint{.21, .95})
	p.poly(53, portraitPoint{.31, .39}, portraitPoint{.58, .37}, portraitPoint{.74, .89}, portraitPoint{.2, .9})
	p.poly(96, portraitPoint{.31, .39}, portraitPoint{.42, .43}, portraitPoint{.35, .85}, portraitPoint{.22, .87})
	p.poly(177, portraitPoint{.27, .37}, portraitPoint{.43, .05}, portraitPoint{.63, .35}, portraitPoint{.57, .48}, portraitPoint{.35, .47})
	p.poly(53, portraitPoint{.33, .34}, portraitPoint{.44, .2}, portraitPoint{.57, .35}, portraitPoint{.53, .44}, portraitPoint{.36, .44})
	p.poly(107, portraitPoint{.37, .3}, portraitPoint{.5, .3}, portraitPoint{.55, .39}, portraitPoint{.47, .49}, portraitPoint{.38, .43})
	p.poly(150, portraitPoint{.48, .34}, portraitPoint{.63, .42}, portraitPoint{.51, .42})
	p.stroke(223, 2, portraitPoint{.38, .35}, portraitPoint{.43, .36})
	p.stroke(58, 2, portraitPoint{.39, .43}, portraitPoint{.5, .44})
	p.stroke(94, 5, portraitPoint{.82, .96}, portraitPoint{.78, .26}, portraitPoint{.87, .14}, portraitPoint{.94, .24})
	p.stroke(137, 1, portraitPoint{.83, .94}, portraitPoint{.8, .27}, portraitPoint{.87, .16})
	p.smoothOval(213, .88, .24, .035, .06)
	p.ring(96, .88, .24, .065, .1)
	p.poly(107, portraitPoint{.58, .51}, portraitPoint{.65, .52}, portraitPoint{.8, .57}, portraitPoint{.79, .63}, portraitPoint{.63, .6})
	p.poly(130, portraitPoint{.21, .67}, portraitPoint{.47, .61}, portraitPoint{.69, .68}, portraitPoint{.67, .81}, portraitPoint{.46, .76}, portraitPoint{.23, .82})
	p.poly(223, portraitPoint{.23, .66}, portraitPoint{.46, .62}, portraitPoint{.46, .74}, portraitPoint{.25, .8})
	p.poly(180, portraitPoint{.47, .62}, portraitPoint{.67, .68}, portraitPoint{.65, .78}, portraitPoint{.47, .74})
	for i := 0; i < 3; i++ {
		v := .67 + float64(i)*.027
		p.stroke(94, 1, portraitPoint{.29, v}, portraitPoint{.42, v - .023})
		p.stroke(94, 1, portraitPoint{.51, v - .01}, portraitPoint{.62, v + .015})
	}
	for _, u := range []float64{.09, .15} {
		filmEarthRect(p, 180, u, .76, u+.025, .94)
		p.smoothOval(220, u+.012, .72, .013, .035)
		p.smoothOval(231, u+.012, .72, .005, .018)
	}
}
