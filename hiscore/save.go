package hiscore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func legacyPath(path string) string {
	return filepath.Join(filepath.Dir(path), strings.Replace(filepath.Base(path), ".termtd-", ".tdef-", 1))
}

func readSaveSource(path string) ([]byte, string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// A dangling symlink is an existing save, not a new player.
		if _, statErr := os.Lstat(path); statErr == nil || !os.IsNotExist(statErr) {
			return nil, path, err
		}
		path = legacyPath(path)
		data, err = os.ReadFile(path)
	}
	return data, path, err
}

// loadSave never exposes a partially decoded value. Callers decode into a
// temporary value and discard it on error. Missing files are normal first runs.
func loadSave(path string, value any) error {
	data, source, err := readSaveSource(path)
	if err != nil {
		if os.IsNotExist(err) {
			if _, statErr := os.Lstat(source); os.IsNotExist(statErr) {
				return nil
			}
		}
		return fmt.Errorf("cannot read %s: %w", source, err)
	}
	if err = decodeSave(data, value); err == nil {
		return nil
	}
	backup, backupErr := preserveSave(source, data)
	if backupErr != nil {
		return fmt.Errorf("cannot load %s: %w; recovery copy failed: %v; original untouched", source, err, backupErr)
	}
	return fmt.Errorf("cannot load %s: %w; recovery copy: %s; original untouched", source, err, backup)
}

func decodeSave(data []byte, value any) error {
	if !bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
		return fmt.Errorf("expected a JSON object")
	}
	return json.Unmarshal(data, value)
}

// Use an exclusive random name: previous recovery copies are never overwritten.
func preserveSave(source string, data []byte) (string, error) {
	f, err := os.CreateTemp(filepath.Dir(source), filepath.Base(source)+".corrupt-*")
	if err != nil {
		return "", err
	}
	name := f.Name()
	ok := false
	defer func() {
		f.Close()
		if !ok {
			os.Remove(name)
		}
	}()
	if _, err = f.Write(data); err != nil {
		return "", err
	}
	if err = f.Sync(); err != nil {
		return "", err
	}
	if err = f.Close(); err != nil {
		return "", err
	}
	ok = true
	return name, nil
}

// Refuse to replace an unreadable or invalid current/legacy save, including
// through the compatibility APIs that do not return load errors.
func checkSave(path string, value any) error {
	data, source, err := readSaveSource(path)
	if err != nil {
		if os.IsNotExist(err) {
			if _, statErr := os.Lstat(source); os.IsNotExist(statErr) {
				return nil
			}
		}
		return fmt.Errorf("save blocked: cannot read %s: %w", source, err)
	}
	if err := decodeSave(data, value); err != nil {
		return fmt.Errorf("save blocked: repair %s or reset progress: %w", source, err)
	}
	return nil
}
