package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tdef/game"
	"tdef/render"
	"tdef/tui"
)

func main() {
	if len(os.Args) < 2 {
		play(nil)
		return
	}
	switch os.Args[1] {
	case "play":
		play(os.Args[2:])
	case "bench":
		bench(os.Args[2:])
	case "headless":
		headless(os.Args[2:])
	case "capture":
		capture(os.Args[2:])
	case "maps":
		mapsCmd()
	case "help", "-h", "--help":
		usage()
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `tdef - terminal tower defense

usage:
  tdef [play flags]        play (default)
  tdef bench [flags]       run autoplay balance benchmark
  tdef headless [flags]    run one autoplay game, print result
  tdef capture [flags]     render headless game frames to files
  tdef maps                list built-in levels
  tdef help

play flags:
  -level string   level name (default "winding"; see tdef maps)
  -maze int       use a procedural maze with this seed (0 = random)
  -seed int64     game seed (default: time-based)

bench flags:
  -n int          games per level (default 40)
  -levels string  comma list of levels, or "all" (default "all")
  -maze int       include N procedural mazes (default 40)

capture flags:
  -out string     output dir (default "frames")
  -every int      render every N ticks (default 10)
  -limit int      stop after N ticks (default: game end)
  -text           write plain text instead of ANSI
`)
}

func parseLevelArgs(fs *flag.FlagSet, args []string, autoWave bool) (*game.Map, string, int64, *flag.FlagSet) {
	level := fs.String("level", "winding", "")
	mazeSeed := fs.Int64("maze", -1, "")
	seed := fs.Int64("seed", 0, "")
	_ = fs.Parse(args)
	if *mazeSeed >= 0 {
		m, err := game.GenerateMap(*mazeSeed, 45, 13)
		if err != nil {
			die("maze: %v", err)
		}
		return m, fmt.Sprintf("maze%d", *mazeSeed), *seed, fs
	}
	m, err := game.LoadLevel(*level)
	if err != nil {
		die("%v", err)
	}
	return m, *level, *seed, fs
}

func pickSeed(seed int64) int64 {
	if seed != 0 {
		return seed
	}
	return time.Now().UnixNano()
}

func parseDiff(fs *flag.FlagSet) *flag.FlagSet {
	fs.String("diff", "normal", "difficulty: easy|normal|hard")
	return fs
}

func diffFrom(fs *flag.FlagSet) game.Difficulty {
	s := fs.Lookup("diff").Value.String()
	switch s {
	case "easy":
		return game.Easy
	case "hard":
		return game.Hard
	}
	return game.Normal
}

func play(args []string) {
	fs := flag.NewFlagSet("play", flag.ExitOnError)
	parseDiff(fs)
	m, name, seed, _ := parseLevelArgs(fs, args, false)
	seed = pickSeed(seed)
	if err := tui.Run(m, name, seed, diffFrom(fs)); err != nil {
		die("%v", err)
	}
}

func bench(args []string) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)
	parseDiff(fs)
	n := fs.Int("n", 40, "")
	levels := fs.String("levels", "all", "")
	nMaze := fs.Int("maze", 40, "")
	_ = fs.Parse(args)
	diff := diffFrom(fs)

	names := game.LevelNames()
	if *levels != "all" {
		names = strings.Split(*levels, ",")
	}
	total := len(names)*(*n) + *nMaze
	fmt.Printf("bench: %d games\n", total)
	start := time.Now()
	all := []game.SimResult{}
	for _, name := range names {
		m, err := game.LoadLevel(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "skip:", err)
			continue
		}
		for i := 0; i < *n; i++ {
			all = append(all, game.RunAutoplayDiff(m, name, int64(i+1), diff))
		}
	}
	for i := 0; i < *nMaze; i++ {
		s := int64(100000 + i)
		m, err := game.GenerateMap(s, 45, 13)
		if err != nil {
			continue
		}
		all = append(all, game.RunAutoplayDiff(m, "maze", s, diff))
	}
	byLevel := map[string][]game.SimResult{}
	for _, r := range all {
		byLevel[r.Level] = append(byLevel[r.Level], r)
	}
	fmt.Printf("%-10s %8s %10s %10s %8s %8s %8s\n", "level", "winrate", "avgWave", "avgLeaks", "towers", "gold", "time")
	wins, count := 0, 0
	for _, name := range append(append([]string{}, names...), "maze") {
		rs, ok := byLevel[name]
		if !ok || len(rs) == 0 {
			continue
		}
		w, waves, leaks, towers, gold, t := stats(rs)
		fmt.Printf("%-10s %7.0f%% %10.1f %10.2f %8.1f %8.0f %8.0f\n",
			name, 100*float64(w)/float64(len(rs)), waves, leaks, towers, gold, t)
		wins += w
		count += len(rs)
	}
	fmt.Printf("%-10s %7.0f%%\n", "ALL", 100*float64(wins)/float64(count))
	fmt.Printf("elapsed %v\n", time.Since(start).Round(time.Millisecond))
}

