package app

import (
	"fmt"

	gomodadapter "github.com/S4M73l09/compver/internal/adapters/gomod"
	"github.com/S4M73l09/compver/internal/engine"
)

func NewAnalyzer(tool string) (*engine.Engine, error) {
	if tool != "" && tool != "go" {
		return nil, fmt.Errorf("herramienta no soportada: %s", tool)
	}

	analyzer := engine.New(
		gomodadapter.New(),
	)

	analyzer.AddProvider(
		gomodadapter.NewProvider(
			"https://proxy.golang.org",
			nil,
		),
	)

	return analyzer, nil
}
