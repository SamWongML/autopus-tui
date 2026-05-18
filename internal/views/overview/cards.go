package overview

import (
	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
	"autopus-tui/internal/views/ctx"
)

func overviewPanel(m Model, c ctx.Ctx, id string, w, h int) string {
	idx := cardIdx(id)
	focused := m.Focus == idx
	title := Cards[idx].Title
	right := ""
	if focused {
		right = "↵ " + theme.SAccent.Render(Cards[idx].Jump)
	} else {
		switch id {
		case "activity":
			right = "last 5m"
		case "sessions":
			right = "press " + theme.SAccent.Render("2")
		}
	}
	body := ""
	switch id {
	case "daemon":
		body = renderDaemonCard(w-4, h-4)
	case "env":
		body = renderEnvCard(w-4, h-4)
	case "server":
		body = renderServerCard(w-4, h-4)
	case "sessions":
		body = renderSessionsSummary(w-4, h-4)
	case "activity":
		body = renderActivityStream(c, w-4, h-4)
	}
	return ui.Panel(title, right, body, w, h, focused, false)
}
