package terraform

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/S4M73l09/compver/internal/model"
)

func TestParseTerraformLifecycle(t *testing.T) {
	records := []terraformLifecycleRecord{
		{
			Cycle:             "1.8",
			Latest:            "1.8.5",
			LatestReleaseDate: "2024-06-05",
			EOL:               json.RawMessage(`"2024-11-26"`),
		},
		{
			Cycle:             "1.9",
			Latest:            "1.9.8",
			LatestReleaseDate: "2024-10-16",
			EOL:               json.RawMessage(`false`),
		},
	}

	releases, err := parseTerraformLifecycle(
		records,
		time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(releases) != 2 {
		t.Fatalf("se esperaban dos ciclos, se obtuvieron %d", len(releases))
	}

	if releases[0].Support != model.SupportEOL {
		t.Fatalf("se esperaba EOL, se obtuvo %s", releases[0].Support)
	}

	if releases[0].Cycle != "1.8" {
		t.Fatalf("ciclo inesperado: %s", releases[0].Cycle)
	}

	if releases[1].Support != model.SupportActive {
		t.Fatalf("se esperaba soporte activo, se obtuvo %s", releases[1].Support)
	}
}
