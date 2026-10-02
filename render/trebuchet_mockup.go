package render

// RenderTrebuchetMockup shows a single definitive siege-engine lore portrait.
// The open triangles expose the axle, suspended weight, sling and winding line.
func RenderTrebuchetMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view the Trebuchet", 180, 233, false)
	} else {
		ph := min(h-7, (w-6)*2/5)
		pw := ph * 5 / 2
		drawTrebuchetPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-9)/2, 1, "TREBUCHET", 180, 233, true)
		putString(f, (w-25)/2, 2, "a little help with gravity", 240, 233, false)
	}
	hint := "tab Necromancer · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab portraits"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

func drawTrebuchetPortrait(f *Frame, x0, y0, w, h int) {
	p := portraitPainter{f, x0, y0, w, h}
	poly, oval := p.poly, p.smoothOval
	poly(236, portraitPoint{.03, .945}, portraitPoint{.91, .945}, portraitPoint{.96, .975}, portraitPoint{.03, .975})
	// Rear timbers are deliberately darker: two A frames have real depth.
	poly(94, portraitPoint{.57, .345}, portraitPoint{.605, .365}, portraitPoint{.43, .875}, portraitPoint{.375, .875})
	poly(95, portraitPoint{.592, .35}, portraitPoint{.622, .35}, portraitPoint{.87, .855}, portraitPoint{.82, .875})
	poly(130, portraitPoint{.606, .383}, portraitPoint{.62, .38}, portraitPoint{.85, .835}, portraitPoint{.835, .84})
	poly(94, portraitPoint{.39, .675}, portraitPoint{.75, .675}, portraitPoint{.765, .72}, portraitPoint{.385, .72})
	poly(130, portraitPoint{.39, .675}, portraitPoint{.75, .675}, portraitPoint{.755, .69}, portraitPoint{.39, .69})
	poly(94, portraitPoint{.42, .88}, portraitPoint{.84, .835}, portraitPoint{.89, .86}, portraitPoint{.48, .918})
	poly(137, portraitPoint{.42, .88}, portraitPoint{.84, .835}, portraitPoint{.87, .85}, portraitPoint{.48, .90})
	// Long skids, end grain and low cross ties make the engine feel grounded.
	poly(94, portraitPoint{.225, .905}, portraitPoint{.745, .905}, portraitPoint{.88, .845}, portraitPoint{.91, .875}, portraitPoint{.76, .955}, portraitPoint{.225, .955})
	poly(130, portraitPoint{.23, .905}, portraitPoint{.745, .905}, portraitPoint{.88, .845}, portraitPoint{.895, .86}, portraitPoint{.752, .928}, portraitPoint{.23, .928})
	poly(137, portraitPoint{.23, .905}, portraitPoint{.745, .905}, portraitPoint{.752, .923}, portraitPoint{.23, .923})
	poly(180, portraitPoint{.235, .905}, portraitPoint{.745, .905}, portraitPoint{.745, .913}, portraitPoint{.235, .913})
	poly(95, portraitPoint{.24, .905}, portraitPoint{.28, .905}, portraitPoint{.28, .953}, portraitPoint{.24, .953})
	poly(239, portraitPoint{.705, .905}, portraitPoint{.738, .905}, portraitPoint{.738, .95}, portraitPoint{.705, .95})
	// The free-hanging sling has two taut sides, a stone and a leather cradle.
	p.stroke(137, 1.35, portraitPoint{.118, .073}, portraitPoint{.068, .385}, portraitPoint{.145, .443})
	p.stroke(180, .75, portraitPoint{.119, .073}, portraitPoint{.073, .385})
	p.stroke(180, 1.05, portraitPoint{.131, .067}, portraitPoint{.223, .385}, portraitPoint{.15, .443})
	oval(237, .146, .395, .073, .067)
	poly(243, portraitPoint{.10, .361}, portraitPoint{.16, .344}, portraitPoint{.20, .378}, portraitPoint{.196, .415}, portraitPoint{.145, .439}, portraitPoint{.096, .407}, portraitPoint{.08, .385})
	poly(247, portraitPoint{.10, .361}, portraitPoint{.16, .344}, portraitPoint{.19, .372}, portraitPoint{.143, .39}, portraitPoint{.082, .385})
	poly(239, portraitPoint{.143, .39}, portraitPoint{.19, .372}, portraitPoint{.196, .415}, portraitPoint{.145, .439}, portraitPoint{.116, .416})
	poly(94, portraitPoint{.06, .387}, portraitPoint{.09, .414}, portraitPoint{.146, .437}, portraitPoint{.20, .412}, portraitPoint{.229, .381}, portraitPoint{.22, .422}, portraitPoint{.159, .462}, portraitPoint{.125, .46}, portraitPoint{.073, .424})
	poly(137, portraitPoint{.074, .414}, portraitPoint{.143, .442}, portraitPoint{.206, .412}, portraitPoint{.19, .438}, portraitPoint{.147, .458}, portraitPoint{.11, .444})
	// A hinged iron-bound ballast box hangs from the short end, not the frame.
	p.stroke(239, 1.4, portraitPoint{.647, .57}, portraitPoint{.666, .535}, portraitPoint{.705, .51}, portraitPoint{.767, .546}, portraitPoint{.78, .57})
	p.stroke(239, 2.3, portraitPoint{.705, .505}, portraitPoint{.703, .56})
	p.stroke(244, .8, portraitPoint{.697, .515}, portraitPoint{.697, .565})
	poly(237, portraitPoint{.635, .56}, portraitPoint{.775, .545}, portraitPoint{.824, .576}, portraitPoint{.82, .721}, portraitPoint{.775, .761}, portraitPoint{.635, .731})
	poly(94, portraitPoint{.635, .58}, portraitPoint{.774, .58}, portraitPoint{.774, .75}, portraitPoint{.635, .726})
	poly(95, portraitPoint{.774, .58}, portraitPoint{.82, .577}, portraitPoint{.815, .72}, portraitPoint{.774, .75})
	poly(137, portraitPoint{.636, .561}, portraitPoint{.774, .547}, portraitPoint{.821, .577}, portraitPoint{.774, .591}, portraitPoint{.636, .586})
	poly(239, portraitPoint{.65, .56}, portraitPoint{.765, .554}, portraitPoint{.79, .568}, portraitPoint{.674, .579})
	poly(243, portraitPoint{.656, .56}, portraitPoint{.68, .55}, portraitPoint{.70, .564}, portraitPoint{.697, .575}, portraitPoint{.674, .578})
	poly(241, portraitPoint{.715, .555}, portraitPoint{.735, .547}, portraitPoint{.759, .56}, portraitPoint{.768, .576}, portraitPoint{.734, .579})
	poly(130, portraitPoint{.642, .606}, portraitPoint{.77, .605}, portraitPoint{.77, .627}, portraitPoint{.642, .627})
	poly(130, portraitPoint{.64, .687}, portraitPoint{.774, .711}, portraitPoint{.774, .732}, portraitPoint{.64, .707})
	poly(239, portraitPoint{.66, .582}, portraitPoint{.682, .585}, portraitPoint{.682, .736}, portraitPoint{.66, .731})
	poly(239, portraitPoint{.733, .587}, portraitPoint{.755, .587}, portraitPoint{.755, .748}, portraitPoint{.733, .744})
	poly(244, portraitPoint{.662, .584}, portraitPoint{.669, .585}, portraitPoint{.669, .731}, portraitPoint{.662, .729})
	poly(244, portraitPoint{.735, .589}, portraitPoint{.742, .589}, portraitPoint{.742, .744}, portraitPoint{.735, .742})
	// The front A frame stays open around the operating rope and windlass.
	poly(94, portraitPoint{.49, .375}, portraitPoint{.535, .395}, portraitPoint{.315, .911}, portraitPoint{.255, .918})
	poly(137, portraitPoint{.49, .385}, portraitPoint{.512, .392}, portraitPoint{.29, .907}, portraitPoint{.26, .907})
	poly(180, portraitPoint{.492, .39}, portraitPoint{.501, .395}, portraitPoint{.275, .907}, portraitPoint{.26, .907})
	poly(94, portraitPoint{.517, .375}, portraitPoint{.56, .375}, portraitPoint{.777, .918}, portraitPoint{.718, .932})
	poly(130, portraitPoint{.534, .399}, portraitPoint{.553, .393}, portraitPoint{.76, .904}, portraitPoint{.735, .91})
	poly(180, portraitPoint{.534, .405}, portraitPoint{.542, .403}, portraitPoint{.745, .905}, portraitPoint{.735, .905})
	poly(137, portraitPoint{.337, .707}, portraitPoint{.674, .707}, portraitPoint{.69, .746}, portraitPoint{.323, .75})
	poly(180, portraitPoint{.337, .707}, portraitPoint{.674, .707}, portraitPoint{.678, .719}, portraitPoint{.331, .719})
	poly(94, portraitPoint{.32, .811}, portraitPoint{.625, .601}, portraitPoint{.639, .637}, portraitPoint{.337, .843})
	poly(130, portraitPoint{.327, .813}, portraitPoint{.628, .605}, portraitPoint{.63, .616}, portraitPoint{.332, .827})
	// A long tapered beam crosses the axle; its short end carries the ballast.
	poly(94, portraitPoint{.108, .075}, portraitPoint{.13, .043}, portraitPoint{.565, .358}, portraitPoint{.73, .478}, portraitPoint{.716, .537}, portraitPoint{.51, .414})
	poly(137, portraitPoint{.108, .075}, portraitPoint{.13, .054}, portraitPoint{.566, .37}, portraitPoint{.724, .486}, portraitPoint{.709, .515}, portraitPoint{.513, .397})
	poly(180, portraitPoint{.114, .069}, portraitPoint{.129, .051}, portraitPoint{.562, .367}, portraitPoint{.714, .474}, portraitPoint{.707, .487}, portraitPoint{.514, .382})
	poly(95, portraitPoint{.108, .075}, portraitPoint{.116, .071}, portraitPoint{.518, .394}, portraitPoint{.716, .525}, portraitPoint{.716, .537}, portraitPoint{.51, .414})
	// Three forged bands and the axle distinguish a machine from bare timbers.
	poly(239, portraitPoint{.272, .15}, portraitPoint{.30, .17}, portraitPoint{.28, .203}, portraitPoint{.253, .18})
	poly(244, portraitPoint{.273, .153}, portraitPoint{.286, .162}, portraitPoint{.266, .189}, portraitPoint{.257, .181})
	poly(239, portraitPoint{.416, .253}, portraitPoint{.45, .277}, portraitPoint{.424, .32}, portraitPoint{.391, .291})
	poly(244, portraitPoint{.417, .257}, portraitPoint{.43, .267}, portraitPoint{.405, .303}, portraitPoint{.394, .293})
	poly(239, portraitPoint{.66, .441}, portraitPoint{.694, .465}, portraitPoint{.675, .513}, portraitPoint{.645, .494})
	poly(244, portraitPoint{.66, .447}, portraitPoint{.674, .456}, portraitPoint{.654, .499}, portraitPoint{.646, .491})
	p.stroke(237, 3.3, portraitPoint{.508, .384}, portraitPoint{.596, .35})
	p.stroke(244, 1.4, portraitPoint{.506, .376}, portraitPoint{.595, .343})
	oval(237, .511, .391, .039, .047)
	oval(137, .508, .386, .027, .033)
	oval(180, .504, .379, .011, .014)
	// The cocking line runs from the long arm to a drum beneath the cross tie.
	p.stroke(94, 1.5, portraitPoint{.345, .235}, portraitPoint{.481, .773})
	p.stroke(180, .8, portraitPoint{.348, .234}, portraitPoint{.484, .773})
	poly(237, portraitPoint{.38, .771}, portraitPoint{.605, .771}, portraitPoint{.605, .799}, portraitPoint{.38, .799})
	poly(244, portraitPoint{.385, .771}, portraitPoint{.605, .771}, portraitPoint{.605, .781}, portraitPoint{.385, .781})
	oval(94, .462, .784, .024, .041)
	poly(137, portraitPoint{.46, .745}, portraitPoint{.545, .745}, portraitPoint{.545, .826}, portraitPoint{.46, .826})
	oval(180, .462, .783, .024, .041)
	oval(94, .542, .783, .024, .041)
	for i := 0; i < 4; i++ {
		u := .48 + float64(i)*.016
		p.stroke(94, .8, portraitPoint{u, .748}, portraitPoint{u, .821})
	}
	p.stroke(137, 1.5, portraitPoint{.387, .786}, portraitPoint{.362, .741}, portraitPoint{.319, .766})
	p.stroke(180, .7, portraitPoint{.387, .783}, portraitPoint{.362, .737}, portraitPoint{.319, .761})
	p.stroke(239, 1.7, portraitPoint{.263, .783}, portraitPoint{.325, .758})
	p.stroke(180, .7, portraitPoint{.263, .778}, portraitPoint{.325, .754})
	// The winding crew leans back on the crank, feet resisting its weight.
	poly(60, portraitPoint{.135, .797}, portraitPoint{.197, .801}, portraitPoint{.163, .881}, portraitPoint{.098, .938}, portraitPoint{.06, .928}, portraitPoint{.101, .866})
	poly(66, portraitPoint{.145, .821}, portraitPoint{.16, .847}, portraitPoint{.11, .894}, portraitPoint{.083, .915}, portraitPoint{.075, .909}, portraitPoint{.125, .862})
	poly(60, portraitPoint{.192, .80}, portraitPoint{.236, .80}, portraitPoint{.241, .86}, portraitPoint{.216, .91}, portraitPoint{.25, .937}, portraitPoint{.204, .937}, portraitPoint{.181, .91}, portraitPoint{.204, .851})
	poly(237, portraitPoint{.08, .90}, portraitPoint{.119, .911}, portraitPoint{.106, .937}, portraitPoint{.04, .947}, portraitPoint{.04, .931})
	poly(94, portraitPoint{.078, .917}, portraitPoint{.105, .92}, portraitPoint{.094, .934}, portraitPoint{.049, .939}, portraitPoint{.048, .932})
	poly(237, portraitPoint{.20, .907}, portraitPoint{.224, .907}, portraitPoint{.26, .935}, portraitPoint{.259, .947}, portraitPoint{.191, .947}, portraitPoint{.185, .927})
	poly(94, portraitPoint{.205, .917}, portraitPoint{.22, .917}, portraitPoint{.25, .938}, portraitPoint{.20, .938})
	poly(66, portraitPoint{.152, .697}, portraitPoint{.198, .693}, portraitPoint{.24, .733}, portraitPoint{.232, .811}, portraitPoint{.135, .829}, portraitPoint{.111, .769}, portraitPoint{.12, .727})
	poly(109, portraitPoint{.15, .71}, portraitPoint{.185, .709}, portraitPoint{.22, .741}, portraitPoint{.215, .78}, portraitPoint{.158, .791}, portraitPoint{.124, .763})
	poly(94, portraitPoint{.17, .715}, portraitPoint{.19, .704}, portraitPoint{.215, .807}, portraitPoint{.16, .839}, portraitPoint{.133, .812}, portraitPoint{.155, .771})
	poly(137, portraitPoint{.167, .75}, portraitPoint{.182, .742}, portraitPoint{.193, .807}, portraitPoint{.164, .824}, portraitPoint{.155, .81})
	poly(107, portraitPoint{.125, .624}, portraitPoint{.17, .64}, portraitPoint{.177, .675}, portraitPoint{.139, .67}, portraitPoint{.103, .607})
	poly(65, portraitPoint{.122, .625}, portraitPoint{.155, .649}, portraitPoint{.16, .663}, portraitPoint{.14, .656})
	poly(107, portraitPoint{.17, .603}, portraitPoint{.219, .609}, portraitPoint{.242, .646}, portraitPoint{.241, .666}, portraitPoint{.27, .68}, portraitPoint{.26, .70}, portraitPoint{.227, .702}, portraitPoint{.213, .73}, portraitPoint{.167, .715}, portraitPoint{.146, .669}, portraitPoint{.151, .634})
	poly(150, portraitPoint{.171, .613}, portraitPoint{.214, .62}, portraitPoint{.23, .646}, portraitPoint{.222, .66}, portraitPoint{.174, .655}, portraitPoint{.155, .64})
	poly(65, portraitPoint{.167, .68}, portraitPoint{.21, .69}, portraitPoint{.241, .681}, portraitPoint{.239, .70}, portraitPoint{.213, .73}, portraitPoint{.17, .715})
	poly(107, portraitPoint{.182, .696}, portraitPoint{.238, .691}, portraitPoint{.221, .712}, portraitPoint{.194, .711})
	poly(22, portraitPoint{.207, .65}, portraitPoint{.24, .65}, portraitPoint{.242, .667}, portraitPoint{.209, .674})
	poly(180, portraitPoint{.228, .709}, portraitPoint{.237, .702}, portraitPoint{.237, .68})
	poly(131, portraitPoint{.148, .63}, portraitPoint{.16, .604}, portraitPoint{.216, .603}, portraitPoint{.233, .627}, portraitPoint{.232, .643}, portraitPoint{.187, .633})
	poly(95, portraitPoint{.155, .629}, portraitPoint{.12, .647}, portraitPoint{.09, .685}, portraitPoint{.12, .677}, portraitPoint{.153, .652})
	// Both bent elbows lead to the same crank, rather than floating nearby.
	poly(65, portraitPoint{.133, .739}, portraitPoint{.16, .735}, portraitPoint{.157, .776}, portraitPoint{.258, .763}, portraitPoint{.27, .785}, portraitPoint{.157, .803}, portraitPoint{.123, .78})
	poly(107, portraitPoint{.137, .751}, portraitPoint{.151, .748}, portraitPoint{.149, .788}, portraitPoint{.256, .772}, portraitPoint{.261, .783}, portraitPoint{.151, .795}, portraitPoint{.13, .78})
	poly(65, portraitPoint{.209, .724}, portraitPoint{.237, .736}, portraitPoint{.256, .754}, portraitPoint{.311, .743}, portraitPoint{.326, .763}, portraitPoint{.255, .78}, portraitPoint{.222, .761})
	poly(107, portraitPoint{.22, .735}, portraitPoint{.234, .741}, portraitPoint{.255, .766}, portraitPoint{.314, .75}, portraitPoint{.318, .76}, portraitPoint{.255, .777}, portraitPoint{.229, .755})
	poly(107, portraitPoint{.256, .762}, portraitPoint{.28, .762}, portraitPoint{.289, .776}, portraitPoint{.271, .793}, portraitPoint{.253, .784})
	poly(150, portraitPoint{.26, .768}, portraitPoint{.276, .766}, portraitPoint{.28, .777}, portraitPoint{.262, .783})
	poly(107, portraitPoint{.303, .744}, portraitPoint{.324, .74}, portraitPoint{.335, .757}, portraitPoint{.324, .777}, portraitPoint{.305, .767})
	poly(150, portraitPoint{.31, .744}, portraitPoint{.323, .746}, portraitPoint{.325, .755}, portraitPoint{.311, .762})
	// A smaller spotter braces a bent knee on the rear rail and sights left.
	poly(60, portraitPoint{.838, .746}, portraitPoint{.885, .745}, portraitPoint{.905, .81}, portraitPoint{.897, .887}, portraitPoint{.874, .91}, portraitPoint{.847, .899}, portraitPoint{.875, .843}, portraitPoint{.86, .80})
	poly(94, portraitPoint{.871, .869}, portraitPoint{.898, .87}, portraitPoint{.91, .905}, portraitPoint{.885, .921}, portraitPoint{.85, .914}, portraitPoint{.854, .902})
	poly(60, portraitPoint{.84, .746}, portraitPoint{.86, .772}, portraitPoint{.815, .815}, portraitPoint{.774, .815}, portraitPoint{.779, .782}, portraitPoint{.823, .78})
	poly(94, portraitPoint{.775, .787}, portraitPoint{.797, .79}, portraitPoint{.804, .817}, portraitPoint{.752, .828}, portraitPoint{.747, .808})
	poly(95, portraitPoint{.832, .655}, portraitPoint{.874, .652}, portraitPoint{.9, .688}, portraitPoint{.883, .772}, portraitPoint{.83, .781}, portraitPoint{.811, .724})
	poly(130, portraitPoint{.844, .668}, portraitPoint{.869, .67}, portraitPoint{.88, .693}, portraitPoint{.864, .756}, portraitPoint{.838, .762}, portraitPoint{.829, .718})
	poly(143, portraitPoint{.85, .617}, portraitPoint{.91, .588}, portraitPoint{.899, .639}, portraitPoint{.86, .652})
	poly(101, portraitPoint{.866, .626}, portraitPoint{.897, .607}, portraitPoint{.885, .635}, portraitPoint{.866, .64})
	poly(143, portraitPoint{.829, .578}, portraitPoint{.865, .58}, portraitPoint{.891, .611}, portraitPoint{.876, .654}, portraitPoint{.842, .679}, portraitPoint{.811, .665}, portraitPoint{.80, .645}, portraitPoint{.772, .637}, portraitPoint{.773, .62}, portraitPoint{.812, .608})
	poly(186, portraitPoint{.83, .588}, portraitPoint{.852, .59}, portraitPoint{.873, .614}, portraitPoint{.859, .636}, portraitPoint{.819, .631}, portraitPoint{.789, .632}, portraitPoint{.8, .62}, portraitPoint{.819, .612})
	poly(101, portraitPoint{.86, .64}, portraitPoint{.88, .635}, portraitPoint{.876, .654}, portraitPoint{.842, .679}, portraitPoint{.811, .665}, portraitPoint{.83, .653})
	poly(143, portraitPoint{.818, .664}, portraitPoint{.83, .67}, portraitPoint{.819, .728}, portraitPoint{.781, .734}, portraitPoint{.746, .713}, portraitPoint{.754, .693}, portraitPoint{.787, .712}, portraitPoint{.806, .704})
	poly(186, portraitPoint{.809, .703}, portraitPoint{.818, .686}, portraitPoint{.816, .718}, portraitPoint{.786, .726}, portraitPoint{.752, .706}, portraitPoint{.757, .699}, portraitPoint{.791, .716})
	poly(143, portraitPoint{.874, .685}, portraitPoint{.896, .672}, portraitPoint{.902, .627}, portraitPoint{.862, .58}, portraitPoint{.818, .573}, portraitPoint{.806, .589}, portraitPoint{.851, .598}, portraitPoint{.879, .631}, portraitPoint{.877, .664}, portraitPoint{.863, .673})
	poly(186, portraitPoint{.894, .672}, portraitPoint{.888, .63}, portraitPoint{.856, .588}, portraitPoint{.82, .581}, portraitPoint{.814, .59}, portraitPoint{.85, .595}, portraitPoint{.88, .636}, portraitPoint{.88, .669})
	// Small pins, grain and the crew's eyes stay on their existing surfaces.
	p.detail(.22, .658, '━', 230)
	p.detail(.205, .703, '─', 180)
	p.detail(.825, .626, '━', 230)
	p.detail(.263, .779, '│', 65)
	p.detail(.316, .758, '│', 65)
	p.detail(.755, .705, '│', 101)
	p.detail(.37, .72, '●', 239)
	p.detail(.652, .723, '●', 239)
	p.detail(.69, .64, '╲', 130)
	p.detail(.71, .691, '╱', 137)
	if h >= 24 {
		p.detail(.195, .681, '╲', 150)
		p.detail(.851, .653, '╯', 94)
		p.detail(.395, .304, '╱', 180)
		p.detail(.396, .591, '╱', 180)
		p.detail(.598, .565, '╲', 180)
		p.detail(.30, .932, '─', 137)
	}
}
