package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/theme"
)

// Glyph returns a state's static glyph (or spinner frame) painted with its
// color. Pass the current spinner frame counter for animated states.
func Glyph(state string, spinFrame int) string {
	m := theme.State(state)
	if m.SpinnerOn {
		return theme.Fg(m.Color).Render(theme.SpinnerFrames[spinFrame%len(theme.SpinnerFrames)])
	}
	return theme.Fg(m.Color).Render(m.Glyph)
}

// StatePill renders an outlined chip with glyph + UPPERCASE label.
func StatePill(state string, spinFrame int) string {
	m := theme.State(state)
	g := Glyph(state, spinFrame)
	lbl := theme.Fg(m.Color).Render(strings.ToUpper(m.Label))
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(m.Color)).
		Padding(0, 1).Render(g + " " + lbl)
}

// Dot renders a colored bullet using the • char.
func Dot(color string) string {
	return theme.Fg(color).Render("●")
}

// KeyChip renders a flat "[k] label" hint. accent uses the accent color for k.
func KeyChip(k, label string, accent bool) string {
	keyCol := theme.TextDim
	if accent {
		keyCol = theme.Accent
	}
	cap := theme.Fg(keyCol).Render("[" + k + "]")
	return cap + " " + theme.SFaint.Render(label)
}

// KeyChipBoxed renders a "boxed" key cap on a surface bg + label. Used in
// status bars where a real border isn't drawable.
func KeyChipBoxed(k, label string, accent bool) string {
	keyCol := theme.TextDim
	if accent {
		keyCol = theme.Accent
	}
	bg := theme.Surface2
	if accent {
		bg = theme.AccentFaint
	}
	cap := lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(keyCol)).
		Padding(0, 1).Render(k)
	return cap + " " + theme.SFaint.Render(label)
}
