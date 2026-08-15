package model

import (
	"time"

	"gorm.io/gorm"
)

type SyncTask struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Key           string         `json:"key" gorm:"uniqueIndex;size:36;not null;default:''"`
	Name          string         `json:"name" gorm:"size:100;not null;default:''"`
	SourceRepoKey string         `json:"source_repo_key" gorm:"size:255;not null;index;default:''"`
	SourceBranch  string         `json:"source_branch" gorm:"size:255;not null;default:''"`
	TargetRepoKey string         `json:"target_repo_key" gorm:"size:255;not null;index;default:''"`
	TargetBranch  string         `json:"target_branch" gorm:"size:255;not null;default:''"`
	SyncMode      string         `json:"sync_mode" gorm:"size:20;default:single"`
	Cron          string         `json:"cron" gorm:"size:100"`
	WebhookToken  string         `json:"webhook_token" gorm:"uniqueIndex;size:36"`
	Enabled       bool           `json:"enabled" gorm:"default:true;index"`
	GitTags       bool           `json:"git_tags" gorm:"default:false"`
	GitForce      bool           `json:"git_force" gorm:"default:false"`
	GitPrune      bool           `json:"git_prune" gorm:"default:false"`
	LastRunAt     *time.Time     `json:"last_run_at"`
	LastStatus    string         `json:"last_status" gorm:"size:20"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (SyncTask) TableName() string {
	return TableSyncTasks
}

type CreateTaskRequest struct {
	Name          string `json:"name"`
	SourceRepoKey string `json:"source_repo_key"`
	SourceBranch  string `json:"source_branch"`
	TargetRepoKey string `json:"target_repo_key"`
	TargetBranch  string `json:"target_branch"`
	SyncMode      string `json:"sync_mode"`
	Cron          string `json:"cron"`
	GitTags       bool   `json:"git_tags"`
	GitForce      bool   `json:"git_force"`
	GitPrune      bool   `json:"git_prune"`
}

type UpdateTaskRequest struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	SyncMode     string `json:"sync_mode"`
	Cron         string `json:"cron"`
	Enabled      bool   `json:"enabled"`
	GitTags      bool   `json:"git_tags"`
	GitForce     bool   `json:"git_force"`
	GitPrune     bool   `json:"git_prune"`
}
