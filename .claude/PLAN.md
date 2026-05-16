# mtop Implementation Plan

Status key: `[ ]` pending · `[~]` in-progress · `[x]` done · `[!]` blocked

This is the durable, cross-session source of truth for the mtop refactor.
Edit checkboxes in place as work completes. Append to **Rules & Tips** when
something is learned. Never delete Rules entries.

## Phase Summary

| Phase | Name                              | Status | Branch                  |
|-------|-----------------------------------|--------|-------------------------|
| A     | De-duplicate widgets (mechanical) | `[x]`  | `refactor/dedup`        |
| B     | Adopt Bubbles library             | `[x]`  | `feat/bubbles`          |
| C     | KPI dashboard Status tab          | `[x]`  | `feat/kpi-dashboard`    |
| D     | Wire real `/health` data source   | `[ ]`  | `feat/health-wire`      |
| E     | Focus / scroll model              | `[ ]`  | `feat/focus-scroll`     |
| F     | Polish, palette, golden tests     | `[ ]`  | `feat/polish`           |

Phases are mostly sequential. A unblocks B (smaller surface to port). B unblocks
C/E (viewport, table, focus come for free). D is independent of B/C but should
come after A so the new `internal/source/` package isn't ported twice. F is last.

---

## Phase A — De-duplicate widgets

> **Goal:** extract shared primitives. **No behavior change.** Target ~−600 LOC.
> **Done when:** `go build ./...` clean, golden diff at 80×24 / 120×40 / 200×60 is empty.

- [x] A.1 Extract `pageHeader(width, left, right string)` into `chrome.go`; replace 4 callers (`renderFilterBar`, `renderLogFilter`, `renderConfigHead`, `renderProfilesHead`)
- [x] A.2 Extract `chip(label string, active bool)` + `chipBar(...)` into `chrome.go`; remove 2 local closures
- [x] A.3 Extract `taskTable(rows []task, w, h int, mode {compact,full}, sel int)` into `widget_table.go`; replace `taskRow` (view_status.go:165) and `taskRowSimple` (view_tasks.go:91)
- [x] A.4 Extract `logTable(lines []logLine, w, h int, mode)` into `widget_logtable.go`; replace `paneLogTail` and `paneLogTable`
- [x] A.5 Extract `kvPane(title string, sections []kvSection)` into `widget_kvpane.go`; replace handwritten `paneDaemon`, `renderTaskSidebar`, `paneTaskPreview`, `renderProfileCard`
- [x] A.6 Move `countBar`, `wrap`, `truncate` into `chrome.go`; replace `truncate`/`wrap` bodies with `github.com/charmbracelet/x/ansi` `Truncate` / `Wrap`
- [x] A.7 Capture golden snapshots: `MTOP_DUMP_PROFILE=1 ./mtop --dump > testdata/<size>.txt` for each of 80×24, 120×40, 200×60 (use **pre-refactor commit** as baseline, then re-run after each task A.1–A.6 — diff must be empty)
- [x] A.8 Tally LOC delta: `git diff --stat main...HEAD` and record in Rules

## Phase B — Adopt Bubbles

> **Goal:** replace hand-rolled list/input/help with `charmbracelet/bubbles`. Another ~−400 LOC, free scroll/animation/focus.
> **Done when:** all 4 widgets adopted, keymap unchanged from user perspective, goldens regenerated.

- [x] B.1 `bubbles/viewport` for Run messages pane and Log tab scroll region; map `j/k/ctrl-d/ctrl-u/g/G` to viewport
- [x] B.2 `bubbles/textinput` for the command bar (`:` and `/` prompts in `chrome.go`)
- [x] B.3 `bubbles/help` for the footer; load keymap from a single `keys.go`
- [x] B.4 `bubbles/table` for the Tasks tab; configure columns from `taskTable` mode=full
- [x] B.5 `bubbles/spinner` for live indicators on status pills (`● online`, in-flight tasks)
- [x] B.6 Regenerate goldens; visual review at 3 sizes

## Phase C — KPI dashboard Status tab

> **Goal:** redesign Status as a real dashboard, not a miniature copy of every other tab. Removes pressure for two-size widgets.
> **Done when:** Status renders 4 KPI tiles + 2 sparklines + recent-failures feed at 120×40; no tables.

