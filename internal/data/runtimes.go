package data

// Runtime is one detected agent CLI on $PATH.
type Runtime struct {
	CLI       string
	Path      string
	Version   string
	Model     string
	Status    string
	Inflight  int
	Queued    int
	Cap       int
	LastUsed  string
	Cost24h   string
	Tokens24h string
	Success   float64
}

var Runtimes = []Runtime{
	{CLI: "claude", Path: "/opt/homebrew/bin/claude", Version: "1.0.42", Model: "sonnet-4.5", Status: "ready", Inflight: 2, Cap: 8, LastUsed: "12s ago", Cost24h: "$4.18", Tokens24h: "412k", Success: 0.94},
	{CLI: "codex", Path: "/opt/homebrew/bin/codex", Version: "0.18.3", Model: "gpt-5.1", Status: "ready", Inflight: 1, Queued: 1, Cap: 6, LastUsed: "1m ago", Cost24h: "$2.81", Tokens24h: "188k", Success: 0.91},
	{CLI: "gemini", Path: "/opt/homebrew/bin/gemini", Version: "0.4.0", Model: "2.5-pro", Status: "ready", Inflight: 0, Cap: 4, LastUsed: "37m ago", Cost24h: "$0.62", Tokens24h: "48k", Success: 0.88},
	{CLI: "claude-opus", Path: "/opt/homebrew/bin/claude", Version: "1.0.42", Model: "opus-4.1", Status: "ready", Inflight: 1, Cap: 2, LastUsed: "4m ago", Cost24h: "$6.40", Tokens24h: "92k", Success: 0.97},
	{CLI: "copilot", Path: "/usr/local/bin/gh-copilot", Version: "0.9.1", Model: "—", Status: "disabled", Inflight: 0, Cap: 4, LastUsed: "never", Cost24h: "$0.00", Tokens24h: "0", Success: 0},
	{CLI: "kiro-cli", Path: "/Applications/Kiro.app/Contents/MacOS/kiro", Version: "0.2.1", Model: "kiro-coder-1", Status: "ready", Inflight: 0, Cap: 2, LastUsed: "2h ago", Cost24h: "$0.18", Tokens24h: "11k", Success: 0.83},
	{CLI: "opencode", Path: "~/.local/bin/opencode", Version: "0.2.7", Model: "auto", Status: "stale", Inflight: 0, Cap: 2, LastUsed: "8d ago", Cost24h: "$0.00", Tokens24h: "0", Success: 0},
	{CLI: "hermes", Path: "—", Version: "—", Model: "—", Status: "not_found", Inflight: 0, Cap: 0, LastUsed: "never", Cost24h: "$0.00", Tokens24h: "0", Success: 0},
	{CLI: "cursor-agent", Path: "/Applications/Cursor.app/.../cursor-agent", Version: "0.7.4", Model: "anthropic/sonnet-4", Status: "ready", Inflight: 0, Cap: 4, LastUsed: "yesterday", Cost24h: "$0.41", Tokens24h: "31k", Success: 0.90},
}
