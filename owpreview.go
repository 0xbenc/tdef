//go:build ignore
package main
import (
	"fmt"
	"tdef/render"
)
func main() {
	pal := render.Palette()
	st := render.NewOWState()
	st.RevealAll = true
	fmt.Println("===== 137x45 reveal-all (frame 30) =====")
	fmt.Println(render.RenderOverworld(137, 45, st, 30, pal).Text())
}
