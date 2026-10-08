package version

type ComparisonStatus string

const (
	StatusCurrent           ComparisonStatus = "current"
	StatusUpdateAvailable   ComparisonStatus = "update-available"
	StatusCurrentPreRelease ComparisonStatus = "current-pre-release"
	StatusNoVersions        ComparisonStatus = "no-versions"
	StatusUnknownCurrent    ComparisonStatus = "unknown-current"
)

type Comparison struct {
	Status  ComparisonStatus
	Current *Version
	Latest  *Version
}

func CompareCurrent(current string, available []Version) Comparison {
	parsedCurrent, err := Parse(current)
	if err != nil {
		return Comparison{Status: StatusUnknownCurrent}
	}

	if len(available) == 0 {
		return Comparison{
			Status:  StatusNoVersions,
			Current: &parsedCurrent,
		}
	}

	latest := available[0]
	for _, candidate := range available[1:] {
		if Compare(candidate, latest) > 0 {
			latest = candidate
		}
	}

	comparison := Comparison{
		Current: &parsedCurrent,
		Latest:  &latest,
	}

	switch {
	case Compare(latest, parsedCurrent) > 0:
		comparison.Status = StatusUpdateAvailable
	case parsedCurrent.Kind() != Stable:
		comparison.Status = StatusCurrentPreRelease
	default:
		comparison.Status = StatusCurrent
	}

	return comparison
}
