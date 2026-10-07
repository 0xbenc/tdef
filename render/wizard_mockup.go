package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderWizardMockup is a static lore study of the guild's frail spellcaster.
func RenderWizardMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.wizard.enlarge_to_view_the_wizard"), 147, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawWizardPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-6)/2, 1, copytext.Text("characters.wizard.wizard"), 147, 233, true)
		subtitle := copytext.Text("characters.wizard.the_shape_of_a_spell")
		putString(f, (w-len(subtitle))/2, 2, subtitle, 252, 233, false)
	}
	hint := copytext.Format("characters.wizard.tab_centurion_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.wizard.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 252, 233, false)
	return f
}

func drawWizardPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly := p.poly
	// A low plinth grounds the robe and the crooked staff.
	poly(236, portraitPoint{.13, .954}, portraitPoint{.87, .954}, portraitPoint{.9, .977}, portraitPoint{.1, .977})
	// The wooden staff bows away from its owner's hunched shoulders.
	shaft := []portraitPoint{{.8, .951}, {.825, .73}, {.837, .47}, {.821, .32}, {.786, .23}}
	p.stroke(94, 3.3, shaft...)
	p.stroke(180, 1.1, portraitPoint{.789, .941}, portraitPoint{.813, .729}, portraitPoint{.825, .469}, portraitPoint{.809, .324}, portraitPoint{.779, .242})
	p.stroke(137, 3.2, portraitCurve(portraitPoint{.786, .23}, portraitPoint{.717, .139}, portraitPoint{.91, .12}, portraitPoint{.873, .259})...)
	p.stroke(180, 1.1, portraitCurve(portraitPoint{.777, .227}, portraitPoint{.721, .143}, portraitPoint{.89, .134}, portraitPoint{.861, .248})...)
	p.stroke(109, .7, portraitPoint{.824, .181}, portraitPoint{.816, .251})
	poly(117, portraitPoint{.816, .238}, portraitPoint{.838, .27}, portraitPoint{.813, .307}, portraitPoint{.791, .273})
	poly(231, portraitPoint{.816, .238}, portraitPoint{.813, .307}, portraitPoint{.801, .273})
	// The trailing robe is a single heavy triangular silhouette, broken by folds.
	poly(236, portraitPoint{.33, .389}, portraitPoint{.476, .348}, portraitPoint{.596, .385}, portraitPoint{.644, .555}, portraitPoint{.635, .714}, portraitPoint{.744, .925}, portraitPoint{.72, .963}, portraitPoint{.22, .959}, portraitPoint{.155, .919}, portraitPoint{.232, .709}, portraitPoint{.251, .477})
	poly(60, portraitPoint{.331, .402}, portraitPoint{.469, .368}, portraitPoint{.572, .411}, portraitPoint{.614, .575}, portraitPoint{.603, .718}, portraitPoint{.701, .928}, portraitPoint{.639, .944}, portraitPoint{.224, .94}, portraitPoint{.186, .909}, portraitPoint{.268, .701}, portraitPoint{.281, .486})
	poly(61, portraitPoint{.31, .44}, portraitPoint{.389, .409}, portraitPoint{.385, .623}, portraitPoint{.3, .787}, portraitPoint{.247, .928}, portraitPoint{.196, .916}, portraitPoint{.279, .686})
	poly(97, portraitPoint{.327, .434}, portraitPoint{.361, .428}, portraitPoint{.336, .673}, portraitPoint{.262, .91}, portraitPoint{.229, .924}, portraitPoint{.295, .681})
	poly(97, portraitPoint{.534, .396}, portraitPoint{.573, .421}, portraitPoint{.604, .573}, portraitPoint{.597, .719}, portraitPoint{.701, .928}, portraitPoint{.631, .938}, portraitPoint{.55, .774}, portraitPoint{.52, .619})
	poly(140, portraitPoint{.552, .432}, portraitPoint{.573, .448}, portraitPoint{.582, .559}, portraitPoint{.552, .605}, portraitPoint{.545, .502})
	poly(238, portraitPoint{.411, .453}, portraitPoint{.538, .427}, portraitPoint{.576, .612}, portraitPoint{.547, .747}, portraitPoint{.605, .943}, portraitPoint{.323, .94}, portraitPoint{.398, .72})
	poly(23, portraitPoint{.453, .499}, portraitPoint{.537, .468}, portraitPoint{.559, .613}, portraitPoint{.525, .729}, portraitPoint{.557, .936}, portraitPoint{.361, .932}, portraitPoint{.435, .709})
	// Two heavy gold facings frame the open chest; neither is a floating outline.
	poly(137, portraitPoint{.412, .423}, portraitPoint{.437, .449}, portraitPoint{.41, .67}, portraitPoint{.338, .935}, portraitPoint{.31, .93}, portraitPoint{.386, .667})
	poly(180, portraitPoint{.412, .447}, portraitPoint{.424, .458}, portraitPoint{.396, .671}, portraitPoint{.324, .928}, portraitPoint{.316, .922}, portraitPoint{.384, .663})
	poly(137, portraitPoint{.535, .426}, portraitPoint{.556, .446}, portraitPoint{.58, .677}, portraitPoint{.653, .935}, portraitPoint{.625, .94}, portraitPoint{.551, .681})
	poly(180, portraitPoint{.539, .45}, portraitPoint{.549, .461}, portraitPoint{.564, .679}, portraitPoint{.637, .932}, portraitPoint{.628, .935}, portraitPoint{.552, .682})
	// Slippers peek from beneath the enormous hem.
	poly(94, portraitPoint{.282, .933}, portraitPoint{.37, .933}, portraitPoint{.368, .965}, portraitPoint{.238, .965})
	poly(137, portraitPoint{.281, .948}, portraitPoint{.332, .942}, portraitPoint{.331, .958}, portraitPoint{.249, .958})
	poly(94, portraitPoint{.549, .939}, portraitPoint{.619, .941}, portraitPoint{.675, .962}, portraitPoint{.55, .967})
	// The near sleeve cups the spell; the other hand grips the staff.
	poly(60, portraitPoint{.298, .442}, portraitPoint{.328, .46}, portraitPoint{.315, .559}, portraitPoint{.401, .628}, portraitPoint{.439, .616}, portraitPoint{.453, .67}, portraitPoint{.396, .726}, portraitPoint{.297, .669}, portraitPoint{.225, .591}, portraitPoint{.241, .492})
	poly(97, portraitPoint{.266, .492}, portraitPoint{.299, .474}, portraitPoint{.28, .571}, portraitPoint{.385, .647}, portraitPoint{.413, .637}, portraitPoint{.425, .665}, portraitPoint{.391, .688}, portraitPoint{.292, .637}, portraitPoint{.248, .582})
	poly(140, portraitPoint{.252, .559}, portraitPoint{.27, .575}, portraitPoint{.306, .619}, portraitPoint{.392, .663}, portraitPoint{.385, .68}, portraitPoint{.29, .637}, portraitPoint{.245, .587})
	poly(180, portraitPoint{.394, .63}, portraitPoint{.41, .625}, portraitPoint{.44, .67}, portraitPoint{.423, .682})
	poly(60, portraitPoint{.559, .413}, portraitPoint{.613, .419}, portraitPoint{.657, .497}, portraitPoint{.734, .477}, portraitPoint{.779, .433}, portraitPoint{.8, .466}, portraitPoint{.775, .558}, portraitPoint{.653, .606}, portraitPoint{.587, .559})
	poly(97, portraitPoint{.599, .434}, portraitPoint{.636, .502}, portraitPoint{.664, .539}, portraitPoint{.745, .515}, portraitPoint{.78, .479}, portraitPoint{.777, .534}, portraitPoint{.65, .577}, portraitPoint{.616, .549}, portraitPoint{.578, .461})
	poly(140, portraitPoint{.626, .501}, portraitPoint{.656, .546}, portraitPoint{.753, .52}, portraitPoint{.741, .54}, portraitPoint{.647, .575}, portraitPoint{.611, .539})
	poly(180, portraitPoint{.764, .457}, portraitPoint{.781, .452}, portraitPoint{.792, .521}, portraitPoint{.775, .527})
	poly(137, portraitPoint{.782, .452}, portraitPoint{.803, .428}, portraitPoint{.835, .433}, portraitPoint{.853, .451}, portraitPoint{.851, .481}, portraitPoint{.827, .495}, portraitPoint{.797, .489})
	poly(223, portraitPoint{.79, .451}, portraitPoint{.804, .441}, portraitPoint{.834, .444}, portraitPoint{.841, .455}, portraitPoint{.832, .468}, portraitPoint{.807, .473}, portraitPoint{.794, .468})
	// A clasped field grimoire hangs clear of the robe, pages facing the viewer.
	p.stroke(94, .9, portraitPoint{.351, .683}, portraitPoint{.294, .744})
	poly(95, portraitPoint{.229, .741}, portraitPoint{.323, .709}, portraitPoint{.361, .826}, portraitPoint{.264, .863})
	poly(180, portraitPoint{.244, .745}, portraitPoint{.317, .723}, portraitPoint{.342, .818}, portraitPoint{.27, .842})
	poly(223, portraitPoint{.25, .75}, portraitPoint{.312, .734}, portraitPoint{.331, .812}, portraitPoint{.277, .832})
	poly(137, portraitPoint{.229, .741}, portraitPoint{.244, .745}, portraitPoint{.27, .842}, portraitPoint{.264, .863})
	poly(95, portraitPoint{.239, .772}, portraitPoint{.33, .747}, portraitPoint{.34, .779}, portraitPoint{.249, .802})
	poly(180, portraitPoint{.282, .758}, portraitPoint{.304, .753}, portraitPoint{.312, .787}, portraitPoint{.29, .793})
	p.detail(.282, .821, '╱', 180)
	// A long, slightly hooked human profile beneath the hat's shadow.
	poly(137, portraitPoint{.397, .282}, portraitPoint{.482, .273}, portraitPoint{.528, .322}, portraitPoint{.533, .348}, portraitPoint{.566, .367}, portraitPoint{.546, .388}, portraitPoint{.517, .38}, portraitPoint{.511, .423}, portraitPoint{.453, .446}, portraitPoint{.405, .388})
	poly(223, portraitPoint{.419, .302}, portraitPoint{.48, .299}, portraitPoint{.506, .329}, portraitPoint{.502, .355}, portraitPoint{.548, .368}, portraitPoint{.536, .377}, portraitPoint{.502, .369}, portraitPoint{.48, .403}, portraitPoint{.431, .386})
	poly(180, portraitPoint{.404, .326}, portraitPoint{.432, .321}, portraitPoint{.439, .376}, portraitPoint{.425, .395}, portraitPoint{.405, .367})
	// Ivory eyebrows, swept whiskers, and a forked beard follow the stooped pose.
	poly(250, portraitPoint{.434, .329}, portraitPoint{.468, .322}, portraitPoint{.494, .338}, portraitPoint{.482, .345}, portraitPoint{.453, .338})
	poly(250, portraitPoint{.421, .37}, portraitPoint{.449, .387}, portraitPoint{.483, .389}, portraitPoint{.505, .382}, portraitPoint{.516, .399}, portraitPoint{.505, .465}, portraitPoint{.467, .531}, portraitPoint{.454, .485}, portraitPoint{.424, .511}, portraitPoint{.408, .449}, portraitPoint{.388, .418})
	poly(231, portraitPoint{.431, .391}, portraitPoint{.454, .409}, portraitPoint{.476, .408}, portraitPoint{.49, .405}, portraitPoint{.486, .457}, portraitPoint{.466, .498}, portraitPoint{.454, .456}, portraitPoint{.433, .482}, portraitPoint{.423, .443})
	poly(246, portraitPoint{.398, .401}, portraitPoint{.424, .434}, portraitPoint{.428, .479}, portraitPoint{.416, .468}, portraitPoint{.405, .439})
	// Huge asymmetric hat: a folded point and broad sagging brim, not a cone.
	poly(60, portraitPoint{.191, .299}, portraitPoint{.309, .252}, portraitPoint{.32, .163}, portraitPoint{.287, .066}, portraitPoint{.377, .028}, portraitPoint{.444, .072}, portraitPoint{.477, .145}, portraitPoint{.533, .231}, portraitPoint{.634, .293}, portraitPoint{.613, .328}, portraitPoint{.483, .303}, portraitPoint{.343, .331}, portraitPoint{.215, .324})
	poly(97, portraitPoint{.323, .248}, portraitPoint{.341, .159}, portraitPoint{.307, .073}, portraitPoint{.372, .047}, portraitPoint{.414, .079}, portraitPoint{.451, .155}, portraitPoint{.508, .235}, portraitPoint{.435, .252})
	poly(140, portraitPoint{.307, .073}, portraitPoint{.372, .047}, portraitPoint{.396, .066}, portraitPoint{.35, .087}, portraitPoint{.375, .173}, portraitPoint{.359, .221}, portraitPoint{.334, .236}, portraitPoint{.341, .159})
	poly(61, portraitPoint{.372, .047}, portraitPoint{.414, .079}, portraitPoint{.438, .131}, portraitPoint{.416, .14}, portraitPoint{.385, .093})
	poly(137, portraitPoint{.31, .248}, portraitPoint{.509, .223}, portraitPoint{.532, .252}, portraitPoint{.306, .281})
	poly(180, portraitPoint{.312, .251}, portraitPoint{.51, .23}, portraitPoint{.516, .241}, portraitPoint{.31, .266})
	poly(97, portraitPoint{.191, .299}, portraitPoint{.324, .278}, portraitPoint{.475, .271}, portraitPoint{.634, .293}, portraitPoint{.602, .307}, portraitPoint{.475, .287}, portraitPoint{.341, .312}, portraitPoint{.218, .314})
	// Spell geometry floats clear of the chest: a hollow diamond around a solid core.
	poly(24, portraitPoint{.537, .453}, portraitPoint{.675, .567}, portraitPoint{.54, .704}, portraitPoint{.402, .576})
	poly(117, portraitPoint{.537, .46}, portraitPoint{.664, .566}, portraitPoint{.54, .693}, portraitPoint{.413, .576})
	poly(23, portraitPoint{.538, .481}, portraitPoint{.638, .568}, portraitPoint{.54, .67}, portraitPoint{.438, .575})
	poly(153, portraitPoint{.539, .514}, portraitPoint{.606, .571}, portraitPoint{.541, .633}, portraitPoint{.477, .575})
	poly(231, portraitPoint{.539, .514}, portraitPoint{.541, .633}, portraitPoint{.505, .576})
	poly(81, portraitPoint{.539, .514}, portraitPoint{.606, .571}, portraitPoint{.565, .581}, portraitPoint{.541, .633})
	// An open hand supports the lower vertex without covering the luminous center.
	poly(137, portraitPoint{.409, .636}, portraitPoint{.434, .632}, portraitPoint{.456, .651}, portraitPoint{.501, .655}, portraitPoint{.526, .637}, portraitPoint{.542, .651}, portraitPoint{.522, .68}, portraitPoint{.459, .691}, portraitPoint{.422, .674})
	poly(223, portraitPoint{.419, .641}, portraitPoint{.433, .641}, portraitPoint{.457, .661}, portraitPoint{.501, .665}, portraitPoint{.526, .649}, portraitPoint{.523, .669}, portraitPoint{.459, .678}, portraitPoint{.429, .665})
	// Half-cell palm edges keep the supporting gesture at the minimum size.
	p.stroke(137, 1.5, portraitPoint{.419, .648}, portraitPoint{.46, .676}, portraitPoint{.5, .675}, portraitPoint{.527, .65})
	p.stroke(223, .7, portraitPoint{.425, .645}, portraitPoint{.462, .665}, portraitPoint{.501, .664}, portraitPoint{.521, .65})
	// Round the eye to the face row; flooring would put it on the brim
	// in compact terminals. Only stamp it onto a skin-colored surface.
	ex, ey := x0+int(.473*float64(w-1)+.5), y0+int(.357*float64(h-1)+.5)
	if c := f.C[ey*f.W+ex]; c.R == '█' && (c.FG == 223 || c.FG == 137 || c.FG == 180) {
		f.Set(ex, ey, Cell{R: '━', FG: 234, BG: c.FG})
	}
	if h >= 18 {
		p.detail(.495, .406, '─', 239)
	}
	p.detail(.817, .462, '│', 137)
	p.detail(.469, .677, '│', 137)
	if h >= 24 {
		p.detail(.367, .198, '◆', 180)
		p.detail(.279, .79, '╱', 140)
		p.detail(.618, .836, '╲', 140)
		p.detail(.437, .461, '╲', 250)
		p.detail(.832, .715, '╱', 180)
	}
}
