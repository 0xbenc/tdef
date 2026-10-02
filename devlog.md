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
- 2026-09-19 (title screen: scripted battle + shockwave reset, branch title-screen)
  - The title is now one 545-frame (18s) deterministic loop — a pure
    function of the frame counter, so it is unit-testable and the loop
    point is exact (frame 544 renders byte-identical to frame 0;
    TestTitleLoopSeam checks it at three sizes).
  - 0-360: a scripted battle in the [ BATTLE ] box. Five different
    towers (gunner/cannon/frost below the path, sniper/mortar above)
    fire on fixed cooldowns at scripted enemy waves: wave 1 (minions)
    is held completely, wave 2 (runners) mostly held — one leaks and
    flashes the exit red — and wave 3's boss walks through the exit
    while its minions are picked off. Beams, shells, splash rings,
    death bursts, and a wave label in the box border (WAVE 1/2/3 ->
    BREACH). An energy packet flows the path when no enemy is around,
    so the standby screen is never static.
  - 360-373: the exit overloads — a heat wash blooms out from E.
  - 373-445: full-frame whiteout static, then a shockwave from the
    exit (white front, hot flickering trail, re-sparks) eats the frame
    and cools into a glowing grid; 445-475 the grid holds, then a
    black beat with one ignition spark.
  - 475-544: reboot — the border draws itself from the top-left (pen
    tip glows), the TDEF slab drops in row by row with a white flash
    on arrival, the tagline types on, then demo/roster/best/footer
    return; the battle comes back at frame 0 state, closing the loop.
  - Two reset effects were prototyped side by side: this shockwave
    (radial detonation from the breach point) and a grid-surge
    scanner wipe. The shockwave won on drama and narrative grounding
    (the breach is what detonates the screen); the surge variant was
    deleted after the comparison.
  - Verification: gofmt/vet/test -race green; new tests cover the
    phase sequence (TestTitlePhases), the blast leaving a complete
    grid lattice (TestTitleBlastCoversFrame), the wave script
    (TestTitleBattleScript: deaths, leak flash, tower rows), and the
    seam; PTY end-to-end (BATTLE box -> BREACH at ~12s -> q) and all
    seven flows pass on the fresh binary.
- 2026-09-19 (title battle retimed to a true 1x opening, branch title-screen)
  - The demo read as a sped-up highlight reel: a minion crossed the path at
    3x 1x-gameplay speed and the towers fired 1.5-2.6x faster than their
    real RoT. Checked against game/balance.go: at scale 2 (a ~100-wide
    terminal) a minion crawls 6.4 screen cells/s; the old demo ran 19.
  - Fix: wave 1 minions now cross in 445 frames — the true 1x scale-2
    speed — spawned 1.5s apart, so the ~19s opening feels like the real
    game. Tower cooldowns are the real level-1 RoT in frames (Gunner 23,
    Cannon 55, Frost 33, Sniper 86, Mortar 67; was 9/30/18/38/46).
  - The relative speeds were also backwards: in-game the Runner (2.6) is
    slower than the Minion (3.2), but the demo had the runner zipping by.
    Waves 2-3 are compressed (runners 300f cross, boss 600f — slower than
    a minion so it still reads as the slow threat) and the boss breaches
    at frame 1500.
  - A true 1x boss walk is ~1 minute (0.8 cells/s), so the boss is
    necessarily compressed; the "1x" promise applies to the opening, which
    is what most viewers will see before pressing enter.
  - Loop is now 1685 frames (~56s, was 545/18s). The ambient path packet
    slows to one cell per frame. Phase/seam/blast/script tests retimed;
    seam, blast coverage, -race suite, and a PTY run to WAVE 3 all green.
- 2026-09-20 (title boot cinematic + weapon fan + signature, branch title-screen)
  - The title opened straight into the battle; it now has a one-shot
    cinematic on every visit. RenderTitle gained a `boot` parameter
    (frames since this title visit started; tui/app.go stamps
    titleBootAt on screen entry), so the whole screen is a pure
    function of (w, h, frame, boot, scores, pal) and the boot replays
    each time you return to the title.
  - Boot v1 (124f, ~4s): ignition point + three expanding rings on
    black (0-14); a light pen traces the TDEF slab out of digital
    noise, white tip with a cooling trail, each letter flashing white
    as its trace completes (15-75); the slab ignites white over a dim
    grid (76-79); the tagline decodes left to right with a caret
    (80-93); chrome fades in in three waves — border/footer, empty
    battlefield, roster/best (94-123).
  - The attract loop was restructured around a 15s idle gate: 450f of
    standby (full UI, empty battlefield — the demo no longer starts
    mid-breach) then 120f of "BATTLE"/"WAVE 1" decoding on with the
    five towers powering up left to right, then the unchanged 1685f
    battle script. The reboot now ends exactly on the standby frame,
    so the loop (2255f) is seamless.
  - Boot v2 (273f, ~9s), per the user: the ignition takes twice as
    long and the second half is a weapon fan where each letter fires
    a different weapon, left to right, at a lock-on reticle on the
    frame border, in the letter's own color:
      - 15-19  rings freeze and settle dim; a crosshair zaps out
      - 20-29  the cross collapses into the slab's bounding box, which
               draws itself (10 cells/frame) over a scan flicker
      - 95-224 T Gunner (5 tracer shots, letter recoil), D Cannon
               (6f recoil + slow shell + AOE starburst), E Sniper
               (14f charge — the letter's top row fills with a ramp —
               then a 3f full-length beam), F Tesla (hash-jittered
               chain arc, two branch strikes, sparks)
      - each shot locks a reticle (cross + 4 dots) at ray ∩ frame
               edge and shatters it into 6 fragments on impact; a
               settle beat flickers the box before the flash
      - 225+ flash, subtitle, chrome fade (as in v1)
  - The signature "by 0xbenc" is embedded in the bottom border's right
    section (FG 238, mirroring the TDEF embed in the top border) on the
    standby, the battle, the boot fade-in and the reboot; it never
    collides with the centered footer.
  - Verification: gofmt/vet/test -race green; TestTitleBootSequence
    covers every phase checkpoint plus the boot->standby seam,
    TestTitleStandbyEmpty, TestTitleSigInBottomBorder, TestTitleIdleGate
    (chrome up at 449, battle text at 451, seam 450 vs 2705); PTY: all
    seven original flows + a new title_idle_battle flow (waits for the
    idle gate to fire the battle), all passing on the fresh binary.
    Harness notes: the PTY marker moved from "TDEF" (now only visible
    during the boot) to "by 0xbenc" (first full-chrome frame), and the
    ptydrv smoke now waits for the in-game "winding" HUD line because
    the title intro no longer says "press any key to start".

