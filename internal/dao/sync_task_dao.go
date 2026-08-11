package dao

import (
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
	query := d.db.Where("source_repo_key = ? OR target_repo_key = ?", repoKey, repoKey)
	total, err := Paginate(query, page, &tasks)
	return tasks, total, err
}

func (d *SyncTaskDAO) FindByKey(key string) (*model.SyncTask, error) {
	return FindByKey[model.SyncTask](d.db, key)
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
	total, err := Paginate(d.db, page, &tasks)
	return tasks, total, err
}
