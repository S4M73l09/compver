package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/S4M73l09/compver/internal/version"
)

func TestConfigDetector(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.tf")
	content := `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    random = {
      version = ">= 3.0, < 4.0"
    }
  }
}
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	dependencies, err := NewConfigDetector().Detect(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(dependencies) != 2 {
		t.Fatalf("se esperaban 2 providers, se encontraron %d", len(dependencies))
	}

	if dependencies[0].Name != "registry.terraform.io/hashicorp/aws" {
		t.Fatalf("provider inesperado: %s", dependencies[0].Name)
	}
	if dependencies[0].Constraint.Empty() {
		t.Fatal("la restricción no debería estar vacía")
	}
	if !dependencies[0].Constraint.Matches(mustParse("5.9.0")) {
		t.Fatal("5.9.0 debería cumplir ~> 5.0")
	}
	if dependencies[0].Constraint.Matches(mustParse("6.0.0")) {
		t.Fatal("6.0.0 no debería cumplir ~> 5.0")
	}
	if dependencies[1].Name != "registry.terraform.io/hashicorp/random" {
		t.Fatalf("provider inesperado: %s", dependencies[1].Name)
	}
}

func mustParse(value string) version.Version {
	parsed, err := version.Parse(value)
	if err != nil {
		panic(err)
	}

	return parsed
}
