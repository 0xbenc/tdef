package game

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
	"strings"
)

// Campaign relic bonuses are also used when describing the reward in the UI.
const (
	RelicGoldBonus  = 60
	RelicLivesBonus = 1
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
	Pull    [3]float64
	Mark    [3]float64
	MarkDur [3]float64
	MineCap [3]int
}

var TowerSpecs = [TowerCount]TowerSpec{
	{
		Name: copytext.Text("defenders.gunner.name"), Short: 'G',
		Cost: [3]int{50, 45, 70}, Dmg: [3]float64{9, 14, 22},
		Range: [3]float64{3.0, 3.4, 3.8}, RoT: [3]float64{1.3, 1.5, 1.8},
	},
	{
		Name: copytext.Text("defenders.cannonier.name"), Short: 'C',
		Cost: [3]int{100, 90, 140}, Dmg: [3]float64{24, 40, 64},
		Range: [3]float64{2.8, 3.2, 3.6}, RoT: [3]float64{0.55, 0.6, 0.65},
		Splash: 1.3,
	},
	{
		Name: copytext.Text("defenders.frost.name"), Short: 'F',
		Cost: [3]int{75, 65, 100}, Dmg: [3]float64{5, 8, 12},
		Range: [3]float64{2.6, 2.9, 3.2}, RoT: [3]float64{0.9, 1.0, 1.1},
		SlowPct: 0.45, SlowDur: 1.6,
	},
	{
		Name: copytext.Text("defenders.ranger.name"), Short: 'S',
		Cost: [3]int{150, 130, 200}, Dmg: [3]float64{70, 120, 190},
		Range: [3]float64{6.0, 6.6, 7.2}, RoT: [3]float64{0.35, 0.4, 0.45},
	},
	{
		Name: copytext.Text("defenders.lightning.name"), Short: 'T',
		Cost: [3]int{200, 160, 240}, Dmg: [3]float64{30, 48, 75},
		Range: [3]float64{3.2, 3.5, 3.9}, RoT: [3]float64{0.8, 0.9, 1.0},
		Chain: [3]int{2, 3, 4},
	},
	{
		Name: copytext.Text("defenders.trebuchet.name"), Short: 'M',
		Cost: [3]int{250, 220, 320}, Dmg: [3]float64{70, 105, 160},
		Range: [3]float64{5.0, 5.5, 6.0}, RoT: [3]float64{0.45, 0.5, 0.55},
		Splash: 2.2,
	},
	{
		Name: copytext.Text("defenders.gnolls.name"), Short: 'L',
		Cost: [3]int{80, 70, 110}, Dmg: [3]float64{6, 9, 14},
		Range: [3]float64{2.2, 2.5, 2.8}, RoT: [3]float64{3.0, 3.5, 4.0},
	},
	{
		Name: copytext.Text("defenders.runeforge.name"), Short: 'R',
		Cost: [3]int{180, 140, 220}, Dmg: [3]float64{22, 34, 50},
		Range: [3]float64{7, 8, 9}, RoT: [3]float64{1, 1, 1},
	},
	{
		Name: copytext.Text("defenders.hookmaster.name"), Short: 'H',
		Cost: [3]int{140, 100, 160}, Dmg: [3]float64{8, 12, 18},
		Range: [3]float64{3.0, 3.4, 3.8}, RoT: [3]float64{.2, .23, .26},
		Pull: [3]float64{2, 2.6, 3.2},
	},
	{
		Name: copytext.Text("defenders.sappers.name"), Short: 'P',
		Cost: [3]int{125, 100, 160}, Dmg: [3]float64{90, 145, 220},
		Range: [3]float64{2.8, 3.2, 3.6}, RoT: [3]float64{.4, .45, .5},
		Splash: 1.1, MineCap: [3]int{3, 4, 5},
	},
	{
		Name: copytext.Text("defenders.witch.name"), Short: 'W',
		Cost: [3]int{160, 120, 180}, Dmg: [3]float64{0, 0, 0},
		Range: [3]float64{3.5, 3.9, 4.3}, RoT: [3]float64{.5, .5, .5},
		Mark: [3]float64{.25, .35, .45}, MarkDur: [3]float64{3, 3.5, 4},
	},
}

