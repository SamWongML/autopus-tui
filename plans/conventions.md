# Conventions (read once per phase)

Cross-phase invariants. Every phase file assumes you know these. ~1k tokens.

## Palette (`theme.go`)

| Var      | Hex       | Used for                            |
| -------- | --------- | ----------------------------------- |
| `bg`     | `#19140f` | canvas (paint blanks with this)     |
| `bg1`    | `#221c16` | footer / cmdbar surface             |
| `bg2`    | `#2c241d` | chips, key caps, agent pills        |
| `bgSel`  | `#3a2f25` | selected row                        |
| `fg`     | `#ebe3d5` | primary text                        |
| `fg1`    | `#c4baa9` | secondary text (values in kv rows)  |
| `dim`    | `#8a7f72` | metadata, labels                    |
| `faint`  | `#574d43` | hairlines, borders                  |
| `accent` | `#e5b563` | working/active state, claude        |
| `ok`     | `#a7c47c` | healthy dots, success               |
| `warn`   | `#e5b563` | needs-input, dirty                  |
| `errCol` | `#d97a6b` | failed                              |
| `info`   | `#8eb0d5` | search/reply mode                   |
| `violet` | `#b29ad5` | codex                               |
| `rose`   | `#d59b9a` | (unused, reserved)                  |

Styles (`sFg`, `sFg1`, `sDim`, `sFaint`, `sAccent`, `sOk`, `sWarn`, `sErr`,
`sInfo`, `sViolet`) wrap these as `lipgloss.NewStyle().Foreground(...)`.

## Status glyphs (`theme.go: statusGlyph`)

| State     | Glyph |
| --------- | ----- |
| working   | `▶`   |
| waiting   | `!`   |
| done      | `✓`   |
| failed    | `✗`   |
| idle      | `·`   |
| stopped   | `▪`   |
| queued    | `…`   |

The healthy dot in chrome is `●` (rendered with `sOk`). After P0, "dirty/unsaved"
uses `▲` (rendered with `sWarn`) — never `●`.

## Chrome helpers (`chrome.go`)

| Function                                                        | Use                                            |
| --------------------------------------------------------------- | ---------------------------------------------- |
| `pane(title, hint, body, w, h, accentBorder bool) string`       | bordered panel; `accentBorder` colors top rule |
| `kv(k, v, tone, width) string`                                  | one key/value row inside a pane                |
| `hr(label, width) string`                                       | dim section divider                            |
| `joinV(parts...) / joinH(parts...) string`                      | join with `bg` fill so gaps stay painted       |
| `bar(value 0..1, width) string`                                 | filled+unfilled bar segments                   |
| `fillBg(s, color) string` / `bgPad(n)` / `bgPadV(width, n)`     | paint empties so the canvas bg shows through   |
| `padR(s, w) string`                                             | right-pad with width-aware truncation          |
| `tabBar(active, width) string`                                  | top tab row                                    |
| `footer(items []KeyHint, width) string`                         | bottom key hint row                            |
| `cmdBar(mode, placeholder, value, width) string`                | `:` / `/` / `reply` prompt above the footer    |

`KeyHint{K, L string}` is the type for footer entries (`{K: "⏎", L: "open"}`).

## Layout invariants

- Every view function takes `(width, height int, ...selection)` and returns a
  string of exactly `height` rows.
- Pad short bodies with `bgPadV(width, n)` so the canvas bg fills the gap.
- Use `lipgloss.Width(s)` for terminal cell width (not `len(s)`); ANSI escapes
  are not characters.
- Side gutters (1 cell each) are added by `main.go`'s `View()` — don't add your
  own at the panel level.

## Data sources (`data.go`)

```go
var daemon Daemon        // single struct, fields: Profile, PID, Uptime, Version,
                         // Server, Connected, Socket, Log, WsRoot, Device,
                         // MemMB, MemMax, CPU, LastPoll, LastHB,
                         // PollsToday, HeartbeatsToday
var runtimes []Runtime   // Name, Bin, Ver, Model, Busy, Max, TasksToday, ErrsToday
var tasks []Task         // ID, Status, Runtime, WS, Issue, Started, Last,
                         // Title, Cost, Seq
var msgs []Msg           // Seq, T, Type, Body, Tool, Args, Lines, Ok
var logLines []LogLine   // T, Lvl, Src, Msg
var profiles []Profile   // Name, State, Server, WS, Tasks, Uptime, PID, Port, Runtimes
var cfgDaemon/cfgClaude/cfgCodex/cfgLogging []CfgRow  // K, V, Env, Hint, Tone, Dirty, Readonly
```

When a phase adds a field, append at the end of the struct in `data.go` and
update the seed values.

## Verification commands

After any change:

```bash
go build ./...                                       # must pass
go run . --dump 2>&1 | sed -E 's/\x1B\[[0-9;]*[mK]//g' > /tmp/after.txt
diff /tmp/before.txt /tmp/after.txt | head -200      # eyeball regressions
```

To capture a "before" baseline before editing:

```bash
go run . --dump 2>&1 | sed -E 's/\x1B\[[0-9;]*[mK]//g' > /tmp/before.txt
```

The interactive form runs with `go run .` (no flags).

## Code style (Go specific to this repo)

- Don't add comments unless the *why* is non-obvious.
- Don't introduce new packages; the repo is single-package `main`.
- Don't add dependencies; `lipgloss`, `bubbletea`, `termenv` are the only ones.
- Don't refactor unrelated code in the same phase.
