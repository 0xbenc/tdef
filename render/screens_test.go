package render

import (
	"reflect"
	"strings"
	"testing"

	"tdef/game"
)

// titleAt renders the title at attract-battle internal frame fr
// (0 = the BATTLE text starts to decode), with the visit clock set so the
// boot cinematic and idle wait are already behind us.
func titleAt(w, h, frame, fr int, scores map[string]int) *Frame {
	boot := titleBootLen + titleIdleWait + titleBattleLead + fr
	return RenderTitle(w, h, frame, boot, scores, Palette())
}

// titleStandbyAt renders the idle title (full UI, empty battlefield) at
// visit-boot frame boot.
func titleStandbyAt(w, h, frame, boot int, scores map[string]int) *Frame {
	return RenderTitle(w, h, frame, titleBootLen+boot, scores, Palette())
}

func TestRenderTitleDeterministic(t *testing.T) {
	scores := map[string]int{"winding": 100, "hub": 90}
	for _, fr := range []int{0, 1, 7, 30, 90, 1234} {
		a := RenderTitle(62, 19, fr, fr, scores, Palette())
		b := RenderTitle(62, 19, fr, fr, scores, Palette())
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("title frame %d is not deterministic", fr)
		}
	}
}

// The title must animate: distinct frames must differ somewhere.
func TestRenderTitleAnimates(t *testing.T) {
	if reflect.DeepEqual(RenderTitle(80, 24, 0, 0, nil, Palette()), RenderTitle(80, 24, 4, 4, nil, Palette())) {
		t.Fatal("title boot does not animate (boot 0 and 4 identical)")
	}
	if a, b := titleStandbyAt(80, 24, 0, 10, nil), titleStandbyAt(80, 24, 15, 25, nil); reflect.DeepEqual(a, b) {
		t.Fatal("title standby does not animate (blink/packet frozen)")
	}
}

// At the minimum 62x19 frame the whole card (logo, tagline, demo, roster,
// best, prompt) must be present and nothing clipped.
func TestRenderTitleFitsFrame(t *testing.T) {
	scores := map[string]int{"maze12345678901234567890": 9999999999}
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		w, h := size[0], size[1]
		// Boot cinematic: size check only, it owns the whole frame.
		for boot := 0; boot < titleBootLen; boot++ {
			if f := RenderTitle(w, h, boot, boot, scores, Palette()); f.W != w || f.H != h {
				t.Fatalf("w=%d h=%d boot %d: size = %dx%d", w, h, boot, f.W, f.H)
			}
		}
		// One full attract cycle: the chrome only has to be intact in the
		// standby, the intro and the battle; the effect segments own the
		// frame.
		for local := 0; local < titleAttractCycle; local++ {
			f := RenderTitle(w, h, local, titleBootLen+local, scores, Palette())
			if f.W != w || f.H != h {
				t.Fatalf("w=%d h=%d local %d: size = %dx%d", w, h, local, f.W, f.H)
			}
			if local < titleIdleWait+titleBattleLead+titleOverloadEnd &&
				(rowWidth(f, 0) != w || rowWidth(f, h-1) != w) {
				t.Fatalf("w=%d h=%d local %d: border row clipped", w, h, local)
			}
		}
		text := titleStandbyAt(w, h, 100, 100, scores).Text()
		for _, want := range []string{
			"█████████", // logo top bevel (T)
			"▒▒▒",       // logo bottom bevel
			"— terminal tower defense —",
			"▶", "E", // battlefield spawn / exit
			" towers ", " enemies ",
			"★ best 9999999999", // best line (longest plausible key)
			"[enter] start", "q quit",
			"by 0xbenc", // signature in the bottom border
		} {
			if !strings.Contains(text, want) {
				t.Errorf("w=%d h=%d: standby missing %q:\n%s", w, h, want, text)
			}
		}
		// The standby battlefield must be empty: no BATTLE/WAVE text.
		for _, absent := range []string{"BATTLE", "WAVE"} {
			if strings.Contains(text, absent) {
				t.Errorf("w=%d h=%d: standby shows %q:\n%s", w, h, absent, text)
			}
		}
	}
	// Empty hiscore shows the placeholder.
	if !strings.Contains(titleStandbyAt(62, 19, 0, 10, nil).Text(), "no scores yet") {
		t.Error("title missing 'no scores yet' for an empty table")
	}
}

