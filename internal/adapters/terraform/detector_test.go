package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".terraform-version")

	if err := os.WriteFile(path, []byte("1.8.5\n"), 0644); err != nil {
		t.Fatal(err)
	}

	dependencies, err := New().Detect(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(dependencies) != 1 {
		t.Fatalf("se esperaba 1 dependencia, se encontraron %d", len(dependencies))
	}

	if dependencies[0].Name != "terraform" {
		t.Fatalf("nombre inesperado: %s", dependencies[0].Name)
	}

	if dependencies[0].CurrentVersion != "1.8.5" {
		t.Fatalf("versión inesperada: %s", dependencies[0].CurrentVersion)
	}
}
