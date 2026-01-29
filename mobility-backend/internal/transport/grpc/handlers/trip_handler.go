package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/trip"
	"github.com/zeng/mobility-backend/internal/proto/gen"
)

// TripHandler implements gen.TripServiceServer. No business logic in transport—delegates to trip.Service.
type TripHandler struct {
	gen.UnimplementedTripServiceServer
	tripSvc *trip.Service
}

func NewTripHandler(tripSvc *trip.Service) *TripHandler {
	return &TripHandler{tripSvc: tripSvc}
}

func stopsFromProto(ss []*gen.Stop) []trip.Stop {
	out := make([]trip.Stop, 0, len(ss))
	for _, s := range ss {
		out = append(out, trip.Stop{Lat: s.Lat, Lng: s.Lng, Label: s.Label, Sequence: int(s.Sequence)})
	}
	return out
}

func stopsToProto(ss []trip.Stop) []*gen.Stop {
	out := make([]*gen.Stop, 0, len(ss))
	for _, s := range ss {
		out = append(out, &gen.Stop{Lat: s.Lat, Lng: s.Lng, Label: s.Label, Sequence: int32(s.Sequence)})
	}
	return out
}

func (h *TripHandler) CreateTrip(ctx context.Context, req *gen.CreateTripRequest) (*gen.CreateTripResponse, error) {
	if len(req.Stops) == 0 {
		return nil, nil // or validation error
	}
	tripID := fmt.Sprintf("trip-%d", time.Now().UnixNano())
	driverID := "driver-1" // TODO: assign from nearby driver pool
	e := &trip.Entity{
		ID:              tripID,
		RiderID:         req.RiderId,
		DriverID:        driverID,
		Status:          "requested",
		OrderedStops:   stopsFromProto(req.Stops),
		CurrentStopIndex: 0,
		ETASeconds:     600,
	}
	if len(e.OrderedStops) > 0 {
		e.PickupLat = e.OrderedStops[0].Lat
		e.PickupLng = e.OrderedStops[0].Lng
	}
	if len(e.OrderedStops) > 1 {
		e.DropoffLat = e.OrderedStops[len(e.OrderedStops)-1].Lat
		e.DropoffLng = e.OrderedStops[len(e.OrderedStops)-1].Lng
	}
	if err := h.tripSvc.Create(ctx, e); err != nil {
		return nil, err
	}
	return &gen.CreateTripResponse{
		TripId:       tripID,
		DriverId:     driverID,
		OrderedStops: stopsToProto(e.OrderedStops),
		EtaSeconds:   int32(e.ETASeconds),
	}, nil
}

func (h *TripHandler) GetTrip(ctx context.Context, req *gen.GetTripRequest) (*gen.GetTripResponse, error) {
	e, err := h.tripSvc.GetByID(ctx, req.TripId)
	if err != nil || e == nil {
		return nil, err
	}
	return &gen.GetTripResponse{
		TripId:           e.ID,
		Status:           e.Status,
		DriverId:         e.DriverID,
		OrderedStops:     stopsToProto(e.OrderedStops),
		CurrentStopIndex: int32(e.CurrentStopIndex),
		EtaSeconds:      int32(e.ETASeconds),
	}, nil
}

func (h *TripHandler) StreamTripUpdates(req *gen.StreamTripUpdatesRequest, stream gen.TripService_StreamTripUpdatesServer) error {
	ctx := stream.Context()
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	send := func() error {
		e, err := h.tripSvc.GetByID(ctx, req.TripId)
		if err != nil || e == nil {
			return err
		}
		return stream.Send(&gen.TripUpdate{
			TripId:           e.ID,
			Status:           e.Status,
			CurrentStopIndex: int32(e.CurrentStopIndex),
			OrderedStops:     stopsToProto(e.OrderedStops),
			EtaSeconds:       int32(e.ETASeconds),
			UpdatedAtMs:      time.Now().UnixMilli(),
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
