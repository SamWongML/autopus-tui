package overview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderSessionsSummary(w, h int) string {
	order := []string{"needs_input", "working", "running", "paused", "idle", "completed", "failed"}
	counts := data.CountSessionStates()
	cols := ui.Max(1, w/22)
	if cols > len(order) {
		cols = len(order)
	}
	cellW := (w - (cols - 1)) / cols
	if cellW < 18 {
		cellW = 18
		cols = ui.Max(1, (w+1)/(cellW+1))
	}
	var b strings.Builder
	for i := 0; i < len(order); i += cols {
		row := make([]string, 0, 2*cols-1)
		for j := 0; j < cols; j++ {
			if j > 0 {
				row = append(row, " ")
			}
			if i+j < len(order) {
				row = append(row, summaryChip(order[i+j], counts[order[i+j]], cellW))
			} else {
				row = append(row, strings.Repeat(" ", cellW))
			}
		}
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, row...) + "\n")
	}
	b.WriteString(ui.Dashed(w) + "\n")
	b.WriteString(theme.SWarn.Render("2 sessions need you") + " " + theme.SFaint.Render("· oldest waiting ") + theme.SText.Render("02:11"))
	_ = h
	return b.String()
}

func summaryChip(state string, count, w int) string {
	meta := theme.State(state)
	active := count > 0
	borderCol := theme.Border
	countCol := theme.TextMute
	if active {
		borderCol = meta.Color
		countCol = meta.Color
	}
	g := theme.Fg(meta.Color).Render(meta.Glyph)
	label := theme.SFaint.Render(strings.ToUpper(meta.Label))
	countStr := theme.Fg(countCol).Bold(true).Render(fmt.Sprintf("%d", count))

	innerW := w - 2
	if innerW < 4 {
		innerW = 4
	}
	leftPart := " " + g + " " + label
	leftW := lipgloss.Width(leftPart)
	rightW := lipgloss.Width(countStr) + 1
	gap := innerW - leftW - rightW
	if gap < 1 {
		gap = 1
	}
	content := leftPart + strings.Repeat(" ", gap) + countStr + " "

	bs := theme.Fg(borderCol)
	top := bs.Render("┌" + strings.Repeat("─", innerW) + "┐")
	mid := bs.Render("│") + ui.PadOrClipANSI(content, innerW) + bs.Render("│")
	bot := bs.Render("└" + strings.Repeat("─", innerW) + "┘")
	return top + "\n" + mid + "\n" + bot
}
