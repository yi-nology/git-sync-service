package model

import (
	"time"

	"gorm.io/gorm"
)

type WebhookRule struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"size:100;not null"`
	RepoKey       string         `json:"repoKey" gorm:"size:255;not null;index"`
	EventType     string         `json:"eventType" gorm:"size:100;default:push"`
	BranchPattern string         `json:"branchPattern" gorm:"size:255"`
	Action        string         `json:"action" gorm:"size:50;default:sync"`
	SyncTaskKeys  string         `json:"syncTaskKeys" gorm:"type:text"`
	MinInterval   int            `json:"minInterval" gorm:"default:60"`
	Enabled       bool           `json:"enabled" gorm:"default:true;index"`
	Description   string         `json:"description" gorm:"type:text"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (WebhookRule) TableName() string {
	return "webhook_rules"
}

type CreateRuleRequest struct {
	Name          string `json:"name"`
	RepoKey       string `json:"repoKey"`
	EventType     string `json:"eventType"`
	BranchPattern string `json:"branchPattern"`
	Action        string `json:"action"`
	SyncTaskKeys  string `json:"syncTaskKeys"`
	MinInterval   int    `json:"minInterval"`
	Enabled       bool   `json:"enabled"`
	Description   string `json:"description"`
}

type UpdateRuleRequest struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	EventType     string `json:"eventType"`
	BranchPattern string `json:"branchPattern"`
	Action        string `json:"action"`
	SyncTaskKeys  string `json:"syncTaskKeys"`
	MinInterval   int    `json:"minInterval"`
	Enabled       bool   `json:"enabled"`
	Description   string `json:"description"`
}
