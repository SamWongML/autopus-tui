package main

// Sample data ported from artboards.jsx. All values match the design so the
// screens render the same content as the mock.

type Daemon struct {
	Profile, Server, Socket, Log, WsRoot, Device, Version string
	PID                                                   int
	Uptime                                                string
	Connected                                             bool
	MemMB, MemMax                                         int
	CPU                                                   float64
	LastPoll, LastHB                                      string
	PollsToday, HeartbeatsToday                           int
}

var daemon = Daemon{
	Profile: "default", PID: 14882, Uptime: "14h 22m 09s", Version: "0.42.1",
	Server: "wss://api.multica.ai", Connected: true,
	Socket: "~/.multica/daemon.sock", Log: "~/.multica/daemon.log",
	WsRoot: "~/multica_workspaces", Device: "andriis-mbp",
	MemMB: 186, MemMax: 512, CPU: 3.2,
	LastPoll: "0.8s ago", LastHB: "4s ago",
	PollsToday: 28412, HeartbeatsToday: 3369,
}

type Runtime struct {
	Name, Bin, Ver, Model string
	Busy, Max             int
	TasksToday, ErrsToday int
}

var runtimes = []Runtime{
	{"claude", "/opt/homebrew/bin/claude", "2.1.142", "sonnet-4.6", 2, 3, 47, 1},
	{"codex", "/opt/homebrew/bin/codex", "0.34.0", "gpt-5.5", 1, 4, 19, 0},
}

type Task struct {
	ID, Status, Runtime, WS, Issue, Started, Last, Title, Cost string
	Seq                                                        int
}

var tasks = []Task{
	{"t-1284", "working", "claude", "blackhole-os", "#4127", "3m 12s", "editing app/models/webhook_event.rb",
		"refactor billing webhook → idempotency keys", "$1.42", 84},
	{"t-1283", "waiting", "claude", "blackhole-os", "#4126", "12s", "awaiting confirm: drop unique index?",
		"migrate analytics events to ClickHouse", "$0.41", 12},
	{"t-1282", "working", "codex", "blackhole-os", "#4118", "1m 04s", "rspec iteration 27 / 40",
		"investigate flaky PaymentSheet integration", "$0.18", 71},
	{"t-1281", "done", "claude", "blackhole-os", "#4115", "8m", "PR #4130 opened",
		"generate openapi client for billing-v2", "$0.92", 142},
	{"t-1280", "failed", "codex", "blackhole-os", "#4114", "21m", "exit 1: bundler resolver",
		"update stripe sdk to 17.0.0", "$0.05", 23},
	{"t-1279", "queued", "claude", "blackhole-os", "#4113", "—", "waiting on runtime claude (3/3)",
		"implement webhook retry exponential backoff", "—", 0},
}

type Msg struct {
	Seq                          int
	T, Type, Body, Tool, Args    string
	Lines                        int
	Ok                           bool
}

var msgs = []Msg{
	{Seq: 78, T: "14:31:02", Type: "text",
		Body: "Plan: read handler + spec; add webhook_events table; wrap with idempotency lookup behind feature flag billing_v2_idempotency."},
	{Seq: 79, T: "14:31:11", Type: "thinking",
		Body: "The existing migration already creates webhook_events with (provider, external_id) unique. Composite needs the new idempotency_key column. Probably safest to keep old index renamed during cutover."},
	{Seq: 80, T: "14:31:18", Type: "tool_call", Tool: "bash", Args: "rails db:migrate:status | tail -3"},
	{Seq: 81, T: "14:31:19", Type: "tool_result",
		Body: "up    20260512_add_org_seats\nup    20260513_create_webhook_events"},
	{Seq: 82, T: "14:31:24", Type: "tool_call", Tool: "edit", Args: "db/migrate/20260514_add_idempotency_key.rb", Lines: 42},
	{Seq: 83, T: "14:31:41", Type: "tool_call", Tool: "bash", Args: "bundle exec rspec spec/models/webhook_event_spec.rb"},
	{Seq: 84, T: "14:31:58", Type: "tool_result", Body: "14 examples, 0 failures · 0.84s", Ok: true},
	{Seq: 85, T: "14:32:02", Type: "thinking",
		Body: "Reviewing the conflict: the existing unique index on (provider, external_id) will collide with the new composite. Either drop it pre-cutover or rename for safety."},
	{Seq: 86, T: "14:32:09", Type: "text",
		Body: "I drafted the migration but the webhook_events table already has a unique index on (provider, external_id). Should I drop it before adding the new composite key, or rename it for a safer cutover?"},
}

type LogLine struct {
	T, Lvl, Src, Msg string
}

// 24h throughput series shown on the Status pulse pane.
var tasksPerHour = []int{2, 1, 0, 0, 1, 2, 3, 4, 5, 4, 6, 7, 6, 5, 7, 8, 6, 9, 10, 8, 6, 5, 4, 3}
var errsPerHour = []int{0, 0, 0, 0, 0, 1, 0, 1, 2, 1, 0, 1, 0, 2, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1}

