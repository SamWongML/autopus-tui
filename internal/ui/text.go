// Package ui provides generic, stateless rendering primitives: text helpers,
// bordered panels, chips, bars, k/v rows, frame painting. It depends only on
// theme; it has no knowledge of any view or data type.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Truncate clips s to n display columns, appending "…" when it overflows.
func Truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= n {
		return s
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n == 1 {
		return "…"
	}
	return string(r[:n-1]) + "…"
}

// PadRight pads s with trailing spaces to width n; truncates if longer.
func PadRight(s string, n int) string {
	s = Truncate(s, n)
	w := lipgloss.Width(s)
	if w >= n {
		return s
	}
	return s + strings.Repeat(" ", n-w)
}

// PadLeft pads s with leading spaces to width n; truncates if longer.
func PadLeft(s string, n int) string {
	s = Truncate(s, n)
	w := lipgloss.Width(s)
	if w >= n {
		return s
	}
	return strings.Repeat(" ", n-w) + s
}

// JoinRight joins a left-aligned label and right-aligned value in a row of
// width w with a single space minimum gap.
func JoinRight(left, right string, w int) string {
	gap := w - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

// Wrap breaks s into lines no wider than w using word-wrap.
func Wrap(s string, w int) string {
	if w < 4 {
		w = 4
	}
	words := strings.Fields(s)
	var b strings.Builder
	col := 0
	for _, wd := range words {
		wlen := lipgloss.Width(wd)
		if col == 0 {
			b.WriteString(wd)
			col = wlen
			continue
		}
		if col+1+wlen > w {
			b.WriteByte('\n')
			b.WriteString(wd)
			col = wlen
		} else {
			b.WriteByte(' ')
			b.WriteString(wd)
			col += 1 + wlen
		}
	}
	return b.String()
}

// SplitLines splits s on '\n', returning nil for the empty string.
func SplitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// Commafy formats an int with thousands separators.
func Commafy(n int) string {
	if n < 0 {
		return "-" + Commafy(-n)
	}
	s := itoaPlain(n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	first := len(s) % 3
	if first > 0 {
		b.WriteString(s[:first])
		if len(s) > first {
			b.WriteByte(',')
		}
	}
	for i := first; i < len(s); i += 3 {
		b.WriteString(s[i : i+3])
		if i+3 < len(s) {
			b.WriteByte(',')
		}
	}
	return b.String()
}

// Itoa is Commafy for nonzero integers; "0" otherwise. Convenience used by
// table cells where 0 should render plainly.
func Itoa(n int) string {
	if n == 0 {
		return "0"
	}
	return Commafy(n)
}

func itoaPlain(n int) string {
	if n == 0 {
		return "0"
	}
	var b strings.Builder
	for _, d := range intDigits(n) {
		b.WriteByte('0' + byte(d))
	}
	return b.String()
}

func intDigits(n int) []int {
	if n == 0 {
		return []int{0}
	}
	var out []int
	for n > 0 {
		out = append([]int{n % 10}, out...)
		n /= 10
	}
	return out
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
