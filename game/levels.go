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

// Health follows actual exposure and wave pressure, rather than road length
// alone. The tactical benchmark also checks focused, upgraded defenses.
var levelHP = map[string]float64{
	"hub":     .52,
	"canyon":  .82,
	"winding": 1.15,
	"garden":  .90,
	"heart":   1.70,
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
	"garden":  1.35,
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
	m.Encounter = map[string]Encounter{"hub": EncounterRotunda, "canyon": EncounterRift, "winding": EncounterHalls, "garden": EncounterGarden}[name]
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
	m.Encounter = EncounterHeart
	return m, nil
}
