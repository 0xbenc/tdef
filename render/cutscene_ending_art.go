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
	case "ending-return":
		filmReturnToVault(p)
	case "ending-hoard":
		filmDragonWide(p, 120, true)
	}
	if art == "ending-hoard" {
		// The old bowl motif is replaced by the unmistakable treasure pile.
		p.poly(236, portraitPoint{0, .76}, portraitPoint{.26, .76}, portraitPoint{.28, 1}, portraitPoint{0, 1})
		for i := 0; i < 12; i++ {
			filmCoin(p, .04+float64(i%4)*.065, .83+float64(i/4)*.06, .035)
		}
	}
}

func filmReturnToVault(p portraitPainter) {
	// A wide arrival shot: Grak on bare stone at left, Malgrath and the
	// treasure within their own recessed vault at right. The central aisle
	// stays empty, separating the figures even at compact terminal sizes.
	filmEarthRect(p, 234, 0, 0, 1, 1)
	p.poly(236, portraitPoint{0, .73}, portraitPoint{.55, .66}, portraitPoint{1, .76}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.stroke(238, 1, portraitPoint{.32, 1}, portraitPoint{.45, .69})
	p.stroke(238, 1, portraitPoint{.72, 1}, portraitPoint{.62, .69})
	// Cold light through the entrance frames the returning builder.
	p.poly(60, portraitPoint{.025, .78}, portraitPoint{.025, .23}, portraitPoint{.15, .07}, portraitPoint{.28, .23}, portraitPoint{.28, .78})
	p.poly(24, portraitPoint{.045, .78}, portraitPoint{.045, .25}, portraitPoint{.15, .12}, portraitPoint{.26, .25}, portraitPoint{.26, .78})
	p.stroke(103, 1, portraitPoint{.045, .75}, portraitPoint{.045, .25}, portraitPoint{.15, .12}, portraitPoint{.26, .25})
	p.poly(237, portraitPoint{.045, .78}, portraitPoint{.26, .78}, portraitPoint{.36, 1}, portraitPoint{.015, 1})
	drawDragonTableau(p.f, p.x0+p.w*40/100, p.y0+p.h*9/100, p.w*58/100, p.h*70/100, 5, 233, owPadView{status: OWOpen}, 120, true)
	filmGrakRear(p, .17, .94)
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
	if art == "ending-earth" {
		far := filmCachedArtwork("earth-skyline", stage.W*2, stage.H, false, filmEarthBackdrop)
		near := filmCachedArtwork("earth-street", stage.W*2, stage.H, true, filmEarthForeground)
		x := filmPanColumn(art, frame, stage.W)
		filmCompositePan(dst, stage, far, near, x/4, x)
		return
	}
	far := filmFrame(stage.W*2, stage.H)
	near := filmTransparentFrame(stage.W*2, stage.H)
	switch art {
	case "ending-pursuit":
		filmOutdoorBackdrop(far, "ruined-road")
		filmPursuitForeground(near)
	case "ending-depths":
		filmDepthsBackdrop(far)
		filmDepthsForeground(near, true)
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

// Earth opens onto a wet, luminous city, with the distant skyline moving more
// slowly than the portal and storefronts. Keep all light steady while reading.
func filmEarthBackdrop(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	filmEarthRect(p, 17, 0, 0, 1, 1)
	filmEarthRect(p, 18, 0, .19, 1, .72)
	filmEarthRect(p, 24, 0, .45, 1, .85)
	p.smoothOval(153, .34, .16, .025, .07)
	p.smoothOval(189, .336, .15, .018, .05)
	// A jagged distant skyline, antennae and hundreds of inhabited windows.
	for i := 0; i < 22; i++ {
		u := float64(i) * .047
		top := .24 + float64((i*7)%9)*.034
		filmEarthRect(p, 60, u, top, u+.037, .84)
		p.stroke(103, .6, portraitPoint{u + .018, top}, portraitPoint{u + .018, top - .05})
		for row := 0; row < 8; row++ {
			v := top + .035 + float64(row)*.045
			if v > .79 {
				break
			}
			for col := 0; col < 3; col++ {
				if (i+row+col)%4 == 0 {
					continue
				}
				c := 110
				if (i+row)%3 == 0 {
					c = 180
				}
				filmEarthRect(p, c, u+.005+float64(col)*.01, v, u+.009+float64(col)*.01, v+.014)
			}
		}
	}
	// A landmark glass tower and a sweeping elevated rail silhouette.
	filmEarthRect(p, 23, .43, .08, .49, .8)
	p.poly(30, portraitPoint{.43, .08}, portraitPoint{.46, .025}, portraitPoint{.49, .08})
	for i := 0; i < 5; i++ {
		u := .435 + float64(i)*.011
		p.stroke(80, .6, portraitPoint{u, .11}, portraitPoint{u, .73})
	}
	p.stroke(117, 1, portraitPoint{.46, .025}, portraitPoint{.46, 0})
	p.stroke(239, 2, portraitPoint{0, .7}, portraitPoint{.25, .65}, portraitPoint{.6, .67}, portraitPoint{1, .6})
	p.stroke(110, .6, portraitPoint{0, .68}, portraitPoint{.25, .63}, portraitPoint{.6, .65}, portraitPoint{1, .58})
}

func filmEarthRect(p portraitPainter, color int, x0, y0, x1, y1 float64) {
	p.poly(color, portraitPoint{x0, y0}, portraitPoint{x1, y0}, portraitPoint{x1, y1}, portraitPoint{x0, y1})
}

func filmEarthForeground(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	filmEarthRect(p, 234, 0, .8, 1, 1)
	filmEarthRect(p, 60, 0, .8, 1, .83)
	filmEarthRect(p, 103, .38, .83, 1, .845)
	// Wet asphalt catches broken columns of cyan, pink and amber light.
	for i := 0; i < 28; i++ {
		u := .4 + float64((i*13)%59)*.01
		v := .86 + float64(i%5)*.026
		c := []int{24, 60, 30, 96, 131, 95}[i%6]
		p.stroke(c, 1, portraitPoint{u, v}, portraitPoint{u + .02 + float64(i%3)*.012, v - .008})
	}
	for i := 0; i < 5; i++ {
		u := .43 + float64(i)*.13
		filmEarthRect(p, 180, u, .965, u+.055, .975)
	}
	// The old world's heavy stone arch silhouettes the electric city beyond.
	p.poly(235, portraitPoint{0, 0}, portraitPoint{.075, 0}, portraitPoint{.085, .84}, portraitPoint{0, 1})
	p.poly(235, portraitPoint{.25, 0}, portraitPoint{.4, 0}, portraitPoint{.42, 1}, portraitPoint{.25, .84})
	p.poly(60, portraitPoint{.075, 0}, portraitPoint{.25, 0}, portraitPoint{.25, .25}, portraitPoint{.17, .14}, portraitPoint{.075, .25})
	p.stroke(141, 2, portraitPoint{.08, .84}, portraitPoint{.08, .25}, portraitPoint{.17, .14}, portraitPoint{.25, .25}, portraitPoint{.255, .84})
	p.stroke(189, .6, portraitPoint{.092, .83}, portraitPoint{.092, .27}, portraitPoint{.17, .18}, portraitPoint{.24, .27})
	// Deep facades, warm apartments, bright awnings and plate-glass shops.
	for i, u := range []float64{.45, .63, .82} {
		top := .14 + float64(i%2)*.075
		filmEarthRect(p, []int{235, 95, 236}[i], u, top, u+.155, .8)
		filmEarthRect(p, 239, u, top, u+.155, top+.025)
		filmEarthRect(p, 233, u+.135, top+.025, u+.155, .8)
		for row := 0; row < 3; row++ {
			v := top + .065 + float64(row)*.11
			for col := 0; col < 3; col++ {
				x := u + .016 + float64(col)*.04
				filmEarthRect(p, 180, x, v, x+.022, v+.06)
				filmEarthRect(p, 223, x, v, x+.009, v+.06)
				p.stroke(95, .5, portraitPoint{x, v + .03}, portraitPoint{x + .022, v + .03})
			}
		}
		neon := []int{81, 213, 220}[i]
		filmEarthRect(p, 233, u+.005, .56, u+.135, .625)
		p.stroke(neon, 1, portraitPoint{u + .005, .56}, portraitPoint{u + .135, .56})
		p.stroke(neon, 1, portraitPoint{u + .005, .625}, portraitPoint{u + .135, .625})
		filmEarthRect(p, 24, u+.012, .65, u+.087, .79)
		filmEarthRect(p, 110, u+.02, .66, u+.035, .77)
		p.stroke(153, .7, portraitPoint{u + .017, .66}, portraitPoint{u + .08, .76})
		filmEarthRect(p, 233, u+.1, .65, u+.13, .8)
		p.poly(neon, portraitPoint{u, .63}, portraitPoint{u + .14, .63}, portraitPoint{u + .155, .66}, portraitPoint{u - .008, .66})
		filmNativeSign(f, u+.018, .585, []string{"MARKET", "CAFE", "TOOLS"}[i], 252, 233)
	}
	// A vertical neon blade reads even when the shop lettering is too small.
	filmEarthRect(p, 53, .603, .24, .623, .52)
	p.stroke(213, 1, portraitPoint{.603, .24}, portraitPoint{.603, .52})
	for i := 0; i < 4; i++ {
		filmEarthRect(p, 225, .609, .27+float64(i)*.055, .617, .29+float64(i)*.055)
	}
	// Streetlamps, cones of light and silhouettes on the pavement.
	for _, u := range []float64{.43, .79} {
		p.poly(24, portraitPoint{u, .38}, portraitPoint{u - .04, .81}, portraitPoint{u + .065, .81})
		p.stroke(109, 1, portraitPoint{u, .83}, portraitPoint{u, .34}, portraitPoint{u + .025, .34})
		p.smoothOval(223, u+.025, .35, .012, .025)
		filmEarthResident(p, u+.045, .82)
	}
	filmEarthCar(p, .48, .91, 110)
	filmEarthCar(p, .71, .95, 167)
	// The delivery truck is visibly loaded rather than a blank white rectangle.
	filmEarthRect(p, 239, .87, .66, .97, .85)
	filmEarthRect(p, 252, .87, .66, .97, .69)
	filmEarthRect(p, 24, .87, .7, .95, .81)
	for i := 0; i < 4; i++ {
		v := .72 + float64(i)*.02
		filmEarthRect(p, 137, .88, v, .948, v+.012)
	}
	p.poly(110, portraitPoint{.97, .73}, portraitPoint{.99, .76}, portraitPoint{1, .84}, portraitPoint{.97, .85})
	p.smoothOval(233, .89, .86, .012, .028)
	p.smoothOval(233, .98, .86, .012, .028)
	filmEvacuation(p, .19, .85)
}

func filmEarthCar(p portraitPainter, u, v float64, color int) {
	p.poly(color, portraitPoint{u, v}, portraitPoint{u + .012, v - .04}, portraitPoint{u + .035, v - .09}, portraitPoint{u + .075, v - .09}, portraitPoint{u + .095, v - .035}, portraitPoint{u + .12, v - .02}, portraitPoint{u + .12, v + .012})
	p.poly(24, portraitPoint{u + .04, v - .08}, portraitPoint{u + .07, v - .08}, portraitPoint{u + .083, v - .04}, portraitPoint{u + .025, v - .04})
	p.stroke(153, .6, portraitPoint{u + .04, v - .08}, portraitPoint{u + .07, v - .08})
	p.smoothOval(233, u+.025, v+.01, .012, .027)
	p.smoothOval(233, u+.095, v+.01, .012, .027)
	p.smoothOval(223, u+.115, v-.018, .007, .013)
	p.poly(60, portraitPoint{u + .12, v - .02}, portraitPoint{u + .19, v - .01}, portraitPoint{u + .2, v + .035}, portraitPoint{u + .12, v + .005})
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
	// Three readable masses: tool display, complete truck silhouette, timber pile.
	// Large negative spaces and separated wheels survive the smallest viewport.
	filmEarthRect(p, 235, 0, 0, 1, 1)
	filmEarthRect(p, 236, 0, .03, .48, .71)
	filmEarthRect(p, 233, .48, 0, 1, .71)
	for row := 0; row < 5; row++ {
		v := .12 + float64(row)*.12
		p.stroke(95, .5, portraitPoint{0, v}, portraitPoint{.47, v})
	}
	// Recessed warm shop window beneath a projecting teal awning.
	filmEarthRect(p, 233, .035, .23, .445, .67)
	filmEarthRect(p, 94, .05, .245, .43, .655)
	filmEarthRect(p, 180, .065, .26, .415, .635)
	filmEarthRect(p, 223, .075, .275, .405, .61)
	filmEarthRect(p, 23, .025, .08, .455, .19)
	p.stroke(81, 1.2, portraitPoint{.025, .08}, portraitPoint{.455, .08})
	filmNativeSign(p.f, .125, .115, "HARDWARE", 252, 23)
	p.poly(30, portraitPoint{.025, .19}, portraitPoint{.455, .19}, portraitPoint{.48, .24}, portraitPoint{.005, .24})
	p.stroke(117, 1, portraitPoint{.005, .24}, portraitPoint{.48, .24})
	// Hammer: heavy metal head, slender wooden handle and hooked claw.
	filmEarthRect(p, 94, .135, .365, .15, .575)
	filmEarthRect(p, 252, .105, .325, .18, .37)
	p.poly(246, portraitPoint{.18, .325}, portraitPoint{.21, .35}, portraitPoint{.19, .385}, portraitPoint{.178, .36})
	// Handsaw: broad taper, unmistakable teeth and a dark cutout in its handle.
	p.poly(246, portraitPoint{.235, .36}, portraitPoint{.29, .34}, portraitPoint{.32, .56}, portraitPoint{.22, .56})
	for i := 0; i < 5; i++ {
		u := .22 + float64(i)*.02
		p.poly(252, portraitPoint{u, .55}, portraitPoint{u + .01, .58}, portraitPoint{u + .02, .55})
	}
	p.poly(130, portraitPoint{.235, .31}, portraitPoint{.285, .3}, portraitPoint{.3, .365}, portraitPoint{.24, .385})
	filmEarthRect(p, 233, .25, .325, .277, .35)
	// Open-ended wrench is visibly different from the hammer and saw.
	p.stroke(246, 3, portraitPoint{.36, .41}, portraitPoint{.365, .57})
	p.poly(252, portraitPoint{.335, .325}, portraitPoint{.348, .38}, portraitPoint{.372, .38}, portraitPoint{.385, .325}, portraitPoint{.4, .37}, portraitPoint{.38, .425}, portraitPoint{.345, .425}, portraitPoint{.323, .37})
	p.stroke(94, 1, portraitPoint{.065, .635}, portraitPoint{.415, .635})
	p.stroke(252, .6, portraitPoint{.082, .28}, portraitPoint{.105, .6})
	// Sidewalk and road recede into the loading alley behind the truck.
	p.poly(239, portraitPoint{0, .69}, portraitPoint{1, .63}, portraitPoint{1, .77}, portraitPoint{0, .83})
	p.stroke(252, 1, portraitPoint{0, .82}, portraitPoint{1, .76})
	filmEarthRect(p, 234, 0, .83, 1, 1)
	for i := 0; i < 7; i++ {
		u := float64(i) * .15
		p.stroke([]int{24, 60, 95}[i%3], 1, portraitPoint{u, .95}, portraitPoint{u + .07, .945})
	}
	// Truck faces right: open flatbed at left, raised cab and hood at right.
	// Both wheels, bumper, windshield and cargo rails fit fully within the shot.
	filmEarthRect(p, 233, .52, .74, .94, .84)
	filmEarthRect(p, 109, .51, .655, .775, .745)
	filmEarthRect(p, 153, .51, .655, .775, .678)
	filmEarthRect(p, 239, .52, .71, .775, .76)
	p.poly(110, portraitPoint{.775, .74}, portraitPoint{.775, .39}, portraitPoint{.87, .39}, portraitPoint{.91, .57}, portraitPoint{.955, .59}, portraitPoint{.955, .76})
	p.poly(153, portraitPoint{.775, .39}, portraitPoint{.87, .39}, portraitPoint{.88, .42}, portraitPoint{.775, .42})
	p.poly(24, portraitPoint{.793, .44}, portraitPoint{.853, .44}, portraitPoint{.885, .555}, portraitPoint{.793, .555})
	p.stroke(189, 1, portraitPoint{.805, .45}, portraitPoint{.832, .54})
	p.stroke(67, 1, portraitPoint{.785, .58}, portraitPoint{.785, .73})
	filmEarthRect(p, 252, .795, .595, .82, .61)
	filmEarthRect(p, 223, .932, .62, .955, .66)
	filmEarthRect(p, 252, .925, .73, .97, .76)
	for _, u := range []float64{.57, .87} {
		p.smoothOval(233, u, .79, .043, .103)
		p.smoothOval(239, u, .79, .025, .06)
		p.smoothOval(252, u, .79, .011, .028)
	}
	// Cargo: individually edged beams, secured with a dark strap.
	for i := 0; i < 4; i++ {
		v := .49 + float64(i)*.039
		p.poly(137, portraitPoint{.5, v}, portraitPoint{.75, v - .025}, portraitPoint{.75, v + .003}, portraitPoint{.5, v + .028})
		p.stroke(223, .8, portraitPoint{.5, v}, portraitPoint{.75, v - .025})
	}
	p.stroke(94, 2, portraitPoint{.635, .48}, portraitPoint{.645, .65})
	// Foreground timber has lit end grain and dark side faces: solid beams,
	// not lines floating in front of the display. Leave the truck wheels clear.
	for i := 0; i < 3; i++ {
		v := .81 + float64(i)*.05
		p.poly(94, portraitPoint{.065, v}, portraitPoint{.34, v - .08}, portraitPoint{.43, v - .045}, portraitPoint{.15, v + .04})
		p.poly(180, portraitPoint{.065, v}, portraitPoint{.34, v - .08}, portraitPoint{.41, v - .065}, portraitPoint{.135, v + .015})
		p.poly(137, portraitPoint{.065, v}, portraitPoint{.135, v + .015}, portraitPoint{.15, v + .04}, portraitPoint{.075, v + .025})
		p.stroke(223, .6, portraitPoint{.085, v + .005}, portraitPoint{.123, v + .016})
	}
	// A worker beside the bed gives the vehicle and material piles human scale.
	p.poly(214, portraitPoint{.46, .54}, portraitPoint{.49, .54}, portraitPoint{.505, .67}, portraitPoint{.445, .67})
	p.smoothOval(223, .476, .505, .018, .037)
	p.poly(220, portraitPoint{.452, .488}, portraitPoint{.455, .464}, portraitPoint{.49, .464}, portraitPoint{.5, .488})
	p.stroke(223, 2, portraitPoint{.49, .57}, portraitPoint{.525, .59})
	p.stroke(24, 3, portraitPoint{.46, .66}, portraitPoint{.45, .77})
	p.stroke(24, 3, portraitPoint{.49, .66}, portraitPoint{.495, .77})
}
