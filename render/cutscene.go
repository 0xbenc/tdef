package render

import (
	"fmt"
	"github.com/0xbenc/termtd/internal/copytext"
	"math"
	"strings"
)

// Films are authored shots, not timed gameplay events. Dialogue waits for the
// reader. Establishing shots pan across wider scenes; dialogue holds still.
type Film int

const (
	FilmOpening Film = iota
	FilmEnding
)

type filmView struct{ x, y, w, h float64 }
type FilmShot struct {
	Name, Speaker, Dialogue string
	art                     string
	view                    filmView
}

var openingShots = withOpeningCopy([]FilmShot{
	{art: "village", view: filmView{0, 0, 1, 1}},
	{art: "village-search", view: filmView{0, 0, 1, 1}},
	{art: "ruined-road", view: filmView{0, 0, 1, 1}},
	{art: "capital-vault", view: filmView{0, 0, 1, 1}},
	{art: "eye", view: filmView{0, 0, 1, 1}},
	{art: "grak", view: filmView{0, 0, 1, 1}},
	{art: "dragon", view: filmView{0, 0, 1, 1}},
	{art: "grak", view: filmView{0, 0, 1, 1}},
	{art: "eye", view: filmView{0, 0, 1, 1}},
	{art: "dragon", view: filmView{0, 0, 1, 1}},
	{art: "grak", view: filmView{0, 0, 1, 1}},
	{art: "eye", view: filmView{0, 0, 1, 1}},
	{art: "mallet", view: filmView{0, 0, 1, 1}},
	{art: "hoard", view: filmView{0, 0, 1, 1}},
	{art: "grak", view: filmView{0, 0, 1, 1}},
	{art: "eye", view: filmView{0, 0, 1, 1}},
	{art: "resolve", view: filmView{0, 0, 1, 1}},
})

var endingShots = withEndingCopy([]FilmShot{
	{art: "ending-fallen", view: filmView{0, 0, 1, 1}},
	{art: "ending-rescue", view: filmView{0, 0, 1, 1}},
	{art: "ending-healer", view: filmView{0, 0, 1, 1}},
	{art: "ending-evacuation", view: filmView{0, 0, 1, 1}},
	{art: "ending-pursuit", view: filmView{0, 0, 1, 1}},
	{art: "ending-maze", view: filmView{0, 0, 1, 1}},
	{art: "ending-depths", view: filmView{0, 0, 1, 1}},
	{art: "ending-earth", view: filmView{0, 0, 1, 1}},
	{art: "ending-supplies", view: filmView{0, 0, 1, 1}},
	{art: "ending-return", view: filmView{0, 0, 1, 1}},
	{art: "ending-hoard", view: filmView{0, 0, 1, 1}},
	{art: "grak", view: filmView{0, 0, 1, 1}},
	{art: "ending-maze", view: filmView{0, 0, 1, 1}},
	{art: "eye", view: filmView{0, 0, 1, 1}},
	{art: "ending-supplies", view: filmView{0, 0, 1, 1}},
	{art: "grak", view: filmView{0, 0, 1, 1}},
	{art: "ending-hoard", view: filmView{0, 0, 1, 1}},
})

func FilmShots(film Film) []FilmShot {
	if film == FilmEnding {
		return endingShots
	}
	return openingShots
}
func FilmTitle(film Film) string {
	if film == FilmEnding {
		return copytext.Text("ui.film_title.the_other_side")
	}
	return copytext.Text("ui.film_title.the_last_monster")
}

type CutsceneState struct {
	Film        Film
	Shot, Frame int // Frame is local to this shot, at 30fps.
	Revealed    bool
}

func CutsceneVisible(st CutsceneState) int {
	shots := FilmShots(st.Film)
	shot := shots[max(0, min(st.Shot, len(shots)-1))]
	n := len([]rune(shot.Dialogue))
	if st.Revealed {
		return n
	}
	return min(n, max(0, st.Frame-18))
}

