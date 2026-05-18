// Package issues renders the issues view: filter strip + table on the left,
// detail panel on the right.
package issues

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds the issues-view local state.
type Model struct {
	Sel      int
	Filter   string
	PendingG bool
}

// New returns a model with the "active" filter selected.
func New() Model { return Model{Filter: "active"} }

// Filtered returns the issues matching the current filter.
func (m Model) Filtered() []data.Issue {
	out := []data.Issue{}
	for _, it := range data.Issues {
		if m.Filter == "all" {
			out = append(out, it)
			continue
		}
		if m.Filter == "active" {
			if it.Status == "todo" || it.Status == "in_progress" || it.Status == "in_review" || it.Status == "blocked" {
				out = append(out, it)
			}
			continue
		}
		if it.Status == m.Filter {
			out = append(out, it)
		}
	}
	return out
}

// Update handles filter cycling, navigation, and "gg"/"G".
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	rows := m.Filtered()
	n := len(rows)
	switch k.String() {
	case "[":
		i := indexOf(data.IssueFilters, m.Filter)
		m.Filter = data.IssueFilters[(i-1+len(data.IssueFilters))%len(data.IssueFilters)]
		m.Sel = 0
	case "]":
		i := indexOf(data.IssueFilters, m.Filter)
		m.Filter = data.IssueFilters[(i+1)%len(data.IssueFilters)]
		m.Sel = 0
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
	}
	return m, nil
}

// View renders the issues view: filter strip + table on left, detail on right.
func (m Model) View(c ctx.Ctx, w, h int) string {
	gap := 1
	leftW := (w - gap) * 64 / 100
	rightW := w - gap - leftW

	filter := renderFilterStrip(m, leftW)
	rows := m.Filtered()
	tableH := h - 2
	if tableH < 6 {
		tableH = 6
	}
	tablePanel := ui.Panel(fmt.Sprintf("issues · %d", len(rows)), "",
		renderTable(m, rows, leftW-4, tableH-4), leftW, tableH, false, false)
	leftCol := filter + "\n" + tablePanel

	rightCol := renderDetail(m, rows, c, rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, ui.VGap(gap, h), rightCol)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"j k", "row"}, {"[ ]", "filter"}, {"↵", "open"}, {"n", "new"}, {"a", "assign"}, {"?", "help"},
	}
}

func indexOf(xs []string, v string) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return 0
}
