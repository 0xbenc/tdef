package render

import (
	"math"
	"strings"

	"tdef/game"
)

// The overworld is the lair's map: a dark void with lit corridors connecting
// the floors of the lair (the game's levels), rendered Mario-world-map style.
// It is a 45x13 grid — identical to the built-in levels — so the overworld
// reuses the exact playfield scale logic (ComputeScale / GameLayout /
// Layout.block): the largest integer scale 1-4 that fits, centered in the
// playfield region, min frame 62x19, live reflow on resize.

const (
	OWW = 45
	OWH = 13
)

// owNode is one floor of the lair: a room on the map, wired to a game level
// id. X,Y is the room centre in grid coords; PW,PH its size. Chrome is the
// room frame colour; FG is the landmark glyph colour.
type owNode struct {
	ID     string
	Name   string
	Level  string
	X, Y   int
	PW, PH int
	Chrome int // room frame colour
	FG     int // landmark colour
	Glyph  rune
	Water  bool // Sunken Garden: the room is water
	Fog    bool // Unmapped Depths: the room is fog
}

// The five floors. The Rotunda is the hub every corridor runs through; the
// Rift is where Grak starts. Order is the lair's "depth": the Rift is the
// mouth, the Unmapped Depths are the far dark.
var owNodes = []owNode{
	{ID: "rift", Name: "the Rift", Level: "canyon", X: 6, Y: 9, PW: 3, PH: 3, Chrome: 208, FG: 214, Glyph: '◈'},
	{ID: "rotunda", Name: "the Rotunda", Level: "hub", X: 22, Y: 6, PW: 5, PH: 5, Chrome: 178, FG: 196, Glyph: '♥'},
	{ID: "halls", Name: "the Long Halls", Level: "winding", X: 35, Y: 3, PW: 3, PH: 3, Chrome: 110, FG: 111, Glyph: '≡'},
	{ID: "garden", Name: "the Sunken Garden", Level: "garden", X: 36, Y: 10, PW: 5, PH: 3, Chrome: 45, FG: 45, Glyph: 'Ω', Water: true},
	{ID: "depths", Name: "the Unmapped Depths", Level: "maze", X: 9, Y: 2, PW: 5, PH: 3, Chrome: 98, FG: 99, Glyph: '?', Fog: true},
}

func owNodeByID(id string) *owNode {
	for i := range owNodes {
		if owNodes[i].ID == id {
			return &owNodes[i]
		}
	}
	return nil
}

// owRoutes are the corridor centre-lines (grid points, orthogonal segments),
// each running through the Rotunda hub. The corridor cell set is the union of
// every segment; pad cells overdraw them, so routes can run pad-to-pad.
var owRoutes = [][][]int{
	{{6, 9}, {6, 6}, {22, 6}},    // Rift -> Rotunda
	{{22, 6}, {35, 6}, {35, 3}},  // Rotunda -> Long Halls
	{{22, 6}, {36, 6}, {36, 10}}, // Rotunda -> Sunken Garden
	{{22, 6}, {9, 6}, {9, 2}},    // Rotunda -> Unmapped Depths
}

// owCorridor is the set of grid cells a corridor runs through.
var owCorridor = buildOWCorridor()

func buildOWCorridor() map[game.Vec]bool {
	set := map[game.Vec]bool{}
	for _, r := range owRoutes {
		for i := 1; i < len(r); i++ {
			x0, y0 := r[i-1][0], r[i-1][1]
			x1, y1 := r[i][0], r[i][1]
			sx, sy := 0, 0
			switch {
			case x1 > x0:
				sx = 1
			case x1 < x0:
				sx = -1
			case y1 > y0:
				sy = 1
			case y1 < y0:
				sy = -1
			}
			x, y := x0, y0
			set[game.Vec{X: x, Y: y}] = true
			for x != x1 || y != y1 {
				x += sx
				y += sy
				set[game.Vec{X: x, Y: y}] = true
			}
		}
	}
	return set
}

