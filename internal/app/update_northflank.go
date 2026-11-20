package app

import (
	"context"
	"sync/atomic"

	"github.com/Mag1cFall/AIStudio2API/internal/api"
)

type northflankUpdateService struct {
	version string
}

type northflankUpdateConflictError struct{}

func (northflankUpdateConflictError) Error() string {
	return "online updates are disabled in the Northflank container build; deploy a newer image"
}

func (northflankUpdateConflictError) HTTPStatus() int { return 409 }

func newSelfUpdater(version, _ string, _ func() string, _ *requestRegistry, _ func()) api.UpdateService {
	return northflankUpdateService{version: version}
}

func (service northflankUpdateService) UpdateStatus(ctx context.Context, _ bool) (api.UpdateStatus, error) {
	if err := ctx.Err(); err != nil {
		return api.UpdateStatus{}, err
	}
	return api.UpdateStatus{Current: service.version, Mode: "container"}, nil
}

func (northflankUpdateService) ApplyUpdate(ctx context.Context) (api.UpdateStatus, error) {
	if err := ctx.Err(); err != nil {
		return api.UpdateStatus{}, err
	}
	return api.UpdateStatus{}, northflankUpdateConflictError{}
}

func updateRestart(_ func(), _ *atomic.Bool, _ func(level, message string, attrs ...any)) func() {
	return func() {}
}

func removeReplacedExecutables(string) error { return nil }
