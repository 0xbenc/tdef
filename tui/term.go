package tui

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

type Terminal struct {
	out   *os.File
	in    *os.File
	old   syscall.Termios
	winch chan struct{}
	size  [2]int
	origW int
}

func ioctl(fd int, req uintptr, arg unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(arg))
	if errno != 0 {
		return errno
	}
	return nil
}

func Open() (*Terminal, error) {
	t := &Terminal{out: os.Stdout, in: os.Stdin, winch: make(chan struct{}, 1)}
	termios, err := getTermios(int(t.in.Fd()))
	if err != nil {
		return nil, err
	}
	t.old = *termios
	raw := *termios
	raw.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if err := setTermios(int(t.in.Fd()), &raw); err != nil {
		return nil, err
	}
	t.size = t.querySize()
	t.origW = t.size[0]
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	go func() {
		for range sig {
			t.size = t.querySize()
			select {
			case t.winch <- struct{}{}:
			default:
			}
		}
	}()
	return t, nil
}

func (t *Terminal) Close() {
	setTermios(int(t.in.Fd()), &t.old)
}

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

func (t *Terminal) Mouse(on bool) {
	if on {
		t.Write([]byte("\x1b[?1000h\x1b[?1002h\x1b[?1005h\x1b[?1006h"))
	} else {
		t.Write([]byte("\x1b[?1006l\x1b[?1005l\x1b[?1002l\x1b[?1000l"))
	}
}

func getTermios(fd int) (*syscall.Termios, error) {
	var t syscall.Termios
	if err := ioctl(fd, syscall.TCGETS, unsafe.Pointer(&t)); err != nil {
		return nil, err
	}
	return &t, nil
}

func setTermios(fd int, t *syscall.Termios) error {
	return ioctl(fd, syscall.TCSETS, unsafe.Pointer(t))
}

func termSize(fd int) (int, int, error) {
	var ws struct{ Row, Col, X, Y uint16 }
	if err := ioctl(fd, syscall.TIOCGWINSZ, unsafe.Pointer(&ws)); err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), nil
}
