# Playable game lore and copy review

Inventory date: 2026-10-05. Scope: writing displayed to a player in the implemented game, including campaign, Quick Play, journals, cinematics, help, and UI feedback. Repository lore documents and store copy are outside this inventory.

This preserves the 29 categories identified in the initial game-content audit and adds source references for reviewing them one at a time. It is a category inventory, not a transcript of every string. Counts describe the current authored content.

## Review workflow

Pick an unchecked numbered item and inspect all its referenced sources. Review the existing text, make the needed copy changes, check it in its player-facing context, then mark the item complete and record the outcome below. A review can finish with no changes.

File links are relative to this document. Symbols are the primary lookup anchors because line numbers move as code changes. Some sources appear in several items because they contain several kinds of copy.

For each change, check character names and terminology, mechanical accuracy, unlock conditions, and fit in the terminal layout. Update any tests that assert the changed text. Preserve IDs, gameplay values, and discovery logic unless the review explicitly includes changing them.

## Checklist and source map

| Done | ID | Type | What players encounter | Sources and lookup anchors |
|---|---|---|---|---|
| [ ] | 01 | Opening cinematic narration | The guild's hunt, Grak's arrival, and his decision to protect Malgrath. | [render/cutscene.go](render/cutscene.go): `openingShots`, `Dialogue` on shots without a speaker. |
| [ ] | 02 | Opening cinematic dialogue | Spoken exchanges between the Guildmaster, Grak, and Malgrath. | [render/cutscene.go](render/cutscene.go): `openingShots`, `Speaker` and `Dialogue`. |
| [ ] | 03 | Ending cinematic narration and dialogue | The aftermath of the siege and Grak and Malgrath's closing conversation. | [render/cutscene.go](render/cutscene.go): `endingShots`. |
| [ ] | 04 | Defender lore | Seven journal entries describing tower allies, personalities, and relationships. | [render/journal.go](render/journal.go): `TowerJournalEntries`, `Lore`. |
| [ ] | 05 | Enemy lore / bestiary | Eight journal entries describing guild fighters, including The Player. | [render/journal_entries.go](render/journal_entries.go): `EnemyJournalEntries`, lore field. |
| [ ] | 06 | Place lore | Six journal entries about chambers, their history, and personal significance. | [render/journal_entries.go](render/journal_entries.go): `PlaceJournalEntries`, lore field. |
| [ ] | 07 | Unlockable story memories | Eighteen vignettes earned by holding each of six floors at each of three difficulties. | [render/journal_entries.go](render/journal_entries.go): `mainJournalMoments`, `buildStoryJournalEntries`. |
| [ ] | 08 | Postgame story memories | Four Afterward vignettes earned through additional Heart / Depths victories. | [render/journal_entries.go](render/journal_entries.go): `bonus` inside `buildStoryJournalEntries`; thresholds 1, 3, 5, and 10 wins. |
| [ ] | 09 | Journal titles and flavor subtitles | Entry names, poetic descriptors, chapter headings, and collection labels. | [render/journal.go](render/journal.go): `TowerJournalEntries`, `JournalChapter`, `RenderJournal`, `drawJournalCollection`, `drawJournalPage`; [render/journal_entries.go](render/journal_entries.go): entry names/subtitles and story titles. |
| [ ] | 10 | Journal tactical notes | Practical advice accompanying defender, enemy, and place lore. | [render/journal.go](render/journal.go): `TowerJournalEntries`, `Note`, `journalLines`; [render/journal_entries.go](render/journal_entries.go): notes in `EnemyJournalEntries` and `PlaceJournalEntries`. |
| [ ] | 11 | Journal discovery and unlock copy | Discovery requirements, memory requirements, and collection progress. | [render/journal_entries.go](render/journal_entries.go): `JournalRequirement`, generated notes in `buildStoryJournalEntries`; [render/journal.go](render/journal.go): `drawJournalCollection`, `drawJournalPage`. Check conditions against [tui/journal.go](tui/journal.go). |
| [ ] | 12 | Overworld ambient narration | The lair waking, corridors lighting, descending, and contextual room text. | [render/overworld.go](render/overworld.go): `drawOWVoice`, `owFloorFlavor`, `owNodeLine`. |
| [ ] | 13 | Overworld character voice | Malgrath's contextual lines in the lair. | [render/overworld.go](render/overworld.go): `owDragonVoice`, `owDragonVoiceDone`, `drawOWVoice`, `owNodeLine`. |
| [ ] | 14 | Campaign progression copy | Sealed-floor explanations, opening doors, the Heart unsealing, and Depths requirements. | [tui/overworld.go](tui/overworld.go): `depthsLockedMessage`, `owSealedMessage`, `owTick`; [tui/app.go](tui/app.go): `owUnsealCheck`, `owHeartBlastCheck`; [render/overworld.go](render/overworld.go): `drawOWVoice`, `owNodeLine`, `drawOWLabel`, `drawOWChromeRows`. |
| [ ] | 15 | Campaign return / result messages | Feedback after a defense, including whether a floor held or fell. | [tui/app.go](tui/app.go): `owSetBanner`; [render/overworld.go](render/overworld.go): `drawOWVoice`, `owNodeLine`, `owInfoLine`, `drawOWChromeRows`. |
| [ ] | 16 | Relic descriptions and reward messages | Relic spending for gold, a free gunner, or additional heart health. | [tui/overworld.go](tui/overworld.go): `owSpendRelic`; [render/overworld.go](render/overworld.go): relic menu in `drawOWVoice`, bonus/relic labels in `drawOWChromeRows`. |
| [ ] | 17 | Floor tactical briefings | Short battlefield-specific advice. | [game/encounters.go](game/encounters.go): `TacticalBrief`; [render/draw.go](render/draw.go): rendering of `game.TacticalBrief`. |
| [ ] | 18 | Incoming-wave warnings / rumors | New-threat introductions and expedition advice, including floor-specific variants. | [game/balance.go](game/balance.go): `WaveTelegraph`; [game/encounters.go](game/encounters.go): `WaveTelegraphFor`; [tui/app.go](tui/app.go): `stepGame` for display timing. |
| [ ] | 19 | Wave names and announcements | Expedition themes, wave-start banners, and incoming formation previews. | [game/balance.go](game/balance.go): `waves`, `WaveTheme`, `WavePreview`; [game/encounters.go](game/encounters.go): `WavePreviewFor`; [render/draw.go](render/draw.go): wave/header/preview/banner text; [tui/app.go](tui/app.go): `startWave`. |
| [ ] | 20 | Character identity captions | Names and descriptors on Grak and Malgrath portrait screens. | [render/grak_mockup.go](render/grak_mockup.go), [render/dragon_mockup.go](render/dragon_mockup.go): portrait renderers and captions. |
| [ ] | 21 | Victory and defeat flavor text | “The lair is held. Malgrath endures.” / “Malgrath has fallen. The lair is clean.” | [render/draw.go](render/draw.go): `drawGameOver`, `lore`. |
| [ ] | 22 | Performance verdicts | Flawless holds, hard-fought victories, and early/late defeats. | [render/draw.go](render/draw.go): `endVerdict`. |
| [ ] | 23 | Results and high-score copy | Run statistics, best scores, record announcements, and leaderboard labels. | [render/draw.go](render/draw.go): `drawGameOver`; [render/screens.go](render/screens.go): `RenderHighScores`, `drawTitleBest`; [tui/app.go](tui/app.go): `owSetBanner`; [render/overworld.go](render/overworld.go): `owInfoLine`, `drawOWChromeRows`. |
| [ ] | 24 | Title / branding copy | TERMTD, terminal tower defense tagline, creator signature, and THE SIEGE attract-demo label. | [render/screens.go](render/screens.go): title renderers, `drawTitleSig`, `drawTitleTagline`, `drawTitleEmptyBox`, `drawTitleDemo`, boot subtitle/UI. |
| [ ] | 25 | Menu and setup copy | Start, Quick Play, Help, High Scores, floor selection, difficulty/renown, and maze seeds. | [render/screens.go](render/screens.go): `MenuItems`, `RenderMenu`, `DiffName`, `RenderLevelSelect`; [render/draw.go](render/draw.go): `diffName`, `levelDisplayName`; [render/overworld.go](render/overworld.go): floor names and `drawOWChromeRows`; [tui/app.go](tui/app.go): `parseSeed`, `startGame`, `lsView`. |
| [ ] | 26 | Help / instructional copy | Grak's Ledger, controls, mechanics hints, and first-run guidance. | [render/screens.go](render/screens.go): `RenderHelp`; [render/overworld.go](render/overworld.go): first-run branch in `drawOWVoice`, `drawOWChromeRows`; [render/draw.go](render/draw.go): gameplay hints; [tui/app.go](tui/app.go): `startWave` first-tower prompt. |
| [ ] | 27 | HUD and tower information | Resources, health, score, countdowns, tower names/costs/levels/stats, targeting, upgrades, and sale values. | [render/draw.go](render/draw.go): `towerInfo`, header and tower-slot rendering; [game/balance.go](game/balance.go): `TowerSpecs`, `EnemySpecs` display names; [game/types.go](game/types.go): `TargetMode.Name`, `TargetMode.Short`; [render/overworld.go](render/overworld.go): campaign HUD in `drawOWChromeRows`. |
| [ ] | 28 | Interaction and status messages | Action feedback, pause/speed states, unavailable actions, and input errors. | [tui/app.go](tui/app.go): `msg` callers, `handleGame`, `place`, `upgradeSelected`, `sellSelected`, `cycleTarget`, `startWave`, `parseSeed`; [tui/overworld.go](tui/overworld.go): `Msg` assignments, `owLaunch`, `owSealedMessage`; [render/draw.go](render/draw.go): message and pause/status rendering; [render/screens.go](render/screens.go): level-select error display. |
| [ ] | 29 | Navigation and system prompts | Hotkey footers, back/restart/quit/replay, scrolling, and terminal-size prompts. | [render/screens.go](render/screens.go): `titleFooter`, screen-specific footer segments, `RenderHighScores`, `RenderLevelSelect`; [render/cutscene.go](render/cutscene.go): `FilmTitle`, `RenderCutscene`; [render/journal.go](render/journal.go): collection/page navigation; [render/draw.go](render/draw.go): `RenderTooSmall`, `drawGameOver`, gameplay controls; [render/overworld.go](render/overworld.go): `drawOWChromeRows`; [render/grak_mockup.go](render/grak_mockup.go), [render/dragon_mockup.go](render/dragon_mockup.go): navigation and size hints. |

## Review log

Add a row for each completed review or outstanding decision. Use the checklist ID to keep future bot sessions scoped.

| ID | Date | Outcome / changes | Verification | Follow-up |
|---|---|---|---|---|
| — | — | — | — | — |

## Suggested handoff to a bot

> Review item NN in COPY-REVIEW.md. Read its source references and the current player-facing text. Make the needed copy changes, verify layout and mechanical accuracy, update relevant text assertions, then record the outcome in the review log and mark the item complete. Keep the work scoped to that item and any shared wording it directly affects.

## Context references

[LORE.md](LORE.md) contains theming and draft-copy context; its statements about implementation status are historical. [README.md](README.md) describes the current game and basic controls. Use the playable source above to establish what players actually see.
