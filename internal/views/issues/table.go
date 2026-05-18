package issues

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderTable(m Model, rows []data.Issue, w, h int) string {
	const (
		cPri  = 2
		cID   = 8
		cStat = 13
		cPri2 = 8
		cAssn = 18
		cUpd  = 4
	)
	gaps := 6
	cTitle := w - cPri - cID - cStat - cPri2 - cAssn - cUpd - gaps
	if cTitle < 12 {
		cTitle = 12
	}
	header := []string{
		theme.SFaint.Render(ui.PadRight("", cPri)),
		theme.SFaint.Render(ui.PadRight("ID", cID)),
		theme.SFaint.Render(ui.PadRight("TITLE", cTitle)),
		theme.SFaint.Render(ui.PadRight("STATUS", cStat)),
		theme.SFaint.Render(ui.PadRight("PRIORITY", cPri2)),
		theme.SFaint.Render(ui.PadRight("ASSIGNEE", cAssn)),
		theme.SFaint.Render(ui.PadLeft("UPD", cUpd)),
	}
	out := []string{strings.Join(header, " "), theme.SBorder.Render(strings.Repeat("─", w))}

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
		line := marker + " " +
			theme.Fg(pm.Color).Render(pm.Glyph) + " " +
			theme.SDim.Render(ui.PadRight(it.ID, cID)) + " " +
			theme.SText.Render(ui.PadRight(it.Title, cTitle)) + " " +
			theme.Fg(sc).Render(ui.PadRight(it.Status, cStat)) + " " +
			theme.Fg(pm.Color).Render(ui.PadRight(it.Priority, cPri2)) + " " +
			theme.SFaint.Render(ui.PadRight(it.Assignee, cAssn)) + " " +
			theme.SFaint.Render(ui.PadLeft(it.Updated, cUpd))
		if selected {
			line = theme.WithBg(line, theme.AccentFaint)
		}
		out = append(out, line)
	}
	for len(out)-2 < visibleH {
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}