const ChainRange = 2.6
const ChainFalloff = 0.6

type EnemySpec struct {
	Name       string
	Short      rune
	HP         float64
	Speed      float64
	Bounty     int
	Lives      int
	SplitN     int
	SplitHP    float64
	Plate      float64 // reduction against small shot
	ArrowPlate float64 // reduction against Ranger arrows
	Ward       float64 // reduction against frost and lightning
	Evasion    float64 // reduction against ordinary shots while not slowed
}

var EnemySpecs = [EnemyCount]EnemySpec{
	{Name: copytext.Text("enemies.squire.name"), Short: 'o', HP: 16, Speed: 3.2, Bounty: 2},
	{Name: copytext.Text("enemies.rogue.name"), Short: 'r', HP: 24, Speed: 2.6, Bounty: 5, Evasion: 0.5},
	{Name: copytext.Text("enemies.mercenary.name"), Short: 'g', HP: 65, Speed: 1.6, Bounty: 8},
	{Name: copytext.Text("enemies.paladin.name"), Short: 't', HP: 240, Speed: 1.05, Bounty: 20, Plate: 0.75, ArrowPlate: 0.4},
	{Name: copytext.Text("enemies.necromancer.name"), Short: 's', HP: 110, Speed: 1.4, Bounty: 12, SplitN: 2, SplitHP: 0.5},
	{Name: copytext.Text("enemies.player.name"), Short: 'B', HP: 1600, Speed: 0.8, Bounty: 140, Lives: 6},
	{Name: copytext.Text("enemies.wizard.name"), Short: 'w', HP: 9, Speed: 2.3, Bounty: 3},
	{Name: copytext.Text("enemies.centurion.name"), Short: 'D', HP: 180, Speed: 1.15, Bounty: 24, Ward: 0.8},
}

