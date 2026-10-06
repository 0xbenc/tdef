package render

import "math"

func paintEndingStill(p portraitPainter, st CutsceneState, art string) {
	switch art {
	case "ending-fallen", "ending-rescue":
		filmHeartBattlefield(p)
		filmFallenChampion(p, .12, .83, .34)
		filmFallenChampion(p, .58, .92, .32)
		if art == "ending-rescue" {
			filmHealer(p, .34, .23, .26, .66, true)
			// A restrained healing ward appears beside the fallen face.
			pulse := .008 * math.Sin(float64(st.Frame)/18)
			p.ring(109, .28, .73, .065+pulse, .095+pulse)
			p.smoothOval(153, .28, .73, .012, .023)
		}
	case "ending-healer":
		filmHeartBattlefield(p)
		filmHealer(p, .21, .03, .55, .96, false)
	case "ending-evacuation":
		filmHeartBattlefield(p)
		filmEvacuation(p, .5, .85)
	case "ending-maze":
		filmDepthsBackdrop(p.f)
		filmDepthsForeground(p.f, false)
		filmGrakRear(p, .38, .96)
	case "ending-supplies":
		filmEarthSupplies(p)
	case "ending-return", "ending-hoard":
		filmDragonWide(p, 120, true)
	}
	if art == "ending-return" || art == "ending-hoard" {
		// The old bowl motif is replaced by the unmistakable treasure pile.
		p.poly(236, portraitPoint{0, .76}, portraitPoint{.26, .76}, portraitPoint{.28, 1}, portraitPoint{0, 1})
		for i := 0; i < 12; i++ {
			filmCoin(p, .04+float64(i%4)*.065, .83+float64(i/4)*.06, .035)
		}
		if art == "ending-return" {
			// Place Grak in front of the coins, keeping his boots and stance.
			drawGrakPortrait(p.f, p.w/30, p.h/4, p.w/3, p.h*3/4)
		}
	}
}

