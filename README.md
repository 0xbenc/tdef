# termtd

Tower defense in your terminal.

You’re Grak, the Last Monster. Heroes are coming for the hoard, and the dragon
who lives on it is dying. Build a defense. Keep him alive.

Eleven defenders, twenty waves, a lair to explore, and a journal full of lore.
Written in Go. Zero dependencies.

## Get it

Download it from [itch.io](https://kairuku-studios.itch.io/termtd) and use the Play
launcher, or install it from my Homebrew tap:

```sh
brew install 0xbenc/tap/termtd
termtd
```

For Windows 10/11, Linux, and macOS (12+), on Intel/AMD and ARM.

## Build it

With Go 1.26.3 or newer:

```sh
go build -o termtd .
./termtd
```

On Windows, build and run from PowerShell:

```powershell
go build -o termtd.exe .
.\termtd.exe
```

## Play

Use the arrows or WASD to move. Pick a tower with `1`–`7`, then press Enter
or click to place it. Green dots mark legal sites. Each build exits placement;
pick a tower again to build another. Escape cancels placement.

- `[` / `]` — switch between the original defenders and four specialists
- `r` — rotate a Runeforge while placing it or after selecting it
- `u` / `x` — upgrade / sell
- `Tab` — cycle through placed defenders; moving onto one also shows its controls
- `t` — cycle targeting priority; the dotted guide shows its current target
- `n` — send the next wave
- `p` / `f` — pause / change speed
- `j` — open the journal in the lair or while paused
- `h` — help · `esc` — back · `q` — quit

Each specialist is limited to one placed per level. Selling one frees its slot.

A fresh campaign teaches the main roster during the first Rotunda defense,
then introduces the specialists during the second defense. Recruitment pauses
the action and supplies gold for the new defender. Placing a defender resumes
the automatic wave countdown; `p` pauses and `n` sends a wave early. Earned
recruits stay unlocked on retries. Quick Play keeps the full roster.

Replay either guided defense without changing your saves:

```powershell
.\termtd.exe tutorial
.\termtd.exe tutorial -stage 2
```

Completed defenses, recruits, discoveries, and scores save automatically.
An unfinished defense restarts when you return. Use a terminal of at least
62 columns by 19 rows; more room gives the artwork and controls extra space.

**Start** advances your campaign in the lair. **Quick Play** lets you choose any
regular map or a procedural maze and records high scores, without changing
campaign progress, relics, journal discoveries, or story progress.

Open **Journal** from the main menu to browse discoveries earned in the campaign.

To start over, choose **Reset progress** in the main menu, then confirm
**Reset all progress**. This clears your campaign, relics, journal, and high scores.

[Story](LORE.md) · [Lore sources](https://github.com/0xbenc/termtd/blob/main/LORE-SOURCES.md) · [Field guide](PLAYING.md) · [Release setup](https://github.com/0xbenc/termtd/blob/main/RELEASING.md) · [MIT license](LICENSE)
