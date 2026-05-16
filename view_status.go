package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

// renderStatus builds the default `mtop status` screen.
func renderStatus(width, height int, selectedTask int) string {
	_ = selectedTask // status no longer surfaces task selection
	colGap := 2
	leftW := (width - colGap) / 2
	// rightW absorbs the rounding remainder so the two columns plus the gap
	// always sum to exactly `width`; otherwise an odd `width` leaves a 1-cell
	// unpainted strap at the right edge.
	rightW := width - leftW - colGap

	leftTop := height * 60 / 100
	leftBot := height - leftTop
	rightTop := height * 55 / 100
	rightBot := height - rightTop

	left := joinV(
		paneDaemon(leftW, leftTop),
		paneRuntimes(leftW, leftBot),
	)
	right := joinV(
		panePulse(rightW, rightTop),
		paneLastEvents(rightW, rightBot),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, bgPad(colGap), right)
}

func paneDaemon(width, height int) string {
	inner := width - 4
	half := inner / 2
	type pair struct{ k, v, tone string }
	rows := [][2]pair{
		{{"state", "● running", "ok"}, {"pid", fmt.Sprintf("%d", daemon.PID), ""}},
		{{"uptime", daemon.Uptime, ""}, {"version", "v" + daemon.Version, ""}},
		{{"server", daemon.Server, ""}, {"connection", "● connected", "ok"}},
		{{"device", daemon.Device, ""}, {"mem", fmt.Sprintf("%d MB / %d", daemon.MemMB, daemon.MemMax), ""}},
		{{"cpu", fmt.Sprintf("%.1f%%", daemon.CPU), ""}, {"socket", daemon.Socket, "dim"}},
		{{"log", daemon.Log, "dim"}, {"workspaces", daemon.WsRoot, "dim"}},
	}
	var lines []string
	for _, row := range rows {
		l := kv(row[0].k, row[0].v, row[0].tone, half)
		r := kv(row[1].k, row[1].v, row[1].tone, half)
		lines = append(lines, l+"  "+r)
	}
	body := strings.Join(lines, "\n")
	return pane("daemon", "profile · "+daemon.Profile, body, width, height, false)
}

func ticker(label string, n int, last, tone string, width int) string {
	pulse := sAccent.Render("●")
	if tone == "ok" {
		pulse = sOk.Render("●")
	}
	num := sAccent.Render(fmt.Sprintf("%6s", commafy(n)))
	if tone == "ok" {
		num = sOk.Render(fmt.Sprintf("%6s", commafy(n)))
	}
	lbl := sFg1.Render(label)
	tail := sDim.Render(last)
	left := pulse + "  " + lbl + "  " + num
	gap := width - lipgloss.Width(left) - lipgloss.Width(tail)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + tail
}

func commafy(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}
	// Insert commas.
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// paneRuntimes renders one row per runtime: a utilization bar plus today's
// counts. Bin path and full version live in Config (tab 5).
func paneRuntimes(width, height int) string {
	inner := width - 4
	var lines []string
	for _, r := range runtimes {
		lines = append(lines, runtimeRow(r, inner))
	}
	body := strings.Join(lines, "\n")
	return pane("runtimes", fmt.Sprintf("%d registered · auto-detected on $PATH", len(runtimes)), body, width, height, false)
}

func runtimeRow(r Runtime, width int) string {
	marker := sFaint.Render("▎")
	dot := sOk.Render("●")
	name := lipgloss.NewStyle().Foreground(agentTone(r.Name)).Render(padR(r.Name, 8))
	util := bar(float64(r.Busy)/float64(r.Max), 8)
	conc := sDim.Render(fmt.Sprintf("%d/%d", r.Busy, r.Max))
	model := sInfo.Render(padR(r.Model, 12))
	stats := sDim.Render(fmt.Sprintf("%d tasks · %d err", r.TasksToday, r.ErrsToday))

	left := marker + " " + dot + " " + name + " " + util + " " + conc + "  " + model
	gap := width - lipgloss.Width(left) - lipgloss.Width(stats)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + stats
}

