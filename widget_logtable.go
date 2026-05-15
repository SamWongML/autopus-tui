package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	logTail = iota // capped at min(8, len(lines)); fits Status tab's miniature
	logFull        // fills height-2 rows; fits Log tab's main table
)

// logTable renders a sequence of log rows for either the Status tab tail
// (mode=logTail, max 8 rows) or the Log tab table (mode=logFull, fills the
// pane's interior). The first row is highlighted with an accent bar; all
// others get a leading space so columns align. Caller wraps with pane().
func logTable(lines []LogLine, w, h, mode int) string {
	inner := w - 4
	var limit int
	if mode == logTail {
		limit = len(lines)
		if limit > 8 {
			limit = 8
		}
	} else {
		limit = h - 2
		if limit > len(lines) {
			limit = len(lines)
		}
	}
	if limit < 0 {
		limit = 0
	}
	out := make([]string, 0, limit)
	for i, l := range lines[:limit] {
		out = append(out, logRowText(l, i == 0, inner))
	}
	return strings.Join(out, "\n")
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
