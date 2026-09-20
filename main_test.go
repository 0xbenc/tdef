package main

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// capture must be environment-deterministic: the -scale flag pins the
// virtual terminal size, so text frames have one exact width and height.
func captureFrameSizes(t *testing.T, scale, wantW, wantH int) {
	t.Helper()
	dir := t.TempDir()
	capture([]string{"-level", "winding", "-out", dir, "-every", "10", "-limit", "30", "-text", "-scale", itoa(scale)})
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("wrote %d frames, want 3", len(entries))
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	sort.Strings(names)
	f, err := os.Open(filepath.Join(dir, names[0]))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	lines := 0
	for scan.Scan() {
		if got := len([]rune(scan.Text())); got != wantW {
			t.Fatalf("scale %d frame row = %d runes, want %d", scale, got, wantW)
		}
		lines++
	}
	if lines != wantH {
		t.Fatalf("scale %d frame = %d rows, want %d", scale, lines, wantH)
	}
}

func itoa(n int) string {
	return string(rune('0' + n))
}

func TestCaptureScale1FrameSize(t *testing.T) {
	captureFrameSizes(t, 1, 62, 19)
}

func TestCaptureScale2FrameSize(t *testing.T) {
	captureFrameSizes(t, 2, 92, 32)
}
