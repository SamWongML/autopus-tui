package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderRun(width, height int) string {
	headH := 3
	body := runBody(width, height-headH)
	return joinV(runHeader(width), body)
}

func runHeader(width int) string {
	id := lipgloss.NewStyle().Foreground(accent).Bold(true).Render("t-1284")
	title := sFg.Render("refactor billing webhook → idempotency keys")
	agent := agentChip("claude")
	meta := sDim.Render("ws ") + sFg1.Render("blackhole-os") + "  " +
		sDim.Render("issue ") + sFg1.Render("#4127")
	live := lipgloss.NewStyle().Foreground(ok).Background(bg2).
		Padding(0, 1).Render("● LIVE")
	state := sWarn.Render("! waiting")

	left := id + "  " + title + "  " + agent + "  " + meta
	right := live + "  " + state
	avail := width - 4 // padding 0,2
	gap := avail - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		left = truncate(left, avail-lipgloss.Width(right)-1)
		gap = 1
	}
	row := fillBg(left+strings.Repeat(" ", gap)+right, bg)
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		MaxWidth(width).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

func runBody(width, height int) string {
	gap := 2
	sidebarW := 38
	if sidebarW > width/3 {
		sidebarW = width / 3
	}
	msgsW := width - sidebarW - gap

	msgsBody := renderMessages(msgsW - 4)
	messagesPane := pane("messages · run-messages t-1284 --since 0", "9 of 86 · filter: all", msgsBody, msgsW, height, false)

	sidebarBody := renderTaskSidebar(sidebarW - 4)
	sidebarPane := pane("task", "t-1284", sidebarBody, sidebarW, height, false)

	return joinH(messagesPane, bgPad(gap), sidebarPane)
}

func renderMessages(width int) string {
	var lines []string
	lines = append(lines, sFaint.Render("— start of stream · seq 78 "+strings.Repeat("─", maxInt(2, width-30))))
	for _, m := range msgs {
		lines = append(lines, msgRow(m, width)...)
	}
	lines = append(lines, sOk.Render("— waiting on user · live tail ● "+strings.Repeat("─", maxInt(2, width-32))))
	return strings.Join(lines, "\n")
}

func msgRow(m Msg, width int) []string {
	seqW, tW, typeW := 4, 9, 14
	bodyW := width - seqW - tW - typeW - 4
	if bodyW < 10 {
		bodyW = 10
	}

	seq := sFaint.Render(padR(fmt.Sprintf("%3d", m.Seq), seqW))
	t := sDim.Render(padR(m.T, tW))
	typeStyled := msgTypeTag(m.Type, typeW)

	var bodyText string
	switch m.Type {
	case "tool_call":
		bodyText = sInfo.Render(m.Tool) + " " + sFg1.Render(m.Args)
		if m.Lines > 0 {
			bodyText += sDim.Render(fmt.Sprintf(" · +%d lines", m.Lines))
		}
	case "tool_result":
		bodyText = sOk.Render(m.Body)
	case "thinking":
		bodyText = lipgloss.NewStyle().Foreground(fg1).Italic(true).Render(m.Body)
	case "text":
		bodyText = sFg.Render(m.Body)
	case "error":
		bodyText = sErr.Render(m.Body)
	}
	if m.Ok {
		bodyText += " " + sOk.Render("✓")
	}

	// Wrap body to bodyW.
	wrapped := wrap(bodyText, bodyW)
	var rows []string
	prefix := seq + " " + t + " " + typeStyled + " "
	cont := strings.Repeat(" ", lipgloss.Width(prefix))

	border := msgBorder(m.Type)
	for i, l := range wrapped {
		if i == 0 {
			rows = append(rows, border+prefix+l)
		} else {
			rows = append(rows, border+cont+l)
		}
	}
	return rows
}

func msgBorder(t string) string {
	switch t {
	case "thinking":
		return sFaint.Render("▎")
	case "tool_call":
		return sInfo.Render("▎")
	case "tool_result":
		return sOk.Render("▎")
	case "text":
		return sAccent.Render("▎")
	case "error":
		return sErr.Render("▎")
	}
	return " "
}

func msgTypeTag(t string, w int) string {
	pad := padR(t, w)
	switch t {
	case "thinking":
		return sDim.Render(pad)
	case "tool_call":
		return sInfo.Render(pad)
	case "tool_result":
		return sOk.Render(pad)
	case "text":
		return sAccent.Render(pad)
	case "error":
		return sErr.Render(pad)
	}
	return sDim.Render(pad)
}

func wrap(s string, width int) []string {
	if width < 1 {
		return []string{s}
	}
	// Word-aware wrap that ignores ANSI escapes for measurement.
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

func renderTaskSidebar(width int) string {
	rows := [][2]string{
		{"status", "● waiting"},
		{"runtime", "claude"},
		{"model", "sonnet-4.6"},
		{"pid", "18421"},
		{"started", "14:28:50"},
		{"elapsed", "3m 12s"},
		{"messages", "86"},
		{"tokens in", "124,118"},
		{"tokens out", "59,902"},
		{"cost", "$1.42"},
		{"branch", "feat/idempotency-keys"},
	}
	var lines []string
	for _, r := range rows {
		tone := ""
		if r[0] == "status" {
			tone = "warn"
		}
		lines = append(lines, kv(r[0], r[1], tone, width))
	}
	lines = append(lines, hr("files touched", width))
	lines = append(lines, fileRow("db/migrate/20260514_add_idempotency_key.rb", "+42", width))
	lines = append(lines, fileRow("app/models/webhook_event.rb", "+8 −2", width))
	lines = append(lines, fileRow("spec/models/webhook_event_spec.rb", "+24", width))
	lines = append(lines, hr("msg type counts", width))
	lines = append(lines, countBar("thinking", 32, 86, "dim", width))
	lines = append(lines, countBar("tool_call", 28, 86, "info", width))
	lines = append(lines, countBar("tool_result", 20, 86, "ok", width))
	lines = append(lines, countBar("text", 5, 86, "amber", width))
	lines = append(lines, countBar("error", 1, 86, "err", width))
	return strings.Join(lines, "\n")
}

func fileRow(path, diff string, width int) string {
	gap := width - lipgloss.Width(path) - lipgloss.Width(diff)
	if gap < 1 {
		gap = 1
		path = truncate(path, width-lipgloss.Width(diff)-1)
	}
	return sFg1.Render(path) + strings.Repeat(" ", gap) + sOk.Render(diff)
}

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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
