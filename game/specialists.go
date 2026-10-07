package game

import "math"

const (
	ForgeHalfWidth = .48
	HookRecovery   = 2.5
	MineArmDelay   = .6
)

// ForgeEnd is shared by the simulation and the placement preview. Walls do
// not block attacks in TERMTD; the ray ends at its range or the map boundary.
func ForgeEnd(m *Map, cell Vec, facing Facing, reach float64) Pos {
	p := cell.Center()
	d := facing.Vector()
	return Pos{math.Max(.5, math.Min(float64(m.W)-.5, p.X+float64(d.X)*reach)),
		math.Max(.5, math.Min(float64(m.H)-.5, p.Y+float64(d.Y)*reach))}
}

func (s *State) Aim(t *Tower, facing Facing) {
	t.mineSitesReady = false
	t.Facing = facing
	t.BeamEnd = ForgeEnd(s.Map, t.Cell, facing, t.Range())
}

func ForgeContains(start, end, p Pos) bool {
	dx, dy := end.X-start.X, end.Y-start.Y
	length := math.Hypot(dx, dy)
	if length < .01 {
		return false
	}
	along := ((p.X-start.X)*dx + (p.Y-start.Y)*dy) / length
	across := math.Abs((p.X-start.X)*dy-(p.Y-start.Y)*dx) / length
	return along > .25 && along <= length && across <= ForgeHalfWidth
}

func HookBudget(k EnemyKind) float64 {
	switch k {
	case EnemyTank, EnemyShield:
		return 2
	case EnemyBoss:
		return .5
	}
	return 6
}

func (s *State) fireSpecialist(t *Tower, dt float64) bool {
	switch t.Kind {
	case TowerRuneforge:
		t.Active = false
		for _, e := range s.Enemies {
			if !e.Dead && !e.Leaked && ForgeContains(t.Pos(), t.BeamEnd, e.Pos) {
				t.Active = true
				s.applyDamage(e, t.Dmg()*dt, t.Kind)
			}
		}
		return true
	case TowerHookmaster, TowerWitch, TowerSappers:
		t.CD = math.Max(0, t.CD-dt)
		t.Flash = math.Max(0, t.Flash-dt)
		if t.CD > 0 {
			return true
		}
		if t.Kind == TowerSappers {
			if s.plantMine(t) {
				t.CD = 1 / t.RoT()
			}
			return true
		}
		e := s.acquireTarget(t)
		if e == nil {
			return true
		}
		t.CD, t.Flash, t.FlashTo = 1/t.RoT(), .2, e.Pos
		if t.Kind == TowerWitch {
			bonus := t.Spec().Mark[t.Level-1]
			if e.HexUntil <= s.Time || bonus >= e.HexBonus {
				e.HexBonus = bonus
				e.HexUntil = s.Time + t.Spec().MarkDur[t.Level-1]
			}
			s.Beams = append(s.Beams, &Beam{From: []Pos{t.Pos(), e.Pos}, TTL: .18, Max: .18, Kind: t.Kind})
		} else {
			old := e.Pos
			pull := t.Spec().Pull[t.Level-1]
			switch e.Kind {
			case EnemyTank, EnemyShield:
				pull *= .4
			case EnemyBoss:
				pull *= .1
			}
			pull = math.Min(pull, HookBudget(e.Kind)-e.PulledDistance)
			pull = math.Min(pull, e.Prog)
			e.PulledDistance += pull
			e.Prog = math.Max(0, e.Prog-pull)
			e.Pos = s.Map.PointAt(e.Prog)
			e.HookUntil = s.Time + HookRecovery
			s.applyDamage(e, t.Dmg(), t.Kind)
			s.Beams = append(s.Beams, &Beam{From: []Pos{t.Pos(), old}, TTL: .22, Max: .22, Kind: t.Kind})
		}
		return true
	}
	return false
}

// MineSites are actual route positions, rather than nearby floor decoration.
// Exposing the candidates also lets selection previews and builders use the
// same legal placement geometry as the simulation.
func (s *State) MineSites(t *Tower) []Mine {
	if t.mineSitesReady {
		return t.mineSites
	}
	var sites []Mine
	for i, cell := range s.Map.Path {
		if i == 0 || i == len(s.Map.Path)-1 || t.Pos().Dist(cell.Center()) > t.Range() {
			continue
		}
		sites = append(sites, Mine{Prog: s.Map.Dist[i], Pos: cell.Center()})
	}
	t.mineSites, t.mineSitesReady = sites, true
	return sites
}

func (s *State) plantMine(t *Tower) bool {
	count := 0
	for _, mine := range s.Mines {
		if mine.Owner == t.ID {
			count++
		}
	}
	if count >= t.Spec().MineCap[t.Level-1] {
		return false
	}
	var best *Mine
	bestScore := -math.MaxFloat64
	for _, site := range s.MineSites(t) {
		spacing := 4.0
		for _, mine := range s.Mines {
			spacing = math.Min(spacing, math.Abs(mine.Prog-site.Prog))
		}
		if spacing < 1.5 {
			continue
		}
		// Spread charges, favoring a short walk from the crew.
		score := spacing - t.Pos().Dist(site.Pos)*.15
		if score > bestScore {
			copy := site
			best, bestScore = &copy, score
		}
	}
	if best == nil {
		return false
	}
	best.Owner, best.ArmAt = t.ID, s.Time+MineArmDelay
	best.Damage, best.Radius = t.Dmg(), t.Spec().Splash
	s.Mines = append(s.Mines, best)
	t.Flash, t.FlashTo = .2, best.Pos
	return true
}

func (s *State) triggerMines() {
	n := 0
	for _, mine := range s.Mines {
		trigger := false
		if mine.ArmAt <= s.Time {
			for _, e := range s.Enemies {
				lo, hi := math.Min(e.PrevProg, e.Prog), math.Max(e.PrevProg, e.Prog)
				if !e.Dead && !e.Leaked && lo <= mine.Prog+.35 && hi >= mine.Prog-.35 {
					trigger = true
					break
				}
			}
		}
		if !trigger {
			s.Mines[n] = mine
			n++
			continue
		}
		for _, e := range s.Enemies {
			if !e.Dead && !e.Leaked && (e.Pos.Dist(mine.Pos) <= mine.Radius || (e.PrevProg <= mine.Prog && e.Prog >= mine.Prog)) {
				s.applySplashDamage(e, mine.Damage, TowerSappers)
			}
		}
		s.Fx = append(s.Fx, &Fx{Pos: mine.Pos, TTL: .3, Max: .3, Color: 214, Ring: mine.Radius})
	}
	s.Mines = s.Mines[:n]
}
