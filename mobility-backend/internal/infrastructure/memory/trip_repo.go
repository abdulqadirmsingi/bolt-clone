package memory

import (
	"context"
	"sync"

	"github.com/zeng/mobility-backend/internal/domain/trip"
)

// TripRepo is an in-memory implementation of trip.Repository for development.
// Replace with postgres implementation for production.
type TripRepo struct {
	mu    sync.RWMutex
	trips map[string]*trip.Entity
}

func NewTripRepo() *TripRepo {
	return &TripRepo{trips: make(map[string]*trip.Entity)}
}

func (r *TripRepo) Create(ctx context.Context, e *trip.Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trips[e.ID] = e
	return nil
}

func (r *TripRepo) GetByID(ctx context.Context, id string) (*trip.Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.trips[id]
	if !ok {
		return nil, nil
	}
	return e, nil
}

func (r *TripRepo) UpdateStatus(ctx context.Context, id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.trips[id]; ok {
		e.Status = status
	}
	return nil
}

func (r *TripRepo) UpdateCurrentStop(ctx context.Context, id string, currentStopIndex int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.trips[id]; ok {
		e.CurrentStopIndex = currentStopIndex
	}
	return nil
}

func (r *TripRepo) UpdateETA(ctx context.Context, id string, etaSeconds int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.trips[id]; ok {
		e.ETASeconds = etaSeconds
	}
	return nil
}

var _ trip.Repository = (*TripRepo)(nil)
