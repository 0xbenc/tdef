package game

import (
	"math"
	"testing"
)

func TestSpecialistPlacementLimit(t *testing.T) {
	s := NewState(loadTestMap(t))
	s.Gold = 100000
	var cells []Vec
	for y := 1; y < s.Map.H-1; y++ {
		for x := 1; x < s.Map.W-1; x++ {
			v := Vec{x, y}
			if s.Map.At(v) == CellGrass {
				cells = append(cells, v)
			}
		}
	}
	for k := TowerRuneforge; k < TowerCount; k++ {
		first := s.Build(cells[0], k)
		if first == nil {
			t.Fatal("first specialist rejected")
		}
		if s.CanBuild(cells[1], k) || s.Build(cells[1], k) != nil {
			t.Fatal("duplicate specialist allowed")
		}
		if !s.Upgrade(first) {
			t.Fatal("limit blocked upgrading")
		}
		s.Sell(first)
		replacement := s.Build(cells[1], k)
		if replacement == nil {
			t.Fatal("selling did not release specialist slot")
		}
		s.Sell(replacement)
	}
	if s.Build(cells[0], TowerGunner) == nil || s.Build(cells[1], TowerGunner) == nil {
		t.Fatal("original defenders incorrectly limited")
	}
	fresh := NewState(s.Map)
	for k := TowerRuneforge; k < TowerCount; k++ {
		if fresh.SpecialistPlaced(k) {
			t.Fatal("specialist limit carried to next level")
		}
	}
}

func specialistEnemy(s *State, k EnemyKind, progress float64) *Enemy {
	e := &Enemy{ID: s.NextID, Kind: k, HP: 10000, MaxHP: 10000, Speed: 1, SlowFactor: 1, Prog: progress, PrevProg: progress, Pos: s.Map.PointAt(progress)}
	s.NextID++
	s.Enemies = append(s.Enemies, e)
	return e
}

func TestForgeContinuousDamageAndRotation(t *testing.T) {
	for _, dt := range []float64{.01, .05, .1} {
		s := NewState(loadTestMap(t))
		tower := &Tower{Kind: TowerRuneforge, Level: 1, Cell: Vec{2, 2}}
		s.Aim(tower, FacingEast)
		for _, k := range []EnemyKind{EnemyGrunt, EnemyTank, EnemyShield} {
			e := specialistEnemy(s, k, 1)
			e.Pos = Pos{4.5, 2.5}
		}
		off := specialistEnemy(s, EnemyGrunt, 1)
		off.Pos = Pos{4.5, 3.5}
		for elapsed := 0.; elapsed < 1-dt/2; elapsed += dt {
			s.fireSpecialist(tower, dt)
		}
		for i, e := range s.Enemies {
			want := 22.
			if i == 2 {
				want *= .2
			}
			if i == 3 {
				want = 0
			}
			if math.Abs(10000-e.HP-want) > 1e-7 {
				t.Fatalf("dt %v enemy %d damage %v want %v", dt, i, 10000-e.HP, want)
			}
		}
		s.Aim(tower, FacingNorth)
		s.fireSpecialist(tower, .1)
		if tower.Active {
			t.Fatal("rotating did not move the damaging ray")
		}
	}
}

func TestHooksFollowRoadAndCannotTrapForever(t *testing.T) {
	for _, kind := range []EnemyKind{EnemyGrunt, EnemyTank, EnemyShield, EnemyBoss} {
		s := NewState(loadTestMap(t))
		e := specialistEnemy(s, kind, 10)
		hook := &Tower{Kind: TowerHookmaster, Level: 3, Cell: s.Map.Path[9]}
		hook2 := *hook
		s.fireSpecialist(hook, .1)
		if e.Prog >= 10 || e.Pos != s.Map.PointAt(e.Prog) {
			t.Fatal("hook failed to follow path")
		}
		progress := e.Prog
		s.fireSpecialist(&hook2, .1)
		if e.Prog != progress {
			t.Fatal("second hook ignored shared recovery")
		}
		for i := 0; i < 20; i++ {
			s.Time += 10
			hook.CD = 0
			s.fireSpecialist(hook, .1)
		}
		if e.PulledDistance > HookBudget(kind)+1e-9 {
			t.Fatal("exceeded lifetime pull allowance")
		}
		before := e.Prog
		s.Time += 10
		hook.CD = 0
		s.fireSpecialist(hook, .1)
		if before != e.Prog {
			t.Fatal("exhausted allowance still permits pulling")
		}
	}
}

