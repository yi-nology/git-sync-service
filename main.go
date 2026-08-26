package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yi-nology/git-sync-service/biz/handler/git_sync"
	"github.com/yi-nology/git-sync-service/sync"

	// Register all platform backends (GitHub, GitLab, Gitea, etc.)
	_ "github.com/yi-nology/git-platform-sdk/backends/all"
)

var syncSvc *sync.Service

func main() {
	cfg, err := sync.LoadConfig("conf/config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	// Configure logger based on config
	var logLevel slog.Level
	switch cfg.Log.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	if cfg.Log.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}
	slog.SetDefault(slog.New(handler))

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

	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)),
		server.WithMaxRequestBodySize(cfg.Webhook.MaxBodySize),
	)

	// Recovery 中间件:handler panic 时返回 500 而非崩进程
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("handler panic recovered", "error", r, "path", string(ctx.Path()))
				ctx.JSON(500, map[string]string{"error": "internal server error"})
				ctx.Abort()
			}
		}()
		ctx.Next(c)
	})

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
