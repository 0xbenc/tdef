package render

import "testing"

func TestDrawRoundedBox(t *testing.T) {
	f := &Frame{W: 20, H: 10, C: make([]Cell, 20*10)}
	drawRoundedBox(f, 2, 3, 10, 5, 240)
	want := map[[2]int]rune{
		{2, 3}: '╭', {11, 3}: '╮', {2, 7}: '╰', {11, 7}: '╯',
	}
	for p, r := range want {
		if got := f.C[p[1]*20+p[0]].R; got != r {
			t.Fatalf("corner at %v: got %q want %q", p, got, r)
		}
	}
	for x := 3; x <= 10; x++ {
		if r := f.C[3*20+x].R; r != '─' {
			t.Fatalf("top edge x=%d: got %q want ─", x, r)
		}
		if r := f.C[7*20+x].R; r != '─' {
			t.Fatalf("bottom edge x=%d: got %q want ─", x, r)
		}
	}
	for y := 4; y <= 6; y++ {
		if r := f.C[y*20+2].R; r != '│' {
			t.Fatalf("left edge y=%d: got %q want │", y, r)
		}
		if r := f.C[y*20+11].R; r != '│' {
			t.Fatalf("right edge y=%d: got %q want │", y, r)
		}
	}
	if c := f.C[5*20+6]; c.R != 0 {
		t.Fatalf("interior (6,5) must be untouched, got %q", c.R)
	}
	if f.C[3*20+2].FG != 240 {
		t.Fatalf("border FG: got %d want 240", f.C[3*20+2].FG)
	}
}

func TestEmbedSegmentReplacesBorder(t *testing.T) {
	f := &Frame{W: 20, H: 1, C: make([]Cell, 20)}
	for x := 0; x < 20; x++ {
		f.C[x] = Cell{R: '─', FG: 240}
	}
	end := embedSegment(f, 0, 6, "TITLE", '┐', '┌', 240, 254, true)
	if end != 13 {
		t.Fatalf("next free column: got %d want 13", end)
	}
	if r := f.C[6].R; r != '┐' {
		t.Fatalf("left bracket: got %q", r)
	}
	for i, want := range []rune("TITLE") {
		c := f.C[7+i]
		if c.R != want || c.FG != 254 || !c.Bold {
			t.Fatalf("title char %d: got %q fg=%d bold=%v", i, c.R, c.FG, c.Bold)
		}
	}
	if r := f.C[12].R; r != '┌' {
		t.Fatalf("right bracket: got %q", r)
	}
	if r, fg := f.C[5].R, f.C[5].FG; r != '─' || fg != 240 {
		t.Fatalf("cell before span must stay ─ 240, got %q %d", r, fg)
	}
	if r, fg := f.C[13].R, f.C[13].FG; r != '─' || fg != 240 {
		t.Fatalf("cell after span must stay ─ 240, got %q %d", r, fg)
	}
}

func TestCenterEmbedClearsCorners(t *testing.T) {
	f := &Frame{W: 20, H: 1, C: make([]Cell, 20)}
	for x := 0; x < 20; x++ {
		f.C[x] = Cell{R: '─', FG: 240}
	}
	centerEmbed(f, 0, "TITLE", '┐', '┌', 240, 254, true)
	// (20 - 7) / 2 = 6; must not touch the corner columns 0 and 19.
	if r := f.C[6].R; r != '┐' {
		t.Fatalf("centered left bracket at x=6: got %q", r)
	}
	if r := f.C[0].R; r != '─' {
		t.Fatalf("corner column must be untouched: got %q", r)
	}
}

func TestTruncateRunesRuneSafe(t *testing.T) {
	if got := truncateRunes("héllo", 3); got != "hél" {
		t.Fatalf("got %q want hél", got)
	}
	if got := truncateRunes("ab", 5); got != "ab" {
		t.Fatalf("got %q want ab", got)
	}
}
