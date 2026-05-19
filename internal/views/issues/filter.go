package issues

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// filterHitBoxes returns the click hit-boxes for the issue filter chips,
// mirroring the widths used in renderFilterStrip. Chips are padded(0,1) with
// a 1-cell border, so each chip width = len(label) + 4. They render with no
// separator (JoinHorizontal), starting at x=0.
func filterHitBoxes() []ui.HitBox {
	hits := make([]ui.HitBox, 0, len(data.IssueFilters))
	x := 0
	for _, f := range data.IssueFilters {
		chipW := lipgloss.Width(f) + 4
		hits = append(hits, ui.HitBox{X1: x, X2: x + chipW - 1, ID: f})
		x += chipW
	}
	return hits
}

func renderFilterStrip(m Model, w int) string {
	var parts []string
	for _, f := range data.IssueFilters {
		active := f == m.Filter
		col, border := theme.TextDim, theme.Border
		if active {
			col, border = theme.Accent, theme.AccentDim
		}
		styled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(col)).
			Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(border)).
			Padding(0, 1).Render(f)
		parts = append(parts, styled)
	}
	all := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	lines := strings.Split(all, "\n")
	pillsLine := all
	if len(lines) >= 2 {
		pillsLine = lines[1]
	}

	cycle := theme.SMute.Render("[") + " " + theme.SMute.Render("]") + " " + theme.SFaint.Render("cycle")
	left := pillsLine + "  " + cycle

	actions := ui.KeyChip("n", "new", true) + "  " +
		ui.KeyChip("a", "assign", false) + "  " +
		ui.KeyChip("s", "status", false)

	gap := w - lipgloss.Width(left) - lipgloss.Width(actions)
	if gap < 2 {
		gap = 2
	}
	return left + strings.Repeat(" ", gap) + actions
}
