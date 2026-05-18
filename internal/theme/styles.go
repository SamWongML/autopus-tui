package theme

import "github.com/charmbracelet/lipgloss"

// Fg returns a style with foreground color c and the canvas background. Always
// setting bg on every styled span prevents the terminal default from bleeding
// through gaps between styled spans (see bg.go for the full explanation).
func Fg(c string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c)).
		Background(lipgloss.Color(Bg))
}

// FgOn paints fg c on explicit bg b. Used by chrome strips (tabbar, statusbar)
// where the surface bg differs from the canvas bg.
func FgOn(c, b string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c)).
		Background(lipgloss.Color(b))
}

// BG returns a style with only background set, useful for painting empty cells.
func BG(c string) lipgloss.Style { return lipgloss.NewStyle().Background(lipgloss.Color(c)) }

// Precomputed foreground styles for the palette colors used everywhere.
var (
	SText    = Fg(Text)
	SDim     = Fg(TextDim)
	SFaint   = Fg(TextFaint)
	SMute    = Fg(TextMute)
	SAccent  = Fg(Accent)
	SOK      = Fg(OK)
	SWarn    = Fg(Warn)
	SErr     = Fg(Err)
	SInfo    = Fg(Info)
	SViolet  = Fg(Violet)
	SBorder  = Fg(Border)
	SBorderH = Fg(BorderHi)
)
