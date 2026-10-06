# termtd

Tower defense in your terminal.

You’re Grak, the Last Monster. Heroes are coming for the hoard, and the dragon
who lives on it is dying. Build a defense. Keep him alive.

Seven towers, twenty waves, a lair to explore, and a journal full of lore.
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
or click to place it.

- `u` / `x` — upgrade / sell
- `n` — send the next wave
- `p` / `f` — pause / change speed
- `j` — open the journal in the lair or while paused
- `h` — help · `esc` — back · `q` — quit

Your campaign saves automatically. Bring a terminal with a little room to spare.

**Start** advances your campaign in the lair. **Quick Play** lets you choose any
regular map or a procedural maze and records high scores, without changing
campaign progress, relics, journal discoveries, or story progress.

Open **Journal** from the main menu to browse discoveries earned in the campaign.

To start over, choose **Reset progress** in the main menu, then confirm
**Reset all progress**. This clears your campaign, relics, journal, and high scores.

[Story](LORE.md) · [Tactics](TACTICS.md) · [Release notes & setup](RELEASING.md) · [MIT license](LICENSE)
