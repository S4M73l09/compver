package terraform

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/S4M73l09/compver/internal/cache"
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
)

var releasePattern = regexp.MustCompile(
	`terraform_([0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.]+)?)`,
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
		client = &http.Client{Timeout: 5 * time.Second}
	}

	return &Provider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		cache:   providerCache,
	}
}

func (Provider) CanHandle(dependency model.Dependency) bool {
	return dependency.Name == "terraform" &&
		strings.HasSuffix(dependency.Source, ".terraform-version")
}

func (p *Provider) AvailableVersions(
	ctx context.Context,
	dependency model.Dependency,
	options providers.QueryOptions,
) (providers.Result, error) {
	const cacheKey = "terraform:terraform"

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
		cacheErr == nil && found &&
		!entry.RetrievedAt.IsZero() &&
		time.Since(entry.RetrievedAt) <= cacheTTL {
		return cacheResult(entry), nil
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		p.baseURL,
		nil,
	)
	if err != nil {
		return providers.Result{}, err
	}

	response, err := p.client.Do(request)
	if err != nil {
		return providers.Result{}, fmt.Errorf(
			"consultando versiones de Terraform: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return providers.Result{}, fmt.Errorf(
			"el registro de Terraform respondió con HTTP %d",
			response.StatusCode,
		)
	}

	versions, err := parseVersions(response.Body)
	if err != nil {
		return providers.Result{}, err
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

func parseVersions(scannerSource interface{ Read([]byte) (int, error) }) ([]version.Version, error) {
	var versions []version.Version
	scanner := bufio.NewScanner(scannerSource)
	seen := make(map[string]struct{})

	for scanner.Scan() {
		matches := releasePattern.FindAllStringSubmatch(scanner.Text(), -1)
		for _, match := range matches {
			parsed, err := version.Parse(match[1])
			if err != nil {
				continue
			}
			if _, exists := seen[parsed.String()]; exists {
				continue
			}
			seen[parsed.String()] = struct{}{}
			versions = append(versions, parsed)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func (p *Provider) readCache(key string) (cache.Entry, bool, error) {
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
