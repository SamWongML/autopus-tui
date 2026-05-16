# P4 — Unified detail pane + peek overlay

Items: **c** (unify Run + Tasks sidebars into one component), **h** (Run's
sidebar becomes a Space-toggleable peek overlay instead of a permanent column).

## Goal

Today there are two near-duplicate detail panels for a task:
- `view_run.go: renderTaskSidebar(width)` shows status, runtime, model, pid,
  started, elapsed, messages, tokens, cost, branch + files touched + msg-type
  counts.
- `view_tasks.go: paneTaskPreview(width, height)` shows status, runtime, pid,
  workspace, cwd, issue, branch, last, last-5 messages, actions.

Both repeat: status, runtime, pid, branch. Both belong to the same task.

Replace with one component that has a shared header and two body modes
(`stream` for Run, `summary` for Tasks). On Run, the panel is hidden by default
and only appears when the user presses Space (peek). On Tasks, it stays visible
on the right as today.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `view_run.go` (all 259 lines)
3. `view_tasks.go` (all 187 lines)
4. `data.go` (you need `Task`, `Msg` shapes)
5. `chrome.go` lines 230–310 (`pane` and `buildTopBorder`)

That's all.

## Files to modify / add

| File                  | Why                                                                      |
| --------------------- | ------------------------------------------------------------------------ |
| **`view_detail.go`**  | **new file** — `renderTaskDetail(t Task, mode string, w, h int) string`  |
| `view_run.go`         | delete `renderTaskSidebar`, `fileRow`, `countBar`; call `renderTaskDetail` |
| `view_tasks.go`       | delete `paneTaskPreview`, `msgRowCompact`, `actionRow`; call `renderTaskDetail` |
| `main.go`             | add `peek bool` to model; Space toggles it on tab 1                      |

## Tasks

1. **Create `view_detail.go` with one exported function.**
   ```go
   func renderTaskDetail(t Task, mode string, width, height int) string
   ```
   `mode` is `"stream"` or `"summary"`. Body composition:

   **Shared header (always):**
   - title: pane title is `"TASK " + t.ID`, hint = `t.Title` (truncated to ~32
     cells)
   - kv rows: status (with `statusGlyph` + `statusStyle`), runtime (with
     `agentChip`), model (claude→sonnet-4.6, codex→gpt-5.5 — derive from
     `runtimes` slice or hardcode for the mock), pid, branch, cost, elapsed.

   **mode == "summary"** (used by Tasks tab):
   - `hr("LAST 5 MESSAGES", w)`
   - 5 lines from the global `msgs` slice tail, using a compact one-line format:
     `T  Type  Body[:width-marshalling]` styled `sDim` for `Type` and `sFg1`
     for body.
   - `hr("ACTIONS", w)`
   - rows: `⏎ open run viewer`, `r reply`, `k kill (SIGTERM)`,
     `K force kill`, `o open workspace`, `c copy id`.

   **mode == "stream"** (used by Run peek):
   - `hr("FILES TOUCHED", w)`
   - 3 file rows (hardcode for the mock; the real implementation will read
     `t.FilesTouched` once that field exists):
     `db/migrate/20260514_add_idemp… +42`
     `app/models/webhook_event.rb +8 −2`
     `spec/models/webhook_event_spec… +24`
   - `hr("MSG TYPE COUNTS", w)`
   - bars: `thinking 32 · tool_call 28 · tool_result 20 · text 5 · error 1`,
     each rendered with `bar(n/max, w-10)`.

   Pad to `height` with `bgPadV` so the pane always fills its slot.

2. **Wire `view_tasks.go` to call it.**
   In `renderTasks`, where it currently calls `paneTaskPreview(rightW, bodyH)`,
   replace with:
   ```go
   right := pane("PREVIEW · "+tasks[selected].ID, "",
                 renderTaskDetail(tasks[selected], "summary", rightW-2, bodyH-2),
                 rightW, bodyH, true)
   ```
   Delete `paneTaskPreview`, `msgRowCompact`, `actionRow` from `view_tasks.go`.

3. **Wire `view_run.go` to call it as an *overlay*.**
   Currently `renderRun` produces a 2-column layout (messages + sidebar). Make
   it 1-column normally and 2-column when `peek` is true. To pass `peek` from
   `model` into `renderRun`, change the signature:
   ```go
   func renderRun(width, height int, peek bool) string
   ```
   and pass `m.peek` from `main.go: View()`.

   When `peek == false`: full-width messages pane, no sidebar.
   When `peek == true`: messages on left at `width - 38`, detail on right at
   width 38, using `renderTaskDetail(tasks[0], "stream", 36, height-2)`. (For
   the mock there's no real selection on Run — always show the active task,
   which is `tasks[0]`/t-1284.)

   Delete `renderTaskSidebar`, `fileRow`, `countBar` from `view_run.go`. The
   `maxInt` helper at the bottom can stay if other code uses it; check first.

4. **`main.go`: add peek state.**
   Add `peek bool` to `model`. In `Update`, handle `" "` (space) on tab 1
   (run): toggle `m.peek`. Default `false`.

5. **`main.go`: footer hints reflect peek.**
   If P1 is merged, extend the tab-1 keys: add `{"␣", "peek"}` when
   `!m.peek`, and `{"␣", "close peek"}` when `m.peek`.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump` tab 2 (run) shows ONE column (full-width messages).
      The right sidebar that used to show "TASK t-1284" with status/runtime/
      tokens/etc. is GONE by default.
- [ ] Interactive: `go run .` → press `2` → press space → the peek panel
      slides in on the right showing the same TASK details.
- [ ] Press space again → peek closes; messages reflow full-width.
- [ ] `go run . --dump` tab 3 (tasks) PREVIEW pane shows the summary-mode body
      (status/runtime/branch + LAST 5 MESSAGES + ACTIONS). Header fields
      identical to the Run peek (no divergence).
- [ ] `view_detail.go` is the only place that knows a task's identity fields.
      `grep -n 'kv.*"pid"' view_*.go` returns hits only in `view_detail.go`.

## Dependencies

None on prior phases (but coordinate with P1 if both edit `main.go`'s `Update`
keymap — easy to merge, separate keys).

## Handoff (paste into next session)

```
P4 complete. New file view_detail.go: renderTaskDetail(t, mode "stream|summary",
w, h). view_run.go: deleted renderTaskSidebar/fileRow/countBar; renderRun now
takes peek bool and renders 1-col default / 2-col when peek. view_tasks.go:
deleted paneTaskPreview/msgRowCompact/actionRow; calls renderTaskDetail in
"summary" mode. main.go: model.peek; Space toggles on tab 1. Both detail panes
share identical header fields — duplication removed.
```
