package sessions

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func renderTable(m Model, c ctx.Ctx, rows []data.Session, w, h int) string {
	if w < 50 {
		w = 50
	}
	const (
		cState = 2
		cIssue = 8
		cAgent = 14
		cWS    = 12
		cElap  = 9
		cCost  = 7
	)
	gaps := 6
	cTitle := w - cState - cIssue - cAgent - cWS - cElap - cCost - gaps
	if cTitle < 12 {
		cTitle = 12
	}

	header := []string{
		theme.SFaint.Render(ui.PadRight("", cState)),
		theme.SFaint.Render(ui.PadRight("ISSUE", cIssue)),
		theme.SFaint.Render(ui.PadRight("TITLE · ACTIVITY", cTitle)),
		theme.SFaint.Render(ui.PadRight("AGENT", cAgent)),
		theme.SFaint.Render(ui.PadRight("WORKSPACE", cWS)),
		theme.SFaint.Render(ui.PadLeft("ELAPSED", cElap)),
		theme.SFaint.Render(ui.PadLeft("COST", cCost)),
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
		s := rows[i]
		selected := i == m.Sel
		marker := " "
		if selected {
			marker = theme.SAccent.Render("▎")
		}
		g := ui.Glyph(s.State, c.Spin)
		issue := theme.SDim.Render(ui.PadRight(s.Issue, cIssue))
		title := ui.PadRight(theme.SText.Render(s.Title)+"  "+activityHint(s), cTitle)
		agent := theme.SDim.Render(ui.PadRight(s.Agent+" · "+s.Model, cAgent))
		ws := theme.SFaint.Render(ui.PadRight(s.Workspace, cWS))
		elap := theme.SDim.Render(ui.PadLeft(s.Elapsed, cElap))
		cost := theme.SDim.Render(ui.PadLeft(s.Cost, cCost))

		row := marker + " " + g + " " + issue + " " + title + " " + agent + " " + ws + " " + elap + " " + cost
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

func activityHint(s data.Session) string {
	col := theme.TextFaint
	if s.State == "needs_input" {
		col = theme.Warn
	}
	return theme.Fg(col).Render(s.Activity)
}
