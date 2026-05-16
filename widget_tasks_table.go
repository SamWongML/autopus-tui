package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// newTasksTable builds the bubbles/table for the Tasks tab. Rows are plain
// strings — bubbles/table calls runewidth.Truncate on raw values inside
// renderRow, which would corrupt pre-rendered ANSI sequences. Per-row color
// (status glyph tone, runtime tone) therefore disappears in this widget; the
// trade-off buys us free cursor handling, header rendering, and scroll. The
// status column is reduced to plain text ("● working", "✓ done") so the
// status glyph still encodes intent at single-tone resolution.
func newTasksTable() table.Model {
	t := table.New(
		table.WithColumns(tasksTableCols(120)),
		table.WithRows(tasksTableRows()),
		table.WithFocused(true),
		table.WithHeight(12),
	)
	s := table.Styles{
		Header:   lipgloss.NewStyle().Foreground(dim).Bold(false),
		Cell:     lipgloss.NewStyle(),
		Selected: lipgloss.NewStyle().Background(bgSel).Foreground(fg).Bold(true),
	}
	t.SetStyles(s)
	return t
}

// tasksTableCols sizes columns to fit `inner` cells across. TITLE flexes; the
// rest are fixed at widths matching the legacy taskTable(tableFull) layout. A
// trailing space is baked into each cell value so adjacent columns visually
// separate without Cell padding — bubbles/table renders columns edge-to-edge
// and applying Padding to Cell would break the Inline-Width math inside
// renderRow.
func tasksTableCols(inner int) []table.Column {
	statusW, idW, runtimeW, seqW, elapsedW, costW := 3, 8, 11, 8, 9, 6
	fixed := statusW + idW + runtimeW + seqW + elapsedW + costW
	titleW := inner - fixed
	if titleW < 12 {
		titleW = 12
	}
	return []table.Column{
		{Title: "", Width: statusW},
		{Title: "ID", Width: idW},
		{Title: "TITLE", Width: titleW},
		{Title: "RUNTIME", Width: runtimeW},
		{Title: "SEQ", Width: seqW},
		{Title: "ELAPSED", Width: elapsedW},
		{Title: "COST", Width: costW},
	}
}

func tasksTableRows() []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, t := range tasks {
		rows[i] = table.Row{
			statusGlyph(t.Status),
			t.ID,
			t.Title,
			t.Runtime,
			fmt.Sprintf("seq %d", t.Seq),
			t.Started,
			t.Cost,
		}
	}
	return rows
}
