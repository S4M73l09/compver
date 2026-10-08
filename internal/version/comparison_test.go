package version

import "testing"

func TestCompareCurrentWithUpdate(t *testing.T) {
	current, err := Parse("1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	latest, err := Parse("1.4.0")
	if err != nil {
		t.Fatal(err)
	}

	result := CompareCurrent("1.2.3", []Version{current, latest})

	if result.Status != StatusUpdateAvailable {
		t.Fatalf("estado inesperado: %s", result.Status)
	}

	if result.Latest == nil || result.Latest.String() != "v1.4.0" {
		t.Fatal("la última versión no fue detectada correctamente")
	}
}

func TestCompareCurrentIsUpToDate(t *testing.T) {
	result := CompareCurrent("1.4.0", []Version{
		{Major: 1, Minor: 4, Patch: 0},
	})

	if result.Status != StatusCurrent {
		t.Fatalf("estado inesperado: %s", result.Status)
	}
}

func TestCompareCurrentPreRelease(t *testing.T) {
	result := CompareCurrent("1.4.0-rc.1", []Version{
		{Major: 1, Minor: 4, Patch: 0, PreRelease: "rc.1"},
	})

	if result.Status != StatusCurrentPreRelease {
		t.Fatalf("estado inesperado: %s", result.Status)
	}
}

func TestCompareUnknownCurrent(t *testing.T) {
	result := CompareCurrent("not-a-version", nil)

	if result.Status != StatusUnknownCurrent {
		t.Fatalf("estado inesperado: %s", result.Status)
	}
}