func TestHexDoesNotStackAndExpiresAfterResistance(t *testing.T) {
	s := NewState(loadTestMap(t))
	e := specialistEnemy(s, EnemyTank, 5)
	witch := &Tower{Kind: TowerWitch, Level: 3, Cell: s.Map.Path[4], TargetMode: TargetStrongest}
	s.fireSpecialist(witch, .1)
	strongUntil := e.HexUntil
	weak := *witch
	weak.Level = 1
	weak.CD = 0
	s.fireSpecialist(&weak, .1)
	if e.HexBonus != .45 || e.HexUntil != strongUntil {
		t.Fatal("weaker witch replaced stronger hex")
	}
	s.applyDamage(e, 100, TowerGunner)
	if math.Abs(10000-e.HP-36.25) > 1e-9 {
		t.Fatal("hex did not amplify resisted damage")
	}
	s.Time = e.HexUntil
	before := e.HP
	s.applyDamage(e, 100, TowerGunner)
	if before-e.HP != 25 {
		t.Fatal("hex did not expire")
	}
}

func TestMinesArmCapSweepAndSale(t *testing.T) {
	s := NewState(loadTestMap(t))
	s.Gold = 10000
	var crew *Tower
	for y := 1; y < s.Map.H-1 && crew == nil; y++ {
		for x := 1; x < s.Map.W-1; x++ {
			v := Vec{x, y}
			if s.Map.At(v) == CellGrass {
				candidate := &Tower{Kind: TowerSappers, Level: 1, Cell: v}
				if len(s.MineSites(candidate)) >= 6 {
					crew = s.Build(v, TowerSappers)
					break
				}
			}
		}
	}
	if crew == nil {
		t.Fatal("no mine placement")
	}
	for i := 0; i < 20; i++ {
		s.plantMine(crew)
	}
	if len(s.Mines) == 0 || len(s.Mines) > 3 {
		t.Fatal("mine cap or planting failed")
	}
	mine := s.Mines[0]
	e := specialistEnemy(s, EnemyRunner, mine.Prog)
	s.triggerMines()
	if e.HP != 10000 {
		t.Fatal("mine detonated before arming")
	}
	s.Time = MineArmDelay
	e.PrevProg = mine.Prog - 2
	e.Prog = mine.Prog + 2
	e.Pos = s.Map.PointAt(e.Prog)
	s.triggerMines()
	if e.HP > 9910 {
		t.Fatal("fast rogue crossed mine without full splash damage")
	}
	s.Sell(crew)
	if len(s.Mines) != 0 {
		t.Fatal("selling crew left mines")
	}
}

func TestMineCanStopLastTickLeak(t *testing.T) {
	s := NewState(loadTestMap(t))
	e := specialistEnemy(s, EnemyGrunt, s.Map.TotalLen-.6)
	e.HP = 1
	e.Speed = 20
	s.Mines = []*Mine{{Prog: s.Map.TotalLen - .5, Pos: s.Map.PointAt(s.Map.TotalLen - .5), Damage: 90, Radius: 1.1}}
	s.moveEnemies(.1)
	s.triggerMines()
	s.finishLeaks()
	if !e.Dead || e.Leaked || s.TotalLeaks != 0 {
		t.Fatal("mine crossing lost to premature leak")
	}
}
