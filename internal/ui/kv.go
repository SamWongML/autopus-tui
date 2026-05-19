package ui

import (
	"autopus-tui/internal/theme"
)

// KVRow returns "left ............... right" plus a dashed bottom-border line
// below it. The two lines emulate HTML's border-bottom: 1px dashed pattern.
// If vc is empty, theme.Text is used for the right value color.
func KVRow(k, v, vc string, w int) string {
	if vc == "" {
		vc = theme.Text
	}
	left := theme.SFaint.Render(k)
	right := theme.Fg(vc).Render(v)
	return JoinRight(left, right, w) + "\n" + Dashed(w)
}

// KVRowDashed is the explicitly-named, dashed-separator form of KVRow. Use it
// in detail panels where the dashed separator below each row is intentional
// (issues/runtimes/workspaces/attach details).
func KVRowDashed(k, v, vc string, w int) string {
	return KVRow(k, v, vc, w)
}

// KVRowFlat is KVRow without the trailing dashed line. Used when the caller
// supplies its own divider or packs rows tightly.
func KVRowFlat(k, v, vc string, w int) string {
	if vc == "" {
		vc = theme.Text
	}
	left := theme.SFaint.Render(k)
	right := theme.Fg(vc).Render(v)
	return JoinRight(left, right, w)
}
