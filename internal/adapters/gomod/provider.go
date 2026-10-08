package gomod

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/S4M73l09/compver/internal/cache"
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
)

type Provider struct {
	baseURL string
	client  *http.Client
	cache   cache.Cache
}

func NewProvider(
	baseURL string,
	client *http.Client,
	providerCache cache.Cache,
) *Provider {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Provider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		cache:   providerCache,
	}
}

func (Provider) CanHandle(dependency model.Dependency) bool {
	return dependency.Name != ""
}

func (p *Provider) AvailableVersions(
	ctx context.Context,
	dependency model.Dependency,
	options providers.QueryOptions,
) (providers.Result, error) {
	cacheKey := "gomod:" + dependency.Name
	entry, found, cacheErr := p.readCache(cacheKey)

	if options.Mode == providers.NetworkOffline {
		if cacheErr != nil {
			return providers.Result{}, cacheErr
		}
		if !found {
			return providers.Result{}, fmt.Errorf(
				"no hay datos en caché para %s",
				dependency.Name,
			)
		}

		return cacheResult(entry), nil
	}

	cacheTTL := options.CacheTTL
	if cacheTTL <= 0 {
		cacheTTL = 24 * time.Hour
	}

	if options.Mode != providers.NetworkRefresh &&
		cacheErr == nil &&
		found &&
		!entry.RetrievedAt.IsZero() &&
		time.Since(entry.RetrievedAt) <= cacheTTL {
		return cacheResult(entry), nil
	}

	url := fmt.Sprintf(
		"%s/%s/@v/list",
		p.baseURL,
		dependency.Name,
	)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return providers.Result{}, err
	}

	response, err := p.client.Do(request)
	if err != nil {
		return providers.Result{}, fmt.Errorf(
			"consultando versiones de %s: %w",
			dependency.Name,
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return providers.Result{}, fmt.Errorf(
			"el proxy respondió con HTTP %d",
			response.StatusCode,
		)
	}

	var versions []version.Version
	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		rawVersion := strings.TrimSpace(scanner.Text())
		if rawVersion == "" {
			continue
		}

		parsedVersion, err := version.Parse(rawVersion)
		if err != nil {
			// Se ignoran las versiones que no reconoce el parser.
			continue
		}

		versions = append(versions, parsedVersion)
	}

	if err := scanner.Err(); err != nil {
		return providers.Result{}, err
	}

	result := providers.Result{
		Versions:    versions,
		Source:      p.baseURL,
		RetrievedAt: time.Now(),
		FromCache:   false,
	}

	if p.cache != nil {
		_ = p.cache.Set(cacheKey, cache.Entry{
			Versions:    result.Versions,
			Source:      result.Source,
			RetrievedAt: result.RetrievedAt,
		})
	}

	return result, nil
}

func (p *Provider) readCache(
	key string,
) (cache.Entry, bool, error) {
	if p.cache == nil {
		return cache.Entry{}, false, nil
	}

	return p.cache.Get(key)
}

func cacheResult(entry cache.Entry) providers.Result {
	return providers.Result{
		Versions:    entry.Versions,
		Source:      entry.Source,
		RetrievedAt: entry.RetrievedAt,
		FromCache:   true,
	}
}
