package tui

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/render"
)

// RunDragonMockup opens the portrait directly, returning to the lair on esc.
func RunDragonMockup() error {
	return runPortraitMockup(ScreenDragonMockup)
}

// RunGrakMockup opens the companion character study directly.
func RunGrakMockup() error {
	return runPortraitMockup(ScreenGrakMockup)
}

// RunGnollMockup opens the definitive Gnoll Slingers lore study.
func RunGnollMockup() error {
	return runPortraitMockup(ScreenGnollMockup)
}

// RunGunnerMockup opens the definitive Orc Gunner lore study.
func RunGunnerMockup() error {
	return runPortraitMockup(ScreenGunnerMockup)
}

// RunFrostMockup opens the Frost Mage lore portrait.
func RunFrostMockup() error {
	return runPortraitMockup(ScreenFrostMockup)
}

// RunPlayerMockup opens the guild champion lore portrait.
func RunPlayerMockup() error {
	return runPortraitMockup(ScreenPlayerMockup)
}

// RunCannonierMockup opens the monster artillery crew lore portrait.
func RunCannonierMockup() error {
	return runPortraitMockup(ScreenCannonierMockup)
}

// RunRangerMockup opens the monster archer lore portrait.
func RunRangerMockup() error {
	return runPortraitMockup(ScreenRangerMockup)
}

// RunLightningMockup opens the Lightning Mage lore portrait.
func RunLightningMockup() error {
	return runPortraitMockup(ScreenLightningMockup)
}

// RunTrebuchetMockup opens the monster siege crew lore portrait.
func RunTrebuchetMockup() error {
	return runPortraitMockup(ScreenTrebuchetMockup)
}

// RunNecromancerMockup opens the guild summoner lore portrait.
func RunNecromancerMockup() error {
	return runPortraitMockup(ScreenNecromancerMockup)
}

// RunPaladinMockup opens the guild vanguard lore portrait.
func RunPaladinMockup() error {
	return runPortraitMockup(ScreenPaladinMockup)
}

// RunCenturionMockup opens the warded guild veteran lore portrait.
func RunCenturionMockup() error {
	return runPortraitMockup(ScreenCenturionMockup)
}

// RunSquireMockup opens the young guild recruit lore portrait.
func RunSquireMockup() error {
	return runPortraitMockup(ScreenSquireMockup)
}

// RunWizardMockup opens the guild spellcaster lore portrait.
func RunWizardMockup() error {
	return runPortraitMockup(ScreenWizardMockup)
}

// RunMercenaryMockup opens the weary guild hireling lore portrait.
func RunMercenaryMockup() error {
	return runPortraitMockup(ScreenMercenaryMockup)
}

// RunRogueMockup opens the guild knife-runner lore portrait.
func RunRogueMockup() error {
	return runPortraitMockup(ScreenRogueMockup)
}

func runPortraitMockup(screen Screen) error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, game.Normal)
	a.ow = render.NewOWState()
	a.owRefresh()
	a.owBootArmed = true
	a.screen = screen
	return a.run()
}

func (a *App) handlePortraitMockup(e Event) {
	if e.Key == KeyCtrlC || (!e.Mouse && (e.Rune == 'q' || e.Rune == 'Q')) {
		a.quit()
	} else if !e.Mouse && e.Key == KeyTab {
		switch a.screen {
		case ScreenDragonMockup:
			a.toScreen(ScreenGrakMockup)
		case ScreenGrakMockup:
			a.toScreen(ScreenGnollMockup)
		case ScreenGnollMockup:
			a.toScreen(ScreenGunnerMockup)
		case ScreenGunnerMockup:
			a.toScreen(ScreenFrostMockup)
		case ScreenFrostMockup:
			a.toScreen(ScreenPlayerMockup)
		case ScreenPlayerMockup:
			a.toScreen(ScreenCannonierMockup)
		case ScreenCannonierMockup:
			a.toScreen(ScreenRangerMockup)
		case ScreenRangerMockup:
			a.toScreen(ScreenLightningMockup)
		case ScreenLightningMockup:
			a.toScreen(ScreenTrebuchetMockup)
		case ScreenTrebuchetMockup:
			a.toScreen(ScreenNecromancerMockup)
		case ScreenNecromancerMockup:
			a.toScreen(ScreenPaladinMockup)
		case ScreenPaladinMockup:
			a.toScreen(ScreenRogueMockup)
		case ScreenRogueMockup:
			a.toScreen(ScreenMercenaryMockup)
		case ScreenMercenaryMockup:
			a.toScreen(ScreenWizardMockup)
		case ScreenWizardMockup:
			a.toScreen(ScreenCenturionMockup)
		case ScreenCenturionMockup:
			a.toScreen(ScreenSquireMockup)
		default:
			a.toScreen(ScreenDragonMockup)
		}
	} else if !e.Mouse && (e.Key == KeyEscape || e.Key == KeyEnter || e.Rune == 'v' || e.Rune == 'V') {
		// Preserve the exact lair state; viewing the study spends no time,
		// refreshes no progression and does not replay the arrival sequence.
		a.screen = ScreenOverworld
		a.prev = nil
	}
}
