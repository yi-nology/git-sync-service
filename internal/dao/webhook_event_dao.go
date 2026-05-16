package dao

import (
	"errors"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

type WebhookEventDAO struct {
	db *gorm.DB
}

func NewWebhookEventDAO(db *gorm.DB) *WebhookEventDAO {
	return &WebhookEventDAO{db: db}
}

func (d *WebhookEventDAO) FindByEventID(eventID string) (*model.WebhookEvent, error) {
	var event model.WebhookEvent
	err := d.db.Where("event_id = ?", eventID).First(&event).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &event, err
}

func (d *WebhookEventDAO) FindByRepoKey(repoKey string, limit int) ([]*model.WebhookEvent, error) {
	var events []*model.WebhookEvent
	err := d.db.Where("repo_key = ?", repoKey).Order("id DESC").Limit(limit).Find(&events).Error
	return events, err
}

func (d *WebhookEventDAO) Create(event *model.WebhookEvent) error {
	return d.db.Create(event).Error
}

func (d *WebhookEventDAO) Update(event *model.WebhookEvent) error {
	return d.db.Save(event).Error
}

func (d *WebhookEventDAO) FindByID(id uint) (*model.WebhookEvent, error) {
	var event model.WebhookEvent
	err := d.db.First(&event, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &event, err
}

func (d *WebhookEventDAO) FindRecent(repoKey string, limit int) ([]*model.WebhookEvent, error) {
	var events []*model.WebhookEvent
	query := d.db.Order("id DESC").Limit(limit)
	if repoKey != "" {
		query = query.Where("repo_key = ?", repoKey)
	}
	err := query.Find(&events).Error
	return events, err
}
