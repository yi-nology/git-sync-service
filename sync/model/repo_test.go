package model

import (
	"testing"
	"time"
)

func TestRepo_TableName(t *testing.T) {
	r := Repo{}
	expected := "repos"

	if r.TableName() != expected {
		t.Errorf("expected table name '%s', got '%s'", expected, r.TableName())
	}
}

func TestRepo_Fields(t *testing.T) {
	now := time.Now()
	r := Repo{
		ID:            1,
		Key:           "test-repo",
		Name:          "Test Repo",
		Platform:      "github",
		PlatformID:    1,
		PlatformOwner: "testuser",
		PlatformRepo:  "test-repo",
		CloneURL:      "https://github.com/testuser/test-repo.git",
		SSHURL:        "git@github.com:testuser/test-repo.git",
		DefaultBranch: "main",
		AccessToken:   "test-token",
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if r.ID != 1 {
		t.Errorf("expected ID 1, got %d", r.ID)
	}

	if r.Key != "test-repo" {
		t.Errorf("expected key 'test-repo', got '%s'", r.Key)
	}

	if r.Name != "Test Repo" {
		t.Errorf("expected name 'Test Repo', got '%s'", r.Name)
	}

	if r.Platform != "github" {
		t.Errorf("expected platform 'github', got '%s'", r.Platform)
	}

	if r.PlatformID != 1 {
		t.Errorf("expected platform ID 1, got %d", r.PlatformID)
	}

	if r.PlatformOwner != "testuser" {
		t.Errorf("expected platform owner 'testuser', got '%s'", r.PlatformOwner)
	}

	if r.PlatformRepo != "test-repo" {
		t.Errorf("expected platform repo 'test-repo', got '%s'", r.PlatformRepo)
	}

	if r.CloneURL != "https://github.com/testuser/test-repo.git" {
		t.Errorf("expected clone URL 'https://github.com/testuser/test-repo.git', got '%s'", r.CloneURL)
	}

	if r.SSHURL != "git@github.com:testuser/test-repo.git" {
		t.Errorf("expected SSH URL 'git@github.com:testuser/test-repo.git', got '%s'", r.SSHURL)
	}

	if r.DefaultBranch != "main" {
		t.Errorf("expected default branch 'main', got '%s'", r.DefaultBranch)
	}

	if r.AccessToken != "test-token" {
		t.Errorf("expected access token 'test-token', got '%s'", r.AccessToken)
	}

	if r.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", r.Status)
	}
}
