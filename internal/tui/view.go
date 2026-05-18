package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/chrome"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

// View renders the full terminal frame: 2-line top bar, 2-line tab bar, body,
// 2-line status bar. Each chrome strip already paints its own bg; the body
// uses canvas bg. paintFrame pads every line to width with the right bg
// color so the terminal default never bleeds through.
func (m Model) View() string {
	if m.quitting || m.W == 0 || m.H == 0 {
		return ""
	}

	top := chrome.TopBar(m.W, m.Now)
	tab := chrome.TabBar(m.W, m.Route, m.Attach, m.Overlay == "help")
	stat := chrome.StatusBar(m.W, m.activeKeyHints())

	// Chrome occupies 6 visual rows total (each region is content + border).
	bodyH := m.H - 6
	if bodyH < 4 {
		bodyH = 4
	}

	body := m.renderBody(bodyH)

	full := theme.WithBg(top, theme.Bg) + "\n" +
		theme.WithBg(tab, theme.Bg2) + "\n" +
		theme.WithBg(body, theme.Bg) + "\n" +
		theme.WithBg(stat, theme.Bg2)
	return paintFrame(full, m.W, m.H)
}

// renderBody picks the right view (or overlay) and returns its frame.
func (m Model) renderBody(bodyH int) string {
	c := m.frameCtx()

	switch m.Overlay {
	case "onboarding":
		return m.Onboarding.View(c, m.W, bodyH)
	case "help":
		return m.Help.View(c, m.W, bodyH)
	}

	base := m.routeView(bodyH)

	if m.Overlay == "palette" {
		over := m.Palette.View(c)
		return overlayCenter(base, over, m.W, bodyH)
	}
	return base
}

// routeView renders the body for the active route (or attach overlay-route).
func (m Model) routeView(bodyH int) string {
	c := m.frameCtx()
	if m.Attach != "" {
		return m.AttachM.View(c, m.W, bodyH)
	}
	switch m.Route {
	case "overview":
		return m.Overview.View(c, m.W, bodyH)
	case "sessions":
		return m.Sessions.View(c, m.W, bodyH)
	case "issues":
		return m.Issues.View(c, m.W, bodyH)
	case "runtimes":
		return m.Runtimes.View(c, m.W, bodyH)
	case "workspaces":
		return m.Workspaces.View(c, m.W, bodyH)
	case "logs":
		return m.Logs.View(c, m.W, bodyH)
	case "config":
		return m.Config.View(c, m.W, bodyH)
	}
	return ""
}

// paintFrame splits full into lines and pads each to exactly w columns. Chrome
// rows have already been bg-patched; we just need to fill any trailing
// horizontal gap so the terminal default doesn't show through.
func paintFrame(full string, w, h int) string {
	lines := strings.Split(full, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for len(lines) < h {
		lines = append(lines, "")
	}
	for i, ln := range lines {
		lines[i] = ui.PaintLine(ln, w, theme.Bg)
	}
	return strings.Join(lines, "\n")
}

// overlayCenter places `over` centered on top of `base` (both already
// rendered). Both are line-major terminal strings; we splice line-by-line.
func overlayCenter(base, over string, w, h int) string {
	overLines := strings.Split(over, "\n")
	overH := len(overLines)
	overW := 0
	for _, l := range overLines {
		if x := lipgloss.Width(l); x > overW {
			overW = x
		}
	}
	if overH > h {
		overH = h
		overLines = overLines[:overH]
	}
	if overW > w {
		overW = w
	}
	top := (h - overH) / 2
	left := (w - overW) / 2
	if top < 0 {
		top = 0
	}
	if left < 0 {
		left = 0
	}
	baseLines := strings.Split(base, "\n")
	for len(baseLines) < h {
		baseLines = append(baseLines, "")
	}
	for i, ln := range overLines {
		row := top + i
		if row >= len(baseLines) {
			break
		}
		baseLines[row] = spliceLine(baseLines[row], ln, left, w)
	}
	return strings.Join(baseLines[:h], "\n")
}

// spliceLine writes `over` into `base` starting at column `col`, padding to
// width `w` with canvas-bg spaces afterwards.
func spliceLine(base, over string, col, w int) string {
	leftStyle := lipgloss.NewStyle().Width(col).MaxWidth(col)
	leftPart := leftStyle.Render(base)
	rest := strings.Repeat(" ", ui.Max(0, w-col-lipgloss.Width(over)))
	return leftPart + over + rest
}
