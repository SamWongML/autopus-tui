package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// kpi is the data shape feeding a KPI tile. Title sits in the pane top
// border; value is the big number; delta is the secondary line. deltaDir
// picks the arrow glyph and tone (↑ ok, ↓ err, → dim).
type kpi struct {
	title    string
	value    string
	unit     string // optional, dim, rendered after value
	delta    string // e.g. "+2 vs 24h"; empty hides the delta line
	deltaDir int    // -1 down, 0 flat, +1 up
	tone     string // "accent" (default), "ok", "err", "warn", "info", "dim"
}

// kpiTile renders a single big-number panel sized to width × height. The
// title is uppercased into the pane top border by pane(); the body holds a
// centered big number with a small delta line below.
func kpiTile(k kpi, width, height int) string {
	inner := width - 2 // pane border eats 2 cells
	innerH := height - 2
	if inner < 4 {
		inner = 4
	}
	if innerH < 1 {
		innerH = 1
	}

	valStyle := sAccent.Bold(true)
	switch k.tone {
	case "ok":
		valStyle = sOk.Bold(true)
	case "err":
		valStyle = sErr.Bold(true)
	case "warn":
		valStyle = sWarn.Bold(true)
	case "info":
		valStyle = sInfo.Bold(true)
	case "dim":
		valStyle = sFg.Bold(true)
	}

	valLine := valStyle.Render(k.value)
	if k.unit != "" {
		valLine += sDim.Render(" " + k.unit)
	}

	rows := []string{valLine}
	if k.delta != "" {
		var arrow string
		dStyle := sDim
		switch k.deltaDir {
		case +1:
			arrow, dStyle = "↑", sOk
		case -1:
			arrow, dStyle = "↓", sErr
		default:
			arrow = "→"
		}
		rows = append(rows, dStyle.Render(arrow+" ")+sDim.Render(k.delta))
	}

	// Center each row inside inner before vertical Place. Place centers the
	// rendered block as a unit and uses the longest line for block width, so
	// shorter lines (e.g. a tiny value above a long delta) would otherwise
	// be left-aligned. Pre-padding to inner forces each line to center.
	centered := make([]string, len(rows))
	rowStyle := lipgloss.NewStyle().Width(inner).Align(lipgloss.Center).Background(bg)
	for i, r := range rows {
		centered[i] = rowStyle.Render(r)
	}
	body := lipgloss.Place(
		inner, innerH,
		lipgloss.Center, lipgloss.Center,
		strings.Join(centered, "\n"),
		lipgloss.WithWhitespaceBackground(bg),
	)
	return pane(k.title, "", body, width, height, false)
}

// kpiRow renders `tiles` side-by-side, sharing the same height and dividing
// `width` (minus inter-tile gaps) evenly. Any rounding remainder is absorbed
// by the last tile so the row width matches `width` exactly.
func kpiRow(tiles []kpi, width, height int) string {
	if len(tiles) == 0 {
		return ""
	}
	gap := 1
	totalGap := gap * (len(tiles) - 1)
	tileW := (width - totalGap) / len(tiles)
	if tileW < 12 {
		tileW = 12
	}
	parts := make([]string, 0, len(tiles)*2-1)
	for i, t := range tiles {
		w := tileW
		if i == len(tiles)-1 {
			w = width - (tileW+gap)*(len(tiles)-1)
		}
		if i > 0 {
			parts = append(parts, bgPadV(gap, height))
		}
		parts = append(parts, kpiTile(t, w, height))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
