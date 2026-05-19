package runtimes

import (
	"fmt"
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// colSpecs is the responsive column layout for the runtimes table.
//
// MinW: minimum cell width before the column is dropped.
// Weight: share of leftover space (0 = fixed at MinW).
// Priority: higher value = dropped first when width is tight.
//
// "marker" and "state" are pinned (Priority 0). Drop order at narrow widths:
// PATH → LAST USED → MODEL.
var colSpecs = []ui.Col{
	{ID: "marker", MinW: 1, Priority: 0},
	{ID: "state", MinW: 2, Priority: 0},
	{ID: "cli", Header: "CLI", MinW: 10, Weight: 2, Priority: 1},
	{ID: "ver", Header: "VER", MinW: 7, Priority: 3},
	{ID: "path", Header: "PATH", MinW: 12, Weight: 5, Priority: 6},
	{ID: "model", Header: "MODEL", MinW: 12, Priority: 4},
	{ID: "status", Header: "STATUS", MinW: 9, Priority: 2},
	{ID: "incap", Header: "IN/CAP", MinW: 7, Priority: 5, AlignR: true},
	{ID: "last", Header: "LAST USED", MinW: 10, Priority: 7},
}

func renderTable(m Model, w, h int) string {
	widths, visible := ui.LayoutCols(colSpecs, w)

	header := make(map[string]string, len(visible))
	for _, col := range visible {
		header[col.ID] = theme.SFaint.Render(col.Header)
	}
	headerRow := ui.RowFromCols(visible, widths, header, 1)

	out := []string{
		headerRow,
		theme.SBorder.Render(strings.Repeat("─", w)),
	}

	visH := h - 2
	if visH < 1 {
		visH = 1
	}
	for i, r := range data.Runtimes {
		selected := i == m.Sel
		marker := " "
		if selected {
			marker = theme.SAccent.Render("▎")
		}
		var col, dotG string
		switch r.Status {
		case "ready":
			col, dotG = theme.OK, "●"
		case "stale":
			col, dotG = theme.Warn, "◐"
		case "not_found":
			col, dotG = theme.TextMute, "✕"
		default:
			col, dotG = theme.TextFaint, "○"
		}
		cliCol := theme.Text
		if r.Status == "not_found" {
			cliCol = theme.TextMute
		}
		statusText := r.Status
		if r.Status == "not_found" {
			statusText = "missing"
		}
		text := map[string]string{
			"marker": marker,
			"state":  theme.Fg(col).Render(dotG),
			"cli":    theme.Fg(cliCol).Render(r.CLI),
			"ver":    theme.SDim.Render(r.Version),
			"path":   theme.SFaint.Render(r.Path),
			"model":  theme.SDim.Render(r.Model),
			"status": theme.Fg(col).Render(statusText),
			"incap":  theme.SDim.Render(fmt.Sprintf("%d/%d", r.Inflight, r.Cap)),
			"last":   theme.SFaint.Render(r.LastUsed),
		}
		row := ui.RowFromCols(visible, widths, text, 1)
		if selected {
			row = theme.WithBg(row, theme.AccentFaint)
		}
		out = append(out, row)
	}
	for len(out)-2 < visH {
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}
