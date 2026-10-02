package render

import "math"

// RenderCannonierMockup depicts a practiced monster crew loading their cannon.
// A foreground bore, round wheels and connected work poses carry the story.
func RenderCannonierMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the crew", 180, 233, false)
	} else {
		ph := min(h-7, (w-6)/3)
		pw := ph * 3
		drawCannonierPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-9)/2, 1, "CANNONIER", 180, 233, true)
		putString(f, (w-22)/2, 2, "one more for the guild", 240, 233, false)
	}
	hint := "tab Ranger · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawCannonierPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.smoothOval
	// A long, low shelf supports the whole team rather than isolated plinths.
	poly(236, portraitPoint{.025, .94}, portraitPoint{.94, .94}, portraitPoint{.97, .975}, portraitPoint{.025, .975})
	// The loader's stance leans toward the cannon; one foot resists the shove.
	poly(60, portraitPoint{.13, .66}, portraitPoint{.22, .68}, portraitPoint{.17, .80}, portraitPoint{.14, .90}, portraitPoint{.07, .91}, portraitPoint{.08, .78})
	poly(59, portraitPoint{.13, .76}, portraitPoint{.17, .80}, portraitPoint{.14, .90}, portraitPoint{.10, .90}, portraitPoint{.115, .81})
	poly(60, portraitPoint{.21, .66}, portraitPoint{.29, .66}, portraitPoint{.30, .79}, portraitPoint{.26, .87}, portraitPoint{.30, .92}, portraitPoint{.22, .92}, portraitPoint{.20, .85}, portraitPoint{.23, .77})
	poly(67, portraitPoint{.255, .72}, portraitPoint{.285, .72}, portraitPoint{.295, .79}, portraitPoint{.26, .855}, portraitPoint{.24, .835}, portraitPoint{.26, .78})
	poly(237, portraitPoint{.075, .865}, portraitPoint{.145, .865}, portraitPoint{.16, .90}, portraitPoint{.13, .94}, portraitPoint{.03, .94}, portraitPoint{.03, .915})
	poly(94, portraitPoint{.065, .895}, portraitPoint{.14, .895}, portraitPoint{.13, .925}, portraitPoint{.04, .925}, portraitPoint{.04, .915})
	poly(237, portraitPoint{.225, .875}, portraitPoint{.265, .875}, portraitPoint{.31, .92}, portraitPoint{.31, .94}, portraitPoint{.215, .94}, portraitPoint{.20, .915})
	poly(94, portraitPoint{.225, .905}, portraitPoint{.265, .905}, portraitPoint{.30, .925}, portraitPoint{.23, .925})
	// A broad back, rolled sleeve, and apron give the loader a working build.
	poly(65, portraitPoint{.145, .39}, portraitPoint{.24, .395}, portraitPoint{.29, .47}, portraitPoint{.29, .61}, portraitPoint{.255, .71}, portraitPoint{.13, .70}, portraitPoint{.075, .58}, portraitPoint{.09, .47})
	poly(107, portraitPoint{.15, .40}, portraitPoint{.215, .405}, portraitPoint{.265, .46}, portraitPoint{.26, .565}, portraitPoint{.18, .595}, portraitPoint{.10, .54}, portraitPoint{.11, .465})
	poly(66, portraitPoint{.105, .46}, portraitPoint{.155, .425}, portraitPoint{.21, .45}, portraitPoint{.195, .52}, portraitPoint{.13, .55}, portraitPoint{.09, .505})
	poly(109, portraitPoint{.10, .505}, portraitPoint{.135, .49}, portraitPoint{.195, .475}, portraitPoint{.195, .52}, portraitPoint{.13, .55})
	poly(94, portraitPoint{.14, .545}, portraitPoint{.245, .52}, portraitPoint{.28, .62}, portraitPoint{.255, .735}, portraitPoint{.20, .72}, portraitPoint{.165, .74}, portraitPoint{.12, .68})
	poly(137, portraitPoint{.16, .56}, portraitPoint{.21, .548}, portraitPoint{.225, .67}, portraitPoint{.205, .715}, portraitPoint{.17, .70})
	poly(130, portraitPoint{.23, .405}, portraitPoint{.25, .415}, portraitPoint{.185, .69}, portraitPoint{.16, .675})
	// Loader's neck and exaggerated jaw sit just above the bronze bore.
	poly(65, portraitPoint{.17, .34}, portraitPoint{.24, .34}, portraitPoint{.245, .43}, portraitPoint{.18, .44})
	poly(107, portraitPoint{.075, .275}, portraitPoint{.175, .29}, portraitPoint{.19, .355}, portraitPoint{.135, .365}, portraitPoint{.10, .325})
	poly(65, portraitPoint{.115, .305}, portraitPoint{.16, .315}, portraitPoint{.17, .345}, portraitPoint{.14, .345})
	poly(107, portraitPoint{.145, .235}, portraitPoint{.215, .23}, portraitPoint{.255, .265}, portraitPoint{.27, .305}, portraitPoint{.315, .335}, portraitPoint{.31, .365}, portraitPoint{.28, .38}, portraitPoint{.265, .42}, portraitPoint{.195, .415}, portraitPoint{.145, .365}, portraitPoint{.12, .30})
	poly(150, portraitPoint{.15, .25}, portraitPoint{.215, .245}, portraitPoint{.248, .275}, portraitPoint{.25, .30}, portraitPoint{.17, .31}, portraitPoint{.135, .295})
	poly(65, portraitPoint{.145, .33}, portraitPoint{.195, .355}, portraitPoint{.265, .355}, portraitPoint{.28, .38}, portraitPoint{.265, .42}, portraitPoint{.195, .415}, portraitPoint{.145, .365})
	poly(107, portraitPoint{.185, .365}, portraitPoint{.275, .36}, portraitPoint{.265, .39}, portraitPoint{.22, .40}, portraitPoint{.19, .39})
	poly(150, portraitPoint{.265, .305}, portraitPoint{.307, .332}, portraitPoint{.30, .347}, portraitPoint{.27, .34})
	poly(22, portraitPoint{.215, .30}, portraitPoint{.26, .30}, portraitPoint{.27, .325}, portraitPoint{.225, .33})
	poly(22, portraitPoint{.215, .365}, portraitPoint{.275, .36}, portraitPoint{.275, .375}, portraitPoint{.22, .385})
	poly(180, portraitPoint{.245, .382}, portraitPoint{.265, .382}, portraitPoint{.267, .338}, portraitPoint{.25, .357})
	// A rusty-red tied cloth and sweat bead express effort, not a uniform.
	poly(95, portraitPoint{.13, .27}, portraitPoint{.135, .23}, portraitPoint{.20, .215}, portraitPoint{.24, .245}, portraitPoint{.25, .275}, portraitPoint{.18, .27})
	poly(131, portraitPoint{.13, .25}, portraitPoint{.20, .235}, portraitPoint{.24, .25}, portraitPoint{.25, .275}, portraitPoint{.15, .283})
	poly(95, portraitPoint{.14, .255}, portraitPoint{.11, .28}, portraitPoint{.08, .31}, portraitPoint{.065, .355}, portraitPoint{.10, .345}, portraitPoint{.13, .31})
	// The smaller spotter is visible above the breech, checking the line.
	poly(101, portraitPoint{.65, .30}, portraitPoint{.75, .29}, portraitPoint{.81, .35}, portraitPoint{.80, .47}, portraitPoint{.69, .515}, portraitPoint{.615, .45}, portraitPoint{.61, .37})
	poly(66, portraitPoint{.64, .34}, portraitPoint{.74, .32}, portraitPoint{.80, .375}, portraitPoint{.77, .45}, portraitPoint{.66, .46}, portraitPoint{.61, .40})
	poly(60, portraitPoint{.74, .35}, portraitPoint{.80, .375}, portraitPoint{.77, .45}, portraitPoint{.72, .44})
	poly(143, portraitPoint{.665, .275}, portraitPoint{.735, .275}, portraitPoint{.75, .35}, portraitPoint{.69, .365}, portraitPoint{.645, .33})
	poly(143, portraitPoint{.74, .22}, portraitPoint{.84, .16}, portraitPoint{.825, .245}, portraitPoint{.76, .275})
	poly(101, portraitPoint{.77, .224}, portraitPoint{.815, .195}, portraitPoint{.80, .238}, portraitPoint{.765, .25})
	poly(143, portraitPoint{.67, .175}, portraitPoint{.735, .175}, portraitPoint{.77, .21}, portraitPoint{.76, .29}, portraitPoint{.71, .33}, portraitPoint{.655, .32}, portraitPoint{.63, .285}, portraitPoint{.59, .275}, portraitPoint{.595, .25}, portraitPoint{.64, .225})
	poly(186, portraitPoint{.675, .19}, portraitPoint{.72, .19}, portraitPoint{.745, .215}, portraitPoint{.735, .25}, portraitPoint{.66, .25}, portraitPoint{.64, .235})
	poly(101, portraitPoint{.735, .245}, portraitPoint{.765, .23}, portraitPoint{.76, .29}, portraitPoint{.71, .33}, portraitPoint{.66, .32}, portraitPoint{.69, .29})
	poly(186, portraitPoint{.645, .24}, portraitPoint{.66, .25}, portraitPoint{.64, .275}, portraitPoint{.598, .27}, portraitPoint{.61, .255})
	poly(94, portraitPoint{.65, .19}, portraitPoint{.67, .15}, portraitPoint{.73, .14}, portraitPoint{.775, .185}, portraitPoint{.77, .22}, portraitPoint{.715, .20}, portraitPoint{.66, .21})
	poly(130, portraitPoint{.67, .165}, portraitPoint{.73, .155}, portraitPoint{.76, .18}, portraitPoint{.755, .195}, portraitPoint{.695, .18})
	poly(180, portraitPoint{.64, .21}, portraitPoint{.755, .215}, portraitPoint{.76, .238}, portraitPoint{.64, .235})
	poly(237, portraitPoint{.64, .24}, portraitPoint{.68, .24}, portraitPoint{.68, .265}, portraitPoint{.637, .265})
	// A connecting timber bed and trail distribute the machine's weight.
	poly(94, portraitPoint{.285, .61}, portraitPoint{.75, .54}, portraitPoint{.84, .66}, portraitPoint{.79, .79}, portraitPoint{.42, .875}, portraitPoint{.275, .80})
	poly(137, portraitPoint{.30, .625}, portraitPoint{.74, .56}, portraitPoint{.81, .65}, portraitPoint{.77, .675}, portraitPoint{.31, .71})
	poly(130, portraitPoint{.375, .69}, portraitPoint{.405, .67}, portraitPoint{.76, .865}, portraitPoint{.745, .91}, portraitPoint{.685, .91})
	poly(94, portraitPoint{.45, .73}, portraitPoint{.49, .70}, portraitPoint{.83, .85}, portraitPoint{.84, .885}, portraitPoint{.765, .89})
	poly(239, portraitPoint{.75, .875}, portraitPoint{.835, .85}, portraitPoint{.86, .87}, portraitPoint{.84, .915}, portraitPoint{.76, .915})
	// The far wheel is shaded and partly hidden by the cannon and loader.
	cannonWheel(p, .375, .795, .09, .123, false)
	// A cast iron barrel swells toward the viewer, with three broad planes.
	oval(237, .76, .43, .074, .098)
	poly(237, portraitPoint{.355, .385}, portraitPoint{.75, .335}, portraitPoint{.82, .365}, portraitPoint{.825, .465}, portraitPoint{.77, .525}, portraitPoint{.36, .69})
	poly(242, portraitPoint{.37, .402}, portraitPoint{.75, .35}, portraitPoint{.80, .375}, portraitPoint{.805, .435}, portraitPoint{.76, .48}, portraitPoint{.385, .615})
	poly(247, portraitPoint{.385, .405}, portraitPoint{.75, .35}, portraitPoint{.80, .375}, portraitPoint{.80, .395}, portraitPoint{.40, .49})
	poly(239, portraitPoint{.40, .49}, portraitPoint{.80, .395}, portraitPoint{.80, .45}, portraitPoint{.765, .49}, portraitPoint{.40, .635})
	poly(236, portraitPoint{.39, .595}, portraitPoint{.765, .49}, portraitPoint{.77, .525}, portraitPoint{.36, .69})
	// Reinforcing hoops and a round trunnion connect metal to the carriage.
	poly(180, portraitPoint{.535, .38}, portraitPoint{.565, .375}, portraitPoint{.60, .58}, portraitPoint{.565, .595})
	poly(137, portraitPoint{.58, .50}, portraitPoint{.60, .50}, portraitPoint{.60, .58}, portraitPoint{.565, .595}, portraitPoint{.565, .565})
	poly(180, portraitPoint{.715, .348}, portraitPoint{.738, .343}, portraitPoint{.77, .52}, portraitPoint{.74, .535})
	poly(137, portraitPoint{.74, .465}, portraitPoint{.765, .458}, portraitPoint{.77, .52}, portraitPoint{.74, .535})
	oval(237, .625, .575, .048, .065)
	oval(137, .618, .566, .031, .042)
	oval(180, .612, .554, .014, .019)
	poly(94, portraitPoint{.59, .58}, portraitPoint{.65, .58}, portraitPoint{.65, .655}, portraitPoint{.59, .65})
	// Breech sight: one small raised block below the spotter's searching eye.
	poly(239, portraitPoint{.74, .352}, portraitPoint{.74, .30}, portraitPoint{.775, .30}, portraitPoint{.79, .347})
	poly(180, portraitPoint{.745, .30}, portraitPoint{.775, .30}, portraitPoint{.78, .32}, portraitPoint{.745, .32})
	// Foreground muzzle: warm lip, recessed bevel and a deep circular bore.
	oval(237, .38, .54, .128, .175)
	oval(137, .377, .535, .121, .165)
	oval(180, .369, .524, .111, .151)
	oval(239, .387, .548, .091, .127)
	oval(232, .388, .550, .082, .114)
	oval(236, .399, .564, .064, .092)
	oval(232, .406, .565, .056, .081)
	// The front wheel masks the carriage but leaves the whole muzzle visible.
	cannonWheel(p, .635, .77, .128, .175, true)
	// Ramrod enters the bore. Its tip is a cloth wad, not a projectile.
	poly(94, portraitPoint{.04, .61}, portraitPoint{.403, .523}, portraitPoint{.41, .555}, portraitPoint{.04, .643})
	poly(137, portraitPoint{.04, .61}, portraitPoint{.40, .524}, portraitPoint{.405, .54}, portraitPoint{.04, .627})
	oval(244, .409, .544, .023, .032)
	oval(250, .405, .534, .018, .021)
	// Two bent arms end in grips on the same shaft, spelling out the push.
	poly(65, portraitPoint{.15, .49}, portraitPoint{.18, .49}, portraitPoint{.185, .56}, portraitPoint{.24, .53}, portraitPoint{.26, .55}, portraitPoint{.23, .61}, portraitPoint{.155, .63}, portraitPoint{.115, .57}, portraitPoint{.12, .535})
	poly(107, portraitPoint{.13, .53}, portraitPoint{.165, .515}, portraitPoint{.175, .575}, portraitPoint{.225, .55}, portraitPoint{.24, .565}, portraitPoint{.215, .60}, portraitPoint{.16, .603}, portraitPoint{.13, .565})
	poly(150, portraitPoint{.17, .575}, portraitPoint{.225, .55}, portraitPoint{.24, .565}, portraitPoint{.215, .583}, portraitPoint{.175, .59})
	poly(65, portraitPoint{.24, .43}, portraitPoint{.275, .455}, portraitPoint{.30, .515}, portraitPoint{.29, .555}, portraitPoint{.265, .565}, portraitPoint{.22, .535}, portraitPoint{.225, .515}, portraitPoint{.27, .52})
	poly(107, portraitPoint{.255, .445}, portraitPoint{.27, .457}, portraitPoint{.286, .51}, portraitPoint{.277, .535}, portraitPoint{.25, .53}, portraitPoint{.237, .515}, portraitPoint{.27, .515})
	poly(107, portraitPoint{.21, .557}, portraitPoint{.235, .545}, portraitPoint{.255, .555}, portraitPoint{.25, .59}, portraitPoint{.23, .607}, portraitPoint{.205, .59})
	poly(150, portraitPoint{.215, .56}, portraitPoint{.235, .55}, portraitPoint{.244, .563}, portraitPoint{.235, .585}, portraitPoint{.213, .58})
	poly(107, portraitPoint{.255, .518}, portraitPoint{.283, .505}, portraitPoint{.30, .53}, portraitPoint{.295, .56}, portraitPoint{.27, .57}, portraitPoint{.25, .55})
	poly(150, portraitPoint{.264, .522}, portraitPoint{.282, .515}, portraitPoint{.287, .53}, portraitPoint{.277, .55}, portraitPoint{.26, .54})
	// The spotter's hands curl over the far barrel, separated from the loading.
	poly(143, portraitPoint{.605, .36}, portraitPoint{.64, .35}, portraitPoint{.665, .373}, portraitPoint{.66, .40}, portraitPoint{.62, .415}, portraitPoint{.595, .393})
	poly(186, portraitPoint{.61, .36}, portraitPoint{.64, .36}, portraitPoint{.65, .375}, portraitPoint{.62, .39}, portraitPoint{.60, .38})
	poly(143, portraitPoint{.78, .37}, portraitPoint{.80, .38}, portraitPoint{.82, .403}, portraitPoint{.82, .435}, portraitPoint{.79, .445}, portraitPoint{.765, .42}, portraitPoint{.765, .392})
	poly(186, portraitPoint{.778, .375}, portraitPoint{.80, .39}, portraitPoint{.804, .415}, portraitPoint{.78, .41})
	// Stacked cannonballs echo the muzzle and wheel circles at a small scale.
	cannonBall(p, .835, .905, .036, .049)
	cannonBall(p, .906, .925, .036, .049)
	cannonBall(p, .873, .847, .036, .049)
	p.detail(.235, .315, '━', 230)
	p.detail(.655, .252, '━', 230)
	p.detail(.67, .30, '─', 94)
	p.detail(.695, .295, '╯', 94)
	p.detail(.165, .32, '╲', 150)
	p.detail(.225, .572, '│', 65)
	p.detail(.275, .54, '│', 65)
	p.detail(.625, .38, '│', 101)
	p.detail(.79, .407, '│', 101)
	if h >= 22 {
		p.detail(.20, .335, '·', 255)
		p.detail(.73, .275, '╱', 186)
		p.detail(.46, .47, '╲', 247)
	}
	if h < 22 {
		p.detail(.26, .375, '▴', 180)
	}
}

