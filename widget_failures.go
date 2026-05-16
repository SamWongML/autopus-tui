package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// recentFailures renders a fixed-column list of failed tasks: time, ID,
// runtime, issue, message. The first entry gets an accent vertical bar
// (matches the convention in logRowText for "active/focused row"); the rest
// get a leading space so columns align. Output is the pane body — callers
// wrap with `pane()` for the chrome.
//
// Up to 5 rows render at full height; tighter heights clamp further.
func recentFailures(items []Failure, width, height int) string {
	inner := width - 4
	if inner < 30 {
		inner = 30
	}
	limit := len(items)
	if limit > 5 {
		limit = 5
	}
	if h := height - 2; h > 0 && limit > h {
		limit = h
	}
	if limit <= 0 {
		return sDim.Render("no failures in the last 24h")
	}

	tW, idW, agentW, issueW := 8, 7, 8, 6
	msgW := inner - tW - idW - agentW - issueW - 5
	if msgW < 10 {
		msgW = 10
	}

	lines := make([]string, 0, limit)
	for i, f := range items[:limit] {
		t := sDim.Render(padR(f.T, tW))
		id := sFg1.Render(padR(f.ID, idW))
		agent := lipgloss.NewStyle().Foreground(agentTone(f.Runtime)).
			Render(padR(f.Runtime, agentW))
		issue := sDim.Render(padR(f.Issue, issueW))
		msg := sErr.Render(truncate(f.Msg, msgW))

		row := t + " " + id + " " + agent + " " + issue + " " + msg
		if i == 0 {
			row = sErr.Render("▎") + row
		} else {
			row = " " + row
		}
		lines = append(lines, row)
	}
	return strings.Join(lines, "\n")
}
