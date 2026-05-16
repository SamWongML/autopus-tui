package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
)

func renderLog(width, height int, vp viewport.Model) string {
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

	left := paneLogTable(leftW, bodyH, vp)
	right := paneLogSummary(rightW, bodyH)
	body := joinH(left, bgPad(gap), right)
	return joinV(filter, body)
}

func renderLogFilter(width int) string {
	left := chipBar("level",
		chipSpec{"trace", true},
		chipSpec{"info", true},
		chipSpec{"warn", true},
		chipSpec{"error", true},
	)
	right := sDim.Render("src: ") + sFg.Render("any") +
		"   " + sDim.Render("follow: ") + sOk.Render("● on")
	return pageHeader(width, left, right)
}

func paneLogTable(width, height int, vp viewport.Model) string {
	// Body comes pre-rendered from the viewport (see resizeViewports() in
	// main.go). pane() paints chrome around the visible scroll window.
	return pane(daemon.Log, "34,182 lines · 12.4 MB · rotating @ 50 MB", vp.View(), width, height, false)
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
