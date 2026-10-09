package cli

import (
	"fmt"
	"os"

	"github.com/S4M73l09/compver/internal/app"
	"github.com/S4M73l09/compver/internal/providers"
)

func (a *App) runScan(args []string) int {
	options, root, err := parseScanOptions(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	analyzer, err := app.NewAnalyzer(options.Tool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	spinner := NewSpinner(
		os.Stderr,
		"Analizando dependencias y consultando versiones...",
		isTerminal(os.Stderr),
	)

	spinner.Start()

	result, err := analyzer.AnalyzeWithOptions(
		root,
		providers.QueryOptions{
			Mode: options.NetworkMode,
		},
	)

	spinner.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Ruta analizada: %s\n", result.Path)
	fmt.Printf("Dependencias encontradas: %d\n", len(result.Dependencies))

	if len(result.Dependencies) == 0 {
		return 0
	}

	printScanResult(result, options)

	return 0
}
