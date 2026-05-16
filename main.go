package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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

	help help.Model
	ti   textinput.Model

	// runVp scrolls the Run-tab messages pane; logVp scrolls the Log-tab
	// left pane. Their dimensions + content are refreshed on WindowSizeMsg.
	runVp, logVp viewport.Model

	// tasksTbl owns the Tasks-tab roster (cursor, scroll, key handling). On
	// tab 2, Up/Down/k/j/g/G route into the table instead of moveSelection.
	// Resized in resizeViewports() alongside the viewports.
	tasksTbl table.Model

	// cmdOverride is the mode string set by the user pressing `:` or `/`,
	// which takes precedence over the tab's default mode (defaultCmdMode).
	// Cleared on Esc / Enter.
	cmdOverride string
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 256
	ti.TextStyle = sFg
	ti.PlaceholderStyle = sDim
	m := model{
		width: 160, height: 48,
		tab: 0, selTask: 0, selCfg: 3, selProf: 0, // selCfg=3 ⇒ max-concurrent-tasks (dirty, selected in mock)
		help:     footerHelp(),
		ti:       ti,
		runVp:    viewport.New(0, 0),
		logVp:    viewport.New(0, 0),
		tasksTbl: newTasksTable(),
	}
	m.resizeViewports()
	return m
}

// chromeHeight is the number of rows occupied by tabs (1+1 border) + cmdbar
// (1+1 border) + footer (1+1 border) — the fixed chrome that wraps the body
// in main.View. Kept in sync with View()'s `chromeH` literal.
const chromeHeight = 6

// resizeViewports rebuilds both per-tab viewports' dimensions and content
// from the current model width/height. Called on WindowSizeMsg and once at
// initial-model construction so --dump (which never receives a resize event)
// renders deterministic frame 0.
func (m *model) resizeViewports() {
	contentW := m.width - 2
	bodyH := m.height - chromeHeight
	if bodyH < 8 {
		bodyH = 8
	}

	// Run tab — same geometry as runBody().
	sidebarW := 38
	if sidebarW > contentW/3 {
		sidebarW = contentW / 3
	}
	msgsW := contentW - sidebarW - 2
	runPaneH := bodyH - 3 // runHeader is 3 rows
	m.runVp.Width = msgsW - 2
	m.runVp.Height = runPaneH - 2
	m.runVp.SetContent(renderMessages(m.runVp.Width - 2))

	// Log tab — same geometry as renderLog().
	rightW := 38
	if rightW > contentW*30/100 {
		rightW = contentW * 30 / 100
	}
	leftW := contentW - rightW - 2
	logPaneH := bodyH - 3 // filter row is 3 rows
	if logPaneH < 6 {
		logPaneH = 6
	}
	m.logVp.Width = leftW - 2
	m.logVp.Height = logPaneH - 2
	// Feed every row (no h-2 cap): pass a height tall enough that logTable
	// emits all `tripled` rows; viewport picks the visible window.
	tripled := append(append(append([]LogLine{}, logLines...), logLines...), logLines[:4]...)
	m.logVp.SetContent(logTable(tripled, leftW, len(tripled)+2, logFull))

	// Tasks tab — same geometry as renderTasks(). Right preview pane mirrors
	// view_tasks.go: width*40% (capped at 48). Table fills remaining left
	// pane minus the 2-cell pane padding on each side.
	tRightW := 48
	if tRightW > contentW*40/100 {
		tRightW = contentW * 40 / 100
	}
	tLeftW := contentW - tRightW - 2
	tPaneH := bodyH - 3 // filter row is 3 rows
	if tPaneH < 6 {
		tPaneH = 6
	}
	m.tasksTbl.SetWidth(tLeftW - 4)
	m.tasksTbl.SetHeight(tPaneH - 2)
	m.tasksTbl.SetColumns(tasksTableCols(tLeftW - 4))
}

