package version

import "testing"

func TestConstrainMatches(t *testing.T) {
	constraint := Constraint{
		Alternatives: [][]Condition{
			{
				{
					Operator: GreaterEqual,
					Version:  mustParse("5.0.0"),
				},
				{
					Operator: Less,
					Version:  mustParse("6.0.0"),
				},
			},
		},
	}

	if !constraint.Matches(mustParse("5.2.0")) {
		t.Fatal("5.2.0 deberia cumplir la restriccion")

	}

	if constraint.Matches(mustParse("6.0.0")) {
		t.Fatal("6.0.0 no deberia cumplir la restricción")
	}
}

func mustParse(value string) Version {
	parsed, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return parsed
}
