package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"tdef/game"
)

// The overworld is the lair's map: a dark void with lit corridors connecting
// the floors of the lair (the game's levels), rendered Mario-world-map style.
// It is a 45x13 grid — identical to the built-in levels — so the overworld
// reuses the exact playfield scale logic (ComputeScale / GameLayout /
// Layout.block): the largest integer scale 1-4 that fits, centered in the
// playfield region, min frame 62x19, live reflow on resize.
//
// The map is the lair's memory: each floor's pad carries the result of its
// last defense (a check, a scar, a best wave), the corridors run brighter as
// floors are held, the Rotunda's heart unseals into a door once every floor
// is held, and Malgrath himself speaks to Grak at the hub.

const (
	OWW = 45
	OWH = 13
)

// OWDescendFrames is the length of the descent transition (frames).
const OWDescendFrames = 24

// OWUnsealFrames is the length of a floor's unseal cascade (frames).
const OWUnsealFrames = 60

// OWBootFrames is the length of the arrival cinematic, "the lair wakes".
const OWBootFrames = 90

// OWBlastFrames is the length of the heart-unseal shockwave.
const OWBlastFrames = 72

// OWTrailMaxAge is how long a walk-trail cell lingers (frames).
const OWTrailMaxAge = 24

// HeartFloorID is the endgame defense's floor id (the heart chamber).
const HeartFloorID = "heart"

// owGlyph is one static landmark glyph, placed at (dx,dy) grid cells from the
// room centre, in colour c.
type owGlyph struct {
	dx, dy int
	r      rune
	c      int
}

// owNode is one floor of the lair: a room on the map, wired to a game level
// id. X,Y is the room centre in grid coords; PW,PH its size (both odd, so the
// centre is exact). Chrome is the room frame colour; Glyph is the living
// centre landmark; Deco is the static structure around it.
type owNode struct {
	ID     string
	Name   string
	Level  string
	X, Y   int
	PW, PH int
	Chrome int  // room frame colour
	Glyph  rune // animated centre landmark
	Water  bool // Sunken Garden: the room interior is water
	Fog    bool // Unmapped Depths: the room interior is fog
	Deco   []owGlyph
}

// The five floors. The Rotunda is the hub every corridor runs through; the
// Rift is where Grak starts. Order is the lair's "depth": the Rift is the
// mouth, the Unmapped Depths are the far dark.
var owNodes = []owNode{
	{
		ID: "rift", Name: "the Rift", Level: "canyon", X: 6, Y: 10, PW: 3, PH: 5,
		Chrome: 208, Glyph: '◈',
		Deco: []owGlyph{{0, -1, '╎', 208}, {0, 1, '╎', 208}},
	},
	{
		ID: "rotunda", Name: "the Rotunda", Level: "hub", X: 22, Y: 6, PW: 7, PH: 5,
		Chrome: 178, Glyph: '♥',
		Deco: []owGlyph{
			{-2, -1, '┃', 178}, {0, -1, '┃', 178}, {2, -1, '┃', 178},
			{-1, -1, '°', 236}, {1, -1, '°', 236},
			{-2, 0, '°', 236}, {2, 0, '°', 236},
			{-2, 1, '┃', 178}, {0, 1, '┃', 178}, {2, 1, '┃', 178},
			{-1, 1, '°', 236}, {1, 1, '°', 236},
		},
	},
	{
		ID: "halls", Name: "the Long Halls", Level: "winding", X: 35, Y: 3, PW: 5, PH: 3,
		Chrome: 110, Glyph: '≡',
		Deco: []owGlyph{{-1, 0, '|', 110}, {1, 0, '|', 110}},
	},
	{
		ID: "garden", Name: "the Sunken Garden", Level: "garden", X: 36, Y: 10, PW: 5, PH: 5,
		Chrome: 45, Glyph: 'Ω', Water: true,
		Deco: []owGlyph{
			{-1, -1, '~', 31}, {0, -1, '∧', 45}, {1, -1, '~', 31},
			{-1, 0, '~', 31}, {1, 0, '~', 31},
			{-1, 1, '~', 31}, {0, 1, '~', 31}, {1, 1, '~', 31},
		},
	},
	{
		ID: "depths", Name: "the Unmapped Depths", Level: "maze", X: 9, Y: 2, PW: 5, PH: 5,
		Chrome: 98, Glyph: '?', Fog: true,
		Deco: []owGlyph{
			{-1, -1, '░', 60}, {0, -1, 'Ø', 99}, {1, -1, '░', 60},
			{-1, 0, '░', 60}, {1, 0, '░', 60},
			{-1, 1, '░', 60}, {0, 1, '░', 60}, {1, 1, '░', 60},
		},
	},
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
	{{6, 10}, {6, 6}, {22, 6}},   // Rift -> Rotunda
	{{22, 6}, {35, 6}, {35, 3}},  // Rotunda -> Long Halls
	{{22, 6}, {36, 6}, {36, 10}}, // Rotunda -> Sunken Garden
	{{22, 6}, {9, 6}, {9, 2}},    // Rotunda -> Unmapped Depths
}

// owRouteFloors maps each corridor (owRouteCells order) to the floor it
// serves, so a route carries the state of its floor.
var owRouteFloors = []string{"rift", "halls", "garden", "depths"}

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

// owRouteCells expands each route into an ordered cell list, oriented from the
// Rotunda hub outward, so an energy pulse can travel the corridors.
var owRouteCells = buildOWRouteCells()

func buildOWRouteCells() [][]game.Vec {
	const hubX, hubY = 22, 6
	out := make([][]game.Vec, len(owRoutes))
	for i, route := range owRoutes {
		r := make([][]int, len(route))
		copy(r, route)
		if r[0][0] != hubX || r[0][1] != hubY {
			for a, b := 0, len(r)-1; a < b; a, b = a+1, b-1 {
				r[a], r[b] = r[b], r[a]
			}
		}
		cells := []game.Vec{{X: r[0][0], Y: r[0][1]}}
		for j := 1; j < len(r); j++ {
			x0, y0 := r[j-1][0], r[j-1][1]
			x1, y1 := r[j][0], r[j][1]
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
			for x != x1 || y != y1 {
				x += sx
				y += sy
				cells = append(cells, game.Vec{X: x, Y: y})
			}
		}
		out[i] = cells
	}
	return out
}

// hoardGlints are fixed gold sparkles in the void near the hoard (the Rotunda
// core and the Rift), the fallen heroes' loot. Purely atmospheric.
var hoardGlints = [][]int{
	{15, 3}, {29, 3}, {15, 9}, {29, 9}, {22, 11}, {22, 1}, {4, 10}, {41, 10},
}

// owProcessionHeroes are the Long Halls' ghosts: a line of fallen heroes
// (their glyphs and colours, from the guild's roster) that marches the
// corridor when the halls are open.
var owProcessionHeroes = [3]struct {
	g rune
	c int
}{
	{'t', 251}, // a paladin
	{'g', 244}, // a mercenary
	{'r', 244}, // a rogue
}

