package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/S4M73l09/compver/internal/cache"
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
)

type LifecycleProvider struct {
	baseURL string
	client  *http.Client
	cache   cache.Cache
}

func NewLifecycleProvider(
	baseURL string,
	client *http.Client,
	providerCache cache.Cache,
) *LifecycleProvider {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	return &LifecycleProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		cache:   providerCache,
	}
}

func (*LifecycleProvider) CanHandle(dependency model.Dependency) bool {
	return dependency.Tool == "terraform" && dependency.Name == "terraform"
}

func (p *LifecycleProvider) Lifecycle(
	ctx context.Context,
	dependency model.Dependency,
	options providers.QueryOptions,
) (providers.LifecycleResult, error) {
	cacheKey := "terraform-lifecycle:terraform"
	entry, found, cacheErr := p.readCache(cacheKey)

	if options.Mode == providers.NetworkOffline {
		if cacheErr != nil {
			return providers.LifecycleResult{}, cacheErr
		}
		if !found || len(entry.Releases) == 0 {
			return providers.LifecycleResult{}, fmt.Errorf(
				"no hay datos de lifecycle en caché para %s",
				dependency.Name,
			)
		}

		return lifecycleCacheResult(entry), nil
	}

	cacheTTL := options.CacheTTL
	if cacheTTL <= 0 {
		cacheTTL = 24 * time.Hour
	}

	if options.Mode != providers.NetworkRefresh &&
		cacheErr == nil && found &&
		!entry.RetrievedAt.IsZero() &&
		time.Since(entry.RetrievedAt) <= cacheTTL &&
		len(entry.Releases) > 0 {
		return lifecycleCacheResult(entry), nil
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL, nil)
	if err != nil {
		return providers.LifecycleResult{}, err
	}

	response, err := p.client.Do(request)
	if err != nil {
		return providers.LifecycleResult{}, fmt.Errorf(
			"consultando lifecycle de Terraform: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return providers.LifecycleResult{}, fmt.Errorf(
			"la fuente de lifecycle respondió con HTTP %d",
			response.StatusCode,
		)
	}

	var records []terraformLifecycleRecord
	if err := json.NewDecoder(response.Body).Decode(&records); err != nil {
		return providers.LifecycleResult{}, fmt.Errorf(
			"decodificando lifecycle de Terraform: %w",
			err,
		)
	}

	releases, err := parseTerraformLifecycle(records, time.Now())
	if err != nil {
		return providers.LifecycleResult{}, err
	}

	result := providers.LifecycleResult{
		Releases:    releases,
		Source:      p.baseURL,
		RetrievedAt: time.Now(),
	}

	if p.cache != nil {
		_ = p.cache.Set(cacheKey, cache.Entry{
			Releases:    result.Releases,
			Source:      result.Source,
			RetrievedAt: result.RetrievedAt,
		})
	}

	return result, nil
}

type terraformLifecycleRecord struct {
	Cycle             string          `json:"cycle"`
	Latest            string          `json:"latest"`
	LatestReleaseDate string          `json:"latestReleaseDate"`
	EOL               json.RawMessage `json:"eol"`
}

func parseTerraformLifecycle(
	records []terraformLifecycleRecord,
	now time.Time,
) ([]model.ReleaseInfo, error) {
	releases := make([]model.ReleaseInfo, 0, len(records))

	for _, record := range records {
		latest, err := version.Parse(record.Latest)
		if err != nil {
			continue
		}

		release := model.ReleaseInfo{
			Version: latest,
			Cycle:   record.Cycle,
			Support: model.SupportActive,
		}

		if record.LatestReleaseDate != "" {
			publishedAt, err := parseLifecycleDate(record.LatestReleaseDate)
			if err != nil {
				return nil, fmt.Errorf(
					"fecha de release inválida para Terraform %s: %w",
					record.Cycle,
					err,
				)
			}
			release.PublishedAt = &publishedAt
		}

		supportEnd, hasEOL, err := parseEOL(record.EOL)
		if err != nil {
			return nil, fmt.Errorf(
				"fecha EOL inválida para Terraform %s: %w",
				record.Cycle,
				err,
			)
		}

		if hasEOL {
			release.SupportEnd = &supportEnd
			if now.After(supportEnd) {
				release.Support = model.SupportEOL
			}
		}

		releases = append(releases, release)
	}

	return releases, nil
}

func parseLifecycleDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func parseEOL(raw json.RawMessage) (time.Time, bool, error) {
	if string(raw) == "false" || string(raw) == "null" || len(raw) == 0 {
		return time.Time{}, false, nil
	}

	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return time.Time{}, false, err
	}

	date, err := parseLifecycleDate(value)
	if err != nil {
		return time.Time{}, false, err
	}

	return date, true, nil
}

func (p *LifecycleProvider) readCache(
	key string,
) (cache.Entry, bool, error) {
	if p.cache == nil {
		return cache.Entry{}, false, nil
	}

	return p.cache.Get(key)
}

func lifecycleCacheResult(entry cache.Entry) providers.LifecycleResult {
	return providers.LifecycleResult{
		Releases:    entry.Releases,
		Source:      entry.Source,
		RetrievedAt: entry.RetrievedAt,
		FromCache:   true,
	}
}
