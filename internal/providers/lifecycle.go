package providers

import (
	"context"
	"time"

	"github.com/S4M73l09/compver/internal/model"
)

type LifecycleResult struct {
	Releases    []model.ReleaseInfo
	Source      string
	RetrievedAt time.Time
	FromCache   bool
}

type LifecycleProvider interface {
	CanHandle(dependency model.Dependency) bool

	Lifecycle(
		ctx context.Context,
		dependency model.Dependency,
		options QueryOptions,
	) (LifecycleResult, error)
}
