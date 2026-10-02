package render

import "math"

// The films share visual motifs, not identical slides: the bowl, the broken
// banner, the builder's mallet, an immense claw, and a small green hand.
func paintFilmShot(f *Frame, st CutsceneState, art string) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	ending := st.Film == FilmEnding
	t := math.Min(1, float64(st.Frame)/120)
	switch art {
	case "gate", "morning":
		filmGate(p, st.Frame, art == "morning")
	case "contract":
		filmContract(p)
	case "bowl":
		filmStillLife(p, st.Frame, ending)
	case "threshold":
		filmThreshold(p, st.Frame)
	case "dragon":
		filmDragonWide(p, st.Frame, ending)
	case "grak":
		filmGrakClose(f, ending)
	case "eye":
		filmDragonFace(p, st.Frame, ending, ending && st.Shot >= 6)
	case "offering":
		filmOffering(p, t)
	case "together", "rest":
		filmTogether(p, st.Frame, ending, art == "rest")
	case "door":
		filmDoor(p, st.Frame)
	case "mallet":
		filmMallet(p, t)
	case "resolve":
		filmResolve(p, st.Frame)
	case "aftermath":
		filmAftermath(p)
	case "ember":
		filmEmber(p, st.Frame, st.Shot >= 7)
	case "touch":
		filmTouch(p, t)
	}
}

