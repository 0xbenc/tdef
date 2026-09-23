//go:build ignore

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"tdef/game"
	"tdef/render"
)

// throwaway look-dev harness: renders the overworld for a named state at a
// named size and frame. Usage: go run owview_tmp.go <state> <size> <frame> [heat]

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: owview <state> <size> <frame> [heat]")
		os.Exit(1)
	}
	st := mkState(os.Args[1])
	w, h, ok := mkSize(os.Args[2])
	if !ok {
		fmt.Fprintln(os.Stderr, "bad size: "+os.Args[2])
		os.Exit(1)
	}
	frame, _ := strconv.Atoi(os.Args[3])
	f := render.RenderOverworld(w, h, st, frame, render.Palette())
	fmt.Printf("===== %s %dx%d frame %d =====\n", os.Args[1], w, h, frame)
	fmt.Println(f.Text())
	if len(os.Args) > 4 && os.Args[4] == "heat" {
		fmt.Println("---- BG heat per grid cell (2-hex) ----")
		l := render.GameLayout(45, 13, w, h)
		for gy := 0; gy < 13; gy++ {
			var b strings.Builder
			for gx := 0; gx < 45; gx++ {
				c := f.C[(l.Oy+gy*l.Scale)*f.W+l.Ox+gx*l.Scale]
				fmt.Fprintf(&b, "%02x ", c.BG)
			}
			fmt.Println(b.String())
		}
	}
}

func mkSize(s string) (int, int, bool) {
	switch s {
	case "min":
		return 62, 19, true
	case "92":
		return 92, 32, true
	case "137":
		return 137, 45, true
	case "182":
		return 182, 58, true
	}
	return 0, 0, false
}

func vec(x, y int) game.Vec { return game.Vec{X: x, Y: y} }

func floorCenter(id string) game.Vec {
	fl, _ := render.OWFloorOf(id)
	return fl.Center
}

func mkState(name string) render.OWState {
	st := render.NewOWState()
	switch name {
	case "fresh":
		// as-is
	case "rift":
		st.Records["rift"] = render.OWRec{Cleared: true, BestWave: 20, LastWave: 20, LastWon: true}
		st.Scores["rift"] = 9100
		st.Unlocked["halls"] = true
		st.Hearts = 1
		st.FirstRun = false
	case "unseal":
		st = mkState("rift")
		st.Unsealing["halls"] = 30
	case "allheld":
		for _, id := range []string{"rift", "halls", "garden", "rotunda"} {
			st.Records[id] = render.OWRec{Cleared: true, BestWave: 20, LastWave: 20, LastWon: true}
			st.Scores[id] = 9000
		}
		st.Unlocked["halls"] = true
		st.Unlocked["garden"] = true
		st.Unlocked["depths"] = true
		st.Hearts = 4
		st.BossReady = true
		st.Tokens = 2
		st.FirstRun = false
	case "bossdone":
		st = mkState("allheld")
		st.BossDone = true
	case "descend":
		st.Descending = "rotunda"
		st.DescendTTL = 14
	case "reveal":
		st.RevealAll = true
	case "relic":
		st.Cursor = floorCenter("rotunda")
		st.Tokens = 3
		st.RelicMenu = true
	case "banner":
		st = mkState("rift")
		st.ReturnMsg = "the Long Halls breaks open"
		st.ReturnTTL = 60
	case "banner2":
		st = mkState("rift")
		st.ReturnMsg = "the Unmapped Depths held — 20/20 · the lair stands steadier"
		st.ReturnTTL = 100
	case "walk":
		st.PushTrail(vec(12, 6))
		st.PushTrail(vec(13, 6))
		st.TrailAge = []int{6, 18}
		st.Cursor = vec(14, 6)
	case "returnfx":
		st = mkState("rift")
		st.ReturnFX = render.OWReturnFX{Floor: "rift", Won: true}
		st.ReturnTTL = 170 // q~0.05: strong flash, ring near the hub
	case "blast":
		st = mkState("allheld")
		st.BlastTTL = 36 // mid-blast: front ~d12
	case "blast2":
		st = mkState("allheld")
		st.BlastTTL = 10 // late blast: front near the far corner
	case "boot1":
		st.BootTTL = 89 // first frame: only the Rotunda lit
	case "boot3":
		st.BootTTL = 40 // mid-sweep
	case "boot7":
		st.BootTTL = 15 // last narration frame
	case "boot9":
		st.BootTTL = 1 // the seam: must equal "fresh"
	case "broken":
		st = mkState("rift")
		st.Records["rift"] = render.OWRec{BestWave: 7, LastWave: 7, LastWon: false}
	case "halls":
		st = mkState("rift")
		st.Records["halls"] = render.OWRec{Cleared: true, BestWave: 20, LastWave: 20, LastWon: true}
		st.Scores["halls"] = 8700
		st.Unlocked["garden"] = true
		st.Hearts = 2
	}
	return st
}
