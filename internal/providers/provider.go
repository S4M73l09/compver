package providers

import (
	"context"
	"time"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

type NetworkMode string

const (
	NetworkAuto    NetworkMode = "auto"
	NetworkOffline NetworkMode = "offline"
	NetworkRefresh NetworkMode = "refresh"
)

type QueryOptions struct {
	Mode     NetworkMode
	CacheTTL time.Duration
}

type Result struct {
	Versions    []version.Version
	Source      string
	RetrievedAt time.Time
	FromCache   bool
}

type Provider interface {
	CanHandle(dependency model.Dependency) bool

	AvailableVersions(
		ctx context.Context,
		dependency model.Dependency,
		options QueryOptions,
	) (Result, error)
}
