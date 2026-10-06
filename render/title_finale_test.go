package render

import (
	"reflect"
	"testing"
)

func TestSiegeStoneLeavesFrameThenHitsFrozenFormation(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {100, 30}, {160, 50}} {
		w, h := size[0], size[1]
		x0, off, cx, floor := max(0, (w-42)/2), screenOff(h), w/2, h-3
		leavesFrame, returns := false, false
		for s := bootComboLaunch; s < bootComboImpact; s++ {
			_, y := bootSiegeStone(x0, off, cx, floor, s)
			if y < 0 {
				leavesFrame = true
			} else if leavesFrame {
				returns = true
			}
		}
		if !leavesFrame || !returns {
			t.Fatalf("%dx%d: stone must leave the frame and fall back into view", w, h)
		}
		x, y := bootSiegeStone(x0, off, cx, floor, bootComboImpact)
		if x != cx || y != floor-2 {
			t.Fatalf("%dx%d: stone misses the center champion", w, h)
		}
		if bootComboImpact-42 < 12 {
			t.Fatal("formation needs a visible frozen pause before impact")
		}
	}
}

func TestFrozenSiegeTargetsHoldTheirPose(t *testing.T) {
	const w, h = 80, 24
	for i, dx := range []int{-9, 0, 9} {
		freezeAt := 42 - i*8
		paint := func(s int) *Frame {
			f := blankFrame(w, h)
			drawBootSiegeTarget(f, w/2+dx, h-3, s, freezeAt, i)
			return f
		}
		if reflect.DeepEqual(paint(0), paint(16)) {
			t.Fatal("approaching formation does not move")
		}
		if !reflect.DeepEqual(paint(50), paint(51)) {
			t.Fatalf("frozen target %d keeps moving", i)
		}
		if reflect.DeepEqual(paint(51), paint(bootComboImpact-5)) {
			t.Fatal("ice must visibly crack before impact")
		}
	}
}

func TestSiegeFinaleIsDeterministicAcrossTerminalSizes(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {100, 30}, {160, 50}} {
		w, h := size[0], size[1]
		for s := 0; s < bootComboLen; s++ {
			boot := titleBootComboStart + s
			a := RenderTitle(w, h, boot, boot, nil, Palette())
			b := RenderTitle(w, h, boot, boot, nil, Palette())
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("%dx%d: finale frame %d is not deterministic", w, h, s)
			}
			if len(a.C) != w*h {
				t.Fatal("finale changed the frame dimensions")
			}
			for _, c := range a.C {
				if c.FG < 0 || c.FG > 255 || c.BG < 0 || c.BG > 255 {
					t.Fatal("finale produced an invalid terminal color")
				}
			}
		}
	}
}

func TestFinaleIgnitionRestoresAllSixLetters(t *testing.T) {
	const w, h = 80, 24
	off, x0 := screenOff(h), (w-42)/2
	for boot := titleBootFlashStart; boot < titleBootFlashEnd; boot++ {
		f := RenderTitle(w, h, boot, boot, nil, Palette())
		for _, step := range titlePenPath {
			c := f.C[(off+2+step.row)*w+x0+step.relX]
			if c.FG != 255 || !c.Bold {
				t.Fatal("ignition must include the two new weapon letters")
			}
		}
	}
}
