package render

import (
	"fmt"
	"strings"
)

type Cell struct {
	R    rune
	FG   int
	BG   int
	Bold bool
}

type Frame struct {
	W, H int
	C    []Cell
}

func (f *Frame) Set(x, y int, c Cell) {
	if x < 0 || y < 0 || x >= f.W || y >= f.H {
		return
	}
	f.C[y*f.W+x] = c
}

func (f *Frame) Put(x, y int, r rune, fg, bg int) {
	f.Set(x, y, Cell{R: r, FG: fg, BG: bg})
}

func (f *Frame) Text() string {
	lines := make([]string, f.H)
	for y := 0; y < f.H; y++ {
		var b strings.Builder
		for x := 0; x < f.W; x++ {
			c := f.C[y*f.W+x]
			if c.R == 0 {
				b.WriteRune(' ')
			} else {
				b.WriteRune(c.R)
			}
		}
		lines[y] = b.String()
	}
	return strings.Join(lines, "\n")
}

func (f *Frame) ANSI() string {
	var b strings.Builder
	prevFG, prevBG := -1, -1
	prevBold := false
	b.WriteString("\x1b[?25l\x1b[2J\x1b[H")
	for y := 0; y < f.H; y++ {
		// Position each row with CUP. A bare cursor-up (\x1b[M) keeps the
		// column, so on a terminal wider than the frame every row after
		// the first would start at column W+1 and the dump would smear.
		if y > 0 {
			fmt.Fprintf(&b, "\x1b[%d;1H", y+1)
		}
		for x := 0; x < f.W; x++ {
			c := f.C[y*f.W+x]
			fg := c.FG
			if fg == 0 {
				fg = 255
			}
			bg := c.BG
			bold := c.Bold
			if fg != prevFG || bg != prevBG || bold != prevBold {
				if fg != prevFG {
					fmt.Fprintf(&b, "\x1b[38;5;%dm", fg)
					prevFG = fg
				}
				if bg != prevBG {
					if bg == 0 {
						b.WriteString("\x1b[49m")
					} else {
						fmt.Fprintf(&b, "\x1b[48;5;%dm", bg)
					}
					prevBG = bg
				}
				if bold != prevBold {
					if bold {
						b.WriteString("\x1b[1m")
					} else {
						b.WriteString("\x1b[22m")
					}
					prevBold = bold
				}
			}
			if c.R == 0 {
				b.WriteRune(' ')
			} else {
				b.WriteRune(c.R)
			}
		}
	}
	b.WriteString("\x1b[0m\x1b[?25h")
	return b.String()
}

// drawRoundedBox draws a rounded box from (x,y) to (x+w-1,y+h-1) in border
// color fg: ╭╮╰╯ corners, ─ edges, │ sides. Interior cells are untouched.
func drawRoundedBox(f *Frame, x, y, w, h, fg int) {
	if w < 2 || h < 2 {
		return
	}
	for i := 1; i < w-1; i++ {
		f.Set(x+i, y, Cell{R: '─', FG: fg})
		f.Set(x+i, y+h-1, Cell{R: '─', FG: fg})
	}
	for j := 1; j < h-1; j++ {
		f.Set(x, y+j, Cell{R: '│', FG: fg})
		f.Set(x+w-1, y+j, Cell{R: '│', FG: fg})
	}
	f.Set(x, y, Cell{R: '╭', FG: fg})
	f.Set(x+w-1, y, Cell{R: '╮', FG: fg})
	f.Set(x, y+h-1, Cell{R: '╰', FG: fg})
	f.Set(x+w-1, y+h-1, Cell{R: '╯', FG: fg})
}

// embedSegment writes a btop-style border title starting at (x,y): the
// brackets (left/right, e.g. ┐┌ on the top border, ┘└ on the bottom) in the
// border color, the text in fg (bold when bold), and returns the next free
// column. It overwrites whatever was in the span, so the border row is
// usually drawn first.
func embedSegment(f *Frame, y, x int, text string, left, right rune, border, fg int, bold bool) int {
	f.Set(x, y, Cell{R: left, FG: border})
	x++
	for _, r := range text {
		f.Set(x, y, Cell{R: r, FG: fg, Bold: bold})
		x++
	}
	f.Set(x, y, Cell{R: right, FG: border})
	return x + 1
}

// centerEmbed centers a border title on row y of the frame (one column of
// clearance from each corner) and returns the next free column after it.
func centerEmbed(f *Frame, y int, text string, left, right rune, border, fg int, bold bool) int {
	w := len([]rune(text)) + 2
	x := (f.W - w) / 2
	if x < 1 {
		x = 1
	}
	if x+w > f.W-1 {
		x = f.W - 1 - w
	}
	if x < 1 {
		x = 1
	}
	return embedSegment(f, y, x, text, left, right, border, fg, bold)
}

// truncateRunes returns the first n runes of s (s itself when shorter).
func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
