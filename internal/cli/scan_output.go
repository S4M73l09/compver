package cli

import (
	"fmt"
	"strings"
	"time"

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

	compatibleVersions := version.FilterByConstraint(
		dependency.AvailableVersions,
		dependency.Constraint,
	)

	selectedVersions := version.Select(
		compatibleVersions,
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
		comparisonVersion(dependency),
		selectedVersions,
	)
	status := string(comparison.Status)
	latest := comparison.Latest
	if dependency.ProviderError != "" {
		availableVersions = "error"
		selectedVersionCount = "-"
	}
	if dependency.Evaluation != nil {
		status = string(dependency.Evaluation.Status)
		if dependency.Evaluation.LatestPublished != nil {
			latest = &dependency.Evaluation.LatestPublished.Version
		}
	}

	fmt.Printf(
		"%-45s %-15s %-12s %-12s %-12s %s\n",
		dependency.Name,
		displayVersion(dependency),
		dependencyType,
		dependency.VersionKind,
		availableVersions,
		selectedVersionCount,
	)

	if len(selectedVersions) > 0 {
		if options.Show[showReleases] {
			printSelectedReleases(dependency, selectedVersions)
		}
	}

	if options.Show[showSources] {
		printSources(dependency)
	}
	if options.Show[showLifecycle] {
		printLifecycle(dependency)
	}
	if options.Show[showVulnerabilities] {
		printVulnerabilities(dependency)
	}

	fmt.Printf("  Estado: %s\n", status)
	if latest != nil {
		fmt.Printf("  Última publicada: %s\n", latest.String())
	}
	if dependency.Evaluation != nil &&
		dependency.Evaluation.LatestCompatible != nil {
		fmt.Printf(
			"  Última compatible: %s\n",
			dependency.Evaluation.LatestCompatible.Version.String(),
		)
	}
}

func printSelectedReleases(
	dependency model.Dependency,
	selectedVersions []version.Version,
) {
	fmt.Println(" Versiones seleccionadas:")

	for _, selectedVersion := range selectedVersions {
		release := findRelease(dependency.AvailableReleases, selectedVersion)
		if release == nil {
			fmt.Printf("  - %s\n", selectedVersion.String())
			continue
		}

		metadata := []string{}
		if release.Cycle != "" {
			metadata = append(metadata, "ciclo "+release.Cycle)
		}
		if release.Support != model.SupportUnknown {
			metadata = append(metadata, "soporte "+string(release.Support))
		}
		if len(release.Vulnerabilities) > 0 {
			metadata = append(metadata,
				fmt.Sprintf("%d vulnerabilidades", len(release.Vulnerabilities)),
			)
		}

		if len(metadata) == 0 {
			fmt.Printf("  - %s\n", selectedVersion.String())
		} else {
			fmt.Printf("  - %s (%s)\n", selectedVersion.String(), strings.Join(metadata, ", "))
		}
	}
}

func printSources(dependency model.Dependency) {
	fmt.Printf("  Fuente del detector: %s\n", displaySource(dependency.Source))
	printVersionObservations(dependency)
	fmt.Printf(
		"  Fuente de versiones: %s%s\n",
		displaySource(dependency.ProviderSource),
		cacheSuffix(dependency.ProviderSource, dependency.ProviderFromCache),
	)
	fmt.Printf(
		"  Fuente de lifecycle: %s%s\n",
		displaySource(dependency.LifecycleSource),
		cacheSuffix(dependency.LifecycleSource, dependency.LifecycleFromCache),
	)
	fmt.Printf(
		"  Fuente de seguridad: %s%s\n",
		displaySource(dependency.VulnerabilitySource),
		cacheSuffix(dependency.VulnerabilitySource, dependency.VulnerabilityFromCache),
	)
}

func printLifecycle(dependency model.Dependency) {
	if dependency.LifecycleError != "" {
		fmt.Printf("  Lifecycle: error (%s)\n", dependency.LifecycleError)
		return
	}

	status := model.SupportUnknown
	if dependency.Evaluation != nil {
		status = dependency.Evaluation.Support
	}

	fmt.Printf("  Lifecycle: %s\n", status)
	if supportEnd := currentSupportEnd(dependency); supportEnd != nil {
		fmt.Printf("  Fin de soporte: %s\n", supportEnd.Format("2006-01-02"))
	}
}

func currentSupportEnd(dependency model.Dependency) *time.Time {
	current := comparisonVersion(dependency)
	parsed, err := version.Parse(current)
	if err != nil {
		return nil
	}

	release := findRelease(dependency.AvailableReleases, parsed)
	if release == nil {
		return nil
	}
	return release.SupportEnd
}

func findRelease(releases []model.ReleaseInfo, target version.Version) *model.ReleaseInfo {
	for index := range releases {
		if version.Compare(releases[index].Version, target) == 0 {
			return &releases[index]
		}
	}
	return nil
}

func displaySource(source string) string {
	if source == "" {
		return "no disponible"
	}
	return source
}

func cacheSuffix(source string, fromCache bool) string {
	if source == "" {
		return ""
	}
	if fromCache {
		return " (caché)"
	}
	return " (red)"
}

func printVulnerabilities(dependency model.Dependency) {
	if dependency.VulnerabilityError != "" {
		fmt.Printf(
			"  Vulnerabilidades: error (%s)\n",
			dependency.VulnerabilityError,
		)
		return
	}
	if dependency.VulnerabilitySource == "" {
		fmt.Println("  Vulnerabilidades: no consultadas")
		return
	}

	if len(dependency.Vulnerabilities) == 0 {
		fmt.Println("  Vulnerabilidades: ninguna conocida")
		return
	}

	fmt.Printf(
		"  Vulnerabilidades: %d\n",
		len(dependency.Vulnerabilities),
	)
	for _, vulnerability := range dependency.Vulnerabilities {
		fmt.Printf(
			"  - %s [%s]: %s\n",
			vulnerability.ID,
			vulnerability.Severity,
			vulnerability.Summary,
		)
		if len(vulnerability.FixedIn) > 0 {
			fixedVersions := make([]string, 0, len(vulnerability.FixedIn))
			for _, fixedVersion := range vulnerability.FixedIn {
				fixedVersions = append(fixedVersions, fixedVersion.String())
			}
			fmt.Printf("    Corregida en: %s\n", strings.Join(fixedVersions, ", "))
		}
	}
}

func displayVersion(dependency model.Dependency) string {
	return dependency.CurrentVersion
}

func comparisonVersion(dependency model.Dependency) string {
	for _, role := range []string{"installed", "locked", "current"} {
		for _, observation := range dependency.Observations {
			if observation.Role == role {
				return observation.Version
			}
		}
	}

	return dependency.CurrentVersion
}

func printVersionObservations(dependency model.Dependency) {
	for _, role := range []string{"constraint", "locked", "installed"} {
		for _, observation := range dependency.Observations {
			if observation.Role != role {
				continue
			}

			fmt.Printf(
				"  %s: %s\n",
				observationLabel(role),
				observation.Version,
			)
		}
	}
}

func observationLabel(role string) string {
	switch role {
	case "constraint":
		return "Restricción"
	case "locked":
		return "Bloqueada"
	case "installed":
		return "Instalada"
	default:
		return role
	}
}
