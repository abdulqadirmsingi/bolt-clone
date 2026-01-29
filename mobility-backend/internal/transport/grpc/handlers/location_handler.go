package handlers

import (
	"context"
	"io"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/fireball"
	"github.com/zeng/mobility-backend/internal/domain/location"
	"github.com/zeng/mobility-backend/internal/proto/gen"
)

// LocationHandler implements gen.LocationServiceServer (streaming location RPCs).
type LocationHandler struct {
	gen.UnimplementedLocationServiceServer
	locSvc   *location.Service
	fireball *fireball.Service
	pub      fireball.LocationPublisher
}

func NewLocationHandler(locSvc *location.Service, fireballSvc *fireball.Service, pub fireball.LocationPublisher) *LocationHandler {
	return &LocationHandler{locSvc: locSvc, fireball: fireballSvc, pub: pub}
}

// StreamDriverLocation: driver sends GPS stream; Fireball filters and we push only significant updates.
func (h *LocationHandler) StreamDriverLocation(stream gen.LocationService_StreamDriverLocationServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		ctx := stream.Context()
		if err := h.fireball.OnDriverLocation(ctx, req.DriverId, req.Lat, req.Lng, req.HeadingDegrees); err != nil {
			continue
		}
		_ = stream.Send(&gen.DriverLocationAck{
			Accepted:           true,
			ServerTimestampMs: time.Now().UnixMilli(),
		})
	}
}

// SubscribeDriverLocation: rider subscribes to driver's location stream (push-based).
func (h *LocationHandler) SubscribeDriverLocation(req *gen.SubscribeDriverLocationRequest, stream gen.LocationService_SubscribeDriverLocationServer) error {
	ctx := stream.Context()
	ch, cancel := h.pub.SubscribeToDriver(ctx, req.DriverId)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case u, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(&gen.DriverLocationUpdate{
				DriverId:        u.DriverID,
				Lat:             u.Lat,
				Lng:             u.Lng,
				HeadingDegrees:  u.Heading,
				UpdatedAtMs:     u.UpdatedAt,
			}); err != nil {
				return err
			}
		}
	}
}
