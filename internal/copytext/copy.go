// Package copytext loads authored game copy shared by simulation and screens.
package copytext

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"strings"
)

//go:embed data/*.json
var files embed.FS

var placeholders = regexp.MustCompile(`\{([a-z][a-z0-9_]*)\}`)
var catalog = mustLoad()

func mustLoad() map[string]string {
	c, err := load(files, required)
	if err != nil {
		panic("game copy: " + err.Error())
	}
	return c
}

func load(source fs.FS, contracts map[string][]string) (map[string]string, error) {
	names, err := fs.Glob(source, "data/*.json")
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, name := range names {
		data, err := fs.ReadFile(source, name)
		if err != nil {
			return nil, err
		}
		dec := json.NewDecoder(strings.NewReader(string(data)))
		domain := strings.TrimSuffix(path.Base(name), ".json")
		if err := readObject(dec, out, domain); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if err := dec.Decode(new(any)); err != io.EOF {
			return nil, fmt.Errorf("%s: expected one JSON object", name)
		}
	}
	for key, wanted := range contracts {
		text, ok := out[key]
		if !ok || strings.TrimSpace(text) == "" {
			return nil, fmt.Errorf("%s: missing or empty copy", key)
		}
		seen := map[string]bool{}
		for _, match := range placeholders.FindAllStringSubmatch(text, -1) {
			seen[match[1]] = true
		}
		actual := make([]string, 0, len(seen))
		for param := range seen {
			actual = append(actual, param)
		}
		slices.Sort(actual)
		if !slices.Equal(actual, wanted) {
			return nil, fmt.Errorf("%s: placeholders must be %v, got %v", key, wanted, actual)
		}
		remaining := placeholders.ReplaceAllString(text, "")
		if strings.ContainsAny(remaining, "{}") {
			return nil, fmt.Errorf("%s: invalid placeholder", key)
		}
	}
	for key := range out {
		if _, ok := contracts[key]; !ok {
			return nil, fmt.Errorf("%s: unknown copy key", key)
		}
	}
	return out, nil
}

func readObject(dec *json.Decoder, out map[string]string, prefix string) error {
	opening, err := dec.Token()
	if err != nil || opening != json.Delim('{') {
		return fmt.Errorf("%s: expected an object", prefix)
	}
	return readMembers(dec, out, prefix)
}

func readMembers(dec *json.Decoder, out map[string]string, prefix string) error {
	seen := map[string]bool{}
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok || key == "" || strings.Contains(key, ".") {
			return fmt.Errorf("%s: invalid key %v", prefix, token)
		}
		name := prefix + "." + key
		if seen[key] {
			return fmt.Errorf("%s: duplicate key", name)
		}
		seen[key] = true
		value, err := dec.Token()
		if err != nil {
			return err
		}
		switch value := value.(type) {
		case string:
			out[name] = value
		case json.Delim:
			if value != '{' {
				return fmt.Errorf("%s: expected text or an object", name)
			}
			if err := readMembers(dec, out, name); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: expected text or an object", name)
		}
	}
	closing, err := dec.Token()
	if err != nil || closing != json.Delim('}') {
		return fmt.Errorf("%s: expected closing brace", prefix)
	}
	return nil
}

// Text returns a required static entry. Bad keys fail loudly during development.
func Text(key string) string {
	text, ok := catalog[key]
	if !ok {
		panic("unknown game copy: " + key)
	}
	if len(required[key]) != 0 {
		panic("game copy needs placeholders: " + key)
	}
	return text
}

// Format replaces named placeholders using pairs of names and preformatted
// values. Gameplay owns the values and numeric formatting; JSON owns wording.
func Format(key string, pairs ...string) string {
	text, ok := catalog[key]
	if !ok || len(pairs)%2 != 0 {
		panic("invalid game copy format: " + key)
	}
	values := make(map[string]string, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		if _, duplicate := values[pairs[i]]; duplicate || !slices.Contains(required[key], pairs[i]) {
			panic("unexpected game copy placeholder: " + key + "." + pairs[i])
		}
		values[pairs[i]] = pairs[i+1]
	}
	if len(values) != len(required[key]) {
		panic("missing game copy placeholder: " + key)
	}
	// Replace only the template, never braces that occur in supplied values.
	return placeholders.ReplaceAllStringFunc(text, func(match string) string {
		return values[match[1:len(match)-1]]
	})
}
