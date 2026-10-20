package cli

import "testing"

func TestParseScanOptionsAcceptsRouteBeforeFlags(t *testing.T) {
	options, root, err := parseScanOptions([]string{
		"/tmp/compver-fixture",
		"--format",
		"json",
		"--show",
		"releases,sources",
		"--limit",
		"6",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if root != "/tmp/compver-fixture" {
		t.Fatalf("ruta inesperada: %s", root)
	}
	if options.OutputFormat != "json" {
		t.Fatalf("formato inesperado: %s", options.OutputFormat)
	}
	if options.Limit != 6 {
		t.Fatalf("límite inesperado: %d", options.Limit)
	}
	if !options.Show[showReleases] || !options.Show[showSources] {
		t.Fatalf("secciones inesperadas: %#v", options.Show)
	}
}

func TestParseScanOptionsAcceptsFlagsBeforeRoute(t *testing.T) {
	options, root, err := parseScanOptions([]string{
		"--format=json",
		"--offline",
		"/tmp/compver-fixture",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if root != "/tmp/compver-fixture" {
		t.Fatalf("ruta inesperada: %s", root)
	}
	if options.OutputFormat != "json" {
		t.Fatalf("formato inesperado: %s", options.OutputFormat)
	}
}

func TestParseScanOptionsRejectsMultipleRoutes(t *testing.T) {
	_, _, err := parseScanOptions([]string{
		"/tmp/one",
		"--format",
		"json",
		"/tmp/two",
	})
	if err == nil {
		t.Fatal("se esperaba un error para múltiples rutas")
	}
}

func TestParseScanOptionsVerboseEnablesAllSections(t *testing.T) {
	options, _, err := parseScanOptions([]string{"--verbose"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	for _, section := range []string{
		showLifecycle,
		showVulnerabilities,
		showSources,
		showReleases,
	} {
		if !options.Show[section] {
			t.Fatalf("sección no activada por verbose: %s", section)
		}
	}
}

func TestParseScanOptionsRejectsConflictingPresentationModes(t *testing.T) {
	_, _, err := parseScanOptions([]string{"--summary", "--verbose"})
	if err == nil {
		t.Fatal("se esperaba un error para modos de presentación incompatibles")
	}
}
