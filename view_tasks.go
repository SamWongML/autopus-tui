package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func renderTasks(width, height int, tbl table.Model) string {
	filterRow := renderFilterBar(width)
	bodyH := height - 3
	if bodyH < 6 {
		bodyH = 6
	}
	gap := 2
	rightW := 48
	if rightW > width*40/100 {
		rightW = width * 40 / 100
	}
	leftW := width - rightW - gap

	left := paneTasksTable(leftW, bodyH, tbl)
	right := paneTaskPreview(rightW, bodyH)
	body := joinH(left, bgPad(gap), right)

	return joinV(filterRow, body)
}

func renderFilterBar(width int) string {
	left := chipBar("filter",
		chipSpec{"status: in-flight", true},
		chipSpec{"runtime: any", false},
		chipSpec{"workspace: blackhole-os", false},
		chipSpec{"since: today", false},
	)
	right := sDim.Render("sort: ") + sFg.Render("started ↓")
	return pageHeader(width, left, right)
}

func paneTasksTable(width, height int, tbl table.Model) string {
	return pane(fmt.Sprintf("%d tasks", len(tasks)), "", tbl.View(), width, height, false)
}

func paneTaskPreview(width, height int) string {
	inner := width - 4
	rows := [][2]string{
		{"status", "● working"},
		{"runtime", "claude · sonnet-4.6"},
		{"pid", "18421"},
		{"workspace", "blackhole-os"},
		{"cwd", "~/multica_workspaces/blackhole-os/t-1284"},
		{"issue", "#4127"},
		{"branch", "feat/idempotency-keys"},
		{"last", "editing webhook_event.rb"},
	}
	kvs := make([]string, 0, len(rows))
	for _, r := range rows {
		tone := ""
		if r[0] == "status" {
			tone = "accent"
		}
		if r[0] == "cwd" {
			tone = "dim"
		}
		kvs = append(kvs, kv(r[0], r[1], tone, inner))
	}
	msgLines := make([]string, 0, 5)
	for _, m := range msgs[len(msgs)-5:] {
		msgLines = append(msgLines, msgRowCompact(m, inner))
	}
	body := kvPane(inner, []kvSection{
		{"", kvs},
		{"last 5 messages", msgLines},
		{"actions", []string{
			actionRow("⏎", "open run viewer"),
			actionRow("r", "reply (waiting tasks only)"),
			actionRow("k", "kill child process (SIGTERM)"),
			actionRow("K", "force kill (SIGKILL)"),
			actionRow("o", "open workspace in $EDITOR"),
			actionRow("c", "copy task id"),
		}},
	})
	return pane("preview · t-1284", "", body, width, height, false)
}

func msgRowCompact(m Msg, width int) string {
	t := sDim.Render(m.T)
	tag := msgTypeTag(m.Type, 12)
	var body string
	switch m.Type {
	case "tool_call":
		body = sInfo.Render(m.Tool) + " " + sFg1.Render(m.Args)
	case "tool_result":
		body = sOk.Render(m.Body)
	case "thinking":
		body = lipgloss.NewStyle().Foreground(fg1).Italic(true).Render(m.Body)
	case "text":
		body = sFg.Render(m.Body)
	case "error":
		body = sErr.Render(m.Body)
	}
	prefix := t + " " + tag + " "
	avail := width - lipgloss.Width(prefix)
	if avail < 4 {
		avail = 4
	}
	return prefix + truncate(stripColor(body), avail)
}

// stripColor — no-op for now; bodies are short enough that lipgloss.Width works.
func stripColor(s string) string { return s }

func actionRow(k, v string) string {
	kb := lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render(k)
	return kb + " " + sFg1.Render(v)
}
