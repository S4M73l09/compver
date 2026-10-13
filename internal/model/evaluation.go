package model

import "github.com/S4M73l09/compver/internal/version"

type EvaluationStatus string

const (
	StatusUnknown          EvaluationStatus = "unknown"
	StatusCurrent          EvaluationStatus = "current"
	StatusUpdateAvailable  EvaluationStatus = "update-available"
	StatusCompatibleUpdate EvaluationStatus = "compatible-update-available"
	StatusDeprecated       EvaluationStatus = "deprecated"
	StatusEndOfLife        EvaluationStatus = "end-of-life"
	StatusVulnerable       EvaluationStatus = "vulnerable"
	StatusVulnerableUpdate EvaluationStatus = "vulnerable-update-available"
)

type Evaluation struct {
	Status           EvaluationStatus
	Current          *version.Version
	LatestPublished  *ReleaseInfo
	LatestCompatible *ReleaseInfo
	Support          SupportStatus
	Vulnerabilities  []Vulnerability
}
