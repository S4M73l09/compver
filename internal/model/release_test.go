package model

import (
	"testing"

	"github.com/S4M73l09/compver/internal/version"
)

func TestReleaseInfoSupportsLifecycleAndVulnerabilities(t *testing.T) {
	releaseVersion, err := version.Parse("1.2.3")
	if err != nil {
		t.Fatal(err)
	}

	release := ReleaseInfo{
		Version: releaseVersion,
		Support: SupportActive,
		Vulnerabilities: []Vulnerability{{
			ID:       "CVE-2026-0001",
			Summary:  "Example vulnerability",
			Severity: "high",
		}},
	}

	if release.Support != SupportActive {
		t.Fatalf("estado de soporte inesperado: %s", release.Support)
	}

	if len(release.Vulnerabilities) != 1 {
		t.Fatalf(
			"se esperaba una vulnerabilidad, se obtuvieron %d",
			len(release.Vulnerabilities),
		)
	}
}
