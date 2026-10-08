package hiscore

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestSaveRecoveryPreservesOriginals(t *testing.T) {
	cases := []struct {
		name    string
		load    func() (any, error)
		save    func() error
		invalid []string
	}{
		{"hiscores", func() (any, error) { return LoadWithError() }, func() error { return Save(Table{"hub": 99}) }, []string{`{broken`, `null`, `[]`, `{"hub":42,"rift":"bad"}`, ""}},
		{"journal", func() (any, error) { return LoadJournalWithError() }, func() error { return SaveJournal(&Journal{EndgameWins: 99}) }, []string{`{broken`, `null`, `[]`, `{"endgame_wins":42,"towers":false}`, ""}},
		{"lair", func() (any, error) { return LoadLairWithError() }, func() error { return SaveLair(&Lair{Tokens: 99}) }, []string{`{broken`, `null`, `[]`, `{"tokens":42,"floors":false}`, ""}},
	}
	for _, c := range cases {
		for _, prefix := range []string{".termtd-", ".tdef-"} {
			for _, invalid := range c.invalid {
				t.Run(c.name+prefix+invalid, func(t *testing.T) {
					home := isolateHome(t)
					empty, err := c.load()
					if err != nil {
						t.Fatal(err)
					}
					path := filepath.Join(home, prefix+c.name+".json")
					if err := os.WriteFile(path, []byte(invalid), 0600); err != nil {
						t.Fatal(err)
					}
					got, err := c.load()
					if err == nil || !strings.Contains(err.Error(), "recovery copy:") {
						t.Fatalf("missing recovery report: %v", err)
					}
					if !reflect.DeepEqual(got, empty) {
						t.Fatalf("partial progress escaped decode: %#v", got)
					}
					copies, err := filepath.Glob(path + ".corrupt-*")
					if err != nil || len(copies) != 1 {
						t.Fatalf("copies: %v, %v", copies, err)
					}
					copyData, err := os.ReadFile(copies[0])
					if err != nil || string(copyData) != invalid {
						t.Fatalf("recovery bytes differ: %q, %v", copyData, err)
					}
					if err := c.save(); err == nil {
						t.Fatal("overwrote invalid save")
					}
					data, err := os.ReadFile(path)
					if err != nil || string(data) != invalid {
						t.Fatalf("original changed: %q, %v", data, err)
					}
					if prefix == ".tdef-" {
						if _, err := os.Stat(filepath.Join(home, ".termtd-"+c.name+".json")); !os.IsNotExist(err) {
							t.Fatal("masked damaged legacy save")
						}
					}
					// A second load must never overwrite an existing recovery copy.
					c.load()
					copies, _ = filepath.Glob(path + ".corrupt-*")
					if len(copies) != 2 {
						t.Fatal("recovery copy was reused")
					}
					// Repair is explicit; old recovery copies remain available.
					if err := os.WriteFile(path, []byte(`{}`), 0600); err != nil {
						t.Fatal(err)
					}
					if _, err := c.load(); err != nil {
						t.Fatal(err)
					}
					if err := c.save(); err != nil {
						t.Fatal(err)
					}
					if _, err := os.Stat(copies[0]); err != nil {
						t.Fatal("repair removed recovery copy")
					}
				})
			}
		}
	}
}

func TestUnreadableSavesDoNotFallBackOrGetReplaced(t *testing.T) {
	home := isolateHome(t)
	for _, c := range []struct {
		name string
		load func() error
		save func() error
	}{
		{"lair", func() error { _, err := LoadLairWithError(); return err }, func() error { return SaveLair(emptyLair()) }},
		{"journal", func() error { _, err := LoadJournalWithError(); return err }, func() error { return SaveJournal(&Journal{}) }},
		{"hiscores", func() error { _, err := LoadWithError(); return err }, func() error { return Save(Table{}) }},
	} {
		path := filepath.Join(home, ".termtd-"+c.name+".json")
		if err := os.Mkdir(path, 0700); err != nil {
			t.Fatal(err)
		}
		legacy := filepath.Join(home, ".tdef-"+c.name+".json")
		if err := os.WriteFile(legacy, []byte(`{}`), 0600); err != nil {
			t.Fatal(err)
		}
		if err := c.load(); err == nil || !strings.Contains(err.Error(), "cannot read") {
			t.Fatalf("read error hidden: %v", err)
		}
		if err := c.save(); err == nil {
			t.Fatal("unreadable original replaced")
		}
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			t.Fatal("original changed")
		}
	}
}

func TestRecoveryCopyFailureKeepsOriginal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX directory permissions")
	}
	home := isolateHome(t)
	path := filepath.Join(home, ".termtd-lair.json")
	if err := os.WriteFile(path, []byte(`{broken`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(home, 0500); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(home, 0700)
	probe, err := os.CreateTemp(home, "probe")
	if err == nil {
		probe.Close()
		os.Remove(probe.Name())
		t.Skip("user bypasses directory permissions")
	}
	_, err = LoadLairWithError()
	if err == nil || !strings.Contains(err.Error(), "recovery copy failed") {
		t.Fatalf("copy failure hidden: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != `{broken` {
		t.Fatal("original changed")
	}
}

func TestDanglingCurrentSaveDoesNotFallBack(t *testing.T) {
	home := isolateHome(t)
	path := filepath.Join(home, ".termtd-lair.json")
	if err := os.Symlink(filepath.Join(home, "missing"), path); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.WriteFile(legacyPath(path), []byte(`{"tokens":42}`), 0600); err != nil {
		t.Fatal(err)
	}
	value, err := LoadLairWithError()
	if err == nil || value.Tokens != 0 {
		t.Fatalf("dangling save treated as absent: %+v, %v", value, err)
	}
	if err := SaveLair(emptyLair()); err == nil {
		t.Fatal("dangling save overwritten")
	}
}
