package dao

import (
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

// PlatformDAO 平台数据访问对象
type PlatformDAO struct {
	db *gorm.DB
}

// NewPlatformDAO 创建 PlatformDAO
func NewPlatformDAO(db *gorm.DB) *PlatformDAO {
	return &PlatformDAO{db: db}
}

// Create 创建平台
func (d *PlatformDAO) Create(platform *model.Platform) error {
	return d.db.Create(platform).Error
}

// FindByKey 根据 Key 查找平台
func (d *PlatformDAO) FindByKey(key string) (*model.Platform, error) {
	var platform model.Platform
	err := d.db.Where("key = ?", key).First(&platform).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

// FindByID 根据 ID 查找平台
func (d *PlatformDAO) FindByID(id uint) (*model.Platform, error) {
	var platform model.Platform
	err := d.db.Where("id = ?", id).First(&platform).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

// FindAll 查找所有平台
func (d *PlatformDAO) FindAll() ([]*model.Platform, error) {
	var platforms []*model.Platform
	err := d.db.Order("is_default DESC, name ASC").Find(&platforms).Error
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

// FindByType 根据类型查找平台
func (d *PlatformDAO) FindByType(platformType string) ([]*model.Platform, error) {
	var platforms []*model.Platform
	err := d.db.Where("type = ?", platformType).Find(&platforms).Error
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

// FindDefault 查找默认平台
func (d *PlatformDAO) FindDefault() (*model.Platform, error) {
	var platform model.Platform
	err := d.db.Where("is_default = ?", true).First(&platform).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

// Update 更新平台
func (d *PlatformDAO) Update(platform *model.Platform) error {
	return d.db.Save(platform).Error
}

// UpdateFields 更新平台指定字段
func (d *PlatformDAO) UpdateFields(key string, fields map[string]interface{}) error {
	return d.db.Model(&model.Platform{}).Where("key = ?", key).Updates(fields).Error
}

// Delete 删除平台（软删除）
func (d *PlatformDAO) Delete(key string) error {
	return d.db.Where("key = ?", key).Delete(&model.Platform{}).Error
}

// SetDefault 设置默认平台
func (d *PlatformDAO) SetDefault(key string) error {
	// 先取消所有默认
	err := d.db.Model(&model.Platform{}).Where("is_default = ?", true).Update("is_default", false).Error
	if err != nil {
		return err
	}
	// 设置新的默认
	return d.db.Model(&model.Platform{}).Where("key = ?", key).Update("is_default", true).Error
}

// UpdateRepoCount 更新平台关联的仓库数量
func (d *PlatformDAO) UpdateRepoCount(platformID uint) error {
	var count int64
	err := d.db.Model(&model.Repo{}).Where("platform_id = ? AND deleted_at IS NULL", platformID).Count(&count).Error
	if err != nil {
		return err
	}
	return d.db.Model(&model.Platform{}).Where("id = ?", platformID).Update("repo_count", count).Error
}

// UpdateStatus 更新平台状态
func (d *PlatformDAO) UpdateStatus(key, status, testResult string) error {
	return d.db.Model(&model.Platform{}).Where("key = ?", key).Updates(map[string]interface{}{
		"status":          status,
		"last_test_at":    gorm.Expr("CURRENT_TIMESTAMP"),
		"last_test_result": testResult,
	}).Error
}

// Count 统计平台数量
func (d *PlatformDAO) Count() (int64, error) {
	var count int64
	err := d.db.Model(&model.Platform{}).Count(&count).Error
	return count, err
}
