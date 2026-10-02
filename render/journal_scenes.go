package render

import "math"

func drawJournalRotunda(f *Frame, x, y, w, h int) {
	drawDragonTableau(f, x, y, w, h, max(1, h/7), 233, owPadView{status: OWOpen}, 90, true)
}

func drawJournalRift(f *Frame, x, y, w, h int) {
	p := portraitPainter{f, x, y, w, h}
	poly := p.poly
	// Two severed bridge shoulders frame a single incandescent wound.
	poly(236, portraitPoint{.05, .85}, portraitPoint{.12, .36}, portraitPoint{.3, .16}, portraitPoint{.43, .25}, portraitPoint{.405, .43}, portraitPoint{.48, .52}, portraitPoint{.4, .66}, portraitPoint{.46, .79}, portraitPoint{.34, .94})
	poly(239, portraitPoint{.13, .81}, portraitPoint{.18, .38}, portraitPoint{.3, .2}, portraitPoint{.39, .28}, portraitPoint{.37, .44}, portraitPoint{.435, .53}, portraitPoint{.36, .67}, portraitPoint{.4, .79}, portraitPoint{.32, .88})
	poly(246, portraitPoint{.18, .38}, portraitPoint{.3, .2}, portraitPoint{.32, .24}, portraitPoint{.25, .51}, portraitPoint{.18, .7})
	poly(236, portraitPoint{.58, .19}, portraitPoint{.74, .12}, portraitPoint{.88, .38}, portraitPoint{.95, .87}, portraitPoint{.65, .94}, portraitPoint{.57, .81}, portraitPoint{.63, .67}, portraitPoint{.55, .52}, portraitPoint{.62, .4})
	poly(239, portraitPoint{.63, .22}, portraitPoint{.73, .18}, portraitPoint{.83, .4}, portraitPoint{.89, .81}, portraitPoint{.68, .88}, portraitPoint{.63, .8}, portraitPoint{.69, .66}, portraitPoint{.6, .51}, portraitPoint{.67, .4})
	poly(109, portraitPoint{.73, .18}, portraitPoint{.83, .4}, portraitPoint{.89, .81}, portraitPoint{.81, .76}, portraitPoint{.75, .42}, portraitPoint{.68, .27})
	poly(52, portraitPoint{.46, .24}, portraitPoint{.51, .37}, portraitPoint{.48, .45}, portraitPoint{.55, .55}, portraitPoint{.48, .7}, portraitPoint{.57, .81}, portraitPoint{.52, .95}, portraitPoint{.43, .8}, portraitPoint{.46, .67}, portraitPoint{.43, .51}, portraitPoint{.46, .42})
	poly(173, portraitPoint{.49, .35}, portraitPoint{.505, .43}, portraitPoint{.493, .48}, portraitPoint{.53, .55}, portraitPoint{.472, .7}, portraitPoint{.54, .81}, portraitPoint{.51, .93}, portraitPoint{.475, .8}, portraitPoint{.484, .68}, portraitPoint{.467, .52})
	p.stroke(215, 1.2, portraitPoint{.49, .44}, portraitPoint{.507, .55}, portraitPoint{.476, .69}, portraitPoint{.519, .805}, portraitPoint{.506, .895})
	// Stumps of the old bridge face one another above the fissure.
	poly(94, portraitPoint{.09, .36}, portraitPoint{.38, .36}, portraitPoint{.41, .4}, portraitPoint{.39, .45}, portraitPoint{.09, .43})
	poly(180, portraitPoint{.1, .36}, portraitPoint{.37, .36}, portraitPoint{.385, .387}, portraitPoint{.1, .38})
	poly(94, portraitPoint{.65, .37}, portraitPoint{.87, .37}, portraitPoint{.91, .45}, portraitPoint{.635, .43})
	poly(137, portraitPoint{.65, .37}, portraitPoint{.87, .37}, portraitPoint{.883, .391}, portraitPoint{.642, .393})
	p.detail(.24, .57, '╱', 240)
	p.detail(.75, .63, '╲', 246)
}

