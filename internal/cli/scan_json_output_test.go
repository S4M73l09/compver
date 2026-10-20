package cli

import (
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

func TestBuildJSONDependencyKeepsDefaultOutputCompact(t *testing.T) {
	availableVersions := mustParseVersions(t, []string{
		"v1.3.0",
		"v1.2.0",
		"v1.1.0",
	})

	dependency := model.Dependency{
		Tool:              "go",
		Name:              "example.com/library",
		CurrentVersion:    "v1.1.0",
		AvailableVersions: availableVersions,
	}

	result := buildJSONDependency(dependency, scanOptions{Limit: 2})

	if result.AvailableVersionsCount != 3 {
		t.Fatalf("contador inesperado: %d", result.AvailableVersionsCount)
	}
	if len(result.AvailableVersions) != 0 {
		t.Fatalf("no se esperaban todas las versiones: %#v", result.AvailableVersions)
	}
	if len(result.SelectedVersions) != 2 {
		t.Fatalf("versiones seleccionadas inesperadas: %#v", result.SelectedVersions)
	}
}

func TestBuildJSONDependencyIncludesAllVersionsWhenRequested(t *testing.T) {
	availableVersions := mustParseVersions(t, []string{
		"v1.3.0",
		"v1.2.0",
		"v1.1.0",
	})

	dependency := model.Dependency{
		Tool:              "go",
		Name:              "example.com/library",
		CurrentVersion:    "v1.1.0",
		AvailableVersions: availableVersions,
	}

	result := buildJSONDependency(dependency, scanOptions{
		Limit:       2,
		AllVersions: true,
	})

	if len(result.AvailableVersions) != 3 {
		t.Fatalf("se esperaban todas las versiones: %#v", result.AvailableVersions)
	}
}

func mustParseVersions(t *testing.T, values []string) []version.Version {
	t.Helper()

	versions := make([]version.Version, 0, len(values))
	for _, value := range values {
		parsed, err := version.Parse(value)
		if err != nil {
			t.Fatalf("parseando versión %s: %v", value, err)
		}
		versions = append(versions, parsed)
	}
	return versions
}
