// Package palette renders the floating ":" command palette. Pressing ↵ on
// "replay onboarding" emits an app.NavigateMsg{Overlay:"onboarding"}; other
// items just close the palette.
package palette

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds palette local state.
type Model struct {
	Query string
	Sel   int
}

// New returns an empty palette.
func New() Model { return Model{} }

// Reset clears the query and selection.
func (m *Model) Reset() {
	m.Query = ""
	m.Sel = 0
}

// Filtered returns the palette entries matching the current query.
func (m Model) Filtered() []data.PaletteItem {
	q := strings.ToLower(m.Query)
	if q == "" {
		return data.PaletteItems
	}
	out := []data.PaletteItem{}
	for _, it := range data.PaletteItems {
		if strings.Contains(strings.ToLower(it.K+" "+it.Label), q) {
			out = append(out, it)
		}
	}
	return out
}

// Update handles palette-local keys. Returns model + optional cmd. Returning
// (m, nil) for "enter" with no nav-emitting command leaves the root model to
// close the overlay.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	s := k.String()
	items := m.Filtered()
	switch s {
	case "esc":
		return m, func() tea.Msg { return app.NavigateMsg{Overlay: ""} }
	case "enter":
		if m.Sel >= 0 && m.Sel < len(items) {
			it := items[m.Sel]
			if it.K == "replay onboarding" {
				return m, func() tea.Msg { return app.NavigateMsg{Overlay: "onboarding"} }
			}
		}
		return m, func() tea.Msg { return app.NavigateMsg{Overlay: ""} }
	case "up":
		if m.Sel > 0 {
			m.Sel--
		}
		return m, nil
	case "down":
		if m.Sel < len(items)-1 {
			m.Sel++
		}
		return m, nil
	case "backspace":
		if len(m.Query) > 0 {
			m.Query = m.Query[:len(m.Query)-1]
			m.Sel = 0
		}
		return m, nil
	}
	if len(s) == 1 {
		m.Query += s
		m.Sel = 0
	}
	return m, nil
}

// View renders the floating palette box. The root model overlays it centered.
func (m Model) View(_ ctx.Ctx) string {
	w := 70
	items := m.Filtered()
	maxItems := 12
	if len(items) > maxItems {
		items = items[:maxItems]
	}

	header := theme.SAccent.Render(":") + " " +
		theme.SText.Render(m.Query) +
		theme.SAccent.Render("▌") +
		strings.Repeat(" ", ui.Max(1, w-6-lipgloss.Width(m.Query)-15)) +
		theme.SFaint.Render("esc to close")

	var rowsB strings.Builder
	for i, it := range items {
		active := i == m.Sel
		marker := "  "
		if active {
			marker = theme.SAccent.Render("▎ ")
		}
		kindCol := theme.SFaint.Render(ui.PadRight(strings.ToUpper(it.Kind), 8))
		cmd := theme.SAccent.Render(":") + theme.SText.Render(ui.PadRight(it.K, 30))
		desc := theme.SFaint.Render(ui.Truncate(it.Label, w-50))
		row := marker + kindCol + " " + cmd + " " + desc
		if active {
			row = theme.WithBg(row, theme.AccentFaint)
		}
		rowsB.WriteString(row + "\n")
	}
	if len(items) == 0 {
		rowsB.WriteString(theme.SFaint.Render("no command matches \"" + m.Query + "\""))
	}

	footer := theme.SFaint.Render("↑↓ select   ↵ run   esc close")

	inner := header + "\n" + theme.SBorder.Render(strings.Repeat("─", w-2)) + "\n" +
		strings.TrimRight(rowsB.String(), "\n") + "\n" +
		theme.SBorder.Render(strings.Repeat("─", w-2)) + "\n" + footer

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.AccentDim)).
		Background(lipgloss.Color(theme.Bg)).
		Padding(0, 1).Width(w).Render(inner)
}
