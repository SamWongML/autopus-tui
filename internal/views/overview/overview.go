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

// View renders the overview grid. Layout adapts to c.Bp:
//   - BPxs:  1 column stack — daemon, server, sessions, env, activity.
//   - BPsm:  2 columns — left stacks daemon/env/sessions; right stacks server/activity.
//   - BPmd+: 3 columns — colA/colB at 35% each, activity gets the rest.
func (m Model) View(c ctx.Ctx, w, h int) string {
	gap := 1
	switch c.Bp {
	case ui.BPxs:
		return m.viewSingleCol(c, w, h)
	case ui.BPsm:
		return m.viewTwoCol(c, w, h, gap)
	default:
		return m.viewThreeCol(c, w, h, gap)
	}
}

func (m Model) viewThreeCol(c ctx.Ctx, w, h, gap int) string {
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

func (m Model) viewTwoCol(c ctx.Ctx, w, h, gap int) string {
	usable := w - gap
	colA := usable / 2
	colB := usable - colA

	// Left column: daemon / env / sessions stacked.
	leftThirds := splitH(h, 3)
	pDaemon := overviewPanel(m, c, "daemon", colA, leftThirds[0])
	pEnv := overviewPanel(m, c, "env", colA, leftThirds[1])
	pSessions := overviewPanel(m, c, "sessions", colA, leftThirds[2])
	colALines := ui.JoinVertical(pDaemon, pEnv, pSessions)

	// Right column: server / activity.
	rightServerH := ui.Min(h/2, 14)
	if rightServerH < 8 {
		rightServerH = 8
	}
	rightActivityH := h - rightServerH
	if rightActivityH < 6 {
		rightActivityH = 6
	}
	pServer := overviewPanel(m, c, "server", colB, rightServerH)
	pActivity := overviewPanel(m, c, "activity", colB, rightActivityH)
	colBLines := ui.JoinVertical(pServer, pActivity)

	gapCol := ui.VGap(gap, h)
	return ui.JoinHorizontal(colALines, gapCol, colBLines)
}

func (m Model) viewSingleCol(c ctx.Ctx, w, h int) string {
	order := []string{"daemon", "server", "sessions", "env", "activity"}
	base := h / len(order)
	if base < 6 {
		base = 6
	}
	leftover := h - base*(len(order)-1)
	if leftover < base {
		leftover = base
	}
	parts := make([]string, 0, len(order))
	for i, id := range order {
		ph := base
		if i == len(order)-1 {
			ph = leftover
		}
		parts = append(parts, overviewPanel(m, c, id, w, ph))
	}
	return ui.JoinVertical(parts...)
}

func splitH(h, n int) []int {
	out := make([]int, n)
	base := h / n
	for i := 0; i < n; i++ {
		out[i] = base
	}
	out[n-1] += h - base*n
	return out
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
