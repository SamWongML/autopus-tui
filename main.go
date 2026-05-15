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
	selTask int // selection index in tasks roster
	selCfg  int // selection index in config rows
	selProf int // selection index in profiles
}

func initialModel() model {
	return model{
		width: 160, height: 48,
		tab: 0, selTask: 0, selCfg: 3, selProf: 0, // selCfg=3 ⇒ max-concurrent-tasks (dirty, selected in mock)
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
			m.tab = 0
		case "2":
			m.tab = 1
		case "3":
			m.tab = 2
		case "4":
			m.tab = 3
		case "5":
			m.tab = 4
		case "6":
			m.tab = 5
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

func (m *model) moveSelection(delta int) {
	switch m.tab {
	case 0:
		m.selTask = clamp(m.selTask+delta, 0, minInt(5, len(tasks)-1))
	case 2:
		m.selTask = clamp(m.selTask+delta, 0, len(tasks)-1)
	case 4:
		m.selCfg = clamp(m.selCfg+delta, 0, len(cfgDaemon)-1)
	case 5:
		m.selProf = clamp(m.selProf+delta, 0, len(profiles)-1)
	}
}

func (m *model) setSelection(v int) {
	switch m.tab {
	case 0, 2:
		m.selTask = clamp(v, 0, len(tasks)-1)
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

	tabs := tabBar(m.tab, m.width)

	// Compute body height.
	chromeH := lipgloss.Height(tabs) + 2 /*cmdbar*/ + 2 /*footer*/
	bodyH := m.height - chromeH
	if bodyH < 8 {
		bodyH = 8
	}

	contentW := m.width - 2
	var body string
	var keys []KeyHint
	var cmdMode, cmdPlaceholder string

	switch m.tab {
	case 0:
		body = renderStatus(contentW, bodyH, m.selTask)
		keys = []KeyHint{
			{"1-6", "tabs"}, {"↑↓", "select"}, {"⏎", "open"},
			{"r", "restart"}, {"f", "follow log"}, {"n", "new task"}, {":", "command"},
		}
		cmdMode = ":"
		cmdPlaceholder = "command — :restart  :pause  :follow  :workspace  :profile  :help"

	case 1:
		body = renderRun(contentW, bodyH)
		keys = []KeyHint{
			{"esc", "back"}, {"j/k", "scroll"}, {"g/G", "top/end"},
			{"f", "follow"}, {"/", "filter"}, {"y", "yank seq"}, {"k", "kill"},
		}
		cmdMode = "reply"
		cmdPlaceholder = "reply inline — ⏎ send to claude · ↑ history · esc detach"

	case 2:
		body = renderTasks(contentW, bodyH, m.selTask)
		keys = []KeyHint{
			{"↑↓", "select"}, {"⏎", "open run"}, {"k", "kill"},
			{"p", "pause runtime"}, {"s", "sort"}, {"/", "filter"}, {"a", "show all"},
		}
		cmdMode = ":"
		cmdPlaceholder = "command — :kill t-1284  :pause runtime claude  :ws open"

	case 3:
		body = renderLog(contentW, bodyH)
		keys = []KeyHint{
			{"f", "follow"}, {"/", "search"}, {"1-4", "level"},
			{"t", "hide ticks"}, {"j/k", "scroll"}, {"y", "yank"}, {"o", "$PAGER"},
		}
		cmdMode = "/"
		cmdPlaceholder = `search — regex ok · ↑ history · Esc clear   task t-128\d`

	case 4:
		body = renderConfig(contentW, bodyH, m.selCfg)
		keys = []KeyHint{
			{"↑↓", "select"}, {"⏎/e", "edit"}, {"⌃s", "save"},
			{"r", "reload"}, {"d", "defaults"}, {"v", "view file"},
		}
		cmdMode = ":"
		cmdPlaceholder = "command — :save  :reload  :reset max-concurrent-tasks"

	case 5:
		body = renderProfiles(contentW, bodyH, m.selProf)
		keys = []KeyHint{
			{"↑↓", "select"}, {"⏎", "attach"}, {"s", "start"}, {"S", "stop"},
			{"r", "restart"}, {"n", "new"}, {"d", "set default"},
		}
		cmdMode = ":"
		cmdPlaceholder = "command — :profile staging  :start  :stop  :default selfhost"
	}

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

	return joinV(tabs, bodyBlock, cmdrow, ft)
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
