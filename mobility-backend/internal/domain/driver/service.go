package driver

import (
	"context"

	"github.com/zeng/mobility-backend/internal/domain/location"
)

// NearbyDriverProvider returns nearby drivers (from Redis/H3). Implemented by location.Service.
type NearbyDriverProvider interface {
	NearbyDrivers(ctx context.Context, lat, lng float64, k int) ([]location.DriverSnapshot, error)
	GetDriver(ctx context.Context, driverID string) (*location.DriverSnapshot, error)
}

// Service holds driver business logic. No transport or infra here.
type Service struct {
	repo     Repository
	geo      NearbyDriverProvider
}

func NewService(repo Repository, geo NearbyDriverProvider) *Service {
	return &Service{repo: repo, geo: geo}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Entity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*Entity, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetNearby returns drivers in the K-ring around (lat, lng). O(K² + M).
func (s *Service) GetNearby(ctx context.Context, lat, lng float64, k int) ([]location.DriverSnapshot, error) {
	return s.geo.NearbyDrivers(ctx, lat, lng, k)
}

func (s *Service) GetLiveSnapshot(ctx context.Context, driverID string) (*location.DriverSnapshot, error) {
	return s.geo.GetDriver(ctx, driverID)
}
