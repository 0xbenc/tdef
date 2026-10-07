package tui

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
)

func TestMainMenuJournalReturnsWithoutChangingProgress(t *testing.T) {
	for _, method := range []string{"Enter", "shortcut", "mouse"} {
		t.Run(method, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)
			if err := hiscore.SaveJournal(&hiscore.Journal{Towers: map[string]bool{"frost": true}}); err != nil {
				t.Fatal(err)
			}
			if err := hiscore.SaveLair(&hiscore.Lair{IntroSeen: true, Tokens: 4}); err != nil {
				t.Fatal(err)
			}
			a := newApp(nil, game.Normal)
			a.screen, a.menuSel = ScreenMenu, render.MenuJournal
			jp, _ := hiscore.JournalPath()
			lp, _ := hiscore.LairPath()
			before := map[string][]byte{}
			for _, path := range []string{jp, lp} {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				before[path] = data
			}
			switch method {
			case "Enter":
				a.handle(Event{Key: KeyEnter})
			case "shortcut":
				a.handle(Event{Rune: rune('1' + render.MenuJournal)})
			case "mouse":
				w, h := a.termSize()
				r := render.MenuRects(w, h)[render.MenuJournal]
				a.handle(Event{Mouse: true, Btn: 0, Press: true, X: r.X + 1, Y: r.Y})
			}
			if a.screen != ScreenJournal || a.journalReturn != ScreenMenu || a.g != nil || !a.journalUI.Unlocked["frost"] {
				t.Fatal("menu Journal must open saved discoveries without starting gameplay")
			}
			a.handle(Event{Rune: '3'}) // Read the discovered Frost Mage entry.
			if !a.journalUI.Reading {
				t.Fatal("saved discovery could not be read")
			}
			a.handle(Event{Key: KeyEscape})
			if a.screen != ScreenJournal || a.journalUI.Reading {
				t.Fatal("Escape from a page must return to the collection")
			}
			a.handle(Event{Key: KeyEscape})
			if a.screen != ScreenMenu || a.menuSel != render.MenuJournal || a.g != nil {
				t.Fatal("journal must return to the menu selection")
			}
			for path, data := range before {
				got, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(got, data) {
					t.Fatal("browsing the journal changed campaign progress")
				}
			}
		})
	}
}

func journalGameApp(t *testing.T) *App {
	a := owTestApp(t)
	a.fromOW = true
	a.lair.Training = &hiscore.TrainingProgress{Stage: 3, Unlocked: game.AllTowersMask}
	a.owBossSeen = map[int]bool{}
	m, err := game.LoadLevel("hub")
	if err != nil {
		t.Fatal(err)
	}
	a.enterGame(m, "hub", game.Normal)
	a.journal = hiscore.LoadJournal()
	a.g.Gold = 10000
	return a
}

func TestTowerDiscoveryIsSilentAndSurvivesRuns(t *testing.T) {
	for _, mouse := range []bool{false, true} {
		a := journalGameApp(t)
		cell, ok := firstGrass(a.g.Map)
		if !ok {
			t.Fatal("no build site")
		}
		a.ui.Cursor = cell
		a.ui.Placing = game.TowerFrost
		a.ui.PlacingOn = true
		a.ui.Message = "existing message"
		a.msgTTL = 1.5
		frame, clock, paused := a.frameNo, a.g.Time, a.ui.Paused
		if mouse {
			ox, oy, sc := a.mapBounds()
			a.handle(Event{Mouse: true, Press: true, Btn: 0, X: ox + cell.X*sc, Y: oy + cell.Y*sc})
		} else {
			a.handle(Event{Key: KeyEnter})
		}
		if len(a.g.Towers) != 1 || !hiscore.LoadJournal().Towers["frost"] {
			t.Fatal("placement failed to persist its page")
		}
		if a.screen != ScreenGame || a.frameNo != frame || a.g.Time != clock || a.ui.Paused != paused || a.ui.Message != "existing message" || a.msgTTL != 1.5 {
			t.Fatal("discovery interrupted play")
		}
		a.sellSelected() // no selection: does not affect the discovery
		a.enterGame(a.g.Map, "hub", game.Hard)
		if !a.journal.Towers["frost"] {
			t.Fatal("starting another defense erased discovery")
		}
	}
}

