# tdef

Tower defense in your terminal.

You’re Grak, the Last Monster. Heroes are coming for the hoard, and the dragon
who lives on it is dying. Build a defense. Keep him alive.

Seven towers, twenty waves, a lair to explore, and a journal full of lore.
Written in Go. Zero dependencies.

## Get it

Download it from [itch.io](https://kairuku-studios.itch.io/tdef) and use the Play
launcher, or install it from my Homebrew tap:

```sh
brew install 0xbenc/tap/tdef
tdef
```

For Linux and macOS (12+), on Intel/AMD and ARM.

## Build it

With Go 1.26.3 or newer:

```sh
go build -o tdef .
./tdef
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

[Story](LORE.md) · [Tactics](TACTICS.md) · [Release notes & setup](RELEASING.md) · [MIT license](LICENSE)
