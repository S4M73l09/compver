package engine

import (
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

func TestMergeDependencies(t *testing.T) {
	constraintVersion, err := version.Parse("5.0.0")
	if err != nil {
		t.Fatal(err)
	}
	constraint := version.Constraint{
		Alternatives: [][]version.Condition{{{
			Operator: version.GreaterEqual,
			Version:  constraintVersion,
		}}},
	}

	dependencies := []model.Dependency{
		{
			Tool:           "terraform",
			Name:           "registry.terraform.io/hashicorp/aws",
			CurrentVersion: "~> 5.0",
			Constraint:     constraint,
			Observations: []model.VersionObservation{{
				Role:    "constraint",
				Version: "~> 5.0",
			}},
		},
		{
			Tool:           "terraform",
			Name:           "registry.terraform.io/hashicorp/aws",
			CurrentVersion: "5.31.0",
			Observations: []model.VersionObservation{{
				Role:    "locked",
				Version: "5.31.0",
			}},
		},
	}

	merged := MergeDependencies(dependencies)
	if len(merged) != 1 {
		t.Fatalf("se esperaba una dependencia unificada, se obtuvieron %d", len(merged))
	}

	if len(merged[0].Observations) != 2 {
		t.Fatalf("se esperaban dos observaciones, se obtuvieron %d", len(merged[0].Observations))
	}

	if merged[0].Constraint.Empty() {
		t.Fatal("se perdió la restricción durante la correlación")
	}
}
