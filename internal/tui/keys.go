package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/app"
)

// handleKey routes a key event to the active overlay/view, applying global
// shortcuts first.
func (m Model) handleKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	s := k.String()

	// Overlays capture all keys first. Palette + onboarding manage their own
	// input fields; help is a pure consumer (we close it inline).
	if m.Overlay == "palette" {
		var cmd tea.Cmd
		m.Palette, cmd = m.Palette.Update(k)
		return m, cmd
	}
	if m.Overlay == "onboarding" {
		var cmd tea.Cmd
		m.Onboarding, cmd = m.Onboarding.Update(k)
		return m, cmd
	}

	// sessions' / search mode: don't let global keys eat input. The sessions
	// model owns the toggle; ask it directly when searching.
	if m.Route == "sessions" && m.Sessions.Searching {
		var cmd tea.Cmd
		m.Sessions, cmd = m.Sessions.Update(k)
		return m, cmd
	}

	// Global keys.
	switch s {
	case "ctrl+c", "q":
		if m.Overlay == "help" || m.Attach != "" {
			break
		}
		m.quitting = true
		return m, tea.Quit
	case "?":
		if m.Overlay == "help" {
			m.Overlay = ""
		} else {
			m.Overlay = "help"
		}
		return m, nil
	case ":":
		m.Overlay = "palette"
		m.Palette.Reset()
		return m, nil
	case "esc":
		if m.Overlay == "help" {
			m.Overlay = ""
			return m, nil
		}
		if m.Attach != "" {
			m.Attach = ""
			m.AttachM.ID = ""
			return m, nil
		}
	case "tab":
		idx := app.RouteIdx(m.Route)
		nx := (idx + 1) % len(app.Routes)
		m.Route = app.Routes[nx].ID
		m.Attach = ""
		return m, nil
	case "shift+tab":
		idx := app.RouteIdx(m.Route)
		nx := (idx - 1 + len(app.Routes)) % len(app.Routes)
		m.Route = app.Routes[nx].ID
		m.Attach = ""
		return m, nil
	}

	// Number keys 1–7 jump to a route.
	if len(s) == 1 && s >= "1" && s <= "7" {
		for _, r := range app.Routes {
			if r.Key == s {
				m.Route = r.ID
				m.Attach = ""
				return m, nil
			}
		}
	}

	if m.Overlay == "help" {
		return m, nil
	}

	if m.Attach != "" {
		var cmd tea.Cmd
		m.AttachM, cmd = m.AttachM.Update(k)
		return m, cmd
	}

	var cmd tea.Cmd
	switch m.Route {
	case "overview":
		m.Overview, cmd = m.Overview.Update(k)
	case "sessions":
		m.Sessions, cmd = m.Sessions.Update(k)
	case "issues":
		m.Issues, cmd = m.Issues.Update(k)
	case "runtimes":
		m.Runtimes, cmd = m.Runtimes.Update(k)
	case "workspaces":
		m.Workspaces, cmd = m.Workspaces.Update(k)
	case "logs":
		m.Logs, cmd = m.Logs.Update(k)
	case "config":
		m.Config, cmd = m.Config.Update(k)
	}
	return m, cmd
}
