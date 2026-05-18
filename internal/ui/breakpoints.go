package ui

// BP is the responsive breakpoint a view should target. The thresholds match
// the index.html mockup's CSS @media boundaries.
type BP int

const (
	BPxs BP = iota // < 100 cols — single column everywhere
	BPsm            // 100–139 — two cols, no detail rail beyond required
	BPmd            // 140–179 — full layout, smaller paddings
	BPlg            // ≥ 180 — full layout
)

// For returns the breakpoint for width w.
func For(w int) BP {
	switch {
	case w < 100:
		return BPxs
	case w < 140:
		return BPsm
	case w < 180:
		return BPmd
	default:
		return BPlg
	}
}
