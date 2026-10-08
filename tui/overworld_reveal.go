package tui

import (
	"github.com/0xbenc/termtd/internal/copytext"
	"github.com/0xbenc/termtd/render"
	"slices"
)

func (a *App) owLocateGrak() {
	if !a.ow.PlayerIntro {
		return
	}
	a.ow.PlayerIntro = false
	if a.lair != nil {
		a.lair.GrakLocated = true
		a.saveCampaign()
	}
}

func (a *App) owQueueReveal(id string) {
	st := &a.ow
	if st.RevealFloor == id || slices.Contains(st.RevealQueue, id) {
		return
	}
	if st.Unsealing == nil {
		st.Unsealing = map[string]int{}
	}
	st.Unsealing[id] = render.OWUnsealFrames
	st.RevealQueue = append(st.RevealQueue, id)
}

func (a *App) owTickReveals() {
	st := &a.ow
	// Advance standalone seals; queued rooms advance only in their own beat.
	for id, ttl := range st.Unsealing {
		if id == st.RevealFloor || slices.Contains(st.RevealQueue, id) {
			continue
		}
		if ttl > 1 {
			st.Unsealing[id] = ttl - 1
		} else {
			delete(st.Unsealing, id)
			st.Unlocked[id] = true
		}
	}
	if st.RevealFloor == "" {
		if len(st.RevealQueue) == 0 || st.BootTTL > 0 || st.Descending != "" || st.BlastTTL > 0 {
			return
		}
		if st.ReturnFX.Floor != "" && st.ReturnTTL > 180-render.OWVictoryBeatFrames {
			return
		}
		st.RevealFloor, st.RevealQueue = st.RevealQueue[0], st.RevealQueue[1:]
		st.RevealTTL = render.OWRevealFrames
		st.ReturnFX = render.OWReturnFX{}
		st.ReturnTTL, st.ReturnMsg = 0, ""
		return
	}
	st.RevealTTL--
	age := render.OWRevealFrames - st.RevealTTL
	if age >= render.OWRevealPanFrames && st.RevealTTL >= render.OWRevealHoldFrames {
		id := st.RevealFloor
		if id == render.HeartFloorID && age == render.OWRevealPanFrames {
			st.BlastTTL = render.OWBlastFrames
		}
		if ttl := st.Unsealing[id]; ttl > 1 {
			st.Unsealing[id] = ttl - 1
		} else {
			delete(st.Unsealing, id)
			st.Unlocked[id] = true
		}
	}
	if st.RevealTTL <= 0 {
		id := st.RevealFloor
		delete(st.Unsealing, id)
		st.Unlocked[id] = true
		st.RevealFloor, st.RevealTTL = "", 0
		st.ReturnMsg = copytext.Format("overworld.progress.opened", "floor", render.OWFloorName(id))
		st.ReturnKind, st.ReturnTTL = render.OWMessageSuccess, 90
	}
}
