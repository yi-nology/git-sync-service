# Task 2: 修改加密逻辑

**Files:**
- Modify: `internal/dao/repo_dao.go:1-85`

**Interfaces:**
- 依赖: `credential.NewCryptoManager()` 返回 `*CryptoManager`
- 提供: `RepoDAO` 结构体，包含 `cm *credential.CryptoManager` 字段

## 步骤

### Step 1: 添加 CryptoManager 字段到 RepoDAO 结构体

```go
type RepoDAO struct {
    db  *gorm.DB
    cm  *credential.CryptoManager
}
```

### Step 2: 更新 NewRepoDAO 函数

```go
func NewRepoDAO(db *gorm.DB) *RepoDAO {
    cm, err := credential.NewCryptoManager()
    if err != nil {
        panic(err)
    }
    return &RepoDAO{db: db, cm: cm}
}
```

### Step 3: 替换 FindByKey 中的解密逻辑

```go
func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
    var repo model.Repo
    err := d.db.Where("`key` = ?", key).First(&repo).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    decrypted, err := d.cm.Decrypt(repo.AccessToken)
    if err != nil {
        return &repo, nil
    }
    repo.AccessToken = decrypted
    return &repo, nil
}
```

### Step 4: 替换 Create 中的加密逻辑

```go
func (d *RepoDAO) Create(repo *model.Repo) error {
    encrypted, err := d.cm.Encrypt(repo.AccessToken)
    if err != nil {
        return err
    }
    repo.AccessToken = encrypted
    return d.db.Create(repo).Error
}
```

### Step 5: 替换 Update 中的加密逻辑

```go
func (d *RepoDAO) Update(repo *model.Repo) error {
    if repo.AccessToken != "" {
        encrypted, err := d.cm.Encrypt(repo.AccessToken)
        if err != nil {
            return err
        }
        repo.AccessToken = encrypted
    }
    return d.db.Save(repo).Error
}
```

### Step 6: 运行测试验证加密逻辑

```bash
go test ./internal/dao/... -v
```

### Step 7: 提交更改

```bash
git add internal/dao/repo_dao.go
git commit -m "feat: replace EncryptGCM/DecryptGCM with CryptoManager"
```

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过