# Lore and copy sources

## JSON

| File | Contents |
|---|---|
| [render/copy/opening.json](render/copy/opening.json) | Opening scene titles, narration and dialogue. |
| [render/copy/ending.json](render/copy/ending.json) | Ending scene titles, narration and dialogue. |
| [defenders.json](internal/copytext/data/defenders.json) | Defender names, subtitles, journal lore and tactical notes. |
| [enemies.json](internal/copytext/data/enemies.json) | Guild fighter names, subtitles, journal lore and tactical notes. |
| [places.json](internal/copytext/data/places.json) | Place names, subtitles, journal lore, tactical notes and map/ledger labels. |
| [memories.json](internal/copytext/data/memories.json) | Story memory titles and lore, by floor and difficulty. |
| [afterward.json](internal/copytext/data/afterward.json) | Postgame memory titles and lore, by total Heart/Depths wins. |
| [overworld.json](internal/copytext/data/overworld.json) | Malgrath dialogue, room narration, progression, results, relic rewards and lair prompts. |
| [encounters.json](internal/copytext/data/encounters.json) | Wave themes, tactical briefings and incoming-wave warnings. |
| [characters.json](internal/copytext/data/characters.json) | Portrait names, captions, navigation and size prompts. |
| [outcomes.json](internal/copytext/data/outcomes.json) | Victory/defeat flavor, performance verdicts, results and result-screen controls. |
| [branding.json](internal/copytext/data/branding.json) | Title, tagline, signature, attract-demo labels, records and credits. |
| [ui.json](internal/copytext/data/ui.json) | Journal headings/requirements, menus, help, recruitment lessons, HUD, status, controls and system prompts. |

## Code

| File or location | Lore/copy role and lookup anchors |
|---|---|
| [render/cutscene_copy.go](render/cutscene_copy.go) | Embeds cinematic JSON and binds it to shots: `withOpeningCopy`, `withEndingCopy`, `withFilmCopy`. |
| [render/cutscene.go](render/cutscene.go) | Scene order, staging, film titles, captions and controls: `openingShots`, `endingShots`, `FilmTitle`, `RenderCutscene`. |
| [render/cutscene_art.go](render/cutscene_art.go), [cutscene_ending_art.go](render/cutscene_ending_art.go), [cutscene_panorama.go](render/cutscene_panorama.go), [cutscene_layers.go](render/cutscene_layers.go), [cutscene_grak.go](render/cutscene_grak.go) | Cinematic scene artwork and composition. |
| [tui/cutscene.go](tui/cutscene.go) | Opening/ending triggers, replay and input: `startCutscene`, `handleCutscene`, `maybeOpening`. |
| [render/journal.go](render/journal.go) | Defender entries, chapters, lore/note layout and collection navigation: `TowerJournalEntries`, `JournalChapter`, `RenderJournal`, `journalLines`. |
| [render/journal_entries.go](render/journal_entries.go) | Enemy/place entries, story and Afterward bindings, unlock text: `EnemyJournalEntries`, `PlaceJournalEntries`, `mainJournalMoments`, `buildStoryJournalEntries`, `JournalRequirement`. |
| [render/journal_scenes.go](render/journal_scenes.go) | Place, memory and Afterward illustrations: `drawJournalMoment`, `drawJournalAfterward`, place painters. |
| [tui/journal.go](tui/journal.go), [hiscore/journal.go](hiscore/journal.go) | Discovery, saved journal progress, memory unlocks and Afterward milestones: `ensureJournal`, `recordEnemyDiscoveries`, `RecordVictory`, `UnlockAfterward`. |
| [render/overworld.go](render/overworld.go) | Character voice, ambient room states, floor labels, progression/results/relic presentation: `owDragonVoice`, `owFloorFlavor`, `drawOWVoice`, `owNodeLine`, `drawOWChromeRows`. |
| [tui/overworld.go](tui/overworld.go) | Campaign seals, descent, renown and relic messages: `owSealedMessage`, `owLaunch`, `owRefresh`, `owSpendRelic`. |
| [tui/app.go](tui/app.go) | Return banners, unseal/Heart events, wave announcements and action feedback: `owSetBanner`, `owUnsealCheck`, `owHeartBlastCheck`, `stepGame`, `startWave`, `handleGame`. |
| [render/training.go](render/training.go), [tui/training.go](tui/training.go), [game/training.go](game/training.go) | Recruitment lesson selection, presentation, funding and timing: `recruitLesson`, `drawRecruitment`, `acknowledgeRecruit`, `lockedRecruit`, `RecruitmentOrder`. |
| [game/balance.go](game/balance.go) | Shared names, counter values, wave themes, warnings and generated formation previews: `TowerSpecs`, `EnemySpecs`, `WaveTheme`, `WaveTelegraph`, `WavePreview`. |
| [game/encounters.go](game/encounters.go) | Floor briefings, floor-specific warnings and formations: `TacticalBrief`, `WaveTelegraphFor`, `WavePreviewFor`, `BuildWaveFor`. |
| [game/types.go](game/types.go) | Targeting names and short labels: `TargetMode.Name`, `TargetMode.Short`. |
| [render/draw.go](render/draw.go) | Gameplay HUD, defender stats, wave banners, result flavor and verdicts: `towerInfo`, `headerSegments`, `drawBeats`, `drawGameOver`, `endVerdict`, `drawEndBeat`, `RenderTooSmall`. |
| [render/keyboard.go](render/keyboard.go), [render/specialists.go](render/specialists.go) | Keyboard action hints, targeting explanations and specialist HUD/roster copy. |
| [render/screens.go](render/screens.go) | Title/demo/records, menus, help, high scores and setup: `titleSig`, `drawTitleBest`, `drawTitleDemo`, `MenuItems`, `RenderHelp`, `RenderHighScores`, `RenderLevelSelect`. |
| [render/credits.go](render/credits.go), [render/reset.go](render/reset.go) | Credits and reset confirmation/result copy: `RenderCredits`, `RenderResetProgress`. |
| `render/*_mockup.go`, [render/specialist_portraits.go](render/specialist_portraits.go) | Character/defender/enemy portrait artwork; mockup renderers bind standalone portrait captions and controls from `characters.json`. |
| [internal/copytext/copy.go](internal/copytext/copy.go), [internal/copytext/contracts.go](internal/copytext/contracts.go) | Shared JSON embedding, lookup, templates, required keys and placeholders. |

## Editing references

| File | Reference |
|---|---|
| [internal/copytext/README.txt](internal/copytext/README.txt) | Shared JSON editing, placeholders, validation and rebuilding. |
| [render/copy/README.txt](render/copy/README.txt) | Cinematic JSON fields, shot ordering and preview commands. |
| [LORE.md](LORE.md) | Historical premise, terminology and draft theming; implementation claims may be outdated. |