- [x] C.1 Decide tile layout: 4 KPIs (active tasks, runtimes online, tokens today, errors today)
- [x] C.2 Build `widget_kpi.go` — big-number tile with delta arrow
- [x] C.3 Build `widget_sparkline.go` — single sparkline renderer reused for tasks/hr and lines/min
- [x] C.4 Build "recent failures" feed (max 5, click-through to Run tab)
- [x] C.5 Delete `paneLogTail` and the embedded mini-task-table from Status (no longer needed after A.4 + this phase)
- [x] C.6 Update goldens

## Phase D — Wire real `/health` data source

> **Goal:** stop hardcoding `data.go`. Real HTTP + log tail. Mock gated by `MTOP_MOCK=1` for dev/screenshots/tests.
> **Done when:** running daemon → mtop shows live data; daemon down → "● connecting…" state, no crash.

- [ ] D.1 `internal/source/source.go` — `type Source interface { Snapshot(ctx) (Snapshot, error); Stream(ctx) <-chan Event }`
- [ ] D.2 `internal/source/health.go` — poll `http://127.0.0.1:<port>/health` every 1s, emit `snapshotMsg` as `tea.Msg`
- [ ] D.3 `internal/source/runmessages.go` — wraps `multica issue run-messages <id> --since <seq> --output json`, polls every 500 ms while Run tab is active
- [ ] D.4 `internal/source/logtail.go` — `tail -F ~/.multica/daemon.log`, parses `level`/`src`/`msg`
- [ ] D.5 `internal/source/mock.go` — move current `data.go` contents here, gated by `MTOP_MOCK=1`
- [ ] D.6 Remove package-level `var daemon` from `data.go`; model holds `Snapshot` instead; views read from model
- [ ] D.7 Empty/loading/error states: show `● connecting…` when health request fails; retry with backoff
- [ ] D.8 Integration smoke test against a running daemon (manual checklist; add to Rules)

## Phase E — Focus / scroll model

> **Goal:** multi-pane focus with `Tab`/`Shift-Tab` cycling; wire `pane()`'s `accentBorder` to focus.
> **Done when:** every multi-pane tab has a `focusedPane` field; Tab cycles; accent border visible on focused pane; scrolling routed to focused pane.

- [ ] E.1 Add `focusedPane int` per tab to model (replace global `selTask` / `selCfg` / `selProf`)
- [ ] E.2 Map `Tab` / `Shift-Tab` to cycle focus on the active tab
- [ ] E.3 Wire `pane(..., accentBorder bool)` to `focusedPane` (already takes the param, never set)
- [ ] E.4 Route `j/k/g/G/ctrl-d/u` to the focused pane's viewport
- [ ] E.5 Update goldens with focused-pane border on Status (default), Tasks, Log

## Phase F — Polish, palette, persistence, tests

> **Goal:** production-grade rough edges. Pick from menu; not all items need to ship.
> **Done when:** items below are either checked or explicitly deferred (move to Rules with reason).

- [ ] F.1 Command palette: `:tasks status=working`, `:kill t-1284`, `:profile staging` actually dispatch via `bubbles/textinput`
- [ ] F.2 Search: `/` with regex + incremental highlight in Log and Tasks
- [ ] F.3 Mouse: click tabs to switch, click rows to select (already enable `tea.WithMouseCellMotion`, ignore events)
- [ ] F.4 Resize robustness: < 40×12 currently bails; collapse panes (drop sidebar → drop preview) before bailing; fix Profiles grid at narrow widths
- [ ] F.5 Persistence: `~/.multica/mtop.state.json` for last tab, last filter, last sort
- [ ] F.6 `?` help overlay using `bubbles/help`
- [ ] F.7 Theming: load palette from `~/.multica/mtop.toml`; respect `NO_COLOR`
- [ ] F.8 Golden tests in CI: snapshot at 80×24 / 120×40 / 200×60 via `--dump`, diff on PR
- [ ] F.9 `MTOP_LOG=~/.mtop.log` writes via `log.Default()` (stderr is taken by Bubbletea)
- [ ] F.10 If `!isatty(stdout)`, exit with hint to use `--dump` (no alt-screen into a pipe)
- [ ] F.11 `version` subcommand + `-ldflags` embed of git SHA
- [ ] F.12 Write repo README

---

## Rules & Tips

<!-- Append discoveries. Never delete. Date entries when non-obvious. -->

