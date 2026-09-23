package tui

import (
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/0xbenc/tdef/game"
	"github.com/0xbenc/tdef/render"
)

// capWriter captures the bytes a blit writes, so a test can replay them on a
// virtual terminal and check what the user would actually see.
type capWriter struct{ data []byte }

func (c *capWriter) Write(p []byte) { c.data = append(c.data, p...) }

// vt is a minimal virtual terminal: it replays the blit's ANSI (cursor moves,
// clear, SGR) onto a grid of runes, so a test can read back the pixels.
type vt struct {
	w, h int
	r    []rune // w*h
	row  int    // 1-based, as in the escape sequences
	col  int    // 1-based
}

func newVT(w, h int) *vt {
	return &vt{w: w, h: h, r: make([]rune, w*h), row: 1, col: 1}
}

func (t *vt) clear() {
	for i := range t.r {
		t.r[i] = ' '
	}
	t.row, t.col = 1, 1
}

func (t *vt) set(r rune) {
	if t.row < 1 || t.row > t.h || t.col < 1 || t.col > t.w {
		return
	}
	t.r[(t.row-1)*t.w+(t.col-1)] = r
	t.col++
}

// replay feeds a captured blit chunk to the virtual terminal.
func (t *vt) replay(data []byte) {
	for i := 0; i < len(data); {
		b := data[i]
		if b != 0x1b {
			r, size := utf8.DecodeRune(data[i:])
			t.set(r)
			i += size
			continue
		}
		i++ // consume ESC
		if i >= len(data) || data[i] != '[' {
			continue // ignore non-CSI escapes
		}
		i++ // consume '['
		j := i
		for j < len(data) && (data[j] < 0x20 || (data[j] >= 0x30 && data[j] <= 0x7e)) && !isCSIEnd(data[j]) {
			j++
		}
		params := string(data[i:j])
		if j >= len(data) {
			break
		}
		switch data[j] {
		case 'H':
			parts := strings.Split(params, ";")
			row, _ := strconv.Atoi(parts[0])
			col := 1
			if len(parts) > 1 {
				col, _ = strconv.Atoi(parts[1])
			}
			t.row, t.col = row, col
		case 'J':
			t.clear()
		case 'm':
			// SGR: ignored for a rune grid.
		}
		i = j + 1
	}
}

func isCSIEnd(b byte) bool {
	switch b {
	case 'H', 'J', 'm', 'A', 'B', 'C', 'D', 'd', 'G':
		return true
	}
	return false
}

// headerGold returns the integer that follows "⛁ " on the first row.
func (t *vt) headerGold() (int, bool) {
	row := ""
	for x := 0; x < t.w; x++ {
		row += string(t.r[x])
	}
	const mark = "⛁ "
	idx := strings.Index(row, mark)
	if idx < 0 {
		return 0, false
	}
	rest := row[idx+len(mark):]
	n := 0
	for _, c := range rest {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

// The gold in the header must track the state exactly across frames where the
// value changes length in both directions: a stale tail would leave a digit
// from a longer value (e.g. 1500 -> 950 must not read 9500).
func TestBlitHeaderGoldTracksState(t *testing.T) {
	m, err := game.LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	const W, H = 62, 19
	g := game.NewState(m)
	g.Status = game.StatusRunning
	g.WaveActive = true
	g.Wave = 3
	ui := render.UI{Level: "winding", Selected: render.NoSelection, Cursor: game.Vec{X: 2, Y: 8}}
	term := &Terminal{size: [2]int{W, H}, winch: make(chan struct{}, 1)}
	a := &App{g: g, ui: ui, term: term, pal: render.Palette()}
	screen := newVT(W, H)

	check := func(frameNo int, want int) {
		t.Helper()
		var cap capWriter
		a.blitTo(render.Render(g, &ui, a.pal, W, H, frameNo), &cap)
		screen.replay(cap.data)
		got, ok := screen.headerGold()
		if !ok {
			t.Fatalf("frame %d: no gold in header: %q", frameNo, screen.row0())
		}
		if got != want {
			t.Fatalf("frame %d: header gold = %d, want %d (row0 %q)", frameNo, got, want, screen.row0())
		}
	}

	g.Gold = 1500
	check(0, 1500) // baseline: 4-digit gold
	g.Gold = 950
	check(1, 950) // shrink across a digit boundary (1500 -> 950)
	g.Gold = 10000
	check(2, 10000) // grow back past it
	g.Gold = 42
	check(3, 42) // shrink to two digits
	g.Gold = 1000000
	check(4, 1000000) // grow to seven digits, then shrink again
	g.Gold = 1
	check(5, 1)
}

func (t *vt) row0() string {
	s := ""
	for x := 0; x < t.w; x++ {
		s += string(t.r[x])
	}
	return s
}
