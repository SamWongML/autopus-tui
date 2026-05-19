// Package sessions renders the sessions list: filter strip, table, and right-
// side peek panel. Pressing ↵ emits an app.NavigateMsg{Attach: selectedID}.
package sessions

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds sessions-view local state.
type Model struct {
	Sel       int
	Filter    string
	Query     string
	Searching bool
	PendingG  bool
	LastW     int // last body width — used by mouse handler to derive layout
}

// New returns a fresh model with Filter="all".
func New() Model { return Model{Filter: "all"} }

// Filtered returns the sessions matching the filter+query, sorted with
// needs_input first (per data.States Sort).
func (m Model) Filtered() []data.Session {
	out := []data.Session{}
	q := strings.ToLower(m.Query)
	for _, s := range data.Sessions {
		if m.Filter != "all" && s.State != m.Filter {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(s.Title+s.Issue+s.ID), q) {
			continue
		}
		out = append(out, s)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return theme.State(out[i].State).Sort < theme.State(out[j].State).Sort
	})
	return out
}

// Update handles search input, filter cycling, navigation, and ↵ (attach).
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if mouse, ok := msg.(tea.MouseMsg); ok {
		return m.handleMouse(mouse)
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	s := k.String()

	if m.Searching {
		switch s {
		case "esc":
			m.Searching = false
			m.Query = ""
			return m, nil
		case "enter":
			m.Searching = false
			return m, nil
		case "backspace":
			if len(m.Query) > 0 {
				m.Query = m.Query[:len(m.Query)-1]
			}
			return m, nil
		}
		if len(s) == 1 {
			m.Query += s
		}
		return m, nil
	}

	rows := m.Filtered()
	n := len(rows)
	switch s {
	case "/":
		m.Searching = true
		return m, nil
	case "[":
		i := indexOf(data.SessionFilters, m.Filter)
		m.Filter = data.SessionFilters[(i-1+len(data.SessionFilters))%len(data.SessionFilters)]
		m.Sel = 0
	case "]":
		i := indexOf(data.SessionFilters, m.Filter)
		m.Filter = data.SessionFilters[(i+1)%len(data.SessionFilters)]
		m.Sel = 0
	case "*":
		for i := m.Sel + 1; i < n; i++ {
			if rows[i].State == "needs_input" {
				m.Sel = i
				return m, nil
			}
		}
		for i := 0; i < n; i++ {
			if rows[i].State == "needs_input" {
				m.Sel = i
				return m, nil
			}
		}
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
		if n > 0 {
			m.Sel = n - 1
		}
		m.PendingG = false
	case "ctrl+d":
		m.Sel = ui.Min(n-1, m.Sel+ui.Max(1, n/8))
	case "ctrl+u":
		m.Sel = ui.Max(0, m.Sel-ui.Max(1, n/8))
	case "enter", "right":
		if m.Sel >= 0 && m.Sel < n {
			id := rows[m.Sel].ID
			return m, func() tea.Msg { return app.NavigateMsg{Attach: id} }
		}
	case "esc":
		if m.Query != "" {
			m.Query = ""
		} else if m.Filter != "all" {
			m.Filter = "all"
		}
	}
	return m, nil
}

// View renders the sessions view. At BPxs/BPsm the peek stacks under the
// table; at BPmd+ it sits to the right.
func (m Model) View(c ctx.Ctx, w, h int) string {
	gap := 1
	rows := m.Filtered()
	stacked := ui.For(w) <= ui.BPsm

	if stacked {
		peekH := ui.Max(8, h/3)
		topH := h - peekH
		if topH < 8 {
			topH = 8
			peekH = ui.Max(4, h-topH)
		}
		filterLine := renderFilterStrip(m, c.Spin, w)
		tableH := topH - 2
		if tableH < 6 {
			tableH = 6
		}
		tablePanel := ui.Panel(fmt.Sprintf("sessions · %d", len(rows)), "sort · needs-you first",
			renderTable(m, c, rows, w-4, tableH-4), w, tableH, false, false)
		topCol := filterLine + "\n" + tablePanel
		peek := renderPeek(m, c, rows, w, peekH)
		return lipgloss.JoinVertical(lipgloss.Left, topCol, peek)
	}

	usable := w - gap
	leftW := usable * 62 / 100
	rightW := usable - leftW

	filterLine := renderFilterStrip(m, c.Spin, leftW)
	tableH := h - 2
	if tableH < 6 {
		tableH = 6
	}
	tablePanel := ui.Panel(fmt.Sprintf("sessions · %d", len(rows)), "sort · needs-you first",
		renderTable(m, c, rows, leftW-4, tableH-4), leftW, tableH, false, false)
	leftCol := filterLine + "\n" + tablePanel

	rightCol := renderPeek(m, c, rows, rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, ui.VGap(gap, h), rightCol)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"j k", "row"}, {"[ ]", "filter"}, {"*", "needs-input"}, {"↵", "attach"}, {"/", "search"}, {"?", "help"},
	}
}

// handleMouse maps a body-relative click into a filter change, a row select +
// attach, or a no-op. Row 0 = filter strip; rows 6+ in the left column = table
// data rows (panel top + title + divider + table header + separator = 5 rows).
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
		leftW = (w - 1) * 62 / 100
	}
	if mouse.X >= leftW {
		return m, nil
	}
	if mouse.Y == 0 {
		hits := filterHitBoxes()
		if id := ui.Hits(hits, mouse.X); id != "" {
			m.Filter = id
			m.Sel = 0
		}
		return m, nil
	}
	rowIdx := mouse.Y - 6
	rows := m.Filtered()
	if rowIdx < 0 || rowIdx >= len(rows) {
		return m, nil
	}
	m.Sel = rowIdx
	id := rows[rowIdx].ID
	return m, func() tea.Msg { return app.NavigateMsg{Attach: id} }
}

func indexOf(xs []string, v string) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return 0
}
