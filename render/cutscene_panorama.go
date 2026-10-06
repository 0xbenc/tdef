package render

import "math"

const (
	filmPanStart = 48                 // Establish the first composition for 1.6 seconds.
	filmPanEnd   = filmPanStart + 150 // Five seconds of travel: 20% slower.
	filmPanHold  = 60                 // Let the final reveal land for two seconds.
)

// CutsceneCanAdvance gives each panorama time to establish, travel, and land.
// Revealing text and skipping the entire film remain available throughout.
func CutsceneCanAdvance(st CutsceneState) bool {
	shots := FilmShots(st.Film)
	shot := shots[max(0, min(st.Shot, len(shots)-1))]
	return !filmIsPanorama(shot.art) || st.Frame >= filmPanEnd+filmPanHold
}

func filmIsPanorama(art string) bool {
	switch art {
	case "village", "village-search", "ruined-road", "capital-vault", "ending-pursuit", "ending-depths", "ending-earth":
		return true
	}
	return false
}

// Move the viewport by whole terminal columns. A shape never gets resampled
// during travel: its exact cells slide together, preserving edges and texture.
func filmPanColumn(art string, frame, travel int) int {
	t := float64(max(0, min(frame-filmPanStart, filmPanEnd-filmPanStart))) / float64(filmPanEnd-filmPanStart)
	t = t * t * (3 - 2*t)
	x := int(math.Round(t * float64(travel)))
	if art == "village-search" {
		return travel - x
	}
	return x
}

func paintFilmPanorama(dst *Frame, stage Rect, art string, frame int) {
	if art == "ending-pursuit" || art == "ending-depths" || art == "ending-earth" {
		paintEndingPanorama(dst, stage, art, frame)
		return
	}
	if filmIsOutdoor(art) {
		far := filmFrame(stage.W*2, stage.H)
		filmOutdoorBackdrop(far, art)
		near := filmTransparentFrame(stage.W*2, stage.H)
		if art == "ruined-road" {
			filmCapitalForeground(near, frame)
		} else {
			filmVillageForeground(near, art == "village-search", frame)
		}
		x0 := filmPanColumn(art, frame, stage.W)
		filmCompositePan(dst, stage, far, near, x0/4, x0)
		return
	}
	wide := filmFrame(stage.W*2, stage.H)
	filmVaultPanorama(wide)
	x0 := filmPanColumn(art, frame, stage.W)
	for y := 0; y < stage.H; y++ {
		for x := 0; x < stage.W; x++ {
			dst.Set(stage.X+x, stage.Y+y, wide.C[y*wide.W+x0+x])
		}
	}
}

func filmVillageForeground(f *Frame, searching bool, frame int) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	poly := p.poly
	poly(236, portraitPoint{0, .68}, portraitPoint{1, .61}, portraitPoint{1, 1}, portraitPoint{0, 1})
	poly(239, portraitPoint{0, .9}, portraitPoint{.43, .75}, portraitPoint{1, .82}, portraitPoint{1, .98}, portraitPoint{.46, .89}, portraitPoint{0, 1})
	// Separate landmarks across two screens: Grak's wrecked home, a burned
	// row of houses, and the village well with nobody waiting beside it.
	for i, u := range []float64{.02, .2, .44, .63, .86} {
		v := .68 + float64(i%2)*.07
		if i%2 == 0 {
			filmSmoke(p, u+.065, v-.29, frame, i)
		}
		poly(94, portraitPoint{u, v}, portraitPoint{u, v - .3}, portraitPoint{u + .025, v - .34}, portraitPoint{u + .04, v - .2}, portraitPoint{u + .115, v - .29}, portraitPoint{u + .135, v})
		poly(238, portraitPoint{u + .04, v}, portraitPoint{u + .04, v - .13}, portraitPoint{u + .077, v - .13}, portraitPoint{u + .077, v})
		p.stroke(137, 1.3, portraitPoint{u - .01, v + .02}, portraitPoint{u + .13, v - .11})
		p.stroke(58, 1.5, portraitPoint{u + .014, v - .27}, portraitPoint{u + .13, v + .035})
	}
	// A snapped signpost and discarded axe punctuate the search.
	p.stroke(137, 1.2, portraitPoint{.45, .91}, portraitPoint{.46, .61})
	poly(94, portraitPoint{.425, .62}, portraitPoint{.49, .66}, portraitPoint{.48, .73}, portraitPoint{.42, .69})
	p.stroke(94, 1.2, portraitPoint{.3, .94}, portraitPoint{.37, .89})
	poly(244, portraitPoint{.355, .875}, portraitPoint{.385, .875}, portraitPoint{.39, .93}, portraitPoint{.365, .915})
	p.smoothOval(240, .79, .82, .065, .09)
	poly(240, portraitPoint{.725, .82}, portraitPoint{.855, .82}, portraitPoint{.85, .94}, portraitPoint{.73, .94})
	p.smoothOval(246, .79, .81, .065, .05)
	p.smoothOval(233, .79, .81, .045, .028)
	p.stroke(94, 1.5, portraitPoint{.733, .81}, portraitPoint{.733, .52}, portraitPoint{.85, .52}, portraitPoint{.85, .81})
	p.stroke(137, 1, portraitPoint{.79, .52}, portraitPoint{.79, .75})
	if searching {
		// The return pan follows the search past fallen kin and back to Grak.
		for _, u := range []float64{.56, .9} {
			poly(58, portraitPoint{u - .025, .95}, portraitPoint{u + .025, .91}, portraitPoint{u + .06, .94}, portraitPoint{u + .06, .97}, portraitPoint{u - .02, .98})
			p.smoothOval(65, u-.022, .95, .014, .028)
		}
		filmGrakRear(p, .16, .93)
	} else {
		// Render the established portrait at native size in the first screen.
		drawGrakPortrait(f, f.W/20, f.H/5, f.W/5, f.H*4/5)
	}
}

