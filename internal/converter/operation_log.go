package converter

import (
	"github.com/yi-nology/git-sync-service/sync/model"
)

// OperationLogInfo 是审计日志对外的 JSON 结构，字段名对齐前端 OperationLog。
type OperationLogInfo struct {
	ID       uint   `json:"id"`
	Time     string `json:"time"`
	User     string `json:"user"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Details  string `json:"details"`
	IP       string `json:"ip"`
}

// ToOperationLogInfo 将审计日志模型转为对外结构。
func ToOperationLogInfo(l *model.OperationLog) *OperationLogInfo {
	if l == nil {
		return nil
	}
	return &OperationLogInfo{
		ID:       l.ID,
		Time:     l.CreatedAt.Format("2006-01-02 15:04:05"),
		User:     l.Actor,
		Action:   l.Action,
		Resource: l.Resource,
		Details:  l.Detail,
		IP:       l.IP,
	}
}

// ToOperationLogList 批量转换审计日志。
func ToOperationLogList(logs []*model.OperationLog) []*OperationLogInfo {
	result := make([]*OperationLogInfo, 0, len(logs))
	for _, l := range logs {
		result = append(result, ToOperationLogInfo(l))
	}
	return result
}
