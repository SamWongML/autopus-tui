package issues

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func renderDetail(m Model, rows []data.Issue, c ctx.Ctx, w, h int) string {
	var it data.Issue
	if len(rows) == 0 {
		it = data.Issues[0]
	} else {
		sel := m.Sel
		if sel >= len(rows) {
			sel = len(rows) - 1
		}
		it = rows[sel]
	}
	sc := theme.Statuses[it.Status]
	if sc == "" {
		sc = theme.TextDim
	}
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Text)).Render(it.Title) + "\n\n")
	b.WriteString(ui.KVRow("status", it.Status, sc, w-4) + "\n")
	b.WriteString(ui.KVRow("priority", it.Priority, theme.Priorities[it.Priority].Color, w-4) + "\n")
	b.WriteString(ui.KVRow("assignee", it.Assignee, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRow("workspace", it.Workspace, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRow("updated", it.Updated+" ago", theme.Text, w-4) + "\n")
	b.WriteString("\n" + theme.SFaint.Render("RECENT RUNS") + "\n")
	runRows := []struct {
		state, id, sub string
	}{
		{"running", "t_8fa1", "claude · started 16:02 · 39m"},
		{"failed", "t_8c92", "claude · 14:18 → 14:38 · timeout"},
		{"completed", "t_8bc0", "codex · 11:30 → 11:42 · pushed"},
	}
	for _, r := range runRows {
		line := ui.Glyph(r.state, c.Spin) + " " + theme.SDim.Render(r.id) + " " + theme.SFaint.Render(r.sub)
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + ui.KeyChip("↵", "open run", true) + "  " +
		ui.KeyChip("a", "assign", false) + "  " +
		ui.KeyChip("s", "status", false) + "  " +
		ui.KeyChip("c", "comment", false))

	return ui.Panel(it.ID, theme.Fg(sc).Render(it.Status), b.String(), w, h, false, false)
}
