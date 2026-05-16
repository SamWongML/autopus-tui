# P1 — Context-sensitive chrome

Items: **d** (top status strip becomes per-tab breadcrumb), **i** (footer hints
follow the focused panel).

## Goal

Stop repeating the same global "daemon ● online · profile default · pid · uptime"
strip on every tab. On Status it's redundant; on other tabs it should be a
breadcrumb describing the current scope. Make the bottom footer hints reflect
the focused panel, not the tab as a whole.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `main.go` (all 243 lines — the `View()` function is where the strip + footer
   are assembled; `Update()` is where focus would change)
3. `chrome.go` lines 334–395 (the `tabBar` and `footer` functions)
4. `chrome.go` lines 12–60 (helpers `padR`, `fillBg`)

That's all. Don't read the view files.

## Files to modify

| File         | Why                                                                                       |
| ------------ | ----------------------------------------------------------------------------------------- |
| `chrome.go`  | refactor `tabBar` to accept a right-side string (the breadcrumb), defaulting to status    |
| `main.go`    | add a `focus int` field to `model`, compute breadcrumb per tab, choose footer hints by focus |

No new files.

## Tasks

1. **`chrome.go`: split `tabBar` into `tabBar(active, width)` (tabs only) and
   `topStrip(tabs string, right string, width int) string`.**
   Currently `tabBar` bakes the daemon meta string in. Pull that out:
   - `tabBar` returns just the colored tab labels (left half).
   - New `topStrip` takes the tab string and a `right` parameter and renders
     them at width with the same border-bottom hairline.
   - On the right edge always include a one-cell `●` connection dot
     (`sOk.Render("●")` if connected, `sErr.Render("●")` if not) followed by a
     dim profile hint (`sDim.Render("default")`) — that's all that's *global*.

2. **`main.go`: add a focus model.**
   Extend `model` with `focus int` (0 = primary list, 1 = sidebar/detail) and
   a no-op `Tab` keystroke that toggles `m.focus = 1 - m.focus`. Default focus
   per tab:
   - tab 0 (status): 0
   - tab 1 (run): 0 (message stream)
   - tab 2 (tasks): 0 (list)
   - tab 3 (log): 0
   - tab 4 (config): 0
   - tab 5 (profiles): 0

3. **`main.go`: compute the breadcrumb in `View()`.**
   Replace the existing `topStrip := tabBar(m.tab, m.width)` (the assembled
   strip) with:
   ```go
   tabs := tabBar(m.tab, m.width)
   right := breadcrumb(m.tab, m.selTask)   // new helper, in main.go
   topStrip := topStrip(tabs, right, m.width)
   ```
   `breadcrumb(tab, sel int) string` returns:
   - tab 0: `""` (Status owns the daemon info; no breadcrumb needed)
   - tab 1: `"tasks › " + tasks[selTask].ID + " · ws " + tasks[selTask].WS + " · seq " + strconv.Itoa(tasks[selTask].Seq)`
   - tab 2: `"tasks · " + filterSummary()`  where `filterSummary` returns
     `"in-flight · sort started↓"` for now (literal until a real filter exists)
   - tab 3: `"log · level=info · follow ●"`
   - tab 4: `"config · " + cfgDaemon[selCfg].K + (" ▲" if dirty else "")`
   - tab 5: `"profiles · active=" + daemon.Profile`
   All styled with `sDim` except status-bearing glyphs.

4. **`main.go`: footer hints follow focus.**
   The switch in `View()` already builds `keys []KeyHint` per tab. Extend it to
   pick a different `keys` slice depending on `m.focus`:
   - tab 1 (run), focus 0: keys for the message stream (j/k scroll, f follow,
     y yank, / filter)
   - tab 1 (run), focus 1: keys for the sidebar (k kill, r reply, o $EDITOR,
     c copy)
   - tab 2 (tasks), focus 0: list keys (↑↓ select, ⏎ open, /, sort)
   - tab 2 (tasks), focus 1: preview keys (r reply, k kill, c copy, o open)
   - other tabs: same keys regardless of focus.
   Always append `{"Tab", "switch focus"}` at the end of the row when both
   panels exist (tabs 1, 2).

5. **`main.go`: dim chrome on unfocused side.**
   No code change needed yet — the actual visual treatment of focus lands in
   P4. P1 only wires the model and the footer hints.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump 2>&1 | sed -E 's/\x1B\[[0-9;]*[mK]//g'` shows on each
      tab a top strip with:
      - tab 1 (status): `[tabs]                                        ● default`
        (no pid/uptime in the strip; that info lives in the body card on Status)
      - tab 2 (run): `[tabs]              tasks › t-1284 · ws blackhole-os · seq 84      ● default`
      - tab 3 (tasks): `[tabs]              tasks · in-flight · sort started↓             ● default`
      - tab 4 (log): `[tabs]               log · level=info · follow ●                    ● default`
      - tab 5 (config): `[tabs]            config · max-concurrent-tasks ▲                ● default`
      - tab 6 (profiles): `[tabs]          profiles · active=default                      ● default`
- [ ] Pressing `Tab` in the interactive form (`go run .`) on tabs 2 or 3 swaps
      footer hint set; the right side reads "Tab switch focus".
- [ ] No regression: pressing 1–6 still switches tabs; j/k/g/G still moves
      selection; q quits.

## Dependencies

None. `[P]` — can run in parallel with P0, P4, P5, P6.

## Handoff (paste into next session)

```
P1 complete. chrome.go: split tabBar → tabBar (labels) + topStrip(tabs, right, w).
main.go: added model.focus (Tab toggles); breadcrumb(tab, sel) computed per
tab; footer hints split by focus on run/tasks tabs. Global right edge is just
"● <profile>". --dump diff shows the new strips per tab, no other panel changes.
```