// owOutcrops are fixed rock formations in the void (grid coords), purely
// atmospheric — they never sit on a corridor or a pad. A few hang from the
// ceiling (top rows) or rise from the floor (bottom rows) like stalactites
// and stalagmites.
var owOutcrops = [][]int{
	{2, 1}, {2, 5}, {2, 9}, {12, 0}, {18, 1}, {24, 1}, {28, 2},
	{38, 0}, {42, 3}, {42, 8}, {12, 12}, {18, 11}, {30, 12},
	{38, 11}, {41, 10},
}

var owOutcropSet = func() map[game.Vec]bool {
	m := map[game.Vec]bool{}
	for _, o := range owOutcrops {
		m[game.Vec{X: o[0], Y: o[1]}] = true
	}
	return m
}()

func (n *owNode) contains(x, y int) bool {
	return x >= n.X-n.PW/2 && x <= n.X+n.PW/2 && y >= n.Y-n.PH/2 && y <= n.Y+n.PH/2
}

// OWWalkable reports whether a grid cell can be stood on: a corridor cell or
// a node pad cell. The void between is unlit dark and impassable.
func OWWalkable(x, y int) bool {
	if x < 0 || y < 0 || x >= OWW || y >= OWH {
		return false
	}
	if owCorridor[game.Vec{X: x, Y: y}] {
		return true
	}
	for i := range owNodes {
		if owNodes[i].contains(x, y) {
			return true
		}
	}
	return false
}

// OWState is the overworld view state. Cursor is Grak's position (grid
// coords); Unlocked holds the floor ids that are open. A floor not in
// Unlocked is rendered sealed (dimmed, locked).
type OWState struct {
	Cursor   game.Vec
	Unlocked map[string]bool
	Msg      string // transient line shown in row 1 (e.g. "sealed")
	// RevealAll shows every floor at full brightness regardless of Unlocked —
	// a look-dev view for evaluating all the landmark art at once.
	RevealAll bool
}

// NewOWState is the starting look: Grak at the Rift, the Rift cleared and the
// Rotunda open ahead, the rest of the lair still sealed.
func NewOWState() OWState {
	return OWState{
		Cursor: game.Vec{X: 6, Y: 9},
		Unlocked: map[string]bool{
			"rift":    true,
			"rotunda": true,
		},
	}
}

// OWFloor is the public identity of a floor: enough for the TUI to launch the
// matching level and name it, without exposing the renderer's node table.
type OWFloor struct {
	ID    string
	Name  string
	Level string
}

// OWFloorAt reports which floor pad (if any) covers grid cell (x,y).
func OWFloorAt(x, y int) (OWFloor, bool) {
	for i := range owNodes {
		if owNodes[i].contains(x, y) {
			return OWFloor{ID: owNodes[i].ID, Name: owNodes[i].Name, Level: owNodes[i].Level}, true
		}
	}
	return OWFloor{}, false
}

// owHash is a small deterministic 0..255 hash for ambient texture.
func owHash(x, y, t int) int {
	h := x*73856093 ^ y*19349663 ^ t*83492791
	h = (h ^ (h >> 13)) * 1274126177
	return (h ^ (h >> 16)) & 0xff
}

// overworldFooter is the control hint group for the screen's bottom border.
func overworldFooter() []fseg {
	return []fseg{
		{key: "wasd", text: " walk"},
		{key: "⏎", text: " descend"},
		{key: "r", text: " reveal"},
		{key: "q", text: " quit"},
	}
}

