package tui

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode   = kernel32.NewProc("SetConsoleMode")
	getScreenInfo    = kernel32.NewProc("GetConsoleScreenBufferInfo")
	readConsoleInput = kernel32.NewProc("ReadConsoleInputW")
	getOutputCP      = kernel32.NewProc("GetConsoleOutputCP")
	setOutputCP      = kernel32.NewProc("SetConsoleOutputCP")
	waitConsoleInput = kernel32.NewProc("WaitForSingleObject")
)

type Terminal struct {
	out, in       *os.File
	oldIn, oldOut uint32
	oldCP         uintptr
	winch         chan struct{}
	size          [2]int
	done          chan struct{}
	readerDone    chan struct{}
	closeOnce     sync.Once
}

// Mouse input is enabled by Open and read as native console records.
func (t *Terminal) Mouse(on bool) {}

func consoleMode(file *os.File, mode uint32) error {
	ok, _, err := setConsoleMode.Call(file.Fd(), uintptr(mode))
	if ok == 0 {
		return err
	}
	return nil
}

func Open() (*Terminal, error) {
	t := &Terminal{in: os.Stdin, out: os.Stdout, winch: make(chan struct{}, 1), done: make(chan struct{})}
	if err := syscall.GetConsoleMode(syscall.Handle(t.in.Fd()), &t.oldIn); err != nil {
		return nil, fmt.Errorf("TERMTD needs an interactive Windows console: %w", err)
	}
	if err := syscall.GetConsoleMode(syscall.Handle(t.out.Fd()), &t.oldOut); err != nil {
		return nil, err
	}
	// Native input records provide keyboard, mouse, and resize events in both
	// Windows Terminal and the classic console. Disable Quick Edit and Ctrl+C
	// processing so clicks and quitting reach the game.
	if err := consoleMode(t.in, (t.oldIn&^(0x1|0x2|0x4|0x40|0x200))|0x80|0x8|0x10); err != nil {
		return nil, err
	}
	if err := consoleMode(t.out, t.oldOut|0x1|0x4|0x8); err != nil {
		consoleMode(t.in, t.oldIn)
		return nil, fmt.Errorf("TERMTD requires Windows 10/11 virtual terminal output: %w", err)
	}
	t.oldCP, _, _ = getOutputCP.Call()
	if ok, _, err := setOutputCP.Call(65001); ok == 0 {
		consoleMode(t.in, t.oldIn)
		consoleMode(t.out, t.oldOut)
		return nil, err
	}
	t.size = t.querySize()
	return t, nil
}

func (t *Terminal) Close() {
	t.closeOnce.Do(func() {
		close(t.done)
		if t.readerDone != nil {
			<-t.readerDone
		}
		consoleMode(t.in, t.oldIn)
		consoleMode(t.out, t.oldOut)
		setOutputCP.Call(t.oldCP)
	})
}

type consoleScreenInfo struct {
	Size, Cursor [2]int16
	Attributes   uint16
	Window       [4]int16
	Maximum      [2]int16
}

func termSize(fd int) (int, int, error) {
	var info consoleScreenInfo
	ok, _, err := getScreenInfo.Call(uintptr(fd), uintptr(unsafe.Pointer(&info)))
	if ok == 0 {
		return 0, 0, err
	}
	return int(info.Window[2] - info.Window[0] + 1), int(info.Window[3] - info.Window[1] + 1), nil
}

// INPUT_RECORD contains a WORD event type, two padding bytes, and a
// 16-byte union. Decode explicitly to preserve the Windows ABI on both CPUs.
type consoleInputRecord struct {
	Kind    uint16
	Padding uint16
	Data    [16]byte
}

func windowsKey(data [16]byte) (Event, int) {
	if binary.LittleEndian.Uint32(data[:4]) == 0 {
		return Event{}, 0
	}
	repeat := int(binary.LittleEndian.Uint16(data[4:6]))
	vk := binary.LittleEndian.Uint16(data[6:8])
	r := rune(binary.LittleEndian.Uint16(data[10:12]))
	keys := map[uint16]Key{0x0d: KeyEnter, 0x1b: KeyEscape, 0x26: KeyUp, 0x28: KeyDown, 0x25: KeyLeft, 0x27: KeyRight, 0x08: KeyBackspace, 0x09: KeyTab}
	if k := keys[vk]; k != KeyNone {
		return Event{Key: k}, repeat
	}
	if r == 3 {
		return Event{Key: KeyCtrlC}, repeat
	}
	if r == 12 {
		return Event{Key: KeyCtrlL}, repeat
	}
	if r >= 32 {
		return Event{Rune: r}, repeat
	}
	return Event{}, 0
}

func windowsMouse(data [16]byte, previous uint32) ([]Event, uint32) {
	buttons := binary.LittleEndian.Uint32(data[4:8])
	flags := binary.LittleEndian.Uint32(data[12:16])
	if flags == 4 {
		btn := 64
		if int16(buttons>>16) < 0 {
			btn = 65
		}
		return []Event{{Mouse: true, Btn: btn, Press: true, X: -1, Y: -1}}, previous
	}
	// Ignore movement, horizontal wheel, and duplicate double-click records.
	if flags != 0 {
		return nil, previous
	}
	var events []Event
	for _, b := range []struct {
		mask uint32
		btn  int
	}{{1, 0}, {4, 1}, {2, 2}} {
		if (buttons^previous)&b.mask != 0 {
			events = append(events, Event{Mouse: true, Btn: b.btn, Press: buttons&b.mask != 0,
				X: int(int16(binary.LittleEndian.Uint16(data[:2]))), Y: int(int16(binary.LittleEndian.Uint16(data[2:4])))})
		}
	}
	return events, buttons
}

func (t *Terminal) startInput(ch chan Event) {
	t.readerDone = make(chan struct{})
	go func() {
		defer close(t.readerDone)
		emit := func(e Event) bool {
			select {
			case ch <- e:
				return true
			case <-t.done:
				return false
			}
		}
		var buttons uint32
		for {
			select {
			case <-t.done:
				return
			default:
			}
			// Waiting with a timeout lets Close stop the reader before restoring
			// or closing console handles. A blocking ReadConsoleInput would keep
			// the input handle alive indefinitely after quitting.
			status, _, _ := waitConsoleInput.Call(t.in.Fd(), 100)
			if status == 258 {
				continue
			} // WAIT_TIMEOUT
			if status != 0 {
				return
			}
			var record consoleInputRecord
			var count uint32
			ok, _, _ := readConsoleInput.Call(t.in.Fd(), uintptr(unsafe.Pointer(&record)), 1, uintptr(unsafe.Pointer(&count)))
			if ok == 0 {
				return
			}
			if count == 0 {
				continue
			}
			switch record.Kind {
			case 1:
				e, repeat := windowsKey(record.Data)
				for i := 0; i < repeat; i++ {
					if !emit(e) {
						return
					}
				}
			case 2:
				events, next := windowsMouse(record.Data, buttons)
				buttons = next
				for _, e := range events {
					if !emit(e) {
						return
					}
				}
			case 4:
				select {
				case t.winch <- struct{}{}:
				default:
				}
			}
		}
	}()
}
