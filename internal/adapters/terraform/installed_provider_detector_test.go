package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstalledProviderDetector(t *testing.T) {
	root := t.TempDir()

	providerPath := filepath.Join(
		root,
		".terraform",
		"providers",
		"registry.terraform.io",
		"hashicorp",
		"aws",
		"5.31.0",
		"linux_amd64",
		"terraform-provider-aws_v5.31.0_x5",
	)

	if err := os.MkdirAll(filepath.Dir(providerPath), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(providerPath, []byte("provider"), 0o644); err != nil {
		t.Fatal(err)
	}

	detector := NewInstalledProviderDetector()

	if !detector.CanHandle(providerPath) {
		t.Fatalf("el detector no reconoció %s", providerPath)
	}

	dependencies, err := detector.Detect(providerPath)
	if err != nil {
		t.Fatalf("error detectando provider: %v", err)
	}

	if len(dependencies) != 1 {
		t.Fatalf("se esperaba una dependencia, se obtuvieron %d", len(dependencies))
	}

	dependency := dependencies[0]

	if dependency.Name != "registry.terraform.io/hashicorp/aws" {
		t.Fatalf("provider incorrecto: %s", dependency.Name)
	}

	if dependency.CurrentVersion != "5.31.0" {
		t.Fatalf("versión incorrecta: %s", dependency.CurrentVersion)
	}
}
