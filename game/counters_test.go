package game

import (
	"math"
	"strings"
	"testing"
)

func TestGuildCounterDamage(t *testing.T) {
	cases := []struct {
		kind EnemyKind
		want [TowerCount]float64 // gun, cannon, frost, ranger, lightning, treb, sling
	}{
		{EnemyTank, [TowerCount]float64{25, 100, 100, 60, 100, 100, 25, 100, 100, 100, 100}},
		{EnemyShield, [TowerCount]float64{100, 100, 20, 100, 20, 100, 100, 20, 100, 100, 100}},
		{EnemyRunner, [TowerCount]float64{50, 100, 100, 50, 100, 100, 50, 100, 100, 100, 100}},
		{EnemyGrunt, [TowerCount]float64{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}},
		{EnemyBoss, [TowerCount]float64{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}},
	}
	for _, tc := range cases {
		for k := TowerKind(0); k < TowerCount; k++ {
			t.Run(EnemySpecs[tc.kind].Name+"/"+TowerSpecs[k].Name, func(t *testing.T) {
				s := NewState(loadTestMap(t))
				e := &Enemy{Kind: tc.kind, HP: 1000, MaxHP: 1000, SlowFactor: 1}
				s.applyDamage(e, 100, k)
				if got := 1000 - e.HP; math.Abs(got-tc.want[k]) > 1e-9 {
					t.Fatalf("damage = %v, want %v", got, tc.want[k])
				}
			})
		}
	}
}

func TestRogueFrostWindowAndSplashImpact(t *testing.T) {
	s := NewState(loadTestMap(t))
	s.Time = 10
	e := &Enemy{ID: 1, Kind: EnemyRunner, HP: 1000, MaxHP: 1000, Speed: 2.6, SlowFactor: 1, Pos: Pos{5, 5}}
	s.Enemies = []*Enemy{e}
	s.applyDamage(e, 0, TowerFrost)
	for _, k := range []TowerKind{TowerGunner, TowerSniper, TowerFlak} {
		before := e.HP
		s.applyDamage(e, 100, k)
		if before-e.HP != 100 {
			t.Fatal("frost did not disable evasion")
		}
	}
	s.Time = e.SlowUntil // Evasion returns at expiry, even before movement cleanup.
	before := e.HP
	s.applyDamage(e, 100, TowerSniper)
	if before-e.HP != 50 {
		t.Fatal("evasion failed to return when frost expired")
	}
	// Exercise real projectile impact, including a rogue caught beside its target.
	other := &Enemy{ID: 2, Kind: EnemyGrunt, HP: 1000, MaxHP: 1000, Pos: Pos{5.5, 5}}
	s.Enemies = append(s.Enemies, other)
	for _, k := range []TowerKind{TowerCannon, TowerMortar} {
		before = e.HP
		s.Projectiles = []*Projectile{{Kind: k, Target: other.ID, Pos: other.Pos, LastPos: other.Pos, Dmg: 100, Splash: TowerSpecs[k].Splash}}
		s.moveProjectiles(.1)
		if before-e.HP != 100 {
			t.Fatal("splash did not bypass the nearby rogue's evasion")
		}
	}
	// The impact flag also bypasses evasion independently of the tower kind.
	before = e.HP
	s.applySplashDamage(e, 100, TowerGunner)
	if before-e.HP != 100 {
		t.Fatal("splash inherited ordinary-shot evasion")
	}
}

func TestWardDoesNotBlockFrostSlow(t *testing.T) {
	s := NewState(loadTestMap(t))
	e := &Enemy{Kind: EnemyShield, HP: 1000, MaxHP: 1000, SlowFactor: 1}
	s.applyDamage(e, 100, TowerFrost)
	if !e.Slowed(s.Time) || e.SlowFactor != TowerSpecs[TowerFrost].SlowPct {
		t.Fatal("ward blocked frost's slow effect")
	}
}

func TestCounterWarningsFitCompactBattlefield(t *testing.T) {
	for _, name := range []string{"hub", "canyon", "winding", "garden", "heart"} {
		m := authoredMap(t, name)
		for _, tc := range []struct {
			wave   int
			answer string
		}{{3, "frost"}, {7, "frost"}, {13, "Cannon"}, {14, "Ranger"}} {
			warning := WaveTelegraphFor(m, tc.wave)
			if !strings.Contains(strings.ToLower(warning), strings.ToLower(tc.answer)) || len([]rune(warning)) > 58 {
				t.Fatalf("%s wave %d lacks a readable counter warning: %q", name, tc.wave, warning)
			}
		}
	}
}
