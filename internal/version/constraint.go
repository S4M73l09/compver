package version

type Operator string

const (
	Equal        Operator = "="
	NotEqual     Operator = "!="
	Greater      Operator = ">"
	GreaterEqual Operator = ">="
	Less         Operator = "<"
	LessEqual    Operator = "<="
)

type Condition struct {
	Operator Operator
	Version  Version
}

type Constraint struct {
	// Cada grupo representa condiciones AND.
	// Los grupos entre sí representan alternativas OR.
	Alternatives [][]Condition
}

func (c Constraint) Empty() bool {
	return len(c.Alternatives) == 0
}

func (c Constraint) Matches(candidate Version) bool {
	for _, alternative := range c.Alternatives {
		matches := true

		for _, condition := range alternative {
			if !matchesCondition(candidate, condition) {
				matches = false
				break
			}
		}

		if matches {
			return true
		}
	}

	return false
}

func FilterByConstraint(
	versions []Version,
	constraint Constraint,
) []Version {
	if constraint.Empty() {
		return versions
	}

	filtered := make([]Version, 0)
	for _, candidate := range versions {
		if constraint.Matches(candidate) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}

func matchesCondition(candidate Version, condition Condition) bool {
	comparison := Compare(candidate, condition.Version)

	switch condition.Operator {
	case Equal:
		return comparison == 0
	case NotEqual:
		return comparison != 0
	case Greater:
		return comparison > 0
	case GreaterEqual:
		return comparison >= 0
	case Less:
		return comparison < 0
	case LessEqual:
		return comparison <= 0
	default:
		return false
	}
}
