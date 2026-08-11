package dao

import (
	"errors"
	"fmt"

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
	cm *credential.CryptoManager
}

func NewRepoDAO(db *gorm.DB) (*RepoDAO, error) {
	cm, err := credential.NewCryptoManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create CryptoManager: %w", err)
	}
	return &RepoDAO{db: db, cm: cm}, nil
}

func (d *RepoDAO) FindAll(page Pagination) ([]*model.Repo, int64, error) {
	var repos []*model.Repo
	total, err := Paginate(d.db, page, &repos)
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
	// Decrypt access token
	if repo.AccessToken != "" {
		if decrypted, err := d.cm.Decrypt(repo.AccessToken); err == nil {
			repo.AccessToken = decrypted
		}
	}
	// Decrypt webhook secret
	if repo.WebhookSecret != "" {
		if decrypted, err := d.cm.Decrypt(repo.WebhookSecret); err == nil {
			repo.WebhookSecret = decrypted
		}
	}
	return &repo, nil
}

func (d *RepoDAO) Create(repo *model.Repo) error {
	if repo.AccessToken != "" {
		encrypted, err := d.cm.Encrypt(repo.AccessToken)
		if err != nil {
			return err
		}
		repo.AccessToken = encrypted
	}
	if repo.WebhookSecret != "" {
		encrypted, err := d.cm.Encrypt(repo.WebhookSecret)
		if err != nil {
			return err
		}
		repo.WebhookSecret = encrypted
	}
	return d.db.Create(repo).Error
}

func (d *RepoDAO) Update(repo *model.Repo) error {
	if repo.AccessToken != "" {
		encrypted, err := d.cm.Encrypt(repo.AccessToken)
		if err != nil {
			return err
		}
		repo.AccessToken = encrypted
	}
	if repo.WebhookSecret != "" {
		encrypted, err := d.cm.Encrypt(repo.WebhookSecret)
		if err != nil {
			return err
		}
		repo.WebhookSecret = encrypted
	}
	return d.db.Save(repo).Error
}

func (d *RepoDAO) Delete(key string) error {
	return d.db.Where("`key` = ?", key).Delete(&model.Repo{}).Error
}
