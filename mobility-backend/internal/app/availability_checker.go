package app

import (
	"context"

	"github.com/zeng/mobility-backend/internal/domain/driver"
	"github.com/zeng/mobility-backend/internal/domain/fireball"
)

// availabilityCheckerAdapter implements fireball.AvailabilityChecker using driver.AvailabilityService.
type availabilityCheckerAdapter struct {
	availSvc *driver.AvailabilityService
}

func newAvailabilityCheckerAdapter(availSvc *driver.AvailabilityService) fireball.AvailabilityChecker {
	return &availabilityCheckerAdapter{availSvc: availSvc}
}

func (a *availabilityCheckerAdapter) IsAvailableForDiscovery(ctx context.Context, driverID string) (bool, error) {
	status, err := a.availSvc.GetStatus(ctx, driverID)
	if err != nil {
		return false, err
	}
	return driver.IsAvailableForDiscovery(status), nil
}
