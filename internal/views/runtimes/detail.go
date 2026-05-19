package runtimes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderDetail(m Model, w, h int) string {
	r := data.Runtimes[m.Sel]
	col := theme.TextFaint
	switch r.Status {
	case "ready":
		col = theme.OK
	case "stale":
		col = theme.Warn
	}

	var b strings.Builder
	statusHead := "● " + strings.ToUpper(r.Status)
	lastUsed := "last used " + r.LastUsed
	headPad := ui.Max(1, w-4-lipgloss.Width(statusHead)-lipgloss.Width(lastUsed))
	b.WriteString(ui.Dot(col) + " " +
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(col)).Render(strings.ToUpper(r.Status)) +
		strings.Repeat(" ", headPad) +
		theme.SFaint.Render(lastUsed) + "\n\n")
	b.WriteString(ui.KVRowDashed("cli", r.CLI, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRowDashed("version", r.Version, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRowDashed("path", ui.Truncate(r.Path, w-12), theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRowDashed("model", r.Model, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRowDashed("concurrency",
		fmt.Sprintf("%d in-flight · cap %d", r.Inflight, r.Cap),
		theme.Text, w-4) + "\n")
	b.WriteString("\n" + ui.Bar(float64(r.Inflight)/float64(ui.Max(1, r.Cap))*100, ui.Min(w-4, 28), col) + "\n")
	b.WriteString("\n" + theme.SFaint.Render("LAST 24H") + "\n")
	sc := theme.Text
	if r.Success > 0.9 {
		sc = theme.OK
	}
	b.WriteString(ui.KVRowDashed("success", fmt.Sprintf("%d%%", int(r.Success*100)), sc, w-4) + "\n")
	b.WriteString(ui.KVRowDashed("tokens", r.Tokens24h, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRowDashed("cost", r.Cost24h, theme.Accent, w-4) + "\n")
	b.WriteString("\n" + theme.SFaint.Render("ENV OVERRIDES") + "\n")
	cli := strings.ToUpper(r.CLI)
	b.WriteString(theme.SFaint.Render(fmt.Sprintf("AUTOPUS_%s_PATH=", cli)) +
		theme.SDim.Render(ui.Truncate(r.Path, w-22)) + "\n")
	b.WriteString(theme.SFaint.Render(fmt.Sprintf("AUTOPUS_%s_MODEL=", cli)) +
		theme.SDim.Render(r.Model) + "\n")
	b.WriteString("\n" + ui.KeyChip("↵", "edit", true) + "  " +
		ui.KeyChip("t", "test handshake", false) + "  " +
		ui.KeyChip("d", "disable", false) + "  " +
		ui.KeyChip("↺", "rescan", false))

	return ui.Panel("runtime · "+r.CLI, "", b.String(), w, h, false, false)
}
