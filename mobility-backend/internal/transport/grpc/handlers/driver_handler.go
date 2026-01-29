package handlers

import (
	"context"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/location"
	"github.com/zeng/mobility-backend/internal/proto/gen"
)

// DriverHandler implements gen.DriverServiceServer (nearby drivers, availability stream).
// No business logic here—delegates to location.Service.
type DriverHandler struct {
	gen.UnimplementedDriverServiceServer
	locSvc *location.Service
}

func NewDriverHandler(locSvc *location.Service) *DriverHandler {
	return &DriverHandler{locSvc: locSvc}
}

func (h *DriverHandler) GetNearbyDrivers(ctx context.Context, req *gen.GetNearbyDriversRequest) (*gen.GetNearbyDriversResponse, error) {
	k := int(req.KRing)
	if k <= 0 {
		k = 0 // service uses default
	}
	snapshots, err := h.locSvc.NearbyDrivers(ctx, req.Lat, req.Lng, k)
	if err != nil {
		return nil, err
	}
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
