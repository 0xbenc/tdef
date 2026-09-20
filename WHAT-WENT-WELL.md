# What Went Well

A retrospective on the title-screen iteration (2026-09-19, branch
`title-screen`, `a2183f1..a4fe0eb`), written by the AI model that
implemented it, for human engineers and other models to review. `devlog.md`
records the *what*; this file is the *why it worked* — with the caveats,
because a retrospective without caveats is just a press release.

## Scorecard

| metric | result |
|---|---|
| commits | 5, each independently green (game / render / tui / main / docs) |
| engine regression | none — bench identical before/after (canyon/garden/hub/winding 100%, maze 62%, ALL 66%) |
| layout defects caught before first commit | 4, all by eyeballing rendered frames (tests saw none of them) |
| test-side bugs | 2 (both in my own assertions, not the code) |
| race detector | clean, including a PTY resize storm through the new screens |
| PTY end-to-end flows | title→quit, title→menu→level→game, typed maze seed→`maze1234` game, help, hiscores — all exit 0 |

## The moves that did the heavy lifting

### 1. The design phase was real, parallel, and contractual

Three read-only subagents (architecture / visual-UX / impl+tests) ran
concretely against the repo before any code was written. The output was
integrated into a **contract**, not prose: the `Screen` enum, exact
renderer signatures (`RenderTitle(w, h, frame, scores, pal) *Frame`,
…), layout formulas, the CLI contract (bare `tdef` → menu flow, explicit
`-level`/`-maze` → direct start), and a named test list per screen.

Why it worked: implementation became *transcription, not exploration*.
The ~860-line tui commit was written in a single pass and its only
failures were two bugs in my own tests.

Reuse: when a feature is mostly design, spend the context budget on
parallel design agents with **disjoint scopes**, and require each to emit
signatures + invariants + a test list. "Here's a cool idea" is not a
design deliverable; "here is the API and what must stay true" is.

### 2. The interesting parts are pure functions of explicit inputs

Every screen renderer is `f(w, h, state, frame, palette) → *Frame`.
The title animation is driven by an **integer 30fps tick counter passed
as an argument** — no wall clock, no globals, no I/O. The only impure
layer is the blit.

Why it worked: determinism tests (`same frame index ⇒ identical frame`),
animation tests (`frames 0 and 15 differ`), and fits-frame tests across
terminal sizes are all trivial. A throwaway harness can render any frame
on demand and diff it. The "graphic and awesome" ask reduced to
"make this pure function pretty."

Reuse: the clock is an argument, not a dependency. If a piece of UI is
hard to test, the usual cause is that time or randomness is hiding
inside it.

### 3. Verification meant *looking*, not just asserting

This is the single highest-yield habit of the iteration. After the render
layer compiled and its tests passed, I ran the renderers in a standalone
harness and **eyeballed the ASCII frames**. That found four defects no
test had seen:

- menu block top-heavy (centering formula was off by a row);
- level-select: the maze row (list item `n`) collided with the seed row —
  the layout math reserved `n` rows for `n+1` rows;
- the map preview wouldn't appear until ~40 rows tall; the chrome was
  too loose to fit a 13-row preview on a common 32-row terminal;
- roster lines misaligned + the seed caret drifted a column on odd-width
  frames.

In the PTY phase the same habit paid again in reverse: my grep-based flow
checks "failed," and instead of patching the app I **wrote a 60-line ANSI
screen-state reconstructor** and discovered the app was correct — the
blit only writes *changed* cells, so literal strings get fragmented in
the capture. The check was wrong; the code was right.

Reuse: for anything user-visible, build a cheap "see the output" tool
*before* the first commit. And when a check fails, interrogate the check
before the code — a false negative sent me debugging a non-bug for a
while, and the reconstructor is what settled it.

### 4. One source of truth for geometry: render and hit-test share a layout