2026-09-21
- Lore skin pass per LORE.md ("The Last Monster" / Malgrath): names and copy
  only, zero gameplay/stat changes. Done as two parallel subagents
  (game/ package; render/ + tui/).
- Towers: Orc Gunner, Cannonier, Frost Mage, Ranger, Lightning Mage,
  Trebuchet, Gnoll Slingers (TowerSpecs.Name; the UI flows from there).
- Enemies: Squire, Rogue, Mercenary, Wizard, Necromancer, Paladin,
  Centurion, The Player (EnemySpecs.Name; Short glyphs unchanged).
- Wave themes reskinned: scouts/raid/column/assault/the coven/the risen/
  the vanguard/the wall/the player/siege/the end; 10 new telegraphs in
  guild-rumor voice ("The Player has set out. If they reach the heart,
  it ends.").
- UI: title demo label IDLE START SCREEN -> THE SIEGE (wave segment
  x20->14, WAVE 1 decode t34->18); help -> GRAK'S LEDGER + subtitle line;
  game-over/victory lore lines in the stats box ("Malgrath has fallen.
  The lair is clean." / "The lair is held. Malgrath endures."); leak
  toast -> "breach! -N ♥" / "the Player breached! -6 ♥" (header red-bold
  style hook leak! -> breach).
- Level display names: the Rotunda / the Long Halls / the Sunken Garden /
  the Rift / the Unmapped Depths (header + level select only; CLI level
  ids and hiscore keys stay the short forms).
- Slot widths re-verified at 62 cols: row 0 max "1 Orc Gunner 50" = 15 =
  cell 15, row 1 max "5 Lightning Mage 200" = 18 <= cell 20. Header
  elision boundaries shift with the longer names (existing drop logic).
- Tests updated across game/render/tui; gofmt/vet/test -race green;
  README + LORE.md updated.
- 2026-09-22 (overworld look-dev pass, branch overworld-lookdev)
- The lair — the 45x13 overworld map Grak walks before descending — got a
  full imagineer-grade look-dev pass on top of the "lair as front door" base
  (a5d54ee). One commit per concern, each green; RenderOverworld stays a
  pure function of (w, h, state, frame, palette), and the bench pass results
  stay byte-identical (the pass touches only render/ + tui/).
- NOTE: the user asked for the design to be done with parallel subagents
  (the repo's design->contract->implement culture). The subagent API was
  unreachable for the whole session ("Cannot connect to API" on every task
  launch, including connectivity probes), so the three design contracts
  (visual / narrative / impl) were written solo and consolidated into one
  spec; implementation then proceeded solo.
- Audit-first (C1): a frame-by-frame look at the rendered lair found concrete
  defects, fixed before any new work — pad edge-glow leaking outside the
  playfield (stray dots on the ledger row) clamped to the map rect; sealed
  rooms too dark because dim() had no floor (floored at 24); the first-run
  hint never showed (Grak starts on the Rift, so the old "empty floor line"
  trigger never fired) now shows whenever FirstRun; the flat one-colour
  ledger now wears each room's accent (rift 208 / halls 110 / garden 45 /
  rotunda 178 / heart 220).
- State plumbing (C2): OWState gained the frame-counted TTLs the pass drives
  — BootTTL (the arrival), BlastTTL (the heart-unseal shockwave), ReturnFX
  (which floor a finished defense answers on, won/broke), and a three-cell
  walk trail (Trail/TrailAge) replacing the old single Prev/StepTTL step.
  The boot is armed once per session on first entry and gates input while it
  plays (walk/descend/mouse; esc and q still pass). The heart-unseal blast
  arms exactly once (BossReady false->true), never re-firing on an
  already-unsealed heart or a renown switch.
- Arrival cinematic (C3): on first entry the map is a dark void and only the
  Rotunda (the dragon's heart) is lit; a light front sweeps out from the
  heart across the corridors, unmasking the map as it passes, and Grak wakes
  with it (masked until the light reaches him). The waking narrates itself —
  "the lair stirs in the dark" -> "the heart beats — light runs the
  corridors" -> "Malgrath: …grak. the guild still hunts." — then hands back
  to the steady state. Seam holds: the last frame is the steady state at
  62x19 and 137x45.
- Per-room ambient (C4): the Rift's fissure (╎) flickers with rising embers;
  the Rotunda's pillars are ┃ with a gold hoard glinting either side of the
  heart; the Long Halls' static fallen-hero flicker becomes a procession of
  three heroes (t/g/r) marching the corridor on a 240-frame cycle, fading at
  the edges; the Garden's ring steps through a depth gradient (31/27/23) with
  a sunbeam on the peak; the Depths' ░ mists swirl and the Ø alternates Ø/ø;
  the void's grain twinkles and outcrops wear a lighter cap. Grak's flat aura
  becomes a two-tone lantern (orthogonal bright 214, diagonals 180) tinting
  every texture subcell. Landmark not occluded: when Grak stands on a room's
  centre the landmark shifts to the first free interior cell (rift -> (5,10),
  rotunda -> (21,6) over the hoard) — the audit's "Grak hides the landmark"
  defect.
- The beats (C5): returning from a defense runs a ring of light down the
  floor's route (bright front, two-cell tail) and flashes the room's chrome
  toward the result colour (gold 220 held / scarlet 167 broke), the ledger
  entry flashing white with it; unsealing the heart fires a shockwave across
  the whole map (bright front ring, warm trail) with the voice "the heart has
  unsealed — the final expedition stirs"; the descent flood gains a hot rim.
- Copy (C6): the result banner reads like the lair — "<floor> held — 20/20 ·
  the lair stands steadier" / "<floor> broke at N · the lair will mend" (the
  heart-held line "the heart is held — Malgrath endures" and the "· new best
  N" suffix kept). Longest banner 59 runes: full at normal widths, graceful
  ellipsis at the 62-col minimum.
- Verification: gofmt+vet+test -race+build green on every commit. New tests:
  TestOWGlowClamped, TestOWSealedDim, TestOWFirstRunHint, TestOwLedgerAccents,
  TestOWTrailFade, TestOWBootPhases, TestOWBootSeam, TestOWInputGateDuringBoot,
  TestOWBootArmedOnce, TestOWHeartBlastCheck, TestOWLandmarkNotOccluded,
   TestOWRotundaHoard, TestOWProcession, TestOWReturnFX, TestOWBlast,
   TestOWBlastDecay, TestOWRichDeterministic. Bench pass results byte-identical
   to baseline (canyon/garden/hub/winding 100%, heart 0%, maze 100%, ALL 83%).
- 2026-09-22 (gameplay look-dev pass, branch gameplay-lookdev)
- The in-game defense screen — the actual tower-defense playfield — got a
  massive imagineer-grade look-dev pass, bigger than the overworld one: the
  user asked to reconsider everything, including how towers, enemies and maps
  look, while keeping the integer scaling system (ComputeScale/GameLayout/
  MinFrame/CaptureSize, scales 1-4) intact. One commit per concern, each green;
  Render stays a pure function of (state, ui, palette, size, frame) and the
  whole pass touches only render/ + tui/ (the game/ engine is untouched, so the
  bench pass results stay byte-identical).
- NOTE: as with the overworld pass, the user wanted parallel subagents for the
  design (the repo's design->contract->implement culture). The subagent API was
  unreachable for the whole session ("Cannot connect to API" on every launch,
  including a connectivity probe), so the design contract was written solo to
  /tmp and implementation proceeded solo.
- Frame clock (C0, 97a353e): Render/GameFrame gained a `frame` int — the 30fps
  ambient tick that keeps advancing while paused and after game-over (unlike
  g.Time, which freezes). This is the clock the beats and end cinematics run
  off. Callers: tui (a.frameNo), capture (tick), all tests. Verified visually
  inert (captured frames byte-identical).
- The battlefield as a place (C1, 76ace3c): the terrain rework. The road is
  now drawn as directional box-drawing connectors (straight/corner/tee/cross
  from neighbour connectivity) on a distinct surface, so the route reads at a
  glance instead of a dotted line; walls are mottled rock (a stable per-cell
  hash picks a light/base/dark shade + speckles at 2x+); grass clearings are
  dark green with sparse tufts; the spawn is a breathing rift (magenta) and the
  exit is the lair's beating heart (a double-thump on the frame clock). The
  cursor treats road connectors as ground (isGroundRune). Scale-aware: 1x = one
  glyph/cell, 2x+ = textured blocks.
- Towers with presence (C2, 7b7d659): each tower is its coloured glyph on a
  dark pad (the pad fills the block at 2x+, so a tower reads as a structure,
  not a floating letter); the selected tower gets bright corner brackets in
  addition to the inverted glyph; firing towers flare their glyph white and kick
  a short tracer toward the target (t.Flash/t.FlashTo — the tracer preserves the
  terrain background and never overwrites the tower glyph); the range indicator
  is the tower's own beam colour on the boundary ring with a faint interior,
  sitting on the terrain instead of punching black holes. Level pips and their
  no-clobber invariant are unchanged.
- Enemies with identity (C3, d815544): an enemy's glyph keeps its own kind
  colour (identity over the old 3-step HP recolor that made every wounded enemy
  look the same); HP is read off the bar. Frost-slowed enemies tint bright cyan
  (e.Slowed), hit flashes stay white (hit wins the brief flash, then cyan
  shows). The boss sits on a dark pad with a cage of rails either side. Tanks
  (Paladin, Centurion, Necromancer, boss) always show their HP bar; lighter
  enemies only when wounded. The spawn rift moved to magenta to stay distinct
  from the frost cyan.
- Projectiles & beams with weight (C4, 022d08b): projectiles are no longer a
  uniform '+' — cannon and mortar fire heavy filled shells (●, mortar bolder)
  while gunner and flak fire light tracers (·), each in its tower's beam colour
  with a short dim trail so the motion reads; beams (sniper/frost/tesla chain)
  are bright continuous energy lines that overwrite terrain glyphs but preserve
  the background.
- The beats (C5, 5072f77): all pure functions of (state, frame). The leak — the
  frame's left/right edges throb red while the lair heart is struck (g.LeakFlash);
  a wave start — a centred banner ("THE SIEGE BEGINS" for wave 1, else "WAVE N")
  that fades as the wave gets under way; the boss entrance — a "THE PLAYER"
  banner and a regal purple edge pulse while the boss is young; the end — when
  the caller records the end frame (UI.EndAtFrame, set by the tui the moment the
  status flips) a cinematic plays before the stats box: a bright sweep races the
  playfield while the lair's verdict decodes in (gold "THE LAIR HOLDS" / red
  "THE LAIR FALLS"). Headless captures, which don't record the end frame, show
  the box at once as before. The end cinematics run off the ambient frame clock
  (g.Time is frozen at game over), so they play to the end.
- The end box reads like the lair (C6, 20ffbd5): the game-over box gains a
  verdict line (the lair's assessment, keyed off leaks for a hold and how far
  the defense got for a fall) under the existing lore, with the restart/quit
  hint moved down a row. The box keeps the VICTORY/DEFEAT titles, stats and lore
  the tests assert and still fits the minimum frame.
- Verification: gofmt+vet+test -race+build green on every commit; the existing
  render/tui tests were updated to the new glyphs/colours where the look
  changed, and the invariants (scale stepping, header elision, slot layout,
  footer fit, pip no-clobber, frame corners) are preserved. Bench pass results
  byte-identical to baseline (canyon/garden/hub/winding 100%, heart 0%, maze
   100%, ALL 83%). Every beat, the boss entrance, the frost tint and the end
   cinematics were eyeballed via `tdef capture` at 1x and 2x and via temporary
   eyeball harnesses (deleted after use).
- Road-corner fix (762b613): the user spotted that the road's corner connectors
  (┌┐└┘) read "backwards" and didn't form one continuous line. roadGlyph had
  mapped all four corners to the diagonally-opposite box-drawing glyph (n+w->┐
  instead of ┘, etc.), so every turn kinked. Each corner now connects the two
  directions the road actually goes (n+w->┘, n+e->└, s+w->┐, s+e->┌); straights,
  tees and the cross were already correct. Bench byte-identical.
- 2026-09-22 (per-level theming pass, branch gameplay-lookdev)
- The user: the maps shouldn't all share one background colour scheme, and each
  should have its own stinger/fascination and light animation. Two commits, each
  green; render-only, so the bench stays byte-identical (83% ALL).
- Each floor wears its own stone (C8, b75c2a5): a `Theme` (terrain palette,
  accent, stinger id) + `themeForLevel(id)` give each level a distinct look —
  warm bronze torchlit Long Halls, dragon-purple Rotunda with a gold road, mossy
  green Sunken Garden, volcanic black-red Rift, blood-soaked heart chamber, and
  cold teal Unmapped Depths (maze by prefix; heart is the boss floor). The
  drawMapPreview wall/grass/road block helpers now take a Theme; the level-select
  preview themes off the selected row (the maze row -> Depths). A test asserts
  no two levels share a wall colour.
- Each floor has a stinger and its own light (C9, 6257ed1): `drawTheme`
  dispatches on Theme.Stinger and paints a landmark + animated light over the
  terrain but under every entity (Render draws towers/enemies/cursor/spawn/exit
  after), so a stinger never hides gameplay. All are pure functions of (map,
  frame) — deterministic and scale-aware (placed at block centres): the Halls'
  amber torch sconces flickering on the road walls; the Rotunda's pile of dragon
  gold with a glint rolling across it; the Garden's glowing pool + drifting
  fireflies; the Rift's wall-splitting fissure with sparks rising out; the heart
  chamber's light beating (two staggered rings swelling out of the lair, ○); and
  the Depths' cold mist drifting down the room with runes (◈) waking on the
  walls. A test asserts each stinger's signature glyph renders (checked over a
  few frames, since the light animates). Eyeballed via `tdef capture` at 1x and
   2x; the boss (heart) isn't reachable via `capture` (it's LoadBoss, not in the
   level list), so its ring is covered by the test instead.
- 2026-09-22 (tower level tile, branch gameplay-lookdev)
- A long look at how a tower shows its level, driven by screenshots:
  1. pips trailing LEFT crowded a neighbour -> 2. a badge ABOVE the glyph, but
     the cell above a tower is usually the ROAD, so it dropped a foreign block
     into the lane -> 3. the level on the tower's OWN pedestal (glyph white,
     base = the full tower colour) -> 4. but a maxed tower was now a SOLID
     bright block, so the G had nowhere to pop (green on green, hard to read).
  Final: a tower is a dark TILE. A neutral dark body (235) keeps the glyph
  readable at every scale; a frame in the tower's colour gives it shape (a cap
  at 2x, a full ring at 3x+); and both step with level — frame brightens
  dim->bright->full, glyph steps dim->full->white. A maxed tower reads as a
  glowing white letter inside a full-colour frame. Level stays on the tower's
  own block, never the road; selection is the corner bracket.
