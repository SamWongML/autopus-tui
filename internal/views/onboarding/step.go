package onboarding

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderStep(m Model, w, h int) string {
	cur := data.OnbSteps[m.Step]
	right := ""
	if m.Step > 0 {
		right += ui.KeyChip("←", "back", false) + "  "
	}
	if m.Step < len(data.OnbSteps)-1 {
		right += ui.KeyChip("→", "next", true)
	} else {
		right += ui.KeyChip("↵", "launch daemon", true)
	}

	var body string
	switch cur.ID {
	case "server":
		body = renderServer(m, w-4)
	case "auth":
		body = renderAuth(w - 4)
	case "workspaces":
		body = renderWorkspaces(m, w-4)
	case "runtimes":
		body = renderRuntimes(w - 4)
	case "daemon":
		body = renderDaemon(w - 4)
	case "review":
		body = renderReview(w - 4)
	}

	title := "step " + ui.Itoa(m.Step+1) + " / " + ui.Itoa(len(data.OnbSteps)) + "  ·  " + cur.Title
	return ui.Panel(title, right, body, w, h, false, false)
}

func renderServer(m Model, w int) string {
	type opt struct {
		id, title, sub, env string
	}
	opts := []opt{
		{"cloud", "Autopus Cloud", "app.autopus.ai · default", ""},
		{"self", "Self-hosted", "point at your own server", "AUTOPUS_APP_URL=https://app.acme.dev\nAUTOPUS_SERVER_URL=wss://api.acme.dev/ws"},
		{"air", "Air-gapped (local)", "no remote · daemon-only", "AUTOPUS_OFFLINE=1"},
	}
	var b strings.Builder
	b.WriteString(theme.SFaint.Render("Pick one. You can change this any time with ") +
		theme.SAccent.Render("autopus config set") + theme.SFaint.Render(".") + "\n\n")
	for _, o := range opts {
		selected := o.id == m.Server
		marker := "○"
		titleCol, borderCol := theme.Text, theme.Border
		if selected {
			marker = "◉"
			titleCol, borderCol = theme.Accent, theme.AccentDim
		}
		content := theme.Fg(titleCol).Render(marker+" "+o.title) + "\n" +
			theme.SFaint.Render("  "+o.sub)
		if o.env != "" {
			content += "\n" + theme.SDim.Render("  "+strings.ReplaceAll(o.env, "\n", "\n  "))
		}
		styled := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(borderCol)).
			Padding(0, 1).Width(w-2).Render(content)
		b.WriteString(styled + "\n")
	}
	return b.String()
}

func renderAuth(w int) string {
	var b strings.Builder
	b.WriteString(theme.SFaint.Render("Two ways to authenticate. The daemon creates a 90-day token either way.") + "\n\n")
	cardW := (w - 3) / 2
	browser := theme.SAccent.Render("↗ Open browser") + "\n\n" +
		theme.SDim.Render("opens ") + theme.SText.Render("app.autopus.ai/cli/auth") + theme.SDim.Render(" with a one-time code.") + "\n" +
		theme.SFaint.Render("Recommended.")
	paste := theme.SText.Render("⎘ Paste token") + "\n\n" +
		theme.SDim.Render("for headless / remote machines. Create one at") + "\n" +
		theme.SAccent.Render("app.autopus.ai/settings/tokens") + "\n\n" +
		lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.Border)).
			Padding(0, 1).Render(theme.SAccent.Render("›")+" apo_••••••••••••"+theme.SAccent.Render("▌"))
	left := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.AccentDim)).Padding(0, 1).Width(cardW).Render(browser)
	right := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.Border)).Padding(0, 1).Width(cardW).Render(paste)
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right))
	return b.String()
}

