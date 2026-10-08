# Grak's field guide

## Get started

Choose **Start** for the campaign. The first two guided defenses introduce the
roster and provide gold for each new recruit. **Quick Play** offers the full
roster and records scores without changing campaign discoveries or progress.

Hold all twenty waves to win a defense. Every fighter who reaches the exit
strikes Malgrath. Ordinary fighters cost one heart; a champion costs six.
Gold pays for defenders and upgrades. Kills and cleared waves earn more.

## Build a defense

Start with affordable shooters covering the first part of the road. Look for
bends or returning lanes where one defender can fire on several passes. Check
range before placing. Walls restrict building, but do not block attacks.

Put frost where your shooters can use the extra time. Use blasts against
crowds, long-range arrows against spaced ranks, and a second firing position
to catch survivors. Upgrading a useful position is often better than filling
unhelpful ground.

| Fighter | What to do |
|---|---|
| Rogue | Frost disables evasion; blasts and magic bypass it. |
| Paladin | Use blasts or magic. Plate resists bullets, stones, and arrows. |
| Centurion | Use physical attacks. Wards resist magic, but frost still slows. |
| Necromancer | Keep fire ready for the two Squires raised when it dies. |
| Champion | Save damage for the final approach. A breach costs six hearts. |

## Specialists

Use `[` and `]` to switch roster pages. You may place one of each specialist
per defense; selling frees its slot.

- **Runeforge:** rotate with `r` and aim its beam down a straight road.
- **Ogre Hookmaster:** pull fighters back through your strongest firing area.
- **Goblin Sappers:** prepare mines near the road. Early fighters can use up
  charges meant for heavier ones. Selling the crew removes its mines.
- **Hex Witch:** put her near damage dealers. Her mark increases their damage;
  she deals none herself. Marks do not stack.

## Controls

| Key | During a defense |
|---|---|
| Arrows / WASD | Move the cursor. |
| 1-7, then Enter / click | Choose and place a defender. Choose again for another. |
| Tab | Select the next placed defender. |
| u / x | Upgrade / sell the selected defender. |
| t | Cycle first, strongest, or closest targeting. |
| r | Rotate a selected Runeforge or its placement preview. |
| [ / ] | Switch roster pages. |
| n | Send the next wave; sending early earns bonus gold. |
| p / f | Pause / cycle speed. |
| h | Open Help; the defense pauses while you read. |
| j | Open the journal while paused. |
| Esc | Cancel placement or selection; return after a finished defense. |
| q | Quit. |

In the lair, walk onto a floor and press Enter to defend it. Tab changes renown.
At the Rotunda, `t` spends relics on the next defense, `i` replays the opening,
and `e` replays an earned ending. On the Depths, type a seed for a repeatable
route; Backspace deletes a digit and `c` clears it.

## Progress and saves

Floor victories are recorded separately for each renown. Hold the Rotunda,
Rift, Long Halls, and Garden on one renown to open the Heart; hold the Heart too
to open the Depths. Recruits stay unlocked on retries. The journal tracks
campaign discoveries and memories.

Campaign records, recruits, discoveries, and scores save automatically in your
home folder. A defense in progress is not saved: quitting means starting that
defense again. If a save fails, keep the game open and fix the reported problem.
It retries automatically; Ctrl+L retries immediately. Quitting with unsaved
changes requires a second quit.

If existing progress cannot be loaded, a persistent warning names the affected
saves. Their saving stays blocked for that session; other saves still work.
Press Ctrl+L to inspect the file paths and error, and use Up/Down to scroll.
Opening these details during a defense pauses it.

Malformed saves get an exact recovery copy beside the original, named
`<save filename>.corrupt-<unique suffix>`. The original stays untouched even if
making the copy fails. Unreadable files also stay untouched. Files live in your
home folder as `.termtd-lair.json`, `.termtd-journal.json`, and
`.termtd-hiscores.json`; older `.tdef-` files are read only when the corresponding
current file is absent.

To recover, quit, keep copies of your files, repair the reported JSON file or
restore a known-good copy at that path, then restart the game. Temporary progress
played with saving blocked does not merge into restored progress. To deliberately
start over instead, use the confirmed Reset progress action. Recovery copies are
kept through reset and are never loaded automatically.

Reset progress in the main menu clears campaign progress, relics, the journal,
and high scores after confirmation.
