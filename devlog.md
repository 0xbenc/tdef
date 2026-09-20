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
 - 2026-09-19 (second pass — review items 11-16, branch review-fixes)
 - hiscore: Save now writes temp-file + rename (atomic) — a crash mid-write
   can no longer truncate the file; Load treats a corrupt file as empty and
   the next Update self-heals. Maze runs are keyed by seed, so the table
   grew unbounded; pruneMaze keeps the top 32 maze scores (deterministic
   ties), hand-crafted levels always kept. First tests for the package.
 - NEW BEST bug: Update returned the PRE-update best, so a record run
   displayed 'NEW BEST <old score>'. Returns post-update best now.
 - capture ANSI dumps used bare cursor-up (\x1b[M, column-preserving): on a
   terminal wider than the frame, row 2+ started at column W+1 and smeared.
   Now CUP per row (\x1b[row;1H), position-correct at any width; verified on
   real dumps (0 cursor-ups, 1 CUP per row).
 - blit diffed against a map[int]Cell rebuilt every frame (up to ~12k
   entries at 4x, hashing + per-entry allocs at 30fps). Now a flat []Cell
   reused across frames, snapshotted with one copy; nil/size-mismatch still
   forces the 2J full redraw (startup/winch/ctrl-L/restart/too-small).
 - Shrinking the terminal below the frame used to keep the sim running blind
   behind the 'enlarge' notice (silent life drain). renderGuarded now pauses
   a running game alongside the notice; 'p' resumes. Notice shortened to
   27 cols so it survives the smallest windows that show it.
 - Enemies in the top map row painted their HP bar over the HUD stat line;
   bar now skipped when it would land above the map origin. Selected-tower
   info line hit ~71 cols at 1x and clipped; compressed to max 58 and
   guarded by TestTowerInfoFitsFrame (every tower/level/mode).
 - Wheel-up jumped straight to 4x; wheel now cycles 1->2->4->1 like 'f'
   (shared cycleSpeed, wrap both ways, tested). Leak message pluralizes
   ('-6 lives'). Hand-rolled itoa (returned "" for negatives) -> strconv.
  - Verification: gofmt+vet+test -race+build green (51 tests, hiscore now
    covered); PTY smoke exit 0; capture dump checked byte-level; README
    mouse/pause lines updated to the new behavior.
  - 2026-09-19 (third pass — review items 1-5, branch review-fixes)
  - Intro screen: the four help lines + prompt did not fit the 62-col
    frame at 1x; the prompt overwrote help line 3. Refit to H-5..H-1
    (TestRenderIntroFitsFrame, later removed with the intro itself).
  - Tesla chain: the chain bolt kept its original target even when it
    died mid-chain, so the chain never propagated. Now re-acquires from
    the killed target (game/step.go); TestTeslaChainFiresWhenTargetDies.
    Bench unchanged — the autoplay AI never buys a Tesla (gold floor
    stays under its 200 cost).
  - UI.Selected zero value: a fresh render.UI zeroed to tower ID 0, the
    first tower built — it read as "selected" and hid the cursor.
    render.NoSelection = -1 now pinned in every UI constructor
    (Run/restart/capture); tests.
  - Data race: the winch goroutine wrote Terminal.size while the main
    loop read it (-race fires under a resize storm). The notifier is
    now pure; the main loop calls RefreshSize() after consuming the
    token. Old build fails the PTY resize-storm probe, new build clean.
  - Level pips drew before the tower glyphs, so two adjacent upgraded
    towers erased each other's pips. Pips now draw after the glyphs,
    skipping occupied cells.
  - 1x clipping: the selected-tower info line (up to ~71 cols at Lv3)
    and the wave/break lines exceeded the 62-col frame. Compressed and
    restructured; TestMenuLinesFitFrame/TestHUDLinesFitFrame.
  - Mouse after game over: a click on grass behind the overlay built a
    tower (handle routed mouse before the status check). Guarded.
  - Still open from pass 3: UTF-8 reader (incomplete sequence emits the
    lead byte raw), LoadMap validation (missing S/E -> zero vec /
    panic), unbounded capture -scale, Path/TotalLen spawn-cell
    semantics, and the LOW sweep (hiscore type guards, SGR reset,
    wheel modifiers, dead code, silent flag fallbacks).
  - 2026-09-19 (title screen + start menu, branch title-screen)
  - New pre-game flow: animated title screen (beveled TDEF slab logo,
    a miniature battle demo, tower/enemy roster, best score, blinking
    prompt) -> main menu (Start / Help / High Scores / Quit) -> level
    select (built-in maps, a procedural maze row with an editable
    seed, difficulty, and a live map preview on tall terminals).
  - Design via three read-only subagents (architecture / visual /
    impl+test), consolidated into a Screen state machine on App.
    Renderers are pure functions (render/screens.go) of
    (w, h, state, 30fps frame counter, palette), so the animation is
    deterministic and unit-testable; Rects helpers share the layout
    functions with the renderers, so mouse hit-testing cannot drift.
  - tui: the App now runs six screens (title/menu/help/hiscores/level
    select/game). The loop advances the sim only on ScreenGame
    (stepGame) and ticks the title animation counter on ScreenTitle.
    Bare `tdef` enters the title flow (RunMenu); explicit -level/-maze
    still start directly (Run); -diff preselects the difficulty. The
    old press-any-key intro and its dead overT field are gone.
  - game: MazeFromSeed centralizes the procedural maze size
    (MazeW/MazeH) so main, bench and the level select agree.
  - Verified: gofmt+vet+test -race+build green; PTY smoke of the full
    flow (title -> menu -> level select -> typed maze seed 1234 ->
    maze1234 game; hiscores with a real table) all exit 0; bench
    identical to baseline (canyon/garden/hub/winding 100%, maze 62%,
    ALL 66%).
  - 2026-09-19 (terminal-sized frame + btop screens, branch title-screen)
  - In-game frame now fills the whole terminal: one rounded box edge to
    edge, playfield centered at the largest integer scale (1-4) that
    fits. The scale formula is unchanged (min((tw-2)/mW,(th-6)/mH)) but
    is now recomputed every frame, so a live resize re-scales and
    re-centers on the next draw instead of waiting for a restart.
  - UI moved into the frame chrome. Top border carries embedded
    segments (btop grammar: ┐title┌ flush at x+2, bold title, border
    color on the brackets): brand, level·diff, wave, pause, stats. Drop
    rule when the frame is narrow: level·diff, then brand, then pause —
    wave and stats never drop, so 62 cols always keep the essentials.
    Row 1 below the border is the message row (placement hints,
    telegraphs, break previews). Bottom border: the seven tower slots
    (4+3, cost-colored, selected highlighted with btop's selected-bg
    95), a context hint line, and the selected tower's stats.
  - Box-drawing helpers live in render/frame.go (drawRoundedBox,
    embedSegment, centerEmbed, truncateRunes); the ┐title┌/┘title└
    grammar and palette were verified against the btop++ source
    (btop_draw.cpp / btop_theme.cpp) rather than guessed.
  - Game layout is computed once per frame by render.GameLayout — the
    single source of truth shared by the renderer and the mouse
    hit-testing, so they cannot drift. Old ComputeLayout/MenuSlots/
    drawHUD/UI.Scale deleted at the cutover.
  - Game-over restyled to a 46x12 rounded box: VICTORY/DEFEAT title in
    the border, two-column stats, NEW BEST callout, restart/quit hint
    in the footer.
  - The five pre-game screens share the same framing: full-window
    rounded box, screen title embedded in the top border, footer
    segment group (┘key text ─ key text└, hotkeys bold red) centered in
    the bottom border, content in a 19-row band centered at any height.
    Title's prompt moved to the footer; help is a fixed two-column
    table; hiscores is an aligned table (gold/silver/bronze rank
    colors, scroll arrows); level select gets bracketed difficulty
    labels and a content-sized [ PREVIEW ] sub-box for the map (scale
    2 when it fits, 1 below that, hint line below that — boundary
    tested at 80x29 vs 80x30).
  - capture -scale 1-4 pins the render to the virtual terminal size
    that yields the scale (62x19 / 92x32 / 137x45 / 182x58); frame
    dumps are byte-identical in any environment.
  - Below 62x19 the game forces a pause and draws a too-small notice;
    growing back resumes on the next frame.
  - PTY verification found the bare-ESC deadline never armed: pty file
    descriptors do not support SetReadDeadline (the error was ignored),
    so a lone ESC left the parser waiting for a CSI until the next
    keypress — an arrow typed within the 50ms window was misparsed and
    lost. Replaced with a mutex-guarded time.AfterFunc fallback, and a
    second ESC inside the window now settles the first as a definite
    Escape while keeping the parser armed (TestBareEscapeThenArrow,
    TestEscapeThenArrowInWindow). Verified on a real pty: esc-back in
    the menu flow, resize storms, all clean under -race.
