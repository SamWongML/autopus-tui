package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderTasks(width, height int, selected int) string {
	filterRow := renderFilterBar(width)
	bodyH := height - 3
	if bodyH < 6 {
		bodyH = 6
	}
	gap := 2
	rightW := 48
	if rightW > width*40/100 {
		rightW = width * 40 / 100
	}
	leftW := width - rightW - gap

	sel := clamp(selected, 0, len(tasks)-1)
	left := paneTasksTable(leftW, bodyH, sel)
	right := pane("preview · "+tasks[sel].ID, "",
		renderTaskDetail(tasks[sel], "summary", rightW-4, bodyH-2),
		rightW, bodyH, true)
	body := joinH(left, bgPad(gap), right)

	return joinV(filterRow, body)
}

func renderFilterBar(width int) string {
	chip := func(s string, on bool) string {
		st := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, false, false)
		if on {
			st = st.Foreground(accent).Background(bg2)
		} else {
			st = st.Foreground(dim).Background(bg1)
		}
		return st.Render(s)
	}
	left := sDim.Render("filter") + "  " +
		chip("status: in-flight", true) + " " +
		chip("runtime: any", false) + " " +
		chip("workspace: blackhole-os", false) + " " +
		chip("since: today", false)
	right := sDim.Render("sort: ") + sFg.Render("started ↓")
	gap := width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	row := fillBg(left+strings.Repeat(" ", gap)+right, bg)
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

// taskSections is the fixed display order for grouped task rows. The mapping
// from a Task to its section is computed by groupedTasks.
var taskSections = []string{
	"Needs input",
	"Working",
	"Ready for review",
	"Completed",
	"Failed",
	"Queued",
}

// visualTaskOrder returns task indices in the order rows are rendered on the
// Tasks tab. Keyboard navigation walks this slice so j/k move to the next
// row on screen, not the next entry in the raw `tasks` declaration order.
func visualTaskOrder() []int {
	sections, by := groupedTasks()
	out := make([]int, 0, len(tasks))
	for _, sec := range sections {
		out = append(out, by[sec]...)
	}
	return out
}

// groupedTasks classifies each task into exactly one section. The "Ready for
// review" bucket wins over "Working" / "Completed" when a PR is open; merged
// PRs fall through to the underlying status.
func groupedTasks() ([]string, map[string][]int) {
	by := map[string][]int{}
	for i, t := range tasks {
		var sec string
		switch {
		case t.Status == "waiting":
			sec = "Needs input"
		case t.PRState != "" && t.PRState != "merged":
			sec = "Ready for review"
		case t.Status == "working":
			sec = "Working"
		case t.Status == "done":
			sec = "Completed"
		case t.Status == "failed":
			sec = "Failed"
		case t.Status == "queued":
			sec = "Queued"
		}
		if sec != "" {
			by[sec] = append(by[sec], i)
		}
	}
	return taskSections, by
}

func paneTasksTable(width, height int, selected int) string {
	inner := width - 4
	statusW, idW, agentW, seqW, elapsedW, costW, prW := 2, 7, 10, 7, 8, 7, 2
	titleW := inner - statusW - idW - agentW - seqW - elapsedW - costW - prW - 7
	if titleW < 12 {
		titleW = 12
	}

	head := sDim.Render(padR("", statusW) + " " +
		padR("ID", idW) + " " +
		padR("TITLE", titleW) + " " +
		padR("RUNTIME", agentW) + " " +
		padR("SEQ", seqW) + " " +
		padR("ELAPSED", elapsedW) + " " +
		padR("COST", costW) + " " +
		padR("PR", prW))

	sections, byName := groupedTasks()

	var lines []string
	lines = append(lines, head)
	lines = append(lines, sFaint.Render(strings.Repeat("─", inner)))
	for _, sec := range sections {
		idxs := byName[sec]
		if len(idxs) == 0 {
			continue
		}
		lines = append(lines, hr(fmt.Sprintf("%s (%d)", sec, len(idxs)), inner))
		for _, idx := range idxs {
			lines = append(lines, taskRowSimple(tasks[idx], idx == selected,
				statusW, idW, titleW, agentW, seqW, elapsedW, costW, prW))
		}
	}
	body := strings.Join(lines, "\n")
	title := fmt.Sprintf("%d tasks · %d need input · %d working",
		len(tasks), len(byName["Needs input"]), len(byName["Working"]))
	return pane(title, "", body, width, height, false)
}

// taskRowSimple — cleaner than taskRow in view_status.go; used by Tasks view.
// Non-selected rows sit directly on the canvas bg (no row stripe); the
// selected row is painted with bgSel end-to-end via fillBg, so the spaces
// between fields don't fall back to canvas bg and produce a vertical strap.
func taskRowSimple(t Task, selected bool, statusW, idW, titleW, agentW, seqW, elapsedW, costW, prW int) string {
	pr := taskPRCell(t, prW)
	if selected {
		bgOn := func(c lipgloss.Color) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(c).Background(bgSel)
		}
		seg := bgOn(accent).Render("▎") +
			lipgloss.NewStyle().
				Foreground(statusStyle(t.Status).GetForeground()).
				Background(bgSel).Render(padR(statusGlyph(t.Status), statusW)) +
			bgOn(fg1).Render(" "+padR(t.ID, idW)) +
			bgOn(fg).Bold(true).Render(" "+padR(truncate(t.Title, titleW), titleW)) +
			bgOn(agentTone(t.Runtime)).Render(" "+padR(t.Runtime, agentW)) +
			bgOn(dim).Render(" "+padR(fmt.Sprintf("seq %d", t.Seq), seqW)) +
			bgOn(dim).Render(" "+padR(t.Started, elapsedW)) +
			bgOn(dim).Render(" "+padR(t.Cost, costW)) +
			" " + pr
		return fillBg(seg, bgSel)
	}
	status := statusStyle(t.Status).Render(padR(statusGlyph(t.Status), statusW))
	id := sFg1.Render(padR(t.ID, idW))
	title := sFg.Render(padR(truncate(t.Title, titleW), titleW))
	runtime := lipgloss.NewStyle().Foreground(agentTone(t.Runtime)).
		Render(padR(t.Runtime, agentW))
	seq := sDim.Render(padR(fmt.Sprintf("seq %d", t.Seq), seqW))
	elapsed := sDim.Render(padR(t.Started, elapsedW))
	cost := sDim.Render(padR(t.Cost, costW))
	return " " + status + " " + id + " " + title + " " + runtime + " " + seq + " " + elapsed + " " + cost + " " + pr
}

// taskPRCell returns a prW-wide cell — the PR dot followed by trailing
// padding, or blanks if the task has no PR. Background painting (selected vs
// canvas) is delegated to the row's outer fillBg / pane bodyStyle.
func taskPRCell(t Task, prW int) string {
	if prW < 1 {
		return ""
	}
	if t.PRState == "" {
		return strings.Repeat(" ", prW)
	}
	return prDot(t.PRState) + strings.Repeat(" ", prW-1)
}

