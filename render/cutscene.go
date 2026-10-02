package render

import (
	"fmt"
	"math"
	"strings"
)

// Films are authored shots, not timed gameplay events. Dialogue waits for the
// reader; the camera, breath, fire and water continue moving during each hold.
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
	driftX, driftY, push    float64
}

var openingShots = []FilmShot{
	{Name: "The story they tell", Dialogue: "They told it simply: a dragon, a hoard, a hero.", art: "gate", view: filmView{0, 0, 1, 1}, driftX: .018, push: .035},
	{Name: "The commission", Speaker: "GUILDMASTER", Dialogue: "Twenty expeditions. Kill the beast. Bring back the gold.", art: "contract", view: filmView{0, 0, 1, 1}, push: .06},
	{Name: "What they left out", Dialogue: "Nobody mentioned the bowl of water.", art: "bowl", view: filmView{0, 0, 1, 1}, driftX: -.02, push: .04},
	{Name: "A visitor", Dialogue: "Grak came for the same reason as everyone else.", art: "threshold", view: filmView{0, 0, 1, 1}, push: .025},
	{Name: "The dragon", Speaker: "MALGRATH", Dialogue: "You came for the gold.", art: "dragon", view: filmView{0, 0, 1, 1}, push: .035},
	{Name: "The thief", Speaker: "GRAK", Dialogue: "I did.", art: "grak", view: filmView{0, 0, 1, 1}, push: .025},
	{Name: "An old bargain", Speaker: "MALGRATH", Dialogue: "Then take it. There is not much time.", art: "eye", view: filmView{0, 0, 1, 1}, push: .025},
	{Name: "A small decision", Speaker: "GRAK", Dialogue: "Gold doesn't get thirsty.", art: "offering", view: filmView{0, 0, 1, 1}, driftY: .015, push: .025},
	{Name: "The bowl", Dialogue: "He set down the coin. He brought the water closer.", art: "together", view: filmView{0, 0, 1, 1}, driftX: .025, push: .02},
	{Name: "The question", Speaker: "MALGRATH", Dialogue: "You know what is coming?", art: "eye", view: filmView{.015, 0, .985, 1}, push: .02},
	{Name: "At the door", Speaker: "GRAK", Dialogue: "Twenty expeditions.", art: "door", view: filmView{0, 0, 1, 1}, push: .06},
	{Name: "A builder's answer", Speaker: "GRAK", Dialogue: "Then I'd better fix the door.", art: "mallet", view: filmView{0, 0, 1, 1}, driftY: -.018, push: .035},
	{Name: "The last monster", Dialogue: "Hold the lair. Keep the fire.", art: "resolve", view: filmView{0, 0, 1, 1}, push: .025},
}

var endingShots = []FilmShot{
	{Name: "After the twentieth", Dialogue: "The last expedition leaves its banner in the dust.", art: "aftermath", view: filmView{0, 0, 1, 1}, driftX: .02, push: .035},
	{Name: "A small fire", Speaker: "MALGRATH", Dialogue: "Grak?", art: "ember", view: filmView{0, 0, 1, 1}, push: .035},
	{Name: "The answer", Speaker: "GRAK", Dialogue: "Here.", art: "grak", view: filmView{0, 0, 1, 1}, push: .02},
	{Name: "What remains", Speaker: "MALGRATH", Dialogue: "How much did they take?", art: "eye", view: filmView{0, 0, 1, 1}, driftY: .01, push: .03},
	{Name: "An inventory", Speaker: "GRAK", Dialogue: "Nothing we need.", art: "bowl", view: filmView{0, 0, 1, 1}, driftX: .015, push: .025},
	{Name: "The promise he can keep", Speaker: "GRAK", Dialogue: "I cannot give you back your wings.", art: "touch", view: filmView{0, 0, 1, 1}, push: .025},
	{Name: "The promise he was asked for", Speaker: "MALGRATH", Dialogue: "I wasn't asking.", art: "eye", view: filmView{.015, 0, .985, 1}, push: .02},
	{Name: "Another breath", Dialogue: "The ember catches. A breath follows it. Then another.", art: "ember", view: filmView{0, 0, 1, 1}, push: -.025},
	{Name: "A familiar question", Speaker: "MALGRATH", Dialogue: "Is that water?", art: "dragon", view: filmView{0, 0, 1, 1}, push: .035},
	{Name: "An honest answer", Speaker: "GRAK", Dialogue: "Mostly. The bowl leaks.", art: "grak", view: filmView{0, 0, 1, 1}, driftX: -.01, push: .025},
	{Name: "A little work", Speaker: "MALGRATH", Dialogue: "Then we will need another.", art: "together", view: filmView{0, 0, 1, 1}, push: -.025},
	{Name: "Tomorrow", Speaker: "GRAK", Dialogue: "Tomorrow.", art: "morning", view: filmView{0, 0, 1, 1}, driftY: -.015, push: -.02},
	{Name: "For now", Speaker: "MALGRATH", Dialogue: "For now, stay.", art: "together", view: filmView{0, 0, 1, 1}, driftX: .012, push: .025},
	{Name: "Still here", Speaker: "GRAK", Dialogue: "I'm here.", art: "rest", view: filmView{0, 0, 1, 1}, push: -.025},
	{Name: "The heart held", Dialogue: "The fire is small. There is someone to tend it.", art: "rest", view: filmView{0, 0, 1, 1}, push: -.04},
}

func FilmShots(film Film) []FilmShot {
	if film == FilmEnding {
		return endingShots
	}
	return openingShots
}
func FilmTitle(film Film) string {
	if film == FilmEnding {
		return "THE HEART HELD"
	}
	return "THE LAST MONSTER"
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
	return min(n, max(0, (st.Frame-18)/2))
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
		journalCenter(f, h/2, "enlarge to 62 × 19 to view", 180, true)
		journalCenter(f, h-2, "esc skip · q quit", 240, false)
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
	master := filmFrame(240, 80)
	paintFilmShot(master, st, shot.art)
	// Ease once, then settle. Camera movement never loops through dialogue.
	t := math.Min(1, float64(max(0, st.Frame))/210)
	t = t * t * (3 - 2*t)
	v := shot.view
	v.w *= 1 - shot.push*t
	v.h *= 1 - shot.push*t
	v.x += (shot.view.w-v.w)/2 + shot.driftX*t
	v.y += (shot.view.h-v.h)/2 + shot.driftY*t
	cinemaProject(f, master, Rect{stageX, stageY, stageW, stageH}, v)
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
	// The storyboard title is subtle; it lends each cut an authored beat.
	meta := fmt.Sprintf("%02d / %02d  ·  %s", st.Shot+1, len(shots), shot.Name)
	journalCenter(f, 2, fitMsg(meta, w-6), 240, false)
	col := 252
	if shot.Speaker == "GRAK" {
		col = 150
	} else if shot.Speaker == "MALGRATH" {
		col = 180
	} else if shot.Speaker == "GUILDMASTER" {
		col = 174
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
	hint := "enter / space reveal · ← previous · esc skip · q quit"
	if CutsceneVisible(st) >= len([]rune(shot.Dialogue)) {
		hint = "enter / space next · ← previous · esc skip · q quit"
	}
	if st.Shot == len(shots)-1 && CutsceneVisible(st) >= len([]rune(shot.Dialogue)) {
		hint = "enter / space return · ← previous · esc skip · q quit"
	}
	journalCenter(f, h-2, hint, 240, false)
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
