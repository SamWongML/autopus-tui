package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Cell is one entry in a fixed-width row. Text may already contain ANSI SGR
// styling — Row preserves it. AlignR=true right-aligns the cell within its
// width budget.
type Cell struct {
	Text   string
	AlignR bool
}

// Row builds a single line composed of the given cells, each clipped or
// padded to its target width, joined by gap-cell spaces.
//
// The returned line is guaranteed to be exactly
//
//	sum(widths) + max(0, len(cells)-1) * gap
//
// cells wide — never more, never less. This is the safe input for a Panel
// body: any oversize cell is hard-clipped via ClipANSI rather than wrapping
// onto a new line (which would break the panel's right border).
//
// Mismatched lengths between cells and widths are tolerated; extras on either
// side are ignored.
func Row(cells []Cell, widths []int, gap int) string {
	n := len(cells)
	if len(widths) < n {
		n = len(widths)
	}
	if n == 0 {
		return ""
	}
	if gap < 0 {
		gap = 0
	}

	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		w := widths[i]
		if w <= 0 {
			continue
		}
		t := cells[i].Text
		tw := lipgloss.Width(t)
		switch {
		case tw == w:
			parts = append(parts, t)
		case tw < w:
			pad := strings.Repeat(" ", w-tw)
			if cells[i].AlignR {
				parts = append(parts, pad+t)
			} else {
				parts = append(parts, t+pad)
			}
		default:
			parts = append(parts, ClipANSI(t, w))
		}
	}
	if gap == 0 {
		return strings.Join(parts, "")
	}
	return strings.Join(parts, strings.Repeat(" ", gap))
}

// RowFromCols is a convenience for the common case of building a row from a
// []Col layout and a map of widths produced by LayoutCols. text is keyed by
// Col.ID. Missing keys render as empty cells of the proper width.
func RowFromCols(cols []Col, widths map[string]int, text map[string]string, gap int) string {
	cells := make([]Cell, 0, len(cols))
	ws := make([]int, 0, len(cols))
	for _, c := range cols {
		w, ok := widths[c.ID]
		if !ok || w <= 0 {
			continue
		}
		cells = append(cells, Cell{Text: text[c.ID], AlignR: c.AlignR})
		ws = append(ws, w)
	}
	return Row(cells, ws, gap)
}
