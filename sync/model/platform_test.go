package model

import (
	"testing"
	"time"
)

func TestPlatform_TableName(t *testing.T) {
	p := Platform{}
	expected := "platforms"

	if p.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, p.TableName())
	}
}

func TestPlatform_Fields(t *testing.T) {
	now := time.Now()
	p := Platform{
		ID:             1,
		Key:            "github",
		Name:           "GitHub",
		Type:           "github",
		InstanceURL:    "https://github.com",
		APIURL:         "https://api.github.com",
		AccessToken:    "test-token",
		SkipTLSVerify:  false,
		CACertPath:     "/path/to/cert",
		ProxyURL:       "http://proxy:8080",
		IsDefault:      true,
		Status:         "active",
		LastTestResult: "success",
		LastTestAt:     &now,
		RepoCount:      10,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}

	if p.Key != "github" {
		t.Errorf("expected key 'github', got '%s'", p.Key)
	}

	if p.Name != "GitHub" {
		t.Errorf("expected name 'GitHub', got '%s'", p.Name)
	}

	if p.Type != "github" {
		t.Errorf("expected type 'github', got '%s'", p.Type)
	}

	if p.InstanceURL != "https://github.com" {
		t.Errorf("expected instance URL 'https://github.com', got '%s'", p.InstanceURL)
	}

	if p.APIURL != "https://api.github.com" {
		t.Errorf("expected API URL 'https://api.github.com', got '%s'", p.APIURL)
	}

	if p.AccessToken != "test-token" {
		t.Errorf("expected access token 'test-token', got '%s'", p.AccessToken)
	}

	if p.SkipTLSVerify {
		t.Error("expected skip TLS verify to be false")
	}

	if p.CACertPath != "/path/to/cert" {
		t.Errorf("expected CA cert path '/path/to/cert', got '%s'", p.CACertPath)
	}

	if p.ProxyURL != "http://proxy:8080" {
		t.Errorf("expected proxy URL 'http://proxy:8080', got '%s'", p.ProxyURL)
	}

	if !p.IsDefault {
		t.Error("expected is default to be true")
	}

	if p.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", p.Status)
	}

	if p.LastTestResult != "success" {
		t.Errorf("expected last test result 'success', got '%s'", p.LastTestResult)
	}

	if p.LastTestAt == nil {
		t.Error("expected last test at to be set")
	}

	if p.RepoCount != 10 {
		t.Errorf("expected repo count 10, got %d", p.RepoCount)
	}
}

func TestGetAPIURL(t *testing.T) {
	tests := []struct {
		name         string
		platformType string
		instanceURL  string
		expected     string
	}{
		{"github with default", "github", "", "https://github.com/api/v3"},
		{"github with custom", "github", "github.example.com", "https://github.example.com/api/v3"},
		{"gitlab with default", "gitlab", "", "https://gitlab.com/api/v4"},
		{"gitlab with custom", "gitlab", "gitlab.example.com", "https://gitlab.example.com/api/v4"},
		{"gitea with default", "gitea", "", "https://gitea.com/api/v1"},
		{"gitea with custom", "gitea", "gitea.example.com", "https://gitea.example.com/api/v1"},
		{"gitee with default", "gitee", "", "https://gitee.com/api/v5"},
		{"gitee with custom", "gitee", "gitee.example.com", "https://gitee.example.com/api/v5"},
		{"unknown type", "unknown", "", ""},
		{"unknown type with custom", "unknown", "example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAPIURL(tt.platformType, tt.instanceURL)
			if result != tt.expected {
				t.Errorf("GetAPIURL(%q, %q) = %q, want %q", tt.platformType, tt.instanceURL, result, tt.expected)
			}
		})
	}
}
