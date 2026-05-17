package dao

import (
	"time"

	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

type SyncRunDAO struct {
	db *gorm.DB
}

func NewSyncRunDAO(db *gorm.DB) *SyncRunDAO {
	return &SyncRunDAO{db: db}
}

func (d *SyncRunDAO) FindByTaskKey(taskKey string, page Pagination) ([]*model.SyncRun, int64, error) {
	var runs []*model.SyncRun
	var total int64
	query := d.db.Where("task_key = ?", taskKey)
	query.Model(&model.SyncRun{}).Count(&total)
	err := query.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(&runs).Error
	return runs, total, err
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

func (d *SyncRunDAO) FindRecent(page Pagination) ([]*model.SyncRun, error) {
	var runs []*model.SyncRun
	err := d.db.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(&runs).Error
	return runs, err
}

func (d *SyncRunDAO) CleanupOlderThan(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result := d.db.Where("created_at < ?", cutoff).Delete(&model.SyncRun{})
	return result.RowsAffected, result.Error
}

func (d *SyncRunDAO) CountByTaskKey(taskKey string) (int64, error) {
	var count int64
	err := d.db.Model(&model.SyncRun{}).Where("task_key = ?", taskKey).Count(&count).Error
	return count, err
}
