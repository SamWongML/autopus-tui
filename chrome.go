package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

// ─── Generic helpers ────────────────────────────────────────────────────────

func padR(s string, w int) string { return lipgloss.NewStyle().Width(w).Render(s) }

// truncate clips s to a visible width of w, appending "…" if clipped.
// ANSI-aware: escape sequences are copied verbatim and don't count toward
// width. A trailing "\x1b[0m" closes any style left open mid-render so the
// next cell on the line doesn't inherit a stray fg/bg.
//
// Note: a swap to `ansi.Truncate` was attempted but emits redundant
// open/close SGR spans around the truncation point. Visible output is
// identical, but byte-equal golden snapshots break; deferred to Phase B.
func truncate(s string, w int) string {
	if lipgloss.Width(s) <= w {
		return s
	}
	if w < 2 {
		return ""
	}
	var b strings.Builder
	visible := 0
	target := w - 1 // room for "…"
	for i := 0; i < len(s); {
		if s[i] == '\x1b' {
			end := i + 1
			if end < len(s) && s[end] == '[' {
				end++
				for end < len(s) && (s[end] < 0x40 || s[end] > 0x7e) {
					end++
				}
				if end < len(s) {
					end++
				}
			}
			b.WriteString(s[i:end])
			i = end
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		rW := lipgloss.Width(string(r))
		if visible+rW > target {
			break
		}
		b.WriteString(s[i : i+size])
		visible += rW
		i += size
	}
	b.WriteString("…\x1b[0m")
	return b.String()
}

// wrap word-wraps s to lines of at most `width` visible cells. Splits on ASCII
// whitespace via strings.Fields, so runs of whitespace collapse to single
// spaces — matches the design's tight message bodies.
//
// Note: a swap to `ansi.Wrap` was attempted but its handling of multi-space
// inputs produces byte-different output; deferred to Phase B.
func wrap(s string, width int) []string {
	if width < 1 {
		return []string{s}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	var cur string
	for _, w := range words {
		if cur == "" {
			cur = w
			continue
		}
		if lipgloss.Width(cur+" "+w) > width {
			lines = append(lines, cur)
			cur = w
		} else {
			cur += " " + w
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

// countBar renders "label  [██░░] N" — a fixed-width label, an ANSI block-
// progress bar, and a right-aligned count colored by `tone`.
func countBar(label string, n, total int, tone string, width int) string {
	labelW := 12
	numW := 6
	barW := width - labelW - numW - 2
	if barW < 4 {
		barW = 4
	}
	var nStyle lipgloss.Style
	switch tone {
	case "dim":
		nStyle = sDim
	case "info":
		nStyle = sInfo
	case "ok":
		nStyle = sOk
	case "amber":
		nStyle = sAccent
	case "err":
		nStyle = sErr
	default:
		nStyle = sFg
	}
	return sFg1.Render(padR(label, labelW)) + " " +
		bar(float64(n)/float64(total), barW) + " " +
		nStyle.Render(fmt.Sprintf("%*d", numW, n))
}

// bgSetSeq returns the ANSI prefix that lipgloss emits for Background(color),
// e.g. "\x1b[48;2;25;20;15m" in truecolor mode, "\x1b[40m" in 16-color mode,
// or "" when the active termenv profile renders without color.
func bgSetSeq(color lipgloss.Color) string {
	sample := lipgloss.NewStyle().Background(color).Render("X")
	idx := strings.Index(sample, "X")
	if idx <= 0 {
		return ""
	}
	return sample[:idx]
}

// fillBg ensures every cell of s renders with `color` as the background.
//
// Lipgloss emits ANSI SGR resets ("\x1b[0m", "\x1b[m", or "\x1b[49m") between
// styled spans, and the bg state established by an explicit Background(...) on
// one row persists past a "\n" into the next row unless something clears it.
// Either case produces visible "straps" of the wrong color. fillBg walks every
// CSI SGR sequence in s, re-emits the bg setter after each one that clears
// the bg, re-emits it after every newline so each line starts fresh, and
// prepends the setter at the very start so cells before any escape are
// painted too.
//
// SGR semantics handled:
//
//	[m        — reset all (default)
//	[0m       — reset all
//	[…;0;…m   — reset all anywhere in the parameter list
//	[49m      — reset bg only (default bg)
//	[…;49;…m  — reset bg as part of a compound
//
// Other SGR codes (foreground changes, attribute toggles, explicit bg sets)
// don't clear the bg and are passed through unchanged. A different inner
// Background(...) set by the caller does override fillBg's paint for the cells
// it covers — by design, since fillBg is the fallback paint, not the only one.
func fillBg(s string, color lipgloss.Color) string {
	bgSet := bgSetSeq(color)
	if bgSet == "" {
		return s
	}

	var b strings.Builder
	b.Grow(len(s) + 64)
	b.WriteString(bgSet)

	i := 0
	for i < len(s) {
		// Detect a CSI sequence: ESC '[' params final-byte(0x40..0x7e).
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			end := i + 2
			for end < len(s) {
				c := s[end]
				if c >= 0x40 && c <= 0x7e {
					end++
					break
				}
				end++
			}
			seq := s[i:end]
			b.WriteString(seq)
			// If final byte is 'm' (SGR) and the params clear the bg, re-emit.
			if len(seq) >= 3 && seq[len(seq)-1] == 'm' && sgrClearsBg(seq[2:len(seq)-1]) {
				b.WriteString(bgSet)
			}
			i = end
			continue
		}
		// Line boundary: re-establish bg so a row whose styling left a different
		// bg (e.g. a selected row painted with bgSel) doesn't leak into the
		// first cells of the next row, which often start with a foreground-only
		// SGR and would otherwise inherit the prior bg.
		if s[i] == '\n' {
			b.WriteByte('\n')
			b.WriteString(bgSet)
			i++
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	// Close the styling state so callers that compose fillBg output into a
	// larger string (e.g. lipgloss bodyStyle wrapping a body that contains a
	// selected-row segment painted with bgSel) don't inherit the bg we set
	// here. Without this, lipgloss carries the open bg state through into the
	// padding it adds for subsequent lines, painting their leading cells with
	// the wrong color — the exact "selection bleed" seen on Status tab.
	b.WriteString("\x1b[0m")
	return b.String()
}

// sgrClearsBg reports whether the given SGR parameter list clears the
// background. Empty params (i.e. "\x1b[m") and "0" both mean "reset all";
// "49" specifically resets the bg to default. Any of those appearing in a
// semicolon-separated compound is enough.
func sgrClearsBg(params string) bool {
	if params == "" {
		return true
	}
	start := 0
	for i := 0; i <= len(params); i++ {
		if i == len(params) || params[i] == ';' {
			p := params[start:i]
			if p == "" || p == "0" || p == "49" {
				return true
			}
			start = i + 1
		}
	}
	return false
}

// joinH joins horizontally with a single-space gap.
func joinH(parts ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func joinV(parts ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// bgPad renders `n` cells of canvas-bg-painted whitespace. Used to fill the
// gap between two regions in a horizontal join so the column between them
// matches the canvas instead of falling back to the terminal default bg.
func bgPad(n int) string {
	if n <= 0 {
		return ""
	}
	return lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", n))
}

// bgPadV renders a column of `n` blank lines, each cell painted with canvas
// bg. Used to vertically separate panes inside a grid.
func bgPadV(width, n int) string {
	if n <= 0 || width <= 0 {
		return ""
	}
	row := lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", width))
	rows := make([]string, n)
	for i := range rows {
		rows[i] = row
	}
	return strings.Join(rows, "\n")
}

// bar renders [████░░░░] with `width` cells.
func bar(value float64, width int) string {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	filled := int(value*float64(width) + 0.5)
	if filled > width {
		filled = width
	}
	return sAccent.Render(strings.Repeat("█", filled)) +
		sFaint.Render(strings.Repeat("░", width-filled))
}

// kv renders "key                  value" with key dim and value bright.
// Total width clamps the value to the available space.
func kv(k, v, tone string, width int) string {
	keyW := len(k)
	gap := width - keyW - lipgloss.Width(v)
	if gap < 1 {
		gap = 1
	}
	var vStyle lipgloss.Style
	switch tone {
	case "accent":
		vStyle = sAccent
	case "ok":
		vStyle = sOk
	case "warn":
		vStyle = sWarn
	case "err":
		vStyle = sErr
	case "dim":
		vStyle = sDim
	default:
		vStyle = sFg
	}
	return sDim.Render(k) + strings.Repeat(" ", gap) + vStyle.Render(v)
}

// hr — dashed divider with optional centred label.
func hr(label string, width int) string {
	line := sFaint.Render(strings.Repeat("─", width))
	if label == "" {
		return line
	}
	lab := sDim.Render(" " + strings.ToUpper(label) + " ")
	half := (width - lipgloss.Width(lab)) / 2
	if half < 1 {
		half = 1
	}
	left := sFaint.Render(strings.Repeat("─", half))
	right := sFaint.Render(strings.Repeat("─", width-half-lipgloss.Width(lab)))
	return left + lab + right
}

// ─── Pane ───────────────────────────────────────────────────────────────────

// pane renders a rounded box with `title` (optionally `hint`) embedded in the
// top border, like:
//
//	╭─ TITLE · hint ─────────────╮
//	│ body                       │
//	╰────────────────────────────╯
//
// Every cell of the returned region — body, borders, and the padding that
// lipgloss adds when content is shorter than `inner` — is explicitly painted
// with the canvas `bg`. A final fillBg pass re-emits the bg setter after each
// SGR reset so internal styled spans don't leak the terminal default through
// when they reset.
func pane(title, hint, body string, width, height int, accentBorder bool) string {
	border := lipgloss.RoundedBorder()
	borderColor := faint
	if accentBorder {
		borderColor = accent
	}

	inner := width - 2 // minus two for border columns
	innerH := height - 2
	if innerH < 1 {
		innerH = 1
	}
	if inner < 4 {
		inner = 4
	}

	// Render body inside inner×innerH with bg-painted padding. lipgloss pads
	// each line to `inner` with bg-painted spaces and clips to innerH rows
	// because we set Background here.
	bodyStyle := lipgloss.NewStyle().
		Width(inner).
		Height(innerH).
		Background(bg)
	rendered := bodyStyle.Render(body)

	borderStyle := lipgloss.NewStyle().Foreground(borderColor).Background(bg)
	top := buildTopBorder(border, title, hint, width, accentBorder)
	bottom := borderStyle.Render(border.BottomLeft + strings.Repeat(border.Bottom, width-2) + border.BottomRight)
	left := borderStyle.Render(border.Left)
	right := borderStyle.Render(border.Right)

	var rows []string
	rows = append(rows, top)
	for _, line := range strings.Split(rendered, "\n") {
		// bodyStyle pads to `inner` already; defensive pad in case it doesn't.
		if w := lipgloss.Width(line); w < inner {
			line += lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", inner-w))
		}
		rows = append(rows, left+line+right)
	}
	rows = append(rows, bottom)
	return fillBg(strings.Join(rows, "\n"), bg)
}

func buildTopBorder(b lipgloss.Border, title, hint string, width int, accentBorder bool) string {
	borderColor := faint
	if accentBorder {
		borderColor = accent
	}
	titleStyled := lipgloss.NewStyle().Foreground(accent).Background(bg).Bold(true).Render(strings.ToUpper(title))
	label := titleStyled
	if hint != "" {
		label += lipgloss.NewStyle().Background(bg).Render(" ") +
			lipgloss.NewStyle().Foreground(dim).Background(bg).Render(hint)
	}
	labelWidth := lipgloss.Width(label)

	// `╭─ label ─...─╮`
	dashes := width - 2 // total interior dashes
	left := 2          // "╭─"
	right := dashes - left - labelWidth - 2
	if right < 1 {
		right = 1
		// Truncate label if too wide.
		over := labelWidth - (dashes - left - 2 - 1)
		if over > 0 && labelWidth > over+1 {
			label = lipgloss.NewStyle().MaxWidth(labelWidth - over).Render(label)
		}
	}

	bs := lipgloss.NewStyle().Foreground(borderColor).Background(bg)
	gap := lipgloss.NewStyle().Background(bg).Render(" ")
	return bs.Render(b.TopLeft+strings.Repeat(b.Top, left)) +
		gap + label + gap +
		bs.Render(strings.Repeat(b.Top, right)+b.TopRight)
}

// chip renders a small padded label pill. Active chips paint with accent
// fg on bg2; inactive chips use dim fg on bg1.
func chip(label string, active bool) string {
	st := lipgloss.NewStyle().Padding(0, 1)
	if active {
		st = st.Foreground(accent).Background(bg2)
	} else {
		st = st.Foreground(dim).Background(bg1)
	}
	return st.Render(label)
}

type chipSpec struct {
	label  string
	active bool
}

// chipBar renders `<label>  <chip1> <chip2> …`. The double space between the
// dim label and the first chip matches both prior local-closure callsites.
func chipBar(label string, items ...chipSpec) string {
	var b strings.Builder
	b.WriteString(sDim.Render(label))
	b.WriteString("  ")
	for i, c := range items {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(chip(c.label, c.active))
	}
	return b.String()
}

// pageHeader renders a tab's top filter/breadcrumb strip: a single row with
// canvas bg, padding (0,2), and a faint bottom border. When `right` is empty,
// `left` is rendered as-is; otherwise the two are separated by stretch space
// so `right` hugs the right edge.
func pageHeader(width int, left, right string) string {
	var row string
	if right == "" {
		row = fillBg(left, bg)
	} else {
		gap := width - lipgloss.Width(left) - lipgloss.Width(right) - 4
		if gap < 1 {
			gap = 1
		}
		row = fillBg(left+strings.Repeat(" ", gap)+right, bg)
	}
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

// ─── Top chrome ─────────────────────────────────────────────────────────────

func titleBar(title string, width int) string {
	dots := sErr.Render("●") + " " + sWarn.Render("●") + " " + sOk.Render("●")
	t := sFg1.Render(title)
	pwd := sDim.Render("workspace: ") + sFg1.Render("blackhole-os")
	left := dots + "  " + t
	avail := width - 2 // padding 0,1
	gap := avail - lipgloss.Width(left) - lipgloss.Width(pwd)
	if gap < 1 {
		// Drop pwd if there's no room.
		gap = 1
		pwd = ""
		gap = avail - lipgloss.Width(left)
		if gap < 1 {
			gap = 1
		}
	}
	bar := fillBg(left+strings.Repeat(" ", gap)+pwd, bg1)
	return lipgloss.NewStyle().
		Background(bg1).
		Foreground(fg1).
		MaxWidth(width).
		Padding(0, 1).
		Render(bar)
}

func tabBar(active int, width int) string {
	tabs := []string{"status", "run", "tasks", "log", "config", "profiles"}
	var parts []string
	for i, t := range tabs {
		num := fmt.Sprintf("%d", i+1)
		if i == active {
			numStyled := lipgloss.NewStyle().Foreground(bg).Background(accent).Bold(true).Padding(0, 1).Render(num)
			lblStyled := sAccent.Render(t)
			parts = append(parts, numStyled+" "+lblStyled)
		} else {
			numStyled := lipgloss.NewStyle().Foreground(dim).Background(bg2).Padding(0, 1).Render(num)
			lblStyled := sDim.Render(t)
			parts = append(parts, numStyled+" "+lblStyled)
		}
	}
	left := strings.Join(parts, "  ")
	meta := sDim.Render("daemon ") + sOk.Render("●") + sDim.Render(" online · profile ") +
		sFg1.Render(daemon.Profile) + sDim.Render(fmt.Sprintf(" · pid %d · %s", daemon.PID, daemon.Uptime))
	gap := width - lipgloss.Width(left) - lipgloss.Width(meta) - 2
	if gap < 1 {
		gap = 1
	}
	row := fillBg(left+strings.Repeat(" ", gap)+meta, bg)
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

func footer(items []KeyHint, width int) string {
	var parts []string
	for _, it := range items {
		k := lipgloss.NewStyle().
			Foreground(fg1).
			Background(bg2).
			Padding(0, 1).
			Render(it.K)
		parts = append(parts, k+" "+sDim.Render(it.L))
	}
	left := strings.Join(parts, "   ")
	right := lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render("?") +
		" " + sDim.Render("help   ") +
		lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render("q") +
		" " + sDim.Render("quit")
	gap := width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}
	row := fillBg(left+strings.Repeat(" ", gap)+right, bg1)
	return lipgloss.NewStyle().
		Background(bg1).
		Width(width).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

type KeyHint struct{ K, L string }

// cmdBar renders the bottom `:` / `/` / reply prompt.
func cmdBar(mode, placeholder, value string, width int) string {
	var modeBlock string
	var arrow string
	switch mode {
	case ":":
		modeBlock = lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Bold(true).Render("cmd")
		arrow = sAccent.Render(":")
	case "/":
		modeBlock = lipgloss.NewStyle().Foreground(info).Background(bg2).Padding(0, 1).Bold(true).Render("search")
		arrow = sInfo.Render("/")
	case "reply":
		modeBlock = lipgloss.NewStyle().Foreground(info).Background(bg2).Padding(0, 1).Bold(true).Render("reply")
		arrow = sInfo.Render("↳")
	default:
		modeBlock = lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render(mode)
		arrow = sAccent.Render(">")
	}
	var v string
	if value != "" {
		v = sFg.Render(value)
	} else {
		v = sDim.Render(placeholder)
	}
	caret := sAccent.Render("█")
	row := modeBlock + " " + arrow + " " + v + " " + caret
	// Truncate to width.
	if lipgloss.Width(row) > width-2 {
		row = lipgloss.NewStyle().MaxWidth(width - 2).Render(row)
	}
	row = fillBg(row, bg)
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}
