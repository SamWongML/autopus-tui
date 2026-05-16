package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderTaskDetail is the single source of truth for the task-detail body
// rendered by the Tasks-tab preview pane and the Run-tab peek overlay. It
// returns the inner body only — the caller wraps it in `pane(...)` with the
// title/hint appropriate to its surface.
//
//	mode == "summary"  → header + LAST 5 MESSAGES + ACTIONS  (Tasks tab)
//	mode == "stream"   → header + FILES TOUCHED + MSG TYPE COUNTS  (Run peek)
//
// The shared header (status, runtime, model, pid, branch, cost, elapsed) is
// identical across both modes — no other view should re-render a task's
// identity fields. Pads to `height` so the pane always fills its slot.
func renderTaskDetail(t Task, mode string, width, height int) string {
	lines := taskHeaderRows(t, width)

	switch mode {
	case "summary":
		lines = append(lines, hr("last 5 messages", width))
		tail := msgs
		if len(tail) > 5 {
			tail = tail[len(tail)-5:]
		}
		for _, m := range tail {
			lines = append(lines, msgRowCompact(m, width))
		}
		lines = append(lines, hr("actions", width))
		for _, a := range taskActions {
			lines = append(lines, actionRow(a.K, a.L))
		}

	case "stream":
		lines = append(lines, hr("files touched", width))
		for _, f := range taskFilesTouched {
			lines = append(lines, fileRow(f.path, f.diff, width))
		}
		lines = append(lines, hr("msg type counts", width))
		for _, c := range taskMsgCounts {
			lines = append(lines, countBar(c.label, c.n, taskMsgCountTotal, c.tone, width))
		}
	}

	body := strings.Join(lines, "\n")
	if pad := height - lipgloss.Height(body); pad > 0 {
		body = body + "\n" + bgPadV(width, pad)
	}
	return body
}

// taskHeaderRows renders the kv rows shared by every detail view. Pre-styled
// values (status glyph, agent chip) are passed through kv with tone="" so the
// inner styling wins.
func taskHeaderRows(t Task, width int) []string {
	statusVal := statusStyle(t.Status).Render(statusGlyph(t.Status) + " " + t.Status)
	return []string{
		kv("status", statusVal, "", width),
		kv("runtime", agentChip(t.Runtime), "", width),
		kv("model", modelForRuntime(t.Runtime), "", width),
		kv("pid", taskPID(t), "", width),
		kv("branch", taskBranch(t), "", width),
		kv("cost", t.Cost, "", width),
		kv("elapsed", t.Started, "", width),
	}
}

// modelForRuntime looks up a runtime's default model by name. Empty if the
// runtime isn't registered.
func modelForRuntime(name string) string {
	for _, r := range runtimes {
		if r.Name == name {
			return r.Model
		}
	}
	return ""
}

// taskPID returns the OS pid for a running task. The mock has no per-task pid
// field; show the seeded value for working/waiting tasks and "—" otherwise so
// the row reads honestly for queued/failed entries.
func taskPID(t Task) string {
	if t.Status == "working" || t.Status == "waiting" {
		return "18421"
	}
	return "—"
}

// taskBranch returns the git branch a task is working on. Mock-only; replace
// with t.Branch once that field exists in data.go.
func taskBranch(t Task) string {
	return "feat/" + strings.TrimPrefix(t.Issue, "#") + "-wip"
}

// taskActions is the fixed action menu shown under a task's detail. The keys
// here mirror what main.go binds in tab-2 focus=1 (footer hints).
var taskActions = []KeyHint{
	{"⏎", "open run viewer"},
	{"r", "reply"},
	{"k", "kill (SIGTERM)"},
	{"K", "force kill"},
	{"o", "open workspace"},
	{"c", "copy id"},
}

// taskFilesTouched is the mock list of files a working task has edited. Will
// be replaced by t.FilesTouched once the daemon reports it.
var taskFilesTouched = []struct{ path, diff string }{
	{"db/migrate/20260514_add_idempotency_key.rb", "+42"},
	{"app/models/webhook_event.rb", "+8 −2"},
	{"spec/models/webhook_event_spec.rb", "+24"},
}

// taskMsgCounts is the mock breakdown of message types in the current run.
// Total is the sum used as the bar denominator.
var taskMsgCounts = []struct {
	label, tone string
	n           int
}{
	{"thinking", "dim", 32},
	{"tool_call", "info", 28},
	{"tool_result", "ok", 20},
	{"text", "amber", 5},
	{"error", "err", 1},
}

const taskMsgCountTotal = 86

// fileRow renders "path                       +diff" with the diff right-aligned
// in ok-green. Truncates the path when the row would overflow.
func fileRow(path, diff string, width int) string {
	gap := width - lipgloss.Width(path) - lipgloss.Width(diff)
	if gap < 1 {
		gap = 1
		path = truncate(path, width-lipgloss.Width(diff)-1)
	}
	return sFg1.Render(path) + strings.Repeat(" ", gap) + sOk.Render(diff)
}

// msgRowCompact renders a single message as one line: "T  type  body" with the
// type dim and the body fg1, truncated to fit. Used inside the detail pane's
// LAST 5 MESSAGES section.
func msgRowCompact(m Msg, width int) string {
	const typeW = 12
	t := sDim.Render(m.T)
	typ := sDim.Render(padR(m.Type, typeW))
	prefix := t + "  " + typ + "  "
	avail := width - lipgloss.Width(prefix)
	if avail < 4 {
		avail = 4
	}
	body := m.Body
	if m.Type == "tool_call" {
		body = m.Tool + " " + m.Args
	}
	return prefix + sFg1.Render(truncate(body, avail))
}

// actionRow renders "[key] description" — a key-cap pill plus a dim label.
func actionRow(k, v string) string {
	cap := lipgloss.NewStyle().Foreground(fg1).Background(bg2).Padding(0, 1).Render(k)
	return cap + " " + sFg1.Render(v)
}