- baseColor darkens a 256-colour toward black while keeping its hue (a naive
  index halving drifts hue — gunner green 46 would become a dark blue, 25);
  indexRGB/cubeLevel/rgbIndex support it. A neutral body (not a tint) is used
  so desaturated colours (cannon/gold) don't muddy into grey under the glyph.
  drawLevelPips + its pass are gone. Tests: pedestal glyph-stepping, no-clobber
  neighbour invariant, and TestTowerFrameChargesByLevel (frame ring at 3x).
  Glyph-vs-body contrast holds >= 2.1, maxed towers 4.77. Bench byte-identical
  (83% ALL). 9dbc023 -> 1669db7 -> f53db5c.
- 2026-09-22 (siege waits for a tower, branch gameplay-lookdev)
- The user: it's clunky that the siege just begins on a timer three seconds
  after the map loads, before they've committed to a defense. Now the first wave
  is held until at least one tower is placed — by auto timer or the early-start
  key alike; later waves are unchanged. StartWave refuses to open wave 1 with no
  towers (the rule lives with the wave-start logic, so both entry points respect
  it); the header shows a gold "build a tower" while the siege is held (kept short
  so it doesn't trip the header's elision at 80 cols), and the early-start key
  answers "the siege waits — build a tower to begin" instead of faking a wave 0.
  The autoplay AI builds on its first tick (before the first Step), so the bench
  stays byte-identical (83% ALL). TestSiegeWaitsForTower locks the gate in;
  TestWaveCompletion now commits a tower first (and its gold check accounts for
  the build). a5f7a79.

