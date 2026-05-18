// Package onboarding renders the first-run wizard: a step rail on the left
// and a per-step body panel on the right.
package onboarding

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds onboarding local state.
type Model struct {
	Step    int
	Server  string
	Auth    bool
	Watched map[string]bool
}

// New returns a fresh wizard with the cloud server option pre-selected and the
// three "main" workspaces toggled on.
func New() Model {
	return Model{
		Server:  "cloud",
		Watched: map[string]bool{"ws_core": true, "ws_platform": true, "ws_docs": true},
	}
}

// Update handles step navigation and dismissal.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "esc", "q":
		return m, func() tea.Msg { return app.NavigateMsg{Overlay: ""} }
	case "left", "h":
		if m.Step > 0 {
			m.Step--
		}
	case "right", "l", "enter":
		if m.Step < len(data.OnbSteps)-1 {
			m.Step++
		} else {
			return m, func() tea.Msg { return app.NavigateMsg{Overlay: ""} }
		}
	}
	return m, nil
}

// View renders the wizard as left rail + right body.
func (m Model) View(_ ctx.Ctx, w, h int) string {
	gap := 1
	railW := 28
	if railW > w/3 {
		railW = w / 3
	}
	bodyW := w - railW - gap

	rail := renderRail(m, railW, h)
	body := renderStep(m, bodyW, h)
	return lipgloss.JoinHorizontal(lipgloss.Top, rail, ui.VGap(gap, h), body)
}

// KeyHints returns the status-bar key hint slice for this overlay.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"← →", "step"}, {"space", "toggle"}, {"↵", "next"}, {"esc", "skip"},
	}
}
