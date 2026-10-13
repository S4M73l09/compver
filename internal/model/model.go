package model

import "github.com/S4M73l09/compver/internal/version"

type VersionObservation struct {
	Role    string
	Version string
	Source  string
}

type Dependency struct {
	Tool                   string
	Name                   string
	CurrentVersion         string
	Constraint             version.Constraint
	Source                 string
	Indirect               bool
	VersionKind            version.Kind
	AvailableVersions      []version.Version
	AvailableReleases      []ReleaseInfo
	SelectedVersions       []version.Version
	ProviderSource         string
	ProviderError          string
	ProviderFromCache      bool
	LifecycleSource        string
	LifecycleError         string
	LifecycleFromCache     bool
	Vulnerabilities        []Vulnerability
	VulnerabilitySource    string
	VulnerabilityError     string
	VulnerabilityFromCache bool
	Observations           []VersionObservation
	Evaluation             *Evaluation
}

type AnalysisResult struct {
	Path         string
	Dependencies []Dependency
}