func TestRenderMenuShowsAllItemsAndSelection(t *testing.T) {
	for sel := 0; sel < len(MenuItems); sel++ {
		f := RenderMenu(62, 19, sel, Palette())
		text := f.Text()
		for _, it := range MenuItems {
			if !strings.Contains(text, it) {
				t.Errorf("sel=%d: menu missing item %q", sel, it)
			}
		}
		if !strings.Contains(text, " ▸ "+MenuItems[sel]) {
			t.Errorf("sel=%d: selection marker not on %q", sel, MenuItems[sel])
		}
	}
}

// MenuRects must land on the rendered item rows.
func TestMenuRectsMatchRenderedRows(t *testing.T) {
	w, h := 62, 19
	f := RenderMenu(w, h, 0, Palette())
	rects := MenuRects(w, h)
	if len(rects) != len(MenuItems) {
		t.Fatalf("rects = %d, want %d", len(rects), len(MenuItems))
	}
	lines := strings.Split(f.Text(), "\n")
	for i, r := range rects {
		row := strings.TrimRight(lines[r.Y], " ")
		if !strings.Contains(row, MenuItems[i]) {
			t.Errorf("item %d rect at row %d does not hit its text: %q", i, r.Y, row)
		}
	}
}

// The menu block must sit centered in the content band at any height.
func TestMenuLayoutCentered(t *testing.T) {
	for _, h := range []int{19, 24, 32, 50} {
		items := menuLayout(h)
		top, bottom := screenOff(h)+1, screenOff(h)+17
		above := items[0] - top
		below := bottom - items[len(items)-1]
		if above < 0 || below < 0 {
			t.Fatalf("h=%d: menu block outside region (above=%d below=%d)", h, above, below)
		}
		if d := above - below; d > 1 || d < -1 {
			t.Errorf("h=%d: menu block not centered (above=%d below=%d)", h, above, below)
		}
	}
}

func TestRenderHelpFits(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}} {
		f := RenderHelp(size[0], size[1], Palette())
		text := f.Text()
		for _, want := range []string{"HELP", "arrows / wasd", "1-7 pick", "esc back"} {
			if !strings.Contains(text, want) {
				t.Errorf("help missing %q:\n%s", want, text)
			}
		}
	}
}

func TestRenderHighScoresEmpty(t *testing.T) {
	f := RenderHighScores(62, 19, 0, map[string]int{}, Palette())
	if !strings.Contains(f.Text(), "no scores yet") {
		t.Errorf("empty hiscores missing placeholder:\n%s", f.Text())
	}
}

func TestRenderHighScoresSorted(t *testing.T) {
	scores := map[string]int{"a": 10, "b": 30, "c": 20}
	lines := strings.Split(RenderHighScores(62, 19, 0, scores, Palette()).Text(), "\n")
	// Names sit in a fixed column (x0+5 of the 37-wide table), so match on
	// that column rather than on the row prefix.
	const w, tableW, nameCol = 62, 37, 6
	x0 := (w - tableW) / 2
	rowOf := func(name string) int {
		for i, ln := range lines {
			if r := []rune(ln); len(r) > x0+nameCol && string(r[x0+nameCol:x0+nameCol+len(name)]) == name {
				return i
			}
		}
		return -1
	}
	rb, rc, ra := rowOf("b"), rowOf("c"), rowOf("a")
	if !(rb >= 0 && rb < rc && rc < ra) {
		t.Errorf("entries not sorted by score desc (rows b=%d c=%d a=%d)", rb, rc, ra)
	}
	text := strings.Join(lines, "\n")
	for _, s := range []string{"30", "20", "10"} {
		if !strings.Contains(text, s) {
			t.Errorf("score %s missing", s)
		}
	}
}

func TestSortScoresDeterministic(t *testing.T) {
	scores := map[string]int{"z": 5, "a": 5, "m": 9}
	got := SortScores(scores)
	want := []ScoreEntry{{"m", 9}, {"a", 5}, {"z", 5}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SortScores = %+v, want %+v", got, want)
	}
}

