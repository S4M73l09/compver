package engine

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/S4M73l09/compver/internal/evaluator"
	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
	"github.com/S4M73l09/compver/internal/version"
)

type Detector interface {
	CanHandle(path string) bool
	Detect(path string) ([]model.Dependency, error)
}

type Engine struct {
	detectors              []Detector
	providers              []providers.Provider
	lifecycleProviders     []providers.LifecycleProvider
	vulnerabilityProviders []providers.VulnerabilityProvider
}

func New(detectors ...Detector) *Engine {
	return &Engine{
		detectors: detectors,
	}
}

func (e *Engine) Analyze(root string) (model.AnalysisResult, error) {
	return e.AnalyzeContextWithOptions(
		context.Background(),
		root,
		providers.QueryOptions{Mode: providers.NetworkAuto},
	)
}

func (e *Engine) AddProvider(provider providers.Provider) {
	e.providers = append(e.providers, provider)
}

func (e *Engine) AddLifecycleProvider(
	provider providers.LifecycleProvider,
) {
	e.lifecycleProviders = append(e.lifecycleProviders, provider)
}

func (e *Engine) AddVulnerabilityProvider(
	provider providers.VulnerabilityProvider,
) {
	e.vulnerabilityProviders = append(e.vulnerabilityProviders, provider)
}

func (e *Engine) AnalyzeContext(
	ctx context.Context,
	root string,
) (model.AnalysisResult, error) {
	return e.AnalyzeContextWithOptions(
		ctx,
		root,
		providers.QueryOptions{Mode: providers.NetworkAuto},
	)
}

func (e *Engine) AnalyzeWithOptions(
	root string,
	options providers.QueryOptions,
) (model.AnalysisResult, error) {
	return e.AnalyzeContextWithOptions(
		context.Background(),
		root,
		options,
	)
}

func (e *Engine) AnalyzeContextWithOptions(
	ctx context.Context,
	root string,
	options providers.QueryOptions,
) (model.AnalysisResult, error) {
	info, err := os.Stat(root)
	if err != nil {
		return model.AnalysisResult{}, err
	}

	if !info.IsDir() {
		return model.AnalysisResult{}, fmt.Errorf("%s no es un directorio", root)
	}

	result := model.AnalysisResult{
		Path: root,
	}

	ignoreMatcher, err := NewIgnoreMatcher(root)
	if err != nil {
		return model.AnalysisResult{}, err
	}

	var detectedDependencies []model.Dependency

	err = filepath.WalkDir(root, func(
		path string,
		entry fs.DirEntry,
		walkErr error,
	) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if walkErr != nil {
			return walkErr
		}

		if ignoreMatcher.ShouldSkip(path, entry) {
			if entry.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if entry.IsDir() {
			return nil
		}

		for _, detector := range e.detectors {
			if !detector.CanHandle(path) {
				continue
			}

			dependencies, err := detector.Detect(path)
			if err != nil {
				return fmt.Errorf("analizando %s: %w", path, err)
			}

			detectedDependencies = append(
				detectedDependencies,
				dependencies...,
			)

			break
		}

		return nil
	})

	if err != nil {
		return model.AnalysisResult{}, err
	}

	result.Dependencies = MergeDependencies(detectedDependencies)

	for index := range result.Dependencies {
		dependency := &result.Dependencies[index]
		currentVersion := comparableVersion(*dependency)

		parsedVersion, parseErr := version.Parse(currentVersion)
		if parseErr != nil {
			dependency.VersionKind = version.Unknown
		} else {
			dependency.VersionKind = parsedVersion.Kind()
		}

		for _, provider := range e.providers {
			if !provider.CanHandle(*dependency) {
				continue
			}

			providerResult, providerErr := provider.AvailableVersions(
				ctx,
				*dependency,
				options,
			)
			if providerErr != nil {
				dependency.ProviderError = providerErr.Error()
				break
			}

			dependency.AvailableVersions = providerResult.Versions
			dependency.AvailableReleases = providerResult.Releases
			dependency.ProviderSource = providerResult.Source
			dependency.ProviderFromCache = providerResult.FromCache
			break
		}

		for _, lifecycleProvider := range e.lifecycleProviders {
			if !lifecycleProvider.CanHandle(*dependency) {
				continue
			}

			lifecycleResult, lifecycleErr := lifecycleProvider.Lifecycle(
				ctx,
				*dependency,
				options,
			)
			if lifecycleErr != nil {
				dependency.LifecycleError = lifecycleErr.Error()
				break
			}

			dependency.AvailableReleases = mergeReleaseInfo(
				dependency.AvailableVersions,
				dependency.AvailableReleases,
				lifecycleResult.Releases,
			)
			dependency.LifecycleSource = lifecycleResult.Source
			dependency.LifecycleFromCache = lifecycleResult.FromCache
			break
		}

		for _, vulnerabilityProvider := range e.vulnerabilityProviders {
			if !vulnerabilityProvider.CanHandle(*dependency) {
				continue
			}

			vulnerabilityResult, vulnerabilityErr := vulnerabilityProvider.Check(
				ctx,
				*dependency,
				options,
			)
			if vulnerabilityErr != nil {
				dependency.VulnerabilityError = vulnerabilityErr.Error()
				break
			}

			dependency.Vulnerabilities = vulnerabilityResult.Vulnerabilities
			dependency.VulnerabilitySource = vulnerabilityResult.Source
			dependency.VulnerabilityFromCache = vulnerabilityResult.FromCache
			attachVulnerabilities(dependency)
			break
		}

		computedEvaluation := evaluator.Evaluate(*dependency)
		dependency.Evaluation = &computedEvaluation
	}

	return result, nil
}

func comparableVersion(dependency model.Dependency) string {
	for _, role := range []string{"installed", "locked", "current"} {
		for _, observation := range dependency.Observations {
			if observation.Role == role {
				return observation.Version
			}
		}
	}

	return dependency.CurrentVersion
}
