package model

import "time"

type SyncRun struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	TaskKey       string     `json:"taskKey" gorm:"size:36;not null;index"`
	TriggerSource string     `json:"triggerSource" gorm:"size:20;not null"`
	Status        string     `json:"status" gorm:"size:20;not null;index"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	CommitRange   string     `json:"commitRange" gorm:"size:255"`
	Details       string     `json:"details" gorm:"type:text"`
	ErrorMessage  string     `json:"errorMessage" gorm:"type:text"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func (SyncRun) TableName() string {
	return "sync_runs"
}
