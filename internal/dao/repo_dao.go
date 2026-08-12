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

// FindByCloneURL 根据 Clone URL 查找仓库
func (d *RepoDAO) FindByCloneURL(cloneURL string) (*model.Repo, error) {
	var repo model.Repo
	err := d.db.Where("clone_url = ?", cloneURL).First(&repo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// FindByPlatformID 根据平台 ID 查找仓库
func (d *RepoDAO) FindByPlatformID(platformID uint) ([]*model.Repo, error) {
	var repos []*model.Repo
	err := d.db.Where("platform_id = ?", platformID).Find(&repos).Error
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// RepoFilter contains optional filters for listing repositories.
type RepoFilter struct {
	Search   string // search by name or clone_url (LIKE)
	Platform string // exact match on platform
	Status   string // exact match on status
	SortBy   string // column to sort by (default: created_at)
	OrderBy  string // sort direction: asc or desc (default: desc)
}

// ListWithFilter returns a filtered, sorted, paginated list of repos and the total count.
func (d *RepoDAO) ListWithFilter(page Pagination, filter *RepoFilter) ([]*model.Repo, int64, error) {
	var repos []*model.Repo
	var total int64

	query := d.db.Model(&model.Repo{})

	// Apply search filter (LIKE on name or clone_url)
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("(name LIKE ? OR clone_url LIKE ?)", like, like)
	}

	// Apply exact-match filters
	if filter.Platform != "" {
		query = query.Where("platform = ?", filter.Platform)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Determine sort column and order
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	order := "DESC"
	if filter.OrderBy == "asc" {
		order = "ASC"
	}

	// Apply pagination and sorting
	if err := query.Offset(page.Offset).Limit(page.Limit).Order(sortBy + " " + order).Find(&repos).Error; err != nil {
		return nil, 0, err
	}

	return repos, total, nil
}
