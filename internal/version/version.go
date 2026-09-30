package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Kind string

const (
	Unknown    Kind = "unknown"
	Stable     Kind = "stable"
	Alpha      Kind = "alpha"
	Beta       Kind = "beta"
	RC         Kind = "rc"
	PreRelease Kind = "pre-release"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
}

func Parse(value string) (Version, error) {
	value = strings.TrimPrefix(value, "v")
	value = strings.SplitN(value, "+", 2)[0]

	parts := strings.SplitN(value, "-", 2)
	numbers := strings.Split(parts[0], ".")
	if len(numbers) != 3 {
		return Version{}, fmt.Errorf("versión inválida: %s", value)
	}

	major, err := strconv.Atoi(numbers[0])
	if err != nil {
		return Version{}, fmt.Errorf("major inválido en %s: %w", value, err)
	}

	minor, err := strconv.Atoi(numbers[1])
	if err != nil {
		return Version{}, fmt.Errorf("minor inválido en %s: %w", value, err)
	}

	patch, err := strconv.Atoi(numbers[2])
	if err != nil {
		return Version{}, fmt.Errorf("patch inválido en %s: %w", value, err)
	}

	result := Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	if len(parts) == 2 {
		result.PreRelease = parts[1]
	}

	return result, nil
}

func (v Version) Kind() Kind {
	if v.PreRelease == "" {
		return Stable
	}

	label := strings.ToLower(strings.SplitN(v.PreRelease, ".", 2)[0])

	switch label {
	case "alpha":
		return Alpha
	case "beta":
		return Beta
	case "rc":
		return RC
	}

	return PreRelease
}

func Compare(left, right Version) int {
	if left.Major != right.Major {
		return compareInt(left.Major, right.Major)
	}

	if left.Minor != right.Minor {
		return compareInt(left.Minor, right.Minor)
	}

	if left.Patch != right.Patch {
		return compareInt(left.Patch, right.Patch)
	}

	if left.PreRelease == right.PreRelease {
		return 0
	}

	if left.PreRelease == "" {
		return 1
	}

	if right.PreRelease == "" {
		return -1
	}

	return strings.Compare(left.PreRelease, right.PreRelease)
}

func compareInt(left, right int) int {
	if left < right {
		return -1
	}

	return 1
}