func filmHeartBattlefield(p portraitPainter) {
	p.poly(234, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	// Breached masonry frames the closed inner gate; no more attackers enter.
	p.poly(238, portraitPoint{.15, .77}, portraitPoint{.15, .24}, portraitPoint{.47, .03}, portraitPoint{.81, .23}, portraitPoint{.84, .78})
	p.poly(233, portraitPoint{.23, .78}, portraitPoint{.23, .31}, portraitPoint{.47, .12}, portraitPoint{.74, .3}, portraitPoint{.76, .78})
	for _, u := range []float64{.29, .38, .47, .56, .65} {
		p.stroke(94, 3, portraitPoint{u, .33}, portraitPoint{u, .77})
	}
	p.poly(236, portraitPoint{0, .75}, portraitPoint{1, .72}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.stroke(95, 4, portraitPoint{.04, .9}, portraitPoint{.5, .97})
	p.poly(131, portraitPoint{.1, .91}, portraitPoint{.23, .91}, portraitPoint{.2, .97}, portraitPoint{.16, .95}, portraitPoint{.13, .99})
	p.poly(239, portraitPoint{.84, .75}, portraitPoint{.94, .69}, portraitPoint{.99, .8}, portraitPoint{.9, .87})
}

func filmFallenChampion(scene portraitPainter, u, v, width float64) {
	p := portraitPainter{scene.f, scene.x0 + int(u*float64(scene.w)), scene.y0 + int((v-.25)*float64(scene.h)), int(width * float64(scene.w)), max(3, scene.h/4)}
	// A knight lying across the frame, helmet left, plated boots right.
	p.poly(95, portraitPoint{.18, .38}, portraitPoint{.68, .26}, portraitPoint{.97, .8}, portraitPoint{.2, .86})
	p.poly(239, portraitPoint{.29, .22}, portraitPoint{.6, .18}, portraitPoint{.73, .53}, portraitPoint{.52, .7}, portraitPoint{.27, .62})
	p.poly(246, portraitPoint{.29, .22}, portraitPoint{.59, .18}, portraitPoint{.64, .33}, portraitPoint{.31, .41})
	p.poly(244, portraitPoint{.64, .42}, portraitPoint{.91, .49}, portraitPoint{.97, .72}, portraitPoint{.71, .66})
	p.smoothOval(244, .16, .45, .14, .27)
	p.stroke(233, 1.5, portraitPoint{.06, .47}, portraitPoint{.23, .47})
	// Shield and sword are separated from the body instead of implying gore.
	p.poly(137, portraitPoint{.35, .48}, portraitPoint{.6, .45}, portraitPoint{.66, .72}, portraitPoint{.49, .94}, portraitPoint{.32, .7})
	p.stroke(246, 1, portraitPoint{.45, .91}, portraitPoint{.94, .92})
}

func filmHealer(scene portraitPainter, u, v, width, height float64, kneeling bool) {
	p := portraitPainter{scene.f, scene.x0 + int(u*float64(scene.w)), scene.y0 + int(v*float64(scene.h)), int(width * float64(scene.w)), int(height * float64(scene.h))}
	// Ivory hood, teal stole, leather medical bag; distinct from armored guild
	// fighters. Bent posture and an extended sleeve read even at small sizes.
	p.poly(246, portraitPoint{.35, .32}, portraitPoint{.66, .34}, portraitPoint{.81, .93}, portraitPoint{.17, .94})
	p.poly(240, portraitPoint{.6, .36}, portraitPoint{.73, .44}, portraitPoint{.81, .93}, portraitPoint{.52, .92})
	p.poly(23, portraitPoint{.36, .32}, portraitPoint{.48, .35}, portraitPoint{.42, .92}, portraitPoint{.29, .92})
	p.poly(109, portraitPoint{.37, .34}, portraitPoint{.41, .35}, portraitPoint{.35, .89}, portraitPoint{.31, .89})
	p.poly(246, portraitPoint{.25, .19}, portraitPoint{.46, .03}, portraitPoint{.68, .18}, portraitPoint{.64, .38}, portraitPoint{.33, .39})
	p.poly(236, portraitPoint{.32, .21}, portraitPoint{.46, .13}, portraitPoint{.59, .21}, portraitPoint{.56, .33}, portraitPoint{.36, .34})
	p.poly(223, portraitPoint{.37, .21}, portraitPoint{.53, .22}, portraitPoint{.53, .33}, portraitPoint{.39, .33})
	if kneeling {
		p.poly(246, portraitPoint{.31, .4}, portraitPoint{.4, .5}, portraitPoint{.09, .75}, portraitPoint{0, .7}, portraitPoint{.07, .61})
		p.smoothOval(223, .04, .73, .065, .06)
	} else {
		p.poly(246, portraitPoint{.67, .4}, portraitPoint{.82, .43}, portraitPoint{.75, .68}, portraitPoint{.54, .63}, portraitPoint{.55, .55})
		p.smoothOval(223, .6, .58, .07, .05)
	}
	p.poly(94, portraitPoint{.58, .7}, portraitPoint{.81, .7}, portraitPoint{.83, .86}, portraitPoint{.6, .88})
	p.stroke(137, 1.2, portraitPoint{.42, .36}, portraitPoint{.69, .71})
	p.stroke(223, 1.5, portraitPoint{.65, .76}, portraitPoint{.76, .76})
}

func filmEvacuation(p portraitPainter, u, v float64) {
	// Human proportions and litter length use native cell dimensions, not
	// fractions of the two-screen panorama. The party fits through the door.
	h := max(5, p.h*38/100)
	cx := p.x0 + int(u*float64(p.w-1))
	feet := p.y0 + int(v*float64(p.h-1))
	span := h * 2
	left, right := cx-span/2, cx+span/2
	filmCarrier(p.f, left, feet, h)
	filmCarrier(p.f, right, feet, h)
	healH := max(6, p.h*48/100)
	healW := healH
	healer := portraitPainter{p.f, left - healW - h/3, feet - healH + 1, healW, healH}
	filmHealer(healer, 0, 0, 1, 1, false)
	// Litter, lying champion and helmet drawn as one readable horizontal mass.
	bed := feet - h/2
	for x := left; x <= right; x++ {
		p.f.Put(x, bed, '?', 94, 233)
		if x > left+1 && x < right-1 {
			p.f.Put(x, bed-1, '?', 95, 233)
		}
	}
	knight := portraitPainter{p.f, left + h/4, bed - max(2, h/4), max(6, span-h/2), max(2, h/4)}
	knight.poly(244, portraitPoint{.18, .15}, portraitPoint{.8, .15}, portraitPoint{.94, .9}, portraitPoint{.13, .9})
	knight.smoothOval(246, .07, .45, .11, .46)
}

func filmCarrier(f *Frame, cx, feet, h int) {
	w := max(6, h)
	p := portraitPainter{f, cx - w/2, feet - h + 1, w, h}
	p.poly(95, portraitPoint{.2, .3}, portraitPoint{.72, .3}, portraitPoint{.85, .78}, portraitPoint{.14, .78})
	p.poly(239, portraitPoint{.3, .3}, portraitPoint{.64, .3}, portraitPoint{.66, .64}, portraitPoint{.27, .64})
	p.poly(246, portraitPoint{.3, .31}, portraitPoint{.64, .31}, portraitPoint{.65, .42}, portraitPoint{.3, .45})
	p.poly(94, portraitPoint{.25, .68}, portraitPoint{.43, .68}, portraitPoint{.42, .99}, portraitPoint{.2, .99})
	p.poly(94, portraitPoint{.56, .68}, portraitPoint{.7, .68}, portraitPoint{.8, .99}, portraitPoint{.55, .99})
	p.poly(244, portraitPoint{.32, .03}, portraitPoint{.65, .03}, portraitPoint{.7, .21}, portraitPoint{.6, .31}, portraitPoint{.32, .27}, portraitPoint{.27, .13})
	p.poly(233, portraitPoint{.33, .16}, portraitPoint{.63, .16}, portraitPoint{.63, .2}, portraitPoint{.33, .2})
	p.poly(244, portraitPoint{.19, .4}, portraitPoint{.36, .45}, portraitPoint{.6, .48}, portraitPoint{.6, .56}, portraitPoint{.21, .53})
}

func paintEndingPanorama(dst *Frame, stage Rect, art string, frame int) {
	far := filmFrame(stage.W*2, stage.H)
	near := filmTransparentFrame(stage.W*2, stage.H)
	switch art {
	case "ending-pursuit":
		filmOutdoorBackdrop(far, "ruined-road")
		filmPursuitForeground(near)
	case "ending-depths":
		filmDepthsBackdrop(far)
		filmDepthsForeground(near, true)
	case "ending-earth":
		filmEarthBackdrop(far)
		filmEarthForeground(near)
	}
	x := filmPanColumn(art, frame, stage.W)
	filmCompositePan(dst, stage, far, near, x/4, x)
}

func filmPursuitForeground(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.poly(236, portraitPoint{0, .73}, portraitPoint{1, .73}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.poly(239, portraitPoint{0, .86}, portraitPoint{1, .85}, portraitPoint{1, .94}, portraitPoint{0, .95})
	// Ruined capital to a separate western entrance, with a broad stretch of
	// open road between them. The Depths are not a hole beneath the vault.
	for _, u := range []float64{.04, .2} {
		p.poly(94, portraitPoint{u, .74}, portraitPoint{u, .33}, portraitPoint{u + .04, .39}, portraitPoint{u + .1, .27}, portraitPoint{u + .13, .74})
		p.stroke(137, 1, portraitPoint{u, .63}, portraitPoint{u + .13, .7})
	}
	p.poly(237, portraitPoint{.65, 1}, portraitPoint{.66, .36}, portraitPoint{.82, .08}, portraitPoint{1, .35}, portraitPoint{1, 1})
	p.poly(60, portraitPoint{.74, .91}, portraitPoint{.74, .41}, portraitPoint{.83, .23}, portraitPoint{.93, .42}, portraitPoint{.94, .91})
	p.poly(233, portraitPoint{.77, .91}, portraitPoint{.77, .43}, portraitPoint{.83, .3}, portraitPoint{.89, .45}, portraitPoint{.91, .91})
	filmGrakRear(p, .17, .93)
	filmEvacuation(p, .8, .87)
}

func filmDepthsBackdrop(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.poly(234, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	// Impossible overlapping doorways float beyond the solid walking surface.
	for _, u := range []float64{.08, .3, .53, .77} {
		p.poly(60, portraitPoint{u, .87}, portraitPoint{u, .32}, portraitPoint{u + .07, .14}, portraitPoint{u + .14, .3}, portraitPoint{u + .15, .88})
		p.poly(235, portraitPoint{u + .026, .88}, portraitPoint{u + .026, .34}, portraitPoint{u + .07, .23}, portraitPoint{u + .112, .34}, portraitPoint{u + .125, .88})
		p.stroke(103, 1, portraitPoint{u + .024, .8}, portraitPoint{u + .024, .34}, portraitPoint{u + .07, .23})
	}
}

func filmDepthsForeground(f *Frame, exit bool) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.poly(237, portraitPoint{0, .83}, portraitPoint{.43, .72}, portraitPoint{1, .83}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.poly(60, portraitPoint{0, .94}, portraitPoint{.42, .83}, portraitPoint{.7, .92}, portraitPoint{1, .85}, portraitPoint{1, .94}, portraitPoint{.72, 1}, portraitPoint{.43, .91}, portraitPoint{0, 1})
	for _, u := range []float64{.06, .36, .62, .97} {
		p.poly(238, portraitPoint{u, 0}, portraitPoint{u + .023, 0}, portraitPoint{u + .025, .84}, portraitPoint{u - .005, .88})
		p.poly(60, portraitPoint{u + .023, 0}, portraitPoint{u + .043, 0}, portraitPoint{u + .055, .88}, portraitPoint{u + .025, .84})
	}
	if exit {
		p.poly(24, portraitPoint{.76, .82}, portraitPoint{.76, .29}, portraitPoint{.82, .16}, portraitPoint{.89, .29}, portraitPoint{.9, .82})
		p.poly(153, portraitPoint{.78, .8}, portraitPoint{.78, .3}, portraitPoint{.82, .23}, portraitPoint{.865, .31}, portraitPoint{.88, .8})
		p.poly(223, portraitPoint{.8, .74}, portraitPoint{.8, .51}, portraitPoint{.85, .51}, portraitPoint{.86, .74})
		filmEvacuation(p, .8, .86)
		filmGrakRear(p, .18, .96)
	}
}

func filmEarthBackdrop(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.poly(24, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	for i, u := range []float64{.08, .23, .38, .52, .68, .83} {
		top := .31 + float64(i%3)*.065
		p.poly(238, portraitPoint{u, .84}, portraitPoint{u, top}, portraitPoint{u + .095, top}, portraitPoint{u + .095, .84})
		for row := 0; row < 3; row++ {
			v := top + .08 + float64(row)*.12
			p.poly(180, portraitPoint{u + .026, v}, portraitPoint{u + .04, v}, portraitPoint{u + .04, v + .035}, portraitPoint{u + .026, v + .035})
		}
	}
}

func filmEarthForeground(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	// We start inside the purple passage looking through a lit opening, then
	// pan out onto an ordinary street: apartments, a shop, car and delivery van.
	p.poly(235, portraitPoint{0, .8}, portraitPoint{1, .8}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.poly(244, portraitPoint{.4, .86}, portraitPoint{1, .86}, portraitPoint{1, .88}, portraitPoint{.4, .88})
	p.poly(60, portraitPoint{0, 0}, portraitPoint{.085, 0}, portraitPoint{.085, .84}, portraitPoint{0, 1})
	p.poly(60, portraitPoint{.25, 0}, portraitPoint{.42, 0}, portraitPoint{.43, 1}, portraitPoint{.25, .84})
	p.poly(60, portraitPoint{.085, 0}, portraitPoint{.25, 0}, portraitPoint{.25, .25}, portraitPoint{.17, .15}, portraitPoint{.085, .25})
	p.poly(103, portraitPoint{.08, .84}, portraitPoint{.08, .24}, portraitPoint{.17, .13}, portraitPoint{.25, .24}, portraitPoint{.255, .84}, portraitPoint{.24, .83}, portraitPoint{.235, .27}, portraitPoint{.17, .19}, portraitPoint{.097, .27}, portraitPoint{.096, .84})
	p.poly(95, portraitPoint{.56, .8}, portraitPoint{.56, .12}, portraitPoint{.76, .12}, portraitPoint{.76, .8})
	for _, u := range []float64{.59, .66, .72} {
		for _, v := range []float64{.21, .38} {
			p.poly(223, portraitPoint{u, v}, portraitPoint{u + .025, v}, portraitPoint{u + .025, v + .065}, portraitPoint{u, v + .065})
		}
	}
	p.poly(23, portraitPoint{.57, .56}, portraitPoint{.75, .56}, portraitPoint{.75, .64}, portraitPoint{.57, .64})
	p.poly(180, portraitPoint{.585, .67}, portraitPoint{.685, .67}, portraitPoint{.685, .8}, portraitPoint{.585, .8})
	p.poly(233, portraitPoint{.71, .67}, portraitPoint{.74, .67}, portraitPoint{.74, .8}, portraitPoint{.71, .8})
	filmNativeSign(f, .58, .575, "MARKET", 223, 23)
	// Flat-sided modern delivery van and a low car distinguish Earth from
	// another fantasy settlement without explaining the technology in text.
	p.poly(246, portraitPoint{.78, .82}, portraitPoint{.78, .62}, portraitPoint{.91, .62}, portraitPoint{.92, .74}, portraitPoint{.955, .77}, portraitPoint{.955, .87}, portraitPoint{.78, .87})
	p.poly(24, portraitPoint{.895, .67}, portraitPoint{.916, .675}, portraitPoint{.922, .744}, portraitPoint{.895, .744})
	p.smoothOval(233, .81, .87, .018, .045)
	p.smoothOval(233, .933, .87, .018, .045)
	p.poly(110, portraitPoint{.46, .95}, portraitPoint{.47, .9}, portraitPoint{.495, .85}, portraitPoint{.54, .85}, portraitPoint{.565, .91}, portraitPoint{.59, .925}, portraitPoint{.59, .96})
	p.poly(24, portraitPoint{.5, .86}, portraitPoint{.535, .86}, portraitPoint{.55, .9}, portraitPoint{.49, .9})
	p.smoothOval(233, .49, .96, .012, .03)
	p.smoothOval(233, .565, .96, .012, .03)
	// A lamp, a customer, and loaded crates establish an inhabited place.
	p.stroke(244, 1, portraitPoint{.52, .82}, portraitPoint{.52, .34}, portraitPoint{.545, .34})
	p.smoothOval(223, .545, .35, .012, .025)
	p.poly(94, portraitPoint{.76, .81}, portraitPoint{.76, .73}, portraitPoint{.785, .73}, portraitPoint{.785, .81})
	filmEarthResident(p, .697, .81)
	filmEvacuation(p, .19, .85)
}

func filmNativeSign(f *Frame, u, v float64, text string, fg, bg int) {
	x, y := int(u*float64(f.W)), int(v*float64(f.H))
	if f.W >= 100 {
		putString(f, x, y, text, fg, bg, false)
	}
}

func filmEarthResident(p portraitPainter, u, v float64) {
	p.poly(110, portraitPoint{u - .012, v}, portraitPoint{u - .012, v - .11}, portraitPoint{u + .012, v - .11}, portraitPoint{u + .018, v})
	p.smoothOval(223, u, v-.14, .008, .026)
	p.poly(94, portraitPoint{u + .02, v - .02}, portraitPoint{u + .02, v - .075}, portraitPoint{u + .036, v - .075}, portraitPoint{u + .036, v - .02})
}

func filmEarthSupplies(p portraitPainter) {
	p.poly(95, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.poly(23, portraitPoint{.08, .17}, portraitPoint{.68, .17}, portraitPoint{.68, .71}, portraitPoint{.08, .71})
	p.poly(109, portraitPoint{.1, .19}, portraitPoint{.66, .19}, portraitPoint{.66, .68}, portraitPoint{.1, .68})
	// Tools in a shop window; lumber and a loaded crate beside the door.
	for i := 0; i < 4; i++ {
		u := .18 + float64(i)*.125
		p.stroke(94, 3, portraitPoint{u, .32}, portraitPoint{u, .62})
		p.poly(246, portraitPoint{u - .036, .28}, portraitPoint{u + .047, .28}, portraitPoint{u + .047, .35}, portraitPoint{u - .036, .35})
	}
	p.poly(233, portraitPoint{.77, .24}, portraitPoint{.94, .24}, portraitPoint{.94, .91}, portraitPoint{.77, .91})
	p.poly(236, portraitPoint{0, .76}, portraitPoint{1, .76}, portraitPoint{1, 1}, portraitPoint{0, 1})
	for i := 0; i < 3; i++ {
		v := .78 + float64(i)*.065
		p.poly(137, portraitPoint{.05, v}, portraitPoint{.62, v - .04}, portraitPoint{.62, v + .015}, portraitPoint{.05, v + .055})
	}
	p.poly(94, portraitPoint{.64, .69}, portraitPoint{.8, .69}, portraitPoint{.8, .93}, portraitPoint{.64, .93})
	p.stroke(137, 3, portraitPoint{.65, .71}, portraitPoint{.79, .91})
	filmNativeSign(p.f, .16, .08, "HARDWARE", 223, 95)
}
