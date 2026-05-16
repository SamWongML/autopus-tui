package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderLog(width, height int) string {
	filter := renderLogFilter(width)
	bodyH := height - 3
	if bodyH < 6 {
		bodyH = 6
	}
	gap := 2
	rightW := 38
	if rightW > width*30/100 {
		rightW = width * 30 / 100
	}
	leftW := width - rightW - gap

	left := paneLogTable(leftW, bodyH)
	right := paneLogSummary(rightW, bodyH)
	body := joinH(left, bgPad(gap), right)
	return joinV(filter, body)
}

func renderLogFilter(width int) string {
	chip := func(s string, on bool) string {
		st := lipgloss.NewStyle().Padding(0, 1)
		if on {
			st = st.Foreground(accent).Background(bg2)
		} else {
			st = st.Foreground(dim).Background(bg1)
		}
		return st.Render(s)
	}
	left := sDim.Render("level") + "  " +
		chip("trace", true) + " " +
		chip("info", true) + " " +
		chip("warn", true) + " " +
		chip("error", true)
	right := sDim.Render("src: ") + sFg.Render("any") +
		"   " + sDim.Render("follow: ") + sOk.Render("● on")
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

func paneLogTable(width, height int) string {
	inner := width - 4
	// Pretend we have a lot of log lines by tripling the sample, like the design.
	tripled := append(append(append([]LogLine{}, logLines...), logLines...), logLines[:4]...)

	rowsAvail := height - 2 // top + bottom border
	if rowsAvail > len(tripled) {
		rowsAvail = len(tripled)
	}
	var lines []string
	for i, l := range tripled[:rowsAvail] {
		lines = append(lines, logRowText(l, i == 0, inner))
	}
	body := strings.Join(lines, "\n")
	return pane(daemon.Log, "34,182 lines · 12.4 MB · rotating @ 50 MB", body, width, height, false)
}

func paneLogSummary(width, height int) string {
	inner := width - 4
	var lines []string
	lines = append(lines, kv("lines", "34,182", "", inner))
	lines = append(lines, kv("errors", "3", "err", inner))
	lines = append(lines, kv("warns", "12", "warn", inner))
	lines = append(lines, kv("infos", "2,891", "", inner))
	lines = append(lines, kv("traces", "31,276", "dim", inner))
	lines = append(lines, hr("lines / minute", inner))
	lines = append(lines, sparklineBars([]int{8, 9, 8, 12, 10, 9, 11, 13, 9, 8, 11, 14, 10, 9, 12, 11, 8, 9, 10, 12}, inner-2))
	lines = append(lines, hr("top sources", inner))
	lines = append(lines, countBar("poll", 11240, 34182, "dim", inner))
	lines = append(lines, countBar("heartbeat", 3369, 34182, "ok", inner))
	lines = append(lines, countBar("rt.claude", 9114, 34182, "amber", inner))
	lines = append(lines, countBar("rt.codex", 5102, 34182, "info", inner))
	lines = append(lines, countBar("server.evt", 4118, 34182, "dim", inner))
	lines = append(lines, countBar("other", 1239, 34182, "dim", inner))
	body := strings.Join(lines, "\n")
	return pane("summary · today", "", body, width, height, false)
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

// sparklineBars renders a single-line bar sparkline using block characters.
func sparklineBars(data []int, width int) string {
	if len(data) == 0 || width < 1 {
		return ""
	}
	// Map each data point to a block height char.
	blocks := []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	max := 1
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	// Pick stride so all data fits in `width`.
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
	return sAccent.Render(b.String())
}