func TestRenderLevelSelectFits(t *testing.T) {
	names := []string{"canyon", "garden", "hub", "winding"}
	m, err := game.LoadLevel("canyon")
	if err != nil {
		t.Fatal(err)
	}
	v := LSState{Levels: names, Cursor: 1, Diff: 1, Preview: m}
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		f := RenderLevelSelect(v, size[0], size[1], Palette())
		text := f.Text()
		for _, want := range append([]string{lsMazeRow, "easy", "normal", "hard", "enter start"}, names...) {
			if !strings.Contains(text, want) {
				t.Errorf("%dx%d: level select missing %q", size[0], size[1], want)
			}
		}
		if !strings.Contains(text, " ▸ garden") {
			t.Errorf("selection marker not on 'garden':\n%s", text)
		}
	}
}

func TestRenderLevelSelectMazeSeedAndPreview(t *testing.T) {
	names := []string{"canyon", "garden"}
	m, err := game.MazeFromSeed(1234)
	if err != nil {
		t.Fatal(err)
	}
	v := LSState{Levels: names, Cursor: len(names), Seed: "1234", Preview: m}
	text := RenderLevelSelect(v, 120, 40, Palette()).Text()
	if !strings.Contains(text, "seed: 1234") {
		t.Errorf("seed row missing: %s", text)
	}
	// A tall frame must show the map preview (spawn + exit markers).
	if !strings.Contains(text, "▶") || !strings.Contains(text, "E") {
		t.Errorf("map preview missing spawn/exit:\n%s", text)
	}
	// The minimum frame has no room for the preview but must still render.
	small := RenderLevelSelect(v, 62, 19, Palette())
	if small.W != 62 || small.H != 19 {
		t.Errorf("small frame size = %dx%d", small.W, small.H)
	}
}

// A 32-row terminal is the target size for the scale-1 preview: it must fit
// there and must not appear in the minimum 19-row frame.
func TestRenderLevelSelectPreviewAt32Rows(t *testing.T) {
	m, err := game.MazeFromSeed(7)
	if err != nil {
		t.Fatal(err)
	}
	names := []string{"canyon", "garden", "hub", "winding"}
	v := LSState{Levels: names, Cursor: len(names), Seed: "7", Preview: m}
	text := RenderLevelSelect(v, 80, 32, Palette()).Text()
	if !strings.Contains(text, "▶") || !strings.Contains(text, "E") {
		t.Errorf("80x32 frame missing the map preview:\n%s", text)
	}
	if strings.Contains(RenderLevelSelect(v, 62, 19, Palette()).Text(), "▶") {
		t.Error("map preview leaked into the 19-row frame")
	}
}

// Every pre-game screen is a full-window rounded box: corners, side rails,
// the screen title embedded in the top border, and its footer group in the
// bottom border.
func TestScreenBoxFraming(t *testing.T) {
	names := []string{"canyon", "garden"}
	v := LSState{Levels: names, Cursor: 0, Seed: "42"}
	screens := map[string]*Frame{
		"title":     titleStandbyAt(62, 19, 0, 0, nil),
		"menu":      RenderMenu(62, 19, 0, Palette()),
		"help":      RenderHelp(62, 19, Palette()),
		"hiscores":  RenderHighScores(62, 19, 0, map[string]int{"a": 1}, Palette()),
		"select":    RenderLevelSelect(v, 62, 19, Palette()),
		"title-80x": titleStandbyAt(80, 24, 0, 0, nil),
	}
	for name, f := range screens {
		w, h := f.W, f.H
		corners := []struct {
			x, y int
			r    rune
		}{
			{0, 0, '╭'}, {w - 1, 0, '╮'}, {0, h - 1, '╰'}, {w - 1, h - 1, '╯'},
		}
		for _, c := range corners {
			if got := f.C[c.y*w+c.x].R; got != c.r {
				t.Errorf("%s: corner (%d,%d) = %q, want %q", name, c.x, c.y, got, c.r)
			}
		}
		for y := 1; y < h-1; y++ {
			if got := f.C[y*w].R; got != '│' {
				t.Errorf("%s: left rail row %d = %q, want │", name, y, got)
			}
			if got := f.C[y*w+w-1].R; got != '│' {
				t.Errorf("%s: right rail row %d = %q, want │", name, y, got)
			}
		}
		if !strings.Contains(f.Text(), "┐") || !strings.Contains(f.Text(), "┘") {
			t.Errorf("%s: missing embedded title brackets\n%s", name, f.Text())
		}
	}
}

