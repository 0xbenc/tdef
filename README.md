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
./tdef journal               # browse Grak's lore journal
./tdef dragon                # view the static Malgrath portrait mockup
./tdef grak                  # view Grak, the Last Monster
./tdef gnolls                # view the Gnoll Slingers lore portrait
./tdef gunner                # view the Orc Gunner lore portrait
./tdef frost                 # view the Frost Mage lore portrait
./tdef player                # view the guild champion lore portrait
./tdef cannonier             # view the monster artillery crew portrait
./tdef ranger                # view the monster archer portrait
./tdef lightning             # view the Lightning Mage portrait
./tdef trebuchet             # view the monster siege crew portrait
./tdef necromancer           # view the guild summoner portrait
./tdef paladin               # view the guild vanguard portrait
./tdef rogue                 # view the guild knife-runner portrait
./tdef mercenary             # view the weary guild hireling portrait
./tdef wizard                # view the guild spellcaster portrait
./tdef centurion             # view the warded guild veteran portrait
./tdef squire                # view the young guild recruit portrait
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

The **lair** (overworld) is a wide cavern you walk, with a scrolling view
that follows Grak along stone walkways and bridges, past Malgrath's hoard,
vaulted halls, and a sunken temple. The crossing forks into an upper stair
and a lower waterside approach. The ledger and controls stay in place as
the world moves. Sealed rooms block entry at the doorway until they open.
Hold a floor and it unseals; step onto it and press `enter` to descend into
its defense, and the result comes back to the map. The four built-in floors are a difficulty ladder —
**the Rotunda** (easiest), **the Long Halls**, **the Sunken Garden**,
**the Rift** (tightest) — and holding all four opens **the Heart**, the
endgame chamber behind the Rotunda. `tab` cycles the renown (easy / normal /
hard); the lair remembers what you've held, and a held floor earns a
**dragon heart** (+1 starting life on later defenses).

A fresh lair starts Grak in **the Rotunda**, the only open floor. Holding it
opens **the Rift**; holding the Rift opens **the Long Halls**, then holding the
Halls opens **the Sunken Garden**. Beat all four floors and then **the Heart**
to unlock **the Unmapped Depths** on that renown. This also unlocks procedural
play from Quick Play and the command line.

Press `v` in the Rotunda to view a larger, static portrait of Malgrath in his
curled resting pose. `esc` or `enter` returns to the lair. This is a visual
mockup for a future cutscene. Press `g` in the Rotunda for Grak’s companion
portrait. `tab` cycles Malgrath → Grak → Gnoll Slingers → Orc Gunner → Frost Mage → the Player → Cannonier → Ranger → Lightning Mage → Trebuchet → Necromancer → Paladin → Rogue → Mercenary → Wizard → Centurion → Squire. The slingers have
one definitive lore portrait: a throwing gnoll and a crouched partner
supplying stones. Each tower portrait shows one definitive form, with no
upgrade variants. The Ranger holds a full draw on a tall recurved bow, with
a swept ear, narrow eye, feathered quiver and wind-torn cloak. The Lightning
Mage catches a descending fork in one claw and directs it with the other,
with swept horns and a flaring violet robe. The Trebuchet exposes its loaded
sling, hanging counterweight and windlass inside a timber frame, with one
monster winding the crank and another sighting from the rear rail. The
Necromancer looms over two returning armored squires beneath a skull-hung
crook and an extended skeletal hand. The Paladin braces an enormous pointed
shield beneath a closed helm, with a flanged mace held upright, overlapping
plate armor and an ivory guild tabard. The Rogue crouches between two hooked
blades, with a pointed hood, masked human face, light leather armor and two
long trailing lengths of red scarf. The Mercenary rests both bare hands on a
planted cleaver, wearing a dented open helmet, salvaged iron on one shoulder
and patched cloth on the other; a scar and short beard frame a weary stare.
The Wizard wears an enormous crooked hat and gold-faced violet robes, with
a forked ivory beard, clasped field grimoire, bowed staff and suspended crystal.
His free hand cups a luminous core inside a hollow diamond spell. The Centurion
stands behind a bronze-bound rectangular shield beneath a broad red helmet
crest, with segmented iron armor, short sword and a broken cyan ward arch.
The Squire leans into an enormous two-handed sword, wearing a tilted kettle
helm, loose mail and a short guild surcoat; a round shield hangs on his back.
All seven towers and all eight enemy types now have one definitive portrait.

Press `j` in the lair, or while a defense is paused, to open **Grak's Journal**.
The journal has four chapters: **Towers**, **Enemies**, **Places**, and
**Story Moments**. Artwork uses the page width; lore and field notes sit below
it and scroll independently. `[` and `]` change chapters. Arrow keys browse the
collection; `enter` reads a page. Left/right turn to another discovered page
and up/down scroll its text. `esc` returns to the collection, then to the exact
lair or paused defense you left. `./tdef journal` opens it directly.

Discovery is silent and permanent across defenses and difficulties: first
successful placement unlocks a tower (including a relic's free gunner), first
encounter unlocks an enemy, and visiting a room unlocks its place illustration.
Failed placements and sealed rooms reveal nothing. There are no discovery
popups, unread badges or upgrade variants. Existing victory records backfill
proven visits, encounters and memories without replaying the campaign.

Every defense has a distinct illustrated story moment on each of the three
renowns: eighteen memories across the four floors, Heart and Depths. Any maze
seed can earn the Depths memory for its difficulty. Undiscovered moments show
the required floor and renown when selected, while keeping their title and
artwork hidden. Four bonus **Afterward** memories unlock at 1, 3, 5 and 10
combined Heart or Depths victories, beyond the main story. Each finished
victory counts once; defeats earn no victory memory. Older saves retain a
minimum repeat count supported by their distinct held endgame records.

Each defense has its own route and build terrain. Inner bends can cover several
passes; open stretches reward long range, while narrow ledges demand careful
placement. Roomier floors fund more towers through starting gold and rewards.
Enemy health is tuned separately for each route.

| floor | layout | build tiles | starting gold (normal) |
|-------|--------|-------------|------------------------|
| Rotunda | coiling route around the hoard, broad inner courts | 231 | 220 |
| Long Halls | staggered chambers, pillars and returning corridors | 234 | 275 |
| Sunken Garden | looping clearings around planted islands | 230 | 308 |
| Rift | cliffside switchbacks and narrow firing ledges | 170 | 286 |
| Heart | long inward spiral with shared firing lanes | 189 | 308 |

Blocked terrain carries the scenery: dragon reliefs and treasure vaults in the
Rotunda, vaulted masonry and banners in the Halls, trees and drowned arches in
the Garden, basalt and flowing lava in the Rift, bone masks and chained ribs in
the Heart, and crystal seams in the Depths. Larger display scales reveal more
detail; roads and tower pads remain clear.

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
to strongest or closest. The Lightning Mage needs only its first target
inside that range: each bounce searches within 2.6 cells of the last hero hit,
so a chain can reach beyond the mage's range.

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
  held, hearts, relics) in `~/.tdef-lair.json`; journal discoveries and endgame wins in
  `~/.tdef-journal.json`.

## License

MIT — see [LICENSE](LICENSE).
