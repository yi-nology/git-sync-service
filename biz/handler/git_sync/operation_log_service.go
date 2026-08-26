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

	// ListOperations 与 OperationStats 互不依赖,并行执行缩短响应时间。
	type statsResult struct {
		today, week, total int64
		err                error
	}
	svc := GetSyncService()
	statsCh := make(chan statsResult, 1)
	go func() {
		t, w, tot, e := svc.OperationStats(ctx)
		statsCh <- statsResult{t, w, tot, e}
	}()

	list, total, err := svc.ListOperations(ctx, offset, limit, &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	sr := <-statsCh
	if sr.err != nil {
		response.InternalError(c, sr.err.Error())
		return
	}
	today, week, statsTotal := sr.today, sr.week, sr.total

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
