package render

import "testing"

func TestOWSealsBlockRoomInteriors(t *testing.T) {
	for _, id := range []string{"halls", "garden", "depths"} {
		st := NewOWState()
		st.RevealAll = true // a visual preview cannot unlock a floor
		n := owNodeByID(id)
		approach, ok := OWFloorApproach(id)
		if !ok || !OWCanWalk(approach.X, approach.Y, st) {
			t.Fatalf("%s: the doorstep must be reachable", id)
		}
		for x := n.X - n.PW/2; x <= n.X+n.PW/2; x++ {
			for y := n.Y - n.PH/2; y <= n.Y+n.PH/2; y++ {
				if OWCanWalk(x, y, st) {
					t.Fatalf("%s: sealed interior remains walkable at %d,%d", id, x, y)
				}
			}
		}
		st.Unlocked[id] = true
		st.Unsealing[id] = 1
		if OWCanWalk(n.X, n.Y, st) {
			t.Fatalf("%s: entered before the unseal completed", id)
		}
		delete(st.Unsealing, id)
		if !OWCanWalk(n.X, n.Y, st) {
			t.Fatalf("%s: unsealed room remains blocked", id)
		}
	}
}
