package game

type State struct {
	TrainingStage  int
	UnlockedTowers uint32
	Lesson         TowerKind
	LessonPending  bool
	Map            *Map
	Time           float64
	Diff           Difficulty

	Gold   int
	Lives  int
	Wave   int
	Status GameStatus

	// SeenEnemies also remembers enemies killed on their spawning tick.
	SeenEnemies [EnemyCount]bool
	Enemies     []*Enemy
	Towers      []*Tower
	Projectiles []*Projectile
	Beams       []*Beam
	Mines       []*Mine
	Fx          []*Fx
	LeakFlash   float64

	SpawnQueue []SpawnEntry
	SpawnIdx   int
	NextWaveAt float64
	WaveActive bool
	WaveStart  float64

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

// The simulation is fully deterministic: a run is determined entirely by
// (map, difficulty). There is no in-engine randomness, so no seed is kept
// here — the only seed in the game is the maze generator's, which happens
// before a State exists.
func NewState(m *Map) *State {
	return NewStateDiff(m, Normal)
}

func NewStateDiff(m *Map, diff Difficulty) *State {
	return &State{
		Map:        m,
		Gold:       m.ScaleGold(diff.Gold()),
		Lives:      diff.Lives(),
		Diff:       diff,
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

// TargetFor exposes the same priority used for a defender's next shot.
func (s *State) TargetFor(t *Tower) *Enemy {
	if t == nil || t.Kind == TowerRuneforge || t.Kind == TowerSappers {
		return nil
	}
	return s.acquireTarget(t)
}

// SpecialistPlaced reports the one-per-kind limit for the current level.
// Selling a specialist makes its slot available again.
func (s *State) SpecialistPlaced(k TowerKind) bool {
	if k < TowerRuneforge || !k.Valid() {
		return false
	}
	for _, t := range s.Towers {
		if t.Kind == k {
			return true
		}
	}
	return false
}

func (s *State) CanBuild(v Vec, k TowerKind) bool {
	return s.TowerAvailable(k) && !s.SpecialistPlaced(k) && s.Map.InBounds(v) && s.Map.At(v) == CellGrass && s.TowerAt(v) == nil && s.Gold >= TowerSpecs[k].Cost[0]
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
	if k == TowerWitch {
		t.TargetMode = TargetStrongest
	}
	s.Aim(t, FacingEast)
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
	s.Aim(t, t.Facing)
	return true
}

func (s *State) Sell(t *Tower) int {
	refund := int(float64(t.Invested) * SellRefund)
	s.Gold += refund
	n := 0
	for _, mine := range s.Mines {
		if mine.Owner != t.ID {
			s.Mines[n] = mine
			n++
		}
	}
	s.Mines = s.Mines[:n]
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
//
// The siege itself — the very first wave — will not open until the player has
// committed to the defense: at least one tower must be placed before it may
// start, whether by the auto timer or the early-start key.
func (s *State) StartWave() int {
	if s.LessonPending {
		return 0
	}
	if s.WaveActive || s.Status != StatusRunning {
		return 0
	}
	if s.Wave >= MaxWaves {
		return 0
	}
	if s.Wave == 0 && len(s.Towers) == 0 {
		return 0
	}
	bonus := 0
	if s.Wave > 0 && s.Time < s.NextWaveAt {
		bonus = s.Map.ScaleGold(EarlyBonus(s.Wave))
		s.Gold += bonus
		s.Score += bonus * 10
	}
	s.Wave++
	s.SpawnQueue = BuildWaveFor(s.Map, s.Wave)
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
	s.SeenEnemies[EnemyMinion] = true
	child := &Enemy{
		ID:         s.NextID,
		Kind:       EnemyMinion,
		HP:         parent.MaxHP * EnemySpecs[EnemySplitter].SplitHP,
		MaxHP:      parent.MaxHP * EnemySpecs[EnemySplitter].SplitHP,
		Prog:       parent.Prog,
		Speed:      EnemySpecs[EnemyMinion].Speed * SpeedScale(s.Wave),
		Pos:        parent.Pos,
		SlowFactor: 1.0,
		Bounty:     s.Map.ScaleGold(EnemySpecs[EnemyMinion].Bounty),
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
	s.damageEnemy(e, d, k, false)
}

func (s *State) applySplashDamage(e *Enemy, d float64, k TowerKind) {
	s.damageEnemy(e, d, k, true)
}

func (s *State) damageEnemy(e *Enemy, d float64, k TowerKind, splash bool) {
	if e.Dead || e.Leaked {
		return
	}
	if k == TowerFrost {
		sp := TowerSpecs[TowerFrost]
		e.SlowUntil = s.Time + sp.SlowDur
		if sp.SlowPct < e.SlowFactor {
			e.SlowFactor = sp.SlowPct
		}
	}
	d *= e.DamageMultiplier(k, s.Time, splash)
	if e.HexUntil > s.Time {
		d *= 1 + e.HexBonus
	}
	if e.Damage(d) {
		s.killEnemy(e)
		return
	}
	e.HitTTL = s.Time + 0.09
}

// DamageMultiplier keeps counters deterministic. Evasion represents glancing
// shots, not random misses; frost disables it only for the active slow window.
// Splash is a property of the actual impact, not merely the tower's name.
func (e *Enemy) DamageMultiplier(k TowerKind, time float64, splash bool) float64 {
	sp := EnemySpecs[e.Kind]
	mul := 1.0
	switch k {
	case TowerGunner, TowerFlak:
		mul *= 1 - sp.Plate
	case TowerSniper:
		mul *= 1 - sp.ArrowPlate
	case TowerFrost, TowerTesla, TowerRuneforge:
		mul *= 1 - sp.Ward
	}
	ordinaryShot := k == TowerGunner || k == TowerSniper || k == TowerFlak
	if ordinaryShot && !splash && !e.Slowed(time) {
		mul *= 1 - sp.Evasion
	}
	return mul
}
