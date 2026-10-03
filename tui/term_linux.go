package tui

import (
	"syscall"
	"unsafe"
)

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
