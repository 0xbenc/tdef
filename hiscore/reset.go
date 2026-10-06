package hiscore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ResetProgress removes all current and legacy saves so old progress cannot
// reappear through the rename fallback. Unrelated files are left untouched.
func ResetProgress() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	home, err = filepath.Abs(home)
	if err != nil {
		return err
	}
	var paths []string
	for _, prefix := range []string{".termtd-", ".tdef-"} {
		for _, name := range []string{"hiscores", "journal", "lair"} {
			path := filepath.Join(home, prefix+name+".json")
			info, err := os.Lstat(path)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			if err == nil && info.IsDir() {
				return fmt.Errorf("save path is a directory: %s", path)
			}
			paths = append(paths, path)
		}
	}
	var failures []error
	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}