func drawJournalHalls(f *Frame, x, y, w, h int) {
	p := portraitPainter{f, x, y, w, h}
	poly := p.poly
	// Three receding doorways turn the silhouette into a face of stone.
	for i := 0; i < 3; i++ {
		inset := float64(i) * .1
		top := .07 + float64(i)*.12
		l, r := .12+inset, .88-inset
		b := .94 - float64(i)*.065
		poly([]int{94, 239, 137}[i], portraitPoint{l, b}, portraitPoint{l, .29 + float64(i)*.11}, portraitPoint{.34 + inset*.45, top}, portraitPoint{.66 - inset*.45, top}, portraitPoint{r, .29 + float64(i)*.11}, portraitPoint{r, b})
		poly(233, portraitPoint{l + .045, b}, portraitPoint{l + .045, .32 + float64(i)*.11}, portraitPoint{.36 + inset*.45, top + .045}, portraitPoint{.64 - inset*.45, top + .045}, portraitPoint{r - .045, .32 + float64(i)*.11}, portraitPoint{r - .045, b})
		p.stroke([]int{180, 246, 223}[i], 1.0, portraitPoint{l + .012, .8}, portraitPoint{l + .012, .3 + float64(i)*.11}, portraitPoint{.35 + inset*.45, top + .018}, portraitPoint{.65 - inset*.45, top + .018})
	}
	poly(236, portraitPoint{.16, .94}, portraitPoint{.42, .63}, portraitPoint{.58, .63}, portraitPoint{.84, .94})
	for i := 0; i < 5; i++ {
		v := .7 + float64(i)*.051
		half := (v - .63) * 1.1
		p.stroke(239, .9, portraitPoint{.5 - half, v}, portraitPoint{.5 + half, v})
	}
	for _, u := range []float64{.18, .81} {
		p.stroke(94, 2.5, portraitPoint{u, .55}, portraitPoint{u, .67})
		poly(173, portraitPoint{u, .42}, portraitPoint{u + .03, .52}, portraitPoint{u, .58}, portraitPoint{u - .025, .52})
		poly(223, portraitPoint{u, .49}, portraitPoint{u + .011, .53}, portraitPoint{u, .562}, portraitPoint{u - .01, .53})
	}
}

func drawJournalGarden(f *Frame, x, y, w, h int) {
	p := portraitPainter{f, x, y, w, h}
	poly, oval := p.poly, p.smoothOval
	// A pale branching crown rises from roots cupping a blue underground pool.
	poly(94, portraitPoint{.45, .31}, portraitPoint{.52, .29}, portraitPoint{.56, .59}, portraitPoint{.69, .78}, portraitPoint{.81, .86}, portraitPoint{.58, .8}, portraitPoint{.51, .67}, portraitPoint{.42, .8}, portraitPoint{.2, .86}, portraitPoint{.36, .76}, portraitPoint{.44, .57})
	poly(137, portraitPoint{.46, .32}, portraitPoint{.484, .33}, portraitPoint{.495, .59}, portraitPoint{.427, .742}, portraitPoint{.327, .79}, portraitPoint{.461, .597})
	for _, branch := range [][]portraitPoint{
		{{.478, .49}, {.34, .29}, {.2, .22}}, {{.496, .39}, {.57, .22}, {.69, .19}}, {{.486, .35}, {.44, .18}, {.5, .08}}, {{.516, .52}, {.65, .4}, {.78, .38}},
	} {
		p.stroke(137, max(1.8, float64(w)*.025), branch...)
	}
	for _, leaf := range [][4]float64{{.22, .18, .15, .075}, {.35, .28, .15, .085}, {.43, .12, .12, .06}, {.57, .16, .15, .077}, {.71, .22, .15, .08}, {.7, .36, .14, .07}, {.2, .35, .1, .06}} {
		oval(65, leaf[0], leaf[1], leaf[2], leaf[3])
		oval(151, leaf[0]-.008, leaf[1]-.025, leaf[2]*.88, leaf[3]*.57)
	}
	oval(239, .5, .853, .36, .107)
	oval(109, .5, .832, .31, .073)
	oval(23, .5, .833, .27, .048)
	oval(117, .47, .82, .2, .023)
	poly(239, portraitPoint{.155, .82}, portraitPoint{.27, .87}, portraitPoint{.72, .87}, portraitPoint{.848, .82}, portraitPoint{.76, .943}, portraitPoint{.25, .943})
	poly(246, portraitPoint{.23, .868}, portraitPoint{.5, .9}, portraitPoint{.5, .934}, portraitPoint{.264, .934})
	for _, u := range []float64{.19, .29, .73, .81} {
		p.stroke(65, .9, portraitPoint{u, .8}, portraitPoint{u, .765})
		p.smoothOval(223, u, .755, .014, .018)
	}
}

