package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Homebrew runs --version outside a TTY. Verify release metadata is embedded
// in the executable and none of the version aliases opens the terminal.
func TestReleaseVersionCLI(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "termtd")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if out, err := exec.Command("go", "build", "-ldflags=-X main.version=1.0.0", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("release build: %v\n%s", err, out)
	}
	for _, arg := range []string{"version", "--version", "-version"} {
		out, err := exec.Command(binary, arg).CombinedOutput()
		if err != nil || string(out) != "termtd 1.0.0\n" {
			t.Errorf("%s: output %q, error %v", arg, out, err)
		}
	}
}
