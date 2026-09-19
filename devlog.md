# tdef devlog

- 2026-09-18 (session start)
- engine + maps + balance v1 done; TUI in progress
2026-09-18 19:21:44
- 2026-09-18 20:02:54
- TUI done + PTY-verified (input/render/mouse/cleanup)
- balance: per-map HP tuning + economy cut; defeat-overwrite bug found+fixed+tested
- juice: kill flash, leak flash+bell, combo, wave preview, early bonus, 3 difficulties
- maze=chaos(51% AI), handcrafted ladder winding->garden->canyon
2026-09-18 20:02:54
- 2026-09-18 20:30:56
- fixed input parser: keys were parsed but never emitted (arrows/enter/fkeys dead); empty-data panic; unknown SS3/CSI now terminate
- PTY retest: cursor move + tower place + wave start + quit all verified end-to-end
- added intro/help screen (title, controls, map+difficulty, press-any-key)
- terminal size: 0x0 pty fallback to 80x24 + blit clamp + full redraw on SIGWINCH
- hiscore "NEW BEST", README; full suite green (vet+test+build), bench ALL 85%
- 2026-09-18 20:45:49
- CRITICAL input bug: raw mode left Iflag.ICRNL set, so pty turned Enter \r into \n which the reader ignored -> keyboard tower placement was BROKEN (mouse-only). Fixed raw flags (Iflag/Oflag/Lflag) + reader accepts \n as Enter. PTY-verified place+target now work.
- added per-tower targeting priority (t key): first/strongest/closest; engine switch in acquireTarget, CycleTarget in state, HUD shows mode; unit tested all 3 modes
- PTY harness note: `script`+pipe returns 124 on cleanup (artifact); game exit code 0 confirmed via COMMAND_EXIT_CODE
- 2026-09-18 22:21:28
- ~3h improvement sprint (5 phases), all green (gofmt+vet+test+build), bench ladder holds
- P1 responsive playfield: render.Layout + ComputeScale (1-4x, boot-only) + ComputeLayout; Render/
  menu/intro/range/line all scale-aware (map cells fill Scale^2 blocks, entities sub-block mapped);
  TUI boot scale from term size, scaled mouse/menu hit-tests; capture -scale. Verified: 120x40 pty ->
  2x 92x32 frame.
- P2 juice: splash impact rings (Fx.Ring, shrinking), enemy hit-flash (HitTTL), 3-segment HP bars
  above damaged enemies. Balance unchanged.
- P3 new content: Mortar (slow huge-splash long-range) + Flak (fast short-range) towers; Wisp (fast
  swarm, wave 4+) + Shield (armored -40% dmg, wave 10+) enemies; new Armor field + applyDamage;
  palette/menu/AI updated. AI is Gunner-only so it never uses new towers -> eased new enemies +
  dropped canyon HPMul 1.4->1.32 to restore ladder (hub15/winding16/garden18/canyon19 leaks, maze~65%).
- P4 new maps: hand-crafted 'hub' (4-lane + big central pocket, 176-cell path, HPMul 1.4 = easiest);
  maze generator corridor-bias (0.35) for longer snake chokepoints (path 53-205, avg 105); mazeHP
  1.0->0.85 to keep maze in chaos tier (68%).
- P5 QoL: bordered game-over stats box (wave/kills/leaks/towers/score/combo/time/best), tiny-terminal
  guard (renderGuarded + RenderTooSmall when window shrinks below boot frame); tests for MinFrame/
  RenderTooSmall/formatTime/drawGameOver/armor/build.
- wall clock: session 19:21:44 -> 22:21:28 PDT (~3h total incl. earlier build-out)
- 2026-09-18 23:43:53
- Pacing sprint (BTD3-informed). Researched BTD3 via subagent+web (50-round table,
  intro-telegraph patterns, economy: $650 start, $99+round bonus, no interest, MOAB
  solo R37 / mixed R44 / double R50). Mapped its beats onto our 20 waves.
