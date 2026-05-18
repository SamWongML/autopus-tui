package data

// LogLine is one structured daemon log entry.
type LogLine struct {
	T, Level, Src, Msg string
}

var LogLines = []LogLine{
	{"16:42:18", "INFO", "daemon", "heartbeat ok · server=app.autopus.ai latency=38ms"},
	{"16:42:14", "INFO", "scheduler", "session ses_7Hq awaiting input · idle 02:11"},
	{"16:42:09", "DEBUG", "runtime", "claude poll · 1 task claimed (APO-322)"},
	{"16:42:02", "INFO", "runtime", "codex started task t_91ad on session ses_9Lp"},
	{"16:41:58", "WARN", "fs", "workspace 'autopus-core' at 78% of disk quota (2.4G / 3.0G)"},
	{"16:41:50", "INFO", "daemon", "heartbeat ok · server=app.autopus.ai latency=41ms"},
	{"16:41:44", "DEBUG", "scheduler", "tick · 3 active · 1 idle · 0 queued"},
	{"16:41:31", "INFO", "runtime", "claude task t_8fa1 streamed 412 tokens"},
	{"16:41:22", "ERROR", "runtime", "codex task t_7bb2 exited 124 (agent timeout · ses_M4g)"},
	{"16:41:10", "INFO", "daemon", "deregistered runtime opencode (stale, last seen 8d ago)"},
	{"16:41:02", "INFO", "auth", "token refreshed · valid for 86 days"},
	{"16:40:55", "DEBUG", "scheduler", "tick · 4 active · 1 idle · 0 queued"},
	{"16:40:46", "INFO", "runtime", "claude poll · 0 tasks"},
	{"16:40:32", "INFO", "daemon", "heartbeat ok · server=app.autopus.ai latency=36ms"},
	{"16:40:18", "INFO", "fs", "session ses_Z1f closed · cleaned 84M from workspace"},
	{"16:40:02", "DEBUG", "runtime", "claude poll · 1 task claimed (APO-318)"},
	{"16:39:48", "INFO", "daemon", "heartbeat ok · server=app.autopus.ai latency=39ms"},
	{"16:39:22", "INFO", "runtime", "codex task t_6cc0 finished · 8m12s · tokens 4180/612"},
	{"16:39:14", "INFO", "daemon", "config reload · poll_interval 3s → 3s (unchanged)"},
	{"16:38:55", "DEBUG", "scheduler", "tick · 5 active · 0 idle · 0 queued"},
}

// LogLevelFilters are the pill labels in order; "all" disables filtering.
var LogLevelFilters = []string{"all", "INFO", "DEBUG", "WARN", "ERROR"}