- Phase A LOC delta (2026-05-15): code went **+55 lines** (1868 → 1923 across `chrome.go main.go view_*.go widget_*.go`). The −600 target was optimistic; the mechanical extraction collapses ~30-line repetitions into ~25-line shared functions with mode branches, plus larger doc comments — net wash. Real wins are in maintainability (one source of truth for table/log/section rendering) and unblocking Phase B (Bubbles drop-in points isolated).
- Phase A goldens: `testdata/{80x24,120x40,200x60}.txt` captured at HEAD, byte-identical to pre-refactor (`refactor/dedup` baseline). Re-generate with `for sz in 80x24 120x40 200x60; do W=${sz%x*}; H=${sz#*x}; MTOP_DUMP_PROFILE=truecolor MTOP_DUMP_WIDTH=$W MTOP_DUMP_HEIGHT=$H ./mtop --dump > testdata/$sz.txt; done`
- Dump sizing: `MTOP_DUMP_WIDTH` / `MTOP_DUMP_HEIGHT` env vars override the model's default 160×48 in the `--dump` code path (added in A baseline capture).
- A.6 deferred swap: `ansi.Truncate` / `ansi.Wrap` produce visually identical output but emit redundant SGR transitions, breaking byte-equal goldens. Move happened; body swap deferred to Phase B with a regenerated goldens diff.
- Build: `go build -o mtop ./cmd/mtop` (verify path — currently `main.go` may live at repo root, not `cmd/mtop`)
- Run with mock data: `MTOP_MOCK=1 ./mtop` (gating exists once D.5 lands)
- Dump goldens: `MTOP_DUMP_PROFILE=1 ./mtop --dump > testdata/<W>x<H>.txt`
- ANSI utils: `github.com/charmbracelet/x/ansi` (`Truncate`, `Wrap`) — already a transitive dep, no go.mod change needed
- Bubbles import path: `github.com/charmbracelet/bubbles/{viewport,textinput,help,table,spinner}`
- Daemon health endpoint: `http://127.0.0.1:7717/health` returns `HealthResponse` JSON (see `multica/server/internal/daemon/health.go:17`): `status`, `pid`, `uptime`, `daemon_id`, `device_name`, `server_url`, `cli_version`, `active_task_count`, `agents`, `workspaces`
- Daemon log path: `~/.multica/daemon.log`
- Daemon socket: `~/.multica/daemon.sock` (currently unused by mtop — `/health` HTTP is the integration point)
- Mtop debug screenshot: `MTOP_DUMP_PROFILE=1` already detected by current code
- Phase B LOC delta (2026-05-16): code went **+353 lines** net (+174 in 5 existing files, +179 across new `keys.go`/`spinner.go`/`widget_tasks_table.go`). The −400 LOC target was optimistic for the same reason as A: Bubbles widgets save behavior code (free scroll, blink, cursor) but need glue code (resize wiring, key routing, factories, per-widget config). The wins are functional, not LOC: viewports actually scroll, the cmd bar accepts text, spinners animate, the table cursor navigates.
- Phase B golden diffs (intentional, captured in `testdata/{80x24,120x40,200x60}.txt`):
  - `liveSpin` (bubbles/spinner Pulse) replaces the static `●` in the tab-bar daemon-online dot and the Run-tab LIVE pill. Dump captures frame 0 = `█`.
  - Tasks-tab roster is now `bubbles/table`: header is uppercased and dim (no bold), status column is the glyph alone (no per-row color — bubbles/table runs `runewidth.Truncate` on raw cell values, which corrupts pre-rendered ANSI), selected row uses a solid `bgSel` background instead of the accent `▎` bar.
  - Footer keycaps now come from `bubbles/help.ShortHelpView`; literal padding spaces inside `key.WithHelp` keep the pre-Phase-B pill-shaped look since `Inline(true)` strips `style.Padding`.
  - Command bar prompt is a `bubbles/textinput.View()` (always renders a 1-cell caret/cursor instead of the static `█` placeholder used pre-B).