// RenderOverworld draws the lair map for a tw×th terminal. It is a pure
// function of (w, h, state, frame, palette) so it is deterministic and
// unit-testable; frame (the 30fps counter) drives the ambient animation.
func RenderOverworld(w, h int, st OWState, frame int, pal Colors) *Frame {
	l := GameLayout(OWW, OWH, w, h)
	f := screenBox(w, h, "THE LAIR", overworldFooter(), frame%30 < 22, pal)

	// 1. The void: the dark of the lair, lit faintly from the Rotunda core —
	//    a radial fall-off to near-black at the edges, with rock grain and
	//    outcrops so it reads as a cave, not a flat black field.
	for y := 0; y < OWH; y++ {
		for x := 0; x < OWW; x++ {
			dx := float64(x) - 22
			dy := float64(y) - 6
			d := math.Sqrt(dx*dx+dy*dy) / 24.0
			if d > 1 {
				d = 1
			}
			bg := 235 - int(d*4) // 235 at the core, 231 at the far edge
			c := Cell{R: ' ', BG: bg}
			if owOutcropSet[game.Vec{X: x, Y: y}] {
				c.R, c.FG = '▒', bg+2
			} else {
				switch (x*13 + y*7) % 29 {
				case 0:
					c.R, c.FG = '·', bg+1
				case 4:
					c.R, c.FG = ':', bg+1
				case 9:
					c.R, c.FG = '▒', bg+2
				}
			}
			l.block(f, x, y, c)
		}
	}

	// 2. Corridors: a lit stone band with a centre trail, drawn under the
	//    pads so the roads visibly run into every floor.
	for v := range owCorridor {
		l.block(f, v.X, v.Y, Cell{R: ' ', FG: 0, BG: 238})
	}
	for v := range owCorridor {
		cx, cy := l.center(v.X, v.Y)
		f.Set(cx, cy, Cell{R: '·', FG: 245, BG: 238})
	}

	// 3. Node pads + landmarks + name tags.
	for i := range owNodes {
		n := &owNodes[i]
		status := owNodeStatus(st, n.ID)
		drawOWPad(f, l, n, status, frame)
		drawOWLabel(f, l, n, status)
	}

	// 4. Grak, at the cursor.
	drawOWPlayer(f, l, st.Cursor, frame)

	// 5. A transient message in the row under the top border.
	if st.Msg != "" {
		fg := 244
		if strings.Contains(st.Msg, "sealed") {
			fg = 174
		}
		putString(f, 2, 1, fitMsg(st.Msg, w-4), fg, 0, false)
	}
	return f
}

// OWStatus is a node's presentation state.
type OWStatus int

const (
	OWSealed  OWStatus = iota // not unlocked: dimmed, locked
	OWOpen                    // unlocked: full art
	OWCurrent                 // unlocked and Grak is on it
)

func owNodeStatus(st OWState, id string) OWStatus {
	n := owNodeByID(id)
	if n != nil && st.Cursor.X >= n.X-n.PW/2 && st.Cursor.X <= n.X+n.PW/2 && st.Cursor.Y >= n.Y-n.PH/2 && st.Cursor.Y <= n.Y+n.PH/2 {
		if st.Unlocked[id] || st.RevealAll {
			return OWCurrent
		}
	}
	if st.Unlocked[id] || st.RevealAll {
		return OWOpen
	}
	return OWSealed
}

// dim halves a 256-colour intensity toward the void for sealed nodes.
func dim(c int) int {
	if c < 232 {
		return c / 2
	}
	return c - (c-232)/2 - 1
}

// owInterior is a room's dark floor colour (the old solid slab is gone — a
// room is lit only by its chrome frame and its landmark).
func owInterior(n *owNode) int {
	switch {
	case n.Water:
		return 23 // deep water
	case n.Fog:
		return 53 // fog
	default:
		return 234 // dark stone
	}
}

