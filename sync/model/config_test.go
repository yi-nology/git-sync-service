package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: 0.0.0.0
  port: 8890
  mode: debug

database:
  driver: sqlite
  dsn: test.db

git:
  backend: gogit
  temp_dir: /tmp/test

sync:
  max_concurrent: 5
  default_timeout: 300
  retry_count: 3
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host 0.0.0.0, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8890 {
		t.Errorf("Expected port 8890, got %d", cfg.Server.Port)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("Expected driver sqlite, got %s", cfg.Database.Driver)
	}
	if cfg.Database.DSN != "test.db" {
		t.Errorf("Expected dsn test.db, got %s", cfg.Database.DSN)
	}
}

func TestLoadConfig_Validation(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: 0.0.0.0
  port: 8890

database:
  driver: sqlite
  dsn: test.db

git:
  temp_dir: /tmp/test
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Sync.MaxConcurrent != 5 {
		t.Errorf("Expected default MaxConcurrent 5, got %d", cfg.Sync.MaxConcurrent)
	}
	if cfg.Sync.DefaultTimeout != 300 {
		t.Errorf("Expected default DefaultTimeout 300, got %d", cfg.Sync.DefaultTimeout)
	}
	if cfg.Sync.RetryCount != 3 {
		t.Errorf("Expected default RetryCount 3, got %d", cfg.Sync.RetryCount)
	}
}

func TestLoadConfig_MissingDriver(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: 0.0.0.0
  port: 8890

database:
  dsn: test.db
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("Expected error for missing database driver")
	}
}

func TestLoadConfig_InvalidDriver(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: 0.0.0.0
  port: 8890

database:
  driver: postgres
  dsn: test.db
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("Expected error for invalid database driver")
	}
}

func TestLoadConfig_MissingDSN(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  host: 0.0.0.0
  port: 8890

database:
  driver: sqlite
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("Expected error for missing database dsn")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Expected error for nonexistent config file")
	}
}

func TestInitDB_SQLite(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := InitDB("sqlite", dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	if db == nil {
		t.Fatal("expected non-nil db")
	}

	// Verify tables were created
	if !db.Migrator().HasTable(&Platform{}) {
		t.Error("expected Platform table to exist")
	}

	if !db.Migrator().HasTable(&Repo{}) {
		t.Error("expected Repo table to exist")
	}

	if !db.Migrator().HasTable(&SyncTask{}) {
		t.Error("expected SyncTask table to exist")
	}

	if !db.Migrator().HasTable(&SyncRun{}) {
		t.Error("expected SyncRun table to exist")
	}

	if !db.Migrator().HasTable(&SyncRunStep{}) {
		t.Error("expected SyncRunStep table to exist")
	}

	if !db.Migrator().HasTable(&WebhookRule{}) {
		t.Error("expected WebhookRule table to exist")
	}

	if !db.Migrator().HasTable(&WebhookRuleTask{}) {
		t.Error("expected WebhookRuleTask table to exist")
	}

	if !db.Migrator().HasTable(&WebhookEvent{}) {
		t.Error("expected WebhookEvent table to exist")
	}

	if !db.Migrator().HasTable(&OperationLog{}) {
		t.Error("expected OperationLog table to exist")
	}
}

func TestInitDB_UnsupportedDriver(t *testing.T) {
	_, err := InitDB("postgres", "test.db")
	if err == nil {
		t.Fatal("Expected error for unsupported driver")
	}
}

func TestInitDB_Memory(t *testing.T) {
	db, err := InitDB("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	if db == nil {
		t.Fatal("expected non-nil db")
	}
}
