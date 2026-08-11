# Task 11: 修复 CI 覆盖率条件

**Files:**
- Modify: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的 CI 配置

## 步骤

### Step 1: 修改 CI 配置

```yaml
- name: Upload coverage
  if: matrix.go-version == '1.26'
  uses: codecov/codecov-action@v3
```

### Step 2: 测试 CI 流程

```bash
# 推送到 GitHub 并检查 CI 运行
```

### Step 3: 提交更改

```bash
git add .github/workflows/ci.yml
git commit -m "ci: fix coverage upload condition to match Go version"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查