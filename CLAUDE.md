# autopus-tui

A Bubble Tea TUI prototype for the (fictional) `autopus` agent daemon. Multi-package Go layout designed so you can find and edit any single view without loading the rest.

## Verifying changes — use `--dump`

This binary has no test suite. The way to verify a rendering change is:

```sh
go run . --dump
```

That renders every top-level route + every overlay + the attach view to stdout and exits. No real terminal needed. Use it from CI / Claude Code to confirm a change before declaring done.

Env vars that control the dump (see `dump.go`):

- `AUTOPUS_DUMP_PROFILE` = `truecolor` | `256` | `ascii`   (default `truecolor`)
- `AUTOPUS_DUMP_WIDTH`   = int                              (default 160)
- `AUTOPUS_DUMP_HEIGHT`  = int                              (default 44)
- `AUTOPUS_DUMP_ONLY`    = comma-separated view names       (default all)

Recognized names: `overview, sessions, issues, runtimes, workspaces, logs, config, attach, help, palette, onboarding`.

Examples:

```sh
# Just the attach view, ANSI-stripped, narrow:
AUTOPUS_DUMP_PROFILE=ascii AUTOPUS_DUMP_ONLY=attach AUTOPUS_DUMP_WIDTH=120 go run . --dump

# Sessions + palette overlay, 256-color, default size:
AUTOPUS_DUMP_PROFILE=256 AUTOPUS_DUMP_ONLY=sessions,palette go run . --dump
```

## Project layout

```
main.go                 entrypoint — arg parsing, launches tea.Program or runDump
dump.go                 --dump driver (env vars above)

internal/
  app/                  cross-package message types (NavigateMsg, ClockMsg, SpinMsg) + route table
  theme/                colors, glyphs, lipgloss styles, bg-leak helpers — zero deps
  ui/                   reusable primitives: text, frame, panel, chip, kv, bars
  data/                 static fixtures (sessions, issues, runtimes, workspaces, logs, profiles, help, palette, onboarding, daemon, activity)
  chrome/               top bar, tab bar, status bar
  views/
    ctx/                per-frame context (Now, Spin) passed to every View()
    overview/           dashboard (5 stat cards in a grid)
    sessions/           sessions list + peek panel; owns its own search input
    issues/             issues table + filter chips + detail
    runtimes/           runtime CLIs (claude/codex/cursor/etc.)
    workspaces/         workspaces list + detail
    logs/               log tail + level filter
    config/             profiles + server/daemon + agents (3-column)
    attach/             attached-session screen (transcript + reply + side panels)
    help/               full-screen help overlay
    palette/            command palette overlay
    onboarding/         multi-step onboarding overlay
  tui/                  root Model — composes all view models, dispatches keys, renders frame
```

### Dependency rule

The import graph is acyclic by design:

```
theme  ←  ui, data, chrome, views/*
app    ←  views/* (for NavigateMsg only)
ctx    ←  views/*
views/* and chrome  ←  tui
tui    ←  main / dump
```

Views never import `tui`. Cross-view navigation is done by emitting `app.NavigateMsg{To, Attach, Overlay}` as a `tea.Cmd`; the root `tui.Model` reads it in `update.go:handleNavigate`.

## How a view is wired in

1. `internal/views/<name>/<name>.go` defines `Model`, `Update(tea.KeyMsg) (Model, tea.Cmd)`, `View(c ctx.Ctx, w, bodyH int) string`, and `KeyHints() []string`.
2. Add the field to `internal/tui/model.go:Model` and initialize it in `New()`.
3. Add a case to `internal/tui/view.go:routeView` and `internal/tui/keys.go` (the per-route Update dispatch).
4. If it should appear in the tab bar, add an entry to `internal/app/routes.go:Routes`.
5. Add a frame to `dump.go:runDump` so it shows up in `--dump`.

## Conventions to keep

- Chrome strips (top/tab/status) are 2 visual rows each — `bodyH := m.H - 6` in `view.go`.
- Every line must be bg-painted, otherwise the terminal default bleeds through dark themes. Use `theme.WithBg` and `ui.PaintLine` (already wired in `paintFrame`).
- The spinner frame and wall clock come from `ctx.Ctx`; don't call `time.Now()` from a view.
- `j/k` move, `g`/`G` jump (vim-style — `gg` is a two-keystroke jump, handled by per-view `PendingG`).
- `enter` activates / drills in; `esc` backs out (closes overlay or detaches).
- Overlay precedence in `view.go:renderBody`: `onboarding` > `help` > body+`palette`.

## Quick commands

```sh
go build ./...                          # whole module
go vet ./...
go run .                                # interactive
go run . --dump                         # render everything to stdout
go run . --help                         # arg/env summary
```
