package data

// OnbStep is one step of the first-run wizard.
type OnbStep struct {
	ID, Title, Sub string
}

var OnbSteps = []OnbStep{
	{"server", "Server", "Where does your daemon talk to?"},
	{"auth", "Authenticate", "Sign in to discover your workspaces"},
	{"workspaces", "Workspaces", "Pick the ones this daemon should watch"},
	{"runtimes", "Agent CLIs", "We found these on your PATH"},
	{"daemon", "Daemon tuning", "Sensible defaults — tweak if you like"},
	{"review", "Review & start", "We'll write ~/.autopus and launch the daemon"},
}
