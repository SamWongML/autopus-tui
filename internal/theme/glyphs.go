package theme

// StateMeta describes a session state's visual identity. SpinnerOn marks states
// that animate with a braille frame instead of using the static glyph.
type StateMeta struct {
	Glyph     string
	Label     string
	Color     string
	Sort      int
	SpinnerOn bool
}

// State returns the meta for a state name, falling back to "idle".
func State(name string) StateMeta {
	if m, ok := States[name]; ok {
		return m
	}
	return States["idle"]
}

var States = map[string]StateMeta{
	"needs_input": {Glyph: "◆", Label: "needs input", Color: Warn, Sort: 0},
	"working":     {Glyph: "◐", Label: "working", Color: Info, Sort: 1, SpinnerOn: true},
	"running":     {Glyph: "●", Label: "running", Color: OK, Sort: 2},
	"paused":      {Glyph: "‖", Label: "paused", Color: TextDim, Sort: 3},
	"idle":        {Glyph: "○", Label: "idle", Color: TextFaint, Sort: 4},
	"completed":   {Glyph: "✓", Label: "completed", Color: OK, Sort: 5},
	"failed":      {Glyph: "✕", Label: "failed", Color: Err, Sort: 6},
}

type PriorityMeta struct {
	Glyph string
	Color string
}

var Priorities = map[string]PriorityMeta{
	"urgent": {Glyph: "◆", Color: Err},
	"high":   {Glyph: "◆", Color: Warn},
	"medium": {Glyph: "◇", Color: TextDim},
	"low":    {Glyph: "◇", Color: TextFaint},
}

var Statuses = map[string]string{
	"backlog":     TextFaint,
	"todo":        TextDim,
	"in_progress": Info,
	"in_review":   Violet,
	"done":        OK,
	"blocked":     Err,
	"cancelled":   TextFaint,
}

var LogLevels = map[string]string{
	"INFO":  TextDim,
	"DEBUG": TextFaint,
	"WARN":  Warn,
	"ERROR": Err,
}

// SpinnerFrames are the braille dots used as the spinner animation. Same
// vocabulary Bubble Tea and Ratatui use.
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
