package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderLightningMockup shows one definitive Lightning Mage lore portrait.
// The two open claws turn a descending fork into an outward discharge.
func RenderLightningMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.lightning.enlarge_to_view_lightning_mage"), 147, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawLightningPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-14)/2, 1, copytext.Text("characters.lightning.lightning_mage"), 183, 233, true)
		putString(f, (w-20)/2, 2, copytext.Text("characters.lightning.the_storm_takes_hold"), 252, 233, false)
	}
	hint := copytext.Format("characters.lightning.tab_trebuchet_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.lightning.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 252, 233, false)
	return f
}

func drawLightningPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly := p.poly
	// A quiet stone plane anchors the feet while the cloth lifts away from it.
	poly(236, portraitPoint{.24, .94}, portraitPoint{.73, .94}, portraitPoint{.80, .965}, portraitPoint{.20, .965})
	// The airborne coat tails diverge from a narrow, backward-curving waist.
	poly(54, portraitPoint{.44, .465}, portraitPoint{.57, .48}, portraitPoint{.63, .62}, portraitPoint{.75, .73}, portraitPoint{.90, .79}, portraitPoint{.82, .83}, portraitPoint{.94, .87}, portraitPoint{.73, .89}, portraitPoint{.63, .84}, portraitPoint{.54, .69}, portraitPoint{.45, .65}, portraitPoint{.31, .75}, portraitPoint{.12, .795}, portraitPoint{.19, .72}, portraitPoint{.11, .715}, portraitPoint{.29, .64}, portraitPoint{.38, .535})
	poly(97, portraitPoint{.445, .49}, portraitPoint{.50, .55}, portraitPoint{.44, .64}, portraitPoint{.31, .73}, portraitPoint{.15, .77}, portraitPoint{.25, .70}, portraitPoint{.185, .71}, portraitPoint{.315, .66}, portraitPoint{.39, .55})
	poly(60, portraitPoint{.31, .67}, portraitPoint{.405, .56}, portraitPoint{.445, .595}, portraitPoint{.40, .655}, portraitPoint{.25, .735}, portraitPoint{.20, .74})
	poly(139, portraitPoint{.445, .49}, portraitPoint{.455, .525}, portraitPoint{.385, .61}, portraitPoint{.29, .67}, portraitPoint{.15, .77}, portraitPoint{.18, .735}, portraitPoint{.315, .65})
	poly(97, portraitPoint{.57, .515}, portraitPoint{.61, .63}, portraitPoint{.72, .75}, portraitPoint{.88, .80}, portraitPoint{.78, .805}, portraitPoint{.90, .855}, portraitPoint{.75, .855}, portraitPoint{.665, .79}, portraitPoint{.575, .635})
	poly(139, portraitPoint{.60, .61}, portraitPoint{.72, .74}, portraitPoint{.88, .80}, portraitPoint{.80, .80}, portraitPoint{.70, .765}, portraitPoint{.61, .675})
	poly(236, portraitPoint{.67, .75}, portraitPoint{.76, .81}, portraitPoint{.89, .85}, portraitPoint{.75, .85}, portraitPoint{.66, .80})
	// Two bent, narrow legs stay planted below the split front of the robe.
	poly(60, portraitPoint{.46, .645}, portraitPoint{.55, .66}, portraitPoint{.48, .775}, portraitPoint{.39, .91}, portraitPoint{.32, .915}, portraitPoint{.375, .765})
	poly(97, portraitPoint{.46, .685}, portraitPoint{.48, .715}, portraitPoint{.405, .82}, portraitPoint{.365, .89}, portraitPoint{.34, .89}, portraitPoint{.395, .775})
	poly(60, portraitPoint{.535, .65}, portraitPoint{.61, .67}, portraitPoint{.59, .79}, portraitPoint{.65, .90}, portraitPoint{.59, .925}, portraitPoint{.525, .80}, portraitPoint{.515, .735})
	poly(97, portraitPoint{.575, .70}, portraitPoint{.595, .72}, portraitPoint{.568, .79}, portraitPoint{.61, .885}, portraitPoint{.585, .885}, portraitPoint{.54, .79})
	poly(236, portraitPoint{.355, .835}, portraitPoint{.41, .85}, portraitPoint{.375, .905}, portraitPoint{.385, .94}, portraitPoint{.255, .94}, portraitPoint{.26, .92}, portraitPoint{.325, .895})
	poly(60, portraitPoint{.36, .85}, portraitPoint{.395, .855}, portraitPoint{.355, .915}, portraitPoint{.285, .925}, portraitPoint{.335, .885})
	poly(147, portraitPoint{.35, .85}, portraitPoint{.365, .855}, portraitPoint{.33, .895}, portraitPoint{.275, .92}, portraitPoint{.265, .92}, portraitPoint{.32, .885})
	poly(236, portraitPoint{.565, .835}, portraitPoint{.61, .825}, portraitPoint{.645, .90}, portraitPoint{.725, .925}, portraitPoint{.72, .945}, portraitPoint{.595, .945}, portraitPoint{.58, .90})
	poly(60, portraitPoint{.59, .84}, portraitPoint{.61, .84}, portraitPoint{.63, .90}, portraitPoint{.69, .925}, portraitPoint{.62, .925}, portraitPoint{.60, .885})
	poly(147, portraitPoint{.585, .84}, portraitPoint{.595, .838}, portraitPoint{.618, .90}, portraitPoint{.655, .917}, portraitPoint{.63, .917}, portraitPoint{.602, .90})
	// Exposed throat and chest arch back against the descending current.
	poly(103, portraitPoint{.49, .365}, portraitPoint{.59, .385}, portraitPoint{.61, .445}, portraitPoint{.58, .51}, portraitPoint{.565, .575}, portraitPoint{.49, .605}, portraitPoint{.435, .56}, portraitPoint{.43, .46})
	poly(139, portraitPoint{.49, .38}, portraitPoint{.54, .39}, portraitPoint{.565, .45}, portraitPoint{.54, .51}, portraitPoint{.525, .57}, portraitPoint{.49, .57}, portraitPoint{.455, .515}, portraitPoint{.46, .45})
	poly(182, portraitPoint{.49, .40}, portraitPoint{.515, .405}, portraitPoint{.525, .45}, portraitPoint{.492, .48}, portraitPoint{.47, .465}, portraitPoint{.475, .435})
	poly(60, portraitPoint{.535, .44}, portraitPoint{.56, .425}, portraitPoint{.57, .475}, portraitPoint{.535, .51}, portraitPoint{.52, .50})
	// Raised, asymmetric lapels and an angular sash repeat the bolt's zigzag.
	poly(54, portraitPoint{.44, .40}, portraitPoint{.485, .44}, portraitPoint{.455, .50}, portraitPoint{.49, .60}, portraitPoint{.465, .655}, portraitPoint{.40, .58}, portraitPoint{.395, .465}, portraitPoint{.365, .405})
	poly(97, portraitPoint{.365, .405}, portraitPoint{.44, .40}, portraitPoint{.485, .44}, portraitPoint{.455, .475}, portraitPoint{.415, .445})
	poly(147, portraitPoint{.365, .405}, portraitPoint{.44, .40}, portraitPoint{.48, .435}, portraitPoint{.463, .445}, portraitPoint{.432, .415})
	poly(97, portraitPoint{.575, .425}, portraitPoint{.625, .42}, portraitPoint{.65, .485}, portraitPoint{.61, .545}, portraitPoint{.60, .62}, portraitPoint{.545, .665}, portraitPoint{.555, .565}, portraitPoint{.58, .49})
	poly(139, portraitPoint{.595, .43}, portraitPoint{.624, .425}, portraitPoint{.644, .475}, portraitPoint{.615, .505}, portraitPoint{.60, .56}, portraitPoint{.58, .575}, portraitPoint{.603, .49})
	poly(180, portraitPoint{.485, .59}, portraitPoint{.595, .575}, portraitPoint{.608, .612}, portraitPoint{.477, .632})
	poly(101, portraitPoint{.48, .615}, portraitPoint{.605, .598}, portraitPoint{.608, .621}, portraitPoint{.472, .645})
	poly(183, portraitPoint{.55, .577}, portraitPoint{.572, .601}, portraitPoint{.55, .634}, portraitPoint{.525, .611})
	poly(195, portraitPoint{.55, .582}, portraitPoint{.55, .612}, portraitPoint{.529, .611})
	// The receiving arm bends upward: palm separated from face by black sky.
	poly(60, portraitPoint{.415, .435}, portraitPoint{.455, .465}, portraitPoint{.405, .505}, portraitPoint{.34, .475}, portraitPoint{.31, .395}, portraitPoint{.275, .335}, portraitPoint{.285, .28}, portraitPoint{.335, .30}, portraitPoint{.35, .385}, portraitPoint{.365, .42})
	poly(103, portraitPoint{.40, .442}, portraitPoint{.43, .467}, portraitPoint{.395, .478}, portraitPoint{.355, .447}, portraitPoint{.33, .38}, portraitPoint{.30, .32}, portraitPoint{.322, .305}, portraitPoint{.353, .365}, portraitPoint{.37, .423})
	poly(182, portraitPoint{.30, .32}, portraitPoint{.318, .305}, portraitPoint{.345, .37}, portraitPoint{.356, .422}, portraitPoint{.375, .443}, portraitPoint{.355, .438}, portraitPoint{.33, .385})
	poly(180, portraitPoint{.299, .336}, portraitPoint{.331, .319}, portraitPoint{.346, .354}, portraitPoint{.315, .375})
	poly(101, portraitPoint{.315, .36}, portraitPoint{.345, .344}, portraitPoint{.355, .375}, portraitPoint{.326, .39})
	// A second arm reaches down and outward, directing rather than catching.
	poly(60, portraitPoint{.625, .425}, portraitPoint{.67, .455}, portraitPoint{.71, .505}, portraitPoint{.795, .475}, portraitPoint{.82, .50}, portraitPoint{.795, .545}, portraitPoint{.70, .565}, portraitPoint{.65, .51}, portraitPoint{.60, .485})
	poly(139, portraitPoint{.647, .453}, portraitPoint{.67, .46}, portraitPoint{.71, .518}, portraitPoint{.788, .49}, portraitPoint{.8, .512}, portraitPoint{.704, .54}, portraitPoint{.655, .488})
	poly(182, portraitPoint{.655, .456}, portraitPoint{.67, .46}, portraitPoint{.711, .518}, portraitPoint{.785, .49}, portraitPoint{.79, .50}, portraitPoint{.708, .53}, portraitPoint{.656, .48})
	poly(180, portraitPoint{.764, .492}, portraitPoint{.784, .48}, portraitPoint{.803, .519}, portraitPoint{.783, .532})
	poly(101, portraitPoint{.774, .51}, portraitPoint{.798, .50}, portraitPoint{.81, .525}, portraitPoint{.789, .54})
	// Wind-swept hair and two ivory horns give the head its own hooked contour.
	poly(236, portraitPoint{.52, .235}, portraitPoint{.625, .27}, portraitPoint{.71, .315}, portraitPoint{.67, .33}, portraitPoint{.72, .385}, portraitPoint{.66, .375}, portraitPoint{.68, .44}, portraitPoint{.60, .41}, portraitPoint{.55, .355})
	poly(60, portraitPoint{.59, .275}, portraitPoint{.64, .29}, portraitPoint{.70, .32}, portraitPoint{.65, .32}, portraitPoint{.69, .38}, portraitPoint{.65, .36}, portraitPoint{.66, .41}, portraitPoint{.60, .385}, portraitPoint{.575, .34})
	poly(103, portraitPoint{.585, .26}, portraitPoint{.62, .225}, portraitPoint{.69, .21}, portraitPoint{.74, .155}, portraitPoint{.72, .23}, portraitPoint{.66, .27}, portraitPoint{.625, .30})
	poly(183, portraitPoint{.602, .262}, portraitPoint{.628, .238}, portraitPoint{.696, .22}, portraitPoint{.74, .155}, portraitPoint{.715, .225}, portraitPoint{.655, .265}, portraitPoint{.62, .285})
	poly(195, portraitPoint{.628, .238}, portraitPoint{.69, .216}, portraitPoint{.74, .155}, portraitPoint{.70, .225}, portraitPoint{.647, .249})
	poly(103, portraitPoint{.625, .31}, portraitPoint{.69, .28}, portraitPoint{.755, .275}, portraitPoint{.825, .205}, portraitPoint{.81, .28}, portraitPoint{.775, .325}, portraitPoint{.70, .35}, portraitPoint{.64, .35})
	poly(183, portraitPoint{.644, .314}, portraitPoint{.69, .294}, portraitPoint{.76, .293}, portraitPoint{.825, .205}, portraitPoint{.795, .288}, portraitPoint{.765, .317}, portraitPoint{.699, .338}, portraitPoint{.645, .336})
	poly(195, portraitPoint{.69, .294}, portraitPoint{.755, .281}, portraitPoint{.825, .205}, portraitPoint{.776, .297})
	// Blade ear, tilted forehead, long upturned nose, dark open mouth.
	poly(139, portraitPoint{.59, .305}, portraitPoint{.68, .34}, portraitPoint{.65, .37}, portraitPoint{.595, .35})
	poly(60, portraitPoint{.608, .325}, portraitPoint{.656, .348}, portraitPoint{.638, .356}, portraitPoint{.603, .34})
	poly(103, portraitPoint{.49, .252}, portraitPoint{.555, .265}, portraitPoint{.605, .305}, portraitPoint{.614, .36}, portraitPoint{.585, .405}, portraitPoint{.53, .405}, portraitPoint{.49, .371}, portraitPoint{.47, .327}, portraitPoint{.435, .289}, portraitPoint{.463, .282})
	poly(139, portraitPoint{.49, .257}, portraitPoint{.54, .271}, portraitPoint{.573, .30}, portraitPoint{.56, .332}, portraitPoint{.517, .352}, portraitPoint{.485, .32}, portraitPoint{.45, .292}, portraitPoint{.468, .288})
	poly(182, portraitPoint{.49, .259}, portraitPoint{.535, .274}, portraitPoint{.552, .293}, portraitPoint{.519, .307}, portraitPoint{.481, .296}, portraitPoint{.445, .289}, portraitPoint{.461, .277})
	poly(60, portraitPoint{.575, .312}, portraitPoint{.595, .325}, portraitPoint{.595, .365}, portraitPoint{.568, .392}, portraitPoint{.55, .371}, portraitPoint{.56, .35})
	poly(139, portraitPoint{.517, .349}, portraitPoint{.55, .34}, portraitPoint{.572, .373}, portraitPoint{.56, .394}, portraitPoint{.533, .386}, portraitPoint{.505, .36})
	poly(234, portraitPoint{.482, .333}, portraitPoint{.515, .343}, portraitPoint{.538, .37}, portraitPoint{.522, .38}, portraitPoint{.496, .359})
	poly(195, portraitPoint{.483, .331}, portraitPoint{.505, .34}, portraitPoint{.51, .351}, portraitPoint{.492, .346})
	poly(183, portraitPoint{.52, .379}, portraitPoint{.535, .367}, portraitPoint{.519, .345})
	poly(54, portraitPoint{.505, .292}, portraitPoint{.533, .294}, portraitPoint{.543, .311}, portraitPoint{.516, .318}, portraitPoint{.499, .306})
	p.stroke(195, 1.0, portraitPoint{.505, .307}, portraitPoint{.53, .294})
	poly(60, portraitPoint{.535, .396}, portraitPoint{.557, .403}, portraitPoint{.545, .44}, portraitPoint{.525, .41})
	poly(103, portraitPoint{.55, .405}, portraitPoint{.574, .40}, portraitPoint{.563, .434})
	// Broad dim edges and fine white cores keep the forks bright and legible.
	strength := min(1.0, float64(h)/30)
	bolt := func(width float64, path ...portraitPoint) {
		width *= strength
		p.stroke(54, width+2.8*strength, path...)
		p.stroke(99, width+1.3*strength, path...)
		p.stroke(147, max(.9, width), path...)
		p.stroke(195, max(.7, width*.40), path...)
	}
	bolt(2.3, portraitPoint{.405, .025}, portraitPoint{.33, .13}, portraitPoint{.405, .12}, portraitPoint{.275, .285})
	bolt(1.2, portraitPoint{.369, .073}, portraitPoint{.555, .055}, portraitPoint{.495, .13}, portraitPoint{.58, .14}, portraitPoint{.545, .20})
	bolt(1.0, portraitPoint{.343, .145}, portraitPoint{.21, .12}, portraitPoint{.24, .18}, portraitPoint{.16, .215})
	bolt(1.8, portraitPoint{.838, .493}, portraitPoint{.89, .43}, portraitPoint{.85, .432}, portraitPoint{.94, .34}, portraitPoint{.967, .265})
	bolt(.9, portraitPoint{.891, .429}, portraitPoint{.95, .49}, portraitPoint{.918, .493}, portraitPoint{.968, .56})
	// Open receiving claw lies over the bolt so the contact has visible fingers.
	poly(139, portraitPoint{.277, .33}, portraitPoint{.315, .32}, portraitPoint{.323, .28}, portraitPoint{.355, .25}, portraitPoint{.347, .235}, portraitPoint{.313, .258}, portraitPoint{.305, .207}, portraitPoint{.29, .201}, portraitPoint{.29, .263}, portraitPoint{.264, .202}, portraitPoint{.247, .19}, portraitPoint{.253, .22}, portraitPoint{.271, .278}, portraitPoint{.231, .237}, portraitPoint{.211, .234}, portraitPoint{.225, .265}, portraitPoint{.262, .307})
	poly(182, portraitPoint{.275, .30}, portraitPoint{.291, .281}, portraitPoint{.30, .265}, portraitPoint{.296, .22}, portraitPoint{.306, .223}, portraitPoint{.313, .267}, portraitPoint{.339, .249}, portraitPoint{.344, .253}, portraitPoint{.316, .283}, portraitPoint{.308, .309}, portraitPoint{.288, .319})
	poly(195, portraitPoint{.254, .21}, portraitPoint{.266, .23}, portraitPoint{.282, .278}, portraitPoint{.274, .286}, portraitPoint{.261, .245})
	// Long directing fingers form a second fan, with sky between the tips.
	poly(139, portraitPoint{.795, .49}, portraitPoint{.82, .465}, portraitPoint{.827, .427}, portraitPoint{.839, .411}, portraitPoint{.844, .423}, portraitPoint{.842, .467}, portraitPoint{.873, .443}, portraitPoint{.9, .441}, portraitPoint{.892, .455}, portraitPoint{.853, .484}, portraitPoint{.90, .485}, portraitPoint{.917, .498}, portraitPoint{.91, .508}, portraitPoint{.85, .507}, portraitPoint{.88, .535}, portraitPoint{.873, .549}, portraitPoint{.84, .532}, portraitPoint{.80, .53}, portraitPoint{.787, .514})
	poly(182, portraitPoint{.798, .494}, portraitPoint{.824, .478}, portraitPoint{.832, .435}, portraitPoint{.84, .422}, portraitPoint{.835, .48}, portraitPoint{.873, .45}, portraitPoint{.885, .45}, portraitPoint{.842, .491}, portraitPoint{.902, .49}, portraitPoint{.909, .501}, portraitPoint{.838, .50}, portraitPoint{.815, .519}, portraitPoint{.80, .516})
	poly(195, portraitPoint{.817, .489}, portraitPoint{.835, .484}, portraitPoint{.847, .495}, portraitPoint{.833, .508}, portraitPoint{.817, .505})
	// Short reflected seams and features stay attached to their material planes.
	p.detail(.54, .345, '╲', 183)
	p.detail(.578, .388, '╱', 182)
	p.detail(.285, .30, '╱', 103)
	p.detail(.823, .518, '╲', 103)
	if h >= 24 {
		p.detail(.52, .454, '╲', 183)
		p.detail(.395, .692, '╱', 139)
		p.detail(.755, .803, '╲', 139)
	}
}
