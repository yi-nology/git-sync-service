# 前端页面重构设计

## 问题诊断

### 1. 同步历史 vs 同步执行日志 — 功能重叠
- `/sync/history` (SyncHistory): 基于 store，有删除/重跑，无详情 drawer
- `/logs/sync` (SyncLogs): 直接调 API，有 drawer 详情（步骤时间线），有重试
- 两者都展示同步执行记录，用户困惑该看哪个

### 2. Webhook 事件归属不清
- 放在"触发配置"下，但本质是事件日志，跟"日志管理"平行

### 3. 侧边栏冗余
- "新建任务"占一个导航位，但它是个动作不是页面
- "触发配置"只有两个子项，其中一个(事件日志)不该在这

### 4. 设置页面半成品
- 系统设置、高级配置都是 placeholder，只有平台管理有实际功能

## 重构方案

### 侧边栏新结构

```
📊 仪表盘                    → /dashboard

🔄 同步管理
   ├─ 任务列表               → /sync
   └─ 执行记录 (NEW)         → /sync/records

📁 仓库管理                   → /repos

📋 日志中心
   ├─ 操作日志               → /logs/operations
   ├─ 系统日志               → /logs/system
   └─ Webhook 事件           → /logs/webhook-events

⚙️ 平台管理                   → /settings/platforms
```

### 核心改动

#### 改动 1: 合并同步历史 + 同步执行日志 → 「执行记录」

**新页面:** `views/sync/SyncRecords.vue` (路由 `/sync/records`)

**合并策略:** 取两者精华
- 保留 SyncLogs 的 drawer 详情（步骤时间线、错误信息、commit range）
- 保留 SyncHistory 的删除记录、重跑操作
- 保留 SyncHistory 的源→目标 flow 可视化列
- 统一使用服务端分页（不再客户端聚合遍历所有任务）

**Stats 卡片 (4 个):**
- 今日执行数 (蓝色)
- 成功率 (绿色)
- 失败数 (红色)
- 平均耗时 (橙色)

**过滤栏:**
- 任务选择器 (下拉)
- 状态筛选 (success/failed/running)
- 触发方式筛选 (manual/cron/webhook/retry)
- 时间范围选择器
- 搜索 + 重置按钮

**表格列:**
- 时间 (start_time → end_time)
- 任务名称 (可点击打开 drawer)
- 源仓库/分支 → 目标仓库/分支 (flow 可视化)
- 触发方式 (彩色 tag)
- 状态 (StatusBadge)
- 耗时
- 操作 (查看详情 / 重试 / 删除)

**Drawer 详情 (640px):**
- 基本信息: 任务名称、触发方式、状态、开始/结束时间、耗时
- 同步信息: 源仓库、源分支、目标分支、commit range
- 执行日志: 结构化步骤时间线 (彩色圆点 + 连线)
- 错误信息: 失败时显示 alert
- 底部操作: 重试按钮 (失败时)、关闭按钮

**删除的页面:**
- `views/sync/SyncHistory.vue` — 删除
- `views/logs/SyncLogs.vue` — 删除
- 路由 `/sync/history` 和 `/logs/sync` — 删除

#### 改动 2: Webhook 事件移入日志中心

**变更:** 将 `WebhookEvents.vue` 的路由从 `/webhook/events` 改为 `/logs/webhook-events`

**侧边栏:** 从"触发配置"移到"日志中心"下

**"触发配置"子菜单删除** — 只剩 Webhook 规则，提升为独立菜单项

#### 改动 3: 新建任务改为页面内按钮

**变更:** 从侧边栏删除"新建任务"入口

**替代:** 在任务列表页 (`/sync`) 顶部添加醒目的"新建任务"按钮（已有，确认保留）

#### 改动 4: 设置精简

**变更:** 删除半成品页面
- 删除 `views/settings/Settings.vue` (系统设置)
- 删除 `views/settings/AdvancedSettings.vue` (高级配置)
- 保留 `views/settings/PlatformSettings.vue` (平台管理) — 唯一有实际功能的

**路由:** 删除 `/settings` 和 `/settings/advanced`，`/settings/platforms` 保持不变

**侧边栏:** "平台管理"提升为顶级菜单项

#### 改动 5: Webhook 规则提升

**变更:** 从"触发配置"子菜单提升为独立菜单项

**侧边栏:**
```
🔗 Webhook 规则              → /webhook/rules
```

### 路由变更汇总

| 操作 | 旧路由 | 新路由 | 说明 |
|---|---|---|---|
| 删除 | `/sync/history` | — | 合并到执行记录 |
| 删除 | `/logs/sync` | — | 合并到执行记录 |
| 新增 | — | `/sync/records` | 统一执行记录页 |
| 移动 | `/webhook/events` | `/logs/webhook-events` | 归入日志中心 |
| 删除 | `/settings` | — | 半成品 |
| 删除 | `/settings/advanced` | — | 半成品 |

### 文件变更汇总

| 操作 | 文件 |
|---|---|
| 新建 | `views/sync/SyncRecords.vue` |
| 删除 | `views/sync/SyncHistory.vue` |
| 删除 | `views/logs/SyncLogs.vue` |
| 删除 | `views/settings/Settings.vue` |
| 删除 | `views/settings/AdvancedSettings.vue` |
| 修改 | `components/layout/AppSider.vue` — 导航结构 |
| 修改 | `router/index.ts` — 路由配置 |
| 修改 | `components/layout/AppHeader.vue` — 面包屑映射 |

### 不变的部分

- Dashboard、仓库管理、操作日志、系统日志 — 保持不变
- WebhookEvents 组件本身不变，只改路由路径
- PlatformSettings 不变
- 所有 API 调用和 store 不变（新页面复用现有 API）
