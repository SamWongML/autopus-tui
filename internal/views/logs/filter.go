package logs

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
)

func renderFilterStrip(m Model, w int) string {
	_ = w
	var parts []string
	for _, l := range data.LogLevelFilters {
		active := l == m.Level
		col := theme.LogLevels[l]
		if col == "" {
			col = theme.TextDim
		}
		border := theme.Border
		if active {
			col = theme.Accent
			border = theme.AccentDim
		}
		styled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(col)).
			Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(border)).
			Padding(0, 1).Render(l)
		parts = append(parts, styled)
	}
	followCol, followBorder := theme.TextDim, theme.Border
	followText := "○ paused"
	if m.Follow {
		followCol, followBorder = theme.Accent, theme.AccentDim
		followText = "● following"
	}
	follow := lipgloss.NewStyle().
		Foreground(lipgloss.Color(followCol)).
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(followBorder)).
		Padding(0, 1).Render(followText)
	parts = append(parts, follow)

	all := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	lines := strings.Split(all, "\n")
	if len(lines) >= 2 {
		return lines[1]
	}
	return all
}
