package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderProfiles(width, height int, selected int) string {
	head := renderProfilesHead(width)
	bodyH := height - 3
	if bodyH < 8 {
		bodyH = 8
	}
	grid := renderProfilesGrid(width, bodyH, selected)
	return joinV(head, grid)
}

func renderProfilesHead(width int) string {
	running := 0
	for _, p := range profiles {
		if p.State == "running" {
			running++
		}
	}
	parts := []string{
		kv("active", "default", "accent", 30),
		kv("profiles", fmt.Sprintf("%d", len(profiles)), "", 22),
		kv("running", fmt.Sprintf("%d/%d", running, len(profiles)), "ok", 22),
		kv("state dir", "~/.multica/profiles", "dim", 40),
	}
	row := fillBg(strings.Join(parts, "   "), bg)
	return lipgloss.NewStyle().
		Background(bg).
		Width(width).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(faint).
		BorderBackground(bg).
		Render(row)
}

func renderProfilesGrid(width, height int, selected int) string {
	// 2 columns × 2 rows: 3 profiles + 1 "new profile" placeholder.
	gap := 2
	leftW := (width - gap) / 2
	// rightW absorbs the rounding remainder when width is odd, so the grid
	// exactly fills `width` and leaves no unpainted column at the right edge.
	rightW := width - leftW - gap
	rowH := (height - gap) / 2
	if rowH < 8 {
		rowH = 8
	}

	cardW := func(col int) int {
		if col == 0 {
			return leftW
		}
		return rightW
	}

	cards := make([]string, 4)
	for i, p := range profiles {
		cards[i] = renderProfileCard(p, cardW(i%2), rowH, i == selected, i == 0)
	}
	cards[3] = renderNewProfileCard(cardW(1), rowH)

	row1 := joinH(cards[0], bgPad(gap), cards[1])
	row2 := joinH(cards[2], bgPad(gap), cards[3])
	rowWidth := lipgloss.Width(row1)
	return joinV(row1, bgPadV(rowWidth, 1), row2)
}

func renderProfileCard(p Profile, width, height int, selected, isDefault bool) string {
	inner := width - 4
	dotColor := ok
	stateColor := ok
	if p.State == "stopped" {
		dotColor = faint
		stateColor = dim
	}

	// Header line: dot name [DEFAULT] state
	head := lipgloss.NewStyle().Foreground(dotColor).Render("●") + "  " +
		lipgloss.NewStyle().Foreground(fg).Bold(true).Render(p.Name)
	if isDefault {
		head += "  " + lipgloss.NewStyle().
			Foreground(bg).Background(accent).Bold(true).Padding(0, 1).Render("DEFAULT")
	}
	stateTag := lipgloss.NewStyle().Foreground(stateColor).Background(bg2).Padding(0, 1).Render(strings.ToUpper(p.State))
	gap := inner - lipgloss.Width(head) - lipgloss.Width(stateTag)
	if gap < 1 {
		gap = 1
	}
	headRow := head + strings.Repeat(" ", gap) + stateTag

	pidStr := "—"
	if p.PID > 0 {
		pidStr = fmt.Sprintf("%d", p.PID)
	}
	kvs := []string{
		kv("pid", pidStr, "", inner/2),
		kv("uptime", p.Uptime, "", inner/2),
		kv("server", p.Server, "dim", inner),
		kv("health port", fmt.Sprintf("127.0.0.1:%d", p.Port), "", inner),
		kv("workspaces", p.WS, "dim", inner),
		kv("runtimes", strings.Join(p.Runtimes, ", "), "", inner),
		kv("tasks", p.Tasks, "", inner),
	}

	var actions string
	kb := func(s string) string {
		return lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render(s)
	}
	if p.State == "running" {
		actions = kb("⏎") + " " + sFg1.Render("attach") + "   " +
			kb("S") + " " + sFg1.Render("stop") + "   " +
			kb("r") + " " + sFg1.Render("restart") + "   " +
			kb("l") + " " + sFg1.Render("log")
	} else {
		actions = kb("s") + " " + sFg1.Render("start") + "   " +
			kb("e") + " " + sFg1.Render("edit config") + "   " +
			kb("x") + " " + sFg1.Render("delete")
	}

	body := headRow + "\n" +
		sFaint.Render(strings.Repeat("─", inner)) + "\n" +
		strings.Join(kvs, "\n") + "\n" +
		hr("actions", inner) + "\n" +
		actions

	accentB := selected
	p2 := pane(p.Name, "", body, width, height, accentB)
	if p.State == "stopped" {
		// Dim the whole card.
		p2 = lipgloss.NewStyle().Foreground(dim).Render(p2)
	}
	return p2
}

func renderNewProfileCard(width, height int) string {
	inner := width - 4
	body := sDim.Render("+  new profile") + "\n" +
		sFaint.Render(strings.Repeat("─", inner)) + "\n" +
		sDim.Render("press ") + lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render("n") +
		sDim.Render(" to create.") + "\n\n" +
		sDim.Render("profiles let one machine connect to multiple servers (cloud +\nself-host) or run isolated daemons per workspace, each with its\nown state dir, health port, and workspaces root.") + "\n" +
		hr("example", inner) + "\n" +
		lipgloss.NewStyle().Foreground(fg1).Background(bg).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(info).Padding(0, 1).Render(
			"$ multica --profile staging \\\n    daemon start \\\n    --server wss://stg.multica.ai \\\n    --health-port 7718")

	return pane("new profile", "", body, width, height, false)
}
