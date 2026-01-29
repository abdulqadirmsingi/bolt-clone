package redis

import (
	"context"

	"github.com/zeng/mobility-backend/internal/domain/driver"
)

// NoopPresenceNotifier implements driver.PresenceNotifier with no-op.
// Replace with Redis Pub/Sub or message bus for instant availability push.
type NoopPresenceNotifier struct{}

func (NoopPresenceNotifier) NotifyAvailabilityChange(ctx context.Context, driverID string, status driver.DriverStatus) error {
	return nil
}

var _ driver.PresenceNotifier = (*NoopPresenceNotifier)(nil)
