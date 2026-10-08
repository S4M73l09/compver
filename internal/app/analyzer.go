package app

import (
	"fmt"

	gomodadapter "github.com/S4M73l09/compver/internal/adapters/gomod"
	"github.com/S4M73l09/compver/internal/cache"
	"github.com/S4M73l09/compver/internal/engine"
)

func NewAnalyzer(tool string) (*engine.Engine, error) {
	if tool != "" && tool != "go" {
		return nil, fmt.Errorf("herramienta no soportada: %s", tool)
	}

	analyzer := engine.New(
		gomodadapter.New(),
	)
	providerCache, err := cache.DefaultFileCache()
	if err != nil {
		return nil, fmt.Errorf("creando caché: %w", err)
	}

	analyzer.AddProvider(
		gomodadapter.NewProvider(
			"https://proxy.golang.org",
			nil,
			providerCache,
		),
	)

	return analyzer, nil
}
