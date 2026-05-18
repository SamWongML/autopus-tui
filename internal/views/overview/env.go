package overview

import (
	"math"
	"strings"

	"autopus-tui/internal/theme"
	"autopus-tui/internal/ui"
)

func renderEnvCard(w, h int) string {
	cpuS := genSpark(28, 12, 4)
	memS := genSpark(28, 32, 4)
	rows := []struct {
		label, val, color string
		pct               int
		s                 []float64
	}{
		{"cpu", "18%", theme.Accent, 18, cpuS},
		{"memory", "6.5G / 16G", theme.Info, 41, memS},
		{"disk · ~/autopus_workspaces", "18.6G / 30G", theme.Violet, 62, memS},
		{"network · out", "124 kbps", theme.OK, 8, cpuS},
	}
	var b strings.Builder
	for _, r := range rows {
		head := ui.JoinRight(theme.SDim.Render(r.label), theme.SText.Render(r.val), w)
		barW := ui.Min(20, w/2)
		sparkW := ui.Min(20, w-barW-2)
		row := ui.Bar(float64(r.pct), barW, r.color) + "  " + ui.Spark(r.s, sparkW, r.color)
		b.WriteString(head + "\n" + row + "\n")
	}
	b.WriteString(ui.Dashed(w) + "\n")
	b.WriteString(theme.SFaint.Render("host ") + theme.SDim.Render("luna.local") + "   " +
		theme.SFaint.Render("os ") + theme.SDim.Render("darwin · arm64") + "   " +
		theme.SFaint.Render("battery ") + theme.SOK.Render("94%"))
	_ = h
	return b.String()
}

// genSpark builds a deterministic sine sparkline.
func genSpark(n int, mid, amp float64) []float64 {
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = mid + amp*math.Sin(float64(i)/2.0)
	}
	return out
}
