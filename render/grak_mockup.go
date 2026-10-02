package render

// RenderGrakMockup presents a static character study in the same terminal
// tile language as Malgrath: large masses, a few planes and sparse details.
func RenderGrakMockup(w, h int) *Frame {
	w, h = max(0, w), max(0, h)
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	if w < 32 || h < 16 {
		putString(f, 1, h/2, "enlarge to view Grak", 150, 233, false)
	} else {
		ph := min(h-7, (w-6)/2)
		pw := ph * 2
		drawGrakPortrait(f, (w-pw)/2, 4+(h-7-ph)/2, pw, ph)
		putString(f, (w-4)/2, 1, "GRAK", 150, 233, true)
		putString(f, (w-16)/2, 2, "the last monster", 240, 233, false)
	}
	hint := "tab Gnolls · esc return · q quit"
	if w < 36 {
		hint = "esc return · tab Gnolls"
	}
	putString(f, max(0, (w-len([]rune(hint)))/2), h-2, hint, 240, 233, false)
	return f
}

type grakPoint = portraitPoint

// drawGrakPortrait uses a hunched three-quarter stance. The enormous head,
// blade ears and uneven tusks carry the identity; the mantle and grounded
// mallet give him the weight of someone who builds and holds the defenses.
func drawGrakPortrait(f *Frame, x0, y0, w, h int) {
	if w < 2 || h < 2 {
		return
	}
	p := portraitPainter{f, x0, y0, w, h}
	poly := p.poly
	// A quiet stone plinth anchors the stance without a second silhouette.
	poly(236, grakPoint{.12, .91}, grakPoint{.81, .91}, grakPoint{.89, .94}, grakPoint{.08, .94})
	// Torn mantle: broad burgundy wedges echo the dragon's palette.
	poly(95, grakPoint{.34, .39}, grakPoint{.68, .39}, grakPoint{.78, .46}, grakPoint{.84, .76}, grakPoint{.73, .73}, grakPoint{.69, .78}, grakPoint{.36, .73}, grakPoint{.23, .77}, grakPoint{.28, .48})
	poly(52, grakPoint{.67, .43}, grakPoint{.77, .48}, grakPoint{.80, .72}, grakPoint{.68, .69})
	// Widely planted legs and blunt, reinforced leather boots.
	poly(58, grakPoint{.36, .67}, grakPoint{.49, .68}, grakPoint{.47, .83}, grakPoint{.31, .86})
	poly(58, grakPoint{.56, .68}, grakPoint{.68, .67}, grakPoint{.74, .86}, grakPoint{.60, .86})
	poly(94, grakPoint{.31, .79}, grakPoint{.47, .79}, grakPoint{.46, .89}, grakPoint{.26, .89}, grakPoint{.26, .85})
	poly(94, grakPoint{.59, .79}, grakPoint{.73, .79}, grakPoint{.79, .85}, grakPoint{.79, .89}, grakPoint{.60, .89})
	poly(130, grakPoint{.31, .79}, grakPoint{.47, .79}, grakPoint{.47, .82}, grakPoint{.30, .82})
	poly(130, grakPoint{.59, .79}, grakPoint{.73, .79}, grakPoint{.74, .82}, grakPoint{.60, .82})
	// The tunic is a tapered slab, crossed by one strong diagonal strap.
	poly(65, grakPoint{.35, .42}, grakPoint{.64, .42}, grakPoint{.73, .50}, grakPoint{.68, .72}, grakPoint{.35, .72}, grakPoint{.31, .52})
	poly(107, grakPoint{.36, .43}, grakPoint{.59, .43}, grakPoint{.64, .62}, grakPoint{.35, .61}, grakPoint{.32, .52})
	poly(130, grakPoint{.37, .42}, grakPoint{.43, .42}, grakPoint{.64, .69}, grakPoint{.57, .70})
	poly(94, grakPoint{.34, .66}, grakPoint{.69, .66}, grakPoint{.69, .72}, grakPoint{.35, .72})
	poly(178, grakPoint{.52, .665}, grakPoint{.59, .665}, grakPoint{.59, .71}, grakPoint{.52, .71})
	poly(94, grakPoint{.535, .68}, grakPoint{.575, .68}, grakPoint{.575, .695}, grakPoint{.535, .695})
	// His free arm ends in a compact fist; the other hand rests on the tool.
	poly(65, grakPoint{.66, .44}, grakPoint{.78, .47}, grakPoint{.83, .57}, grakPoint{.77, .64}, grakPoint{.67, .62}, grakPoint{.65, .55})
	poly(107, grakPoint{.69, .46}, grakPoint{.76, .48}, grakPoint{.79, .57}, grakPoint{.72, .60}, grakPoint{.68, .54})
	poly(107, grakPoint{.70, .57}, grakPoint{.79, .58}, grakPoint{.79, .66}, grakPoint{.70, .66}, grakPoint{.67, .62})
	poly(107, grakPoint{.32, .45}, grakPoint{.40, .49}, grakPoint{.35, .57}, grakPoint{.27, .59}, grakPoint{.18, .54}, grakPoint{.20, .48})
	poly(65, grakPoint{.25, .53}, grakPoint{.35, .53}, grakPoint{.35, .57}, grakPoint{.27, .59}, grakPoint{.18, .54})
	// A long wooden shaft and a chamfered stone head read as a builder's mallet.
	poly(94, grakPoint{.19, .42}, grakPoint{.225, .42}, grakPoint{.225, .89}, grakPoint{.19, .89})
	poly(180, grakPoint{.19, .47}, grakPoint{.20, .47}, grakPoint{.20, .86}, grakPoint{.19, .86})
	poly(238, grakPoint{.11, .36}, grakPoint{.28, .36}, grakPoint{.31, .39}, grakPoint{.30, .46}, grakPoint{.10, .46}, grakPoint{.09, .40})
	poly(244, grakPoint{.11, .36}, grakPoint{.28, .36}, grakPoint{.29, .39}, grakPoint{.10, .39})
	poly(236, grakPoint{.25, .39}, grakPoint{.31, .39}, grakPoint{.30, .46}, grakPoint{.25, .46})
	poly(107, grakPoint{.18, .48}, grakPoint{.25, .48}, grakPoint{.26, .53}, grakPoint{.23, .55}, grakPoint{.17, .53})
	// The collar and neck sit under a huge angular jaw, not a human helmet.
	poly(95, grakPoint{.33, .41}, grakPoint{.43, .38}, grakPoint{.63, .39}, grakPoint{.72, .44}, grakPoint{.58, .49}, grakPoint{.41, .46})
	poly(65, grakPoint{.42, .34}, grakPoint{.66, .34}, grakPoint{.63, .44}, grakPoint{.45, .44})
	// Blade ears break the silhouette; a missing corner makes them uneven.
	poly(107, grakPoint{.12, .16}, grakPoint{.36, .21}, grakPoint{.40, .30}, grakPoint{.25, .28}, grakPoint{.20, .23})
	poly(65, grakPoint{.18, .19}, grakPoint{.34, .23}, grakPoint{.34, .27}, grakPoint{.25, .25})
	poly(107, grakPoint{.66, .20}, grakPoint{.87, .12}, grakPoint{.83, .20}, grakPoint{.85, .215}, grakPoint{.74, .29}, grakPoint{.66, .28})
	poly(65, grakPoint{.70, .23}, grakPoint{.82, .17}, grakPoint{.77, .24}, grakPoint{.71, .27})
	poly(107, grakPoint{.33, .17}, grakPoint{.63, .15}, grakPoint{.73, .22}, grakPoint{.76, .32}, grakPoint{.67, .40}, grakPoint{.42, .41}, grakPoint{.29, .32}, grakPoint{.29, .23})
	// Three large facial planes: lit forehead, dark temple, broad lower jaw.
	poly(150, grakPoint{.34, .18}, grakPoint{.61, .17}, grakPoint{.67, .22}, grakPoint{.36, .23}, grakPoint{.31, .26})
	poly(65, grakPoint{.65, .22}, grakPoint{.73, .23}, grakPoint{.76, .32}, grakPoint{.67, .40}, grakPoint{.62, .34})
	poly(65, grakPoint{.31, .32}, grakPoint{.44, .36}, grakPoint{.67, .35}, grakPoint{.67, .40}, grakPoint{.42, .41})
	poly(107, grakPoint{.36, .34}, grakPoint{.63, .34}, grakPoint{.67, .37}, grakPoint{.59, .39}, grakPoint{.42, .38})
	// Heavy eyebrows and a projecting nose set a wary, stubborn expression.
	poly(58, grakPoint{.35, .25}, grakPoint{.48, .27}, grakPoint{.48, .295}, grakPoint{.35, .28})
	poly(58, grakPoint{.57, .27}, grakPoint{.67, .245}, grakPoint{.69, .28}, grakPoint{.57, .295})
	poly(150, grakPoint{.53, .24}, grakPoint{.59, .25}, grakPoint{.69, .31}, grakPoint{.56, .32}, grakPoint{.51, .29})
	poly(65, grakPoint{.51, .30}, grakPoint{.68, .31}, grakPoint{.64, .33}, grakPoint{.53, .33})
	poly(58, grakPoint{.40, .345}, grakPoint{.66, .345}, grakPoint{.65, .36}, grakPoint{.42, .36})
	// One short tusk and one long tusk are the character's simplest signature.
	poly(180, grakPoint{.405, .36}, grakPoint{.44, .36}, grakPoint{.415, .32})
	poly(180, grakPoint{.625, .36}, grakPoint{.66, .36}, grakPoint{.66, .295}, grakPoint{.635, .32})
	// Small marks inherit the surface beneath them, retaining solid masses.
	detail := p.detail
	detail(.405, .28, '━', 180)
	detail(.60, .28, '━', 180)
	detail(.34, .30, '╲', 150)
	if h >= 24 {
		detail(.355, .32, '╲', 150)
	} else {
		// Preserve the tusk signature when the filled wedges are sub-cell.
		detail(.43, .35, '▴', 180)
		detail(.65, .34, '▴', 180)
	}
	detail(.72, .615, '│', 65)
	detail(.75, .615, '│', 65)
	detail(.20, .515, '━', 65)
}