// The scripted loop must show each phase's signature content in order.
// (fr is the internal battle frame; the intro frames precede it.)
func TestTitlePhases(t *testing.T) {
	scores := map[string]int{"ada": 42}
	text := func(fr int) string {
		return titleAt(100, 30, 0, fr, scores).Text()
	}
	cases := []struct {
		fr   int
		want []string
		abs  []string
	}{
		{200, []string{"BATTLE", "WAVE 1"}, nil},
		{600, []string{"WAVE 2"}, nil},
		{1100, []string{"WAVE 3", "B"}, nil},
		{1505, []string{"BREACH"}, nil},
		{1514, nil, []string{"BATTLE", "TDEF", "·"}},                  // static whiteout
		{1525, []string{"█"}, []string{"BATTLE"}},                     // shockwave front
		{1590, nil, []string{"BATTLE", "█"}},                          // cooled grid only
		{1640, nil, []string{"BATTLE", "[enter]"}},                    // mid-reboot
		{1683, []string{"[enter] start"}, []string{"BATTLE", "WAVE"}}, // rebooted into the standby
	}
	for _, c := range cases {
		got := text(c.fr)
		for _, s := range c.want {
			if !strings.Contains(got, s) {
				t.Errorf("frame %d: missing %q", c.fr, s)
			}
		}
		for _, s := range c.abs {
			if strings.Contains(got, s) {
				t.Errorf("frame %d: should not contain %q:\n%s", c.fr, s, got)
			}
		}
	}
	// The void is exactly one ignition spark at the center.
	f := titleAt(100, 30, 0, 1608, scores)
	n := 0
	for _, c := range f.C {
		if c.R != ' ' && c.R != 0 {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("void frame 1608: %d non-space cells, want 1", n)
	}
	if c := f.C[15*100+50]; c.R != '·' || c.FG != 231 {
		t.Fatalf("ignition spark = %+v, want ·/231 at the center", c)
	}
}

// The reboot's last frame must be exactly the standby, so the attract loop
// (standby -> battle -> shockwave -> reboot -> standby) is seamless.
func TestTitleRebootSeam(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {100, 30}, {137, 45}} {
		scores := map[string]int{"ada": 42}
		a := titleAt(size[0], size[1], 37, titleCycle-1, scores)
		b := titleStandbyAt(size[0], size[1], 37, titleIdleWait-1, scores)
		if !reflect.DeepEqual(a, b) {
			t.Errorf("w=%d h=%d: reboot end differs from the standby", size[0], size[1])
		}
	}
}

// The last blast frame must be a full-frame grid: nothing but grid points
// (·) on the 3x2 lattice and spaces.
func TestTitleBlastCoversFrame(t *testing.T) {
	f := titleAt(100, 30, 0, titleBlastEnd-1, nil)
	for y := 0; y < f.H; y++ {
		for x := 0; x < f.W; x++ {
			c := f.C[y*f.W+x]
			if c.R != ' ' && c.R != '·' {
				t.Fatalf("cell %d,%d = %q, want space or grid point", x, y, c.R)
			}
			if c.R == ' ' && x%3 == 0 && y%2 == 0 {
				t.Fatalf("cell %d,%d: missing grid point", x, y)
			}
		}
	}
}

// The battle script: waves arrive, die, and leak on schedule, and the five
// towers sit on the right rows.
func TestTitleBattleScript(t *testing.T) {
	const w, h = 100, 30
	off := screenOff(h)
	pathY, topY, botY := off+demoPath, off+demoUpper, off+demoLower
	pathRow := func(fr int) string {
		f := titleAt(w, h, 0, fr, nil)
		var b strings.Builder
		for x := 0; x < w; x++ {
			b.WriteRune(f.C[pathY*w+x].R)
		}
		return b.String()
	}
	if got := pathRow(180); strings.Count(got, "o") < 3 {
		t.Errorf("frame 180: want >=3 wave-1 minions on the path, got %q", got)
	}
	// The first minion dies at frame 370 and bursts at its death point
	// (u = 340/445, inside the mortar's range).
	f371 := titleAt(w, h, 0, 371, nil)
	if r := f371.C[pathY*w+demoX(w, 340.0/445)].R; r != '*' {
		t.Errorf("frame 371: death burst rune = %q, want *", r)
	}
	if got := pathRow(600); !strings.Contains(got, "r") {
		t.Errorf("frame 600: want a wave-2 runner on the path, got %q", got)
	}
	if got := pathRow(1120); !strings.Contains(got, "B") || strings.Count(got, "o") < 3 {
		t.Errorf("frame 1120: want boss + 3 wave-3 minions, got %q", got)
	}
	// The leaked runner (exits at frame 855) flashes the exit marker red.
	fl := titleAt(w, h, 0, 856, nil)
	if c := fl.C[pathY*w+(w-3)]; c.BG != 167 {
		t.Errorf("frame 856: exit cell BG = %d, want 167 (leak flash)", c.BG)
	}
	// Tower rows: sniper and mortar above, gunner/cannon/frost below.
	f0 := titleAt(w, h, 0, 0, nil)
	for _, tw := range demoTowers {
		tx := demoX(w, tw.u)
		row := topY
		if !tw.above {
			row = botY
		}
		want := game.TowerSpecs[tw.kind].Short
		if r := f0.C[row*w+tx].R; r != want {
			t.Errorf("frame 0: tower %d at %d,%d = %q, want %q", tw.kind, tx, row, r, want)
		}
	}
}

