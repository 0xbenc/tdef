package tui

import (
	"os"
	"testing"
	"time"
)

func feed(t *testing.T, data []byte, maxEvents int) []Event {
	t.Helper()
	ch := make(chan Event, 64)
	r := &reader{ch: ch}
	r.feed(data)
	out := []Event{}
	deadline := time.Now().Add(200 * time.Millisecond)
	for len(out) < maxEvents {
		select {
		case e := <-ch:
			out = append(out, e)
		case <-time.After(time.Until(deadline)):
			return out
		}
	}
	return out
}

func TestArrowsAndKeys(t *testing.T) {
	data := []byte{0x1b, '[', 'D', 0x1b, '[', 'A', '1', '\r', 'n'}
	evs := feed(t, data, 5)
	if len(evs) != 5 {
		t.Fatalf("got %d events, want 5: %+v", len(evs), evs)
	}
	if evs[0].Key != KeyLeft {
		t.Errorf("ev0 = %v, want KeyLeft", evs[0].Key)
	}
	if evs[1].Key != KeyUp {
		t.Errorf("ev1 = %v, want KeyUp", evs[1].Key)
	}
	if evs[2].Rune != '1' {
		t.Errorf("ev2 = %v, want '1'", evs[2])
	}
	if evs[3].Key != KeyEnter {
		t.Errorf("ev3 = %v, want KeyEnter", evs[3].Key)
	}
	if evs[4].Rune != 'n' {
		t.Errorf("ev4 = %v, want 'n'", evs[4])
	}
}

func TestSplitEscape(t *testing.T) {
	ch := make(chan Event, 64)
	r := &reader{ch: ch}
	r.feed([]byte{0x1b})
	r.feed([]byte{'['})
	r.feed([]byte{'C'})
	select {
	case e := <-ch:
		if e.Key != KeyRight {
			t.Fatalf("got %v, want KeyRight", e.Key)
		}
	default:
		t.Fatal("no event for split escape sequence")
	}
}

func TestDoubleEscape(t *testing.T) {
	evs := feed(t, []byte{0x1b, 0x1b}, 1)
	if len(evs) != 1 || evs[0].Key != KeyEscape {
		t.Fatalf("got %+v, want single KeyEscape", evs)
	}
}

// A real terminal sends a single 0x1b byte for the Escape key. The reader
// must emit KeyEscape after escWindow without any further input.
func TestBareEscapeTimeout(t *testing.T) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer pr.Close()
	defer pw.Close()
	ch := make(chan Event, 64)
	startReader(pr, ch)
	if _, err := pw.Write([]byte{0x1b}); err != nil {
		t.Fatal(err)
	}
	select {
	case e := <-ch:
		if e.Key != KeyEscape {
			t.Fatalf("got %+v, want KeyEscape", e)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no KeyEscape emitted for bare ESC")
	}
}

// A bare ESC must not eat the next keypress: after the timeout fires, 'x'
// arrives as a normal rune event.
func TestBareEscapeDoesNotSwallowNextKey(t *testing.T) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer pr.Close()
	defer pw.Close()
	ch := make(chan Event, 64)
	startReader(pr, ch)
	if _, err := pw.Write([]byte{0x1b}); err != nil {
		t.Fatal(err)
	}
	select {
	case <-ch: // the KeyEscape
	case <-time.After(2 * time.Second):
		t.Fatal("no KeyEscape emitted for bare ESC")
	}
	if _, err := pw.Write([]byte{'x'}); err != nil {
		t.Fatal(err)
	}
	select {
	case e := <-ch:
		if e.Rune != 'x' {
			t.Fatalf("got %+v, want rune 'x'", e)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("'x' swallowed after bare ESC")
	}
}

func TestMouseSGR(t *testing.T) {
	evs := feed(t, []byte("\x1b[<0;15;5M"), 1)
	if len(evs) != 1 {
		t.Fatalf("got %d events", len(evs))
	}
	e := evs[0]
	if !e.Mouse || e.Btn != 0 || !e.Press || e.X != 14 || e.Y != 4 {
		t.Fatalf("bad mouse event: %+v", e)
	}
}

func TestMouseWheel(t *testing.T) {
	evs := feed(t, []byte("\x1b[<64;1;1M"), 1)
	if len(evs) != 1 || evs[0].Btn != 64 || !evs[0].Press {
		t.Fatalf("bad wheel event: %+v", evs)
	}
}

func TestEnterVariants(t *testing.T) {
	for _, b := range []byte{'\r', '\n'} {
		evs := feed(t, []byte{b}, 1)
		if len(evs) != 1 || evs[0].Key != KeyEnter {
			t.Fatalf("byte %q: got %+v, want KeyEnter", b, evs)
		}
	}
}

func TestArrowBurst(t *testing.T) {
	var data []byte
	for i := 0; i < 12; i++ {
		data = append(data, 0x1b, '[', 'D')
	}
	evs := feed(t, data, 12)
	if len(evs) != 12 {
		t.Fatalf("got %d events, want 12", len(evs))
	}
	for i, e := range evs {
		if e.Key != KeyLeft {
			t.Fatalf("ev%d = %v, want KeyLeft", i, e.Key)
		}
	}
}

func TestUnknownSequenceConsumed(t *testing.T) {
	for _, seq := range []string{"\x1bOP", "\x1b[15~", "\x1b[5;5~"} {
		data := append([]byte(seq), 'a')
		evs := feed(t, data, 1)
		if len(evs) != 1 || evs[0].Rune != 'a' {
			t.Fatalf("seq %q: got %+v, want single 'a' after consumed seq", seq, evs)
		}
	}
}

var _ = os.Stdin
