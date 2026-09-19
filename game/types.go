package game

import "math"

type Vec struct {
	X, Y int
}

func (a Vec) Sub(b Vec) Vec { return Vec{a.X - b.X, a.Y - b.Y} }

func (a Vec) Man(b Vec) int {
	dx, dy := a.X-b.X, a.Y-b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

type Pos struct {
	X, Y float64
}

func (p Pos) Center() Vec { return Vec{int(math.Floor(p.X + 0.5)), int(math.Floor(p.Y + 0.5))} }

func (p Pos) Dist(q Pos) float64 {
	dx, dy := p.X-q.X, p.Y-q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (c Vec) Center() Pos { return Pos{float64(c.X) + 0.5, float64(c.Y) + 0.5} }

type TowerKind int

const (
	TowerGunner TowerKind = iota
	TowerCannon
	TowerFrost
	TowerSniper
	TowerTesla
	TowerMortar
	TowerFlak
	TowerCount
)

func (k TowerKind) Valid() bool { return k >= 0 && k < TowerCount }

// TargetMode controls which in-range enemy a tower attacks.
type TargetMode int

const (
	TargetFirst     TargetMode = iota // furthest along the path
	TargetStrongest                   // highest current HP
	TargetClosest                     // nearest to the tower
	TargetModeCount
)

func (m TargetMode) Next() TargetMode {
	if m+1 >= TargetModeCount {
		return TargetFirst
	}
	return m + 1
}

func (m TargetMode) Short() string {
	switch m {
	case TargetFirst:
		return "F"
	case TargetStrongest:
		return "S"
	case TargetClosest:
		return "C"
	}
	return "?"
}

func (m TargetMode) Name() string {
	switch m {
	case TargetFirst:
		return "first"
	case TargetStrongest:
		return "strongest"
	case TargetClosest:
		return "closest"
	}
	return "?"
}

type EnemyKind int

const (
	EnemyMinion EnemyKind = iota
	EnemyRunner
	EnemyGrunt
	EnemyTank
	EnemySplitter
	EnemyBoss
	EnemyWisp
	EnemyShield
	EnemyCount
)

type GameStatus int

const (
	StatusRunning GameStatus = iota
	StatusVictory
	StatusDefeat
)

type CellKind int

const (
	CellWall CellKind = iota
	CellGrass
	CellPath
)
