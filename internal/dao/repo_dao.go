package dao

import (
	"errors"

	"github.com/yi-nology/git-sync-service/sync/model"
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
