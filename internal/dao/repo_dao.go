package dao

import (
	"errors"

	"github.com/yi-nology/git-platform-sdk/pkg/credential"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

type Pagination struct {
	Offset int
	Limit  int
}

func DefaultPagination(offset, limit int) Pagination {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return Pagination{Offset: offset, Limit: limit}
}

type RepoDAO struct {
	db *gorm.DB
}

func NewRepoDAO(db *gorm.DB) *RepoDAO {
	return &RepoDAO{db: db}
}

func (d *RepoDAO) FindAll(page Pagination) ([]*model.Repo, int64, error) {
	var repos []*model.Repo
	var total int64
	d.db.Model(&model.Repo{}).Count(&total)
	err := d.db.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(&repos).Error
	return repos, total, err
}

func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
	var repo model.Repo
	err := d.db.Where("`key` = ?", key).First(&repo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	decrypted, err := credential.DecryptGCM(repo.AccessToken)
	if err != nil {
		return &repo, nil
	}
	repo.AccessToken = decrypted
	return &repo, nil
}

func (d *RepoDAO) Create(repo *model.Repo) error {
	encrypted, err := credential.EncryptGCM(repo.AccessToken)
	if err != nil {
		return err
	}
	repo.AccessToken = encrypted
	return d.db.Create(repo).Error
}

func (d *RepoDAO) Update(repo *model.Repo) error {
	if repo.AccessToken != "" {
		encrypted, err := credential.EncryptGCM(repo.AccessToken)
		if err != nil {
			return err
		}
		repo.AccessToken = encrypted
	}
	return d.db.Save(repo).Error
}

func (d *RepoDAO) Delete(key string) error {
	return d.db.Where("`key` = ?", key).Delete(&model.Repo{}).Error
}