func drawJournalHeart(f *Frame, x, y, w, h int) {
	p := portraitPainter{f, x, y, w, h}
	poly, oval := p.poly, p.smoothOval
	// Paired ribs enclose a single warm ember: the lair's reason made figurative.
	oval(52, .5, .52, .22, .31)
	oval(95, .5, .53, .175, .25)
	poly(173, portraitPoint{.49, .24}, portraitPoint{.56, .4}, portraitPoint{.6, .5}, portraitPoint{.59, .64}, portraitPoint{.52, .79}, portraitPoint{.43, .74}, portraitPoint{.395, .6}, portraitPoint{.425, .49}, portraitPoint{.445, .35}, portraitPoint{.46, .48})
	poly(215, portraitPoint{.5, .38}, portraitPoint{.54, .53}, portraitPoint{.55, .65}, portraitPoint{.505, .745}, portraitPoint{.458, .689}, portraitPoint{.446, .581}, portraitPoint{.477, .49})
	poly(223, portraitPoint{.499, .5}, portraitPoint{.518, .602}, portraitPoint{.505, .708}, portraitPoint{.482, .65}, portraitPoint{.48, .583})
	for i := 0; i < 3; i++ {
		t := float64(i)
		start := .17 + t*.155
		edge := .22 - t*.035
		end := .69 + t*.075
		for _, sign := range []float64{-1, 1} {
			pts := portraitCurve(portraitPoint{.5 + sign*.11, start}, portraitPoint{.5 + sign*(.5-edge), start + .025}, portraitPoint{.5 + sign*(.5-edge), end - .08}, portraitPoint{.5 + sign*.07, end})
			p.stroke(239, max(2.5, float64(w)*.045), pts...)
			p.stroke(180, .8, pts...)
		}
	}
	poly(94, portraitPoint{.26, .854}, portraitPoint{.73, .854}, portraitPoint{.8, .932}, portraitPoint{.2, .932})
	poly(137, portraitPoint{.32, .872}, portraitPoint{.67, .872}, portraitPoint{.73, .915}, portraitPoint{.26, .915})
}

