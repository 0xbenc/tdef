package render

import "sync"

// Retain only the latest size for each static layer. Resizing replaces that
// layer rather than accumulating full-size images. Cached frames are read-only.
var filmArtworkCache = struct {
	sync.Mutex
	layers map[string]*Frame
}{layers: make(map[string]*Frame)}

func filmCachedArtwork(key string, w, h int, transparent bool, paint func(*Frame)) *Frame {
	filmArtworkCache.Lock()
	defer filmArtworkCache.Unlock()
	if f := filmArtworkCache.layers[key]; f != nil && f.W == w && f.H == h {
		return f
	}
	var f *Frame
	if transparent {
		f = filmTransparentFrame(w, h)
	} else {
		f = filmFrame(w, h)
	}
	paint(f)
	filmArtworkCache.layers[key] = f
	return f
}
