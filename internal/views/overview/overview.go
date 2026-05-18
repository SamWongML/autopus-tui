// Package overview renders the home dashboard: 5 stat cards in a 2×3-style
// grid (activity spans the right column full-height). h/j/k/l moves focus;
// ↵ jumps to the focused card's underlying route via app.NavigateMsg.
package overview

import (
	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/app"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Card describes one stat card's position in the grid and its jump-to route.
type Card struct {
	ID, Jump, Title string
	Row, Col        int
}

// Cards is the static layout: 2 columns of stacked stat panels, plus a
// full-height activity column.
var Cards = []Card{
	{"daemon", "config", "daemon", 0, 0},
	{"env", "config", "local environment", 1, 0},
	{"server", "config", "remote · app.autopus.ai", 0, 1},
	{"sessions", "sessions", "sessions", 1, 1},
	{"activity", "logs", "recent activity", 0, 2},
}

// Model holds overview-view local state.
type Model struct {
	Focus int // index into Cards
}

// New returns a model with the daemon card focused.
func New() Model { return Model{} }

// Update handles arrow/hjkl focus moves and ↵ jump.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	s := k.String()
	cur := Cards[m.Focus]
	move := func(dr, dc int) {
		for i, c := range Cards {
			if c.Row == cur.Row+dr && c.Col == cur.Col+dc {
				m.Focus = i
				return
			}
		}
	}
	switch s {
	case "down", "j":
		move(1, 0)
	case "up", "k":
		move(-1, 0)
	case "right", "l":
		move(0, 1)
	case "left", "h":
		move(0, -1)
	case "enter":
		jump := cur.Jump
		return m, func() tea.Msg { return app.NavigateMsg{To: jump} }
	}
	return m, nil
}

// View renders the overview grid.
func (m Model) View(c ctx.Ctx, w, h int) string {
	gap := 1
	usable := w - 2*gap
	if usable < 30 {
		usable = 30
	}
	colA := usable * 35 / 100
	colB := usable * 35 / 100
	colC := usable - colA - colB

	rowH := (h - 1) / 2
	if rowH < 6 {
		rowH = 6
	}
	row2H := h - rowH - 1
	if row2H < 4 {
		row2H = 4
	}

	pDaemon := overviewPanel(m, c, "daemon", colA, rowH)
	pEnv := overviewPanel(m, c, "env", colA, row2H)
	pServer := overviewPanel(m, c, "server", colB, rowH)
	pSessions := overviewPanel(m, c, "sessions", colB, row2H)
	pActivity := overviewPanel(m, c, "activity", colC, h)

	colALines := ui.JoinVertical(pDaemon, ui.Blank(colA), pEnv)
	colBLines := ui.JoinVertical(pServer, ui.Blank(colB), pSessions)

	gapCol := ui.VGap(gap, h)
	return ui.JoinHorizontal(colALines, gapCol, colBLines, gapCol, pActivity)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"h j k l", "focus"}, {"↵", "drill"}, {"?", "help"}, {":", "cmd"},
	}
}

func cardIdx(id string) int {
	for i, c := range Cards {
		if c.ID == id {
			return i
		}
	}
	return 0
}