var logLines = []LogLine{
	{"14:32:14", "info", "runtime.claude", "task t-1284 message seq=84 · result · 1 ms"},
	{"14:32:14", "trace", "poll", "server tick → 0 new, 1 updated, 0 cancelled"},
	{"14:32:11", "info", "runtime.claude", "task t-1284 spawn child pid=18421 cwd=~/multica_workspaces/blackhole-os/t-1284"},
	{"14:32:11", "trace", "poll", "server tick → 1 new (t-1284)"},
	{"14:32:09", "warn", "runtime.codex", `task t-1280 exit=1 stderr_tail="Bundler::VersionConflict: stripe-ruby 17.0"`},
	{"14:32:08", "info", "runtime.codex", "task t-1280 message seq=23 · error"},
	{"14:32:01", "trace", "heartbeat", "pulse → ok (server age 12ms)"},
	{"14:31:58", "info", "runtime.claude", "task t-1284 message seq=83 · tool_result · 17 ms"},
	{"14:31:53", "trace", "poll", "server tick → 0 new"},
	{"14:31:46", "trace", "heartbeat", "pulse → ok (server age 9ms)"},
	{"14:31:42", "info", "runtime.codex", "task t-1282 message seq=71 · tool_call rspec"},
	{"14:31:38", "trace", "poll", "server tick → 0 new"},
}

type Profile struct {
	Name, State, Server, WS, Tasks, Uptime string
	PID, Port                              int
	Runtimes                               []string
}

var profiles = []Profile{
	{Name: "default", State: "running", PID: 14882, Server: "wss://api.multica.ai", Port: 7717,
		WS: "~/multica_workspaces", Runtimes: []string{"claude", "codex"}, Tasks: "3/20", Uptime: "14h 22m"},
	{Name: "staging", State: "running", PID: 14903, Server: "wss://staging.multica.ai", Port: 7718,
		WS: "~/.multica/profiles/staging/ws", Runtimes: []string{"claude"}, Tasks: "0/8", Uptime: "02h 11m"},
	{Name: "selfhost", State: "stopped", PID: 0, Server: "ws://localhost:8080", Port: 7719,
		WS: "~/.multica/profiles/selfhost/ws", Runtimes: []string{"claude", "codex"}, Tasks: "—", Uptime: "—"},
}

// CfgRow values for the Config view.
type CfgRow struct {
	K, V, Env, Hint, Tone string
	Dirty, Readonly       bool
}

var cfgDaemon = []CfgRow{
	{K: "poll-interval", V: "3s", Env: "MULTICA_POLL_INTERVAL", Hint: "server poll cadence (≥ 1s)"},
	{K: "heartbeat-interval", V: "15s", Env: "MULTICA_HEARTBEAT_INTERVAL", Hint: "runtime pulse to server"},
	{K: "agent-timeout", V: "2h", Env: "MULTICA_AGENT_TIMEOUT", Hint: "per-task wall-clock cap"},
	{K: "max-concurrent-tasks", V: "20", Env: "MULTICA_MAX_CONCURRENT_TASKS", Hint: "across all runtimes", Dirty: true},
	{K: "workspaces-root", V: "~/multica_workspaces", Env: "MULTICA_WORKSPACES_ROOT", Hint: "isolated per-task cwd"},
	{K: "device-name", V: "andriis-mbp", Hint: "shows in Settings → Runtimes"},
	{K: "daemon-id", V: "d-mac-andriis-7717", Hint: "globally unique", Readonly: true},
	{K: "server-url", V: "wss://api.multica.ai", Env: "MULTICA_SERVER_URL"},
	{K: "health-port", V: "7717", Env: "MULTICA_HEALTH_PORT", Hint: "loopback only"},
	{K: "auto-update", V: "enabled", Hint: "check brew + pull tap weekly", Dirty: true},
}

var cfgClaude = []CfgRow{
	{K: "bin", V: "/opt/homebrew/bin/claude"},
	{K: "model", V: "sonnet-4.6", Env: "MULTICA_CLAUDE_MODEL"},
	{K: "max-concurrent", V: "3"},
	{K: "extra-args", V: "--dangerously-skip-permissions", Tone: "warn"},
}

var cfgCodex = []CfgRow{
	{K: "bin", V: "/opt/homebrew/bin/codex"},
	{K: "model", V: "gpt-5.5", Env: "MULTICA_CODEX_MODEL"},
	{K: "max-concurrent", V: "4"},
	{K: "extra-args", V: ""},
}

var cfgLogging = []CfgRow{
	{K: "level", V: "info", Hint: "trace|debug|info|warn|error"},
	{K: "file", V: "~/.multica/daemon.log"},
	{K: "rotate-at", V: "50 MB"},
	{K: "keep", V: "5 files"},
	{K: "redact-secrets", V: "on", Tone: "ok"},
}
