package render

import "github.com/0xbenc/termtd/internal/copytext"

// RenderRogueMockup is the guild's knife-runner: a low, forward-leaning crouch
// framed by opposing hooked blades and two trailing lengths of red scarf.
func RenderRogueMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, copytext.Text("characters.rogue.enlarge_to_view_the_rogue"), 131, 233, false)
	} else {
		ph := min(h-7, (w-6)/3)
		pw := ph * 3
		drawRoguePortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		title, subtitle := copytext.Text("characters.rogue.rogue"), copytext.Text("characters.rogue.already_behind_you")
		putString(f, (w-len(title))/2, 1, title, 131, 233, true)
		putString(f, (w-len(subtitle))/2, 2, subtitle, 252, 233, false)
	}
	hint := copytext.Format("characters.rogue.tab_mercenary_esc_return_q_quit", "tab", "tab", "escape", "esc", "quit", "q")
	if w < 36 {
		hint = copytext.Format("characters.rogue.esc_return_tab_portraits", "tab", "tab", "escape", "esc")
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 252, 233, false)
	return f
}

func drawRoguePortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly := p.poly
	poly(236, portraitPoint{.19, .945}, portraitPoint{.80, .945}, portraitPoint{.865, .97}, portraitPoint{.16, .97})
	// Two cloth lengths sweep back through empty sky, with broad angular folds.
	poly(52, portraitPoint{.635, .341}, portraitPoint{.537, .349}, portraitPoint{.442, .267}, portraitPoint{.35, .205}, portraitPoint{.24, .24}, portraitPoint{.143, .225}, portraitPoint{.056, .151}, portraitPoint{.029, .093}, portraitPoint{.101, .145}, portraitPoint{.165, .168}, portraitPoint{.253, .156}, portraitPoint{.367, .155}, portraitPoint{.452, .216}, portraitPoint{.567, .26})
	poly(95, portraitPoint{.62, .321}, portraitPoint{.534, .325}, portraitPoint{.44, .244}, portraitPoint{.35, .184}, portraitPoint{.239, .217}, portraitPoint{.144, .207}, portraitPoint{.06, .153}, portraitPoint{.034, .105}, portraitPoint{.112, .155}, portraitPoint{.171, .185}, portraitPoint{.256, .171}, portraitPoint{.367, .166}, portraitPoint{.448, .228}, portraitPoint{.564, .273})
	poly(131, portraitPoint{.367, .166}, portraitPoint{.448, .228}, portraitPoint{.564, .273}, portraitPoint{.62, .321}, portraitPoint{.57, .31}, portraitPoint{.433, .251}, portraitPoint{.345, .19}, portraitPoint{.247, .224}, portraitPoint{.146, .213}, portraitPoint{.12, .192}, portraitPoint{.238, .205})
	poly(52, portraitPoint{.581, .363}, portraitPoint{.48, .358}, portraitPoint{.376, .32}, portraitPoint{.274, .35}, portraitPoint{.17, .315}, portraitPoint{.10, .282}, portraitPoint{.143, .344}, portraitPoint{.256, .397}, portraitPoint{.373, .354}, portraitPoint{.464, .412}, portraitPoint{.546, .43})
	poly(95, portraitPoint{.576, .371}, portraitPoint{.479, .375}, portraitPoint{.376, .338}, portraitPoint{.274, .369}, portraitPoint{.177, .329}, portraitPoint{.143, .313}, portraitPoint{.168, .342}, portraitPoint{.255, .383}, portraitPoint{.378, .351}, portraitPoint{.47, .401}, portraitPoint{.548, .417})
	poly(131, portraitPoint{.375, .338}, portraitPoint{.474, .372}, portraitPoint{.576, .371}, portraitPoint{.56, .39}, portraitPoint{.474, .387}, portraitPoint{.374, .352}, portraitPoint{.273, .383}, portraitPoint{.251, .378})
	// A stretched rear leg and a deeply bent front knee give the crouch force.
	poly(236, portraitPoint{.414, .581}, portraitPoint{.497, .625}, portraitPoint{.408, .741}, portraitPoint{.301, .768}, portraitPoint{.263, .87}, portraitPoint{.257, .925}, portraitPoint{.186, .934}, portraitPoint{.187, .843}, portraitPoint{.218, .749}, portraitPoint{.33, .676})
	poly(60, portraitPoint{.42, .61}, portraitPoint{.463, .637}, portraitPoint{.381, .723}, portraitPoint{.277, .75}, portraitPoint{.23, .852}, portraitPoint{.22, .907}, portraitPoint{.195, .915}, portraitPoint{.202, .85}, portraitPoint{.235, .76}, portraitPoint{.339, .697})
	poly(67, portraitPoint{.404, .651}, portraitPoint{.42, .669}, portraitPoint{.37, .713}, portraitPoint{.282, .74}, portraitPoint{.26, .774}, portraitPoint{.24, .77}, portraitPoint{.258, .732}, portraitPoint{.344, .699})
	poly(94, portraitPoint{.208, .807}, portraitPoint{.268, .827}, portraitPoint{.245, .897}, portraitPoint{.291, .925}, portraitPoint{.288, .947}, portraitPoint{.183, .947}, portraitPoint{.176, .913})
	poly(137, portraitPoint{.213, .834}, portraitPoint{.242, .842}, portraitPoint{.225, .9}, portraitPoint{.244, .924}, portraitPoint{.204, .924}, portraitPoint{.193, .909})
	poly(236, portraitPoint{.476, .575}, portraitPoint{.542, .572}, portraitPoint{.68, .647}, portraitPoint{.725, .709}, portraitPoint{.697, .772}, portraitPoint{.615, .839}, portraitPoint{.656, .9}, portraitPoint{.618, .932}, portraitPoint{.553, .86}, portraitPoint{.558, .813}, portraitPoint{.628, .73}, portraitPoint{.541, .711}, portraitPoint{.437, .649})
	poly(60, portraitPoint{.482, .602}, portraitPoint{.53, .6}, portraitPoint{.658, .664}, portraitPoint{.69, .708}, portraitPoint{.675, .748}, portraitPoint{.591, .831}, portraitPoint{.627, .899}, portraitPoint{.604, .909}, portraitPoint{.571, .849}, portraitPoint{.58, .813}, portraitPoint{.652, .719}, portraitPoint{.552, .69}, portraitPoint{.459, .641})
	poly(67, portraitPoint{.537, .617}, portraitPoint{.643, .67}, portraitPoint{.67, .702}, portraitPoint{.651, .723}, portraitPoint{.64, .71}, portraitPoint{.553, .675}, portraitPoint{.497, .645})
	poly(94, portraitPoint{.573, .808}, portraitPoint{.621, .822}, portraitPoint{.629, .864}, portraitPoint{.686, .899}, portraitPoint{.781, .924}, portraitPoint{.789, .944}, portraitPoint{.611, .949}, portraitPoint{.583, .921}, portraitPoint{.555, .86})
	poly(137, portraitPoint{.583, .836}, portraitPoint{.607, .843}, portraitPoint{.612, .877}, portraitPoint{.653, .912}, portraitPoint{.705, .929}, portraitPoint{.63, .93}, portraitPoint{.6, .913}, portraitPoint{.574, .861})
	poly(94, portraitPoint{.621, .66}, portraitPoint{.682, .683}, portraitPoint{.709, .713}, portraitPoint{.685, .756}, portraitPoint{.637, .757}, portraitPoint{.623, .733})
	poly(137, portraitPoint{.635, .684}, portraitPoint{.671, .696}, portraitPoint{.69, .715}, portraitPoint{.675, .738}, portraitPoint{.65, .74}, portraitPoint{.639, .72})
	// The far arm reaches backward, with a bent elbow and a reverse knife grip.
	poly(236, portraitPoint{.5, .365}, portraitPoint{.548, .408}, portraitPoint{.474, .51}, portraitPoint{.411, .576}, portraitPoint{.279, .629}, portraitPoint{.23, .675}, portraitPoint{.19, .665}, portraitPoint{.205, .606}, portraitPoint{.365, .536}, portraitPoint{.4, .435})
	poly(23, portraitPoint{.495, .39}, portraitPoint{.525, .42}, portraitPoint{.451, .5}, portraitPoint{.4, .553}, portraitPoint{.269, .602}, portraitPoint{.227, .647}, portraitPoint{.21, .639}, portraitPoint{.237, .604}, portraitPoint{.38, .547}, portraitPoint{.419, .45})
	poly(66, portraitPoint{.426, .457}, portraitPoint{.441, .479}, portraitPoint{.405, .54}, portraitPoint{.279, .592}, portraitPoint{.276, .61}, portraitPoint{.393, .566}, portraitPoint{.422, .535}, portraitPoint{.465, .486})
	poly(94, portraitPoint{.238, .597}, portraitPoint{.272, .593}, portraitPoint{.285, .625}, portraitPoint{.251, .645}, portraitPoint{.224, .621})
	// A leather cuirass narrows from forward shoulders to the low hip.
	poly(236, portraitPoint{.549, .335}, portraitPoint{.65, .379}, portraitPoint{.659, .47}, portraitPoint{.595, .558}, portraitPoint{.508, .636}, portraitPoint{.419, .624}, portraitPoint{.382, .578}, portraitPoint{.415, .477}, portraitPoint{.472, .412})
	poly(94, portraitPoint{.544, .36}, portraitPoint{.618, .386}, portraitPoint{.632, .457}, portraitPoint{.575, .539}, portraitPoint{.49, .608}, portraitPoint{.43, .596}, portraitPoint{.409, .573}, portraitPoint{.437, .487}, portraitPoint{.489, .43})
	poly(137, portraitPoint{.544, .369}, portraitPoint{.571, .388}, portraitPoint{.578, .441}, portraitPoint{.534, .51}, portraitPoint{.463, .568}, portraitPoint{.435, .569}, portraitPoint{.458, .496}, portraitPoint{.504, .443})
	poly(95, portraitPoint{.586, .406}, portraitPoint{.616, .405}, portraitPoint{.623, .455}, portraitPoint{.574, .524}, portraitPoint{.514, .575}, portraitPoint{.514, .549}, portraitPoint{.564, .491})
	poly(130, portraitPoint{.496, .409}, portraitPoint{.521, .387}, portraitPoint{.595, .538}, portraitPoint{.57, .559})
	poly(180, portraitPoint{.537, .465}, portraitPoint{.558, .47}, portraitPoint{.573, .5}, portraitPoint{.552, .516}, portraitPoint{.533, .486})
	poly(94, portraitPoint{.551, .476}, portraitPoint{.559, .481}, portraitPoint{.565, .493}, portraitPoint{.552, .502}, portraitPoint{.544, .484})
	// Short split coat skirts leave the crouching legs unencumbered.
	poly(94, portraitPoint{.421, .552}, portraitPoint{.481, .578}, portraitPoint{.444, .688}, portraitPoint{.366, .731}, portraitPoint{.368, .68}, portraitPoint{.382, .604})
	poly(137, portraitPoint{.416, .578}, portraitPoint{.44, .588}, portraitPoint{.419, .664}, portraitPoint{.379, .706}, portraitPoint{.388, .645})
	poly(95, portraitPoint{.496, .589}, portraitPoint{.556, .563}, portraitPoint{.576, .649}, portraitPoint{.553, .715}, portraitPoint{.516, .675})
	poly(94, portraitPoint{.42, .579}, portraitPoint{.522, .595}, portraitPoint{.565, .572}, portraitPoint{.577, .603}, portraitPoint{.525, .632}, portraitPoint{.41, .613})
	poly(180, portraitPoint{.48, .594}, portraitPoint{.507, .597}, portraitPoint{.509, .625}, portraitPoint{.479, .619})
	poly(239, portraitPoint{.486, .6}, portraitPoint{.5, .603}, portraitPoint{.501, .616}, portraitPoint{.486, .613})
	// A soft pouch follows the rear hip instead of adding another weapon.
	poly(94, portraitPoint{.399, .612}, portraitPoint{.431, .621}, portraitPoint{.412, .701}, portraitPoint{.381, .713}, portraitPoint{.36, .68}, portraitPoint{.367, .641})
	poly(137, portraitPoint{.379, .639}, portraitPoint{.415, .635}, portraitPoint{.401, .664}, portraitPoint{.366, .656})
	// The near arm folds tightly under the forward head, bracer ending at steel.
	poly(236, portraitPoint{.618, .374}, portraitPoint{.659, .403}, portraitPoint{.714, .49}, portraitPoint{.746, .581}, portraitPoint{.722, .617}, portraitPoint{.671, .596}, portraitPoint{.639, .537}, portraitPoint{.577, .489}, portraitPoint{.582, .424})
	poly(23, portraitPoint{.617, .394}, portraitPoint{.644, .416}, portraitPoint{.69, .503}, portraitPoint{.72, .564}, portraitPoint{.71, .589}, portraitPoint{.685, .58}, portraitPoint{.664, .525}, portraitPoint{.603, .479}, portraitPoint{.599, .436})
	poly(66, portraitPoint{.624, .422}, portraitPoint{.64, .439}, portraitPoint{.677, .505}, portraitPoint{.693, .539}, portraitPoint{.683, .554}, portraitPoint{.65, .509}, portraitPoint{.611, .476}, portraitPoint{.615, .445})
	poly(236, portraitPoint{.713, .55}, portraitPoint{.782, .5}, portraitPoint{.825, .5}, portraitPoint{.851, .535}, portraitPoint{.809, .578}, portraitPoint{.739, .613}, portraitPoint{.713, .601}, portraitPoint{.697, .583})
	poly(94, portraitPoint{.737, .552}, portraitPoint{.784, .513}, portraitPoint{.81, .514}, portraitPoint{.83, .535}, portraitPoint{.797, .561}, portraitPoint{.747, .587}, portraitPoint{.728, .58})
	poly(137, portraitPoint{.745, .548}, portraitPoint{.783, .52}, portraitPoint{.799, .521}, portraitPoint{.801, .535}, portraitPoint{.764, .563}, portraitPoint{.739, .573}, portraitPoint{.734, .564})
	poly(180, portraitPoint{.776, .516}, portraitPoint{.791, .51}, portraitPoint{.814, .55}, portraitPoint{.798, .562})
	// A pointed hood frames one human eye; the scarf masks the rest of the face.
	poly(236, portraitPoint{.541, .1}, portraitPoint{.63, .146}, portraitPoint{.688, .205}, portraitPoint{.727, .292}, portraitPoint{.693, .38}, portraitPoint{.639, .419}, portraitPoint{.551, .406}, portraitPoint{.477, .33}, portraitPoint{.473, .242}, portraitPoint{.496, .164})
	poly(23, portraitPoint{.541, .1}, portraitPoint{.63, .146}, portraitPoint{.681, .214}, portraitPoint{.705, .281}, portraitPoint{.681, .318}, portraitPoint{.649, .254}, portraitPoint{.595, .211}, portraitPoint{.551, .185}, portraitPoint{.5, .257}, portraitPoint{.489, .326}, portraitPoint{.477, .33}, portraitPoint{.473, .242}, portraitPoint{.496, .164})
	poly(66, portraitPoint{.541, .1}, portraitPoint{.63, .146}, portraitPoint{.681, .214}, portraitPoint{.684, .246}, portraitPoint{.629, .205}, portraitPoint{.577, .173}, portraitPoint{.546, .153})
	poly(234, portraitPoint{.576, .218}, portraitPoint{.629, .234}, portraitPoint{.692, .277}, portraitPoint{.705, .322}, portraitPoint{.676, .376}, portraitPoint{.625, .395}, portraitPoint{.565, .36}, portraitPoint{.539, .296})
	poly(137, portraitPoint{.611, .229}, portraitPoint{.655, .245}, portraitPoint{.692, .278}, portraitPoint{.708, .31}, portraitPoint{.752, .337}, portraitPoint{.731, .354}, portraitPoint{.691, .349}, portraitPoint{.667, .384}, portraitPoint{.615, .363}, portraitPoint{.59, .314}, portraitPoint{.592, .268})
	poly(223, portraitPoint{.614, .238}, portraitPoint{.652, .254}, portraitPoint{.683, .28}, portraitPoint{.686, .304}, portraitPoint{.639, .317}, portraitPoint{.603, .288})
	poly(180, portraitPoint{.687, .304}, portraitPoint{.741, .335}, portraitPoint{.731, .344}, portraitPoint{.689, .331})
	poly(236, portraitPoint{.652, .285}, portraitPoint{.684, .281}, portraitPoint{.7, .297}, portraitPoint{.676, .315}, portraitPoint{.653, .306})
	p.stroke(230, 1.0, portraitPoint{.664, .302}, portraitPoint{.685, .296})
	poly(52, portraitPoint{.595, .32}, portraitPoint{.659, .329}, portraitPoint{.71, .319}, portraitPoint{.722, .349}, portraitPoint{.691, .379}, portraitPoint{.679, .407}, portraitPoint{.62, .399}, portraitPoint{.579, .366})
	poly(95, portraitPoint{.604, .332}, portraitPoint{.657, .341}, portraitPoint{.709, .332}, portraitPoint{.713, .349}, portraitPoint{.681, .374}, portraitPoint{.666, .39}, portraitPoint{.622, .383}, portraitPoint{.59, .361})
	poly(131, portraitPoint{.606, .334}, portraitPoint{.655, .346}, portraitPoint{.71, .337}, portraitPoint{.701, .35}, portraitPoint{.661, .36}, portraitPoint{.615, .35})
	poly(95, portraitPoint{.553, .365}, portraitPoint{.596, .349}, portraitPoint{.621, .37}, portraitPoint{.603, .406}, portraitPoint{.56, .415}, portraitPoint{.537, .395})
	// Rear knife: reverse-gripped, with a broad hooked cutting edge below hand.
	p.stroke(94, 2, portraitPoint{.235, .612}, portraitPoint{.19, .705})
	p.stroke(137, .8, portraitPoint{.23, .615}, portraitPoint{.184, .699})
	poly(137, portraitPoint{.159, .689}, portraitPoint{.178, .677}, portraitPoint{.226, .713}, portraitPoint{.227, .734}, portraitPoint{.208, .73}, portraitPoint{.182, .71}, portraitPoint{.169, .721})
	poly(239, portraitPoint{.173, .711}, portraitPoint{.217, .738}, portraitPoint{.195, .836}, portraitPoint{.144, .899}, portraitPoint{.079, .928}, portraitPoint{.042, .927}, portraitPoint{.075, .88}, portraitPoint{.135, .828}, portraitPoint{.154, .766})
	poly(251, portraitPoint{.18, .721}, portraitPoint{.198, .737}, portraitPoint{.179, .827}, portraitPoint{.131, .887}, portraitPoint{.051, .924}, portraitPoint{.091, .887}, portraitPoint{.144, .833}, portraitPoint{.163, .772})
	poly(110, portraitPoint{.198, .737}, portraitPoint{.207, .738}, portraitPoint{.189, .831}, portraitPoint{.14, .891}, portraitPoint{.081, .919}, portraitPoint{.131, .887}, portraitPoint{.179, .827})
	p.stroke(251, 1.15, portraitPoint{.18, .721}, portraitPoint{.179, .827}, portraitPoint{.131, .887}, portraitPoint{.051, .924})
	// Forward knife mirrors the sweep above the fist, leaving sky around its tip.
	p.stroke(94, 2, portraitPoint{.782, .569}, portraitPoint{.827, .496})
	p.stroke(137, .8, portraitPoint{.779, .562}, portraitPoint{.825, .492})
	poly(137, portraitPoint{.799, .479}, portraitPoint{.815, .467}, portraitPoint{.861, .51}, portraitPoint{.855, .54}, portraitPoint{.842, .541}, portraitPoint{.846, .518}, portraitPoint{.816, .497}, portraitPoint{.804, .505})
	poly(239, portraitPoint{.819, .488}, portraitPoint{.847, .503}, portraitPoint{.894, .453}, portraitPoint{.955, .37}, portraitPoint{.978, .264}, portraitPoint{.938, .324}, portraitPoint{.881, .38}, portraitPoint{.841, .447})
	poly(251, portraitPoint{.825, .485}, portraitPoint{.839, .49}, portraitPoint{.882, .443}, portraitPoint{.945, .365}, portraitPoint{.975, .275}, portraitPoint{.942, .335}, portraitPoint{.89, .39}, portraitPoint{.852, .458})
	poly(110, portraitPoint{.839, .49}, portraitPoint{.847, .496}, portraitPoint{.889, .449}, portraitPoint{.948, .369}, portraitPoint{.965, .311}, portraitPoint{.945, .365}, portraitPoint{.882, .443})
	p.stroke(251, 1.15, portraitPoint{.825, .485}, portraitPoint{.882, .443}, portraitPoint{.945, .365}, portraitPoint{.975, .275})
	// Exposed fingers close around each handle, layered over the equipment.
	poly(137, portraitPoint{.212, .614}, portraitPoint{.237, .61}, portraitPoint{.251, .633}, portraitPoint{.237, .663}, portraitPoint{.212, .672}, portraitPoint{.196, .648})
	poly(223, portraitPoint{.215, .62}, portraitPoint{.231, .617}, portraitPoint{.238, .634}, portraitPoint{.223, .649}, portraitPoint{.203, .648}, portraitPoint{.205, .634})
	poly(137, portraitPoint{.797, .518}, portraitPoint{.818, .515}, portraitPoint{.837, .535}, portraitPoint{.831, .558}, portraitPoint{.812, .577}, portraitPoint{.784, .565}, portraitPoint{.782, .543})
	poly(223, portraitPoint{.8, .524}, portraitPoint{.815, .522}, portraitPoint{.825, .539}, portraitPoint{.807, .552}, portraitPoint{.79, .553}, portraitPoint{.79, .54})
	// Short seams and one eye-glint sharpen the face, grips and planted boot.
	p.detail(.672, .3, '━', 230)
	p.detail(.646, .281, '╲', 180)
	p.detail(.651, .363, '╲', 131)
	p.detail(.221, .665, '│', 94)
	p.detail(.809, .545, '│', 94)
	p.detail(.65, .72, '╲', 180)
	if h >= 22 {
		p.detail(.531, .45, '╱', 180)
		p.detail(.28, .62, '╱', 137)
		p.detail(.58, .857, '╲', 180)
		p.detail(.679, .918, '─', 180)
	}
}