// owDragonVoice is Malgrath at the Rotunda, tiered by the hearts held at the
// current renown (0..4). He speaks slowly; the line turns every 6s.
var owDragonVoice = [5][]string{
	{
		"Malgrath: …the guild again? Good. I was already dying, Grak.",
		"Malgrath: take the hoard. They came for it — take theirs first.",
	},
	{"Malgrath: my breath is short today. Hold the halls."},
	{"Malgrath: the heart steadies. The hoard is nearly safe."},
	{"Malgrath: I remember when we were the ones they feared."},
	{"Malgrath: before you came I was dying. Now I am only resting."},
}

const owDragonVoiceDone = "Malgrath: they will come again, Grak. We will be ready."

// owFloorFlavor is the lair's voice on each floor, by state: 0 never held,
// 1 the last defense broke, 2 held.
var owFloorFlavor = map[string][3]string{
	"rift": {
		"the rift still burns where the first expedition broke in",
		"the rift smokes — the guild reached the heart here",
		"the rift is held. no blood on the basalt",
	},
	"rotunda": {
		"the heart beats slow in the rotunda",
		"the rotunda still aches where they struck",
		"the heart beats steady in the rotunda",
	},
	"halls": {
		"stone corridors ring hollow with old echoes",
		"the long halls still smell of blood",
		"the halls ring empty. the guild lost its way",
	},
	"garden": {
		"the water remembers a temple's prayers",
		"the garden is stained; the water runs red at the edges",
		"the garden keeps. the temple stones are quiet",
	},
	"depths": {
		"the depths are uncharted; the dark maps itself",
		"something in the depths is still hunting",
		"the depths let you through. barely",
	},
	HeartFloorID: {
		"the heart chamber. the final expedition waits",
		"the guild's champion broke through the heart",
		"the heart is held. the guild's hunt fails",
	},
}

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

// OWRec is the lair's memory of one floor at the current renown.
type OWRec struct {
	Cleared  bool
	BestWave int
	LastWave int
	LastWon  bool
}

// OWReturnFX is the lair's reaction to a defense that just ended: which
// floor to answer on (the heart renders on the Rotunda) and how it went.
type OWReturnFX struct {
	Floor string
	Won   bool
}

// OWState is the overworld view state. Cursor is Grak's position (grid
// coords); Unlocked holds the floor ids that are open (a floor not in
// Unlocked is rendered sealed, unless it is mid-unseal in Unsealing).
//
// The Records/Scores/Hearts/... fields are the lair's memory, rebuilt by the
// TUI from hiscore.Lair; the TTL fields are frame-counted transitions,
// decremented by the TUI's per-frame tick. RenderOverworld stays a pure
// function of (w, h, state, frame).
type OWState struct {
	Cursor    game.Vec
	Unlocked  map[string]bool
	Msg       string // transient line (e.g. "sealed")
	RevealAll bool   // look-dev: every floor at full brightness

	Diff      int              // index into Difficulties (the renown)
	Records   map[string]OWRec // floor id -> record at the current renown
	Scores    map[string]int   // floor id -> best score
	Hearts    int              // dragon hearts at this renown (0..4)
	BossReady bool             // the heart has unsealed
	BossDone  bool             // the heart is held
	Tokens    int              // relic tokens
	Seed      string           // the Depths' maze seed ("" = uncharted)

	BootTTL    int    // frames left on the arrival cinematic (0 = not playing)
	BlastTTL   int    // frames left on the heart-unseal shockwave
	Descending string // floor id under the descent transition
	DescendTTL int
	ReturnMsg  string // the result banner, after a defense
	ReturnTTL  int
	ReturnFX   OWReturnFX     // the result beat, while ReturnTTL > 0
	Unsealing  map[string]int // floor id -> frames until it opens
	Trail      []game.Vec     // walk trail, newest first (cap 3)
	TrailAge   []int          // frames left on each trail cell
	RelicMenu  bool           // the relic-spend line is active
	FirstRun   bool           // show the one-time hint line

	BonusGold  int // relic bonuses for the next defense
	BonusTower bool
	BonusLives int
}

// NewOWState is the starting look: Grak at the Rift, the Rift and the
// Rotunda open, the rest of the lair still sealed, no memory.
func NewOWState() OWState {
	return OWState{
		Cursor:    game.Vec{X: 6, Y: 10},
		Unlocked:  map[string]bool{"rift": true, "rotunda": true},
		Unsealing: map[string]int{},
		Records:   map[string]OWRec{},
		Scores:    map[string]int{},
		Diff:      1, // normal
		FirstRun:  true,
	}
}

// PushTrail records the cell Grak is leaving, keeping the newest three steps.
func (st *OWState) PushTrail(v game.Vec) {
	st.Trail = append([]game.Vec{v}, st.Trail...)
	st.TrailAge = append([]int{OWTrailMaxAge}, st.TrailAge...)
	if len(st.Trail) > 3 {
		st.Trail = st.Trail[:3]
		st.TrailAge = st.TrailAge[:3]
	}
}

// OWFloor is the public identity of a floor: enough for the TUI to launch the
// matching level, name it, and stand Grak on its centre.
type OWFloor struct {
	ID     string
	Name   string
	Level  string
	Center game.Vec
}

// OWFloorAt reports which floor pad (if any) covers grid cell (x,y).
func OWFloorAt(x, y int) (OWFloor, bool) {
	for i := range owNodes {
		if owNodes[i].contains(x, y) {
			n := &owNodes[i]
			return OWFloor{ID: n.ID, Name: n.Name, Level: n.Level, Center: game.Vec{X: n.X, Y: n.Y}}, true
		}
	}
	return OWFloor{}, false
}

// OWFloorOf returns a floor's public identity by id.
func OWFloorOf(id string) (OWFloor, bool) {
	n := owNodeByID(id)
	if n == nil {
		return OWFloor{}, false
	}
	return OWFloor{ID: n.ID, Name: n.Name, Level: n.Level, Center: game.Vec{X: n.X, Y: n.Y}}, true
}

// OWFloorName is a floor's display name; the boss door reads "the Heart".
func OWFloorName(id string) string {
	if id == HeartFloorID {
		return "the Heart"
	}
	n := owNodeByID(id)
	if n == nil {
		return id
	}
	return n.Name
}

// glow is the 0..1 falloff of a light centred on (cx,cy) sampled at (x,y);
// zero beyond radius. Squared distance, so no per-cell sqrt.
func glow(x, y, cx, cy, radius int) float64 {
	dx, dy := x-cx, y-cy
	d2 := dx*dx + dy*dy
	r2 := radius * radius
	if d2 >= r2 {
		return 0
	}
	return float64(r2-d2) / float64(r2)
}

// owHash is a small deterministic 0..255 hash for ambient texture.
func owHash(x, y, t int) int {
	h := x*73856093 ^ y*19349663 ^ t*83492791
	h = (h ^ (h >> 13)) * 1274126177
	return (h ^ (h >> 16)) & 0xff
}

// lerp steps from a to b by t (clamped 0..1).
func lerp(a, b int, t float64) int {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return a + int(float64(b-a)*t)
}

// overworldFooter is the control hint group for the screen's bottom border.
func overworldFooter() []fseg {
	return []fseg{
		{key: "wasd", text: " walk"},
		{key: "⏎", text: " descend"},
		{key: "tab", text: " renown"},
		{key: "esc", text: " back"},
	}
}