// The title's prompt lives in the bottom border (btop style) and blinks.
func TestTitlePromptInBottomBorder(t *testing.T) {
	f := titleStandbyAt(62, 19, 0, 0, nil)
	bottom := f.Text()
	lines := strings.Split(bottom, "\n")
	row := lines[18]
	for _, want := range []string{"[enter] start", "q quit", "┘", "└"} {
		if !strings.Contains(row, want) {
			t.Errorf("bottom border row missing %q: %q", want, row)
		}
	}
}

// The boot cinematic plays in order — black, ignition, ring settle,
// crosshair zap, box self-draw, light-pen trace, the weapon fan, white
// flash, subtitle decode, chrome fade-in — and its last frame is exactly
// the standby.
func TestTitleBootSequence(t *testing.T) {
	const w, h = 100, 30
	off := screenOff(h)
	x0 := (w - 42) / 2
	// Boot 0: pure black.
	f0 := RenderTitle(w, h, 0, 0, nil, Palette())
	for i, c := range f0.C {
		if c.R != ' ' && c.R != 0 {
			t.Fatalf("boot 0 cell %d = %q, want black", i, c.R)
		}
	}
	// Boot 6: the ignition point at the center.
	f6 := RenderTitle(w, h, 6, 6, nil, Palette())
	if c := f6.C[(h/2)*w+w/2]; c.R != '·' || c.FG != 231 {
		t.Fatalf("boot 6 ignition = %+v, want ·/231 at center", c)
	}
	// Boot 16: the rings have settled dim, the center point cools.
	f16 := RenderTitle(w, h, 16, 16, nil, Palette())
	if c := f16.C[(h/2)*w+w/2]; c.R != '·' || c.FG != 236 {
		t.Fatalf("boot 16 settled point = %+v, want ·/236", c)
	}
	// Boot 19: the crosshair zaps out across the whole frame.
	f19 := RenderTitle(w, h, 19, 19, nil, Palette())
	if c := f19.C[(h/2)*w+0]; c.R != '·' || c.FG != 234 {
		t.Fatalf("boot 19 cross left = %+v, want ·/234", c)
	}
	if c := f19.C[0*w+w/2]; c.R != '·' || c.FG != 234 {
		t.Fatalf("boot 19 cross top = %+v, want ·/234", c)
	}
	// Boot 25: the slab box is mid self-draw — corner landed, tip white.
	f25 := RenderTitle(w, h, 25, 25, nil, Palette())
	if c := f25.C[(off+1)*w+x0-1]; c.R != '╭' || c.FG != 234 {
		t.Fatalf("boot 25 box corner = %+v, want ╭/234", c)
	}
	if c := f25.C[(off+7)*w+x0+34]; c.FG != 51 {
		t.Fatalf("boot 25 box tip = %+v, want 51", c)
	}
	// Boot 40: the pen is mid-trace — the tip is white, untraced cells
	// flicker as dim noise, the box stays drawn.
	f40 := RenderTitle(w, h, 40, 40, nil, Palette())
	tip := titlePenPath[bootPenIndex(40)]
	if c := f40.C[(off+2+tip.row)*w+x0+tip.relX]; c.R != '█' || c.FG != 255 || !c.Bold {
		t.Fatalf("boot 40 pen tip = %+v, want █/255/bold", c)
	}
	ghost := titlePenPath[len(titlePenPath)-1]
	if c := f40.C[(off+2+ghost.row)*w+x0+ghost.relX]; (c.R != '·' && c.R != '+' && c.R != '░') || c.FG < 236 || c.FG > 238 {
		t.Fatalf("boot 40 untraced cell = %+v, want dim noise", c)
	}
	// Boot 94: the pen has finished — the whole slab is locked in color.
	f94 := RenderTitle(w, h, 94, 94, nil, Palette())
	if c := f94.C[(off+2)*w+x0]; c.R != '█' || c.FG != titleColors[0] {
		t.Fatalf("boot 94 T cell = %+v, want █/46", c)
	}
	if c := f94.C[(off+6)*w+x0+3]; c.R != '▒' || c.FG != titleColors[0] {
		t.Fatalf("boot 94 T bottom cell = %+v, want ▒/46", c)
	}
	// Boot 100: T's Gunner tracer is in flight, reticle locked at the target.
	f100 := RenderTitle(w, h, 100, 100, nil, Palette())
	if c := f100.C[13*w+25]; c.R != '█' || c.FG != 255 {
		t.Fatalf("boot 100 tracer = %+v, want █/255", c)
	}
	if c := f100.C[27*w+11]; c.R != '+' || c.FG != 244 {
		t.Fatalf("boot 100 reticle = %+v, want +/244", c)
	}
	// Boot 145: the Cannon shell has burst at the target.
	f145 := RenderTitle(w, h, 145, 145, nil, Palette())
	if c := f145.C[2*w+41]; c.R != '·' || c.FG != 240 {
		t.Fatalf("boot 145 burst core = %+v, want ·/240", c)
	}
	if c := f145.C[2*w+39]; c.R != '·' || c.FG != 220 {
		t.Fatalf("boot 145 burst ring = %+v, want ·/220", c)
	}
	// Boot 174: E's Sniper beam is lit full length.
	f174 := RenderTitle(w, h, 174, 174, nil, Palette())
	if c := f174.C[5*w+56]; c.R != '█' || c.FG != 255 {
		t.Fatalf("boot 174 beam = %+v, want █/255", c)
	}
	// Boot 205: F's Tesla chain is live — letter settled, frame busy.
	f205 := RenderTitle(w, h, 205, 205, nil, Palette())
	if c := f205.C[7*w+62]; c.R != '█' || c.FG != 171 {
		t.Fatalf("boot 205 F letter = %+v, want █/171", c)
	}
	n := 0
	for _, c := range f205.C {
		if c.R != ' ' {
			n++
		}
	}
	if n < 200 {
		t.Fatalf("boot 205: %d non-blank cells, want the arc + slab + box", n)
	}
	// Boot 226: the whole slab burns white over the grid.
	f226 := RenderTitle(w, h, 226, 226, nil, Palette())
	for _, s := range titlePenPath {
		if c := f226.C[(off+2+s.row)*w+x0+s.relX]; c.FG != 255 || !c.Bold {
			t.Fatalf("boot 226 slab cell (letter %d) = %+v, want white", s.li, c)
		}
	}
	if c := f226.C[0]; c.R != '·' || c.FG != 234 {
		t.Fatalf("boot 226 grid = %+v, want ·/234", c)
	}
	// Boot 235: the subtitle decodes left to right with a caret.
	f235 := RenderTitle(w, h, 235, 235, nil, Palette())
	tag := []rune(titleTagline)
	sx, sy := (w-len(tag))/2, off+8
	if c := f235.C[sy*w+sx+11]; c.FG != 255 || !c.Bold {
		t.Fatalf("boot 235 subtitle head = %+v, want white bold", c)
	}
	if c := f235.C[sy*w+sx+12]; c.R != '█' || c.FG != 251 {
		t.Fatalf("boot 235 caret = %+v, want █/251", c)
	}
	// Boot 245: the chrome fades in — border and footer still ghosted,
	// the signature not yet up. (x=20 is clear of the TDEF title.)
	f245 := RenderTitle(w, h, 245, 245, nil, Palette())
	if c := f245.C[0*w+20]; c.FG != 234 {
		t.Fatalf("boot 245 border = %+v, want ghost 234", c)
	}
	if c := f245.C[(h-1)*w+40]; c.R != '·' || c.FG != 234 {
		t.Fatalf("boot 245 footer = %+v, want ghost", c)
	}
	if c := f245.C[(h-1)*w+(w-11)]; c.R != '─' {
		t.Fatalf("boot 245 sig slot = %+v, want border dash", c)
	}
	// Boot 259: border and battlefield have landed, roster still ghosting,
	// the signature embedded in the bottom border.
	f259 := RenderTitle(w, h, 259, 259, nil, Palette())
	if c := f259.C[0*w+20]; c.FG != 240 {
		t.Fatalf("boot 259 border = %+v, want 240", c)
	}
	if c := f259.C[(off+14)*w+40]; c.R != '·' || c.FG != 234 {
		t.Fatalf("boot 259 roster = %+v, want ghost", c)
	}
	if c := f259.C[(h-1)*w+(w-11)]; c.R != 'b' || c.FG != 238 {
		t.Fatalf("boot 259 sig start = %+v, want b/238", c)
	}
	if c := f259.C[(h-1)*w+(w-3)]; c.R != 'c' || c.FG != 238 {
		t.Fatalf("boot 259 sig end = %+v, want c/238", c)
	}
	// Boot 272 (the last boot frame) == the standby at the same clock.
	if a, b := RenderTitle(w, h, 272, 272, nil, Palette()), RenderTitle(w, h, 272, 273, nil, Palette()); !reflect.DeepEqual(a, b) {
		t.Fatal("the boot's last frame differs from the standby")
	}
}

