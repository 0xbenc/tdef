package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderRangerMockup is one definitive lore portrait: a lean monster at full
// draw, with the eye, fingers and long arrow composing a single line of aim.
func RenderRangerMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.ranger.enlarge_to_view_the_ranger"), 151, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawRangerPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-6)/2, 1, copytext.Text("characters.ranger.ranger"), 151, 233, true)
		putString(f, (w-20)/2, 2, copytext.Text("characters.ranger.one_breath_one_shot"), 252, 233, false)
	}
	hint := copytext.Format("characters.ranger.tab_lightning_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.ranger.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 252, 233, false)
	return f
}

func drawRangerPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly := p.poly
	poly(236, portraitPoint{.13, .93}, portraitPoint{.63, .93}, portraitPoint{.70, .955}, portraitPoint{.09, .955})
	// A sideways-torn cloak balances the open space inside the drawn bow.
	poly(236, portraitPoint{.33, .285}, portraitPoint{.43, .37}, portraitPoint{.32, .46}, portraitPoint{.30, .60}, portraitPoint{.22, .73}, portraitPoint{.07, .83}, portraitPoint{.09, .72}, portraitPoint{.035, .75}, portraitPoint{.06, .62}, portraitPoint{.14, .55}, portraitPoint{.235, .39})
	poly(65, portraitPoint{.32, .30}, portraitPoint{.40, .37}, portraitPoint{.30, .455}, portraitPoint{.285, .58}, portraitPoint{.20, .72}, portraitPoint{.095, .785}, portraitPoint{.13, .69}, portraitPoint{.055, .70}, portraitPoint{.095, .61}, portraitPoint{.155, .56}, portraitPoint{.245, .40})
	poly(66, portraitPoint{.31, .35}, portraitPoint{.34, .375}, portraitPoint{.22, .575}, portraitPoint{.14, .63}, portraitPoint{.075, .69}, portraitPoint{.095, .61}, portraitPoint{.155, .56})
	poly(236, portraitPoint{.255, .43}, portraitPoint{.30, .435}, portraitPoint{.25, .61}, portraitPoint{.155, .72}, portraitPoint{.13, .75}, portraitPoint{.17, .635})
	// Soft leather boots and bent knees convey balance rather than armor.
	poly(58, portraitPoint{.33, .62}, portraitPoint{.43, .65}, portraitPoint{.34, .79}, portraitPoint{.255, .90}, portraitPoint{.19, .91}, portraitPoint{.25, .77})
	poly(65, portraitPoint{.32, .68}, portraitPoint{.355, .715}, portraitPoint{.30, .80}, portraitPoint{.235, .865}, portraitPoint{.22, .86}, portraitPoint{.275, .77})
	poly(58, portraitPoint{.43, .62}, portraitPoint{.50, .63}, portraitPoint{.54, .765}, portraitPoint{.53, .875}, portraitPoint{.57, .91}, portraitPoint{.50, .915}, portraitPoint{.465, .83}, portraitPoint{.46, .76}, portraitPoint{.40, .72})
	poly(65, portraitPoint{.465, .69}, portraitPoint{.50, .68}, portraitPoint{.53, .77}, portraitPoint{.512, .86}, portraitPoint{.49, .82})
	poly(94, portraitPoint{.245, .805}, portraitPoint{.30, .825}, portraitPoint{.275, .88}, portraitPoint{.255, .90}, portraitPoint{.27, .93}, portraitPoint{.16, .93}, portraitPoint{.145, .915}, portraitPoint{.205, .87})
	poly(137, portraitPoint{.26, .82}, portraitPoint{.287, .83}, portraitPoint{.245, .895}, portraitPoint{.19, .916}, portraitPoint{.17, .916}, portraitPoint{.225, .87})
	poly(94, portraitPoint{.49, .79}, portraitPoint{.54, .79}, portraitPoint{.54, .865}, portraitPoint{.63, .91}, portraitPoint{.64, .93}, portraitPoint{.505, .93}, portraitPoint{.48, .88})
	poly(137, portraitPoint{.505, .80}, portraitPoint{.535, .80}, portraitPoint{.53, .865}, portraitPoint{.58, .90}, portraitPoint{.52, .90}, portraitPoint{.50, .875})
	// A slanted quiver and three broad feather tips sit behind the swept ear.
	poly(94, portraitPoint{.205, .355}, portraitPoint{.27, .375}, portraitPoint{.26, .535}, portraitPoint{.21, .62}, portraitPoint{.17, .60}, portraitPoint{.165, .50})
	poly(137, portraitPoint{.22, .38}, portraitPoint{.252, .39}, portraitPoint{.24, .525}, portraitPoint{.205, .595}, portraitPoint{.19, .585}, portraitPoint{.19, .51})
	poly(180, portraitPoint{.20, .35}, portraitPoint{.27, .375}, portraitPoint{.275, .402}, portraitPoint{.20, .377})
	for i := 0; i < 3; i++ {
		u := .145 + float64(i)*.04
		p.stroke(137, .8, portraitPoint{u, .19 + float64(i)*.015}, portraitPoint{u + .055, .385})
		poly(144, portraitPoint{u - .012, .17}, portraitPoint{u + .014, .19}, portraitPoint{u + .025, .245}, portraitPoint{u - .005, .225})
		poly(109, portraitPoint{u - .012, .17}, portraitPoint{u + .006, .19}, portraitPoint{u + .014, .225}, portraitPoint{u - .005, .225})
	}
	// A long, tapered torso and diagonal strap keep the figure lean.
	poly(65, portraitPoint{.35, .345}, portraitPoint{.46, .355}, portraitPoint{.52, .42}, portraitPoint{.515, .57}, portraitPoint{.48, .67}, portraitPoint{.345, .665}, portraitPoint{.30, .555}, portraitPoint{.305, .43})
	poly(94, portraitPoint{.36, .39}, portraitPoint{.445, .385}, portraitPoint{.49, .445}, portraitPoint{.485, .57}, portraitPoint{.45, .64}, portraitPoint{.35, .62}, portraitPoint{.32, .515}, portraitPoint{.33, .435})
	poly(137, portraitPoint{.36, .405}, portraitPoint{.41, .40}, portraitPoint{.455, .455}, portraitPoint{.44, .565}, portraitPoint{.38, .595}, portraitPoint{.35, .54})
	poly(130, portraitPoint{.445, .36}, portraitPoint{.472, .38}, portraitPoint{.365, .61}, portraitPoint{.34, .60})
	poly(180, portraitPoint{.40, .485}, portraitPoint{.424, .495}, portraitPoint{.405, .53}, portraitPoint{.382, .517})
	poly(65, portraitPoint{.33, .60}, portraitPoint{.48, .605}, portraitPoint{.485, .66}, portraitPoint{.45, .695}, portraitPoint{.40, .67}, portraitPoint{.355, .69})
	poly(94, portraitPoint{.335, .59}, portraitPoint{.485, .595}, portraitPoint{.485, .625}, portraitPoint{.34, .627})
	poly(180, portraitPoint{.415, .591}, portraitPoint{.443, .595}, portraitPoint{.441, .63}, portraitPoint{.414, .625})
	// A single spare knife hangs close to the hip; it does not break the pose.
	poly(94, portraitPoint{.47, .615}, portraitPoint{.495, .615}, portraitPoint{.53, .745}, portraitPoint{.516, .775}, portraitPoint{.493, .742})
	poly(137, portraitPoint{.486, .655}, portraitPoint{.497, .655}, portraitPoint{.522, .742}, portraitPoint{.516, .76})
	// Recurved limbs use continuous half-cell strokes, with a quiet lit edge.
	bow := portraitCurve(portraitPoint{.71, .09}, portraitPoint{.82, .065}, portraitPoint{.79, .15}, portraitPoint{.745, .195})
	bow = append(bow, portraitCurve(portraitPoint{.745, .195}, portraitPoint{.75, .27}, portraitPoint{.835, .31}, portraitPoint{.81, .415})[1:]...)
	bow = append(bow, portraitCurve(portraitPoint{.81, .415}, portraitPoint{.85, .515}, portraitPoint{.765, .57}, portraitPoint{.76, .64})[1:]...)
	bow = append(bow, portraitCurve(portraitPoint{.76, .64}, portraitPoint{.76, .715}, portraitPoint{.82, .795}, portraitPoint{.72, .835})[1:]...)
	p.stroke(94, 2.3, bow...)
	p.stroke(137, 1.5, bow...)
	lit := make([]portraitPoint, len(bow))
	for i, a := range bow {
		lit[i] = portraitPoint{a.x - .006, a.y - .005}
	}
	p.stroke(180, .7, lit...)
	p.stroke(230, .8, portraitPoint{.71, .09}, portraitPoint{.397, .39}, portraitPoint{.72, .835})
	// The hood opening and blade ear frame a narrow monster profile.
	poly(65, portraitPoint{.375, .135}, portraitPoint{.465, .165}, portraitPoint{.535, .245}, portraitPoint{.545, .345}, portraitPoint{.48, .405}, portraitPoint{.375, .40}, portraitPoint{.31, .34}, portraitPoint{.305, .235})
	poly(236, portraitPoint{.375, .135}, portraitPoint{.38, .235}, portraitPoint{.335, .29}, portraitPoint{.365, .35}, portraitPoint{.42, .40}, portraitPoint{.375, .40}, portraitPoint{.31, .34}, portraitPoint{.305, .235})
	poly(107, portraitPoint{.375, .135}, portraitPoint{.465, .165}, portraitPoint{.535, .245}, portraitPoint{.525, .275}, portraitPoint{.49, .225}, portraitPoint{.425, .18})
	poly(234, portraitPoint{.39, .22}, portraitPoint{.47, .225}, portraitPoint{.52, .265}, portraitPoint{.525, .33}, portraitPoint{.47, .38}, portraitPoint{.39, .365}, portraitPoint{.345, .29})
	poly(108, portraitPoint{.19, .19}, portraitPoint{.385, .25}, portraitPoint{.415, .32}, portraitPoint{.31, .295}, portraitPoint{.25, .245})
	poly(65, portraitPoint{.23, .215}, portraitPoint{.38, .265}, portraitPoint{.39, .292}, portraitPoint{.315, .273})
	poly(108, portraitPoint{.39, .235}, portraitPoint{.465, .235}, portraitPoint{.51, .275}, portraitPoint{.525, .305}, portraitPoint{.575, .33}, portraitPoint{.563, .35}, portraitPoint{.52, .35}, portraitPoint{.49, .392}, portraitPoint{.43, .375}, portraitPoint{.38, .325}, portraitPoint{.365, .28})
	poly(151, portraitPoint{.395, .247}, portraitPoint{.455, .248}, portraitPoint{.49, .28}, portraitPoint{.495, .30}, portraitPoint{.435, .31}, portraitPoint{.385, .28})
	poly(65, portraitPoint{.385, .31}, portraitPoint{.44, .33}, portraitPoint{.505, .325}, portraitPoint{.52, .35}, portraitPoint{.49, .392}, portraitPoint{.43, .375})
	poly(108, portraitPoint{.44, .345}, portraitPoint{.505, .337}, portraitPoint{.51, .356}, portraitPoint{.475, .376}, portraitPoint{.44, .366})
	poly(151, portraitPoint{.51, .302}, portraitPoint{.565, .327}, portraitPoint{.556, .34}, portraitPoint{.515, .329})
	poly(22, portraitPoint{.452, .29}, portraitPoint{.495, .283}, portraitPoint{.51, .31}, portraitPoint{.465, .318})
	poly(65, portraitPoint{.475, .346}, portraitPoint{.524, .34}, portraitPoint{.51, .36}, portraitPoint{.48, .36})
	poly(180, portraitPoint{.505, .358}, portraitPoint{.52, .354}, portraitPoint{.528, .325})
	// Pulling elbow is carried backward; fingers anchor at the jaw, not chest.
	poly(65, portraitPoint{.34, .395}, portraitPoint{.36, .445}, portraitPoint{.295, .475}, portraitPoint{.225, .425}, portraitPoint{.225, .385}, portraitPoint{.29, .365}, portraitPoint{.385, .367}, portraitPoint{.416, .39}, portraitPoint{.39, .421}, portraitPoint{.30, .410}, portraitPoint{.27, .408})
	poly(108, portraitPoint{.235, .39}, portraitPoint{.29, .377}, portraitPoint{.378, .377}, portraitPoint{.401, .39}, portraitPoint{.385, .405}, portraitPoint{.294, .392}, portraitPoint{.257, .403}, portraitPoint{.294, .445}, portraitPoint{.287, .452}, portraitPoint{.23, .418})
	poly(151, portraitPoint{.242, .39}, portraitPoint{.29, .378}, portraitPoint{.37, .378}, portraitPoint{.39, .39}, portraitPoint{.29, .388}, portraitPoint{.247, .40})
	poly(94, portraitPoint{.338, .368}, portraitPoint{.363, .37}, portraitPoint{.354, .41}, portraitPoint{.331, .41})
	// Bow arm is fully extended, with a wrapped wrist supporting the grip.
	poly(65, portraitPoint{.475, .395}, portraitPoint{.54, .385}, portraitPoint{.59, .414}, portraitPoint{.75, .417}, portraitPoint{.785, .395}, portraitPoint{.822, .40}, portraitPoint{.836, .445}, portraitPoint{.795, .47}, portraitPoint{.755, .454}, portraitPoint{.588, .461}, portraitPoint{.525, .44})
	poly(108, portraitPoint{.53, .402}, portraitPoint{.56, .403}, portraitPoint{.593, .427}, portraitPoint{.753, .428}, portraitPoint{.786, .412}, portraitPoint{.81, .42}, portraitPoint{.80, .445}, portraitPoint{.755, .438}, portraitPoint{.595, .444}, portraitPoint{.552, .435})
	poly(151, portraitPoint{.588, .427}, portraitPoint{.753, .428}, portraitPoint{.772, .419}, portraitPoint{.78, .431}, portraitPoint{.753, .436}, portraitPoint{.592, .435})
	poly(94, portraitPoint{.755, .411}, portraitPoint{.78, .403}, portraitPoint{.802, .452}, portraitPoint{.775, .461})
	// The shaft is thin and straight, crossing the curved bow above the grip.
	p.stroke(180, 1.0, portraitPoint{.355, .388}, portraitPoint{.973, .428})
	p.stroke(223, .45, portraitPoint{.405, .389}, portraitPoint{.965, .425})
	p.stroke(244, 1.2, portraitPoint{.927, .384}, portraitPoint{.99, .429}, portraitPoint{.926, .475})
	poly(244, portraitPoint{.927, .384}, portraitPoint{.99, .429}, portraitPoint{.926, .475}, portraitPoint{.942, .428})
	poly(110, portraitPoint{.942, .428}, portraitPoint{.99, .429}, portraitPoint{.926, .475})
	poly(109, portraitPoint{.344, .374}, portraitPoint{.387, .38}, portraitPoint{.40, .395}, portraitPoint{.352, .39})
	poly(144, portraitPoint{.352, .392}, portraitPoint{.39, .395}, portraitPoint{.382, .407}, portraitPoint{.34, .408})
	// Both grips sit over the equipment. Three short finger marks tell the draw.
	poly(108, portraitPoint{.385, .371}, portraitPoint{.408, .37}, portraitPoint{.425, .394}, portraitPoint{.413, .412}, portraitPoint{.383, .411}, portraitPoint{.376, .395})
	poly(151, portraitPoint{.386, .375}, portraitPoint{.404, .377}, portraitPoint{.411, .393}, portraitPoint{.387, .398})
	poly(108, portraitPoint{.793, .399}, portraitPoint{.811, .390}, portraitPoint{.828, .404}, portraitPoint{.843, .44}, portraitPoint{.832, .46}, portraitPoint{.81, .45}, portraitPoint{.797, .425})
	poly(151, portraitPoint{.804, .402}, portraitPoint{.816, .399}, portraitPoint{.828, .424}, portraitPoint{.818, .435}, portraitPoint{.807, .421})
	p.detail(.475, .303, '━', 230)
	p.detail(.41, .315, '╲', 151)
	p.detail(.397, .39, '│', 65)
	p.detail(.411, .395, '│', 65)
	p.detail(.818, .427, '│', 65)
	if h >= 24 {
		p.detail(.445, .345, '╱', 151)
		p.detail(.545, .43, '╱', 108)
		p.detail(.25, .865, '╱', 180)
	}
}
