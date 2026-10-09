package terraform

import (
	"bufio"
	"fmt"
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
	return filepath.Base(path) == ".terraform-version"
}

func (Detector) Detect(path string) ([]model.Dependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if value == "" || strings.HasPrefix(value, "#") {
			continue
		}

		return []model.Dependency{{
			Tool:           "terraform",
			Name:           "terraform",
			CurrentVersion: value,
			Source:         path,
		}}, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no se encontró una versión de Terraform en %s", path)
}
