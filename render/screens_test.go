package render

import (
	"reflect"
	"strings"
	"testing"

	"tdef/game"
)

func TestRenderTitleDeterministic(t *testing.T) {
	scores := map[string]int{"winding": 100, "hub": 90}
	for _, fr := range []int{0, 1, 7, 30, 90, 1234} {
		a := RenderTitle(62, 19, fr, scores, Palette())
		b := RenderTitle(62, 19, fr, scores, Palette())
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("title frame %d is not deterministic", fr)
		}
	}
}

// The title must animate: distinct frames must differ somewhere.
func TestRenderTitleAnimates(t *testing.T) {
	if reflect.DeepEqual(RenderTitle(80, 24, 0, nil, Palette()), RenderTitle(80, 24, 4, nil, Palette())) {
		t.Fatal("title does not animate (frames 0 and 4 identical)")
	}
	if reflect.DeepEqual(RenderTitle(80, 24, 0, nil, Palette()), RenderTitle(80, 24, 15, nil, Palette())) {
		t.Fatal("title blink period not moving (frames 0 and 15 identical)")
	}
}

// At the minimum 62x19 frame the whole card (logo, tagline, demo, roster,
// best, prompt) must be present and nothing clipped.
func TestRenderTitleFitsFrame(t *testing.T) {
	scores := map[string]int{"maze12345678901234567890": 9999999999}
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		w, h := size[0], size[1]
		// Cover one full demo cycle (2 ticks per step, walk+leak steps).
		period := 2 * (w - 7 + 6)
		for fr := 0; fr < period; fr++ {
			f := RenderTitle(w, h, fr, scores, Palette())
			if f.W != w || f.H != h {
				t.Fatalf("w=%d h=%d frame %d: size = %dx%d", w, h, fr, f.W, f.H)
			}
			if rowWidth(f, 0) != w || rowWidth(f, h-1) != w {
				t.Fatalf("w=%d h=%d frame %d: border row clipped", w, h, fr)
			}
		}
		text := RenderTitle(w, h, 0, scores, Palette()).Text()
		for _, want := range []string{
			"█████████", // logo top bevel (T)
			"▒▒▒",       // logo bottom bevel
			"— terminal tower defense —",
			"▶", "E", // demo spawn / exit
			" towers ", " enemies ",
			"★ best 9999999999", // best line (longest plausible key)
			"[enter] start", "q quit",
		} {
			if !strings.Contains(text, want) {
				t.Errorf("w=%d h=%d: title missing %q:\n%s", w, h, want, text)
			}
		}
	}
	// Empty hiscore shows the placeholder.
	if !strings.Contains(RenderTitle(62, 19, 0, nil, Palette()).Text(), "no scores yet") {
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
		"title":     RenderTitle(62, 19, 0, nil, Palette()),
		"menu":      RenderMenu(62, 19, 0, Palette()),
		"help":      RenderHelp(62, 19, Palette()),
		"hiscores":  RenderHighScores(62, 19, 0, map[string]int{"a": 1}, Palette()),
		"select":    RenderLevelSelect(v, 62, 19, Palette()),
		"title-80x": RenderTitle(80, 24, 0, nil, Palette()),
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

// The title's prompt lives in the bottom border (btop style) and blinks.
func TestTitlePromptInBottomBorder(t *testing.T) {
	f := RenderTitle(62, 19, 0, nil, Palette())
	bottom := f.Text()
	lines := strings.Split(bottom, "\n")
	row := lines[18]
	for _, want := range []string{"[enter] start", "q quit", "┘", "└"} {
		if !strings.Contains(row, want) {
			t.Errorf("bottom border row missing %q: %q", want, row)
		}
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
