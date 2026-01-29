package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/location"
)

const (
	driverKeyPrefix = "driver:"
	driverTTL       = 10 * time.Minute
	h3CellPrefix    = "h3:"
)

// GeoRepo implements location.GeoRepository using Redis.
// Drivers: driver:{id} -> JSON snapshot; h3:{cell} -> set of driver IDs.
type GeoRepo struct {
	client *Client
}

type driverPayload struct {
	DriverID  string  `json:"driver_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Heading   float64 `json:"heading"`
	H3Index   string  `json:"h3_index"`
	UpdatedAt int64   `json:"updated_at"`
}

func (r *GeoRepo) UpsertDriver(ctx context.Context, driverID, h3Index string, lat, lng, heading float64, updatedAt int64) error {
	key := driverKeyPrefix + driverID
	payload := driverPayload{DriverID: driverID, Lat: lat, Lng: lng, Heading: heading, H3Index: h3Index, UpdatedAt: updatedAt}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	pipe := r.client.Pipeline()
	// Remove from previous H3 cell so driver appears in only one cell (no ghost entries).
	prev, _ := r.DriverByID(ctx, driverID)
	if prev != nil && prev.H3Index != "" && prev.H3Index != h3Index {
		pipe.SRem(ctx, h3CellPrefix+prev.H3Index, driverID)
	}
	pipe.Set(ctx, key, data, driverTTL)
	if h3Index != "" {
		cellKey := h3CellPrefix + h3Index
		pipe.SAdd(ctx, cellKey, driverID)
		pipe.Expire(ctx, cellKey, driverTTL)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (r *GeoRepo) DriversInCells(ctx context.Context, h3Indices []string) ([]location.DriverSnapshot, error) {
	seen := make(map[string]struct{})
	var out []location.DriverSnapshot
	for _, idx := range h3Indices {
		ids, err := r.client.SMembers(ctx, h3CellPrefix+idx).Result()
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			snap, err := r.DriverByID(ctx, id)
			if err != nil || snap == nil {
				continue
			}
			out = append(out, *snap)
		}
	}
	return out, nil
}

func (r *GeoRepo) DriverByID(ctx context.Context, driverID string) (*location.DriverSnapshot, error) {
	data, err := r.client.Get(ctx, driverKeyPrefix+driverID).Bytes()
	if err != nil {
		return nil, err
	}
	var p driverPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &location.DriverSnapshot{
		DriverID:  p.DriverID,
		Lat:       p.Lat,
		Lng:       p.Lng,
		Heading:   p.Heading,
		H3Index:   p.H3Index,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

func (r *GeoRepo) RemoveDriver(ctx context.Context, driverID string) error {
	snap, err := r.DriverByID(ctx, driverID)
	if err != nil {
		return err
	}
	pipe := r.client.Pipeline()
	pipe.Del(ctx, driverKeyPrefix+driverID)
	if snap.H3Index != "" {
		pipe.SRem(ctx, h3CellPrefix+snap.H3Index, driverID)
	}
	_, err = pipe.Exec(ctx)
	return err
}

// NewGeoRepo returns a GeoRepository implementation (for dependency injection).
func NewGeoRepo(client *Client) *GeoRepo {
	return &GeoRepo{client: client}
}

// Ensure GeoRepo implements location.GeoRepository.
var _ location.GeoRepository = (*GeoRepo)(nil)
