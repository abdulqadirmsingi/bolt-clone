package redis

import (
	"context"
	"sync"

	"github.com/zeng/mobility-backend/internal/domain/fireball"
)

// FireballPublisher broadcasts driver location updates to gRPC stream subscribers.
// Uses in-memory pub/sub per driver; in production this can be Redis Pub/Sub or a message bus.
type FireballPublisher struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan fireball.LocationUpdate]struct{}
}

func NewFireballPublisher(_ *Client) *FireballPublisher {
	return &FireballPublisher{
		subscribers: make(map[string]map[chan fireball.LocationUpdate]struct{}),
	}
}

func (p *FireballPublisher) PublishDriverLocation(ctx context.Context, driverID string, lat, lng, heading float64, updatedAt int64) error {
	update := fireball.LocationUpdate{
		DriverID:  driverID,
		Lat:       lat,
		Lng:       lng,
		Heading:   heading,
		UpdatedAt: updatedAt,
	}
	p.mu.RLock()
	subs := p.subscribers[driverID]
	p.mu.RUnlock()
	if subs == nil {
		return nil
	}
	for ch := range subs {
		select {
		case ch <- update:
		case <-ctx.Done():
			return ctx.Err()
		default:
			// non-blocking; skip if channel full
		}
	}
	return nil
}

func (p *FireballPublisher) SubscribeToDriver(ctx context.Context, driverID string) (<-chan fireball.LocationUpdate, func()) {
	ch := make(chan fireball.LocationUpdate, 32)
	p.mu.Lock()
	if p.subscribers[driverID] == nil {
		p.subscribers[driverID] = make(map[chan fireball.LocationUpdate]struct{})
	}
	p.subscribers[driverID][ch] = struct{}{}
	p.mu.Unlock()
	cancel := func() {
		p.mu.Lock()
		delete(p.subscribers[driverID], ch)
		if len(p.subscribers[driverID]) == 0 {
			delete(p.subscribers, driverID)
		}
		p.mu.Unlock()
		close(ch)
	}
	return ch, cancel
}

var _ fireball.LocationPublisher = (*FireballPublisher)(nil)
