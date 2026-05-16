package dao

import (
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

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