func drawJournalDepths(f *Frame, x, y, w, h int) {
	p := portraitPainter{f, x, y, w, h}
	poly := p.poly
	// An impossible descending spiral surrounded by floating broken masonry.
	poly(236, portraitPoint{.14, .2}, portraitPoint{.34, .07}, portraitPoint{.69, .08}, portraitPoint{.86, .28}, portraitPoint{.81, .7}, portraitPoint{.68, .92}, portraitPoint{.28, .91}, portraitPoint{.1, .68})
	poly(234, portraitPoint{.22, .27}, portraitPoint{.39, .15}, portraitPoint{.66, .18}, portraitPoint{.78, .33}, portraitPoint{.71, .68}, portraitPoint{.63, .82}, portraitPoint{.31, .82}, portraitPoint{.2, .62})
	path := []portraitPoint{{.29, .74}, {.22, .37}, {.39, .22}, {.65, .24}, {.73, .41}, {.65, .7}, {.39, .69}, {.33, .43}, {.44, .35}, {.6, .37}, {.62, .52}, {.55, .6}, {.44, .57}, {.46, .47}, {.54, .48}}
	p.stroke(24, max(2.5, float64(w)*.04), path...)
	p.stroke(117, 1.2, path...)
	for _, q := range [][2]float64{{.09, .34}, {.18, .08}, {.88, .6}, {.79, .9}, {.15, .86}} {
		poly(239, portraitPoint{q[0] - .03, q[1]}, portraitPoint{q[0], q[1] - .033}, portraitPoint{q[0] + .05, q[1] - .01}, portraitPoint{q[0] + .04, q[1] + .04}, portraitPoint{q[0] - .02, q[1] + .04})
	}
}

func drawJournalMoment(f *Frame, x, y, w, h, floor, diff int) {
	// A location above, a witness below. Objects record the particular memory.
	PlaceJournalEntries[floor].paint(f, x+w/4, y, w*3/4, h*3/4)
	drawGrakPortrait(f, x, y+h/3, w/2, h*2/3)
	p := portraitPainter{f, x, y, w, h}
	poly := p.poly
	switch floor {
	case 0: // Water beside the hoard: the bowl, cases, and an unbroken pot.
		journalBowl(p, .69, .79, .2)
		if diff == 1 {
			p.stroke(137, 2, portraitPoint{.57, .9}, portraitPoint{.62, .9})
			p.stroke(137, 2, portraitPoint{.74, .91}, portraitPoint{.79, .91})
		}
		if diff == 2 {
			poly(65, portraitPoint{.84, .75}, portraitPoint{.92, .75}, portraitPoint{.91, .88}, portraitPoint{.855, .88})
			p.stroke(151, 1, portraitPoint{.88, .78}, portraitPoint{.86, .7}, portraitPoint{.9, .66})
		}
	case 1:
		p.stroke(94, 2, portraitPoint{.82, .56}, portraitPoint{.82, .9})
		poly(95, portraitPoint{.82, .57}, portraitPoint{.94, .59}, portraitPoint{.895, .64}, portraitPoint{.94, .69}, portraitPoint{.82, .69})
		if diff == 1 {
			for i := 0; i < 3; i++ {
				u := .57 + float64(i)*.055
				p.stroke(180, 1.1, portraitPoint{u, .83}, portraitPoint{u + .025, .9})
			}
		}
		if diff == 2 {
			poly(239, portraitPoint{.61, .84}, portraitPoint{.79, .84}, portraitPoint{.76, .95}, portraitPoint{.57, .95})
		}
	case 2:
		// A borrowed helmet, empty boots, then an open homeward doorway.
		if diff == 1 {
			poly(109, portraitPoint{.64, .79}, portraitPoint{.665, .73}, portraitPoint{.74, .72}, portraitPoint{.79, .8}, portraitPoint{.82, .825}, portraitPoint{.62, .825})
		} else {
			p.stroke(180, 1, portraitPoint{.71, .75}, portraitPoint{.71, .9})
			if diff == 2 {
				p.stroke(223, 1, portraitPoint{.76, .78}, portraitPoint{.85, .78})
				p.stroke(223, 1, portraitPoint{.76, .82}, portraitPoint{.9, .82})
			}
		}
	case 3:
		journalBowl(p, .71, .85, .18)
		for i := 0; i < diff+1; i++ {
			u := .62 + float64(i)*.09
			p.stroke(65, 1.3, portraitPoint{u, .8}, portraitPoint{u, .7})
			p.smoothOval(223, u, .68, .025, .026)
		}
	case 4:
		p.smoothOval(52, .73, .82, .15, .12)
		p.smoothOval(173, .73, .82, .11, .085)
		if diff == 0 {
			p.smoothOval(215, .73, .8, .04, .045)
		}
		if diff == 1 {
			p.stroke(223, 1.2, portraitPoint{.65, .8}, portraitPoint{.75, .8})
			p.detail(.69, .8, '━', 234)
		}
		if diff == 2 {
			p.stroke(180, 1.1, portraitPoint{.65, .83}, portraitPoint{.7, .83}, portraitPoint{.71, .79}, portraitPoint{.73, .87}, portraitPoint{.75, .83}, portraitPoint{.83, .83})
		}
	case 5:
		poly(137, portraitPoint{.57, .71}, portraitPoint{.87, .74}, portraitPoint{.9, .93}, portraitPoint{.6, .9})
		poly(223, portraitPoint{.59, .73}, portraitPoint{.85, .76}, portraitPoint{.88, .91}, portraitPoint{.62, .88})
		p.stroke(94, 1, portraitPoint{.64, .78}, portraitPoint{.74, .78}, portraitPoint{.72, .84}, portraitPoint{.83, .87})
		if diff > 0 {
			p.stroke(65, 1, portraitPoint{.65, .87}, portraitPoint{.78, .81}, portraitPoint{.84, .84})
		}
		if diff == 2 {
			p.smoothOval(239, .92, .91, .035, .025)
		}
	}
}

