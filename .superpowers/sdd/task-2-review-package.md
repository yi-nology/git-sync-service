# Task 2: 修改加密逻辑 - 审查包

## 提交记录

```
commit 9b22948
Author: zhangyi
Date:   2026-08-08

    feat: replace EncryptGCM/DecryptGCM with CryptoManager
    
    1 file changed, 9 insertions(+), 4 deletions(-)
```

## Diff 统计

```
 internal/dao/repo_dao.go | 13 +++++++++----
 1 file changed, 9 insertions(+), 4 deletions(-)
```

## 完整 Diff

```diff
diff --git a/internal/dao/repo_dao.go b/internal/dao/repo_dao.go
index 1234567..abcdefg 100644
--- a/internal/dao/repo_dao.go
+++ b/internal/dao/repo_dao.go
@@ -29,6 +29,7 @@ func DefaultPagination(offset, limit int) Pagination {
 
 type RepoDAO struct {
 	db *gorm.DB
+	cm *credential.CryptoManager
 }
 
 func NewRepoDAO(db *gorm.DB) *RepoDAO {
-	return &RepoDAO{db: db}
+	cm, err := credential.NewCryptoManager()
+	if err != nil {
+		panic(err)
+	}
+	return &RepoDAO{db: db, cm: cm}
 }
 
 func (d *RepoDAO) FindAll(page Pagination) ([]*model.Repo, int64, error) {
@@ -54,7 +55,7 @@ func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
 	if err != nil {
 		return nil, err
 	}
-	decrypted, err := credential.DecryptGCM(repo.AccessToken)
+	decrypted, err := d.cm.Decrypt(repo.AccessToken)
 	if err != nil {
 		return &repo, nil
 	}
@@ -63,7 +64,7 @@ func (d *RepoDAO) FindByKey(key string) (*model.Repo, error) {
 }
 
 func (d *RepoDAO) Create(repo *model.Repo) error {
-	encrypted, err := credential.EncryptGCM(repo.AccessToken)
+	encrypted, err := d.cm.Encrypt(repo.AccessToken)
 	if err != nil {
 		return err
 	}
@@ -73,7 +74,7 @@ func (d *RepoDAO) Create(repo *model.Repo) error {
 
 func (d *RepoDAO) Update(repo *model.Repo) error {
 	if repo.AccessToken != "" {
-		encrypted, err := credential.EncryptGCM(repo.AccessToken)
+		encrypted, err := d.cm.Encrypt(repo.AccessToken)
 		if err != nil {
 			return err
 		}
```

## 任务简报要求

1. 添加 CryptoManager 字段到 RepoDAO 结构体
2. 更新 NewRepoDAO 函数以初始化 CryptoManager
3. 替换 FindByKey 中的解密逻辑
4. 替换 Create 中的加密逻辑
5. 替换 Update 中的加密逻辑
6. 运行测试验证加密逻辑
7. 提交更改

## 全局约束

- 所有加密操作必须使用 `CryptoManager` 替换 `credential.EncryptGCM/DecryptGCM`
- 平台检测错误必须处理 `ErrPlatformNotSupported`
- 必须设置 `ENCRYPTION_KEY` 环境变量
- 所有现有测试必须通过