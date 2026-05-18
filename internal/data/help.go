package data

// HelpGroup is a titled cluster of key/description rows in the help overlay.
type HelpGroup struct {
	Title string
	Rows  [][2]string
}

var HelpGroups = []HelpGroup{
	{"global", [][2]string{
		{"1 – 7", "jump to view by number"},
		{"tab / ⇧tab", "next / previous view"},
		{":", "command palette"},
		{"?", "this help"},
		{"esc", "close overlay · detach"},
		{"q", "quit (the daemon keeps running)"},
	}},
	{"lists · everywhere", [][2]string{
		{"j  /  ↓", "move down"},
		{"k  /  ↑", "move up"},
		{"g g", "first row"},
		{"G", "last row"},
		{"⌃d / ⌃u", "half-page down / up"},
		{"[  ]", "cycle filter pills (left / right)"},
		{"/", "search the list"},
		{"↵  /  →", "open / attach / drill"},
		{"esc", "clear filter, then go back"},
	}},
	{"overview", [][2]string{
		{"h j k l  /  ←↓↑→", "move focus between cards"},
		{"↵", "drill into the focused card's source view"},
	}},
	{"sessions", [][2]string{
		{"*", "jump to next needs-input session"},
		{"r", "reply to the focused session"},
		{"b", "background (/bg) — keep running, return to list"},
		{"t", "tail this session's logs"},
		{"d", "show diff for this run"},
		{"p", "approve / deny pending permissions"},
		{"x", "cancel"},
		{"!", "restart with a new prompt"},
	}},
	{"issues", [][2]string{
		{"n", "new issue"},
		{"a", "assign…"},
		{"s", "set status"},
		{"c", "add comment"},
	}},
	{"workspaces & runtimes", [][2]string{
		{"w", "watch / unwatch this workspace"},
		{"m", "list members"},
		{"o", "open in browser"},
		{"t", "test runtime handshake"},
		{"↺", "rescan $PATH for agent CLIs"},
	}},
	{"logs", [][2]string{
		{"f", "toggle follow"},
		{"c", "clear screen (file is untouched)"},
		{"s", "save current view to a file"},
	}},
	{"command palette · type :", [][2]string{
		{":login", "browser auth"},
		{":daemon stop", "stop the daemon"},
		{":daemon restart", "stop + start"},
		{":profile <name>", "switch profile (relaunches)"},
		{":new issue", "open the issue composer"},
		{":theme", "cycle warm-dark · paper · mono"},
		{":rescan agents", "re-detect agent CLIs on PATH"},
	}},
}