func TestFailedPlacementAndUpgradeCannotDiscoverAnotherTower(t *testing.T) {
	a := journalGameApp(t)
	cell, _ := firstGrass(a.g.Map)
	a.ui.Cursor = game.Vec{X: -1, Y: -1}
	a.ui.Placing = game.TowerSniper
	a.place()
	if len(a.journal.Towers) != 0 {
		t.Fatal("invalid placement unlocked a page")
	}
	a.ui.Cursor = cell
	a.g.Gold = 0
	a.place()
	if len(a.journal.Towers) != 0 {
		t.Fatal("unaffordable placement unlocked a page")
	}
	a.g.Gold = 10000
	a.ui.Placing = game.TowerGunner
	a.place()
	a.ui.Placing = game.TowerFrost
	a.place() // occupied
	if a.journal.Towers["frost"] {
		t.Fatal("occupied placement unlocked a page")
	}
	a.ui.Selected = a.g.Towers[0].ID
	a.upgradeSelected()
	a.sellSelected()
	if len(hiscore.LoadJournal().Towers) != 1 || !a.journal.Towers["gunner"] {
		t.Fatal("selling or upgrading changed the journal roster")
	}
}

func TestRelicTowerAlsoDiscoversItsPage(t *testing.T) {
	a := owTestApp(t)
	a.fromOW = true
	a.bonusTower = true
	m, _ := game.LoadLevel("hub")
	a.enterGame(m, "hub", game.Normal)
	if a.g.Gold != game.NewStateDiff(m, game.Normal).Gold {
		t.Fatal("free relic gunner spent starting gold")
	}
	if len(a.g.Towers) != 1 || !hiscore.LoadJournal().Towers["gunner"] {
		t.Fatal("free placement did not discover gunner")
	}
}

func TestJournalKeepsGameAndLairStateOnReturn(t *testing.T) {
	for _, from := range []Screen{ScreenGame, ScreenOverworld} {
		a := journalGameApp(t)
		a.screen = from
		a.ui.Paused = true
		a.ow.Msg = "still waiting"
		a.ow.BonusGold = 123
		for _, k := range []game.TowerKind{game.TowerGunner, game.TowerFrost} {
			a.discoverTower(k)
		}
		ow, ui, clock, frame := a.ow, a.ui, a.g.Time, a.frameNo
		a.handle(Event{Rune: 'j'})
		if a.screen != ScreenJournal {
			t.Fatal("journal did not open")
		}
		a.handle(Event{Key: KeyEnter})
		a.handle(Event{Key: KeyRight})
		if a.journalUI.Cursor != 2 {
			t.Fatal("page navigation did not skip locked cannonier")
		}
		for i := 0; i < 30; i++ {
			a.handle(Event{Key: KeyDown})
			a.handle(Event{Rune: 'n'})
		}
		if a.g.Time != clock || a.frameNo != frame || !reflect.DeepEqual(ow, a.ow) || !reflect.DeepEqual(ui, a.ui) {
			t.Fatal("journal advanced or changed its caller")
		}
		a.handle(Event{Key: KeyEscape})
		a.handle(Event{Key: KeyEscape})
		if a.screen != from || !reflect.DeepEqual(ui, a.ui) || !reflect.DeepEqual(ow, a.ow) {
			t.Fatal("return changed caller state")
		}
	}
	a := journalGameApp(t)
	a.ui.Paused = false
	a.handle(Event{Rune: 'j'})
	if a.screen != ScreenGame || a.ui.Paused {
		t.Fatal("journal interrupted an active defense")
	}
}

