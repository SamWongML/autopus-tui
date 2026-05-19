package attach

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderReplyBox(m Model, s *data.Session, w int) string {
	border := theme.Border
	if s.State == "needs_input" {
		border = theme.AccentDim
	}
	header := theme.SFaint.Render("›") + " " + theme.SFaint.Render("reply to ") + theme.SDim.Render(s.ID) +
		strings.Repeat(" ", ui.Max(1, w-22-lipgloss.Width("[↵] send"))) + ui.KeyChip("↵", "send", false)
	placeholder := "type a message…"
	if s.State == "needs_input" {
		placeholder = "(a) refactor the dispatcher first, then teardown"
	}
	input := m.Reply
	rendered := theme.SFaint.Render(placeholder)
	if input != "" {
		rendered = theme.SText.Render(input)
	}
	caret := theme.SAccent.Render("▌")
	prompt := theme.SAccent.Render("› ") + rendered + caret

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(border)).
		Padding(0, 1).Width(w-2).Render(header + "\n" + prompt)
}

func renderRunMeta(s *data.Session, w int) string {
	rows := [][3]string{
		{"session", s.ID, theme.Text},
		{"agent", s.Agent + " · " + s.Model, theme.Text},
		{"branch", s.Branch, theme.Violet},
		{"started", s.Started, theme.Text},
		{"elapsed", s.Elapsed, theme.Text},
	}
	var b strings.Builder
	for _, r := range rows {
		b.WriteString(ui.KVRowDashed(r[0], r[1], r[2], w) + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderBudget(s *data.Session, w int) string {
	var b strings.Builder
	b.WriteString(ui.KVRowDashed("in", ui.Commafy(s.TokensIn), theme.Text, w) + "\n")
	b.WriteString(ui.KVRowDashed("out", ui.Commafy(s.TokensOut), theme.Text, w) + "\n")
	b.WriteString(ui.KVRowDashed("cost", s.Cost, theme.Accent, w) + "\n")
	b.WriteString(ui.KVRowDashed("cap", "$5.00", theme.Text, w) + "\n")
	pct := 0.0
	fmt.Sscanf(s.Cost, "$%f", &pct)
	pct = pct / 5.0 * 100
	b.WriteString("\n" + ui.Bar(pct, ui.Min(w, 24), theme.Accent))
	return b.String()
}

func renderActions(w int) string {
	keys := [][2]string{
		{"r", "reply"},
		{"b", "background"},
		{"t", "open log tail"},
		{"d", "diff for this run"},
		{"p", "permission requests"},
		{"x", "cancel"},
		{"!", "restart with new prompt"},
	}
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(ui.KeyChip(k[0], k[1], false) + "\n")
	}
	_ = w
	return strings.TrimRight(b.String(), "\n")
}
