package game

import "math"

func (s *State) Step(dt float64) {
	if s.Status != StatusRunning {
		return
	}
	s.Time += dt

	if !s.WaveActive && s.Wave < MaxWaves && s.Time >= s.NextWaveAt {
		s.StartWave()
	}

	s.spawnDue()
	s.moveEnemies(dt)
	s.fireTowers(dt)
	s.moveProjectiles(dt)
	s.decayBeams(dt)
	s.decayFx(dt)
	s.cleanup()
}

func (s *State) decayFx(dt float64) {
	if s.LeakFlash > 0 {
		s.LeakFlash -= dt
	}
	n := 0
	for _, f := range s.Fx {
		f.TTL -= dt
		if f.TTL > 0 {
			s.Fx[n] = f
			n++
		}
	}
	s.Fx = s.Fx[:n]
	if len(s.Fx) > 64 {
		s.Fx = s.Fx[len(s.Fx)-64:]
	}
}

func (s *State) spawnDue() {
	for s.SpawnIdx < len(s.SpawnQueue) && s.SpawnQueue[s.SpawnIdx].At <= s.waveTime() {
		e := s.SpawnQueue[s.SpawnIdx]
		s.SpawnIdx++
		s.SeenEnemies[e.Kind] = true
		spec := &EnemySpecs[e.Kind]
		hpMul := HPScale(s.Wave) * s.Diff.HPMul() * s.Map.HPMul
		spMul := SpeedScale(s.Wave)
		if e.Kind == EnemyBoss {
			bhp, bsp := BossMod(s.Wave)
			hpMul *= bhp
			spMul *= bsp
		}
		enemy := &Enemy{
			ID:         s.NextID,
			Kind:       e.Kind,
			HP:         spec.HP * hpMul,
			MaxHP:      spec.HP * hpMul,
			Speed:      spec.Speed * spMul,
			Pos:        s.Map.Path[0].Center(),
			SlowFactor: 1.0,
			Bounty:     s.Map.ScaleGold(spec.Bounty),
			Lives:      1,
			Armor:      spec.Armor,
		}
		if spec.Lives > 0 {
			enemy.Lives = spec.Lives
		}
		s.NextID++
		s.Enemies = append(s.Enemies, enemy)
	}
}

func (s *State) waveTime() float64 {
	return s.Time - s.WaveStart
}

func (s *State) moveEnemies(dt float64) {
	for _, e := range s.Enemies {
		if e.Dead || e.Leaked {
			continue
		}
		speed := e.Speed
		if e.Slowed(s.Time) {
			speed *= e.SlowFactor
		} else if e.SlowUntil <= s.Time {
			e.SlowFactor = 1.0
		}
		e.Prog += speed * dt
		if e.Prog >= s.Map.TotalLen {
			e.Leaked = true
			s.leakEnemy(e)
			continue
		}
		e.Pos = s.Map.PointAt(e.Prog)
	}
}

func (s *State) fireTowers(dt float64) {
	for _, t := range s.Towers {
		t.CD -= dt
		if t.Flash > 0 {
			t.Flash -= dt
		}
		if t.CD > 0 {
			continue
		}
		target := s.acquireTarget(t)
		if target == nil {
			continue
		}
		t.CD = 1.0 / t.RoT()
		t.Flash = 0.12
		t.FlashTo = target.Pos
		spec := t.Spec()
		switch t.Kind {
		case TowerSniper, TowerFrost, TowerTesla:
			s.beamShot(t, target)
		default:
			s.Projectiles = append(s.Projectiles, &Projectile{
				ID:      s.NextID,
				Kind:    t.Kind,
				Pos:     t.Pos(),
				Target:  target.ID,
				LastPos: target.Pos,
				Dmg:     t.Dmg(),
				Splash:  spec.Splash,
			})
			s.NextID++
		}
	}
}

func (s *State) acquireTarget(t *Tower) *Enemy {
	var best *Enemy
	bestVal := math.Inf(-1)
	te := t.Pos()
	for _, e := range s.Enemies {
		if e.Dead || e.Leaked {
			continue
		}
		if e.Pos.Dist(te) > t.Range() {
			continue
		}
		var val float64
		switch t.TargetMode {
		case TargetStrongest:
			val = e.HP
		case TargetClosest:
			val = -e.Pos.Dist(te)
		default: // TargetFirst
			val = e.Prog
		}
		if val > bestVal {
			best = e
			bestVal = val
		}
	}
	return best
}