func cannonWheel(p portraitPainter, cx, cy, rx, ry float64, near bool) {
	outer, rim, spoke, shade := 236, 95, 137, 236
	if near {
		outer, rim, spoke, shade = 241, 137, 130, 236
	}
	p.smoothOval(outer, cx, cy, rx, ry)
	p.smoothOval(237, cx+.004, cy+.006, rx*.96, ry*.96)
	p.smoothOval(rim, cx, cy, rx*.87, ry*.87)
	p.smoothOval(shade, cx, cy, rx*.68, ry*.68)
	for i := 0; i < 6; i++ {
		a := float64(i)*math.Pi/3 + .15
		point := func(r, t float64) portraitPoint { return portraitPoint{cx + rx*r*math.Cos(t), cy + ry*r*math.Sin(t)} }
		p.poly(spoke, point(.12, a-.3), point(.78, a-.13), point(.78, a+.13), point(.12, a+.3))
	}
	p.smoothOval(239, cx, cy, rx*.27, ry*.27)
	p.smoothOval(180, cx-rx*.025, cy-ry*.025, rx*.12, ry*.12)
}

func cannonBall(p portraitPainter, cx, cy, rx, ry float64) {
	p.smoothOval(237, cx, cy, rx, ry)
	p.smoothOval(241, cx-rx*.13, cy-ry*.20, rx*.8, ry*.8)
	p.smoothOval(244, cx-rx*.35, cy-ry*.38, rx*.27, ry*.27)
}
