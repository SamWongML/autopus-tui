package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// mtop — a Bubble Tea TUI for the local Multica daemon, ported from
// `Multica Warm.html`. Tabs 1–6 switch between Status, Run, Tasks, Log,
// Config, Profiles. The data and chrome match the design canvas; in a real
// build, these screens would be wired to `~/.multica/daemon.sock` and the
// `multica issue run-messages` stream.

type model struct {
	width, height int

	tab     int // 0..5
	focus   int // 0 = primary list/stream, 1 = sidebar/detail (Tab toggles)
	selTask int // selection index in tasks roster
	selCfg  int // selection index in config rows
	selProf int // selection index in profiles
}

func initialModel() model {
	return model{
		width: 160, height: 48,
		tab: 0, focus: 0, selTask: 0, selCfg: 3, selProf: 0, // selCfg=3 ⇒ max-concurrent-tasks (dirty, selected in mock)
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "1":
			m.setTab(0)
		case "2":
			m.setTab(1)
		case "3":
			m.setTab(2)
		case "4":
			m.setTab(3)
		case "5":
			m.setTab(4)
		case "6":
			m.setTab(5)
		case "tab":
			m.focus = 1 - m.focus
		case "j", "down":
			m.moveSelection(+1)
		case "k", "up":
			m.moveSelection(-1)
		case "g", "home":
			m.setSelection(0)
		case "G", "end":
			m.setSelection(999)
		}
	}
	return m, nil
}

// setTab switches the active tab and resets focus to the primary panel —
// every tab's default focus is the list/stream side (P4 introduces visual
// dimming of the unfocused side; here we only need the model state correct).
func (m *model) setTab(t int) {
	m.tab = t
	m.focus = 0
}

func (m *model) moveSelection(delta int) {
	switch m.tab {
	case 0:
		m.selTask = clamp(m.selTask+delta, 0, minInt(5, len(tasks)-1))
	case 2:
		order := visualTaskOrder()
		pos := 0
		for i, idx := range order {
			if idx == m.selTask {
				pos = i
				break
			}
		}
		m.selTask = order[clamp(pos+delta, 0, len(order)-1)]
	case 4:
		m.selCfg = clamp(m.selCfg+delta, 0, len(cfgDaemon)-1)
	case 5:
		m.selProf = clamp(m.selProf+delta, 0, len(profiles)-1)
	}
}

