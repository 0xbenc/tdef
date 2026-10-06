package tui

import (
	"context"
	"encoding/binary"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestWindowsConsoleLifecycle(t *testing.T) {
	if os.Getenv("TERMTD_CONSOLE_TEST") != "1" {
		binary, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary, "-test.run=^TestWindowsConsoleLifecycle$", "-test.timeout=10s")
		cmd.Env = append(os.Environ(), "TERMTD_CONSOLE_TEST=1")
		// CREATE_NO_WINDOW supplies a real console without a visible window.
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("hidden console: %v\n%s", err, out)
		}
		return
	}
	in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	var oldIn, oldOut uint32
	if err := syscall.GetConsoleMode(syscall.Handle(in.Fd()), &oldIn); err != nil {
		t.Fatal(err)
	}
	if err := syscall.GetConsoleMode(syscall.Handle(out.Fd()), &oldOut); err != nil {
		t.Fatal(err)
	}
	oldCP, _, _ := getOutputCP.Call()
	stdin, stdout := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = in, out
	term, err := Open()
	os.Stdin, os.Stdout = stdin, stdout
	if err != nil {
		t.Fatal(err)
	}
	defer term.Close()
	w, h := term.Size()
	if w <= 0 || h <= 0 {
		t.Fatalf("size: %dx%d", w, h)
	}
	term.AltScreen(true)
	term.Cursor(false)
	term.Write([]byte("Grak ░ █\r\n"))
	term.Cursor(true)
	term.AltScreen(false)
	ch := make(chan Event, 8)
	term.startInput(ch)
	key := consoleInputRecord{Kind: 1}
	binary.LittleEndian.PutUint32(key.Data[:4], 1)
	binary.LittleEndian.PutUint16(key.Data[4:6], 1)
	binary.LittleEndian.PutUint16(key.Data[6:8], 0x26)
	mouse := consoleInputRecord{Kind: 2}
	binary.LittleEndian.PutUint16(mouse.Data[:2], 5)
	binary.LittleEndian.PutUint16(mouse.Data[2:4], 3)
	binary.LittleEndian.PutUint32(mouse.Data[4:8], 1)
	records := []consoleInputRecord{key, mouse, {Kind: 4}}
	var written uint32
	ok, _, err := kernel32.NewProc("WriteConsoleInputW").Call(in.Fd(), uintptr(unsafe.Pointer(&records[0])), uintptr(len(records)), uintptr(unsafe.Pointer(&written)))
	if ok == 0 || written != 3 {
		t.Fatalf("inject console input: %v (%d)", err, written)
	}
	for _, want := range []Event{{Key: KeyUp}, {Mouse: true, Btn: 0, Press: true, X: 5, Y: 3}} {
		select {
		case got := <-ch:
			if got != want {
				t.Fatalf("input: %+v, want %+v", got, want)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("console input timed out")
		}
	}
	select {
	case <-term.Winch():
		term.RefreshSize()
	case <-time.After(3 * time.Second):
		t.Fatal("resize timed out")
	}
	term.Close()
	var restoredIn, restoredOut uint32
	syscall.GetConsoleMode(syscall.Handle(in.Fd()), &restoredIn)
	syscall.GetConsoleMode(syscall.Handle(out.Fd()), &restoredOut)
	cp, _, _ := getOutputCP.Call()
	if restoredIn != oldIn || restoredOut != oldOut || cp != oldCP {
		t.Fatal("console settings were not restored")
	}
}

func TestWindowsConsoleABI(t *testing.T) {
	if unsafe.Sizeof(consoleInputRecord{}) != 20 || unsafe.Sizeof(consoleScreenInfo{}) != 22 {
		t.Fatal("console structs must match the Windows ABI")
	}
}

func TestWindowsKeys(t *testing.T) {
	for _, tc := range []struct {
		vk   uint16
		char uint16
		key  Key
		r    rune
	}{
		{0x26, 0, KeyUp, 0}, {0x1b, 27, KeyEscape, 0},
		{0x0d, 13, KeyEnter, 0}, {0x43, 3, KeyCtrlC, 0},
		{0x4c, 12, KeyCtrlL, 0}, {0x57, 'w', KeyNone, 'w'},
	} {
		var data [16]byte
		binary.LittleEndian.PutUint32(data[:4], 1)
		binary.LittleEndian.PutUint16(data[4:6], 3)
		binary.LittleEndian.PutUint16(data[6:8], tc.vk)
		binary.LittleEndian.PutUint16(data[10:12], tc.char)
		e, repeat := windowsKey(data)
		if e.Key != tc.key || e.Rune != tc.r || repeat != 3 {
			t.Fatalf("unexpected key: %+v repeat=%d", e, repeat)
		}
		data[0] = 0
		if _, repeat := windowsKey(data); repeat != 0 {
			t.Fatal("key release emitted input")
		}
	}
}

func TestWindowsMouse(t *testing.T) {
	var data [16]byte
	binary.LittleEndian.PutUint16(data[:2], 12)
	binary.LittleEndian.PutUint16(data[2:4], 7)
	binary.LittleEndian.PutUint32(data[4:8], 1)
	events, buttons := windowsMouse(data, 0)
	if len(events) != 1 || events[0] != (Event{Mouse: true, Btn: 0, Press: true, X: 12, Y: 7}) {
		t.Fatalf("click: %+v", events)
	}
	binary.LittleEndian.PutUint32(data[12:16], 1)
	if events, _ := windowsMouse(data, buttons); len(events) != 0 {
		t.Fatal("drag emitted click")
	}
	binary.LittleEndian.PutUint32(data[12:16], 0)
	binary.LittleEndian.PutUint32(data[4:8], 0)
	events, _ = windowsMouse(data, buttons)
	if len(events) != 1 || events[0].Press {
		t.Fatalf("release: %+v", events)
	}
	for _, tc := range []struct {
		delta int16
		btn   int
	}{{120, 64}, {-120, 65}} {
		binary.LittleEndian.PutUint32(data[12:16], 4)
		binary.LittleEndian.PutUint32(data[4:8], uint32(uint16(tc.delta))<<16)
		events, _ := windowsMouse(data, 0)
		if len(events) != 1 || events[0].Btn != tc.btn || events[0].X != -1 {
			t.Fatalf("wheel: %+v", events)
		}
	}
}
