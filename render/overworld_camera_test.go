package render

import (
	"strings"
	"testing"

	"github.com/0xbenc/tdef/game"
)

// Camera, renderer and hit testing must agree even when a room is partially
// clipped, a resize changes scale, or the eased camera lags behind Grak.
func TestOWScrollingPicking(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {80, 24}, {92, 32}, {137, 45}, {182, 58}, {400, 70}} {
		w, h := size[0], size[1]
		for _, n := range owNodes {
			st := NewOWState()
			st.Cursor = game.Vec{X: n.X, Y: n.Y}
			st.CameraSet, st.CameraX = true, 6 // simulate a long jump before camera settles
			l := OverworldLayout(w, h, st)
			f := RenderOverworld(w, h, st, 30, Palette())
			px, py := l.center(n.X, n.Y)
			if px <= 1 || px >= w-2 || py < ChromeTop || py >= h-ChromeBot {
				t.Fatalf("%dx%d: %s player outside viewport at %d,%d", w, h, n.ID, px, py)
			}
			if got := f.C[py*w+px].R; got != '@' {
				t.Fatalf("%dx%d: %s player = %q", w, h, n.ID, got)
			}
			for y := ChromeTop; y < h-ChromeBot; y++ {
				for x := 1; x < w-1; x++ {
					v, hit := OWFloorAtFrame(w, h, x, y, st)
					if !hit {
						continue
					}
					fl, ok := OWFloorAt((x-l.Ox)/l.Scale, (y-l.Oy)/l.Scale)
					if !ok || v != fl.Center {
						t.Fatalf("%dx%d: click (%d,%d) disagrees with camera", w, h, x, y)
					}
				}
			}
			for _, p := range [][2]int{{0, py}, {w - 1, py}, {px, 1}, {px, h - 4}, {-1, py}, {w, py}} {
				if _, hit := OWFloorAtFrame(w, h, p[0], p[1], st); hit {
					t.Fatalf("%dx%d: chrome/outside click hits a floor at %v", w, h, p)
				}
			}
		}
	}
}

func TestOWScrollPreservesChrome(t *testing.T) {
	for _, size := range [][2]int{{62, 19}, {92, 32}, {182, 58}} {
		w, h := size[0], size[1]
		st := NewOWState()
		st.FirstRun, st.RevealAll = false, true
		left := RenderOverworld(w, h, st, 30, Palette())
		st.Cursor = game.Vec{X: 65, Y: 9}
		right := RenderOverworld(w, h, st, 30, Palette())
		for y := h - ChromeBot; y < h; y++ {
			// The context row intentionally changes with the occupied floor.
			if y == h-2 {
				continue
			}
			for x := 0; x < w; x++ {
				if left.C[y*w+x] != right.C[y*w+x] {
					t.Fatalf("%dx%d: scrolling changed footer at %d,%d", w, h, x, y)
				}
			}
		}
		for y := 0; y < h; y++ {
			if left.C[y*w] != right.C[y*w] || left.C[y*w+w-1] != right.C[y*w+w-1] {
				t.Fatalf("%dx%d: scrolling overwrote frame on row %d", w, h, y)
			}
		}
	}
}

func TestOWSmallWindow(t *testing.T) {
	for _, size := range [][2]int{{40, 12}, {62, 18}, {61, 19}} {
		f := RenderOverworld(size[0], size[1], NewOWState(), 0, Palette())
		if f.W != size[0] || f.H != size[1] || !strings.Contains(f.Text(), "terminal too small") {
			t.Fatalf("small window %v did not render its resize notice", size)
		}
	}
}

// Extending the cavern must not turn grayscale shade 232 into ANSI white 231.
func TestOWDistantStoneStaysDark(t *testing.T) {
	st := NewOWState()
	st.Cursor = game.Vec{X: OWW - 1, Y: 6}
	f := RenderOverworld(92, 32, st, 30, Palette())
	for y := ChromeTop; y < f.H-ChromeBot; y++ {
		for x := 1; x < f.W-1; x++ {
			if f.C[y*f.W+x].BG == 231 {
				t.Fatalf("white cave background at %d,%d", x, y)
			}
		}
	}
}
