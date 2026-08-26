package model

import "time"

type WebhookEvent struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	EventID      string     `json:"eventId" gorm:"uniqueIndex;size:100;not null"`
	RepoKey      string     `json:"repoKey" gorm:"size:255;index:idx_webhook_events_repo_key_created,priority:1"`
	EventType    string     `json:"eventType" gorm:"size:50;not null"`
	Source       string     `json:"source" gorm:"size:20;not null"`
	ActorName    string     `json:"actorName" gorm:"size:200"`
	Branch       string     `json:"branch" gorm:"size:255"`
	CommitSHA    string     `json:"commitSha" gorm:"size:40"`
	Payload      []byte     `json:"payload" gorm:"type:json"`
	Status       string     `json:"status" gorm:"size:20;default:received;index"`
	ErrorMessage string     `json:"errorMessage" gorm:"type:text"`
	ProcessedAt  *time.Time `json:"processedAt"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"index:idx_webhook_events_created_at;index:idx_webhook_events_repo_key_created,priority:2"`
}

func (WebhookEvent) TableName() string {
	return TableWebhookEvents
}
