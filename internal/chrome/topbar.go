// Package chrome renders the persistent frame around the body: top bar (brand
// + daemon health + counts + clock), tab bar (numbered routes), and status
// bar (contextual key hints + version footer). It is stateless — every
// renderer receives all data it needs as arguments.
package chrome

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// TopBar renders the brand line: ◆ autopus · agent daemon │ ● daemon up …
// │ profile … srv … ◆N need you · M working │ cpu mem │ HH:MM:SS
func TopBar(w int, now time.Time) string {
	if w < 20 {
		w = 20
	}
	needsInput, working := 0, 0
	for _, s := range data.Sessions {
		if s.State == "needs_input" {
			needsInput++
		}
		if s.State == "working" || s.State == "running" {
			working++
		}
	}

	sep := theme.FgOn(theme.TextMute, theme.Bg).Render("│")
	bullet := theme.FgOn(theme.TextFaint, theme.Bg).Render("·")

	brand := theme.FgOn(theme.Accent, theme.Bg).Render("◆") + " " +
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Text)).Background(lipgloss.Color(theme.Bg)).Render(app.Name) +
		" " + theme.FgOn(theme.TextFaint, theme.Bg).Render("· "+app.Tagline)

	daemonStatus := theme.FgOn(theme.OK, theme.Bg).Render("●") + " " +
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.OK)).Background(lipgloss.Color(theme.Bg)).Render("daemon up") +
		" " + theme.FgOn(theme.TextFaint, theme.Bg).Render(data.Daemon.Uptime)

	srv := theme.FgOn(theme.TextFaint, theme.Bg).Render("profile") + " " + theme.FgOn(theme.Accent, theme.Bg).Render(data.Daemon.Profile) + " " + bullet + " " +
		theme.FgOn(theme.TextFaint, theme.Bg).Render("srv") + " " + theme.FgOn(theme.TextDim, theme.Bg).Render(app.AppHost) + " " + theme.FgOn(theme.OK, theme.Bg).Render("↑ 38ms")

	ts := now.Format("15:04:05")

	left := brand + " " + sep + " " + daemonStatus + " " + sep + " " + srv

	var rightParts []string
	if needsInput > 0 {
		rightParts = append(rightParts,
			theme.FgOn(theme.Warn, theme.Bg).Render("◆")+" "+
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Warn)).Background(lipgloss.Color(theme.Bg)).Render(fmt.Sprintf("%d", needsInput))+" "+
				theme.FgOn(theme.Warn, theme.Bg).Render("need you"))
	}
	rightParts = append(rightParts,
		theme.FgOn(theme.TextFaint, theme.Bg).Render("·")+" "+theme.FgOn(theme.Info, theme.Bg).Render(fmt.Sprintf("%d", working))+" "+theme.FgOn(theme.TextFaint, theme.Bg).Render("working"),
		sep,
		theme.FgOn(theme.TextDim, theme.Bg).Render("cpu 18% · mem 41%"),
		sep,
		theme.FgOn(theme.Text, theme.Bg).Render(ts),
	)
	right := strings.Join(rightParts, " ")

	gap := w - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}
	line := theme.BG(theme.Bg).Render(" ") + left + theme.BG(theme.Bg).Render(strings.Repeat(" ", gap)) + right + theme.BG(theme.Bg).Render(" ")
	return ui.PaintLine(line, w, theme.Bg)
}
