package ui

// Col describes one column in a fluid table layout. Higher Priority value =
// less important = dropped first when the available width is too tight. Weight
// is the share of leftover space the column claims (0 = fixed at MinW).
type Col struct {
	ID       string
	Header   string
	MinW     int
	Weight   int
	Priority int
	AlignR   bool
}

// LayoutCols decides which columns survive at width w and how wide each is.
// It reserves a 1-cell gap between every pair of visible columns. Returns the
// per-column widths keyed by ID and the surviving cols in original order.
//
// Drop order: at each step, the col with the highest Priority number is
// removed; ties are broken by the wider MinW (the heavier column goes first).
// At least one column is always kept, even if it does not fit.
//
// Width distribution: every visible col starts at MinW. The leftover is split
// across cols proportional to Weight (integer division). Any remainder goes
// to the highest-Weight column.
func LayoutCols(cols []Col, w int) (widths map[string]int, visible []Col) {
	widths = map[string]int{}
	if len(cols) == 0 || w <= 0 {
		return widths, nil
	}

	// Indices into cols[] for the currently-kept set, in original order.
	keep := make([]int, len(cols))
	for i := range cols {
		keep[i] = i
	}

	minSum := func() int {
		t := 0
		for _, idx := range keep {
			t += cols[idx].MinW
		}
		return t
	}
	gaps := func() int {
		if len(keep) <= 1 {
			return 0
		}
		return len(keep) - 1
	}

	for len(keep) > 1 && minSum()+gaps() > w {
		drop := 0
		for i := 1; i < len(keep); i++ {
			a, b := cols[keep[i]], cols[keep[drop]]
			if a.Priority > b.Priority ||
				(a.Priority == b.Priority && a.MinW > b.MinW) {
				drop = i
			}
		}
		keep = append(keep[:drop], keep[drop+1:]...)
	}

	visible = make([]Col, 0, len(keep))
	for _, idx := range keep {
		visible = append(visible, cols[idx])
		widths[cols[idx].ID] = cols[idx].MinW
	}

	leftover := w - minSum() - gaps()
	if leftover <= 0 {
		return widths, visible
	}

	totalWeight := 0
	for _, c := range visible {
		totalWeight += c.Weight
	}
	if totalWeight == 0 {
		return widths, visible
	}

	used := 0
	maxWi, maxWv := -1, -1
	for i, c := range visible {
		if c.Weight <= 0 {
			continue
		}
		add := leftover * c.Weight / totalWeight
		widths[c.ID] += add
		used += add
		if c.Weight > maxWv {
			maxWv = c.Weight
			maxWi = i
		}
	}
	if maxWi >= 0 && used < leftover {
		widths[visible[maxWi].ID] += leftover - used
	}
	return widths, visible
}
