package ui

import (
	"strings"

	"autopus-tui/internal/theme"
)

// Bar renders a horizontal progress bar of fixed width. pct is 0..100.
func Bar(pct float64, width int, color string) string {
	if width <= 0 {
		return ""
	}
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int(float64(width)*pct/100 + 0.5)
	if filled > width {
		filled = width
	}
	a := theme.Fg(color).Render(strings.Repeat("█", filled))
	b := theme.Fg(theme.Surface3).Render(strings.Repeat("░", width-filled))
	return a + b
}

// Spark renders a unicode block sparkline using the supplied values; only the
// last `width` values are shown.
func Spark(values []float64, width int, color string) string {
	chars := []rune("▁▂▃▄▅▆▇█")
	if width <= 0 || len(values) == 0 {
		return ""
	}
	if len(values) > width {
		values = values[len(values)-width:]
	}
	maxV := 1.0
	for _, v := range values {
		if v > maxV {
			maxV = v
		}
	}
	var b strings.Builder
	for _, v := range values {
		idx := int(v / maxV * 7)
		if idx < 0 {
			idx = 0
		}
		if idx > 7 {
			idx = 7
		}
		b.WriteRune(chars[idx])
	}
	return theme.Fg(color).Render(b.String())
}
