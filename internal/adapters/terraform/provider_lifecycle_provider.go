package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/S4M73l09/compver/internal/cache"
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
)

type ProviderLifecycleProvider struct {
	baseURL string
	client  *http.Client
	cache   cache.Cache
}

func NewProviderLifecycleProvider(
	baseURL string,
	client *http.Client,
	providerCache cache.Cache,
) *ProviderLifecycleProvider {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	return &ProviderLifecycleProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		cache:   providerCache,
	}
}

func (ProviderLifecycleProvider) CanHandle(
	dependency model.Dependency,
) bool {
	parts := strings.Split(dependency.Name, "/")

	return len(parts) == 3 &&
		parts[0] == "registry.terraform.io"
}

func (p *ProviderLifecycleProvider) Lifecycle(
	ctx context.Context,
	dependency model.Dependency,
	options providers.QueryOptions,
) (providers.LifecycleResult, error) {
	cacheKey := "terraform-provider-lifecycle:" + dependency.Name
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

		return providerLifecycleCacheResult(entry), nil
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
		return providerLifecycleCacheResult(entry), nil
	}

	parts := strings.Split(dependency.Name, "/")
	owner := parts[1]
	providerName := parts[2]
	repository := fmt.Sprintf(
		"%s/terraform-provider-%s",
		url.PathEscape(owner),
		url.PathEscape(providerName),
	)

	requestURL := fmt.Sprintf(
		"%s/repos/%s/releases?per_page=100",
		p.baseURL,
		repository,
	)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL,
		nil,
	)
	if err != nil {
		return providers.LifecycleResult{}, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := p.client.Do(request)
	if err != nil {
		return providers.LifecycleResult{}, fmt.Errorf(
			"consultando releases de %s: %w",
			dependency.Name,
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return providers.LifecycleResult{}, fmt.Errorf(
			"GitHub respondió con HTTP %d para %s",
			response.StatusCode,
			repository,
		)
	}

	var releases []githubRelease
	if err := json.NewDecoder(response.Body).Decode(&releases); err != nil {
		return providers.LifecycleResult{}, fmt.Errorf(
			"decodificando releases de %s: %w",
			dependency.Name,
			err,
		)
	}

	releaseInfo, err := parseGitHubReleases(releases)
	if err != nil {
		return providers.LifecycleResult{}, err
	}

	result := providers.LifecycleResult{
		Releases:    releaseInfo,
		Source:      requestURL,
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

type githubRelease struct {
	TagName     string     `json:"tag_name"`
	PublishedAt *time.Time `json:"published_at"`
	Draft       bool       `json:"draft"`
	Prerelease  bool       `json:"prerelease"`
}

func parseGitHubReleases(
	githubReleases []githubRelease,
) ([]model.ReleaseInfo, error) {
	releases := make([]model.ReleaseInfo, 0, len(githubReleases))

	for _, githubRelease := range githubReleases {
		if githubRelease.Draft || githubRelease.TagName == "" {
			continue
		}

		parsedVersion, err := version.Parse(githubRelease.TagName)
		if err != nil {
			continue
		}

		release := model.ReleaseInfo{
			Version:     parsedVersion,
			Cycle:       model.ReleaseCycle(parsedVersion),
			PublishedAt: githubRelease.PublishedAt,
			Support:     model.SupportUnknown,
		}

		releases = append(releases, release)
	}

	return releases, nil
}

func (p *ProviderLifecycleProvider) readCache(
	key string,
) (cache.Entry, bool, error) {
	if p.cache == nil {
		return cache.Entry{}, false, nil
	}

	return p.cache.Get(key)
}

func providerLifecycleCacheResult(
	entry cache.Entry,
) providers.LifecycleResult {
	return providers.LifecycleResult{
		Releases:    entry.Releases,
		Source:      entry.Source,
		RetrievedAt: entry.RetrievedAt,
		FromCache:   true,
	}
}
