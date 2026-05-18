package theme

import (
	"fmt"
	"strings"
)

// bgSGR returns the raw ANSI sequence that sets the background to a #RRGGBB
// color. Used by WithBg to patch ANSI resets emitted by lipgloss spans.
func bgSGR(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return ""
	}
	parse := func(a, b byte) int {
		hx := func(c byte) int {
			switch {
			case c >= '0' && c <= '9':
				return int(c - '0')
			case c >= 'a' && c <= 'f':
				return int(c-'a') + 10
			case c >= 'A' && c <= 'F':
				return int(c-'A') + 10
			}
			return 0
		}
		return hx(a)<<4 | hx(b)
	}
	r := parse(hex[1], hex[2])
	g := parse(hex[3], hex[4])
	b := parse(hex[5], hex[6])
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// WithBg secures the background of a composed region against ANSI-reset bleed.
//
// Lipgloss closes every styled span with a full reset (\x1b[0m), which wipes
// the parent's background too (charmbracelet/lipgloss#209). Anywhere inner
// styled spans, raw spaces, or JoinHorizontal padding sit, the terminal
// default bg shows through as visible horizontal bands ("straps").
//
// WithBg post-processes a region's rendered output so:
//   - the region opens with the desired bg set,
//   - every inner reset is immediately followed by the same bg set,
//   - the region closes with a single reset.
//
// BCE-aware terminals carry the bg across newlines and paint to end-of-line,
// so patching resets alone is enough. WithBg nests safely.
func WithBg(s, hex string) string {
	open := bgSGR(hex)
	if open == "" {
		return s
	}
	s = strings.ReplaceAll(s, "\x1b[0m", "\x1b[0m"+open)
	return open + s + "\x1b[0m"
}
