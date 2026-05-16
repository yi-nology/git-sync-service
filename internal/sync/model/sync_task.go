package model

import (
	"time"

	"gorm.io/gorm"
)

type SyncTask struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Key           string         `json:"key" gorm:"uniqueIndex;size:36;not null"`
	Name          string         `json:"name" gorm:"size:100;not null"`
	SourceRepoKey string         `json:"sourceRepoKey" gorm:"size:255;not null;index"`
	SourceBranch  string         `json:"sourceBranch" gorm:"size:255;not null"`
	TargetRepoKey string         `json:"targetRepoKey" gorm:"size:255;not null;index"`
	TargetBranch  string         `json:"targetBranch" gorm:"size:255;not null"`
	SyncMode      string         `json:"syncMode" gorm:"size:20;default:single"`
	Cron          string         `json:"cron" gorm:"size:100"`
	WebhookToken  string         `json:"webhookToken" gorm:"uniqueIndex;size:36"`
	Enabled       bool           `json:"enabled" gorm:"default:true;index"`
	GitTags       bool           `json:"gitTags" gorm:"default:false"`
	GitForce      bool           `json:"gitForce" gorm:"default:false"`
	GitPrune      bool           `json:"gitPrune" gorm:"default:false"`
	GitNoVerify   bool           `json:"gitNoVerify" gorm:"default:false"`
	PushOptions   string         `json:"pushOptions" gorm:"size:500"`
	LastRunAt     *time.Time     `json:"lastRunAt"`
	LastStatus    string         `json:"lastStatus" gorm:"size:20"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (SyncTask) TableName() string {
	return "sync_tasks"
}

type CreateTaskRequest struct {
	Name          string `json:"name"`
	SourceRepoKey string `json:"sourceRepoKey"`
	SourceBranch  string `json:"sourceBranch"`
	TargetRepoKey string `json:"targetRepoKey"`
	TargetBranch  string `json:"targetBranch"`
	SyncMode      string `json:"syncMode"`
	Cron          string `json:"cron"`
	GitTags       bool   `json:"gitTags"`
	GitForce      bool   `json:"gitForce"`
	GitPrune      bool   `json:"gitPrune"`
	GitNoVerify   bool   `json:"gitNoVerify"`
	PushOptions   string `json:"pushOptions"`
}

type UpdateTaskRequest struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	SourceBranch  string `json:"sourceBranch"`
	TargetBranch  string `json:"targetBranch"`
	SyncMode      string `json:"syncMode"`
	Cron          string `json:"cron"`
	Enabled       bool   `json:"enabled"`
	GitTags       bool   `json:"gitTags"`
	GitForce      bool   `json:"gitForce"`
	GitPrune      bool   `json:"gitPrune"`
	GitNoVerify   bool   `json:"gitNoVerify"`
	PushOptions   string `json:"pushOptions"`
}
