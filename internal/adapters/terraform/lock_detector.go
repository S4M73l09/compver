package terraform

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/S4M73l09/compver/internal/model"
)

var (
	providerBlockPattern = regexp.MustCompile(
		`^provider\s+"([^"]+)"\s*\{`,
	)
	versionPattern = regexp.MustCompile(
		`^version\s*=\s*"([^"]+)"`,
	)
)

type LockDetector struct{}

func NewLockDetector() LockDetector {
	return LockDetector{}
}

func (LockDetector) CanHandle(path string) bool {
	return strings.HasSuffix(path, ".terraform.lock.hcl")
}

func (LockDetector) Detect(path string) ([]model.Dependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dependencies []model.Dependency
	var providerName string
	var inProviderBlock bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if match := providerBlockPattern.FindStringSubmatch(line); match != nil {
			providerName = match[1]
			inProviderBlock = true
			continue
		}

		if inProviderBlock && line == "}" {
			providerName = ""
			inProviderBlock = false
			continue
		}

		if !inProviderBlock {
			continue
		}

		match := versionPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		dependencies = append(dependencies, model.Dependency{
			Tool:           "terraform",
			Name:           providerName,
			CurrentVersion: match[1],
			Source:         path,
		})
		providerName = ""
		inProviderBlock = false
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}
