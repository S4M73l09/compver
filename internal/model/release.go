package model

import (
	"fmt"
	"time"

	"github.com/S4M73l09/compver/internal/version"
)

type SupportStatus string

const (
	SupportUnknown    SupportStatus = "unknown"
	SupportActive     SupportStatus = "active"
	SupportDeprecated SupportStatus = "deprecated"
	SupportEOL        SupportStatus = "eol"
)

type Vulnerability struct {
	ID       string
	Summary  string
	Severity string
	FixedIn  []version.Version
}

type ReleaseInfo struct {
	Version         version.Version
	Cycle           string
	PublishedAt     *time.Time
	Support         SupportStatus
	SupportEnd      *time.Time
	Vulnerabilities []Vulnerability
}

func ReleaseInfosFromVersions(versions []version.Version) []ReleaseInfo {
	releases := make([]ReleaseInfo, 0, len(versions))

	for _, releaseVersion := range versions {
		releases = append(releases, ReleaseInfo{
			Version: releaseVersion,
			Cycle:   ReleaseCycle(releaseVersion),
			Support: SupportUnknown,
		})
	}

	return releases
}

func ReleaseCycle(releaseVersion version.Version) string {
	return fmt.Sprintf("%d.%d", releaseVersion.Major, releaseVersion.Minor)
}
