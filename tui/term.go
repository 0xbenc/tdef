package tui

// RefreshSize re-queries the terminal size. Call it from the main
// goroutine after consuming a Winch token: by then the resize has
// settled, and keeping all t.size writes on one goroutine means Size()
// can read it without synchronization.
func (t *Terminal) RefreshSize() { t.size = t.querySize() }

func (t *Terminal) Size() (int, int) { return t.size[0], t.size[1] }

func (t *Terminal) Winch() <-chan struct{} { return t.winch }

func (t *Terminal) querySize() [2]int {
	if w, h, err := termSize(int(t.out.Fd())); err == nil && w > 0 && h > 0 {
		return [2]int{w, h}
	}
	return [2]int{80, 24}
}

func (t *Terminal) Write(p []byte) { t.out.Write(p) }

func (t *Terminal) AltScreen(on bool) {
	if on {
		t.Write([]byte("\x1b[?1049h"))
	} else {
		t.Write([]byte("\x1b[?1049l"))
	}
}

func (t *Terminal) Cursor(on bool) {
	if on {
		t.Write([]byte("\x1b[?25h"))
	} else {
		t.Write([]byte("\x1b[?25l"))
	}
}

func (t *Terminal) mouseVT(on bool) {
	if on {
		t.Write([]byte("\x1b[?1000h\x1b[?1002h\x1b[?1005h\x1b[?1006h"))
	} else {
		t.Write([]byte("\x1b[?1006l\x1b[?1005l\x1b[?1002l\x1b[?1000l"))
	}
}
