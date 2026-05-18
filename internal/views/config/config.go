// Package config renders the config view: profiles + budgets on the left,
// server + daemon in the middle, per-agent overrides on the right.
package config

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds config-view local state.
type Model struct {
	Profile  string
	AgentSel int
}

// New returns a model with the default profile selected.
func New() Model { return Model{Profile: "default"} }

// Update handles agent-chip navigation.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "right", "l":
		if m.AgentSel < len(data.AgentCfgs)-1 {
			m.AgentSel++
		}
	case "left", "h":
		if m.AgentSel > 0 {
			m.AgentSel--
		}
	}
	return m, nil
}

// View renders the config view.
func (m Model) View(_ ctx.Ctx, w, h int) string {
	gap := 1
	col1 := (w - 2*gap) * 28 / 100
	col2 := (w - 2*gap) * 35 / 100
	col3 := w - 2*gap - col1 - col2

	leftCol := renderProfilesCol(m, col1, h)
	midCol := renderServerDaemonCol(m, col2, h)
	rightCol := renderAgentsPanel(m, col3, h)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftCol, ui.VGap(gap, h),
		midCol, ui.VGap(gap, h),
		rightCol)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"j k", "row"}, {"e", "edit"}, {"↵", "reload"}, {"?", "help"},
	}
}
