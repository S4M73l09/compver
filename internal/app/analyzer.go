package app

import (
	"fmt"

	gomodadapter "github.com/S4M73l09/compver/internal/adapters/gomod"
	terraformadapter "github.com/S4M73l09/compver/internal/adapters/terraform"
	"github.com/S4M73l09/compver/internal/cache"
	"github.com/S4M73l09/compver/internal/engine"
)

func NewAnalyzer(tool string) (*engine.Engine, error) {
	if tool != "" && tool != "go" && tool != "terraform" {
		return nil, fmt.Errorf("herramienta no soportada: %s", tool)
	}

	var detectors []engine.Detector
	if tool == "" || tool == "go" {
		detectors = append(detectors, gomodadapter.New())
	}
	if tool == "" || tool == "terraform" {
		detectors = append(detectors, terraformadapter.New())
	}

	analyzer := engine.New(detectors...)
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
	analyzer.AddProvider(
		terraformadapter.NewProvider(
			"https://releases.hashicorp.com/terraform/",
			nil,
			providerCache,
		),
	)
	analyzer.AddProvider(
		terraformadapter.NewRegistryProvider(
			"https://registry.terraform.io",
			nil,
			providerCache,
		),
	)

	return analyzer, nil
}