- 2026-10-01 (a wider, figurative lair)
- Expanded the overworld from 45×13 to 78×13, with larger room footprints and
  a bridge between the hoard and the far floors. Each floor has its own
  architectural silhouette: a broken basalt breach, Malgrath asleep over gold,
  vaulted halls, a temple with a submerged lower basin, and a timbered mine
  mouth with converging rails. The roof and floor enclose the cavern.
- The viewport uses the existing 1×–4× scale steps, following Grak horizontally
  with an eased camera. World-space rendering is cropped before the fixed
  voice, ledger, controls and frame are drawn. Mouse picking shares the camera
  transform; wheel jumps reveal their destination immediately. Minimum window
  remains 62×19. Arrival and heart-unseal sweeps now cover the wider world.
- Kept terrain shades inside the ANSI grayscale ramp: extending the old
  falloff reached index 231 (white), creating bright patches at the far edge.
  Warm/cool accents now live in the grain over dark stone. The Heart descent
  also resolves to the Rotunda artwork so its transition plays.
- Verification: existing progression and presentation tests adapted to the
  wider world; added camera/resize/picking, fixed-chrome clipping, small-window,
  distant-stone, camera-settling and scrolled-mouse regression coverage.
  Inspected rendered views at the Rift, Rotunda, Halls and Garden.

