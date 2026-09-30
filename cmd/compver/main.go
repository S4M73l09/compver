package main

import (
	"fmt"
	"os"

	gomod "github.com/S4M73l09/compver/internal/adapters/gomod"
	"github.com/S4M73l09/compver/internal/engine"
)

func main() {
	root := "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	analyzer := engine.New(
		gomod.New(),
	)
	result, err := analyzer.Analyze(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ruta analizada: %s\n", result.Path)
	fmt.Printf("Dependencias encontradas: %d\n", len(result.Dependencies))

	if len(result.Dependencies) == 0 {
		return
	}

	fmt.Println()
	fmt.Printf("%-45s %-15s %-12s %s\n", "Nombre", "Versión", "Tipo", "Estado")
	fmt.Printf("%-45s %-15s %-12s %s\n", "------", "-------", "----", "------")

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
}
