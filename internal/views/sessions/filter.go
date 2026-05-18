package sessions

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
)

func renderFilterStrip(m Model, w int) string {
	var parts []string
	counts := map[string]int{"all": len(data.Sessions)}
	for _, s := range data.Sessions {
		counts[s.State]++
	}
	for _, f := range data.SessionFilters {
		active := f == m.Filter
		var label, color, g string
		if f == "all" {
			label = "all"
			color = theme.TextDim
		} else {
			meta := theme.State(f)
			label = meta.Label
			color = meta.Color
			g = theme.Fg(meta.Color).Render(meta.Glyph)
		}
		txt := label + " " + fmt.Sprintf("%d", counts[f])
		var styled string
		if active {
			styled = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Accent)).
				Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.AccentDim)).
				Padding(0, 1).Render(g + " " + txt)
		} else {
			styled = lipgloss.NewStyle().
				Foreground(lipgloss.Color(color)).
				Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.Border)).
				Padding(0, 1).Render(g + " " + txt)
		}
		parts = append(parts, styled)
	}

	searchBorder := theme.Border
	slashCol := theme.TextFaint
	queryStr := theme.SFaint.Render("filter sessions…")
	caret := ""
	if m.Searching {
		searchBorder = theme.AccentDim
		slashCol = theme.Accent
		caret = theme.SAccent.Render("▌")
		if m.Query != "" {
			queryStr = theme.SText.Render(m.Query)
		}
	} else if m.Query != "" {
		queryStr = theme.SText.Render(m.Query)
	}
	searchInner := theme.Fg(slashCol).Render("/ ") + queryStr + caret
	search := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(searchBorder)).
		Padding(0, 1).Render(searchInner)

	all := lipgloss.JoinHorizontal(lipgloss.Top, append(parts, search)...)
	lines := strings.Split(all, "\n")
	pillsLine := all
	if len(lines) >= 2 {
		pillsLine = lines[1]
	}

	hint := theme.SMute.Render("[") + " " + theme.SMute.Render("]") + " " + theme.SFaint.Render("cycle")
	used := lipgloss.Width(pillsLine)
	if used+lipgloss.Width(hint)+2 <= w {
		pillsLine = pillsLine + "  " + hint
	}
	return pillsLine
}
