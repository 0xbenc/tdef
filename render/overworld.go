package render

import (
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
	"strconv"
	"strings"

	"github.com/0xbenc/termtd/game"
)

// The overworld is the lair's map: a dark void with lit corridors connecting
// the floors of the lair (the game's levels), rendered Mario-world-map style.
// The lair is wider than a battlefield. Its horizontal viewport follows Grak,
// while the voice and expedition ledger stay fixed. World rendering and mouse
// picking share OverworldLayout; the full cavern never has to fit the terminal.

const (
	OWW = 78
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

// owNode is one floor of the lair: a room on the map, wired to a game level
// id. X,Y is the room centre in grid coords; PW,PH its size (both odd, so the
// centre is exact). Chrome is the architecture colour; Glyph is the living
// centre landmark. Each floor has its own drawing in overworld_scene.go.
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
}

// The five floors. The Rotunda is the hub every corridor runs through; the
// Rotunda is where Grak starts. Order is the lair's "depth": the Rift is the
// mouth, the Unmapped Depths are the far dark.
var owNodes = []owNode{
	{
		ID: "rift", Name: copytext.Text("places.rift.map_name"), Level: "canyon", X: 6, Y: 10, PW: 9, PH: 5,
		Chrome: 208, Glyph: '◈',
	},
	{
		ID: "rotunda", Name: copytext.Text("places.rotunda.map_name"), Level: "hub", X: 22, Y: 6, PW: 15, PH: 7,
		Chrome: 178, Glyph: '♥',
	},
	{
		ID: "halls", Name: copytext.Text("places.halls.map_name"), Level: "winding", X: 63, Y: 3, PW: 15, PH: 5,
		Chrome: 110, Glyph: '≡',
	},
	{
		ID: "garden", Name: copytext.Text("places.garden.map_name"), Level: "garden", X: 65, Y: 9, PW: 15, PH: 5,
		Chrome: 45, Glyph: 'Ω', Water: true,
	},
	{
		ID: "depths", Name: copytext.Text("places.depths.map_name"), Level: "maze", X: 9, Y: 2, PW: 11, PH: 5,
		Chrome: 98, Glyph: '?', Fog: true,
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
	{{6, 10}, {11, 10}, {11, 6}, {22, 6}},                                                      // Rift -> Rotunda
	{{22, 6}, {30, 6}, {30, 5}, {36, 5}, {36, 6}, {45, 6}, {45, 4}, {53, 4}, {53, 3}, {63, 3}}, // Rotunda -> Long Halls
	{{22, 6}, {30, 6}, {30, 5}, {36, 5}, {36, 6}, {45, 6}, {45, 8}, {55, 8}, {55, 9}, {65, 9}}, // Rotunda -> Sunken Garden
	{{22, 6}, {9, 6}, {9, 2}}, // Rotunda -> Unmapped Depths
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
		copytext.Text("overworld.malgrath.tier_0.line_1"),
		copytext.Text("overworld.malgrath.tier_0.line_2"),
	},
	{copytext.Text("overworld.malgrath.tier_1.line_1")},
	{copytext.Text("overworld.malgrath.tier_2.line_1")},
	{copytext.Text("overworld.malgrath.tier_3.line_1")},
	{copytext.Text("overworld.malgrath.tier_4.line_1")},
}

var owDragonVoiceDone = copytext.Text("overworld.ow_dragon_voice_done.malgrath_they_will_come_again_grak_we")

