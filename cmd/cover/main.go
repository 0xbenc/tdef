// Command cover prints the itch.io cover as a real ANSI terminal frame.
package main

import (
	"fmt"

	"github.com/0xbenc/termtd/render"
)

func main() {
	fmt.Print(render.RenderCover().ANSI())
}
