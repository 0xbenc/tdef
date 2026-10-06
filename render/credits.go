package render

func RenderCredits(w, h int, pal Colors) *Frame {
	pal.Dim = 252
	f := screenBox(w, h, "CREDITS", []fseg{{key: "esc", text: " back"}}, true, pal)
	off := screenOff(h)
	centerPut(f, off+2, "TERMTD", pal.Gold, true)
	centerPut(f, off+5, "A game by", pal.Dim, false)
	centerPut(f, off+6, "Kairuku Studios", pal.Bright, true)
	centerPut(f, off+9, "Created by", pal.Dim, false)
	centerPut(f, off+10, "0xbenc", pal.Bright, true)
	centerPut(f, off+13, "Published by", pal.Dim, false)
	centerPut(f, off+14, "Off Court Creations", pal.Bright, true)
	return f
}
