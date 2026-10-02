package game

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed levels/*.txt
var levelFS embed.FS

//go:embed levels/boss/*.txt
var bossFS embed.FS

func LevelNames() []string {
	entries, err := levelFS.ReadDir("levels")
	if err != nil {
		return nil
	}
	names := []string{}
	for _, e := range entries {
		if e.IsDir() { // the boss subdirectory is not a level
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".txt"))
	}
	sort.Strings(names)
	return names
}

var levelHP = map[string]float64{
	// Tuned with levelGold against the actual bends and firing lanes. The
	// normal autoplay proxy holds the four floors with decreasing spare lives.
	"hub":     0.65,
	"winding": 0.77,
	"garden":  1.02,
	"canyon":  1.04,
	// The heart chamber is the endgame gate behind the unsealed heart: a
	// inward spiral. Autoplay is deterministic, so it is all-or-nothing per
	// map; 2.20 is the edge where the weak-player proxy reaches the final
	// wave and dies there — the final expedition a passive player loses,
	// while a player who arrives with hearts (+lives) and relics can hold it.
	"heart": 2.20,
}

// mazeHP keeps the procedural maze in the "chaos" tier (~60% AI win). The
// corridor-biased generator yields longer chokepoints but fewer effective
// build pockets, so it needs a lower HP than the hand-crafted maps.
const mazeHP = 0.58

// Roomier floors fund more towers; health is tuned independently against
// their actual route coverage, rather than treating every grass tile alike.
var levelGold = map[string]float64{
	"hub":     1.0,
	"winding": 1.25,
	"garden":  1.4,
	"canyon":  1.3,
	"heart":   1.4,
}

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
	m.GoldMul = levelGold[name]
	return m, nil
}

// LoadBoss loads the heart chamber: the endgame defense behind the unsealed
// heart. It is deliberately not in LevelNames, so the level select never
// lists it — the lair unseals it.
func LoadBoss() (*Map, error) {
	data, err := bossFS.ReadFile("levels/boss/heart.txt")
	if err != nil {
		return nil, err
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
	m.HPMul = levelHP["heart"]
	m.GoldMul = levelGold["heart"]
	return m, nil
}