// The standby battlefield is empty: no BATTLE/WAVE text, no towers, no
// enemies — just the path, its ambient packet, and the spawn/exit markers.
func TestTitleStandbyEmpty(t *testing.T) {
	const w, h = 100, 30
	off := screenOff(h)
	f := titleStandbyAt(w, h, 0, 100, nil)
	text := f.Text()
	for _, absent := range []string{"BATTLE", "WAVE"} {
		if strings.Contains(text, absent) {
			t.Errorf("standby shows %q:\n%s", absent, text)
		}
	}
	for _, tw := range demoTowers {
		tx := demoX(w, tw.u)
		ty := off + demoLower
		if tw.above {
			ty = off + demoUpper
		}
		if c := f.C[ty*w+tx]; c.R != ' ' {
			t.Errorf("standby tower cell %d,%d = %q, want empty", tx, ty, c.R)
		}
	}
	var b strings.Builder
	for x := 0; x < w; x++ {
		b.WriteRune(f.C[(off+demoPath)*w+x].R)
	}
	path := b.String()
	for _, r := range []string{"o", "r", "B"} {
		if strings.Contains(path, r) {
			t.Errorf("standby path has an enemy %q: %s", r, path)
		}
	}
	if !strings.Contains(path, "▶") || !strings.Contains(path, "E") {
		t.Errorf("standby path missing spawn/exit: %s", path)
	}
}

