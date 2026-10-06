package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderGunnerMockup depicts one definitive Orc Gunner lore pose. The iron
// barrel, braced limbs and concentrated squint establish the character.
func RenderGunnerMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.gunner.enlarge_to_view_the_gunner"), 114, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawGunnerPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-10)/2, 1, copytext.Text("characters.gunner.orc_gunner"), 114, 233, true)
		putString(f, (w-22)/2, 2, copytext.Text("characters.gunner.steady_hands_hot_iron"), 240, 233, false)
	}
	hint := copytext.Format("characters.gunner.tab_frost_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.gunner.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawGunnerPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.oval
	poly(236, portraitPoint{.09, .92}, portraitPoint{.70, .92}, portraitPoint{.77, .95}, portraitPoint{.06, .95})
	// Quiet smoke above the muzzle balances the red knot behind his head.
	poly(238, portraitPoint{.955, .35}, portraitPoint{.91, .30}, portraitPoint{.81, .28}, portraitPoint{.77, .23}, portraitPoint{.77, .18}, portraitPoint{.82, .14}, portraitPoint{.87, .13}, portraitPoint{.89, .15}, portraitPoint{.83, .17}, portraitPoint{.80, .20}, portraitPoint{.81, .23}, portraitPoint{.85, .25}, portraitPoint{.94, .28}, portraitPoint{.975, .32})
	poly(240, portraitPoint{.82, .14}, portraitPoint{.87, .13}, portraitPoint{.89, .15}, portraitPoint{.83, .17}, portraitPoint{.80, .20}, portraitPoint{.80, .18})
	// The back leg reaches left; the forward knee takes the weapon's weight.
	poly(60, portraitPoint{.29, .64}, portraitPoint{.43, .66}, portraitPoint{.33, .78}, portraitPoint{.27, .88}, portraitPoint{.16, .89}, portraitPoint{.20, .77})
	poly(59, portraitPoint{.29, .72}, portraitPoint{.33, .78}, portraitPoint{.27, .88}, portraitPoint{.22, .88}, portraitPoint{.25, .79}, portraitPoint{.21, .78})
	poly(60, portraitPoint{.42, .64}, portraitPoint{.53, .65}, portraitPoint{.60, .77}, portraitPoint{.57, .87}, portraitPoint{.65, .90}, portraitPoint{.55, .90}, portraitPoint{.49, .84}, portraitPoint{.47, .78}, portraitPoint{.40, .74})
	poly(67, portraitPoint{.49, .70}, portraitPoint{.53, .70}, portraitPoint{.59, .78}, portraitPoint{.55, .82}, portraitPoint{.51, .78})
	// Blunt metal-toed boots and a large knee plate make the stance heavy.
	poly(237, portraitPoint{.17, .84}, portraitPoint{.28, .84}, portraitPoint{.29, .89}, portraitPoint{.25, .92}, portraitPoint{.10, .92}, portraitPoint{.10, .89})
	poly(241, portraitPoint{.11, .89}, portraitPoint{.24, .89}, portraitPoint{.24, .92}, portraitPoint{.10, .92})
	poly(237, portraitPoint{.53, .85}, portraitPoint{.59, .85}, portraitPoint{.67, .89}, portraitPoint{.69, .92}, portraitPoint{.53, .92}, portraitPoint{.51, .89})
	poly(241, portraitPoint{.60, .89}, portraitPoint{.67, .89}, portraitPoint{.69, .92}, portraitPoint{.60, .92})
	poly(239, portraitPoint{.51, .73}, portraitPoint{.58, .74}, portraitPoint{.61, .78}, portraitPoint{.57, .82}, portraitPoint{.51, .80}, portraitPoint{.49, .77})
	poly(244, portraitPoint{.52, .74}, portraitPoint{.58, .75}, portraitPoint{.59, .77}, portraitPoint{.51, .77})
	// A broad barrel chest tapers into a leather apron and cartridge belt.
	poly(65, portraitPoint{.30, .36}, portraitPoint{.45, .36}, portraitPoint{.54, .42}, portraitPoint{.56, .58}, portraitPoint{.53, .69}, portraitPoint{.29, .69}, portraitPoint{.22, .59}, portraitPoint{.23, .44})
	poly(71, portraitPoint{.31, .37}, portraitPoint{.43, .37}, portraitPoint{.49, .44}, portraitPoint{.47, .56}, portraitPoint{.30, .59}, portraitPoint{.25, .49})
	poly(94, portraitPoint{.27, .52}, portraitPoint{.50, .53}, portraitPoint{.53, .66}, portraitPoint{.49, .73}, portraitPoint{.39, .69}, portraitPoint{.32, .72}, portraitPoint{.27, .66})
	poly(130, portraitPoint{.30, .54}, portraitPoint{.35, .54}, portraitPoint{.40, .69}, portraitPoint{.32, .70})
	poly(94, portraitPoint{.37, .38}, portraitPoint{.42, .38}, portraitPoint{.32, .65}, portraitPoint{.28, .63})
	poly(180, portraitPoint{.31, .55}, portraitPoint{.345, .56}, portraitPoint{.335, .61}, portraitPoint{.30, .60})
	poly(94, portraitPoint{.29, .64}, portraitPoint{.52, .64}, portraitPoint{.52, .69}, portraitPoint{.29, .69})
	// Three oversized cartridges read as ammunition rather than decoration.
	for i := 0; i < 3; i++ {
		u := .335 + float64(i)*.05
		poly(180, portraitPoint{u, .63}, portraitPoint{u + .025, .63}, portraitPoint{u + .03, .69}, portraitPoint{u, .69})
		poly(130, portraitPoint{u, .68}, portraitPoint{u + .03, .68}, portraitPoint{u + .03, .70}, portraitPoint{u, .70})
	}
	// The support arm has a dropped elbow, then rises into the fore-end.
	poly(65, portraitPoint{.49, .43}, portraitPoint{.58, .44}, portraitPoint{.64, .55}, portraitPoint{.70, .49}, portraitPoint{.76, .48}, portraitPoint{.77, .53}, portraitPoint{.65, .64}, portraitPoint{.59, .64}, portraitPoint{.52, .56})
	poly(71, portraitPoint{.54, .45}, portraitPoint{.58, .46}, portraitPoint{.64, .56}, portraitPoint{.70, .50}, portraitPoint{.74, .50}, portraitPoint{.74, .53}, portraitPoint{.64, .60}, portraitPoint{.59, .58})
	poly(114, portraitPoint{.64, .56}, portraitPoint{.70, .50}, portraitPoint{.74, .50}, portraitPoint{.74, .52}, portraitPoint{.64, .58})
	poly(94, portraitPoint{.68, .51}, portraitPoint{.72, .49}, portraitPoint{.75, .54}, portraitPoint{.71, .57})
	// Two red cloth tails give the otherwise squared-off silhouette a gesture.
	poly(95, portraitPoint{.31, .21}, portraitPoint{.28, .22}, portraitPoint{.20, .28}, portraitPoint{.14, .29}, portraitPoint{.19, .33}, portraitPoint{.28, .28}, portraitPoint{.33, .25})
	poly(131, portraitPoint{.30, .23}, portraitPoint{.27, .25}, portraitPoint{.23, .34}, portraitPoint{.19, .36}, portraitPoint{.23, .37}, portraitPoint{.29, .31}, portraitPoint{.33, .26})
	// Head in profile: short blade ear, shaved dome, heavy brow and jaw.
	poly(71, portraitPoint{.19, .23}, portraitPoint{.33, .25}, portraitPoint{.36, .32}, portraitPoint{.29, .33}, portraitPoint{.23, .29})
	poly(65, portraitPoint{.24, .26}, portraitPoint{.32, .27}, portraitPoint{.33, .30}, portraitPoint{.28, .30})
	poly(71, portraitPoint{.33, .17}, portraitPoint{.44, .16}, portraitPoint{.50, .19}, portraitPoint{.54, .25}, portraitPoint{.55, .28}, portraitPoint{.62, .30}, portraitPoint{.62, .33}, portraitPoint{.57, .35}, portraitPoint{.56, .40}, portraitPoint{.41, .42}, portraitPoint{.32, .36}, portraitPoint{.30, .26})
	poly(114, portraitPoint{.35, .18}, portraitPoint{.44, .175}, portraitPoint{.49, .20}, portraitPoint{.51, .24}, portraitPoint{.36, .245}, portraitPoint{.32, .27})
	poly(65, portraitPoint{.32, .30}, portraitPoint{.40, .33}, portraitPoint{.48, .34}, portraitPoint{.56, .34}, portraitPoint{.56, .40}, portraitPoint{.41, .42}, portraitPoint{.32, .36})
	poly(71, portraitPoint{.38, .35}, portraitPoint{.48, .355}, portraitPoint{.55, .355}, portraitPoint{.55, .39}, portraitPoint{.42, .40}, portraitPoint{.37, .38})
	poly(114, portraitPoint{.53, .27}, portraitPoint{.60, .295}, portraitPoint{.61, .31}, portraitPoint{.54, .315}, portraitPoint{.51, .30})
	poly(22, portraitPoint{.46, .265}, portraitPoint{.53, .25}, portraitPoint{.55, .28}, portraitPoint{.47, .295})
	poly(22, portraitPoint{.44, .35}, portraitPoint{.56, .345}, portraitPoint{.56, .365}, portraitPoint{.46, .375})
	// A single broken tusk and a full tusk are large, tapered ivory shapes.
	poly(180, portraitPoint{.43, .37}, portraitPoint{.455, .37}, portraitPoint{.45, .34}, portraitPoint{.435, .34})
	poly(180, portraitPoint{.53, .37}, portraitPoint{.565, .365}, portraitPoint{.565, .30}, portraitPoint{.545, .325})
	poly(131, portraitPoint{.315, .22}, portraitPoint{.49, .215}, portraitPoint{.52, .245}, portraitPoint{.33, .255}, portraitPoint{.30, .245})
	poly(95, portraitPoint{.315, .245}, portraitPoint{.52, .235}, portraitPoint{.52, .245}, portraitPoint{.33, .26})
	// The worn shoulder is one thick slab, with a bone stud and brass rivets.
	poly(237, portraitPoint{.22, .39}, portraitPoint{.28, .36}, portraitPoint{.37, .36}, portraitPoint{.45, .43}, portraitPoint{.42, .51}, portraitPoint{.24, .51}, portraitPoint{.19, .46})
	poly(241, portraitPoint{.24, .39}, portraitPoint{.29, .375}, portraitPoint{.36, .375}, portraitPoint{.42, .425}, portraitPoint{.40, .47}, portraitPoint{.23, .47}, portraitPoint{.21, .44})
	poly(244, portraitPoint{.245, .39}, portraitPoint{.29, .375}, portraitPoint{.36, .375}, portraitPoint{.38, .39}, portraitPoint{.25, .415}, portraitPoint{.22, .43})
	poly(180, portraitPoint{.28, .37}, portraitPoint{.28, .31}, portraitPoint{.34, .37})
	// Near forearm runs toward the trigger; the mass remains visible below it.
	poly(65, portraitPoint{.25, .49}, portraitPoint{.34, .49}, portraitPoint{.38, .55}, portraitPoint{.46, .49}, portraitPoint{.52, .50}, portraitPoint{.53, .55}, portraitPoint{.42, .63}, portraitPoint{.34, .63}, portraitPoint{.27, .58})
	poly(71, portraitPoint{.28, .50}, portraitPoint{.33, .50}, portraitPoint{.38, .56}, portraitPoint{.47, .50}, portraitPoint{.50, .52}, portraitPoint{.40, .59}, portraitPoint{.34, .60}, portraitPoint{.28, .55})
	poly(114, portraitPoint{.35, .57}, portraitPoint{.38, .56}, portraitPoint{.47, .50}, portraitPoint{.49, .52}, portraitPoint{.39, .58})
	// A battered walnut butt is pressed into the shoulder. Its iron cap
	// and inset grain are broad planes, not scattered surface noise.
	poly(94, portraitPoint{.30, .45}, portraitPoint{.34, .43}, portraitPoint{.54, .475}, portraitPoint{.57, .515}, portraitPoint{.42, .555}, portraitPoint{.31, .52})
	poly(130, portraitPoint{.34, .45}, portraitPoint{.52, .485}, portraitPoint{.52, .51}, portraitPoint{.41, .53}, portraitPoint{.34, .505})
	poly(237, portraitPoint{.30, .45}, portraitPoint{.33, .44}, portraitPoint{.34, .515}, portraitPoint{.31, .52})
	// Long iron barrel and receiver: the upward diagonal dominates the pose.
	poly(237, portraitPoint{.49, .465}, portraitPoint{.88, .375}, portraitPoint{.94, .425}, portraitPoint{.88, .485}, portraitPoint{.53, .54}, portraitPoint{.49, .52})
	poly(242, portraitPoint{.54, .455}, portraitPoint{.88, .38}, portraitPoint{.90, .435}, portraitPoint{.57, .505}, portraitPoint{.53, .495})
	poly(246, portraitPoint{.56, .445}, portraitPoint{.87, .38}, portraitPoint{.885, .40}, portraitPoint{.575, .465})
	poly(239, portraitPoint{.57, .485}, portraitPoint{.89, .42}, portraitPoint{.89, .455}, portraitPoint{.58, .515})
	poly(180, portraitPoint{.68, .415}, portraitPoint{.705, .41}, portraitPoint{.73, .475}, portraitPoint{.70, .482})
	poly(137, portraitPoint{.71, .46}, portraitPoint{.725, .455}, portraitPoint{.73, .475}, portraitPoint{.70, .482})
	poly(180, portraitPoint{.825, .383}, portraitPoint{.85, .377}, portraitPoint{.875, .45}, portraitPoint{.845, .46})
	// A large brass flare and deep, elliptical bore make the muzzle readable.
	poly(137, portraitPoint{.85, .37}, portraitPoint{.94, .35}, portraitPoint{.975, .39}, portraitPoint{.97, .46}, portraitPoint{.94, .49}, portraitPoint{.865, .49})
	poly(180, portraitPoint{.855, .37}, portraitPoint{.94, .35}, portraitPoint{.965, .38}, portraitPoint{.94, .445}, portraitPoint{.865, .46})
	oval(232, .942, .417, .026, .052)
	poly(94, portraitPoint{.865, .46}, portraitPoint{.92, .445}, portraitPoint{.96, .455}, portraitPoint{.94, .49}, portraitPoint{.865, .49})
	// Sight, striker and trigger guard are a few unambiguous hardware shapes.
	poly(180, portraitPoint{.81, .392}, portraitPoint{.81, .365}, portraitPoint{.835, .36}, portraitPoint{.835, .388})
	poly(239, portraitPoint{.515, .465}, portraitPoint{.51, .425}, portraitPoint{.53, .412}, portraitPoint{.56, .42}, portraitPoint{.54, .443}, portraitPoint{.55, .46})
	oval(180, .495, .553, .038, .037)
	oval(233, .495, .553, .024, .023)
	// Fists and finger planes wrap over the gun, completing both arm gestures.
	poly(71, portraitPoint{.46, .50}, portraitPoint{.49, .485}, portraitPoint{.53, .49}, portraitPoint{.545, .525}, portraitPoint{.52, .55}, portraitPoint{.47, .545})
	poly(114, portraitPoint{.475, .495}, portraitPoint{.50, .49}, portraitPoint{.515, .51}, portraitPoint{.48, .525})
	poly(71, portraitPoint{.695, .48}, portraitPoint{.72, .455}, portraitPoint{.745, .45}, portraitPoint{.77, .47}, portraitPoint{.765, .52}, portraitPoint{.725, .54}, portraitPoint{.70, .52})
	poly(114, portraitPoint{.715, .47}, portraitPoint{.745, .455}, portraitPoint{.758, .47}, portraitPoint{.748, .505}, portraitPoint{.72, .515})
	p.detail(.49, .28, '━', 230)
	p.detail(.375, .29, '╲', 114)
	if h >= 24 {
		p.detail(.39, .32, '╲', 114)
	}
	if h < 24 {
		p.detail(.55, .355, '▴', 180)
	}
	p.detail(.26, .445, '●', 180)
	p.detail(.38, .445, '●', 180)
	p.detail(.485, .515, '│', 65)
	p.detail(.515, .515, '│', 65)
	p.detail(.72, .49, '│', 65)
	p.detail(.745, .49, '│', 65)
	p.detail(.565, .78, '╱', 246)
	if h >= 24 {
		p.detail(.625, .465, '╱', 246)
	}
}