const (
	StartingGold    = 220
	StartingLives   = 20
	MaxWaves        = 20
	SellRefund      = 0.7
	ProjectileSpeed = 14.0
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

// WaveBonus is the gold paid for clearing a wave. Front-loaded enough that
// even a tiny wave funds progress (BTD3's $99+round shape, scaled down).
func WaveBonus(wave int) int { return 35 + 5*(wave-1) }

// EarlyBonus is the gold paid for starting the next wave before the inter-wave
// timer ends. Modest on purpose: early-starting is a tempo choice, not a
// dominant income source.
func EarlyBonus(prevWave int) int { return EarlyBonusBase + prevWave/2 }

// AutoWaveDelayFor returns the break after completing wave w (the delay before
// wave w+1 auto-starts). Tapers 11s -> ~4.5s so tempo rises as the game
// escalates, in place of a fixed 12s.
func AutoWaveDelayFor(w int) float64 {
	d := 11.0 - float64(w-1)*7.0/19.0
	d = math.Round(d*2) / 2
	if d < 4 {
		d = 4
	}
	if d > 11 {
		d = 11
	}
	return d
}

func HPScale(wave int) float64 {
	w := float64(wave - 1)
	return 1 + 0.13*w + 0.005*w*w
}

func SpeedScale(wave int) float64 {
	w := float64(wave - 1)
	return 1 + 0.018*w
}

type SpawnEntry struct {
	At   float64
	Kind EnemyKind
}

type spawnGroup struct {
	kind  EnemyKind
	count int
	gap   float64
}

type WaveDef struct {
	Theme  string
	Groups []spawnGroup
	Bosses int
}

// waves is the hand-shaped 20-wave composition (BTD3-style): one new enemy
// type introduced at a time with a legible near-solo debut, deliberate
// breather dips at introductions, and a back-loaded climax (W15/W18/W20).
// Index 1..MaxWaves; index 0 is unused.
var waves = [...]WaveDef{
	{}, // 0 unused
	{Theme: copytext.Text("encounters.waves.wave_1.theme"), Groups: []spawnGroup{{EnemyMinion, 14, 1.0}}},
	{Theme: copytext.Text("encounters.waves.wave_2.theme"), Groups: []spawnGroup{{EnemyMinion, 30, 0.8}}},
	{Theme: copytext.Text("encounters.waves.wave_3.theme"), Groups: []spawnGroup{{EnemyMinion, 15, 0.8}, {EnemyRunner, 8, 0.8}}},
	{Theme: copytext.Text("encounters.waves.wave_4.theme"), Groups: []spawnGroup{{EnemyMinion, 10, 0.7}, {EnemyRunner, 24, 0.7}}},
	{Theme: copytext.Text("encounters.waves.wave_5.theme"), Groups: []spawnGroup{{EnemyMinion, 22, 0.6}, {EnemyRunner, 28, 0.6}}},
	{Theme: copytext.Text("encounters.waves.wave_6.theme"), Groups: []spawnGroup{{EnemyMinion, 13, 0.7}, {EnemyGrunt, 15, 0.7}}},
	{Theme: copytext.Text("encounters.waves.wave_7.theme"), Groups: []spawnGroup{{EnemyRunner, 50, 0.35}}},
	{Theme: copytext.Text("encounters.waves.wave_8.theme"), Groups: []spawnGroup{{EnemyMinion, 27, 0.5}, {EnemyRunner, 30, 0.5}}},
	{Theme: copytext.Text("encounters.waves.wave_9.theme"), Groups: []spawnGroup{{EnemyRunner, 32, 0.5}, {EnemyGrunt, 15, 0.5}}},
	{Theme: copytext.Text("encounters.waves.wave_10.theme"), Groups: []spawnGroup{{EnemyGrunt, 32, 0.7}}},
	{Theme: copytext.Text("encounters.waves.wave_11.theme"), Groups: []spawnGroup{{EnemyGrunt, 12, 0.5}, {EnemyWisp, 11, 0.5}}},
	{Theme: copytext.Text("encounters.waves.wave_12.theme"), Groups: []spawnGroup{{EnemySplitter, 10, 0.8}}},
	{Theme: copytext.Text("encounters.waves.wave_13.theme"), Groups: []spawnGroup{{EnemyTank, 5, 0.8}, {EnemyRunner, 20, 0.8}}},
	{Theme: copytext.Text("encounters.waves.wave_14.theme"), Groups: []spawnGroup{{EnemyGrunt, 12, 0.6}, {EnemyShield, 3, 0.6}}},
	{Theme: copytext.Text("encounters.waves.wave_15.theme"), Groups: []spawnGroup{{EnemyMinion, 8, 1.5}}, Bosses: 1},
	{Theme: copytext.Text("encounters.waves.wave_16.theme"), Groups: []spawnGroup{{EnemyRunner, 8, 0.5}, {EnemyGrunt, 28, 0.5}, {EnemySplitter, 16, 0.5}, {EnemyShield, 2, 0.5}}},
	{Theme: copytext.Text("encounters.waves.wave_17.theme"), Groups: []spawnGroup{{EnemyRunner, 16, 0.4}, {EnemyGrunt, 38, 0.4}, {EnemySplitter, 6, 0.4}, {EnemyShield, 3, 0.4}, {EnemyTank, 3, 0.4}}},
	{Theme: copytext.Text("encounters.waves.wave_18.theme"), Groups: []spawnGroup{{EnemyGrunt, 30, 0.7}, {EnemyTank, 10, 0.7}, {EnemyShield, 8, 0.7}}, Bosses: 1},
	{Theme: copytext.Text("encounters.waves.wave_19.theme"), Groups: []spawnGroup{{EnemyRunner, 34, 0.3}, {EnemyGrunt, 46, 0.3}, {EnemySplitter, 20, 0.3}, {EnemyShield, 8, 0.3}, {EnemyTank, 6, 0.3}, {EnemyWisp, 12, 0.3}}},
	{Theme: copytext.Text("encounters.waves.wave_20.theme"), Groups: []spawnGroup{{EnemyTank, 12, 1.0}, {EnemyWisp, 24, 1.0}, {EnemySplitter, 15, 1.0}, {EnemyShield, 7, 1.0}}, Bosses: 2},
}

func BuildWave(wave int) []SpawnEntry {
	if wave < 1 || wave > MaxWaves {
		return nil
	}
	def := waves[wave]
	var entries []SpawnEntry
	t := 0.0
	add := func(kind EnemyKind, count int, gap float64) {
		for i := 0; i < count; i++ {
			entries = append(entries, SpawnEntry{At: t, Kind: kind})
			t += gap
		}
		t += 1.0 // short pause between groups
	}
	for _, g := range def.Groups {
		add(g.kind, g.count, g.gap)
	}
	for i := 0; i < def.Bosses; i++ {
		entries = append(entries, SpawnEntry{At: t, Kind: EnemyBoss})
		t += 3.5
	}
	return entries
}

// WaveTheme returns the theme label for a wave (for the HUD).
func WaveTheme(wave int) string {
	if wave < 1 || wave > MaxWaves {
		return ""
	}
	return waves[wave].Theme
}

// WaveTelegraph returns a short pre-wave message for the upcoming wave,
// shown during the break before it (BTD3's pre-round comment habit). It
// telegraphs new enemy types and the boss waves. Empty for unremarkable waves.
func WaveTelegraph(wave int) string {
	switch wave {
	case 3:
		return copytext.Text("encounters.wave_telegraph.the_guild_has_sent_rogues")
	case 6:
		return copytext.Text("encounters.wave_telegraph.mercenaries_and_they_hit_harder")
	case 7:
		return copytext.Text("encounters.wave_telegraph.a_swift_raid_loose_arrows_fly_true")
	case 11:
		return copytext.Text("encounters.wave_telegraph.wizards_fast_and_frail_but_deadly")
	case 12:
		return copytext.Text("encounters.wave_telegraph.a_necromancer_walks_among_them_they_rise")
	case 13:
		return copytext.Text("encounters.wave_telegraph.a_paladin_leads_the_vanguard")
	case 14:
		return copytext.Text("encounters.wave_telegraph.centurions_warded_they_shrug_off_your_blows")
	case 15:
		return copytext.Text("encounters.wave_telegraph.the_player_has_set_out_if_they")
	case 18:
		return copytext.Text("encounters.wave_telegraph.the_player_returns_stronger")
	case 20:
		return copytext.Text("encounters.wave_telegraph.the_final_expedition_two_players_hold_the")
	}
	return ""
}

func WavePreview(wave int) string {
	entries := BuildWave(wave)
	counts := map[EnemyKind]int{}
	for _, e := range entries {
		counts[e.Kind]++
	}
	parts := []string{}
	order := []EnemyKind{EnemyMinion, EnemyRunner, EnemyGrunt, EnemyTank, EnemyWisp, EnemySplitter, EnemyShield, EnemyBoss}
	for _, k := range order {
		if counts[k] > 0 {
			parts = append(parts, fmt.Sprintf("%d%c", counts[k], EnemySpecs[k].Short))
		}
	}
	return strings.Join(parts, " ")
}

// BossMod tweaks boss stats per wave: returns hp and speed multipliers.
// Bosses appear on wave 15 (debut), 18 (stronger), and 20 (finale, two of
// them). The finale bosses are the run's climax.
func BossMod(wave int) (hp, speed float64) {
	switch wave {
	case 15:
		return 1.0, 1.0
	case 18:
		return 1.75, 1.0
	case 20:
		return 2.5, 1.0
	}
	return 1.0, 1.0
}
