package render

// RenderMercenaryMockup shows a weary guild hireling resting both hands on a
// planted cleaver. Unequal armor and a scarred open face tell the same story.
func RenderMercenaryMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the Mercenary", 180, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawMercenaryPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		title, subtitle := "MERCENARY", "paid to be here"
		putString(f, (w-len(title))/2, 1, title, 180, 233, true)
		putString(f, (w-len(subtitle))/2, 2, subtitle, 240, 233, false)
	}
	hint := "tab Wizard · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawMercenaryPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.smoothOval
	poly(236, portraitPoint{.17, .944}, portraitPoint{.79, .944}, portraitPoint{.84, .973}, portraitPoint{.14, .973})
	// A worn coat pools behind the wide legs, rather than billowing heroically.
	poly(52, portraitPoint{.335, .354}, portraitPoint{.64, .342}, portraitPoint{.726, .444}, portraitPoint{.74, .668}, portraitPoint{.784, .837}, portraitPoint{.703, .843}, portraitPoint{.68, .886}, portraitPoint{.62, .826}, portraitPoint{.563, .65}, portraitPoint{.37, .684}, portraitPoint{.269, .81}, portraitPoint{.234, .789}, portraitPoint{.242, .559})
	poly(95, portraitPoint{.645, .399}, portraitPoint{.702, .459}, portraitPoint{.712, .652}, portraitPoint{.755, .821}, portraitPoint{.703, .814}, portraitPoint{.68, .864}, portraitPoint{.648, .814}, portraitPoint{.62, .617})
	// One steel greave and one bound trouser leg are deliberately mismatched.
	poly(236, portraitPoint{.348, .591}, portraitPoint{.487, .625}, portraitPoint{.446, .745}, portraitPoint{.388, .879}, portraitPoint{.279, .919}, portraitPoint{.272, .855}, portraitPoint{.297, .733})
	poly(60, portraitPoint{.356, .64}, portraitPoint{.439, .656}, portraitPoint{.405, .752}, portraitPoint{.368, .867}, portraitPoint{.31, .889}, portraitPoint{.294, .849}, portraitPoint{.318, .749})
	poly(67, portraitPoint{.356, .676}, portraitPoint{.39, .685}, portraitPoint{.363, .752}, portraitPoint{.328, .841}, portraitPoint{.312, .838}, portraitPoint{.332, .755})
	poly(239, portraitPoint{.299, .738}, portraitPoint{.386, .75}, portraitPoint{.386, .806}, portraitPoint{.337, .901}, portraitPoint{.275, .895}, portraitPoint{.28, .831})
	poly(109, portraitPoint{.312, .769}, portraitPoint{.365, .782}, portraitPoint{.345, .838}, portraitPoint{.316, .88}, portraitPoint{.29, .877}, portraitPoint{.296, .831})
	poly(246, portraitPoint{.315, .772}, portraitPoint{.335, .777}, portraitPoint{.32, .834}, portraitPoint{.302, .875}, portraitPoint{.292, .872})
	poly(239, portraitPoint{.31, .703}, portraitPoint{.387, .713}, portraitPoint{.414, .749}, portraitPoint{.385, .781}, portraitPoint{.305, .774}, portraitPoint{.288, .743})
	poly(244, portraitPoint{.317, .717}, portraitPoint{.372, .724}, portraitPoint{.392, .747}, portraitPoint{.371, .762}, portraitPoint{.31, .755}, portraitPoint{.305, .738})
	poly(236, portraitPoint{.506, .608}, portraitPoint{.64, .59}, portraitPoint{.681, .73}, portraitPoint{.704, .877}, portraitPoint{.763, .923}, portraitPoint{.646, .932}, portraitPoint{.61, .852}, portraitPoint{.572, .756})
	poly(58, portraitPoint{.541, .656}, portraitPoint{.615, .642}, portraitPoint{.647, .738}, portraitPoint{.67, .859}, portraitPoint{.654, .9}, portraitPoint{.627, .849}, portraitPoint{.6, .765})
	poly(101, portraitPoint{.606, .691}, portraitPoint{.631, .715}, portraitPoint{.65, .808}, portraitPoint{.652, .844}, portraitPoint{.637, .849}, portraitPoint{.612, .768})
	poly(95, portraitPoint{.626, .741}, portraitPoint{.664, .732}, portraitPoint{.682, .791}, portraitPoint{.639, .805})
	p.stroke(144, 1.2, portraitPoint{.632, .754}, portraitPoint{.667, .747})
	p.stroke(144, 1.0, portraitPoint{.639, .778}, portraitPoint{.673, .772})
	poly(94, portraitPoint{.619, .808}, portraitPoint{.686, .8}, portraitPoint{.7, .866}, portraitPoint{.729, .906}, portraitPoint{.771, .932}, portraitPoint{.76, .948}, portraitPoint{.634, .948}, portraitPoint{.609, .911})
	poly(137, portraitPoint{.63, .823}, portraitPoint{.66, .817}, portraitPoint{.674, .874}, portraitPoint{.686, .908}, portraitPoint{.731, .933}, portraitPoint{.65, .933}, portraitPoint{.625, .902})
	poly(239, portraitPoint{.288, .865}, portraitPoint{.36, .878}, portraitPoint{.389, .922}, portraitPoint{.376, .95}, portraitPoint{.214, .95}, portraitPoint{.203, .928})
	poly(94, portraitPoint{.282, .888}, portraitPoint{.341, .896}, portraitPoint{.367, .925}, portraitPoint{.35, .936}, portraitPoint{.218, .936}, portraitPoint{.229, .926})
	poly(137, portraitPoint{.277, .903}, portraitPoint{.313, .912}, portraitPoint{.297, .93}, portraitPoint{.219, .936}, portraitPoint{.229, .926})
	// Broad sagging leather shoulders sit above a patched, slightly bowed belly.
	poly(236, portraitPoint{.324, .335}, portraitPoint{.414, .319}, portraitPoint{.581, .322}, portraitPoint{.692, .373}, portraitPoint{.732, .475}, portraitPoint{.686, .583}, portraitPoint{.641, .653}, portraitPoint{.336, .665}, portraitPoint{.259, .564}, portraitPoint{.256, .44})
	poly(95, portraitPoint{.345, .356}, portraitPoint{.414, .342}, portraitPoint{.565, .346}, portraitPoint{.654, .386}, portraitPoint{.696, .475}, portraitPoint{.655, .573}, portraitPoint{.617, .626}, portraitPoint{.352, .64}, portraitPoint{.299, .554}, portraitPoint{.285, .457})
	poly(137, portraitPoint{.364, .371}, portraitPoint{.429, .362}, portraitPoint{.494, .389}, portraitPoint{.507, .488}, portraitPoint{.47, .575}, portraitPoint{.384, .603}, portraitPoint{.338, .55}, portraitPoint{.318, .465})
	poly(94, portraitPoint{.494, .389}, portraitPoint{.55, .372}, portraitPoint{.625, .401}, portraitPoint{.66, .473}, portraitPoint{.615, .571}, portraitPoint{.543, .605}, portraitPoint{.47, .575}, portraitPoint{.507, .488})
	poly(130, portraitPoint{.579, .349}, portraitPoint{.614, .367}, portraitPoint{.398, .596}, portraitPoint{.361, .574})
	poly(180, portraitPoint{.487, .459}, portraitPoint{.508, .434}, portraitPoint{.534, .461}, portraitPoint{.51, .49})
	poly(94, portraitPoint{.502, .457}, portraitPoint{.51, .447}, portraitPoint{.52, .461}, portraitPoint{.51, .476})
	poly(94, portraitPoint{.345, .603}, portraitPoint{.615, .584}, portraitPoint{.652, .607}, portraitPoint{.641, .643}, portraitPoint{.342, .663})
	poly(137, portraitPoint{.35, .607}, portraitPoint{.61, .591}, portraitPoint{.626, .608}, portraitPoint{.346, .636})
	poly(180, portraitPoint{.374, .617}, portraitPoint{.418, .613}, portraitPoint{.42, .65}, portraitPoint{.373, .652})
	poly(95, portraitPoint{.385, .627}, portraitPoint{.408, .625}, portraitPoint{.408, .641}, portraitPoint{.385, .641})
	// A rectangular sewn patch replaces any decorative chest emblem.
	poly(65, portraitPoint{.327, .511}, portraitPoint{.372, .494}, portraitPoint{.401, .54}, portraitPoint{.351, .571})
	poly(101, portraitPoint{.335, .515}, portraitPoint{.368, .509}, portraitPoint{.389, .536}, portraitPoint{.353, .557})
	// The salvaged left pauldron is broad iron; the right is quilted cloth.
	poly(239, portraitPoint{.244, .354}, portraitPoint{.322, .316}, portraitPoint{.407, .34}, portraitPoint{.437, .401}, portraitPoint{.417, .454}, portraitPoint{.308, .483}, portraitPoint{.242, .443}, portraitPoint{.218, .399})
	poly(244, portraitPoint{.255, .367}, portraitPoint{.32, .334}, portraitPoint{.388, .351}, portraitPoint{.415, .398}, portraitPoint{.396, .425}, portraitPoint{.303, .447}, portraitPoint{.253, .42}, portraitPoint{.237, .393})
	poly(246, portraitPoint{.255, .369}, portraitPoint{.32, .339}, portraitPoint{.36, .35}, portraitPoint{.355, .381}, portraitPoint{.284, .418}, portraitPoint{.246, .405}, portraitPoint{.237, .393})
	poly(109, portraitPoint{.313, .434}, portraitPoint{.397, .41}, portraitPoint{.417, .435}, portraitPoint{.393, .463}, portraitPoint{.305, .474}, portraitPoint{.262, .453}, portraitPoint{.257, .434})
	poly(239, portraitPoint{.386, .382}, portraitPoint{.416, .399}, portraitPoint{.414, .423}, portraitPoint{.381, .412})
	poly(65, portraitPoint{.602, .338}, portraitPoint{.659, .354}, portraitPoint{.72, .408}, portraitPoint{.741, .462}, portraitPoint{.7, .504}, portraitPoint{.62, .484}, portraitPoint{.576, .43}, portraitPoint{.565, .374})
	poly(101, portraitPoint{.608, .355}, portraitPoint{.653, .374}, portraitPoint{.7, .418}, portraitPoint{.709, .455}, portraitPoint{.685, .475}, portraitPoint{.627, .456}, portraitPoint{.593, .405})
	poly(144, portraitPoint{.615, .356}, portraitPoint{.632, .367}, portraitPoint{.603, .416}, portraitPoint{.59, .405})
	poly(58, portraitPoint{.664, .399}, portraitPoint{.692, .419}, portraitPoint{.703, .452}, portraitPoint{.672, .468}, portraitPoint{.654, .446})
	// Elbows are carried outward; both forearms slope toward the same support.
	poly(239, portraitPoint{.267, .43}, portraitPoint{.31, .46}, portraitPoint{.337, .518}, portraitPoint{.309, .561}, portraitPoint{.276, .575}, portraitPoint{.235, .53}, portraitPoint{.227, .477})
	poly(109, portraitPoint{.26, .457}, portraitPoint{.29, .473}, portraitPoint{.313, .517}, portraitPoint{.292, .545}, portraitPoint{.269, .533}, portraitPoint{.246, .486})
	poly(239, portraitPoint{.285, .525}, portraitPoint{.331, .487}, portraitPoint{.426, .462}, portraitPoint{.477, .484}, portraitPoint{.47, .528}, portraitPoint{.359, .565}, portraitPoint{.298, .579}, portraitPoint{.272, .56})
	poly(244, portraitPoint{.3, .527}, portraitPoint{.336, .506}, portraitPoint{.425, .479}, portraitPoint{.45, .489}, portraitPoint{.447, .514}, portraitPoint{.349, .542}, portraitPoint{.303, .556}, portraitPoint{.288, .545})
	poly(109, portraitPoint{.334, .541}, portraitPoint{.451, .509}, portraitPoint{.469, .525}, portraitPoint{.357, .561}, portraitPoint{.299, .575}, portraitPoint{.296, .563})
	poly(180, portraitPoint{.409, .465}, portraitPoint{.425, .461}, portraitPoint{.445, .526}, portraitPoint{.427, .535})
	poly(95, portraitPoint{.696, .455}, portraitPoint{.729, .471}, portraitPoint{.733, .533}, portraitPoint{.707, .575}, portraitPoint{.67, .565}, portraitPoint{.645, .518}, portraitPoint{.66, .487})
	poly(137, portraitPoint{.69, .478}, portraitPoint{.711, .49}, portraitPoint{.712, .533}, portraitPoint{.696, .547}, portraitPoint{.68, .533}, portraitPoint{.674, .505})
	poly(94, portraitPoint{.696, .526}, portraitPoint{.654, .504}, portraitPoint{.573, .492}, portraitPoint{.533, .515}, portraitPoint{.552, .553}, portraitPoint{.625, .576}, portraitPoint{.698, .579}, portraitPoint{.719, .553})
	poly(137, portraitPoint{.688, .536}, portraitPoint{.649, .518}, portraitPoint{.585, .506}, portraitPoint{.566, .521}, portraitPoint{.586, .541}, portraitPoint{.654, .555}, portraitPoint{.69, .558}, portraitPoint{.707, .548})
	poly(180, portraitPoint{.583, .499}, portraitPoint{.598, .501}, portraitPoint{.588, .551}, portraitPoint{.573, .545})
	// A short wrapped grip supports a cleaver with a chipped, blunt lower edge.
	p.stroke(94, 3, portraitPoint{.517, .468}, portraitPoint{.511, .641})
	p.stroke(137, 1.6, portraitPoint{.51, .468}, portraitPoint{.504, .639})
	oval(239, .517, .465, .027, .02)
	oval(180, .514, .461, .016, .012)
	poly(239, portraitPoint{.467, .603}, portraitPoint{.552, .603}, portraitPoint{.565, .645}, portraitPoint{.462, .65})
	poly(137, portraitPoint{.469, .61}, portraitPoint{.55, .61}, portraitPoint{.553, .63}, portraitPoint{.466, .634})
	poly(239, portraitPoint{.433, .631}, portraitPoint{.623, .65}, portraitPoint{.661, .86}, portraitPoint{.634, .91}, portraitPoint{.485, .957}, portraitPoint{.393, .926}, portraitPoint{.402, .873}, portraitPoint{.421, .86}, portraitPoint{.405, .844})
	poly(244, portraitPoint{.44, .645}, portraitPoint{.605, .661}, portraitPoint{.641, .856}, portraitPoint{.621, .894}, portraitPoint{.483, .935}, portraitPoint{.41, .916}, portraitPoint{.419, .877}, portraitPoint{.431, .863}, portraitPoint{.418, .837})
	poly(251, portraitPoint{.44, .645}, portraitPoint{.469, .648}, portraitPoint{.457, .818}, portraitPoint{.453, .86}, portraitPoint{.426, .907}, portraitPoint{.483, .929}, portraitPoint{.483, .942}, portraitPoint{.407, .919}, portraitPoint{.417, .877}, portraitPoint{.433, .863}, portraitPoint{.418, .837})
	poly(109, portraitPoint{.559, .657}, portraitPoint{.605, .661}, portraitPoint{.641, .856}, portraitPoint{.621, .894}, portraitPoint{.541, .918}, portraitPoint{.56, .864})
	oval(239, .587, .7, .019, .023)
	oval(234, .584, .699, .011, .014)
	// Coins hang off the free hip, worn rather than bright treasure.
	p.stroke(94, .8, portraitPoint{.628, .596}, portraitPoint{.688, .648}, portraitPoint{.675, .686})
	oval(137, .655, .623, .014, .018)
	oval(180, .652, .618, .007, .009)
	oval(137, .693, .652, .014, .018)
	oval(180, .69, .647, .007, .009)
	oval(137, .674, .688, .014, .018)
	oval(180, .671, .683, .007, .009)
	// The hands are bare beyond fingerless cuffs, one resting over the other.
	poly(137, portraitPoint{.437, .47}, portraitPoint{.467, .456}, portraitPoint{.502, .46}, portraitPoint{.531, .477}, portraitPoint{.553, .487}, portraitPoint{.548, .511}, portraitPoint{.497, .517}, portraitPoint{.457, .508}, portraitPoint{.432, .495})
	poly(223, portraitPoint{.448, .472}, portraitPoint{.469, .467}, portraitPoint{.497, .47}, portraitPoint{.524, .484}, portraitPoint{.545, .489}, portraitPoint{.542, .499}, portraitPoint{.494, .502}, portraitPoint{.46, .494}, portraitPoint{.438, .486})
	poly(137, portraitPoint{.53, .501}, portraitPoint{.556, .501}, portraitPoint{.581, .517}, portraitPoint{.578, .544}, portraitPoint{.55, .552}, portraitPoint{.51, .54}, portraitPoint{.484, .522}, portraitPoint{.491, .505})
	poly(223, portraitPoint{.532, .508}, portraitPoint{.552, .511}, portraitPoint{.569, .521}, portraitPoint{.567, .533}, portraitPoint{.548, .539}, portraitPoint{.512, .528}, portraitPoint{.495, .514}, portraitPoint{.509, .507})
	// A heavy neck and rolled collar lead up to an openly weary human face.
	poly(95, portraitPoint{.435, .301}, portraitPoint{.553, .3}, portraitPoint{.59, .369}, portraitPoint{.559, .408}, portraitPoint{.459, .403}, portraitPoint{.404, .366})
	poly(137, portraitPoint{.45, .322}, portraitPoint{.536, .32}, portraitPoint{.558, .36}, portraitPoint{.536, .391}, portraitPoint{.472, .387}, portraitPoint{.43, .36})
	poly(94, portraitPoint{.407, .343}, portraitPoint{.451, .368}, portraitPoint{.494, .406}, portraitPoint{.46, .423}, portraitPoint{.394, .385})
	poly(95, portraitPoint{.536, .377}, portraitPoint{.57, .333}, portraitPoint{.603, .37}, portraitPoint{.553, .42}, portraitPoint{.51, .406})
	poly(137, portraitPoint{.443, .185}, portraitPoint{.501, .169}, portraitPoint{.56, .201}, portraitPoint{.585, .257}, portraitPoint{.579, .309}, portraitPoint{.553, .355}, portraitPoint{.504, .379}, portraitPoint{.455, .351}, portraitPoint{.424, .295}, portraitPoint{.419, .237})
	poly(223, portraitPoint{.45, .197}, portraitPoint{.494, .181}, portraitPoint{.546, .207}, portraitPoint{.56, .246}, portraitPoint{.541, .283}, portraitPoint{.49, .3}, portraitPoint{.445, .269}, portraitPoint{.435, .236})
	poly(173, portraitPoint{.443, .273}, portraitPoint{.48, .29}, portraitPoint{.497, .322}, portraitPoint{.468, .334}, portraitPoint{.447, .315}, portraitPoint{.428, .282})
	poly(95, portraitPoint{.554, .256}, portraitPoint{.571, .25}, portraitPoint{.574, .299}, portraitPoint{.553, .333}, portraitPoint{.532, .321}, portraitPoint{.539, .29})
	poly(223, portraitPoint{.521, .265}, portraitPoint{.553, .292}, portraitPoint{.559, .309}, portraitPoint{.533, .318}, portraitPoint{.511, .3})
	// A short dark beard squares the jaw without hiding the mouth or expression.
	poly(238, portraitPoint{.448, .309}, portraitPoint{.47, .32}, portraitPoint{.496, .336}, portraitPoint{.54, .331}, portraitPoint{.567, .314}, portraitPoint{.563, .347}, portraitPoint{.536, .378}, portraitPoint{.51, .367}, portraitPoint{.496, .387}, portraitPoint{.474, .366}, portraitPoint{.454, .364})
	poly(95, portraitPoint{.474, .32}, portraitPoint{.508, .313}, portraitPoint{.542, .325}, portraitPoint{.546, .342}, portraitPoint{.49, .343})
	p.stroke(137, .9, portraitPoint{.485, .343}, portraitPoint{.54, .342})
	p.stroke(234, 1.05, portraitPoint{.456, .274}, portraitPoint{.49, .279})
	p.stroke(234, 1.05, portraitPoint{.532, .279}, portraitPoint{.558, .273})
	// A scar crosses one brow and cheek; its pale edge remains part of the skin.
	p.stroke(131, 1.15, portraitPoint{.454, .232}, portraitPoint{.472, .266}, portraitPoint{.499, .314})
	p.stroke(223, .65, portraitPoint{.461, .237}, portraitPoint{.481, .269}, portraitPoint{.506, .315})
	// Dented open helmet: a chipped crown, broad brow and unequal cheek guards.
	poly(239, portraitPoint{.388, .224}, portraitPoint{.396, .159}, portraitPoint{.439, .107}, portraitPoint{.463, .134}, portraitPoint{.485, .095}, portraitPoint{.556, .12}, portraitPoint{.603, .183}, portraitPoint{.616, .247}, portraitPoint{.594, .288}, portraitPoint{.568, .228}, portraitPoint{.449, .215}, portraitPoint{.419, .285}, portraitPoint{.396, .269})
	poly(241, portraitPoint{.405, .207}, portraitPoint{.412, .165}, portraitPoint{.438, .128}, portraitPoint{.462, .151}, portraitPoint{.487, .111}, portraitPoint{.542, .134}, portraitPoint{.584, .185}, portraitPoint{.592, .216}, portraitPoint{.568, .207}, portraitPoint{.454, .194}, portraitPoint{.427, .232})
	poly(246, portraitPoint{.438, .128}, portraitPoint{.462, .151}, portraitPoint{.487, .111}, portraitPoint{.516, .124}, portraitPoint{.514, .18}, portraitPoint{.454, .194}, portraitPoint{.414, .206}, portraitPoint{.412, .165})
	poly(109, portraitPoint{.516, .124}, portraitPoint{.542, .134}, portraitPoint{.584, .185}, portraitPoint{.592, .216}, portraitPoint{.568, .207}, portraitPoint{.514, .18})
	poly(239, portraitPoint{.404, .212}, portraitPoint{.453, .196}, portraitPoint{.566, .205}, portraitPoint{.595, .224}, portraitPoint{.6, .247}, portraitPoint{.553, .231}, portraitPoint{.452, .227}, portraitPoint{.417, .244})
	poly(137, portraitPoint{.412, .221}, portraitPoint{.453, .206}, portraitPoint{.557, .213}, portraitPoint{.591, .229}, portraitPoint{.591, .237}, portraitPoint{.551, .224}, portraitPoint{.455, .22}, portraitPoint{.42, .235})
	poly(244, portraitPoint{.4, .244}, portraitPoint{.429, .228}, portraitPoint{.433, .282}, portraitPoint{.451, .314}, portraitPoint{.43, .326}, portraitPoint{.405, .295})
	poly(94, portraitPoint{.569, .25}, portraitPoint{.584, .259}, portraitPoint{.589, .3}, portraitPoint{.566, .335}, portraitPoint{.554, .325}, portraitPoint{.57, .295})
	// Surface-bound seams, scuffs and fingers complete the practical details.
	// Eye marks also cover half-block cells left by the scar strokes.
	for _, u := range []float64{.474, .545} {
		f.Set(x0+int(u*float64(w-1)), y0+int(.28*float64(h-1)), Cell{R: '━', FG: 234, BG: 223})
	}
	p.detail(.52, .344, '─', 137)
	p.detail(.478, .487, '│', 137)
	p.detail(.532, .54, '│', 137)
	p.detail(.275, .386, '●', 137)
	p.detail(.36, .39, '╲', 246)
	p.detail(.346, .541, '╳', 144)
	p.detail(.657, .433, '╲', 144)
	if h >= 24 {
		p.detail(.549, .155, '╲', 239)
		p.detail(.452, .32, '╱', 223)
		p.detail(.604, .853, '╲', 244)
		p.detail(.653, .86, '│', 180)
		p.detail(.303, .93, '─', 137)
	}
}
