package dao

import (
	"errors"
	"time"

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

func (d *WebhookEventDAO) FindByRepoKey(repoKey string, page Pagination) ([]*model.WebhookEvent, int64, error) {
	var events []*model.WebhookEvent
	query := d.db.Where("repo_key = ?", repoKey)
	total, err := Paginate(query, page, &events)
	return events, total, err
}

func (d *WebhookEventDAO) Create(event *model.WebhookEvent) error {
	return d.db.Create(event).Error
}

func (d *WebhookEventDAO) Update(event *model.WebhookEvent) error {
	return d.db.Save(event).Error
}

func (d *WebhookEventDAO) FindByID(id uint) (*model.WebhookEvent, error) {
	return FindByID[model.WebhookEvent](d.db, id)
}

func (d *WebhookEventDAO) FindRecent(repoKey string, page Pagination) ([]*model.WebhookEvent, error) {
	var events []*model.WebhookEvent
	query := d.db.Offset(page.Offset).Limit(page.Limit).Order("id DESC")
	if repoKey != "" {
		query = query.Where("repo_key = ?", repoKey)
	}
	err := query.Find(&events).Error
	return events, err
}

func (d *WebhookEventDAO) CleanupOlderThan(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result := d.db.Where("created_at < ?", cutoff).Delete(&model.WebhookEvent{})
	return result.RowsAffected, result.Error
}

func (d *WebhookEventDAO) CountByRepoKey(repoKey string) (int64, error) {
	var count int64
	err := d.db.Model(&model.WebhookEvent{}).Where("repo_key = ?", repoKey).Count(&count).Error
	return count, err
}