func (s *State) beamShot(t *Tower, target *Enemy) {
	switch t.Kind {
	case TowerTesla:
		hops := TowerSpecs[TowerTesla].Chain[t.Level-1]
		dmg := t.Dmg()
		chain := []Pos{t.Pos(), target.Pos}
		s.applyDamage(target, dmg, t.Kind)
		hit := map[int]bool{target.ID: true}
		cur := target
		// The chain propagates from the last hit enemy's position even if
		// that hit killed it — lightning jumps off the corpse. Skipping the
		// hop when the target dies would neuter the tower against exactly
		// the squishy swarms it is built for. Dead enemies are still
		// skipped as jump destinations below.
		for i := 0; i < hops; i++ {
			dmg *= ChainFalloff
			var next *Enemy
			nextD := ChainRange
			for _, e := range s.Enemies {
				if e.Dead || e.Leaked || hit[e.ID] {
					continue
				}
				d := e.Pos.Dist(cur.Pos)
				if d < nextD {
					next = e
					nextD = d
				}
			}
			if next == nil {
				break
			}
			hit[next.ID] = true
			s.applyDamage(next, dmg, t.Kind)
			chain = append(chain, next.Pos)
			cur = next
		}
		s.Beams = append(s.Beams, &Beam{From: chain, TTL: 0.1, Max: 0.1, Kind: t.Kind})
	default:
		s.applyDamage(target, t.Dmg(), t.Kind)
		s.Beams = append(s.Beams, &Beam{
			From: []Pos{t.Pos(), target.Pos},
			TTL:  0.08, Max: 0.08, Kind: t.Kind,
		})
	}
}

func (s *State) moveProjectiles(dt float64) {
	for _, p := range s.Projectiles {
		if p.Dead {
			continue
		}
		target := s.Enemy(p.Target)
		if target != nil && !target.Dead && !target.Leaked {
			p.LastPos = target.Pos
		}
		d := p.Pos.Dist(p.LastPos)
		step := ProjectileSpeed * dt
		if d <= step || d < 0.05 {
			p.Dead = true
			if p.Splash > 0 {
				for _, e := range s.Enemies {
					if e.Dead || e.Leaked {
						continue
					}
					if e.Pos.Dist(p.LastPos) <= p.Splash {
						s.applyDamage(e, p.Dmg, p.Kind)
					}
				}
				s.Fx = append(s.Fx, &Fx{
					Pos: p.LastPos, TTL: 0.25, Max: 0.25,
					Color: 203, Ring: p.Splash,
				})
			} else if target != nil && !target.Dead && !target.Leaked {
				s.applyDamage(target, p.Dmg, p.Kind)
			}
			continue
		}
		ux := (p.LastPos.X - p.Pos.X) / d
		uy := (p.LastPos.Y - p.Pos.Y) / d
		p.Pos = Pos{p.Pos.X + ux*step, p.Pos.Y + uy*step}
	}
}

func (s *State) decayBeams(dt float64) {
	for i := range s.Beams {
		s.Beams[i].TTL -= dt
	}
	n := 0
	for _, b := range s.Beams {
		if b.TTL > 0 {
			s.Beams[n] = b
			n++
		}
	}
	s.Beams = s.Beams[:n]
}

func (s *State) cleanup() {
	n := 0
	for _, e := range s.Enemies {
		if !e.Dead && !e.Leaked {
			s.Enemies[n] = e
			n++
		}
	}
	s.Enemies = s.Enemies[:n]
	m := 0
	for _, p := range s.Projectiles {
		if !p.Dead {
			s.Projectiles[m] = p
			m++
		}
	}
	s.Projectiles = s.Projectiles[:m]
	if s.Status != StatusRunning || !s.WaveActive {
		return
	}
	if s.SpawnIdx >= len(s.SpawnQueue) && len(s.Enemies) == 0 {
		s.WaveActive = false
		s.Gold += s.Map.ScaleGold(WaveBonus(s.Wave))
		s.Score += WaveBonus(s.Wave) * 100
		if s.Wave >= MaxWaves {
			s.Status = StatusVictory
		} else {
			s.NextWaveAt = s.Time + AutoWaveDelayFor(s.Wave)
		}
	}
}
