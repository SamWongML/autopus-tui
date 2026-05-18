// Package data holds the static (mocked) entities the views render. In a real
// build these would be loaded from ~/.autopus/profiles/<name>/state.db or
// streamed from the daemon socket. They live here so views can import a
// stable, dependency-free types layer without coupling to a data fetcher.
package data

// Session is a running or completed agent run on a single issue.
type Session struct {
	ID        string
	Issue     string
	Title     string
	State     string
	Agent     string
	Model     string
	Workspace string
	Branch    string
	Started   string
	Elapsed   string
	TokensIn  int
	TokensOut int
	Cost      string
	Question  string
	Activity  string
	Log       [][3]string // kind, label, body
}

// Sessions is the mock roster.
var Sessions = []Session{
	{
		ID: "ses_7Hq", Issue: "APO-318",
		Title:     "Refactor task scheduler to use cooperative cancellation",
		State:     "needs_input", Agent: "claude", Model: "sonnet-4.5",
		Workspace: "autopus-core", Branch: "feat/coop-cancel",
		Started:   "16:02", Elapsed: "00:39:14",
		TokensIn:  84120, TokensOut: 12840, Cost: "$1.34",
		Question:  "Two valid approaches surfaced. Should I (a) refactor the dispatcher first, then teardown — or (b) add the context.CancelCause shim first and migrate per-runtime? (a) is cleaner; (b) ships faster.",
		Activity:  "asked a question · waiting 02:11",
		Log: [][3]string{
			{"tool", "Read", "internal/daemon/scheduler.go"},
			{"text", "thinking", "The current ctx is propagated through three layers without a cancel reason…"},
			{"tool", "Grep", "ctx.Err in 14 files"},
			{"tool", "Edit", "internal/daemon/runtime.go (+38 −12)"},
			{"text", "thinking", "I should confirm with the user before touching the dispatcher signature — both approaches viable."},
			{"ask", "question", "Two valid approaches surfaced…"},
		},
	},
	{
		ID: "ses_2Kc", Issue: "APO-322",
		Title:     "Investigate flaky e2e test in heartbeat integration",
		State:     "working", Agent: "claude", Model: "sonnet-4.5",
		Workspace: "autopus-core", Branch: "fix/heartbeat-flake",
		Started:   "16:11", Elapsed: "00:30:48",
		TokensIn:  41220, TokensOut: 6140, Cost: "$0.62",
		Activity:  "running tests · ./e2e/heartbeat (4/12)",
		Log: [][3]string{
			{"tool", "Read", "e2e/heartbeat_test.go"},
			{"tool", "Bash", "go test ./e2e/heartbeat -run TestRecover -count=20"},
			{"text", "thinking", "Got 3/20 failures. Looks like a race on the deregister path…"},
			{"tool", "Bash", "go test ./e2e/heartbeat -race -count=20"},
		},
	},
	{
		ID: "ses_9Lp", Issue: "APO-310",
		Title:     "Add disk quota enforcement per workspace",
		State:     "running", Agent: "codex", Model: "gpt-5.1",
		Workspace: "autopus-core", Branch: "feat/disk-quota",
		Started:   "15:18", Elapsed: "01:24:02",
		TokensIn:  162400, TokensOut: 24810, Cost: "$2.18",
		Activity:  "writing migration · 011_add_quota.sql",
	},
	{
		ID: "ses_Q3d", Issue: "APO-301",
		Title:     "Wire OTLP exporter behind feature flag",
		State:     "needs_input", Agent: "claude", Model: "opus-4.1",
		Workspace: "autopus-platform", Branch: "feat/otlp",
		Started:   "15:54", Elapsed: "00:48:09",
		TokensIn:  58900, TokensOut: 9440, Cost: "$1.88",
		Question:  "Permission requested: write to /etc/autopus/otlp.yaml — outside the workspace sandbox. Allow once / always / deny?",
		Activity:  "permission prompt · 00:31",
	},
	{
		ID: "ses_Z1f", Issue: "APO-298",
		Title:     "Update docs: agent runtime configuration",
		State:     "completed", Agent: "claude", Model: "sonnet-4.5",
		Workspace: "autopus-docs", Branch: "docs/runtime",
		Started:   "14:30", Elapsed: "00:42:11",
		TokensIn:  28400, TokensOut: 18220, Cost: "$0.78",
		Activity:  "PR opened · autopus-ai/autopus#1284",
	},
	{
		ID: "ses_M4g", Issue: "APO-289",
		Title:     "Investigate memory leak in long-running runs",
		State:     "failed", Agent: "codex", Model: "gpt-5.1",
		Workspace: "autopus-core", Branch: "fix/mem-leak",
		Started:   "13:11", Elapsed: "00:18:42",
		TokensIn:  12400, TokensOut: 1980, Cost: "$0.14",
		Activity:  "agent timeout · 18m without progress",
	},
	{
		ID: "ses_T8s", Issue: "APO-330",
		Title:     "Sketch: cooperative cache eviction for runtime",
		State:     "idle", Agent: "claude", Model: "haiku-4.5",
		Workspace: "autopus-core", Branch: "exp/cache-evict",
		Started:   "12:00", Elapsed: "00:00:00",
		TokensIn:  0, TokensOut: 0, Cost: "$0.00",
		Activity:  "queued · waiting for runtime slot",
	},
	{
		ID: "ses_R5n", Issue: "APO-275",
		Title:     "Render flame graph in dashboard",
		State:     "paused", Agent: "claude", Model: "sonnet-4.5",
		Workspace: "autopus-platform", Branch: "feat/flamegraph",
		Started:   "11:24", Elapsed: "02:08:30",
		TokensIn:  96440, TokensOut: 11300, Cost: "$1.42",
		Activity:  "backgrounded by you · 01:14:08 ago",
	},
}

