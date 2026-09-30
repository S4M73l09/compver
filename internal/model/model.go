package model

type Dependency struct {
	Name           string
	CurrentVersion string
	Source         string
}

type AnalysisResult struct {
	Path         string
	Dependencies []Dependency
}
