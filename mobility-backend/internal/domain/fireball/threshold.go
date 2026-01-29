package fireball

import (
	"sync"
	"time"

	"github.com/zeng/mobility-backend/internal/domain/location"
)

// Threshold decides when a GPS update is "significant" enough to push.
// Reduces bandwidth and server load by not pushing every raw sample.
type Threshold struct {
	minDistanceMeters float64
	minHeadingDegrees float64
	maxSilenceSeconds int
	throttleMs        int

	mu         sync.Mutex
	last       map[string]*lastState
	lastPurge  time.Time
}

type lastState struct {
	lat       float64
	lng       float64
	heading   float64
	updatedAt int64
	pushedAt  time.Time
}

func NewThreshold(minDistM, minHeadingDeg, maxSilenceSec, throttleMs int) *Threshold {
	t := &Threshold{
		minDistanceMeters: float64(minDistM),
		minHeadingDegrees: float64(minHeadingDeg),
		maxSilenceSeconds: maxSilenceSec,
		throttleMs:        throttleMs,
		last:              make(map[string]*lastState),
	}
	return t
}

// ShouldPush returns true if we should push this update (significant movement or keep-alive).
// Fireball calls this before writing to Redis and pushing to gRPC streams.
func (t *Threshold) ShouldPush(driverID string, lat, lng, heading float64, updatedAt int64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.maybePurge()
	now := time.Now()
	prev, ok := t.last[driverID]

	// First update: always push
	if !ok {
		t.last[driverID] = &lastState{lat, lng, heading, updatedAt, now}
		return true
	}

	// Throttle: don't push more often than throttleMs
	if now.Sub(prev.pushedAt) < time.Duration(t.throttleMs)*time.Millisecond {
		return false
	}

	// Keep-alive: push if too long since last push
	if t.maxSilenceSeconds > 0 && now.Sub(prev.pushedAt) >= time.Duration(t.maxSilenceSeconds)*time.Second {
		t.last[driverID] = &lastState{lat, lng, heading, updatedAt, now}
		return true
	}

	// Significant distance?
	curr := location.GeoPoint{Lat: lat, Lng: lng}
	prevPt := location.GeoPoint{Lat: prev.lat, Lng: prev.lng}
	dist := location.HaversineDistanceMeters(prevPt, curr)
	if dist >= t.minDistanceMeters {
		t.last[driverID] = &lastState{lat, lng, heading, updatedAt, now}
		return true
	}

	// Significant heading change?
	diff := location.AngleDiffDegrees(prev.heading, heading)
	if diff < 0 {
		diff = -diff
	}
	if diff >= t.minHeadingDegrees {
		t.last[driverID] = &lastState{lat, lng, heading, updatedAt, now}
		return true
	}

	return false
}

func (t *Threshold) maybePurge() {
	if time.Since(t.lastPurge) < 5*time.Minute {
		return
	}
	t.lastPurge = time.Now()
	// Optional: remove entries older than e.g. 1 hour to avoid unbounded map growth
	// For now we keep all; in production purge stale driver IDs.
}
