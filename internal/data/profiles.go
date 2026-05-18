package data

// Profile is a named daemon configuration (server endpoint + auth).
type Profile struct {
	ID     string
	Server string
	WS     string
	Active bool
	PID    int
	Uptime string
}

var Profiles = []Profile{
	{ID: "default", Server: "app.autopus.ai", WS: "wss://api.autopus.ai/ws", Active: true, PID: 48132, Uptime: "4h 32m"},
	{ID: "staging", Server: "staging.autopus.ai", WS: "wss://staging-api.autopus.ai/ws", Uptime: "—"},
	{ID: "self-hosted", Server: "app.acme.dev", WS: "wss://api.acme.dev/ws", Uptime: "—"},
}