`MenuRects` / `LSRects` are computed by the *same* layout functions the
renderers use, so a click can never drift from a pixel. Tests assert the
rects land on the rendered text (e.g. the difficulty rect's cell holds
the label's first character). Mouse support then cost almost nothing and
the mouse tests are two lines each.

Reuse: whenever "what the user sees" and "what the user can hit" are two
concepts, derive both from one function. The drift class of bug doesn't
get fixed by tests; it gets fixed by construction.

### 5. Baseline before, regression after, negative controls

Bench numbers were captured **before** the feature and re-run after:
identical. A UI feature "shouldn't change engine behavior" — proving that
with one command is cheaper than assuming it. The same pattern ran
through the review-fix passes: the Tesla-chain test fails on the old
code (negative control), the `Terminal.size` race fix was verified with a
PTY resize storm where the *old* build trips the race detector and the
new one is clean.

Reuse: for any change, write down what must stay the same *before*
changing, then re-measure and report the identity. "Nothing changed" is a
claim that needs evidence.

### 6. Small commits, green tree, an honest debt ledger

Five commits, each green on its own; the diff of each is one concern
(engine data, render, state machine, CLI, docs). And the work that was
*not* done — review items 6–8 (UTF-8 reader, `LoadMap` validation,
`capture -scale` clamp, path semantics, LOW sweep) — was written into the
devlog and the final report as still open.

Reuse: a commit that isn't independently green isn't done. Parked debt
written down is an asset; parked debt in a head is a regression waiting
to happen.

### 7. Mistakes became process, and process survived session boundaries

The session summaries carried forward hard-won gotchas verbatim:
`git checkout -- <file>` reverts *all* uncommitted changes in that file
(happened twice — the fix is surgical edits plus full-file `/tmp`
backups); the bash tool validates `workdir` existence before executing
(you can't `mkdir` your own workdir inside the command); zsh mangles
`grep --include` globs (use the dedicated tool). Each of these cost a
rework cycle the first time and zero the second.

Reuse: log *process* failures (tooling, environment, workflow), not just
code bugs. They compound silently across sessions, and a frontier model
with good memory of them outperforms a smarter model without.

## Credit where due

The human's constraints did half the work, and the record should say so:

- **Zero external dependencies + deterministic engine** (set in session
  one) *forced* the pure-function renderer design. With an allowed GUI
  toolkit, the testable design above would never have been the path of
  least resistance.
- **The prompt allocated the design phase explicitly** — "do bunch of
  subagents about how best to design and implement this. *then* execute
  on a new branch." That single instruction is probably the biggest
  quality driver in this iteration.
- The standing preferences — one commit per fix, verify empirically,
  devlog entries — gave the model a quality bar to converge on instead of
  inventing one per session.

The "10x" here is a division of labor: the human set the constraints and
the verification culture; the model executed against them and kept the
evidence trail.

## What I'd do differently (the nuance)

- **The layout math took 2–3 iterations.** The row collision and the
  centering off-by-one were both visible on paper before a single frame
  was drawn. For grid layouts, write the vertical budget as a test with
  exact expected rows *first*; the test would have caught the maze-row /
  seed-row collision on the first try.
- **The visual harness and the ANSI reconstructor were throwaway**
  (`/tmp`). A committed `tdef screens` debug command — render every
  screen to text/ANSI at a given size and frame — would make visual
  regression cheap, shareable, and part of the review loop.
- **The first PTY test was a false negative**: all keys arrived in one
  batch, so the app quit before drawing the intermediate screens.
  Spaced input plus state reconstruction is what actually verified the
  flow. PTY tests need timing, and "it exited 0" is a weak assertion.
- **The stale-binary handoff gap**: the feature was verified in
  `/tmp/tdef-bin` while the user's `./tdef` in the repo dir was built
  hours earlier and predated the work. Verification artifacts should be
  rebuilt *where the user runs them*, or the report should say exactly
  which binary to run.
- The title's letter grid is a hand-tuned constant. A small generator
  (or image→block-element conversion) would own it better than a
  `[4][5]string` literal.

## Reuse card for other models

1. Design phase: parallel read-only agents, disjoint scopes, deliver a
   contract (signatures, invariants, test list) — not prose.
2. Make the new hard part a pure function; the clock is an argument.
3. Before the first commit, build a tool that lets you *see* the output.
4. When a check fails, interrogate the check before the code.
5. Render geometry and hit-test geometry come from one function; test
   the derivation, not the coincidence.
6. Capture the baseline before the change; re-run after; report the
   identity.
7. One commit per concern; green at every commit; write the parked debt
   down where the next session (human or AI) will read it.
8. Log process failures into session memory; they compound.
