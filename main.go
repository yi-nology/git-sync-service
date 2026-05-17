package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yi-nology/git-sync-service/biz/handler/git_sync"
	"github.com/yi-nology/git-sync-service/sync"
)

var syncSvc *sync.Service

func main() {
	cfg, err := sync.LoadConfig("conf/config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	syncSvc, err = sync.NewService(cfg)
	if err != nil {
		slog.Error("init sync service failed", "error", err)
		os.Exit(1)
	}

	if err := syncSvc.Start(); err != nil {
		slog.Error("start sync service failed", "error", err)
		os.Exit(1)
	}
	defer syncSvc.Stop()

	git_sync.SetSyncServiceGetter(func() *sync.Service {
		return syncSvc
	})

	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))

	register(h)

	go func() {
		h.Spin()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
