package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

// renderStatus builds the default `mtop status` screen.
func renderStatus(width, height int, selectedTask int) string {
	colGap := 2
	leftW := (width - colGap) / 2
	// rightW absorbs the rounding remainder so the two columns plus the gap
	// always sum to exactly `width`; otherwise an odd `width` leaves a 1-cell
	// unpainted strap at the right edge.
	rightW := width - leftW - colGap

	// Body height available after chrome subtracted by caller.
	// Split each column 1.05 / 0.95 like the design.
	leftTop := height * 53 / 100
	leftBot := height - leftTop
	rightTop := height / 2
	rightBot := height - rightTop

	left := joinV(
		paneDaemon(leftW, leftTop),
		paneRuntimes(leftW, leftBot),
	)
	right := joinV(
		paneTasksInFlight(rightW, rightTop, selectedTask),
		paneLogTail(rightW, rightBot),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, bgPad(colGap), right)
}

func paneDaemon(width, height int) string {
	inner := width - 4
	half := inner / 2
	// 2-column KV grid: 12 KVs in 6 rows.
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
	lines = append(lines, hr("tickers", inner))
	lines = append(lines, ticker("poll       3s", daemon.PollsToday, daemon.LastPoll+" · 0 new", "accent", inner))
	lines = append(lines, ticker("heartbeat 15s", daemon.HeartbeatsToday, daemon.LastHB+" · ok", "ok", inner))

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

func paneRuntimes(width, height int) string {
	inner := width - 4
	var lines []string
	for i, r := range runtimes {
		marker := sAccent.Render("▎ ")
		if i > 0 {
			marker = sFaint.Render("▎ ")
		}
		head := marker + sOk.Render("●") + " " + sFg.Render(r.Name) +
			"  " + sDim.Render("v"+r.Ver) + "  " + sInfo.Render("model="+r.Model)
		conc := bar(float64(r.Busy)/float64(r.Max), 8) + " " + sDim.Render(fmt.Sprintf(" %d/%d", r.Busy, r.Max))
		// Align bar to right
		gap := inner - lipgloss.Width(head) - lipgloss.Width(conc)
		if gap < 1 {
			gap = 1
		}
		lines = append(lines, head+strings.Repeat(" ", gap)+conc)

		foot := marker + sDim.Render("$ ") + sFg1.Render(r.Bin)
		stat := sDim.Render(fmt.Sprintf("%d tasks · %d err", r.TasksToday, r.ErrsToday))
		gap = inner - lipgloss.Width(foot) - lipgloss.Width(stat)
		if gap < 1 {
			gap = 1
		}
		lines = append(lines, foot+strings.Repeat(" ", gap)+stat)
		if i < len(runtimes)-1 {
			lines = append(lines, "")
		}
	}
	body := strings.Join(lines, "\n")
	return pane("runtimes", "2 of 2 registered · auto-detected on $PATH", body, width, height, false)
}

func paneTasksInFlight(width, height int, selected int) string {
	inner := width - 4
	// Columns: status(2) id(7) title(rest) runtime(8) seq(8) elapsed(8) cost(6)
	statusW, idW, agentW, seqW, elapsedW, costW := 2, 7, 8, 8, 8, 6
	titleW := inner - statusW - idW - agentW - seqW - elapsedW - costW - 6 // gaps
	if titleW < 12 {
		titleW = 12
	}

	head := sDim.Render(padR("", statusW) + " " +
		padR("ID", idW) + " " +
		padR("TITLE", titleW) + " " +
		padR("RUNTIME", agentW) + " " +
		padR("SEQ", seqW) + " " +
		padR("ELAPSED", elapsedW) + " " +
		padR("COST", costW))

	var lines []string
	lines = append(lines, head)
	lines = append(lines, sFaint.Render(strings.Repeat("─", inner)))
	limit := len(tasks)
	if limit > 6 {
		limit = 6
	}
	for i, t := range tasks[:limit] {
		lines = append(lines, taskRow(t, i == selected, statusW, idW, titleW, agentW, seqW, elapsedW, costW))
	}
	body := strings.Join(lines, "\n")
	return pane("tasks · in flight", "3/20 busy · 1 queued · sort: started ↓", body, width, height, false)
}

func taskRow(t Task, selected bool, statusW, idW, titleW, agentW, seqW, elapsedW, costW int) string {
	if selected {
		// Compose every cell — including the spaces between fields and the
		// pane padding that follows — under a single bgSel paint. fillBg
		// re-emits bgSel after each inner SGR reset so unstyled spaces
		// between segments don't fall back to canvas bg and produce a strap.
		bgOn := func(fg lipgloss.Color) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(fg).Background(bgSel)
		}
		seg := bgOn(accent).Render("▎") +
			lipgloss.NewStyle().
				Foreground(statusStyle(t.Status).GetForeground()).
				Background(bgSel).Render(padR(statusGlyph(t.Status), statusW-1)) +
			bgOn(fg).Render(" "+padR(t.ID, idW)) +
			bgOn(fg).Bold(true).Render(" "+padR(truncate(t.Title, titleW), titleW)) +
			bgOn(agentTone(t.Runtime)).Render(" "+padR(t.Runtime, agentW)) +
			bgOn(dim).Render(" "+padR(fmt.Sprintf("seq %d", t.Seq), seqW)) +
			bgOn(dim).Render(" "+padR(t.Started, elapsedW)) +
			bgOn(dim).Render(" "+padR(t.Cost, costW))
		return fillBg(seg, bgSel)
	}
	status := statusStyle(t.Status).Render(padR(statusGlyph(t.Status), statusW))
	id := sFg1.Render(padR(t.ID, idW))
	title := sFg.Render(padR(truncate(t.Title, titleW), titleW))
	runtime := lipgloss.NewStyle().Foreground(agentTone(t.Runtime)).
		Render(padR(t.Runtime, agentW))
	seq := sDim.Render(padR(fmt.Sprintf("seq %d", t.Seq), seqW))
	elapsed := sDim.Render(padR(t.Started, elapsedW))
	cost := sDim.Render(padR(t.Cost, costW))
	return status + " " + id + " " + title + " " + runtime + " " + seq + " " + elapsed + " " + cost
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

func paneLogTail(width, height int) string {
	inner := width - 4
	limit := len(logLines)
	if limit > 8 {
		limit = 8
	}
	var lines []string
	for i, l := range logLines[:limit] {
		lines = append(lines, logRowText(l, i == 0, inner))
	}
	body := strings.Join(lines, "\n")
	return pane("daemon log · tail", "follow ●LIVE · q to detach", body, width, height, true)
}

func logRowText(l LogLine, hl bool, width int) string {
	tW, lvlW, srcW := 8, 6, 16
	msgW := width - tW - lvlW - srcW - 3
	if msgW < 8 {
		msgW = 8
	}

	tStyled := sDim.Render(padR(l.T, tW))
	lvlStyled := levelTag(l.Lvl, lvlW)
	srcStyled := sInfo.Render(padR(l.Src, srcW))
	msgStyled := sFg.Render(truncate(l.Msg, msgW))
	if l.Src == "poll" || l.Src == "heartbeat" {
		srcStyled = sFaint.Render(padR(l.Src, srcW))
		msgStyled = sDim.Render(truncate(l.Msg, msgW))
	}
	if l.Lvl == "warn" {
		msgStyled = sWarn.Render(truncate(l.Msg, msgW))
	}
	if l.Lvl == "error" {
		msgStyled = sErr.Render(truncate(l.Msg, msgW))
	}
	row := tStyled + " " + lvlStyled + " " + srcStyled + " " + msgStyled
	if hl {
		row = sAccent.Render("▎") + row
	} else {
		row = " " + row
	}
	return row
}

func levelTag(lvl string, w int) string {
	pad := padR(strings.ToUpper(lvl), w)
	switch lvl {
	case "trace":
		return sFaint.Render(pad)
	case "info":
		return lipgloss.NewStyle().Foreground(info).Background(bg2).Render(pad)
	case "warn":
		return lipgloss.NewStyle().Foreground(warn).Background(bg2).Render(pad)
	case "error":
		return lipgloss.NewStyle().Foreground(errCol).Background(bg2).Render(pad)
	}
	return sDim.Render(pad)
}
