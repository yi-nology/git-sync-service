package git_sync

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

// BatchRequest represents a batch operation request.
type BatchRequest struct {
	Action string   `json:"action" vd:"required"`
	Keys   []string `json:"keys" vd:"required"`
}

// BatchRepos handles batch operations on repositories.
func BatchRepos(ctx context.Context, c *app.RequestContext) {
	var req BatchRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if len(req.Keys) == 0 {
		response.BadRequest(c, "keys must not be empty")
		return
	}

	switch req.Action {
	case "delete":
		batchDelete(ctx, c, req.Keys)
	default:
		response.BadRequest(c, fmt.Sprintf("unsupported action: %s", req.Action))
	}
}

func batchDelete(ctx context.Context, c *app.RequestContext, keys []string) {
	total := len(keys)
	success := 0
	failed := 0
	var errors []string

	for _, key := range keys {
		if err := GetSyncService().DeleteRepo(ctx, key); err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s: %v", key, err))
		} else {
			success++
		}
	}

	response.Success(c, &response.BatchResponse{
		Total:   total,
		Success: success,
		Failed:  failed,
		Errors:  errors,
	})
}