func filmGate(p portraitPainter, frame int, morning bool) {
	poly, oval := p.poly, p.smoothOval
	sky, sun, light := 235, 131, 95
	if morning {
		sky, sun, light = 237, 223, 109
	}
	poly(sky, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	oval(sun, .77, .27, .13, .2)
	// Angular mountain layers. The gate is a single pointed mouth in the cliff.
	poly(236, portraitPoint{0, .48}, portraitPoint{.13, .21}, portraitPoint{.26, .43}, portraitPoint{.45, .08}, portraitPoint{.67, .43}, portraitPoint{.83, .35}, portraitPoint{1, .52}, portraitPoint{1, 1}, portraitPoint{0, 1})
	poly(238, portraitPoint{0, .66}, portraitPoint{.12, .42}, portraitPoint{.25, .52}, portraitPoint{.45, .15}, portraitPoint{.58, .46}, portraitPoint{.8, .62}, portraitPoint{1, .54}, portraitPoint{1, 1}, portraitPoint{0, 1})
	poly(light, portraitPoint{.45, .15}, portraitPoint{.58, .46}, portraitPoint{.67, .58}, portraitPoint{.55, .55}, portraitPoint{.47, .29}, portraitPoint{.4, .43})
	poly(239, portraitPoint{.26, .85}, portraitPoint{.29, .47}, portraitPoint{.43, .32}, portraitPoint{.56, .46}, portraitPoint{.6, .85})
	poly(244, portraitPoint{.29, .78}, portraitPoint{.32, .49}, portraitPoint{.43, .37}, portraitPoint{.47, .43}, portraitPoint{.35, .52}, portraitPoint{.32, .8})
	poly(233, portraitPoint{.34, .84}, portraitPoint{.355, .52}, portraitPoint{.43, .43}, portraitPoint{.50, .53}, portraitPoint{.53, .84})
	if morning {
		poly(180, portraitPoint{.35, .82}, portraitPoint{.365, .54}, portraitPoint{.43, .455}, portraitPoint{.442, .48}, portraitPoint{.375, .57}, portraitPoint{.373, .82})
	}
	poly(236, portraitPoint{.35, .77}, portraitPoint{.53, .77}, portraitPoint{.87, 1}, portraitPoint{.14, 1})
	for i := 0; i < 5; i++ {
		v := .8 + float64(i)*.04
		half := .1 + (v-.8)*1.2
		p.stroke(240, 1.1, portraitPoint{.44 - half, v}, portraitPoint{.44 + half, v})
	}
	// The guild's procession climbs from the right; after victory, two kin
	// stand in the doorway and the severed banner lies on the far bank.
	if !morning {
		for i := 0; i < 5; i++ {
			u := .64 + float64(i)*.062
			v := .81 + float64(i)*.033
			filmSmallHero(p, u, v, .032, i == 0, frame)
		}
	} else {
		filmSmallKin(p, .405, .75, .024)
		filmSmallKin(p, .472, .76, .022)
		p.stroke(94, 3, portraitPoint{.7, .87}, portraitPoint{.88, .95})
		poly(95, portraitPoint{.725, .87}, portraitPoint{.79, .87}, portraitPoint{.84, .935}, portraitPoint{.79, .92}, portraitPoint{.75, .94})
		p.stroke(65, 1.2, portraitPoint{.15, .84}, portraitPoint{.15, .73}, portraitPoint{.125, .71})
		oval(151, .125, .71, .026, .015)
	}
}

func filmSmallHero(p portraitPainter, u, v, s float64, banner bool, frame int) {
	poly := p.poly
	poly(95, portraitPoint{u - s, v}, portraitPoint{u - s*.7, v - s*2.8}, portraitPoint{u + s*.4, v - s*2.8}, portraitPoint{u + s, v})
	p.smoothOval(246, u, v-s*3.2, s*.65, s*.8)
	poly(239, portraitPoint{u - s, v - s*2.2}, portraitPoint{u, v - s*2.3}, portraitPoint{u + s*.2, v - s*.7}, portraitPoint{u - s*.7, v - s*.45})
	p.stroke(180, 1, portraitPoint{u + s*.7, v - s*3}, portraitPoint{u + s*.7, v})
	if banner {
		p.stroke(94, 1.6, portraitPoint{u - s*.2, v}, portraitPoint{u - s*.2, v - s*6.5})
		sway := .009 * math.Sin(float64(frame)/28)
		poly(131, portraitPoint{u - s*.2, v - s*6.3}, portraitPoint{u + s*2.8 + sway, v - s*6}, portraitPoint{u + s*2.5, v - s*4.8}, portraitPoint{u - s*.2, v - s*5.1})
		p.stroke(223, 1, portraitPoint{u + s*.5, v - s*5.8}, portraitPoint{u + s*1.9, v - s*5.1})
	}
}
func filmSmallKin(p portraitPainter, u, v, s float64) {
	p.poly(65, portraitPoint{u - s, v}, portraitPoint{u - s*.8, v - s*2}, portraitPoint{u + s*.8, v - s*2}, portraitPoint{u + s, v})
	p.poly(107, portraitPoint{u - s*1.2, v - s*3}, portraitPoint{u - s*.4, v - s*2.7}, portraitPoint{u + s*.4, v - s*2.7}, portraitPoint{u + s*1.2, v - s*3.1}, portraitPoint{u + s*.6, v - s*2}, portraitPoint{u - s*.6, v - s*2})
}

func filmContract(p portraitPainter) {
	poly := p.poly
	poly(236, portraitPoint{.08, .14}, portraitPoint{.87, .1}, portraitPoint{.93, .86}, portraitPoint{.12, .93})
	poly(137, portraitPoint{.16, .09}, portraitPoint{.8, .13}, portraitPoint{.84, .82}, portraitPoint{.12, .86})
	poly(223, portraitPoint{.18, .12}, portraitPoint{.78, .16}, portraitPoint{.81, .79}, portraitPoint{.145, .825})
	// Printed dragon: wing, horn, muzzle, curled tail. A sword crosses it.
	poly(94, portraitPoint{.36, .36}, portraitPoint{.52, .22}, portraitPoint{.64, .48}, portraitPoint{.52, .6}, portraitPoint{.38, .55}, portraitPoint{.25, .47}, portraitPoint{.25, .42}, portraitPoint{.35, .42})
	poly(94, portraitPoint{.32, .43}, portraitPoint{.34, .3}, portraitPoint{.38, .4})
	p.stroke(94, 5, portraitCurve(portraitPoint{.55, .52}, portraitPoint{.76, .8}, portraitPoint{.25, .78}, portraitPoint{.3, .57})...)
	p.stroke(233, 3, portraitPoint{.35, .63}, portraitPoint{.61, .27})
	p.stroke(233, 3, portraitPoint{.35, .5}, portraitPoint{.49, .61})
	for i := 0; i < 3; i++ {
		v := .66 + float64(i)*.048
		p.stroke(137, 1.6, portraitPoint{.19, v}, portraitPoint{.38 - float64(i)*.02, v})
	}
	// Wax seal catches the eye against the ivory commission.
	p.smoothOval(52, .69, .7, .085, .12)
	p.smoothOval(131, .685, .68, .07, .096)
	p.stroke(223, 2, portraitPoint{.655, .7}, portraitPoint{.68, .625}, portraitPoint{.716, .698})
	p.stroke(223, 1.3, portraitPoint{.663, .675}, portraitPoint{.707, .674})
	filmCoin(p, .85, .88, .054)
	filmCoin(p, .915, .82, .036)
	// Cold, square gauntlet on the paper's edge.
	poly(239, portraitPoint{0, .47}, portraitPoint{.16, .43}, portraitPoint{.22, .5}, portraitPoint{.22, .59}, portraitPoint{.15, .63}, portraitPoint{0, .62})
	for i := 0; i < 3; i++ {
		v := .475 + float64(i)*.05
		p.stroke(246, 2, portraitPoint{.1, v}, portraitPoint{.195, v + .017})
	}
}
func filmCoin(p portraitPainter, u, v, r float64) {
	p.smoothOval(94, u, v+.012, r, r*.7)
	p.smoothOval(178, u, v, r, r*.7)
	p.smoothOval(223, u-r*.12, v-r*.17, r*.75, r*.35)
	p.stroke(130, 1.5, portraitPoint{u - r*.2, v - r*.23}, portraitPoint{u + r*.2, v + r*.23})
}

func filmStillLife(p portraitPainter, frame int, ending bool) {
	poly := p.poly
	poly(236, portraitPoint{0, .72}, portraitPoint{.73, .5}, portraitPoint{1, .66}, portraitPoint{1, 1}, portraitPoint{0, 1})
	poly(94, portraitPoint{.66, .39}, portraitPoint{.79, .3}, portraitPoint{1, .36}, portraitPoint{1, .78}, portraitPoint{.82, .79}, portraitPoint{.68, .7})
	poly(131, portraitPoint{.66, .4}, portraitPoint{.79, .31}, portraitPoint{1, .37}, portraitPoint{1, .64}, portraitPoint{.82, .69}, portraitPoint{.69, .6})
	poly(173, portraitPoint{.69, .42}, portraitPoint{.79, .35}, portraitPoint{.98, .4}, portraitPoint{.96, .44}, portraitPoint{.78, .4})
	for i := 0; i < 3; i++ {
		u := .73 + float64(i)*.075
		poly(180, portraitPoint{u, .65}, portraitPoint{u + .065, .67}, portraitPoint{u + .02, .78})
	}
	for i := 0; i < 7; i++ {
		u := .08 + float64(i)*.08
		v := .78 - float64(i%3)*.06
		filmCoin(p, u, v, .034)
	}
	filmBowl(p, .42, .51, .265, frame, true)
	// One small drop. Its downward travel makes the leaky bowl a literal detail.
	if ending {
		t := float64(frame%75) / 75
		p.smoothOval(109, .66, .665+t*.13, .005, .012)
	}
}
func filmBowl(p portraitPainter, u, v, r float64, frame int, cracked bool) {
	poly := p.poly
	p.smoothOval(236, u, v+r*.65, r*1.05, r*.23)
	poly(239, portraitPoint{u - r, v}, portraitPoint{u + r, v}, portraitPoint{u + r*.76, v + r*.68}, portraitPoint{u - r*.62, v + r*.68})
	poly(246, portraitPoint{u - r*.94, v + .02}, portraitPoint{u - r*.15, v + r*.21}, portraitPoint{u - r*.19, v + r*.6}, portraitPoint{u - r*.63, v + r*.6})
	p.smoothOval(246, u, v, r, r*.32)
	p.smoothOval(23, u, v, r*.84, r*.24)
	p.smoothOval(109, u-.025, v-.007, r*.71, r*.15)
	ripple := .008 * math.Sin(float64(frame)/23)
	p.stroke(153, 1, portraitPoint{u - r*.55, v + ripple}, portraitPoint{u - r*.03, v + ripple})
	p.stroke(117, 1, portraitPoint{u + r*.1, v + .025 - ripple}, portraitPoint{u + r*.54, v + .025 - ripple})
	if cracked {
		p.stroke(236, 1.5, portraitPoint{u + r*.68, v + r*.2}, portraitPoint{u + r*.51, v + r*.4}, portraitPoint{u + r*.62, v + r*.6})
	}
}

func filmThreshold(p portraitPainter, frame int) {
	poly := p.poly
	drawDragonTableau(p.f, 85, 8, 153, 65, 6, 233, owPadView{status: OWOpen}, frame, true)
	// Foreground jamb frames the small visitor; the dragon fills the far room.
	poly(236, portraitPoint{0, 0}, portraitPoint{.15, 0}, portraitPoint{.13, 1}, portraitPoint{0, 1})
	poly(240, portraitPoint{.14, 0}, portraitPoint{.18, 0}, portraitPoint{.17, 1}, portraitPoint{.13, 1})
	poly(236, portraitPoint{.18, .87}, portraitPoint{.83, .9}, portraitPoint{1, 1}, portraitPoint{.13, 1})
	drawGrakPortrait(p.f, 22, 22, 94, 55)
	filmGrakEyes(p, 22, 22, 94, 55)
}
func filmDragonWide(p portraitPainter, frame int, ending bool) {
	breath := .006 * math.Sin(float64(frame)/43)
	drawDragonTableau(p.f, 9, 3, 222, 73, 8, 233, owPadView{status: OWOpen}, frame, true)
	// Only the soft wing fold moves with his breathing; masonry stays still.
	p.stroke(131, 3, portraitPoint{.665, .34}, portraitPoint{.70, .43 + breath}, portraitPoint{.735, .52 + breath})
	filmBowl(p, .16, .84, .09, frame, true)
	p.stroke(180, 4, portraitPoint{.243, .46}, portraitPoint{.27, .46})
	if ending {
		p.smoothOval(223, .253, .46, .018, .009)
		p.smoothOval(233, .254, .46, .006, .008)
	}
}

func filmGrakClose(f *Frame, ending bool) {
	// A true close-up of the established face. The mallet and feet fall out
	// of the frame; blade ears and uneven tusks retain the same character.
	head := filmFrame(280, 140)
	drawGrakPortrait(head, 0, 0, 280, 140)
	if ending {
		p := portraitPainter{head, 0, 0, 280, 140}
		p.stroke(95, 1.5, portraitPoint{.337, .3}, portraitPoint{.351, .324})
		p.poly(58, portraitPoint{.355, .266}, portraitPoint{.462, .274}, portraitPoint{.455, .281}, portraitPoint{.358, .278})
	}
	cinemaProject(f, head, Rect{0, 0, f.W, f.H}, filmView{.085, .095, .84, .34})
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.stroke(223, 4, portraitPoint{.37, .56}, portraitPoint{.402, .575})
	p.stroke(223, 4, portraitPoint{.605, .565}, portraitPoint{.64, .545})
}

func filmDragonFace(p portraitPainter, frame int, ending, awake bool) {
	poly := p.poly
	breath := .007 * math.Sin(float64(frame)/43)
	// Snout left, massive brow and cheek right; the horns lean back exactly
	// as in the curled portrait, enlarged into separate ivory planes.
	poly(95, portraitPoint{.1, .56}, portraitPoint{.3, .44}, portraitPoint{.34, .3}, portraitPoint{.7, .23}, portraitPoint{.92, .42}, portraitPoint{1, .91}, portraitPoint{.49, .96}, portraitPoint{.32, .77}, portraitPoint{.085, .74}, portraitPoint{.04, .64})
	poly(131, portraitPoint{.07, .58 + breath}, portraitPoint{.3, .48}, portraitPoint{.34, .33}, portraitPoint{.68, .27}, portraitPoint{.85, .42}, portraitPoint{.82, .66}, portraitPoint{.62, .81}, portraitPoint{.33, .68}, portraitPoint{.08, .69})
	poly(173, portraitPoint{.34, .34}, portraitPoint{.68, .28}, portraitPoint{.78, .39}, portraitPoint{.52, .37}, portraitPoint{.31, .51}, portraitPoint{.16, .57}, portraitPoint{.09, .58}, portraitPoint{.3, .46})
	poly(180, portraitPoint{.41, .34}, portraitPoint{.61, .02}, portraitPoint{.57, .24}, portraitPoint{.49, .35})
	poly(223, portraitPoint{.43, .32}, portraitPoint{.57, .08}, portraitPoint{.51, .27})
	poly(180, portraitPoint{.65, .3}, portraitPoint{.85, .07}, portraitPoint{.8, .32}, portraitPoint{.71, .41})
	poly(137, portraitPoint{.72, .34}, portraitPoint{.85, .07}, portraitPoint{.8, .32})
	poly(95, portraitPoint{.36, .46}, portraitPoint{.52, .4}, portraitPoint{.62, .45}, portraitPoint{.52, .49}, portraitPoint{.4, .5})
	eyeH := .009
	if awake {
		eyeH = .022 + math.Min(1, float64(frame)/90)*.008
	}
	if ending && !awake {
		eyeH = .006
	}
	p.smoothOval(223, .457, .486, .073, eyeH)
	if awake {
		p.smoothOval(52, .46, .486, .009, eyeH*.9)
		p.smoothOval(230, .433, .478, .012, .004)
	}
	p.smoothOval(52, .15, .618, .02, .012)
	p.stroke(95, 2, portraitPoint{.12, .693}, portraitPoint{.28, .711}, portraitPoint{.4, .691})
	poly(180, portraitPoint{.19, .7}, portraitPoint{.22, .7}, portraitPoint{.205, .75})
	p.stroke(131, 2, portraitPoint{.8, .56}, portraitPoint{.78, .64}, portraitPoint{.7, .69})
}

func filmOffering(p portraitPainter, t float64) {
	poly := p.poly
	filmBowl(p, .68, .68, .215, int(t*120), true)
	// A heavy sleeve, diagonal forearm, broad palm, three articulated fingers.
	shift := (1 - t) * .035
	poly(95, portraitPoint{0, .1}, portraitPoint{.24, .18}, portraitPoint{.32, .34}, portraitPoint{.22, .5}, portraitPoint{0, .42})
	poly(65, portraitPoint{.14, .26}, portraitPoint{.26, .3}, portraitPoint{.44, .55 + shift}, portraitPoint{.35, .67 + shift}, portraitPoint{.18, .47})
	poly(107, portraitPoint{.19, .29}, portraitPoint{.27, .32}, portraitPoint{.41, .54 + shift}, portraitPoint{.36, .58 + shift}, portraitPoint{.23, .43})
	poly(107, portraitPoint{.35, .53 + shift}, portraitPoint{.46, .52 + shift}, portraitPoint{.51, .58 + shift}, portraitPoint{.48, .7 + shift}, portraitPoint{.34, .69 + shift}, portraitPoint{.29, .61 + shift})
	for i := 0; i < 3; i++ {
		u := .37 + float64(i)*.036
		poly(150, portraitPoint{u, .6 + shift}, portraitPoint{u + .03, .59 + shift}, portraitPoint{u + .024, .74 + shift}, portraitPoint{u, .74 + shift})
	}
	filmCoin(p, .46, .805, .045)
	p.stroke(65, 1.5, portraitPoint{.356, .594 + shift}, portraitPoint{.4, .66 + shift})
}

func filmTogether(p portraitPainter, frame int, ending, rest bool) {
	poly := p.poly
	// Establish both figures in the same space. The gold is background; the
	// small bowl bridges their silhouettes in the foreground.
	drawDragonTableau(p.f, 64, 0, 173, 70, 6, 233, owPadView{status: OWOpen}, frame, true)
	poly(236, portraitPoint{0, .89}, portraitPoint{1, .89}, portraitPoint{1, 1}, portraitPoint{0, 1})
	if !rest {
		drawGrakPortrait(p.f, 0, 12, 102, 64)
		filmGrakEyes(p, 0, 12, 102, 64)
	} else {
		filmSeatedGrak(p)
	}
	filmBowl(p, .46, .86, .085, frame, true)
	p.stroke(180, 3, portraitPoint{.427, .405}, portraitPoint{.453, .405})
	if ending {
		p.smoothOval(223, .43, .405, .016, .013)
		p.smoothOval(233, .432, .405, .004, .005)
	}
	if rest {
		// Light from the open door falls across the floor, not a magical cure.
		poly(239, portraitPoint{.05, .94}, portraitPoint{.45, .98}, portraitPoint{.53, 1}, portraitPoint{0, 1})
		p.stroke(109, 1, portraitPoint{.4, .952}, portraitPoint{.57, .952})
	}
}

func filmSeatedGrak(p portraitPainter) {
	poly := p.poly
	// A new resting pose: mantle settles onto a crate, knees forward, hand
	// loose on the mallet laid across them. Face is lifted toward the dragon.
	poly(94, portraitPoint{.09, .74}, portraitPoint{.26, .74}, portraitPoint{.26, .96}, portraitPoint{.085, .96})
	p.stroke(137, 1.3, portraitPoint{.1, .77}, portraitPoint{.245, .77})
	poly(95, portraitPoint{.12, .45}, portraitPoint{.27, .45}, portraitPoint{.32, .61}, portraitPoint{.29, .86}, portraitPoint{.09, .88}, portraitPoint{.085, .69})
	poly(65, portraitPoint{.14, .5}, portraitPoint{.255, .49}, portraitPoint{.31, .64}, portraitPoint{.25, .75}, portraitPoint{.13, .73})
	poly(107, portraitPoint{.15, .51}, portraitPoint{.235, .5}, portraitPoint{.267, .61}, portraitPoint{.14, .61})
	poly(58, portraitPoint{.16, .72}, portraitPoint{.34, .72}, portraitPoint{.345, .82}, portraitPoint{.25, .83}, portraitPoint{.25, .94}, portraitPoint{.18, .94})
	poly(58, portraitPoint{.24, .76}, portraitPoint{.39, .78}, portraitPoint{.36, .9}, portraitPoint{.33, .96}, portraitPoint{.27, .96}, portraitPoint{.29, .86})
	poly(94, portraitPoint{.175, .89}, portraitPoint{.253, .89}, portraitPoint{.28, .955}, portraitPoint{.17, .96})
	poly(94, portraitPoint{.29, .92}, portraitPoint{.37, .92}, portraitPoint{.4, .965}, portraitPoint{.28, .965})
	// Blade ears and uneven tusks match Grak's portrait, at a new angle.
	poly(107, portraitPoint{.13, .31}, portraitPoint{.285, .3}, portraitPoint{.315, .37}, portraitPoint{.29, .49}, portraitPoint{.165, .48}, portraitPoint{.12, .4})
	poly(107, portraitPoint{.07, .3}, portraitPoint{.155, .34}, portraitPoint{.165, .4}, portraitPoint{.115, .38})
	poly(107, portraitPoint{.27, .32}, portraitPoint{.355, .26}, portraitPoint{.328, .355}, portraitPoint{.29, .4})
	poly(150, portraitPoint{.155, .315}, portraitPoint{.265, .31}, portraitPoint{.282, .35}, portraitPoint{.14, .36})
	poly(65, portraitPoint{.29, .355}, portraitPoint{.317, .375}, portraitPoint{.29, .49}, portraitPoint{.25, .45})
	poly(150, portraitPoint{.23, .375}, portraitPoint{.27, .385}, portraitPoint{.315, .431}, portraitPoint{.235, .43})
	poly(65, portraitPoint{.239, .433}, portraitPoint{.312, .433}, portraitPoint{.298, .452}, portraitPoint{.243, .451})
	p.stroke(58, 5, portraitPoint{.16, .388}, portraitPoint{.205, .397})
	p.stroke(223, 3, portraitPoint{.17, .405}, portraitPoint{.196, .41})
	p.stroke(58, 5, portraitPoint{.25, .39}, portraitPoint{.28, .38})
	p.stroke(223, 3, portraitPoint{.252, .41}, portraitPoint{.282, .4})
	poly(65, portraitPoint{.15, .433}, portraitPoint{.29, .433}, portraitPoint{.29, .485}, portraitPoint{.165, .48})
	p.stroke(58, 3, portraitPoint{.17, .451}, portraitPoint{.281, .451})
	poly(180, portraitPoint{.17, .46}, portraitPoint{.192, .46}, portraitPoint{.174, .421})
	poly(180, portraitPoint{.265, .46}, portraitPoint{.29, .458}, portraitPoint{.29, .395})
	poly(107, portraitPoint{.255, .56}, portraitPoint{.3, .6}, portraitPoint{.32, .7}, portraitPoint{.275, .73}, portraitPoint{.255, .68})
	p.stroke(94, 4, portraitPoint{.09, .68}, portraitPoint{.385, .75})
	poly(239, portraitPoint{.06, .63}, portraitPoint{.12, .64}, portraitPoint{.13, .73}, portraitPoint{.055, .73})
	poly(246, portraitPoint{.06, .63}, portraitPoint{.12, .64}, portraitPoint{.12, .66}, portraitPoint{.06, .65})
	poly(107, portraitPoint{.285, .68}, portraitPoint{.326, .69}, portraitPoint{.329, .74}, portraitPoint{.29, .75})
}

func filmDoor(p portraitPainter, frame int) {
	poly := p.poly
	// Reverse angle: the narrow light beyond the shut gate carries moving
	// human shadows. Grak sees exactly what the first exterior shot promised.
	poly(236, portraitPoint{.12, 1}, portraitPoint{.14, .23}, portraitPoint{.49, .03}, portraitPoint{.86, .23}, portraitPoint{.89, 1})
	poly(239, portraitPoint{.17, 1}, portraitPoint{.2, .26}, portraitPoint{.49, .08}, portraitPoint{.81, .27}, portraitPoint{.84, 1})
	poly(94, portraitPoint{.24, 1}, portraitPoint{.26, .31}, portraitPoint{.49, .15}, portraitPoint{.75, .33}, portraitPoint{.77, 1})
	for i := 0; i < 7; i++ {
		u := .28 + float64(i)*.067
		p.stroke(58, 2, portraitPoint{u, .36}, portraitPoint{u, .96})
	}
	poly(223, portraitPoint{.49, .165}, portraitPoint{.514, .179}, portraitPoint{.53, 1}, portraitPoint{.488, 1})
	poly(137, portraitPoint{.49, .88}, portraitPoint{.53, .88}, portraitPoint{.83, 1}, portraitPoint{.19, 1})
	moving := .015 * math.Sin(float64(frame)/52)
	poly(233, portraitPoint{.49, .49 + moving}, portraitPoint{.513, .49 + moving}, portraitPoint{.524, .77 + moving}, portraitPoint{.49, .77 + moving})
	p.stroke(239, 4, portraitPoint{.25, .65}, portraitPoint{.73, .68})
	p.smoothOval(94, .46, .66, .035, .05)
	p.smoothOval(180, .455, .65, .015, .02)
}

func filmMallet(p portraitPainter, t float64) {
	poly := p.poly
	poly(236, portraitPoint{.04, .78}, portraitPoint{.86, .76}, portraitPoint{1, 1}, portraitPoint{0, 1})
	p.stroke(94, 11, portraitPoint{.64, .17}, portraitPoint{.27, .94})
	p.stroke(180, 2, portraitPoint{.62, .19}, portraitPoint{.26, .91})
	poly(238, portraitPoint{.39, .16}, portraitPoint{.51, .045}, portraitPoint{.88, .32}, portraitPoint{.91, .43}, portraitPoint{.8, .57}, portraitPoint{.67, .57})
	poly(246, portraitPoint{.4, .17}, portraitPoint{.51, .055}, portraitPoint{.87, .325}, portraitPoint{.78, .39})
	poly(240, portraitPoint{.78, .39}, portraitPoint{.87, .325}, portraitPoint{.9, .43}, portraitPoint{.8, .56}, portraitPoint{.73, .5})
	shift := (1 - t) * .06
	poly(95, portraitPoint{1, .54}, portraitPoint{.81, .58}, portraitPoint{.8, .84}, portraitPoint{1, .9})
	poly(65, portraitPoint{.84, .61}, portraitPoint{.59, .57 + shift}, portraitPoint{.48, .72 + shift}, portraitPoint{.55, .82 + shift}, portraitPoint{.83, .77})
	poly(107, portraitPoint{.78, .6}, portraitPoint{.58, .57 + shift}, portraitPoint{.5, .7 + shift}, portraitPoint{.6, .74 + shift}, portraitPoint{.83, .69})
	for i := 0; i < 3; i++ {
		v := .625 + float64(i)*.048 + shift
		p.stroke(150, 5, portraitPoint{.46, v}, portraitPoint{.57, v + .035})
	}
	p.stroke(65, 1.4, portraitPoint{.52, .65 + shift}, portraitPoint{.56, .755 + shift})
}
func filmResolve(p portraitPainter, frame int) {
	filmDoor(p, frame)
	// Lower the background into silhouette so the builder is the final mass.
	var colors [256]int
	for c := range colors {
		colors[c] = filmFade(c, .45)
	}
	for i, c := range p.f.C {
		if c.FG != 233 {
			c.FG = colors[c.FG]
			p.f.C[i] = c
		}
	}
	drawGrakPortrait(p.f, 68, 0, 126, 79)
	filmGrakEyes(p, 68, 0, 126, 79)
	journalCenter(p.f, 75, "G R A K", 180, true)
}

func filmAftermath(p portraitPainter) {
	drawJournalHalls(p.f, 55, -5, 130, 83)
	p.poly(236, portraitPoint{0, .79}, portraitPoint{1, .79}, portraitPoint{1, 1}, portraitPoint{0, 1})
	// Guild crest face-down, a shield split at its rim, a spent arrow.
	p.stroke(94, 5, portraitPoint{.14, .62}, portraitPoint{.78, .95})
	p.poly(95, portraitPoint{.22, .665}, portraitPoint{.41, .71}, portraitPoint{.38, .84}, portraitPoint{.31, .79}, portraitPoint{.26, .845})
	p.poly(239, portraitPoint{.66, .76}, portraitPoint{.85, .7}, portraitPoint{.92, .8}, portraitPoint{.82, .94}, portraitPoint{.69, .92})
	p.poly(137, portraitPoint{.675, .77}, portraitPoint{.84, .72}, portraitPoint{.87, .79}, portraitPoint{.80, .89}, portraitPoint{.71, .88})
	p.stroke(233, 3, portraitPoint{.78, .73}, portraitPoint{.76, .81}, portraitPoint{.815, .865}, portraitPoint{.795, .918})
	p.stroke(180, 1.5, portraitPoint{.14, .91}, portraitPoint{.52, .85})
	p.poly(246, portraitPoint{.52, .85}, portraitPoint{.48, .865}, portraitPoint{.5, .83})
}
func filmEmber(p portraitPainter, frame int, steady bool) {
	drawJournalHeart(p.f, 20, -1, 200, 82)
	// The same ribs carry a weak fire, then the rhythm of breathing. A warm
	// core brightens without erasing the dark stone around it.
	pulse := math.Sin(float64(frame) / 19)
	glow := .03
	if steady {
		glow = .055 + .007*pulse
	}
	p.smoothOval(173, .5, .535, .047+glow, .085+glow)
	p.smoothOval(215, .5, .54, .035+glow*.35, .066+glow*.4)
	p.smoothOval(223, .498, .552, .015, .045)
	for i := 0; i < 4; i++ {
		t := float64((frame+i*41)%180) / 180
		u := .48 + .04*math.Sin(float64(i)*2)
		p.smoothOval(180, u, .57-t*.24, .003, .005)
	}
}
func filmTouch(p portraitPainter, t float64) {
	poly := p.poly
	poly(236, portraitPoint{0, .86}, portraitPoint{1, .86}, portraitPoint{1, 1}, portraitPoint{0, 1})
	// Dragon claw is enormous and still. The small hand travels down to rest
	// between its knuckles. No weapon or gold occupies this composition.
	poly(95, portraitPoint{1, .3}, portraitPoint{.73, .3}, portraitPoint{.54, .45}, portraitPoint{.43, .59}, portraitPoint{.45, .8}, portraitPoint{.88, .9}, portraitPoint{1, .86})
	poly(131, portraitPoint{1, .35}, portraitPoint{.74, .35}, portraitPoint{.57, .48}, portraitPoint{.47, .61}, portraitPoint{.5, .75}, portraitPoint{.88, .82}, portraitPoint{1, .76})
	poly(173, portraitPoint{.76, .37}, portraitPoint{.58, .5}, portraitPoint{.49, .61}, portraitPoint{.6, .59}, portraitPoint{.82, .43})
	for i := 0; i < 3; i++ {
		u := .53 + float64(i)*.13
		poly(180, portraitPoint{u, .76}, portraitPoint{u + .09, .78}, portraitPoint{u + .045, .9})
	}
	dy := (1 - t) * .075
	poly(95, portraitPoint{0, .09}, portraitPoint{.22, .11}, portraitPoint{.26, .32 - dy}, portraitPoint{.08, .38 - dy}, portraitPoint{0, .31 - dy})
	poly(65, portraitPoint{.12, .27 - dy}, portraitPoint{.26, .28 - dy}, portraitPoint{.46, .5 - dy}, portraitPoint{.37, .63 - dy}, portraitPoint{.18, .43 - dy})
	poly(107, portraitPoint{.16, .28 - dy}, portraitPoint{.25, .3 - dy}, portraitPoint{.45, .5 - dy}, portraitPoint{.39, .555 - dy}, portraitPoint{.2, .4 - dy})
	poly(107, portraitPoint{.37, .49 - dy}, portraitPoint{.5, .48 - dy}, portraitPoint{.58, .56 - dy}, portraitPoint{.53, .68 - dy}, portraitPoint{.37, .65 - dy}, portraitPoint{.33, .58 - dy})
	for i := 0; i < 3; i++ {
		u := .425 + float64(i)*.037
		p.stroke(150, 5, portraitPoint{u, .56 - dy}, portraitPoint{u + .05, .65 - dy})
	}
	p.stroke(65, 1.5, portraitPoint{.4, .53 - dy}, portraitPoint{.41, .6 - dy})
}

// Filled eye marks survive camera sampling where single native glyphs would
// disappear between rows. They retain the established brow and tusk shapes.
func filmGrakEyes(p portraitPainter, x, y, w, h int) {
	for _, u := range []float64{.405, .60} {
		cx := (float64(x) + u*float64(w-1)) / float64(p.w-1)
		cy := (float64(y) + .282*float64(h-1)) / float64(p.h-1)
		p.stroke(180, 3.5, portraitPoint{cx - .006, cy}, portraitPoint{cx + .008, cy})
	}
}
