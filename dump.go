// Dump-mode renderer: feed a Model the same WindowSizeMsg the runtime would,
// then call View() and print the result. Iterates every route + every overlay
// + the attach view, so a single `go run . --dump` is a full visual snapshot.
//
// Controlled by env vars:
//
//   - AUTOPUS_DUMP_PROFILE = truecolor | 256 | ascii      (default: truecolor)
//   - AUTOPUS_DUMP_WIDTH   = <int>                        (default: 160)
//   - AUTOPUS_DUMP_HEIGHT  = <int>                        (default: 44)
//   - AUTOPUS_DUMP_ONLY    = comma-sep names (filters)    (default: all)
//
// Names recognized by AUTOPUS_DUMP_ONLY: overview, sessions, issues, runtimes,
// workspaces, logs, config, attach, help, palette, onboarding.
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"autopus-tui/internal/app"
	"autopus-tui/internal/data"
	"autopus-tui/internal/tui"
	"autopus-tui/internal/ui"
)

type dumpFrame struct {
	name string
	mut  func(m *tui.Model)
}

func runDump(out io.Writer) {
	setProfile(os.Getenv("AUTOPUS_DUMP_PROFILE"))
	w := envInt("AUTOPUS_DUMP_WIDTH", 160)
	h := envInt("AUTOPUS_DUMP_HEIGHT", 44)
	only := parseOnly(os.Getenv("AUTOPUS_DUMP_ONLY"))

	frames := []dumpFrame{
		// One frame per top-level route.
		{"overview", func(m *tui.Model) { m.Route = "overview" }},
		{"sessions", func(m *tui.Model) { m.Route = "sessions" }},
		{"issues", func(m *tui.Model) { m.Route = "issues" }},
		{"runtimes", func(m *tui.Model) { m.Route = "runtimes" }},
		{"workspaces", func(m *tui.Model) { m.Route = "workspaces" }},
		{"logs", func(m *tui.Model) { m.Route = "logs" }},
		{"config", func(m *tui.Model) { m.Route = "config" }},

		// Attach view: pick the first session for a deterministic snapshot.
		{"attach", func(m *tui.Model) {
			m.Route = "sessions"
			m.Attach = data.Sessions[0].ID
			m.AttachM.ID = m.Attach
		}},

		// Overlays.
		{"help", func(m *tui.Model) { m.Route = "overview"; m.Overlay = "help" }},
		{"palette", func(m *tui.Model) { m.Route = "overview"; m.Overlay = "palette" }},
		{"onboarding", func(m *tui.Model) { m.Overlay = "onboarding" }},
	}

	first := true
	for _, f := range frames {
		if len(only) > 0 && !only[f.name] {
			continue
		}
		if !first {
			fmt.Fprintln(out)
		}
		first = false
		dumpFrameOne(out, f, w, h)
	}
}

func dumpFrameOne(out io.Writer, f dumpFrame, w, h int) {
	m := tui.New()
	m.W, m.H = w, h
	m.Bp = ui.For(w)
	f.mut(&m)

	fmt.Fprintln(out, strings.Repeat("=", w))
	fmt.Fprintf(out, "=== %s · route=%s overlay=%q attach=%q · %dx%d\n",
		f.name, m.Route, m.Overlay, m.Attach, w, h)
	fmt.Fprintln(out, strings.Repeat("=", w))
	fmt.Fprintln(out, m.View())
}

func setProfile(v string) {
	switch v {
	case "", "truecolor":
		lipgloss.SetColorProfile(termenv.TrueColor)
	case "256":
		lipgloss.SetColorProfile(termenv.ANSI256)
	case "ascii":
		lipgloss.SetColorProfile(termenv.Ascii)
	}
}

func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func parseOnly(v string) map[string]bool {
	out := map[string]bool{}
	for _, name := range strings.Split(v, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			out[name] = true
		}
	}
	return out
}

// Guard against an unused import if Init/Update ever shed their use of these
// types — they're load-bearing for the rest of the binary.
var _ tea.Model = tui.Model{}
var _ = app.Routes
