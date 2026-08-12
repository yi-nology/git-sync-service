package service

import (
	"context"
	"time"

	"github.com/yi-nology/git-sync-service/internal/dao"
	"github.com/yi-nology/git-sync-service/sync/model"
)

// OperationLogService 处理审计日志相关操作。
type OperationLogService struct {
	opLogDAO *dao.OperationLogDAO
}

// NewOperationLogService 创建新的 OperationLogService 实例。
func NewOperationLogService(opLogDAO *dao.OperationLogDAO) *OperationLogService {
	return &OperationLogService{opLogDAO: opLogDAO}
}

// Record 记录一条审计日志（best-effort，由调用方决定如何处理错误）。
func (s *OperationLogService) Record(ctx context.Context, entry *model.OperationLog) error {
	if entry.Status == "" {
		entry.Status = model.StatusSuccess
	}
	return s.opLogDAO.Create(entry)
}

// List 按过滤条件分页返回审计日志。
func (s *OperationLogService) List(ctx context.Context, offset, limit int, filter dao.OperationLogFilter) ([]*model.OperationLog, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return s.opLogDAO.List(page, filter)
}

// Stats 返回今日、近 7 天（本周）、总操作数。
func (s *OperationLogService) Stats(ctx context.Context) (today, week, total int64, err error) {
	total, err = s.opLogDAO.CountAll()
	if err != nil {
		return 0, 0, 0, err
	}
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	today, err = s.opLogDAO.CountSince(startOfToday)
	if err != nil {
		return 0, 0, total, err
	}
	// “本周”按滚动 7 天统计（含今日），避免跨地区周一/周日起点歧义。
	week, err = s.opLogDAO.CountSince(startOfToday.AddDate(0, 0, -6))
	if err != nil {
		return today, 0, total, err
	}
	return today, week, total, nil
}
