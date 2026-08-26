package model

import (
	"fmt"
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
	Driver         string `yaml:"driver"`
	DSN            string `yaml:"dsn"`
	MaxIdleConns   int    `yaml:"max_idle_conns"`
	MaxOpenConns   int    `yaml:"max_open_conns"`
	ConnMaxLifeSec int    `yaml:"conn_max_life_sec"` // 连接最大存活秒数,0 用默认 5 分钟
	ConnMaxIdleSec int    `yaml:"conn_max_idle_sec"` // 空闲连接最大存活秒数,0 用默认 2 分钟
}

type RedisConfig struct {
	Addr           string `yaml:"addr"`
	Password       string `yaml:"password"`
	DB             int    `yaml:"db"`
	PoolSize       int    `yaml:"pool_size"`        // 连接池大小,0 用 go-redis 默认(10*GOMAXPROCS)
	MinIdleConns   int    `yaml:"min_idle_conns"`   // 最小空闲连接数,0 不预热
	DialTimeoutSec int    `yaml:"dial_timeout_sec"` // 建连超时秒数,0 不设超时
	ReadTimeoutSec int    `yaml:"read_timeout_sec"` // 读超时秒数,0 不设超时
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

// LoadConfig 从指定路径加载配置文件。
// path 参数由调用方控制（通常是启动参数或环境变量），不存在用户注入风险。
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // 配置文件路径由程序启动参数决定，非用户可控输入
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		c.Server.Host = DefaultHost
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		c.Server.Port = DefaultPort
	}
	// 拒绝已知的测试默认 API Key,防止裸部署
	if c.Server.APIKey == "test-api-key-123" {
		return fmt.Errorf("refusing to start with default test API key; set a real api_key in config")
	}
	if c.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}
	if c.Database.DSN == "" {
		return fmt.Errorf("database dsn is required")
	}
	if c.Database.Driver != DriverMySQL && c.Database.Driver != DriverSQLite {
		return fmt.Errorf("unsupported database driver: %s", c.Database.Driver)
	}
	if c.Database.MaxIdleConns <= 0 {
		c.Database.MaxIdleConns = DefaultMaxIdleConns
	}
	if c.Database.MaxOpenConns <= 0 {
		c.Database.MaxOpenConns = DefaultMaxOpenConns
	}
	if c.Database.ConnMaxLifeSec <= 0 {
		c.Database.ConnMaxLifeSec = 300 // 5 分钟
	}
	if c.Database.ConnMaxIdleSec <= 0 {
		c.Database.ConnMaxIdleSec = 120 // 2 分钟
	}
	if c.Git.TempDir == "" {
		c.Git.TempDir = DefaultTempDir
	}
	if c.Sync.MaxConcurrent <= 0 {
		c.Sync.MaxConcurrent = DefaultMaxConcurrent
	}
	if c.Sync.DefaultTimeout <= 0 {
		c.Sync.DefaultTimeout = DefaultTimeout
	}
	if c.Sync.RetryCount <= 0 {
		c.Sync.RetryCount = DefaultRetryCount
	}
	if c.Webhook.RateLimit <= 0 {
		c.Webhook.RateLimit = 10
	}
	if c.Webhook.MaxBodySize <= 0 {
		c.Webhook.MaxBodySize = 10 << 20 // 10MB
	}
	// Redis 超时默认值:0 意味着无限阻塞,生产环境必须有超时
	if c.Redis.Addr != "" {
		if c.Redis.DialTimeoutSec <= 0 {
			c.Redis.DialTimeoutSec = 5
		}
		if c.Redis.ReadTimeoutSec <= 0 {
			c.Redis.ReadTimeoutSec = 3
		}
	}
	return nil
}