// owFloorFlavor is the lair's voice on each floor, by state: 0 never held,
// 1 the last defense broke, 2 held.
var owFloorFlavor = map[string][3]string{
	"rift": {
		copytext.Text("overworld.rooms.rift.unheld"),
		copytext.Text("overworld.rooms.rift.lost"),
		copytext.Text("overworld.rooms.rift.held"),
	},
	"rotunda": {
		copytext.Text("overworld.rooms.rotunda.unheld"),
		copytext.Text("overworld.rooms.rotunda.lost"),
		copytext.Text("overworld.rooms.rotunda.held"),
	},
	"halls": {
		copytext.Text("overworld.rooms.halls.unheld"),
		copytext.Text("overworld.rooms.halls.lost"),
		copytext.Text("overworld.rooms.halls.held"),
	},
	"garden": {
		copytext.Text("overworld.rooms.garden.unheld"),
		copytext.Text("overworld.rooms.garden.lost"),
		copytext.Text("overworld.rooms.garden.held"),
	},
	"depths": {
		copytext.Text("overworld.rooms.depths.unheld"),
		copytext.Text("overworld.rooms.depths.lost"),
		copytext.Text("overworld.rooms.depths.held"),
	},
	HeartFloorID: {
		copytext.Text("overworld.rooms.heart.unheld"),
		copytext.Text("overworld.rooms.heart.lost"),
		copytext.Text("overworld.rooms.heart.held"),
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

// OWCanWalk applies progression to the static cavern geometry. Room bounds
// take precedence over corridors that continue into a room's interior.
func OWCanWalk(x, y int, st OWState) bool {
	if !OWWalkable(x, y) {
		return false
	}
	if fl, room := OWFloorAt(x, y); room {
		return OWFloorOpen(fl.ID, st)
	}
	return true
}

// OWFloorOpen is the shared access rule, independent of visual previews.
func OWFloorOpen(id string, st OWState) bool {
	return st.Unlocked[id] && st.Unsealing[id] == 0
}

// OWFloorApproach is the corridor cell immediately outside a floor's doorway.
// Browser shortcuts can visit a sealed floor here without entering its pad.
func OWFloorApproach(id string) (game.Vec, bool) {
	n := owNodeByID(id)
	if n == nil {
		return game.Vec{}, false
	}
	for i, cells := range owRouteCells {
		if id != "rotunda" && owRouteFloors[i] != id {
			continue
		}
		for j := 1; j < len(cells); j++ {
			a, b := cells[j-1], cells[j]
			inA, inB := n.contains(a.X, a.Y), n.contains(b.X, b.Y)
			if inA == inB {
				continue
			}
			if inA {
				return b, true
			}
			return a, true
		}
	}
	return game.Vec{}, false
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
// OWMessageKind controls presentation independently of editable wording.
type OWMessageKind int

const (
	OWMessageNeutral OWMessageKind = iota
	OWMessageSuccess
	OWMessageDefeat
	OWMessageLocked
)

func (kind OWMessageKind) color() int {
	switch kind {
	case OWMessageSuccess:
		return 114
	case OWMessageDefeat:
		return 167
	case OWMessageLocked:
		return 174
	default:
		return 244
	}
}

type OWState struct {
	Cursor    game.Vec
	Unlocked  map[string]bool
	Msg       string // transient line (e.g. "sealed")
	MsgKind   OWMessageKind
	RevealAll bool // look-dev: every floor at full brightness

	CameraX   float64 // horizontal focus in world cells, eased by the TUI
	CameraSet bool    // false lets static/headless callers focus directly on Cursor

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
	ReturnKind OWMessageKind
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

// NewOWState starts Grak at the Rotunda, with only the hub open and no memory.
func NewOWState() OWState {
	return OWState{
		Cursor:    game.Vec{X: 22, Y: 6},
		Unlocked:  map[string]bool{"rotunda": true},
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
		return copytext.Text("places.heart.map_name")
	}
	n := owNodeByID(id)
	if n == nil {
		return id
	}
	return n.Name
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
		{key: "wasd", text: copytext.Text("overworld.overworld_footer.walk")},
		{key: "⏎", text: copytext.Text("overworld.overworld_footer.descend")},
		{key: "tab", text: copytext.Text("overworld.overworld_footer.renown")},
		{key: "j", text: copytext.Text("overworld.overworld_footer.journal")},
		{key: "esc", text: copytext.Text("overworld.overworld_footer.back")},
	}
}

// RenderOverworld draws the lair map for a tw×th terminal. It is a pure
// function of (w, h, state, frame, palette) so it is deterministic and
// unit-testable; frame (the 30fps counter) drives the ambient animation.
func RenderOverworld(w, h int, st OWState, frame int, pal Colors) *Frame {
	if w < FrameW || h < ChromeTop+OWH+ChromeBot {
		return RenderTooSmall(w, h, FrameW, ChromeTop+OWH+ChromeBot)
	}
	view := OverworldLayout(w, h, st)
	// Paint in world space first. Cropping afterwards prevents partially
	// visible rooms, labels and cinematics from writing into the screen chrome.
	l := view
	l.Ox = 1
	l.W = OWW*l.Scale + 2
	f := &Frame{W: l.W, H: h, C: make([]Cell, l.W*h)}
	title := copytext.Text("overworld.render_overworld.the_lair")
	if st.BossDone {
		title = copytext.Text("overworld.render_overworld.the_lair_held")
	}
	screen := screenBox(w, h, title, overworldFooter(), frame%30 < 22, pal)

	// 1. The cavern is a composed scene: roof, distant masonry, water,
	// rock silhouettes and the supports beneath the crossing.
	drawOWCavern(f, l, frame)

	// 2. One connected walkway, drawn from the same cells Grak walks.
	// Shared branches are paved once, with turns and junctions intact.
	drawOWWalkway(f, l, st)

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
	}
	// Labels follow every building so a neighbouring room's glow cannot
	// punctuate a name. Grak is drawn later and stays visible on crossing roads.
	for i := range owNodes {
		n := &owNodes[i]
		v := owPadView{status: owNodeStatus(st, n.ID), rec: st.Records[n.ID], cursor: st.Cursor}
		if ttl := st.Unsealing[n.ID]; ttl > 0 {
			v.unseal = 1 - float64(ttl)/OWUnsealFrames
		}
		if n.ID == "rotunda" {
			v.boss, v.done = st.BossReady, st.BossDone
		}
		drawOWLabel(f, l, n, v)
	}

	// 4. The Long Halls' echo: a procession of fallen heroes marches the
	//    walkway (drawn after the paving, so it rides on top).
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
	for y := ChromeTop; y < h-ChromeBot; y++ {
		for x := 1; x < w-1; x++ {
			sx := x - view.Ox + l.Ox
			if sx >= 0 && sx < f.W {
				screen.Set(x, y, f.C[y*f.W+sx])
			}
		}
	}
	// Small cut-edge arrows signal that the cavern continues offscreen.
	if view.Ox < 1 {
		screen.Set(1, ChromeTop+OWH*view.Scale/2, Cell{R: '‹', FG: 245, BG: 233})
	}
	if view.X(OWW) > w-1 {
		screen.Set(w-2, ChromeTop+OWH*view.Scale/2, Cell{R: '›', FG: 245, BG: 233})
	}
	drawOWVoice(screen, w, st, frame)

	// 9. The bottom chrome: expedition ledger, renown + hearts + relics,
	//    and the context line (seed, relics, pending bonuses).
	drawOWChromeRows(screen, st)
	return screen
}

// drawOWBoot is the arrival cinematic, "the lair wakes": the Rotunda (the
// dragon's heart) is the ignition source and stays lit, while a light front
// sweeps out from it — the unlit void is masked as dark until the front
// passes. A pure function of (st, frame): progress comes from st.BootTTL.
// At BootTTL <= 1 the sweep has crossed the far corner of the full world,
// so nothing is masked and the frame is the steady state (the seam).
func drawOWBoot(f *Frame, l Layout, st OWState) {
	if st.BootTTL <= 1 {
		return
	}
	b := OWBootFrames - st.BootTTL // frames into the waking
	// The sweep leaves the heart at frame 15 and crosses the full cavern
	// by the end, so the far rooms wake during
	// the "light runs the corridors" line.
	r := 0.0
	if b > 15 {
		r = float64(b-15) * owSweepRadius() / 75.0
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
	if r > 0 && r <= owSweepRadius() {
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

// indexRGB expands a 256-colour index to (r,g,b) in 0..255: the 16 base
// colours are left as-is (callers here only pass cube/grey indices), the
// 6x6x6 cube, and the 232..255 grey ramp.
func indexRGB(c int) (r, g, b int) {
	if c < 16 {
		return 0, 0, 0
	}
	if c < 232 {
		n := c - 16
		return cubeLevel(n / 36), cubeLevel((n / 6) % 6), cubeLevel(n % 6)
	}
	v := 8 + 10*(c-232)
	return v, v, v
}

func cubeLevel(i int) int {
	switch {
	case i == 0:
		return 0
	case i == 1:
		return 95
	default:
		return 135 + 40*(i-2)
	}
}

// rgbIndex picks the nearest 256-colour to (r,g,b) by squared distance.
func rgbIndex(r, g, b int) int {
	best, bd := 0, 1<<30
	for c := 0; c < 256; c++ {
		cr, cg, cb := indexRGB(c)
		d := (r-cr)*(r-cr) + (g-cg)*(g-cg) + (b-cb)*(b-cb)
		if d < bd {
			bd, best = d, c
		}
	}
	return best
}

// baseColor darkens a 256-colour toward black while keeping its hue (a naive
// index halving drifts hue — e.g. gunner green 46 -> 25, a dark blue). lit is
// the retained lightness, 0..100; it is the "charged" base for a level-2 tower.
func baseColor(c, lit int) int {
	r, g, b := indexRGB(c)
	return rgbIndex(r*lit/100, g*lit/100, b*lit/100)
}

// owInterior is a room's dark floor colour.
func owInterior(n *owNode) int {
	switch {
	case n.Water:
		return 233 // temple stone; water has its own lower basin
	case n.Fog, n.ID == "rotunda", n.ID == "rift":
		return 233 // darkness behind the mine timbers
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

// drawOWPad renders a floor as an illustrated room with a dark interior
// and an interactive landmark
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
					if y < n.Y {
						if owHash(x, y, 0)%11 == 0 {
							c.R, c.FG = '·', 237
						}
						break // the pediment and columns stand above the water
					}
					c.BG = 233
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
					// A few low wisps drift behind the rails, without turning
					// the doorway into a field of square fog tiles.
					if v.fogLift < 4 && y >= n.Y && owHash(x-frame/12, y, 0)%7 == 0 {
						c.R, c.FG = '~', 60
					}
				default:
					if owHash(x*7, y*11, 0)%31 == 0 {
						c.R, c.FG = '·', 236
					}
				}
			}
			base := c
			base.R = ' '
			l.block(f, x, y, base)
			if c.R != ' ' {
				cx, cy := l.center(x, y)
				f.Set(cx, cy, c)
			}
		}
	}

	// The room's living details over the interior (embers, hoard glints).
	drawOWRoomAmbient(f, l, n, open, frame, bg)

	// Architecture: each room has its own silhouette in the accent colour.
	fx0, fy0 := l.X(x0), l.Y(y0)
	fx1, fy1 := l.X(x1)+l.Scale-1, l.Y(y1)+l.Scale-1
	drawOWBuilding(f, l, n, accent, bg, v, frame)

	if v.status == OWSealed {
		if approach, ok := OWFloorApproach(n.ID); ok {
			for _, d := range []game.Vec{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}} {
				door := game.Vec{X: approach.X + d.X, Y: approach.Y + d.Y}
				if !n.contains(door.X, door.Y) {
					continue
				}
				l.block(f, door.X, door.Y, Cell{R: '▒', FG: 240, BG: 234})
				x, y := l.center(door.X, door.Y)
				f.Put(x, y, '╳', 245, 234)
				break
			}
		}
	}

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

// drawOWRoomAmbient stamps a room's living details over its interior: the
// rift's rising embers and the rotunda's hoard glints. Texture-guarded so
// they never hide a landmark or a road.
func drawOWRoomAmbient(f *Frame, l Layout, n *owNode, open bool, frame int, bg int) {
	switch n.ID {
	case "rift":
		// Embers rising from the fissure: two columns, climbing and fading.
		for k := 0; k < 2; k++ {
			x := n.X - 2 + (k%2)*4
			y := n.Y + 1 - ((frame/6 + k*3) % 4)
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
			if c := f.C[cy*f.W+cx]; c.R == ' ' || c.R == '·' || c.R == ':' || c.R == '▒' || strings.ContainsRune("═║╔╗╚╝╠╣╦╩╬", c.R) {
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

// owLandmarkPos keeps a floor's interaction mark beside Grak when he stands
// at its centre. The wider footprints always leave a cell to the left.
func owLandmarkPos(l Layout, n *owNode, v owPadView) (int, int) {
	if v.cursor.X == n.X && v.cursor.Y == n.Y {
		return l.center(n.X-1, n.Y)
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
	if l.Scale == 1 {
		switch n.ID {
		case "rift":
			label = copytext.Text("overworld.draw_owlabel.rift")
		case "depths":
			label = copytext.Text("overworld.draw_owlabel.depths")
		case "halls":
			label = copytext.Text("overworld.draw_owlabel.long_halls")
		case "garden":
			label = copytext.Text("overworld.draw_owlabel.sunken_garden")
		}
	}
	if n.ID == "rotunda" && v.boss {
		label = copytext.Text("overworld.draw_owlabel.the_heart")
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
	if n.ID == "depths" {
		row = n.Y - n.PH/2 // a mine sign over its timbered entrance
	}
	if row >= OWH {
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
			if c.R == ' ' || c.R == '·' || c.R == ':' || c.R == '▒' || strings.ContainsRune("═║╔╗╚╝╠╣╦╩╬", c.R) {
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
						if c.R == ' ' || c.R == '·' || c.R == ':' || c.R == '▒' || strings.ContainsRune("═║╔╗╚╝╠╣╦╩╬", c.R) {
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
	// Grak occupies one terminal character at every scale. A recent step
	// alternates the glyph briefly, then settles back to the idle marker.
	glyph := '@'
	if len(st.TrailAge) > 0 && st.TrailAge[0] > 0 && (frame/4)%2 == 1 {
		glyph = '&'
	}
	f.Set(cx, cy, Cell{R: glyph, FG: fg, BG: bg, Bold: true})
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
		cell := f.C[cy*f.W+cx]
		cell.R, cell.FG, cell.Bold = h.g, c, i == 0
		f.Set(cx, cy, cell)
		if l.Scale >= 2 {
			putOver(f, cx, cy-1, 'o', c)
			putOver(f, cx, cy+1, '┴', dim(c))
		}
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
	R := p * owSweepRadius()
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
	id := st.Descending
	if id == HeartFloorID {
		id = "rotunda"
	}
	n := owNodeByID(id)
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
			s, fg = copytext.Text("overworld.draw_owvoice.the_lair_stirs_in_the_dark"), 244
		case b < 50:
			s, fg = copytext.Text("overworld.draw_owvoice.the_heart_beats_light_runs_the_corridors"), 244
		default:
			s, fg = copytext.Text("overworld.draw_owvoice.malgrath_grak_the_guild_still_hunts"), 214
		}
	case st.BlastTTL > 0:
		s, fg = copytext.Text("overworld.draw_owvoice.the_heart_has_unsealed_the_final_expedition"), 220
	case st.ReturnTTL > 0 && st.ReturnMsg != "":
		s = st.ReturnMsg
		fg = st.ReturnKind.color()
	case st.Descending != "":
		s = copytext.Format("overworld.progress.descending", "floor", OWFloorName(st.Descending))
		fg = 244
	case st.RelicMenu:
		s = copytext.Format("overworld.draw_owvoice.relics_1_60g_2_free_gunner_3", "tokens", strconv.Itoa(st.Tokens), "gold_key", "1", "tower_key", "2", "lives_key", "3", "escape", "esc", "gold", strconv.Itoa(game.RelicGoldBonus), "lives", strconv.Itoa(game.RelicLivesBonus))
		fg = 220
	case st.Msg != "":
		s = st.Msg
		fg = st.MsgKind.color()
	default:
		if st.FirstRun {
			// Shown even on a pad: a first-time Grak starts on the Rotunda, so
			// the old "only when the line is empty" trigger never fired.
			s, fg = copytext.Format("overworld.draw_owvoice.wasd_walk_the_lair_enter_descend_tab", "walk", "wasd", "enter", "enter", "tab", "tab"), 240
		} else {
			s, fg = owNodeLine(st, frame)
			if s == "" && st.BossDone {
				s, fg = copytext.Text("overworld.draw_owvoice.malgrath_endures"), 255
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
		return copytext.Text("overworld.ow_node_line.the_heart_has_unsealed_enter_to_face"), 220
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
		result = owInfoLine(fl.Name, best, copytext.Text("overworld.ow_node_line.held_20_20"))
		rfg = 114
	case rec.BestWave > 0:
		result = owInfoLine(fl.Name, best, copytext.Format("overworld.results.broke_at", "wave", strconv.Itoa(rec.BestWave)))
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
		return copytext.Format("overworld.results.with_best", "floor", name, "score", best, "result", result)
	}
	return copytext.Format("overworld.results.without_best", "floor", name, "result", result)
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
		rs = append(rs, owRun{"  " + copytext.Text("places."+id+".ledger_name") + " ", fg, b})
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
	rs = append(rs, owRun{copytext.Text("overworld.draw_owchrome_rows.heart"), hfg, hb})
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
		rs = append(rs, owRun{copytext.Text("overworld.draw_owchrome_rows.relics") + strconv.Itoa(st.Tokens), 220, false})
	}
	putCenteredRuns(f, h-3, rs)

	// Row h-2: context — the Depths' seed, pending bonuses, the relic offer.
	var ctx string
	fl, onPad := OWFloorAt(st.Cursor.X, st.Cursor.Y)
	switch {
	case onPad && fl.ID == "depths":
		if st.Seed != "" {
			ctx = copytext.Text("overworld.draw_owchrome_rows.seed") + st.Seed + copytext.Format("overworld.draw_owchrome_rows.cut_c_clear", "backspace", "⌫", "clear", "c")
		} else {
			ctx = copytext.Text("overworld.draw_owchrome_rows.no_seed_enter_rolls_the_uncharted")
		}
	case st.BonusGold > 0 || st.BonusTower || st.BonusLives > 0:
		var b []string
		if st.BonusGold > 0 {
			b = append(b, "+"+strconv.Itoa(st.BonusGold)+"g")
		}
		if st.BonusTower {
			b = append(b, copytext.Text("overworld.draw_owchrome_rows.free_gunner"))
		}
		if st.BonusLives > 0 {
			b = append(b, "+"+strconv.Itoa(st.BonusLives)+"♥")
		}
		ctx = copytext.Text("overworld.draw_owchrome_rows.next_defense") + strings.Join(b, " · ")
	case onPad && fl.ID == "rotunda" && st.Tokens > 0 && !st.RelicMenu:
		ctx = copytext.Format("overworld.draw_owchrome_rows.t_spend_a_relic", "target", "t")
	}
	if onPad && fl.ID == "rotunda" && !st.RelicMenu {
		if ctx != "" {
			ctx += " · "
		}
		ctx += copytext.Format("overworld.draw_owchrome_rows.i_opening", "opening", "i")
		if st.BossDone && len([]rune(ctx))+len([]rune(copytext.Format("overworld.draw_owchrome_rows.e_ending", "ending", "e"))) <= f.W-4 {
			ctx += copytext.Format("overworld.draw_owchrome_rows.e_ending", "ending", "e")
		}
		if len([]rune(ctx))+len([]rune(copytext.Format("overworld.draw_owchrome_rows.v_malgrath", "portrait", "v"))) <= f.W-4 {
			ctx += copytext.Format("overworld.draw_owchrome_rows.v_malgrath", "portrait", "v")
		}
		if len([]rune(ctx))+len([]rune(copytext.Format("overworld.draw_owchrome_rows.g_grak", "grak", "g"))) <= f.W-4 {
			ctx += copytext.Format("overworld.draw_owchrome_rows.g_grak", "grak", "g")
		}
	}
	if ctx != "" {
		putCenteredRuns(f, h-2, []owRun{{fitMsg(ctx, f.W-4), 240, false}})
	}
}

// pal_Bright returns the palette's bright white; RenderOverworld's inner
// helpers do not take the palette, so this reads the fixed bright value.
func pal_Bright() int { return 255 }

// OWFloorAtFrame maps a frame-space click to the floor pad under it, returning
// the pad's grid centre and whether a pad was hit. It shares OverworldLayout with
// the renderer, so a click lands exactly where the pad is drawn.
func OWFloorAtFrame(w, h, fx, fy int, st OWState) (game.Vec, bool) {
	l := OverworldLayout(w, h, st)
	if fx < 1 || fx >= w-1 || fy < ChromeTop || fy >= h-ChromeBot {
		return game.Vec{}, false
	}
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
// mouse. The renderer and this share the node table + OverworldLayout, so clicks
// cannot drift from the pads.
func OWRects(w, h int, st OWState) []Rect {
	l := OverworldLayout(w, h, st)
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
