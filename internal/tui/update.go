package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/app"
	"autopus-tui/internal/chrome"
	"autopus-tui/internal/ui"
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
		m.Bp = ui.For(m.W)
		m.Sessions.LastW = m.W
		m.Issues.LastW = m.W
		m.Runtimes.LastW = m.W
		m.Workspaces.LastW = m.W
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
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

// handleMouse dispatches a mouse event. Wheel events are translated into
// synthetic up/down keys and dispatched through the existing routing. Click
// presses on row 0 (top), 1 (tabs), or H-1 (status) hit chrome; everything
// else is forwarded to the active overlay/attach/route with body-relative
// coords. The palette overlay receives box-relative coords so its handler
// only has to think in inner geometry.
func (m Model) handleMouse(mouse tea.MouseMsg) (tea.Model, tea.Cmd) {
	if mouse.Action == tea.MouseActionPress {
		switch mouse.Button {
		case tea.MouseButtonWheelUp:
			return m.handleKey(tea.KeyMsg{Type: tea.KeyUp})
		case tea.MouseButtonWheelDown:
			return m.handleKey(tea.KeyMsg{Type: tea.KeyDown})
		}
	}
	if mouse.Action != tea.MouseActionPress || mouse.Button != tea.MouseButtonLeft {
		return m, nil
	}
	if m.W == 0 || m.H == 0 {
		return m, nil
	}

	if mouse.Y == 0 {
		return m, nil
	}
	if mouse.Y == 1 {
		_, hits := chrome.TabBar(m.W, m.Route, m.Attach, m.Overlay == "help")
		id := ui.Hits(hits, mouse.X)
		if id == "" {
			return m, nil
		}
		if id == "help" {
			if m.Overlay == "help" {
				m.Overlay = ""
			} else {
				m.Overlay = "help"
			}
			return m, nil
		}
		m.Route = id
		m.Attach = ""
		m.Overlay = ""
		return m, nil
	}
	if mouse.Y >= m.H-1 {
		return m, nil
	}

	bodyMouse := mouse
	bodyMouse.Y = mouse.Y - 2

	if m.Overlay == "palette" {
		bodyH := m.H - 3
		items := m.Palette.Filtered()
		const maxItems = 12
		if len(items) > maxItems {
			items = items[:maxItems]
		}
		boxH := 6 + len(items)
		const boxW = 70
		topY := (bodyH - boxH) / 2
		leftX := (m.W - boxW) / 2
		if topY < 0 {
			topY = 0
		}
		if leftX < 0 {
			leftX = 0
		}
		boxMouse := bodyMouse
		boxMouse.X = bodyMouse.X - leftX
		boxMouse.Y = bodyMouse.Y - topY
		var cmd tea.Cmd
		m.Palette, cmd = m.Palette.Update(boxMouse)
		return m, cmd
	}
	if m.Overlay == "help" || m.Overlay == "onboarding" {
		return m, nil
	}

	if m.Attach != "" {
		return m, nil
	}

	var cmd tea.Cmd
	switch m.Route {
	case "sessions":
		m.Sessions, cmd = m.Sessions.Update(bodyMouse)
	case "issues":
		m.Issues, cmd = m.Issues.Update(bodyMouse)
	case "runtimes":
		m.Runtimes, cmd = m.Runtimes.Update(bodyMouse)
	case "workspaces":
		m.Workspaces, cmd = m.Workspaces.Update(bodyMouse)
	case "logs":
		m.Logs, cmd = m.Logs.Update(bodyMouse)
	}
	return m, cmd
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
