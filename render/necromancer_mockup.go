package render

// RenderNecromancerMockup depicts the guild's death-worker above two returning
// squires. A bent back, hooked staff and suspended hand frame the summons.
func RenderNecromancerMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the Necromancer", 151, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawNecromancerPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		title, subtitle := "NECROMANCER", "they rise when he falls"
		putString(f, (w-len(title))/2, 1, title, 151, 233, true)
		putString(f, (w-len(subtitle))/2, 2, subtitle, 240, 233, false)
	}
	hint := "tab Paladin · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawNecromancerPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.smoothOval
	poly(236, portraitPoint{.09, .945}, portraitPoint{.87, .945}, portraitPoint{.92, .971}, portraitPoint{.07, .971})
	// A bowed spine and pooled coat make a heavy arch over the smaller figures.
	poly(236, portraitPoint{.43, .24}, portraitPoint{.55, .27}, portraitPoint{.66, .36}, portraitPoint{.705, .49}, portraitPoint{.68, .60}, portraitPoint{.72, .72}, portraitPoint{.83, .91}, portraitPoint{.85, .95}, portraitPoint{.69, .96}, portraitPoint{.58, .92}, portraitPoint{.43, .955}, portraitPoint{.31, .95}, portraitPoint{.34, .81}, portraitPoint{.375, .63}, portraitPoint{.32, .47}, portraitPoint{.35, .345})
	poly(60, portraitPoint{.435, .265}, portraitPoint{.55, .295}, portraitPoint{.64, .37}, portraitPoint{.68, .485}, portraitPoint{.645, .60}, portraitPoint{.685, .74}, portraitPoint{.79, .925}, portraitPoint{.69, .945}, portraitPoint{.58, .90}, portraitPoint{.46, .94}, portraitPoint{.34, .932}, portraitPoint{.38, .77}, portraitPoint{.415, .64}, portraitPoint{.35, .475})
	poly(237, portraitPoint{.435, .36}, portraitPoint{.515, .42}, portraitPoint{.48, .61}, portraitPoint{.46, .76}, portraitPoint{.41, .93}, portraitPoint{.34, .932}, portraitPoint{.38, .77}, portraitPoint{.415, .64}, portraitPoint{.35, .475})
	poly(96, portraitPoint{.54, .315}, portraitPoint{.625, .38}, portraitPoint{.66, .48}, portraitPoint{.623, .61}, portraitPoint{.66, .75}, portraitPoint{.746, .921}, portraitPoint{.69, .937}, portraitPoint{.605, .83}, portraitPoint{.57, .64}, portraitPoint{.595, .485})
	poly(95, portraitPoint{.62, .402}, portraitPoint{.638, .42}, portraitPoint{.656, .482}, portraitPoint{.633, .545}, portraitPoint{.62, .535}, portraitPoint{.634, .475})
	poly(236, portraitPoint{.535, .58}, portraitPoint{.585, .57}, portraitPoint{.555, .71}, portraitPoint{.58, .905}, portraitPoint{.53, .921}, portraitPoint{.515, .75})
	poly(64, portraitPoint{.675, .768}, portraitPoint{.695, .795}, portraitPoint{.79, .925}, portraitPoint{.755, .936}, portraitPoint{.7, .846})
	// Broad cuffs and wine-red lining belong to a human guild robe.
	poly(236, portraitPoint{.355, .39}, portraitPoint{.405, .425}, portraitPoint{.385, .49}, portraitPoint{.29, .55}, portraitPoint{.19, .53}, portraitPoint{.15, .485}, portraitPoint{.22, .465}, portraitPoint{.295, .474})
	poly(60, portraitPoint{.353, .41}, portraitPoint{.39, .437}, portraitPoint{.36, .48}, portraitPoint{.283, .52}, portraitPoint{.203, .511}, portraitPoint{.18, .481}, portraitPoint{.224, .48}, portraitPoint{.292, .489})
	poly(95, portraitPoint{.182, .482}, portraitPoint{.227, .481}, portraitPoint{.239, .527}, portraitPoint{.191, .53}, portraitPoint{.162, .49})
	poly(236, portraitPoint{.535, .365}, portraitPoint{.61, .38}, portraitPoint{.695, .433}, portraitPoint{.705, .506}, portraitPoint{.675, .556}, portraitPoint{.61, .53}, portraitPoint{.525, .455})
	poly(60, portraitPoint{.55, .381}, portraitPoint{.61, .398}, portraitPoint{.672, .442}, portraitPoint{.68, .491}, portraitPoint{.651, .527}, portraitPoint{.603, .487}, portraitPoint{.547, .434})
	poly(96, portraitPoint{.573, .395}, portraitPoint{.612, .409}, portraitPoint{.668, .444}, portraitPoint{.659, .473}, portraitPoint{.611, .45})
	poly(95, portraitPoint{.668, .445}, portraitPoint{.697, .455}, portraitPoint{.698, .497}, portraitPoint{.675, .547}, portraitPoint{.647, .528}, portraitPoint{.658, .491})
	// A narrow, faded stole and simple bone clasp emphasize the hunched chest.
	poly(95, portraitPoint{.428, .37}, portraitPoint{.46, .39}, portraitPoint{.455, .495}, portraitPoint{.49, .57}, portraitPoint{.47, .65}, portraitPoint{.444, .628}, portraitPoint{.459, .573}, portraitPoint{.424, .495})
	poly(137, portraitPoint{.45, .391}, portraitPoint{.459, .399}, portraitPoint{.45, .49}, portraitPoint{.48, .571}, portraitPoint{.475, .595}, portraitPoint{.44, .49})
	oval(144, .441, .419, .028, .028)
	poly(230, portraitPoint{.42, .407}, portraitPoint{.445, .393}, portraitPoint{.461, .413}, portraitPoint{.453, .438}, portraitPoint{.434, .437})
	poly(236, portraitPoint{.428, .415}, portraitPoint{.439, .412}, portraitPoint{.44, .423}, portraitPoint{.429, .426})
	poly(236, portraitPoint{.446, .414}, portraitPoint{.456, .417}, portraitPoint{.452, .429}, portraitPoint{.444, .425})
	// A sloping cowl opens around a human face: bald brow, hooked nose, thin jaw.
	poly(236, portraitPoint{.405, .13}, portraitPoint{.485, .145}, portraitPoint{.56, .208}, portraitPoint{.589, .298}, portraitPoint{.552, .385}, portraitPoint{.464, .417}, portraitPoint{.373, .374}, portraitPoint{.327, .30}, portraitPoint{.35, .202})
	poly(60, portraitPoint{.405, .13}, portraitPoint{.485, .145}, portraitPoint{.56, .208}, portraitPoint{.58, .288}, portraitPoint{.555, .327}, portraitPoint{.53, .257}, portraitPoint{.47, .216}, portraitPoint{.413, .199}, portraitPoint{.357, .275}, portraitPoint{.35, .202})
	poly(96, portraitPoint{.405, .13}, portraitPoint{.485, .145}, portraitPoint{.548, .201}, portraitPoint{.549, .233}, portraitPoint{.481, .185}, portraitPoint{.424, .18})
	poly(234, portraitPoint{.41, .213}, portraitPoint{.48, .213}, portraitPoint{.541, .26}, portraitPoint{.557, .32}, portraitPoint{.527, .366}, portraitPoint{.465, .388}, portraitPoint{.391, .339}, portraitPoint{.376, .277})
	poly(108, portraitPoint{.43, .222}, portraitPoint{.48, .227}, portraitPoint{.523, .253}, portraitPoint{.536, .291}, portraitPoint{.564, .322}, portraitPoint{.55, .336}, portraitPoint{.529, .336}, portraitPoint{.514, .374}, portraitPoint{.477, .387}, portraitPoint{.43, .357}, portraitPoint{.416, .305}, portraitPoint{.407, .267})
	poly(144, portraitPoint{.432, .228}, portraitPoint{.478, .233}, portraitPoint{.515, .258}, portraitPoint{.517, .294}, portraitPoint{.495, .316}, portraitPoint{.439, .306}, portraitPoint{.421, .269})
	poly(230, portraitPoint{.435, .236}, portraitPoint{.472, .24}, portraitPoint{.498, .259}, portraitPoint{.49, .277}, portraitPoint{.438, .275}, portraitPoint{.42, .261})
	poly(65, portraitPoint{.425, .307}, portraitPoint{.463, .322}, portraitPoint{.49, .32}, portraitPoint{.494, .342}, portraitPoint{.48, .362}, portraitPoint{.435, .347})
	poly(144, portraitPoint{.486, .343}, portraitPoint{.527, .331}, portraitPoint{.517, .37}, portraitPoint{.49, .378}, portraitPoint{.478, .361})
	poly(230, portraitPoint{.519, .287}, portraitPoint{.555, .317}, portraitPoint{.549, .327}, portraitPoint{.522, .318})
	poly(234, portraitPoint{.472, .286}, portraitPoint{.502, .279}, portraitPoint{.521, .296}, portraitPoint{.495, .312}, portraitPoint{.473, .306})
	p.stroke(151, 1.3, portraitPoint{.482, .295}, portraitPoint{.507, .293})
	poly(236, portraitPoint{.485, .352}, portraitPoint{.522, .342}, portraitPoint{.519, .354}, portraitPoint{.489, .364})
	// The crooked staff has a bone hook and a skull hung inside its opening.
	shaft := []portraitPoint{{.135, .929}, {.159, .68}, {.183, .44}, {.173, .293}, {.135, .211}}
	p.stroke(95, 3, shaft...)
	p.stroke(137, 1.8, shaft...)
	lit := make([]portraitPoint, len(shaft))
	for i, a := range shaft {
		lit[i] = portraitPoint{a.x - .007, a.y}
	}
	p.stroke(180, .7, lit...)
	crook := portraitCurve(portraitPoint{.135, .211}, portraitPoint{.025, .045}, portraitPoint{.227, -.025}, portraitPoint{.30, .083})
	crook = append(crook, portraitCurve(portraitPoint{.30, .083}, portraitPoint{.346, .125}, portraitPoint{.322, .179}, portraitPoint{.284, .18})[1:]...)
	p.stroke(94, 3.3, crook...)
	p.stroke(144, 2.3, crook...)
	p.stroke(230, 1.0, crook...)
	p.stroke(137, .8, portraitPoint{.286, .175}, portraitPoint{.274, .231})
	oval(144, .264, .243, .052, .065)
	poly(230, portraitPoint{.223, .214}, portraitPoint{.269, .19}, portraitPoint{.299, .215}, portraitPoint{.307, .254}, portraitPoint{.288, .278}, portraitPoint{.28, .307}, portraitPoint{.242, .305}, portraitPoint{.236, .278}, portraitPoint{.216, .255})
	poly(144, portraitPoint{.279, .207}, portraitPoint{.299, .215}, portraitPoint{.307, .254}, portraitPoint{.288, .278}, portraitPoint{.28, .307}, portraitPoint{.265, .303}, portraitPoint{.273, .26})
	if h >= 18 {
		oval(234, .239, .246, .014, .017)
		oval(234, .28, .245, .014, .017)
		poly(236, portraitPoint{.258, .258}, portraitPoint{.269, .274}, portraitPoint{.251, .276})
		p.stroke(95, .9, portraitPoint{.24, .288}, portraitPoint{.281, .289})
	} else {
		// A two-cell skull needs one dark feature, retaining its pale jaw.
		p.detail(.255, .28, '▪', 234)
	}
	// Staff fingers curl around the wood; the hand retains an exposed wrist.
	poly(108, portraitPoint{.194, .477}, portraitPoint{.21, .478}, portraitPoint{.213, .506}, portraitPoint{.182, .512}, portraitPoint{.161, .495}, portraitPoint{.16, .467}, portraitPoint{.174, .457}, portraitPoint{.182, .482})
	poly(230, portraitPoint{.161, .468}, portraitPoint{.174, .459}, portraitPoint{.177, .482}, portraitPoint{.194, .487}, portraitPoint{.205, .485}, portraitPoint{.203, .495}, portraitPoint{.176, .498}, portraitPoint{.165, .488})
	// Two loose threads of green light link the hand to returning helmets.
	thread := func(path ...portraitPoint) {
		p.stroke(23, 1.8, path...)
		p.stroke(65, 1.1, path...)
		p.stroke(151, .75, path...)
	}
	thread(portraitPoint{.719, .554}, portraitPoint{.65, .606}, portraitPoint{.55, .581}, portraitPoint{.43, .647}, portraitPoint{.341, .71})
	thread(portraitPoint{.785, .553}, portraitPoint{.805, .614}, portraitPoint{.737, .634}, portraitPoint{.708, .71})
	oval(235, .326, .932, .173, .027)
	oval(235, .709, .939, .16, .027)
	poly(65, portraitPoint{.21, .913}, portraitPoint{.23, .888}, portraitPoint{.243, .925}, portraitPoint{.22, .941})
	poly(65, portraitPoint{.82, .905}, portraitPoint{.836, .876}, portraitPoint{.848, .927}, portraitPoint{.823, .941})
	// The near palm and hooked fingers frame sky above the two subordinate heads.
	poly(108, portraitPoint{.674, .461}, portraitPoint{.706, .456}, portraitPoint{.743, .477}, portraitPoint{.81, .487}, portraitPoint{.837, .512}, portraitPoint{.827, .557}, portraitPoint{.81, .563}, portraitPoint{.814, .526}, portraitPoint{.791, .516}, portraitPoint{.788, .56}, portraitPoint{.763, .605}, portraitPoint{.745, .6}, portraitPoint{.763, .554}, portraitPoint{.756, .518}, portraitPoint{.742, .55}, portraitPoint{.719, .583}, portraitPoint{.703, .576}, portraitPoint{.721, .535}, portraitPoint{.715, .512}, portraitPoint{.686, .531}, portraitPoint{.662, .516}, portraitPoint{.66, .489})
	poly(230, portraitPoint{.676, .47}, portraitPoint{.701, .467}, portraitPoint{.736, .488}, portraitPoint{.81, .498}, portraitPoint{.827, .516}, portraitPoint{.821, .542}, portraitPoint{.82, .521}, portraitPoint{.785, .508}, portraitPoint{.777, .553}, portraitPoint{.753, .594}, portraitPoint{.766, .555}, portraitPoint{.767, .512}, portraitPoint{.749, .511}, portraitPoint{.73, .549}, portraitPoint{.715, .572}, portraitPoint{.726, .54}, portraitPoint{.733, .508}, portraitPoint{.703, .497}, portraitPoint{.681, .512}, portraitPoint{.67, .50})
	poly(144, portraitPoint{.68, .481}, portraitPoint{.706, .48}, portraitPoint{.725, .498}, portraitPoint{.712, .511}, portraitPoint{.688, .513}, portraitPoint{.676, .50})
	necromancerSquire(p, .33, .80, true)
	necromancerSquire(p, .70, .79, false)
	// Fine features deepen the faces and grips without punching holes in them.
	p.detail(.445, .32, '╲', 144)
	p.detail(.49, .357, '─', 108)
	p.detail(.18, .49, '│', 108)
	p.detail(.702, .491, '╲', 108)
	p.detail(.765, .525, '│', 108)
	if h >= 18 {
		p.detail(.254, .297, '│', 144)
		p.detail(.272, .297, '│', 144)
	}
	if h >= 24 {
		p.detail(.481, .263, '╲', 144)
		p.detail(.468, .345, '╱', 108)
		p.detail(.604, .453, '╲', 96)
		p.detail(.553, .769, '│', 60)
	}
}

