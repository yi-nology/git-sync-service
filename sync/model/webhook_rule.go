package model

import (
	"time"

	"gorm.io/gorm"
)

type WebhookRule struct {
	ID            uint               `json:"id" gorm:"primaryKey"`
	Name          string             `json:"name" gorm:"size:100;not null"`
	RepoKey       string             `json:"repoKey" gorm:"size:255;not null;index"`
	EventType     string             `json:"eventType" gorm:"size:100;default:push"`
	BranchPattern string             `json:"branchPattern" gorm:"size:255"`
	Action        string             `json:"action" gorm:"size:50;default:sync"`
	MinInterval   int                `json:"minInterval" gorm:"default:60"`
	Enabled       bool               `json:"enabled" gorm:"default:true;index"`
	Description   string             `json:"description" gorm:"type:text"`
	Tasks         []WebhookRuleTask  `json:"tasks,omitempty" gorm:"foreignKey:RuleID"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt     `json:"-" gorm:"index"`
}

func (WebhookRule) TableName() string {
	return TableWebhookRules
}

type WebhookRuleTask struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	RuleID     uint      `json:"ruleId" gorm:"not null;index;uniqueIndex:idx_rule_task"`
	TaskKey    string    `json:"taskKey" gorm:"size:36;not null;uniqueIndex:idx_rule_task"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (WebhookRuleTask) TableName() string {
	return TableWebhookRuleTasks
}

type CreateRuleRequest struct {
	Name          string   `json:"name"`
	RepoKey       string   `json:"repoKey"`
	EventType     string   `json:"eventType"`
	BranchPattern string   `json:"branchPattern"`
	Action        string   `json:"action"`
	TaskKeys      []string `json:"taskKeys"`
	MinInterval   int      `json:"minInterval"`
	Enabled       bool     `json:"enabled"`
	Description   string   `json:"description"`
}

type UpdateRuleRequest struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	EventType     string   `json:"eventType"`
	BranchPattern string   `json:"branchPattern"`
	Action        string   `json:"action"`
	TaskKeys      []string `json:"taskKeys"`
	MinInterval   int      `json:"minInterval"`
	Enabled       bool     `json:"enabled"`
	Description   string   `json:"description"`
}

// GetTaskKeys returns the task keys from the Tasks relationship
func (r *WebhookRule) GetTaskKeys() []string {
	keys := make([]string, len(r.Tasks))
	for i, t := range r.Tasks {
		keys[i] = t.TaskKey
	}
	return keys
}

// SetTaskKeys updates the Tasks relationship
func (r *WebhookRule) SetTaskKeys(keys []string) {
	tasks := make([]WebhookRuleTask, len(keys))
	for i, key := range keys {
		tasks[i] = WebhookRuleTask{
			RuleID:  r.ID,
			TaskKey: key,
		}
	}
	r.Tasks = tasks
}
