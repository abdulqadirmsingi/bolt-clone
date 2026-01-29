package rider

import "context"

// Service holds rider business logic. No transport or infra here.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Entity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*Entity, error) {
	return s.repo.GetByUserID(ctx, userID)
}
