package hiscore

import (
	"os"
	"path/filepath"
	"strings"
)

// readSave preserves progress from before the project rename. New saves always
// use termtd paths; an existing new save takes precedence, even if corrupt.
func readSave(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		legacy := strings.Replace(filepath.Base(path), ".termtd-", ".tdef-", 1)
		return os.ReadFile(filepath.Join(filepath.Dir(path), legacy))
	}
	return data, err
}
