package model

import "time"

// SyncRunStep records a single step within a sync run.
type SyncRunStep struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	RunID      uint       `json:"runId" gorm:"not null;index"`
	StepName   string     `json:"stepName" gorm:"size:50;not null"`
	Status     string     `json:"status" gorm:"size:20;not null"`
	StartTime  time.Time  `json:"startTime"`
	EndTime    *time.Time `json:"endTime"`
	DurationMs int64      `json:"durationMs"`
	ErrorMsg   string     `json:"errorMsg" gorm:"type:text"`
	ErrorType  string     `json:"errorType" gorm:"size:30"`
	Output     string     `json:"output" gorm:"type:text"`
	RetryCount int        `json:"retryCount"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"index"`
}

func (SyncRunStep) TableName() string {
	return TableSyncRunSteps
}
