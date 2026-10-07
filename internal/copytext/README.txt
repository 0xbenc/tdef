Editable game copy

The JSON files in data/ are the actual text used by the game. The shared
copytext package makes them available to game, render, and tui without an
import cycle. Cinematic scene copy stays in render/copy/opening.json and
render/copy/ending.json, beside its existing editing instructions.

Where to edit
  defenders.json  Eleven allies: names, subtitles, lore, tactical notes.
  enemies.json    Eight guild fighters: names, subtitles, lore, notes.
  places.json     Six places: names, lore, notes, and map/ledger labels.
  memories.json   Main story memories, grouped by floor and difficulty.
  afterward.json Four bonus memories, grouped by the required win count.
  overworld.json Ambient narration, Malgrath's voice, progression, results,
                 relic rewards, and overworld instructions.
  encounters.json Floor briefings, wave themes, and incoming-wave warnings.
  characters.json Portrait titles, descriptors, and navigation hints.
  outcomes.json  Victory/defeat flavor, verdicts, statistics, and controls.
  branding.json  Title copy, attract-demo labels, studio/publisher credits.
  ui.json        Journal headings and requirements, menus, help, setup,
                 HUD labels, feedback, reset confirmation, and footers.

COPY-REVIEW.md links each checklist section to its JSON sources and the Go
code that selects, unlocks, or lays out that text. Extraction preserves the
current writing; the existing lore still needs the planned copy reviews.

Editing rules
  Edit string VALUES. Keep object keys, file names, and structure intact.
  Keys identify the role of an entry; do not rename them when rewriting it.
  Use valid JSON. Escape quotes inside a string as \" and newlines as \n.
  Keep every {placeholder} with its exact name. You can move it or repeat it.
  Go supplies the real numbers, names, and control keys. For example:
      "{floor} broke at {wave} · the lair will mend"
    can become:
      "Wave {wave} broke the defense at {floor}."
  Do not insert percent-based formatting such as %d; use existing named
  placeholders. Compact number formatting and alignment remain in Go.
  Keep short labels and prompts short. Journal lore wraps; compact menus,
  HUD rows, captions, and footers have limited room in a small terminal.
  Wording does not control IDs, unlock rules, rewards, saved progress,
  keyboard bindings, difficulty settings, or message colors.

Defender/enemy display names are shared by gameplay and the journal. Other
name variants are separate because their capitalization or wording differs
by context. Map and ledger labels are display text, never saved floor IDs.
Art glyphs, terminal escape sequences, generated formations, and technical
errors remain code-owned. Changing game rules or adding an entry requires
Go changes too; contracts.go records the required keys and placeholders.

Check and play
  go test ./...
  go vet ./...
  go build -o termtd.exe .
  .\termtd.exe

JSON is embedded when building. Saving a file does not update an existing
executable or running game. Rebuild after editing, then review the text in
its real screen. Tests reject missing/empty entries, unknown/duplicate keys,
invalid JSON, and changed placeholders. They also check shared names and
persistent journal IDs; existing game tests cover progression and rewards.

Preview the films with .\termtd.exe intro or .\termtd.exe ending.
