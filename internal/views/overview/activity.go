package overview

import (
	"strings"

	"autopus-tui/internal/data"
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func renderActivityStream(c ctx.Ctx, w, h int) string {
	items := data.ActivityStream
	maxItems := h / 2
	if maxItems < 1 {
		maxItems = 1
	}
	if len(items) > maxItems {
		items = items[:maxItems]
	}
	var b strings.Builder
	for _, it := range items {
		var g string
		switch it.Kind {
		case "warn":
			g = theme.SWarn.Render("!")
		case "info":
			g = theme.SInfo.Render("↻")
		case "dot":
			g = theme.SFaint.Render("·")
		default:
			g = ui.Glyph(it.Kind, c.Spin)
		}
		idCol := ui.PadRight(it.ID, 8)
		bodyW := w - 5 - 1 - 1 - 8 - 1
		if bodyW < 4 {
			bodyW = 4
		}
		body := ui.Truncate(it.Body, bodyW)
		line := theme.SMute.Render(it.T) + " " + g + " " + theme.SFaint.Render(idCol) + " " + theme.SDim.Render(body)
		b.WriteString(line + "\n")
		b.WriteString(ui.Dashed(w) + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}
