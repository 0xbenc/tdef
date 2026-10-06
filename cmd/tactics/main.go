package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/0xbenc/termtd/game"
)

type result struct {
	Map, Plan                                string
	Won                                      bool
	Wave, Lives, Leaks, Towers, Kills, Spent int
	Time                                     float64
	Mix                                      [game.TowerCount]int
	FirstAttack                              float64
	TimedOut                                 bool
}
type plan struct {
	name    string
	kinds   []game.TowerKind
	cluster bool
}

var plans = []plan{
	{"greedy", nil, false},
	{"gunline", []game.TowerKind{game.TowerGunner}, false},
	{"rangers", []game.TowerKind{game.TowerGunner, game.TowerSniper, game.TowerSniper}, false},
	{"frost-rangers", []game.TowerKind{game.TowerGunner, game.TowerSniper, game.TowerFrost, game.TowerSniper}, false},
	{"cannons", []game.TowerKind{game.TowerGunner, game.TowerCannon, game.TowerCannon}, false},
	{"frost-chain", []game.TowerKind{game.TowerGunner, game.TowerTesla, game.TowerFrost, game.TowerTesla}, false},
	{"siege", []game.TowerKind{game.TowerGunner, game.TowerMortar, game.TowerMortar}, false},
	{"frost-siege", []game.TowerKind{game.TowerGunner, game.TowerMortar, game.TowerFrost, game.TowerMortar}, false},
	{"slingers", []game.TowerKind{game.TowerGunner, game.TowerFlak, game.TowerFlak}, false},
	{"ranged-siege", []game.TowerKind{game.TowerGunner, game.TowerSniper, game.TowerMortar}, false},
	{"one-camp", []game.TowerKind{game.TowerGunner, game.TowerCannon, game.TowerFrost, game.TowerSniper, game.TowerTesla}, true},
}

func run(m *game.Map, name string, p plan, diff game.Difficulty, hearts int) result {
	s := game.NewStateDiff(m, diff)
	s.Lives += hearts
	firstAttack := 0.0
	var ai *game.Autoplay
	if p.kinds == nil {
		ai = game.NewAutoplay(s)
	}
	// Weighted coverage shares spending across the track. A deliberately
	// fixed small camp is the same builder with a spatial constraint.
	samples := m.Samples
	cover := make([][][]int, game.TowerCount)
	for k := game.TowerKind(0); k < game.TowerCount; k++ {
		cover[k] = make([][]int, m.W*m.H)
		for y := 1; y < m.H-1; y++ {
			for x := 1; x < m.W-1; x++ {
				v := game.Vec{X: x, Y: y}
				if m.At(v) != game.CellGrass {
					continue
				}
				for i, sp := range samples {
					if v.Center().Dist(sp) <= game.TowerSpecs[k].Range[0] {
						cover[k][y*m.W+x] = append(cover[k][y*m.W+x], i)
					}
				}
			}
		}
	}
	anchor := game.Vec{X: -99, Y: -99}
	bestCover := 0
	for y := 1; y < m.H-1; y++ {
		for x := 1; x < m.W-1; x++ {
			n := len(cover[game.TowerGunner][y*m.W+x])
			if n > bestCover {
				anchor = game.Vec{X: x, Y: y}
				bestCover = n
			}
		}
	}
	next, seq := 0., 0
	for s.Status == game.StatusRunning && s.Time < 4000 {
		if ai != nil {
			ai.Tick()
		} else if s.Time >= next {
			next = s.Time + .5
			// Start with affordable defenders, then follow the selected composition.
			k := game.TowerGunner
			if seq >= 2 {
				k = p.kinds[(seq-2)%len(p.kinds)]
			}
			cost := game.TowerSpecs[k].Cost[0]
			weights := make([]float64, len(samples))
			for i := range weights {
				weights[i] = 1
			}
			for _, tw := range s.Towers {
				for i, sp := range samples {
					if tw.Cell.Center().Dist(sp) <= tw.Range() {
						if tw.Kind != game.TowerFrost {
							weights[i] /= 1.4
						}
					}
				}
			}
			best, score := game.Vec{X: -1, Y: -1}, 0.
			for y := 1; y < m.H-1; y++ {
				for x := 1; x < m.W-1; x++ {
					v := game.Vec{X: x, Y: y}
					if m.At(v) != game.CellGrass || s.TowerAt(v) != nil {
						continue
					}
					if p.cluster && v.Center().Dist(anchor.Center()) > 3.2 {
						continue
					}
					n := 0.
					for _, i := range cover[k][y*m.W+x] {
						// Open near the entrance, as a player would, rather than
						// waiting for the first expedition to reach a distant loop.
						if seq < 2 && !p.cluster && i >= 36 {
							continue
						}
						weight := weights[i]
						if k == game.TowerFrost {
							weight = 1
							for _, tw := range s.Towers {
								if tw.Kind != game.TowerFrost && tw.Cell.Center().Dist(samples[i]) <= tw.Range() {
									weight += .7
								}
							}
						} else {
							for _, tw := range s.Towers {
								if tw.Kind == game.TowerFrost && tw.Cell.Center().Dist(samples[i]) <= tw.Range() {
									weight *= 1.6
								}
							}
						}
						n += weight
					}
					if n > score {
						score = n
						best = v
					}
				}
			}
			// Mature a four-tower foundation before continually buying new towers.
			var upgrade *game.Tower
			upScore := 0.
			if seq >= 2 && s.Wave >= 1 {
				for _, tw := range s.Towers {
					if tw.Level >= 3 || (s.Wave < 3 && tw.Level >= 2) {
						continue
					}
					v := float64(len(cover[tw.Kind][tw.Cell.Y*m.W+tw.Cell.X])) / float64(s.UpgradeCost(tw))
					if v > upScore {
						upgrade = tw
						upScore = v
					}
				}
			}
			if upgrade != nil && s.Gold >= s.UpgradeCost(upgrade)+12 {
				s.Upgrade(upgrade)
			} else if (upgrade == nil || s.Wave < 3) && len(s.Towers) < 24 && best.X >= 0 && s.Gold >= cost+12 {
				if s.Build(best, k) != nil {
					seq++
				}
			}

		}
		s.Step(.05)
		if firstAttack == 0 && (len(s.Projectiles) > 0 || len(s.Beams) > 0 || s.TotalKills > 0) {
			firstAttack = math.Round(s.Time*10) / 10
		}

	}
	r := result{Map: name, Plan: p.name, Won: s.Status == game.StatusVictory, Wave: s.Wave, Lives: s.Lives, Leaks: s.TotalLeaks, Towers: len(s.Towers), Kills: s.TotalKills, Time: math.Round(s.Time)}
	r.FirstAttack = firstAttack
	r.TimedOut = s.Status == game.StatusRunning
	for _, tw := range s.Towers {
		r.Mix[tw.Kind]++
		r.Spent += tw.Invested
	}
	return r
}

