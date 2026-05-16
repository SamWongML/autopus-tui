# P2 — Status tab: pulse + runtime utilization

Items: **a** (demote mini-tasks + mini-log; replace with system-pulse content),
**e** (runtime panel becomes a thin utilization strip; remove duplication from
Profiles in P5).

## Goal

Make Status show information that *only* Status can usefully show: a system
pulse (throughput sparkline, error-rate trend, last-events ribbon) and a thin
runtime utilization strip. Stop re-rendering heads of the Tasks and Log tabs.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `view_status.go` (all 297 lines)
3. `view_log.go` lines 99–125 (the `sparklineBars` function — reuse it)
4. `data.go` (all of it; you'll add two seed series)

That's all.

## Files to modify

| File             | Why                                                                                  |
| ---------------- | ------------------------------------------------------------------------------------ |
| `view_status.go` | rewrite `renderStatus` body; replace `paneTasksInFlight` and `paneLogTail`; reshape `paneRuntimes` |
| `data.go`        | add `tasksPerHour []int` (24 values) and `errsPerHour []int` (24 values)             |

No new files.

## Tasks

1. **`data.go`: add throughput data.**
   ```go
   var tasksPerHour = []int{ /* 24 ints, mock values 0–10 */ }
   var errsPerHour  = []int{ /* 24 ints, mock values 0–4 */ }
   ```
   Use varied values that read well as a sparkline (e.g., a rising curve with a
   small dip).

2. **`view_status.go`: rewrite the 4-panel grid.**
   Current layout (in `renderStatus`):
   ```
   ┌ DAEMON ──────┐  ┌ TASKS · IN FLIGHT ─┐
   │              │  │ ... tasks roster... │   ← drop
   ├ RUNTIMES ───┤  ├ DAEMON LOG ─────────┤
   │              │  │ ... log tail ...    │   ← drop
   └──────────────┘  └─────────────────────┘
   ```
   New layout (still 2-column grid, still uses `pane(...)` + `joinV/joinH`):
   ```
   ┌ DAEMON ────────────┐  ┌ PULSE · today ─────────────┐
   │ (paneDaemon, unchanged) │ tasks/h sparkline         │
   │                         │ errors/h sparkline        │
   │                         │ poll & heartbeat tickers  │
   └─────────────────────┘  │ (moved from DAEMON pane)  │
   ┌ RUNTIMES ──────────┐  ├─ LAST EVENTS ──────────────┤
   │ claude ████░ 2/3   │  │ 5 most-recent rows mixing  │
   │ codex  ██░░░ 1/4   │  │ task transitions + warnings│
   └─────────────────────┘  └────────────────────────────┘
   ```

3. **`view_status.go`: shrink `paneDaemon`.**
   Currently it renders 12 kv rows *plus* the TICKERS block at the bottom. Move
   the tickers (`poll` and `heartbeat`) to the new Pulse pane. Daemon pane keeps
   just identity rows: state, pid, uptime, version, server, connection, device,
   mem, cpu, socket, log, workspaces.

4. **`view_status.go`: write `panePulse(width, height int) string`.**
   Body, in order, separated by `hr(label, w)`:
   - `kv("tasks/h", fmt.Sprintf("%d today", sum(tasksPerHour)), "", w)`
   - one row: `sparklineBars(tasksPerHour, w-2)` styled `sFg1`
   - `kv("errors/h", fmt.Sprintf("%d today", sum(errsPerHour)), errTone(), w)`
   - one row: `sparklineBars(errsPerHour, w-2)` styled `sErr` if total > 0
     else `sDim`
   - `hr("TICKERS", w)`
   - the two ticker rows (`ticker(...)` already exists in `view_status.go` —
     just move the calls here)

5. **`view_status.go`: write `paneLastEvents(width, height int) string`.**
   Body: 5 most-recent rows derived from a merge of `tasks` (status transitions
   you can fake by reading task status) and `logLines` filtered to `warn`/`err`.
   For the mock, just produce 5 hardcoded one-line entries like:
   ```
   ▶ t-1284  started · 3m 12s ago
   ! t-1283  waiting on user · 12s ago
   ✗ t-1280  failed · bundler resolver · 21m ago
   ⚠ runtime.codex  Bundler::VersionConflict
   ✓ t-1281  PR #4130 opened · 8m ago
   ```
   Each row uses `statusGlyph` + `statusStyle` for the leading glyph; rest is
   `sDim`. Pad to `width` and to 5 lines (use `bgPadV` for the remainder).

6. **`view_status.go`: reshape `paneRuntimes`.**
   Currently each runtime takes 3+ lines (name + bin/tasks). Compress to one
   row per runtime:
   ```
   ▎ ● claude  ████▒░░  2/3   sonnet-4.6   47 tasks · 1 err
   ▎ ● codex   ██░░░░░  1/4   gpt-5.5      19 tasks · 0 err
   ```
   `bar(float64(r.Busy)/float64(r.Max), 7)` gives the bar. No version string,
   no full bin path — those live in Config (tab 5).

7. **`view_status.go`: wire it.**
   `renderStatus` now composes:
   ```go
   left  := joinV(paneDaemon(leftW, leftH), paneRuntimes(leftW, restH))
   right := joinV(panePulse(rightW, pulseH), paneLastEvents(rightW, restH))
   return joinH(left, right)
   ```
   Delete `paneTasksInFlight` and `paneLogTail` (and their helpers `taskRow`
   and `logRowText` *if not used elsewhere* — check `view_tasks.go` and
   `view_log.go` first; `taskRowSimple` lives in `view_tasks.go` separately).

   Specifically: `taskRow` is *only* used by `paneTasksInFlight` (Status's mini
   list) — safe to delete. `logRowText` is *only* used by `paneLogTail` — same.
   `paneLogTable` in `view_log.go` has its own row renderer.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump` tab 1 shows: DAEMON + RUNTIMES on left, PULSE + LAST
      EVENTS on right.
- [ ] No "TASKS · IN FLIGHT" panel anywhere on tab 1.
- [ ] No "DAEMON LOG · TAIL" panel anywhere on tab 1.
- [ ] The tasks roster on tab 3 is unchanged.
- [ ] The log tail on tab 4 is unchanged.
- [ ] Runtimes panel on tab 1 fits both runtimes in 4 visible lines (header +
      2 rows + spacer).

## Dependencies

P0 (uses no new glyphs from P0 directly, but `dirtyMark` is referenced in
"errors today" if you want to highlight non-zero errors — optional). Can start
once P0 is merged. `[P]` with P3 (different file) once P0 lands.

## Handoff (paste into next session)

```
P2 complete. view_status.go reshaped: DAEMON+RUNTIMES (left), PULSE+LAST EVENTS
(right). Deleted paneTasksInFlight, paneLogTail, taskRow, logRowText.
Added: panePulse, paneLastEvents in view_status.go.
data.go: added tasksPerHour[24], errsPerHour[24].
Tab 3 (tasks) and tab 4 (log) outputs identical to before — confirmed via dump
diff scoped to those tabs.
```
