package fireball

import (
	"context"
	"time"
)

// LocationUpdater writes driver location to Redis (H3 + geo). Implemented by location.Service.
type LocationUpdater interface {
	UpdateDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64, updatedAt int64) error
}

// AvailabilityChecker returns whether the driver is available for discovery (ONLINE or ON_TRIP).
// If nil, location updates are always applied (backward compatible).
type AvailabilityChecker interface {
	IsAvailableForDiscovery(ctx context.Context, driverID string) (bool, error)
}

// FireballService listens to location updates, compares with previous state,
// and only when movement is significant: updates Redis and pushes to gRPC streams.
// Raw GPS is never written to PostgreSQL.
// When AvailabilityChecker is set, location is only written when driver is ONLINE or ON_TRIP.
type Service struct {
	threshold  *Threshold
	loc        LocationUpdater
	pub        LocationPublisher
	availCheck AvailabilityChecker
}

func NewService(threshold *Threshold, loc LocationUpdater, pub LocationPublisher) *Service {
	return &Service{threshold: threshold, loc: loc, pub: pub}
}

// SetAvailabilityChecker sets optional checker; only update H3 when driver is ONLINE/ON_TRIP.
func (s *Service) SetAvailabilityChecker(c AvailabilityChecker) {
	s.availCheck = c
}

// OnDriverLocation is called for every incoming GPS sample (e.g. from gRPC stream).
// It filters by threshold; only when significant: writes Redis and publishes. Returns whether a push was emitted.
// When AvailabilityChecker is set, skips update when driver is OFFLINE (no H3 write, no publish).
func (s *Service) OnDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64) (pushed bool, err error) {
	updatedAt := time.Now().UnixMilli()

	if s.availCheck != nil {
		ok, err := s.availCheck.IsAvailableForDiscovery(ctx, driverID)
		if err != nil || !ok {
			return false, nil
		}
	}

	if !s.threshold.ShouldPush(driverID, lat, lng, heading, updatedAt) {
		return false, nil // no push; bandwidth saved
	}

	if err := s.loc.UpdateDriverLocation(ctx, driverID, lat, lng, heading, updatedAt); err != nil {
		return false, err
	}
	if err := s.pub.PublishDriverLocation(ctx, driverID, lat, lng, heading, updatedAt); err != nil {
		return false, err
	}
	return true, nil
}
