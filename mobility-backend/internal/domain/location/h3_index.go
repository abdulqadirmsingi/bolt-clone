package location

import (
	"github.com/uber/h3-go/v4"
)

// H3Index wraps Uber H3 for spatial partitioning.
// Why hexagons > squares:
// - No corner bias: distance from center to edge is uniform (squares have 4 corners farther).
// - 6 neighbors per cell; adjacency is consistent.
// - K-ring (GridDisk) gives symmetric, predictable neighborhoods (O(K²) cells).
// - H3 is hierarchical: same index at different resolutions aligns.
type H3Index struct {
	resolution int // 0–15; 9 ≈ 0.1 km², 10 ≈ 0.03 km²
	defaultK   int // K-ring layers for "nearby drivers"
}

func NewH3Index(resolution, defaultK int) *H3Index {
	if resolution < 0 || resolution > 15 {
		resolution = 9
	}
	if defaultK < 1 {
		defaultK = 2
	}
	return &H3Index{resolution: resolution, defaultK: defaultK}
}

// LatLngToCell returns the H3 cell containing the given point at configured resolution.
func (h *H3Index) LatLngToCell(lat, lng float64) (h3.Cell, error) {
	return h3.LatLngToCell(h3.NewLatLng(lat, lng), h.resolution)
}

// CellToLatLng returns the center of the cell.
func (h *H3Index) CellToLatLng(c h3.Cell) (h3.LatLng, error) {
	return c.LatLng()
}

// GridDisk returns the cell and all cells within K steps (K-ring / grid disk).
// Complexity: O(K²) cells returned; matching drivers is O(K² + M) where M = drivers in those cells.
// K=0: 1 cell, K=1: 7 cells, K=2: 19 cells, K=3: 37 cells.
func (h *H3Index) GridDisk(center h3.Cell, k int) ([]h3.Cell, error) {
	if k < 0 {
		k = h.defaultK
	}
	return center.GridDisk(k)
}

// KRingStrings returns H3 index strings for Redis/PostgreSQL storage.
// Store driver_id -> h3_index in Redis; query by set of h3 indices for nearby drivers.
func (h *H3Index) KRingStrings(lat, lng float64, k int) ([]string, error) {
	center, err := h.LatLngToCell(lat, lng)
	if err != nil {
		return nil, err
	}
	cells, err := h.GridDisk(center, k)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(cells))
	for _, c := range cells {
		if c.IsValid() {
			out = append(out, c.String())
		}
	}
	return out, nil
}

// DefaultKRing returns K-ring indices for (lat, lng) using default K.
func (h *H3Index) DefaultKRing(lat, lng float64) ([]string, error) {
	return h.KRingStrings(lat, lng, h.defaultK)
}

func (h *H3Index) Resolution() int { return h.resolution }
func (h *H3Index) DefaultK() int   { return h.defaultK }
