package fireball

import (
	"context"
	"time"
)

// LocationUpdater writes driver location to Redis (H3 + geo). Implemented by location.Service.
type LocationUpdater interface {
	UpdateDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64, updatedAt int64) error
}

// FireballService listens to location updates, compares with previous state,
// and only when movement is significant: updates Redis and pushes to gRPC streams.
// Raw GPS is never written to PostgreSQL.
type Service struct {
	threshold *Threshold
	loc      LocationUpdater
	pub      LocationPublisher
}

func NewService(threshold *Threshold, loc LocationUpdater, pub LocationPublisher) *Service {
	return &Service{threshold: threshold, loc: loc, pub: pub}
}

// OnDriverLocation is called for every incoming GPS sample (e.g. from gRPC stream).
// It filters by threshold, then writes to Redis (via location service) and publishes to subscribers.
func (s *Service) OnDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64) error {
	updatedAt := time.Now().UnixMilli()

	if !s.threshold.ShouldPush(driverID, lat, lng, heading, updatedAt) {
		return nil // no push; bandwidth saved
	}

	if err := s.loc.UpdateDriverLocation(ctx, driverID, lat, lng, heading, updatedAt); err != nil {
		return err
	}
	return s.pub.PublishDriverLocation(ctx, driverID, lat, lng, heading, updatedAt)
}
