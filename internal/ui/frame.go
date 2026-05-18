package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/theme"
)

// PaintLine pads s to width w with a styled background fill so an empty cell
// shows the supplied bg color instead of the terminal default.
func PaintLine(s string, w int, bgColor string) string {
	width := lipgloss.Width(s)
	if width >= w {
		return s
	}
	return s + theme.BG(bgColor).Render(strings.Repeat(" ", w-width))
}

// VGap returns a multi-line bg-filled rectangle (w wide, h tall) used as a
// styled gutter between columns so the terminal default doesn't bleed through.
func VGap(w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	line := theme.BG(theme.Bg).Render(strings.Repeat(" ", w))
	lines := make([]string, h)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

// PadOrClipANSI pads or clips a (possibly ANSI-colored) string to exactly
// width n. Padding spaces are styled with the canvas bg so the terminal
// default color doesn't show through.
func PadOrClipANSI(s string, n int) string {
	w := lipgloss.Width(s)
	if w == n {
		return s
	}
	if w < n {
		return s + theme.BG(theme.Bg).Render(strings.Repeat(" ", n-w))
	}
	return lipgloss.NewStyle().Width(n).MaxWidth(n).Background(lipgloss.Color(theme.Bg)).Render(s)
}

// JoinVertical joins parts with newlines (lipgloss-free for performance).
func JoinVertical(parts ...string) string { return strings.Join(parts, "\n") }

// JoinHorizontal is lipgloss.JoinHorizontal with top alignment — the only mode
// this UI uses.
func JoinHorizontal(blocks ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, blocks...)
}

// Blank returns w spaces painted with the canvas bg.
func Blank(w int) string {
	return theme.BG(theme.Bg).Render(strings.Repeat(" ", w))
}

// Dashed renders a horizontal dashed separator the width of w.
func Dashed(w int) string {
	if w < 1 {
		return ""
	}
	return theme.SBorder.Render(strings.Repeat("╌", w))
}
