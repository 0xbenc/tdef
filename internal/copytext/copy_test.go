package copytext

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestCopyValidation(t *testing.T) {
	contracts := map[string][]string{"sample.result": {"floor", "wave"}}
	for _, tc := range []struct {
		name, json, want string
	}{
		{"valid", `{"result":"{floor} held wave {wave}. {floor} is safe."}`, ""},
		{"missing entry", `{}`, "missing or empty"},
		{"empty entry", `{"result":" "}`, "missing or empty"},
		{"missing placeholder", `{"result":"{floor} is safe."}`, "placeholders"},
		{"misspelled placeholder", `{"result":"{room} held {wave}"}`, "placeholders"},
		{"broken placeholder", `{"result":"{floor} held {wave} {oops"}`, "invalid placeholder"},
		{"extra entry", `{"result":"{floor} {wave}","typo":"Unused"}`, "unknown copy key"},
		{"duplicate entry", `{"result":"first","result":"{floor} {wave}"}`, "duplicate key"},
		{"duplicate group", `{"nested":{"line":"first"},"nested":{"line":"second"}}`, "duplicate key"},
		{"number", `{"result":5}`, "expected text or an object"},
		{"null", `{"result":null}`, "expected text or an object"},
		{"array", `{"result":["{floor} {wave}"]}`, "expected text or an object"},
		{"trailing object", `{"result":"{floor} {wave}"} {}`, "one JSON object"},
		{"invalid JSON", `{"result":"{floor} {wave}`, "unexpected EOF"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := fstest.MapFS{"data/sample.json": {Data: []byte(tc.json)}}
			_, err := load(source, contracts)
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestEmbeddedCopy(t *testing.T) {
	if _, err := load(files, required); err != nil {
		t.Fatal(err)
	}
}

func TestFormatAllowsReorderedAndRepeatedValues(t *testing.T) {
	const key = "overworld.results.lost"
	old := catalog[key]
	catalog[key] = "Wave {wave}: {floor}. Hold {floor} next time."
	t.Cleanup(func() { catalog[key] = old })
	got := Format(key, "floor", "a room called {wave}", "wave", "17")
	want := "Wave 17: a room called {wave}. Hold a room called {wave} next time."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatRejectsBadBindings(t *testing.T) {
	for _, pairs := range [][]string{
		{"floor", "Rotunda"},
		{"floor", "Rotunda", "typo", "17"},
		{"floor", "Rotunda", "floor", "Rift"},
		{"floor"},
	} {
		t.Run(strings.Join(pairs, "/"), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("bad bindings must fail instead of displaying broken copy")
				}
			}()
			Format("overworld.results.lost", pairs...)
		})
	}
}
