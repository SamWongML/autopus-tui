package main

import "github.com/charmbracelet/bubbles/spinner"

// liveSpin is the single bubbles/spinner instance shared across views for
// "live" status pills (daemon online dot in the tab bar, the LIVE pill in the
// Run header). Centralizing it ensures every animated indicator stays in
// frame-lock and a single tea.Cmd advances them all.
//
// View functions read `liveSpin.View()` directly; main.Update forwards
// spinner.TickMsg back into liveSpin and re-emits the resulting Tick command.
// In --dump mode no Update ever runs, so View() emits frame 0 deterministically
// — that's the snapshotted character in the golden files.
var liveSpin spinner.Model

func init() {
	liveSpin = spinner.New(spinner.WithSpinner(spinner.Pulse))
	liveSpin.Style = sOk
}
