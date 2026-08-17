package git_sync

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	operation_log "github.com/yi-nology/git-sync-service/biz/model/operation_log"
	"github.com/yi-nology/git-sync-service/internal/converter"
	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

// ListOperationLogs GET /api/v1/logs/operations
// 返回 {list, pagination, stats}（stats.total 为全量总数，pagination.total 为过滤后的结果数）。
func ListOperationLogs(ctx context.Context, c *app.RequestContext) {
	var req operation_log.ListOperationLogsReq
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

	offset, limit := converter.PageToOffset(req.Page, req.PageSize)

	filter := dao.OperationLogFilter{
		Search:    req.Search,
		Action:    req.Action,
		Actor:     req.User,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	list, total, err := GetSyncService().ListOperations(ctx, offset, limit, &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	today, week, statsTotal, err := GetSyncService().OperationStats(ctx)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	totalPages := int32(0)
	if req.PageSize > 0 && total > 0 {
		totalPages = converter.SafeInt64ToInt32((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	response.Success(c, &operation_log.ListOperationLogsResp{
		List: converter.ToOperationLogList(list),
		Pagination: &operation_log.OperationLogPagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		Stats: &operation_log.OperationLogStats{
			Today: today,
			Week:  week,
			Total: statsTotal,
		},
	})
}
