# mtop — Claude Code project memory

A Bubble Tea TUI dashboard for the multica daemon. ~2.1 KLOC Go, currently
mocked against hardcoded data in `data.go`.

## Active work

@.claude/PLAN.md

When starting a session, read the imported plan above, find the first `[~]` or
`[ ]` task in the active phase, and continue from there. Update checkboxes in
place as work completes. Append discoveries to the **Rules & Tips** section of
`.claude/PLAN.md` — never delete entries from it.

## Conventions

- Go module: see `go.mod`. TUI uses `github.com/charmbracelet/bubbletea` +
  `lipgloss`. After Phase B, also `bubbles/*`.
- Branch naming matches the **Branch** column of the Phase Summary table.
- One phase per branch; merge to `main` between phases.
- Each phase ends with regenerated golden snapshots — never check in code that
  changes goldens silently. If goldens change, the diff goes in the PR
  description.
- No behavior changes inside refactor phases (A). Behavior changes belong to
  feature phases (C, D, E, F).
- Do not introduce abstractions that aren't earning their keep across at least
  two callers in this repo.

## Build & test

- Build: `go build -o mtop .`
- Run (mock data): `MTOP_MOCK=1 ./mtop` (gating lands in Phase D.5; before then
  the binary uses `data.go` unconditionally)
- Dump for golden snapshots: `MTOP_DUMP_PROFILE=1 ./mtop --dump > testdata/<W>x<H>.txt`
- Goldens covered sizes: 80×24, 120×40, 200×60

## What lives where

- `main.go` — entry point, tea program setup, alt-screen, mouse motion
- `chrome.go` — shared chrome (tab bar, footer, cmd bar, `pane`, `kv`, `hr`, `bar`)
- `theme.go` — palette
- `data.go` — currently hardcoded fixtures (to be moved under `internal/source/mock.go` in Phase D.5)
- `view_<tab>.go` — one file per tab (status, run, tasks, log, config, profiles)
