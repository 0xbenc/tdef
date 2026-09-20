package tui

import (
	"os"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

type Key int

const (
	KeyNone Key = iota
	KeyEnter
	KeyEscape
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyBackspace
	KeyCtrlC
	KeyCtrlL
)

type Event struct {
	Key  Key
	Rune rune

	Mouse bool
	Btn   int // 0 left, 1 middle, 2 right, 64 wheel up, 65 wheel down
	Press bool
	X, Y  int
}

type reader struct {
	ch    chan Event
	state int
	buf   []byte

	// escT implements the bare-ESC fallback. A pty does not support read
	// deadlines (SetReadDeadline returns an error), so a lone 0x1b cannot
	// be resolved by the reader unblocking; a timer resolves it instead.
	mu   sync.Mutex
	escT *time.Timer
}

// escWindow is how long to wait after a bare ESC for a following '[' or 'O'
// before treating it as an Escape keypress. Real terminals emit the full CSI
// sequence in one burst; a lone ESC means the user pressed Escape.
const escWindow = 50 * time.Millisecond

func startReader(in *os.File, ch chan Event) {
	r := &reader{ch: ch}
	go func() {
		tmp := make([]byte, 512)
		for {
			n, err := in.Read(tmp)
			if n > 0 {
				r.feed(tmp[:n])
			}
			if err != nil {
				return
			}
		}
	}()
}

// armEsc arms the bare-ESC fallback: if no further byte arrives within
// escWindow, the pending 0x1b is treated as the Escape key.
func (r *reader) armEsc() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.escT != nil {
		r.escT.Stop()
	}
	r.escT = time.AfterFunc(escWindow, func() {
		r.mu.Lock()
		fire := r.state == 1 && r.escT != nil
		if fire {
			r.state = 0
			r.escT = nil
		}
		r.mu.Unlock()
		if fire {
			r.emit(Event{Key: KeyEscape})
		}
	})
}

func (r *reader) feed(data []byte) {
	for len(data) > 0 {
		data = r.step(data)
	}
}

func (r *reader) step(data []byte) []byte {
	for len(data) > 0 {
		r.mu.Lock()
		st := r.state
		r.mu.Unlock()
		b := data[0]
		data = data[1:]
		switch st {
		case 0:
			switch {
			case b == 0x1b:
				r.setState(1)
				if len(data) == 0 {
					// Bare ESC: the reader is blocked in Read and the fd
					// may not support deadlines, so a timer resolves it.
					r.armEsc()
					return data
				}
			case b == '\r' || b == '\n':
				r.emit(Event{Key: KeyEnter})
			case b == 0x7f || b == 0x08:
				r.emit(Event{Key: KeyBackspace})
			case b == 0x03:
				r.emit(Event{Key: KeyCtrlC})
			case b == 0x0c:
				r.emit(Event{Key: KeyCtrlL})
			case b < 0x20:
			default:
				if b < 0x80 {
					r.emit(Event{Rune: rune(b)})
				} else {
					r.buf = append(r.buf, b)
					_, size := utf8.DecodeRune(r.buf)
					for size > len(r.buf) && len(data) > 0 {
						r.buf = append(r.buf, data[0])
						data = data[1:]
						_, size = utf8.DecodeRune(r.buf)
					}
					if size <= len(r.buf) {
						_, sz := utf8.DecodeRune(r.buf)
						r.emit(Event{Rune: rune(r.buf[0])})
						r.buf = r.buf[sz:]
					}
				}
			}
		case 1:
			// A byte arrived for the pending bare ESC. Settle it under the
			// lock so the fallback timer cannot fire between the check and
			// the transition.
			r.mu.Lock()
			if r.state != 1 {
				// The fallback already emitted KeyEscape; parse this byte
				// fresh in state 0.
				r.mu.Unlock()
				data = append([]byte{b}, data...)
				continue
			}
			if r.escT != nil {
				r.escT.Stop()
				r.escT = nil
			}
			if b == 0x1b {
				// A second ESC settles the first one: a CSI sequence never
				// starts with ESC, so the pending 0x1b was the Escape key.
				// Emit it and keep watching — this byte may start the next
				// sequence (a fast Escape+arrow must not lose the arrow).
				r.mu.Unlock()
				r.emit(Event{Key: KeyEscape})
				if len(data) == 0 {
					r.armEsc()
					return data
				}
				continue
			}
			switch b {
			case '[', 'O':
				r.state = 2
				r.buf = r.buf[:0]
				r.mu.Unlock()
			default:
				r.state = 0
				r.mu.Unlock()
				r.emit(Event{Rune: rune(b)})
			}
		case 2:
			if b == '<' {
				r.setState(3)
				r.buf = r.buf[:0]
				continue
			}
			r.buf = append(r.buf, b)
			if k, done := r.interpret(); done {
				r.setState(0)
				if k != KeyNone {
					r.emit(Event{Key: k})
				}
			}
		case 3:
			if b == 'M' || b == 'm' {
				r.setState(0)
				r.mouse(b == 'M')
				continue
			}
			r.buf = append(r.buf, b)
		}
	}
	return data
}

// setState stores a parser state under the lock.
func (r *reader) setState(s int) {
	r.mu.Lock()
	r.state = s
	r.mu.Unlock()
}

func (r *reader) interpret() (Key, bool) {
	s := string(r.buf)
	if s == "" {
		return KeyNone, false
	}
	last := s[len(s)-1]
	if last < 0x40 || last > 0x7e {
		return KeyNone, false
	}
	if last == '~' {
		switch s {
		case "1~", "7~", "8~":
			return KeyUp, true
		case "2~":
			return KeyDown, true
		case "4~", "5~":
			return KeyLeft, true
		case "6~":
			return KeyRight, true
		}
		return KeyNone, true
	}
	switch last {
	case 'A':
		return KeyUp, true
	case 'B':
		return KeyDown, true
	case 'C':
		return KeyRight, true
	case 'D':
		return KeyLeft, true
	case 'H':
		return KeyUp, true
	case 'F':
		return KeyDown, true
	case 'Z':
		return KeyBackspace, true
	}
	return KeyNone, true
}

// press is true for 'M' (button down) and false for 'm' (button up); the
// terminator byte is consumed by step() before this runs.
func (r *reader) mouse(press bool) {
	s := string(r.buf)
	parts := splitN(s, ';')
	if len(parts) < 3 {
		return
	}
	btn, _ := strconv.Atoi(parts[0])
	x, _ := strconv.Atoi(parts[1])
	y, _ := strconv.Atoi(parts[2])
	switch {
	case btn == 64 || btn == 65: // wheel (press events only)
		if press {
			r.emit(Event{Mouse: true, Btn: btn, Press: true, X: -1, Y: -1})
		}
	case btn >= 0 && btn <= 2: // real buttons
		// Button-event tracking (?1002) also reports drag motions with
		// codes 32-35; those are not clicks and must be ignored, or a
		// click-drag would place a tower (or toggle a menu slot) per
		// motion step.
		if x == 0 || y == 0 {
			return
		}
		r.emit(Event{Mouse: true, Btn: btn, Press: press, X: x - 1, Y: y - 1})
	}
}

func splitN(s string, sep byte) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

func (r *reader) emit(e Event) {
	r.ch <- e
}