// The signature sits embedded in the bottom border's right section (mirroring
// the TDEF embed in the top border), muted, and clear of the centered footer.
func TestTitleSigInBottomBorder(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {100, 30}} {
		w, h := size[0], size[1]
		f := titleStandbyAt(w, h, 0, 0, nil)
		text := f.Text()
		if !strings.Contains(text, "by 0xbenc") {
			t.Errorf("w=%d h=%d: standby missing the signature:\n%s", w, h, text)
		}
		x0 := w - len(titleSig) - 2
		if c := f.C[(h-1)*w+x0]; c.R != 'b' || c.FG != 238 {
			t.Errorf("w=%d h=%d: sig start = %+v, want b/238", w, h, c)
		}
		if c := f.C[(h-1)*w+x0+len(titleSig)-1]; c.R != 'c' || c.FG != 238 {
			t.Errorf("w=%d h=%d: sig end = %+v, want c/238", w, h, c)
		}
		// The border cell between the footer and the signature is intact.
		lines := strings.Split(text, "\n")
		row := lines[h-1]
		if !strings.Contains(row, "[enter] start") || !strings.Contains(row, "q quit") {
			t.Errorf("w=%d h=%d: footer broken by the signature: %q", w, h, row)
		}
		if c := f.C[(h-1)*w+x0-3]; c.R != '─' || c.FG != 240 {
			t.Errorf("w=%d h=%d: border before the signature = %+v, want ─/240", w, h, c)
		}
	}
}

