package render

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"strings"

	"github.com/0xbenc/termtd/game"
)

var EnemyJournalEntries = []TowerJournalEntry{
	{"squire", -1, copytext.Text("enemies.squire.name"), copytext.Text("enemies.squire.subtitle"), copytext.Text("enemies.squire.lore"), copytext.Text("enemies.squire.note"), 223, drawSquirePortrait},
	{"rogue", -1, copytext.Text("enemies.rogue.name"), copytext.Text("enemies.rogue.subtitle"), copytext.Text("enemies.rogue.lore"), copytext.Text("enemies.rogue.note"), 109, drawRoguePortrait},
	{"mercenary", -1, copytext.Text("enemies.mercenary.name"), copytext.Text("enemies.mercenary.subtitle"), copytext.Text("enemies.mercenary.lore"), copytext.Text("enemies.mercenary.note"), 180, drawMercenaryPortrait},
	{"paladin", -1, copytext.Text("enemies.paladin.name"), copytext.Text("enemies.paladin.subtitle"), copytext.Text("enemies.paladin.lore"), copytext.Text("enemies.paladin.note"), 223, drawPaladinPortrait},
	{"necromancer", -1, copytext.Text("enemies.necromancer.name"), copytext.Text("enemies.necromancer.subtitle"), copytext.Text("enemies.necromancer.lore"), copytext.Text("enemies.necromancer.note"), 147, drawNecromancerPortrait},
	{"player", -1, copytext.Text("enemies.player.name"), copytext.Text("enemies.player.subtitle"), copytext.Text("enemies.player.lore"), copytext.Text("enemies.player.note"), 180, drawPlayerPortrait},
	{"wizard", -1, copytext.Text("enemies.wizard.name"), copytext.Text("enemies.wizard.subtitle"), copytext.Text("enemies.wizard.lore"), copytext.Text("enemies.wizard.note"), 147, drawWizardPortrait},
	{"centurion", -1, copytext.Text("enemies.centurion.name"), copytext.Text("enemies.centurion.subtitle"), copytext.Text("enemies.centurion.lore"), copytext.Text("enemies.centurion.note"), 153, drawCenturionPortrait},
}

func EnemyJournalID(k game.EnemyKind) string {
	if k < 0 || int(k) >= len(EnemyJournalEntries) {
		return ""
	}
	return EnemyJournalEntries[k].ID
}

var JournalFloorIDs = []string{"rotunda", "rift", "halls", "garden", "heart", "depths"}
var JournalFloorNames = []string{copytext.Text("places.rotunda.short_name"), copytext.Text("places.rift.short_name"), copytext.Text("places.halls.short_name"), copytext.Text("places.garden.short_name"), copytext.Text("places.heart.short_name"), copytext.Text("places.depths.short_name")}

var PlaceJournalEntries = []TowerJournalEntry{
	{"rotunda", -1, copytext.Text("places.rotunda.name"), copytext.Text("places.rotunda.subtitle"), copytext.Text("places.rotunda.lore"), copytext.Text("places.rotunda.note"), 180, drawJournalRotunda},
	{"rift", -1, copytext.Text("places.rift.name"), copytext.Text("places.rift.subtitle"), copytext.Text("places.rift.lore"), copytext.Text("places.rift.note"), 173, drawJournalRift},
	{"halls", -1, copytext.Text("places.halls.name"), copytext.Text("places.halls.subtitle"), copytext.Text("places.halls.lore"), copytext.Text("places.halls.note"), 180, drawJournalHalls},
	{"garden", -1, copytext.Text("places.garden.name"), copytext.Text("places.garden.subtitle"), copytext.Text("places.garden.lore"), copytext.Text("places.garden.note"), 151, drawJournalGarden},
	{"heart", -1, copytext.Text("places.heart.name"), copytext.Text("places.heart.subtitle"), copytext.Text("places.heart.lore"), copytext.Text("places.heart.note"), 215, drawJournalHeart},
	{"depths", -1, copytext.Text("places.depths.name"), copytext.Text("places.depths.subtitle"), copytext.Text("places.depths.lore"), copytext.Text("places.depths.note"), 117, drawJournalDepths},
}

type journalMoment struct{ title, lore string }

