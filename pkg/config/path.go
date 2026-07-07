package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// IsValidName checks if the name contains only alphanumeric characters, underscores, or hyphens.
func IsValidName(name string) bool {
	return nameRegex.MatchString(name)
}

// SafeJoin joins baseDir and subPath while ensuring that target is a subdirectory of baseDir.
func SafeJoin(baseDir, subPath string) (string, error) {
	if filepath.IsAbs(subPath) {
		return "", fmt.Errorf("security error: absolute paths not allowed: %s", subPath)
	}

	cleanBase := filepath.Clean(baseDir)
	joined := filepath.Join(cleanBase, subPath)
	cleanJoined := filepath.Clean(joined)

	baseWithSep := cleanBase
	if !strings.HasSuffix(baseWithSep, string(filepath.Separator)) {
		baseWithSep += string(filepath.Separator)
	}

	if !strings.HasPrefix(cleanJoined, cleanBase) {
		return "", fmt.Errorf("security error: path traversal detected: %s resolves outside %s", subPath, baseDir)
	}

	if cleanJoined != cleanBase && !strings.HasPrefix(cleanJoined, baseWithSep) {
		return "", fmt.Errorf("security error: path traversal detected: %s resolves outside %s", subPath, baseDir)
	}

	return cleanJoined, nil
}
