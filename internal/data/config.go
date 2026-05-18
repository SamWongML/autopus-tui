package data

import "autopus-tui/internal/theme"

// CfgRow is one row of the config view. Tag is the badge ("ENV", "FLAG·ENV",
// "TUI") shown after the key. Color overrides the value color when set.
type CfgRow struct {
	K, V, Tag string
	Color     string
}

var CfgServer = []CfgRow{
	{K: "app_url", V: "https://app.autopus.ai"},
	{K: "server_url", V: "wss://api.autopus.ai/ws"},
	{K: "user", V: "@you · org-acme", Color: theme.Accent},
	{K: "token", V: "apo_•••2f8a · 86d left"},
}

var CfgDaemon = []CfgRow{
	{K: "poll_interval", V: "3s", Tag: "FLAG·ENV"},
	{K: "heartbeat_interval", V: "15s", Tag: "FLAG·ENV"},
	{K: "agent_timeout", V: "2h", Tag: "FLAG·ENV"},
	{K: "max_concurrent_tasks", V: "20", Tag: "FLAG·ENV"},
	{K: "daemon_id", V: "luna.local"},
	{K: "device_name", V: "luna.local"},
	{K: "runtime_name", V: "Local Agent"},
	{K: "workspaces_root", V: "~/autopus_workspaces"},
	{K: "log_level", V: "info", Tag: "FLAG·ENV"},
	{K: "log_path", V: "~/.autopus/profiles/default/daemon.log"},
}

var CfgTUI = []CfgRow{
	{K: "default_view", V: "overview", Tag: "TUI"},
	{K: "notify_on_needs_input", V: "desktop + sound · 10s", Tag: "TUI"},
	{K: "autostart", V: "on login (launchd)", Tag: "TUI"},
	{K: "theme", V: "warm-dark", Tag: "TUI"},
	{K: "follow_focus_on_alert", V: "true", Tag: "TUI"},
	{K: "confirm_destructive", V: "true", Tag: "TUI"},
	{K: "show_thinking", V: "collapsed", Tag: "TUI"},
	{K: "transcript_density", V: "comfortable", Tag: "TUI"},
}

// AgentCfg is one per-agent override block in the config view.
type AgentCfg struct {
	Name        string
	Path        string
	Model       string
	Env         [][2]string
	Concurrency int
}

var AgentCfgs = []AgentCfg{
	{Name: "claude", Path: "/opt/homebrew/bin/claude", Model: "sonnet-4.5", Concurrency: 8, Env: [][2]string{
		{"AUTOPUS_CLAUDE_PATH", "(default)"}, {"AUTOPUS_CLAUDE_MODEL", "(default)"},
	}},
	{Name: "codex", Path: "/opt/homebrew/bin/codex", Model: "gpt-5.1", Concurrency: 6, Env: [][2]string{
		{"AUTOPUS_CODEX_PATH", "(default)"}, {"AUTOPUS_CODEX_MODEL", "(default)"},
	}},
	{Name: "gemini", Path: "/opt/homebrew/bin/gemini", Model: "2.5-pro", Concurrency: 4},
	{Name: "claude-opus", Path: "/opt/homebrew/bin/claude", Model: "opus-4.1", Concurrency: 2, Env: [][2]string{
		{"AUTOPUS_CLAUDE_MODEL", "opus-4.1 (pinned)"},
	}},
}

// Budgets is the daily/per-run spend cap block.
var Budgets = struct {
	DailyCap, DailyUsed, PerRunCap float64
	WarnPct                        int
}{DailyCap: 50.00, DailyUsed: 14.20, PerRunCap: 5.00, WarnPct: 80}
