package terraform

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

type InstalledProviderDetector struct{}

func NewInstalledProviderDetector() InstalledProviderDetector {
	return InstalledProviderDetector{}
}

func (InstalledProviderDetector) CanHandle(path string) bool {
	parts := pathParts(path)

	if len(parts) < 7 {
		return false
	}

	if !containsTerraformProvidersDirectory(parts) {
		return false
	}

	filename := parts[len(parts)-1]

	return strings.HasPrefix(filename, "terraform-provider-")
}

func (InstalledProviderDetector) Detect(
	path string,
) ([]model.Dependency, error) {
	parts := pathParts(path)

	if len(parts) < 7 {
		return nil, fmt.Errorf(
			"ruta de provider instalada no válida: %s",
			path,
		)
	}

	versionIndex := len(parts) - 3
	if versionIndex < 3 {
		return nil, fmt.Errorf(
			"no se pudo localizar la version instalada: %s",
			path,
		)
	}

	currentVersion := parts[versionIndex]

	if _, err := version.Parse(currentVersion); err != nil {
		return nil, fmt.Errorf(
			"versión instalada no válida %q: %w",
			currentVersion,
			err,
		)
	}

	providersIndex := indexOf(parts, "providers")
	if providersIndex == -1 || providersIndex+1 >= versionIndex {
		return nil, fmt.Errorf(
			"no se pudo localizar el provider: %s",
			path,
		)
	}

	providerSource := strings.Join(
		parts[providersIndex+1:versionIndex],
		"/",
	)

	return []model.Dependency{{
		Tool:           "terraform",
		Name:           providerSource,
		CurrentVersion: currentVersion,
		Source:         path,
		Observations: []model.VersionObservation{{
			Role:    "installed",
			Version: currentVersion,
			Source:  path,
		}},
	}}, nil
}

func pathParts(path string) []string {
	normalized := filepath.ToSlash(filepath.Clean(path))

	return strings.Split(normalized, "/")
}

func containsTerraformProvidersDirectory(parts []string) bool {
	for index := 0; index < len(parts)-1; index++ {
		if parts[index] == ".terraform" &&
			parts[index+1] == "providers" {
			return true
		}
	}

	return false
}

func indexOf(parts []string, value string) int {
	for index, part := range parts {
		if part == value {
			return index
		}
	}

	return -1
}
