# P7 — Tab consolidation (6 → 4)

Item: **g** (reduce the tab count).

## Goal

After P2 (Status now owns the daemon/runtime/pulse) and P5 (Profiles becomes a
switcher), the Profiles tab no longer needs its own slot — it can live as a
horizontal control on Status. After P4 (Run is a peek/attach over Tasks), the
Run tab is also redundant as a top-level destination — `Enter` on a task row
becomes the attach, exactly like Claude Code Agent View.

Final tab set: **tasks · run · log · config** — where "run" is the attached
view of a single task, reached by `Enter` from tasks, not by a number key.

Two sub-options for representation:
- (A) Keep 4 numbered tabs `tasks · run · log · config`; "run" is empty until
  attached, like a docking slot. Pressing `Enter` on a task fills it.
- (B) Three top-level tabs `tasks · log · config`; "run" is a *modal* full-
  screen takeover, no tab number. Esc to return.

Pick (B). It mirrors Agent View more closely — attaching takes over the screen.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `main.go` (all 243 lines)
3. `chrome.go` lines 334–395 (`tabBar`, `footer`)
4. `view_status.go` lines 12–40 (`renderStatus` — you'll fold the profile
   switcher row in here)
5. `view_profiles.go` (after P5 — the `renderProfileSwitcher` function)

Skip view_run.go, view_tasks.go, view_log.go, view_config.go — their internals
don't change in this phase.

## Files to modify

| File              | Why                                                       |
| ----------------- | --------------------------------------------------------- |
| `main.go`         | reduce tab count; introduce `attached bool`; route Esc/← |
| `chrome.go`       | `tabBar` accepts 4 tab names                              |
| `view_status.go`  | host the profile switcher row at the top                  |
| `view_profiles.go`| keep `renderProfileSwitcher` exported; the rest can stay  |
|                   | but isn't reachable from the tab bar                      |

**Renames / deletions:**
- The Profiles "tab" goes away. Don't delete `view_profiles.go`; just stop
  routing to it. Its `renderProfileSwitcher` is reused by Status.
- The Run "tab" goes away. `renderRun` is invoked when `m.attached == true`,
  replacing the whole body.

## Tasks

1. **`chrome.go: tabBar` — drop to 4 labels.**
   Change the slice to `[]string{"tasks", "log", "config", "status"}`. Yes,
   Status moves to position 4 — Tasks is the primary entry per the Agent View
   philosophy. Number keys: `1 tasks · 2 log · 3 config · 4 status`.
   (If you prefer Status first, swap; just keep Tasks reachable as `1`.)

2. **`main.go`: rebind tab numbers.**
   ```
   "1" → m.tab = 0  (tasks)
   "2" → m.tab = 1  (log)
   "3" → m.tab = 2  (config)
   "4" → m.tab = 3  (status)
   ```
   Delete handling for "5" and "6".

3. **`main.go`: add attach state.**
   - `attached bool` (default false). When true, the body is `renderRun(w, h,
     m.peek)` regardless of `m.tab`.
   - On the tasks tab, `"enter"` sets `m.attached = true`.
   - In any state, `"esc"` (and `"left"` on an empty input — out of scope for
     mock) sets `m.attached = false`. The existing `"esc"` → `tea.Quit` must
     change: `"esc"` now goes back to tasks if attached, only quits if `q` /
     `ctrl+c`.

4. **`main.go`: View() dispatch.**
   ```go
   if m.attached {
       body = renderRun(contentW, bodyH, m.peek)
       keys = []KeyHint{
           {"esc", "back to tasks"}, {"j/k", "scroll"}, {"f", "follow"},
           {"␣", "peek"}, {"y", "yank seq"}, {"r", "reply"}, {"k", "kill"},
       }
       cmdMode = "reply"
       cmdPlaceholder = "reply inline — ⏎ send · esc detach"
   } else {
       switch m.tab {
           case 0: // tasks (was 2)
           case 1: // log    (was 3)
           case 2: // config (was 4)
           case 3: // status (was 0; integrated profile switcher)
       }
   }
   ```

5. **`view_status.go`: fold in the profile switcher.**
   Above the existing 2x2 grid (DAEMON, RUNTIMES, PULSE, LAST EVENTS), insert
   one line containing `renderProfileSwitcher(width, m.selProf)`. `renderStatus`
   gains a `selProf int` argument; `main.go` passes it.

   This means Status now shows: switcher + DAEMON + RUNTIMES + PULSE + LAST
   EVENTS. The switcher gets `↑↓` and `←/→` no longer; it's a passive header.
   (Switching profiles becomes a `:profile <name>` command — out of scope for
   mock.)

6. **`main.go`: breadcrumb update (if P1 merged).**
   - tab 0 (tasks): same as before
   - tab 1 (log): same
   - tab 2 (config): same
   - tab 3 (status): `""` (Status owns its own header now)
   - attached: `"task › " + tasks[selTask].ID + " · ws " + tasks[selTask].WS`

7. **Cleanup.**
   - Confirm `view_profiles.go` still compiles (it isn't routed but its
     functions are referenced by `view_status.go`). If `renderProfiles` /
     `renderProfileDetail` are unreachable, leave them — they could be
     re-exposed via a `:profile detail <name>` command. If you'd rather
     prune, delete everything except `renderProfileSwitcher` and `stateHint`.
   - `view_run.go` is reached only via attach; its signature remains
     `renderRun(w, h int, peek bool)` from P4.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump` shows 4 tabs in the tab bar: tasks, log, config, status.
- [ ] No "run" or "profiles" tab in the bar.
- [ ] Status (tab 4) shows the profile switcher chip row above the daemon
      grid.
- [ ] Interactive: `go run .` opens on tasks (default `m.tab = 0`). Press
      Enter on a task → screen is replaced by the Run view (messages stream
      + cmdbar `reply` prompt). Press Esc → back to tasks list, selection
      preserved.
- [ ] `q` and `ctrl+c` still quit. `esc` only quits when `!m.attached` AND
      no other modal — keep this conservative; quitting via `q` is the
      documented path.
- [ ] No duplicate panel remains across the 4 tabs (verify by `grep -n 'pane('
      view_*.go` — the daemon-identity rows appear only in `view_status.go`,
      the task list only in `view_tasks.go`, the log table only in
      `view_log.go`, the runtime config only in `view_config.go`).

## Dependencies

P2, P4, P5 must be merged. P0, P1, P3, P6 may or may not be; they don't
affect the consolidation mechanics.

## Handoff (paste into next session)

```
P7 complete. Tab count is 4: tasks · log · config · status. Number keys
1–4. Profiles tab removed; switcher moved into Status header. Run tab
removed; attach via Enter on a task row, detach via Esc. main.go has
attached bool + esc/enter handling. view_status.go takes selProf and renders
the profile switcher above the grid. go build clean; --dump diff confirms
4 tabs, no orphans.
```
