package main

import (
	"flag"
	"fmt"
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

func parseLevelArgs(fs *flag.FlagSet, args []string) (*game.Map, string, *flag.FlagSet) {
	level := fs.String("level", "winding", "")
	mazeSeed := fs.Int64("maze", -1, "")
	_ = fs.Parse(args)
	if *mazeSeed >= 0 {
		// -maze 0 means "random maze": resolve the seed here (and name the
		// level after the real seed, so hiscores stay per-maze).
		ms := pickSeed(*mazeSeed)
		m, err := game.MazeFromSeed(ms)
		if err != nil {
			die("maze: %v", err)
		}
		return m, fmt.Sprintf("maze%d", ms), fs
	}
	m, err := game.LoadLevel(*level)
	if err != nil {
		die("%v", err)
	}
	return m, *level, fs
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
	m, name, _ := parseLevelArgs(fs, args)
	if err := tui.Run(m, name, diffFrom(fs)); err != nil {
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
			all = append(all, game.RunAutoplayDiff(m, name, diff))
		}
	}
	for i := 0; i < *nMaze; i++ {
		s := int64(100000 + i)
		m, err := game.MazeFromSeed(s)
		if err != nil {
			continue
		}
		all = append(all, game.RunAutoplayDiff(m, "maze", diff))
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
	if count > 0 {
		fmt.Printf("%-10s %7.0f%%\n", "ALL", 100*float64(wins)/float64(count))
	} else {
		fmt.Printf("%-10s %7s\n", "ALL", "n/a")
	}
	fmt.Printf("elapsed %v\n", time.Since(start).Round(time.Millisecond))
}

func stats(rs []game.SimResult) (wins int, waves, leaks, towers, gold, t float64) {
	for _, r := range rs {
		if r.Won {
			wins++
		}
		waves += float64(r.Wave)
		// r.Leaks counts actual leaks. 20-r.Lives would overcount boss
		// leaks (6 lives each) and be wrong on hard (15 starting lives).
		leaks += float64(r.Leaks)
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
	m, name, _ := parseLevelArgs(fs, args)
	r := game.RunAutoplayDiff(m, name, diffFrom(fs))
	fmt.Printf("%s: %s wave=%d lives=%d gold=%d towers=%d kills=%d time=%.0fs\n",
		name, status(r.Won), r.Wave, r.Lives, r.Gold, r.Towers, r.Kills, r.Time)
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
	m, name, _ := parseLevelArgs(fs, args)
	if *every < 1 {
		die("capture: -every must be >= 1")
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		die("%v", err)
	}
	s := game.NewStateDiff(m, diffFrom(fs))
	ai := game.NewAutoplay(s)
	pal := render.Palette()
	ui := render.UI{Cursor: game.Vec{X: m.W / 2, Y: m.H / 2}, Placing: game.TowerGunner, Selected: render.NoSelection, Speed: 1, Scale: *scale}
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
