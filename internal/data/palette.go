package data

// PaletteItem is one entry in the `:` command palette.
type PaletteItem struct {
	K, Label, Kind string
}

var PaletteItems = []PaletteItem{
	{"replay onboarding", "Re-open the first-run setup wizard", "tui"},
	{"login", "Authenticate via browser", "auth"},
	{"logout", "Sign out", "auth"},
	{"daemon start", "Start the agent daemon", "daemon"},
	{"daemon stop", "Stop the agent daemon", "daemon"},
	{"daemon restart", "Restart the agent daemon", "daemon"},
	{"daemon logs -f", "Tail daemon logs", "daemon"},
	{"profile default", "Switch to profile · default", "profile"},
	{"profile staging", "Switch to profile · staging", "profile"},
	{"new issue", "Compose a new issue", "issue"},
	{"watch autopus-platform", "Watch workspace", "ws"},
	{"open app.autopus.ai/issues", "Open in browser", "ext"},
	{"theme cycle", "Cycle theme · warm dark · paper · mono", "tui"},
	{"export config", "Export profile config to ./autopus.toml", "cfg"},
	{"rescan agents", "Rescan $PATH for agent CLIs", "daemon"},
}
