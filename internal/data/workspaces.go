package data

// Workspace is one project/repo the daemon may watch.
type Workspace struct {
	ID       string
	Name     string
	Members  int
	Watch    bool
	Issues   int
	Sessions int
	Role     string
}

var Workspaces = []Workspace{
	{ID: "ws_core", Name: "autopus-core", Members: 7, Watch: true, Issues: 38, Sessions: 4, Role: "Maintainer"},
	{ID: "ws_platform", Name: "autopus-platform", Members: 4, Watch: true, Issues: 21, Sessions: 3, Role: "Maintainer"},
	{ID: "ws_docs", Name: "autopus-docs", Members: 12, Watch: true, Issues: 9, Sessions: 1, Role: "Editor"},
	{ID: "ws_pinpoint", Name: "pinpoint-eng", Members: 3, Watch: false, Issues: 6, Sessions: 0, Role: "Viewer"},
	{ID: "ws_org", Name: "org-internal", Members: 22, Watch: false, Issues: 14, Sessions: 0, Role: "Member"},
}
