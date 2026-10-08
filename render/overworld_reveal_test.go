package render

import (
	"reflect"
	"strings"
	"testing"
)

func TestGrakIntroVisibleAtEverySupportedSize(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {100, 32}, {137, 45}, {236, 58}} {
		st := NewOWState()
		st.PlayerIntro = true
		f := RenderOverworld(size[0], size[1], st, 0, Palette())
		if !strings.Contains(f.Text(), "YOU — GRAK") || !strings.Contains(f.Text(), "Arrows / WASD to move") {
			t.Fatalf("%v hides player introduction", size)
		}
		l := OverworldLayout(size[0], size[1], st)
		x, y := l.center(st.Cursor.X, st.Cursor.Y)
		player := f.C[y*f.W+x]
		if player.R != '@' || player.BG != 24 || !player.Bold {
			t.Fatal("Grak lost his distinctive single-cell marker")
		}
		st.BootTTL = OWBootFrames
		if strings.Contains(RenderOverworld(size[0], size[1], st, 0, Palette()).Text(), "YOU — GRAK") {
			t.Fatal("orientation interrupted arrival animation")
		}
	}
}

func TestAreaRevealsAreVisibleDeterministicAndKeepChrome(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {100, 32}, {137, 45}, {236, 58}} {
		for _, id := range []string{"rift", "halls", "garden", "depths", HeartFloorID} {
			st := NewOWState()
			st.RevealFloor = id
			st.RevealTTL = OWRevealFrames - 50
			st.Unsealing[id] = OWUnsealFrames - 26
			for _, node := range owNodes {
				st.Unlocked[node.ID] = true
			}
			nID := id
			if nID == HeartFloorID {
				nID = "rotunda"
				st.BossReady = true
			}
			n := owNodeByID(nID)
			f := RenderOverworld(size[0], size[1], st, 50, Palette())
			l := OverworldLayout(size[0], size[1], st)
			x, y := l.center(n.X, n.Y)
			if x < 1 || x >= f.W-1 || y < ChromeTop || y >= f.H-ChromeBot {
				t.Fatalf("%v: %s remains offscreen", size, id)
			}
			if !strings.Contains(f.Text(), "NEW AREA:") {
				t.Fatal("unlock lacks a dedicated announcement")
			}
			again := RenderOverworld(size[0], size[1], st, 50, Palette())
			if !reflect.DeepEqual(f, again) {
				t.Fatal("reveal is nondeterministic")
			}
			for yy := ChromeTop; yy < f.H-ChromeBot; yy++ {
				if f.C[yy*f.W].R != '│' || f.C[yy*f.W+f.W-1].R != '│' {
					t.Fatal("reveal escaped into frame chrome")
				}
			}
		}
	}
}

func TestPendingRoomRemainsSealedDuringOtherReveals(t *testing.T) {
	st := NewOWState()
	st.RevealQueue = []string{"halls"}
	st.Unlocked["halls"] = true // save refresh already knows the floor has been unlocked
	st.Unsealing["halls"] = OWUnsealFrames
	if owNodeStatus(st, "halls") != OWSealed || OWFloorOpen("halls", st) {
		t.Fatal("queued room revealed before its own beat")
	}
}
