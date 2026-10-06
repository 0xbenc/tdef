package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderCenturionMockup is a static lore study of the guild fighter.
func RenderCenturionMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.centurion.enlarge_to_view_the_centurion"), 180, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawCenturionPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-9)/2, 1, copytext.Text("characters.centurion.centurion"), 180, 233, true)
		subtitle := copytext.Text("characters.centurion.hold_the_line")
		putString(f, (w-len(subtitle))/2, 2, subtitle, 240, 233, false)
	}
	hint := copytext.Format("characters.centurion.tab_squire_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.centurion.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawCenturionPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly := p.poly
	poly(236, portraitPoint{.16, .948}, portraitPoint{.83, .948}, portraitPoint{.88, .976}, portraitPoint{.12, .976})
	// A broken, architectural ward surrounds the shield, with deliberate air gaps.
	p.stroke(24, 2.3, portraitPoint{.397, .858}, portraitPoint{.352, .824}, portraitPoint{.352, .403}, portraitPoint{.432, .323}, portraitPoint{.728, .323}, portraitPoint{.826, .413}, portraitPoint{.826, .64})
	p.stroke(24, 2.3, portraitPoint{.826, .785}, portraitPoint{.826, .822}, portraitPoint{.767, .862})
	p.stroke(117, 1.0, portraitPoint{.361, .759}, portraitPoint{.361, .403}, portraitPoint{.431, .333}, portraitPoint{.724, .333}, portraitPoint{.815, .417}, portraitPoint{.815, .64})
	p.stroke(117, 1.0, portraitPoint{.816, .785}, portraitPoint{.816, .817}, portraitPoint{.764, .85})
	if h >= 24 {
		// Center the ward node on a tile so its diamond has matching tips.
		u := float64(int(.825*float64(w-1)+.5)) / float64(w-1)
		v := float64(int(.713*float64(h-1)+.5)) / float64(h-1)
		poly(153, portraitPoint{u, v - .038}, portraitPoint{u + .025, v}, portraitPoint{u, v + .038}, portraitPoint{u - .025, v})
	} else {
		poly(153, portraitPoint{.825, .688}, portraitPoint{.844, .713}, portraitPoint{.825, .738}, portraitPoint{.805, .713})
	}
	// A heavy mantle falls straight instead of flying behind him.
	poly(52, portraitPoint{.364, .301}, portraitPoint{.574, .317}, portraitPoint{.607, .664}, portraitPoint{.522, .869}, portraitPoint{.419, .85}, portraitPoint{.302, .904}, portraitPoint{.248, .872}, portraitPoint{.274, .604}, portraitPoint{.265, .395})
	poly(95, portraitPoint{.329, .36}, portraitPoint{.366, .359}, portraitPoint{.371, .647}, portraitPoint{.321, .862}, portraitPoint{.281, .87}, portraitPoint{.3, .615})
	poly(131, portraitPoint{.291, .387}, portraitPoint{.313, .374}, portraitPoint{.304, .607}, portraitPoint{.284, .821}, portraitPoint{.274, .838}, portraitPoint{.281, .598})
	// Braced legs, layered greaves and two broad sandals give the mass a footing.
	poly(236, portraitPoint{.363, .585}, portraitPoint{.469, .603}, portraitPoint{.445, .76}, portraitPoint{.429, .911}, portraitPoint{.34, .932}, portraitPoint{.32, .879}, portraitPoint{.341, .727})
	poly(239, portraitPoint{.356, .725}, portraitPoint{.438, .727}, portraitPoint{.422, .891}, portraitPoint{.341, .912}, portraitPoint{.336, .875})
	poly(109, portraitPoint{.372, .749}, portraitPoint{.422, .751}, portraitPoint{.4, .881}, portraitPoint{.351, .893})
	poly(246, portraitPoint{.372, .749}, portraitPoint{.391, .749}, portraitPoint{.373, .886}, portraitPoint{.351, .893})
	poly(137, portraitPoint{.343, .818}, portraitPoint{.423, .818}, portraitPoint{.42, .842}, portraitPoint{.34, .842})
	poly(236, portraitPoint{.531, .626}, portraitPoint{.629, .635}, portraitPoint{.664, .799}, portraitPoint{.675, .922}, portraitPoint{.581, .93}, portraitPoint{.565, .812})
	poly(109, portraitPoint{.59, .78}, portraitPoint{.64, .778}, portraitPoint{.658, .896}, portraitPoint{.599, .911})
	poly(137, portraitPoint{.584, .839}, portraitPoint{.65, .833}, portraitPoint{.654, .856}, portraitPoint{.588, .862})
	poly(94, portraitPoint{.343, .894}, portraitPoint{.418, .897}, portraitPoint{.451, .956}, portraitPoint{.284, .956}, portraitPoint{.276, .939})
	poly(137, portraitPoint{.329, .924}, portraitPoint{.402, .921}, portraitPoint{.416, .943}, portraitPoint{.29, .943})
	poly(94, portraitPoint{.592, .897}, portraitPoint{.668, .899}, portraitPoint{.722, .948}, portraitPoint{.712, .962}, portraitPoint{.573, .962})
	poly(137, portraitPoint{.6, .921}, portraitPoint{.65, .917}, portraitPoint{.695, .943}, portraitPoint{.584, .944})
	// Segmented iron cuirass has a square waist and broad, low shoulders.
	poly(239, portraitPoint{.326, .349}, portraitPoint{.422, .307}, portraitPoint{.565, .31}, portraitPoint{.65, .368}, portraitPoint{.622, .579}, portraitPoint{.585, .686}, portraitPoint{.372, .673}, portraitPoint{.32, .551})
	poly(109, portraitPoint{.351, .366}, portraitPoint{.432, .329}, portraitPoint{.556, .333}, portraitPoint{.61, .377}, portraitPoint{.585, .543}, portraitPoint{.375, .553})
	poly(246, portraitPoint{.367, .366}, portraitPoint{.444, .332}, portraitPoint{.459, .4}, portraitPoint{.402, .509}, portraitPoint{.375, .53})
	poly(239, portraitPoint{.467, .332}, portraitPoint{.552, .337}, portraitPoint{.594, .38}, portraitPoint{.571, .53}, portraitPoint{.455, .53})
	for _, v := range []float64{.41, .46, .51} {
		poly(239, portraitPoint{.349, v}, portraitPoint{.603, v + .007}, portraitPoint{.598, v + .024}, portraitPoint{.353, v + .023})
		poly(251, portraitPoint{.352, v}, portraitPoint{.453, v + .003}, portraitPoint{.449, v + .01}, portraitPoint{.353, v + .009})
	}
	poly(94, portraitPoint{.367, .546}, portraitPoint{.601, .546}, portraitPoint{.598, .587}, portraitPoint{.364, .587})
	poly(180, portraitPoint{.423, .549}, portraitPoint{.466, .549}, portraitPoint{.466, .58}, portraitPoint{.423, .58})
	// Leather strips of the military skirt are broad enough to stay distinct.
	for i := 0; i < 5; i++ {
		u := .365 + float64(i)*.041
		poly(94, portraitPoint{u, .581}, portraitPoint{u + .035, .581}, portraitPoint{u + .041, .689}, portraitPoint{u + .023, .708}, portraitPoint{u - .003, .687})
		poly(137, portraitPoint{u + .004, .587}, portraitPoint{u + .019, .587}, portraitPoint{u + .025, .684}, portraitPoint{u + .012, .684})
	}
	// Sword-side shoulder and arm remain exposed beside the immense shield.
	poly(239, portraitPoint{.302, .356}, portraitPoint{.367, .338}, portraitPoint{.415, .377}, portraitPoint{.397, .431}, portraitPoint{.303, .442}, portraitPoint{.274, .398})
	poly(246, portraitPoint{.305, .368}, portraitPoint{.36, .354}, portraitPoint{.387, .379}, portraitPoint{.379, .405}, portraitPoint{.302, .411}, portraitPoint{.288, .392})
	poly(137, portraitPoint{.288, .43}, portraitPoint{.359, .428}, portraitPoint{.355, .536}, portraitPoint{.312, .589}, portraitPoint{.264, .55})
	poly(239, portraitPoint{.281, .493}, portraitPoint{.35, .5}, portraitPoint{.332, .573}, portraitPoint{.285, .581}, portraitPoint{.259, .552})
	poly(109, portraitPoint{.285, .507}, portraitPoint{.329, .511}, portraitPoint{.317, .555}, portraitPoint{.283, .559})
	// A compact gladius points down and out, separate from the stance.
	poly(239, portraitPoint{.253, .604}, portraitPoint{.299, .622}, portraitPoint{.23, .842}, portraitPoint{.185, .89}, portraitPoint{.19, .821})
	poly(251, portraitPoint{.253, .604}, portraitPoint{.271, .612}, portraitPoint{.21, .828}, portraitPoint{.185, .89}, portraitPoint{.19, .821})
	poly(109, portraitPoint{.271, .612}, portraitPoint{.299, .622}, portraitPoint{.23, .842}, portraitPoint{.185, .89}, portraitPoint{.21, .828})
	p.stroke(94, 2.8, portraitPoint{.293, .554}, portraitPoint{.27, .623})
	p.stroke(180, 1.7, portraitPoint{.231, .606}, portraitPoint{.314, .634})
	poly(137, portraitPoint{.278, .548}, portraitPoint{.305, .542}, portraitPoint{.322, .564}, portraitPoint{.306, .603}, portraitPoint{.276, .606}, portraitPoint{.26, .58})
	poly(223, portraitPoint{.278, .557}, portraitPoint{.299, .554}, portraitPoint{.304, .571}, portraitPoint{.291, .591}, portraitPoint{.273, .588})
	// The scutum is a flat-topped wall with bevels, not a pointed knight's shield.
	poly(137, portraitPoint{.442, .36}, portraitPoint{.721, .357}, portraitPoint{.774, .393}, portraitPoint{.777, .843}, portraitPoint{.736, .891}, portraitPoint{.443, .89}, portraitPoint{.413, .849}, portraitPoint{.415, .399})
	poly(180, portraitPoint{.451, .373}, portraitPoint{.715, .371}, portraitPoint{.757, .404}, portraitPoint{.757, .83}, portraitPoint{.725, .869}, portraitPoint{.451, .869}, portraitPoint{.43, .841}, portraitPoint{.433, .408})
	poly(52, portraitPoint{.459, .388}, portraitPoint{.707, .388}, portraitPoint{.739, .414}, portraitPoint{.739, .822}, portraitPoint{.715, .851}, portraitPoint{.462, .851}, portraitPoint{.449, .831}, portraitPoint{.449, .419})
	poly(95, portraitPoint{.466, .4}, portraitPoint{.571, .394}, portraitPoint{.567, .841}, portraitPoint{.468, .841}, portraitPoint{.461, .823})
	poly(131, portraitPoint{.468, .405}, portraitPoint{.495, .403}, portraitPoint{.492, .834}, portraitPoint{.472, .832})
	// A bronze spine and square boss make the shield read as forged architecture.
	poly(137, portraitPoint{.575, .386}, portraitPoint{.609, .386}, portraitPoint{.609, .851}, portraitPoint{.575, .851})
	poly(180, portraitPoint{.583, .387}, portraitPoint{.594, .387}, portraitPoint{.594, .851}, portraitPoint{.583, .851})
	poly(137, portraitPoint{.459, .589}, portraitPoint{.737, .589}, portraitPoint{.737, .618}, portraitPoint{.459, .618})
	poly(239, portraitPoint{.557, .55}, portraitPoint{.624, .55}, portraitPoint{.641, .578}, portraitPoint{.641, .631}, portraitPoint{.621, .656}, portraitPoint{.558, .656}, portraitPoint{.541, .631}, portraitPoint{.541, .579})
	poly(180, portraitPoint{.561, .565}, portraitPoint{.619, .565}, portraitPoint{.626, .588}, portraitPoint{.626, .624}, portraitPoint{.612, .642}, portraitPoint{.562, .642}, portraitPoint{.554, .622}, portraitPoint{.554, .587})
	poly(137, portraitPoint{.587, .565}, portraitPoint{.619, .565}, portraitPoint{.626, .588}, portraitPoint{.626, .624}, portraitPoint{.612, .642}, portraitPoint{.587, .642})
	// Narrow human face, framed by the brow and two long cheek plates.
	poly(94, portraitPoint{.425, .231}, portraitPoint{.57, .232}, portraitPoint{.559, .322}, portraitPoint{.5, .353}, portraitPoint{.441, .326})
	poly(137, portraitPoint{.444, .243}, portraitPoint{.549, .243}, portraitPoint{.542, .305}, portraitPoint{.509, .329}, portraitPoint{.461, .308})
	poly(223, portraitPoint{.458, .247}, portraitPoint{.507, .245}, portraitPoint{.506, .279}, portraitPoint{.526, .291}, portraitPoint{.51, .302}, portraitPoint{.472, .29})
	poly(239, portraitPoint{.404, .242}, portraitPoint{.436, .231}, portraitPoint{.45, .28}, portraitPoint{.469, .307}, portraitPoint{.446, .332}, portraitPoint{.415, .301})
	poly(109, portraitPoint{.556, .232}, portraitPoint{.591, .245}, portraitPoint{.583, .305}, portraitPoint{.548, .329}, portraitPoint{.531, .31}, portraitPoint{.549, .282})
	// The helmet dome supports a broad horsehair crest, deliberately transverse.
	poly(239, portraitPoint{.395, .242}, portraitPoint{.406, .179}, portraitPoint{.447, .14}, portraitPoint{.525, .132}, portraitPoint{.576, .171}, portraitPoint{.604, .241}, portraitPoint{.57, .261}, portraitPoint{.43, .26})
	poly(246, portraitPoint{.414, .226}, portraitPoint{.422, .182}, portraitPoint{.454, .157}, portraitPoint{.496, .15}, portraitPoint{.496, .223})
	poly(109, portraitPoint{.496, .15}, portraitPoint{.526, .149}, portraitPoint{.56, .182}, portraitPoint{.582, .23}, portraitPoint{.496, .223})
	poly(180, portraitPoint{.401, .225}, portraitPoint{.587, .226}, portraitPoint{.596, .246}, portraitPoint{.401, .247})
	poly(137, portraitPoint{.46, .136}, portraitPoint{.46, .089}, portraitPoint{.532, .089}, portraitPoint{.535, .141})
	poly(52, portraitPoint{.304, .137}, portraitPoint{.332, .082}, portraitPoint{.402, .039}, portraitPoint{.505, .029}, portraitPoint{.6, .048}, portraitPoint{.668, .095}, portraitPoint{.682, .146}, portraitPoint{.627, .16}, portraitPoint{.579, .123}, portraitPoint{.51, .101}, portraitPoint{.432, .109}, portraitPoint{.363, .148})
	poly(131, portraitPoint{.323, .13}, portraitPoint{.349, .088}, portraitPoint{.41, .054}, portraitPoint{.503, .045}, portraitPoint{.588, .063}, portraitPoint{.647, .106}, portraitPoint{.655, .135}, portraitPoint{.628, .141}, portraitPoint{.582, .109}, portraitPoint{.514, .088}, portraitPoint{.432, .095}, portraitPoint{.362, .134})
	poly(167, portraitPoint{.349, .088}, portraitPoint{.41, .054}, portraitPoint{.503, .045}, portraitPoint{.504, .06}, portraitPoint{.414, .072}, portraitPoint{.362, .108}, portraitPoint{.347, .126}, portraitPoint{.332, .13})
	// Facial marks stay on the face even when the helmet occupies only two rows.
	if h < 20 {
		// At this scale the brow and eyes share a row: retain one clear opening.
		x, y := x0+int(.5*float64(w-1)+.5), y0+int(.3*float64(h-1)+.5)
		f.Set(x, y, Cell{R: '━', FG: 234, BG: 223})
	} else {
		for _, u := range []float64{.475, .531} {
			x, y := x0+int(u*float64(w-1)+.5), y0+int(.269*float64(h-1)+.5)
			c := f.C[y*f.W+x]
			if c.R == '█' && (c.FG == 223 || c.FG == 137) {
				f.Set(x, y, Cell{R: '━', FG: 234, BG: c.FG})
			}
		}
		p.detail(.503, .319, '─', 239)
	}
	p.detail(.292, .583, '│', 137)
	p.detail(.59, .6, '◆', 223)
	if h >= 24 {
		for _, u := range []float64{.46, .726} {
			for _, v := range []float64{.419, .815} {
				p.detail(u, v, '●', 180)
			}
		}
		p.detail(.398, .812, '─', 180)
		p.detail(.519, .071, '│', 52)
		p.detail(.575, .084, '╲', 52)
		p.detail(.412, .077, '╱', 52)
		p.detail(.478, .756, '╱', 137)
	}
}