func journalBowl(p portraitPainter, u, v, r float64) {
	p.smoothOval(239, u, v, r, r*.25)
	p.smoothOval(117, u, v-.006, r*.8, r*.14)
	p.poly(109, portraitPoint{u - r, v}, portraitPoint{u + r, v}, portraitPoint{u + r*.7, v + r*.33}, portraitPoint{u - r*.7, v + r*.33})
	p.poly(246, portraitPoint{u - r, v}, portraitPoint{u - r*.7, v + r*.33}, portraitPoint{u, v + r*.33}, portraitPoint{u, v + r*.17})
}

func drawJournalAfterward(f *Frame, x, y, w, h, moment int) {
	drawJournalRotunda(f, x+w/4, y, w*3/4, h*3/4)
	drawGrakPortrait(f, x, y+h/3, w/2, h*2/3)
	p := portraitPainter{f, x, y, w, h}
	poly := p.poly
	switch moment {
	case 0:
		journalBowl(p, .7, .81, .2)
		for _, u := range []float64{.6, .71, .82} {
			p.smoothOval(173, u, .775, .033, .037)
			p.stroke(65, 1, portraitPoint{u, .744}, portraitPoint{u + .018, .73})
		}
	case 1:
		p.stroke(94, 2, portraitPoint{.62, .74}, portraitPoint{.84, .74}, portraitPoint{.84, .91})
		p.stroke(94, 2, portraitPoint{.65, .9}, portraitPoint{.65, .61}, portraitPoint{.81, .61}, portraitPoint{.81, .74})
		p.stroke(180, .9, portraitPoint{.65, .63}, portraitPoint{.81, .63})
	case 2:
		poly(95, portraitPoint{.59, .7}, portraitPoint{.86, .7}, portraitPoint{.86, .93}, portraitPoint{.59, .93})
		poly(131, portraitPoint{.615, .73}, portraitPoint{.75, .73}, portraitPoint{.75, .9}, portraitPoint{.615, .9})
		p.stroke(223, .9, portraitPoint{.7, .72}, portraitPoint{.7, .78}, portraitPoint{.74, .81}, portraitPoint{.71, .85}, portraitPoint{.71, .91})
	case 3:
		journalBowl(p, .72, .85, .18)
		p.smoothOval(239, .58, .93, .035, .019)
		p.smoothOval(239, .87, .92, .04, .023)
		for i := 0; i < 5; i++ {
			u := .48 + float64(i)*.085
			v := .09 + math.Sin(float64(i))*.025
			p.detail(u, v, '·', 180)
		}
	}
}
