package render

import (
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/internal/copytext"
	"strings"
	"testing"
)

func TestRecruitmentRosterLeavesLockedSlotsEmpty(t *testing.T) {
	m, _ := game.LoadLevel("hub")
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		for _, stage := range []int{1, 2} {
			s := game.NewState(m)
			mask := uint32(1 << game.TowerGunner)
			page := 0
			if stage == 2 {
				mask = game.MainTowersMask | 1<<game.TowerRuneforge
				page = 1
			}
			s.ConfigureTraining(stage, mask)
			s.LessonPending = false
			ui := UI{Selected: -1, RosterPage: page}
			check := func() {
				f := Render(s, &ui, Palette(), size[0], size[1], 0)
				for _, slot := range TowerSlots(size[0], size[1], page) {
					key := f.C[slot.Y*f.W+slot.X].R
					if s.TowerAvailable(slot.Kind) {
						if key != rune('0'+slot.Key) {
							t.Fatalf("available defender lost its fixed key: %+v", slot)
						}
					} else {
						for x := slot.X; x < slot.X+slot.W; x++ {
							if r := f.C[slot.Y*f.W+x].R; r != 0 && r != ' ' {
								t.Fatalf("locked defender visible at %v: %+v", size, slot)
							}
						}
					}
				}
			}
			check()
			for _, slot := range TowerSlots(size[0], size[1], page) {
				s.UnlockedTowers |= 1 << slot.Kind
				check()
			}
		}
	}
}

func TestRecruitmentCardsFit(t *testing.T) {
	m, _ := game.LoadLevel("hub")
	for _, size := range [][2]int{{62, 19}, {80, 24}, {120, 40}} {
		for k := game.TowerKind(0); k < game.TowerCount; k++ {
			s := game.NewState(m)
			s.ConfigureTraining(1, 1<<game.TowerGunner)
			s.Lesson = k
			s.LessonPending = true
			f := Render(s, &UI{Selected: -1, Level: "hub", Paused: true}, Palette(), size[0], size[1], 0)
			for _, text := range []string{game.TowerSpecs[k].Name, copytext.Text("ui.training.continue")} {
				if !strings.Contains(f.Text(), text) {
					t.Fatalf("%v %v clipped %q", size, k, text)
				}
			}
			lines := recruitmentLines(recruitLesson(k), min(64, size[0]-4)-6)
			lines = append(lines, recruitmentLines(copytext.Text("ui.training.safe"), min(64, size[0]-4)-6)...)
			for _, line := range lines {
				if !strings.Contains(f.Text(), line) {
					t.Fatalf("%v %v clipped lesson line", size, k)
				}
			}
		}
	}
}
