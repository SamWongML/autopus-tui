package data

// DaemonInfo is the running-daemon metadata shown in the top bar and overview.
type DaemonInfo struct {
	PID     int
	Profile string
	Uptime  string
	Version string
}

var Daemon = DaemonInfo{
	PID:     48132,
	Profile: "default",
	Uptime:  "4h 32m 18s",
	Version: "0.18.2",
}
