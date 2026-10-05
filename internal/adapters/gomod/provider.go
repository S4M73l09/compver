package gomod

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
)

type Provider struct {
	baseURL string
	client  *http.Client
}

func NewProvider(baseURL string, client *http.Client) *Provider {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Provider{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
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
	if options.Mode == providers.NetworkOffline {
		return providers.Result{}, fmt.Errorf("modo offline activo")
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

	return providers.Result{
		Versions:    versions,
		Source:      p.baseURL,
		RetrievedAt: time.Now(),
		FromCache:   false,
	}, nil
}
