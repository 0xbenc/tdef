//go:build linux || darwin

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
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	go func() {
		// Pure notifier: the size is re-queried by the main goroutine
		// (RefreshSize) after it sees the token. Writing t.size here
		// would race with Size() on the main loop.
		for range sig {
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

func (t *Terminal) startInput(ch chan Event) { startReader(t.in, ch) }

func (t *Terminal) Mouse(on bool) { t.mouseVT(on) }

func termSize(fd int) (int, int, error) {
	var ws struct{ Row, Col, X, Y uint16 }
	if err := ioctl(fd, syscall.TIOCGWINSZ, unsafe.Pointer(&ws)); err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), nil
}
