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

func (v Version) String() string {
	result := fmt.Sprintf(
		"v%d.%d.%d",
		v.Major,
		v.Minor,
		v.Patch,
	)

	if v.PreRelease != "" {
		result += "-" + v.PreRelease
	}

	return result
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

	if left.PreRelease == "" {
		if right.PreRelease == "" {
			return 0
		}

		return 1
	}

	if right.PreRelease == "" {
		return -1
	}

	leftIdentifiers := strings.Split(left.PreRelease, ".")
	rightIdentifiers := strings.Split(right.PreRelease, ".")
	identifierCount := len(leftIdentifiers)
	if len(rightIdentifiers) < identifierCount {
		identifierCount = len(rightIdentifiers)
	}

	for index := 0; index < identifierCount; index++ {
		comparison := comparePreReleaseIdentifier(
			leftIdentifiers[index],
			rightIdentifiers[index],
		)
		if comparison != 0 {
			return comparison
		}
	}

	if len(leftIdentifiers) < len(rightIdentifiers) {
		return -1
	}

	if len(leftIdentifiers) > len(rightIdentifiers) {
		return 1
	}

	return 0
}

func comparePreReleaseIdentifier(left, right string) int {
	leftNumeric := isNumeric(left)
	rightNumeric := isNumeric(right)

	if leftNumeric && rightNumeric {
		leftNumber, leftErr := strconv.Atoi(left)
		rightNumber, rightErr := strconv.Atoi(right)
		if leftErr == nil && rightErr == nil {
			return compareInt(leftNumber, rightNumber)
		}
	}

	if leftNumeric != rightNumeric {
		if leftNumeric {
			return -1
		}

		return 1
	}

	return strings.Compare(left, right)
}

func isNumeric(value string) bool {
	if value == "" {
		return false
	}

	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}

	return true
}

func compareInt(left, right int) int {
	if left < right {
		return -1
	}

	if left > right {
		return 1
	}

	return 0
}
