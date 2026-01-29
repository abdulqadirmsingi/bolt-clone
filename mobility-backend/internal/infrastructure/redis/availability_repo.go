package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeng/mobility-backend/internal/domain/driver"
)

const (
	availabilityKeyPrefix   = "driver:availability:"
	availabilityTTL         = 15 * time.Minute
	heartbeatSortedSetKey   = "drivers:heartbeat" // ZSET: driverID -> last_heartbeat_ms (for timeout scan)
)

type availabilityPayload struct {
	Status        string `json:"status"`
	LastHeartbeat int64  `json:"last_heartbeat_ms"`
}

// AvailabilityRepo implements driver.AvailabilityStore using Redis.
type AvailabilityRepo struct {
	client *Client
}

func NewAvailabilityRepo(client *Client) *AvailabilityRepo {
	return &AvailabilityRepo{client: client}
}

func (r *AvailabilityRepo) key(driverID string) string {
	return availabilityKeyPrefix + driverID
}

func (r *AvailabilityRepo) SetStatus(ctx context.Context, driverID string, status driver.DriverStatus) error {
	data, err := r.getPayload(ctx, driverID)
	if err != nil {
		data = &availabilityPayload{}
	}
	data.Status = string(status)
	if err := r.setPayload(ctx, driverID, data); err != nil {
		return err
	}
	// Keep sorted set for heartbeat timeout: ONLINE/ON_TRIP = track, OFFLINE = remove
	if driver.IsAvailableForDiscovery(status) {
		r.client.ZAdd(ctx, heartbeatSortedSetKey, redis.Z{Score: float64(time.Now().UnixMilli()), Member: driverID})
	} else {
		r.client.ZRem(ctx, heartbeatSortedSetKey, driverID)
	}
	return nil
}

func (r *AvailabilityRepo) GetStatus(ctx context.Context, driverID string) (driver.DriverStatus, error) {
	data, err := r.getPayload(ctx, driverID)
	if err != nil {
		return driver.DriverOffline, nil
	}
	switch data.Status {
	case string(driver.DriverOnline):
		return driver.DriverOnline, nil
	case string(driver.DriverOnTrip):
		return driver.DriverOnTrip, nil
	default:
		return driver.DriverOffline, nil
	}
}

func (r *AvailabilityRepo) SetLastHeartbeat(ctx context.Context, driverID string, atUnixMs int64) error {
	data, err := r.getPayload(ctx, driverID)
	if err != nil {
		data = &availabilityPayload{}
	}
	data.LastHeartbeat = atUnixMs
	if err := r.setPayload(ctx, driverID, data); err != nil {
		return err
	}
	// Only track heartbeat in ZSET when driver is available (ONLINE/ON_TRIP)
	if data.Status == string(driver.DriverOnline) || data.Status == string(driver.DriverOnTrip) {
		r.client.ZAdd(ctx, heartbeatSortedSetKey, redis.Z{Score: float64(atUnixMs), Member: driverID})
	}
	return nil
}

func (r *AvailabilityRepo) GetLastHeartbeat(ctx context.Context, driverID string) (int64, error) {
	data, err := r.getPayload(ctx, driverID)
	if err != nil {
		return 0, nil
	}
	return data.LastHeartbeat, nil
}

func (r *AvailabilityRepo) getPayload(ctx context.Context, driverID string) (*availabilityPayload, error) {
	data, err := r.client.Get(ctx, r.key(driverID)).Bytes()
	if err != nil {
		return nil, err
	}
	var p availabilityPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AvailabilityRepo) setPayload(ctx context.Context, driverID string, p *availabilityPayload) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(driverID), data, availabilityTTL).Err()
}

// DriverIDsWithHeartbeatBefore returns driver IDs whose last heartbeat is before the given unix ms.
// Used by heartbeat checker to mark them OFFLINE.
func (r *AvailabilityRepo) DriverIDsWithHeartbeatBefore(ctx context.Context, beforeUnixMs int64) ([]string, error) {
	ids, err := r.client.ZRangeByScore(ctx, heartbeatSortedSetKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(beforeUnixMs, 10),
	}).Result()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// RemoveFromHeartbeatSet removes driver from the heartbeat ZSET (after marking OFFLINE).
func (r *AvailabilityRepo) RemoveFromHeartbeatSet(ctx context.Context, driverID string) error {
	return r.client.ZRem(ctx, heartbeatSortedSetKey, driverID).Err()
}

var _ driver.AvailabilityStore = (*AvailabilityRepo)(nil)
