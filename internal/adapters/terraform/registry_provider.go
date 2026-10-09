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

type RegistryProvider struct {
	baseURL string
	client  *http.Client
	cache   cache.Cache
}

func NewRegistryProvider(
	baseURL string,
	client *http.Client,
	providerCache cache.Cache,
) *RegistryProvider {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	return &RegistryProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		cache:   providerCache,
	}
}

func (RegistryProvider) CanHandle(dependency model.Dependency) bool {
	return strings.HasPrefix(
		dependency.Name,
		"registry.terraform.io/",
	) && strings.HasSuffix(dependency.Source, ".terraform.lock.hcl")
}

func (p *RegistryProvider) AvailableVersions(
	ctx context.Context,
	dependency model.Dependency,
	options providers.QueryOptions,
) (providers.Result, error) {
	parts := strings.Split(dependency.Name, "/")
	if len(parts) != 3 || parts[0] != "registry.terraform.io" {
		return providers.Result{}, fmt.Errorf(
			"provider Terraform no compatible: %s",
			dependency.Name,
		)
	}

	cacheKey := "terraform-provider:" + dependency.Name
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

		return registryCacheResult(entry), nil
	}

	cacheTTL := options.CacheTTL
	if cacheTTL <= 0 {
		cacheTTL = 24 * time.Hour
	}
	if options.Mode != providers.NetworkRefresh &&
		cacheErr == nil && found &&
		!entry.RetrievedAt.IsZero() &&
		time.Since(entry.RetrievedAt) <= cacheTTL {
		return registryCacheResult(entry), nil
	}

	url := fmt.Sprintf(
		"%s/v1/providers/%s/%s/versions",
		p.baseURL,
		parts[1],
		parts[2],
	)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return providers.Result{}, err
	}

	response, err := p.client.Do(request)
	if err != nil {
		return providers.Result{}, fmt.Errorf(
			"consultando provider %s: %w",
			dependency.Name,
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return providers.Result{}, fmt.Errorf(
			"el Registry respondió con HTTP %d",
			response.StatusCode,
		)
	}

	var payload struct {
		Versions []struct {
			Version string `json:"version"`
		} `json:"versions"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return providers.Result{}, fmt.Errorf("decodificando respuesta del Registry: %w", err)
	}

	var versions []version.Version
	for _, item := range payload.Versions {
		parsed, err := version.Parse(item.Version)
		if err != nil {
			continue
		}
		versions = append(versions, parsed)
	}

	result := providers.Result{
		Versions:    versions,
		Source:      p.baseURL,
		RetrievedAt: time.Now(),
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

func (p *RegistryProvider) readCache(key string) (cache.Entry, bool, error) {
	if p.cache == nil {
		return cache.Entry{}, false, nil
	}

	return p.cache.Get(key)
}

func registryCacheResult(entry cache.Entry) providers.Result {
	return providers.Result{
		Versions:    entry.Versions,
		Source:      entry.Source,
		RetrievedAt: entry.RetrievedAt,
		FromCache:   true,
	}
}
