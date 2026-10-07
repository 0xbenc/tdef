package render

import (
	"fmt"
	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/internal/copytext"
	"strings"
)

func recruitLesson(kind game.TowerKind) string {
	switch kind {
	case game.TowerGunner:
		return copytext.Text("ui.training.gunner")
	case game.TowerFrost:
		return copytext.Text("ui.training.frost")
	case game.TowerCannon:
		return copytext.Text("ui.training.cannon")
	case game.TowerSniper:
		return copytext.Text("ui.training.ranger")
	case game.TowerFlak:
		return copytext.Text("ui.training.slingers")
	case game.TowerTesla:
		return copytext.Text("ui.training.lightning")
	case game.TowerMortar:
		return copytext.Text("ui.training.trebuchet")
	case game.TowerRuneforge:
		return copytext.Text("ui.training.runeforge")
	case game.TowerHookmaster:
		return copytext.Text("ui.training.hookmaster")
	case game.TowerSappers:
		return copytext.Text("ui.training.sappers")
	default:
		return copytext.Text("ui.training.witch")
	}
}

func recruitmentLines(text string, width int) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(text) {
		if line != "" && len([]rune(line+" "+word)) > width {
			lines = append(lines, line)
			line = ""
		}
		if line != "" {
			line += " "
		}
		line += word
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func drawRecruitment(f *Frame, g *game.State, pal Colors) {
	if !g.LessonPending {
		return
	}
	w := min(64, f.W-4)
	lines := recruitmentLines(recruitLesson(g.Lesson), w-6)
	funding := copytext.Text("ui.training.starting_gold")
	if g.Lesson != game.TowerGunner {
		funding = copytext.Format("ui.training.funding", "gold", fmt.Sprint(game.TowerSpecs[g.Lesson].Cost[0]))
	}
	lines = append(lines, "", funding)
	lines = append(lines, recruitmentLines(copytext.Text("ui.training.safe"), w-6)...)
	h := len(lines) + 7
	x, y := (f.W-w)/2, (f.H-h)/2
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			f.Set(xx, yy, Cell{R: ' ', BG: 233})
		}
	}
	drawRoundedBox(f, x, y, w, h, 180)
	putString(f, x+3, y+1, copytext.Text("ui.training.title"), 180, 233, true)
	putString(f, x+3, y+2, game.TowerSpecs[g.Lesson].Name, pal.Tower[g.Lesson], 233, true)
	for i, line := range lines {
		putString(f, x+3, y+4+i, line, 252, 233, false)
	}
	putString(f, x+3, y+h-2, copytext.Text("ui.training.continue"), 252, 233, true)
}
