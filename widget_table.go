package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	tableCompact = iota
	tableFull
)

// taskTable renders the body (header + divider + rows) of a tasks table for
// either the Status tab miniature (mode=tableCompact, up to 6 rows, narrower
// columns) or the Tasks tab full roster (mode=tableFull, every row, wider
// columns). The caller wraps the returned body with pane().
func taskTable(rows []Task, w, h, mode, sel int) string {
	_ = h
	inner := w - 4
	var statusW, idW, agentW, seqW, elapsedW, costW int
	if mode == tableFull {
		statusW, idW, agentW, seqW, elapsedW, costW = 2, 7, 10, 7, 8, 7
	} else {
		statusW, idW, agentW, seqW, elapsedW, costW = 2, 7, 8, 8, 8, 6
	}
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

	lines := make([]string, 0, 2+len(rows))
	lines = append(lines, head)
	lines = append(lines, sFaint.Render(strings.Repeat("─", inner)))

	limit := len(rows)
	if mode == tableCompact && limit > 6 {
		limit = 6
	}
	for i, t := range rows[:limit] {
		lines = append(lines, taskRowRender(t, i == sel, mode, statusW, idW, titleW, agentW, seqW, elapsedW, costW))
	}
	return strings.Join(lines, "\n")
}

func taskRowRender(t Task, selected bool, mode, statusW, idW, titleW, agentW, seqW, elapsedW, costW int) string {
	if selected {
		bgOn := func(c lipgloss.Color) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(c).Background(bgSel)
		}
		// Compact mode: the accent bar absorbs one cell of the status column,
		// so the glyph pads to statusW-1. Full mode: the bar is prefixed
		// before a full-width status pad so unselected rows (which gain a
		// leading-space gutter) stay column-aligned.
		glyphW := statusW
		idFg := fg1
		if mode == tableCompact {
			glyphW = statusW - 1
			idFg = fg
		}
		seg := bgOn(accent).Render("▎") +
			lipgloss.NewStyle().
				Foreground(statusStyle(t.Status).GetForeground()).
				Background(bgSel).Render(padR(statusGlyph(t.Status), glyphW)) +
			bgOn(idFg).Render(" "+padR(t.ID, idW)) +
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
	prefix := ""
	if mode == tableFull {
		prefix = " "
	}
	return prefix + status + " " + id + " " + title + " " + runtime + " " + seq + " " + elapsed + " " + cost
}
