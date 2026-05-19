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
	LastW    int  // last body width — used by mouse handler to derive layout
}

// New returns a zero-value model (first runtime selected).
func New() Model { return Model{} }

// Update handles row navigation and "gg"/"G" jump-to-top/bottom.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if mouse, ok := msg.(tea.MouseMsg); ok {
		return m.handleMouse(mouse)
	}
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

// View renders the runtimes view. At BPxs/BPsm the detail stacks under the
// table; at BPmd+ it sits on the right.
func (m Model) View(_ ctx.Ctx, w, h int) string {
	gap := 1
	stacked := ui.For(w) <= ui.BPsm

	if stacked {
		detailH := ui.Max(10, h/3)
		topH := h - detailH
		if topH < 8 {
			topH = 8
			detailH = ui.Max(4, h-topH)
		}
		leftPanel := ui.Panel("agent runtimes", "auto-detected on $PATH",
			renderTable(m, w-4, topH-4), w, topH, false, false)
		rightPanel := renderDetail(m, w, detailH)
		return lipgloss.JoinVertical(lipgloss.Left, leftPanel, rightPanel)
	}

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

// handleMouse maps a body-relative click to a runtime row. Layout: panel only,
// data rows start at body Y=5 (top border + title + divider + table header +
// separator). Side-by-side at BPmd+: left column is 64% of (w-1).
func (m Model) handleMouse(mouse tea.MouseMsg) (Model, tea.Cmd) {
	if mouse.Action != tea.MouseActionPress || mouse.Button != tea.MouseButtonLeft {
		return m, nil
	}
	w := m.LastW
	if w <= 0 {
		return m, nil
	}
	leftW := w
	if ui.For(w) > ui.BPsm {
		leftW = (w - 1) * 64 / 100
	}
	if mouse.X >= leftW {
		return m, nil
	}
	rowIdx := mouse.Y - 5
	if rowIdx < 0 || rowIdx >= len(data.Runtimes) {
		return m, nil
	}
	m.Sel = rowIdx
	return m, nil
}
