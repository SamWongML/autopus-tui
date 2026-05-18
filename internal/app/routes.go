package app

// Route is one entry in the top tab bar — a numbered, jumpable view.
type Route struct {
	ID    string
	Key   string
	Label string
}

// Routes is the canonical tab order. Add a new tab here AND add the matching
// child model in app.go (Model struct + buildModel + activeView + viewFor).
var Routes = []Route{
	{"overview", "1", "Overview"},
	{"sessions", "2", "Sessions"},
	{"issues", "3", "Issues"},
	{"runtimes", "4", "Runtimes"},
	{"workspaces", "5", "Workspaces"},
	{"logs", "6", "Logs"},
	{"config", "7", "Config"},
}

// RouteIdx returns the index of a route by id, or 0 if not found.
func RouteIdx(id string) int {
	for i, r := range Routes {
		if r.ID == id {
			return i
		}
	}
	return 0
}
