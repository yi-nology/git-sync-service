package git_sync

import (
	"context"
	"log/slog"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	git_sync "github.com/yi-nology/git-sync-service/biz/handler/git_sync"
)

var apiKey string

func SetAPIKey(key string) {
	apiKey = key
}

func rootMw() []app.HandlerFunc {
	return nil
}

func _apiMw() []app.HandlerFunc {
	return []app.HandlerFunc{corsMiddleware()}
}

func _v1Mw() []app.HandlerFunc {
	return []app.HandlerFunc{authMiddleware()}
}

func corsMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		c.Header("Access-Control-Max-Age", "86400")

		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}

		c.Next(ctx)
	}
}

func authMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Path())
		if strings.HasPrefix(path, "/api/v1/webhook/receive/") {
			c.Next(ctx)
			return
		}

		key := string(c.GetHeader("X-API-Key"))
		if key == "" {
			key = c.Query("api_key")
		}

		if key == "" {
			c.AbortWithStatusJSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "API key required. Provide X-API-Key header or api_key query parameter.",
			})
			return
		}

		svc := git_sync.GetSyncService()
		configuredKey := ""
		if svc != nil {
			configuredKey = svc.GetAPIKey()
		}
		if configuredKey != "" && key != configuredKey {
			slog.Warn("invalid API key attempt", "ip", c.ClientIP(), "path", path)
			c.AbortWithStatusJSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "invalid API key",
			})
			return
		}

		c.Next(ctx)
	}
}

func _repoMw() []app.HandlerFunc          { return nil }
func _getrepoMw() []app.HandlerFunc       { return nil }
func _listbranchesMw() []app.HandlerFunc  { return nil }
func _createrepoMw() []app.HandlerFunc    { return nil }
func _deleterepoMw() []app.HandlerFunc    { return nil }
func _testconnectionMw() []app.HandlerFunc { return nil }
func _updaterepoMw() []app.HandlerFunc    { return nil }
func _listreposMw() []app.HandlerFunc     { return nil }
func _syncMw() []app.HandlerFunc          { return nil }
func _listhistoryMw() []app.HandlerFunc   { return nil }
func _previewsyncMw() []app.HandlerFunc   { return nil }
func _taskMw() []app.HandlerFunc          { return nil }
func _gettaskMw() []app.HandlerFunc       { return nil }
func _createtaskMw() []app.HandlerFunc    { return nil }
func _deletetaskMw() []app.HandlerFunc    { return nil }
func _runtaskMw() []app.HandlerFunc       { return nil }
func _updatetaskMw() []app.HandlerFunc    { return nil }
func _listtasksMw() []app.HandlerFunc     { return nil }
func _webhookMw() []app.HandlerFunc       { return nil }
func _listeventsMw() []app.HandlerFunc    { return nil }
func _ruleMw() []app.HandlerFunc          { return nil }
func _getruleMw() []app.HandlerFunc       { return nil }
func _createruleMw() []app.HandlerFunc    { return nil }
func _deleteruleMw() []app.HandlerFunc    { return nil }
func _updateruleMw() []app.HandlerFunc    { return nil }
func _listrulesMw() []app.HandlerFunc     { return nil }
func _eventMw() []app.HandlerFunc         { return nil }
func _retryeventMw() []app.HandlerFunc    { return nil }
