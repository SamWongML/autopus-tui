package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderTasks(width, height int, selected int) string {
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

	left := paneTasksTable(leftW, bodyH, selected)
	right := paneTaskPreview(rightW, bodyH)
	body := joinH(left, bgPad(gap), right)

	return joinV(filterRow, body)
}

func renderFilterBar(width int) string {
	chip := func(s string, on bool) string {
		st := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, false, false)
		if on {
			st = st.Foreground(accent).Background(bg2)
		} else {
			st = st.Foreground(dim).Background(bg1)
		}
		return st.Render(s)
	}
	left := sDim.Render("filter") + "  " +
		chip("status: in-flight", true) + " " +
		chip("runtime: any", false) + " " +
		chip("workspace: blackhole-os", false) + " " +
		chip("since: today", false)
	right := sDim.Render("sort: ") + sFg.Render("started ↓")
	gap := width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	row := fillBg(left+strings.Repeat(" ", gap)+right, bg)
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

func paneTasksTable(width, height int, selected int) string {
	inner := width - 4
	statusW, idW, agentW, seqW, elapsedW, costW := 2, 7, 10, 7, 8, 7
	titleW := inner - statusW - idW - agentW - seqW - elapsedW - costW - 6
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
	for i, t := range tasks {
		lines = append(lines, taskRowSimple(t, i == selected, statusW, idW, titleW, agentW, seqW, elapsedW, costW))
	}
	body := strings.Join(lines, "\n")
	return pane(fmt.Sprintf("%d tasks", len(tasks)), "", body, width, height, false)
}

// taskRowSimple — cleaner than taskRow in view_status.go; used by Tasks view.
// Non-selected rows sit directly on the canvas bg (no row stripe); the
// selected row is painted with bgSel end-to-end via fillBg, so the spaces
// between fields don't fall back to canvas bg and produce a vertical strap.
func taskRowSimple(t Task, selected bool, statusW, idW, titleW, agentW, seqW, elapsedW, costW int) string {
	if selected {
		bgOn := func(c lipgloss.Color) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(c).Background(bgSel)
		}
		seg := bgOn(accent).Render("▎") +
			lipgloss.NewStyle().
				Foreground(statusStyle(t.Status).GetForeground()).
				Background(bgSel).Render(padR(statusGlyph(t.Status), statusW)) +
			bgOn(fg1).Render(" "+padR(t.ID, idW)) +
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
	return " " + status + " " + id + " " + title + " " + runtime + " " + seq + " " + elapsed + " " + cost
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
	var lines []string
	for _, r := range rows {
		tone := ""
		if r[0] == "status" {
			tone = "accent"
		}
		if r[0] == "cwd" {
			tone = "dim"
		}
		lines = append(lines, kv(r[0], r[1], tone, inner))
	}
	lines = append(lines, hr("last 5 messages", inner))
	for _, m := range msgs[len(msgs)-5:] {
		lines = append(lines, msgRowCompact(m, inner))
	}
	lines = append(lines, hr("actions", inner))
	lines = append(lines, actionRow("⏎", "open run viewer"))
	lines = append(lines, actionRow("r", "reply (waiting tasks only)"))
	lines = append(lines, actionRow("k", "kill child process (SIGTERM)"))
	lines = append(lines, actionRow("K", "force kill (SIGKILL)"))
	lines = append(lines, actionRow("o", "open workspace in $EDITOR"))
	lines = append(lines, actionRow("c", "copy task id"))
	body := strings.Join(lines, "\n")
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
