package issues

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// colSpecs is the responsive column layout for the issues table.
//
// MinW: minimum cell width before the column is dropped.
// Weight: share of leftover space (0 = fixed at MinW).
// Priority: higher value = dropped first when width is tight.
//
// "marker" carries the selection bar and is pinned (Priority 0). "pri" is the
// priority glyph; "priority" is the spelled-out word column.
var colSpecs = []ui.Col{
	{ID: "marker", MinW: 1, Priority: 0},
	{ID: "pri", MinW: 2, Priority: 1},
	{ID: "id", Header: "ID", MinW: 7, Priority: 3},
	{ID: "title", Header: "TITLE", MinW: 20, Weight: 5, Priority: 1},
	{ID: "status", Header: "STATUS", MinW: 11, Priority: 2},
	{ID: "priority", Header: "PRIORITY", MinW: 8, Priority: 4},
	{ID: "assignee", Header: "ASSIGNEE", MinW: 12, Priority: 6},
	{ID: "upd", Header: "UPD", MinW: 4, Priority: 7, AlignR: true},
}

func renderTable(m Model, rows []data.Issue, w, h int) string {
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

	visibleH := h - 2
	if visibleH < 1 {
		visibleH = 1
	}
	start := 0
	if m.Sel >= visibleH {
		start = m.Sel - visibleH + 1
	}
	if start+visibleH > len(rows) {
		start = ui.Max(0, len(rows)-visibleH)
	}
	end := ui.Min(start+visibleH, len(rows))

	for i := start; i < end; i++ {
		it := rows[i]
		selected := i == m.Sel
		marker := " "
		if selected {
			marker = theme.SAccent.Render("▎")
		}
		pm := theme.Priorities[it.Priority]
		sc := theme.Statuses[it.Status]
		if sc == "" {
			sc = theme.TextDim
		}
		text := map[string]string{
			"marker":   marker,
			"pri":      theme.Fg(pm.Color).Render(pm.Glyph),
			"id":       theme.SDim.Render(it.ID),
			"title":    theme.SText.Render(it.Title),
			"status":   theme.Fg(sc).Render(it.Status),
			"priority": theme.Fg(pm.Color).Render(it.Priority),
			"assignee": theme.SFaint.Render(it.Assignee),
			"upd":      theme.SFaint.Render(it.Updated),
		}
		row := ui.RowFromCols(visible, widths, text, 1)
		if selected {
			row = theme.WithBg(row, theme.AccentFaint)
		}
		out = append(out, row)
	}
	for len(out)-2 < visibleH {
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}
