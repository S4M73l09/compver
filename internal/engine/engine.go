package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/version"
)

type Detector interface {
	CanHandle(path string) bool
	Detect(path string) ([]model.Dependency, error)
}

type Engine struct {
	detectors []Detector
}

func New(detectors ...Detector) *Engine {
	return &Engine{
		detectors: detectors,
	}
}

func (e *Engine) Analyze(root string) (model.AnalysisResult, error) {
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

	err = filepath.WalkDir(root, func(
		path string,
		entry fs.DirEntry,
		walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			if path != root && entry.Name() == ".git" {
				return filepath.SkipDir
			}

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

			for index := range dependencies {
				parsedVersion, parseErr := version.Parse(
					dependencies[index].CurrentVersion,
				)
				if parseErr != nil {
					dependencies[index].VersionKind = version.Unknown
					continue
				}

				dependencies[index].VersionKind = parsedVersion.Kind()
			}

			result.Dependencies = append(
				result.Dependencies,
				dependencies...,
			)

			break
		}

		return nil
	})

	if err != nil {
		return model.AnalysisResult{}, err
	}

	return result, nil
}
