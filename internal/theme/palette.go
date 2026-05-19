// Package theme owns the visual vocabulary: colors, lipgloss styles, glyph
// tables, and the bg-leak protection helpers. It has no dependencies on any
// other package in this tree.
package theme

// Warm-dark Anthropic palette. All UI colors must come from this list — never
// inline a hex value in a view package.
const (
	Bg          = "#1a1814"
	Bg2         = "#15130f"
	Surface     = "#211e1a"
	Surface2    = "#2a2620"
	Surface3    = "#322d26"
	Border      = "#3a342c"
	BorderHi    = "#4a4238"
	Text        = "#ebe5d8"
	TextDim     = "#a39885"
	TextFaint   = "#756c5d"
	TextMute    = "#574f43"
	Accent      = "#d97757"
	AccentDim   = "#ab4929"
	AccentFaint = "#2c2018"
	OK          = "#8aa67a"
	Warn        = "#d4a857"
	Err         = "#c46b5e"
	Info        = "#6f93b8"
	Violet      = "#9b87b8"
)
