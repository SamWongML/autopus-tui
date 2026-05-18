package app

import "time"

// ClockMsg ticks once per second; the root model uses it to refresh the wall
// clock shown in the top bar.
type ClockMsg time.Time

// SpinMsg ticks every ~90ms to advance the braille spinner frame counter
// shared by all views.
type SpinMsg time.Time

// NavigateMsg asks the root model to switch routes. Children emit this via a
// tea.Cmd; the root reads To/Attach/Overlay and updates its own state.
//
//   - To       — route ID to switch to (empty = leave route unchanged)
//   - Attach   — session ID to attach to (sets the attach view's session and
//                routes to it; the actual top tab is unchanged)
//   - Overlay  — "", "help", "palette", or "onboarding"
type NavigateMsg struct {
	To      string
	Attach  string
	Overlay string
}
