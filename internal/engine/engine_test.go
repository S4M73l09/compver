package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

type fakeDetector struct{}

func (fakeDetector) CanHandle(path string) bool {
	return filepath.Base(path) == "example.txt"
}

func (fakeDetector) Detect(path string) ([]model.Dependency, error) {
	return []model.Dependency{
		{
			Name:           "example",
			CurrentVersion: "1.0.0",
			Source:         path,
		},
	}, nil
}

func TestAnalyzeUsesDetectors(t *testing.T) {
	root := t.TempDir()

	err := os.WriteFile(
		filepath.Join(root, "example.txt"),
		[]byte("example"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	engine := New(fakeDetector{})

	result, err := engine.Analyze(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Dependencies) != 1 {
		t.Fatalf(
			"se esperaba 1 dependencia, se encontraron %d",
			len(result.Dependencies),
		)
	}

	if result.Dependencies[0].VersionKind != version.Stable {
		t.Fatalf(
			"se esperaba una versión estable, se obtuvo %s",
			result.Dependencies[0].VersionKind,
		)
	}
}
