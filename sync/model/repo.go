package model

import (
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Key           string         `json:"key" gorm:"uniqueIndex;size:255;not null"`
	Name          string         `json:"name" gorm:"size:255;not null"`
	PlatformID    uint           `json:"platform_id" gorm:"index"`                // 关联平台 ID
	Platform      string         `json:"platform" gorm:"size:50;not null;default:'unknown'"`
	PlatformOwner string         `json:"platform_owner" gorm:"size:200;not null;default:''"`
	PlatformRepo  string         `json:"platform_repo" gorm:"size:200;not null;default:''"`
	CloneURL      string         `json:"clone_url" gorm:"size:500"`
	SSHURL        string         `json:"ssh_url" gorm:"size:500"`
	DefaultBranch string         `json:"default_branch" gorm:"size:100;default:main"`
	AccessToken   string         `json:"-" gorm:"type:text"`                     // 保留用于兼容，优先使用 Platform 的 Token
	WebhookSecret string         `json:"-" gorm:"size:255"`
	WebhookID     int64          `json:"webhook_id"`
	Status        string         `json:"status" gorm:"size:20;default:active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Repo) TableName() string {
	return TableRepos
}
