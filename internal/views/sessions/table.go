package sessions

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// colSpecs is the responsive column layout for the sessions table.
//
// MinW: minimum cell width before the column is dropped.
// Weight: share of leftover space (0 = fixed at MinW).
// Priority: higher value = dropped first when width is tight.
//
// The "marker" and "state" columns are pinned (Priority 0) so they're never
// dropped — they carry the selection bar and state glyph.
var colSpecs = []ui.Col{
	{ID: "marker", MinW: 1, Priority: 0},
	{ID: "state", MinW: 2, Priority: 0},
	{ID: "issue", Header: "ISSUE", MinW: 8, Priority: 3},
	{ID: "title", Header: "TITLE · ACTIVITY", MinW: 20, Weight: 5, Priority: 1},
	{ID: "agent", Header: "AGENT", MinW: 10, Priority: 4},
	{ID: "ws", Header: "WORKSPACE", MinW: 8, Priority: 5},
	{ID: "elap", Header: "ELAPSED", MinW: 8, Priority: 6, AlignR: true},
	{ID: "cost", Header: "COST", MinW: 6, Priority: 7, AlignR: true},
}

// renderTable renders the sessions table inside a content area of (w, h) cells.
// Every emitted line is exactly w cells wide — built via ui.Row which clips
// oversize cells rather than wrapping them onto a new line (which would break
// the panel's right border).
func renderTable(m Model, c ctx.Ctx, rows []data.Session, w, h int) string {
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
		s := rows[i]
		marker := " "
		if i == m.Sel {
			marker = theme.SAccent.Render("▎")
		}
		text := map[string]string{
			"marker": marker,
			"state":  ui.Glyph(s.State, c.Spin),
			"issue":  theme.SDim.Render(s.Issue),
			"title":  theme.SText.Render(s.Title) + "  " + activityHint(s),
			"agent":  theme.SDim.Render(s.Agent + " · " + s.Model),
			"ws":     theme.SFaint.Render(s.Workspace),
			"elap":   theme.SDim.Render(s.Elapsed),
			"cost":   theme.SDim.Render(s.Cost),
		}
		row := ui.RowFromCols(visible, widths, text, 1)
		if i == m.Sel {
			row = theme.WithBg(row, theme.AccentFaint)
		}
		out = append(out, row)
	}
	for len(out)-2 < visibleH {
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}

func activityHint(s data.Session) string {
	col := theme.TextFaint
	if s.State == "needs_input" {
		col = theme.Warn
	}
	return theme.Fg(col).Render(s.Activity)
}
