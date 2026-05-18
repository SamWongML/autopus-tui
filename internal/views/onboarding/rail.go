package onboarding

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderRail(m Model, w, h int) string {
	var b strings.Builder
	for i, s := range data.OnbSteps {
		done := i < m.Step
		active := i == m.Step
		var marker, label string
		switch {
		case done:
			marker = theme.SOK.Render("✓")
			label = theme.SDim.Render(s.Title)
		case active:
			marker = theme.SAccent.Render("▸")
			label = theme.SAccent.Render(s.Title)
		default:
			marker = theme.SMute.Render(ui.Itoa(i + 1))
			label = theme.SFaint.Render(s.Title)
		}
		line := marker + " " + label
		if active {
			line = theme.WithBg(ui.PadRight(line, w-4), theme.AccentFaint)
		}
		b.WriteString(line + "\n")
		if active {
			b.WriteString("    " + theme.SFaint.Render(ui.Truncate(s.Sub, w-8)) + "\n")
		}
	}
	b.WriteString("\n" + theme.SFaint.Render("writes to ") +
		theme.SDim.Render("~/.autopus/profiles/default/config.toml"))
	return ui.Panel("first-run setup", theme.SAccent.Render("◆"), b.String(), w, h, false, false)
}
