package version

import "sort"

type SelectionOptions struct {
	Limit              int
	IncludePreReleases bool
	AllVersions        bool
}

func Select(
	versions []Version,
	options SelectionOptions,
) []Version {
	selected := make([]Version, 0, len(versions))

	for _, candidate := range versions {
		if !options.IncludePreReleases &&
			candidate.Kind() != Stable {
			continue
		}

		selected = append(selected, candidate)
	}

	sort.Slice(selected, func(i, j int) bool {
		return Compare(selected[i], selected[j]) > 0
	})

	if options.AllVersions {
		return selected
	}

	limit := options.Limit
	if limit <= 0 {
		limit = 3
	}

	if len(selected) > limit {
		return selected[:limit]
	}

	return selected
}