func renderWorkspaces(m Model, w int) string {
	var b strings.Builder
	b.WriteString(theme.SFaint.Render("Pick the workspaces this machine should handle.") + "\n\n")
	const cBox = 4
	const cName = 22
	const cMem = 10
	const cIss = 12
	cRole := w - cBox - cName - cMem - cIss - 5
	if cRole < 6 {
		cRole = 6
	}
	for _, x := range data.Workspaces {
		on := m.Watched[x.ID]
		var box string
		col := theme.Text
		if on {
			box = theme.SAccent.Render("[x]")
			col = theme.Accent
		} else {
			box = theme.SMute.Render("[ ]")
		}
		line := box + " " +
			theme.Fg(col).Render(ui.PadRight(x.Name, cName)) + " " +
			theme.SFaint.Render(ui.PadRight(ui.Itoa(x.Members)+" ppl", cMem)) + " " +
			theme.SFaint.Render(ui.PadRight(ui.Itoa(x.Issues)+" issues", cIss)) + " " +
			theme.SFaint.Render(ui.PadRight(x.Role, cRole))
		if on {
			line = theme.WithBg(line, theme.AccentFaint)
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + theme.SFaint.Render("space toggle · a select all · A clear all"))
	return b.String()
}

func renderRuntimes(w int) string {
	var b strings.Builder
	b.WriteString(theme.SFaint.Render("Detected agent CLIs on ") + theme.SText.Render("$PATH") + theme.SFaint.Render(".") + "\n\n")
	const cBox = 4
	const cCli = 12
	const cVer = 10
	const cModel = 14
	const cStat = 10
	cPath := w - cBox - cCli - cVer - cModel - cStat - 5
	if cPath < 10 {
		cPath = 10
	}
	for _, r := range data.Runtimes {
		ok := r.Status == "ready"
		var box string
		col := theme.TextDim
		if ok {
			box = theme.SOK.Render("[x]")
			col = theme.Text
		} else {
			box = theme.SMute.Render("[ ]")
		}
		statusCol := theme.TextFaint
		if ok {
			statusCol = theme.OK
		}
		line := box + " " +
			theme.Fg(col).Render(ui.PadRight(r.CLI, cCli)) + " " +
			theme.SFaint.Render(ui.PadRight(r.Version, cVer)) + " " +
			theme.SFaint.Render(ui.PadRight(r.Path, cPath)) + " " +
			theme.SDim.Render(ui.PadRight(r.Model, cModel)) + " " +
			theme.Fg(statusCol).Render(ui.PadRight(r.Status, cStat))
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + theme.SFaint.Render("Pin custom paths via ") +
		theme.SAccent.Render("AUTOPUS_<CLI>_PATH") + theme.SFaint.Render(" / ") +
		theme.SAccent.Render("AUTOPUS_<CLI>_MODEL") + theme.SFaint.Render("."))
	return b.String()
}

func renderDaemon(w int) string {
	type f struct{ k, v, desc string }
	fields := []f{
		{"poll_interval", "3s", "how often to ask the server for new tasks"},
		{"heartbeat_interval", "15s", "how often this daemon pings the server"},
		{"agent_timeout", "2h", "kill an agent that goes silent this long"},
		{"max_concurrent_tasks", "20", "ceiling across all runtimes"},
	}
	tuiFields := [][2]string{
		{"notify_on_needs_input", "desktop + sound · 10s"},
		{"autostart", "on login (launchd)"},
		{"default view", "Sessions"},
		{"theme", "warm dark"},
	}
	var b strings.Builder
	b.WriteString(theme.SAccent.Render("DAEMON") + "\n")
	for _, fi := range fields {
		line := theme.SText.Render(ui.PadRight(fi.k, 22)) + " " +
			theme.SFaint.Render(ui.PadRight(fi.desc, w-22-1-10)) + " " +
			theme.SAccent.Render(ui.PadLeft(fi.v, 9))
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + theme.SAccent.Render("THIS TUI · COMPLEMENTARY SETTINGS") + "\n")
	for _, t := range tuiFields {
		line := theme.SText.Render(ui.PadRight(t[0], 22)) + " " + theme.SDim.Render(t[1])
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + theme.SFaint.Render("← prev field · → next · tab edit · ↵ next step"))
	return b.String()
}

func renderReview(w int) string {
	var b strings.Builder
	b.WriteString(theme.SFaint.Render("About to write ") +
		theme.SText.Render("~/.autopus/profiles/default/config.toml") +
		theme.SFaint.Render(" and start the daemon.") + "\n\n")
	toml := `# autopus · profile=default
[server]
  app_url    = "https://app.autopus.ai"
  server_url = "wss://api.autopus.ai/ws"

[daemon]
  poll_interval        = "3s"
  heartbeat_interval   = "15s"
  agent_timeout        = "2h"
  max_concurrent_tasks = 20
  workspaces_root      = "~/autopus_workspaces"

[watch]
  workspaces = ["autopus-core", "autopus-platform", "autopus-docs"]

[agents.claude]
  path  = "/opt/homebrew/bin/claude"
  model = "sonnet-4.5"

[agents.codex]
  path  = "/opt/homebrew/bin/codex"
  model = "gpt-5.1"

[tui]
  default_view          = "sessions"
  notify_on_needs_input = "desktop+sound:10s"
  autostart             = "launchd"
  theme                 = "warm-dark"`
	tomlStyled := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Surface)).
		Foreground(lipgloss.Color(theme.TextDim)).
		Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(theme.Border)).
		Padding(0, 1).Width(w-2).Render(toml)
	b.WriteString(tomlStyled + "\n\n")
	b.WriteString(ui.KeyChip("↵", "write & launch", true) + "  " +
		ui.KeyChip("e", "edit toml", false) + "  " +
		ui.KeyChip("←", "back", false))
	return b.String()
}
