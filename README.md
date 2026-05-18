# autopus-tui

> Terminal UI prototype for the **autopus** agent daemon — monitor and interact with AI coding agents directly from your terminal.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) · Go 1.26 · truecolor theming via [Lip Gloss](https://github.com/charmbracelet/lipgloss)

---

<!-- Add a terminal recording here: vhs/asciinema → README.gif -->

## Views

| Key | View | Description |
|-----|------|-------------|
| `1` | **Overview** | Dashboard — daemon status, env, server health, session counts, recent activity |
| `2` | **Sessions** | Agent-run roster with live state, cost, and peek panel |
| `3` | **Issues** | Issue tracker with filter chips and detail pane |
| `4` | **Runtimes** | Installed agent CLIs (Claude, Codex, Cursor…) |
| `5` | **Workspaces** | Repo workspaces and their active branches |
| `6` | **Logs** | Live log tail with level filter |
| `7` | **Config** | Profiles · server/daemon · agents — three-column layout |
| `↵` | **Attach** | Full attach screen: transcript, reply, side panels |
| `?` | **Help** | Full-screen key reference overlay |
| `:` | **Palette** | Command palette overlay |

## Quick start

```sh
git clone <repo>
cd autopus-tui
go run .
```

Requires Go 1.22+. No external services — all data is static fixtures.

## Navigation

```
1-7         jump to tab
h j k l     move focus (vim-style)
g / G       jump to top / bottom  (gg = top)
enter       drill in / attach
esc         back out / close overlay
/ or f      search / filter (sessions, issues, logs)
: or space  command palette
q           quit
```

## Verify rendering without a terminal

The `--dump` flag renders every view and overlay to stdout, then exits. Use it from CI or Claude Code to confirm visual changes:

```sh
go run . --dump
```

Control the output with env vars:

| Variable | Values | Default |
|----------|--------|---------|
| `AUTOPUS_DUMP_PROFILE` | `truecolor` · `256` · `ascii` | `truecolor` |
| `AUTOPUS_DUMP_WIDTH` | integer | `160` |
| `AUTOPUS_DUMP_HEIGHT` | integer | `44` |
| `AUTOPUS_DUMP_ONLY` | comma-separated view names | all views |

```sh
# ASCII, narrow, attach view only
AUTOPUS_DUMP_PROFILE=ascii AUTOPUS_DUMP_ONLY=attach AUTOPUS_DUMP_WIDTH=120 go run . --dump

# Sessions + palette, 256-color
AUTOPUS_DUMP_PROFILE=256 AUTOPUS_DUMP_ONLY=sessions,palette go run . --dump
```

Recognized names: `overview sessions issues runtimes workspaces logs config attach help palette onboarding`

## Project layout

```
main.go                 entrypoint — arg parsing, launches tea.Program or runDump
dump.go                 --dump driver

internal/
  app/                  NavigateMsg, ClockMsg, SpinMsg + route table
  theme/                colors, glyphs, lipgloss styles, bg-leak helpers
  ui/                   primitives: text, frame, panel, chip, kv, bars
  data/                 static fixtures (sessions, issues, runtimes, …)
  chrome/               top bar, tab bar, status bar
  views/
    ctx/                per-frame context (Now, Spin) passed to every View()
    overview/           dashboard grid
    sessions/           list + peek panel + search
    issues/             table + filter chips + detail
    runtimes/           runtime CLI cards
    workspaces/         list + detail
    logs/               tail + level filter
    config/             3-column: profiles · server/daemon · agents
    attach/             transcript + reply + side panels
    help/               full-screen overlay
    palette/            command palette overlay
    onboarding/         multi-step onboarding overlay
  tui/                  root Model — composes all views, dispatches keys, renders frame
```

**Dependency rule** — the import graph is acyclic by design:

```
theme  ←  ui, data, chrome, views/*
app    ←  views/*
ctx    ←  views/*
views/*, chrome  ←  tui  ←  main / dump
```

Views never import `tui`. Cross-view navigation is done by emitting `app.NavigateMsg`; the root `tui.Model` handles it in `update.go`.

## Adding a view

1. Create `internal/views/<name>/<name>.go` with `Model`, `Update`, `View`, and `KeyHints`.
2. Add the field to `internal/tui/model.go` and initialize it in `New()`.
3. Register a case in `internal/tui/view.go:routeView` and `internal/tui/keys.go`.
4. Append an entry to `internal/app/routes.go:Routes` to add a tab.
5. Add a frame to `dump.go:runDump`.

## Common commands

```sh
go build ./...      # compile
go vet ./...        # static analysis
go run .            # interactive TUI
go run . --dump     # render all views to stdout
go run . --help     # flag/env reference
```

## Tech stack

- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) — Elm-inspired TUI framework
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) — layout and styling
- [muesli/termenv](https://github.com/muesli/termenv) — terminal color profile detection
