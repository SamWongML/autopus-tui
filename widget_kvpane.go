package main

import "strings"

// kvSection is one block of pre-rendered lines inside a kvPane body. The
// `title` is passed to hr() as the divider label; an empty title yields a
// plain dashed line (or no divider at all when it's the first section).
type kvSection struct {
	title string
	lines []string
}

// kvPane renders the body of a section-based pane: each section is preceded
// by an `hr(title, inner)` divider, except the first section when its title
// is empty (so callers can lead with an untitled block such as a KV grid or
// a custom header row). The caller wraps the returned body with `pane()`.
func kvPane(inner int, sections []kvSection) string {
	var out []string
	for i, s := range sections {
		if !(i == 0 && s.title == "") {
			out = append(out, hr(s.title, inner))
		}
		out = append(out, s.lines...)
	}
	return strings.Join(out, "\n")
}
