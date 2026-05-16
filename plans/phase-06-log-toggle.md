# P6 — Log Summary toggle

Item: **j** (btop-style toggle for the Summary sidebar on Tab 4).

## Goal

The Summary panel on Tab 4 (lines / errors / warns / sparkline / top sources)
takes ~36 cells of width permanently. Many users live in tail-follow mode and
want the full width for log lines. Bind `s` to toggle Summary visibility.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `view_log.go` (all 125 lines)
3. `main.go` lines 37–95 (the `Update` and `moveSelection` blocks)

That's all.

## Files to modify

| File         | Why                                                                  |
| ------------ | -------------------------------------------------------------------- |
| `view_log.go`| `renderLog` becomes `renderLog(w, h int, showSummary bool)`          |
| `main.go`    | `logSummary bool` on model (default true); `s` toggles it on tab 3   |

No new files.

## Tasks

1. **`view_log.go`: change signature.**
   ```go
   func renderLog(width, height int, showSummary bool) string
   ```
   When `showSummary == true`: existing 2-column layout (table on left,
   `paneLogSummary` on right).
   When `showSummary == false`: full-width log table, no summary pane.

2. **`view_log.go`: title hint when toggled off.**
   On the log table pane, append `· s to show summary` to the existing pane
   `hint` argument when `!showSummary`. Pattern: change `pane(title, "", body,
   w, h, true)` to `pane(title, hint, body, w, h, true)` where
   `hint := "" ; if !showSummary { hint = "s to show summary" }`.

3. **`main.go`: add state + keybinding.**
   - Add `logSummary bool` to `model`. Initialize `true` in `initialModel`.
   - In `Update`, on tab 3 (log), intercept `"s"` to toggle
     `m.logSummary = !m.logSummary`.
   - In `View`, change the tab-3 dispatch to:
     ```go
     body = renderLog(contentW, bodyH, m.logSummary)
     ```

4. **`main.go`: footer key hint.**
   The tab-3 footer (`keys` slice) currently includes `{"t", "hide ticks"}`.
   Add `{"s", "toggle summary"}` next to it. If P1 is merged, this hint shows
   up in the focused-panel footer rule already.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump` tab 4 still shows the 2-column layout (Summary on
      right) because `logSummary` defaults to true. No visual regression.
- [ ] Interactive: `go run .` → press `4` → press `s` → Summary panel
      disappears; log table fills full width.
- [ ] Press `s` again → Summary returns; layout reverts.
- [ ] When Summary is hidden, the log pane top-right hint reads
      `s to show summary`.

## Dependencies

None. `[P]` with everything.

## Handoff (paste into next session)

```
P6 complete. view_log.go: renderLog now takes showSummary bool; full-width
when off. main.go: model.logSummary (default true); s toggles on tab 3.
Footer hint added. No regression in default rendering.
```
