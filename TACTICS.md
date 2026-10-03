# Fixed-floor tactical redesign

The five authored maps now teach different placement and composition choices.
Walls restrict tower positions; attacks still travel over walls. No new line-of-sight rule is implied.

| Floor | Placement problem | Enemy pressure |
|---|---|---|
| Rotunda | Inner courts cover several passes with one compact defense | Original formations; forgiving opening |
| Rift | Isolated banks require spreading investment | Raiding parties separated by short regrouping gaps |
| Long Halls | Long returning lanes give Rangers sustained exposure | Spaced columns, fewer bodies in the largest groups |
| Sunken Garden | Short hairpins reward slow fields and splash overlap | Packed raids, larger crowds |
| Heart | Outer interception plus an inner finishing defense | Late rushing fronts followed by armor; separated bosses |

Tower costs and abilities are shared. Floor health and gold were tuned against actual exposure and upgraded defenses. Enemy debuts and boss milestones remain unchanged, and Unmapped Depths keeps its original waves and generation.

## Measured comparisons

`go run ./cmd/tactics` runs deterministic builders, not human players. Focused policies share two affordable starting Gunners, foundation upgrades, weighted route coverage, a twelve-gold reserve and a 24-tower cap. One-camp additionally limits placement to a radius of 3.2 cells around its best coverage site. It therefore also has less spending capacity; this comparison measures the practical placement restriction, not equal spending. The greedy policy uses the original autoplay instead.

Normal difficulty, without campaign hearts or relics:

| Floor | Rangers | Frost + siege | One compact camp |
|---|---|---|---|
| Rotunda | Held, 20 lives | Held, 20 lives | Held, 20 lives |
| Rift | Held, 13 lives | Held, 20 lives | Lost on wave 20 |
| Long Halls | Held, 20 lives | Held, 20 lives | Held, 8 lives |
| Sunken Garden | Lost on wave 7 | Held, 20 lives | Lost on wave 17 |
| Heart | Held, 20 lives | Held, 20 lives | Lost on wave 18 |

The shape-only control (`-neutral`) gives every map the same original waves, 0.9 enemy health multiplier and 1.25 gold multiplier. The compact camp holds the Rotunda but loses in the Rift and Heart. Rangers hold the Halls without leaks, while the same policy leaks in the Garden. Geometry therefore changes outcomes before the authored formations are added.

Hard campaign comparisons include only the hearts earned on preceding floors; no relics:

| Floor | Successful policies |
|---|---|
| Rotunda | rangers, frost-chain, siege, frost-siege, ranged-siege, one-camp |
| Rift | frost-siege, ranged-siege |
| Long Halls | rangers, frost-rangers, ranged-siege |
| Sunken Garden | siege, frost-siege |
| Heart | rangers, siege, frost-siege, ranged-siege |

The complete check covers 55 runs each on normal, easy campaign, hard campaign and the neutral control: 220 runs without simulation timeouts. Every authored floor has multiple successful tested compositions on each difficulty. Selected tactical counterexamples and hard alternatives are regression tests in `cmd/tactics`.

## Playability checks and limits

Cheap opening placements engage within eight simulation seconds in the tested focused policies. Build tiles can reach the road; every painted road belongs to one unbranched route, with no hairpin shortcuts. Initial hints, accurate previews and pre-wave warnings make the intended problems legible without modal interruptions. Early enemy introductions remain familiar. Several damage compositions work, so the Garden is not gated on buying a Frost Mage.

Live terminal checks cover placement, wave start, hints and quitting on all four regular floors. Render previews also cover the Heart. Automated results establish tactical differences and catch stalled or unaffordable policies; whether the pacing and choices feel fun still needs human playtesting. The useful next check is a fresh run through the campaign, especially the Garden raid and the Heart's transition from outer defense to inner reserve.
