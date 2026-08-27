package dao

import (
	"github.com/yi-nology/git-platform-sdk/pkg/credential"

	errors "github.com/cockroachdb/errors"
	"github.com/yi-nology/git-sync-service/sync/model"
	"gorm.io/gorm"
)

// PlatformDAO 平台数据访问对象。
// AccessToken 与 repo_dao 一致采用加密存储:写入加密、读取解密;
// 读取解密失败视为存量明文,原样返回(下次写入时自动转为密文)。
type PlatformDAO struct {
	db *gorm.DB
	cm *credential.CryptoManager
}

// NewPlatformDAO 创建 PlatformDAO
func NewPlatformDAO(db *gorm.DB) (*PlatformDAO, error) {
	cm, err := credential.NewCryptoManager()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CryptoManager")
	}
	return &PlatformDAO{db: db, cm: cm}, nil
}

// decrypt 就地解密平台 AccessToken(解密失败保留原文,兼容存量明文)
func (d *PlatformDAO) decrypt(p *model.Platform) {
	if p != nil && p.AccessToken != "" {
		if decrypted, err := d.cm.Decrypt(p.AccessToken); err == nil {
			p.AccessToken = decrypted
		}
	}
}

// decryptAll 就地解密平台列表
func (d *PlatformDAO) decryptAll(platforms []*model.Platform) {
	for _, p := range platforms {
		d.decrypt(p)
	}
}

// encrypt 就地加密平台 AccessToken(空 token 跳过)
func (d *PlatformDAO) encrypt(p *model.Platform) error {
	if p.AccessToken == "" {
		return nil
	}
	encrypted, err := d.cm.Encrypt(p.AccessToken)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt access token")
	}
	p.AccessToken = encrypted
	return nil
}

// Create 创建平台
func (d *PlatformDAO) Create(platform *model.Platform) error {
	if err := d.encrypt(platform); err != nil {
		return err
	}
	return d.db.Create(platform).Error
}

// FindByKey 根据 Key 查找平台(not found 返回 nil,nil,与 RepoDAO 行为一致)。
func (d *PlatformDAO) FindByKey(key string) (*model.Platform, error) {
	var platform model.Platform
	err := d.db.Where("key = ?", key).First(&platform).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.decrypt(&platform)
	return &platform, nil
}

// FindByID 根据 ID 查找平台(not found 返回 nil,nil,与 RepoDAO 行为一致)。
func (d *PlatformDAO) FindByID(id uint) (*model.Platform, error) {
	var platform model.Platform
	err := d.db.Where("id = ?", id).First(&platform).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.decrypt(&platform)
	return &platform, nil
}

// FindAll 查找所有平台
func (d *PlatformDAO) FindAll() ([]*model.Platform, error) {
	var platforms []*model.Platform
	err := d.db.Order("is_default DESC, name ASC").Find(&platforms).Error
	if err != nil {
		return nil, err
	}
	d.decryptAll(platforms)
	return platforms, nil
}

// FindByType 根据类型查找平台
func (d *PlatformDAO) FindByType(platformType string) ([]*model.Platform, error) {
	var platforms []*model.Platform
	err := d.db.Where("type = ?", platformType).Find(&platforms).Error
	if err != nil {
		return nil, err
	}
	d.decryptAll(platforms)
	return platforms, nil
}

// FindDefault 查找默认平台(not found 返回 nil,nil)。
func (d *PlatformDAO) FindDefault() (*model.Platform, error) {
	var platform model.Platform
	err := d.db.Where("is_default = ?", true).First(&platform).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.decrypt(&platform)
	return &platform, nil
}

// Update 更新平台
func (d *PlatformDAO) Update(platform *model.Platform) error {
	if err := d.encrypt(platform); err != nil {
		return err
	}
	return d.db.Save(platform).Error
}

// UpdateFields 更新平台指定字段
// 注意:字段值原样写入;如需更新 access_token 请走 Update 以确保加密。
func (d *PlatformDAO) UpdateFields(key string, fields map[string]interface{}) error {
	return d.db.Model(&model.Platform{}).Where("key = ?", key).Updates(fields).Error
}

// Delete 删除平台（软删除）
func (d *PlatformDAO) Delete(key string) error {
	return d.db.Where("key = ?", key).Delete(&model.Platform{}).Error
}

// SetDefault 设置默认平台(事务保证原子性:先清后设,不会出现无默认的窗口)
func (d *PlatformDAO) SetDefault(key string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Platform{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&model.Platform{}).Where("key = ?", key).Update("is_default", true).Error
	})
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
