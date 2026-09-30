package version

import "testing"

func TestParseStableVersion(t *testing.T) {
	parsed, err := Parse("v1.2.3")
	if err != nil {
		t.Fatal(err)
	}

	if parsed.Major != 1 || parsed.Minor != 2 || parsed.Patch != 3 {
		t.Fatalf("versión inesperada: %+v", parsed)
	}

	if parsed.Kind() != Stable {
		t.Fatalf("se esperaba una versión estable, se obtuvo %s", parsed.Kind())
	}
}

func TestParsePreReleaseVersion(t *testing.T) {
	parsed, err := Parse("1.2.3-rc.1")
	if err != nil {
		t.Fatal(err)
	}

	if parsed.PreRelease != "rc.1" {
		t.Fatalf("pre-release inesperada: %s", parsed.PreRelease)
	}

	if parsed.Kind() != RC {
		t.Fatalf("se esperaba una versión rc, se obtuvo %s", parsed.Kind())
	}
}

func TestClassifyPreReleaseVersions(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected Kind
	}{
		{name: "alpha", value: "1.2.3-alpha.1", expected: Alpha},
		{name: "beta", value: "1.2.3-beta.1", expected: Beta},
		{name: "rc", value: "1.2.3-rc.1", expected: RC},
		{name: "other", value: "1.2.3-nightly.1", expected: PreRelease},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := Parse(test.value)
			if err != nil {
				t.Fatal(err)
			}

			if parsed.Kind() != test.expected {
				t.Fatalf(
					"se esperaba %s, se obtuvo %s",
					test.expected,
					parsed.Kind(),
				)
			}
		})
	}
}

func TestStableIsNewerThanPreRelease(t *testing.T) {
	stable, err := Parse("1.2.3")
	if err != nil {
		t.Fatal(err)
	}

	preRelease, err := Parse("1.2.3-rc.1")
	if err != nil {
		t.Fatal(err)
	}

	if Compare(stable, preRelease) <= 0 {
		t.Fatal("una versión estable debería ser posterior a su pre-release")
	}
}
