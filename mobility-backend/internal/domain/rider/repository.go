package rider

import "context"

// Repository abstracts persistence for riders (PostgreSQL).
type Repository interface {
	GetByID(ctx context.Context, id string) (*Entity, error)
	GetByUserID(ctx context.Context, userID string) (*Entity, error)
}
