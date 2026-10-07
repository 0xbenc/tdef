# Fixed-floor tactics and guild counters

The five authored maps teach different placement and composition choices. Walls restrict tower positions; attacks still travel over walls.

| Floor | Placement problem | Enemy pressure |
|---|---|---|
| Rotunda | Inner courts cover several passes with one compact defense | Original formations; forgiving opening |
| Rift | Isolated banks require spreading investment | Raiding parties separated by short regrouping gaps |
| Long Halls | Long returning lanes give Rangers sustained exposure | Spaced columns, fewer bodies in the largest groups |
| Sunken Garden | Short hairpins reward slow fields and splash overlap | Packed raids, larger crowds |
| Heart | Outer interception plus an inner finishing defense | Late rushing fronts followed by plate; separated bosses |

## Guild counters

| Fighter | Defense | Preferred response | Alternatives |
|---|---|---|---|
| Paladin | Plate reduces gun and sling damage by 75%, arrows by 40% | Cannonier or Trebuchet | Frost, lightning and Runeforge bypass plate |
| Centurion | Ward reduces frost, lightning and Runeforge damage by 80% | Ranger | Guns, slings and siege deal full damage; frost still slows |
| Rogue | Evasion reduces gun, arrow and sling damage by 50% while unfrosted | Frost with a damage tower | Splash, lightning and Runeforge bypass evasion |

Evasion uses deterministic damage reduction, not random misses. It returns as soon as frost expires. Splash bypass applies to the actual explosion, including enemies beside its target. Rogues have 24 base health, reduced from 34 to keep their early introductions manageable with the new evasion. Other enemy health, speeds, costs, waves, rewards and progression are unchanged.

Paladins use brass background pads, Centurions blue, and evasive Rogues purple. Frosted Rogues lose the purple pad and use the existing cyan slow indicator. Debut and raid warnings explain the counters; bestiary and defender tactical notes describe the same rules.

## Measured comparisons

`go run ./cmd/tactics` runs deterministic builders, not human players. Focused policies share two affordable starting Gunners, weighted route coverage, a twelve-gold reserve and a 24-tower cap. They buy their first counter before upgrading the foundation: frost plans open with a Frost Mage, while siege plans open with an affordable Cannonier before saving for Trebuchets. Pure gun, slinger and Ranger plans retain their composition weaknesses. One-camp additionally limits placement to a radius of 3.2 cells around its best coverage site. Greedy uses the original autoplay and does not adapt to the new traits.

`mixed-counters` cycles Frost, Cannonier, Ranger and Trebuchet after the opening. Focused policies remain simple fixed builders; losing runs can reflect placement, spending or composition.

Normal difficulty, without campaign hearts or relics:

| Floor | Rangers | Frost + siege | Mixed counters | One compact camp |
|---|---|---|---|---|
| Rotunda | Held, 17 lives | Held, 20 lives | Held, 20 lives | Held, 20 lives |
| Rift | Held, 14 lives | Held, 20 lives | Held, 20 lives | Lost on wave 20 |
| Long Halls | Held, 20 lives | Held, 20 lives | Held, 20 lives | Held, 8 lives |
| Sunken Garden | Lost on wave 5 | Held, 9 lives | Held, 18 lives | Lost on wave 16 |
| Heart | Held, 20 lives | Held, 18 lives | Held, 20 lives | Lost on wave 18 |

Hard campaign checks include only hearts earned on preceding floors, without relics. Verified alternatives include:

| Floor | Successful tested policies |
|---|---|
| Rotunda | siege, frost-siege, ranged-siege, mixed-counters, one-camp |
| Rift | frost-siege, ranged-siege, mixed-counters |
| Long Halls | rangers, frost-rangers, frost-siege, ranged-siege, mixed-counters |
| Sunken Garden | siege, frost-siege, ranged-siege |
| Heart | ranged-siege, mixed-counters |

The complete pass covers twelve policies across five floors on normal, easy campaign, hard campaign and a neutral control: 240 runs, without simulation timeouts. Every authored floor has multiple successful tested compositions on each difficulty. The neutral control shares the original waves, 0.9 health multiplier and 1.25 gold multiplier; its compact camp still holds the Rotunda and loses in the Rift and Heart. Long Halls still gives the Ranger policy better exposure than the Garden.

## Playability checks and limits

Simulation tests cover each damage matchup, frost application and expiry, splash impact, unchanged ordinary targets, and multiple viable Hard campaign compositions. Presentation tests check the trait backgrounds and frost tell. These establish mechanical behavior and practical build options, not human win rates. The next useful playtest is a fresh campaign through the Rogue introduction and the first Paladin/Centurion waves, checking whether the warnings arrive early enough and the responses feel worthwhile.


## Specialists

Each of the four specialists is limited to one placed per level. Selling one frees its slot; upgrades do not consume another slot. Original defenders remain unrestricted.

The second roster page uses `[ ]` and local keys `1` to `4`. A page change cancels placement and preserves the selected defender. `r` turns a Runeforge preview or a selected forge through four directions. The bottom-border pager also accepts clicks. Existing defender keys and journal discovery IDs stay stable.

