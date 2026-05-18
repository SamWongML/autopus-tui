package workspaces

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// colSpecs is the responsive column layout for the workspaces table.
//
// "marker" and "watch" dots are pinned (Priority 0). Drop order at narrow
// widths: ROLE → MEMBERS.
var colSpecs = []ui.Col{
	{ID: "marker", MinW: 1, Priority: 0},
	{ID: "watch", MinW: 2, Priority: 0},
	{ID: "name", Header: "NAME", MinW: 12, Weight: 5, Priority: 1},
	{ID: "members", Header: "MEMBERS", MinW: 7, Priority: 4, AlignR: true},
	{ID: "issues", Header: "ISSUES", MinW: 6, Priority: 3, AlignR: true},
	{ID: "sessions", Header: "SESSIONS", MinW: 8, Priority: 2, AlignR: true},
	{ID: "role", Header: "ROLE", MinW: 10, Priority: 5},
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

	for i, x := range data.Workspaces {
		selected := i == m.Sel
		marker := " "
		if selected {
			marker = theme.SAccent.Render("▎")
		}
		var watchGlyph, watchCol string
		if x.Watch {
			watchGlyph, watchCol = "◆", theme.Accent
		} else {
			watchGlyph, watchCol = "◇", theme.TextMute
		}
		nameCol := theme.Text
		if !x.Watch {
			nameCol = theme.TextDim
		}
		sessCol := theme.TextMute
		if x.Sessions > 0 {
			sessCol = theme.OK
		}
		text := map[string]string{
			"marker":   marker,
			"watch":    theme.Fg(watchCol).Render(watchGlyph),
			"name":     theme.Fg(nameCol).Render(x.Name),
			"members":  theme.SDim.Render(ui.Itoa(x.Members)),
			"issues":   theme.SDim.Render(ui.Itoa(x.Issues)),
			"sessions": theme.Fg(sessCol).Render(ui.Itoa(x.Sessions)),
			"role":     theme.SFaint.Render(x.Role),
		}
		row := ui.RowFromCols(visible, widths, text, 1)
		if selected {
			row = theme.WithBg(row, theme.AccentFaint)
		}
		out = append(out, row)
	}
	_ = h
	return strings.Join(out, "\n")
}
