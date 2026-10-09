package cli

import (
	"fmt"
	"strings"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

func printScanResult(
	result model.AnalysisResult,
	options scanOptions,
) {
	groups := make(map[string][]model.Dependency)
	toolOrder := make([]string, 0)

	for _, dependency := range result.Dependencies {
		tool := dependency.Tool
		if tool == "" {
			tool = "general"
		}

		if _, exists := groups[tool]; !exists {
			toolOrder = append(toolOrder, tool)
		}

		groups[tool] = append(groups[tool], dependency)
	}

	for _, tool := range toolOrder {
		fmt.Println()
		fmt.Printf("=== %s ===\n", strings.ToUpper(tool))
		printTableHeader()

		for _, dependency := range groups[tool] {
			printDependency(dependency, options)
		}
	}
}

func printTableHeader() {
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
}

func printDependency(
	dependency model.Dependency,
	options scanOptions,
) {
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
