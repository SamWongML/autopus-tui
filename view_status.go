package main

import (
	"fmt"
	"strings"

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
	grid := make([]string, 0, len(rows))
	for _, row := range rows {
		l := kv(row[0].k, row[0].v, row[0].tone, half)
		r := kv(row[1].k, row[1].v, row[1].tone, half)
		grid = append(grid, l+"  "+r)
	}
	body := kvPane(inner, []kvSection{
		{"", grid},
		{"tickers", []string{
			ticker("poll       3s", daemon.PollsToday, daemon.LastPoll+" · 0 new", "accent", inner),
			ticker("heartbeat 15s", daemon.HeartbeatsToday, daemon.LastHB+" · ok", "ok", inner),
		}},
	})
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
	body := taskTable(tasks, width, height, tableCompact, selected)
	return pane("tasks · in flight", "3/20 busy · 1 queued · sort: started ↓", body, width, height, false)
}

func paneLogTail(width, height int) string {
	body := logTable(logLines, width, height, logTail)
	return pane("daemon log · tail", "follow ●LIVE · q to detach", body, width, height, true)
}
