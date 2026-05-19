package chrome

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// TabBar renders the numbered tab strip + attached/help right-side indicator.
// It also returns hit-boxes (route IDs + "help") so the mouse handler can map
// a click X to a target.
//
//   - activeRoute  — id of the currently focused route (e.g. "sessions")
//   - attachID     — non-empty when in the attach view (suppresses tab focus)
//   - helpOpen     — whether the help overlay is currently visible
func TabBar(w int, activeRoute, attachID string, helpOpen bool) (string, []ui.HitBox) {
	parts := make([]string, 0, len(app.Routes))
	widths := make([]int, 0, len(app.Routes))
	for _, r := range app.Routes {
		active := !helpOpen && attachID == "" && r.ID == activeRoute
		var keyCol, lblCol, keyBorder string
		if active {
			keyCol, lblCol, keyBorder = theme.Accent, theme.Accent, theme.AccentDim
		} else {
			keyCol, lblCol, keyBorder = theme.TextMute, theme.TextFaint, theme.Border
		}
		cap := theme.FgOn(keyBorder, theme.Bg2).Render("[") +
			theme.FgOn(keyCol, theme.Bg2).Render(r.Key) +
			theme.FgOn(keyBorder, theme.Bg2).Render("]")
		lblStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(lblCol)).
			Background(lipgloss.Color(theme.Bg2))
		if active {
			lblStyle = lblStyle.Underline(true)
		}
		lbl := lblStyle.Render(r.Label)
		s := cap + " " + lbl
		if r.ID == "sessions" {
			n := 0
			for _, ss := range data.Sessions {
				if ss.State == "needs_input" {
					n++
				}
			}
			if n > 0 {
				s += " " + theme.FgOn(theme.Warn, theme.Bg2).Render(fmt.Sprintf("◆%d", n))
			}
		}
		parts = append(parts, s)
		widths = append(widths, lipgloss.Width(s))
	}

	left := strings.Join(parts, theme.BG(theme.Bg2).Render("  "))

	// Tab hit-boxes. The line begins with a 1-cell pad, so tabs start at col 1.
	hits := make([]ui.HitBox, 0, len(app.Routes)+1)
	x := 1
	for i, r := range app.Routes {
		hits = append(hits, ui.HitBox{X1: x, X2: x + widths[i] - 1, ID: r.ID})
		x += widths[i]
		if i < len(app.Routes)-1 {
			x += 2
		}
	}

	attachPart := ""
	if attachID != "" {
		attachPart = theme.FgOn(theme.Accent, theme.Bg2).Render("● attached "+attachID) + " " + theme.FgOn(theme.TextFaint, theme.Bg2).Render("(esc detach)") + "   "
	}
	keyCol := theme.TextMute
	lblCol := theme.TextFaint
	keyBorder := theme.Border
	if helpOpen {
		keyCol, lblCol, keyBorder = theme.Accent, theme.Accent, theme.AccentDim
	}
	helpLblStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(lblCol)).
		Background(lipgloss.Color(theme.Bg2))
	if helpOpen {
		helpLblStyle = helpLblStyle.Underline(true)
	}
	helpPart := theme.FgOn(keyBorder, theme.Bg2).Render("[") +
		theme.FgOn(keyCol, theme.Bg2).Render("?") +
		theme.FgOn(keyBorder, theme.Bg2).Render("]") + " " +
		helpLblStyle.Render("Help")

	right := attachPart + helpPart

	gap := w - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	helpX1 := 1 + lipgloss.Width(left) + gap + lipgloss.Width(attachPart)
	helpX2 := helpX1 + lipgloss.Width(helpPart) - 1
	hits = append(hits, ui.HitBox{X1: helpX1, X2: helpX2, ID: "help"})

	line := theme.BG(theme.Bg2).Render(" ") + left + theme.BG(theme.Bg2).Render(strings.Repeat(" ", gap)) + right + theme.BG(theme.Bg2).Render(" ")
	return ui.PaintLine(line, w, theme.Bg2), hits
}
