# autopus-tui — Polish Plan (Items a–m)

Master index for the layout/duplication-removal work identified in the audit. Each
phase below is a separate file in this directory, sized so a fresh Claude Code
session can load **one phase + `conventions.md` + the named source files** and
implement it without seeing this master.

If you're the orchestrator, start here. If you're an implementer, open the phase
file directly and ignore this one.

---

## Reading order for an implementer (always)

1. `plans/conventions.md` — cross-phase invariants (theme tokens, chrome
   helpers, glyph table). Tiny; ~1k tokens.
2. `plans/phase-NN-<slug>.md` — the one phase you're assigned.
3. The source files that phase names under "Context bootstrap" — nothing else.

Do not load other phases. Do not load this PLAN.md. Each phase is self-contained.

---

## Phase status

| ID  | Phase                                | Items   | Deps        | Size  | Status  |
| --- | ------------------------------------ | ------- | ----------- | ----- | ------- |
| P0  | Visual primitives                    | k, l, m\* | —         | small | done    |
| P1  | Context-sensitive chrome             | d, i    | —           | med   | done    |
| P2  | Status tab: pulse + utilization      | a, e    | P0          | med   | done    |
| P3  | Tasks tab: state groups + PR dots    | f, m    | P0          | med   | done    |
| P4  | Unified detail pane + peek overlay   | c, h    | —           | large | pending |
| P5  | Profiles → switcher                  | b       | —           | med   | pending |
| P6  | Log Summary toggle                   | j       | —           | small | pending |
| P7  | Tab consolidation (6 → 4)            | g       | P2, P4, P5  | large | pending |

`m*` in P0 = add the PR-dot helper. P3 applies it on task rows.

## Dependency graph

```
P0 ─┬─→ P2 ─┐
    └─→ P3  │
            ├─→ P7
P4 ─────────┤
P5 ─────────┘
P1   (parallel-safe with everything; pick first or last)
P6   (parallel-safe; small)
```

Parallel-safe pairs (different Claude Code sessions at the same time):
- P0 + P1 + P4 + P5 + P6 — all touch disjoint files.
- After P0 lands: P2 + P3 in parallel (different views).
- P7 must run last (it edits `main.go` after the views are reshaped).

## What each item maps to (cross-reference for the audit)

```
a → P2  demote Status mini-tasks/mini-log → pulse + last-events ribbon
b → P5  Profiles → horizontal switcher segmented control
c → P4  unify Run + Tasks detail sidebars into one component
d → P1  context-sensitive top status strip (breadcrumb per tab)
e → P2  runtimes → utilization strip on Status; remove from Profiles
f → P3  group tasks by state on Tab 3
g → P7  reduce tab count 6 → 4
h → P4  Tab 2 sidebar becomes Space-toggleable peek overlay
i → P1  per-panel footer hints (replace static per-view footers)
j → P6  btop-style toggle for Log Summary sidebar
k → P0  distinct glyph for dirty/unsaved vs healthy
l → P0  fix "DAEMON multica daemon config" copy
m → P0  PR-dot helper · P3 applies it to task rows
```

## Conventions (excerpt — full version in `conventions.md`)

- All colors come from `theme.go` package-level vars (`bg`, `fg1`, `accent`,
  `ok`, `warn`, `errCol`, `info`, `violet`). Do **not** introduce new colors.
- All panel layout uses `pane(title, hint, body, w, h, accentBorder)` from
  `chrome.go`. Do **not** roll your own border rendering.
- All vertical/horizontal joining uses `joinV` / `joinH` from `chrome.go`.
- All key–value rows use `kv(k, v, tone, width)` from `chrome.go`.
- All rendered panels must paint the canvas bg (`bg`) on otherwise-empty lines —
  this is what `fillBg` and `bgPad`/`bgPadV` do. Don't return shorter strings.

## Acceptance for the whole plan

When P0–P7 are merged:
- `go build ./...` passes.
- `go run . --dump` produces output with no duplicated panels (the duplications
  in the audit table — daemon log tail, tasks roster, daemon identity, runtime
  registry — appear at most once each).
- Tab count is 4 (`tasks · run · log · config`).
- Top status strip is breadcrumb-style on tabs 2–4, single global dot.
- No regression: every action key listed in any phase's "Acceptance" still
  works after the merge.