func (m *model) setSelection(v int) {
	switch m.tab {
	case 0:
		m.selTask = clamp(v, 0, len(tasks)-1)
	case 2:
		order := visualTaskOrder()
		m.selTask = order[clamp(v, 0, len(order)-1)]
	case 4:
		m.selCfg = clamp(v, 0, len(cfgDaemon)-1)
	case 5:
		m.selProf = clamp(v, 0, len(profiles)-1)
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m model) View() string {
	if m.width < 40 || m.height < 12 {
		return sDim.Render("terminal too small — resize to at least 80×24")
	}

	tabs := tabBar(m.tab)
	top := topStrip(tabs, m.breadcrumb(), m.width)

	// Compute body height.
	chromeH := lipgloss.Height(top) + 2 /*cmdbar*/ + 2 /*footer*/
	bodyH := m.height - chromeH
	if bodyH < 8 {
		bodyH = 8
	}

	contentW := m.width - 2
	var body string
	var cmdMode, cmdPlaceholder string

	switch m.tab {
	case 0:
		body = renderStatus(contentW, bodyH, m.selTask)
		cmdMode = ":"
		cmdPlaceholder = "command — :restart  :pause  :follow  :workspace  :profile  :help"

	case 1:
		body = renderRun(contentW, bodyH)
		cmdMode = "reply"
		cmdPlaceholder = "reply inline — ⏎ send to claude · ↑ history · esc detach"

	case 2:
		body = renderTasks(contentW, bodyH, m.selTask)
		cmdMode = ":"
		cmdPlaceholder = "command — :kill t-1284  :pause runtime claude  :ws open"

	case 3:
		body = renderLog(contentW, bodyH)
		cmdMode = "/"
		cmdPlaceholder = `search — regex ok · ↑ history · Esc clear   task t-128\d`

	case 4:
		body = renderConfig(contentW, bodyH, m.selCfg)
		cmdMode = ":"
		cmdPlaceholder = "command — :save  :reload  :reset max-concurrent-tasks"

	case 5:
		body = renderProfiles(contentW, bodyH, m.selProf)
		cmdMode = ":"
		cmdPlaceholder = "command — :profile staging  :start  :stop  :default selfhost"
	}

	keys := m.footerKeys()

	// Pad body to bodyH so chrome sits at the bottom. Empty pad lines are
	// rendered as bg-painted blank rows so the canvas stays consistent below
	// short content.
	bodyLines := strings.Split(body, "\n")
	blankBodyLine := lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", contentW))
	for len(bodyLines) < bodyH {
		bodyLines = append(bodyLines, blankBodyLine)
	}
	if len(bodyLines) > bodyH {
		bodyLines = bodyLines[:bodyH]
	}
	body = fillBg(strings.Join(bodyLines, "\n"), bg)
	// Side gutter (one cell on each side) painted with canvas bg so the body
	// is centered without leaving unpainted cells at the edges.
	gutter := lipgloss.NewStyle().Background(bg).Render(" ")
	bodyLines = strings.Split(body, "\n")
	for i, line := range bodyLines {
		bodyLines[i] = gutter + line + gutter
	}
	bodyBlock := strings.Join(bodyLines, "\n")

	cmdrow := cmdBar(cmdMode, cmdPlaceholder, "", m.width)
	ft := footer(keys, m.width)

	return joinV(top, bodyBlock, cmdrow, ft)
}

// breadcrumb renders the right side of the top strip — a per-tab scope hint.
// Status owns the daemon card in its body, so it returns "" and lets topStrip
// fall back to just `● <profile>`. Every other tab gets a literal breadcrumb
// describing the current selection or filter.
func (m model) breadcrumb() string {
	switch m.tab {
	case 1:
		t := tasks[clamp(m.selTask, 0, len(tasks)-1)]
		return sDim.Render("tasks › ") + sFg1.Render(t.ID) +
			sDim.Render(" · ws ") + sFg1.Render(t.WS) +
			sDim.Render(" · seq ") + sFg1.Render(fmt.Sprintf("%d", t.Seq))
	case 2:
		_, by := groupedTasks()
		return sDim.Render("tasks · ") +
			sFg1.Render(fmt.Sprintf("%d", len(by["Needs input"]))) +
			sDim.Render(" needs input · ") +
			sFg1.Render(fmt.Sprintf("%d", len(by["Working"]))) +
			sDim.Render(" working")
	case 3:
		return sDim.Render("log · level=info · follow ") + sOk.Render("●")
	case 4:
		row := cfgDaemon[clamp(m.selCfg, 0, len(cfgDaemon)-1)]
		s := sDim.Render("config · ") + sFg1.Render(row.K)
		if row.Dirty {
			s += " " + sWarn.Render("▲")
		}
		return s
	case 5:
		return sDim.Render("profiles · active=") + sFg1.Render(daemon.Profile)
	}
	return ""
}

// footerKeys returns the key-hint row for the bottom footer. For tabs that
// expose two focus zones (run, tasks) the hint set switches with m.focus and
// is suffixed with a Tab/switch-focus reminder. Other tabs are focus-agnostic.
func (m model) footerKeys() []KeyHint {
	switch m.tab {
	case 0:
		return []KeyHint{
			{"1-6", "tabs"}, {"↑↓", "select"}, {"⏎", "open"},
			{"r", "restart"}, {"f", "follow log"}, {"n", "new task"}, {":", "command"},
		}
	case 1:
		var ks []KeyHint
		if m.focus == 0 {
			ks = []KeyHint{
				{"j/k", "scroll"}, {"g/G", "top/end"},
				{"f", "follow"}, {"y", "yank"}, {"/", "filter"},
			}
		} else {
			ks = []KeyHint{
				{"k", "kill"}, {"r", "reply"}, {"o", "$EDITOR"}, {"c", "copy"},
			}
		}
		return append(ks, KeyHint{"Tab", "switch focus"})
	case 2:
		var ks []KeyHint
		if m.focus == 0 {
			ks = []KeyHint{
				{"↑↓", "select"}, {"⏎", "open run"}, {"/", "filter"}, {"s", "sort"},
			}
		} else {
			ks = []KeyHint{
				{"r", "reply"}, {"k", "kill"}, {"c", "copy"}, {"o", "open"},
			}
		}
		return append(ks, KeyHint{"Tab", "switch focus"})
	case 3:
		return []KeyHint{
			{"f", "follow"}, {"/", "search"}, {"1-4", "level"},
			{"t", "hide ticks"}, {"j/k", "scroll"}, {"y", "yank"}, {"o", "$PAGER"},
		}
	case 4:
		return []KeyHint{
			{"↑↓", "select"}, {"⏎/e", "edit"}, {"⌃s", "save"},
			{"r", "reload"}, {"d", "defaults"}, {"v", "view file"},
		}
	case 5:
		return []KeyHint{
			{"↑↓", "select"}, {"⏎", "attach"}, {"s", "start"}, {"S", "stop"},
			{"r", "restart"}, {"n", "new"}, {"d", "set default"},
		}
	}
	return nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--dump" {
		// --dump is the offline rendering path: stdout is a pipe, so lipgloss
		// would normally disable color. Honor MTOP_DUMP_PROFILE=truecolor|256
		// to force a richer profile for visual inspection of the ANSI stream.
		switch os.Getenv("MTOP_DUMP_PROFILE") {
		case "truecolor":
			lipgloss.SetColorProfile(termenv.TrueColor)
		case "256":
			lipgloss.SetColorProfile(termenv.ANSI256)
		}
		for tab := 0; tab < 6; tab++ {
			m := initialModel()
			m.tab = tab
			fmt.Println(strings.Repeat("=", 80))
			fmt.Printf("=== tab %d\n", tab+1)
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println(m.View())
			fmt.Println()
		}
		return
	}
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
