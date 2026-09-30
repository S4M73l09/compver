package gomod

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/S4M73l09/compver/internal/model"
)

type Detector struct{}

func New() Detector {
	return Detector{}
}

func (Detector) CanHandle(path string) bool {
	return filepath.Base(path) == "go.mod"
}

func (Detector) Detect(path string) ([]model.Dependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dependencies []model.Dependency
	inRequireBlock := false

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "require (" {
			inRequireBlock = true
			continue
		}

		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		if strings.HasPrefix(line, "require ") {
			line = strings.TrimSpace(
				strings.TrimPrefix(line, "require "),
			)
		} else if !inRequireBlock {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		dependencies = append(dependencies, model.Dependency{
			Name:           fields[0],
			CurrentVersion: fields[1],
			Source:         path,
			Indirect:       strings.Contains(line, "// indirect"),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}