// The attract loop: 15s (450 frames) of standby, then the BATTLE/WAVE
// decode and tower power-up, then the battle script; it repeats.
func TestTitleIdleGate(t *testing.T) {
	const w, h = 100, 30
	text := func(boot int) string {
		return RenderTitle(w, h, boot, boot, nil, Palette()).Text()
	}
	// The frame before the 15s deadline: still idle.
	if got := text(titleBootLen + titleIdleWait - 1); strings.Contains(got, "BATTLE") {
		t.Fatalf("1 frame before the idle deadline the battle text is up:\n%s", got)
	}
	// 1s past: BATTLE has decoded, WAVE 1 is still arriving.
	at := titleBootLen + titleIdleWait + 20
	if got := text(at); !strings.Contains(got, "BATTLE") || strings.Contains(got, "WAVE 1") {
		t.Fatalf("1s past the deadline want BATTLE up and WAVE 1 decoding:\n%s", text(at))
	}
	// 2s past: WAVE 1 is on.
	at += 20
	if got := text(at); !strings.Contains(got, "WAVE 1") {
		t.Fatalf("2s past the deadline want WAVE 1:\n%s", text(at))
	}
	// Past the lead-in: the wave-1 script is running.
	at = titleBootLen + titleIdleWait + titleBattleLead + 200
	if got := text(at); !strings.Contains(got, "WAVE 1") {
		t.Fatalf("battle frame 200: want the WAVE 1 label:\n%s", text(at))
	}
	// One full attract cycle later the battle is running again.
	at = titleBootLen + titleAttractCycle + titleIdleWait + titleBattleLead + 600
	if got := text(at); !strings.Contains(got, "WAVE 2") {
		t.Fatalf("second cycle, battle frame 600: want WAVE 2:\n%s", text(at))
	}
}

// The level-select map preview appears in a [ PREVIEW ] sub-box exactly when
// the frame has H*2+2 rows (scale 2) or H+2 rows (scale 1) below the chrome.
func TestLevelSelectPreviewBoundary(t *testing.T) {
	m, err := game.MazeFromSeed(7)
	if err != nil {
		t.Fatal(err)
	}
	names := []string{"canyon", "garden", "hub", "winding"}
	v := LSState{Levels: names, Cursor: len(names), Seed: "7", Preview: m}
	// 80x29 leaves 14 rows below the chrome: one short of the 15 the
	// scale-1 preview needs, so it degrades to the hint line.
	small := RenderLevelSelect(v, 80, 29, Palette()).Text()
	if strings.Contains(small, "PREVIEW") {
		t.Errorf("80x29 should not fit the preview:\n%s", small)
	}
	if !strings.Contains(small, "preview needs more room") {
		t.Errorf("80x29 missing the preview hint:\n%s", small)
	}
	// 80x30 leaves exactly 15: the scale-1 preview with its sub-box label.
	big := RenderLevelSelect(v, 80, 30, Palette()).Text()
	if !strings.Contains(big, "PREVIEW") || !strings.Contains(big, "▶") {
		t.Errorf("80x30 should show the preview sub-box:\n%s", big)
	}
}

func TestLSRectsMatchRenderedRows(t *testing.T) {
	w, h := 62, 19
	names := []string{"canyon", "garden", "hub", "winding"}
	v := LSState{Levels: names, Cursor: 0, Seed: "42"}
	f := RenderLevelSelect(v, w, h, Palette())
	lines := strings.Split(f.Text(), "\n")
	rows, diffs, seedRect := LSRects(v, w, h)
	if len(rows) != len(names)+1 {
		t.Fatalf("rows = %d, want %d", len(rows), len(names)+1)
	}
	for i, r := range rows {
		label := lsMazeRow
		if i < len(names) {
			label = names[i]
		}
		if !strings.Contains(lines[r.Y], label) {
			t.Errorf("row %d rect at %d misses %q: %q", i, r.Y, label, lines[r.Y])
		}
	}
	for i, name := range []string{"easy", "normal", "hard"} {
		y, x := diffs[i].Y, diffs[i].X+1 // label text starts one column in
		if x >= w || y >= h {
			t.Fatalf("difficulty %q rect out of frame", name)
		}
		if got := f.C[y*w+x].R; got != rune(name[0]) {
			t.Errorf("difficulty %q rect at (%d,%d) hits %q, want %q", name, x, y, got, name[0])
		}
	}
	if !strings.Contains(lines[seedRect.Y], "seed:") {
		t.Errorf("seed rect at row %d misses the seed line: %q", seedRect.Y, lines[seedRect.Y])
	}
}
