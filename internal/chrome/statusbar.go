package chrome

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// StatusBar renders the contextual key-hint strip + version/pid footer. hints
// are de-duplicated on key before rendering. Each hint is one [key, label] pair.
func StatusBar(w int, hints [][2]string) string {
	seen := map[string]bool{}
	var parts []string
	for _, h := range hints {
		if seen[h[0]] {
			continue
		}
		seen[h[0]] = true
		cap := lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Surface2)).
			Foreground(lipgloss.Color(theme.Text)).Bold(true).
			Padding(0, 1).Render(h[0])
		parts = append(parts, cap+" "+theme.FgOn(theme.TextFaint, theme.Bg2).Render(h[1]))
	}
	left := strings.Join(parts, theme.BG(theme.Bg2).Render("  "))
	right := theme.FgOn(theme.TextMute, theme.Bg2).Render(fmt.Sprintf("autopus-tui v%s · pid %d", data.Daemon.Version, data.Daemon.PID))

	gap := w - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}
	top := theme.FgOn(theme.Border, theme.Bg2).Render(strings.Repeat("─", w))
	line := theme.BG(theme.Bg2).Render(" ") + left + theme.BG(theme.Bg2).Render(strings.Repeat(" ", gap)) + right + theme.BG(theme.Bg2).Render(" ")
	return ui.PaintLine(top, w, theme.Bg2) + "\n" + ui.PaintLine(line, w, theme.Bg2)
}
