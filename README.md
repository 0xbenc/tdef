# tdef

> — terminal tower defense —

You are **Grak, the Last Monster**. The Dragon **Malgrath** is dying at the
bottom of his lair, and the guild has marked the hoard: **twenty
expeditions** are marching in to put the beast down and take the gold. You
took the hoard. Now hold the lair against all twenty — or watch the heart
go cold.

`tdef` is a tower defense game that runs entirely in your terminal. No
engine, no assets, no dependencies — Go's standard library and a dying
dragon.

## Build & play

```sh
go build -o tdef .

./tdef                       # title screen → main menu → the lair
./tdef overworld             # jump straight to the lair map
./tdef play -level canyon    # skip the menus: hub | winding | garden | canyon
./tdef play -maze 123        # a procedural maze (seeded, reproducible)
./tdef play -diff hard -maze 7
```

Needs a real TTY (raw mode, alternate screen, mouse).

## Holding the lair

Pick a tower with `1`–`7`, move the cursor with the arrows (or `wasd`), and
place it on the grass with `enter`. Heroes walk a fixed path toward the
heart; every one that gets through strikes the dying Dragon. If his HP
reaches zero the lair falls. Hold all twenty expeditions and Malgrath
endures.

| key | action |
|-----|--------|
| `1`–`7` | pick a tower |
| arrows / `wasd` | move the cursor |
| `enter` / click | place a tower |
| `u` / `x` | upgrade / sell the selected tower (70% refund) |
| `t` | cycle targeting: first → strongest → closest |
| `n` | call the next wave early (bonus gold) |
| `p` / `f` | pause / cycle speed 1× – 2× – 4× |
| `h` / `esc` / `q` | help / back out / quit |

The **lair** (overworld) is a map you walk. Hold a floor and it unseals;
step onto it and press `enter` to descend into its defense, and the result
comes back to the map. The four built-in floors are a difficulty ladder —
**the Rotunda** (easiest), **the Long Halls**, **the Sunken Garden**,
**the Rift** (tightest) — and holding all four opens **the Heart**, the
endgame chamber behind the Rotunda. `tab` cycles the renown (easy / normal /
hard); the lair remembers what you've held, and a held floor earns a
**dragon heart** (+1 starting life on later defenses).

## Your kin (towers)

The last monster kin you command. Each has three levels; upgrading deepens
the defense.

| # | tower | cost | role |
|---|-------|------|------|
| 1 | Orc Gunner | 50 | cheap, fast, single target |
| 2 | Cannonier | 100 | slow, splash |
| 3 | Frost Mage | 75 | slows enemies in range |
| 4 | Ranger | 150 | long range, hard hits |
| 5 | Lightning Mage | 200 | chains lightning to nearby heroes |
| 6 | Trebuchet | 250 | huge splash, very slow — the crowd nuke |
| 7 | Gnoll Slingers | 80 | very fast, short range — shreds the quick |

By default a tower fires on the hero furthest along the path; `t` cycles it
to strongest or closest.

## The guild (heroes)

Eight kinds of hero, introduced one at a time across the twenty expeditions
then mixed in. Squires and rogues at first; necromancers (who rise again
when they fall), warded centurions, and paladins later. On waves 15, 18, and
20 the guild's champion arrives — **the Player**. Leaking one costs six of
the Dragon's heart.

## Under the hood

- **Zero dependencies.** The engine, renderer, and terminal driver are all
  Go's standard library. No `go.sum`, no `vendor/`, no CGO.

- **Deterministic engine.** `game/` is pure — no I/O, no wall clock. The
  whole game is a function of its map (or maze seed) and its input, so the
  same seed and the same moves replay to the same frame. That's what makes
  balance tuning and the headless tools below possible.

- **The terminal is the frame.** The whole screen is one rounded box, edge
  to edge, with the playfield centered inside at the largest integer scale
  (1×–4×) that fits. The HUD lives in the border itself — status segments in
  the top edge, tower slots and the selected tower's stats in the bottom.
  Resize any time and the board re-scales on the next frame; below 62×19 the
  game pauses and asks you to enlarge the terminal.

- **Headless.** Drive the engine without a terminal:

  ```sh
  ./tdef bench -n 40 -maze 60        # 40 autoplay games per level + 60 mazes
  ./tdef headless -level garden      # simulate one game, print the result
  ./tdef capture -level canyon -text # dump rendered frames to files
  ./tdef maps                        # list the built-in maps
  ```

  `bench` plays many games with a simple greedy AI and reports win rate,
  average wave, leaks, and towers built per map — the tool behind the
  per-map HP and economy tuning. `capture -scale 1–4` pins the render to a
  fixed virtual terminal size, so frame dumps are byte-identical in any
  environment.

- **Layout.**

  ```
  main.go            CLI (play / bench / headless / capture / overworld / maps)
  game/              pure, deterministic engine (no I/O)
  render/            pure frame model + ANSI/text exporters
  tui/               terminal driver, input parser, app loop
  hiscore/           persistent high scores + the lair's memory
  ```

  High scores live in `~/.tdef-hiscores.json`; the lair's memory (floors
  held, hearts, relics) in `~/.tdef-lair.json`.

## License

MIT — see [LICENSE](LICENSE).
