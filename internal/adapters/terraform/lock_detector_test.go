package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLockDetector(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".terraform.lock.hcl")
	content := `provider "registry.terraform.io/hashicorp/aws" {
  version     = "5.42.0"
  constraints = "~> 5.0"
}

provider "registry.terraform.io/hashicorp/random" {
  version = "3.6.0"
}
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	dependencies, err := NewLockDetector().Detect(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(dependencies) != 2 {
		t.Fatalf("se esperaban 2 providers, se encontraron %d", len(dependencies))
	}

	if dependencies[0].Name != "registry.terraform.io/hashicorp/aws" {
		t.Fatalf("provider inesperado: %s", dependencies[0].Name)
	}
	if dependencies[0].CurrentVersion != "5.42.0" {
		t.Fatalf("versión inesperada: %s", dependencies[0].CurrentVersion)
	}
}
