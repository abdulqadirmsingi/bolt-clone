package location

import (
	"context"
)

// GeoRepository abstracts Redis/DB for driver locations and H3 index lookups.
type GeoRepository interface {
	UpsertDriver(ctx context.Context, driverID, h3Index string, lat, lng, heading float64, updatedAt int64) error
	DriversInCells(ctx context.Context, h3Indices []string) ([]DriverSnapshot, error)
	DriverByID(ctx context.Context, driverID string) (*DriverSnapshot, error)
	RemoveDriver(ctx context.Context, driverID string) error
}

type Service struct {
	h3    *H3Index
	geo   GeoRepository
}

func NewService(h3 *H3Index, geo GeoRepository) *Service {
	return &Service{h3: h3, geo: geo}
}

// UpdateDriverLocation updates Redis with driver position and H3 cell.
// Call this from Fireball after threshold check; do not write raw GPS to PostgreSQL.
func (s *Service) UpdateDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64, updatedAt int64) error {
	cell, err := s.h3.LatLngToCell(lat, lng)
	if err != nil {
		return err
	}
	return s.geo.UpsertDriver(ctx, driverID, cell.String(), lat, lng, heading, updatedAt)
}

// NearbyDrivers returns drivers in the K-ring around (lat, lng). O(K² + M).
func (s *Service) NearbyDrivers(ctx context.Context, lat, lng float64, k int) ([]DriverSnapshot, error) {
	if k <= 0 {
		k = s.h3.DefaultK()
	}
	indices, err := s.h3.KRingStrings(lat, lng, k)
	if err != nil {
		return nil, err
	}
	return s.geo.DriversInCells(ctx, indices)
}

// GetDriver returns current snapshot for one driver.
func (s *Service) GetDriver(ctx context.Context, driverID string) (*DriverSnapshot, error) {
	return s.geo.DriverByID(ctx, driverID)
}

// DriverOffline removes driver from spatial index (e.g. trip ended, app closed).
func (s *Service) DriverOffline(ctx context.Context, driverID string) error {
	return s.geo.RemoveDriver(ctx, driverID)
}