var mainJournalMoments = [][3]journalMoment{
	{
		{copytext.Text("memories.rotunda.easy.title"), copytext.Text("memories.rotunda.easy.lore")},
		{copytext.Text("memories.rotunda.normal.title"), copytext.Text("memories.rotunda.normal.lore")},
		{copytext.Text("memories.rotunda.hard.title"), copytext.Text("memories.rotunda.hard.lore")},
	},
	{
		{copytext.Text("memories.rift.easy.title"), copytext.Text("memories.rift.easy.lore")},
		{copytext.Text("memories.rift.normal.title"), copytext.Text("memories.rift.normal.lore")},
		{copytext.Text("memories.rift.hard.title"), copytext.Text("memories.rift.hard.lore")},
	},
	{
		{copytext.Text("memories.halls.easy.title"), copytext.Text("memories.halls.easy.lore")},
		{copytext.Text("memories.halls.normal.title"), copytext.Text("memories.halls.normal.lore")},
		{copytext.Text("memories.halls.hard.title"), copytext.Text("memories.halls.hard.lore")},
	},
	{
		{copytext.Text("memories.garden.easy.title"), copytext.Text("memories.garden.easy.lore")},
		{copytext.Text("memories.garden.normal.title"), copytext.Text("memories.garden.normal.lore")},
		{copytext.Text("memories.garden.hard.title"), copytext.Text("memories.garden.hard.lore")},
	},
	{
		{copytext.Text("memories.heart.easy.title"), copytext.Text("memories.heart.easy.lore")},
		{copytext.Text("memories.heart.normal.title"), copytext.Text("memories.heart.normal.lore")},
		{copytext.Text("memories.heart.hard.title"), copytext.Text("memories.heart.hard.lore")},
	},
	{
		{copytext.Text("memories.depths.easy.title"), copytext.Text("memories.depths.easy.lore")},
		{copytext.Text("memories.depths.normal.title"), copytext.Text("memories.depths.normal.lore")},
		{copytext.Text("memories.depths.hard.title"), copytext.Text("memories.depths.hard.lore")},
	},
}

var StoryJournalEntries = buildStoryJournalEntries()

func buildStoryJournalEntries() []TowerJournalEntry {
	var entries []TowerJournalEntry
	for floor, moments := range mainJournalMoments {
		for diff, m := range moments {
			fi, di := floor, diff
			entries = append(entries, TowerJournalEntry{fmt.Sprintf("%s:%d", JournalFloorIDs[floor], diff), -1, m.title, JournalFloorNames[floor] + " · " + DiffName(diff), m.lore, copytext.Format("ui.journal.earned_floor", "floor", JournalFloorNames[floor], "renown", DiffName(diff)), PlaceJournalEntries[floor].Accent, func(f *Frame, x, y, w, h int) { drawJournalMoment(f, x, y, w, h, fi, di) }})
		}
	}
	bonus := []journalMoment{
		{copytext.Text("afterward.wins_1.title"), copytext.Text("afterward.wins_1.lore")},
		{copytext.Text("afterward.wins_3.title"), copytext.Text("afterward.wins_3.lore")},
		{copytext.Text("afterward.wins_5.title"), copytext.Text("afterward.wins_5.lore")},
		{copytext.Text("afterward.wins_10.title"), copytext.Text("afterward.wins_10.lore")},
	}
	for i, m := range bonus {
		bi := i
		n := []int{1, 3, 5, 10}[i]
		entries = append(entries, TowerJournalEntry{fmt.Sprintf("afterward:%d", n), -1, m.title, copytext.Text("ui.build_story_journal_entries.afterward_2"), m.lore, copytext.Format("ui.build_story_journal_entries.earned_after_victories_in_the_heart_or", "wins", fmt.Sprintf("%d", n)), 180, func(f *Frame, x, y, w, h int) { drawJournalAfterward(f, x, y, w, h, bi) }})
	}
	return entries
}

func JournalRequirement(section, cursor int) string {
	switch section {
	case JournalEnemies:
		return copytext.Text("ui.journal_requirement.encounter_this_guild_fighter_to_write_its")
	case JournalPlaces:
		return copytext.Text("ui.journal_requirement.visit_this_part_of_the_lair_to")
	case JournalStories:
		e := StoryJournalEntries[cursor]
		if strings.HasPrefix(e.ID, "afterward:") {
			return copytext.Format("ui.journal.afterward_requirement", "wins", strings.TrimPrefix(e.ID, "afterward:"))
		}
		return copytext.Format("ui.journal.floor_requirement", "floor", JournalFloorNames[cursor/3], "renown", DiffName(cursor%3))
	default:
		return copytext.Text("ui.journal_requirement.place_a_tower_to_write_its_page")
	}
}
