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
	style terminalStyle,
) {
	groups, toolOrder := groupDependencies(result.Dependencies)

	for _, tool := range toolOrder {
		rows := make([]scanTableRow, 0, len(groups[tool]))
		for _, dependency := range groups[tool] {
			rows = append(rows, buildScanTableRow(dependency, options))
		}

		fmt.Println()
		fmt.Printf("=== %s ===\n", strings.ToUpper(tool))
		printToolSummary(groups[tool], style)

		widths := calculateTableWidths(rows)
		printTableHeader(widths)

		for _, row := range rows {
			printDependency(row, options, style, widths)
		}
	}
}

func printToolSummaries(
	result model.AnalysisResult,
	style terminalStyle,
) {
	groups, toolOrder := groupDependencies(result.Dependencies)

	for _, tool := range toolOrder {
		fmt.Println()
		fmt.Printf("=== %s ===\n", strings.ToUpper(tool))
		printToolSummary(groups[tool], style)
	}
}

func groupDependencies(
	dependencies []model.Dependency,
) (map[string][]model.Dependency, []string) {
	groups := make(map[string][]model.Dependency)
	toolOrder := make([]string, 0)

	for _, dependency := range dependencies {
		tool := dependency.Tool
		if tool == "" {
			tool = "general"
		}

		if _, exists := groups[tool]; !exists {
			toolOrder = append(toolOrder, tool)
		}

		groups[tool] = append(groups[tool], dependency)
	}

	return groups, toolOrder
}

func printToolSummary(
	dependencies []model.Dependency,
	style terminalStyle,
) {
	updates := 0
	vulnerabilities := 0
	endOfLife := 0

	for _, dependency := range dependencies {
		switch dependencyStatus(dependency) {
		case string(model.StatusUpdateAvailable),
			string(model.StatusCompatibleUpdate),
			string(model.StatusVulnerableUpdate):
			updates++
		case string(model.StatusEndOfLife):
			endOfLife++
		}

		vulnerabilities += len(dependency.Vulnerabilities)
	}

	fmt.Printf("Dependencias: %d\n", len(dependencies))
	fmt.Printf("Actualizaciones: %s\n", style.count(updates))
	fmt.Printf("Vulnerabilidades: %s\n", style.count(vulnerabilities))
	fmt.Printf("EOL: %s\n", style.count(endOfLife))
	fmt.Println()
}

func printScanSummary(
	result model.AnalysisResult,
	style terminalStyle,
) {
	tools := make(map[string]bool)
	updates := 0
	vulnerabilities := 0
	endOfLife := 0

	for _, dependency := range result.Dependencies {
		tool := dependency.Tool
		if tool == "" {
			tool = "general"
		}
		tools[tool] = true

		switch dependencyStatus(dependency) {
		case string(model.StatusUpdateAvailable),
			string(model.StatusCompatibleUpdate),
			string(model.StatusVulnerableUpdate):
			updates++
		case string(model.StatusEndOfLife):
			endOfLife++
		}

		vulnerabilities += len(dependency.Vulnerabilities)
	}

	fmt.Println(style.heading("Compver scan"))
	fmt.Println("────────────")
	fmt.Printf("Ruta: %s\n", result.Path)
	fmt.Printf("Herramientas: %d\n", len(tools))
	fmt.Printf("Dependencias: %d\n", len(result.Dependencies))
	fmt.Printf("Actualizaciones: %s\n", style.count(updates))
	fmt.Printf("Vulnerabilidades: %s\n", style.count(vulnerabilities))
	fmt.Printf("EOL: %s\n", style.count(endOfLife))
}

type scanTableRow struct {
	dependency      model.Dependency
	dependencyType  string
	version         string
	status          string
	available       string
	selected        string
	selectedVersion []version.Version
	latest          *version.Version
}

type tableWidths struct {
	name      int
	version   int
	typeName  int
	status    int
	available int
	selected  int
}

func buildScanTableRow(
	dependency model.Dependency,
	options scanOptions,
) scanTableRow {
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

	return scanTableRow{
		dependency:      dependency,
		dependencyType:  dependencyType,
		version:         displayVersion(dependency),
		status:          status,
		available:       availableVersions,
		selected:        selectedVersionCount,
		selectedVersion: selectedVersions,
		latest:          latest,
	}
}

