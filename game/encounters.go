package game

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
	"sort"
	"strings"
)

// Encounter is attached only by authored level loaders. A zero encounter
// keeps the original wave schedule, including every procedural maze.
type Encounter string

const (
	EncounterRotunda Encounter = "rotunda"
	EncounterRift    Encounter = "rift"
	EncounterHalls   Encounter = "halls"
	EncounterGarden  Encounter = "garden"
	EncounterHeart   Encounter = "heart"
)

func TacticalBrief(m *Map) string {
	if m == nil {
		return ""
	}
	switch m.Encounter {
	case EncounterRotunda:
		return copytext.Text("encounters.tactical_brief.inner_courts_can_cover_the_road_more")
	case EncounterRift:
		return copytext.Text("encounters.tactical_brief.separate_banks_need_separate_defenses")
	case EncounterHalls:
		return copytext.Text("encounters.tactical_brief.long_bows_reach_across_the_returning_halls")
	case EncounterGarden:
		return copytext.Text("encounters.tactical_brief.slow_the_front_ranks_catch_the_crowd")
	case EncounterHeart:
		return copytext.Text("encounters.tactical_brief.thin_the_outer_ranks_finish_the_survivors")
	}
	return ""
}

// waveGroups copies the base introductions before shaping a floor's columns.
// New kinds still debut on the same waves; changes are legible through the
// live preview, and never mutate the shared/procedural wave definitions.
func waveGroups(m *Map, wave int) []spawnGroup {
	groups := append([]spawnGroup(nil), waves[wave].Groups...)
	if m == nil || wave < 3 {
		return groups
	}
	for i := range groups {
		g := &groups[i]
		switch m.Encounter {
		case EncounterRift:
			if g.kind == EnemyRunner {
				g.gap *= .8
			} else {
				g.gap *= .95
			}
		case EncounterHalls:
			gap := .95
			switch g.kind {
			case EnemyRunner:
				gap = 1.25
			case EnemyGrunt:
				gap = 1.45
			case EnemyTank, EnemyShield, EnemySplitter:
				gap = 1.85
			}
			g.gap = math.Max(g.gap, gap)
			if g.count > 24 {
				g.count = int(math.Ceil(float64(g.count) * .75))
			}
		case EncounterGarden:
			g.gap *= .55
			if wave >= 5 && g.count >= 12 {
				g.count = int(math.Ceil(float64(g.count) * 1.1))
			}
		case EncounterHeart:
			if g.kind == EnemyRunner || g.kind == EnemyMinion || g.kind == EnemyWisp {
				g.gap *= .65
			} else {
				g.gap *= .85
			}
		}
	}
	if m.Encounter == EncounterHeart && wave >= 16 {
		rank := func(k EnemyKind) int {
			switch k {
			case EnemyRunner, EnemyWisp, EnemyMinion:
				return 0
			case EnemyGrunt:
				return 1
			case EnemySplitter:
				return 2
			default:
				return 3
			}
		}
		sort.SliceStable(groups, func(i, j int) bool { return rank(groups[i].kind) < rank(groups[j].kind) })
	}
	return groups
}

func BuildWaveFor(m *Map, wave int) []SpawnEntry {
	if wave < 1 || wave > MaxWaves {
		return nil
	}
	if m == nil || m.Encounter == "" || m.Encounter == EncounterRotunda {
		return BuildWave(wave)
	}
	groups := waveGroups(m, wave)
	entries := []SpawnEntry{}
	at := 0.
	for _, g := range groups {
		for i := 0; i < g.count; i++ {
			entries = append(entries, SpawnEntry{At: at, Kind: g.kind})
			at += g.gap
			if m.Encounter == EncounterRift && i+1 < g.count && (i+1)%8 == 0 {
				at += 1.4
			}
		}
		pause := 1.
		if m.Encounter == EncounterGarden {
			pause = .7
		}
		if m.Encounter == EncounterHeart {
			pause = 2.0
		}
		at += pause
	}
	for i := 0; i < waves[wave].Bosses; i++ {
		entries = append(entries, SpawnEntry{At: at, Kind: EnemyBoss})
		gap := 3.5
		if m.Encounter == EncounterHeart {
			gap = 6
		}
		at += gap
	}
	return entries
}

func WavePreviewFor(m *Map, wave int) string {
	counts := [EnemyCount]int{}
	for _, e := range BuildWaveFor(m, wave) {
		counts[e.Kind]++
	}
	parts := []string{}
	for _, k := range []EnemyKind{EnemyMinion, EnemyRunner, EnemyGrunt, EnemyTank, EnemyWisp, EnemySplitter, EnemyShield, EnemyBoss} {
		if counts[k] > 0 {
			parts = append(parts, fmt.Sprintf("%d%c", counts[k], EnemySpecs[k].Short))
		}
	}
	return strings.Join(parts, " ")
}

func WaveTelegraphFor(m *Map, wave int) string {
	if m != nil && wave == 7 {
		switch m.Encounter {
		case EncounterRift:
			return copytext.Text("encounters.wave_telegraph_for.rogues_in_bursts_cover_the_crossings_and")
		case EncounterHalls:
			return copytext.Text("encounters.wave_telegraph_for.a_spaced_rogue_raid_long_bows_can")
		case EncounterGarden:
			return copytext.Text("encounters.wave_telegraph_for.a_packed_rogue_raid_slow_the_front")
		case EncounterHeart:
			return copytext.Text("encounters.wave_telegraph_for.a_rogue_rush_let_the_outer_defense")
		}
	}
	if text := WaveTelegraph(wave); text != "" {
		return text
	}
	if m == nil {
		return ""
	}
	if wave == 5 || wave == 10 || wave == 16 {
		switch m.Encounter {
		case EncounterRift:
			return copytext.Text("encounters.wave_telegraph_for.raiding_parties_cross_in_bursts_hold_more")
		case EncounterHalls:
			return copytext.Text("encounters.wave_telegraph_for.a_spaced_column_long_bows_can_cover")
		case EncounterGarden:
			return copytext.Text("encounters.wave_telegraph_for.they_arrive_shoulder_to_shoulder_catch_them")
		case EncounterHeart:
			return copytext.Text("encounters.wave_telegraph_for.a_rushing_front_then_heavy_armor_keep")
		}
	}
	return ""
}
