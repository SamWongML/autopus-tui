// Package workspaces renders the workspaces table on the left and a
// per-workspace detail panel on the right.
package workspaces

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/data"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds the current row selection.
type Model struct {
	Sel int
}

// New returns a zero-valued workspaces model (top row selected).
func New() Model { return Model{} }

// Update handles row navigation and the (w)atch toggle.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	n := len(data.Workspaces)
	switch k.String() {
	case "j", "down":
		if m.Sel < n-1 {
			m.Sel++
		}
	case "k", "up":
		if m.Sel > 0 {
			m.Sel--
		}
	case "w":
		data.Workspaces[m.Sel].Watch = !data.Workspaces[m.Sel].Watch
	}
	return m, nil
}

// View renders the workspaces view as left table + right detail.
func (m Model) View(c ctx.Ctx, w, h int) string {
	gap := 1
	leftW := (w - gap) * 62 / 100
	rightW := w - gap - leftW

	leftPanel := ui.Panel("workspaces", "w · watch toggle",
		renderTable(m, leftW-4, h-4), leftW, h, false, false)
	rightPanel := renderDetail(m, c, rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, ui.VGap(gap, h), rightPanel)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"j k", "row"}, {"w", "watch"}, {"m", "members"}, {"o", "browser"}, {"?", "help"},
	}
}
