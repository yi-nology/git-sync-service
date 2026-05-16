package sync

import (
	"errors"

	"github.com/yi-nology/git-sync-service/internal/sync/model"
	"gorm.io/gorm"
)

type RepoDAO struct {
	db *gorm.DB
}

func NewRepoDAO(db *gorm.DB) *RepoDAO {
	return &RepoDAO{db: db}
}

func (d *RepoDAO) FindAll() ([]*model.Repo, error) {
	var repos []*model.Repo
	err := d.db.Find(&repos).Error
	return repos, err
}

func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
	var repo model.Repo
	err := d.db.Where("`key` = ?", key).First(&repo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &repo, err
}

func (d *RepoDAO) Create(repo *model.Repo) error {
	return d.db.Create(repo).Error
}

func (d *RepoDAO) Update(repo *model.Repo) error {
	return d.db.Save(repo).Error
}

func (d *RepoDAO) Delete(key string) error {
	return d.db.Where("`key` = ?", key).Delete(&model.Repo{}).Error
}

type SyncTaskDAO struct {
	db *gorm.DB
}

func NewSyncTaskDAO(db *gorm.DB) *SyncTaskDAO {
	return &SyncTaskDAO{db: db}
}

func (d *SyncTaskDAO) FindByRepoKey(repoKey string) ([]*model.SyncTask, error) {
	var tasks []*model.SyncTask
	err := d.db.Where("source_repo_key = ? OR target_repo_key = ?", repoKey, repoKey).Find(&tasks).Error
	return tasks, err
}

func (d *SyncTaskDAO) FindByKey(key string) (*model.SyncTask, error) {
	var task model.SyncTask
	err := d.db.Where("`key` = ?", key).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

func (d *SyncTaskDAO) FindAllEnabled() ([]*model.SyncTask, error) {
	var tasks []*model.SyncTask
	err := d.db.Where("enabled = ? AND cron IS NOT NULL AND cron != ?", true, "").Find(&tasks).Error
	return tasks, err
}

func (d *SyncTaskDAO) Create(task *model.SyncTask) error {
	return d.db.Create(task).Error
}

func (d *SyncTaskDAO) Update(task *model.SyncTask) error {
	return d.db.Save(task).Error
}

func (d *SyncTaskDAO) Delete(key string) error {
	return d.db.Where("`key` = ?", key).Delete(&model.SyncTask{}).Error
}

type SyncRunDAO struct {
	db *gorm.DB
}

func NewSyncRunDAO(db *gorm.DB) *SyncRunDAO {
	return &SyncRunDAO{db: db}
}

func (d *SyncRunDAO) FindByTaskKey(taskKey string, limit int) ([]*model.SyncRun, error) {
	var runs []*model.SyncRun
	err := d.db.Where("task_key = ?", taskKey).Order("id DESC").Limit(limit).Find(&runs).Error
	return runs, err
}

func (d *SyncRunDAO) Create(run *model.SyncRun) error {
	return d.db.Create(run).Error
}

func (d *SyncRunDAO) Update(run *model.SyncRun) error {
	return d.db.Save(run).Error
}

func (d *SyncRunDAO) Delete(id uint) error {
	return d.db.Delete(&model.SyncRun{}, id).Error
}

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