func stats(rs []game.SimResult) (wins int, waves, leaks, towers, gold, t float64) {
	for _, r := range rs {
		if r.Won {
			wins++
		}
		waves += float64(r.Wave)
		leaks += float64(20 - r.Lives)
		towers += float64(r.Towers)
		gold += float64(r.Gold)
		t += r.Time
	}
	n := float64(len(rs))
	waves, leaks, towers, gold, t = waves/n, leaks/n, towers/n, gold/n, t/n
	return
}

func headless(args []string) {
	fs := flag.NewFlagSet("headless", flag.ExitOnError)
	parseDiff(fs)
	m, name, seed, _ := parseLevelArgs(fs, args, true)
	seed = pickSeed(seed)
	r := game.RunAutoplayDiff(m, name, seed, diffFrom(fs))
	fmt.Printf("%s seed=%d: %s wave=%d lives=%d gold=%d towers=%d kills=%d time=%.0fs\n",
		name, seed, status(r.Won), r.Wave, r.Lives, r.Gold, r.Towers, r.Kills, r.Time)
}

func status(won bool) string {
	if won {
		return "WIN"
	}
	return "LOSE"
}

func capture(args []string) {
	fs := flag.NewFlagSet("capture", flag.ExitOnError)
	parseDiff(fs)
	out := fs.String("out", "frames", "")
	every := fs.Int("every", 10, "")
	limit := fs.Int64("limit", 0, "")
	text := fs.Bool("text", false, "")
	scale := fs.Int("scale", 1, "playfield scale 1-4")
	m, name, seed, _ := parseLevelArgs(fs, args, true)
	seed = pickSeed(seed)
	if err := os.MkdirAll(*out, 0o755); err != nil {
		die("%v", err)
	}
	s := game.NewStateDiff(m, seed, true, diffFrom(fs))
	ai := game.NewAutoplay(s)
	pal := render.Palette()
	ui := render.UI{Cursor: game.Vec{X: m.W / 2, Y: m.H / 2}, Placing: game.TowerGunner, Speed: 1, Scale: *scale}
	dt := 1.0 / 20.0
	tick := int64(0)
	for s.Status == game.StatusRunning {
		ai.Tick()
		s.Step(dt)
		tick++
		if tick%int64(*every) == 0 {
			f := render.Render(s, &ui, pal)
			p := filepath.Join(*out, fmt.Sprintf("%s_%06d", name, tick))
			var content string
			if *text {
				content = f.Text() + "\n"
			} else {
				content = f.ANSI() + "\n"
			}
			os.WriteFile(p, []byte(content), 0o644)
		}
		if *limit > 0 && tick >= *limit {
			break
		}
	}
	st := "running"
	if s.Status == game.StatusVictory {
		st = "WIN"
	} else if s.Status == game.StatusDefeat {
		st = "LOSE"
	}
	fmt.Printf("wrote frames to %s (ticks %d, wave %d, %s)\n", *out, tick, s.Wave, st)
}

func mapsCmd() {
	fmt.Println("built-in levels:")
	for _, n := range game.LevelNames() {
		m, err := game.LoadLevel(n)
		if err != nil {
			fmt.Printf("  %-12s %v\n", n, err)
			continue
		}
		fmt.Printf("  %-12s path %.0f cells\n", n, m.TotalLen)
	}
	fmt.Println("procedural: -maze <seed>")
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "tdef: "+format+"\n", args...)
	os.Exit(1)
}

var _ = rand.Int64N
