package driver

import "context"

// Repository abstracts persistence for drivers (PostgreSQL).
// Live location and H3 index are in location.Service / Redis, not here.
type Repository interface {
	GetByID(ctx context.Context, id string) (*Entity, error)
	GetByUserID(ctx context.Context, userID string) (*Entity, error)
}
