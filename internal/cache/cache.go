package cache

import (
	"time"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

type Entry struct {
	Key             string
	Versions        []version.Version
	Releases        []model.ReleaseInfo
	Vulnerabilities []model.Vulnerability
	Source          string
	RetrievedAt     time.Time
}

type Cache interface {
	Get(key string) (Entry, bool, error)
	Set(key string, entry Entry) error
	Delete(key string) error
}
