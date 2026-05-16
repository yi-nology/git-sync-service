package dao

import (
	"errors"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

type WebhookRuleDAO struct {
	db *gorm.DB
}

func NewWebhookRuleDAO(db *gorm.DB) *WebhookRuleDAO {
	return &WebhookRuleDAO{db: db}
}

func (d *WebhookRuleDAO) FindByRepoKey(repoKey string) ([]*model.WebhookRule, error) {
	var rules []*model.WebhookRule
	err := d.db.Where("repo_key = ?", repoKey).Find(&rules).Error
	return rules, err
}

func (d *WebhookRuleDAO) FindByID(id uint) (*model.WebhookRule, error) {
	var rule model.WebhookRule
	err := d.db.First(&rule, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rule, err
}

func (d *WebhookRuleDAO) Create(rule *model.WebhookRule) error {
	return d.db.Create(rule).Error
}

func (d *WebhookRuleDAO) Update(rule *model.WebhookRule) error {
	return d.db.Save(rule).Error
}

func (d *WebhookRuleDAO) Delete(id uint) error {
	return d.db.Delete(&model.WebhookRule{}, id).Error
}