// The returning squires share guild kit, but one is still hauling a knee up.
func necromancerSquire(p portraitPainter, cx, cy float64, kneeling bool) {
	poly := func(fg int, points ...portraitPoint) {
		for i, a := range points {
			points[i] = portraitPoint{cx + a.x*.16, cy + a.y*.28}
		}
		p.poly(fg, points...)
	}
	stroke := func(fg int, width float64, points ...portraitPoint) {
		for i, a := range points {
			points[i] = portraitPoint{cx + a.x*.16, cy + a.y*.28}
		}
		p.stroke(fg, width, points...)
	}
	if kneeling {
		poly(239, portraitPoint{-.18, .12}, portraitPoint{.13, .18}, portraitPoint{-.06, .35}, portraitPoint{-.39, .43}, portraitPoint{-.55, .49}, portraitPoint{-.62, .45}, portraitPoint{-.32, .29})
		poly(246, portraitPoint{-.14, .19}, portraitPoint{-.06, .24}, portraitPoint{-.36, .39}, portraitPoint{-.48, .43}, portraitPoint{-.46, .36})
		poly(239, portraitPoint{.10, .12}, portraitPoint{.32, .17}, portraitPoint{.49, .31}, portraitPoint{.39, .43}, portraitPoint{.61, .46}, portraitPoint{.61, .52}, portraitPoint{.23, .51}, portraitPoint{.19, .41}, portraitPoint{.28, .30})
		poly(109, portraitPoint{.24, .22}, portraitPoint{.34, .22}, portraitPoint{.42, .32}, portraitPoint{.29, .44}, portraitPoint{.23, .39}, portraitPoint{.31, .29})
	} else {
		poly(239, portraitPoint{-.25, .10}, portraitPoint{-.02, .13}, portraitPoint{-.10, .35}, portraitPoint{-.23, .47}, portraitPoint{-.46, .51}, portraitPoint{-.5, .47}, portraitPoint{-.26, .40})
		poly(246, portraitPoint{-.18, .18}, portraitPoint{-.08, .19}, portraitPoint{-.17, .34}, portraitPoint{-.25, .43}, portraitPoint{-.35, .45}, portraitPoint{-.25, .34})
		poly(239, portraitPoint{.02, .10}, portraitPoint{.23, .10}, portraitPoint{.30, .33}, portraitPoint{.40, .44}, portraitPoint{.53, .46}, portraitPoint{.55, .52}, portraitPoint{.23, .52}, portraitPoint{.11, .34})
		poly(109, portraitPoint{.10, .19}, portraitPoint{.19, .19}, portraitPoint{.20, .32}, portraitPoint{.30, .46}, portraitPoint{.23, .47}, portraitPoint{.12, .34})
	}
	poly(95, portraitPoint{-.20, -.30}, portraitPoint{.20, -.29}, portraitPoint{.35, -.11}, portraitPoint{.36, .20}, portraitPoint{.27, .24}, portraitPoint{.12, .13}, portraitPoint{-.23, .18}, portraitPoint{-.39, .26}, portraitPoint{-.33, -.08})
	poly(239, portraitPoint{-.22, -.28}, portraitPoint{.17, -.29}, portraitPoint{.33, -.16}, portraitPoint{.23, .12}, portraitPoint{-.20, .16}, portraitPoint{-.37, -.10})
	poly(246, portraitPoint{-.20, -.26}, portraitPoint{.07, -.28}, portraitPoint{.12, -.06}, portraitPoint{-.05, .12}, portraitPoint{-.19, .13}, portraitPoint{-.29, -.08})
	poly(109, portraitPoint{.07, -.28}, portraitPoint{.25, -.20}, portraitPoint{.23, .09}, portraitPoint{-.05, .12}, portraitPoint{.12, -.06})
	poly(239, portraitPoint{-.25, -.56}, portraitPoint{.04, -.62}, portraitPoint{.25, -.52}, portraitPoint{.28, -.33}, portraitPoint{.17, -.24}, portraitPoint{-.21, -.25}, portraitPoint{-.31, -.38})
	poly(246, portraitPoint{-.25, -.53}, portraitPoint{.04, -.58}, portraitPoint{.22, -.50}, portraitPoint{.20, -.43}, portraitPoint{-.27, -.42})
	poly(144, portraitPoint{-.23, -.41}, portraitPoint{.20, -.41}, portraitPoint{.17, -.22}, portraitPoint{.07, -.15}, portraitPoint{-.14, -.18}, portraitPoint{-.22, -.25})
	poly(230, portraitPoint{-.19, -.40}, portraitPoint{.10, -.40}, portraitPoint{.12, -.32}, portraitPoint{-.02, -.24}, portraitPoint{-.17, -.27})
	stroke(234, 1.1, portraitPoint{-.18, -.35}, portraitPoint{.12, -.35})
	stroke(151, .7, portraitPoint{-.15, -.35}, portraitPoint{-.06, -.35})
	stroke(151, .7, portraitPoint{.04, -.35}, portraitPoint{.10, -.35})
	poly(236, portraitPoint{-.04, -.26}, portraitPoint{.03, -.23}, portraitPoint{-.06, -.22})
	stroke(108, .65, portraitPoint{-.13, -.20}, portraitPoint{.09, -.19})
	// A planted sword and a large pointed shield make these clearly squires.
	stroke(239, 1.5, portraitPoint{.57, -.23}, portraitPoint{.65, .50})
	stroke(246, .7, portraitPoint{.565, -.23}, portraitPoint{.635, .49})
	stroke(137, 1.3, portraitPoint{.46, -.16}, portraitPoint{.71, -.18})
	stroke(94, 1.5, portraitPoint{.57, -.30}, portraitPoint{.57, -.19})
	poly(144, portraitPoint{.20, -.18}, portraitPoint{.31, -.12}, portraitPoint{.43, -.05}, portraitPoint{.55, -.21}, portraitPoint{.65, -.20}, portraitPoint{.54, .01}, portraitPoint{.39, .08}, portraitPoint{.22, -.05})
	poly(230, portraitPoint{.41, -.04}, portraitPoint{.55, -.20}, portraitPoint{.60, -.20}, portraitPoint{.51, -.02}, portraitPoint{.40, .04})
	poly(109, portraitPoint{-.62, -.15}, portraitPoint{-.23, -.20}, portraitPoint{-.09, -.09}, portraitPoint{-.14, .21}, portraitPoint{-.38, .37}, portraitPoint{-.65, .15})
	poly(95, portraitPoint{-.56, -.10}, portraitPoint{-.25, -.14}, portraitPoint{-.16, -.06}, portraitPoint{-.20, .17}, portraitPoint{-.37, .28}, portraitPoint{-.58, .13})
	poly(131, portraitPoint{-.54, -.09}, portraitPoint{-.38, -.11}, portraitPoint{-.37, .28}, portraitPoint{-.58, .13})
	stroke(144, .8, portraitPoint{-.38, -.08}, portraitPoint{-.36, .19})
	stroke(144, .8, portraitPoint{-.48, .02}, portraitPoint{-.25, .01})
}
