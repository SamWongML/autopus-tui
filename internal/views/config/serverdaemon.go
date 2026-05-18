package config

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderServerDaemonCol(_ Model, w, h int) string {
	serverH := 8
	daemonH := h - serverH - 1

	var sb strings.Builder
	for _, r := range data.CfgServer {
		sb.WriteString(renderCfgRow(r, w-4) + "\n")
	}
	serverPanel := ui.Panel("server", theme.SOK.Render("● reachable"),
		strings.TrimRight(sb.String(), "\n"), w, serverH, false, false)

	var db strings.Builder
	for _, r := range data.CfgDaemon {
		db.WriteString(renderCfgRow(r, w-4) + "\n")
	}
	db.WriteString("\n" + theme.SAccent.Render("TUI · COMPLEMENTARY") + "\n")
	for _, r := range data.CfgTUI {
		db.WriteString(renderCfgRow(r, w-4) + "\n")
	}
	daemonPanel := ui.Panel("daemon", theme.SOK.Render("● up · pid 48132"),
		strings.TrimRight(db.String(), "\n"), w, daemonH, false, false)

	return serverPanel + "\n" + daemonPanel
}

func renderCfgRow(r data.CfgRow, w int) string {
	tag := ""
	if r.Tag != "" {
		c := theme.TextMute
		if r.Tag == "TUI" {
			c = theme.AccentDim
		}
		tag = " " + theme.Fg(c).Render(r.Tag)
	}
	left := theme.SFaint.Render(r.K) + tag
	col := theme.Text
	if r.Color != "" {
		col = r.Color
	}
	right := theme.Fg(col).Render(r.V)
	return ui.JoinRight(left, right, w)
}
