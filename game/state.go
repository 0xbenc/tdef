package game

import (
	"math/rand/v2"
)

type State struct {
	Map  *Map
	Rng  *rand.Rand
	Time float64
	Diff Difficulty

	Gold   int
	Lives  int
	Wave   int
	Status GameStatus

	Enemies     []*Enemy
	Towers      []*Tower
	Projectiles []*Projectile
	Beams       []*Beam
	Fx          []*Fx
	LeakFlash   float64

	SpawnQueue []SpawnEntry
	SpawnIdx   int
	NextWaveAt float64
	WaveActive bool
	WaveStart  float64
	AutoWave   bool

	Score      int
	TotalKills int
	TotalLeaks int
	Combo      int
	MaxCombo   int

	NextID int

	// UI hints (headless-safe)
	Placing   TowerKind
	PlacingOn bool
	Cursor    Vec
	Selected  int
}

func NewState(m *Map, seed int64, autoWave bool) *State {
	return NewStateDiff(m, seed, autoWave, Normal)
}

func NewStateDiff(m *Map, seed int64, autoWave bool, diff Difficulty) *State {
	return &State{
		Map:        m,
		Rng:        rand.New(rand.NewPCG(uint64(seed), uint64(seed<<32))),
		Gold:       diff.Gold(),
		Lives:      diff.Lives(),
		Diff:       diff,
		AutoWave:   autoWave,
		NextWaveAt: 3.0,
		Placing:    TowerGunner,
		Cursor:     Vec{X: m.W / 2, Y: m.H / 2},
		Selected:   -1,
	}
}

func (s *State) Enemy(id int) *Enemy {
	for _, e := range s.Enemies {
		if e.ID == id {
			return e
		}
	}
	return nil
}

func (s *State) Tower(id int) *Tower {
	for _, t := range s.Towers {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func (s *State) TowerAt(v Vec) *Tower {
	for _, t := range s.Towers {
		if t.Cell == v {
			return t
		}
	}
	return nil
}

func (s *State) CanBuild(v Vec, k TowerKind) bool {
	return s.Map.InBounds(v) && s.Map.At(v) == CellGrass && s.TowerAt(v) == nil && s.Gold >= TowerSpecs[k].Cost[0]
}

func (s *State) Build(v Vec, k TowerKind) *Tower {
	if !s.CanBuild(v, k) {
		return nil
	}
	s.Gold -= TowerSpecs[k].Cost[0]
	t := &Tower{
		ID:       s.NextID,
		Kind:     k,
		Level:    1,
		Cell:     v,
		Invested: TowerSpecs[k].Cost[0],
	}
	s.NextID++
	s.Towers = append(s.Towers, t)
	return t
}

func (s *State) UpgradeCost(t *Tower) int {
	if t.Level >= 3 {
		return 0
	}
	return TowerSpecs[t.Kind].Cost[t.Level]
}

func (s *State) Upgrade(t *Tower) bool {
	c := s.UpgradeCost(t)
	if c == 0 || s.Gold < c {
		return false
	}
	s.Gold -= c
	t.Invested += c
	t.Level++
	return true
}

func (s *State) Sell(t *Tower) int {
	refund := int(float64(t.Invested) * SellRefund)
	s.Gold += refund
	for i, x := range s.Towers {
		if x.ID == t.ID {
			s.Towers = append(s.Towers[:i], s.Towers[i+1:]...)
			break
		}
	}
	return refund
}

// CycleTarget cycles the targeting mode of the tower at v, if any.
// Returns the tower (possibly nil) and its new mode.
func (s *State) CycleTarget(v Vec) (*Tower, TargetMode) {
	t := s.TowerAt(v)
	if t == nil {
		return nil, TargetFirst
	}
	t.TargetMode = t.TargetMode.Next()
	return t, t.TargetMode
}

// StartWave begins the next wave. If it was started before the auto timer
// expired, the player gets an early-start bonus. Returns bonus gold.
func (s *State) StartWave() int {
	if s.WaveActive || s.Status != StatusRunning {
		return 0
	}
	if s.Wave >= MaxWaves {
		return 0
	}
	bonus := 0
	if s.Wave > 0 && s.Time < s.NextWaveAt {
		bonus = EarlyBonusBase + s.Wave
		s.Gold += bonus
		s.Score += bonus * 10
	}
	s.Wave++
	s.SpawnQueue = BuildWave(s.Wave)
	s.SpawnIdx = 0
	s.WaveActive = true
	s.WaveStart = s.Time
	return bonus
}

func (s *State) killEnemy(e *Enemy) {
	s.Gold += e.Bounty
	s.TotalKills++
	s.Combo++
	if s.Combo > s.MaxCombo {
		s.MaxCombo = s.Combo
	}
	mul := 1 + 0.1*float64(s.Combo/5)
	s.Score += int(float64(e.Bounty) * 10 * mul)
	s.Fx = append(s.Fx, &Fx{
		Pos: e.Pos, TTL: 0.3, Max: 0.3,
		R: '✦', Color: killColor(e.Kind),
	})
	spec := &EnemySpecs[e.Kind]
	if spec.SplitN > 0 {
		for i := 0; i < spec.SplitN; i++ {
			s.splitChild(e)
		}
	}
}

func killColor(k EnemyKind) int {
	switch k {
	case EnemyBoss:
		return 204
	case EnemyTank:
		return 180
	case EnemySplitter:
		return 171
	}
	return 220
}

func (s *State) splitChild(parent *Enemy) {
	child := &Enemy{
		ID:         s.NextID,
		Kind:       EnemyMinion,
		HP:         parent.MaxHP * EnemySpecs[EnemySplitter].SplitHP,
		MaxHP:      parent.MaxHP * EnemySpecs[EnemySplitter].SplitHP,
		Prog:       parent.Prog,
		Speed:      EnemySpecs[EnemyMinion].Speed * SpeedScale(s.Wave),
		Pos:        parent.Pos,
		SlowFactor: 1.0,
		Bounty:     EnemySpecs[EnemyMinion].Bounty,
		Lives:      1,
	}
	s.NextID++
	s.Enemies = append(s.Enemies, child)
}

func (s *State) leakEnemy(e *Enemy) {
	s.Lives -= e.Lives
	s.TotalLeaks++
	s.Combo = 0
	s.LeakFlash = 0.6
	if s.Lives <= 0 {
		s.Lives = 0
		s.Status = StatusDefeat
	}
}

func (s *State) applyDamage(e *Enemy, d float64, k TowerKind) {
	if k == TowerFrost {
		sp := TowerSpecs[TowerFrost]
		e.SlowUntil = s.Time + sp.SlowDur
		if sp.SlowPct < e.SlowFactor {
			e.SlowFactor = sp.SlowPct
		}
	}
	if e.Armor > 0 {
		d *= 1 - e.Armor
	}
	if e.Damage(d) {
		s.killEnemy(e)
		return
	}
	e.HitTTL = s.Time + 0.09
}
