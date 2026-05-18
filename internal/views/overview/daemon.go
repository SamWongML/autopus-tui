package overview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderDaemonCard(w, h int) string {
	d := data.Daemon
	headLeft := ui.Dot(theme.OK) + " " + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.OK)).Render("UP") +
		" " + theme.SFaint.Render("· healthy")
	headRight := theme.SFaint.Render("v" + d.Version)
	pad := ui.Max(1, w-lipgloss.Width("● UP · healthy")-lipgloss.Width("v"+d.Version)-1)
	header := headLeft + strings.Repeat(" ", pad) + headRight
	rows := [][2]string{
		{"pid", fmt.Sprintf("%d", d.PID)},
		{"profile", d.Profile},
		{"uptime", d.Uptime},
		{"poll interval", "3.0s"},
		{"heartbeat", "15.0s"},
		{"max concurrent", "20"},
		{"workspaces root", "~/autopus_workspaces"},
		{"config", "~/.autopus/profiles/default"},
	}
	var b strings.Builder
	b.WriteString(header + "\n")
	for _, r := range rows {
		vc := theme.Text
		if r[0] == "profile" {
			vc = theme.Accent
		}
		b.WriteString(ui.KVRow(r[0], r[1], vc, w) + "\n")
	}
	_ = h
	return strings.TrimRight(b.String(), "\n")
}
