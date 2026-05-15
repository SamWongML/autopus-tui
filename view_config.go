package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderConfig(width, height int, selected int) string {
	head := renderConfigHead(width)
	bodyH := height - 3
	if bodyH < 6 {
		bodyH = 6
	}
	gap := 2
	leftW := (width - gap) * 52 / 100
	rightW := width - leftW - gap

	// Left column: daemon pane (full height).
	left := paneCfgDaemon(leftW, bodyH, selected)
	// Right column: agents pane + logging pane stacked.
	rightTop := bodyH * 60 / 100
	rightBot := bodyH - rightTop
	right := joinV(paneCfgAgents(rightW, rightTop), paneCfgLogging(rightW, rightBot))

	body := joinH(left, bgPad(gap), right)
	return joinV(head, body)
}

func renderConfigHead(width int) string {
	left := sFg1.Render("config file") + "  " + sOk.Render("~/.multica/config.toml")
	right := sDim.Render("profile: ") + sFg.Render("default") + "   " +
		sDim.Render("last-applied: ") + sFg1.Render("14:12:08") + "   " +
		sDim.Render("dirty: ") + sWarn.Render("● 2 unsaved")
	return pageHeader(width, left, right)
}

func paneCfgDaemon(width, height int, selected int) string {
	inner := width - 4
	var lines []string
	for i, r := range cfgDaemon {
		lines = append(lines, cfgRow(r, i == selected, inner))
	}
	body := strings.Join(lines, "\n")
	return pane("daemon", "multica daemon config", body, width, height, false)
}

func cfgRow(r CfgRow, selected bool, width int) string {
	kW := 22
	// Auto-fit env column: hide when there's not enough room.
	envW := 30
	if width < 70 {
		envW = 0
	}
	vW := width - kW - envW - 2
	if vW < 16 {
		vW = 16
	}

	k := sFg1.Render(padR(r.K, kW))
	v := r.V
	if v == "" {
		v = "—"
	}
	var vColor lipgloss.Color
	vBold := false
	switch r.Tone {
	case "warn":
		vColor = warn
	case "ok":
		vColor = ok
	default:
		vColor = fg
		vBold = true
	}
	if r.Readonly {
		vColor = dim
		vBold = false
	}
	if r.Dirty {
		v = "● " + v
		vColor = warn
		vBold = false
	}
	vStyle := lipgloss.NewStyle().Foreground(vColor)
	if vBold {
		vStyle = vStyle.Bold(true)
	}
	vRendered := vStyle.Render(padR(truncate(v, vW), vW))
	var envRendered string
	if envW > 0 {
		envRendered = sInfo.Render(padR(r.Env, envW))
	}

	row := " " + k + " " + vRendered + " " + envRendered
	if selected {
		bgOn := func(c lipgloss.Color) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(c).Background(bgSel)
		}
		vSel := bgOn(vColor)
		if vBold {
			vSel = vSel.Bold(true)
		}
		seg := bgOn(accent).Render("▎") +
			bgOn(fg1).Render(padR(r.K, kW)) +
			vSel.Render(" "+padR(truncate(v, vW), vW)+" █")
		if envW > 0 {
			seg += bgOn(info).Render(padR(r.Env, envW-2))
		}
		row = fillBg(seg, bgSel)
	}

	if r.Hint != "" || r.Readonly || r.Dirty {
		hint := r.Hint
		if r.Readonly {
			hint = "[read-only] " + hint
		}
		if r.Dirty {
			hint = "● unsaved · " + hint
		}
		hintRow := strings.Repeat(" ", kW+2) + sDim.Render(hint)
		row += "\n" + hintRow
	}
	return row
}

func paneCfgAgents(width, height int) string {
	inner := width - 4
	var lines []string
	lines = append(lines, sAccent.Render("── claude "+strings.Repeat("─", maxInt(2, inner-11))))
	for _, r := range cfgClaude {
		lines = append(lines, cfgRow(r, false, inner))
	}
	lines = append(lines, "")
	lines = append(lines, sAccent.Render("── codex "+strings.Repeat("─", maxInt(2, inner-10))))
	for _, r := range cfgCodex {
		lines = append(lines, cfgRow(r, false, inner))
	}
	body := strings.Join(lines, "\n")
	return pane("agents", "auto-detected on $PATH · override per-agent", body, width, height, false)
}

func paneCfgLogging(width, height int) string {
	inner := width - 4
	var lines []string
	for _, r := range cfgLogging {
		lines = append(lines, cfgRow(r, false, inner))
	}
	body := strings.Join(lines, "\n")
	return pane("logging", "", body, width, height, false)
}
