package cli

import (
	"fmt"
	"os"

	gomod "github.com/S4M73l09/compver/internal/adapters/gomod"
	"github.com/S4M73l09/compver/internal/engine"
)

func (a *App) runScan(args []string) int {
	root := "."

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "scan acepta como máximo una ruta")
		return 1
	}

	if len(args) == 1 {
		root = args[0]
	}

	analyzer := engine.New(
		gomod.New(),
	)

	result, err := analyzer.Analyze(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Ruta analizada: %s\n", result.Path)
	fmt.Printf("Dependencias encontradas: %d\n", len(result.Dependencies))

	if len(result.Dependencies) == 0 {
		return 0
	}

	fmt.Println()
	fmt.Printf(
		"%-45s %-15s %-12s %s\n",
		"Nombre",
		"Version",
		"Tipo",
		"Estado",
	)

	fmt.Printf(
		"%-45s %-15s %-12s %s\n",
		"------",
		"-------",
		"----",
		"------",
	)

	for _, dependency := range result.Dependencies {
		dependencyType := "directa"

		if dependency.Indirect {
			dependencyType = "indirecta"
		}

		fmt.Printf(
			"%-45s %-15s %-12s %s\n",
			dependency.Name,
			dependency.CurrentVersion,
			dependencyType,
			dependency.VersionKind,
		)
	}

	return 0
}
