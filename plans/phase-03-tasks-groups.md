# P3 — Tasks tab: state groups + PR dots

Items: **f** (group tasks by state on Tab 3), **m** (apply the PR-dot helper
added in P0 to task rows).

## Goal

Replace the flat-chronological task list on Tab 3 with state-grouped sections
(Claude Code Agent View pattern: `Needs input → Working → Completed → Failed →
Queued`). Render a right-edge PR dot per row using `prDot` from P0.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `view_tasks.go` (all 187 lines)
3. `data.go` lines 38–56 (the `Task` struct and seed `tasks` slice)
4. `theme.go` (you'll call `statusGlyph`, `statusStyle`, `prDot`, `agentChip`)

That's all.

## Files to modify

| File           | Why                                                                                |
| -------------- | ---------------------------------------------------------------------------------- |
| `data.go`      | extend `Task` with `PRState string` (`"" | "waiting" | "passed" | "merged" | "draft"`) and `PRNum int` |
| `view_tasks.go`| group rows by state in `paneTasksTable`; render the PR dot at right edge          |
| `theme.go`     | (no changes — uses helpers from P0)                                                |

No new files.

## Tasks

1. **`data.go`: extend Task.**
   Append two fields to the `Task` struct:
   ```go
   type Task struct {
       // existing fields...
       PRState string  // "", "waiting", "passed", "merged", "draft"
       PRNum   int     // 0 if none
   }
   ```
   Update seed `tasks`:
   - `t-1284`: PRState `""`, PRNum 0 (still working, no PR yet)
   - `t-1283`: PRState `""`, 0
   - `t-1282`: PRState `""`, 0
   - `t-1281`: PRState `"waiting"`, PRNum 4130   (PR opened, checks running)
   - `t-1280`: PRState `""`, 0
   - `t-1279`: PRState `""`, 0

2. **`view_tasks.go`: group function.**
   Add `groupedTasks() map[string][]Task` returning tasks keyed by group name
   in this *fixed display order*:
   ```
   "Needs input"     ← status == "waiting"
   "Working"         ← status == "working"
   "Ready for review"← PRState != "" && PRState != "merged"
   "Completed"       ← status == "done"
   "Failed"          ← status == "failed"
   "Queued"          ← status == "queued"
   ```
   A task in both "Working" and "Ready for review" goes in "Ready for review"
   (PR status wins). Return `[]string` of section names in display order so the
   caller can iterate deterministically.

3. **`view_tasks.go`: rewrite `paneTasksTable`.**
   Current body is a header row + flat row list. New body:
   - one column-header row (same as today, dim)
   - for each non-empty group in display order:
     - one section header: `hr(name+" ("+count+")", width)`
     - that group's rows via `taskRowSimple`
   - blank-pad to `height` with `bgPadV`.

   Pane title stays `"6 TASKS"` but make it dynamic:
   `fmt.Sprintf("%d TASKS · %d need input · %d working", total, needsInput, working)`.

4. **`view_tasks.go`: PR dot on rows.**
   Modify `taskRowSimple(t Task, selected bool, ...widths)` to:
   - shrink `costW` by 2 cells (give that width to a new `prW = 2`)
   - append at row end: `"  " + prDot(t.PRState)` when `t.PRState != ""`,
     otherwise `"  " + " "` (so alignment stays). The `prDot` helper from P0
     returns a single styled cell — `lipgloss.Width(...)` is 1.

   If row width math is fiddly, easier: leave the existing widths alone and
   tack the dot onto the *right padding* after `padR`'ing the cost column.

5. **`view_tasks.go`: column header.**
   Add a `PR` heading right of `COST`. Right-align dim. Width 2 (`" P"` or just
   `"PR"`).

6. **Update `filterSummary` in P1's main.go** — *only* if you already merged
   P1 — to read `len(groupedTasks()["Needs input"])` so the breadcrumb says
   e.g. `"tasks · 1 needs input · 3 working"`. If P1 is not yet merged, skip
   this; P1's task list will pick it up.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump` tab 3 shows section headers in this order:
      `Needs input (1) · Working (2) · Ready for review (1) · Completed (1) · Failed (1) · Queued (1)`.
      Empty groups do not appear.
- [ ] Row `t-1281` has a green PR dot at right edge (passed=ok); wait, seed
      `t-1281` is `PRState: "waiting"`, so dot is **warn-yellow**. Adjust to
      taste, but match the seed.
- [ ] Rows without a PR have a blank-cell at the right edge — same column
      width — so the table aligns.
- [ ] Pane title reflects counts: e.g. `6 TASKS · 1 NEED INPUT · 2 WORKING`.

## Dependencies

P0 (uses `prDot`). `[P]` with P2 once P0 lands.

## Handoff (paste into next session)

```
P3 complete. data.go: Task gained PRState string + PRNum int (seeded t-1281).
view_tasks.go: paneTasksTable now renders groups (Needs input → Working →
Ready for review → Completed → Failed → Queued); empty groups hidden; counts
in pane title. taskRowSimple renders a right-edge PR dot via prDot().
go build clean. --dump shows the new grouping.
```
