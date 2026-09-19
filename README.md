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

./tdef play                      # play (default: winding, normal)
./tdef play -level canyon        # pick a map: winding | garden | canyon
./tdef play -maze -seed 123      # procedural maze (seeded, reproducible)
./tdef play -diff easy           # easy | normal | hard
./tdef play -diff hard -maze -seed 7
```

Requires a real TTY (raw mode, alternate screen, mouse).

## Controls

| key | action |
|-----|--------|
| `1`-`5` | select tower (Gunner, Cannon, Frost, Sniper, Tesla) |
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
| `q` / `esc` | quit |

Mouse works too: click a tower in the menu, then click the map; scroll
wheel changes speed.

## Towers

| # | tower | cost | role |
|---|-------|------|------|
| 1 | Gunner | 50 | cheap, fast, single-target |
| 2 | Cannon | 100 | slow, splash damage |
| 3 | Frost | 75 | slows enemies in range |
| 4 | Sniper | 150 | long range, high single damage |
| 5 | Tesla | 200 | chains lightning to nearby enemies |

Each tower has 3 levels; upgrading increases damage/range.

By default a tower attacks the enemy **furthest along the path** (first).
Press `t` on a tower to cycle its targeting priority to **strongest**
(highest HP) or **closest** — useful for focusing tanks or finishing off
stragglers. The current mode shows in the tower info line.

## Enemies

Minion, Runner (fast), Grunt (tough), Tank (very tough), Splitter
(splits into two minions on death), and a Boss every 5 waves (two on the
final wave). Enemy HP and speed scale with the wave number; the boss gets
stronger each time it appears.

## Maps

Three hand-crafted 45×13 maps form a difficulty ladder — `winding`
(easiest), `garden`, `canyon` (tightest). Each has a different path
layout and buildable pockets, so the right tower placement differs per
map. `-maze` generates a random-but-solvable maze from a seed.

## Balance / headless tools

The engine is deterministic, so it can be driven without a terminal:

```sh
./tdef bench -n 40 -maze 60   # run autoplay (greedy AI) games, report stats
./tdef headless -level garden -seed 5   # simulate one game, print result
./tdef capture -level canyon -seed 5 -text -out /tmp/frames   # dump frames
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
