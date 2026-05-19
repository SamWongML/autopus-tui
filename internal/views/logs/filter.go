package logs

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// filterHitBoxes returns the click hit-boxes for the level chips and the
// follow toggle, mirroring renderFilterStrip. Each chip is padded(0,1) with a
// 1-cell border (= label + 4 cells). Chips are JoinHorizontal'd starting at
// x=0 with no separator. The follow toggle ID is "__follow".
func filterHitBoxes() []ui.HitBox {
	hits := make([]ui.HitBox, 0, len(data.LogLevelFilters)+1)
	x := 0
	for _, l := range data.LogLevelFilters {
		chipW := lipgloss.Width(l) + 4
		hits = append(hits, ui.HitBox{X1: x, X2: x + chipW - 1, ID: l})
		x += chipW
	}
	follow := "● following"
	chipW := lipgloss.Width(follow) + 4
	hits = append(hits, ui.HitBox{X1: x, X2: x + chipW - 1, ID: "__follow"})
	return hits
}

func renderFilterStrip(m Model, spin, w int) string {
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
		dot := "●"
		if !ui.BlinkOn(spin) {
			dot = "○"
		}
		followText = dot + " following"
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
