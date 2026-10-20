package cli

import (
	"fmt"
	"strings"

	"github.com/S4M73l09/compver/internal/model"
)

const (
	failOnVulnerable = "vulnerable"
	failOnEndOfLife  = "end-of-life"
)

func parseFailOn(value string) (map[string]bool, error) {
	conditions := make(map[string]bool)
	if value == "" {
		return conditions, nil
	}

	for _, rawCondition := range strings.Split(value, ",") {
		condition := strings.TrimSpace(strings.ToLower(rawCondition))
		if condition == "" {
			continue
		}
		if condition == "eol" {
			condition = failOnEndOfLife
		}

		switch condition {
		case failOnVulnerable, failOnEndOfLife:
			conditions[condition] = true
		default:
			return nil, fmt.Errorf(
				"criterio no soportado en --fail-on: %s (válidos: vulnerable, end-of-life)",
				condition,
			)
		}
	}

	return conditions, nil
}

func scanPolicyViolated(result model.AnalysisResult, conditions map[string]bool) bool {
	if conditions[failOnVulnerable] {
		for _, dependency := range result.Dependencies {
			if len(dependency.Vulnerabilities) > 0 ||
				dependencyStatus(dependency) == string(model.StatusVulnerable) ||
				dependencyStatus(dependency) == string(model.StatusVulnerableUpdate) {
				return true
			}
		}
	}

	if conditions[failOnEndOfLife] {
		for _, dependency := range result.Dependencies {
			if dependencyStatus(dependency) == string(model.StatusEndOfLife) {
				return true
			}
		}
	}

	return false
}
