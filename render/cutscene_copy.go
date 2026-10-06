package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Cinematic copy is embedded so packaged games need no external content files.
// Edit copy/opening.json or copy/ending.json and rebuild to change the writing.
//
//go:embed copy/opening.json
var openingCopy []byte

//go:embed copy/ending.json
var endingCopy []byte

func withOpeningCopy(shots []FilmShot) []FilmShot {
	return withFilmCopy("opening", openingCopy, shots)
}

func withEndingCopy(shots []FilmShot) []FilmShot {
	return withFilmCopy("ending", endingCopy, shots)
}

func withFilmCopy(name string, data []byte, shots []FilmShot) []FilmShot {
	var copy []struct {
		Title   string `json:"title"`
		Speaker string `json:"speaker"`
		Text    string `json:"text"`
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&copy); err != nil {
		panic(fmt.Sprintf("%s cinematic copy: %v", name, err))
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		panic(fmt.Sprintf("%s cinematic copy: unexpected content after entries", name))
	}
	if len(copy) != len(shots) {
		panic(fmt.Sprintf("%s cinematic copy: expected %d entries, got %d", name, len(shots), len(copy)))
	}
	for i, entry := range copy {
		if strings.TrimSpace(entry.Title) == "" || strings.TrimSpace(entry.Text) == "" {
			panic(fmt.Sprintf("%s cinematic copy: entry %d needs a title and text", name, i+1))
		}
		switch entry.Speaker {
		case "", "GRAK", "MALGRATH", "GUILDMASTER", "HEALER":
		default:
			panic(fmt.Sprintf("%s cinematic copy: entry %d has unknown speaker %q", name, i+1, entry.Speaker))
		}
		shots[i].Name, shots[i].Speaker, shots[i].Dialogue = entry.Title, entry.Speaker, entry.Text
	}
	return shots
}
