package runtimes

import (
	"fmt"
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderTable(m Model, w, h int) string {
	const (
		cDot   = 2
		cCli   = 13
		cVer   = 8
		cModel = 14
		cStat  = 11
		cInCap = 7
		cLastU = 12
	)
	gaps := 6
	cPath := w - cDot - cCli - cVer - cModel - cStat - cInCap - cLastU - gaps
	if cPath < 10 {
		cPath = 10
	}
	header := []string{
		theme.SFaint.Render(ui.PadRight("", cDot)),
		theme.SFaint.Render(ui.PadRight("CLI", cCli)),
		theme.SFaint.Render(ui.PadRight("VER", cVer)),
		theme.SFaint.Render(ui.PadRight("PATH", cPath)),
		theme.SFaint.Render(ui.PadRight("MODEL", cModel)),
		theme.SFaint.Render(ui.PadRight("STATUS", cStat)),
		theme.SFaint.Render(ui.PadLeft("IN/CAP", cInCap)),
		theme.SFaint.Render(ui.PadRight("LAST USED", cLastU)),
	}
	out := []string{strings.Join(header, " "), theme.SBorder.Render(strings.Repeat("─", w))}

	visH := h - 2
	for i, r := range data.Runtimes {
		selected := i == m.Sel
		marker := " "
		if selected {
			marker = theme.SAccent.Render("▎")
		}
		var col, dotG string
		switch r.Status {
		case "ready":
			col, dotG = theme.OK, "●"
		case "stale":
			col, dotG = theme.Warn, "◐"
		case "not_found":
			col, dotG = theme.TextMute, "✕"
		default:
			col, dotG = theme.TextFaint, "○"
		}
		cliCol := theme.Text
		if r.Status == "not_found" {
			cliCol = theme.TextMute
		}
		statusText := r.Status
		if r.Status == "not_found" {
			statusText = "missing"
		}
		line := marker + " " +
			theme.Fg(col).Render(dotG) + " " +
			theme.Fg(cliCol).Render(ui.PadRight(r.CLI, cCli)) + " " +
			theme.SDim.Render(ui.PadRight(r.Version, cVer)) + " " +
			theme.SFaint.Render(ui.PadRight(r.Path, cPath)) + " " +
			theme.SDim.Render(ui.PadRight(r.Model, cModel)) + " " +
			theme.Fg(col).Render(ui.PadRight(statusText, cStat)) + " " +
			theme.SDim.Render(ui.PadLeft(fmt.Sprintf("%d/%d", r.Inflight, r.Cap), cInCap)) + " " +
			theme.SFaint.Render(ui.PadRight(r.LastUsed, cLastU))
		if selected {
			line = theme.WithBg(line, theme.AccentFaint)
		}
		out = append(out, line)
	}
	for len(out)-2 < visH {
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}
