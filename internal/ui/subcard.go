package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/theme"
)

// SubCard renders a small bordered callout nested inside another Panel. It's
// the WAITING-ON-YOU / inline-note shape used in peek / attach / palette.
//
// outerW is the total cell width the sub-card occupies (border + padding +
// content + padding + border). The internal arithmetic produces a body wrap
// width of outerW-4 — the value lipgloss.Wrap actually needs given Padding(0,
// 1) and a 1-cell border on each side. Callers should NOT pre-wrap body;
// pass the raw string and SubCard does the right thing.
//
// borderHex selects the border color. Pass "" for theme.AccentDim.
func SubCard(title, body string, outerW int, borderHex string) string {
	if outerW < 6 {
		outerW = 6
	}
	if borderHex == "" {
		borderHex = theme.AccentDim
	}
	innerW := outerW - 4 // 2 border + 2 padding
	if innerW < 1 {
		innerW = 1
	}

	wrapped := Wrap(body, innerW)
	header := theme.SAccent.Render(title)
	content := header
	if wrapped != "" {
		content += "\n" + wrapped
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(borderHex)).
		Background(lipgloss.Color(theme.Bg)).
		Foreground(lipgloss.Color(theme.Text)).
		Padding(0, 1).
		Width(outerW - 2). // lipgloss adds border on top → total = outerW
		Render(content)
}

// SubCardLines is SubCard but returns the lines already split, so callers
// who are composing a multi-line panel body can append them line-by-line.
func SubCardLines(title, body string, outerW int, borderHex string) []string {
	return strings.Split(SubCard(title, body, outerW, borderHex), "\n")
}
