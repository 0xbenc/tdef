# tdef

A tower defense game that runs entirely in your terminal. Enemies walk a
fixed path; you place towers on the grass to stop them before they reach
the exit. Survive all 20 waves.

Built in Go with **zero external dependencies** — the engine, renderer,
and terminal driver are all standard library. The engine is fully
deterministic, so it can be simulated headlessly (used for balance
tuning) or driven by the interactive TUI.

## Build & play

```sh
go build -o tdef .

./tdef                           # title screen -> main menu -> level select
./tdef play -level canyon        # jump straight in: hub | winding | garden | canyon
./tdef play -maze 123            # procedural maze (seeded, reproducible)
./tdef play -diff easy           # easy | normal | hard
./tdef play -diff hard -maze 7
```

Bare `./tdef` (or `play` with no level) starts at the animated title
screen: **main menu** (Start / Help / High Scores / Quit) → **level
select** with the built-in maps, a procedural maze row (type a seed,
empty = random), and a difficulty pick. Tall terminals get a live map
preview. Explicit `-level`/`-maze` flags skip the menus and start the
game directly; `-diff` preselects the difficulty.

Requires a real TTY (raw mode, alternate screen, mouse).

## Controls

| key | action |
|-----|--------|
| `1`-`7` | select tower (Gunner, Cannon, Frost, Sniper, Tesla, Mortar, Flak) |
| arrows / `wasd` | move cursor |
| `enter` / click | place tower on grass |
| `u` | upgrade the selected tower (on a tower) |
| `x` | sell the selected tower (70% refund) |
| `t` | cycle tower targeting: first / strongest / closest (on a tower) |
| `n` | start the next wave early (bonus gold) |
| `p` | pause |
| `f` | cycle speed 1x / 2x / 4x |
| `h` | toggle help |
| `r` | restart (on game over) |
| `esc` | cancel placement; back out of menus (not in-game) |
| `q` | quit (any screen) |

In the menus: `enter` activates, `1`-`4` pick a menu item, and the mouse
works too (click items, wheel scrolls). Mouse in-game: click a tower in
the menu, then click the map; scroll wheel cycles speed (1x/2x/4x, like
`f`).

## Towers

| # | tower | cost | role |
|---|-------|------|------|
| 1 | Gunner | 50 | cheap, fast, single-target |
| 2 | Cannon | 100 | slow, splash damage |
| 3 | Frost | 75 | slows enemies in range |
| 4 | Sniper | 150 | long range, high single damage |
| 5 | Tesla | 200 | chains lightning to nearby enemies |
| 6 | Mortar | 250 | very slow, huge splash, long range — crowd nuker |
| 7 | Flak | 80 | very fast, short range — shreds fast enemies |

Each tower has 3 levels; upgrading increases damage/range.

By default a tower attacks the enemy **furthest along the path** (first).
Press `t` on a tower to cycle its targeting priority to **strongest**
(highest HP) or **closest** — useful for focusing tanks or finishing off
stragglers. The current mode shows in the tower info line.

## Enemies & waves

Eight enemy types, introduced one at a time across 20 hand-shaped waves
(Bloons TD 3 style — a new type debuts in a legible near-solo wave, then
gets mixed in):

| type | role | debuts |
|------|------|--------|
| Minion | basic | wave 1 |
| Runner | fast | wave 3 |
| Grunt | tough | wave 6 |
| Wisp | very fast, weak | wave 11 |
| Splitter | splits into two Minions | wave 12 |
| Tank | very tough, slow | wave 13 |
| Shield | armored (−40% damage) | wave 14 |
| Boss | high HP; leaking one costs 6 lives | wave 15, 18, 20 (two) |

Each wave has a theme (`warmup`, `runner`, `gauntlet`, `boss`, `finale`, …)
shown in the HUD, and a one-line **telegraph** warns you before a
mechanically new wave (e.g. "Beware the Boss — leaking it costs 6 lives").
Early waves are light with breather dips at each introduction; the last
five waves (W16–W20) are a back-loaded climax. Enemy HP and speed scale
with the wave number; the boss gets stronger each time it appears.

### Pacing

The inter-wave break tapers from 11s (after wave 1) down to ~4.5s (after
wave 19), so tempo rises as the game escalates. Starting a wave early
(`n`) pays a small bonus gold — a tempo choice, not a dominant income
source. Clearing a wave pays a bonus that grows with the wave number.

## Maps

Four hand-crafted 45×13 maps form a difficulty ladder — `hub` (easiest,
open central pocket), `winding`, `garden`, `canyon` (tightest). Each has
a different path layout and buildable pockets, so the right tower
placement differs per map. `-maze` generates a random-but-solvable maze
from a seed; the generator biases toward long, snake-like chokepoints.

## Rendering

The playfield auto-scales to your terminal at startup (1×–4×), so a big
window shows a bigger board. If the window is shrunk below the board, an
"enlarge your terminal" notice is shown instead of a clipped frame (the
game pauses while it is up; `p` resumes). When
a game ends, a stats box summarizes the run (waves, kills, leaks, towers,
score, best combo, time) and your best score for that map.

## Balance / headless tools

The engine is deterministic, so it can be driven without a terminal:

```sh
./tdef bench -n 40 -maze 60   # run autoplay (greedy AI) games, report stats
./tdef headless -level garden # simulate one game, print result
./tdef capture -level canyon -text -out /tmp/frames   # dump frames
./tdef maps                   # list built-in maps
```

`bench` plays many games with a simple greedy AI and reports win rate,
average wave reached, leaks, and towers built per map — used to tune the
per-map HP multipliers and economy so that a competent player can win
while a passive one loses.

High scores are stored per-map in `~/.tdef-hiscores.json`.

## Layout

```
main.go            CLI (play / bench / headless / capture / maps)
game/              pure, deterministic engine (no I/O)
  balance.go       tower/enemy specs, wave scaling, difficulty
  mapgen.go        procedural maze generator
  sim.go           autoplay AI + headless simulation
  levels/          built-in .txt maps
render/            pure frame model + ANSI/text exporters
tui/               terminal driver, input parser, app loop
hiscore/           persistent high scores
```
