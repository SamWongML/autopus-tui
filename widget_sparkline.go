package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// sparkBars are the 8 one-cell vertical-bar glyphs from ▁ to █.
var sparkBars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// sparkline renders `values` as a single-line bar series of exactly `width`
// cells. The series is resampled (block-averaged or nearest-neighbor) when
// width ≠ len(values); peaks normalize to the tallest bar. Output is wrapped
// in sAccent for the warm-palette amber tone.
func sparkline(values []float64, width int) string {
	if width <= 0 || len(values) == 0 {
		return ""
	}
	samples := resample(values, width)
	maxV := samples[0]
	minV := samples[0]
	for _, v := range samples {
		if v > maxV {
			maxV = v
		}
		if v < minV {
			minV = v
		}
	}
	span := maxV - minV
	if span == 0 {
		span = 1
	}
	var b strings.Builder
	b.Grow(width * 4)
	for _, v := range samples {
		idx := int((v-minV)/span*7 + 0.5)
		if idx < 0 {
			idx = 0
		}
		if idx > 7 {
			idx = 7
		}
		b.WriteRune(sparkBars[idx])
	}
	return sAccent.Render(b.String())
}

// resample returns a slice of length `n` derived from `src`. If len(src) == n
// the slice is returned unchanged. If shrinking, contiguous src buckets are
// averaged; if growing, src values are nearest-neighbor stretched. Either way
// the first and last samples are anchored to the start/end of src.
func resample(src []float64, n int) []float64 {
	if n == len(src) || n <= 0 {
		return src
	}
	out := make([]float64, n)
	if n == 1 {
		var sum float64
		for _, v := range src {
			sum += v
		}
		out[0] = sum / float64(len(src))
		return out
	}
	for i := 0; i < n; i++ {
		start := i * len(src) / n
		end := (i + 1) * len(src) / n
		if end <= start {
			end = start + 1
		}
		if end > len(src) {
			end = len(src)
		}
		var sum float64
		for j := start; j < end; j++ {
			sum += src[j]
		}
		out[i] = sum / float64(end-start)
	}
	return out
}

// sparklinePane renders a `pane(title, hint, …)` whose body is a sparkline
// row followed by a stats row (min · avg · max · now). The pane sizes its
// own inner content from `width` and `height`; `unit` is a short suffix
// shown after each stat (e.g. "/h", " l/min", "").
func sparklinePane(title, hint string, values []float64, unit string, width, height int) string {
	inner := width - 2
	innerH := height - 2
	if inner < 8 {
		inner = 8
	}
	if innerH < 1 {
		innerH = 1
	}
	sparkW := inner - 4
	if sparkW < 6 {
		sparkW = inner - 2
	}
	if sparkW < 2 {
		sparkW = 2
	}

	line := sparkline(values, sparkW)

	var sum, minV, maxV float64
	if len(values) > 0 {
		minV, maxV = values[0], values[0]
	}
	for _, v := range values {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
		sum += v
	}
	avg := 0.0
	if len(values) > 0 {
		avg = sum / float64(len(values))
	}
	now := 0.0
	if len(values) > 0 {
		now = values[len(values)-1]
	}

	stat := func(label string, v float64) string {
		return sDim.Render(label+" ") + sFg1.Render(fmtSparkNum(v)+unit)
	}
	sep := sDim.Render("   ")
	full := strings.Join([]string{
		stat("min", minV),
		stat("avg", avg),
		stat("max", maxV),
		stat("now", now),
	}, sep)
	short := strings.Join([]string{
		stat("min", minV),
		stat("max", maxV),
		stat("now", now),
	}, sep)
	tight := stat("now", now)
	// Pick the widest variant that fits inside the 2-cell-padded inner band.
	bandW := inner - 4
	stats := full
	switch {
	case lipgloss.Width(full) <= bandW:
		stats = full
	case lipgloss.Width(short) <= bandW:
		stats = short
	default:
		stats = tight
	}

	// Vertical layout inside innerH: blank · sparkline · blank · stats, with
	// extra blank padding distributed above/below for taller panes.
	var rows []string
	if innerH >= 4 {
		topPad := (innerH - 3) / 2
		for i := 0; i < topPad; i++ {
			rows = append(rows, "")
		}
		rows = append(rows, "  "+line)
		rows = append(rows, "")
		rows = append(rows, "  "+stats)
		for len(rows) < innerH {
			rows = append(rows, "")
		}
	} else if innerH >= 2 {
		rows = []string{"  " + line, "  " + stats}
	} else {
		rows = []string{"  " + line}
	}

	body := lipgloss.NewStyle().Width(inner).Render(strings.Join(rows, "\n"))
	return pane(title, hint, body, width, height, false)
}

// fmtSparkNum formats a number compactly. Integers print plain; values with
// a fractional component print with one decimal. Used by sparklinePane stat
// rows so "min 0 · avg 7.3 · max 14 · now 6" lines up regardless of scale.
func fmtSparkNum(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%.1f", v)
}