// This is a deterministic balance probe, not a claim about human win rates.
// Focused plans share foundations, upgrade rules, a 24-tower cap and reserve.
func main() {
	diffFlag := flag.String("diff", "normal", "easy | normal | hard")
	neutral := flag.Bool("neutral", false, "compare shapes with identical waves, health and economy")
	campaign := flag.Bool("campaign", false, "include the hearts earned before each floor")
	floor := flag.String("floor", "all", "hub | canyon | winding | garden | heart | all")
	only := flag.String("plan", "all", "one strategy name, or all")
	out := flag.String("out", "", "optional JSON report")
	flag.Parse()
	diff := game.Normal
	switch strings.ToLower(*diffFlag) {
	case "easy":
		diff = game.Easy
	case "normal":
	case "hard":
		diff = game.Hard
	default:
		fmt.Fprintln(os.Stderr, "invalid difficulty")
		os.Exit(2)
	}
	names := []string{"hub", "canyon", "winding", "garden", "heart"}
	if *floor != "all" {
		names = []string{*floor}
	}
	results := []result{}
	for _, name := range names {
		var m *game.Map
		var err error
		if name == "heart" {
			m, err = game.LoadBoss()
		} else {
			m, err = game.LoadLevel(name)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if *neutral {
			m.Encounter = ""
			m.HPMul = .9
			m.GoldMul = 1.25
		}
		hearts := 0
		if *campaign {
			hearts = map[string]int{"hub": 0, "canyon": 1, "winding": 2, "garden": 3, "heart": 4}[name]
		}
		for _, p := range plans {
			if *only != "all" && *only != p.name {
				continue
			}
			r := run(m, name, p, diff, hearts)
			results = append(results, r)
			fmt.Printf("%-8s %-14s win=%-5v wave=%2d lives=%2d leaks=%2d towers=%2d spent=%5d time=%4.0f\n", name, p.name, r.Won, r.Wave, r.Lives, r.Leaks, r.Towers, r.Spent, r.Time)
		}
	}
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "unknown strategy")
		os.Exit(2)
	}
	if *out != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err == nil {
			err = os.WriteFile(*out, append(data, '\n'), 0644)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
