# P0 — Visual primitives

Items: **k** (dirty-state glyph), **l** (copy fix), **m** (PR-dot helper).

## Goal

Add the visual building blocks the later phases consume: a dirty/unsaved glyph
that does not collide with the healthy `●`, a PR-status dot helper, and a copy
fix in the config view. No structural changes; tiny diff.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `theme.go` (all of it — 97 lines)
3. `view_config.go` lines 30–50 (the `renderConfigHead` function)
4. `view_run.go` lines 137–155 (look at `msgTypeTag` for a style chip pattern)

That's all. Do not read other files.

## Files to modify

| File             | Why                                                                                      |
| ---------------- | ---------------------------------------------------------------------------------------- |
| `theme.go`       | add `dirtyGlyph()`, `prDot(state string)` helpers + a `pr*` color set if needed          |
| `view_config.go` | replace `●` with the dirty glyph on `Dirty: true` rows; fix the duplicate-word subtitle  |

No new files.

## Tasks

1. **`theme.go`: add a `dirtyGlyph` constant or helper.**
   Convention says `▲` rendered with `sWarn`. Add either:
   ```go
   var dirtyMark = sWarn.Render("▲")
   ```
   at package level, near the existing styles block.

2. **`theme.go`: add `prDot(state string) string`.**
   States: `"waiting"`, `"passed"`, `"merged"`, `"draft"`. Returns one styled
   cell. Mapping:
   - `waiting` → `sWarn.Render("●")`
   - `passed`  → `sOk.Render("●")`
   - `merged`  → `sViolet.Render("●")`
   - `draft`   → `sFaint.Render("●")`
   - default   → `sFaint.Render("·")`
   Place it next to `statusGlyph`.

3. **`view_config.go`: fix the subtitle.**
   In `renderConfigHead` (around line 30–50), the daemon pane title currently
   reads `"DAEMON multica daemon config"` (search for that exact string — it is
   passed to `pane(...)` as the title for the daemon column). Change to
   `"DAEMON · multica config"`.

4. **`view_config.go`: replace the dirty dot.**
   Find every place a `CfgRow` with `Dirty == true` renders a `●`. The function
   that draws rows is `cfgRow` (around line 60). Where the current code shows
   `sWarn.Render("●")` or `sAccent.Render("●")` for dirty state, swap for the
   package-level `dirtyMark` from task 1. Hint lines like `"● unsaved · ..."`
   should also use `dirtyMark` for the leading marker, not `●`.

   Healthy connection / running state dots stay `sOk.Render("●")`. Do not touch
   them.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump 2>&1 | sed -E 's/\x1B\[[0-9;]*[mK]//g'` shows
      `▲` in front of every dirty row on tab 5 (the `max-concurrent-tasks` row
      and `auto-update` row).
- [ ] The same dump shows `● online`, `● running`, `● connected` unchanged —
      `●` still means healthy.
- [ ] On tab 5 the daemon pane title is `DAEMON · multica config` (no
      duplicated "daemon" word).
- [ ] `prDot("passed")` returns a green-tone single-cell string when called
      manually (a quick `fmt.Println(lipgloss.Width(prDot("passed")))` should
      print `1`).

## Dependencies

None. `[P]` — can run in parallel with P1, P4, P5, P6.

## Handoff (paste into next session)

```
P0 complete. Added in theme.go:
  - dirtyMark = sWarn.Render("▲")    (package-level)
  - prDot(state string) string       (waiting=warn, passed=ok, merged=violet,
                                      draft=faint, default=faint·)
Fixed view_config.go subtitle: "DAEMON · multica config".
All dirty markers in view_config.go use dirtyMark; healthy ● unchanged.
go build clean. --dump diff shows only the intended glyph swaps.
```
