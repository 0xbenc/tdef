package render

// Grak's rear view shares the portrait's hunched shoulders, oversized head,
// uneven blade ears, torn burgundy mantle, leather boots and builder's mallet.
// Size in terminal cells rather than panorama fractions so a wide landscape
// cannot stretch him sideways. The feet anchor him to the same ground point.
func filmGrakRear(scene portraitPainter, u, v float64) {
	h := max(5, scene.h*42/100)
	w := max(7, h*3/2)
	x := scene.x0 + int(u*float64(scene.w-1)) - w/2
	y := scene.y0 + int(v*float64(scene.h-1)) - h + 1
	p := portraitPainter{scene.f, x, y, w, h}
	poly := p.poly
	// Separate planted legs survive the small silhouette.
	poly(58, portraitPoint{.35, .65}, portraitPoint{.48, .65}, portraitPoint{.46, .9}, portraitPoint{.31, .9})
	poly(58, portraitPoint{.58, .65}, portraitPoint{.7, .66}, portraitPoint{.75, .9}, portraitPoint{.6, .9})
	poly(94, portraitPoint{.31, .84}, portraitPoint{.46, .84}, portraitPoint{.46, .97}, portraitPoint{.25, .97})
	poly(94, portraitPoint{.6, .84}, portraitPoint{.75, .84}, portraitPoint{.8, .97}, portraitPoint{.6, .97})
	// A broad shoulder line tapers into the ragged mantle, rather than a
	// single green triangle widening all the way to the ground.
	poly(95, portraitPoint{.29, .4}, portraitPoint{.67, .38}, portraitPoint{.77, .46}, portraitPoint{.72, .74}, portraitPoint{.64, .7}, portraitPoint{.59, .77}, portraitPoint{.47, .72}, portraitPoint{.34, .76}, portraitPoint{.29, .64}, portraitPoint{.24, .49})
	poly(52, portraitPoint{.62, .43}, portraitPoint{.74, .47}, portraitPoint{.69, .72}, portraitPoint{.62, .69})
	poly(130, portraitPoint{.35, .43}, portraitPoint{.4, .42}, portraitPoint{.63, .69}, portraitPoint{.59, .72})
	poly(65, portraitPoint{.25, .46}, portraitPoint{.34, .48}, portraitPoint{.3, .64}, portraitPoint{.2, .61})
	poly(107, portraitPoint{.72, .46}, portraitPoint{.8, .49}, portraitPoint{.82, .65}, portraitPoint{.75, .67}, portraitPoint{.7, .6})
	// Tall shaft and chamfered stone head match the established tool.
	poly(94, portraitPoint{.15, .43}, portraitPoint{.19, .43}, portraitPoint{.19, .97}, portraitPoint{.15, .97})
	poly(238, portraitPoint{.04, .35}, portraitPoint{.25, .35}, portraitPoint{.29, .39}, portraitPoint{.27, .46}, portraitPoint{.03, .46}, portraitPoint{.02, .39})
	poly(244, portraitPoint{.04, .35}, portraitPoint{.25, .35}, portraitPoint{.25, .39}, portraitPoint{.03, .39})
	poly(107, portraitPoint{.16, .53}, portraitPoint{.25, .53}, portraitPoint{.27, .61}, portraitPoint{.17, .62})
	// The ears keep their different angles and the right ear's missing corner.
	poly(107, portraitPoint{.1, .14}, portraitPoint{.36, .2}, portraitPoint{.4, .3}, portraitPoint{.24, .26}, portraitPoint{.17, .2})
	poly(107, portraitPoint{.65, .2}, portraitPoint{.9, .1}, portraitPoint{.86, .2}, portraitPoint{.88, .22}, portraitPoint{.73, .3}, portraitPoint{.65, .28})
	poly(107, portraitPoint{.33, .17}, portraitPoint{.63, .15}, portraitPoint{.73, .22}, portraitPoint{.74, .33}, portraitPoint{.65, .41}, portraitPoint{.42, .41}, portraitPoint{.29, .32}, portraitPoint{.29, .23})
	poly(150, portraitPoint{.34, .18}, portraitPoint{.61, .17}, portraitPoint{.67, .22}, portraitPoint{.34, .25})
	poly(65, portraitPoint{.65, .22}, portraitPoint{.73, .23}, portraitPoint{.74, .33}, portraitPoint{.65, .41}, portraitPoint{.61, .34})
}
