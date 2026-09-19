package game

import (
	"fmt"
	"strings"
)

type TowerSpec struct {
	Name    string
	Short   rune
	Cost    [3]int
	Dmg     [3]float64
	Range   [3]float64
	RoT     [3]float64
	Splash  float64
	SlowPct float64
	SlowDur float64
	Chain   [3]int
}

var TowerSpecs = [TowerCount]TowerSpec{
	{
		Name: "Gunner", Short: 'G',
		Cost: [3]int{50, 45, 70}, Dmg: [3]float64{9, 14, 22},
		Range: [3]float64{3.0, 3.4, 3.8}, RoT: [3]float64{1.3, 1.5, 1.8},
	},
	{
		Name: "Cannon", Short: 'C',
		Cost: [3]int{100, 90, 140}, Dmg: [3]float64{24, 40, 64},
		Range: [3]float64{2.8, 3.2, 3.6}, RoT: [3]float64{0.55, 0.6, 0.65},
		Splash: 1.3,
	},
	{
		Name: "Frost", Short: 'F',
		Cost: [3]int{75, 65, 100}, Dmg: [3]float64{5, 8, 12},
		Range: [3]float64{2.6, 2.9, 3.2}, RoT: [3]float64{0.9, 1.0, 1.1},
		SlowPct: 0.45, SlowDur: 1.6,
	},
	{
		Name: "Sniper", Short: 'S',
		Cost: [3]int{150, 130, 200}, Dmg: [3]float64{70, 120, 190},
		Range: [3]float64{6.0, 6.6, 7.2}, RoT: [3]float64{0.35, 0.4, 0.45},
	},
	{
		Name: "Tesla", Short: 'T',
		Cost: [3]int{200, 160, 240}, Dmg: [3]float64{30, 48, 75},
		Range: [3]float64{3.2, 3.5, 3.9}, RoT: [3]float64{0.8, 0.9, 1.0},
		Chain: [3]int{2, 3, 4},
	},
}

const ChainRange = 2.6
const ChainFalloff = 0.6

type EnemySpec struct {
	Name    string
	Short   rune
	HP      float64
	Speed   float64
	Bounty  int
	Lives   int
	SplitN  int
	SplitHP float64
}

var EnemySpecs = [EnemyCount]EnemySpec{
	{Name: "Minion", Short: 'o', HP: 16, Speed: 3.2, Bounty: 2},
	{Name: "Runner", Short: 'r', HP: 34, Speed: 2.6, Bounty: 5},
	{Name: "Grunt", Short: 'g', HP: 65, Speed: 1.6, Bounty: 8},
	{Name: "Tank", Short: 't', HP: 240, Speed: 1.05, Bounty: 20},
	{Name: "Splitter", Short: 's', HP: 110, Speed: 1.4, Bounty: 12, SplitN: 2, SplitHP: 0.5},
	{Name: "Boss", Short: 'B', HP: 1600, Speed: 0.8, Bounty: 140, Lives: 3},
}

const (
	StartingGold    = 220
	StartingLives   = 20
	MaxWaves        = 20
	SellRefund      = 0.7
	ProjectileSpeed = 14.0
	AutoWaveDelay   = 12.0
	EarlyBonusBase  = 8
)

type Difficulty int

const (
	Normal Difficulty = iota
	Easy
	Hard
)

func (d Difficulty) HPMul() float64 {
	switch d {
	case Easy:
		return 0.7
	case Hard:
		return 1.45
	}
	return 1.0
}

func (d Difficulty) Lives() int {
	if d == Hard {
		return 15
	}
	return StartingLives
}

func (d Difficulty) Gold() int {
	if d == Easy {
		return 260
	}
	if d == Hard {
		return 200
	}
	return StartingGold
}

func WaveBonus(wave int) int { return 30 + 4*wave }

func HPScale(wave int) float64 {
	w := float64(wave - 1)
	return 1 + 0.24*w + 0.012*w*w
}

func SpeedScale(wave int) float64 {
	w := float64(wave - 1)
	return 1 + 0.018*w
}

type SpawnEntry struct {
	At   float64
	Kind EnemyKind
}

func BuildWave(wave int) []SpawnEntry {
	var entries []SpawnEntry
	t := 0.0
	add := func(kind EnemyKind, count int, gap float64) {
		for i := 0; i < count; i++ {
			entries = append(entries, SpawnEntry{At: t, Kind: kind})
			t += gap
		}
		t += 1.4
	}
	add(EnemyGrunt, 4+2*wave, 0.85)
	if wave >= 2 {
		add(EnemyRunner, 2+2*wave, 0.5)
	}
	if wave >= 3 {
		add(EnemyTank, (wave+2)/4, 1.5)
	}
	if wave >= 6 {
		add(EnemySplitter, (wave-2)/2, 1.1)
	}
	if wave%5 == 0 {
		n := 1
		if wave == MaxWaves {
			n = 2
		}
		for i := 0; i < n; i++ {
			entries = append(entries, SpawnEntry{At: t, Kind: EnemyBoss})
			t += 3.5
		}
	}
	return entries
}

func WavePreview(wave int) string {
	entries := BuildWave(wave)
	counts := map[EnemyKind]int{}
	for _, e := range entries {
		counts[e.Kind]++
	}
	parts := []string{}
	order := []EnemyKind{EnemyGrunt, EnemyRunner, EnemyTank, EnemySplitter, EnemyBoss}
	for _, k := range order {
		if counts[k] > 0 {
			parts = append(parts, fmt.Sprintf("%d%c", counts[k], EnemySpecs[k].Short))
		}
	}
	return strings.Join(parts, " ")
}

// BossMod tweaks boss stats per wave: returns hp and speed multipliers.
func BossMod(wave int) (hp, speed float64) {
	switch {
	case wave >= 20:
		return 1.6, 1.0
	case wave >= 15:
		return 1.5, 1.15
	case wave >= 10:
		return 1.0, 1.5
	}
	return 1.0, 1.0
}
