package onboarding

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// renderRailTop renders the step rail as a single horizontal strip used when
// the layout is too narrow for a left rail.
func renderRailTop(m Model, w int) string {
	parts := make([]string, 0, len(data.OnbSteps))
	for i, s := range data.OnbSteps {
		done := i < m.Step
		active := i == m.Step
		var marker string
		switch {
		case done:
			marker = theme.SOK.Render("✓")
		case active:
			marker = theme.SAccent.Render(ui.Itoa(i + 1))
		default:
			marker = theme.SMute.Render(ui.Itoa(i + 1))
		}
		label := theme.SFaint.Render(s.Title)
		if active {
			label = theme.SAccent.Render(s.Title)
		} else if done {
			label = theme.SDim.Render(s.Title)
		}
		parts = append(parts, marker+" "+label)
	}
	line := strings.Join(parts, theme.SFaint.Render(" › "))
	return theme.WithBg(ui.PadOrClipANSI(" "+line, w), theme.Bg2)
}

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
