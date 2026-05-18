// autopus-tui is a Bubble Tea TUI prototype for the (fictional) autopus agent
// daemon. The binary supports two modes:
//
//   - default: interactive terminal UI (alt-screen + mouse).
//   - --dump:  render every view (and overlay) to stdout, then exit. Used by
//     Claude Code / CI to validate rendering changes without a real terminal.
//
// See dump.go for the env vars that control --dump (AUTOPUS_DUMP_PROFILE,
// AUTOPUS_DUMP_WIDTH, AUTOPUS_DUMP_HEIGHT).
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"autopus-tui/internal/tui"
)

func main() {
	for _, a := range os.Args[1:] {
		switch a {
		case "--dump", "-d":
			runDump(os.Stdout)
			return
		case "--help", "-h":
			fmt.Println("autopus-tui — agent daemon TUI prototype")
			fmt.Println("Usage: autopus-tui [--dump]")
			fmt.Println()
			fmt.Println("  --dump     render every view + overlay to stdout, then exit.")
			fmt.Println("             env: AUTOPUS_DUMP_PROFILE={truecolor,256,ascii}")
			fmt.Println("                  AUTOPUS_DUMP_WIDTH=<int>  AUTOPUS_DUMP_HEIGHT=<int>")
			fmt.Println("                  AUTOPUS_DUMP_ONLY=<view>  (e.g. sessions,attach,help)")
			return
		}
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
