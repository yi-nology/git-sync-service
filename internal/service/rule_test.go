package service

import (
	"testing"
)

func TestMatchBranch(t *testing.T) {
	tests := []struct {
		pattern string
		branch  string
		want    bool
	}{
		{"", "main", true},
		{"*", "main", true},
		{"main", "main", true},
		{"main", "develop", false},
		{"feature/*", "feature/123", true},
		{"feature/*", "main", false},
		{"release/*,hotfix/*", "release/v1.0", true},
		{"release/*,hotfix/*", "hotfix/fix-123", true},
		{"release/*,hotfix/*", "develop", false},
		{"main,develop", "main", true},
		{"main,develop", "develop", true},
		{"main,develop", "feature", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"/"+tt.branch, func(t *testing.T) {
			got := matchBranch(tt.pattern, tt.branch)
			if got != tt.want {
				t.Errorf("matchBranch(%q, %q) = %v, want %v", tt.pattern, tt.branch, got, tt.want)
			}
		})
	}
}

func TestMatchEventType(t *testing.T) {
	tests := []struct {
		pattern string
		actual  string
		want    bool
	}{
		{"", "push", true},
		{"*", "push", true},
		{"push", "push", true},
		{"push", "pull_request", false},
		{"push,pull_request", "push", false},
	}

	for _, tt := range tests {
		got := matchEventType(tt.pattern, tt.actual)
		if got != tt.want {
			t.Errorf("matchEventType(%q, %q) = %v, want %v", tt.pattern, tt.actual, got, tt.want)
		}
	}
}
