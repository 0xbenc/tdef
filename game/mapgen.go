package game

import (
	"math/rand/v2"
)

// corridorBias is the chance to keep carving straight, which produces longer
// snake-like chokepoints instead of tight grid mazes.
const corridorBias = 0.35

// GenerateMap builds a perfect maze and uses its unique route from the left
// edge to the right edge as the enemy path. All other open cells are buildable.
func GenerateMap(seed int64, w, h int) (*Map, error) {
	if w%2 == 0 {
		w++
	}
	if h%2 == 0 {
		h++
	}
	if w < 9 || h < 7 {
		w, h = 45, 13
	}
	rng := rand.New(rand.NewPCG(uint64(seed), 0x9e3779b97f4a7c15))
	open := make([]bool, w*h)
	start := Vec{1, 1}
	open[start.Y*w+start.X] = true
	stack := []Vec{start}
	headings := []Vec{{0, 0}}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		heading := headings[len(headings)-1]
		dirs := [4]Vec{{2, 0}, {-2, 0}, {0, 2}, {0, -2}}
		rng.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
		// Corridor bias: keep going straight to carve longer chokepoints.
		if heading != (Vec{}) && rng.Float64() < corridorBias {
			for i := range dirs {
				if dirs[i] == heading {
					dirs[0], dirs[i] = dirs[i], dirs[0]
					break
				}
			}
		}
		advanced := false
		for _, d := range dirs {
			n := Vec{cur.X + d.X, cur.Y + d.Y}
			if n.X < 1 || n.Y < 1 || n.X >= w-1 || n.Y >= h-1 {
				continue
			}
			if open[n.Y*w+n.X] {
				continue
			}
			mid := Vec{cur.X + d.X/2, cur.Y + d.Y/2}
			open[mid.Y*w+mid.X] = true
			open[n.Y*w+n.X] = true
			stack = append(stack, n)
			headings = append(headings, d)
			advanced = true
			break
		}
		if !advanced {
			stack = stack[:len(stack)-1]
			headings = headings[:len(headings)-1]
		}
	}
	y0 := (h / 2) | 1
	if y0 == 0 {
		y0 = 1
	}
	y1 := ((h / 2) + 3) | 1
	if y1 > h-2 {
		y1 = h - 2
	}
	open[y0*w+0] = true
	open[y1*w+w-1] = true

	entrance := Vec{0, y0}
	exit := Vec{w - 1, y1}
	prev := map[Vec]Vec{}
	seen := map[Vec]bool{entrance: true}
	queue := []Vec{entrance}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == exit {
			break
		}
		for _, d := range [4]Vec{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := Vec{cur.X + d.X, cur.Y + d.Y}
			if n.X < 0 || n.Y < 0 || n.X >= w || n.Y >= h || seen[n] || !open[n.Y*w+n.X] {
				continue
			}
			seen[n] = true
			prev[n] = cur
			queue = append(queue, n)
		}
	}
	if !seen[exit] {
		return nil, errNoRoute
	}
	pathSet := map[Vec]bool{}
	for cur := exit; cur != entrance; cur = prev[cur] {
		pathSet[cur] = true
	}
	pathSet[entrance] = true
	pathSet[exit] = true

	rows := make([]string, h)
	for y := 0; y < h; y++ {
		b := make([]byte, w)
		for x := 0; x < w; x++ {
			v := Vec{x, y}
			if !open[y*w+x] {
				b[x] = '#'
				continue
			}
			if v == entrance {
				b[x] = 'S'
				continue
			}
			if v == exit {
				b[x] = 'E'
				continue
			}
			if pathSet[v] {
				b[x] = '-'
				continue
			}
			b[x] = '.'
		}
		rows[y] = string(b)
	}
	m, err := LoadMap(w, h, rows)
	if err != nil {
		return nil, err
	}
	m.HPMul = mazeHP
	return m, nil
}
