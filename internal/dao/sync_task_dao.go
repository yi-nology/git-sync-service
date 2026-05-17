package dao

import (
	"errors"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

type SyncTaskDAO struct {
	db *gorm.DB
}

func NewSyncTaskDAO(db *gorm.DB) *SyncTaskDAO {
	return &SyncTaskDAO{db: db}
}

func (d *SyncTaskDAO) FindByRepoKey(repoKey string, page Pagination) ([]*model.SyncTask, int64, error) {
	var tasks []*model.SyncTask
	var total int64
	query := d.db.Where("source_repo_key = ? OR target_repo_key = ?", repoKey, repoKey)
	query.Model(&model.SyncTask{}).Count(&total)
	err := query.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(&tasks).Error
	return tasks, total, err
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

func (d *SyncTaskDAO) FindAll(page Pagination) ([]*model.SyncTask, int64, error) {
	var tasks []*model.SyncTask
	var total int64
	d.db.Model(&model.SyncTask{}).Count(&total)
	err := d.db.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(&tasks).Error
	return tasks, total, err
}
