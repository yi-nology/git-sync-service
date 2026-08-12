package git_sync

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/internal/converter"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

// ListOperationLogsRequest 操作日志查询请求。
type ListOperationLogsRequest struct {
	Page      int    `query:"page" default:"1"`
	PageSize  int    `query:"page_size" default:"10"`
	Search    string `query:"search"`
	Action    string `query:"action"`
	User      string `query:"user"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// OperationLogStats 操作统计。
type OperationLogStats struct {
	Today int64 `json:"today"`
	Week  int64 `json:"week"`
	Total int64 `json:"total"`
}

// ListOperationLogs GET /api/v1/logs/operations
// 返回 {list, pagination, stats}（stats.total 为全量总数，pagination.total 为过滤后的结果数）。
func ListOperationLogs(ctx context.Context, c *app.RequestContext) {
	var req ListOperationLogsRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset, limit := converter.PageToOffset(int32(req.Page), int32(req.PageSize))

	filter := dao.OperationLogFilter{
		Search:    req.Search,
		Action:    req.Action,
		Actor:     req.User,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	list, total, err := GetSyncService().ListOperations(ctx, offset, limit, filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	today, week, statsTotal, err := GetSyncService().OperationStats(ctx)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	totalPages := 0
	if req.PageSize > 0 && total > 0 {
		totalPages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	response.Success(c, map[string]interface{}{
		"list": converter.ToOperationLogList(list),
		"pagination": map[string]interface{}{
			"page":        req.Page,
			"page_size":   req.PageSize,
			"total":       total,
			"total_pages": totalPages,
		},
		"stats": OperationLogStats{
			Today: today,
			Week:  week,
			Total: statsTotal,
		},
	})
}
