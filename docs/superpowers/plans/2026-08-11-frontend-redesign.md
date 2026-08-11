# Frontend Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete frontend redesign of Git Sync Service to achieve modern SaaS-style UI with improved usability, including dashboard, sync management, repo management, trigger config, logging, and platform management.

**Architecture:** Vue 3 + Composition API + Ant Design Vue 4.x frontend with Go + Hertz backend API restructuring for unified response format, pagination, filtering, and batch operations.

**Tech Stack:** Vue 3, TypeScript, Ant Design Vue 4.x, Pinia, Axios, Go, Hertz, GORM

## Global Constraints

- Design System: Use colors (#1677FF primary, #52C41A success, #FAAD14 warning, #FF4D4F error)
- Spacing: 4px base unit (xs: 4px, sm: 8px, md: 16px, lg: 24px, xl: 32px)
- Border Radius: sm: 4px, md: 8px, lg: 12px, xl: 16px
- Layout: Remove sidebar, use top navigation (56px height)
- Components: Ant Design Vue 4.x with consistent styling
- API Response Format: Unified `{code, message, data, timestamp}` format
- Pagination: `{page, page_size, total, total_pages}` in response

---

## Phase 1: Foundation (Week 1-2)

### Task 1: Design System Setup

**Files:**
- Modify: `frontend/src/styles/variables.scss`
- Create: `frontend/src/styles/mixins.scss`
- Create: `frontend/src/styles/global.scss`

**Interfaces:**
- Produces: Design tokens (colors, spacing, typography, shadows)

- [ ] **Step 1: Update variables.scss with design tokens**

```scss
// Colors
$primary: #1677FF;
$success: #52C41A;
$warning: #FAAD14;
$error: #FF4D4F;

$text-primary: #1F1F1F;
$text-secondary: #8C8C8C;
$text-tertiary: #BFBFBF;

$bg-primary: #FFFFFF;
$bg-secondary: #F5F5F5;
$border: #E8E8E8;

// Spacing
$spacing-xs: 4px;
$spacing-sm: 8px;
$spacing-md: 16px;
$spacing-lg: 24px;
$spacing-xl: 32px;

// Border Radius
$radius-sm: 4px;
$radius-md: 8px;
$radius-lg: 12px;
$radius-xl: 16px;

// Shadows
$shadow-card: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
$shadow-card-hover: 0 6px 16px 0 rgba(0, 0, 0, 0.08);
```

- [ ] **Step 2: Create mixins.scss with reusable patterns**

```scss
@mixin card-style {
  background: $bg-primary;
  border-radius: $radius-md;
  box-shadow: $shadow-card;
  transition: box-shadow 0.2s ease;
  
  &:hover {
    box-shadow: $shadow-card-hover;
  }
}

@mixin flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

@mixin text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
```

- [ ] **Step 3: Create global.scss with base styles**

```scss
@import './variables.scss';
@import './mixins.scss';

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  color: $text-primary;
  background: $bg-secondary;
}

.page-container {
  background: $bg-secondary;
  min-height: 100vh;
  padding: $spacing-lg;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: $spacing-lg;
  
  .page-title {
    font-size: 22px;
    font-weight: 600;
    color: $text-primary;
    margin: 0;
    line-height: 1.3;
  }
  
  .page-subtitle {
    font-size: 14px;
    color: $text-secondary;
    margin: 4px 0 0 0;
  }
}
```

- [ ] **Step 4: Verify styles compile**

Run: `cd frontend && npm run dev`
Expected: No compilation errors

- [ ] **Step 5: Commit design system**

```bash
git add frontend/src/styles/
git commit -m "feat: add design system with variables, mixins, and global styles"
```

---

### Task 2: Top Navigation Layout

**Files:**
- Modify: `frontend/src/components/layout/AppLayout.vue`
- Create: `frontend/src/components/layout/TopNav.vue`

**Interfaces:**
- Consumes: Vue Router routes
- Produces: Top navigation component with logo, menu, and user actions

- [ ] **Step 1: Create TopNav.vue component**

```vue
<template>
  <header class="top-nav">
    <div class="nav-left">
      <div class="logo" @click="router.push('/dashboard')">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
        <span class="logo-text">Git Sync</span>
      </div>
      
      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="horizontal"
        @click="handleMenuClick"
      >
        <a-menu-item key="/dashboard">
          <DashboardOutlined />
          <span>仪表盘</span>
        </a-menu-item>
        
        <a-sub-menu key="sync">
          <template #title><SyncOutlined /> 同步管理</template>
          <a-menu-item key="/sync">同步任务</a-menu-item>
          <a-menu-item key="/sync/history">同步历史</a-menu-item>
          <a-menu-item key="/sync/new">新建任务</a-menu-item>
        </a-sub-menu>
        
        <a-menu-item key="/repos">
          <FolderOutlined />
          <span>仓库管理</span>
        </a-menu-item>
        
        <a-sub-menu key="webhook">
          <template #title><ApiOutlined /> 触发配置</template>
          <a-menu-item key="/webhook/rules">Webhook 规则</a-menu-item>
          <a-menu-item key="/webhook/events">事件日志</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="logs">
          <template #title><FileTextOutlined /> 日志管理</template>
          <a-menu-item key="/logs/operations">操作日志</a-menu-item>
          <a-menu-item key="/logs/sync">同步执行日志</a-menu-item>
          <a-menu-item key="/logs/system">系统运行日志</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="system">
          <template #title><SettingOutlined /> 系统</template>
          <a-menu-item key="/settings">系统设置</a-menu-item>
          <a-menu-item key="/settings/platforms">平台管理</a-menu-item>
          <a-menu-item key="/settings/advanced">高级配置</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </div>
    
    <div class="nav-right">
      <a-tooltip title="刷新页面">
        <a-button type="text" @click="handleRefresh">
          <ReloadOutlined />
        </a-button>
      </a-tooltip>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DashboardOutlined,
  SyncOutlined,
  FolderOutlined,
  ApiOutlined,
  FileTextOutlined,
  SettingOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'

const route = useRoute()
const router = useRouter()

const selectedKeys = ref<string[]>([])

watch(() => route.path, (path) => {
  let highlightPath = path
  if (path.startsWith('/repos/config/') || path.startsWith('/local-repos/')) {
    highlightPath = '/repos'
  }
  selectedKeys.value = [highlightPath]
}, { immediate: true })

const handleMenuClick = ({ key }: { key: string }) => {
  router.push(key)
}

const handleRefresh = () => {
  router.go(0)
}
</script>

<style scoped lang="scss">
.top-nav {
  height: 56px;
  background: #FFFFFF;
  border-bottom: 1px solid #F0F0F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  color: #1677FF;
  
  &:hover {
    opacity: 0.8;
  }
  
  .logo-text {
    font-size: 18px;
    font-weight: 700;
    letter-spacing: -0.5px;
  }
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
```

- [ ] **Step 2: Update AppLayout.vue to use TopNav**

```vue
<template>
  <a-layout class="app-layout">
    <TopNav />
    <a-layout-content class="app-content">
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </a-layout-content>
  </a-layout>
</template>

<script setup lang="ts">
import TopNav from './TopNav.vue'
</script>

<style scoped lang="scss">
.app-layout {
  min-height: 100vh;
}

.app-content {
  padding: 24px;
  min-height: calc(100vh - 56px);
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.15s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
```

- [ ] **Step 3: Test navigation works**

Run: `cd frontend && npm run dev`
Expected: Top navigation displays correctly, menu items navigate to correct routes

- [ ] **Step 4: Commit layout changes**

```bash
git add frontend/src/components/layout/
git commit -m "feat: implement top navigation layout replacing sidebar"
```

---

### Task 3: Router Updates for New Pages

**Files:**
- Modify: `frontend/src/router/index.ts`

**Interfaces:**
- Produces: Routes for all new pages including logs

- [ ] **Step 1: Add log routes to router**

```typescript
import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: AppLayout,
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      { path: '/dashboard', component: () => import('@/views/dashboard/Dashboard.vue') },
      
      // Sync Management
      { path: '/sync', component: () => import('@/views/sync/SyncTaskList.vue') },
      { path: '/sync/history', component: () => import('@/views/sync/SyncHistory.vue') },
      { path: '/sync/new', component: () => import('@/views/sync/NewSyncTask.vue') },
      
      // Repo Management
      { path: '/repos', component: () => import('@/views/repos/RepoList.vue') },
      { path: '/repos/config/:id', component: () => import('@/views/repos/RepoUnifiedConfig.vue') },
      { path: '/local-repos/:id', component: () => import('@/views/repos/LocalRepoDetail.vue') },
      
      // Trigger Config
      { path: '/webhook/rules', component: () => import('@/views/webhook/WebhookRules.vue') },
      { path: '/webhook/events', component: () => import('@/views/webhook/WebhookEvents.vue') },
      
      // Logs (NEW)
      { path: '/logs/operations', component: () => import('@/views/logs/OperationLogs.vue') },
      { path: '/logs/sync', component: () => import('@/views/logs/SyncLogs.vue') },
      { path: '/logs/system', component: () => import('@/views/logs/SystemLogs.vue') },
      
      // System
      { path: '/settings', component: () => import('@/views/settings/Settings.vue') },
      { path: '/settings/platforms', component: () => import('@/views/settings/PlatformSettings.vue') },
      { path: '/settings/advanced', component: () => import('@/views/settings/AdvancedSettings.vue') },
    ]
  }
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
```

- [ ] **Step 2: Verify routes compile**

Run: `cd frontend && npm run dev`
Expected: No TypeScript errors

- [ ] **Step 3: Commit router updates**

```bash
git add frontend/src/router/index.ts
git commit -m "feat: add routes for log management pages"
```

---

## Phase 2: Backend API Restructuring (Week 2-3)

### Task 4: Unified API Response Format

**Files:**
- Modify: `internal/pkg/response/response.go`
- Create: `internal/pkg/response/types.go`

**Interfaces:**
- Produces: Unified response format `{code, message, data, timestamp}`

- [ ] **Step 1: Create response types**

```go
package response

import (
	"time"
)

// Response 统一响应格式
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// PaginatedData 分页数据
type PaginatedData struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// BatchResponse 批量操作响应
type BatchResponse struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}
```

- [ ] **Step 2: Update response helpers**

```go
package response

import (
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func Success(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

func Created(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:      0,
		Message:   "created",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

func BadRequest(c *app.RequestContext, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:      400,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

func NotFound(c *app.RequestContext, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:      404,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

func InternalError(c *app.RequestContext, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:      500,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}

func Paginated(c *app.RequestContext, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PaginatedData{
			List: list,
			Pagination: Pagination{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
		Timestamp: time.Now().Unix(),
	})
}
```

- [ ] **Step 3: Test response format**

Run: `go build ./...`
Expected: No compilation errors

- [ ] **Step 4: Commit response format**

```bash
git add internal/pkg/response/
git commit -m "feat: implement unified API response format with pagination"
```

---

### Task 5: Repo API with Pagination and Filtering

**Files:**
- Modify: `biz/handler/git_sync/repo_service.go`
- Modify: `internal/service/repo.go`
- Modify: `internal/dao/repo_dao.go`

**Interfaces:**
- Consumes: `ListReposRequest{page, page_size, search, platform, status, sort_by, sort_order}`
- Produces: Paginated repo list with filtering

- [ ] **Step 1: Add ListReposRequest struct**

```go
// biz/handler/git_sync/repo_service.go

type ListReposRequest struct {
	Page      int    `query:"page" default:"1"`
	PageSize  int    `query:"page_size" default:"10"`
	Search    string `query:"search"`
	Platform  string `query:"platform"`
	Status    string `query:"status"`
	SortBy    string `query:"sort_by" default:"created_at"`
	SortOrder string `query:"sort_order" default:"desc"`
}
```

- [ ] **Step 2: Update ListRepos handler**

```go
func ListRepos(ctx context.Context, c *app.RequestContext) {
	var req ListReposRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	repos, total, err := GetSyncService().ListReposWithFilter(ctx, req.Page, req.PageSize, req.Search, req.Platform, req.Status, req.SortBy, req.SortOrder)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Paginated(c, repos, total, req.Page, req.PageSize)
}
```

- [ ] **Step 3: Add DAO method with filtering**

```go
// internal/dao/repo_dao.go

func (d *RepoDAO) ListWithFilter(page, pageSize int, search, platform, status, sortBy, sortOrder string) ([]*model.Repo, int64, error) {
	var repos []*model.Repo
	var total int64

	query := d.db.Model(&model.Repo{})

	if search != "" {
		query = query.Where("name LIKE ? OR clone_url LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order(sortBy + " " + sortOrder).Find(&repos).Error; err != nil {
		return nil, 0, err
	}

	return repos, total, nil
}
```

- [ ] **Step 4: Test API**

Run: `go build ./...`
Expected: No compilation errors

- [ ] **Step 5: Commit repo API changes**

```bash
git add biz/handler/git_sync/repo_service.go internal/service/repo.go internal/dao/repo_dao.go
git commit -m "feat: add pagination and filtering to repo list API"
```

---

### Task 6: Batch Operations API

**Files:**
- Modify: `biz/handler/git_sync/repo_service.go`
- Create: `biz/handler/git_sync/batch_service.go`

**Interfaces:**
- Consumes: `BatchRequest{action, keys}`
- Produces: `BatchResponse{total, success, failed, errors}`

- [ ] **Step 1: Create batch service handler**

```go
// biz/handler/git_sync/batch_service.go

package git_sync

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

type BatchRequest struct {
	Action string   `json:"action" vd:"required"`
	Keys   []string `json:"keys" vd:"required"`
}

func BatchRepos(ctx context.Context, c *app.RequestContext) {
	var req BatchRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result := &response.BatchResponse{
		Total: len(req.Keys),
	}

	for _, key := range req.Keys {
		var err error
		switch req.Action {
		case "delete":
			err = GetSyncService().DeleteRepo(ctx, key)
		default:
			err = fmt.Errorf("unsupported action: %s", req.Action)
		}

		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", key, err.Error()))
		} else {
			result.Success++
		}
	}

	response.Success(c, result)
}
```

- [ ] **Step 2: Register batch route**

```go
// biz/handler/git_sync/init.go

// Add to route registration
// POST /api/v1/repos/batch -> BatchRepos
```

- [ ] **Step 3: Test compilation**

Run: `go build ./...`
Expected: No compilation errors

- [ ] **Step 4: Commit batch operations**

```bash
git add biz/handler/git_sync/batch_service.go
git commit -m "feat: add batch operations API for repos"
```

---

## Phase 3: Frontend Core Modules (Week 3-5)

### Task 7: Dashboard Redesign

**Files:**
- Modify: `frontend/src/views/dashboard/Dashboard.vue`

**Interfaces:**
- Consumes: `repoApi.list()`, `syncTaskApi.list()`, `systemApi.status()`
- Produces: Dashboard with stats, quick actions, recent tasks/repos, system status

- [ ] **Step 1: Create system API**

```typescript
// frontend/src/api/index.ts

export const systemApi = {
  status: () =>
    api.get<any, SystemStatusResp>('/system/status'),
  health: () =>
    api.get<any, HealthResp>('/system/health'),
}

export interface SystemStatusResp {
  status: string
  version: string
  uptime: number
  repoCount: number
  taskCount: number
  runningTask: number
  lastSyncAt: string
}

export interface HealthResp {
  status: string
  database: {
    status: string
    size: number
  }
}
```

- [ ] **Step 2: Update Dashboard.vue with new design**

[Include full implementation as designed in spec section 4]

- [ ] **Step 3: Test dashboard**

Run: `cd frontend && npm run dev`
Expected: Dashboard displays with stats, quick actions, recent tasks/repos, system status

- [ ] **Step 4: Commit dashboard**

```bash
git add frontend/src/views/dashboard/Dashboard.vue frontend/src/api/index.ts
git commit -m "feat: redesign dashboard with modern SaaS style and system status"
```

---

### Task 8: Operation Logs Page

**Files:**
- Create: `frontend/src/views/logs/OperationLogs.vue`
- Modify: `frontend/src/api/index.ts`

**Interfaces:**
- Consumes: `logApi.listOperations()`
- Produces: Operation logs list with filtering

- [ ] **Step 1: Add log API**

```typescript
// frontend/src/api/index.ts

export const logApi = {
  listOperations: (params?: OperationLogRequest) =>
    api.get<any, OperationLogResp>('/logs/operations', { params }),
  listSync: (params?: SyncLogRequest) =>
    api.get<any, SyncLogResp>('/logs/sync', { params }),
  listSystem: (params?: SystemLogRequest) =>
    api.get<any, SystemLogResp>('/logs/system', { params }),
}

export interface OperationLogRequest {
  page?: number
  page_size?: number
  search?: string
  action?: string
  user?: string
  start_date?: string
  end_date?: string
}

export interface OperationLog {
  id: number
  time: string
  user: string
  action: string
  resource: string
  details: string
  ip: string
}

export interface OperationLogResp {
  list: OperationLog[]
  pagination: Pagination
}

export interface SyncLogRequest {
  page?: number
  page_size?: number
  task_key?: string
  status?: string
  start_date?: string
  end_date?: string
}

export interface SyncLog {
  id: number
  taskName: string
  trigger: string
  status: string
  startTime: string
  endTime: string
  duration: number
  sourceRepo: string
  targetRepo: string
  commits: number
  error?: string
}

export interface SyncLogResp {
  list: SyncLog[]
  pagination: Pagination
}

export interface SystemLogRequest {
  page?: number
  page_size?: number
  level?: string
  module?: string
  search?: string
  start_date?: string
  end_date?: string
}

export interface SystemLog {
  id: number
  time: string
  level: string
  module: string
  message: string
  details?: string
  stack?: string
}

export interface SystemLogResp {
  list: SystemLog[]
  pagination: Pagination
}
```

- [ ] **Step 2: Create OperationLogs.vue**

[Include full implementation as designed in spec section 8.2]

- [ ] **Step 3: Test logs page**

Run: `cd frontend && npm run dev`
Expected: Operation logs page displays with filtering and pagination

- [ ] **Step 4: Commit logs page**

```bash
git add frontend/src/views/logs/OperationLogs.vue frontend/src/api/index.ts
git commit -m "feat: add operation logs page with filtering and pagination"
```

---

### Task 9: Sync Execution Logs Page

**Files:**
- Create: `frontend/src/views/logs/SyncLogs.vue`

**Interfaces:**
- Consumes: `logApi.listSync()`
- Produces: Sync execution logs with status and details

- [ ] **Step 1: Create SyncLogs.vue**

[Include full implementation as designed in spec section 8.3]

- [ ] **Step 2: Test sync logs page**

Run: `cd frontend && npm run dev`
Expected: Sync logs page displays with filtering, pagination, and detail drawer

- [ ] **Step 3: Commit sync logs**

```bash
git add frontend/src/views/logs/SyncLogs.vue
git commit -m "feat: add sync execution logs page with detail drawer"
```

---

### Task 10: System Logs Page

**Files:**
- Create: `frontend/src/views/logs/SystemLogs.vue`

**Interfaces:**
- Consumes: `logApi.listSystem()`
- Produces: System logs with level filtering and auto-refresh

- [ ] **Step 1: Create SystemLogs.vue**

[Include full implementation as designed in spec section 8.4]

- [ ] **Step 2: Test system logs page**

Run: `cd frontend && npm run dev`
Expected: System logs page displays with level filtering, auto-refresh, and detail drawer

- [ ] **Step 3: Commit system logs**

```bash
git add frontend/src/views/logs/SystemLogs.vue
git commit -m "feat: add system logs page with auto-refresh and level filtering"
```

---

## Phase 4: Remaining Modules (Week 5-6)

### Task 11: Sync Task List Redesign

**Files:**
- Modify: `frontend/src/views/sync/SyncTaskList.vue`

**Interfaces:**
- Consumes: `syncTaskApi.list()` with pagination and filtering
- Produces: Redesigned sync task list with search, filter, batch operations

- [ ] **Step 1: Update SyncTaskList.vue with new design**

[Include full implementation as designed in spec section 5.1]

- [ ] **Step 2: Update syncTask store for pagination**

```typescript
// frontend/src/stores/syncTask.ts

export const useSyncTaskStore = defineStore('syncTask', () => {
  const tasks = ref<SyncTask[]>([])
  const loading = ref(false)
  const total = ref(0)

  async function fetchTasks(params?: any) {
    loading.value = true
    try {
      const resp = await syncTaskApi.list(params)
      tasks.value = resp.data.tasks || []
      total.value = resp.data.total || 0
    } finally {
      loading.value = false
    }
  }

  // ... other methods

  return {
    tasks,
    loading,
    total,
    fetchTasks,
    // ... other methods
  }
})
```

- [ ] **Step 3: Test sync task list**

Run: `cd frontend && npm run dev`
Expected: Sync task list displays with filtering, pagination, and batch selection

- [ ] **Step 4: Commit sync task list**

```bash
git add frontend/src/views/sync/SyncTaskList.vue frontend/src/stores/syncTask.ts
git commit -m "feat: redesign sync task list with filtering and batch operations"
```

---

### Task 12: Repo List Redesign

**Files:**
- Modify: `frontend/src/views/repos/RepoList.vue`

**Interfaces:**
- Consumes: `repoApi.list()` with pagination and filtering
- Produces: Redesigned repo list with card layout

- [ ] **Step 1: Update RepoList.vue with card layout**

[Include full implementation as designed in spec section 6.1]

- [ ] **Step 2: Update repo store for pagination**

```typescript
// frontend/src/stores/repo.ts

export const useRepoStore = defineStore('repo', () => {
  const repos = ref<Repo[]>([])
  const loading = ref(false)
  const total = ref(0)

  async function fetchRepos(params?: any) {
    loading.value = true
    try {
      const resp = await repoApi.list(params)
      repos.value = resp.data.repos || []
      total.value = resp.data.total || 0
    } finally {
      loading.value = false
    }
  }

  // ... other methods

  return {
    repos,
    loading,
    total,
    fetchRepos,
    // ... other methods
  }
})
```

- [ ] **Step 3: Test repo list**

Run: `cd frontend && npm run dev`
Expected: Repo list displays with card layout, filtering, and pagination

- [ ] **Step 4: Commit repo list**

```bash
git add frontend/src/views/repos/RepoList.vue frontend/src/stores/repo.ts
git commit -m "feat: redesign repo list with card layout and filtering"
```

---

## Phase 5: Polish and Testing (Week 6-7)

### Task 13: Platform Settings Redesign

**Files:**
- Modify: `frontend/src/views/settings/PlatformSettings.vue`

**Interfaces:**
- Consumes: `platformApi.list()`, `platformApi.create()`, `platformApi.update()`, `platformApi.delete()`
- Produces: Platform management with card layout and two-step creation

- [ ] **Step 1: Update PlatformSettings.vue with new design**

[Include full implementation as designed in spec section 9]

- [ ] **Step 2: Test platform settings**

Run: `cd frontend && npm run dev`
Expected: Platform settings displays with card layout, two-step creation, and sync functionality

- [ ] **Step 3: Commit platform settings**

```bash
git add frontend/src/views/settings/PlatformSettings.vue
git commit -m "feat: redesign platform settings with card layout and sync"
```

---

### Task 14: Final Integration and Testing

**Files:**
- All modified files

**Interfaces:**
- End-to-end testing of all features

- [ ] **Step 1: Run frontend build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 2: Run backend build**

Run: `go build ./...`
Expected: Build succeeds without errors

- [ ] **Step 3: Test all pages manually**

- Dashboard loads correctly with stats
- Navigation works for all routes
- Sync task list with filtering and pagination
- Repo list with card layout
- Operation logs page
- Sync logs page
- System logs page
- Platform settings with card layout

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "feat: complete frontend redesign with modern SaaS style"
```

---

## Summary

This plan covers:
1. **Foundation**: Design system, top navigation, router updates
2. **Backend API**: Unified response format, pagination, filtering, batch operations
3. **Frontend Core**: Dashboard, operation logs, sync logs, system logs
4. **Remaining Modules**: Sync task list, repo list redesign
5. **Polish**: Platform settings, integration testing

Total estimated time: 6-7 weeks for full implementation