// SessionFilters are the pill labels in order, including "all".
var SessionFilters = []string{"all", "needs_input", "working", "running", "paused", "idle", "completed", "failed"}

// FindSession returns the session matching id, falling back to the first
// session if none match. Returning a pointer lets attach views reflect updates.
func FindSession(id string) *Session {
	for i := range Sessions {
		if Sessions[i].ID == id {
			return &Sessions[i]
		}
	}
	return &Sessions[0]
}

// CountSessionStates returns a map[state]count over all sessions.
func CountSessionStates() map[string]int {
	out := map[string]int{"all": len(Sessions)}
	for _, s := range Sessions {
		out[s.State]++
	}
	return out
}

// TranscriptEntry is one row of the attach view's chronological transcript.
type TranscriptEntry struct {
	Kind, Time, Body, Tool, Arg string
}

var SessionTranscript = []TranscriptEntry{
	{Kind: "user", Time: "16:02:14", Body: "Refactor the task scheduler to use cooperative cancellation — keep the existing tests green."},
	{Kind: "plan", Time: "16:02:18", Body: "Plan:\n  1. Map the current cancellation surface\n  2. Pick between context.CancelCause vs a per-runtime shim\n  3. Migrate dispatcher → runtimes → tasks\n  4. Keep wire compat for heartbeat"},
	{Kind: "tool", Time: "16:02:21", Tool: "Read", Arg: "internal/daemon/scheduler.go"},
	{Kind: "tool", Time: "16:02:24", Tool: "Read", Arg: "internal/daemon/runtime.go"},
	{Kind: "tool", Time: "16:02:30", Tool: "Grep", Arg: "ctx.Err — 14 matches across 9 files"},
	{Kind: "thinking", Time: "16:02:38", Body: "Three layers propagate ctx without a cancel reason. The wire format includes a status int, so I can carry the cause without breaking heartbeat compat."},
	{Kind: "tool", Time: "16:02:52", Tool: "Edit", Arg: "internal/daemon/runtime.go (+38 −12)"},
	{Kind: "tool", Time: "16:03:08", Tool: "Bash", Arg: "go test ./internal/daemon/... -run TestCancel"},
	{Kind: "tool", Time: "16:03:18", Tool: "Bash", Arg: "go vet ./..."},
	{Kind: "thinking", Time: "16:03:24", Body: "Cleaner to refactor the dispatcher first — but that's a wider change. Shipping the shim first means we get observability now and can migrate per-runtime."},
	{Kind: "ask", Time: "16:39:10", Body: "Two valid approaches surfaced. Should I:\n  (a) refactor the dispatcher first, then teardown — cleaner, ~600 LOC.\n  (b) add the context.CancelCause shim first and migrate per-runtime — ships faster, less risk."},
}
