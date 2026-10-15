package engine

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var defaultIgnoredDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".cache":       true,
	"dist":         true,
	"build":        true,
}

type IgnoreMatcher struct {
	root     string
	patterns []string
}

func NewIgnoreMatcher(root string) (*IgnoreMatcher, error) {
	matcher := &IgnoreMatcher{
		root: root,
	}

	ignorePath := filepath.Join(root, ".compverignore")
	file, err := os.Open(ignorePath)
	if os.IsNotExist(err) {
		return matcher, nil
	}

	if err != nil {
		return nil, fmt.Errorf("abriendo %s: %w", ignorePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		pattern := strings.TrimSpace(scanner.Text())
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		pattern = filepath.ToSlash(pattern)
		pattern = strings.TrimSuffix(pattern, "/")

		if _, err := path.Match(pattern, "example"); err != nil {
			return nil, fmt.Errorf(
				"patrón inválido en %s:%d: %w",
				ignorePath,
				lineNumber,
				err,
			)
		}

		matcher.patterns = append(matcher.patterns, pattern)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("leyendo %s: %w", ignorePath, err)
	}

	return matcher, nil
}

func (m *IgnoreMatcher) ShouldSkip(
	pathname string,
	entry fs.DirEntry,
) bool {
	if entry.IsDir() && defaultIgnoredDirectories[entry.Name()] {
		return true
	}

	relativePath, err := filepath.Rel(m.root, pathname)
	if err != nil {
		return false
	}

	relativePath = filepath.ToSlash(relativePath)

	for _, pattern := range m.patterns {
		if matchesIgnorePattern(pattern, relativePath, entry.Name()) {
			return true
		}
	}

	return false
}

func matchesIgnorePattern(
	pattern string,
	relativePath string,
	entryName string,
) bool {
	if !strings.Contains(pattern, "/") {
		matched, _ := path.Match(pattern, entryName)
		return matched
	}

	matched, _ := path.Match(pattern, relativePath)
	return matched
}
