package ui

// HitBox is a 1D click range on a single rendered row. ID identifies what the
// range corresponds to — a route ID, a filter name, a row index encoded as a
// string. Views compute these alongside rendering so the mouse handler can map
// a click back to a stable target without re-parsing the rendered string.
type HitBox struct {
	X1, X2 int
	ID     string
}

// Hits is a small helper: returns the ID of the first hit-box containing x, or
// "" if none match.
func Hits(boxes []HitBox, x int) string {
	for _, h := range boxes {
		if x >= h.X1 && x <= h.X2 {
			return h.ID
		}
	}
	return ""
}
