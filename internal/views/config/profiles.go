package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderProfilesCol(m Model, w, h int) string {
	profilesH := 4 + len(data.Profiles)*4 + 2
	if profilesH > h*2/3 {
		profilesH = h * 2 / 3
	}
	if profilesH < 10 {
		profilesH = 10
	}
	budgetsH := h - profilesH - 1
	if budgetsH < 8 {
		budgetsH = 8
	}

	var pb strings.Builder
	for _, p := range data.Profiles {
		active := p.ID == m.Profile
		marker := "○"
		col := theme.TextMute
		if active {
			marker = "◉"
			col = theme.Accent
		}
		nameCol := theme.Text
		if active {
			nameCol = theme.Accent
		}
		status := theme.SMute.Render("○ stopped")
		if p.Active {
			status = theme.SOK.Render("● running")
		}
		line1 := theme.Fg(col).Render(marker) + " " + theme.Fg(nameCol).Render(p.ID) +
			strings.Repeat(" ", ui.Max(1, w-6-lipgloss.Width(p.ID)-lipgloss.Width("○ running"))) +
			status
		line2 := "  " + theme.SFaint.Render(ui.Truncate(p.Server, w-6)) +
			strings.Repeat(" ", ui.Max(1, w-4-lipgloss.Width("  "+ui.Truncate(p.Server, w-6))-lipgloss.Width(p.Uptime))) +
			theme.SFaint.Render(p.Uptime)
		pb.WriteString(line1 + "\n" + line2 + "\n")
		if active {
			pb.WriteString(ui.Dashed(w-4) + "\n")
		} else {
			pb.WriteString("\n")
		}
	}
	pb.WriteString(ui.KeyChip("n", "new profile", false) + "  " +
		ui.KeyChip("D", "delete", false) + "  " +
		ui.KeyChip("↵", "switch", true))
	profilesPanel := ui.Panel("profiles", "~/.autopus/profiles", pb.String(), w, profilesH, false, false)

	bg := data.Budgets
	var bb strings.Builder
	bb.WriteString(ui.KVRow("daily cap", fmt.Sprintf("$%.2f", bg.DailyCap), theme.Text, w-4) + "\n")
	bb.WriteString(ui.KVRow("daily used", fmt.Sprintf("$%.2f", bg.DailyUsed), theme.Accent, w-4) + "\n")
	bb.WriteString(ui.KVRow("per-run cap", fmt.Sprintf("$%.2f", bg.PerRunCap), theme.Text, w-4) + "\n")
	bb.WriteString(ui.KVRow("warn at", fmt.Sprintf("%d%%", bg.WarnPct), theme.Text, w-4) + "\n")
	bb.WriteString("\n" + ui.Bar(bg.DailyUsed/bg.DailyCap*100, ui.Min(w-4, 26), theme.Accent) + "\n\n")
	bb.WriteString(theme.SOK.Render("● within budget") + " " + theme.SFaint.Render(fmt.Sprintf("· %d%% used today", int(bg.DailyUsed/bg.DailyCap*100))))
	budgetsPanel := ui.Panel("budgets", "", bb.String(), w, budgetsH, false, false)

	return profilesPanel + "\n" + budgetsPanel
}
