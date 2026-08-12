package model

import "time"

// OperationLog 记录用户对系统发起的写操作（审计日志）。
type OperationLog struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Action       string    `json:"action" gorm:"size:32;not null;index"`  // create/update/delete/run/retry/sync
	ResourceType string    `json:"resource_type" gorm:"size:32;not null"` // repo/task/rule/platform/event
	ResourceKey  string    `json:"resource_key" gorm:"size:255;index"`    // 资源标识（key/name/id）
	Resource     string    `json:"resource" gorm:"size:500"`              // 中文摘要，如 "创建仓库 repo-main"
	Actor        string    `json:"actor" gorm:"size:128;index"`           // 操作者（X-User，缺省 admin）
	IP           string    `json:"ip" gorm:"size:64"`
	Status       string    `json:"status" gorm:"size:16;default:success"` // success/failed
	Detail       string    `json:"detail" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"index"`
}

func (OperationLog) TableName() string {
	return TableOperationLogs
}
