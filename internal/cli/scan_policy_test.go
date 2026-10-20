package cli

import (
	"testing"

	"github.com/S4M73l09/compver/internal/model"
)

func TestParseFailOnAcceptsMultipleConditions(t *testing.T) {
	conditions, err := parseFailOn("vulnerable, eol")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if !conditions[failOnVulnerable] || !conditions[failOnEndOfLife] {
		t.Fatalf("criterios inesperados: %#v", conditions)
	}
}

func TestParseFailOnRejectsUnknownCondition(t *testing.T) {
	_, err := parseFailOn("unsupported")
	if err == nil {
		t.Fatal("se esperaba un error para un criterio desconocido")
	}
}

func TestScanPolicyViolatedByVulnerability(t *testing.T) {
	result := model.AnalysisResult{
		Dependencies: []model.Dependency{{
			Vulnerabilities: []model.Vulnerability{{ID: "TEST-1"}},
		}},
	}

	if !scanPolicyViolated(result, map[string]bool{failOnVulnerable: true}) {
		t.Fatal("se esperaba una violación por vulnerabilidad")
	}
}

func TestScanPolicyViolatedByEndOfLife(t *testing.T) {
	result := model.AnalysisResult{
		Dependencies: []model.Dependency{{
			Evaluation: &model.Evaluation{Status: model.StatusEndOfLife},
		}},
	}

	if !scanPolicyViolated(result, map[string]bool{failOnEndOfLife: true}) {
		t.Fatal("se esperaba una violación por fin de vida")
	}
}

func TestScanPolicyPassesWithoutMatchingCondition(t *testing.T) {
	result := model.AnalysisResult{
		Dependencies: []model.Dependency{{
			Evaluation: &model.Evaluation{Status: model.StatusCurrent},
		}},
	}

	if scanPolicyViolated(result, map[string]bool{failOnVulnerable: true}) {
		t.Fatal("no se esperaba una violación de política")
	}
}
