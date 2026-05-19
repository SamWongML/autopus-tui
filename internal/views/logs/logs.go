// Package logs renders the daemon.log tail view: a level-filter strip and a
// scrolling body of structured log lines.
package logs

import (
	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/data"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds the logs view's local UI state.
type Model struct {
	Follow bool
	Level  string
}

// New returns a fresh model with follow=true and level="all".
func New() Model {
	return Model{Follow: true, Level: "all"}
}

// Update handles keys local to the logs view.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if mouse, ok := msg.(tea.MouseMsg); ok {
		return m.handleMouse(mouse)
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "f":
		m.Follow = !m.Follow
	case "[":
		i := indexOf(data.LogLevelFilters, m.Level)
		m.Level = data.LogLevelFilters[(i-1+len(data.LogLevelFilters))%len(data.LogLevelFilters)]
	case "]":
		i := indexOf(data.LogLevelFilters, m.Level)
		m.Level = data.LogLevelFilters[(i+1)%len(data.LogLevelFilters)]
	}
	return m, nil
}

// View renders the logs view as a strip + bordered body.
func (m Model) View(c ctx.Ctx, w, h int) string {
	strip := renderFilterStrip(m, c.Spin, w)
	rows := filtered(m)
	panelH := h - 2
	if panelH < 6 {
		panelH = 6
	}
	body := renderBody(rows, m.Follow, w-4, panelH-4)
	p := ui.Panel("daemon.log", "~/.autopus/daemon.log", body, w, panelH, false, false)
	return strip + "\n" + p
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"f", "follow"}, {"[ ]", "level"}, {"c", "clear"}, {"s", "save"}, {"?", "help"},
	}
}

// handleMouse maps a body-relative click on the filter strip (Y=0) to a level
// chip or the follow toggle. Rows below the strip aren't selectable.
func (m Model) handleMouse(mouse tea.MouseMsg) (Model, tea.Cmd) {
	if mouse.Action != tea.MouseActionPress || mouse.Button != tea.MouseButtonLeft {
		return m, nil
	}
	if mouse.Y != 0 {
		return m, nil
	}
	hits := filterHitBoxes()
	id := ui.Hits(hits, mouse.X)
	if id == "" {
		return m, nil
	}
	if id == "__follow" {
		m.Follow = !m.Follow
		return m, nil
	}
	m.Level = id
	return m, nil
}

func filtered(m Model) []data.LogLine {
	if m.Level == "all" {
		return data.LogLines
	}
	out := []data.LogLine{}
	for _, l := range data.LogLines {
		if l.Level == m.Level {
			out = append(out, l)
		}
	}
	return out
}

func indexOf(xs []string, v string) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return 0
}
