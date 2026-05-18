package overview

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderServerCard(w, h int) string {
	headLeft := ui.Dot(theme.OK) + " " + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.OK)).Render("REACHABLE")
	headRight := theme.SFaint.Render("last hb 4s ago")
	pad := ui.Max(1, w-lipgloss.Width("● REACHABLE")-lipgloss.Width("last hb 4s ago")-1)
	header := headLeft + strings.Repeat(" ", pad) + headRight
	rows := [][3]string{
		{"app url", "app.autopus.ai", theme.Text},
		{"ws url", "api.autopus.ai/ws", theme.Text},
		{"region", "us-west · iad1", theme.Text},
		{"latency p50", "38ms", theme.OK},
		{"latency p99", "92ms", theme.Text},
		{"token", "valid · 86d left", theme.Text},
		{"user", "@you · org-acme", theme.Accent},
	}
	var b strings.Builder
	b.WriteString(header + "\n")
	for _, r := range rows {
		b.WriteString(ui.KVRow(r[0], r[1], r[2], w) + "\n")
	}
	hb := genSpark(48, 35, 5)
	b.WriteString(theme.SFaint.Render("heartbeat trace · last 60s") + "\n")
	b.WriteString(ui.Spark(hb, ui.Min(w, 48), theme.OK))
	_ = h
	return b.String()
}