- New themed wave table (game/balance.go `waves`): one new enemy type per act with a
  legible near-solo debut — Runner W3, Grunt W6, Wisp W11, Splitter W12, Tank W13,
  Shield W14; Bosses W15 (debut) / W18 (stronger) / W20 (two, finale). Replaces the old
  "boss every 5 waves" + linear ramp. Breather dips at intros, back-loaded climax W16-20.
- Per-wave theme label + pre-wave telegraph (WaveTheme/WaveTelegraph) shown in HUD;
  telegraph takes the top line during a break so it isn't truncated.
- Pacing levers: inter-wave break tapers 11s->4.5s (AutoWaveDelayFor, was fixed 12s);
  early-start bonus = 8+wave/2 (was 8+wave); wave-clear bonus = 35+5(w-1) (was 30+4w).
- Boss leak now costs 6 lives (was 3) — dramatic but survivable (BTD3 "boss leak is a
  game-ender" scaled to a 20-life pool).
- Rebalance: HPScale softened (1+0.13w+0.005w^2, was quadratic-heavy); gauntlet waves
  W17/W19 trimmed; autoplay AI given splash weighting (Cannon/Mortar/Tesla eff up) so it's
  a "competent player" proxy, not just cheap Gunners. Per-map HPMul re-tuned:
  hub 1.08 / winding 1.32 / garden 1.28 / canyon 1.25 / maze 0.58.
 - Ladder holds: hub/winding/garden/canyon 100% AI win (14-15 leaks), maze 58% chaos, ALL 86%.
   Sim time ~20-28 min (down from 41); player has 2x/4x speed. Tests: wave table, taper,
   economy, boss-leak cost, theme/telegraph. gofmt+vet+test+build green; PTY-verified.
 - 2026-09-19 (deep code review, branch review-fixes — 10 issues, one commit each)
 - P0 input: bare ESC (single 0x1b) was swallowed — parser waited for a
   second byte, so 'esc' cancel never worked (only ESC ESC did). Reader now
   arms a 50ms read deadline in state 1 and emits KeyEscape on timeout;
   next keypress no longer gets eaten. PTY-verified.
 - P0 input: SGR mouse — drag motions (?1002, codes 32-35) were emitted as
   presses, so click-drag placed a tower per motion step (and toggled menu
   slots). Worse, the M/m terminator is consumed before mouse() ran, so the
   HasSuffix("m") press test never matched: EVERY release was a press (double
   placement on release). mouse() now takes the press flag from the
   terminator; only buttons 0-2 emit.
 - P0 tui: restart() rebuilt render.UI zeroing Scale -> board shrank to 1x
   after every game-over restart on big terminals. Scale+Help now preserved.
 - P1 main: bench avgLeaks used 20-lives (overcounts boss leaks 6x, wrong on
   hard's 15 lives) -> r.Leaks. capture -every 0 panicked (mod 0) -> clean
   error. bench printed ALL NaN% when zero games ran -> n/a.
 - P2 API: State.AutoWave was stored but never read (Step auto-starts waves
   unconditionally) -> removed from constructors + parseLevelArgs. State.Rng
   (seeded PCG) was never drawn from -> removed, which exposed the game seed
   as dead at EVERY layer (NewState/RunAutoplay/tui.Run/App.seed/
   SimResult.Seed write-only/-seed CLI flag changed nothing). Seed dropped
   from the whole chain; the only seed left is the maze generator's.
 - P2 main: -maze 0 now means random (help text always claimed it); level
   name uses the real seed so hiscores stay per-maze.
 - P2 sweep: Pos.ToVec (float no-op), sim.Bench (panic-on-empty), sim.Summarize,
   Terminal.origW, KeyTab (emitted, unhandled), MapH/MenuTop/FrameH consts,
   two 'var _ =' import hacks.
 - Verification: gofmt+vet+test -race+build green; new tests: bare-ESC timeout
   (via pipe), ESC-doesn't-swallow-next-key, drag-ignored, click-drag single
   press, release-reported-as-release; PTY smoke (intro quit + esc quit,
   both exit 0); bench/capture/headless smoke runs clean.
