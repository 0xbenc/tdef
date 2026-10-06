package render

import "github.com/0xbenc/termtd/internal/copytext"

func RenderCredits(w, h int, pal Colors) *Frame {
	pal.Dim = 252
	f := screenBox(w, h, copytext.Text("branding.render_credits.credits"), []fseg{{key: "esc", text: copytext.Text("branding.render_credits.back")}}, true, pal)
	off := screenOff(h)
	centerPut(f, off+2, copytext.Text("branding.render_credits.termtd"), pal.Gold, true)
	centerPut(f, off+5, copytext.Text("branding.render_credits.a_game_by"), pal.Dim, false)
	centerPut(f, off+6, copytext.Text("branding.render_credits.kairuku_studios"), pal.Bright, true)
	centerPut(f, off+9, copytext.Text("branding.render_credits.created_by"), pal.Dim, false)
	centerPut(f, off+10, copytext.Text("branding.render_credits.0xbenc"), pal.Bright, true)
	centerPut(f, off+13, copytext.Text("branding.render_credits.published_by"), pal.Dim, false)
	centerPut(f, off+14, copytext.Text("branding.render_credits.off_court_creations"), pal.Bright, true)
	return f
}
