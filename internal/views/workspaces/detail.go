package workspaces

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func renderDetail(m Model, c ctx.Ctx, w, h int) string {
	x := data.Workspaces[m.Sel]
	var b strings.Builder
	watchLabel, watchCol := "◇ unwatched", theme.TextMute
	if x.Watch {
		watchLabel, watchCol = "◆ watched", theme.Accent
	}
	headPad := ui.Max(1, w-4-lipgloss.Width(watchLabel)-lipgloss.Width(x.ID))
	b.WriteString(theme.Fg(watchCol).Render(watchLabel) +
		strings.Repeat(" ", headPad) +
		theme.SFaint.Render(x.ID) + "\n\n")
	b.WriteString(ui.KVRow("name", x.Name, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRow("role", x.Role, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRow("members", ui.Itoa(x.Members), theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRow("issues", ui.Itoa(x.Issues), theme.Text, w-4) + "\n")
	sessCol := theme.TextDim
	if x.Sessions > 0 {
		sessCol = theme.OK
	}
	b.WriteString(ui.KVRow("sessions", ui.Itoa(x.Sessions), sessCol, w-4) + "\n")
	b.WriteString(ui.KVRow("root", "~/autopus_workspaces/"+x.Name, theme.Text, w-4) + "\n")
	b.WriteString(ui.KVRow("disk", "2.4G / 3.0G (78%)", theme.Warn, w-4) + "\n")
	b.WriteString("\n" + theme.SFaint.Render("ACTIVE SESSIONS") + "\n")
	count := 0
	for _, s := range data.Sessions {
		if s.Workspace == x.Name && count < 4 {
			line := ui.Glyph(s.State, c.Spin) + " " + theme.SDim.Render(s.Issue) + " " +
				theme.SText.Render(ui.Truncate(s.Title, w-22))
			b.WriteString(line + "\n")
			count++
		}
	}
	if count == 0 {
		b.WriteString(theme.SFaint.Render("no active sessions") + "\n")
	}
	b.WriteString("\n")
	label := "watch"
	if x.Watch {
		label = "unwatch"
	}
	b.WriteString(ui.KeyChip("w", label, true) + "  " +
		ui.KeyChip("m", "members", false) + "  " +
		ui.KeyChip("o", "browser", false) + "  " +
		ui.KeyChip("c", "clean cache", false))
	return ui.Panel(x.Name, "", b.String(), w, h, false, false)
}
