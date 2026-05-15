package main

import "github.com/charmbracelet/lipgloss"

// Warm-dark palette. Approximates the OKLCH tokens from the Multica Warm
// design (`Multica Warm.html`). Lipgloss truecolor; if the terminal can't do
// truecolor it falls back to the nearest 256-color cell.
var (
	bg     = lipgloss.Color("#19140f")
	bg1    = lipgloss.Color("#221c16")
	bg2    = lipgloss.Color("#2c241d")
	bgSel  = lipgloss.Color("#3a2f25")
	fg     = lipgloss.Color("#ebe3d5")
	fg1    = lipgloss.Color("#c4baa9")
	dim    = lipgloss.Color("#8a7f72")
	faint  = lipgloss.Color("#574d43")
	accent = lipgloss.Color("#e5b563")
	ok     = lipgloss.Color("#a7c47c")
	warn   = lipgloss.Color("#e5b563")
	errCol = lipgloss.Color("#d97a6b")
	info   = lipgloss.Color("#8eb0d5")
	violet = lipgloss.Color("#b29ad5")
	rose   = lipgloss.Color("#d59b9a")
)

// Reusable styles.
var (
	sFg     = lipgloss.NewStyle().Foreground(fg)
	sFg1    = lipgloss.NewStyle().Foreground(fg1)
	sDim    = lipgloss.NewStyle().Foreground(dim)
	sFaint  = lipgloss.NewStyle().Foreground(faint)
	sAccent = lipgloss.NewStyle().Foreground(accent)
	sOk     = lipgloss.NewStyle().Foreground(ok)
	sWarn   = lipgloss.NewStyle().Foreground(warn)
	sErr    = lipgloss.NewStyle().Foreground(errCol)
	sInfo   = lipgloss.NewStyle().Foreground(info)
	sViolet = lipgloss.NewStyle().Foreground(violet)
)

// Status glyphs — ASCII-safe so they render even where box-drawing fails.
func statusGlyph(s string) string {
	switch s {
	case "working":
		return "▶"
	case "waiting":
		return "!"
	case "done":
		return "✓"
	case "failed":
		return "✗"
	case "idle":
		return "·"
	case "stopped":
		return "▪"
	case "queued":
		return "…"
	}
	return "·"
}

func statusStyle(s string) lipgloss.Style {
	switch s {
	case "working":
		return sAccent
	case "waiting":
		return sWarn
	case "done":
		return sOk
	case "failed":
		return sErr
	case "idle", "queued":
		return sDim
	case "stopped":
		return sFaint
	}
	return sDim
}

// Agent chip tone, mirroring AGENT_COLORS in tui.jsx.
func agentTone(name string) lipgloss.Color {
	switch name {
	case "claude":
		return accent
	case "codex":
		return violet
	}
	return fg1
}

func agentChip(name string) string {
	tone := agentTone(name)
	return lipgloss.NewStyle().
		Foreground(tone).
		Background(bg2).
		Padding(0, 1).
		Render(name)
}