- Bubbles table styling caveat: `table.Model.renderRow` calls `runewidth.Truncate` on the raw row value *before* applying `Styles.Cell`/`Styles.Selected`. Pre-styled cells get their ANSI sequences width-counted and split mid-CSI. Per-row colors are therefore not viable inside `bubbles/table`; either pick column-uniform `Styles.Cell.Foreground(...)` or accept plain text. Selected styling can preserve underlying foregrounds by setting only `Background(bgSel)` on `Styles.Selected` (default `DefaultStyles().Selected` paints a foreground that erases content tone — override it).
- Bubbles table widths: `Styles.Cell.Padding(...)` is *not* Inline-stripped (the inner `lipgloss.NewStyle().Width().MaxWidth().Inline(true)` only wraps the value; `Cell` wraps outside that). Padding therefore widens each column by 2; easier to bake a trailing space into the cell value (or column title) so column gaps don't break `SetWidth` math.
- Spinner frame in dumps: `liveSpin` is a package-level `spinner.Model` initialized in `spinner.go:init()`. In `--dump`, no `Update` ever runs, so `View()` always returns frame 0. That makes the snapshot deterministic; a different spinner.Style is needed if frame 0 should not be `█` (bubbles/spinner Pulse).
- Tasks-tab key routing: `j/k` on tab 2 forward to `m.tasksTbl.Update(msg)` rather than `moveSelection(±1)`; `g`/`G` call `GotoTop`/`GotoBottom`. Selection on Status tab still uses `selTask` against the manual compact `taskTable`.
- Phase C dashboard layout (2026-05-16): Status tab body splits as `kpiH : sparkH : failH = 21% : 36% : remainder` of bodyH, clamped to a 5-row min on each band. At 120×40 (bodyH 34) that's 7 / 12 / 15; at 80×24 (bodyH 18) that's 5 / 6 / 7; at 200×60 (bodyH 54) that's 11 / 19 / 24.
- KPI tile titles must be ≤ 11 cells wide for the smallest tile (width 18 at 80-cell viewport) to avoid pane top-border truncation. `buildTopBorder` in chrome.go truncates when `right < 1`, i.e. when `labelWidth > width - 7`. Single-word titles ("active", "online", "tokens", "errors") render cleanly; descriptive context (counts, deltas) lives in the body's value / unit / delta lines.
- KPI tile content centering: `lipgloss.Place(w, h, Center, Center, body)` centers the rendered block as a unit using the longest line as block width, so shorter lines are left-aligned within the block. To get per-line centering, pre-wrap each row in `lipgloss.NewStyle().Width(inner).Align(Center).Background(bg)` before Place. Without that, a 2-line tile where line 2 is wider (e.g. "→ claude · codex" under "2 / 2") flush-lefts line 1.
- Sparkline stat bar fallback: at narrow widths the four-stat row `min · avg · max · now` overflows inner. `sparklinePane` measures each variant against `inner - 4` and picks the widest that fits (full → 3-stat short → "now N" only). 80-cell viewport falls through to the short variant for `linesPerMin` (its larger numbers push width past the 36-cell band).
- `bgPadV(width, height)` is the right primitive for a vertical strip of bg-painted blank cells; `bgPad(n)` only produces a single-row gap. `kpiRow` uses `bgPadV(gap, height)` between tiles to avoid an unpainted seam.
- Phase C LOC delta (2026-05-16): code went **+374 lines** net. `view_status.go` shrinks 143 → 125 (−18; paneDaemon/paneRuntimes/paneTasksInFlight/paneLogTail/ticker/commafy all deleted, replaced by `statusKPIs`/`sparkRow`/`failuresPane`), `data.go` grows 155 → 189 (+34; sparkline series + failures fixture), and three new widget files add `widget_kpi.go` (111), `widget_sparkline.go` (191), `widget_failures.go` (56). The wins here are visual, not size — the dashboard reads at a glance, and `Snapshot` from Phase D drops in by replacing the `statusKPIs` reads of `tasks` / `runtimes` / `failures` with live values.
- Phase C golden diffs (intentional, captured in `testdata/{80x24,120x40,200x60}.txt`): all three goldens differ only inside the tab-1 block. tabs 2–6 are byte-identical. The Status body changed from `daemon/runtimes/in-flight-tasks/log-tail` quad-pane → KPI tile row + sparkline row + recent-failures pane. Tab-bar meta still surfaces daemon identity (pid, uptime), so the daemon pane removal didn't lose any header-bar info.

## Compact Instructions

When `/compact` runs, preserve verbatim:
1. The **Phase Summary** table (all rows, status column)
2. Every `[~]` and `[!]` item under any phase, with its surrounding two `[ ]` siblings for context
3. The entire **Rules & Tips** section
4. The currently active branch and worktree path

It is OK to summarize completed (`[x]`) tasks as a count per phase.
