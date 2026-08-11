package git_sync

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// HealthCheck returns the health status of the service and its dependencies.
func HealthCheck(ctx context.Context, c *app.RequestContext) {
	status := GetSyncService().HealthCheck()

	httpStatus := http.StatusOK
	for _, v := range status {
		// Accept "ok", "not configured", and "not checked" as healthy states
		if v != "ok" && v != "not configured" && v != "not checked (lock service removed)" {
			httpStatus = http.StatusServiceUnavailable
			break
		}
	}

	c.JSON(httpStatus, map[string]any{
		"status": status,
	})
}
