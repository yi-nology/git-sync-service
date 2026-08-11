package dao

import (
	"errors"

	"gorm.io/gorm"
)

// Paginate 执行分页查询
func Paginate[T any](db *gorm.DB, page Pagination, dest *[]*T) (int64, error) {
	var total int64
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return 0, err
	}
	err := db.Offset(page.Offset).Limit(page.Limit).Order("id DESC").Find(dest).Error
	return total, err
}

// FindByID 根据 ID 查找记录
func FindByID[T any](db *gorm.DB, id uint) (*T, error) {
	var result T
	err := db.Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindByKey 根据 key 字段查找记录
func FindByKey[T any](db *gorm.DB, key string) (*T, error) {
	var result T
	err := db.Where("`key` = ?", key).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}
