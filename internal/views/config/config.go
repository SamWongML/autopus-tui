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

// View renders the config view. Layout adapts to c.Bp:
//   - BPxs:  single column stack (profiles, server/daemon, agents).
//   - BPsm:  2 columns — profiles on the left, server/daemon stacked over agents on the right.
//   - BPmd+: 3 columns at 30/35/35.
func (m Model) View(c ctx.Ctx, w, h int) string {
	gap := 1
	switch c.Bp {
	case ui.BPxs:
		return m.viewSingleCol(w, h)
	case ui.BPsm:
		return m.viewTwoCol(w, h, gap)
	default:
		return m.viewThreeCol(w, h, gap)
	}
}

func (m Model) viewThreeCol(w, h, gap int) string {
	col1 := (w - 2*gap) * 30 / 100
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

func (m Model) viewTwoCol(w, h, gap int) string {
	colL := (w - gap) * 42 / 100
	colR := w - gap - colL

	rightTopH := h / 2
	if rightTopH < 10 {
		rightTopH = 10
	}
	rightBotH := h - rightTopH
	if rightBotH < 10 {
		rightBotH = 10
	}

	left := renderProfilesCol(m, colL, h)
	rightTop := renderServerDaemonCol(m, colR, rightTopH)
	rightBot := renderAgentsPanel(m, colR, rightBotH)
	right := ui.JoinVertical(rightTop, rightBot)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, ui.VGap(gap, h), right)
}

func (m Model) viewSingleCol(w, h int) string {
	// Profiles+budgets and server+daemon each pack two panels, so they need
	// more room than the agents panel.
	topH := ui.Max(18, h*36/100)
	midH := ui.Max(18, h*36/100)
	botH := h - topH - midH
	if botH < 10 {
		botH = 10
	}
	top := renderProfilesCol(m, w, topH)
	mid := renderServerDaemonCol(m, w, midH)
	bot := renderAgentsPanel(m, w, botH)
	return ui.JoinVertical(top, mid, bot)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"j k", "row"}, {"e", "edit"}, {"↵", "reload"}, {"?", "help"},
	}
}
