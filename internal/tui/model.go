// Package tui owns the root bubbletea Model: it composes one child model per
// view, dispatches keys + navigation messages, and renders the final frame.
//
// Imports flow: tui → app, chrome, ui, theme, data, views/*.
// Views import app (for NavigateMsg etc.) but never import tui, so there is no
// cycle. Adding a new view = add it to app.Routes + add a field here +
// register it in the four switches in this file.
package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/attach"
	"autopus-tui/internal/views/config"
	"autopus-tui/internal/views/ctx"
	"autopus-tui/internal/views/help"
	"autopus-tui/internal/views/issues"
	"autopus-tui/internal/views/logs"
	"autopus-tui/internal/views/onboarding"
	"autopus-tui/internal/views/overview"
	"autopus-tui/internal/views/palette"
	"autopus-tui/internal/views/runtimes"
	"autopus-tui/internal/views/sessions"
	"autopus-tui/internal/views/workspaces"
)

// Model is the application's root state. Each field is either chrome state
// (w/h/now/spin) or a child view's model. Overlay is "" | "help" | "palette" |
// "onboarding". When Attach is non-empty it overrides the active route.
type Model struct {
	W, H int
	Bp   ui.BP
	Now  time.Time
	Spin int

	Route   string
	Attach  string // session ID; empty = no attach
	Overlay string // "", "help", "palette", "onboarding"

	Overview   overview.Model
	Sessions   sessions.Model
	Issues     issues.Model
	Runtimes   runtimes.Model
	Workspaces workspaces.Model
	Logs       logs.Model
	Config     config.Model
	AttachM    attach.Model
	Help       help.Model
	Palette    palette.Model
	Onboarding onboarding.Model

	quitting bool
}

// New returns a freshly initialized root model with all child views at their
// defaults and the overview route selected.
func New() Model {
	return Model{
		Route:      "overview",
		Now:        time.Now(),
		Overview:   overview.New(),
		Sessions:   sessions.New(),
		Issues:     issues.New(),
		Runtimes:   runtimes.New(),
		Workspaces: workspaces.New(),
		Logs:       logs.New(),
		Config:     config.New(),
		AttachM:    attach.New(),
		Help:       help.New(),
		Palette:    palette.New(),
		Onboarding: onboarding.New(),
	}
}

// frameCtx returns the per-frame context every view's View receives.
func (m Model) frameCtx() ctx.Ctx {
	return ctx.Ctx{Now: m.Now, Spin: m.Spin}
}

// activeKeyHints returns the status-bar hints for whatever view currently
// owns the keys. Overlay > attach > route.
func (m Model) activeKeyHints() [][2]string {
	switch m.Overlay {
	case "help":
		return m.Help.KeyHints()
	case "onboarding":
		return m.Onboarding.KeyHints()
	}
	if m.Attach != "" {
		return m.AttachM.KeyHints()
	}
	switch m.Route {
	case "overview":
		return m.Overview.KeyHints()
	case "sessions":
		return m.Sessions.KeyHints()
	case "issues":
		return m.Issues.KeyHints()
	case "runtimes":
		return m.Runtimes.KeyHints()
	case "workspaces":
		return m.Workspaces.KeyHints()
	case "logs":
		return m.Logs.KeyHints()
	case "config":
		return m.Config.KeyHints()
	}
	return nil
}

// clockCmd schedules the next 1s wall-clock tick.
func clockCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return app.ClockMsg(t) })
}

// spinCmd schedules the next ~90ms spinner frame advance.
func spinCmd() tea.Cmd {
	return tea.Tick(90*time.Millisecond, func(t time.Time) tea.Msg { return app.SpinMsg(t) })
}

// _ keeps a reference to data so go's "imported and not used" check stays
// happy if a future refactor temporarily drops the explicit use. Cheap and
// honest — data is the canonical source of truth.
var _ = data.Sessions
var _ = theme.Accent
