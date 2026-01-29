package app

import (
	"context"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/driver"
	"github.com/zeng/mobility-backend/internal/infrastructure/redis"
)

// HeartbeatChecker periodically marks drivers as OFFLINE when no heartbeat received within timeout.
type HeartbeatChecker struct {
	availRepo *redis.AvailabilityRepo
	availSvc  *driver.AvailabilityService
	timeout   time.Duration
	interval  time.Duration
}

func NewHeartbeatChecker(availRepo *redis.AvailabilityRepo, availSvc *driver.AvailabilityService, timeoutSec, intervalSec int) *HeartbeatChecker {
	if intervalSec <= 0 {
		intervalSec = 30
	}
	return &HeartbeatChecker{
		availRepo: availRepo,
		availSvc:  availSvc,
		timeout:   time.Duration(timeoutSec) * time.Second,
		interval:  time.Duration(intervalSec) * time.Second,
	}
}

// Run blocks and runs the checker loop until ctx is cancelled.
func (c *HeartbeatChecker) Run(ctx context.Context) {
	tick := time.NewTicker(c.interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			c.markStaleOffline(ctx)
		}
	}
}

func (c *HeartbeatChecker) markStaleOffline(ctx context.Context) {
	beforeMs := time.Now().Add(-c.timeout).UnixMilli()
	ids, err := c.availRepo.DriverIDsWithHeartbeatBefore(ctx, beforeMs)
	if err != nil || len(ids) == 0 {
		return
	}
	for _, driverID := range ids {
		_ = c.availSvc.SetAvailability(ctx, driverID, driver.DriverOffline)
	}
}
