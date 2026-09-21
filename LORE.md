# LORE — "The Last Monster"

Working theming for tdef. No code changes yet — this is the source of truth for the
text/skin pass. `OPEN` markers need a decision.

## Premise

You are **Grak**, the Last Monster — the last survivor of a dead villain team, its
hero. You find **Malgrath**, the Dragon, your former ally and the team's big bad,
dying in the deep of his lair beside his hoard.

You take the hoard (starting gold) and fortify the lair with what monster kin
remain. The heroes' guild has marked the lair: **twenty expeditions** are coming,
and they want the Dragon dead and the hoard.

- **Lives = the Dragon's HP.** Every hero who reaches the heart strikes the dying
  Dragon. A Boss that leaks hits it for 6.
- **Lose:** the Dragon's HP reaches 0 — the Dragon dies, the lair falls, the hoard
  is taken.
- **Win:** all twenty expeditions turn back. The guild's hunt fails; the Dragon
  endures.

## Malgrath, the Dragon

- The villain team's big bad; ancient and already dying when you arrive.
- The hoard is your starting gold: "They came to take the Dragon's treasure —
  take theirs first."

## Grak, the Last Monster (you)

- The team's last survivor; guardian of the lair and the dying Dragon.
- The title signature `by 0xbenc` = the Last Monster's mark.

## Defenders (towers)

The last monster kin you command. Role and stats unchanged; names only.

| key | old   | new            | role                          |
|-----|-------|----------------|-------------------------------|
| 1   | Gunner| Orc Gunner     | cheap, fast, single target    |
| 2   | Cannon| Cannonier      | slow, splash                  |
| 3   | Frost | Frost Mage     | slows enemies in range        |
| 4   | Sniper| Ranger         | long range, high damage       |
| 5   | Tesla | Lightning Mage | chains lightning              |
| 6   | Mortar| Trebuchet      | very slow, huge splash, long range |
| 7   | Flak  | Gnoll Slingers | very fast, short range, shreds fast enemies |

- Upgrade = the defense deepens (the pit widens, the cold spreads, the ward holds).

## Heroes (enemies)

The guild's expeditions. Glyphs, stats and debut waves unchanged; names only.

| glyph | old      | new         | role                                    |
|-------|----------|-------------|-----------------------------------------|
| `o`   | Minion   | Squire      | basic                                   |
| `r`   | Runner   | Rogue       | fast                                    |
| `g`   | Grunt    | Mercenary   | tough                                   |
| `w`   | Wisp     | Wizard      | very fast, weak                         |
| `s`   | Splitter | Necromancer | on death, two Squires rise              |
| `t`   | Tank     | Paladin     | very tough, slow                        |
| `D`   | Shield   | Centurion   | warded — cuts incoming damage 40%       |
| `B`   | Boss     | **The Player** | the guild's champion; leaking one costs 6 |

- "The Player" is deliberate: the hero is the player; the lair is the game.

## World mapping

| game element    | lore                                                        |
|-----------------|-------------------------------------------------------------|
| the board       | the Dragon's lair                                           |
| the path        | the approach to the heart                                   |
| waves           | expeditions (20)                                            |
| gold            | the fallen heroes' loot, taken from the dead                |
| lives           | the Dragon's HP                                             |
| breach/leak     | a hero has reached the heart — the Dragon is struck         |
| high score      | the lair that held best                                     |
| difficulty      | the lair's renown — Easy: a forgotten lair (weak expeditions); Hard: a legendary hoard (strong expeditions) |
| boot sequence   | the lair waking; the defenses coming online                 |
| attract demo    | the siege, looping                                          |

### Maps = floors of the lair

| old     | new             | note                              |
|---------|-----------------|-----------------------------------|
| hub     | the Rotunda     | central lair hall, easiest        |
| winding | the Long Halls  |                                   |
| garden  | the Sunken Garden | where a temple once stood       |
| canyon  | the Rift        | tightest                          |
| maze    | the Unmapped Depths | seeded, uncharted            |

## Draft copy (pending sign-off, for the code pass)

### Telegraphs (pre-expedition rumors)

| wave | draft |
|------|-------|
| 3    | "The guild has sent rogues." |
| 6    | "Mercenaries — and they hit harder." |
| 7    | "A swift raid. Loose arrows fly true." |
| 11   | "Wizards — fast and frail, but deadly." |
| 12   | "A necromancer walks among them. They rise when he falls." |
| 13   | "A paladin leads the vanguard." |
| 14   | "Centurions — warded, they shrug off your blows." |
| 15   | "The Player has set out. If they reach the heart, it ends." |
| 18   | "The Player returns, stronger." |
| 20   | "The final expedition. Two Players. Hold the lair." |

### HUD / screens

- Lives: shown as `♥ N` (Malgrath's heart) — no text label, the header
  stats segment is never elided and has no room for one. The help
  subtitle carries the mapping.
- Wave themes (decided — debut waves named for the new threat, the rest as
  expedition types): warmup→**scouts**, runner→**raid**, mixed→**column**,
  grunt→**assault**, wisp→**the coven**, splitter→**the risen**, tank→
  **the vanguard**, shield→**the wall**, boss→**the player**, gauntlet→
  **siege**, finale→**the end**.
- Help = "the Last Monster's ledger — how to hold the lair against twenty expeditions."
- Title demo label `IDLE START SCREEN` → `THE SIEGE` (box width adjusted in code).
- Game over: "Malgrath has fallen. The lair is clean."
- Victory: "The lair is held. Malgrath endures." (shortened from "The
  twentieth expedition turns back. Malgrath endures." to fit the 44-col
  box interior)
- Tagline `— terminal tower defense —`: kept as the meta layer (the lair, rendered
  in a terminal). TDEF stays the product name. (decided)

## Decided

- Mortar = **Trebuchet**; Cannon = **Cannonier** (spelling fixed).
- Flak = **Gnoll Slingers**.
- Boss = **The Player**.
- Dragon = **Malgrath**; Last Monster = **Grak**.
- Tagline stays `— terminal tower defense —` (meta layer; TDEF = product name).
- Wave themes reskinned (mapping above).

No open decisions — LORE.md is complete; next step is the code pass.
