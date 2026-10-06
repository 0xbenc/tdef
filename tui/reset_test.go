package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
)

func resetTestApp(t *testing.T) (*App, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if err := hiscore.Save(hiscore.Table{"hub": 42}); err != nil {
		t.Fatal(err)
	}
	if err := hiscore.SaveJournal(&hiscore.Journal{EndgameWins: 3}); err != nil {
		t.Fatal(err)
	}
	if err := hiscore.SaveLair(&hiscore.Lair{IntroSeen: true, Tokens: 7}); err != nil {
		t.Fatal(err)
	}
	return &App{
		screen: ScreenMenu, menuSel: render.MenuResetProgress,
		scores: hiscore.Load(), journal: hiscore.LoadJournal(), lair: hiscore.LoadLair(),
		journalMigrated: true, owBootArmed: true, owBossSeen: map[int]bool{1: true},
		bonusGold: 100, bonusLives: 2, bonusTower: true, heartEndingPending: true,
	}, home
}

func TestResetProgressNeedsSecondConfirmation(t *testing.T) {
	for name, cancel := range map[string]Event{"default Enter": {Key: KeyEnter}, "Escape": {Key: KeyEscape}, "No": {Rune: 'n'}} {
		t.Run(name, func(t *testing.T) {
			a, _ := resetTestApp(t)
			a.handle(Event{Key: KeyEnter})
			if a.screen != ScreenResetProgress || a.resetSel != 0 {
				t.Fatal("menu option must open confirmation with Cancel selected")
			}
			if hiscore.Load()["hub"] != 42 {
				t.Fatal("opening confirmation deleted progress")
			}
			a.handle(cancel)
			if a.screen != ScreenMenu || a.resetDone || hiscore.Load()["hub"] != 42 || a.lair.Tokens != 7 {
				t.Fatal("cancelling must preserve saved and cached progress")
			}
		})
	}
}

func TestConfirmedResetClearsSessionAndRestartsOpening(t *testing.T) {
	a, home := resetTestApp(t)
	if err := os.WriteFile(filepath.Join(home, ".tdef-hiscores.json"), []byte(`{"hub":99}`), 0644); err != nil {
		t.Fatal(err)
	}
	a.handle(Event{Key: KeyEnter})
	a.handle(Event{Key: KeyDown})
	a.handle(Event{Key: KeyEnter})
	if !a.resetDone || a.resetErr != "" || a.screen != ScreenResetProgress {
		t.Fatal("reset must report success")
	}
	if len(a.scores) != 0 || a.lair.IntroSeen || a.lair.Tokens != 0 || a.journal.EndgameWins != 0 ||
		a.bonusGold != 0 || a.bonusLives != 0 || a.bonusTower || a.heartEndingPending || a.owBootArmed || len(a.owBossSeen) != 0 {
		t.Fatal("cached campaign state survived reset")
	}
	if len(hiscore.Load()) != 0 || hiscore.LoadLair().IntroSeen || hiscore.LoadJournal().EndgameWins != 0 {
		t.Fatal("reset did not persist")
	}
	a.handle(Event{Key: KeyEnter})
	if a.screen != ScreenMenu {
		t.Fatal("success must return to menu")
	}
	a.activateMenu(0)
	if a.screen != ScreenCutscene || !a.filmRemember {
		t.Fatal("Start after reset must replay the opening as a new journey")
	}
}

func TestResetProgressMouseRequiresConfirmClick(t *testing.T) {
	a, _ := resetTestApp(t)
	a.handle(Event{Rune: rune('1' + render.MenuResetProgress)})
	w, h := a.termSize()
	r := render.ResetProgressRects(w, h)[1]
	for _, e := range []Event{
		{Mouse: true, Btn: 0, Press: false, X: r.X + 1, Y: r.Y},
		{Mouse: true, Btn: 65, Press: true, X: r.X + 1, Y: r.Y},
	} {
		a.handle(e)
	}
	if hiscore.Load()["hub"] != 42 {
		t.Fatal("release or scroll confirmed a reset")
	}
	a.handle(Event{Mouse: true, Btn: 0, Press: true, X: r.X + 1, Y: r.Y})
	if !a.resetDone || len(hiscore.Load()) != 0 {
		t.Fatal("confirmation click did not reset progress")
	}
}

func TestResetProgressFailureIsVisibleAndRetryable(t *testing.T) {
	a, home := resetTestApp(t)
	path := filepath.Join(home, ".tdef-lair.json")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}
	a.activateMenu(render.MenuResetProgress)
	a.handle(Event{Key: KeyDown})
	a.handle(Event{Key: KeyEnter})
	if a.resetDone || a.resetErr == "" || a.resetSel != 0 || hiscore.Load()["hub"] != 42 {
		t.Fatal("failed reset must show an error, preserve saves, and default back to Cancel")
	}
	if a.bonusGold != 100 || a.bonusLives != 2 || !a.bonusTower {
		t.Fatal("failed reset must preserve pending relic bonuses")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	a.handle(Event{Key: KeyDown})
	a.handle(Event{Key: KeyEnter})
	if !a.resetDone || a.resetErr != "" {
		t.Fatal("retry failed")
	}
}
