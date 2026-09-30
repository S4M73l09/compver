package model

import "github.com/S4M73l09/compver/internal/version"

type Dependency struct {
	Name           string
	CurrentVersion string
	Source         string
	Indirect       bool
	VersionKind    version.Kind
}

type AnalysisResult struct {
	Path         string
	Dependencies []Dependency
}
