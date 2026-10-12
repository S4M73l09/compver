package terraform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/S4M73l09/compver/internal/version"
)

type constraintVersion struct {
	version.Version
	precision int
}

func ParseTerraformConstraint(
	raw string,
) (version.Constraint, error) {
	parts := strings.Split(raw, ",")
	conditions := make([]version.Condition, 0)

	for _, part := range parts {
		operator, value := splitOperator(strings.TrimSpace(part))

		parsed, err := parseConstraintVersion(value)
		if err != nil {
			return version.Constraint{}, err
		}

		if operator == "~>" {
			conditions = append(
				conditions,
				pessimisticConditions(parsed)...,
			)
			continue
		}

		conditions = append(conditions, version.Condition{
			Operator: version.Operator(operator),
			Version:  parsed.Version,
		})
	}

	if len(conditions) == 0 {
		return version.Constraint{}, fmt.Errorf(
			"restricción vacía: %s",
			raw,
		)
	}

	return version.Constraint{
		Alternatives: [][]version.Condition{
			conditions,
		},
	}, nil
}

func splitOperator(value string) (version.Operator, string) {
	for _, operator := range []string{
		"~>", ">=", "<=", "!=", ">", "<", "=",
	} {
		if strings.HasPrefix(value, operator) {
			return version.Operator(operator), strings.TrimSpace(
				strings.TrimPrefix(value, operator),
			)
		}
	}

	return version.Equal, strings.TrimSpace(value)
}

func parseConstraintVersion(value string) (constraintVersion, error) {
	parts := strings.SplitN(strings.TrimSpace(value), "-", 2)
	numbers := strings.Split(parts[0], ".")
	if len(numbers) < 1 || len(numbers) > 3 {
		return constraintVersion{}, fmt.Errorf(
			"versión de restricción inválida: %s",
			value,
		)
	}

	parsedNumbers := []int{0, 0, 0}
	for index, number := range numbers {
		parsed, err := strconv.Atoi(number)
		if err != nil {
			return constraintVersion{}, fmt.Errorf(
				"versión de restricción inválida: %s",
				value,
			)
		}
		parsedNumbers[index] = parsed
	}

	parsed := constraintVersion{
		Version: version.Version{
			Major: parsedNumbers[0],
			Minor: parsedNumbers[1],
			Patch: parsedNumbers[2],
		},
		precision: len(numbers),
	}

	if len(parts) == 2 {
		parsed.PreRelease = parts[1]
	}

	return parsed, nil
}

func pessimisticConditions(
	constraint constraintVersion,
) []version.Condition {
	upperBound := constraint.Version
	if constraint.precision <= 2 {
		upperBound = version.Version{
			Major: constraint.Major + 1,
		}
	} else {
		upperBound = version.Version{
			Major: constraint.Major,
			Minor: constraint.Minor + 1,
		}
	}

	return []version.Condition{
		{
			Operator: version.GreaterEqual,
			Version:  constraint.Version,
		},
		{
			Operator: version.Less,
			Version:  upperBound,
		},
	}
}
