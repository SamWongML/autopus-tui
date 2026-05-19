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
	ID       string
	Reply    string
	Scroll   int
	PeekOpen bool // narrow-width: show right rail when true (toggled with `r`)
}

// New returns a fresh model.
func New() Model { return Model{} }

// Update handles attach-local keys. Single-char keys not bound to an action
// append to the reply buffer; enter flushes the buffer to the transcript.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	s := k.String()
	switch s {
	case "esc", "left":
		m.ID = ""
		m.Reply = ""
		return m, func() tea.Msg { return app.NavigateMsg{Attach: ""} }
	case "j", "down":
		m.Scroll++
		return m, nil
	case "k", "up":
		if m.Scroll > 0 {
			m.Scroll--
		}
		return m, nil
	case "g":
		m.Scroll = 0
		return m, nil
	case "G":
		m.Scroll = 9999
		return m, nil
	case "r":
		m.PeekOpen = !m.PeekOpen
		return m, nil
	case "enter":
		if m.Reply != "" {
			data.AppendReply(m.ID, m.Reply)
			m.Reply = ""
			m.Scroll = 9999
		}
		return m, nil
	case "backspace":
		if len(m.Reply) > 0 {
			m.Reply = m.Reply[:len(m.Reply)-1]
		}
		return m, nil
	case "space":
		m.Reply += " "
		return m, nil
	}
	if len(s) == 1 {
		m.Reply += s
	}
	return m, nil
}

// View renders the attach view. At BPxs the right rail (run/budget/actions)
// is hidden by default — toggle with `r`. The left column always gets the
// full width when the rail is hidden.
func (m Model) View(c ctx.Ctx, w, h int) string {
	s := data.FindSession(m.ID)
	gap := 1

	showRail := true
	if ui.For(w) == ui.BPxs {
		showRail = m.PeekOpen
	}

	leftW := w
	rightW := 0
	if showRail {
		rightW = ui.Min(40, ui.Max(28, w/3))
		if rightW > w-30 {
			rightW = ui.Max(0, w-30-gap)
		}
		if rightW <= 0 {
			showRail = false
			leftW = w
		} else {
			leftW = w - gap - rightW
		}
	}

	pill := ui.StatePill(s.State, c.Spin)
	pillLines := strings.Split(pill, "\n")
	pillLine := pill
	if len(pillLines) >= 2 {
		pillLine = pillLines[1]
	}

	headerRight := theme.SFaint.Render(fmt.Sprintf("%s · elapsed %s", s.ID, s.Elapsed))
	prefix := ui.KeyChip("esc", "detach", true) + "  " +
		pillLine + " " + theme.SFaint.Render(s.Issue) + " " + theme.SText.Render("· ")
	titleW := ui.Max(4, leftW-lipgloss.Width(prefix)-lipgloss.Width(headerRight)-1)
	headerLeft := prefix + theme.SText.Render(ui.Truncate(s.Title, titleW))
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

	if !showRail {
		return leftCol
	}

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
		{"esc", "detach"}, {"j k", "scroll"}, {"r", "rail"}, {"b", "bg"}, {"t", "tail"}, {"?", "help"},
	}
}
