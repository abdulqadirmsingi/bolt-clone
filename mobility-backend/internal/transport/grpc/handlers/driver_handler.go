package handlers

import (
	"context"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/driver"
	"github.com/zeng/mobility-backend/internal/domain/location"
	"github.com/zeng/mobility-backend/internal/proto/gen"
)

// DriverHandler implements gen.DriverServiceServer (nearby drivers, availability stream, SetAvailability, Heartbeat).
type DriverHandler struct {
	gen.UnimplementedDriverServiceServer
	locSvc     *location.Service
	availSvc   *driver.AvailabilityService
}

func NewDriverHandler(locSvc *location.Service, availSvc *driver.AvailabilityService) *DriverHandler {
	return &DriverHandler{locSvc: locSvc, availSvc: availSvc}
}

func (h *DriverHandler) filterByAvailability(ctx context.Context, snapshots []location.DriverSnapshot) ([]location.DriverSnapshot, error) {
	out := make([]location.DriverSnapshot, 0, len(snapshots))
	for _, s := range snapshots {
		status, err := h.availSvc.GetStatus(ctx, s.DriverID)
		if err != nil || !driver.IsAvailableForDiscovery(status) {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}

func (h *DriverHandler) GetNearbyDrivers(ctx context.Context, req *gen.GetNearbyDriversRequest) (*gen.GetNearbyDriversResponse, error) {
	k := int(req.KRing)
	if k <= 0 {
		k = 0
	}
	snapshots, err := h.locSvc.NearbyDrivers(ctx, req.Lat, req.Lng, k)
	if err != nil {
		return nil, err
	}
	snapshots, _ = h.filterByAvailability(ctx, snapshots)
	drivers := make([]*gen.DriverSummary, 0, len(snapshots))
	for _, s := range snapshots {
		drivers = append(drivers, &gen.DriverSummary{
			DriverId:       s.DriverID,
			Lat:            s.Lat,
			Lng:            s.Lng,
			HeadingDegrees: s.Heading,
			UpdatedAtMs:    s.UpdatedAt,
		})
	}
	return &gen.GetNearbyDriversResponse{Drivers: drivers}, nil
}

// StreamDriverAvailability: push-based. Sends initial snapshot, then periodic snapshots (no client polling).
func (h *DriverHandler) StreamDriverAvailability(req *gen.StreamDriverAvailabilityRequest, stream gen.DriverService_StreamDriverAvailabilityServer) error {
	ctx := stream.Context()
	k := int(req.KRing)
	if k <= 0 {
		k = 0
	}
	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()
	send := func() error {
		snapshots, err := h.locSvc.NearbyDrivers(ctx, req.Lat, req.Lng, k)
		if err != nil {
			return err
		}
		snapshots, _ = h.filterByAvailability(ctx, snapshots)
		drivers := make([]*gen.DriverSummary, 0, len(snapshots))
		for _, s := range snapshots {
			drivers = append(drivers, &gen.DriverSummary{
				DriverId:       s.DriverID,
				Lat:            s.Lat,
				Lng:            s.Lng,
				HeadingDegrees: s.Heading,
				UpdatedAtMs:    s.UpdatedAt,
			})
		}
		return stream.Send(&gen.DriverAvailabilityUpdate{
			Drivers:     drivers,
			UpdatedAtMs: time.Now().UnixMilli(),
		})
	}
	if err := send(); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			if err := send(); err != nil {
				return err
			}
		}
	}
}

// SetAvailability transitions driver to ONLINE / ON_TRIP / OFFLINE.
func (h *DriverHandler) SetAvailability(ctx context.Context, req *gen.SetAvailabilityRequest) (*gen.SetAvailabilityResponse, error) {
	status := protoStatusToDomain(req.Status)
	if err := h.availSvc.SetAvailability(ctx, req.DriverId, status); err != nil {
		return &gen.SetAvailabilityResponse{Ok: false}, nil
	}
	return &gen.SetAvailabilityResponse{Ok: true}, nil
}

// Heartbeat updates last heartbeat; used by server to mark OFFLINE if no heartbeat within timeout.
func (h *DriverHandler) Heartbeat(ctx context.Context, req *gen.HeartbeatRequest) (*gen.HeartbeatResponse, error) {
	now := time.Now().UnixMilli()
	if req.ClientTimestampMs > 0 {
		now = req.ClientTimestampMs
	}
	if err := h.availSvc.RecordHeartbeat(ctx, req.DriverId, now); err != nil {
		return &gen.HeartbeatResponse{Ok: false, ServerTimestampMs: time.Now().UnixMilli()}, nil
	}
	return &gen.HeartbeatResponse{Ok: true, ServerTimestampMs: time.Now().UnixMilli()}, nil
}

func protoStatusToDomain(s gen.DriverStatus) driver.DriverStatus {
	switch s {
	case gen.DriverStatus_DRIVER_STATUS_ONLINE:
		return driver.DriverOnline
	case gen.DriverStatus_DRIVER_STATUS_ON_TRIP:
		return driver.DriverOnTrip
	case gen.DriverStatus_DRIVER_STATUS_OFFLINE:
		return driver.DriverOffline
	default:
		return driver.DriverOffline
	}
}
