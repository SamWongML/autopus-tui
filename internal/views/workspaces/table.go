package workspaces

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderTable(m Model, w, h int) string {
	const (
		cDot  = 2
		cMem  = 8
		cIss  = 7
		cSess = 9
		cRole = 12
	)
	gaps := 5
	cName := w - cDot - cMem - cIss - cSess - cRole - gaps
	if cName < 12 {
		cName = 12
	}
	header := []string{
		theme.SFaint.Render(ui.PadRight("", cDot)),
		theme.SFaint.Render(ui.PadRight("NAME", cName)),
		theme.SFaint.Render(ui.PadLeft("MEMBERS", cMem)),
		theme.SFaint.Render(ui.PadLeft("ISSUES", cIss)),
		theme.SFaint.Render(ui.PadLeft("SESSIONS", cSess)),
		theme.SFaint.Render(ui.PadRight("ROLE", cRole)),
	}
	out := []string{strings.Join(header, " "), theme.SBorder.Render(strings.Repeat("─", w))}

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
		line := marker + " " +
			theme.Fg(watchCol).Render(watchGlyph) + " " +
			theme.Fg(nameCol).Render(ui.PadRight(x.Name, cName)) + " " +
			theme.SDim.Render(ui.PadLeft(ui.Itoa(x.Members), cMem)) + " " +
			theme.SDim.Render(ui.PadLeft(ui.Itoa(x.Issues), cIss)) + " " +
			theme.Fg(sessCol).Render(ui.PadLeft(ui.Itoa(x.Sessions), cSess)) + " " +
			theme.SFaint.Render(ui.PadRight(x.Role, cRole))
		if selected {
			line = theme.WithBg(line, theme.AccentFaint)
		}
		out = append(out, line)
	}
	_ = h
	return strings.Join(out, "\n")
}
