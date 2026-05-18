package data

// ActivityItem is one row of the overview's recent-activity stream. Kind is a
// state name OR a generic marker ("warn", "info", "dot").
type ActivityItem struct {
	T        string
	Kind     string
	ID, Body string
}

var ActivityStream = []ActivityItem{
	{"16:42", "needs_input", "ses_7Hq", "asked a question · scheduler refactor approach"},
	{"16:41", "running", "ses_9Lp", "writing migration 011_add_quota.sql"},
	{"16:41", "warn", "fs", "workspace 'autopus-core' at 78% disk quota"},
	{"16:41", "failed", "ses_M4g", "agent timeout · 18m without progress"},
	{"16:41", "completed", "ses_Z1f", "PR opened · autopus#1284 · docs/runtime"},
	{"16:40", "info", "auth", "token refreshed · valid for 86 days"},
	{"16:40", "working", "ses_2Kc", "running tests · ./e2e/heartbeat (4/12)"},
	{"16:39", "running", "ses_7Hq", "edited internal/daemon/runtime.go"},
	{"16:39", "dot", "cfg", "config reload · poll_interval unchanged"},
	{"16:38", "needs_input", "ses_Q3d", "requested permission to /etc/autopus"},
	{"16:37", "running", "ses_9Lp", "ran ./migrate.sh on stage db"},
	{"16:36", "completed", "ses_T2k", "tests passed · 412/412 · 8m32s"},
}
