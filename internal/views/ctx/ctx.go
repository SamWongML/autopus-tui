// Package ctx defines the small render-time context every view receives. It
// carries the values that change every tick (Now, Spin) without forcing them
// onto every view function signature.
package ctx

import "time"

// Ctx is the per-frame render context passed to view View methods.
type Ctx struct {
	Now  time.Time
	Spin int
}
