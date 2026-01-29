package trip

import "context"

// Repository abstracts persistence for trips (PostgreSQL).
type Repository interface {
	Create(ctx context.Context, e *Entity) error
	GetByID(ctx context.Context, id string) (*Entity, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateCurrentStop(ctx context.Context, id string, currentStopIndex int) error
	UpdateETA(ctx context.Context, id string, etaSeconds int) error
}
