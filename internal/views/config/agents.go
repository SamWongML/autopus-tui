package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderAgentsPanel(m Model, w, h int) string {
	var ab strings.Builder
	var chips []string
	for i, a := range data.AgentCfgs {
		col := theme.TextDim
		border := theme.Border
		if i == m.AgentSel {
			col, border = theme.Accent, theme.AccentDim
		}
		styled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(col)).
			Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(border)).
			Padding(0, 1).Render(a.Name)
		chips = append(chips, styled)
	}
	chipRow := lipgloss.JoinHorizontal(lipgloss.Top, chips...)
	lines := strings.Split(chipRow, "\n")
	if len(lines) >= 2 {
		ab.WriteString(lines[1] + "\n\n")
	} else {
		ab.WriteString(chipRow + "\n\n")
	}

	a := data.AgentCfgs[m.AgentSel]
	ab.WriteString(ui.KVRow("path", ui.Truncate(a.Path, w-12), theme.Text, w-4) + "\n")
	ab.WriteString(ui.KVRow("model", a.Model, theme.Text, w-4) + "\n")
	ab.WriteString(ui.KVRow("concurrency cap", fmt.Sprintf("%d", a.Concurrency), theme.Text, w-4) + "\n")

	ab.WriteString("\n" + theme.SFaint.Render("ENV") + "\n")
	if len(a.Env) == 0 {
		ab.WriteString(theme.SMute.Render("none set · using defaults") + "\n")
	} else {
		for _, e := range a.Env {
			line := ui.JoinRight(theme.SFaint.Render(e[0]), theme.SDim.Render(e[1]), w-4)
			ab.WriteString(line + "\n")
		}
	}

	prefix := strings.ToUpper(strings.Split(a.Name, "-")[0])
	ab.WriteString("\n" + theme.SFaint.Render("override with ") +
		theme.SAccent.Render(fmt.Sprintf("AUTOPUS_%s_PATH", prefix)) +
		theme.SFaint.Render(" / ") +
		theme.SAccent.Render(fmt.Sprintf("AUTOPUS_%s_MODEL", prefix)) + "\n")
	ab.WriteString(theme.SFaint.Render("reload with ") + theme.SAccent.Render("autopus daemon restart") + "\n\n")
	ab.WriteString(ui.KeyChip("e", "edit", true) + "  " +
		ui.KeyChip("↵", "reload daemon", false) + "  " +
		ui.KeyChip("⌃R", "reset defaults", false))

	return ui.Panel("per-agent overrides", fmt.Sprintf("%d configured", len(data.AgentCfgs)),
		ab.String(), w, h, false, false)
}