- 2026-10-01 (walkways and figurative scenery)
- Replaced the straight dotted crossing with a connected stone walkway. The
  route rises before the bridge, returns to its deck, then forks into an upper
  approach to the Halls and a lower waterside route to the Garden. The Rift
  approaches from its side. Paving, exposed edges, turns and forks all derive
  from the walkable corridor cells; overlapping branches are drawn once.
- Replaced scattered grain and outcrop tiles with composed cave silhouettes,
  dark arches on the distant wall, bridge piers, a waterfall, and an underground
  river. Torches, a fork signpost and a fallen expedition's bones give the
  crossing concrete landmarks. Stone shades remain within the grayscale ramp.
- Refined the rooms into material-based illustrations: masonry towers with
  banners, a broad temple roof and ivy over its basin, timber mine supports
  with a cart and lantern, and basalt around the Rift's lava. Malgrath has
  treasure chests beside the gold. Sparse fog and ripples replace block-sized
  texture glyphs. Grak gets horns, shoulders and stepping feet at 2×+, and the
  Halls' ghosts get heads and feet. One-cell glyphs remain readable at 1×.
- The mine's sign sits over its entrance; compact labels avoid covering the
  walkway. Regression tests check the final composed path at every scale and
  walking through both branches of the fork, alongside camera, picking,
  progression and transition coverage.
- Verification: go vet, go test -race, build and diff checks pass; inspected
  the composed scenes at 62×19 and 92×32.

