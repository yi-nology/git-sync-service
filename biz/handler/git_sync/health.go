package git_sync

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

var startTime = time.Now()

// SystemStatusData 系统状态数据
type SystemStatusData struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Uptime      int64  `json:"uptime"`
	RepoCount   int    `json:"repoCount"`
	TaskCount   int    `json:"taskCount"`
	RunningTask int    `json:"runningTask"`
	FailedTask  int    `json:"failedTask"`
	GoVersion   string `json:"goVersion"`
	Platform    string `json:"platform"`
}

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

// SystemStatus returns the system status information.
func SystemStatus(ctx context.Context, c *app.RequestContext) {
	svc := GetSyncService()
	
	// Get counts (use large limit to get all)
	repos, _, _ := svc.ListRepos(ctx, 0, 10000)
	tasks, _, _ := svc.ListTasks(ctx, "", 0, 10000)
	
	runningCount := 0
	failedCount := 0
	for _, t := range tasks {
		if t.LastStatus == "running" {
			runningCount++
		}
		if t.LastStatus == "failed" {
			failedCount++
		}
	}

	// 使用统一的响应函数，只需要传入 data 部分
	response.Success(c, &SystemStatusData{
		Status:      "running",
		Version:     "v1.5.0",
		Uptime:      int64(time.Since(startTime).Seconds()),
		RepoCount:   len(repos),
		TaskCount:   len(tasks),
		RunningTask: runningCount,
		FailedTask:  failedCount,
		GoVersion:   runtime.Version(),
		Platform:    runtime.GOOS + "/" + runtime.GOARCH,
	})
}
