package game

type Enemy struct {
	ID             int
	Kind           EnemyKind
	HP             float64
	MaxHP          float64
	Prog           float64
	PrevProg       float64
	Speed          float64
	Pos            Pos
	SlowUntil      float64
	SlowFactor     float64
	Bounty         int
	Lives          int
	Dead           bool
	Leaked         bool
	HitTTL         float64 // >0 while the enemy flashes from a recent hit
	PulledDistance float64
	HookUntil      float64 // shared recovery prevents several Ogres trapping a fighter
	HexUntil       float64
	HexBonus       float64
}

func (e *Enemy) Slowed(time float64) bool { return e.SlowUntil > time }

func (e *Enemy) Damage(d float64) bool {
	if e.Dead {
		return false
	}
	e.HP -= d
	if e.HP <= 0 {
		e.Dead = true
		return true
	}
	return false
}

type Tower struct {
	ID             int
	Kind           TowerKind
	Level          int
	Cell           Vec
	CD             float64
	Invested       int
	Flash          float64
	FlashTo        Pos
	TargetMode     TargetMode
	Facing         Facing
	BeamEnd        Pos
	Active         bool
	mineSites      []Mine
	mineSitesReady bool
}

func (t *Tower) Spec() *TowerSpec { return &TowerSpecs[t.Kind] }

func (t *Tower) Dmg() float64   { return TowerSpecs[t.Kind].Dmg[t.Level-1] }
func (t *Tower) Range() float64 { return TowerSpecs[t.Kind].Range[t.Level-1] }
func (t *Tower) RoT() float64   { return TowerSpecs[t.Kind].RoT[t.Level-1] }
func (t *Tower) Pos() Pos       { return t.Cell.Center() }

type Projectile struct {
	ID      int
	Kind    TowerKind
	Pos     Pos
	Target  int
	LastPos Pos
	Dmg     float64
	Splash  float64
	SlowPct float64
	SlowDur float64
	Dead    bool
}

type Beam struct {
	From []Pos
	TTL  float64
	Max  float64
	Kind TowerKind
}

type Mine struct {
	Owner  int
	Prog   float64
	Pos    Pos
	ArmAt  float64
	Damage float64
	Radius float64
}

type Fx struct {
	Pos   Pos
	TTL   float64
	Max   float64
	R     rune
	Color int
	Ring  float64 // if >0, render as an expanding ring of this radius (map units)
}
