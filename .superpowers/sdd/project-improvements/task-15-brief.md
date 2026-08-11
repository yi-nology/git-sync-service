# Task 15: 修复 README 端口错误

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: 无
- Produces: 修复后的文档

## 步骤

### Step 1: 修改 README.md 中的端口引用

```markdown
API 文档地址: http://localhost:8890
```

### Step 2: 验证端口配置

```bash
grep -n "8080" README.md
```

### Step 3: 提交更改

```bash
git add README.md
git commit -m "docs: fix port reference in README (8080 -> 8890)"
```

## 全局约束

- 所有 API 端点必须有认证保护
- Lock 和 Semaphore 必须是线程安全的
- 所有核心功能必须有测试覆盖
- 文档必须准确完整
- 代码必须通过所有 lint 检查