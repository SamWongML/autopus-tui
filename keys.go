package main

import "github.com/charmbracelet/bubbles/key"

// tabKeyMap satisfies help.KeyMap for one tab. The footer renders ShortHelp;
// FullHelp is reserved for the `?` overlay (Phase F.6) and currently mirrors
// ShortHelp as a single group.
type tabKeyMap struct {
	short []key.Binding
}

func (k tabKeyMap) ShortHelp() []key.Binding  { return k.short }
func (k tabKeyMap) FullHelp() [][]key.Binding { return [][]key.Binding{k.short} }

// bk builds a key binding whose Help.Key is wrapped with one space of padding
// on each side. bubbles/help renders ShortKey with Inline(true), which strips
// style.Padding; the literal spaces preserve the pre-Phase-B pill-shaped
// keycap look when the surrounding ShortKey style paints a Background(bg2).
func bk(keys []string, k, desc string) key.Binding {
	return key.NewBinding(key.WithKeys(keys...), key.WithHelp(" "+k+" ", desc))
}

// helpKey / helpQuit are the universal right-side footer hints, rendered as a
// separate ShortHelpView so footer() can stretch-justify them against the
// per-tab hints on the left.
var (
	helpKey  = bk([]string{"?"}, "?", "help")
	helpQuit = bk([]string{"q", "ctrl+c"}, "q", "quit")
)

// tabKeys[i] is the ShortHelp displayed for tab i (0-indexed: status, run,
// tasks, log, config, profiles). Most bindings are display-only today; the
// Update handler still only consumes 1-6 / j / k / up / down / g / G / home /
// end / q / ctrl+c / esc, matching pre-Phase-B behavior.
var tabKeys = []tabKeyMap{
	{short: []key.Binding{
		bk([]string{"1", "2", "3", "4", "5", "6"}, "1-6", "tabs"),
		bk([]string{"up", "down", "k", "j"}, "↑↓", "select"),
		bk([]string{"enter"}, "⏎", "open"),
		bk([]string{"r"}, "r", "restart"),
		bk([]string{"f"}, "f", "follow log"),
		bk([]string{"n"}, "n", "new task"),
		bk([]string{":"}, ":", "command"),
	}},
	{short: []key.Binding{
		bk([]string{"esc"}, "esc", "back"),
		bk([]string{"j", "k"}, "j/k", "scroll"),
		bk([]string{"g", "G"}, "g/G", "top/end"),
		bk([]string{"f"}, "f", "follow"),
		bk([]string{"/"}, "/", "filter"),
		bk([]string{"y"}, "y", "yank seq"),
		bk([]string{"k"}, "k", "kill"),
	}},
	{short: []key.Binding{
		bk([]string{"up", "down"}, "↑↓", "select"),
		bk([]string{"enter"}, "⏎", "open run"),
		bk([]string{"k"}, "k", "kill"),
		bk([]string{"p"}, "p", "pause runtime"),
		bk([]string{"s"}, "s", "sort"),
		bk([]string{"/"}, "/", "filter"),
		bk([]string{"a"}, "a", "show all"),
	}},
	{short: []key.Binding{
		bk([]string{"f"}, "f", "follow"),
		bk([]string{"/"}, "/", "search"),
		bk([]string{"1", "2", "3", "4"}, "1-4", "level"),
		bk([]string{"t"}, "t", "hide ticks"),
		bk([]string{"j", "k"}, "j/k", "scroll"),
		bk([]string{"y"}, "y", "yank"),
		bk([]string{"o"}, "o", "$PAGER"),
	}},
	{short: []key.Binding{
		bk([]string{"up", "down"}, "↑↓", "select"),
		bk([]string{"enter", "e"}, "⏎/e", "edit"),
		bk([]string{"ctrl+s"}, "⌃s", "save"),
		bk([]string{"r"}, "r", "reload"),
		bk([]string{"d"}, "d", "defaults"),
		bk([]string{"v"}, "v", "view file"),
	}},
	{short: []key.Binding{
		bk([]string{"up", "down"}, "↑↓", "select"),
		bk([]string{"enter"}, "⏎", "attach"),
		bk([]string{"s"}, "s", "start"),
		bk([]string{"S"}, "S", "stop"),
		bk([]string{"r"}, "r", "restart"),
		bk([]string{"n"}, "n", "new"),
		bk([]string{"d"}, "d", "set default"),
	}},
}
