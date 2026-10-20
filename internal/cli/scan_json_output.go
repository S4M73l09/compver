package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

type jsonAnalysisResult struct {
	Path         string           `json:"path"`
	Dependencies []jsonDependency `json:"dependencies"`
}

type jsonDependency struct {
	Tool                   string               `json:"tool"`
	Name                   string               `json:"name"`
	CurrentVersion         string               `json:"current_version"`
	Constraint             string               `json:"constraint,omitempty"`
	Indirect               bool                 `json:"indirect"`
	VersionKind            string               `json:"version_kind"`
	Status                 string               `json:"status,omitempty"`
	LatestPublished        string               `json:"latest_published,omitempty"`
	LatestCompatible       string               `json:"latest_compatible,omitempty"`
	Observations           []jsonObservation    `json:"observations,omitempty"`
	AvailableVersionsCount int                  `json:"available_versions_count"`
	AvailableVersions      []string             `json:"available_versions,omitempty"`
	SelectedVersions       []string             `json:"selected_versions,omitempty"`
	Releases               []jsonRelease        `json:"releases,omitempty"`
	Vulnerabilities        []jsonVulnerability  `json:"vulnerabilities,omitempty"`
	Sources                jsonSources          `json:"sources"`
	Errors                 jsonDependencyErrors `json:"errors,omitempty"`
}

type jsonObservation struct {
	Role    string `json:"role"`
	Version string `json:"version"`
	Source  string `json:"source"`
}

type jsonRelease struct {
	Version         string              `json:"version"`
	Cycle           string              `json:"cycle,omitempty"`
	PublishedAt     *time.Time          `json:"published_at,omitempty"`
	Support         string              `json:"support,omitempty"`
	SupportEnd      *time.Time          `json:"support_end,omitempty"`
	Vulnerabilities []jsonVulnerability `json:"vulnerabilities,omitempty"`
}

type jsonVulnerability struct {
	ID       string   `json:"id"`
	Summary  string   `json:"summary"`
	Severity string   `json:"severity"`
	FixedIn  []string `json:"fixed_in,omitempty"`
}

type jsonSources struct {
	Detector               string `json:"detector,omitempty"`
	Versions               string `json:"versions,omitempty"`
	Lifecycle              string `json:"lifecycle,omitempty"`
	Security               string `json:"security,omitempty"`
	ProviderFromCache      bool   `json:"provider_from_cache,omitempty"`
	LifecycleFromCache     bool   `json:"lifecycle_from_cache,omitempty"`
	VulnerabilityFromCache bool   `json:"vulnerability_from_cache,omitempty"`
}

type jsonDependencyErrors struct {
	Provider      string `json:"provider,omitempty"`
	Lifecycle     string `json:"lifecycle,omitempty"`
	Vulnerability string `json:"vulnerability,omitempty"`
}

func printScanJSON(
	writer io.Writer,
	result model.AnalysisResult,
	options scanOptions,
) error {
	output := buildJSONAnalysisResult(result, options)

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	return encoder.Encode(output)
}

func buildJSONAnalysisResult(
	result model.AnalysisResult,
	options scanOptions,
) jsonAnalysisResult {
	output := jsonAnalysisResult{
		Path:         result.Path,
		Dependencies: make([]jsonDependency, 0, len(result.Dependencies)),
	}

	for _, dependency := range result.Dependencies {
		output.Dependencies = append(
			output.Dependencies,
			buildJSONDependency(dependency, options),
		)
	}

	return output
}

