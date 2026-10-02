package render

// RenderGnollMockup is one definitive lore portrait for the Gnoll Slingers.
// The throwing pose and stone-bearing partner depict their rapid-fire role.
func RenderGnollMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the slingers", 180, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawGnollPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-13)/2, 1, "GNOLL SLINGERS", 180, 233, true)
		putString(f, (w-23)/2, 2, "one winds, one supplies", 240, 233, false)
	}
	hint := "tab Gunner · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawGnollPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.oval
	// A single stone shelf grounds both figures and their shared ammunition.
	poly(236, portraitPoint{.06, .92}, portraitPoint{.86, .92}, portraitPoint{.91, .95}, portraitPoint{.04, .95})
	// A whip tail, with a dark tuft, counterbalances the throwing arm.
	poly(137, portraitPoint{.58, .61}, portraitPoint{.71, .65}, portraitPoint{.83, .72}, portraitPoint{.86, .80}, portraitPoint{.82, .84}, portraitPoint{.79, .82}, portraitPoint{.82, .78}, portraitPoint{.79, .74}, portraitPoint{.69, .70}, portraitPoint{.57, .67})
	poly(237, portraitPoint{.82, .79}, portraitPoint{.86, .79}, portraitPoint{.86, .85}, portraitPoint{.81, .86}, portraitPoint{.79, .83})
	// Wide, bent digitigrade legs: powerful thighs, narrow hocks, clawed feet.
	poly(137, portraitPoint{.45, .62}, portraitPoint{.53, .66}, portraitPoint{.47, .77}, portraitPoint{.49, .85}, portraitPoint{.48, .90}, portraitPoint{.43, .90}, portraitPoint{.44, .83}, portraitPoint{.41, .75})
	poly(95, portraitPoint{.45, .70}, portraitPoint{.47, .77}, portraitPoint{.49, .85}, portraitPoint{.48, .90}, portraitPoint{.45, .90}, portraitPoint{.46, .83}, portraitPoint{.44, .77})
	poly(137, portraitPoint{.54, .63}, portraitPoint{.63, .64}, portraitPoint{.73, .75}, portraitPoint{.68, .85}, portraitPoint{.79, .90}, portraitPoint{.78, .92}, portraitPoint{.64, .90}, portraitPoint{.63, .84}, portraitPoint{.66, .76}, portraitPoint{.55, .72})
	poly(180, portraitPoint{.60, .67}, portraitPoint{.68, .73}, portraitPoint{.68, .77}, portraitPoint{.65, .76}, portraitPoint{.56, .69})
	poly(237, portraitPoint{.445, .87}, portraitPoint{.49, .87}, portraitPoint{.53, .91}, portraitPoint{.42, .91})
	poly(237, portraitPoint{.66, .87}, portraitPoint{.73, .88}, portraitPoint{.80, .90}, portraitPoint{.80, .92}, portraitPoint{.65, .92})
	// The hunched torso and jagged charcoal mane establish hyena anatomy.
	poly(137, portraitPoint{.48, .37}, portraitPoint{.59, .40}, portraitPoint{.62, .51}, portraitPoint{.64, .62}, portraitPoint{.58, .69}, portraitPoint{.43, .65}, portraitPoint{.39, .55}, portraitPoint{.39, .45})
	poly(95, portraitPoint{.39, .48}, portraitPoint{.45, .47}, portraitPoint{.50, .60}, portraitPoint{.61, .60}, portraitPoint{.58, .69}, portraitPoint{.43, .65}, portraitPoint{.39, .55})
	poly(237, portraitPoint{.47, .33}, portraitPoint{.51, .36}, portraitPoint{.47, .42}, portraitPoint{.45, .49}, portraitPoint{.41, .53}, portraitPoint{.41, .49}, portraitPoint{.38, .50}, portraitPoint{.39, .45}, portraitPoint{.36, .46}, portraitPoint{.40, .39})
	// Blue-green wraps are quiet blocks beneath the warm fur.
	poly(66, portraitPoint{.43, .59}, portraitPoint{.61, .58}, portraitPoint{.64, .66}, portraitPoint{.60, .70}, portraitPoint{.52, .67}, portraitPoint{.45, .69}, portraitPoint{.41, .65})
	poly(60, portraitPoint{.54, .62}, portraitPoint{.63, .62}, portraitPoint{.64, .66}, portraitPoint{.60, .70}, portraitPoint{.52, .67})
	poly(94, portraitPoint{.43, .58}, portraitPoint{.61, .57}, portraitPoint{.62, .61}, portraitPoint{.44, .62})
	// The free arm reaches forward, holding a second stone ready to load.
	poly(137, portraitPoint{.56, .44}, portraitPoint{.62, .44}, portraitPoint{.66, .52}, portraitPoint{.76, .51}, portraitPoint{.79, .55}, portraitPoint{.76, .59}, portraitPoint{.63, .59}, portraitPoint{.58, .54})
	poly(180, portraitPoint{.62, .51}, portraitPoint{.66, .52}, portraitPoint{.76, .51}, portraitPoint{.78, .54}, portraitPoint{.65, .55})
	oval(244, .765, .515, .022, .025)
	// One enormous sling loop is the portrait's primary identifying shape.
	p.ring(180, .59, .145, .235, .12)
	poly(94, portraitPoint{.805, .11}, portraitPoint{.84, .13}, portraitPoint{.85, .16}, portraitPoint{.82, .18}, portraitPoint{.80, .15})
	oval(244, .826, .143, .017, .021)
	poly(180, portraitPoint{.37, .23}, portraitPoint{.39, .22}, portraitPoint{.43, .28}, portraitPoint{.41, .29})
	// Raised elbow and grip form a zigzag against the loop's open center.
	poly(137, portraitPoint{.44, .45}, portraitPoint{.38, .42}, portraitPoint{.31, .33}, portraitPoint{.35, .26}, portraitPoint{.39, .23}, portraitPoint{.43, .27}, portraitPoint{.39, .32}, portraitPoint{.38, .34}, portraitPoint{.48, .38})
	poly(180, portraitPoint{.32, .32}, portraitPoint{.35, .27}, portraitPoint{.39, .25}, portraitPoint{.40, .28}, portraitPoint{.36, .33}, portraitPoint{.40, .39}, portraitPoint{.38, .40})
	poly(137, portraitPoint{.365, .23}, portraitPoint{.39, .22}, portraitPoint{.43, .24}, portraitPoint{.43, .28}, portraitPoint{.39, .29}, portraitPoint{.36, .27})
	// Tall rounded ears, a sloping forehead and long muzzle distinguish a
	// hyena from Grak's broad goblin jaw. The face looks toward the throw.
	poly(137, portraitPoint{.47, .36}, portraitPoint{.45, .26}, portraitPoint{.47, .23}, portraitPoint{.50, .25}, portraitPoint{.53, .35})
	poly(95, portraitPoint{.475, .32}, portraitPoint{.47, .26}, portraitPoint{.49, .27}, portraitPoint{.51, .34})
	poly(137, portraitPoint{.55, .34}, portraitPoint{.55, .235}, portraitPoint{.58, .21}, portraitPoint{.61, .23}, portraitPoint{.60, .35})
	poly(95, portraitPoint{.565, .31}, portraitPoint{.565, .25}, portraitPoint{.585, .24}, portraitPoint{.59, .32})
	poly(137, portraitPoint{.47, .33}, portraitPoint{.57, .32}, portraitPoint{.62, .37}, portraitPoint{.71, .40}, portraitPoint{.72, .445}, portraitPoint{.67, .48}, portraitPoint{.57, .47}, portraitPoint{.51, .49}, portraitPoint{.46, .43})
	poly(180, portraitPoint{.51, .34}, portraitPoint{.57, .335}, portraitPoint{.63, .38}, portraitPoint{.70, .40}, portraitPoint{.69, .43}, portraitPoint{.59, .415}, portraitPoint{.53, .40})
	poly(95, portraitPoint{.47, .40}, portraitPoint{.54, .43}, portraitPoint{.62, .44}, portraitPoint{.67, .46}, portraitPoint{.57, .47}, portraitPoint{.51, .49}, portraitPoint{.46, .43})
	poly(232, portraitPoint{.69, .40}, portraitPoint{.72, .405}, portraitPoint{.73, .43}, portraitPoint{.70, .445}, portraitPoint{.68, .425})
	poly(237, portraitPoint{.59, .435}, portraitPoint{.70, .44}, portraitPoint{.67, .46}, portraitPoint{.60, .455})
	poly(180, portraitPoint{.635, .438}, portraitPoint{.655, .44}, portraitPoint{.645, .465})
	// Four deliberate flank spots; no fur texture competing with the pose.
	oval(95, .51, .52, .018, .024)
	oval(95, .555, .55, .015, .021)
	// The crouched partner has a rounder mass and a lowered, alert face.
	poly(95, portraitPoint{.14, .71}, portraitPoint{.09, .75}, portraitPoint{.07, .82}, portraitPoint{.10, .85}, portraitPoint{.17, .84}, portraitPoint{.18, .80}, portraitPoint{.11, .81}, portraitPoint{.11, .78}, portraitPoint{.18, .76})
	poly(237, portraitPoint{.07, .81}, portraitPoint{.11, .82}, portraitPoint{.11, .85}, portraitPoint{.07, .86}, portraitPoint{.05, .84})
	poly(137, portraitPoint{.16, .66}, portraitPoint{.23, .66}, portraitPoint{.29, .73}, portraitPoint{.28, .81}, portraitPoint{.20, .85}, portraitPoint{.11, .81}, portraitPoint{.10, .74})
	poly(95, portraitPoint{.11, .73}, portraitPoint{.17, .73}, portraitPoint{.20, .80}, portraitPoint{.26, .80}, portraitPoint{.20, .85}, portraitPoint{.11, .81})
	poly(237, portraitPoint{.17, .63}, portraitPoint{.21, .66}, portraitPoint{.17, .70}, portraitPoint{.14, .73}, portraitPoint{.11, .76}, portraitPoint{.11, .71}, portraitPoint{.09, .72}, portraitPoint{.13, .66})
	// Folded haunch, forward hock and long toes make the crouch readable.
	oval(137, .21, .81, .074, .058)
	poly(180, portraitPoint{.19, .77}, portraitPoint{.25, .78}, portraitPoint{.26, .81}, portraitPoint{.21, .83}, portraitPoint{.17, .81})
	poly(95, portraitPoint{.24, .82}, portraitPoint{.26, .84}, portraitPoint{.20, .89}, portraitPoint{.13, .90}, portraitPoint{.12, .87}, portraitPoint{.20, .86})
	poly(237, portraitPoint{.12, .87}, portraitPoint{.20, .88}, portraitPoint{.21, .91}, portraitPoint{.09, .91}, portraitPoint{.09, .89})
	poly(66, portraitPoint{.13, .76}, portraitPoint{.17, .75}, portraitPoint{.20, .79}, portraitPoint{.18, .83}, portraitPoint{.12, .81})
	// The feeder's arm connects the pair through the shared stone pouch.
	poly(137, portraitPoint{.24, .69}, portraitPoint{.28, .70}, portraitPoint{.30, .76}, portraitPoint{.39, .78}, portraitPoint{.40, .82}, portraitPoint{.35, .83}, portraitPoint{.27, .80}, portraitPoint{.24, .75})
	poly(180, portraitPoint{.27, .72}, portraitPoint{.29, .72}, portraitPoint{.31, .76}, portraitPoint{.38, .78}, portraitPoint{.38, .80}, portraitPoint{.29, .79})
	// A loose scarf ties the two silhouettes together without matching armor.
	poly(66, portraitPoint{.19, .64}, portraitPoint{.27, .64}, portraitPoint{.30, .69}, portraitPoint{.25, .72}, portraitPoint{.18, .68})
	poly(60, portraitPoint{.25, .68}, portraitPoint{.29, .70}, portraitPoint{.32, .75}, portraitPoint{.29, .77}, portraitPoint{.27, .72})
	poly(137, portraitPoint{.17, .60}, portraitPoint{.145, .49}, portraitPoint{.16, .47}, portraitPoint{.19, .50}, portraitPoint{.22, .60})
	poly(95, portraitPoint{.175, .58}, portraitPoint{.16, .51}, portraitPoint{.18, .53}, portraitPoint{.20, .59})
	poly(137, portraitPoint{.23, .59}, portraitPoint{.25, .49}, portraitPoint{.28, .48}, portraitPoint{.30, .51}, portraitPoint{.27, .61})
	poly(95, portraitPoint{.25, .57}, portraitPoint{.27, .52}, portraitPoint{.28, .53}, portraitPoint{.265, .59})
	poly(137, portraitPoint{.17, .57}, portraitPoint{.25, .56}, portraitPoint{.29, .60}, portraitPoint{.36, .615}, portraitPoint{.37, .65}, portraitPoint{.32, .68}, portraitPoint{.24, .67}, portraitPoint{.20, .69}, portraitPoint{.16, .64})
	poly(180, portraitPoint{.20, .58}, portraitPoint{.25, .58}, portraitPoint{.29, .61}, portraitPoint{.35, .62}, portraitPoint{.34, .645}, portraitPoint{.26, .63}, portraitPoint{.22, .63})
	poly(95, portraitPoint{.17, .63}, portraitPoint{.24, .65}, portraitPoint{.32, .65}, portraitPoint{.32, .68}, portraitPoint{.24, .67}, portraitPoint{.20, .69})
	poly(232, portraitPoint{.345, .615}, portraitPoint{.37, .62}, portraitPoint{.38, .64}, portraitPoint{.355, .66}, portraitPoint{.335, .64})
	poly(237, portraitPoint{.27, .648}, portraitPoint{.35, .65}, portraitPoint{.33, .665}, portraitPoint{.28, .665})
	// Broad pouch lip and three light stones complete the supply gesture.
	poly(94, portraitPoint{0.310, 0.815}, portraitPoint{0.400, 0.805}, portraitPoint{0.440, 0.845}, portraitPoint{0.430, 0.915}, portraitPoint{0.320, 0.915}, portraitPoint{0.290, 0.865})
	poly(130, portraitPoint{0.320, 0.835}, portraitPoint{0.410, 0.835}, portraitPoint{0.420, 0.875}, portraitPoint{0.400, 0.910}, portraitPoint{0.320, 0.910})
	oval(250, .335, .815, .016, .023)
	oval(253, .380, .800, .016, .025)
	oval(247, .421, .827, .016, .023)
	poly(180, portraitPoint{0.305, 0.835}, portraitPoint{0.415, 0.835}, portraitPoint{0.425, 0.855}, portraitPoint{0.310, 0.865})
	poly(237, portraitPoint{.54, .375}, portraitPoint{.59, .375}, portraitPoint{.59, .40}, portraitPoint{.55, .395})
	p.detail(.56, .385, '━', 220)
	p.detail(.645, .45, '▾', 180)
	p.detail(.245, .61, '━', 232)
	p.detail(.245, .595, '━', 180)
	p.detail(.39, .255, '│', 95)
	p.detail(.745, .56, '│', 95)
	if h < 24 {
		p.detail(.65, .45, '▾', 180)
		p.detail(.31, .66, '▾', 180)
	}
	p.detail(.47, .9, '╲', 180)
	p.detail(.70, .905, '╲', 180)
}
