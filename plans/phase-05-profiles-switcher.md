# P5 — Profiles → switcher

Item: **b** (collapse the 4-card Profiles grid into a horizontal switcher
segmented control + a single active-profile detail card).

## Goal

The current Tab 6 renders one full card per profile (default, staging,
selfhost, "new profile" educational card). Each card duplicates 7 of the 12
fields that already live in Tab 1's DAEMON card — the *active* profile is
literally already on screen on Status. Replace with:

- A horizontal segmented switcher at top: one chip per profile, active one is
  highlighted, stopped ones are dim.
- Below it: ONE detail panel for the *selected* profile (which might not be
  the active one — so this view still earns its keep when comparing).
- A small docs / "new profile" hint at the bottom of the screen, not a
  permanent card.

## Context bootstrap (read in order)

1. `plans/conventions.md`
2. `view_profiles.go` (all 156 lines)
3. `data.go` lines 101–114 (the `Profile` struct and seed `profiles` slice)
4. `view_status.go` lines 39–66 (look at `paneDaemon` to match its row order
   for consistency)

That's all.

## Files to modify

| File              | Why                                                        |
| ----------------- | ---------------------------------------------------------- |
| `view_profiles.go`| rewrite — switcher chip row + single detail pane           |
| `main.go`         | `selProf` now selects which profile to inspect, not which card; left/right keys move it |

No new files.

## Tasks

1. **`view_profiles.go`: write `renderProfileSwitcher(width, selected int) string`.**
   Produce one horizontal row of chips:
   ```
   [● default · 3/20]  [● staging · 0/8]  [○ selfhost · stopped]    [+ new profile]
   ```
   - Each chip uses `lipgloss.NewStyle().Background(bg2).Padding(0,1)`.
   - The selected chip uses `Background(bgSel).Foreground(accent).Bold(true)`.
   - State dot is `sOk("●")` for running, `sFaint("○")` for stopped.
   - Suffix is `tasks` field for running profiles, `"stopped"` for stopped.
   - "new profile" is a `+` chip in `sDim` at the right side.
   - Centered? No — left-aligned with a spacer to the right edge.

2. **`view_profiles.go`: write `renderProfileDetail(p Profile, w, h int) string`.**
   Use `pane("PROFILE · "+p.Name, stateHint(p), body, w, h, p.State=="running")`.
   `stateHint(p)` returns `"running · " + p.Uptime` or `"stopped"`.
   Body rows (same order as `paneDaemon` for consistency):
   - state (with statusGlyph)
   - pid (`"—"` if zero)
   - uptime (`"—"` if empty)
   - server
   - health port
   - workspaces
   - runtimes (joined by `, `)
   - tasks
   - then `hr("ACTIONS", w)` and:
     - if running: `⏎ attach · S stop · r restart · l log`
     - if stopped: `s start · e edit config · x delete`

3. **`view_profiles.go`: rewrite `renderProfiles`.**
   Compose:
   ```
   top row:  segmented switcher                                    (1 line)
   spacer:   bgPadV(width, 1)
   middle:   renderProfileDetail(profiles[selected], width, h-3)
   bottom:   one-line dim hint "n new · ←/→ switch · ⏎ attach"     (1 line)
   ```
   `selected` is `m.selProf` passed in from `main.go`.

4. **`view_profiles.go`: delete obsolete functions.**
   Remove `renderProfilesHead`, `renderProfilesGrid`, `renderProfileCard`,
   `renderNewProfileCard`. After this change, the file should be ~80 lines.

5. **`main.go`: ←/→ moves `selProf`.**
   In `Update`, on tab 5, intercept `"left"` and `"right"` (or `"h"`/`"l"`)
   to decrement/increment `m.selProf` clamped to `[0, len(profiles)-1]`. The
   existing j/k keys can also move selection (already wired via
   `moveSelection`); add the horizontal binding so the segmented control feels
   right.

6. **(Optional) Top status strip integration.**
   If P1 is merged, the breadcrumb on tab 5 already shows `profiles ·
   active=<name>`. Adjust it to `profiles · viewing=<name> · active=<name>` so
   the user can tell selection from active.

## Acceptance

- [ ] `go build ./...` passes.
- [ ] `go run . --dump` tab 6 shows:
      - Line 1: switcher chips for default, staging, selfhost, + new.
      - Below: ONE pane titled `PROFILE · default` (or whichever is selected)
        with state/pid/uptime/.../tasks + ACTIONS divider.
      - Bottom: one dim hint line.
- [ ] No grid of 4 cards anywhere.
- [ ] Pressing right-arrow in interactive mode advances the switcher; the
      detail pane re-renders for the new profile.
- [ ] Pane border accent is **on** when viewing the active profile,
      **off-color** (no accent) when viewing a stopped profile.

## Dependencies

None. `[P]` with P0, P1, P4, P6.

## Handoff (paste into next session)

```
P5 complete. view_profiles.go rewritten: renderProfileSwitcher (chip row),
renderProfileDetail (one pane). Deleted renderProfilesHead/Grid/Card,
renderNewProfileCard. main.go: ←/→ moves selProf on tab 5. Detail pane mirrors
paneDaemon's row order so the duplication is now one-direction (Status owns the
fields; Profiles re-presents them for the *selected* profile — no permanent
card grid).
```
