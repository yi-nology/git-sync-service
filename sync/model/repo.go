package model

import (
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Key           string         `json:"key" gorm:"uniqueIndex;size:255;not null"`
	Name          string         `json:"name" gorm:"size:255;not null"`
	PlatformID    uint           `json:"platformId" gorm:"index"`                // 关联平台 ID
	Platform      string         `json:"platform" gorm:"size:50;not null;default:'unknown'"`
	PlatformOwner string         `json:"platformOwner" gorm:"size:200;not null;default:''"`
	PlatformRepo  string         `json:"platformRepo" gorm:"size:200;not null;default:''"`
	CloneURL      string         `json:"cloneUrl" gorm:"size:500"`
	SSHURL        string         `json:"sshUrl" gorm:"size:500"`
	DefaultBranch string         `json:"defaultBranch" gorm:"size:100;default:main"`
	AccessToken   string         `json:"-" gorm:"type:text"`                     // 保留用于兼容，优先使用 Platform 的 Token
	WebhookSecret string         `json:"-" gorm:"size:255"`
	WebhookID     int64          `json:"webhookId"`
	Status        string         `json:"status" gorm:"size:20;default:active"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	PlatformRef *Platform `json:"platformRef,omitempty" gorm:"foreignKey:PlatformID"`
}

func (Repo) TableName() string {
	return TableRepos
}
