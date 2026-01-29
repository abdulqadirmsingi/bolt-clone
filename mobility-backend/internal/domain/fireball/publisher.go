package fireball

import (
	"context"
)

// LocationPublisher pushes driver location updates to subscribers (e.g. gRPC stream manager).
// "At least once": backend may send same update to multiple subscribers; clients dedupe by sequence/ts.
type LocationPublisher interface {
	PublishDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64, updatedAt int64) error
	// Subscribe is used by gRPC handler to receive updates for a trip's driver and forward to client stream
	SubscribeToDriver(ctx context.Context, driverID string) (<-chan LocationUpdate, func())
}

// LocationUpdate is the payload pushed to subscribers.
type LocationUpdate struct {
	DriverID  string
	Lat       float64
	Lng       float64
	Heading   float64
	UpdatedAt int64
}
