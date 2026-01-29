package trip

import "context"

// Service holds trip business logic. No transport or infra here.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, e *Entity) error {
	return s.repo.Create(ctx, e)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Entity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) UpdateCurrentStop(ctx context.Context, id string, currentStopIndex int) error {
	return s.repo.UpdateCurrentStop(ctx, id, currentStopIndex)
}

func (s *Service) UpdateETA(ctx context.Context, id string, etaSeconds int) error {
	return s.repo.UpdateETA(ctx, id, etaSeconds)
}
