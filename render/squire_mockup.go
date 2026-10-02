package render

// RenderSquireMockup is a static lore study of the guild fighter.
func RenderSquireMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the Squire", 223, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawSquirePortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-6)/2, 1, "SQUIRE", 223, 233, true)
		subtitle := "the first expedition"
		putString(f, (w-len(subtitle))/2, 2, subtitle, 240, 233, false)
	}
	hint := "tab Malgrath · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawSquirePortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.smoothOval
	poly(236, portraitPoint{.18, .951}, portraitPoint{.81, .951}, portraitPoint{.86, .978}, portraitPoint{.14, .978})
	// The borrowed round shield is strapped to his back, not held in front.
	oval(94, .739, .491, .115, .16)
	oval(137, .752, .487, .096, .143)
	oval(180, .774, .473, .065, .121)
	poly(94, portraitPoint{.766, .36}, portraitPoint{.785, .369}, portraitPoint{.799, .591}, portraitPoint{.78, .615})
	oval(239, .777, .489, .024, .035)
	oval(109, .783, .494, .014, .022)
	// He plants one oversized boot ahead of the other and bends both knees.
	poly(236, portraitPoint{.455, .657}, portraitPoint{.549, .704}, portraitPoint{.486, .793}, portraitPoint{.403, .854}, portraitPoint{.361, .937}, portraitPoint{.276, .926}, portraitPoint{.318, .808}, portraitPoint{.412, .742})
	poly(60, portraitPoint{.452, .697}, portraitPoint{.508, .717}, portraitPoint{.463, .774}, portraitPoint{.381, .839}, portraitPoint{.345, .904}, portraitPoint{.303, .897}, portraitPoint{.343, .819}, portraitPoint{.428, .754})
	poly(67, portraitPoint{.447, .72}, portraitPoint{.465, .732}, portraitPoint{.438, .772}, portraitPoint{.362, .84}, portraitPoint{.337, .892}, portraitPoint{.317, .889}, portraitPoint{.351, .821})
	poly(236, portraitPoint{.537, .665}, portraitPoint{.662, .651}, portraitPoint{.684, .766}, portraitPoint{.633, .844}, portraitPoint{.711, .922}, portraitPoint{.68, .952}, portraitPoint{.587, .899}, portraitPoint{.548, .842}, portraitPoint{.584, .758})
	poly(60, portraitPoint{.569, .704}, portraitPoint{.635, .697}, portraitPoint{.654, .76}, portraitPoint{.603, .836}, portraitPoint{.657, .914}, portraitPoint{.621, .923}, portraitPoint{.575, .873}, portraitPoint{.566, .839}, portraitPoint{.609, .758})
	poly(67, portraitPoint{.609, .71}, portraitPoint{.635, .71}, portraitPoint{.638, .759}, portraitPoint{.591, .835}, portraitPoint{.614, .874}, portraitPoint{.6, .88}, portraitPoint{.58, .842}, portraitPoint{.623, .757})
	// One cheap kneecap and bindings, rather than a matched suit of plate.
	poly(239, portraitPoint{.601, .727}, portraitPoint{.651, .731}, portraitPoint{.675, .762}, portraitPoint{.636, .799}, portraitPoint{.587, .785}, portraitPoint{.579, .753})
	poly(109, portraitPoint{.609, .744}, portraitPoint{.647, .745}, portraitPoint{.657, .763}, portraitPoint{.632, .782}, portraitPoint{.6, .773}, portraitPoint{.593, .754})
	p.stroke(144, 1.1, portraitPoint{.35, .827}, portraitPoint{.395, .846})
	p.stroke(144, 1.1, portraitPoint{.331, .852}, portraitPoint{.374, .87})
	// Rolled boot tops, thick soles and an exaggerated forward toe.
	poly(94, portraitPoint{.305, .878}, portraitPoint{.367, .895}, portraitPoint{.345, .939}, portraitPoint{.374, .957}, portraitPoint{.235, .957}, portraitPoint{.245, .934})
	poly(137, portraitPoint{.312, .897}, portraitPoint{.342, .909}, portraitPoint{.325, .937}, portraitPoint{.256, .946}, portraitPoint{.26, .934})
	poly(180, portraitPoint{.303, .883}, portraitPoint{.368, .898}, portraitPoint{.361, .914}, portraitPoint{.295, .9})
	poly(94, portraitPoint{.591, .869}, portraitPoint{.652, .85}, portraitPoint{.68, .892}, portraitPoint{.72, .921}, portraitPoint{.789, .94}, portraitPoint{.778, .963}, portraitPoint{.671, .963}, portraitPoint{.615, .923})
	poly(137, portraitPoint{.608, .883}, portraitPoint{.638, .874}, portraitPoint{.66, .906}, portraitPoint{.701, .932}, portraitPoint{.755, .949}, portraitPoint{.68, .948}, portraitPoint{.63, .912})
	poly(180, portraitPoint{.591, .868}, portraitPoint{.652, .85}, portraitPoint{.66, .869}, portraitPoint{.602, .888})
	// A too-large mail shirt sags beneath a short, creased guild surcoat.
	poly(239, portraitPoint{.491, .414}, portraitPoint{.572, .393}, portraitPoint{.67, .418}, portraitPoint{.729, .485}, portraitPoint{.697, .6}, portraitPoint{.664, .714}, portraitPoint{.459, .759}, portraitPoint{.409, .69}, portraitPoint{.445, .538})
	poly(109, portraitPoint{.491, .432}, portraitPoint{.572, .412}, portraitPoint{.644, .436}, portraitPoint{.696, .493}, portraitPoint{.67, .594}, portraitPoint{.63, .697}, portraitPoint{.466, .731}, portraitPoint{.434, .683}, portraitPoint{.466, .54})
	poly(246, portraitPoint{.489, .439}, portraitPoint{.524, .427}, portraitPoint{.532, .493}, portraitPoint{.474, .636}, portraitPoint{.453, .638}, portraitPoint{.47, .535})
	poly(95, portraitPoint{.537, .409}, portraitPoint{.587, .413}, portraitPoint{.609, .457}, portraitPoint{.661, .468}, portraitPoint{.682, .443}, portraitPoint{.693, .473}, portraitPoint{.631, .666}, portraitPoint{.657, .744}, portraitPoint{.575, .767}, portraitPoint{.554, .716}, portraitPoint{.526, .778}, portraitPoint{.472, .75}, portraitPoint{.51, .657}, portraitPoint{.506, .523})
	poly(131, portraitPoint{.551, .43}, portraitPoint{.572, .437}, portraitPoint{.591, .484}, portraitPoint{.613, .487}, portraitPoint{.568, .632}, portraitPoint{.545, .66}, portraitPoint{.515, .723}, portraitPoint{.494, .739}, portraitPoint{.528, .646}, portraitPoint{.529, .521})
	poly(173, portraitPoint{.55, .453}, portraitPoint{.561, .455}, portraitPoint{.568, .508}, portraitPoint{.545, .604}, portraitPoint{.534, .601}, portraitPoint{.547, .512})
	poly(52, portraitPoint{.619, .487}, portraitPoint{.66, .489}, portraitPoint{.601, .659}, portraitPoint{.627, .736}, portraitPoint{.589, .746}, portraitPoint{.574, .698})
	// A plain stitched guild chevron, distinct from the senior fighters' shields.
	poly(223, portraitPoint{.551, .486}, portraitPoint{.574, .513}, portraitPoint{.595, .492}, portraitPoint{.586, .52}, portraitPoint{.572, .544}, portraitPoint{.551, .521})
	poly(94, portraitPoint{.477, .648}, portraitPoint{.65, .61}, portraitPoint{.657, .639}, portraitPoint{.475, .681})
	poly(180, portraitPoint{.574, .632}, portraitPoint{.607, .625}, portraitPoint{.613, .653}, portraitPoint{.58, .659})
	// Oversized shoulder cups emphasize the small body inside the equipment.
	poly(239, portraitPoint{.469, .437}, portraitPoint{.509, .408}, portraitPoint{.552, .432}, portraitPoint{.56, .48}, portraitPoint{.503, .512}, portraitPoint{.444, .482})
	poly(246, portraitPoint{.472, .446}, portraitPoint{.502, .427}, portraitPoint{.532, .443}, portraitPoint{.531, .467}, portraitPoint{.486, .487}, portraitPoint{.46, .471})
	poly(239, portraitPoint{.644, .419}, portraitPoint{.705, .433}, portraitPoint{.75, .481}, portraitPoint{.727, .517}, portraitPoint{.661, .507}, portraitPoint{.632, .463})
	poly(109, portraitPoint{.66, .439}, portraitPoint{.697, .449}, portraitPoint{.732, .483}, portraitPoint{.713, .498}, portraitPoint{.671, .489}, portraitPoint{.65, .461})
	// Both elbows pull toward the sword. Cuffs are steel; fingers remain bare.
	poly(239, portraitPoint{.458, .48}, portraitPoint{.507, .497}, portraitPoint{.474, .559}, portraitPoint{.491, .589}, portraitPoint{.46, .626}, portraitPoint{.417, .594}, portraitPoint{.415, .559})
	poly(109, portraitPoint{.453, .497}, portraitPoint{.486, .506}, portraitPoint{.45, .563}, portraitPoint{.468, .59}, portraitPoint{.451, .603}, portraitPoint{.434, .578}, portraitPoint{.434, .548})
	poly(239, portraitPoint{.698, .506}, portraitPoint{.724, .514}, portraitPoint{.697, .584}, portraitPoint{.613, .642}, portraitPoint{.554, .635}, portraitPoint{.536, .602}, portraitPoint{.568, .573}, portraitPoint{.61, .596}, portraitPoint{.662, .542})
	poly(109, portraitPoint{.696, .52}, portraitPoint{.707, .529}, portraitPoint{.681, .575}, portraitPoint{.611, .621}, portraitPoint{.569, .617}, portraitPoint{.555, .602}, portraitPoint{.576, .592}, portraitPoint{.612, .609}, portraitPoint{.671, .556})
	poly(246, portraitPoint{.61, .606}, portraitPoint{.672, .559}, portraitPoint{.68, .563}, portraitPoint{.614, .618}, portraitPoint{.581, .618}, portraitPoint{.579, .607})
	// An enormous borrowed blade tilts forward. Its long straight planes carry it.
	poly(239, portraitPoint{.152, .071}, portraitPoint{.222, .126}, portraitPoint{.508, .512}, portraitPoint{.468, .552}, portraitPoint{.183, .166})
	poly(251, portraitPoint{.152, .071}, portraitPoint{.196, .149}, portraitPoint{.485, .526}, portraitPoint{.468, .552}, portraitPoint{.183, .166})
	poly(109, portraitPoint{.152, .071}, portraitPoint{.222, .126}, portraitPoint{.508, .512}, portraitPoint{.485, .526}, portraitPoint{.196, .149})
	poly(246, portraitPoint{.174, .108}, portraitPoint{.207, .143}, portraitPoint{.49, .515}, portraitPoint{.484, .522}, portraitPoint{.198, .149})
	// A long grip genuinely joins both hands and the heavy crossguard.
	p.stroke(94, 3.0, portraitPoint{.484, .538}, portraitPoint{.565, .662})
	p.stroke(137, 1.1, portraitPoint{.477, .542}, portraitPoint{.557, .667})
	p.stroke(137, 2.4, portraitPoint{.419, .567}, portraitPoint{.54, .497})
	p.stroke(180, 1.0, portraitPoint{.42, .558}, portraitPoint{.532, .496})
	oval(239, .567, .669, .025, .023)
	oval(180, .562, .665, .014, .013)
	poly(137, portraitPoint{.465, .545}, portraitPoint{.482, .54}, portraitPoint{.508, .558}, portraitPoint{.523, .578}, portraitPoint{.51, .603}, portraitPoint{.485, .603}, portraitPoint{.456, .577})
	poly(223, portraitPoint{.472, .552}, portraitPoint{.48, .55}, portraitPoint{.5, .565}, portraitPoint{.509, .578}, portraitPoint{.499, .589}, portraitPoint{.483, .582}, portraitPoint{.465, .567})
	poly(137, portraitPoint{.521, .587}, portraitPoint{.537, .582}, portraitPoint{.56, .602}, portraitPoint{.572, .622}, portraitPoint{.558, .644}, portraitPoint{.538, .646}, portraitPoint{.511, .619})
	poly(223, portraitPoint{.526, .595}, portraitPoint{.536, .593}, portraitPoint{.553, .607}, portraitPoint{.562, .621}, portraitPoint{.553, .632}, portraitPoint{.537, .62}, portraitPoint{.519, .606})
	// A youthful face looks past the blade; a bare neck makes the head human.
	poly(137, portraitPoint{.585, .352}, portraitPoint{.653, .354}, portraitPoint{.665, .424}, portraitPoint{.624, .447}, portraitPoint{.572, .42})
	poly(223, portraitPoint{.599, .376}, portraitPoint{.64, .374}, portraitPoint{.644, .419}, portraitPoint{.622, .429}, portraitPoint{.591, .412})
	poly(94, portraitPoint{.547, .29}, portraitPoint{.632, .262}, portraitPoint{.7, .308}, portraitPoint{.708, .365}, portraitPoint{.683, .392}, portraitPoint{.651, .423}, portraitPoint{.608, .415}, portraitPoint{.574, .384}, portraitPoint{.549, .352})
	poly(137, portraitPoint{.552, .306}, portraitPoint{.634, .286}, portraitPoint{.682, .315}, portraitPoint{.672, .382}, portraitPoint{.64, .414}, portraitPoint{.608, .402}, portraitPoint{.587, .379}, portraitPoint{.546, .367}, portraitPoint{.528, .351}, portraitPoint{.553, .337})
	poly(223, portraitPoint{.562, .313}, portraitPoint{.627, .301}, portraitPoint{.653, .32}, portraitPoint{.646, .353}, portraitPoint{.638, .38}, portraitPoint{.614, .394}, portraitPoint{.588, .372}, portraitPoint{.549, .361}, portraitPoint{.541, .352}, portraitPoint{.57, .343})
	poly(180, portraitPoint{.645, .324}, portraitPoint{.67, .334}, portraitPoint{.659, .374}, portraitPoint{.641, .39}, portraitPoint{.627, .38})
	// A kettle helm tilted too low, with one visibly loose leather chin strap.
	poly(239, portraitPoint{.509, .296}, portraitPoint{.531, .23}, portraitPoint{.586, .199}, portraitPoint{.667, .216}, portraitPoint{.712, .263}, portraitPoint{.718, .307}, portraitPoint{.775, .33}, portraitPoint{.753, .346}, portraitPoint{.652, .327}, portraitPoint{.551, .327}, portraitPoint{.484, .316})
	poly(246, portraitPoint{.53, .29}, portraitPoint{.548, .239}, portraitPoint{.586, .215}, portraitPoint{.62, .219}, portraitPoint{.609, .289})
	poly(109, portraitPoint{.62, .219}, portraitPoint{.66, .232}, portraitPoint{.693, .27}, portraitPoint{.699, .307}, portraitPoint{.609, .289})
	poly(251, portraitPoint{.548, .239}, portraitPoint{.586, .215}, portraitPoint{.603, .217}, portraitPoint{.558, .248}, portraitPoint{.544, .286}, portraitPoint{.532, .29})
	poly(239, portraitPoint{.489, .302}, portraitPoint{.577, .287}, portraitPoint{.67, .299}, portraitPoint{.775, .33}, portraitPoint{.753, .34}, portraitPoint{.656, .317}, portraitPoint{.555, .316}, portraitPoint{.49, .317})
	p.stroke(94, .8, portraitPoint{.696, .333}, portraitPoint{.688, .395}, portraitPoint{.665, .425})
	// Rounded eye placement prevents the low brim stealing his expression.
	eyeU := .58
	if h < 12 {
		eyeU = .61
	}
	x, y := x0+int(eyeU*float64(w-1)+.5), y0+int(.347*float64(h-1)+.5)
	c := f.C[y*f.W+x]
	if c.R == '█' && (c.FG == 223 || c.FG == 137) {
		f.Set(x, y, Cell{R: '━', FG: 234, BG: c.FG})
	}
	if h >= 18 {
		p.detail(.606, .399, '╱', 95)
	}
	p.detail(.491, .581, '│', 137)
	p.detail(.543, .62, '│', 137)
	if h >= 24 {
		p.detail(.59, .252, '╱', 251)
		p.detail(.717, .468, '●', 137)
		p.detail(.473, .706, '∷', 246)
		p.detail(.651, .818, '╲', 180)
		p.detail(.698, .94, '─', 180)
		p.detail(.564, .698, '╱', 173)
	}
}