// The same wrap is used before and after reveal: lines never jump as words
// appear. Text reserves enough rows for the entire caption at the current width.
func filmLines(s string, width int) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(s) {
		if line != "" && len([]rune(line))+1+len([]rune(word)) > width {
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

func RenderCutscene(w, h int, st CutsceneState) *Frame {
	w, h = max(0, w), max(0, h)
	f := filmFrame(w, h)
	if w < 62 || h < 19 {
		journalCenter(f, h/2, copytext.Text("ui.render_cutscene.enlarge_to_62_19_to_view"), 180, true)
		journalCenter(f, h-2, copytext.Format("ui.render_cutscene.esc_skip_q_quit", "escape", "esc", "quit", "q"), titleTextFG, false)
		return f
	}
	shots := FilmShots(st.Film)
	shot := shots[max(0, min(st.Shot, len(shots)-1))]
	captionW := min(84, w-10)
	lines := filmLines(shot.Dialogue, captionW)
	captionY := h - 4 - len(lines)
	speakerY := captionY - 1
	stageH := min(speakerY-4, (w-6)/3)
	stageW := stageH * 3
	stageX, stageY := (w-stageW)/2, 3+(speakerY-4-stageH)/2
	stage := Rect{stageX, stageY, stageW, stageH}
	if filmIsPanorama(shot.art) {
		paintFilmPanorama(f, stage, shot.art, st.Frame)
	} else {
		// Gestures finish once, leaving a stable composition for reading.
		pose := st
		pose.Frame = min(120, max(0, st.Frame))
		var master *Frame
		if shot.art == "ending-supplies" || shot.art == "ending-return" {
			master = filmCachedArtwork(shot.art, 240, 80, false, func(f *Frame) {
				paintFilmShot(f, pose, shot.art)
			})
		} else {
			master = filmFrame(240, 80)
			paintFilmShot(master, pose, shot.art)
		}
		cinemaProject(f, master, stage, shot.view)
	}
	// A brief fade-in on each cut; no flickering or flashing transitions.
	fade := math.Min(1, float64(max(0, st.Frame))/12)
	if fade < 1 {
		var colors [256]int
		for c := range colors {
			colors[c] = filmFade(c, fade)
		}
		for yy := stageY; yy < stageY+stageH; yy++ {
			for xx := stageX; xx < stageX+stageW; xx++ {
				i := yy*w + xx
				c := f.C[i]
				c.FG = colors[c.FG]
				c.BG = colors[c.BG]
				f.C[i] = c
			}
		}
	}
	journalCenter(f, 1, FilmTitle(st.Film), 180, true)
	// Match storyboard labels and prompts to the other screens' soft white.
	meta := fmt.Sprintf("%02d / %02d  ·  %s", st.Shot+1, len(shots), shot.Name)
	journalCenter(f, 2, fitMsg(meta, w-6), titleTextFG, false)
	col := 252
	if shot.Speaker == "GRAK" {
		col = 150
	} else if shot.Speaker == "MALGRATH" {
		col = 180
	} else if shot.Speaker == "GUILDMASTER" {
		col = 174
	} else if shot.Speaker == "HEALER" {
		col = 117
	}
	journalCenter(f, speakerY, shot.Speaker, col, true)
	left := (w - captionW) / 2
	n := CutsceneVisible(st)
	for i, line := range lines {
		rr := []rune(line)
		visible := min(len(rr), max(0, n))
		putString(f, left, captionY+i, string(rr[:visible]), col, 233, false)
		n -= len(rr) + 1
	}
	hint := copytext.Format("ui.render_cutscene.enter_space_reveal_previous_esc_skip_q", "advance", "enter / space", "escape", "esc", "quit", "q")
	if CutsceneVisible(st) >= len([]rune(shot.Dialogue)) {
		hint = copytext.Format("ui.render_cutscene.enter_space_next_previous_esc_skip_q", "advance", "enter / space", "escape", "esc", "quit", "q")
		if !CutsceneCanAdvance(st) {
			hint = copytext.Format("ui.render_cutscene.scene_playing_previous_esc_skip_q_quit", "escape", "esc", "quit", "q")
		}
	}
	if st.Shot == len(shots)-1 && CutsceneVisible(st) >= len([]rune(shot.Dialogue)) {
		hint = copytext.Format("ui.render_cutscene.enter_space_return_previous_esc_skip_q", "advance", "enter / space", "escape", "esc", "quit", "q")
	}
	journalCenter(f, h-2, hint, titleTextFG, false)
	return f
}

func filmFrame(w, h int) *Frame {
	f := &Frame{W: w, H: h, C: make([]Cell, w*h)}
	for i := range f.C {
		f.C[i] = Cell{R: ' ', FG: 240, BG: 233}
	}
	return f
}
func cinemaProject(dst, src *Frame, r Rect, v filmView) {
	for y := 0; y < r.H; y++ {
		for x := 0; x < r.W; x++ {
			u := v.x + (float64(x)+.5)/float64(r.W)*v.w
			vv := v.y + (float64(y)+.5)/float64(r.H)*v.h
			sx, sy := int(u*float64(src.W)), int(vv*float64(src.H))
			if sx >= 0 && sx < src.W && sy >= 0 && sy < src.H {
				dst.Set(r.X+x, r.Y+y, src.C[sy*src.W+sx])
			}
		}
	}
}
func filmFade(c int, t float64) int {
	if c == 233 {
		return c
	}
	r, g, b := xtermRGB(c)
	// Find the nearest actual palette color, rather than interpolating palette
	// indices (which would flash green/cyan through the red and amber artwork).
	br, bg, bb := xtermRGB(233)
	r = int(float64(br) + float64(r-br)*t)
	g = int(float64(bg) + float64(g-bg)*t)
	b = int(float64(bb) + float64(b-bb)*t)
	best, dist := 233, 1<<30
	for i := 16; i < 256; i++ {
		rr, gg, bbb := xtermRGB(i)
		d := (r-rr)*(r-rr) + (g-gg)*(g-gg) + (b-bbb)*(b-bbb)
		if d < dist {
			best, dist = i, d
		}
	}
	return best
}
func xtermRGB(c int) (int, int, int) {
	if c >= 232 {
		n := 8 + (c-232)*10
		return n, n, n
	}
	if c >= 16 {
		c -= 16
		values := [6]int{0, 95, 135, 175, 215, 255}
		return values[c/36], values[c/6%6], values[c%6]
	}
	values := [16][3]int{{0, 0, 0}, {128, 0, 0}, {0, 128, 0}, {128, 128, 0}, {0, 0, 128}, {128, 0, 128}, {0, 128, 128}, {192, 192, 192}, {128, 128, 128}, {255, 0, 0}, {0, 255, 0}, {255, 255, 0}, {0, 0, 255}, {255, 0, 255}, {0, 255, 255}, {255, 255, 255}}
	v := values[max(0, min(15, c))]
	return v[0], v[1], v[2]
}
