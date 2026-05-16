package service

import (
	"path/filepath"
	"strings"
)

func matchBranch(pattern, branch string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}

	patterns := splitAndTrim(pattern, ",")
	for _, p := range patterns {
		if matchSinglePattern(p, branch) {
			return true
		}
	}

	return false
}

func matchSinglePattern(pattern, branch string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return pattern == branch
	}

	matched, _ := filepath.Match(pattern, branch)
	return matched
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
