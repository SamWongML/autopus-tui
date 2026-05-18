package logs

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderBody(rows []data.LogLine, follow bool, w, h int) string {
	const (
		cT     = 8
		cLevel = 6
		cSrc   = 10
	)
	gaps := 3
	cMsg := w - cT - cLevel - cSrc - gaps
	if cMsg < 12 {
		cMsg = 12
	}
	visH := h
	if follow {
		visH--
	}
	start := 0
	if len(rows) > visH {
		start = len(rows) - visH
	}
	var b strings.Builder
	for i := start; i < len(rows); i++ {
		l := rows[i]
		col := theme.LogLevels[l.Level]
		if col == "" {
			col = theme.TextDim
		}
		line := theme.SMute.Render(ui.PadRight(l.T, cT)) + " " +
			theme.Fg(col).Render(ui.PadRight(l.Level, cLevel)) + " " +
			theme.SFaint.Render(ui.PadRight(l.Src, cSrc)) + " " +
			theme.SDim.Render(ui.Truncate(l.Msg, cMsg))
		b.WriteString(line + "\n")
	}
	if follow {
		b.WriteString(theme.SAccent.Render("▌") + " " + theme.SFaint.Render("tailing…"))
	}
	return strings.TrimRight(b.String(), "\n")
}