func TestJournalLockedPagesAndMouseNavigation(t *testing.T) {
	a := journalGameApp(t)
	a.ui.Paused = true
	a.openJournal()
	a.handle(Event{Key: KeyEnter})
	if a.journalUI.Reading {
		t.Fatal("locked page opened")
	}
	a.discoverTower(game.TowerSniper)
	r := render.JournalCards(80, 24)[3]
	a.handle(Event{Mouse: true, Press: true, Btn: 0, X: r.X + 1, Y: r.Y})
	if !a.journalUI.Reading || a.journalUI.Cursor != 3 {
		t.Fatal("mouse did not open discovered card")
	}
	for i := 0; i < 100; i++ {
		a.handle(Event{Mouse: true, Press: true, Btn: 65})
	}
	if a.journalUI.Scroll != render.JournalMaxScroll(80, 24, 3) {
		t.Fatal("scroll escaped page")
	}
	a.handle(Event{Key: KeyCtrlC})
	if !a.quitting {
		t.Fatal("journal ignored quit")
	}
}

func TestJournalCollectionUsesGridRows(t *testing.T) {
	a := journalGameApp(t)
	a.openJournal()
	a.handle(Event{Key: KeyDown})
	if a.journalUI.Cursor != 4 {
		t.Fatal("down did not move to the second row")
	}
	a.handle(Event{Key: KeyUp})
	if a.journalUI.Cursor != 0 {
		t.Fatal("up did not return to the first row")
	}
	a.handle(Event{Key: KeyRight})
	a.handle(Event{Key: KeyDown})
	if a.journalUI.Cursor != 5 {
		t.Fatal("grid navigation lost its column")
	}
	a.handle(Event{Key: KeyDown})
	if a.journalUI.Cursor != 9 {
		t.Fatal("down did not reach the specialists row")
	}
	a.handle(Event{Key: KeyDown})
	if a.journalUI.Cursor != 9 {
		t.Fatal("down escaped the last row")
	}
}

func TestEnemyEncounterSurvivesSpawningTickAndDefeat(t *testing.T) {
	a := journalGameApp(t)
	sink, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sink.Close() })
	a.term = &Terminal{out: sink}
	a.ui.Paused = false
	a.g.Wave = 1
	a.g.WaveActive = true
	a.g.SpawnQueue = []game.SpawnEntry{{At: 0, Kind: game.EnemyWisp}}
	// The road ends immediately: the enemy spawns and leaks in one simulation
	// step, leaving no live enemy for a renderer or post-step scan to discover.
	a.g.Map.TotalLen = .001
	a.g.Lives = 1
	a.owBossSeen = map[int]bool{}
	a.stepGame(1 / tickRate)
	if len(a.g.Enemies) != 0 || a.g.Status != game.StatusDefeat {
		t.Fatal("fixture did not disappear on the spawning tick")
	}
	if !hiscore.LoadJournal().Enemies["wizard"] {
		t.Fatal("lost first sight before cleanup")
	}
	if len(a.journal.Stories) != 0 || a.journal.EndgameWins != 0 {
		t.Fatal("defeat unlocked a victory memory")
	}
	if a.screen != ScreenGame || a.ui.Paused {
		t.Fatal("enemy discovery interrupted play")
	}
}

func TestStoryVictoryCountsOnlyOnceAndDefeatsNeverCount(t *testing.T) {
	a := journalGameApp(t)
	a.level = "heart"
	a.g.Status = game.StatusVictory
	a.g.Wave = game.MaxWaves
	a.owBossSeen = map[int]bool{}
	for i := 0; i < 120; i++ {
		a.stepGame(1 / frameRate)
	}
	if a.journal.EndgameWins != 1 || !a.journal.Stories["heart:1"] || !a.journal.Stories["afterward:1"] {
		t.Fatal("final frames recounted or lost victory")
	}
	a.enterGame(a.g.Map, "heart", game.Normal)
	a.g.Status = game.StatusDefeat
	a.stepGame(1 / frameRate)
	if a.journal.EndgameWins != 1 {
		t.Fatal("defeat counted as an endgame win")
	}
}

