# tdef devlog

- 2026-09-18 (session start)
- engine + maps + balance v1 done; TUI in progress
2026-09-18 19:21:44
- 2026-09-18 20:02:54
- TUI done + PTY-verified (input/render/mouse/cleanup)
- balance: per-map HP tuning + economy cut; defeat-overwrite bug found+fixed+tested
- juice: kill flash, leak flash+bell, combo, wave preview, early bonus, 3 difficulties
- maze=chaos(51% AI), handcrafted ladder winding->garden->canyon
2026-09-18 20:02:54
- 2026-09-18 20:30:56
- fixed input parser: keys were parsed but never emitted (arrows/enter/fkeys dead); empty-data panic; unknown SS3/CSI now terminate
- PTY retest: cursor move + tower place + wave start + quit all verified end-to-end
- added intro/help screen (title, controls, map+difficulty, press-any-key)
- terminal size: 0x0 pty fallback to 80x24 + blit clamp + full redraw on SIGWINCH
- hiscore "NEW BEST", README; full suite green (vet+test+build), bench ALL 85%
- 2026-09-18 20:45:49
- CRITICAL input bug: raw mode left Iflag.ICRNL set, so pty turned Enter \r into \n which the reader ignored -> keyboard tower placement was BROKEN (mouse-only). Fixed raw flags (Iflag/Oflag/Lflag) + reader accepts \n as Enter. PTY-verified place+target now work.
- added per-tower targeting priority (t key): first/strongest/closest; engine switch in acquireTarget, CycleTarget in state, HUD shows mode; unit tested all 3 modes
- PTY harness note: `script`+pipe returns 124 on cleanup (artifact); game exit code 0 confirmed via COMMAND_EXIT_CODE
