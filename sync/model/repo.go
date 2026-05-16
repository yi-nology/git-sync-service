package model

import (
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Key           string         `json:"key" gorm:"uniqueIndex;size:255;not null"`
	Name          string         `json:"name" gorm:"size:255;not null"`
	Platform      string         `json:"platform" gorm:"size:50;not null"`
	PlatformOwner string         `json:"platformOwner" gorm:"size:200;not null"`
	PlatformRepo  string         `json:"platformRepo" gorm:"size:200;not null"`
	CloneURL      string         `json:"cloneUrl" gorm:"size:500"`
	SSHURL        string         `json:"sshUrl" gorm:"size:500"`
	DefaultBranch string         `json:"defaultBranch" gorm:"size:100;default:main"`
	AccessToken   string         `json:"-" gorm:"type:text"`
	WebhookSecret string         `json:"-" gorm:"size:255"`
	WebhookID     int64          `json:"webhookId"`
	Status        string         `json:"status" gorm:"size:20;default:active"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Repo) TableName() string {
	return "repos"
}