| Defender | Job | Limits and useful partners |
|---|---|---|
| Runeforge | Continuous magic damage to every fighter on a straight ray; 22/34/50 damage per second, 7/8/9 tile reach | Frost extends exposure. Plate does not stop it; Centurion wards do. Corners and poor alignment waste its range. Three Long Halls alcoves provide firing positions without changing the road. |
| Ogre Hookmaster | Pulls fighters backward along their actual route, into an existing kill zone | Shared 2.5 second recovery; heavy fighters take 40% displacement, bosses 10%. Lifetime pull allowances are 6 tiles for light fighters, 2 for heavy fighters and 0.5 for bosses. This prevents indefinite trapping even with frost and several Ogres. |
| Goblin Sappers | Automatically prepare nearby road tiles with 90/145/220 damage splash charges | Stock caps 3/4/5; 0.6 second arming delay; global spacing keeps charges apart. Recruits can waste them. Mines survive wave transitions, disappear with their crew, and bypass plate and evasion. |
| Hex Witch | Amplifies damage taken by one marked target by 25/35/45%, for 3/3.5/4 seconds | Defaults to strongest targeting. Marks never stack; a weaker witch cannot replace a stronger active mark. Damage still passes through the target's resistance. Pair with damage dealers rather than buying a wall of witches. |

Costs and upgrades are in `game/balance.go`. Mines retain their planting-time damage when the crew upgrades. Pulls respect corners and never place a fighter off the road. Fast fighters cannot skip a mine between simulation ticks, including on their final step toward the exit. Beam damage uses elapsed seconds, so faster rendering does not increase damage.

### Measured specialist pass before the one-per-level limit

The expanded probe uses real directional beam coverage and aims each forge. Support placement values nearby damage coverage; the current probe respects the one-per-level limit and buys conventional damage when a specialist is already placed. The recorded reports below predate that limit and allowed duplicates. The ordinary greedy policy remains a rough damage-cost heuristic and does not value hex support. These are deterministic policies, not optimized human play.

Seventeen policies across all five authored floors on Easy, Normal and Hard campaign produced **255 completed runs with no timeouts**. Reports are saved in `artifacts/tactics-specialists-{easy,normal,hard}.json`.

| Composition | Normal campaign wins | Hard campaign wins |
|---|---|---|
| Frost, Runeforge, Ranger | Rotunda, Long Halls, Garden, Heart | Rotunda |
| Cannonier, Hookmaster, Ranger | All five floors | Rotunda, Long Halls |
| Sappers, Ranger, Frost | Rotunda, Rift, Long Halls, Heart | Rift, Long Halls |
| Cannonier, Witch, Ranger | All five floors | Rotunda, Long Halls |
| All four specialists with Frost and Ranger | All five floors | Rotunda |

Every specialist has a tested Hard campaign niche. The combined specialist policy clears every Normal floor, but struggles on Hard outside the Rotunda. Existing siege and mixed-counter alternatives still clear each Hard floor. The Rift challenges beam alignment and magic damage; the Garden can overwhelm a slow, expensive opening. Buying every specialist is therefore an option, rather than a universal solution.

Combat tests cover beam timing and rotation, hex strength and expiry, shared hook recovery and total pull limits, mine arming, stock, swept crossings and sale cleanup. Roster tests cover small terminals, keyboard and mouse paging, preserved selection and preview direction. Journal checks cover all eleven entries at five terminal sizes. The specialist-heavy 24-defender/150-fighter simulation-and-render benchmark measured approximately **0.57 ms/frame** at 120 by 40 on the development machine; terminal output costs vary with the terminal.


## First two guided defenses

New campaigns start the Rotunda with Gunners. Recruits arrive after clearing these waves, before the next expedition:

| First defense | Clear wave |
|---|---|
| Frost Mage | 1 |
| Cannonier | 2 |
| Ranger | 4 |
| Gnoll Slingers | 5 |
| Lightning Mage | 7 |
| Trebuchet | 9 |

The second campaign defense keeps all seven main defenders. Runeforge joins before its first wave; Hookmaster joins after wave 2, Sappers after 4, and Witch after 6. Each specialist retains its one-per-level placement limit.

Every recruit after the initial Gunners grants its base purchase price in gold, once. The lesson freezes the simulation; acknowledging it opens a placement preview on useful ground and leaves planning paused. The player chooses the final position. `n` starts the next wave and resumes play; `p` can run preparations during a break without sending the next wave. Counter warnings remain available during an extended planning break. No enemy, health, damage or cost nerfs are used.

Unlocks save immediately, survive defeat and restarts, and carry between difficulties. Completing the Rotunda advances to specialist training; completing that defense finishes training. Existing campaign records bypass roster locks, while Quick Play remains unrestricted. Journal recruitment pages become available as units join. `termtd tutorial` and `termtd tutorial -stage 2` provide save-free replays; replay progress does not carry into the regular campaign.

The deterministic balance probe accepts `-training 1` (first Rotunda) or `-training 2 -floor canyon -campaign` (second defense in the Rift). Regression checks require two different viable compositions in each guided defense on Easy, Normal and Hard. A broader Hard Rotunda pass cleared with 14 of 17 strategies; Normal cleared with 15 of 17. Reports are in `artifacts/ftue-rotunda-{normal,hard}.json` and `artifacts/ftue-rift-hard.json`. These policies check practical options, not human success rates.

The gold HUD regression also replays incremental terminal updates with symbols rendered at two columns, checking the visible balance as digit lengths change. The blitter anchors the next cell after potentially variable-width symbols so HUD digits cannot drift; box and block artwork retains its contiguous rendering path.
