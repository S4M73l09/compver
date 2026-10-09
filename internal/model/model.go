package model

import "github.com/S4M73l09/compver/internal/version"

type Dependency struct {
	Tool              string
	Name              string
	CurrentVersion    string
	Source            string
	Indirect          bool
	VersionKind       version.Kind
	AvailableVersions []version.Version
	SelectedVersions  []version.Version
	ProviderSource    string
	ProviderError     string
	ProviderFromCache bool
}

type AnalysisResult struct {
	Path         string
	Dependencies []Dependency
}
