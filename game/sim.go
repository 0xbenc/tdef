package game

type SimResult struct {
	Won    bool
	Wave   int
	Lives  int
	Gold   int
	Kills  int
	Leaks  int
	Towers int
	Time   float64
	Seed   int64
	Level  string
}

func RunAutoplay(m *Map, level string, seed int64) SimResult {
	return RunAutoplayDiff(m, level, seed, Normal)
}

func RunAutoplayDiff(m *Map, level string, seed int64, diff Difficulty) SimResult {
	s := NewStateDiff(m, seed, true, diff)
	ai := NewAutoplay(s)
	dt := 1.0 / 20.0
	maxTowers := 0
	for s.Status == StatusRunning {
		ai.tick()
		s.Step(dt)
		if len(s.Towers) > maxTowers {
			maxTowers = len(s.Towers)
		}
	}
	return SimResult{
		Won:    s.Status == StatusVictory,
		Wave:   s.Wave,
		Lives:  s.Lives,
		Gold:   s.Gold,
		Kills:  s.TotalKills,
		Leaks:  s.TotalLeaks,
		Towers: maxTowers,
		Time:   s.Time,
		Seed:   seed,
		Level:  level,
	}
}

type Autoplay struct {
	s       *State
	cover   [TowerCount][][]int
	nextDec float64
}

func NewAutoplay(s *State) *Autoplay {
	a := &Autoplay{s: s, nextDec: 0}
	a.coverage()
	return a
}

func (a *Autoplay) Tick() { a.tick() }

func (a *Autoplay) coverage() {
	s := a.s
	a.cover = [TowerCount][][]int{}
	for k := TowerKind(0); k < TowerCount; k++ {
		r := TowerSpecs[k].Range[0]
		r2 := r * r
		row := [][]int{}
		for y := 0; y < s.Map.H; y++ {
			line := []int{}
			for x := 0; x < s.Map.W; x++ {
				v := Vec{x, y}
				if s.Map.At(v) != CellGrass {
					line = append(line, 0)
					continue
				}
				c := 0
				p := v.Center()
				for _, sp := range s.Map.Samples {
					dx, dy := p.X-sp.X, p.Y-sp.Y
					if dx*dx+dy*dy <= r2 {
						c++
					}
				}
				line = append(line, c)
			}
			row = append(row, line)
		}
		a.cover[k] = row
	}
}

func (a *Autoplay) tick() {
	s := a.s
	if s.Status != StatusRunning || s.Time < a.nextDec {
		return
	}
	a.nextDec = s.Time + 0.5
	a.upgrade()
	a.build()
}

func (a *Autoplay) build() {
	s := a.s
	for {
		bestK, bestX, bestY := TowerGunner, -1, -1
		bestV := 0.0
		for k := TowerKind(0); k < TowerCount; k++ {
			cost := TowerSpecs[k].Cost[0]
			if s.Gold < cost+40 {
				continue
			}
			dps := TowerSpecs[k].Dmg[0] * TowerSpecs[k].RoT[0]
			// eff weights splash/chain towers higher so the AI builds a real
			// mix (a "competent player" proxy), not just cheap single-target
			// Gunners — otherwise it drowns in the dense BTD3-style waves.
			eff := float64(1)
			switch k {
			case TowerCannon:
				eff = 2.5
			case TowerTesla:
				eff = 2.0
			case TowerFrost:
				eff = 1.1
			case TowerMortar:
				eff = 2.8
			case TowerFlak:
				eff = 1.2
			case TowerSniper:
				eff = 1.3
			}
			for y := 0; y < s.Map.H; y++ {
				for x := 0; x < s.Map.W; x++ {
					c := a.cover[k][y][x]
					if c <= 0 {
						continue
					}
					if s.TowerAt(Vec{x, y}) != nil {
						continue
					}
					v := float64(c) * dps * eff / float64(cost)
					if v > bestV {
						bestV = v
						bestK, bestX, bestY = k, x, y
					}
				}
			}
		}
		if bestX < 0 {
			return
		}
		if s.Build(Vec{bestX, bestY}, bestK) == nil {
			return
		}
	}
}

func (a *Autoplay) upgrade() {
	s := a.s
	if s.Wave < 3 {
		return
	}
	for {
		var best *Tower
		bestC := -1
		for _, t := range s.Towers {
			if t.Level >= 3 {
				continue
			}
			c := a.cover[t.Kind][t.Cell.Y][t.Cell.X]
			if c > bestC {
				best = t
				bestC = c
			}
		}
		if best == nil {
			return
		}
		cost := s.UpgradeCost(best)
		if cost == 0 || s.Gold < cost+60 {
			return
		}
		if bestC <= 0 {
			return
		}
		if !s.Upgrade(best) {
			return
		}
	}
}

func Bench(levels []string, games int, seedBase int64) []SimResult {
	out := []SimResult{}
	for i := 0; i < games; i++ {
		level := levels[i%len(levels)]
		m, err := LoadLevel(level)
		if err != nil {
			continue
		}
		out = append(out, RunAutoplay(m, level, seedBase+int64(i)))
	}
	return out
}

func Summarize(rs []SimResult) (wins int, waves, lives, gold, kills int, time float64) {
	for _, r := range rs {
		if r.Won {
			wins++
		}
		waves += r.Wave
		lives += r.Lives
		gold += r.Gold
		kills += r.Kills
		time += r.Time
	}
	n := float64(len(rs))
	if n == 0 {
		return
	}
	waves, lives, gold, kills = int(float64(waves)/n), int(float64(lives)/n), int(float64(gold)/n), int(float64(kills)/n)
	time /= n
	return
}
