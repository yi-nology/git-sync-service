package sync

import (
	"github.com/yi-nology/git-sync-service/internal/service"
	"github.com/yi-nology/git-sync-service/sync/model"
)

type Service = service.Service

type Config = model.Config

func NewService(cfg *Config) (*Service, error) {
	return service.NewService(cfg)
}

func LoadConfig(path string) (*Config, error) {
	return model.LoadConfig(path)
}
