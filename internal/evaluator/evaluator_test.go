package evaluator

import (
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

func TestEvaluateCompatibleUpdate(t *testing.T) {
	current, _ := version.Parse("1.0.0")
	latest, _ := version.Parse("2.0.0")

	result := Evaluate(model.Dependency{
		CurrentVersion: "1.0.0",
		AvailableReleases: []model.ReleaseInfo{
			{Version: current, Support: model.SupportActive},
			{Version: latest, Support: model.SupportActive},
		},
	})

	if result.Status != model.StatusUpdateAvailable {
		t.Fatalf("estado inesperado: %s", result.Status)
	}

	if result.LatestPublished == nil ||
		version.Compare(result.LatestPublished.Version, latest) != 0 {
		t.Fatal("no se detectó la última versión publicada")
	}
}

func TestEvaluateEndOfLife(t *testing.T) {
	current, _ := version.Parse("1.0.0")

	result := Evaluate(model.Dependency{
		CurrentVersion: "1.0.0",
		AvailableReleases: []model.ReleaseInfo{{
			Version: current,
			Support: model.SupportEOL,
		}},
	})

	if result.Status != model.StatusEndOfLife {
		t.Fatalf("estado inesperado: %s", result.Status)
	}
}
