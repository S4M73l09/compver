package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "go.mod")

	content := `module example.com/project

go 1.22

require (
	github.com/example/library v1.2.3
	golang.org/x/text v0.14.0 // indirect
)
`

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	detector := New()

	if !detector.CanHandle(path) {
		t.Fatal("el detector debería aceptar go.mod")
	}

	dependencies, err := detector.Detect(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(dependencies) != 2 {
		t.Fatalf(
			"se esperaban 2 dependencias, se encontraron %d",
			len(dependencies),
		)
	}

	if dependencies[0].Name != "github.com/example/library" {
		t.Fatalf(
			"dependencia inesperada: %s",
			dependencies[0].Name,
		)
	}

	if dependencies[0].Indirect {
		t.Fatal("la primera dependencia no debería ser indirecta")
	}

	if !dependencies[1].Indirect {
		t.Fatal("la segunda dependencia debería ser indirecta")
	}
}
