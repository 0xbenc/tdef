package render

import (
	"strings"
	"testing"
)

func TestHelpContentAccessibleAtEverySize(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {100, 30}, {120, 40}} {
		for page, topic := range helpPages() {
			var seen strings.Builder
			for scroll := 0; scroll <= HelpMaxScroll(size[0], size[1], page); scroll++ {
				f := RenderHelpPage(size[0], size[1], HelpState{Page: page, Scroll: scroll}, Palette())
				seen.WriteString(f.Text())
				if !strings.Contains(f.Text(), "esc back") {
					t.Fatal("help lost its exit control")
				}
			}
			for _, section := range topic.sections {
				if !strings.Contains(seen.String(), section.title) {
					t.Fatalf("%v page %d hides %q", size, page, section.title)
				}
				for _, control := range section.controls {
					if !strings.Contains(seen.String(), control[0]) {
						t.Fatalf("%v hides control %s", size, control[0])
					}
				}
			}
		}
	}
}

func TestHelpTabsMatchMouseTargets(t *testing.T) {
	for _, w := range []int{62, 80, 120} {
		f := RenderHelpPage(w, 24, HelpState{}, Palette())
		for x := 0; x < w; x++ {
			c := f.C[4*w+x]
			if c.R >= '1' && c.R <= '5' {
				if c.FG != 167 || !c.Bold || HelpTopicAt(w, 24, x, 4) != int(c.R-'1') {
					t.Fatal("help tab display and click target disagree")
				}
			}
		}
	}
}