// drawOWPad renders a floor as a room: a double-line chrome frame in the
// floor's accent colour, a dark textured interior, and the glowing landmark
// at the centre. Sealed rooms are dimmed throughout.
func drawOWPad(f *Frame, l Layout, n *owNode, status OWStatus, frame int) {
	open := status != OWSealed
	accent := n.Chrome
	if !open {
		accent = dim(n.Chrome)
	}
	bg := owInterior(n)
	if !open {
		bg = dim(bg)
	}
	x0, y0 := n.X-n.PW/2, n.Y-n.PH/2
	x1, y1 := x0+n.PW-1, y0+n.PH-1

	// Interior: dark, lightly textured (water shimmers, fog drifts, stone
	// has a faint grain).
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x < 0 || y < 0 || x >= OWW || y >= OWH {
				continue
			}
			c := Cell{R: ' ', FG: 0, BG: bg}
			if open {
				switch {
				case n.Water:
					c.R, c.FG = '~', 31
					if owHash(x, y, frame/8)%3 == 0 {
						c.R = '≈'
					}
				case n.Fog:
					if owHash(x*3, y*5, frame/6)%5 == 0 {
						c.R, c.FG = '░', 60
					}
				default:
					if owHash(x*7, y*11, 0)%17 == 0 {
						c.R, c.FG = '·', 236
					}
				}
			}
			l.block(f, x, y, c)
		}
	}

	// Chrome: a double-line frame around the room in the accent colour.
	fx0, fy0 := l.X(x0), l.Y(y0)
	fx1, fy1 := l.X(x1)+l.Scale-1, l.Y(y1)+l.Scale-1
	drawOWChrome(f, fx0, fy0, fx1, fy1, accent)

	// Landmark at the centre.
	cx, cy := l.center(n.X, n.Y)
	drawOWLandmark(f, cx, cy, n, status, frame, bg)

	// A soft glow just outside the frame, so open rooms read as lit beacons.
	if open {
		for _, p := range [4][2]int{
			{cx, fy0 - 1}, {cx, fy1 + 1}, {fx0 - 1, cy}, {fx1 + 1, cy},
		} {
			if p[0] < 0 || p[1] < 0 || p[0] >= f.W || p[1] >= f.H {
				continue
			}
			if r := f.C[p[1]*f.W+p[0]].R; r == ' ' || r == 0 {
				f.Put(p[0], p[1], '·', dim(accent), 0)
			}
		}
	}
}

// drawOWChrome draws a double-line box frame (1 cell thick) around the frame
// rectangle (x0,y0)-(x1,y1) in colour.
func drawOWChrome(f *Frame, x0, y0, x1, y1, color int) {
	if x1-x0 < 1 || y1-y0 < 1 {
		return
	}
	for x := x0 + 1; x < x1; x++ {
		f.Set(x, y0, Cell{R: '═', FG: color})
		f.Set(x, y1, Cell{R: '═', FG: color})
	}
	for y := y0 + 1; y < y1; y++ {
		f.Set(x0, y, Cell{R: '║', FG: color})
		f.Set(x1, y, Cell{R: '║', FG: color})
	}
	f.Set(x0, y0, Cell{R: '╔', FG: color})
	f.Set(x1, y0, Cell{R: '╗', FG: color})
	f.Set(x0, y1, Cell{R: '╚', FG: color})
	f.Set(x1, y1, Cell{R: '╝', FG: color})
}

// drawOWLandmark stamps the floor's glyph at (cx,cy) with its living behaviour:
// the Rift breathes, the Rotunda's heart beats, the rest hold steady.
func drawOWLandmark(f *Frame, cx, cy int, n *owNode, status OWStatus, frame int, bg int) {
	switch {
	case status == OWCurrent:
		f.Set(cx, cy, Cell{R: n.Glyph, FG: 255, BG: bg, Bold: true})
	case n.Glyph == '♥': // the dying Dragon's heart: a slow, gapped pulse
		bold := frame%60 < 34
		c := 196
		if frame%60 < 10 {
			c = 202
		}
		if status == OWSealed {
			c, bold = dim(196), false
		}
		f.Set(cx, cy, Cell{R: n.Glyph, FG: c, BG: bg, Bold: bold})
	case n.Glyph == '◈': // the Rift: a hot, breathing fissure
		c := 208
		if frame%16 < 8 {
			c = 214
		}
		if status == OWSealed {
			c = dim(208)
		}
		f.Set(cx, cy, Cell{R: n.Glyph, FG: c, BG: bg, Bold: true})
	default:
		fg := n.FG
		if status == OWSealed {
			fg = dim(n.FG)
		}
		f.Set(cx, cy, Cell{R: n.Glyph, FG: fg, BG: bg, Bold: status != OWSealed})
	}
}

