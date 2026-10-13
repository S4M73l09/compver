package engine

import (
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

func mergeReleaseInfo(
	availableVersions []version.Version,
	availableReleases []model.ReleaseInfo,
	lifecycleReleases []model.ReleaseInfo,
) []model.ReleaseInfo {
	if len(availableReleases) == 0 {
		availableReleases = model.ReleaseInfosFromVersions(availableVersions)
	}

	merged := append([]model.ReleaseInfo(nil), availableReleases...)

	for _, lifecycleRelease := range lifecycleReleases {
		found := false
		for index := range merged {
			if !sameReleaseOrCycle(merged[index], lifecycleRelease) {
				continue
			}

			mergeReleaseMetadata(&merged[index], lifecycleRelease)
			found = true
		}

		if !found {
			merged = append(merged, lifecycleRelease)
		}
	}

	return merged
}

func sameReleaseOrCycle(
	left model.ReleaseInfo,
	right model.ReleaseInfo,
) bool {
	if version.Compare(left.Version, right.Version) == 0 {
		return true
	}

	return left.Cycle != "" && left.Cycle == right.Cycle
}

func mergeReleaseMetadata(
	target *model.ReleaseInfo,
	source model.ReleaseInfo,
) {
	if target.PublishedAt == nil {
		target.PublishedAt = source.PublishedAt
	}

	if target.Support == model.SupportUnknown {
		target.Support = source.Support
	}

	if target.SupportEnd == nil {
		target.SupportEnd = source.SupportEnd
	}

	if len(source.Vulnerabilities) > 0 {
		target.Vulnerabilities = source.Vulnerabilities
	}
}
