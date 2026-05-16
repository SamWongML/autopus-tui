package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderRun lays out the Run tab. By default the messages stream gets the full
// width; when `peek` is true a 38-cell task-detail overlay is joined on the
// right, sharing the same renderTaskDetail component the Tasks-tab preview
// uses. The toggle is wired to Space in main.go.
func renderRun(width, height int, peek bool) string {
	headH := 3
	body := runBody(width, height-headH, peek)
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

func runBody(width, height int, peek bool) string {
	if !peek {
		msgsBody := renderMessages(width - 4)
		return pane("messages · run-messages t-1284 --since 0", "9 of 86 · filter: all", msgsBody, width, height, false)
	}

	gap := 2
	peekW := 38
	if peekW > width/3 {
		peekW = width / 3
	}
	msgsW := width - peekW - gap

	msgsBody := renderMessages(msgsW - 4)
	messagesPane := pane("messages · run-messages t-1284 --since 0", "9 of 86 · filter: all", msgsBody, msgsW, height, false)

	t := tasks[0]
	detailBody := renderTaskDetail(t, "stream", peekW-4, height-2)
	detailPane := pane("task "+t.ID, "", detailBody, peekW, height, true)

	return joinH(messagesPane, bgPad(gap), detailPane)
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
