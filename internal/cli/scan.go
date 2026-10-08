package cli

import (
	"fmt"
	"os"

	"github.com/S4M73l09/compver/internal/app"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
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

	result, err := analyzer.AnalyzeWithOptions(
		root,
		providers.QueryOptions{
			Mode: options.NetworkMode,
		},
	)
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
		"%-45s %-15s %-12s %-12s %-12s %s\n",
		"Nombre",
		"Version",
		"Tipo",
		"Estado",
		"Disponibles",
		"Mostradas",
	)

	fmt.Printf(
		"%-45s %-15s %-12s %-12s %-12s %s\n",
		"------",
		"-------",
		"----",
		"------",
		"-----------",
		"---------",
	)

	for _, dependency := range result.Dependencies {
		dependencyType := "directa"

		if dependency.Indirect {
			dependencyType = "indirecta"
		}

		selectedVersions := version.Select(
			dependency.AvailableVersions,
			version.SelectionOptions{
				Limit:              options.Limit,
				IncludePreReleases: options.IncludePreReleases,
				AllVersions:        options.AllVersions,
			},
		)

		availableVersions := fmt.Sprintf(
			"%d",
			len(dependency.AvailableVersions),
		)
		selectedVersionCount := fmt.Sprintf(
			"%d",
			len(selectedVersions),
		)
		comparison := version.CompareCurrent(
			dependency.CurrentVersion,
			selectedVersions,
		)
		if dependency.ProviderError != "" {
			availableVersions = "error"
			selectedVersionCount = "-"
		}

		fmt.Printf(
			"%-45s %-15s %-12s %-12s %-12s %s\n",
			dependency.Name,
			dependency.CurrentVersion,
			dependencyType,
			dependency.VersionKind,
			availableVersions,
			selectedVersionCount,
		)

		if len(selectedVersions) > 0 {
			fmt.Println(" Versiones seleccionadas:")

			for _, selectedVersion := range selectedVersions {
				fmt.Printf("  - %s\n", selectedVersion.String())
			}
		}

		fmt.Printf("  Estado: %s\n", comparison.Status)
		if comparison.Latest != nil {
			fmt.Printf("  Última disponible: %s\n", comparison.Latest.String())
		}
	}

	return 0
}
