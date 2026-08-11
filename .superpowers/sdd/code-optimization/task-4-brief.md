# Task 4: 提取 DAO 分页通用函数

**Files:**
- Create: `internal/dao/dao_helper.go`
- Modify: 所有 DAO 文件

**问题:** 分页模式重复 6 次

## 步骤

### Step 1: 创建通用分页函数

创建 `internal/dao/dao_helper.go`：

```go
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
```

### Step 2: 更新所有 DAO 文件使用通用函数

更新以下文件：
- `internal/dao/repo_dao.go`
- `internal/dao/sync_task_dao.go`
- `internal/dao/sync_run_dao.go`
- `internal/dao/webhook_event_dao.go`
- `internal/dao/webhook_rule_dao.go`

### Step 3: 运行测试

```bash
go test ./internal/dao/... -v
```

### Step 4: 提交

```bash
git add internal/dao/dao_helper.go internal/dao/*.go
git commit -m "refactor: extract generic DAO helper functions"
```

## 全局约束

- 所有测试必须通过
- 代码必须通过 lint 检查
- 不引入新的 breaking changes