func filmCapitalForeground(f *Frame, frame int) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	poly := p.poly
	// The distant skyline grows into the Rotunda as the camera travels east.
	for i := 0; i < 8; i++ {
		u := .4 + float64(i)*.066
		v := .43 + float64(i%3)*.055
		if i == 0 || i == 3 || i == 5 {
			filmSmoke(p, u+.026, v, frame, i)
		}
		poly(239, portraitPoint{u, .8}, portraitPoint{u, v}, portraitPoint{u + .012, v - .04}, portraitPoint{u + .03, v + .04}, portraitPoint{u + .048, v}, portraitPoint{u + .053, .8})
	}
	poly(236, portraitPoint{0, .76}, portraitPoint{1, .73}, portraitPoint{1, 1}, portraitPoint{0, 1})
	poly(240, portraitPoint{.11, 1}, portraitPoint{.54, .77}, portraitPoint{.83, .73}, portraitPoint{.87, .78}, portraitPoint{.57, .84}, portraitPoint{.32, 1})
	// Broken round roof and massive columns identify the capital's hall.
	poly(94, portraitPoint{.73, .38}, portraitPoint{.76, .23}, portraitPoint{.82, .16}, portraitPoint{.84, .24}, portraitPoint{.87, .18}, portraitPoint{.93, .36})
	poly(137, portraitPoint{.73, .38}, portraitPoint{.93, .36}, portraitPoint{.93, .43}, portraitPoint{.73, .45})
	for _, u := range []float64{.75, .8, .86, .91} {
		poly(244, portraitPoint{u, .43}, portraitPoint{u + .016, .43}, portraitPoint{u + .016, .78}, portraitPoint{u, .78})
	}
	poly(233, portraitPoint{.817, .79}, portraitPoint{.817, .54}, portraitPoint{.844, .47}, portraitPoint{.872, .54}, portraitPoint{.872, .79})
	filmGrakRear(p, .17, .93)
	p.stroke(94, 1.2, portraitPoint{.06, .94}, portraitPoint{.06, .69}, portraitPoint{.105, .71})
}

func filmVaultPanorama(f *Frame) {
	p := portraitPainter{f, 0, 0, f.W, f.H}
	p.poly(236, portraitPoint{0, 0}, portraitPoint{1, 0}, portraitPoint{1, 1}, portraitPoint{0, 1})
	// Begin on the descent, then reveal the entire dragon and his treasure.
	for i := 0; i < 7; i++ {
		u := float64(i) * .052
		v := .18 + float64(i)*.105
		p.poly(240, portraitPoint{u, v}, portraitPoint{u + .09, v}, portraitPoint{u + .09, 1}, portraitPoint{u, 1})
		p.stroke(244, 1, portraitPoint{u, v}, portraitPoint{u + .09, v})
	}
	filmGrakRear(p, .13, .41)
	// Draw once at a fixed scale; the camera only copies these native cells.
	dragon := filmFrame(240, 80)
	paintFilmShot(dragon, CutsceneState{Film: FilmOpening, Frame: 120}, "dragon")
	cinemaProject(f, dragon, Rect{f.W / 2, 0, f.W / 2, f.H}, filmView{0, 0, 1, 1})
	p.poly(238, portraitPoint{.48, 0}, portraitPoint{.51, 0}, portraitPoint{.51, 1}, portraitPoint{.48, 1})
	for i := 0; i < 5; i++ {
		filmCoin(p, .39+float64(i)*.025, .91-float64(i%2)*.05, .013)
	}
}