// drawOWLabel writes the node's name tag just below its pad (above, for the
// bottom floors, so it never runs into the frame chrome).
func drawOWLabel(f *Frame, l Layout, n *owNode, status OWStatus) {
	label := n.Name
	fg := n.Chrome
	mark := '◆'
	switch status {
	case OWCurrent:
		mark = '▶'
	case OWSealed:
		fg = 240
		mark = '✕'
	}
	s := string(mark) + " " + label
	row := n.Y + n.PH/2 + 1
	if row > OWH-2 {
		row = n.Y - n.PH/2 - 1
	}
	if row < 1 {
		row = 1
	}
	cx, _ := l.center(n.X, n.Y)
	y := l.Y(row)
	x := cx - len([]rune(s))/2
	// A dark backing strip keeps the tag legible over corridor/void.
	for i := 0; i < len(s); i++ {
		if x+i >= 0 && x+i < f.W {
			f.Set(x+i, y, Cell{R: ' ', BG: 233})
		}
	}
	markFG := fg
	if status == OWCurrent {
		markFG = pal_Bright()
	}
	f.Set(x, y, Cell{R: mark, FG: markFG, BG: 233, Bold: status != OWSealed})
	putString(f, x+2, y, label, fg, 233, status == OWCurrent)
}

// drawOWPlayer stamps Grak at the cursor: a bold bright marker with a soft
// pulse, so the Last Monster reads clearly against the dark lair.
func drawOWPlayer(f *Frame, l Layout, cur game.Vec, frame int) {
	cx, cy := l.center(cur.X, cur.Y)
	pulse := frame%30 < 18
	fg, bg := 255, 236
	if !pulse {
		bg = 235
	}
	f.Set(cx, cy, Cell{R: '@', FG: fg, BG: bg, Bold: true})
}

// pal_Bright returns the palette's bright white; RenderOverworld's inner
// helpers do not take the palette, so this reads the fixed bright value.
func pal_Bright() int { return 255 }

// OWFloorAtFrame maps a frame-space click to the floor pad under it, returning
// the pad's grid centre and whether a pad was hit. It shares GameLayout with
// the renderer, so a click lands exactly where the pad is drawn.
func OWFloorAtFrame(w, h, fx, fy int) (game.Vec, bool) {
	l := GameLayout(OWW, OWH, w, h)
	if fx < l.Ox || fy < l.Oy || fx >= l.Ox+OWW*l.Scale || fy >= l.Oy+OWH*l.Scale {
		return game.Vec{}, false
	}
	gx, gy := (fx-l.Ox)/l.Scale, (fy-l.Oy)/l.Scale
	fl, ok := OWFloorAt(gx, gy)
	if !ok {
		return game.Vec{}, false
	}
	n := owNodeByID(fl.ID)
	return game.Vec{X: n.X, Y: n.Y}, true
}

// OWRects returns one hit-test rect per node pad (frame space), for the
// mouse. The renderer and this share the node table + GameLayout, so clicks
// cannot drift from the pads.
func OWRects(w, h int) []Rect {
	l := GameLayout(OWW, OWH, w, h)
	out := make([]Rect, len(owNodes))
	for i := range owNodes {
		n := &owNodes[i]
		out[i] = Rect{
			X: l.X(n.X - n.PW/2), Y: l.Y(n.Y - n.PH/2),
			W: n.PW * l.Scale, H: n.PH * l.Scale,
		}
	}
	return out
}
