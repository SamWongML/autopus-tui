// Package attach renders the per-session "attach" view: transcript on the
// left, run/budget/actions stack on the right, reply box at the bottom. esc
// emits an app.NavigateMsg{Attach:""} to clear the attachment.
package attach

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

// Model holds attach-view local state. The session being attached is owned by
// the root model (passed in via ID); we only track scroll + reply buffer here.
type Model struct {
	ID     string
	Reply  string
	Scroll int
}

// New returns a fresh model.
func New() Model { return Model{} }

// Update handles attach-local keys.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "esc", "left":
		m.ID = ""
		m.Reply = ""
		return m, func() tea.Msg { return app.NavigateMsg{Attach: ""} }
	case "j", "down":
		m.Scroll++
	case "k", "up":
		if m.Scroll > 0 {
			m.Scroll--
		}
	case "g":
		m.Scroll = 0
	case "G":
		m.Scroll = 9999
	case "backspace":
		if len(m.Reply) > 0 {
			m.Reply = m.Reply[:len(m.Reply)-1]
		}
	}
	return m, nil
}

// View renders the attach view.
func (m Model) View(c ctx.Ctx, w, h int) string {
	s := data.FindSession(m.ID)
	gap := 1
	leftW := (w - gap) * 70 / 100
	rightW := w - gap - leftW

	pill := ui.StatePill(s.State, c.Spin)
	pillLines := strings.Split(pill, "\n")
	pillLine := pill
	if len(pillLines) >= 2 {
		pillLine = pillLines[1]
	}

	headerLeft := ui.KeyChip("esc", "detach", true) + "  " +
		pillLine + " " + theme.SFaint.Render(s.Issue) + " " + theme.SText.Render("· "+ui.Truncate(s.Title, leftW-40))
	headerRight := theme.SFaint.Render(fmt.Sprintf("%s · elapsed %s", s.ID, s.Elapsed))
	header := ui.JoinRight(headerLeft, headerRight, leftW)

	replyH := 4
	transH := h - 1 - 1 - replyH
	if transH < 6 {
		transH = 6
	}
	transcriptBody := renderTranscript(m, c, leftW-4, transH-2)
	transcriptPanel := ui.Panel("transcript", "j/k scroll · g top · G bottom", transcriptBody, leftW, transH, false, false)
	replyBox := renderReplyBox(m, s, leftW)
	leftCol := header + "\n" + transcriptPanel + "\n" + replyBox

	runBody := renderRunMeta(s, rightW-4)
	budBody := renderBudget(s, rightW-4)
	actBody := renderActions(rightW - 4)
	runH := 8
	budH := 9
	actH := h - runH - budH - 2

	runP := ui.Panel("run", "", runBody, rightW, runH, false, false)
	budP := ui.Panel("budget", "", budBody, rightW, budH, false, false)
	actP := ui.Panel("actions", "", actBody, rightW, actH, false, false)
	rightCol := runP + "\n" + budP + "\n" + actP

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, ui.VGap(gap, h), rightCol)
}

// KeyHints returns the status-bar key hint slice for this view.
func (m Model) KeyHints() [][2]string {
	return [][2]string{
		{"esc", "detach"}, {"j k", "scroll"}, {"r", "reply"}, {"b", "bg"}, {"t", "tail"}, {"?", "help"},
	}
}
