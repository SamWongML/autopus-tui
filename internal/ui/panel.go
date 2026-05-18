package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/theme"
)

// Panel draws a bordered box with an HTML-style header row.
//
// Layout (w cols wide, h rows tall):
//
//	┌───────────────────────────────┐
//	│ TITLE                  right  │   ← header (accent-faint bg when focused)
//	├───────────────────────────────┤   ← divider
//	│ body line 1                   │
//	│ body line 2                   │
//	└───────────────────────────────┘
//
// Inner content area is h-4 rows by w-4 columns.
func Panel(title, right, body string, w, h int, focused, accent bool) string {
	if w < 4 {
		w = 4
	}
	if h < 4 {
		h = 4
	}
	borderCol := theme.Border
	titleCol := theme.TextDim
	headerBg := theme.Bg
	if focused {
		borderCol = theme.Accent
		titleCol = theme.Accent
		headerBg = theme.AccentFaint
	} else if accent {
		borderCol = theme.AccentDim
		titleCol = theme.Accent
	}

	innerW := w - 2
	bs := theme.Fg(borderCol).Background(lipgloss.Color(theme.Bg))

	top := bs.Render("┌" + strings.Repeat("─", innerW) + "┐")

	titleText := strings.ToUpper(title)
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(titleCol)).
		Background(lipgloss.Color(headerBg))
	titleRendered := titleStyle.Render(titleText)
	rightRendered := ""
	if right != "" {
		rightRendered = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.TextFaint)).
			Background(lipgloss.Color(headerBg)).
			Render(right)
	}
	used := 1 + lipgloss.Width(titleRendered) + lipgloss.Width(rightRendered) + 1
	if rightRendered != "" {
		used += 1
	}
	fill := innerW - used
	if fill < 0 {
		over := -fill
		visibleLen := lipgloss.Width(titleText) - over - 1
		if visibleLen < 1 {
			visibleLen = 1
		}
		titleText = Truncate(titleText, visibleLen)
		titleRendered = titleStyle.Render(titleText)
		used = 1 + lipgloss.Width(titleRendered) + lipgloss.Width(rightRendered) + 1
		if rightRendered != "" {
			used += 1
		}
		fill = innerW - used
		if fill < 0 {
			fill = 0
		}
	}
	headerBgFill := lipgloss.NewStyle().Background(lipgloss.Color(headerBg))
	headerInner := headerBgFill.Render(" ") + titleRendered
	if rightRendered != "" {
		headerInner += headerBgFill.Render(strings.Repeat(" ", fill+1)) + rightRendered + headerBgFill.Render(" ")
	} else {
		headerInner += headerBgFill.Render(strings.Repeat(" ", fill+1))
	}
	header := bs.Render("│") + headerInner + bs.Render("│")

	div := bs.Render("├" + strings.Repeat("─", innerW) + "┤")

	contentW := innerW - 2
	if contentW < 1 {
		contentW = 1
	}
	contentH := h - 4
	if contentH < 0 {
		contentH = 0
	}
	rawLines := SplitLines(body)
	bodyPad := lipgloss.NewStyle().Background(lipgloss.Color(theme.Bg))
	lines := make([]string, 0, contentH)
	for _, ln := range rawLines {
		ln = PadOrClipANSI(ln, contentW)
		lines = append(lines, bs.Render("│")+bodyPad.Render(" ")+ln+bodyPad.Render(" ")+bs.Render("│"))
	}
	emptyLine := bs.Render("│") + bodyPad.Render(strings.Repeat(" ", contentW+2)) + bs.Render("│")
	for len(lines) < contentH {
		lines = append(lines, emptyLine)
	}
	if len(lines) > contentH {
		lines = lines[:contentH]
	}

	bot := bs.Render("└" + strings.Repeat("─", innerW) + "┘")

	parts := []string{top, header, div}
	parts = append(parts, lines...)
	parts = append(parts, bot)
	return strings.Join(parts, "\n")
}
