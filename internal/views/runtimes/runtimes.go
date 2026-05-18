// Package runtimes renders the agent-runtimes view: a table of detected CLIs
// and a detail panel for the selected one.
package runtimes

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/data"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds the runtimes-view local state.
type Model struct {
	Sel      int
	PendingG bool // "gg" → top
}

// New returns a zero-value model (first runtime selected).
func New() Model { return Model{} }

// Update handles row navigation and "gg"/"G" jump-to-top/bottom.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	n := len(data.Runtimes)
	switch k.String() {
	case "j", "down":
		if m.Sel < n-1 {
			m.Sel++
		}
		m.PendingG = false
	case "k", "up":
		if m.Sel > 0 {
			m.Sel--
		}
		m.PendingG = false
	case "g":
		if m.PendingG {
			m.Sel = 0
			m.PendingG = false
		} else {
			m.PendingG = true
		}
	case "G":
		m.Sel = n - 1
		m.PendingG = false
	}
	return m, nil
}

// View renders the runtimes view: table on left, detail on right.
func (m Model) View(_ ctx.Ctx, w, h int) string {
	gap := 1
	leftW := (w - gap) * 64 / 100
	rightW := w - gap - leftW

	leftPanel := ui.Panel("agent runtimes", "auto-detected on $PATH",
		renderTable(m, leftW-4, h-4), leftW, h, false, false)
	rightPanel := renderDetail(m, rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, ui.VGap(gap, h), rightPanel)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"j k", "row"}, {"↵", "config"}, {"t", "test"}, {"d", "disable"}, {"↺", "rescan"}, {"?", "help"},
	}
}