// panePulse renders today's throughput, error trend, and the daemon tickers.
// This is content only the Status tab can usefully show; it replaces the
// duplicated tasks-roster and log-tail panels.
func panePulse(width, height int) string {
	inner := width - 4
	tasksTotal := sumInts(tasksPerHour)
	errsTotal := sumInts(errsPerHour)

	errTone := "dim"
	if errsTotal > 0 {
		errTone = "warn"
	}
	errSpark := sDim.Render(strings.Repeat("▁", min(inner-2, len(errsPerHour))))
	if errsTotal > 0 {
		errSpark = sparklineColored(errsPerHour, inner-2, sErr)
	}

	var lines []string
	lines = append(lines, kv("tasks/h", fmt.Sprintf("%d today", tasksTotal), "accent", inner))
	lines = append(lines, "  "+sparklineColored(tasksPerHour, inner-2, sAccent))
	lines = append(lines, kv("errors/h", fmt.Sprintf("%d today", errsTotal), errTone, inner))
	lines = append(lines, "  "+errSpark)
	lines = append(lines, hr("tickers", inner))
	lines = append(lines, ticker("poll       3s", daemon.PollsToday, daemon.LastPoll+" · 0 new", "accent", inner))
	lines = append(lines, ticker("heartbeat 15s", daemon.HeartbeatsToday, daemon.LastHB+" · ok", "ok", inner))

	body := strings.Join(lines, "\n")
	return pane("pulse · today", "live", body, width, height, false)
}

// paneLastEvents shows the 5 most-recent system events: task transitions and
// runtime warnings, merged. The mock seeds these directly; once a real event
// stream lands, replace `lastEvents` with a derivation from tasks/logLines.
func paneLastEvents(width, height int) string {
	inner := width - 4
	const rows = 5
	events := lastEvents()
	if len(events) > rows {
		events = events[:rows]
	}
	var lines []string
	for _, e := range events {
		lines = append(lines, eventRow(e, inner))
	}
	// Pad to rows so the pane height matches the design even with fewer events.
	for i := len(events); i < rows; i++ {
		lines = append(lines, "")
	}
	body := strings.Join(lines, "\n")
	return pane("last events", "merged · tasks + runtimes", body, width, height, false)
}

type event struct {
	status string // statusGlyph key (working/waiting/done/failed) or "warn"
	id     string // e.g. "t-1284" or "runtime.codex"
	msg    string // human description
	when   string // "3m 12s ago"
}

func lastEvents() []event {
	return []event{
		{"working", "t-1284", "started", "3m 12s ago"},
		{"waiting", "t-1283", "waiting on user", "12s ago"},
		{"failed", "t-1280", "failed · bundler resolver", "21m ago"},
		{"warn", "runtime.codex", "Bundler::VersionConflict", "21m ago"},
		{"done", "t-1281", "PR #4130 opened", "8m ago"},
	}
}

func eventRow(e event, width int) string {
	var glyph string
	if e.status == "warn" {
		glyph = sWarn.Render("⚠")
	} else {
		glyph = statusStyle(e.status).Render(statusGlyph(e.status))
	}
	id := sFg1.Render(padR(e.id, 14))
	when := sDim.Render(e.when)
	head := glyph + " " + id + " "
	avail := width - lipgloss.Width(head) - lipgloss.Width(when) - 1
	if avail < 4 {
		avail = 4
	}
	msg := sDim.Render(truncate(e.msg, avail))
	gap := width - lipgloss.Width(head) - lipgloss.Width(msg) - lipgloss.Width(when)
	if gap < 1 {
		gap = 1
	}
	return head + msg + strings.Repeat(" ", gap) + when
}

// sparklineColored renders a single-line bar sparkline in the given style.
// Mirrors sparklineBars in view_log.go but lets the caller pick the color so
// errors render in sErr while throughput stays sAccent.
func sparklineColored(data []int, width int, style lipgloss.Style) string {
	if len(data) == 0 || width < 1 {
		return ""
	}
	blocks := []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	max := 1
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	stride := 1
	if len(data) > width {
		stride = len(data) / width
	}
	var b strings.Builder
	for i := 0; i < len(data); i += stride {
		idx := int(float64(data[i]) / float64(max) * 8)
		if idx > 8 {
			idx = 8
		}
		b.WriteString(blocks[idx])
	}
	return style.Render(b.String())
}

func sumInts(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// truncate clips s to a visible width of w, appending "…" if clipped.
// ANSI-aware: escape sequences are copied verbatim and don't count toward
// width. A trailing "\x1b[0m" closes any style left open mid-render so the
// next cell on the line doesn't inherit a stray fg/bg.
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