func TestJournalBackfillsProvenProgressWithoutRepeatingWins(t *testing.T) {
	a := owTestApp(t)
	a.lair.Record("rotunda", 2, game.MaxWaves, true)
	a.lair.Record("heart", 1, game.MaxWaves, true)
	a.lair.Record("garden", 0, 2, false)
	a.ensureJournal()
	if !a.journal.Stories["rotunda:2"] || !a.journal.Stories["heart:1"] || a.journal.Stories["garden:0"] {
		t.Fatal("historical victories inferred incorrectly")
	}
	if len(a.journal.Enemies) != int(game.EnemyCount) || !a.journal.Places["garden"] || a.journal.Places["depths"] {
		t.Fatal("historical encounters or visits inferred incorrectly")
	}
	if a.journal.EndgameWins != 1 {
		t.Fatal("historical wins counted incorrectly")
	}
	a.ensureJournal()
	a.journal = nil
	a.journalMigrated = false
	a.ensureJournal()
	if a.journal.EndgameWins != 1 {
		t.Fatal("migration recounted historical wins")
	}
}

func TestPlaceDiscoveryRequiresActualVisit(t *testing.T) {
	a := owTestApp(t)
	a.owRefresh()
	a.discoverCurrentPlace()
	if !a.journal.Places["rotunda"] || a.journal.Places["rift"] || a.journal.Places["heart"] {
		t.Fatal("initial visit spoiled sealed rooms")
	}
	fl, _ := render.OWFloorOf("rift")
	a.ow.Cursor = fl.Center
	a.discoverCurrentPlace()
	if a.journal.Places["rift"] {
		t.Fatal("sealed floor was discovered")
	}
	a.lair.Record("rotunda", 1, 20, true)
	a.owRefresh()
	a.discoverCurrentPlace()
	if !a.journal.Places["rift"] {
		t.Fatal("open floor visit not discovered")
	}
	m, _ := game.LoadBoss()
	a.fromOW = true
	a.enterGame(m, "heart", game.Normal)
	if !hiscore.LoadJournal().Places["heart"] {
		t.Fatal("heart entry did not discover place")
	}
}

func TestJournalChaptersAndCollectionPagination(t *testing.T) {
	a := journalGameApp(t)
	for _, e := range render.StoryJournalEntries {
		a.journal.DiscoverStory(e.ID)
	}
	a.openJournal()
	for i := 0; i < 3; i++ {
		a.handle(Event{Rune: ']'})
	}
	if a.journalUI.Section != render.JournalStories {
		t.Fatal("chapter navigation failed")
	}
	for i := 0; i < 16; i++ {
		a.handle(Event{Key: KeyRight})
	}
	if a.journalUI.Cursor != 16 {
		t.Fatal("collection cannot reach its third sheet")
	}
	a.handle(Event{Rune: '6'})
	if !a.journalUI.Reading || a.journalUI.Cursor != 21 {
		t.Fatal("last bonus memory not reachable")
	}
	a.handle(Event{Key: KeyRight})
	if a.journalUI.Cursor != 0 {
		t.Fatal("memory pages did not wrap")
	}
	a.handle(Event{Key: KeyEscape})
	for i := 0; i < 21; i++ {
		a.handle(Event{Key: KeyRight})
	}
	r := render.JournalCards(80, 24, render.JournalStories)[7] // last sheet has six entries
	a.handle(Event{Mouse: true, Press: true, Btn: 0, X: r.X + 1, Y: r.Y})
	if a.journalUI.Cursor != 21 || a.journalUI.Reading {
		t.Fatal("empty card corrupted selection")
	}
	a.handle(Event{Rune: ']'})
	if a.journalUI.Section != render.JournalTowers || a.journalUI.Cursor != 0 {
		t.Fatal("chapter did not reset to its collection")
	}
}

func TestQueuedEnemiesDoNotUnlockBeforeArrival(t *testing.T) {
	a := journalGameApp(t)
	a.ui.Paused = false
	a.g.Wave = 1
	a.g.WaveActive = true
	a.g.SpawnQueue = []game.SpawnEntry{{At: 0, Kind: game.EnemyMinion}, {At: 99, Kind: game.EnemyShield}}
	a.ui.Message = "existing message"
	a.msgTTL = 5
	a.stepGame(1 / tickRate)
	if !a.journal.Enemies["squire"] || a.journal.Enemies["centurion"] {
		t.Fatal("queued enemy revealed before arrival")
	}
	if a.screen != ScreenGame || a.ui.Paused || a.ui.Message != "existing message" {
		t.Fatal("first encounter interrupted play")
	}
}