// RenderOverworld draws the lair map for a tw×th terminal. It is a pure
// function of (w, h, state, frame, palette) so it is deterministic and
// unit-testable; frame (the 30fps counter) drives the ambient animation.
func RenderOverworld(w, h int, st OWState, frame int, pal Colors) *Frame {
	l := GameLayout(OWW, OWH, w, h)
	title := "THE LAIR"
	if st.BossDone {
		title = "THE LAIR — HELD"
	}
	f := screenBox(w, h, title, overworldFooter(), frame%30 < 22, pal)

	// 1. The void: the dark of the lair, lit faintly from the Rotunda core —
	//    a radial fall-off to near-black at the edges, warm glow pools around
	//    the fire floors (Rift, Rotunda) and cool around the cold ones
	//    (Garden, Depths), with rock grain and outcrops so it reads as a cave.
	//    A held lair glows a shade warmer.
	glowR := 8
	if st.BossDone {
		glowR = 10
	}
	for y := 0; y < OWH; y++ {
		for x := 0; x < OWW; x++ {
			dx := float64(x) - 22
			dy := float64(y) - 6
			d := math.Sqrt(dx*dx+dy*dy) / 24.0
			if d > 1 {
				d = 1
			}
			bg := 235 - int(d*4) // 235 at the core, 231 at the far edge
			if st.BossDone {
				bg = 237 - int(d*4)
			}
			// Warm/cool temperature: glow pools around the fire and cold floors.
			warm, cool := 0.0, 0.0
			for _, c := range [][2]int{{6, 10}, {22, 6}} {
				if g := glow(x, y, c[0], c[1], glowR); g > warm {
					warm = g
				}
			}
			for _, c := range [][2]int{{36, 10}, {9, 2}} {
				if g := glow(x, y, c[0], c[1], glowR); g > cool {
					cool = g
				}
			}
			switch t := warm - cool; {
			case t > 0.55:
				bg = 95 // warm fire glow
			case t > 0.28:
				bg = 88
			case t < -0.55:
				bg = 24 // cold deep glow
			case t < -0.28:
				bg = 17
			}
			c := Cell{R: ' ', BG: bg}
			if owOutcropSet[game.Vec{X: x, Y: y}] {
				c.R, c.FG = '▒', bg+2
				// A lighter cap on top when the cell above is open void, so
				// the formation reads as a ridge.
				if y > 0 && !OWWalkable(x, y-1) && !owOutcropSet[game.Vec{X: x, Y: y - 1}] {
					capx, capy := l.center(x, y-1)
					if cc := f.C[capy*f.W+capx]; cc.R == ' ' || cc.R == '·' || cc.R == ':' {
						f.C[capy*f.W+capx] = Cell{R: '░', FG: bg + 3, BG: cc.BG}
					}
				}
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
			// Grain twinkle: a speck of dust catches the light for a moment.
			if c.R != ' ' && owHash(x, y, frame/16)%23 == 0 {
				c.FG++
			}
			l.block(f, x, y, c)
		}
	}

	// 2. Corridors: a lit stone band with a centre trail, drawn per route so
	//    each road carries the state of the floor it serves: sealed routes
	//    run dark, held routes run bright with a green trail, and a broken
	//    defense leaves scars. Energy pulses travel out from the heart and
	//    back — the lair's circulation.
	for i, cells := range owRouteCells {
		fl := owRouteFloors[i]
		band := 238
		if owRouteSealed(st, fl) {
			band = 236
		}
		for _, v := range cells {
			l.block(f, v.X, v.Y, Cell{R: ' ', FG: 0, BG: band})
		}
		pulse, trail := owRouteColors(st, fl)
		for _, v := range cells {
			cx, cy := l.center(v.X, v.Y)
			f.Set(cx, cy, Cell{R: '·', FG: trail, BG: band})
		}
		if len(cells) == 0 {
			continue
		}
		// Inbound first (dim embers returning to the heart), so the
		// outbound head rides on top.
		half := len(cells) / 2
		head := (frame/2 + half) % len(cells)
		for k := 0; k < 3; k++ {
			idx := head - k
			if idx < 0 {
				idx += len(cells)
			}
			cx, cy := l.center(cells[idx].X, cells[idx].Y)
			f.Set(cx, cy, Cell{R: '·', FG: [3]int{124, 94, 88}[k], BG: band})
		}
		speed := 2
		if st.BossDone {
			speed = 1
		}
		head = (frame / speed) % len(cells)
		for k := 0; k < 3; k++ {
			idx := head - k
			if idx < 0 {
				idx += len(cells)
			}
			cx, cy := l.center(cells[idx].X, cells[idx].Y)
			f.Set(cx, cy, Cell{R: '·', FG: pulse[k], BG: band, Bold: k == 0})
		}
		// Scars where the last defense on this floor broke.
		if rec := st.Records[fl]; !rec.LastWon && rec.LastWave > 0 {
			for j, v := range cells {
				if j%6 != 3 {
					continue
				}
				cx, cy := l.center(v.X, v.Y)
				if f.C[cy*f.W+cx].R == '·' {
					f.Set(cx, cy, Cell{R: '✕', FG: 166, BG: band})
				}
			}
		}
	}

	// 2b. Hoard glints: gold sparkles twinkle in the void near the hoard
	//     (twinkling faster once the lair is held).
	glintMod := 7
	if st.BossDone {
		glintMod = 4
	}
	for _, g := range hoardGlints {
		x, y := g[0], g[1]
		if OWWalkable(x, y) {
			continue // only in the void, not on a road or floor
		}
		if (frame/4+owHash(x, y, 0))%glintMod != 0 {
			continue
		}
		cx, cy := l.center(x, y)
		f.Set(cx, cy, Cell{R: '✦', FG: 220, BG: 0, Bold: true})
	}

	// 3. Node pads + landmarks + name tags (with result badges). The depths'
	//    fog thins as the built-in floors are held; an unsealing room wakes
	//    over OWUnsealFrames.
	fogLift := 0
	for _, id := range []string{"rift", "rotunda", "halls", "garden"} {
		if st.Records[id].Cleared {
			fogLift++
		}
	}
	// The return flash: as Grak comes back from a defense, the room he
	// answers on blazes toward the result colour over the first 30 frames.
	returnFlash := 0.0
	returnFlashC := 220
	if st.ReturnTTL > 0 && st.ReturnFX.Floor != "" {
		q := 1 - float64(st.ReturnTTL)/180
		returnFlash = 1 - q*6
		if returnFlash < 0 {
			returnFlash = 0
		}
		if !st.ReturnFX.Won {
			returnFlashC = 167
		}
	}
	for i := range owNodes {
		n := &owNodes[i]
		v := owPadView{
			status:  owNodeStatus(st, n.ID),
			rec:     st.Records[n.ID],
			fogLift: fogLift,
			cursor:  st.Cursor,
		}
		if ttl, ok := st.Unsealing[n.ID]; ok {
			v.status = OWOpen
			v.unseal = 1 - float64(ttl)/float64(OWUnsealFrames)
		}
		if n.ID == "rotunda" {
			v.boss = st.BossReady
			v.done = st.BossDone
		}
		if n.ID == st.ReturnFX.Floor && returnFlash > 0 {
			v.flash = returnFlash
			v.flashC = returnFlashC
		}
		drawOWPad(f, l, n, v, frame)
		drawOWLabel(f, l, n, v)
	}

	// 4. The Long Halls' echo: a procession of fallen heroes marches the
	//    corridor (drawn after the corridor pulses, so it rides on top).
	drawOWProcession(f, l, st, frame)

	// 4b. The return ring: a pulse of light runs the floor's route as Grak
	//     answers back on the map.
	drawOWReturnRing(f, l, st)

	// 5. Grak, at the cursor, with a fading trail behind his last step.
	drawOWPlayer(f, l, st, frame)

	// 6. The descent: the floor Grak is entering floods in from its centre.
	drawOWDescent(f, l, st)

	// 6b. The heart-unseal blast: a shockwave runs out from the heart across
	//     the whole map.
	drawOWBlast(f, l, st, frame)

	// 7. The arrival: while the lair wakes, the unlit dark masks the map
	//    until the light from the heart sweeps out across it.
	drawOWBoot(f, l, st)

	// 8. The lair's voice line (row under the top border; the waking
	//    narrates itself while the boot plays).
	drawOWVoice(f, w, st, frame)

	// 9. The bottom chrome: expedition ledger, renown + hearts + relics,
	//    and the context line (seed, relics, pending bonuses).
	drawOWChromeRows(f, st)
	return f
}

// drawOWBoot is the arrival cinematic, "the lair wakes": the Rotunda (the
// dragon's heart) is the ignition source and stays lit, while a light front
// sweeps out from it — the unlit void is masked as dark until the front
// passes. A pure function of (st, frame): progress comes from st.BootTTL.
// At BootTTL <= 1 the sweep has crossed the far corner (22.8 grid units),
// so nothing is masked and the frame is the steady state (the seam).
func drawOWBoot(f *Frame, l Layout, st OWState) {
	if st.BootTTL <= 1 {
		return
	}
	b := OWBootFrames - st.BootTTL // frames into the waking
	// The sweep leaves the heart at frame 15 and crosses 27 grid units
	// (just past the far corner) by the end, so the far rooms wake during
	// the "light runs the corridors" line.
	r := 0.0
	if b > 15 {
		r = float64(b-15) * 27.0 / 75.0
	}
	rot := owNodeByID("rotunda")
	for y := 0; y < OWH; y++ {
		for x := 0; x < OWW; x++ {
			if rot != nil && rot.contains(x, y) {
				continue // the ignition source stays lit
			}
			dx, dy := float64(x)-22, float64(y)-6
			if math.Sqrt(dx*dx+dy*dy) <= r {
				continue // the light has reached this cell
			}
			l.block(f, x, y, Cell{R: ' ', BG: 233})
		}
	}
	if r > 0 && r <= 24 {
		for y := 0; y < OWH; y++ {
			for x := 0; x < OWW; x++ {
				dx, dy := float64(x)-22, float64(y)-6
				if math.Abs(math.Sqrt(dx*dx+dy*dy)-r) < 0.75 {
					cx, cy := l.center(x, y)
					f.Set(cx, cy, Cell{R: '·', FG: 255, BG: 233, Bold: true})
				}
			}
		}
	}
}

// OWStatus is a node's presentation state.
type OWStatus int

const (
	OWSealed  OWStatus = iota // not unlocked: dimmed, locked
	OWOpen                    // unlocked: full art
	OWCurrent                 // unlocked and Grak is on it
)

func owNodeStatus(st OWState, id string) OWStatus {
	if st.Unsealing[id] > 0 {
		return OWOpen
	}
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

// dim halves a 256-colour intensity toward the void for sealed nodes, with a
// floor of 24 so a sealed room still reads as a shape (garden chrome 45 ->
// 24, not 22, which vanished into the dark).
func dim(c int) int {
	if c < 232 {
		if v := c/2 + 2; v > 24 {
			return v
		}
		return 24
	}
	return c - (c-232)/2 - 1
}

// owInterior is a room's dark floor colour.
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

// owPadView is everything a pad needs to know to draw itself in one frame.
type owPadView struct {
	status  OWStatus
	rec     OWRec
	unseal  float64  // 0..1 unseal-cascade progress (0 = none)
	fogLift int      // 0..4 built-in floors held (the depths' reveal)
	boss    bool     // the heart has unsealed (the Rotunda is a door)
	done    bool     // the heart is held (the end state)
	cursor  game.Vec // where Grak stands (so the landmark can avoid him)
	flash   float64  // 0..1 return-flash strength on this pad's chrome
	flashC  int      // the flash target colour (220 held / 167 broke)
}

func owRouteSealed(st OWState, fl string) bool {
	return !st.Unlocked[fl] && st.Unsealing[fl] == 0 && !st.RevealAll
}

// owRouteColors is a corridor's pulse palette and trail colour, by the state
// of the floor it serves.
func owRouteColors(st OWState, fl string) (pulse [3]int, trail int) {
	if st.BossDone {
		return [3]int{220, 214, 144}, 220
	}
	if owRouteSealed(st, fl) {
		return [3]int{100, 124, 148}, 238
	}
	if st.Records[fl].Cleared {
		return [3]int{220, 196, 144}, 114
	}
	return [3]int{214, 178, 136}, 245
}

// drawOWPad renders a floor as a room: a double-line chrome frame in the
// floor's accent colour, a dark textured interior, and the glowing landmark
// at the centre. Sealed rooms are dimmed and wear a blinking lock; an
// unsealing room wakes over the cascade.
func drawOWPad(f *Frame, l Layout, n *owNode, v owPadView, frame int) {
	open := v.status != OWSealed
	accent := n.Chrome
	if !open {
		accent = dim(n.Chrome)
	}
	bgFull := owInterior(n)
	bg := bgFull
	if !open {
		bg = dim(bgFull)
	}
	if v.unseal > 0 {
		// The room is waking: interior and frame brighten as it opens.
		open = true
		bg = lerp(dim(bgFull), bgFull, v.unseal)
		accent = lerp(dim(n.Chrome), n.Chrome, v.unseal)
	}
	if v.flash > 0 {
		// The return flash: the room's frame blazes toward the result colour
		// (gold for a hold, scarlet for a break) as Grak comes back.
		accent = lerp(accent, v.flashC, v.flash)
	}
	x0, y0 := n.X-n.PW/2, n.Y-n.PH/2
	x1, y1 := x0+n.PW-1, y0+n.PH-1

	// Interior: dark, lightly textured (water shimmers and ripples, fog
	// drifts and thins, stone has a faint grain).
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
					// A ripple ring sweeps out from the centre.
					dx, dy := x-n.X, y-n.Y
					d := int(math.Round(math.Sqrt(float64(dx*dx + dy*dy))))
					if d > 0 && d == (frame/10)%6 {
						c.R, c.FG = '≈', 45
					}
				case n.Fog:
					if v.fogLift >= 4 {
						// Fully revealed: the stone shows through.
						if owHash(x*7, y*11, 0)%17 == 0 {
							c.R, c.FG = '·', 236
						}
					} else {
						// Drifting fog, thinning as the lair is held.
						if owHash(x*3+frame/4, y*5-frame/4, 7)%(5-v.fogLift/2) == 0 {
							c.R, c.FG = '░', 60
						}
						if v.fogLift >= 1 && owHash(x*11-frame/3, y*13, 9)%9 == 0 {
							c.R, c.FG = '░', 100
						}
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

	// Deco: the landmark structure (pillars, arch, spire, vortex), each room
	// with its own living ambient (the rift's flicker, the garden's depth
	// gradient, the depths' swirl), dimmed with the room when sealed, dark
	// until the room wakes.
	for i, g := range n.Deco {
		if !OWWalkable(n.X+g.dx, n.Y+g.dy) {
			continue
		}
		gx, gy := l.center(n.X+g.dx, n.Y+g.dy)
		r, c := owDecoCell(n, g, i, frame)
		if !open || v.unseal < 0.5 {
			c = dim(c)
		}
		f.Set(gx, gy, Cell{R: r, FG: c, BG: bg, Bold: open})
	}
	// The room's living details over the interior (embers, hoard glints).
	drawOWRoomAmbient(f, l, n, open, frame, bg)

	// Chrome: a double-line frame around the room in the accent colour.
	fx0, fy0 := l.X(x0), l.Y(y0)
	fx1, fy1 := l.X(x1)+l.Scale-1, l.Y(y1)+l.Scale-1
	drawOWChrome(f, fx0, fy0, fx1, fy1, accent)

	// A sealed room wears a blinking lock on its frame.
	if v.status == OWSealed && frame%40 < 24 {
		lx, ly := l.center(n.X, n.Y-n.PH/2)
		f.Set(lx, ly, Cell{R: '╳', FG: dim(167)})
	}

	// Landmark at the centre — but if Grak is standing on it, the landmark
	// shifts to the first free interior cell so the player never occludes it.
	cx, cy := owLandmarkPos(l, n, v)
	drawOWLandmark(f, cx, cy, n, v, frame, bg)

	// A soft glow just outside the frame, so open rooms read as lit beacons.
	// Clamped to the playfield: the top and bottom pads touch the map edges,
	// and an unclamped glow would spill onto the voice row and the ledger.
	if open {
		for _, p := range [4][2]int{
			{cx, fy0 - 1}, {cx, fy1 + 1}, {fx0 - 1, cy}, {fx1 + 1, cy},
		} {
			if p[0] < l.Ox || p[1] < l.Oy || p[0] >= l.Ox+OWW*l.Scale || p[1] >= l.Oy+OWH*l.Scale {
				continue
			}
			if r := f.C[p[1]*f.W+p[0]].R; r == ' ' || r == 0 {
				f.Put(p[0], p[1], '·', dim(accent), 0)
			}
		}
	}
}

// owDecoCell returns the rune and base colour for a deco glyph, applying the
// room's living ambient: the rift's fissure flickers, the garden's ring steps
// through its depth gradient (with a sunbeam on the peak), and the depths'
// swirl cycles its mists. Static rooms (rotunda, halls) return the glyph
// unchanged.
func owDecoCell(n *owNode, g owGlyph, i int, frame int) (rune, int) {
	r, c := g.r, g.c
	switch n.ID {
	case "rift":
		if c == 208 && owHash(0, i, frame/8)%2 == 0 {
			c = 214 // the fissure flickers
		}
	case "garden":
		switch g.dy {
		case -1:
			if g.r == '∧' {
				c = 45
				if frame%10 < 5 {
					c = 87 // a sunbeam
				}
			} else {
				c = 31
			}
		case 0:
			c = 27
		case 1:
			c = 23
		}
	case "depths":
		switch g.r {
		case '░':
			c = 60 + ((i+frame/8)%4 - 2) // a slow swirl through the mists
		case 'Ø':
			if frame%12 < 6 {
				r, c = 'Ø', 99
			} else {
				r, c = 'ø', 107
			}
		}
	}
	return r, c
}

// drawOWRoomAmbient stamps a room's living details over its interior: the
// rift's rising embers and the rotunda's hoard glints. Texture-guarded so
// they never hide a landmark or a road.
func drawOWRoomAmbient(f *Frame, l Layout, n *owNode, open bool, frame int, bg int) {
	switch n.ID {
	case "rift":
		// Embers rising from the fissure: two columns, climbing and fading.
		for k := 0; k < 2; k++ {
			x := 5 + (k%2)*2
			y := 11 - ((frame/6 + k*3) % 4)
			step := (frame/6 + k*3) % 4
			var fg int
			switch step {
			case 0:
				fg = 214
			case 1:
				fg = 202
			case 2:
				fg = 196
			default:
				continue // the ember has risen out and faded
			}
			if !open {
				fg = dim(fg)
			}
			cx, cy := l.center(x, y)
			if cx < 0 || cy < 0 || cx >= f.W || cy >= f.H {
				continue
			}
			if c := f.C[cy*f.W+cx]; c.R == ' ' || c.R == '·' || c.R == ':' || c.R == '▒' {
				f.Set(cx, cy, Cell{R: '·', FG: fg, BG: c.BG})
			}
		}
	case "rotunda":
		// The hoard: gold glints on either side of the heart.
		for _, p := range [2][2]int{{-1, 0}, {1, 0}} {
			x, y := n.X+p[0], n.Y+p[1]
			if !n.contains(x, y) {
				continue
			}
			c, bold := 220, false
			if owHash(frame/4, x, y)%5 == 0 {
				c, bold = 230, true // a glint
			}
			if !open {
				c = dim(c)
			}
			cx, cy := l.center(x, y)
			if cx < 0 || cy < 0 || cx >= f.W || cy >= f.H {
				continue
			}
			if cc := f.C[cy*f.W+cx]; cc.R == ' ' || cc.R == '·' || cc.R == ':' || cc.R == '▒' {
				f.Set(cx, cy, Cell{R: '·', FG: c, BG: cc.BG, Bold: bold})
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

// owLandmarkGlyph is the centre glyph in the "Grak is here" case: the
// unsealed Rotunda shows the portal, not the heart.
func owLandmarkGlyph(n *owNode, v owPadView) rune {
	if n.Glyph == '♥' && v.boss && !v.done {
		return '◉'
	}
	return n.Glyph
}

// owHalo stamps soft light on the four neighbours of (cx,cy), only over
// plain texture so it never hides a landmark or a road.
func owHalo(f *Frame, cx, cy, c int, on bool) {
	if !on {
		return
	}
	for _, p := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		ox, oy := cx+p[0], cy+p[1]
		if ox < 0 || oy < 0 || ox >= f.W || oy >= f.H {
			continue
		}
		r := f.C[oy*f.W+ox].R
		if r == ' ' || r == '·' || r == ':' || r == '▒' {
			f.C[oy*f.W+ox].R = '·'
			f.C[oy*f.W+ox].FG = c
		}
	}
}

// owLandmarkPos returns the frame pixel for a room's landmark. When Grak is
// standing on the pad centre the landmark shifts to the first free interior
// cell in probe order (up, down, left, right, then the diagonals) so the
// player never occludes it; otherwise it stays at the centre. "Free" means
// strictly interior, not the centre, and not a deco cell. If no cell is free
// the landmark stays at the centre (the player occludes it).
func owLandmarkPos(l Layout, n *owNode, v owPadView) (int, int) {
	centre := game.Vec{X: n.X, Y: n.Y}
	if v.cursor != centre {
		return l.center(n.X, n.Y)
	}
	deco := map[game.Vec]bool{}
	for _, g := range n.Deco {
		deco[game.Vec{X: g.dx, Y: g.dy}] = true
	}
	x0, y0 := n.X-n.PW/2, n.Y-n.PH/2
	x1, y1 := x0+n.PW-1, y0+n.PH-1
	probes := [8][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	for _, p := range probes {
		x, y := n.X+p[0], n.Y+p[1]
		if x < x0 || x > x1 || y < y0 || y > y1 {
			continue // not interior
		}
		if deco[game.Vec{X: p[0], Y: p[1]}] {
			continue // a deco cell
		}
		return l.center(x, y)
	}
	return l.center(n.X, n.Y)
}

// drawOWLandmark stamps the floor's glyph at (cx,cy) with its living
// behaviour: the Rift breathes, the Rotunda's heart beats (or becomes the
// unsealed portal, or holds steady and bright once the heart is held), the
// rest hold steady. An unsealing room keeps its landmark dark, then ignites
// it near the end of the cascade.
func drawOWLandmark(f *Frame, cx, cy int, n *owNode, v owPadView, frame int, bg int) {
	if v.unseal > 0 && v.unseal < 0.5 {
		return
	}
	ignite := v.unseal > 0 && v.unseal < 0.85
	switch {
	case v.status == OWCurrent || ignite:
		f.Set(cx, cy, Cell{R: owLandmarkGlyph(n, v), FG: 255, BG: bg, Bold: true})
	case n.Glyph == '♥':
		if v.boss && !v.done {
			// The unsealed heart: a portal that pulls the light in.
			c := 220
			if frame%20 < 10 {
				c = 226
			}
			f.Set(cx, cy, Cell{R: '◉', FG: c, BG: bg, Bold: true})
			owHalo(f, cx, cy, 220, frame%20 < 10)
		} else if v.done {
			// The held heart: steady, bright, unafraid.
			f.Set(cx, cy, Cell{R: '♥', FG: 255, BG: bg, Bold: frame%12 < 8})
			owHalo(f, cx, cy, 255, frame%12 < 8)
		} else {
			// The dying Dragon's heart: a slow, gapped pulse.
			bold := frame%60 < 34
			c := 196
			if frame%60 < 10 {
				c = 202
			}
			if v.status == OWSealed {
				c, bold = dim(196), false
			}
			f.Set(cx, cy, Cell{R: n.Glyph, FG: c, BG: bg, Bold: bold})
		}
	case n.Glyph == '◈':
		c := 208
		if frame%16 < 8 {
			c = 214
		}
		if v.status == OWSealed {
			c = dim(208)
		}
		f.Set(cx, cy, Cell{R: n.Glyph, FG: c, BG: bg, Bold: true})
		if v.status != OWSealed {
			// The fissure breathes: a halo swells and fades.
			owHalo(f, cx, cy, 214, frame%40 < 24)
		}
	default:
		fg := n.Chrome
		if v.status == OWSealed {
			fg = dim(n.Chrome)
		}
		if n.Fog && v.fogLift >= 4 {
			fg = 189 // the revealed depths
		}
		f.Set(cx, cy, Cell{R: n.Glyph, FG: fg, BG: bg, Bold: v.status != OWSealed})
	}
}

// drawOWLabel writes the node's name tag just below its pad (above, for the
// bottom floors, so it never runs into the frame chrome), with the floor's
// result badge: ✓ held, ✗ best wave on a broken floor, ◉ the unsealed heart.
func drawOWLabel(f *Frame, l Layout, n *owNode, v owPadView) {
	label := n.Name
	if n.ID == "rotunda" && v.boss {
		label = "the Heart"
	}
	fg := n.Chrome
	mark := '◆'
	switch v.status {
	case OWCurrent:
		mark = '▶'
	case OWSealed:
		fg = 240
		mark = '✕'
	case OWOpen:
		if v.unseal > 0 {
			mark = '✦'
		}
	}
	badge, badgeFG := "", fg
	switch {
	case n.ID == "rotunda" && v.boss:
		if v.done {
			badge, badgeFG = " ✓", 255
		} else {
			badge, badgeFG = " ◉", 220
		}
	case v.rec.Cleared:
		badge, badgeFG = " ✓", 114
	case v.rec.BestWave > 0:
		badge, badgeFG = " "+strconv.Itoa(v.rec.BestWave), 167
	}
	s := string(mark) + " " + label + badge
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
	if v.status == OWCurrent {
		markFG = pal_Bright()
	}
	f.Set(x, y, Cell{R: mark, FG: markFG, BG: 233, Bold: v.status != OWSealed})
	putString(f, x+2, y, label, fg, 233, v.status == OWCurrent)
	putString(f, x+2+len(label), y, badge, badgeFG, 233, true)
}

// drawOWPlayer stamps Grak at the cursor: a bold bright marker with a soft
// pulse, and a fading footprint behind his last step.
func drawOWPlayer(f *Frame, l Layout, st OWState, frame int) {
	// The trail: Grak's last three steps, newest first, dimming out.
	for i, v := range st.Trail {
		if st.TrailAge[i] <= 0 {
			continue
		}
		px, py := l.center(v.X, v.Y)
		if px >= 0 && py >= 0 && px < f.W && py < f.H {
			c := f.C[py*f.W+px]
			if c.R == ' ' || c.R == '·' || c.R == ':' || c.R == '▒' {
				fg := 174 + st.TrailAge[i]*38/24
				if fg > 214 {
					fg = 214
				}
				f.C[py*f.W+px] = Cell{R: '·', FG: fg, BG: c.BG}
			}
		}
	}
	// The lantern: warm light breathes across the eight neighbour blocks,
	// tinting every texture subcell — orthogonal neighbours bright, diagonals
	// a shade dimmer — so it never hides a landmark or a road.
	if frame%30 < 20 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				fg := 180
				if dx == 0 || dy == 0 {
					fg = 214
				}
				for sy := 0; sy < l.Scale; sy++ {
					for sx := 0; sx < l.Scale; sx++ {
						px, py := l.X(st.Cursor.X+dx)+sx, l.Y(st.Cursor.Y+dy)+sy
						if px < 0 || py < 0 || px >= f.W || py >= f.H {
							continue
						}
						c := f.C[py*f.W+px]
						if c.R == ' ' || c.R == '·' || c.R == ':' || c.R == '▒' {
							c.FG = fg
							f.C[py*f.W+px] = c
						}
					}
				}
			}
		}
	}
	cx, cy := l.center(st.Cursor.X, st.Cursor.Y)
	pulse := frame%30 < 18
	fg, bg := 255, 236
	if !pulse {
		bg = 235
	}
	f.Set(cx, cy, Cell{R: '@', FG: fg, BG: bg, Bold: true})
}

// drawOWProcession marches the Long Halls' ghosts down their corridor when the
// halls are open: the guild's echo of every defense. A line of three fallen
// heroes creeps along the route once a while, fading in and out at the edges
// of its window.
func drawOWProcession(f *Frame, l Layout, st OWState, frame int) {
	open := st.Unlocked["halls"] || st.RevealAll || st.Unsealing["halls"] > 0
	if !open {
		return
	}
	cells := owRouteCells[1]
	if len(cells) == 0 {
		return
	}
	cyc := frame % 240
	if cyc >= 60 {
		return // the procession is between passes
	}
	base := 10 + (frame/4)%5
	for i, h := range owProcessionHeroes {
		idx := base - i*2
		idx %= len(cells)
		if idx < 0 {
			idx += len(cells)
		}
		v := cells[idx]
		cx, cy := l.center(v.X, v.Y)
		if cx < 0 || cy < 0 || cx >= f.W || cy >= f.H {
			continue
		}
		c := h.c
		if cyc < 12 || cyc > 48 {
			c = dim(c) // the echo fades in and out
		}
		f.Set(cx, cy, Cell{R: h.g, FG: c, Bold: i == 0})
	}
}

// owRouteIndexFor is the corridor (owRouteCells index) that serves a floor.
func owRouteIndexFor(floor string) int {
	for i, fl := range owRouteFloors {
		if fl == floor {
			return i
		}
	}
	return -1
}

// drawOWReturnRing runs a pulse of light down the route of the floor Grak
// just answered on: a bright front with a two-cell tail, over the first 60
// frames of the return banner.
func drawOWReturnRing(f *Frame, l Layout, st OWState) {
	if st.ReturnTTL <= 0 || st.ReturnFX.Floor == "" {
		return
	}
	q := 1 - float64(st.ReturnTTL)/180
	if q*3 >= 1 {
		return // the ring has finished its pass
	}
	ri := owRouteIndexFor(st.ReturnFX.Floor)
	if ri < 0 {
		return
	}
	cells := owRouteCells[ri]
	if len(cells) == 0 {
		return
	}
	front := int(q * 3 * float64(len(cells)-1))
	if front > len(cells)-1 {
		front = len(cells) - 1
	}
	for k, fg := range []int{255, 214, 196} {
		idx := front - k
		if idx < 0 {
			break
		}
		v := cells[idx]
		cx, cy := l.center(v.X, v.Y)
		if cx < 0 || cy < 0 || cx >= f.W || cy >= f.H {
			continue
		}
		f.Set(cx, cy, Cell{R: '·', FG: fg, Bold: k == 0})
	}
}

// drawOWBlast is the heart-unseal shockwave: a bright ring runs out from the
// heart across the whole map, leaving a warm trail behind it, and fades as it
// crosses the far corner. A pure function of (st, frame) via st.BlastTTL.
func drawOWBlast(f *Frame, l Layout, st OWState, frame int) {
	if st.BlastTTL <= 0 {
		return
	}
	p := 1 - float64(st.BlastTTL)/OWBlastFrames
	R := p * 24
	k := 1.0
	if p > 0.85 {
		k = (1 - p) / 0.15
	}
	for y := 0; y < OWH; y++ {
		for x := 0; x < OWW; x++ {
			dx, dy := float64(x)-22, float64(y)-6
			d := math.Sqrt(dx*dx + dy*dy)
			cx, cy := l.center(x, y)
			if cx < 0 || cy < 0 || cx >= f.W || cy >= f.H {
				continue
			}
			if math.Abs(d-R) < 0.8*k {
				f.Set(cx, cy, Cell{R: '·', FG: 255, Bold: true}) // the front
				continue
			}
			if d > R-5*k && d < R-0.8*k {
				c := [3]int{214, 208, 196}[owHash(x, y, frame/2)%3] // the trail
				f.Set(cx, cy, Cell{R: '·', FG: c})
			}
		}
	}
}

// drawOWDescent swallows Grak into the floor he is descending into: the pad
// floods in from its centre with the floor's own light, and the last frames
// blaze solid before the screen changes.
func drawOWDescent(f *Frame, l Layout, st OWState) {
	if st.Descending == "" || st.DescendTTL <= 0 {
		return
	}
	n := owNodeByID(st.Descending)
	if n == nil {
		return
	}
	phase := 1 - float64(st.DescendTTL)/float64(OWDescendFrames)
	if phase < 0 {
		phase = 0
	}
	maxD2 := (n.PW/2)*(n.PW/2) + (n.PH/2)*(n.PH/2)
	x0, y0 := n.X-n.PW/2, n.Y-n.PH/2
	r2 := phase * phase * float64(maxD2)
	for y := y0; y <= y0+n.PH-1; y++ {
		for x := x0; x <= x0+n.PW-1; x++ {
			dx, dy := x-n.X, y-n.Y
			d2 := float64(dx*dx + dy*dy)
			if d2 <= r2 {
				c := Cell{R: '·', FG: dim(n.Chrome), BG: n.Chrome}
				if st.DescendTTL <= 6 {
					c = Cell{R: ' ', FG: 0, BG: n.Chrome}
				}
				l.block(f, x, y, c)
			} else if d2 <= r2+3 {
				// The hot rim: a bright edge just ahead of the flood.
				cx, cy := l.center(x, y)
				if cx < 0 || cy < 0 || cx >= f.W || cy >= f.H {
					continue
				}
				if cc := f.C[cy*f.W+cx]; cc.R == ' ' || cc.R == '·' || cc.R == ':' || cc.R == '▒' {
					f.Set(cx, cy, Cell{R: '·', FG: n.Chrome, BG: cc.BG, Bold: true})
				}
			}
		}
	}
}

// drawOWVoice writes the lair's single voice line in the row under the top
// border. One line, by priority: the result banner, the descent, the relic
// menu, a transient input message, the floor under Grak (Malgrath speaks at
// the Rotunda), then first-run help.
func drawOWVoice(f *Frame, w int, st OWState, frame int) {
	var s string
	var fg int
	switch {
	case st.BootTTL > 15:
		// The waking narrates itself (frames 0-74); from frame 75 on the
		// normal voice takes over, so the boot's last frame is the steady
		// state (the seam).
		b := OWBootFrames - st.BootTTL
		switch {
		case b < 25:
			s, fg = "the lair stirs in the dark", 244
		case b < 50:
			s, fg = "the heart beats — light runs the corridors", 244
		default:
			s, fg = "Malgrath: …grak. the guild still hunts.", 214
		}
	case st.BlastTTL > 0:
		s, fg = "the heart has unsealed — the final expedition stirs", 220
	case st.ReturnTTL > 0 && st.ReturnMsg != "":
		s = st.ReturnMsg
		fg = 244
		switch {
		case strings.Contains(s, "held") || strings.Contains(s, "endures") || strings.Contains(s, "breaks open"):
			fg = 114
		case strings.Contains(s, "broke"):
			fg = 167
		}
	case st.Descending != "":
		s = "descending into " + OWFloorName(st.Descending) + "…"
		fg = 244
	case st.RelicMenu:
		s = fmt.Sprintf("relics %d — 1) +60g · 2) free gunner · 3) +1♥ · esc close", st.Tokens)
		fg = 220
	case st.Msg != "":
		s = st.Msg
		fg = 244
		if strings.Contains(s, "sealed") {
			fg = 174
		}
	default:
		if st.FirstRun {
			// Shown even on a pad: a first-time Grak starts on the Rift, so
			// the old "only when the line is empty" trigger never fired.
			s, fg = "wasd walk the lair · enter descend · tab renown", 240
		} else {
			s, fg = owNodeLine(st, frame)
			if s == "" && st.BossDone {
				s, fg = "Malgrath endures.", 255
			}
		}
	}
	if s == "" {
		return
	}
	putString(f, 2, 1, fitMsg(s, w-4), fg, 0, false)
}

// owNodeLine is the floor under Grak: at the Rotunda, Malgrath speaks (tiered
// by the hearts held); on an unsealed heart, the door's warning; elsewhere,
// the floor's result line and its flavor, turning between the two.
func owNodeLine(st OWState, frame int) (string, int) {
	fl, ok := OWFloorAt(st.Cursor.X, st.Cursor.Y)
	if !ok {
		return "", 0
	}
	id := fl.ID
	rec := st.Records[id]
	switch {
	case id == "rotunda" && st.BossReady:
		if st.BossDone {
			return owDragonVoiceDone, 255
		}
		return "the heart has unsealed — enter to face the final expedition", 220
	case id == "rotunda":
		tier := st.Hearts
		if tier < 0 {
			tier = 0
		}
		if tier > 4 {
			tier = 4
		}
		voice := owDragonVoice[tier]
		return voice[(frame/180)%len(voice)], 214
	}
	best := ""
	if s := st.Scores[id]; s > 0 {
		best = strconv.Itoa(s)
	}
	var result string
	rfg := 240
	switch {
	case rec.Cleared:
		result = owInfoLine(fl.Name, best, "held 20/20")
		rfg = 114
	case rec.BestWave > 0:
		result = owInfoLine(fl.Name, best, "broke at "+strconv.Itoa(rec.BestWave))
		rfg = 167
	}
	state := 0
	switch {
	case rec.Cleared:
		state = 2
	case !rec.LastWon && rec.LastWave > 0:
		state = 1
	}
	flavor := owFloorFlavor[id][state]
	if result != "" && (frame/150)%2 == 0 {
		return result, rfg
	}
	return flavor, 244
}

func owInfoLine(name, best, result string) string {
	if best != "" {
		return name + " · best " + best + " · " + result
	}
	return name + " · " + result
}

// owRun is one colored segment of a chrome line.
type owRun struct {
	s  string
	fg int
	b  bool
}

// putCenteredRuns draws a mixed-color line centered on row y.
func putCenteredRuns(f *Frame, y int, runs []owRun) {
	total := 0
	for _, r := range runs {
		total += len([]rune(r.s))
	}
	if total == 0 {
		return
	}
	x := (f.W - total) / 2
	if x < 1 {
		x = 1
	}
	for _, r := range runs {
		for _, ch := range r.s {
			f.Set(x, y, Cell{R: ch, FG: r.fg, Bold: r.b})
			x++
		}
	}
}

// drawOWChromeRows writes the lair's bottom chrome: the expedition ledger
// (each floor's best result at this renown, and the heart), the renown +
// hearts + relics, and a context line (the Depths' seed, the relic offer,
// pending bonuses). The chrome waits for the waking to end (frame 76.5+).
func drawOWChromeRows(f *Frame, st OWState) {
	if st.BootTTL > 13 {
		return
	}
	h := f.H
	// Row h-4: the expedition ledger. Each floor's id wears the accent of
	// its room, so the ledger reads as the map, not a log. A floor Grak just
	// answered on flashes white while the return flash plays.
	accent := map[string]int{"rift": 208, "halls": 110, "garden": 45, "rotunda": 178}
	ledgerFlash := ""
	if st.ReturnTTL > 0 && st.ReturnFX.Floor != "" {
		if q := 1 - float64(st.ReturnTTL)/180; 1-q*6 > 0 {
			ledgerFlash = st.ReturnFX.Floor
		}
	}
	var rs []owRun
	for _, id := range []string{"rift", "halls", "garden", "rotunda"} {
		rec := st.Records[id]
		fg, b := accent[id], false
		if id == ledgerFlash {
			fg, b = 255, true
		}
		rs = append(rs, owRun{"  " + id + " ", fg, b})
		switch {
		case rec.Cleared:
			rs = append(rs, owRun{"✓" + strconv.Itoa(game.MaxWaves), 114, true})
		case rec.BestWave > 0:
			rs = append(rs, owRun{"✗" + strconv.Itoa(rec.BestWave), 167, false})
		default:
			rs = append(rs, owRun{"·", 245, false})
		}
	}
	hfg, hb := 220, false
	if "heart" == ledgerFlash {
		hfg, hb = 255, true
	}
	rs = append(rs, owRun{"  heart ", hfg, hb})
	switch {
	case st.BossDone:
		rs = append(rs, owRun{"✓", 255, true})
	case st.BossReady:
		rs = append(rs, owRun{"◉", 220, true})
	default:
		rs = append(rs, owRun{"·", 245, false})
	}
	putCenteredRuns(f, h-4, rs)

	// Row h-3: the renown, the dragon hearts, the relic tokens.
	rs = rs[:0]
	for i := range Difficulties {
		name := DiffName(i)
		if i == st.Diff {
			rs = append(rs, owRun{"▸" + name + " ", 255, true})
		} else {
			rs = append(rs, owRun{" " + name + " ", 240, false})
		}
	}
	if st.Hearts > 0 {
		rs = append(rs, owRun{"♥×" + strconv.Itoa(st.Hearts), 214, true})
	}
	if st.Tokens > 0 {
		rs = append(rs, owRun{"  relics " + strconv.Itoa(st.Tokens), 220, false})
	}
	putCenteredRuns(f, h-3, rs)

	// Row h-2: context — the Depths' seed, pending bonuses, the relic offer.
	var ctx string
	fl, onPad := OWFloorAt(st.Cursor.X, st.Cursor.Y)
	switch {
	case onPad && fl.ID == "depths":
		if st.Seed != "" {
			ctx = "seed " + st.Seed + " · ⌫ cut · c clear"
		} else {
			ctx = "no seed — enter rolls the uncharted"
		}
	case st.BonusGold > 0 || st.BonusTower || st.BonusLives > 0:
		var b []string
		if st.BonusGold > 0 {
			b = append(b, "+"+strconv.Itoa(st.BonusGold)+"g")
		}
		if st.BonusTower {
			b = append(b, "free gunner")
		}
		if st.BonusLives > 0 {
			b = append(b, "+"+strconv.Itoa(st.BonusLives)+"♥")
		}
		ctx = "next defense: " + strings.Join(b, " · ")
	case onPad && fl.ID == "rotunda" && st.Tokens > 0 && !st.RelicMenu:
		ctx = "t spend a relic"
	}
	if ctx != "" {
		putCenteredRuns(f, h-2, []owRun{{ctx, 240, false}})
	}
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
