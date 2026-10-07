package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderPaladinMockup is the guild's immovable vanguard: a closed helm above
// a monumental pointed shield, with a flanged mace held upright at its side.
func RenderPaladinMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.paladin.enlarge_to_view_the_paladin"), 230, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawPaladinPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		title, subtitle := copytext.Text("characters.paladin.paladin"), copytext.Text("characters.paladin.the_line_does_not_yield")
		putString(f, (w-len(title))/2, 1, title, 230, 233, true)
		putString(f, (w-len(subtitle))/2, 2, subtitle, 252, 233, false)
	}
	hint := copytext.Format("characters.paladin.tab_rogue_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.paladin.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 252, 233, false)
	return f
}

func drawPaladinPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.smoothOval
	poly(236, portraitPoint{.16, .945}, portraitPoint{.80, .945}, portraitPoint{.84, .97}, portraitPoint{.12, .97})
	// A heavy mantle falls straight behind the armor, with one folded edge.
	poly(52, portraitPoint{.44, .315}, portraitPoint{.64, .325}, portraitPoint{.72, .44}, portraitPoint{.79, .875}, portraitPoint{.735, .917}, portraitPoint{.665, .896}, portraitPoint{.63, .61}, portraitPoint{.37, .47})
	poly(95, portraitPoint{.655, .36}, portraitPoint{.69, .415}, portraitPoint{.733, .67}, portraitPoint{.79, .875}, portraitPoint{.753, .906}, portraitPoint{.70, .744})
	poly(144, portraitPoint{.655, .36}, portraitPoint{.67, .38}, portraitPoint{.712, .65}, portraitPoint{.753, .89}, portraitPoint{.737, .901}, portraitPoint{.693, .664})
	// Wide planted legs and slab-like feet resist the whole weight of the shield.
	poly(239, portraitPoint{.37, .585}, portraitPoint{.49, .59}, portraitPoint{.48, .735}, portraitPoint{.435, .885}, portraitPoint{.34, .915}, portraitPoint{.335, .805})
	poly(109, portraitPoint{.405, .65}, portraitPoint{.464, .661}, portraitPoint{.445, .758}, portraitPoint{.417, .869}, portraitPoint{.37, .885}, portraitPoint{.361, .808})
	poly(246, portraitPoint{.398, .688}, portraitPoint{.43, .69}, portraitPoint{.406, .802}, portraitPoint{.386, .868}, portraitPoint{.371, .868})
	poly(239, portraitPoint{.55, .58}, portraitPoint{.655, .58}, portraitPoint{.684, .716}, portraitPoint{.681, .878}, portraitPoint{.731, .916}, portraitPoint{.607, .923}, portraitPoint{.572, .784})
	poly(109, portraitPoint{.581, .647}, portraitPoint{.632, .635}, portraitPoint{.658, .727}, portraitPoint{.65, .862}, portraitPoint{.62, .895}, portraitPoint{.595, .79})
	poly(246, portraitPoint{.582, .66}, portraitPoint{.606, .651}, portraitPoint{.633, .73}, portraitPoint{.63, .867}, portraitPoint{.617, .877}, portraitPoint{.604, .784})
	// An angular knee boss and overlapping greaves expose the right leg's mass.
	poly(239, portraitPoint{.58, .706}, portraitPoint{.656, .699}, portraitPoint{.693, .738}, portraitPoint{.665, .783}, portraitPoint{.593, .779}, portraitPoint{.562, .743})
	poly(251, portraitPoint{.589, .717}, portraitPoint{.638, .713}, portraitPoint{.665, .74}, portraitPoint{.643, .763}, portraitPoint{.595, .762}, portraitPoint{.578, .742})
	poly(230, portraitPoint{.592, .72}, portraitPoint{.617, .717}, portraitPoint{.635, .74}, portraitPoint{.61, .76}, portraitPoint{.595, .76}, portraitPoint{.58, .742})
	poly(239, portraitPoint{.34, .866}, portraitPoint{.438, .877}, portraitPoint{.454, .917}, portraitPoint{.433, .947}, portraitPoint{.269, .947}, portraitPoint{.251, .929})
	poly(246, portraitPoint{.343, .885}, portraitPoint{.421, .89}, portraitPoint{.436, .917}, portraitPoint{.416, .934}, portraitPoint{.27, .934}, portraitPoint{.28, .925})
	poly(230, portraitPoint{.332, .9}, portraitPoint{.379, .905}, portraitPoint{.367, .925}, portraitPoint{.275, .932}, portraitPoint{.28, .925})
	poly(239, portraitPoint{.607, .874}, portraitPoint{.668, .87}, portraitPoint{.73, .913}, portraitPoint{.77, .931}, portraitPoint{.755, .948}, portraitPoint{.613, .948}, portraitPoint{.588, .928})
	poly(246, portraitPoint{.62, .887}, portraitPoint{.657, .884}, portraitPoint{.72, .919}, portraitPoint{.749, .935}, portraitPoint{.62, .935}, portraitPoint{.605, .92})
	poly(230, portraitPoint{.62, .889}, portraitPoint{.642, .888}, portraitPoint{.674, .917}, portraitPoint{.644, .931}, portraitPoint{.62, .931}, portraitPoint{.609, .918})
	// The cuirass forms a broad rigid barrel, layered above the hips.
	poly(239, portraitPoint{.37, .315}, portraitPoint{.586, .311}, portraitPoint{.675, .389}, portraitPoint{.655, .525}, portraitPoint{.629, .619}, portraitPoint{.383, .643}, portraitPoint{.317, .511}, portraitPoint{.323, .396})
	poly(109, portraitPoint{.385, .336}, portraitPoint{.574, .332}, portraitPoint{.629, .397}, portraitPoint{.615, .509}, portraitPoint{.579, .586}, portraitPoint{.4, .611}, portraitPoint{.353, .502}, portraitPoint{.35, .406})
	poly(251, portraitPoint{.384, .341}, portraitPoint{.492, .338}, portraitPoint{.513, .451}, portraitPoint{.49, .568}, portraitPoint{.401, .597}, portraitPoint{.359, .499}, portraitPoint{.359, .409})
	poly(246, portraitPoint{.492, .338}, portraitPoint{.569, .337}, portraitPoint{.621, .398}, portraitPoint{.611, .478}, portraitPoint{.573, .541}, portraitPoint{.49, .568}, portraitPoint{.513, .451})
	poly(180, portraitPoint{.385, .329}, portraitPoint{.575, .326}, portraitPoint{.613, .355}, portraitPoint{.586, .36}, portraitPoint{.566, .343}, portraitPoint{.398, .349})
	poly(239, portraitPoint{.381, .527}, portraitPoint{.624, .506}, portraitPoint{.657, .577}, portraitPoint{.633, .62}, portraitPoint{.375, .643}, portraitPoint{.35, .59})
	poly(246, portraitPoint{.385, .549}, portraitPoint{.62, .53}, portraitPoint{.637, .567}, portraitPoint{.377, .594}, portraitPoint{.366, .576})
	poly(109, portraitPoint{.377, .604}, portraitPoint{.637, .579}, portraitPoint{.639, .601}, portraitPoint{.373, .626})
	// An ivory tabard hangs vertically through the exposed space beside the shield.
	poly(144, portraitPoint{.467, .367}, portraitPoint{.592, .365}, portraitPoint{.628, .426}, portraitPoint{.602, .559}, portraitPoint{.63, .76}, portraitPoint{.616, .833}, portraitPoint{.56, .819}, portraitPoint{.502, .841}, portraitPoint{.479, .783}, portraitPoint{.498, .609}, portraitPoint{.459, .479})
	poly(223, portraitPoint{.47, .377}, portraitPoint{.555, .369}, portraitPoint{.592, .414}, portraitPoint{.565, .554}, portraitPoint{.589, .761}, portraitPoint{.574, .823}, portraitPoint{.529, .81}, portraitPoint{.503, .831}, portraitPoint{.487, .781}, portraitPoint{.51, .599}, portraitPoint{.473, .475})
	poly(230, portraitPoint{.47, .381}, portraitPoint{.508, .377}, portraitPoint{.543, .415}, portraitPoint{.521, .571}, portraitPoint{.545, .765}, portraitPoint{.53, .802}, portraitPoint{.507, .82}, portraitPoint{.493, .78}, portraitPoint{.516, .594}, portraitPoint{.48, .475})
	// The guild seal is a single dark lozenge above a descending stroke.
	poly(236, portraitPoint{.558, .555}, portraitPoint{.589, .598}, portraitPoint{.562, .648}, portraitPoint{.536, .607})
	poly(236, portraitPoint{.552, .634}, portraitPoint{.568, .638}, portraitPoint{.572, .741}, portraitPoint{.557, .763})
	poly(94, portraitPoint{.554, .528}, portraitPoint{.616, .528}, portraitPoint{.616, .559}, portraitPoint{.553, .559})
	poly(180, portraitPoint{.575, .526}, portraitPoint{.61, .526}, portraitPoint{.609, .563}, portraitPoint{.575, .564})
	poly(239, portraitPoint{.584, .536}, portraitPoint{.601, .536}, portraitPoint{.601, .554}, portraitPoint{.584, .554})
	// Massive shoulders have overlapping lower lames, rather than round balloons.
	poly(239, portraitPoint{.303, .335}, portraitPoint{.367, .313}, portraitPoint{.434, .347}, portraitPoint{.432, .409}, portraitPoint{.405, .461}, portraitPoint{.31, .478}, portraitPoint{.263, .431}, portraitPoint{.264, .382})
	poly(246, portraitPoint{.308, .35}, portraitPoint{.36, .329}, portraitPoint{.414, .356}, portraitPoint{.41, .403}, portraitPoint{.377, .427}, portraitPoint{.289, .442}, portraitPoint{.278, .4})
	poly(230, portraitPoint{.308, .351}, portraitPoint{.36, .334}, portraitPoint{.386, .349}, portraitPoint{.373, .38}, portraitPoint{.294, .418}, portraitPoint{.28, .404})
	poly(180, portraitPoint{.28, .435}, portraitPoint{.414, .404}, portraitPoint{.42, .428}, portraitPoint{.3, .468}, portraitPoint{.28, .455})
	poly(239, portraitPoint{.586, .314}, portraitPoint{.653, .319}, portraitPoint{.726, .37}, portraitPoint{.744, .421}, portraitPoint{.721, .48}, portraitPoint{.638, .47}, portraitPoint{.582, .428}, portraitPoint{.562, .365})
	poly(246, portraitPoint{.59, .332}, portraitPoint{.647, .335}, portraitPoint{.711, .38}, portraitPoint{.723, .418}, portraitPoint{.695, .443}, portraitPoint{.638, .425}, portraitPoint{.586, .394}, portraitPoint{.578, .36})
	poly(230, portraitPoint{.591, .334}, portraitPoint{.626, .336}, portraitPoint{.69, .381}, portraitPoint{.675, .397}, portraitPoint{.63, .381}, portraitPoint{.584, .354})
	poly(180, portraitPoint{.6, .414}, portraitPoint{.695, .448}, portraitPoint{.731, .427}, portraitPoint{.727, .451}, portraitPoint{.697, .474}, portraitPoint{.629, .452})
	poly(109, portraitPoint{.621, .454}, portraitPoint{.696, .479}, portraitPoint{.72, .465}, portraitPoint{.717, .49}, portraitPoint{.689, .508}, portraitPoint{.643, .487})
	// The bent mace arm holds its elbow close; the gauntlet meets the shaft.
	poly(239, portraitPoint{.683, .447}, portraitPoint{.721, .451}, portraitPoint{.745, .504}, portraitPoint{.724, .568}, portraitPoint{.696, .587}, portraitPoint{.645, .556}, portraitPoint{.658, .518})
	poly(109, portraitPoint{.692, .464}, portraitPoint{.715, .468}, portraitPoint{.724, .507}, portraitPoint{.704, .551}, portraitPoint{.679, .55}, portraitPoint{.674, .525})
	poly(246, portraitPoint{.692, .472}, portraitPoint{.705, .475}, portraitPoint{.709, .514}, portraitPoint{.69, .541}, portraitPoint{.68, .537})
	poly(239, portraitPoint{.70, .527}, portraitPoint{.755, .502}, portraitPoint{.807, .505}, portraitPoint{.818, .55}, portraitPoint{.766, .574}, portraitPoint{.697, .58}, portraitPoint{.674, .562})
	poly(251, portraitPoint{.708, .537}, portraitPoint{.758, .517}, portraitPoint{.797, .516}, portraitPoint{.806, .536}, portraitPoint{.763, .551}, portraitPoint{.703, .561}, portraitPoint{.689, .55})
	poly(109, portraitPoint{.707, .562}, portraitPoint{.764, .555}, portraitPoint{.814, .537}, portraitPoint{.812, .555}, portraitPoint{.768, .574}, portraitPoint{.707, .58})
	poly(180, portraitPoint{.746, .513}, portraitPoint{.761, .507}, portraitPoint{.778, .562}, portraitPoint{.76, .57})
	// A flanged mace reads as a heavy crown of steel above a straight handle.
	p.stroke(94, 2.3, portraitPoint{.794, .245}, portraitPoint{.794, .83})
	p.stroke(180, 1.0, portraitPoint{.787, .247}, portraitPoint{.787, .83})
	p.stroke(239, 2.8, portraitPoint{.794, .482}, portraitPoint{.794, .691})
	for i := 0; i < 5; i++ {
		v := .50 + float64(i)*.038
		p.stroke(137, .9, portraitPoint{.78, v}, portraitPoint{.809, v + .018})
	}
	oval(137, .794, .834, .023, .022)
	oval(180, .787, .829, .012, .013)
	poly(239, portraitPoint{.777, .07}, portraitPoint{.803, .07}, portraitPoint{.816, .12}, portraitPoint{.811, .231}, portraitPoint{.799, .26}, portraitPoint{.777, .247}, portraitPoint{.762, .128})
	poly(246, portraitPoint{.778, .08}, portraitPoint{.791, .079}, portraitPoint{.795, .221}, portraitPoint{.785, .244}, portraitPoint{.773, .227}, portraitPoint{.77, .13})
	poly(239, portraitPoint{.746, .09}, portraitPoint{.764, .097}, portraitPoint{.78, .139}, portraitPoint{.776, .232}, portraitPoint{.753, .249}, portraitPoint{.703, .21}, portraitPoint{.697, .132})
	poly(251, portraitPoint{.746, .098}, portraitPoint{.759, .11}, portraitPoint{.767, .146}, portraitPoint{.763, .22}, portraitPoint{.75, .233}, portraitPoint{.73, .2}, portraitPoint{.729, .13})
	poly(239, portraitPoint{.819, .09}, portraitPoint{.843, .10}, portraitPoint{.889, .137}, portraitPoint{.886, .214}, portraitPoint{.832, .25}, portraitPoint{.808, .221}, portraitPoint{.813, .141})
	poly(109, portraitPoint{.834, .107}, portraitPoint{.858, .143}, portraitPoint{.852, .204}, portraitPoint{.829, .233}, portraitPoint{.819, .216}, portraitPoint{.821, .143})
	poly(180, portraitPoint{.746, .084}, portraitPoint{.764, .09}, portraitPoint{.773, .13}, portraitPoint{.766, .136}, portraitPoint{.753, .111}, portraitPoint{.733, .13}, portraitPoint{.723, .139}, portraitPoint{.715, .127})
	poly(180, portraitPoint{.819, .084}, portraitPoint{.843, .094}, portraitPoint{.874, .141}, portraitPoint{.862, .143}, portraitPoint{.836, .116}, portraitPoint{.82, .135}, portraitPoint{.813, .128})
	poly(246, portraitPoint{.779, .064}, portraitPoint{.791, .045}, portraitPoint{.803, .064}, portraitPoint{.803, .089}, portraitPoint{.777, .089})
	poly(180, portraitPoint{.773, .24}, portraitPoint{.808, .24}, portraitPoint{.813, .27}, portraitPoint{.769, .27})
	// Steel fingers wrap over the handle; black joints articulate the grip.
	poly(239, portraitPoint{.779, .504}, portraitPoint{.798, .496}, portraitPoint{.825, .511}, portraitPoint{.828, .544}, portraitPoint{.804, .565}, portraitPoint{.781, .552}, portraitPoint{.769, .532})
	poly(251, portraitPoint{.783, .509}, portraitPoint{.797, .504}, portraitPoint{.816, .514}, portraitPoint{.817, .53}, portraitPoint{.796, .543}, portraitPoint{.779, .534})
	poly(230, portraitPoint{.783, .51}, portraitPoint{.797, .507}, portraitPoint{.813, .516}, portraitPoint{.799, .526}, portraitPoint{.779, .528})
	// The neck is sealed by stacked gorget plates. No exposed face softens it.
	poly(239, portraitPoint{.4, .289}, portraitPoint{.566, .286}, portraitPoint{.603, .337}, portraitPoint{.566, .386}, portraitPoint{.41, .384}, portraitPoint{.371, .342})
	poly(109, portraitPoint{.408, .301}, portraitPoint{.558, .298}, portraitPoint{.58, .337}, portraitPoint{.552, .361}, portraitPoint{.422, .361}, portraitPoint{.394, .338})
	poly(246, portraitPoint{.412, .31}, portraitPoint{.55, .307}, portraitPoint{.566, .33}, portraitPoint{.545, .341}, portraitPoint{.425, .341}, portraitPoint{.4, .331})
	poly(180, portraitPoint{.408, .357}, portraitPoint{.554, .357}, portraitPoint{.566, .371}, portraitPoint{.411, .373})
	// A ridged closed helm has a knife-thin visor and a severe pointed jaw.
	poly(239, portraitPoint{.485, .105}, portraitPoint{.548, .13}, portraitPoint{.578, .185}, portraitPoint{.579, .27}, portraitPoint{.548, .313}, portraitPoint{.486, .347}, portraitPoint{.414, .308}, portraitPoint{.388, .247}, portraitPoint{.393, .175}, portraitPoint{.434, .127})
	poly(246, portraitPoint{.485, .117}, portraitPoint{.537, .141}, portraitPoint{.562, .188}, portraitPoint{.559, .268}, portraitPoint{.538, .299}, portraitPoint{.486, .331}, portraitPoint{.429, .301}, portraitPoint{.406, .241}, portraitPoint{.41, .18}, portraitPoint{.444, .143})
	poly(251, portraitPoint{.485, .117}, portraitPoint{.487, .207}, portraitPoint{.463, .246}, portraitPoint{.487, .331}, portraitPoint{.429, .301}, portraitPoint{.406, .241}, portraitPoint{.41, .18}, portraitPoint{.444, .143})
	poly(109, portraitPoint{.487, .207}, portraitPoint{.537, .19}, portraitPoint{.562, .188}, portraitPoint{.559, .268}, portraitPoint{.538, .299}, portraitPoint{.486, .331}, portraitPoint{.463, .246})
	poly(230, portraitPoint{.478, .091}, portraitPoint{.493, .09}, portraitPoint{.507, .15}, portraitPoint{.498, .226}, portraitPoint{.484, .243}, portraitPoint{.477, .188})
	p.stroke(234, 1.45, portraitPoint{.414, .224}, portraitPoint{.473, .235}, portraitPoint{.548, .222})
	poly(251, portraitPoint{.479, .232}, portraitPoint{.494, .226}, portraitPoint{.507, .278}, portraitPoint{.486, .297}, portraitPoint{.47, .277})
	poly(230, portraitPoint{.479, .232}, portraitPoint{.487, .23}, portraitPoint{.489, .278}, portraitPoint{.482, .286}, portraitPoint{.473, .277})
	poly(239, portraitPoint{.439, .284}, portraitPoint{.476, .303}, portraitPoint{.526, .283}, portraitPoint{.514, .313}, portraitPoint{.486, .336}, portraitPoint{.454, .316})
	// The foreground shield is one convex kite, not a flat rectangular panel.
	poly(239, portraitPoint{.155, .43}, portraitPoint{.353, .364}, portraitPoint{.555, .409}, portraitPoint{.574, .662}, portraitPoint{.509, .793}, portraitPoint{.398, .922}, portraitPoint{.248, .832}, portraitPoint{.168, .681})
	poly(137, portraitPoint{.17, .44}, portraitPoint{.354, .378}, portraitPoint{.542, .423}, portraitPoint{.556, .658}, portraitPoint{.495, .784}, portraitPoint{.397, .903}, portraitPoint{.26, .821}, portraitPoint{.184, .676})
	poly(180, portraitPoint{.17, .44}, portraitPoint{.354, .378}, portraitPoint{.542, .423}, portraitPoint{.529, .443}, portraitPoint{.352, .4}, portraitPoint{.193, .457}, portraitPoint{.207, .672}, portraitPoint{.28, .811}, portraitPoint{.397, .889}, portraitPoint{.397, .903}, portraitPoint{.26, .821}, portraitPoint{.184, .676})
	poly(144, portraitPoint{.196, .462}, portraitPoint{.354, .408}, portraitPoint{.52, .448}, portraitPoint{.536, .653}, portraitPoint{.476, .778}, portraitPoint{.396, .869}, portraitPoint{.286, .8}, portraitPoint{.217, .669})
	poly(223, portraitPoint{.196, .462}, portraitPoint{.354, .408}, portraitPoint{.377, .625}, portraitPoint{.396, .869}, portraitPoint{.286, .8}, portraitPoint{.217, .669})
	poly(230, portraitPoint{.203, .466}, portraitPoint{.33, .423}, portraitPoint{.341, .622}, portraitPoint{.365, .832}, portraitPoint{.29, .791}, portraitPoint{.228, .662})
	poly(109, portraitPoint{.377, .625}, portraitPoint{.52, .448}, portraitPoint{.536, .653}, portraitPoint{.476, .778}, portraitPoint{.396, .869})
	poly(144, portraitPoint{.355, .414}, portraitPoint{.52, .448}, portraitPoint{.528, .543}, portraitPoint{.377, .625})
	// The same guild seal appears once, large and stark, across the convex face.
	poly(236, portraitPoint{.349, .484}, portraitPoint{.435, .587}, portraitPoint{.375, .696}, portraitPoint{.286, .6})
	poly(234, portraitPoint{.349, .484}, portraitPoint{.363, .595}, portraitPoint{.375, .696}, portraitPoint{.286, .6})
	poly(236, portraitPoint{.354, .668}, portraitPoint{.391, .665}, portraitPoint{.407, .786}, portraitPoint{.388, .821}, portraitPoint{.371, .793})
	// A bracing gauntlet curls over the inside rim, tying shield to bearer.
	poly(239, portraitPoint{.538, .478}, portraitPoint{.562, .48}, portraitPoint{.581, .505}, portraitPoint{.579, .54}, portraitPoint{.557, .557}, portraitPoint{.539, .542})
	poly(251, portraitPoint{.542, .484}, portraitPoint{.558, .488}, portraitPoint{.574, .505}, portraitPoint{.57, .524}, portraitPoint{.55, .532}, portraitPoint{.54, .516})
	poly(230, portraitPoint{.542, .488}, portraitPoint{.553, .49}, portraitPoint{.565, .505}, portraitPoint{.552, .516}, portraitPoint{.54, .513})
	p.detail(.564, .519, '│', 109)
	// Restrained pins, ventilation and finger seams belong to their materials.
	p.detail(.437, .269, '│', 239)
	p.detail(.457, .282, '│', 239)
	p.detail(.524, .267, '│', 239)
	p.detail(.541, .26, '│', 239)
	p.detail(.792, .521, '│', 239)
	p.detail(.808, .532, '│', 239)
	p.detail(.225, .47, '●', 137)
	p.detail(.496, .462, '●', 137)
	p.detail(.239, .665, '●', 137)
	p.detail(.482, .68, '●', 137)
	if h >= 24 {
		p.detail(.679, .405, '╲', 251)
		p.detail(.645, .745, '╱', 246)
		p.detail(.698, .924, '╲', 251)
		p.detail(.769, .169, '│', 239)
	}
}
