package driver

import "context"

// DriverStatus is the driver presence/availability state.
// State machine: OFFLINE → ONLINE → ON_TRIP → ONLINE → OFFLINE.
type DriverStatus string

const (
	DriverOffline DriverStatus = "OFFLINE"
	DriverOnline  DriverStatus = "ONLINE"
	DriverOnTrip  DriverStatus = "ON_TRIP"
)

// AvailabilityStore persists driver status in Redis (key: driver:availability:{id}).
// Used for presence checks and heartbeat timeout.
type AvailabilityStore interface {
	SetStatus(ctx context.Context, driverID string, status DriverStatus) error
	GetStatus(ctx context.Context, driverID string) (DriverStatus, error)
	SetLastHeartbeat(ctx context.Context, driverID string, atUnixMs int64) error
	GetLastHeartbeat(ctx context.Context, driverID string) (int64, error)
}

// PresenceNotifier fires when a driver's availability changes (e.g. OFFLINE, ONLINE).
// Implemented by infra; used to push availability updates to riders/apps.
type PresenceNotifier interface {
	NotifyAvailabilityChange(ctx context.Context, driverID string, status DriverStatus) error
}

// DriverGeoRemover removes a driver from the spatial index (H3). Implemented by location.Service.
type DriverGeoRemover interface {
	DriverOffline(ctx context.Context, driverID string) error
}

// AvailabilityService applies availability state transitions and notifies presence.
type AvailabilityService struct {
	store   AvailabilityStore
	notify  PresenceNotifier
	geoRem  DriverGeoRemover
}

func NewAvailabilityService(store AvailabilityStore, notify PresenceNotifier, geoRem DriverGeoRemover) *AvailabilityService {
	return &AvailabilityService{store: store, notify: notify, geoRem: geoRem}
}

// SetAvailability transitions driver to the given status. On OFFLINE, removes from H3 and notifies.
func (s *AvailabilityService) SetAvailability(ctx context.Context, driverID string, status DriverStatus) error {
	if err := s.store.SetStatus(ctx, driverID, status); err != nil {
		return err
	}
	if status == DriverOffline {
		_ = s.geoRem.DriverOffline(ctx, driverID)
		_ = s.notify.NotifyAvailabilityChange(ctx, driverID, status)
	} else {
		_ = s.notify.NotifyAvailabilityChange(ctx, driverID, status)
	}
	return nil
}

// RecordHeartbeat updates last heartbeat time for the driver (used for timeout → OFFLINE).
func (s *AvailabilityService) RecordHeartbeat(ctx context.Context, driverID string, atUnixMs int64) error {
	return s.store.SetLastHeartbeat(ctx, driverID, atUnixMs)
}

// GetStatus returns current driver status.
func (s *AvailabilityService) GetStatus(ctx context.Context, driverID string) (DriverStatus, error) {
	return s.store.GetStatus(ctx, driverID)
}

// IsAvailableForDiscovery returns true if driver appears in "nearby" (H3) and availability stream.
func IsAvailableForDiscovery(status DriverStatus) bool {
	return status == DriverOnline || status == DriverOnTrip
}
