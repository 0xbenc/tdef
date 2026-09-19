package tui

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
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
	KeyTab
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
			// State 1 means we saw an ESC and are waiting to see whether a
			// CSI/SS3 sequence follows. Arm a deadline so a bare Escape is
			// not swallowed until the next unrelated keypress.
			if r.state == 1 {
				_ = in.SetReadDeadline(time.Now().Add(escWindow))
			} else {
				_ = in.SetReadDeadline(time.Time{})
			}
			n, err := in.Read(tmp)
			if err == io.EOF {
				return
			}
			if n > 0 {
				r.feed(tmp[:n])
				continue
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				r.state = 0
				r.emit(Event{Key: KeyEscape})
				continue
			}
			if err != nil {
				return
			}
		}
	}()
}

func (r *reader) feed(data []byte) {
	for len(data) > 0 {
		data = r.step(data)
	}
}

func (r *reader) step(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	switch r.state {
	case 0:
		b := data[0]
		data = data[1:]
		switch {
		case b == 0x1b:
			r.state = 1
			if len(data) == 0 {
				return data
			}
			return r.step(data)
		case b == '\r' || b == '\n':
			r.emit(Event{Key: KeyEnter})
		case b == 0x7f || b == 0x08:
			r.emit(Event{Key: KeyBackspace})
		case b == 0x03:
			r.emit(Event{Key: KeyCtrlC})
		case b == 0x0c:
			r.emit(Event{Key: KeyCtrlL})
		case b == '\t':
			r.emit(Event{Key: KeyTab})
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
		if len(data) == 0 {
			return data
		}
		b := data[0]
		data = data[1:]
		if b == 0x1b {
			r.state = 0
			r.emit(Event{Key: KeyEscape})
			return r.step(data)
		}
		if b == '[' || b == 'O' {
			r.state = 2
			r.buf = r.buf[:0]
			return r.step(data)
		}
		r.state = 0
		r.emit(Event{Rune: rune(b)})
		return r.step(data)
	case 2:
		if len(data) == 0 {
			return data
		}
		b := data[0]
		data = data[1:]
		if b == '<' {
			r.state = 3
			r.buf = r.buf[:0]
			return r.step(data)
		}
		r.buf = append(r.buf, b)
		if k, done := r.interpret(); done {
			r.state = 0
			if k != KeyNone {
				r.emit(Event{Key: k})
			}
			return r.step(data)
		}
		return r.step(data)
	case 3:
		if len(data) == 0 {
			return data
		}
		b := data[0]
		data = data[1:]
		if b == 'M' || b == 'm' {
			r.state = 0
			r.mouse()
			return r.step(data)
		}
		r.buf = append(r.buf, b)
		return r.step(data)
	}
	return data
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

func (r *reader) mouse() {
	s := string(r.buf)
	parts := splitN(s, ';')
	if len(parts) < 3 {
		return
	}
	btn, _ := strconv.Atoi(parts[0])
	x, _ := strconv.Atoi(parts[1])
	y, _ := strconv.Atoi(parts[2])
	if x == 0 || y == 0 {
		return
	}
	press := !bytes.HasSuffix([]byte(s), []byte("m"))
	if btn < 64 {
		r.emit(Event{Mouse: true, Btn: btn, Press: press, X: x - 1, Y: y - 1})
	} else if press {
		r.emit(Event{Mouse: true, Btn: btn, Press: true, X: -1, Y: -1})
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
