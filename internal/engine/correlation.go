package engine

import (
	"strings"

	"github.com/S4M73l09/compver/internal/model"
)

// MergeDependencies groups observations that describe the same dependency.
// The tool and dependency name form the identity of a dependency.
func MergeDependencies(dependencies []model.Dependency) []model.Dependency {
	merged := make([]model.Dependency, 0, len(dependencies))
	indexes := make(map[string]int)

	for _, dependency := range dependencies {
		if len(dependency.Observations) == 0 && dependency.CurrentVersion != "" {
			dependency.Observations = []model.VersionObservation{{
				Role:    "current",
				Version: dependency.CurrentVersion,
				Source:  dependency.Source,
			}}
		}

		key := dependencyKey(dependency)
		index, exists := indexes[key]
		if !exists {
			indexes[key] = len(merged)
			merged = append(merged, dependency)
			continue
		}

		mergeDependency(&merged[index], dependency)
	}

	return merged
}

func dependencyKey(dependency model.Dependency) string {
	return strings.Join([]string{dependency.Tool, dependency.Name}, "\x00")
}

func mergeDependency(target *model.Dependency, source model.Dependency) {
	if target.Constraint.Empty() && !source.Constraint.Empty() {
		target.Constraint = source.Constraint
	}

	if target.CurrentVersion == "" {
		target.CurrentVersion = source.CurrentVersion
	}

	target.Indirect = target.Indirect || source.Indirect
	target.Observations = appendUniqueObservations(
		target.Observations,
		source.Observations,
	)
}

func appendUniqueObservations(
	target []model.VersionObservation,
	source []model.VersionObservation,
) []model.VersionObservation {
	for _, candidate := range source {
		duplicate := false
		for _, existing := range target {
			if existing == candidate {
				duplicate = true
				break
			}
		}

		if !duplicate {
			target = append(target, candidate)
		}
	}

	return target
}
