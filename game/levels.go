package game

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed levels/*.txt
var levelFS embed.FS

func LevelNames() []string {
	entries, err := levelFS.ReadDir("levels")
	if err != nil {
		return nil
	}
	names := []string{}
	for _, e := range entries {
		names = append(names, strings.TrimSuffix(e.Name(), ".txt"))
	}
	sort.Strings(names)
	return names
}

var levelHP = map[string]float64{
	// Tuned so the autoplay (weak-player) proxy wins each hand-crafted map at
	// ~100% with a spread of leaks: hub easiest -> canyon hardest.
	"hub":     1.08,
	"winding": 1.32,
	"garden":  1.28,
	"canyon":  1.25,
}

// mazeHP keeps the procedural maze in the "chaos" tier (~60% AI win). The
// corridor-biased generator yields longer chokepoints but fewer effective
// build pockets, so it needs a lower HP than the hand-crafted maps.
const mazeHP = 0.58

func LoadLevel(name string) (*Map, error) {
	data, err := levelFS.ReadFile("levels/" + name + ".txt")
	if err != nil {
		return nil, fmt.Errorf("unknown level %q (have: %s)", name, strings.Join(LevelNames(), ", "))
	}
	rows := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	for i := range rows {
		rows[i] = strings.TrimRight(rows[i], "\r")
	}
	w := len(rows[0])
	m, err := LoadMap(w, len(rows), rows)
	if err != nil {
		return nil, err
	}
	m.HPMul = 1.0
	if v, ok := levelHP[name]; ok {
		m.HPMul = v
	}
	return m, nil
}
