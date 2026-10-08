package tui

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
	"strconv"
	"time"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/render"
)

// RunOverworld starts on the lair map directly (CLI entry point).
func RunOverworld() error {
	term, err := Open()
	if err != nil {
		return err
	}
	a := newApp(term, game.Normal)
	a.ow = render.NewOWState()
	a.owRefresh()
	a.screen = ScreenOverworld
	if a.maybeOpening() {
		return a.run()
	}
	if !a.owBootArmed {
		a.owBootArmed = true
		a.ow.BootTTL = render.OWBootFrames
	}
	return a.run()
}

// owFloorOrder is the lair's depth order — the wheel browser's sequence.
var depthsLockedMessage = copytext.Text("overworld.depths_locked_message.depths_sealed_beat_all_floors_and_the")

var owFloorOrder = []string{"rift", "rotunda", "halls", "garden", "depths"}

// handleOverworld drives Grak around the lair map: arrows/wasd walk (one
// cell per keypress), enter descends into the floor under the cursor, tab
// cycles the renown, 0-9 seed the Depths, t spends a relic at the Rotunda,
// esc returns to the title, q quits. The mouse clicks a floor pad to step
// onto it (or descend if already on it); the wheel hops between floors.
func (a *App) handleOverworld(e Event) {
	if e.Key == KeyCtrlC {
		a.quit()
		return
	}
	if a.ow.BootTTL > 0 {
		// The lair is waking: Grak cannot walk yet. Only nav keys pass.
		if e.Key == KeyEscape {
			a.toScreen(ScreenTitle)
		} else if !e.Mouse && (e.Rune == 'q' || e.Rune == 'Q') {
			a.quit()
		}
		return
	}
	if e.Mouse {
		a.handleOWMouse(e)
		return
	}
	if a.ow.Descending != "" {
		return // the door is closing
	}
	if a.ow.RelicMenu {
		switch {
		case e.Rune == '1':
			a.owSpendRelic(0)
		case e.Rune == '2':
			a.owSpendRelic(1)
		case e.Rune == '3':
			a.owSpendRelic(2)
		case e.Key == KeyEscape:
			a.ow.RelicMenu = false
		case e.Rune == 'q' || e.Rune == 'Q':
			a.quit()
		}
		return
	}
	switch e.Rune {
	case 'q', 'Q':
		a.quit()
	case 'w', 'W':
		a.owWalk(0, -1)
	case 's', 'S':
		a.owWalk(0, 1)
	case 'a', 'A':
		a.owWalk(-1, 0)
	case 'd', 'D':
		a.owWalk(1, 0)
	case 'r', 'R':
		a.ow.RevealAll = !a.ow.RevealAll
		if a.ow.RevealAll {
			a.ow.Msg = copytext.Text("overworld.handle_overworld.the_whole_lair_lit_every_floor_revealed")
			a.ow.MsgKind = render.OWMessageNeutral
		} else {
			a.ow.Msg = ""
			a.ow.MsgKind = render.OWMessageNeutral
		}
	case 'i', 'I':
		if a.owCursorFloor() == "rotunda" {
			a.startCutscene(render.FilmOpening, ScreenOverworld, false)
		}
	case 'e', 'E':
		if a.owCursorFloor() == "rotunda" && a.lair.BossHeld(a.ow.Diff) {
			a.startCutscene(render.FilmEnding, ScreenOverworld, false)
		}
	case 'h', 'H':
		a.openHelp()
	case 'j', 'J':
		a.openJournal()
	case 'g', 'G':
		if a.owCursorFloor() == "rotunda" {
			a.toScreen(ScreenGrakMockup)
		}
	case 'v', 'V':
		if a.owCursorFloor() == "rotunda" {
			a.toScreen(ScreenDragonMockup)
		}
	case 't', 'T':
		if a.owCursorFloor() == "rotunda" && a.ow.Tokens > 0 {
			a.ow.RelicMenu = true
		}
	case 'c', 'C':
		if a.owCursorFloor() == "depths" {
			a.ow.Seed = ""
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if a.owCursorFloor() == "depths" && len(a.ow.Seed) < 18 {
			a.ow.Seed += string(e.Rune)
		}
	}
	switch e.Key {
	case KeyUp:
		a.owWalk(0, -1)
	case KeyDown:
		a.owWalk(0, 1)
	case KeyLeft:
		a.owWalk(-1, 0)
	case KeyRight:
		a.owWalk(1, 0)
	case KeyTab:
		a.owCycleDiff(1)
	case KeyBackspace:
		if a.owCursorFloor() == "depths" && a.ow.Seed != "" {
			a.ow.Seed = a.ow.Seed[:len(a.ow.Seed)-1]
		}
	case KeyEnter:
		a.owEnter()
	case KeyEscape:
		a.toScreen(ScreenTitle)
	}
}

// owWalk is one keypress: one cell. The terminal driver emits no key-release
// events, so a held key cannot be told from a tap — repeating would make a
// single tap stride many cells, so walking is strictly tap-per-step.
func (a *App) owWalk(dx, dy int) {
	a.owStep(dx, dy)
}

// owStep is one deterministic step: walk if the destination is walkable, and
// leave a trail behind the old cell.
func (a *App) owStep(dx, dy int) {
	a.ow.Msg = ""
	a.ow.MsgKind = render.OWMessageNeutral
	n := game.Vec{X: a.ow.Cursor.X + dx, Y: a.ow.Cursor.Y + dy}
	if render.OWCanWalk(n.X, n.Y, a.ow) {
		a.ow.PushTrail(a.ow.Cursor)
		a.ow.Cursor = n
		a.discoverCurrentPlace()
	} else if fl, room := render.OWFloorAt(n.X, n.Y); room && !render.OWFloorOpen(fl.ID, a.ow) {
		a.owSealedMessage(fl.ID)
	}
}

// owTick advances the overworld's frame-counted transitions and the held
// walk. Called once per 30fps frame from the app loop.
func (a *App) owTick() {
	a.discoverCurrentPlace()
	st := &a.ow
	if !st.CameraSet {
		st.CameraX, st.CameraSet = float64(st.Cursor.X), true
	}
	// Ease the viewport toward Grak without adding extra walking steps.
	st.CameraX += (float64(st.Cursor.X) - st.CameraX) * 0.25
	if math.Abs(float64(st.Cursor.X)-st.CameraX) < 0.02 {
		st.CameraX = float64(st.Cursor.X)
	}
	if st.BootTTL > 0 {
		st.BootTTL--
	}
	if st.BlastTTL > 0 {
		st.BlastTTL--
	}
	if st.ReturnTTL > 0 {
		st.ReturnTTL--
		if st.ReturnTTL == 0 {
			st.ReturnMsg = ""
			st.ReturnKind = render.OWMessageNeutral
			st.ReturnFX = render.OWReturnFX{}
		}
	}
	for i := len(st.TrailAge) - 1; i >= 0; i-- {
		if st.TrailAge[i] > 0 {
			st.TrailAge[i]--
		}
		if st.TrailAge[i] <= 0 {
			st.Trail = append(st.Trail[:i], st.Trail[i+1:]...)
			st.TrailAge = append(st.TrailAge[:i], st.TrailAge[i+1:]...)
		}
	}
	for id, ttl := range st.Unsealing {
		ttl--
		if ttl <= 0 {
			delete(st.Unsealing, id)
			st.Unlocked[id] = true
			st.ReturnMsg = copytext.Format("overworld.progress.opened", "floor", render.OWFloorName(id))
			st.ReturnKind = render.OWMessageSuccess
			st.ReturnTTL = 90
		} else {
			st.Unsealing[id] = ttl
		}
	}
	if st.Descending != "" {
		st.DescendTTL--
		if st.DescendTTL <= 0 {
			id := st.Descending
			st.Descending = ""
			st.DescendTTL = 0
			a.owLaunch(id)
		}
	}
}

// owCursorFloor is the floor id under Grak ("" off the pads).
func (a *App) owCursorFloor() string {
	if fl, ok := render.OWFloorAt(a.ow.Cursor.X, a.ow.Cursor.Y); ok {
		return fl.ID
	}
	return ""
}

// owEnter descends into the floor under the cursor (the unsealed Rotunda is
// the heart door). Sealed floors are refused with a message. The descent is
// a short transition: the pad floods in, then owLaunch starts the game.
func (a *App) owEnter() {
	st := &a.ow
	if st.Descending != "" {
		return
	}
	fl, ok := render.OWFloorAt(st.Cursor.X, st.Cursor.Y)
	if !ok {
		if id := a.owBrowseFloor(); id != "" && !render.OWFloorOpen(id, *st) {
			a.owSealedMessage(id)
		}
		return
	}
	id := fl.ID
	if id == "rotunda" && st.BossReady {
		id = render.HeartFloorID
	}
	if id != render.HeartFloorID && !render.OWFloorOpen(id, *st) {
		a.owSealedMessage(id)
		return
	}
	st.RelicMenu = false
	st.Msg = ""
	st.MsgKind = render.OWMessageNeutral
	st.Descending = id
	st.DescendTTL = render.OWDescendFrames
}

// owLaunch descends: load the floor's map, carry the pending relic bonuses
// into the run, and enter the game. The overworld is the source of the run,
// so its result comes back to the lair.
func (a *App) owLaunch(id string) {
	st := &a.ow
	if id == hiscore.DepthsFloor && !a.lair.DepthsReady(st.Diff) {
		st.Descending, st.DescendTTL = "", 0
		st.Msg = depthsLockedMessage
		st.MsgKind = render.OWMessageLocked
		return
	}
	a.diff = render.Difficulties[st.Diff]
	a.fromOW = true
	a.owFloorID = id
	a.bonusGold = st.BonusGold
	a.bonusTower = st.BonusTower
	a.bonusLives = st.BonusLives
	st.BonusGold, st.BonusTower, st.BonusLives = 0, false, 0
	var (
		m    *game.Map
		name string
		err  error
	)
	switch id {
	case render.HeartFloorID:
		m, err = game.LoadBoss()
		name = "heart"
	case "depths":
		seed, ok := parseOWSeed(st.Seed)
		if !ok {
			st.Msg = copytext.Text("overworld.ow_launch.seed_too_large")
			st.MsgKind = render.OWMessageNeutral
			a.fromOW = false
			return
		}
		m, err = game.MazeFromSeed(seed)
		name = fmt.Sprintf("maze%d", seed)
	default:
		fl, ok := render.OWFloorOf(id)
		if !ok {
			st.Msg = copytext.Text("overworld.ow_launch.unknown_floor")
			st.MsgKind = render.OWMessageNeutral
			a.fromOW = false
			return
		}
		m, err = game.LoadLevel(fl.Level)
		name = fl.Level
	}
	if err != nil {
		st.Msg = err.Error()
		st.MsgKind = render.OWMessageNeutral
		a.fromOW = false
		return
	}
	a.enterGame(m, name, a.diff)
}

// parseOWSeed converts the typed Depths seed to an int64. An empty string (or
// a seed of zero) means "uncharted": a fresh random maze.
func parseOWSeed(s string) (int64, bool) {
	if s == "" {
		return time.Now().UnixNano(), true
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	if n == 0 {
		return time.Now().UnixNano(), true
	}
	return n, true
}

// owCycleDiff turns the renown (the difficulty the lair is remembered at):
// records, badges, hearts and unseals all follow it.
func (a *App) owCycleDiff(dir int) {
	st := &a.ow
	n := len(render.Difficulties)
	st.Diff = (st.Diff + dir + n) % n
	st.Unsealing = map[string]int{} // another renown owns different seals
	a.diff = render.Difficulties[st.Diff]
	a.owRefresh()
}

// owSpendRelic trades one relic for a bonus on the next defense.
func (a *App) owSpendRelic(kind int) {
	st := &a.ow
	if st.Tokens <= 0 {
		st.RelicMenu = false
		return
	}
	switch kind {
	case 0:
		st.BonusGold += game.RelicGoldBonus
		st.Msg = copytext.Format("overworld.ow_spend_relic.the_lair_grants_60_gold_to_the", "gold", strconv.Itoa(game.RelicGoldBonus))
		st.MsgKind = render.OWMessageNeutral
	case 1:
		if st.BonusTower {
			st.Msg = copytext.Text("overworld.ow_spend_relic.gunner_already_reserved")
			st.MsgKind = render.OWMessageNeutral
			st.RelicMenu = false
			return
		}
		st.BonusTower = true
		st.Msg = copytext.Format("overworld.ow_spend_relic.an_orc_gunner_will_be_waiting_at", "tower", game.TowerSpecs[game.TowerGunner].Name)
		st.MsgKind = render.OWMessageNeutral
	case 2:
		st.BonusLives += game.RelicLivesBonus
		st.Msg = copytext.Format("overworld.ow_spend_relic.the_heart_steadies_1_to_the_next", "lives", strconv.Itoa(game.RelicLivesBonus))
		st.MsgKind = render.OWMessageNeutral
	}
	st.Tokens--
	a.lair.Tokens = st.Tokens
	a.lair.BonusGold, a.lair.BonusTower, a.lair.BonusLives = st.BonusGold, st.BonusTower, st.BonusLives
	a.saveCampaign()
	st.RelicMenu = false
}

// owRefresh rebuilds the overworld's memory view from the lair's records:
// per-floor results, best scores, the dragon hearts, the unseal chain, the
// relics. Session fields (cursor, transitions, seed) are left alone.
func (a *App) owRefresh() {
	st := &a.ow
	if st.Unsealing == nil {
		st.Unsealing = map[string]int{}
	}
	if st.Diff < 0 || st.Diff >= len(render.Difficulties) {
		st.Diff = diffIndex(a.diff)
	}
	d := st.Diff
	if a.lair == nil || (a.saveErrors[saveCampaign] == nil && a.loadErrors[saveCampaign] == nil) {
		a.reloadCampaign()
	}
	st.BonusGold, st.BonusTower, st.BonusLives = a.lair.BonusGold, a.lair.BonusTower, a.lair.BonusLives
	recs := map[string]render.OWRec{}
	for _, id := range []string{"rift", "rotunda", "halls", "garden", "depths", render.HeartFloorID} {
		r := a.lair.Floor(id, d)
		recs[id] = render.OWRec{Cleared: r.Cleared, BestWave: r.BestWave, LastWave: r.LastWave, LastWon: r.LastWon}
	}
	st.Records = recs
	st.Scores = map[string]int{}
	scores := a.scores
	if a.saveErrors[saveScores] == nil && a.loadErrors[saveScores] == nil {
		a.reloadScores()
		scores = a.scores
	}
	for lvl, id := range map[string]string{
		"canyon": "rift", "hub": "rotunda", "winding": "halls", "garden": "garden",
		"heart": render.HeartFloorID,
	} {
		if s, ok := scores[lvl]; ok {
			st.Scores[id] = s
		}
	}
	// The Depths' best is the best across every maze seed.
	best := 0
	for k, s := range scores {
		if len(k) > 4 && k[:4] == "maze" && s > best {
			best = s
		}
	}
	if best > 0 {
		st.Scores["depths"] = best
	}
	st.Hearts = a.lair.ClearedCount(d)
	st.BossReady = a.lair.BossReady(d)
	st.BossDone = a.lair.BossHeld(d)
	st.Tokens = a.lair.Tokens
	// Mark this renown's heart state as seen, so an already-unsealed heart
	// never re-fires the shockwave on entry.
	if a.owBossSeen == nil {
		a.owBossSeen = map[int]bool{}
	}
	a.owBossSeen[d] = st.BossReady
	// The unseal chain: the Rotunda is always open; the Rift opens after
	// the Rotunda is held; the Halls open after the Rift, the Garden after
	// the Halls. Depths opens after all four and the Heart are held.
	st.Unlocked = map[string]bool{"rotunda": true}
	if recs["rotunda"].Cleared {
		st.Unlocked["rift"] = true
	}
	if recs["rift"].Cleared {
		st.Unlocked["halls"] = true
	}
	if recs["halls"].Cleared {
		st.Unlocked["garden"] = true
	}
	if a.lair.DepthsReady(d) {
		st.Unlocked["depths"] = true
	}
	if fl, room := render.OWFloorAt(st.Cursor.X, st.Cursor.Y); room && !render.OWFloorOpen(fl.ID, *st) {
		a.owVisitFloor(fl.ID)
	}
}

// handleOWMouse: the wheel hops between floors in the lair's depth order; a
// click visits an open floor or a sealed floor's doorstep; a second click
// on an occupied open pad descends.
func (a *App) handleOWMouse(e Event) {
	if !e.Press {
		return
	}
	switch e.Btn {
	case 64:
		a.owWheel(true)
		return
	case 65:
		a.owWheel(false)
		return
	}
	if a.ow.Descending != "" {
		return
	}
	w, h := a.termSize()
	if c, ok := render.OWFloorAtFrame(w, h, e.X, e.Y, a.ow); ok {
		if hit, fok := render.OWFloorAt(c.X, c.Y); fok {
			if on, ok := render.OWFloorAt(a.ow.Cursor.X, a.ow.Cursor.Y); ok && on.ID == hit.ID {
				a.owEnter()
				return
			}
		}
		if fl, ok := render.OWFloorAt(c.X, c.Y); ok {
			a.owVisitFloor(fl.ID)
		}
	}
}

// owWheel visits floors in depth order, stopping at sealed doorways.
func (a *App) owWheel(up bool) {
	st := &a.ow
	if st.Descending != "" {
		return
	}
	dir := 1
	if up {
		dir = -1
	}
	cur := -1
	id := a.owBrowseFloor()
	for i, floor := range owFloorOrder {
		if floor == id {
			cur = i
			break
		}
	}

	n := len(owFloorOrder)
	var idx int
	if cur < 0 {
		idx = 0
		if up {
			idx = n - 1
		}
	} else {
		idx = (cur + dir + n) % n
	}
	a.owVisitFloor(owFloorOrder[idx])
}

// A closed floor is visited at its doorstep; every shortcut shares this rule.
func (a *App) owVisitFloor(id string) {
	fl, ok := render.OWFloorOf(id)
	if !ok {
		return
	}
	st := &a.ow
	target := fl.Center
	st.Msg = ""
	st.MsgKind = render.OWMessageNeutral
	if !render.OWFloorOpen(id, *st) {
		approach, ok := render.OWFloorApproach(id)
		if !ok {
			return
		}
		target = approach
		a.owSealedMessage(id)
	}
	st.Cursor = target
	st.CameraX, st.CameraSet = float64(target.X), true
}

// Include doorsteps in the wheel's floor order without making them interiors
// for descent, relic spending, or seed entry.
func (a *App) owBrowseFloor() string {
	if id := a.owCursorFloor(); id != "" {
		return id
	}
	for _, id := range owFloorOrder {
		if approach, ok := render.OWFloorApproach(id); ok && approach == a.ow.Cursor {
			return id
		}
	}
	return ""
}

func (a *App) owSealedMessage(id string) {
	if a.ow.Unsealing[id] > 0 {
		a.ow.Msg = copytext.Format("overworld.progress.opening", "floor", render.OWFloorName(id))
		a.ow.MsgKind = render.OWMessageNeutral
	} else if id == hiscore.DepthsFloor {
		a.ow.Msg = depthsLockedMessage
		a.ow.MsgKind = render.OWMessageLocked
	} else {
		a.ow.Msg = copytext.Format("overworld.progress.sealed", "floor", render.OWFloorName(id))
		a.ow.MsgKind = render.OWMessageLocked
	}
}
