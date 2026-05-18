package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/app"
)

// Init kicks off the wall-clock and spinner tickers.
func (m Model) Init() tea.Cmd {
	return tea.Batch(clockCmd(), spinCmd())
}

// Update dispatches messages: window resize, clock/spin ticks, NavigateMsg
// from children, and keys (which it routes to the active view).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.W, m.H = msg.Width, msg.Height
		return m, nil
	case app.ClockMsg:
		m.Now = time.Time(msg)
		return m, clockCmd()
	case app.SpinMsg:
		_ = msg
		m.Spin = (m.Spin + 1) % 10
		return m, spinCmd()
	case app.NavigateMsg:
		return m.handleNavigate(msg), nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

// handleNavigate applies a NavigateMsg sent by a child view.
func (m Model) handleNavigate(n app.NavigateMsg) Model {
	if n.To != "" {
		m.Route = n.To
		m.Attach = ""
		return m
	}
	// Empty Attach explicitly clears the attachment; non-empty sets it.
	if n.Attach != m.Attach {
		m.Attach = n.Attach
		m.AttachM.ID = n.Attach
		m.AttachM.Scroll = 0
		m.AttachM.Reply = ""
	}
	if n.Overlay != m.Overlay {
		m.Overlay = n.Overlay
		if n.Overlay == "palette" {
			m.Palette.Reset()
		}
	}
	return m
}
