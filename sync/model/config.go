package model

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Git      GitConfig      `yaml:"git"`
	Sync     SyncConfig     `yaml:"sync"`
	Webhook  WebhookConfig  `yaml:"webhook"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Mode   string `yaml:"mode"`
	APIKey string `yaml:"api_key"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type GitConfig struct {
	Backend string `yaml:"backend"`
	TempDir string `yaml:"temp_dir"`
}

type SyncConfig struct {
	MaxConcurrent  int `yaml:"max_concurrent"`
	DefaultTimeout int `yaml:"default_timeout"`
	RetryCount     int `yaml:"retry_count"`
}

type WebhookConfig struct {
	RateLimit   int `yaml:"rate_limit"`
	MaxBodySize int `yaml:"max_body_size"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
