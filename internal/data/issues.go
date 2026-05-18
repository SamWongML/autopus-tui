package data

// Issue is a single tracked piece of work.
type Issue struct {
	ID        string
	Title     string
	Status    string
	Priority  string
	Assignee  string
	Workspace string
	Updated   string
}

var Issues = []Issue{
	{ID: "APO-318", Title: "Refactor task scheduler to use cooperative cancellation", Status: "in_progress", Priority: "high", Assignee: "Lambda (claude)", Workspace: "autopus-core", Updated: "2m"},
	{ID: "APO-322", Title: "Investigate flaky e2e test in heartbeat integration", Status: "in_progress", Priority: "urgent", Assignee: "Sigma (claude)", Workspace: "autopus-core", Updated: "5m"},
	{ID: "APO-310", Title: "Add disk quota enforcement per workspace", Status: "in_progress", Priority: "high", Assignee: "Theta (codex)", Workspace: "autopus-core", Updated: "12m"},
	{ID: "APO-301", Title: "Wire OTLP exporter behind feature flag", Status: "in_review", Priority: "medium", Assignee: "Lambda (claude)", Workspace: "autopus-platform", Updated: "31m"},
	{ID: "APO-298", Title: "Update docs: agent runtime configuration", Status: "done", Priority: "low", Assignee: "Lambda (claude)", Workspace: "autopus-docs", Updated: "1h"},
	{ID: "APO-289", Title: "Investigate memory leak in long-running runs", Status: "blocked", Priority: "urgent", Assignee: "Theta (codex)", Workspace: "autopus-core", Updated: "2h"},
	{ID: "APO-330", Title: "Sketch: cooperative cache eviction for runtime", Status: "todo", Priority: "medium", Assignee: "—", Workspace: "autopus-core", Updated: "3h"},
	{ID: "APO-275", Title: "Render flame graph in dashboard", Status: "in_progress", Priority: "medium", Assignee: "Lambda (claude)", Workspace: "autopus-platform", Updated: "1d"},
	{ID: "APO-340", Title: "Switch default model envelope to thinking-1m", Status: "backlog", Priority: "low", Assignee: "—", Workspace: "autopus-platform", Updated: "2d"},
	{ID: "APO-285", Title: "Cancel button should propagate to all in-flight subtasks", Status: "todo", Priority: "high", Assignee: "—", Workspace: "autopus-core", Updated: "2d"},
}

// IssueFilters are the pill labels in order; "active" matches todo+in_progress+
// in_review+blocked, anything else is an exact match (or "all").
var IssueFilters = []string{"all", "active", "todo", "in_progress", "in_review", "blocked", "done", "backlog"}