func calculateTableWidths(rows []scanTableRow) tableWidths {
	return tableWidths{
		name:      columnWidth("Nombre", rowValues(rows, func(row scanTableRow) string { return row.dependency.Name }), 45),
		version:   columnWidth("Version", rowValues(rows, func(row scanTableRow) string { return row.version }), 20),
		typeName:  columnWidth("Tipo", rowValues(rows, func(row scanTableRow) string { return row.dependencyType }), 12),
		status:    columnWidth("Estado", rowValues(rows, func(row scanTableRow) string { return row.status }), 30),
		available: columnWidth("Disponibles", rowValues(rows, func(row scanTableRow) string { return row.available }), 12),
		selected:  columnWidth("Mostradas", rowValues(rows, func(row scanTableRow) string { return row.selected }), 9),
	}
}

func rowValues(
	rows []scanTableRow,
	value func(scanTableRow) string,
) []string {
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		values = append(values, value(row))
	}
	return values
}

func columnWidth(header string, values []string, maximum int) int {
	width := len([]rune(header))
	for _, value := range values {
		if length := len([]rune(value)); length > width {
			width = length
		}
	}
	if width > maximum {
		return maximum
	}
	return width
}

func truncate(value string, maximum int) string {
	if maximum <= 0 {
		return ""
	}

	characters := []rune(value)
	if len(characters) <= maximum {
		return value
	}
	if maximum == 1 {
		return "…"
	}

	return string(characters[:maximum-1]) + "…"
}

func printTableHeader(widths tableWidths) {
	fmt.Printf(
		"%-*s %-*s %-*s %-*s %-*s %-*s\n",
		widths.name, "Nombre",
		widths.version, "Version",
		widths.typeName, "Tipo",
		widths.status, "Estado",
		widths.available, "Disponibles",
		widths.selected, "Mostradas",
	)

	fmt.Printf(
		"%-*s %-*s %-*s %-*s %-*s %-*s\n",
		widths.name, strings.Repeat("-", widths.name),
		widths.version, strings.Repeat("-", widths.version),
		widths.typeName, strings.Repeat("-", widths.typeName),
		widths.status, strings.Repeat("-", widths.status),
		widths.available, strings.Repeat("-", widths.available),
		widths.selected, strings.Repeat("-", widths.selected),
	)
}

func printDependency(
	row scanTableRow,
	options scanOptions,
	style terminalStyle,
	widths tableWidths,
) {
	dependency := row.dependency

	fmt.Printf(
		"%-*s %-*s %-*s %-*s %-*s %s\n",
		widths.name, truncate(row.dependency.Name, widths.name),
		widths.version, truncate(row.version, widths.version),
		widths.typeName, truncate(row.dependencyType, widths.typeName),
		widths.status, truncate(row.status, widths.status),
		widths.available, truncate(row.available, widths.available),
		truncate(row.selected, widths.selected),
	)

	if len(row.selectedVersion) > 0 {
		if options.Show[showReleases] {
			printSelectedReleases(dependency, row.selectedVersion)
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

	fmt.Printf("  Estado: %s\n", style.status(row.status))
	if row.latest != nil {
		fmt.Printf("  Última publicada: %s\n", row.latest.String())
	}
	if dependency.Evaluation != nil &&
		dependency.Evaluation.LatestCompatible != nil {
		fmt.Printf(
			"  Última compatible: %s\n",
			dependency.Evaluation.LatestCompatible.Version.String(),
		)
	}
}

func dependencyStatus(dependency model.Dependency) string {
	if dependency.Evaluation != nil {
		return string(dependency.Evaluation.Status)
	}

	selectedVersions := version.Select(
		version.FilterByConstraint(
			dependency.AvailableVersions,
			dependency.Constraint,
		),
		version.SelectionOptions{Limit: 3},
	)

	comparison := version.CompareCurrent(
		comparisonVersion(dependency),
		selectedVersions,
	)
	return string(comparison.Status)
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
