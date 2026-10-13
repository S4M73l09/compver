package evaluator

import (
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

func Evaluate(dependency model.Dependency) model.Evaluation {
	releases := dependency.AvailableReleases
	if len(releases) == 0 {
		releases = model.ReleaseInfosFromVersions(
			dependency.AvailableVersions,
		)
	}

	evaluation := model.Evaluation{
		Status:  model.StatusUnknown,
		Support: model.SupportUnknown,
	}

	current, err := parseCurrentVersion(dependency)
	if err == nil {
		evaluation.Current = &current
	}

	for index := range releases {
		release := releases[index]
		if evaluation.LatestPublished == nil ||
			version.Compare(release.Version, evaluation.LatestPublished.Version) > 0 {
			evaluation.LatestPublished = &release
		}

		if !dependency.Constraint.Empty() &&
			!dependency.Constraint.Matches(release.Version) {
			continue
		}

		if evaluation.LatestCompatible == nil ||
			version.Compare(release.Version, evaluation.LatestCompatible.Version) > 0 {
			evaluation.LatestCompatible = &release
		}
	}

	if current != (version.Version{}) {
		if currentRelease := findRelease(releases, current); currentRelease != nil {
			evaluation.Support = currentRelease.Support
			evaluation.Vulnerabilities = currentRelease.Vulnerabilities
		}
	}

	evaluation.Status = evaluateStatus(
		evaluation,
		!dependency.Constraint.Empty(),
	)

	return evaluation
}

func parseCurrentVersion(dependency model.Dependency) (version.Version, error) {
	for _, role := range []string{"installed", "locked", "current"} {
		for _, observation := range dependency.Observations {
			if observation.Role != role {
				continue
			}

			parsed, err := version.Parse(observation.Version)
			if err == nil {
				return parsed, nil
			}
		}
	}

	return version.Parse(dependency.CurrentVersion)
}

func findRelease(
	releases []model.ReleaseInfo,
	current version.Version,
) *model.ReleaseInfo {
	for index := range releases {
		if version.Compare(releases[index].Version, current) == 0 {
			return &releases[index]
		}
	}

	return nil
}

func evaluateStatus(
	evaluation model.Evaluation,
	hasConstraint bool,
) model.EvaluationStatus {
	if evaluation.Support == model.SupportEOL {
		if len(evaluation.Vulnerabilities) > 0 {
			return model.StatusVulnerable
		}
		return model.StatusEndOfLife
	}

	if evaluation.Support == model.SupportDeprecated {
		return model.StatusDeprecated
	}

	if len(evaluation.Vulnerabilities) > 0 {
		if evaluation.LatestCompatible != nil &&
			evaluation.Current != nil &&
			version.Compare(
				evaluation.LatestCompatible.Version,
				*evaluation.Current,
			) > 0 {
			return model.StatusVulnerableUpdate
		}
		return model.StatusVulnerable
	}

	if evaluation.Current == nil {
		return model.StatusUnknown
	}

	if hasConstraint && evaluation.LatestCompatible != nil &&
		version.Compare(
			evaluation.LatestCompatible.Version,
			*evaluation.Current,
		) > 0 {
		return model.StatusCompatibleUpdate
	}

	if evaluation.LatestPublished != nil &&
		version.Compare(
			evaluation.LatestPublished.Version,
			*evaluation.Current,
		) > 0 {
		return model.StatusUpdateAvailable
	}

	return model.StatusCurrent
}
