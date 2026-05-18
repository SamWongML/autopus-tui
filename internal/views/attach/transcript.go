package attach

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func renderTranscript(m Model, c ctx.Ctx, w, h int) string {
	var lines []string
	for _, e := range data.SessionTranscript {
		switch e.Kind {
		case "user":
			lines = append(lines, theme.SAccent.Render("YOU · "+e.Time))
			lines = append(lines, theme.SText.Render(ui.Wrap(e.Body, w)))
			lines = append(lines, "")
		case "plan":
			lines = append(lines, theme.SViolet.Render("◇ PLAN · "+e.Time))
			for _, l := range strings.Split(e.Body, "\n") {
				lines = append(lines, theme.SDim.Render(l))
			}
			lines = append(lines, "")
		case "tool":
			lines = append(lines, theme.SMute.Render(e.Time)+" "+ui.Dot(theme.Info)+" "+theme.SInfo.Render(e.Tool)+" "+theme.SDim.Render(e.Arg))
		case "thinking":
			for _, l := range strings.Split(ui.Wrap(e.Body, w-8), "\n") {
				prefix := theme.SBorder.Render("│ ")
				lines = append(lines, "       "+prefix+lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(theme.TextFaint)).Render(l))
			}
		case "ask":
			lines = append(lines, "")
			header := theme.SAccent.Render("◆ QUESTION · " + e.Time)
			body := ui.Wrap(e.Body, w-4)
			boxed := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.AccentDim)).
				Foreground(lipgloss.Color(theme.Text)).
				Padding(0, 1).Width(w-2).Render(header + "\n" + body)
			lines = append(lines, strings.Split(boxed, "\n")...)
		}
	}
	lines = append(lines, "")
	lines = append(lines, ui.Glyph("needs_input", c.Spin)+" "+theme.SAccent.Render("waiting for you")+" "+theme.SFaint.Render("· 02:11 idle"))

	start := m.Scroll
	if start > len(lines)-h {
		start = ui.Max(0, len(lines)-h)
	}
	end := ui.Min(start+h, len(lines))
	return strings.Join(lines[start:end], "\n")
}
