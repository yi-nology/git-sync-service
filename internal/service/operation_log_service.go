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
func (s *OperationLogService) List(ctx context.Context, offset, limit int, filter *dao.OperationLogFilter) ([]*model.OperationLog, int64, error) {
	page := dao.DefaultPagination(offset, limit)
	return s.opLogDAO.List(page, filter)
}

// Stats 返回今日、近 7 天（本周）、总操作数。单次 SQL 查询。
func (s *OperationLogService) Stats(ctx context.Context) (today, week, total int64, err error) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	r, err := s.opLogDAO.StatsOnce(startOfToday)
	if err != nil {
		return 0, 0, 0, err
	}
	return r.Today, r.Week, r.Total, nil
}
