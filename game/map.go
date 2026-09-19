package game

import (
	"errors"
	"fmt"
)

type Map struct {
	W, H int
	Cell []CellKind

	Path     []Vec
	Dist     []float64
	TotalLen float64
	Spawn    Vec
	Exit     Vec
	PathSet  map[Vec]bool
	Samples  []Pos

	// HPMul is a per-map enemy health multiplier used to tune difficulty.
	HPMul float64
}

const SampleStep = 0.25

func LoadMap(w, h int, rows []string) (*Map, error) {
	if len(rows) != h {
		return nil, fmt.Errorf("map has %d rows, want %d", len(rows), h)
	}
	m := &Map{W: w, H: h, Cell: make([]CellKind, w*h), PathSet: map[Vec]bool{}, HPMul: 1.0}
	for y, row := range rows {
		if len(row) != w {
			return nil, fmt.Errorf("row %d has width %d, want %d", y, len(row), w)
		}
		for x, ch := range row {
			v := Vec{x, y}
			switch ch {
			case '#':
				m.Cell[y*w+x] = CellWall
			case '.':
				m.Cell[y*w+x] = CellGrass
			case '-':
				m.Cell[y*w+x] = CellPath
			case 'S':
				m.Cell[y*w+x] = CellPath
				m.Spawn = v
			case 'E':
				m.Cell[y*w+x] = CellPath
				m.Exit = v
			default:
				return nil, fmt.Errorf("row %d: bad cell %q", y, ch)
			}
		}
	}
	if err := m.computePath(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Map) At(v Vec) CellKind {
	if v.X < 0 || v.Y < 0 || v.X >= m.W || v.Y >= m.H {
		return CellWall
	}
	return m.Cell[v.Y*m.W+v.X]
}

func (m *Map) InBounds(v Vec) bool {
	return v.X >= 0 && v.Y >= 0 && v.X < m.W && v.Y < m.H
}

func (m *Map) IsPath(v Vec) bool { return m.PathSet[v] }

func (m *Map) computePath() error {
	if !m.InBounds(m.Spawn) || !m.InBounds(m.Exit) || m.At(m.Spawn) != CellPath || m.At(m.Exit) != CellPath {
		return errors.New("map needs an S (spawn) and an E (exit) on path")
	}
	prev := map[Vec]Vec{}
	seen := map[Vec]bool{m.Spawn: true}
	queue := []Vec{m.Spawn}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == m.Exit {
			break
		}
		for _, d := range [4]Vec{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := Vec{cur.X + d.X, cur.Y + d.Y}
			if seen[n] || m.At(n) != CellPath {
				continue
			}
			seen[n] = true
			prev[n] = cur
			queue = append(queue, n)
		}
	}
	if !seen[m.Exit] {
		return errors.New("no path from S to E")
	}
	path := []Vec{}
	for cur := m.Exit; cur != m.Spawn; cur = prev[cur] {
		path = append(path, cur)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	m.Path = path
	m.Dist = make([]float64, len(path))
	for i, v := range path {
		m.PathSet[v] = true
		if i > 0 {
			m.Dist[i] = m.Dist[i-1] + float64(path[i].Man(path[i-1]))
		}
	}
	m.TotalLen = m.Dist[len(path)-1]
	m.Samples = m.SamplePoints()
	return nil
}

func (m *Map) SamplePoints() []Pos {
	pts := []Pos{}
	for i := 1; i < len(m.Path); i++ {
		a, b := m.Dist[i-1], m.Dist[i]
		for d := a + SampleStep; d < b; d += SampleStep {
			pts = append(pts, m.PointAt(d))
		}
	}
	return pts
}

func (m *Map) PointAt(d float64) Pos {
	if d <= 0 {
		return m.Path[0].Center()
	}
	if d >= m.TotalLen {
		return m.Path[len(m.Path)-1].Center()
	}
	lo, hi := 0, len(m.Dist)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if m.Dist[mid] <= d {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	a, b := m.Path[lo], m.Path[lo+1]
	t := (d - m.Dist[lo]) / (m.Dist[lo+1] - m.Dist[lo])
	pa, pb := a.Center(), b.Center()
	return Pos{pa.X + (pb.X-pa.X)*t, pa.Y + (pb.Y-pa.Y)*t}
}
