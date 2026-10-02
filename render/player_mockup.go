package render

// RenderPlayerMockup depicts the guild's champion as the monsters see them:
// bright, composed, and already wearing trophies from someone else's lair.
func RenderPlayerMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the Player", 180, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawPlayerPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-10)/2, 1, "THE PLAYER", 180, 233, true)
		putString(f, (w-20)/2, 2, "the guild's champion", 240, 233, false)
	}
	hint := "tab Cannonier · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawPlayerPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.oval
	poly(236, portraitPoint{.27, .93}, portraitPoint{.79, .93}, portraitPoint{.83, .95}, portraitPoint{.24, .95})
	// The cape falls in two red planes; its edge is a banner, not fur texture.
	poly(95, portraitPoint{.45, .34}, portraitPoint{.63, .34}, portraitPoint{.72, .44}, portraitPoint{.79, .70}, portraitPoint{.84, .88}, portraitPoint{.76, .875}, portraitPoint{.72, .91}, portraitPoint{.62, .87}, portraitPoint{.59, .64}, portraitPoint{.36, .47})
	poly(52, portraitPoint{.63, .37}, portraitPoint{.69, .43}, portraitPoint{.72, .66}, portraitPoint{.78, .875}, portraitPoint{.72, .91}, portraitPoint{.62, .87}, portraitPoint{.59, .64})
	poly(131, portraitPoint{.70, .48}, portraitPoint{.72, .49}, portraitPoint{.79, .71}, portraitPoint{.84, .88}, portraitPoint{.80, .875}, portraitPoint{.75, .70})
	// Long plate-clad legs and pointed sabatons make the stance aristocratic.
	poly(67, portraitPoint{.41, .57}, portraitPoint{.53, .59}, portraitPoint{.51, .74}, portraitPoint{.49, .88}, portraitPoint{.40, .91}, portraitPoint{.38, .84}, portraitPoint{.40, .71})
	poly(251, portraitPoint{.42, .62}, portraitPoint{.49, .64}, portraitPoint{.48, .74}, portraitPoint{.46, .87}, portraitPoint{.41, .88}, portraitPoint{.40, .82})
	poly(255, portraitPoint{.42, .64}, portraitPoint{.445, .65}, portraitPoint{.445, .73}, portraitPoint{.425, .86}, portraitPoint{.41, .87})
	poly(67, portraitPoint{.53, .58}, portraitPoint{.62, .60}, portraitPoint{.66, .75}, portraitPoint{.69, .87}, portraitPoint{.75, .91}, portraitPoint{.64, .91}, portraitPoint{.59, .82}, portraitPoint{.56, .73})
	poly(251, portraitPoint{.57, .64}, portraitPoint{.615, .64}, portraitPoint{.65, .76}, portraitPoint{.67, .87}, portraitPoint{.65, .89}, portraitPoint{.61, .81})
	poly(110, portraitPoint{.615, .64}, portraitPoint{.635, .67}, portraitPoint{.67, .80}, portraitPoint{.69, .87}, portraitPoint{.67, .87}, portraitPoint{.65, .76})
	poly(180, portraitPoint{.41, .72}, portraitPoint{.49, .725}, portraitPoint{.495, .745}, portraitPoint{.41, .74})
	poly(180, portraitPoint{.60, .72}, portraitPoint{.65, .715}, portraitPoint{.66, .74}, portraitPoint{.61, .75})
	poly(239, portraitPoint{.385, .87}, portraitPoint{.48, .87}, portraitPoint{.51, .91}, portraitPoint{.48, .935}, portraitPoint{.31, .935}, portraitPoint{.32, .92})
	poly(251, portraitPoint{.385, .88}, portraitPoint{.465, .88}, portraitPoint{.48, .905}, portraitPoint{.45, .92}, portraitPoint{.32, .92})
	poly(239, portraitPoint{.65, .87}, portraitPoint{.69, .87}, portraitPoint{.78, .91}, portraitPoint{.80, .935}, portraitPoint{.65, .935}, portraitPoint{.62, .91})
	poly(251, portraitPoint{.655, .88}, portraitPoint{.69, .89}, portraitPoint{.765, .91}, portraitPoint{.78, .92}, portraitPoint{.66, .92})
	// The greatsword is a single enormous diagonal behind the shoulder.
	poly(239, portraitPoint{.065, .06}, portraitPoint{.23, .125}, portraitPoint{.425, .32}, portraitPoint{.41, .365}, portraitPoint{.35, .405}, portraitPoint{.12, .165})
	poly(251, portraitPoint{.08, .078}, portraitPoint{.22, .135}, portraitPoint{.41, .325}, portraitPoint{.395, .355}, portraitPoint{.35, .385}, portraitPoint{.135, .165})
	poly(255, portraitPoint{.08, .078}, portraitPoint{.155, .11}, portraitPoint{.375, .35}, portraitPoint{.35, .385}, portraitPoint{.135, .165})
	poly(110, portraitPoint{.155, .11}, portraitPoint{.22, .135}, portraitPoint{.41, .325}, portraitPoint{.395, .355}, portraitPoint{.375, .35})
	poly(180, portraitPoint{.30, .395}, portraitPoint{.31, .36}, portraitPoint{.43, .325}, portraitPoint{.47, .34}, portraitPoint{.485, .37}, portraitPoint{.455, .385}, portraitPoint{.43, .36}, portraitPoint{.335, .395}, portraitPoint{.32, .43})
	poly(130, portraitPoint{.355, .39}, portraitPoint{.385, .38}, portraitPoint{.425, .445}, portraitPoint{.405, .465})
	oval(180, .417, .454, .024, .021)
	// The torso is a polished shield shape with a large, quiet guild emblem.
	poly(239, portraitPoint{.40, .335}, portraitPoint{.57, .33}, portraitPoint{.65, .405}, portraitPoint{.62, .525}, portraitPoint{.56, .605}, portraitPoint{.42, .61}, portraitPoint{.345, .515}, portraitPoint{.35, .415})
	poly(251, portraitPoint{.405, .35}, portraitPoint{.555, .345}, portraitPoint{.61, .405}, portraitPoint{.58, .52}, portraitPoint{.53, .58}, portraitPoint{.435, .585}, portraitPoint{.375, .50}, portraitPoint{.375, .415})
	poly(255, portraitPoint{.405, .36}, portraitPoint{.49, .35}, portraitPoint{.515, .46}, portraitPoint{.49, .54}, portraitPoint{.435, .585}, portraitPoint{.375, .50}, portraitPoint{.385, .425})
	poly(110, portraitPoint{.555, .35}, portraitPoint{.61, .405}, portraitPoint{.58, .52}, portraitPoint{.53, .58}, portraitPoint{.49, .54}, portraitPoint{.515, .46})
	poly(180, portraitPoint{.395, .35}, portraitPoint{.55, .345}, portraitPoint{.585, .375}, portraitPoint{.55, .365}, portraitPoint{.41, .375})
	poly(180, portraitPoint{.49, .425}, portraitPoint{.525, .455}, portraitPoint{.505, .49}, portraitPoint{.47, .46})
	poly(130, portraitPoint{.505, .443}, portraitPoint{.525, .455}, portraitPoint{.505, .49}, portraitPoint{.505, .464})
	// Overlapping waist plates spread above the belt and trophies.
	poly(239, portraitPoint{.39, .53}, portraitPoint{.56, .535}, portraitPoint{.63, .59}, portraitPoint{.60, .635}, portraitPoint{.38, .64}, portraitPoint{.35, .595})
	poly(251, portraitPoint{.39, .55}, portraitPoint{.54, .55}, portraitPoint{.58, .58}, portraitPoint{.55, .60}, portraitPoint{.38, .61}, portraitPoint{.37, .59})
	poly(110, portraitPoint{.55, .55}, portraitPoint{.60, .585}, portraitPoint{.60, .62}, portraitPoint{.55, .62}, portraitPoint{.53, .60})
	poly(94, portraitPoint{.36, .61}, portraitPoint{.62, .61}, portraitPoint{.62, .65}, portraitPoint{.36, .655})
	poly(180, portraitPoint{.465, .60}, portraitPoint{.515, .60}, portraitPoint{.525, .655}, portraitPoint{.46, .66})
	poly(95, portraitPoint{.477, .618}, portraitPoint{.505, .618}, portraitPoint{.508, .64}, portraitPoint{.476, .642})
	// The extended arm opens toward a small fragment of someone else's hoard.
	poly(239, portraitPoint{.60, .38}, portraitPoint{.69, .39}, portraitPoint{.74, .48}, portraitPoint{.84, .515}, portraitPoint{.84, .57}, portraitPoint{.71, .56}, portraitPoint{.66, .50}, portraitPoint{.60, .47})
	poly(251, portraitPoint{.65, .42}, portraitPoint{.69, .43}, portraitPoint{.735, .495}, portraitPoint{.825, .53}, portraitPoint{.825, .55}, portraitPoint{.715, .525}, portraitPoint{.675, .49})
	poly(110, portraitPoint{.715, .525}, portraitPoint{.825, .55}, portraitPoint{.84, .545}, portraitPoint{.84, .57}, portraitPoint{.71, .56}, portraitPoint{.68, .515})
	poly(180, portraitPoint{.815, .513}, portraitPoint{.85, .525}, portraitPoint{.85, .575}, portraitPoint{.82, .567})
	poly(180, portraitPoint{.85, .533}, portraitPoint{.88, .535}, portraitPoint{.92, .52}, portraitPoint{.95, .50}, portraitPoint{.968, .51}, portraitPoint{.95, .54}, portraitPoint{.98, .535}, portraitPoint{.985, .55}, portraitPoint{.955, .565}, portraitPoint{.91, .585}, portraitPoint{.87, .575}, portraitPoint{.845, .56})
	poly(223, portraitPoint{.855, .535}, portraitPoint{.88, .54}, portraitPoint{.925, .525}, portraitPoint{.95, .50}, portraitPoint{.958, .51}, portraitPoint{.938, .55}, portraitPoint{.91, .566}, portraitPoint{.86, .557})
	// Broad pauldrons catch the same light as the blade and breastplate.
	poly(239, portraitPoint{.32, .37}, portraitPoint{.39, .355}, portraitPoint{.435, .39}, portraitPoint{.41, .46}, portraitPoint{.31, .485}, portraitPoint{.275, .44}, portraitPoint{.285, .405})
	poly(251, portraitPoint{.32, .385}, portraitPoint{.385, .37}, portraitPoint{.415, .40}, portraitPoint{.40, .43}, portraitPoint{.305, .45}, portraitPoint{.29, .425})
	poly(255, portraitPoint{.32, .385}, portraitPoint{.385, .37}, portraitPoint{.40, .385}, portraitPoint{.31, .415}, portraitPoint{.295, .44}, portraitPoint{.29, .425})
	poly(180, portraitPoint{.31, .46}, portraitPoint{.41, .435}, portraitPoint{.41, .46}, portraitPoint{.31, .485}, portraitPoint{.29, .46})
	poly(239, portraitPoint{.59, .355}, portraitPoint{.655, .36}, portraitPoint{.72, .405}, portraitPoint{.72, .44}, portraitPoint{.67, .475}, portraitPoint{.58, .44}, portraitPoint{.56, .39})
	poly(251, portraitPoint{.59, .37}, portraitPoint{.645, .375}, portraitPoint{.70, .41}, portraitPoint{.70, .435}, portraitPoint{.67, .445}, portraitPoint{.59, .415}, portraitPoint{.575, .39})
	poly(255, portraitPoint{.59, .37}, portraitPoint{.645, .375}, portraitPoint{.685, .40}, portraitPoint{.615, .395}, portraitPoint{.575, .39})
	poly(180, portraitPoint{.59, .42}, portraitPoint{.67, .45}, portraitPoint{.72, .43}, portraitPoint{.72, .44}, portraitPoint{.67, .475}, portraitPoint{.59, .445})
	// The sword arm bends up casually; its weight rests on the shoulder.
	poly(239, portraitPoint{.30, .45}, portraitPoint{.35, .44}, portraitPoint{.37, .51}, portraitPoint{.40, .45}, portraitPoint{.43, .425}, portraitPoint{.46, .445}, portraitPoint{.44, .49}, portraitPoint{.38, .57}, portraitPoint{.33, .57}, portraitPoint{.285, .51})
	poly(251, portraitPoint{.31, .465}, portraitPoint{.33, .46}, portraitPoint{.36, .53}, portraitPoint{.41, .45}, portraitPoint{.43, .445}, portraitPoint{.43, .475}, portraitPoint{.37, .55}, portraitPoint{.34, .54})
	poly(180, portraitPoint{.39, .445}, portraitPoint{.425, .425}, portraitPoint{.445, .44}, portraitPoint{.425, .48}, portraitPoint{.40, .47})
	// A single curved red crest balances the enormous blade to the left.
	poly(95, portraitPoint{.44, .18}, portraitPoint{.46, .11}, portraitPoint{.515, .095}, portraitPoint{.58, .115}, portraitPoint{.625, .16}, portraitPoint{.61, .19}, portraitPoint{.55, .14}, portraitPoint{.49, .14}, portraitPoint{.48, .19})
	poly(131, portraitPoint{.46, .11}, portraitPoint{.515, .095}, portraitPoint{.58, .115}, portraitPoint{.625, .16}, portraitPoint{.59, .155}, portraitPoint{.54, .125}, portraitPoint{.49, .125}, portraitPoint{.465, .15})
	// An open-faced helmet leaves a human, calmly satisfied expression visible.
	poly(130, portraitPoint{.42, .19}, portraitPoint{.50, .18}, portraitPoint{.57, .22}, portraitPoint{.57, .32}, portraitPoint{.53, .365}, portraitPoint{.43, .345}, portraitPoint{.405, .27})
	poly(180, portraitPoint{.44, .22}, portraitPoint{.52, .21}, portraitPoint{.56, .245}, portraitPoint{.56, .29}, portraitPoint{.59, .30}, portraitPoint{.57, .325}, portraitPoint{.55, .35}, portraitPoint{.49, .37}, portraitPoint{.44, .34}, portraitPoint{.43, .28})
	poly(223, portraitPoint{.45, .23}, portraitPoint{.51, .225}, portraitPoint{.54, .25}, portraitPoint{.54, .295}, portraitPoint{.575, .305}, portraitPoint{.55, .32}, portraitPoint{.53, .345}, portraitPoint{.49, .35}, portraitPoint{.45, .325})
	poly(230, portraitPoint{.46, .235}, portraitPoint{.505, .23}, portraitPoint{.525, .25}, portraitPoint{.515, .28}, portraitPoint{.46, .28})
	poly(130, portraitPoint{.435, .24}, portraitPoint{.47, .255}, portraitPoint{.45, .305}, portraitPoint{.43, .29})
	poly(239, portraitPoint{.40, .22}, portraitPoint{.415, .165}, portraitPoint{.455, .145}, portraitPoint{.525, .15}, portraitPoint{.565, .19}, portraitPoint{.58, .235}, portraitPoint{.55, .25}, portraitPoint{.445, .245})
	poly(251, portraitPoint{.415, .205}, portraitPoint{.43, .175}, portraitPoint{.465, .16}, portraitPoint{.51, .165}, portraitPoint{.545, .20}, portraitPoint{.545, .227}, portraitPoint{.44, .225})
	poly(255, portraitPoint{.435, .175}, portraitPoint{.465, .16}, portraitPoint{.475, .19}, portraitPoint{.46, .22}, portraitPoint{.425, .22})
	poly(110, portraitPoint{.51, .165}, portraitPoint{.545, .20}, portraitPoint{.56, .235}, portraitPoint{.545, .24}, portraitPoint{.52, .225})
	poly(180, portraitPoint{.407, .22}, portraitPoint{.45, .227}, portraitPoint{.55, .227}, portraitPoint{.57, .235}, portraitPoint{.565, .25}, portraitPoint{.445, .25}, portraitPoint{.40, .235})
	// A bright cheek plate protects one side, leaving the reaching side open.
	poly(251, portraitPoint{.405, .24}, portraitPoint{.435, .25}, portraitPoint{.445, .30}, portraitPoint{.435, .335}, portraitPoint{.405, .30})
	poly(110, portraitPoint{.405, .27}, portraitPoint{.425, .28}, portraitPoint{.435, .335}, portraitPoint{.405, .30})
	// Trophy one: a broken curved tusk hanging from a short leather cord.
	poly(94, portraitPoint{.37, .645}, portraitPoint{.38, .645}, portraitPoint{.37, .72}, portraitPoint{.355, .73}, portraitPoint{.35, .715})
	poly(180, portraitPoint{.345, .705}, portraitPoint{.375, .715}, portraitPoint{.37, .765}, portraitPoint{.345, .795}, portraitPoint{.30, .80}, portraitPoint{.33, .77}, portraitPoint{.345, .745})
	poly(223, portraitPoint{.345, .715}, portraitPoint{.36, .72}, portraitPoint{.355, .76}, portraitPoint{.33, .785}, portraitPoint{.315, .79}, portraitPoint{.345, .745})
	// Trophy two: a long-muzzled monster skull, with ears and empty eye sockets.
	poly(94, portraitPoint{.605, .64}, portraitPoint{.62, .64}, portraitPoint{.64, .68}, portraitPoint{.625, .69})
	poly(144, portraitPoint{.605, .69}, portraitPoint{.59, .65}, portraitPoint{.63, .675}, portraitPoint{.67, .675}, portraitPoint{.70, .65}, portraitPoint{.695, .705}, portraitPoint{.73, .73}, portraitPoint{.71, .765}, portraitPoint{.68, .755}, portraitPoint{.65, .775}, portraitPoint{.61, .75}, portraitPoint{.59, .72})
	poly(230, portraitPoint{.615, .69}, portraitPoint{.655, .685}, portraitPoint{.68, .705}, portraitPoint{.72, .73}, portraitPoint{.695, .75}, portraitPoint{.65, .755}, portraitPoint{.61, .73})
	poly(255, portraitPoint{.62, .69}, portraitPoint{.65, .69}, portraitPoint{.67, .71}, portraitPoint{.63, .715}, portraitPoint{.61, .72})
	// Angular sockets, a nasal cut and a broken jaw keep the trophy skeletal.
	poly(232, portraitPoint{.615, .715}, portraitPoint{.637, .705}, portraitPoint{.648, .723}, portraitPoint{.633, .740}, portraitPoint{.619, .733})
	poly(232, portraitPoint{.667, .717}, portraitPoint{.682, .710}, portraitPoint{.692, .725}, portraitPoint{.680, .740}, portraitPoint{.666, .734})
	poly(239, portraitPoint{.704, .735}, portraitPoint{.719, .732}, portraitPoint{.715, .752}, portraitPoint{.700, .751})
	poly(230, portraitPoint{.62, .753}, portraitPoint{.637, .752}, portraitPoint{.638, .783}, portraitPoint{.626, .778})
	poly(230, portraitPoint{.655, .756}, portraitPoint{.671, .755}, portraitPoint{.667, .779}, portraitPoint{.654, .779})
	// A small, low hill of gold completes the gesture without becoming scenery.
	poly(130, portraitPoint{.83, .89}, portraitPoint{.86, .85}, portraitPoint{.90, .84}, portraitPoint{.925, .87}, portraitPoint{.965, .90}, portraitPoint{.975, .94}, portraitPoint{.805, .94})
	poly(178, portraitPoint{.83, .89}, portraitPoint{.86, .85}, portraitPoint{.90, .84}, portraitPoint{.925, .87}, portraitPoint{.965, .90}, portraitPoint{.945, .91}, portraitPoint{.90, .865}, portraitPoint{.865, .875})
	poly(130, portraitPoint{.49, .261}, portraitPoint{.527, .258}, portraitPoint{.536, .276}, portraitPoint{.491, .278})
	p.detail(.505, .28, '─', 232)
	p.detail(.528, .335, '─', 130)
	p.detail(.547, .331, '╯', 130)
	p.detail(.48, .455, '╲', 255)
	// Half-block sockets survive the coarse portrait sizes without becoming
	// round cartoon eyes; their background retains the ivory skull plane.
	p.detail(.63, .72, '▀', 232)
	p.detail(.68, .72, '▀', 232)
	p.detail(.705, .745, '▾', 239)
	p.detail(.895, .90, '◆', 220)
	p.detail(.85, .92, '◆', 178)
	if h >= 24 {
		p.detail(.42, .895, '╱', 255)
		p.detail(.36, .505, '╱', 255)
		p.detail(.55, .50, '╲', 251)
	}
}
