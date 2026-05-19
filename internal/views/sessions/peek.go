package sessions

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func renderPeek(m Model, c ctx.Ctx, rows []data.Session, w, h int) string {
	if len(rows) == 0 {
		return ui.Panel("peek", "", theme.SFaint.Render("no sessions match"), w, h, false, false)
	}
	sel := m.Sel
	if sel >= len(rows) {
		sel = len(rows) - 1
	}
	if sel < 0 {
		sel = 0
	}
	s := rows[sel]
	accent := s.State == "needs_input"

	var b strings.Builder
	pill := ui.StatePill(s.State, c.Spin)
	pillLines := strings.Split(pill, "\n")
	pillLine := pill
	if len(pillLines) >= 2 {
		pillLine = pillLines[1]
	}
	pad := ui.Max(1, w-4-lipgloss.Width(pillLine)-lipgloss.Width(s.ID)-lipgloss.Width(s.Issue)-3)
	b.WriteString(pillLine + "  " + theme.SFaint.Render(s.ID) +
		strings.Repeat(" ", pad) +
		theme.SFaint.Render(s.Issue) + "\n")
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Text)).Render(s.Title) + "\n")
	b.WriteString("\n")

	rows2 := [][3]string{
		{"agent", s.Agent + " · " + s.Model, theme.Text},
		{"workspace", s.Workspace, theme.Text},
		{"branch", s.Branch, theme.Violet},
		{"started", s.Started, theme.Text},
		{"elapsed", s.Elapsed, theme.Text},
		{"tokens", fmt.Sprintf("%s in · %s out", ui.Commafy(s.TokensIn), ui.Commafy(s.TokensOut)), theme.Text},
		{"cost", s.Cost, theme.Text},
	}
	for _, r := range rows2 {
		b.WriteString(ui.KVRow(r[0], r[1], r[2], w-4) + "\n")
	}

	if s.Question != "" {
		b.WriteString("\n")
		b.WriteString(ui.SubCard("◆ WAITING ON YOU", s.Question, w-4, theme.AccentDim) + "\n")
	}

	if len(s.Log) > 0 {
		b.WriteString("\n" + theme.SFaint.Render("RECENT") + "\n")
		for _, l := range s.Log {
			col := theme.TextDim
			switch l[0] {
			case "ask":
				col = theme.Warn
			case "tool":
				col = theme.Info
			}
			line := theme.SFaint.Render(ui.PadRight(l[0], 6)) + " " +
				theme.Fg(col).Render(ui.PadRight(l[1], 10)) + " " +
				theme.SFaint.Render(ui.Truncate(l[2], w-22))
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n" + ui.KeyChip("↵", "attach", true) + "  " +
		ui.KeyChip("r", "reply", false) + "  " +
		ui.KeyChip("b", "background", false) + "  " +
		ui.KeyChip("x", "cancel", false) + "  " +
		ui.KeyChip("*", "next needs-input", false))

	return ui.Panel("peek", "↵ "+theme.SAccent.Render("attach"), b.String(), w, h, false, accent)
}
