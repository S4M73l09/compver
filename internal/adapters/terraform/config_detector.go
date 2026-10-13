package terraform

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/S4M73l09/compver/internal/model"
)

var (
	requiredProvidersPattern = regexp.MustCompile(
		`^required_providers\s*(?:=\s*)?\{`,
	)
	providerEntryPattern = regexp.MustCompile(
		`^([A-Za-z0-9_-]+)\s*=\s*\{`,
	)
	sourcePattern = regexp.MustCompile(
		`^source\s*=\s*"([^"]+)"`,
	)
	constraintPattern = regexp.MustCompile(
		`^version\s*=\s*(.+)$`,
	)
	quotedValuePattern = regexp.MustCompile(`"([^"]+)"`)
)

type ConfigDetector struct{}

func NewConfigDetector() ConfigDetector {
	return ConfigDetector{}
}

func (ConfigDetector) CanHandle(path string) bool {
	return strings.HasSuffix(path, ".tf")
}

func (ConfigDetector) Detect(path string) ([]model.Dependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dependencies []model.Dependency
	var providerName string
	var providerSource string
	var constraint string
	var inRequiredProviders bool
	var inProviderEntry bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(stripComment(scanner.Text()))
		if line == "" {
			continue
		}

		if !inRequiredProviders {
			if requiredProvidersPattern.MatchString(line) {
				inRequiredProviders = true
			}
			continue
		}

		if !inProviderEntry {
			if line == "}" {
				inRequiredProviders = false
				continue
			}

			match := providerEntryPattern.FindStringSubmatch(line)
			if match == nil {
				continue
			}

			providerName = match[1]
			providerSource = ""
			constraint = ""
			inProviderEntry = true
			continue
		}

		if line == "}" {
			if constraint != "" {
				if providerSource == "" {
					providerSource = "hashicorp/" + providerName
				}

				parsedConstraint, err := ParseTerraformConstraint(constraint)
				if err != nil {
					return nil, err
				}

				dependencies = append(dependencies, model.Dependency{
					Tool:           "terraform",
					Name:           normalizeProviderSource(providerSource),
					CurrentVersion: constraint,
					Constraint:     parsedConstraint,
					Source:         path,
					Observations: []model.VersionObservation{{
						Role:    "constraint",
						Version: constraint,
						Source:  path,
					}},
				})
			}

			providerName = ""
			providerSource = ""
			constraint = ""
			inProviderEntry = false
			continue
		}

		if match := sourcePattern.FindStringSubmatch(line); match != nil {
			providerSource = match[1]
			continue
		}

		if match := constraintPattern.FindStringSubmatch(line); match != nil {
			values := quotedValuePattern.FindAllStringSubmatch(match[1], -1)
			constraints := make([]string, 0, len(values))
			for _, value := range values {
				constraints = append(constraints, value[1])
			}
			constraint = strings.Join(constraints, ", ")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}

func normalizeProviderSource(source string) string {
	parts := strings.Split(strings.TrimSpace(source), "/")

	switch len(parts) {
	case 1:
		return "registry.terraform.io/hashicorp/" + parts[0]
	case 2:
		return "registry.terraform.io/" + parts[0] + "/" + parts[1]
	default:
		return strings.Join(parts, "/")
	}
}

func stripComment(line string) string {
	if index := strings.Index(line, "#"); index >= 0 {
		return line[:index]
	}

	return line
}
