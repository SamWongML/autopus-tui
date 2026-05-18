// Package help renders the full-screen Help overlay (replaces the body).
// State is intentionally trivial — the root model owns the open/closed flag.
package help

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model is a marker — the help overlay has no local state.
type Model struct{}

// New returns a zero model.
func New() Model { return Model{} }

// View renders the help overlay.
func (m Model) View(_ ctx.Ctx, w, h int) string {
	gap := 1
	usable := w - gap
	colW := usable / 2
	var leftCol, rightCol strings.Builder

	for i, g := range data.HelpGroups {
		b := &leftCol
		if i%2 == 1 {
			b = &rightCol
		}
		b.WriteString(theme.SAccent.Render(strings.ToUpper(g.Title)) + "\n")
		b.WriteString(ui.Dashed(colW-4) + "\n")
		for _, row := range g.Rows {
			line := theme.SText.Render(ui.PadRight(row[0], 20)) + " " + theme.SDim.Render(row[1])
			b.WriteString(ui.Truncate(line, colW-4) + "\n")
		}
		b.WriteString("\n")
	}

	bodyW := w - 4
	tips := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.Border)).
		Padding(0, 1).Width(bodyW-2).Render(
		theme.SAccent.Render("TIPS") + "\n" +
			theme.SDim.Render("• You can use the mouse — click rows, tabs, filter chips.") + "\n" +
			theme.SDim.Render("• Closing this TUI does ") + theme.SText.Render("not") + theme.SDim.Render(" stop the daemon — use ") + theme.SAccent.Render(":daemon stop") + theme.SDim.Render(".") + "\n" +
			theme.SDim.Render("• ") + theme.SAccent.Render("◆") + theme.SDim.Render(" always means: a session is waiting for you.") + "\n" +
			theme.SDim.Render("• Hold ") + theme.SAccent.Render("⇧") + theme.SDim.Render(" while moving to extend a multi-select."))

	leftH := strings.Count(leftCol.String(), "\n") + 1
	rightH := strings.Count(rightCol.String(), "\n") + 1
	colH := leftH
	if rightH > colH {
		colH = rightH
	}
	twoCol := lipgloss.JoinHorizontal(lipgloss.Top, leftCol.String(), ui.VGap(gap, colH), rightCol.String())
	body := twoCol + "\n" + tips
	return ui.Panel("help · keys & commands", "press ? again to close", body, w, h, false, true)
}

// KeyHints returns the help overlay's hint slice.
func (m Model) KeyHints() [][2]string {
	return [][2]string{{"?", "close"}, {"esc", "close"}}
}
