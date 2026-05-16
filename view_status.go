package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// renderStatus builds the Status tab as a dashboard: a row of 4 big-number
// KPI tiles, a row of 2 sparkline panes, and a recent-failures feed below.
// The old daemon/runtimes/mini-task/log-tail panes were removed in Phase C —
// daemon identity now lives in the tab-bar meta, runtime info surfaces via
// the runtimes-online KPI, and task/log detail lives on its own tabs.
//
// The `_ unused` selectedTask argument is preserved for parity with main.go's
// switch arm; selection no longer drives anything on Status.
func renderStatus(width, height int, _ int) string {
	// Row heights scale with available body height. At 120×40 (bodyH=34) we
	// land on 7 / 12 / 15. At 80×24 (bodyH=18) on 5 / 6 / 7. At 200×60
	// (bodyH=54) on 9 / 18 / 27. Always sums to exactly `height`.
	kpiH := height * 21 / 100
	if kpiH < 5 {
		kpiH = 5
	}
	sparkH := height * 36 / 100
	if sparkH < 6 {
		sparkH = 6
	}
	failH := height - kpiH - sparkH
	if failH < 5 {
		failH = 5
		kpiH = (height - failH) * 37 / 100
		if kpiH < 5 {
			kpiH = 5
		}
		sparkH = height - kpiH - failH
	}

	row1 := kpiRow(statusKPIs(), width, kpiH)
	row2 := sparkRow(width, sparkH)
	row3 := failuresPane(width, failH)
	return joinV(row1, row2, row3)
}

// statusKPIs derives the four Status-tab KPI tiles from current mock fixtures.
// Phase D replaces the call sites with /health-derived values.
func statusKPIs() []kpi {
	active := 0
	for _, t := range tasks {
		switch t.Status {
		case "working", "waiting", "queued":
			active++
		}
	}
	online := 0
	totalCap := 0
	for _, r := range runtimes {
		online++
		totalCap += r.Max
	}
	// Titles are intentionally single words so the pane top border doesn't
	// truncate at 80-cell width. The value+unit+delta lines carry the rest
	// of the context (e.g. "4 / 7 cap", "+2 vs 1h").
	return []kpi{
		{
			title:    "active",
			value:    fmt.Sprintf("%d", active),
			unit:     fmt.Sprintf("/ %d cap", totalCap),
			delta:    "+2 vs 1h",
			deltaDir: +1,
			tone:     "accent",
		},
		{
			title:    "online",
			value:    fmt.Sprintf("%d", online),
			unit:     fmt.Sprintf("/ %d runtimes", len(runtimes)),
			delta:    "claude · codex",
			deltaDir: 0,
			tone:     "ok",
		},
		{
			title:    "tokens",
			value:    "187k",
			unit:     "today",
			delta:    "+18% vs 24h",
			deltaDir: +1,
			tone:     "info",
		},
		{
			title:    "errors",
			value:    fmt.Sprintf("%d", len(failures)),
			unit:     "today",
			delta:    "-2 vs 24h",
			deltaDir: -1,
			tone:     "err",
		},
	}
}

// sparkRow lays out the tasks/hr and lines/min sparkline panes side-by-side.
// Each pane gets half the row minus a 2-cell gap; rounding remainder goes to
// the right pane so the row width matches `width` exactly.
func sparkRow(width, height int) string {
	gap := 2
	leftW := (width - gap) / 2
	rightW := width - leftW - gap
	left := sparklinePane(
		"tasks / hour", "last 24h",
		tasksPerHour, "",
		leftW, height,
	)
	right := sparklinePane(
		"lines / min", "last 30 min",
		linesPerMin, "",
		rightW, height,
	)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, bgPad(gap), right)
}

// failuresPane wraps the recentFailures widget in a pane sized to the
// remaining row of the Status tab.
func failuresPane(width, height int) string {
	body := recentFailures(failures, width, height)
	return pane("recent failures", "last 24h · enter to open task", body, width, height, false)
}