func (m model) Init() tea.Cmd { return tea.Batch(liveSpin.Tick, textinput.Blink) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.resizeViewports()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		liveSpin, cmd = liveSpin.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		// While the command bar is focused, all keys go to the textinput
		// except Esc (cancel) and Enter (would-be execute; just clears).
		if m.ti.Focused() {
			switch msg.String() {
			case "esc":
				m.ti.Blur()
				m.ti.Reset()
				m.cmdOverride = ""
				return m, nil
			case "enter":
				// Command dispatch lands in a later phase; for now treat
				// Enter as commit-and-clear.
				m.ti.Blur()
				m.ti.Reset()
				m.cmdOverride = ""
				return m, nil
			}
			var cmd tea.Cmd
			m.ti, cmd = m.ti.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case ":":
			m.cmdOverride = ":"
			return m, m.ti.Focus()
		case "/":
			m.cmdOverride = "/"
			return m, m.ti.Focus()
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
			if vp := m.scrollVp(); vp != nil {
				vp.LineDown(1)
			} else if m.tab == 2 {
				m.tasksTbl, _ = m.tasksTbl.Update(msg)
			} else {
				m.moveSelection(+1)
			}
		case "k", "up":
			if vp := m.scrollVp(); vp != nil {
				vp.LineUp(1)
			} else if m.tab == 2 {
				m.tasksTbl, _ = m.tasksTbl.Update(msg)
			} else {
				m.moveSelection(-1)
			}
		case "ctrl+d":
			if vp := m.scrollVp(); vp != nil {
				vp.HalfPageDown()
			}
		case "ctrl+u":
			if vp := m.scrollVp(); vp != nil {
				vp.HalfPageUp()
			}
		case "g", "home":
			if vp := m.scrollVp(); vp != nil {
				vp.GotoTop()
			} else if m.tab == 2 {
				m.tasksTbl.GotoTop()
			} else {
				m.setSelection(0)
			}
		case "G", "end":
			if vp := m.scrollVp(); vp != nil {
				vp.GotoBottom()
			} else if m.tab == 2 {
				m.tasksTbl.GotoBottom()
			} else {
				m.setSelection(999)
			}
		}
	}
	return m, nil
}

// scrollVp returns the viewport that consumes j/k/g/G/ctrl-d/ctrl-u on the
// active tab, or nil if the active tab uses moveSelection/setSelection
// instead. Pointer return so callers can mutate scroll state in place.
func (m *model) scrollVp() *viewport.Model {
	switch m.tab {
	case 1:
		return &m.runVp
	case 3:
		return &m.logVp
	}
	return nil
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
	var cmdMode, cmdPlaceholder string

	switch m.tab {
	case 0:
		body = renderStatus(contentW, bodyH, m.selTask)
		cmdMode = ":"
		cmdPlaceholder = "command — :restart  :pause  :follow  :workspace  :profile  :help"

	case 1:
		body = renderRun(contentW, bodyH, m.runVp)
		cmdMode = "reply"
		cmdPlaceholder = "reply inline — ⏎ send to claude · ↑ history · esc detach"

	case 2:
		body = renderTasks(contentW, bodyH, m.tasksTbl)
		cmdMode = ":"
		cmdPlaceholder = "command — :kill t-1284  :pause runtime claude  :ws open"

	case 3:
		body = renderLog(contentW, bodyH, m.logVp)
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

	// User-pressed `:` or `/` overrides the tab's default mode while the
	// input is focused.
	if m.cmdOverride != "" {
		cmdMode = m.cmdOverride
	}
	cmdrow := cmdBar(cmdMode, cmdPlaceholder, m.ti, m.width)
	ft := footer(m.help, tabKeys[m.tab], m.width)

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
		w, h := 160, 48
		if v, err := strconv.Atoi(os.Getenv("MTOP_DUMP_WIDTH")); err == nil && v > 0 {
			w = v
		}
		if v, err := strconv.Atoi(os.Getenv("MTOP_DUMP_HEIGHT")); err == nil && v > 0 {
			h = v
		}
		for tab := 0; tab < 6; tab++ {
			m := initialModel()
			m.width, m.height = w, h
			m.help.Width = w
			m.tab = tab
			// initialModel() sized viewports to the default 160×48; rebuild
			// after the size override so --dump renders the requested W×H.
			m.resizeViewports()
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
