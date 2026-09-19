package game

import (
	"strings"
	"testing"
)

const testMapText = `
##########
#S--.#####
##.-##...#
#---##.###
#-..##.###
#----#.###
#...----E#
##########
`

func loadTestMap(t *testing.T) *Map {
	t.Helper()
	rows := strings.Split(strings.TrimSpace(testMapText), "\n")
	m, err := LoadMap(10, len(rows), rows)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestMapPath(t *testing.T) {
	m := loadTestMap(t)
	if m.Spawn != (Vec{1, 1}) {
		t.Errorf("spawn = %v", m.Spawn)
	}
	if m.Exit != (Vec{8, 6}) {
		t.Errorf("exit = %v", m.Exit)
	}
	if len(m.Path) < 10 {
		t.Errorf("path too short: %d", len(m.Path))
	}
	prev := m.Path[0]
	for _, v := range m.Path[1:] {
		if v.Man(prev) != 1 {
			t.Fatalf("path not contiguous at %v -> %v", prev, v)
		}
		prev = v
	}
	for i := 1; i < len(m.Path); i++ {
		if m.Dist[i] <= m.Dist[i-1] {
			t.Fatalf("dist not monotonic at %d", i)
		}
	}
	p0 := m.PointAt(0)
	if p0 != m.Path[0].Center() {
		t.Errorf("PointAt(0) = %v", p0)
	}
	pEnd := m.PointAt(m.TotalLen)
	if pEnd != m.Path[len(m.Path)-1].Center() {
		t.Errorf("PointAt(end) = %v", pEnd)
	}
	mid := m.PointAt(m.TotalLen / 2)
	if !m.InBounds(mid.Center()) {
		t.Errorf("PointAt(mid) out of bounds: %v", mid)
	}
}

func TestGenerateMap(t *testing.T) {
	for seed := int64(1); seed <= 25; seed++ {
		m, err := GenerateMap(seed, 45, 13)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		if len(m.Path) < 20 {
			t.Fatalf("seed %d: path too short %d", seed, len(m.Path))
		}
	}
}

func TestBuildAndSell(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	cost := TowerSpecs[TowerGunner].Cost[0]
	goldBefore := s.Gold
	if !s.CanBuild(Vec{4, 1}, TowerGunner) {
		t.Fatal("should be able to build on grass")
	}
	tw := s.Build(Vec{4, 1}, TowerGunner)
	if tw == nil {
		t.Fatal("build failed")
	}
	if s.Gold != goldBefore-cost {
		t.Errorf("gold after build = %d", s.Gold)
	}
	if s.CanBuild(Vec{4, 1}, TowerGunner) {
		t.Error("cell occupied")
	}
	if s.CanBuild(Vec{1, 1}, TowerGunner) {
		t.Error("cannot build on spawn path")
	}
	refund := s.Sell(tw)
	if refund != int(float64(cost)*SellRefund) {
		t.Errorf("refund = %d, want %d", refund, int(float64(cost)*SellRefund))
	}
	if s.Gold != goldBefore-cost+refund {
		t.Errorf("gold after sell = %d, want %d", s.Gold, goldBefore-cost+refund)
	}
}

func TestUpgrade(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	tw := s.Build(Vec{4, 1}, TowerGunner)
	c := TowerSpecs[TowerGunner].Cost[1]
	gold := s.Gold
	if !s.Upgrade(tw) {
		t.Fatal("upgrade failed")
	}
	if tw.Level != 2 || s.Gold != gold-c {
		t.Errorf("level=%d gold=%d", tw.Level, s.Gold)
	}
}

func TestEnemyLeak(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	s.StartWave()
	livesBefore := s.Lives
	s.Enemies = append(s.Enemies, &Enemy{
		ID: 999, Kind: EnemyGrunt, HP: 10, MaxHP: 10,
		Prog: m.TotalLen - 0.1, Speed: 50, SlowFactor: 1,
		Pos: m.PointAt(m.TotalLen - 0.1), Lives: 1,
	})
	s.Step(1.0)
	if s.Lives != livesBefore-1 {
		t.Errorf("lives after leak = %d, want %d", s.Lives, livesBefore-1)
	}
}

func TestEnemyDeathAndGold(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	s.Gold = 0
	goldBefore := s.Gold
	e := &Enemy{ID: 999, Kind: EnemyGrunt, HP: 10, MaxHP: 10, Speed: 1, SlowFactor: 1, Lives: 1,
		Bounty: EnemySpecs[EnemyGrunt].Bounty}
	s.Enemies = append(s.Enemies, e)
	if !e.Damage(10) {
		t.Fatal("damage should kill")
	}
	if !e.Dead {
		t.Fatal("enemy not dead")
	}
	s.killEnemy(e)
	if s.Gold != goldBefore+EnemySpecs[EnemyGrunt].Bounty {
		t.Errorf("gold = %d", s.Gold)
	}
}

func TestTargetModes(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	tw := &Tower{ID: 1, Kind: TowerGunner, Level: 3, Cell: Vec{X: 5, Y: 2}}
	s.Towers = append(s.Towers, tw)
	eA := &Enemy{ID: 1, Kind: EnemyGrunt, HP: 500, MaxHP: 500, Speed: 1, SlowFactor: 1, Lives: 1, Pos: Pos{5.5, 3.5}, Prog: 10}
	eB := &Enemy{ID: 2, Kind: EnemyGrunt, HP: 100, MaxHP: 100, Speed: 1, SlowFactor: 1, Lives: 1, Pos: Pos{7.5, 2.5}, Prog: 30}
	eC := &Enemy{ID: 3, Kind: EnemyGrunt, HP: 200, MaxHP: 200, Speed: 1, SlowFactor: 1, Lives: 1, Pos: Pos{5.5, 4.5}, Prog: 20}
	s.Enemies = append(s.Enemies, eA, eB, eC)
	check := func(mode TargetMode, want int) {
		tw.TargetMode = mode
		got := s.acquireTarget(tw)
		if got == nil || got.ID != want {
			t.Errorf("%v: got %v want id %d", mode.Name(), got, want)
		}
	}
	check(TargetFirst, eB.ID)
	check(TargetStrongest, eA.ID)
	check(TargetClosest, eA.ID)
	tw.TargetMode = TargetFirst
	if _, m2 := s.CycleTarget(tw.Cell); m2 != TargetStrongest {
		t.Errorf("cycle: got %v want strongest", m2)
	}
}

func TestArmorReducesDamage(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	e := &Enemy{ID: 1, Kind: EnemyShield, HP: 100, MaxHP: 100, Speed: 1, SlowFactor: 1, Lives: 1, Armor: 0.4}
	s.Enemies = append(s.Enemies, e)
	s.applyDamage(e, 50, TowerGunner)
	// 50 damage * (1-0.4) = 30
	if e.HP != 70 {
		t.Errorf("HP after armored hit = %v, want 70", e.HP)
	}
}

func TestNewTowersBuild(t *testing.T) {
	m := loadTestMap(t)
	cells := findGrassCells(t, m, 2)
	for i, k := range []TowerKind{TowerMortar, TowerFlak} {
		s := NewState(m, 1)
		s.Gold = 1000
		if s.Build(cells[i], k) == nil {
			t.Fatalf("%s should build", TowerSpecs[k].Name)
		}
	}
}

func findGrassCells(t *testing.T, m *Map, n int) []Vec {
	t.Helper()
	var out []Vec
	for y := 0; y < m.H && len(out) < n; y++ {
		for x := 0; x < m.W && len(out) < n; x++ {
			v := Vec{X: x, Y: y}
			if m.At(v) == CellGrass {
				out = append(out, v)
			}
		}
	}
	if len(out) < n {
		t.Fatalf("found %d grass cells, want %d", len(out), n)
	}
	return out
}

func TestFrostSlows(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	s.Gold = 1000
	e := &Enemy{ID: 42, Kind: EnemyGrunt, HP: 1000, MaxHP: 1000, Speed: 2, SlowFactor: 1, Lives: 1}
	s.Enemies = append(s.Enemies, e)
	s.applyDamage(e, 1, TowerFrost)
	if e.SlowFactor >= 1.0 {
		t.Fatalf("slow not applied: %v", e.SlowFactor)
	}
	if e.SlowUntil <= s.Time {
		t.Fatal("slow until not set")
	}
}

func TestWaveCompletion(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	s.StartWave()
	if !s.WaveActive {
		t.Fatal("wave should be active")
	}
	bonus := WaveBonus(1)
	s.Enemies = nil
	s.SpawnIdx = len(s.SpawnQueue)
	s.cleanup()
	if s.WaveActive {
		t.Fatal("wave should be complete")
	}
	if s.Gold != StartingGold+bonus {
		t.Errorf("gold = %d, want %d", s.Gold, StartingGold+bonus)
	}
}

func TestFinalWaveLeakIsDefeat(t *testing.T) {
	m := loadTestMap(t)
	s := NewState(m, 1)
	s.Wave = MaxWaves - 1
	s.StartWave()
	if s.Wave != MaxWaves {
		t.Fatalf("wave = %d", s.Wave)
	}
	s.Lives = 1
	s.SpawnIdx = len(s.SpawnQueue)
	s.Enemies = []*Enemy{{
		ID: 1, Kind: EnemyGrunt, HP: 10, MaxHP: 10,
		Prog: m.TotalLen - 0.05, Speed: 100, SlowFactor: 1,
		Pos: m.PointAt(m.TotalLen - 0.05), Lives: 1,
	}}
	s.Step(1.0)
	if s.Status != StatusDefeat {
		t.Fatalf("status = %v, want defeat", s.Status)
	}
}

func waveHas(entries []SpawnEntry, k EnemyKind) bool {
	for _, e := range entries {
		if e.Kind == k {
			return true
		}
	}
	return false
}

func TestWaveTableStructure(t *testing.T) {
	for w := 1; w <= MaxWaves; w++ {
		entries := BuildWave(w)
		if len(entries) == 0 {
			t.Fatalf("wave %d has no entries", w)
		}
		var wantBosses int
		switch w {
		case 15, 18:
			wantBosses = 1
		case 20:
			wantBosses = 2
		}
		bosses := 0
		for _, e := range entries {
			if e.Kind == EnemyBoss {
				bosses++
			}
		}
		if bosses != wantBosses {
			t.Errorf("wave %d bosses = %d, want %d", w, bosses, wantBosses)
		}
	}
	if !waveHas(BuildWave(3), EnemyRunner) {
		t.Error("Runner should debut on wave 3")
	}
	if !waveHas(BuildWave(11), EnemyWisp) {
		t.Error("Wisp should debut on wave 11")
	}
	if !waveHas(BuildWave(13), EnemyTank) {
		t.Error("Tank should debut on wave 13")
	}
	if !waveHas(BuildWave(14), EnemyShield) {
		t.Error("Shield should debut on wave 14")
	}
	if waveHas(BuildWave(5), EnemyBoss) || waveHas(BuildWave(10), EnemyBoss) {
		t.Error("no boss before wave 15")
	}
}

func TestAutoWaveDelayTaper(t *testing.T) {
	first, last := AutoWaveDelayFor(1), AutoWaveDelayFor(19)
	if first <= last {
		t.Errorf("delay should taper: first=%v last=%v", first, last)
	}
	if first > 11 || first < 10 {
		t.Errorf("first delay = %v, want ~11", first)
	}
	if last < 4 || last > 5 {
		t.Errorf("last delay = %v, want ~4-5", last)
	}
}

func TestEconomyFormulas(t *testing.T) {
	if WaveBonus(1) != 35 || WaveBonus(20) != 130 {
		t.Errorf("WaveBonus = %d..%d, want 35..130", WaveBonus(1), WaveBonus(20))
	}
	if EarlyBonus(0) != EarlyBonusBase || EarlyBonus(19) != EarlyBonusBase+9 {
		t.Errorf("EarlyBonus = %d..%d, want %d..%d", EarlyBonus(0), EarlyBonus(19), EarlyBonusBase, EarlyBonusBase+9)
	}
}

func TestBossLeakLifeCost(t *testing.T) {
	if EnemySpecs[EnemyBoss].Lives != 6 {
		t.Errorf("boss leak cost = %d, want 6", EnemySpecs[EnemyBoss].Lives)
	}
	m := loadTestMap(t)
	s := NewState(m, 1)
	s.Lives = 20
	s.Enemies = []*Enemy{{
		ID: 1, Kind: EnemyBoss, HP: 10, MaxHP: 10,
		Prog: m.TotalLen - 0.05, Speed: 100, SlowFactor: 1,
		Pos: m.PointAt(m.TotalLen - 0.05), Lives: 6,
	}}
	s.Step(1.0)
	if s.Lives != 14 {
		t.Errorf("lives after boss leak = %d, want 14", s.Lives)
	}
}

func TestWaveThemeAndTelegraph(t *testing.T) {
	if WaveTheme(1) != "warmup" {
		t.Errorf("WaveTheme(1) = %q, want warmup", WaveTheme(1))
	}
	if WaveTelegraph(3) == "" || WaveTelegraph(15) == "" || WaveTelegraph(20) == "" {
		t.Error("expected telegraphs on waves 3, 15, 20")
	}
	if WaveTelegraph(1) != "" {
		t.Error("wave 1 should have no telegraph")
	}
}

func TestDeterminism(t *testing.T) {
	m, err := LoadLevel("winding")
	if err != nil {
		t.Fatal(err)
	}
	run := func(seed int64) (int, int, bool) {
		r := RunAutoplay(m, "winding", seed)
		return r.Kills, r.Lives, r.Won
	}
	k1, l1, w1 := run(777)
	k2, l2, w2 := run(777)
	if k1 != k2 || l1 != l2 || w1 != w2 {
		t.Fatalf("nondeterministic: (%d,%d,%v) vs (%d,%d,%v)", k1, l1, w1, k2, l2, w2)
	}
	k3, _, _ := run(778)
	if k3 == k1 {
		t.Log("warning: different seeds gave identical results (possible but suspicious)")
	}
}

func TestFullGame(t *testing.T) {
	for _, level := range []string{"winding", "garden", "canyon"} {
		m, err := LoadLevel(level)
		if err != nil {
			t.Fatal(err)
		}
		r := RunAutoplay(m, level, 42)
		if r.Wave < 1 {
			t.Fatalf("%s: game never started a wave", level)
		}
		if r.Time <= 0 {
			t.Fatalf("%s: no time advanced", level)
		}
	}
}