func buildJSONDependency(
	dependency model.Dependency,
	options scanOptions,
) jsonDependency {
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

	result := jsonDependency{
		Tool:                   dependency.Tool,
		Name:                   dependency.Name,
		CurrentVersion:         dependency.CurrentVersion,
		Constraint:             constraintString(dependency.Constraint),
		Indirect:               dependency.Indirect,
		VersionKind:            string(dependency.VersionKind),
		Observations:           jsonObservations(dependency.Observations),
		AvailableVersionsCount: len(dependency.AvailableVersions),
		SelectedVersions:       versionStrings(selectedVersions),
		Releases:               selectedJSONReleases(dependency.AvailableReleases, selectedVersions),
		Vulnerabilities:        jsonVulnerabilities(dependency.Vulnerabilities),
		Sources: jsonSources{
			Detector:               dependency.Source,
			Versions:               dependency.ProviderSource,
			Lifecycle:              dependency.LifecycleSource,
			Security:               dependency.VulnerabilitySource,
			ProviderFromCache:      dependency.ProviderFromCache,
			LifecycleFromCache:     dependency.LifecycleFromCache,
			VulnerabilityFromCache: dependency.VulnerabilityFromCache,
		},
		Errors: jsonDependencyErrors{
			Provider:      dependency.ProviderError,
			Lifecycle:     dependency.LifecycleError,
			Vulnerability: dependency.VulnerabilityError,
		},
	}
	if options.AllVersions {
		result.AvailableVersions = versionStrings(dependency.AvailableVersions)
	}

	if dependency.Evaluation != nil {
		result.Status = string(dependency.Evaluation.Status)
		if dependency.Evaluation.LatestPublished != nil {
			result.LatestPublished = dependency.Evaluation.LatestPublished.Version.String()
		}
		if dependency.Evaluation.LatestCompatible != nil {
			result.LatestCompatible = dependency.Evaluation.LatestCompatible.Version.String()
		}
	}

	return result
}

func jsonObservations(
	observations []model.VersionObservation,
) []jsonObservation {
	result := make([]jsonObservation, 0, len(observations))
	for _, observation := range observations {
		result = append(result, jsonObservation{
			Role:    observation.Role,
			Version: observation.Version,
			Source:  observation.Source,
		})
	}
	return result
}

func versionStrings(versions []version.Version) []string {
	result := make([]string, 0, len(versions))
	for _, releaseVersion := range versions {
		result = append(result, releaseVersion.String())
	}
	return result
}

func selectedJSONReleases(
	availableReleases []model.ReleaseInfo,
	selectedVersions []version.Version,
) []jsonRelease {
	result := make([]jsonRelease, 0, len(selectedVersions))
	for _, selectedVersion := range selectedVersions {
		for _, release := range availableReleases {
			if version.Compare(release.Version, selectedVersion) != 0 {
				continue
			}

			result = append(result, jsonRelease{
				Version:         release.Version.String(),
				Cycle:           release.Cycle,
				PublishedAt:     release.PublishedAt,
				Support:         string(release.Support),
				SupportEnd:      release.SupportEnd,
				Vulnerabilities: jsonVulnerabilities(release.Vulnerabilities),
			})
			break
		}
	}
	return result
}

func jsonVulnerabilities(
	vulnerabilities []model.Vulnerability,
) []jsonVulnerability {
	result := make([]jsonVulnerability, 0, len(vulnerabilities))
	for _, vulnerability := range vulnerabilities {
		result = append(result, jsonVulnerability{
			ID:       vulnerability.ID,
			Summary:  vulnerability.Summary,
			Severity: vulnerability.Severity,
			FixedIn:  versionStrings(vulnerability.FixedIn),
		})
	}
	return result
}

func constraintString(constraint version.Constraint) string {
	if constraint.Empty() {
		return ""
	}

	alternatives := make([]string, 0, len(constraint.Alternatives))
	for _, conditions := range constraint.Alternatives {
		parts := make([]string, 0, len(conditions))
		for _, condition := range conditions {
			parts = append(parts, fmt.Sprintf(
				"%s %s",
				condition.Operator,
				condition.Version.String(),
			))
		}
		alternatives = append(alternatives, strings.Join(parts, ", "))
	}

	return strings.Join(alternatives, " || ")
}