- 2026-10-01 (single-character Grak)
- Returned Grak to one terminal character at every scale. He idles as @ and
  briefly alternates @ / & after walking; the extra head, shoulders and feet
  are removed. The existing lantern and fading footprints remain.

- 2026-10-01 (verify lightning's per-bounce reach)
- Confirmed that Lightning Mage targeting already uses tower range only for
  the opening target, then searches within ChainRange (2.6 cells) of the last
  hero hit. Added normal-firing regressions at all three levels with every
  bounce target outside the tower's range, verifying damage falloff, beam
  endpoints and bounce limits. A second case verifies that an oversized gap
  stops the chain. Documented the two distinct ranges in the README.

- 2026-10-01 (sealed rooms block entry)
- Movement now checks progression before entering any room footprint, even
  where a corridor continues through it. A sealed room remains blocked until
  its unseal finishes. Door approaches come from the route geometry; clicks
  and wheel visits stop there instead of teleporting inside. A portcullis
  marks the boundary. The reveal preview changes appearance only.
- Switching renown moves Grak out of newly sealed rooms and cancels the old
  renown's unfinished unseals. Tests cover all sealed interiors, doorway
  walking, opening completion, mouse/wheel shortcuts and renown changes.

- 2026-10-02 (floors with distinct routes and room to build)
- Replaced all five built-in defense layouts, borrowing the procedural maps'
  irregular bends and return passes instead of repeating horizontal lanes.
  The Rotunda coils around a broad inner court; the Long Halls step through
  unequal chambers with pillars; the Garden loops around planted islands and
  rounded clearings; the Rift climbs narrow ledges around rock masses; the
  Heart spirals inward through shared firing lanes. All remain 45x13, fitting
  the existing minimum 62x19 terminal and integer scales.
- Build tiles: hub 114 -> 231, winding 40 -> 234, garden 45 -> 230,
  canyon 32 -> 170, heart 60 -> 189. Each has over 100 useful placements
  within a level-one Slinger's range of the route. Bends offer shared coverage,
  while separated stretches, inner courts and ledges vary tower utility.
- Added per-map GoldMul, independently of HPMul, to fund wider defenses.
  Starting gold, enemy bounties (including necromancer children), wave-clear
  income and early-start bonuses scale and round to a whole coin. Tower and
  upgrade prices, refund percentage, lives and wave compositions stay the
  same. Unconfigured/custom maps and procedural mazes retain their economy.
  The HUD and wave-clear messages display the actual scaled rewards.
- Normal floor tuning (HP / gold multipliers): hub .65 / 1,
  winding .77 / 1.25, garden 1.02 / 1.4, canyon 1.04 / 1.3,
  heart 2.20 / 1.4. Autoplay finishes the four progression floors with
  20 / 14 / 12 / 5 lives respectively; the Heart falls on wave 20 without
  earned hearts, but holds with four earned hearts and four lives remaining.
  Autoplay builds 132 / 168 / 188 / 169 towers on the progression floors,
  compared with 114 / 40 / 45 / 32 before this pass. This is a deterministic
  balance proxy, not a substitute for human playtesting.
- Checked all floors on easy/normal/hard: easy holds all five, hard holds
  the Rotunda and falls on the later floors. The same 20 benchmark mazes
  retain their 60% win rate and identical results.
- Verification: new tests check unbranched routes, every painted road being
  used (no BFS shortcuts), useful build space, scaled spawn/split income,
  wave rewards, early bonuses, unchanged prices/refunds and visible payouts.
  Inspected text-rendered empty floors at 1x and 2x. gofmt, go vet,
  go test -race, build and diff checks pass; rebuilt ./tdef for play.

- 2026-10-02 (scenery in the nonfunctional terrain)
- Added a dedicated wall-only scenery pass shared by gameplay and map previews.
  The Rotunda wears vault cornices, a dragon relief and treasure vaults; the
  Halls get bonded masonry, fluted pillars, roof arches and hanging banners;
  the Garden's blocked ground becomes a flooded basin with reeds, composed
  trees and drowned arches; the Rift gets basalt terraces, peaks and an
  animated lava seam; the Heart gets ribbed edges, suspended chains, an altar
  and bone masks; seeded Depths get buried masonry and crystal seams.
- Authored illustrations find a solid wall footprint near their desired
  anchor, reserve it against overlapping scenery, and use compact drawings at
  1x and detailed drawings at 2x+. Materials join in terminal coordinates.
  Drawing is clipped to blocked cells and the playfield; no roads, build pads,
  map geometry or economy values change. Existing ambient floor effects remain.
- Darkened the Halls' and Heart's wall shadows and kept material accents dim
  enough that units and roads retain priority. The cursor now stays visible
  over wall illustrations, even where their glyph is not a ground glyph.
- Inspected all six environments in a coloured 2x contact sheet, and generated
  previews at 1x through 4x. Added exhaustive scenery checks for functional
  terrain/chrome preservation, deterministic output and unchanged map data,
  including two maze seeds and a small generated map. Added cropped-preview
  and decorated-wall cursor coverage. Tests, vet, race checks and build pass;
  autoplay results remain identical to the preceding geometry/economy pass.

- 2026-10-02 (waterfall direction)
- Reversed the waterfall's animation phase: its drops and gaps now travel
  downward rather than upward. Regression coverage follows the rendered
  pattern through two cycles at all four scales. Waterfall and walkway checks
  pass; rebuilt ./tdef.

- 2026-10-02 (start in the hub)
- Grak now starts at the Rotunda, the only initially open overworld floor.
  Holding it unseals the Rift; the existing Rift -> Halls -> Garden chain,
  two-floor Depths gate and four-floor Heart gate continue from there.
- Updated fresh-state and title-to-defense coverage for the hub start, added
  the Rift to sealed doorway/walking/mouse coverage, and checked that a loss
  leaves it sealed while a hub win opens it only after the unseal completes.
  Existing render fixtures explicitly choose their player position and open
  floors. Tests with the race detector, vet and build pass; rebuilt ./tdef.

- 2026-10-02 (a figurative overworld Rotunda)
- Replaced the thin ASCII dragon and scattered gold with three composed
  silhouettes inside the existing hub footprint: a rounded stone vault with
  thick piers, a solid sleeping dragon, and a low mound of gold. Malgrath's
  projecting muzzle, pale swept-back horns, folded triangular wing and curled
  tail identify the figure. A closed eye briefly opens; sparse facets catch
  the light on the hoard.
- Geometry scales with the room's terminal footprint. Rounded small-scale
  feature placement keeps the horns and eye attached to the head at 1x.
  Grak and the heart/portal remain visible over the figure; the room's
  unseal and result colours still apply. Walkways and progression are intact.
- Inspected coloured before/after views and all four scales. Updated the
  hoard check for solid dragon/gold masses at every scale and added a room
  footprint check to guard neighboring corridors. Tests, vet, race checks,
  build and diff checks pass; rebuilt ./tdef.

- 2026-10-02: added a static Malgrath portrait study, accessible with v in the
  Rotunda or `tdef dragon`. Shares the overworld pose at a larger tile footprint,
  with filled tapered horns, broad wing folds, a bent foreleg and a curled tail
  enclosing the hoard. Responsive framing; no cutscene timeline or progression
  effects. Escape/enter restores the lair exactly where it was.

- 2026-10-02: refined the Malgrath portrait: curved torso shading and a warm
  back rim, clearer wing folds, a tapered tail tip separated from the paw,
  chamfered muzzle and an anchored tooth, solid vault piers and irregular gold
  facets. Eyes, claws and gold details keep their surface color beneath them
  instead of punching dark cell-shaped holes into the illustration.

- 2026-10-02: designed Grak as a squat angular defender with blade ears,
  asymmetric tusks, a scarred broad jaw, worn red mantle, diagonal leather
  strap and a grounded builder’s mallet. Added a static responsive companion
  portrait: `tdef grak`, g in the Rotunda, and tab between both character
  studies. Lair progression and time remain unchanged while viewing them.

- 2026-10-02: added the definitive Gnoll Slingers lore portrait: a long-muzzled,
  hunched thrower winding a broad sling loop above an ammunition-feeding partner.
  Warm fur, charcoal manes, digitigrade legs, teal wraps and a shared stone
  pouch establish their shape language. Thin cord uses half-cell geometry.
  `tdef gnolls` opens it directly; tab cycles all three studies. Shared polygon
  drawing preserves Grak’s geometry. No upgrade variants or gameplay changes.

- 2026-10-02: designed the Orc Gunner lore portrait: a braced green orc in a
  red headcloth, iron shoulder slab and heavy boots, gripping a walnut-stocked
  iron gun with a brass flare and a deep muzzle bore. Bent support elbow,
  squint, tusks, three large cartridges and quiet smoke establish the pose.
  `tdef gunner` opens it directly; tab from the Gnoll Slingers reaches it.
  One definitive portrait, with no upgrade variants or simulation changes.

- 2026-10-02: designed two full lore portraits. Frost Mage: tall dark cloak,
  peaked hood, ancient gaunt face and icy beard, forked crystal staff, reflected
  cold planes and a shard above an open long-fingered hand. Player: monumental
  shoulder-rested sword, polished plate and gold trim, red crest and cape,
  composed half-smile, reaching hand, stolen tusk and ivory monster skull.
  Reviewed/refined both at small and large terminal sizes; moved frost glow
  behind the fingers and gave the skull explicit dark half-cell sockets.
  `tdef frost` and `tdef player` open them; tab cycles all six lore studies.
  Both are static definitive forms with no upgrades or simulation changes.

- 2026-10-02: designed the Cannonier lore composition: a sweating, braced
  loader pushes a cloth-tipped ramrod into the foreground bore while a
  smaller, grinning spotter checks over the breech. Cast iron planes, bronze
  hoops, timber braces, two spoked wheels and stacked cannonballs establish
  weight and scale. Added half-cell oval compositing for the round lip, rims
  and balls, preserving the uncovered background/underlying material. Older
  portraits keep their existing geometry. `tdef cannonier` opens it directly;
  tab from the Player reaches it and the seven studies cycle back to Malgrath.

- 2026-10-02: designed the Ranger lore portrait: a lean, long-eared monster
  at full draw, with a blade-shaped ear, narrow eye, angled quiver, torn cloak
  and planted leather boots. A tall recurved bow, taut triangular string,
  feathered arrow and faceted steel arrowhead define the aiming pose. New
  half-cell path strokes preserve the surfaces behind cords and shafts.
  `tdef ranger` opens the static study; tab from Cannonier reaches it, then
  wraps to Malgrath. No upgrade variants or simulation changes.

- 2026-10-02: designed the Lightning Mage lore portrait: a wiry horned monster
  arches backward between a descending fork and an outward discharge. Open
  receiving and directing claws, an upturned face, swept ivory horns, split
  violet robe, angular sash and planted feet define the silhouette. Half-cell
  bolts layer dim purple edges around pale electric cores; their width scales
  down in compact terminals to preserve the figure. `tdef lightning` opens
  the static study; tab from Ranger reaches it, then wraps to Malgrath. One
  definitive form, with no upgrade variants or simulation changes.

- 2026-10-02: completed the tower lore portraits with the Trebuchet: two
  timber A frames, a tapered iron-bound throwing arm, visible axle, free
  loaded sling, hinged ballast box and rope-wound windlass. A green winding
  monster grips the crank while an olive spotter braces on the rear rail.
  Dark rear timbers and lit front planes establish depth; open triangles
  keep the mechanism visible. `tdef trebuchet` opens the static study; tab
  from Lightning Mage reaches it, then wraps to Malgrath. One definitive
  form per tower, with no upgrades or simulation changes.

- 2026-10-02: designed the Necromancer lore portrait: a bent human death-worker
  looms over two returning squires, one hauling a knee beneath himself and
  the other braced on a planted sword. Skull-hung bone crook, sallow face,
  long hooked fingers, wine-lined robe and pale green summoning threads frame
  tarnished guild armor and pointed red shields. Reviewed the large and
  compact compositions; preserved the lit eye in its dark socket and quieted
  the ground light. `tdef necromancer` opens the static study; tab from
  Trebuchet reaches it, then wraps to Malgrath. No simulation changes.

- 2026-10-02: designed the Paladin lore portrait: a severe closed helm over
  an enormous convex pointed shield, with a flanged steel mace held upright.
  Overlapping shoulder plates, sealed gorget, ivory guild tabard, burgundy
  mantle, articulated gauntlets and planted sabatons establish an immovable
  stance. Dark lozenge seals remain stark against pale shield planes; a
  visible bracing grip connects shield and bearer. Reviewed large and compact
  layouts. `tdef paladin` opens the static study; tab from Necromancer reaches
  it, then wraps to Malgrath. No simulation changes.

- 2026-10-02: fixed gameplay's frozen victory/defeat sequence after portrait
  viewer work stopped advancing the game screen's animation clock. The clock
  now ticks during active, paused and completed gameplay; simulation stays
  stopped after the result. Regression coverage clears the final wave,
  reaches both result screens, and returns to the lair without duplicate
  rewards or further simulation.

- 2026-10-02: moved Unmapped Depths to post-campaign: all four fixed floors
  and the Heart must be won on the selected renown. One progression check
  governs overworld refresh/unseal, pending descent, Quick Play mazes and
  direct `play -maze` launches. Locked previews and entries explain the
  requirement. Existing victory records remain intact.

- 2026-10-02: designed the Rogue lore portrait: a masked human crouches
  forward between opposing hooked blades, with a pointed teal hood, two
  trailing lengths of red scarf, leather cuirass, bent braced knees and soft
  boots. Forward and reverse knife grips remain visible; half-cell cutting
  edges keep the steel continuous in compact terminals. `tdef rogue` opens
  the static study; tab from Paladin reaches it, then wraps to Malgrath.
  No simulation or progression changes.

- 2026-10-02: designed the Mercenary lore portrait: a scarred, bearded guild
  hireling rests both bare hands on a broad planted cleaver. Bowed shoulders,
  dented open helmet, salvaged iron on one side and quilted cloth on the other,
  mismatched leg armor, sewn patches and worn belt coins establish a tired
  practical fighter. `tdef mercenary` opens the static study; Tab from Rogue
  reaches it, then wraps to Malgrath. No simulation or progression changes.

- 2026-10-02: designed the Wizard lore portrait with a crooked oversized hat,
  stooped human profile, forked ivory beard and heavy gold-faced violet robe.
  His free hand cups a bright faceted core within a hollow cyan diamond; the
  other grips a bowed staff carrying a suspended crystal. A clasped grimoire
  hangs at his hip. `tdef wizard` opens the static study; Tab from Mercenary
  reaches it, then wraps to Malgrath. No simulation or progression changes.

- 2026-10-02: completed the guild portrait roster with Centurion and Squire.
  Centurion has a transverse red horsehair crest, narrow human helmet opening,
  segmented iron cuirass, military skirt and planted sandals. His immense
  flat-topped bronze-bound shield has a square boss; a broken cyan ward arch
  floats outside its rim, with a short gladius exposed on the other side.
  Squire leans forward behind a huge two-handed blade in an oversized tilted
  kettle helm, loose mail, short red surcoat and rolled boots, with a borrowed
  round shield strapped to his back. `tdef centurion` and `tdef squire` open
  the static studies; the cycle now ends Wizard → Centurion → Squire → Malgrath.
  All seven towers and eight enemy types have one definitive lore portrait.
  No simulation, upgrade or progression changes.
