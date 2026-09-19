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
